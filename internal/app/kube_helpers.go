package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"sigs.k8s.io/yaml"
)

// addNamespaces creates a set of namespaces in your k8s cluster.
// If a namespace with the same name exists, it will skip it.
// If --ns-override flag is used, it only creates the provided namespace in that flag
func addNamespaces(s *State) {
	var wg sync.WaitGroup
	for nsName, ns := range s.Namespaces {
		if ns.disabled {
			continue
		}
		wg.Add(1)
		go func(name string, cfg *Namespace, wg *sync.WaitGroup) {
			defer wg.Done()
			createNamespace(name, cfg.Labels, cfg.Annotations)
			labelNamespace(name, cfg.Labels, s.Settings.NamespaceLabelsAuthoritative)
			annotateNamespace(name, cfg.Annotations)
			if !flags.dryRun {
				setLimits(name, cfg.Limits)
				setQuotas(name, cfg.Quotas)
			}
		}(nsName, ns, &wg)
	}
	wg.Wait()
}

// kubectl prepares a kubectl command to be executed
func kubectl(args []string, desc string) Command {
	return Command{
		Cmd:         kubectlBin,
		Args:        args,
		Description: desc,
	}
}

// createNamespace creates a namespace in the k8s cluster
func createNamespace(ns string, labels map[string]string, annotations map[string]string) {
	checkCmd := kubectl([]string{"get", "namespace", ns}, "Looking for namespace [ "+ns+" ]")
	if _, err := checkCmd.Exec(); err == nil {
		log.Verbose("Namespace [ " + ns + " ] exists")
		return
	}

	targetFile, err := tempBuildNamespaceDefinition(ns, labels, annotations)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer os.Remove(targetFile.Name())

	cmd := kubectl([]string{"apply", "-f", targetFile.Name(), flags.getKubeDryRunFlag("apply")},
		"Creating namespace [ "+ns+" ]")

	if _, err := cmd.Exec(); err != nil {
		log.Fatalf("Failed creating namespace [ "+ns+" ] with error: %v", err)
	}
	log.Info("Namespace [ " + ns + " ] created")
}

// tempBuildNamespaceDefinition creates temporary file with Namespace object definition
func tempBuildNamespaceDefinition(ns string, labels map[string]string, annotations map[string]string) (f *os.File, err error) {
	definition := `
---
apiVersion: v1
kind: Namespace
metadata:
`
	definition += fmt.Sprintf("  name: \"%s\"\n", ns)

	metadata := make(map[string]map[string]string)
	metadata["labels"] = labels
	metadata["annotations"] = annotations
	d, err := yaml.Marshal(&metadata)
	if err != nil {
		return nil, err
	}
	definition += Indent(string(d), strings.Repeat(" ", 2))

	targetFile, err := ioutil.TempFile(execTempDir, "namespace-*.yaml")
	if err != nil {
		return nil, err
	}

	if _, err = targetFile.Write([]byte(definition)); err != nil {
		return nil, err
	}

	return targetFile, nil
}

// labelNamespace labels a namespace with provided labels
func labelNamespace(ns string, labels map[string]string, authoritative bool) {
	var nsLabels map[string]string

	args := []string{"label", "--overwrite", "namespace/" + ns, flags.getKubeDryRunFlag("label")}

	if authoritative {
		cmdGetLabels := kubectl([]string{"get", "namespace", ns, "-o", "jsonpath='{.metadata.labels}'"}, "Getting namespace [ "+ns+" ] current labels")
		res, err := cmdGetLabels.Exec()
		if err != nil {
			log.Error(fmt.Sprintf("Could not get namespace [ %s ] labels. Error message: %v", ns, err))
		}
		if err := json.Unmarshal([]byte(strings.Trim(res.output, `'`)), &nsLabels); err != nil {
			log.Fatal(fmt.Sprintf("failed to unmarshal kubectl get namespace labels output: %s, ended with error: %s", res.output, err))
		}
		// ignore default k8s namespace label from being removed
		delete(nsLabels, "kubernetes.io/metadata.name")
		// ignore every label defined in DSF for the namespace from being removed
		for definedLabelKey := range labels {
			delete(nsLabels, definedLabelKey)
		}
		for label := range nsLabels {
			args = append(args, label+"-")
		}
	}
	if len(labels) == 0 && len(nsLabels) == 0 {
		return
	}

	for k, v := range labels {
		args = append(args, k+"="+v)
	}
	cmd := kubectl(args, "Labeling namespace [ "+ns+" ]")

	if _, err := cmd.Exec(); err != nil && flags.verbose {
		log.Warning(fmt.Sprintf("Could not label namespace [ %s with %v ]. Error message: %v", ns, strings.Join(args[4:], ","), err))
	}
}

// annotateNamespace annotates a namespace with provided annotations
func annotateNamespace(ns string, annotations map[string]string) {
	if len(annotations) == 0 {
		return
	}

	args := []string{"annotate", "--overwrite", "namespace/" + ns, flags.getKubeDryRunFlag("annotate")}
	for k, v := range annotations {
		args = append(args, k+"="+v)
	}
	cmd := kubectl(args, "Annotating namespace [ "+ns+" ]")

	if _, err := cmd.Exec(); err != nil && flags.verbose {
		log.Info(fmt.Sprintf("Could not annotate namespace [ %s with %v ]. Error message: %v", ns, strings.Join(args[4:], ","), err))
	}
}

// setLimits creates a LimitRange resource in the provided Namespace
func setLimits(ns string, lims Limits) {
	if len(lims) == 0 {
		return
	}

	definition := `
---
apiVersion: v1
kind: LimitRange
metadata:
  name: limit-range
spec:
  limits:
`
	d, err := yaml.Marshal(&lims)
	if err != nil {
		log.Fatal(err.Error())
	}

	definition += Indent(string(d), strings.Repeat(" ", 4))

	if err := apply(definition, ns, "LimitRange"); err != nil {
		log.Fatal(err.Error())
	}
}

func setQuotas(ns string, quotas *Quotas) {
	if quotas == nil {
		return
	}

	definition := `
---
apiVersion: v1
kind: ResourceQuota
metadata:
  name: resource-quota
spec:
  hard:
`

	for _, customQuota := range quotas.CustomQuotas {
		definition += Indent(customQuota.Name+": '"+customQuota.Value+"'\n", strings.Repeat(" ", 4))
	}

	// Special formatting for custom quotas so manually write these and then set to nil for marshalling
	quotas.CustomQuotas = nil

	d, err := yaml.Marshal(&quotas)
	if err != nil {
		log.Fatal(err.Error())
	}

	definition += Indent(string(d), strings.Repeat(" ", 4))

	if err := apply(definition, ns, "ResourceQuota"); err != nil {
		log.Fatal(err.Error())
	}
}

func apply(definition, ns, kind string) error {
	targetFile, err := ioutil.TempFile(execTempDir, kind+"-*.yaml")
	if err != nil {
		return err
	}
	defer os.Remove(targetFile.Name())

	if _, err = targetFile.Write([]byte(definition)); err != nil {
		return err
	}

	cmd := kubectl([]string{"apply", "-f", targetFile.Name(), "-n", ns, flags.getKubeDryRunFlag("apply")},
		"Creating "+kind+" in namespace [ "+ns+" ]")

	if _, err := cmd.Exec(); err != nil {
		return fmt.Errorf("error creating %s in namespace [ %s ]: %w", kind, ns, err)
	}

	return nil
}

// createContext creates a context -connecting to a k8s cluster- in kubectl config.
// It returns true if successful, false otherwise
func createContext(s *State) error {
	if s.Settings.BearerToken && s.Settings.BearerTokenPath == "" {
		log.Info("Creating kube context with bearer token from K8S service account.")
		s.Settings.BearerTokenPath = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	} else if s.Settings.BearerToken && s.Settings.BearerTokenPath != "" {
		log.Info("Creating kube context with bearer token from " + s.Settings.BearerTokenPath)
	} else if s.Settings.Password == "" || s.Settings.Username == "" || s.Settings.ClusterURI == "" {
		return errors.New("missing information to create context [ " + s.Settings.KubeContext + " ] " +
			"you are either missing PASSWORD, USERNAME or CLUSTERURI in the Settings section of your desired state file.")
	} else if !s.Settings.BearerToken && (s.Certificates == nil || s.Certificates["caCrt"] == "" || s.Certificates["caKey"] == "") {
		return errors.New("missing information to create context [ " + s.Settings.KubeContext + " ] " +
			"you are either missing caCrt or caKey or both in the Certifications section of your desired state file.")
	} else if s.Settings.BearerToken && (s.Certificates == nil || s.Certificates["caCrt"] == "") {
		return errors.New("missing information to create context [ " + s.Settings.KubeContext + " ] " +
			"caCrt is missing in the Certifications section of your desired state file.")
	}

	// set certs locations (relative filepath, GCS bucket, AWS bucket)
	caCrt := s.Certificates["caCrt"]
	caKey := s.Certificates["caKey"]
	caClient := s.Certificates["caClient"]

	// download certs and keys
	// GCS bucket+file format should be: gs://bucket-name/dir.../filename.ext
	// S3 bucket+file format should be: s3://bucket-name/dir.../filename.ext

	// CA cert
	if caCrt != "" {
		caCrt = downloadFile(caCrt, "", "ca.crt")
	}

	// CA key
	if caKey != "" {
		caKey = downloadFile(caKey, "", "ca.key")
	}

	// client certificate
	if caClient != "" {
		caClient = downloadFile(caClient, "", "client.crt")
	}

	// bearer token
	tokenPath := "bearer.token"
	if s.Settings.BearerToken && s.Settings.BearerTokenPath != "" {
		tokenPath = downloadFile(s.Settings.BearerTokenPath, "", tokenPath)
	}

	// connecting to the cluster
	setCredentialsCmdArgs := []string{}
	if s.Settings.BearerToken {
		token := readFile(tokenPath)
		if s.Settings.Username == "" {
			s.Settings.Username = "helmsman"
		}
		setCredentialsCmdArgs = append(setCredentialsCmdArgs, "config", "set-credentials", s.Settings.Username, "--token="+token)
	} else {
		setCredentialsCmdArgs = append(setCredentialsCmdArgs, "config", "set-credentials", s.Settings.Username, "--username="+s.Settings.Username,
			"--password="+s.Settings.Password, "--client-key="+caKey)
		if caClient != "" {
			setCredentialsCmdArgs = append(setCredentialsCmdArgs, "--client-certificate="+caClient)
		}
	}
	cmd := kubectl(setCredentialsCmdArgs, "Creating kubectl context - setting credentials")

	if _, err := cmd.Exec(); err != nil {
		return fmt.Errorf("failed to create context [ "+s.Settings.KubeContext+" ]: %w", err)
	}

	cmd = kubectl([]string{"config", "set-cluster", s.Settings.KubeContext, "--server=" + s.Settings.ClusterURI, "--certificate-authority=" + caCrt}, "Creating kubectl context - setting cluster")

	if _, err := cmd.Exec(); err != nil {
		return fmt.Errorf("failed to create context [ "+s.Settings.KubeContext+" ]: %w", err)
	}

	cmd = kubectl([]string{"config", "set-context", s.Settings.KubeContext, "--cluster=" + s.Settings.KubeContext, "--user=" + s.Settings.Username}, "Creating kubectl context - setting context")

	if _, err := cmd.Exec(); err != nil {
		return fmt.Errorf("failed to create context [ "+s.Settings.KubeContext+" ]: %w", err)
	}

	if setKubeContext(s.Settings.KubeContext) {
		return nil
	}

	return errors.New("something went wrong while setting the kube context to the newly created one")
}

// setKubeContext sets your kubectl context to the one specified in the desired state file.
// It returns false if it fails to set the context. This means the context does not exist.
func setKubeContext(kctx string) bool {
	if kctx == "" {
		return getKubeContext()
	}

	cmd := kubectl([]string{"config", "use-context", kctx}, "Setting kube context to [ "+kctx+" ]")

	if _, err := cmd.Exec(); err != nil {
		log.Info("Kubectl context [ " + kctx + " ] does not exist. Attempting to create it...")
		return false
	}

	return true
}

// getKubeContext gets your kubectl context.
// It returns false if no context is set.
func getKubeContext() bool {
	cmd := kubectl([]string{"config", "current-context"}, "Getting kubectl context")

	if res, err := cmd.Exec(); err != nil || res.output == "" {
		log.Info("Kubectl context is not set")
		return false
	}

	return true
}

// getReleaseContext extracts the Helmsman release context from the helm storage driver objects (secret or configmap) labels
func getReleaseContext(releaseName, namespace, storageBackend string) string {
	// kubectl get secrets -n test1 -l MANAGED-BY=HELMSMAN -o=jsonpath='{.items[0].metadata.labels.HELMSMAN_CONTEXT}'
	// kubectl get secret sh.helm.release.v1.argo.v1  -n test1  -o=jsonpath='{.metadata.labels.HELMSMAN_CONTEXT}'
	// kubectl get secret -l owner=helm,name=argo -n test1 -o=jsonpath='{.items[-1].metadata.labels.HELMSMAN_CONTEXT}'
	cmd := kubectl([]string{"get", storageBackend, "-n", namespace, "-l", "owner=helm", "-l", "name=" + releaseName, "-o", "jsonpath='{.items[-1].metadata.labels.HELMSMAN_CONTEXT}'"}, "Getting Helmsman context for [ "+releaseName+" ] release")

	res, err := cmd.Exec()
	if err != nil {
		log.Fatal(err.Error())
	}
	rctx := strings.Trim(res.output, `"' `)
	if rctx == "" {
		rctx = defaultContextName
	}
	return rctx
}

// getKubectlVersion returns kubectl client version
func getKubectlVersion() string {
	var kubectlVersion struct {
		ClientVersion map[string]string `json:"clientVersion"`
	}

	cmd := kubectl([]string{"version", "--output=json", "--client"}, "Checking kubectl version")
	ver, err := cmd.Exec()
	if err != nil {
		log.Fatalf("While checking kubectl version: %v", err)
	}

	if err := json.Unmarshal([]byte(ver.output), &kubectlVersion); err != nil {
		log.Fatal(fmt.Sprintf("failed to unmarshal kubectl version CLI output: %s", err))
	}

	version := kubectlVersion.ClientVersion["gitVersion"]
	if version == "" {
		log.Fatal("Could not get 'gitVersion' from kubectl version cmd output")
	}
	return version
}

func checkKubectlVersion(constraint string) bool {
	return checkVersion(getKubectlVersion(), constraint)
}

// getKubeDryRunFlag returns kubectl dry-run flag if helmsman --dry-run flag is enabled
// TODO: this should be cleanup once 1.18 is old enough
func (c *cli) getKubeDryRunFlag(action string) string {
	var flag string
	if c.dryRun {
		flag = "--dry-run"
		recent := checkKubectlVersion(">=v1.18.0")
		switch action {
		case "apply":
			if recent {
				flag += "=server"
			} else {
				flag = "--server-dry-run"
			}
		default:
			if recent {
				flag += "=client"
			}
		}
	}
	return flag
}
