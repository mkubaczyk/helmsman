// Package app contains the main logic for the application.
package app

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

const (
	helmBin            = "helm"
	kubectlBin         = "kubectl"
	defaultContextName = "default"
	resourcePool       = 10
)

// tempFilesDir is overridden in Main() to a PID-scoped path so that concurrent
// helmsman processes don't collide. The default is used by tests.
var tempFilesDir = ".helmsman-tmp"

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
// them concurrently. Each subprocess gets full process isolation: its own PID-scoped
// temp dir, its own kubeconfig copy (so kubectl context changes don't race on the
// shared ~/.kube/config), and its own in-memory globals. The -p flag controls max
// concurrency across DSF files.
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

	// Rebuild args from os.Args, stripping -parallel-files, -kubeconfig, and all
	// -f flags. We add back a single -f and a per-subprocess -kubeconfig below.
	var baseArgs []string
	skipNext := false
	for _, arg := range os.Args[1:] {
		if skipNext {
			skipNext = false
			continue
		}
		if arg == "-parallel-files" || arg == "--parallel-files" {
			continue
		}
		if arg == "-f" || arg == "--f" || arg == "-kubeconfig" || arg == "--kubeconfig" {
			skipNext = true
			continue
		}
		if strings.HasPrefix(arg, "-f=") || strings.HasPrefix(arg, "--f=") ||
			strings.HasPrefix(arg, "-kubeconfig=") || strings.HasPrefix(arg, "--kubeconfig=") {
			continue
		}
		baseArgs = append(baseArgs, arg)
	}

	concurrency := f.parallel
	if concurrency < 1 {
		concurrency = 1
	}
	sem := make(chan struct{}, concurrency)

	var wg sync.WaitGroup
	exitCodes := make([]int, len(f.files))

	for i, file := range f.files {
		wg.Add(1)
		go func(idx int, filename string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			// Each subprocess gets its own kubeconfig copy so that
			// `kubectl config use-context` calls don't collide on disk.
			kubeconfigCopy := filepath.Join(orchTmpDir, fmt.Sprintf("kubeconfig-%d", idx))
			if err := copyKubeconfig(srcKubeconfig, kubeconfigCopy); err != nil {
				log.Fatal(fmt.Sprintf("could not copy kubeconfig for DSF %s: %v", filename, err))
			}

			args := append(append([]string{}, baseArgs...), "-f", filename, "-kubeconfig", kubeconfigCopy)
			cmd := exec.Command(exe, args...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin
			if err := cmd.Run(); err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					exitCodes[idx] = exitErr.ExitCode()
				} else {
					exitCodes[idx] = 1
				}
			}
		}(i, file.name)
	}

	wg.Wait()

	for _, code := range exitCodes {
		if code != 0 {
			return code
		}
	}
	return exitCodeSucceed
}

// Main is the app main function
func Main() int {
	var s State

	flags.parse()

	// When --parallel-files is set, re-exec once per DSF file as isolated subprocesses.
	if flags.parallelFiles && len(flags.files) > 1 {
		return runParallelFiles(&flags)
	}

	// Each process gets its own temp dir so concurrent helmsman runs don't
	// collide — the first process to exit would otherwise wipe a shared dir.
	tempFilesDir = fmt.Sprintf(".helmsman-tmp-%d", os.Getpid())
	defer os.RemoveAll(tempFilesDir)
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
