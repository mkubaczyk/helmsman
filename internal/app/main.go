// Package app contains the main logic for the application.
package app

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/sync/errgroup"
)

const (
	helmBin            = "helm"
	kubectlBin         = "kubectl"
	defaultContextName = "default"
	resourcePool       = 10
	tempFilesDir       = ".helmsman-tmp"
)

var appVersion = "dev"

const (
	exitCodeSucceed            = 0
	exitCodeSucceedWithChanges = 2
)

var (
	flags      cli
	settings   *Config
	curContext string
	log        = &Logger{}
)

func init() {
	// Parse cli flags and read config files
	flags.setup()
}

// resolveKubeconfigPath returns the active kubeconfig path: the -kubeconfig flag
// value if set, then the KUBECONFIG env var, then the default ~/.kube/config.
func resolveKubeconfigPath(f *cli) string {
	if f.kubeconfig != "" {
		return f.kubeconfig
	}
	if v := os.Getenv("KUBECONFIG"); v != "" {
		return v
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".kube", "config")
}

// copyKubeconfig copies src to dst, creating dst if needed.
func copyKubeconfig(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0o700); err != nil {
		return err
	}
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

// runParallelFiles re-execs the current helmsman binary once per DSF file, running
// them concurrently. Each subprocess gets full process isolation: its own uniquely-named
// temp dir (via os.MkdirTemp), its own kubeconfig copy (so kubectl context changes don't
// race on the shared ~/.kube/config), and its own in-memory globals. The -p flag controls
// max concurrency across DSF files.
func runParallelFiles(f *cli) int {
	exe, err := os.Executable()
	if err != nil {
		log.Fatal("could not determine helmsman executable path: " + err.Error())
	}

	// Create a short-lived temp dir for kubeconfig copies. Each subprocess gets
	// its own copy so concurrent `kubectl config use-context` calls don't race.
	orchTmpDir, err := os.MkdirTemp("", "helmsman-parallel-*")
	if err != nil {
		log.Fatal("could not create orchestrator temp dir: " + err.Error())
	}
	defer os.RemoveAll(orchTmpDir)

	srcKubeconfig := resolveKubeconfigPath(f)

	// Reconstruct args from the parsed flag set, skipping flags handled per-subprocess.
	// flag.Visit only iterates flags explicitly set by the user (not defaults).
	// --name=value is valid for all Go flag types (bool, string, int).
	var baseArgs []string
	flag.Visit(func(fl *flag.Flag) {
		switch fl.Name {
		case "parallel-files", "f", "kubeconfig":
			return // each subprocess gets its own -f and -kubeconfig
		}
		// stringArray flags: String() returns space-joined values — iterate the
		// underlying slice directly to avoid splitting values that contain spaces.
		// Note: fileOptionArray (-f) is always skipped above; no case needed here.
		if a, ok := fl.Value.(*stringArray); ok {
			for _, val := range *a {
				baseArgs = append(baseArgs, fmt.Sprintf("--%s=%s", fl.Name, val))
			}
			return
		}
		// All other types (bool, string, int): --name=value works uniformly.
		baseArgs = append(baseArgs, fmt.Sprintf("--%s=%s", fl.Name, fl.Value.String()))
	})

	g := new(errgroup.Group)
	g.SetLimit(f.parallel)

	for i, file := range f.files {
		g.Go(func() error {
			// Each subprocess gets its own kubeconfig copy so that
			// `kubectl config use-context` calls don't collide on disk.
			kubeconfigCopy := filepath.Join(orchTmpDir, fmt.Sprintf("kubeconfig-%d", i))
			if err := copyKubeconfig(srcKubeconfig, kubeconfigCopy); err != nil {
				return fmt.Errorf("could not copy kubeconfig for DSF %s: %w", file.name, err)
			}

			args := append([]string{"-f", file.name, "--kubeconfig=" + kubeconfigCopy}, baseArgs...)
			cmd := exec.Command(exe, args...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin
			return cmd.Run()
		})
	}

	if err := g.Wait(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		log.Error("parallel execution failed: " + err.Error())
		return 1
	}
	return exitCodeSucceed
}

// Main is the app main function
func Main() int {
	var s State

	flags.parse()

	// When --parallel-files is set, re-exec once per DSF file as isolated subprocesses.
	if flags.parallelFiles {
		if len(flags.files) < 2 {
			log.Warning("--parallel-files has no effect with fewer than 2 -f flags; running normally")
		} else {
			return runParallelFiles(&flags)
		}
	}

	// Each process gets its own unique temp dir via os.MkdirTemp so that
	// concurrent helmsman runs don't collide on disk.
	var tmpErr error
	flags.tempDir, tmpErr = os.MkdirTemp("", tempFilesDir+"-*")
	if tmpErr != nil {
		log.Fatal("could not create temp dir: " + tmpErr.Error())
	}
	defer os.RemoveAll(flags.tempDir)
	if !flags.noCleanup {
		defer s.cleanup()
	}

	if err := flags.readState(&s); err != nil {
		log.Fatal(err.Error())
	}

	if len(flags.target) > 0 && len(s.targetMap) == 0 {
		log.Info("No apps defined with -target flag were found, exiting")
		os.Exit(0)
	}

	if len(flags.group) > 0 && len(s.targetMap) == 0 {
		log.Info("No apps defined with -group flag were found, exiting")
		os.Exit(0)
	}

	log.SlackWebhook = s.Settings.SlackWebhook
	log.MSTeamsWebhook = s.Settings.MSTeamsWebhook

	settings = &s.Settings
	curContext = s.Context

	// set the kubecontext to be used Or create it if it does not exist
	log.Info("Setting up kubectl")
	if !setKubeContext(s.Settings.KubeContext) {
		if err := createContext(&s); err != nil {
			log.Fatal(err.Error())
		}
	}

	// add repos -- fails if they are not valid
	log.Info("Setting up helm")
	if err := addHelmRepos(s.HelmRepos); err != nil && !flags.destroy {
		log.Fatal(err.Error())
	}

	if flags.apply || flags.dryRun || flags.destroy {
		// add/validate namespaces
		if !flags.noNs {
			log.Info("Setting up namespaces")
			if flags.nsOverride == "" {
				addNamespaces(&s)
			} else {
				createNamespace(flags.nsOverride, nil, nil)
				s.overrideAppsNamespace(flags.nsOverride)
			}
		}
	}

	log.Info("Getting chart information")

	err := s.getReleaseChartsInfo()
	if flags.skipValidation {
		log.Info("Skipping charts' validation.")
	} else if err != nil {
		log.Fatal(err.Error())
	} else {
		log.Info("Charts validated.")
	}

	if flags.destroy {
		log.Warning("Destroy flag is enabled. Your releases will be deleted!")
	}

	if flags.migrateContext {
		log.Warning("migrate-context flag is enabled. Context will be changed to [ " + s.Context + " ] and Helmsman labels will be applied.")
		s.updateContextLabels()
	}

	if flags.checkForChartUpdates {
		for _, r := range s.Apps {
			r.checkChartForUpdates()
		}
	}

	log.Info("Preparing plan")
	cs := s.getCurrentState()
	p := cs.makePlan(&s)
	if !flags.keepUntrackedReleases {
		cs.cleanUntrackedReleases(&s, p)
	}

	p.sort()
	p.print()
	if flags.debug {
		p.printCmds()
	}
	p.sendToSlack()
	p.sendToMSTeams()

	if flags.apply || flags.dryRun || flags.destroy {
		p.exec()
	}

	exitCode := exitCodeSucceed

	if flags.detailedExitCode && len(p.Commands) > 0 {
		exitCode = exitCodeSucceedWithChanges
	}

	return exitCode
}
