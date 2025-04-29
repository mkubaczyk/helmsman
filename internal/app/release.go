package app

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Release type representing Helm releases which are described in the desired state
type Release struct {
	// Name is the helm release name
	Name string `json:"name"`
	// Description is a user friendly description of the helm release
	Description string `json:"description,omitempty"`
	// Namespace where to deploy the helm release
	Namespace string `json:"namespace"`
	// Enabled can be used to togle a helm release
	Enabled NullBool `json:"enabled"`
	Group   string   `json:"group,omitempty"`
	Chart   string   `json:"chart"`
	// Version of the helm chart to deploy
	Version string `json:"version"`
	// ValuesFile is the path for a values file for the helm release
	ValuesFile string `json:"valuesFile,omitempty"`
	// ValuesFiles is a list of paths a values files for the helm release
	ValuesFiles []string `json:"valuesFiles,omitempty"`
	// SecretsFile is the path for an encrypted values file for the helm release
	SecretsFile string `json:"secretsFile,omitempty"`
	// SecretsFiles is a list of paths for encrypted values files for the helm release
	SecretsFiles []string `json:"secretsFiles,omitempty"`
	// PostRenderer is the path to an executable to be used for post rendering
	PostRenderer string `json:"postRenderer,omitempty"`
	// Test indicates if the chart tests should be executed
	Test NullBool `json:"test,omitempty"`
	// Protected defines if the release should be protected against changes
	Protected NullBool `json:"protected,omitempty"`
	// Wait defines whether helm should block execution until all k8s resources are in a ready state
	Wait NullBool `json:"wait,omitempty"`
	// Priority allows defining the execution order, releases with the same priority can be executed in parallel
	Priority int `json:"priority,omitempty"`
	// Set can be used to overwrite the chart values
	Set map[string]string `json:"set,omitempty"`
	// SetString can be used to overwrite string values
	SetString map[string]string `json:"setString,omitempty"`
	// SetFile can be used to overwrite the chart values
	SetFile map[string]string `json:"setFile,omitempty"`
	// HelmFlags is a list of additional flags to pass to the helm command
	HelmFlags []string `json:"helmFlags,omitempty"`
	// HelmDiffFlags is a list of cli flags to pass to helm diff
	HelmDiffFlags []string `json:"helmDiffFlags,omitempty"`
	// NoHooks can be used to disable the execution of helm hooks
	NoHooks NullBool `json:"noHooks,omitempty"`
	// Timeout is the number of seconds to wait for the release to complete
	Timeout int `json:"timeout,omitempty"`
	// Hooks can be used to define lifecycle hooks specific to this release
	Hooks map[string]interface{} `json:"hooks,omitempty"`
	// MaxHistory is the maximum number of histoical releases to keep
	MaxHistory int `json:"maxHistory,omitempty"`
	disabled   bool
}

func (r *Release) key() string {
	return fmt.Sprintf("%s-%s", r.Name, r.Namespace)
}

func (r *Release) Disable() {
	r.disabled = true
}

// isReleaseConsideredToRun checks if a release is being targeted for operations as specified by user cmd flags (--group or --target)
func (r *Release) isConsideredToRun() bool {
	if r == nil {
		return false
	}
	return !r.disabled
}

// validate validates if a release inside a desired state meets the specifications or not.
// check the full specification @ https://github.com/mkubaczyk/helmsman/blob/master/docs/desired_state_specification.md
func (r *Release) validate(appLabel string, seen map[string]map[string]bool, s *State) error {
	if seen[r.Name][r.Namespace] {
		return errors.New("release name must be unique within a given namespace")
	}

	if flags.nsOverride == "" && r.Namespace == "" {
		return errors.New("release targeted namespace can't be empty")
	} else if flags.nsOverride == "" && r.Namespace != "" && r.Namespace != "kube-system" && !s.isNamespaceDefined(r.Namespace) {
		return errors.New("release " + r.Name + " is using namespace [ " + r.Namespace + " ] which is not defined in the Namespaces section of your desired state file." +
			" Release [ " + r.Name + " ] can't be installed in that Namespace until its defined.")
	}
	_, err := os.Stat(r.Chart)
	if r.Chart == "" || os.IsNotExist(err) && !strings.ContainsAny(r.Chart, "/") {
		return errors.New("chart can't be empty and must be of the format: repo/chart")
	}
	if r.Version == "" {
		return errors.New("version can't be empty")
	}

	if r.ValuesFile != "" && len(r.ValuesFiles) > 0 {
		return errors.New("valuesFile and valuesFiles should not be used together")
	} else if r.ValuesFile != "" {
		if err := isValidFile(r.ValuesFile, validManifestFiles); err != nil {
			return fmt.Errorf("invalid values file: %w", err)
		}
	} else if len(r.ValuesFiles) > 0 {
		for _, filePath := range r.ValuesFiles {
			if err := isValidFile(filePath, validManifestFiles); err != nil {
				return fmt.Errorf("invalid values file: %w", err)
			}
		}
	}

	if r.SecretsFile != "" && len(r.SecretsFiles) > 0 {
		return errors.New("secretsFile and secretsFiles should not be used together")
	} else if r.SecretsFile != "" {
		if err := isValidFile(r.SecretsFile, validManifestFiles); err != nil {
			return fmt.Errorf("invalid secrets file: %w", err)
		}
	} else if len(r.SecretsFiles) > 0 {
		for _, filePath := range r.SecretsFiles {
			if err := isValidFile(filePath, validManifestFiles); err != nil {
				return fmt.Errorf("invalid secrets file: %w", err)
			}
		}
	}

	if r.PostRenderer != "" && !ToolExists(r.PostRenderer) {
		return fmt.Errorf("%s must be executable and available in your PATH", r.PostRenderer)
	}

	if r.Priority != 0 && r.Priority > 0 {
		return errors.New("priority can only be 0 or negative value, positive values are not allowed")
	}

	if (len(r.Hooks)) != 0 {
		if err := validateHooks(r.Hooks); err != nil {
			return err
		}
	}

	if seen[r.Name] == nil {
		seen[r.Name] = make(map[string]bool)
	}
	seen[r.Name][r.Namespace] = true

	return nil
}

// testRelease creates a Helm command to test a particular release.
func (r *Release) test(afterCommands *[]hookCmd) {
	if flags.dryRun {
		log.Verbose("Dry-run, skipping tests:  " + r.Name)
		return
	}
	cmd := helmCmd(r.getHelmArgsFor("test"), "Running tests for release [ "+r.Name+" ] in namespace [ "+r.Namespace+" ]")
	*afterCommands = append([]hookCmd{{Command: cmd, Type: test}}, *afterCommands...)
}

// installRelease creates a Helm command to install a particular release in a particular namespace using a particular Tiller.
func (r *Release) install(p *plan) {
	before, after := r.checkHooks("install")

	if r.Test.Value {
		r.test(&after)
	}

	cmd := helmCmd(r.getHelmArgsFor("install"), "Install release [ "+r.Name+" ] version [ "+r.Version+" ] in namespace [ "+r.Namespace+" ]")
	p.addCommand(cmd, r.Priority, r, before, after)
}

// uninstall uninstalls a release
func (r *Release) uninstall(p *plan, optionalNamespace ...string) {
	ns := r.Namespace
	if len(optionalNamespace) > 0 {
		ns = optionalNamespace[0]
	}
	priority := r.Priority
	if p.ReverseDelete {
		priority *= -1
	}

	before, after := r.checkHooks("delete", ns)

	cmd := helmCmd(r.getHelmArgsFor("uninstall", ns), "Delete release [ "+r.Name+" ] in namespace [ "+ns+" ]")
	p.addCommand(cmd, priority, r, before, after)
}

// diffRelease diffs an existing release with the specified values.yaml
func (r *Release) diff() (string, error) {
	var (
		args        []string
		maxExitCode int
	)

	if !flags.kubectlDiff {
		args = []string{"diff"}
		if !flags.showSecrets {
			args = append(args, "--suppress-secrets")
		}
		if flags.noColors {
			args = append(args, "--no-color")
		}
		if flags.diffContext != -1 {
			args = append(args, "--context", strconv.Itoa(flags.diffContext))
		}
		args = concat(args, r.getHelmArgsFor("diff"))
	} else {
		args = r.getHelmArgsFor("template")
	}

	desc := "Diffing release [ " + r.Name + " ] in namespace [ " + r.Namespace + " ]"
	cmd := CmdPipe{helmCmd(args, desc)}

	if flags.kubectlDiff {
		cmd = append(cmd, kubectl([]string{"diff", "--namespace", r.Namespace, "-f", "-"}, desc))
		maxExitCode = 1
	}

	res, err := cmd.RetryExecWithThreshold(3, maxExitCode)
	if err != nil {
		if flags.kubectlDiff && res.code <= 1 {
			// kubectl diff exit status:
			//   0 No differences were found.
			//   1 Differences were found.
			//   >1 Kubectl or diff failed with an error.
			return res.output, nil
		}
		return "", fmt.Errorf("command failed: %w", err)
	}

	return res.output, nil
}

// upgradeRelease upgrades an existing release with the specified values.yaml
func (r *Release) upgrade(p *plan) {
	before, after := r.checkHooks("upgrade")

	if r.Test.Value {
		r.test(&after)
	}

	cmd := helmCmd(r.getHelmArgsFor("upgrade"), "Upgrade release [ "+r.Name+" ] to version [ "+r.Version+" ] in namespace [ "+r.Namespace+" ]")

	p.addCommand(cmd, r.Priority, r, before, after)
}

// reInstall uninstalls a release and reinstalls it.
// This is used when moving a release to another namespace or when changing the chart used for it.
// When the release is being moved to another namespace, the optionalOldNamespace is used to provide
// the namespace from which the release is deleted.
func (r *Release) reInstall(p *plan, optionalOldNamespace ...string) {
	oldNamespace := ""
	if len(optionalOldNamespace) > 0 {
		oldNamespace = optionalOldNamespace[0]
	}
	r.uninstall(p, oldNamespace)
	r.install(p)
}

// rollbackRelease evaluates if a rollback action needs to be taken for a given release.
// if the release is already deleted but from a different namespace than the one specified in input,
// it purge deletes it and create it in the specified namespace.
func (r *Release) rollback(cs *currentState, p *plan) {
	rs, ok := cs.releases[r.key()]
	if !ok {
		return
	}

	if r.Namespace == rs.Namespace {

		cmd := helmCmd(concat([]string{"rollback", r.Name, rs.getRevision()}, r.getWait(), r.getTimeout(), r.getNoHooks(), flags.getRunFlags()), "Rolling back release [ "+r.Name+" ] in namespace [ "+r.Namespace+" ]")
		p.addCommand(cmd, r.Priority, r, []hookCmd{}, []hookCmd{})
		r.upgrade(p) // this is to reflect any changes in values file(s)
		p.addDecision("Release [ "+r.Name+" ] was deleted and is desired to be rolled back to "+
			"namespace [ "+r.Namespace+" ]", r.Priority, create)
	} else {
		r.reInstall(p)
		p.addDecision("Release [ "+r.Name+" ] is deleted BUT from namespace [ "+rs.Namespace+
			" ]. Will purge delete it from there and install it in namespace [ "+r.Namespace+" ]", r.Priority, create)
		p.addDecision("WARNING: rolling back release [ "+r.Name+" ] from [ "+rs.Namespace+" ] to [ "+r.Namespace+
			" ] might not correctly connect to existing volumes. Check https://github.com/mkubaczyk/helmsman/blob/master/docs/how_to/apps/moving_across_namespaces.md"+
			" for details if this release uses PV and PVC.", r.Priority, create)
	}
}

// mark applies Helmsman specific labels to Helm's state resources (secrets/configmaps)
func (r *Release) mark(storageBackend string) {
	r.label(storageBackend, "MANAGED-BY=HELMSMAN", "NAMESPACE="+r.Namespace, "HELMSMAN_CONTEXT="+curContext)
}

// label labels Helm's state resources (secrets/configmaps)
func (r *Release) label(storageBackend string, labels ...string) {
	if len(labels) == 0 {
		return
	}
	if r.Enabled.Value {

		args := []string{"label", "--overwrite", storageBackend, "-n", r.Namespace, "-l", "owner=helm,name=" + r.Name}
		args = append(args, labels...)
		cmd := kubectl(args, "Applying Helmsman labels to [ "+r.Name+" ] release")

		if _, err := cmd.Exec(); err != nil {
			log.Fatal(err.Error())
		}
	}
}

// annotate annotates Helm's state resources (secrets/configmaps)
func (r *Release) annotate(storageBackend string, annotations ...string) {
	if len(annotations) == 0 {
		return
	}
	if r.Enabled.Value {

		args := []string{"annotate", "--overwrite", storageBackend, "-n", r.Namespace, "-l", "owner=helm,name=" + r.Name}
		args = append(args, annotations...)
		cmd := kubectl(args, "Applying Helmsman annotations to [ "+r.Name+" ] release")

		if _, err := cmd.Exec(); err != nil {
			log.Fatal(err.Error())
		}
	}
}

// isProtected checks if a release is protected or not.
// A protected is release is either: a) deployed in a protected namespace b) flagged as protected in the desired state file
// Any release in a protected namespace is protected by default regardless of its flag
// returns true if a release is protected, false otherwise
func (r *Release) isProtected(cs *currentState, n *Namespace) bool {
	// if the release does not exist in the cluster, it is not protected
	if ok := cs.releaseExists(r, ""); !ok {
		return false
	}
	if n.Protected || r.Protected.Value {
		return true
	}
	return false
}

// getNoHooks returns the no-hooks flag for install/upgrade commands
func (r *Release) getNoHooks() []string {
	if r.NoHooks.Value {
		return []string{"--no-hooks"}
	}
	return []string{}
}

// getTimeout returns the timeout flag for install/upgrade commands
func (r *Release) getTimeout() []string {
	if r.Timeout != 0 {
		return []string{"--timeout", strconv.Itoa(r.Timeout) + "s"}
	}
	return []string{}
}

// getSetValues returns --set params to be used with helm install/upgrade commands
func (r *Release) getSetValues() []string {
	res := []string{}
	for k, v := range r.Set {
		res = append(res, "--set", k+"="+strings.ReplaceAll(v, ",", "\\,")+"")
	}
	return res
}

// getSetStringValues returns --set-string params to be used with helm install/upgrade commands
func (r *Release) getSetStringValues() []string {
	res := []string{}
	for k, v := range r.SetString {
		res = append(res, "--set-string", k+"="+strings.ReplaceAll(v, ",", "\\,")+"")
	}
	return res
}

// getSetFileValues returns --set-file params to be used with helm install/upgrade commands
func (r *Release) getSetFileValues() []string {
	res := []string{}
	for k, v := range r.SetFile {
		res = append(res, "--set-file", k+"="+strings.ReplaceAll(v, ",", "\\,")+"")
	}
	return res
}

// getWait returns a partial helm command containing the helm wait flag (--wait) if the wait flag for the release was set to true
// Otherwise, retruns an empty string
func (r *Release) getWait() []string {
	res := []string{}
	if r.Wait.Value {
		res = append(res, "--wait")
	}
	return res
}

// getDesiredNamespace returns the namespace of a release
func (r *Release) getDesiredNamespace() string {
	return r.Namespace
}

// getMaxHistory returns the max-history flag for upgrade commands
func (r *Release) getMaxHistory() []string {
	if r.MaxHistory != 0 {
		return []string{"--history-max", strconv.Itoa(r.MaxHistory)}
	}
	return []string{}
}

// getHelmFlags returns helm flags
func (r *Release) getHelmFlags() []string {
	var flgs []string
	if flags.forceUpgrades {
		flgs = append(flgs, "--force")
	}

	return concat(r.getNoHooks(), r.getWait(), r.getTimeout(), r.getMaxHistory(), flags.getRunFlags(), r.HelmFlags, flgs)
}

// getPostRenderer returns the post-renderer Helm flag
func (r *Release) getPostRenderer() []string {
	args := []string{}
	if r.PostRenderer != "" {
		args = append(args, "--post-renderer", r.PostRenderer)
	}
	return args
}

// getHelmArgsFor returns helm arguments for a specific helm operation
func (r *Release) getHelmArgsFor(action string, optionalNamespaceOverride ...string) []string {
	ns := r.Namespace
	if len(optionalNamespaceOverride) > 0 {
		ns = optionalNamespaceOverride[0]
	}
	args := []string{r.Name, r.Chart, "--namespace", ns}
	if r.Version != "latest" {
		args = append(args, "--version", r.Version)
	}
	switch action {
	case "template":
		return concat([]string{"template"}, args, []string{"--skip-tests", "--no-hooks"}, r.getValuesFiles(), r.getSetValues(), r.getSetStringValues(), r.getSetFileValues(), r.getPostRenderer())
	case "install", "upgrade":
		return concat([]string{"upgrade", "--install"}, args, r.getValuesFiles(), r.getSetValues(), r.getSetStringValues(), r.getSetFileValues(), r.getHelmFlags(), r.getPostRenderer())
	case "diff":
		return concat([]string{"upgrade"}, args, r.getValuesFiles(), r.getSetValues(), r.getSetStringValues(), r.getSetFileValues(), r.HelmDiffFlags, r.getPostRenderer())
	case "uninstall":
		return concat([]string{action, "--namespace", ns, r.Name}, flags.getRunFlags())
	default:
		return []string{action, "--namespace", ns, r.Name}
	}
}

func (r *Release) checkChartDepUpdate() {
	if !r.isConsideredToRun() {
		return
	}
	if flags.updateDeps && isLocalChart(r.Chart) {
		if err := updateChartDep(r.Chart); err != nil {
			log.Fatal("helm dependency update failed: " + err.Error())
		}
	}
}

func (r *Release) checkChartForUpdates() {
	if !r.isConsideredToRun() {
		return
	}

	if flags.checkForChartUpdates && !isLocalChart(r.Chart) {
		chartInfo, err := getChartInfo(r.Chart, ">= "+r.Version)
		if err != nil {
			log.Fatal("Couldn't check version for " + r.Chart + ": " + err.Error())
		}

		if checkVersion(r.Version, "< "+chartInfo.Version) {
			log.Infof("Newer version for release %s, chart %s found: current %s, latest %s", r.Name, r.Chart, r.Version, chartInfo.Version)
		}
	}
}

// overrideNamespace overrides a release defined namespace with a new given one
func (r *Release) overrideNamespace(newNs string) {
	log.Info("Overriding namespace for app:  " + r.Name)
	r.Namespace = newNs
}

// inheritHooks passes global hooks config from the state to the release hooks if they are unset
// release hooks override the global ones
func (r *Release) inheritHooks(s *State) {
	if len(s.Settings.GlobalHooks) != 0 {
		if len(r.Hooks) == 0 {
			r.Hooks = s.Settings.GlobalHooks
		} else {
			for key := range s.Settings.GlobalHooks {
				if _, ok := r.Hooks[key]; !ok {
					r.Hooks[key] = s.Settings.GlobalHooks[key]
				}
			}
		}
	}
}

// inheritMaxHistory passes global max history from the state to the release if it is unset
func (r *Release) inheritMaxHistory(s *State) {
	if s.Settings.GlobalMaxHistory != 0 {
		if r.MaxHistory == 0 {
			r.MaxHistory = s.Settings.GlobalMaxHistory
		}
	}
}

// checkHooks checks if a hook of certain type exists and creates its command
// if success condition for the hook is defined, a "kubectl wait" command is created
// returns two slices of before and after commands
func (r *Release) checkHooks(action string, optionalNamespace ...string) ([]hookCmd, []hookCmd) {
	ns := r.Namespace
	if len(optionalNamespace) > 0 {
		ns = optionalNamespace[0]
	}
	var beforeCmds, afterCmds []hookCmd

	switch action {
	case "install":
		{
			if _, ok := r.Hooks[preInstall]; ok {
				beforeCmds = append(beforeCmds, r.getHookCommands(preInstall, ns)...)
			}
			if _, ok := r.Hooks[postInstall]; ok {
				afterCmds = append(afterCmds, r.getHookCommands(postInstall, ns)...)
			}
		}
	case "upgrade":
		{
			if _, ok := r.Hooks[preUpgrade]; ok {
				beforeCmds = append(beforeCmds, r.getHookCommands(preUpgrade, ns)...)
			}
			if _, ok := r.Hooks[postUpgrade]; ok {
				afterCmds = append(afterCmds, r.getHookCommands(postUpgrade, ns)...)
			}
		}
	case "delete":
		{
			if _, ok := r.Hooks[preDelete]; ok {
				beforeCmds = append(beforeCmds, r.getHookCommands(preDelete, ns)...)
			}
			if _, ok := r.Hooks[postDelete]; ok {
				afterCmds = append(afterCmds, r.getHookCommands(postDelete, ns)...)
			}
		}
	}
	return beforeCmds, afterCmds
}

func (r *Release) getHookCommands(hookType, ns string) []hookCmd {
	var cmds []hookCmd
	if _, ok := r.Hooks[hookType]; ok {
		hook := r.Hooks[hookType].(string)
		if err := isValidFile(hook, validManifestFiles); err == nil {
			cmd := kubectl([]string{"apply", "-n", ns, "-f", hook, flags.getKubeDryRunFlag("apply")}, "Apply "+hook+" manifest "+hookType)
			cmds = append(cmds, hookCmd{Command: cmd, Type: hookType})
			if wait, waitCmds := r.shouldWaitForHook(hook, hookType, ns); wait {
				cmds = append(cmds, waitCmds...)
			}
		} else { // shell hook
			args := strings.Fields(hook)
			cmds = append(cmds, hookCmd{
				Command: Command{
					Cmd:         args[0],
					Args:        args[1:],
					Description: fmt.Sprintf("%s shell hook", hookType),
				},
				Type: hookType,
			})
		}
	}
	return cmds
}

// shouldWaitForHook checks if there is a success condition to wait for after applying a hook
// returns a boolean and the wait command if applicable
func (r *Release) shouldWaitForHook(hookFile string, hookType string, namespace string) (bool, []hookCmd) {
	var cmds []hookCmd
	if flags.dryRun {
		return false, cmds
	} else if _, ok := r.Hooks["successCondition"]; ok {
		timeoutFlag := ""
		if _, ok := r.Hooks["successTimeout"]; ok {
			timeoutFlag = "--timeout=" + r.Hooks["successTimeout"].(string)
		}
		cmd := kubectl([]string{"wait", "-n", namespace, "-f", hookFile, "--for=condition=" + r.Hooks["successCondition"].(string), timeoutFlag}, "Wait for "+hookType+" : "+hookFile)
		cmds = append(cmds, hookCmd{Command: cmd})
		if _, ok := r.Hooks["deleteOnSuccess"]; ok && r.Hooks["deleteOnSuccess"].(bool) {
			cmd = kubectl([]string{"delete", "-n", namespace, "-f", hookFile}, "Delete "+hookType+" : "+hookFile)
			cmds = append(cmds, hookCmd{Command: cmd})
		}
		return true, cmds
	}
	return false, cmds
}

// print prints the details of the release
func (r Release) print() {
	fmt.Println("")
	fmt.Println("\tname: ", r.Name)
	fmt.Println("\tdescription: ", r.Description)
	fmt.Println("\tnamespace: ", r.Namespace)
	fmt.Println("\tenabled: ", r.Enabled.Value)
	fmt.Println("\tchart: ", r.Chart)
	fmt.Println("\tversion: ", r.Version)
	fmt.Println("\tvaluesFile: ", r.ValuesFile)
	fmt.Println("\tvaluesFiles: ", strings.Join(r.ValuesFiles, ","))
	fmt.Println("\tpostRenderer: ", r.PostRenderer)
	fmt.Println("\ttest: ", r.Test.Value)
	fmt.Println("\tprotected: ", r.Protected.Value)
	fmt.Println("\twait: ", r.Wait.Value)
	fmt.Println("\tpriority: ", r.Priority)
	fmt.Println("\tSuccessCondition: ", r.Hooks["successCondition"])
	fmt.Println("\tSuccessTimeout: ", r.Hooks["successTimeout"])
	fmt.Println("\tDeleteOnSuccess: ", r.Hooks["deleteOnSuccess"])
	fmt.Println("\tpreInstall: ", r.Hooks[preInstall])
	fmt.Println("\tpostInstall: ", r.Hooks[postInstall])
	fmt.Println("\tpreUpgrade: ", r.Hooks[preUpgrade])
	fmt.Println("\tpostUpgrade: ", r.Hooks[postUpgrade])
	fmt.Println("\tpreDelete: ", r.Hooks[preDelete])
	fmt.Println("\tpostDelete: ", r.Hooks[postDelete])
	fmt.Println("\tno-hooks: ", r.NoHooks.Value)
	fmt.Println("\ttimeout: ", r.Timeout)
	fmt.Println("\tvalues to override from env:")
	printMap(r.Set, 2)
	fmt.Println("------------------- ")
}
