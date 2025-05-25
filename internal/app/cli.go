package app

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	banner = "\n" +
		" _          _ \n" +
		"| |        | | \n" +
		"| |__   ___| |_ __ ___  ___ _ __ ___   __ _ _ __\n" +
		"| '_ \\ / _ \\ | '_ ` _ \\/ __| '_ ` _ \\ / _` | '_ \\ \n" +
		"| | | |  __/ | | | | | \\__ \\ | | | | | (_| | | | | \n" +
		"|_| |_|\\___|_|_| |_| |_|___/_| |_| |_|\\__,_|_| |_|"
	slogan = "A Helm-Charts-as-Code tool.\n\n"
)

// Allow parsing of multiple string command line options into an array of strings
type stringArray []string

type fileOptionArray []fileOption

type fileOption struct {
	name     string
	priority int
}

func (f *fileOptionArray) String() string {
	var a []string
	for _, v := range *f {
		a = append(a, v.name)
	}
	return strings.Join(a, " ")
}

func (f *fileOptionArray) Set(value string) error {
	var fo fileOption

	fo.name = value
	*f = append(*f, fo)
	return nil
}

func (i *stringArray) String() string {
	return strings.Join(*i, " ")
}

func (i *stringArray) Set(value string) error {
	*i = append(*i, value)
	return nil
}

type cli struct {
	debug                 bool
	files                 fileOptionArray
	spec                  string
	envFiles              stringArray
	target                stringArray
	targetExcluded        stringArray
	group                 stringArray
	groupExcluded         stringArray
	kubeconfig            string
	apply                 bool
	destroy               bool
	dryRun                bool
	verbose               bool
	noBanner              bool
	noColors              bool
	noFancy               bool
	noNs                  bool
	nsOverride            string
	contextOverride       string
	skipValidation        bool
	keepUntrackedReleases bool
	showDiff              bool
	diffContext           int
	noEnvSubst            bool
	substEnvValues        bool
	noRecursiveEnvExpand  bool
	noSSMSubst            bool
	substSSMValues        bool
	detailedExitCode      bool
	updateDeps            bool
	forceUpgrades         bool
	renameReplace         bool
	version               bool
	noCleanup             bool
	migrateContext        bool
	parallel              int
	alwaysUpgrade         bool
	noUpdate              bool
	kubectlDiff           bool
	downloadCharts        bool
	checkForChartUpdates  bool
	skipIgnoredApps       bool
	skipPendingApps       bool
	pendingAppRetries     int
	showSecrets           bool
	verify                bool
}

func printUsage() {
	fmt.Print(banner)
	fmt.Println("Helmsman version: " + appVersion)
	fmt.Println("Helmsman is a Helm Charts as Code tool which allows you to automate the deployment/management of your Helm charts.")
	fmt.Println("")
	fmt.Printf("Usage: helmsman [options]\n")
	flag.PrintDefaults()
}

func (c *cli) setup() {
	// parsing command line flags
	flag.Var(&c.files, "f", "desired state file name(s), may be supplied more than once to merge state files")
	flag.Var(&c.envFiles, "e", "additional file(s) to load environment variables from, may be supplied more than once, it extends default .env file lookup, every next file takes precedence over previous ones in case of having the same environment variables defined")
	flag.Var(&c.target, "target", "limit execution to specific app.")
	flag.Var(&c.group, "group", "limit execution to specific group of apps.")
	flag.Var(&c.targetExcluded, "exclude-target", "exclude specific app from execution.")
	flag.Var(&c.groupExcluded, "exclude-group", "exclude specific group of apps from execution.")
	flag.IntVar(&c.diffContext, "diff-context", -1, "number of lines of context to show around changes in helm diff output")
	flag.IntVar(&c.parallel, "p", 1, "max number of concurrent helm releases to run")
	flag.StringVar(&c.spec, "spec", "", "specification file name, contains locations of desired state files to be merged")
	flag.StringVar(&c.kubeconfig, "kubeconfig", "", "path to the kubeconfig file to use for CLI requests")
	flag.StringVar(&c.nsOverride, "ns-override", "", "override defined namespaces with this one")
	flag.StringVar(&c.contextOverride, "context-override", "", "override releases context defined in release state with this one")
	flag.BoolVar(&c.apply, "apply", false, "apply the plan directly")
	flag.BoolVar(&c.dryRun, "dry-run", false, "apply the dry-run option for helm commands.")
	flag.BoolVar(&c.destroy, "destroy", false, "delete all deployed releases.")
	flag.BoolVar(&c.version, "v", false, "show the version")
	flag.BoolVar(&c.debug, "debug", false, "show the debug execution logs and actual helm/kubectl commands. This can log secrets and should only be used for debugging purposes.")
	flag.BoolVar(&c.verbose, "verbose", false, "show verbose execution logs.")
	flag.BoolVar(&c.noBanner, "no-banner", false, "don't show the banner")
	flag.BoolVar(&c.noColors, "no-color", false, "don't use colors")
	flag.BoolVar(&c.noFancy, "no-fancy", false, "don't display the banner and don't use colors")
	flag.BoolVar(&c.noNs, "no-ns", false, "don't create namespaces")
	flag.BoolVar(&c.skipValidation, "skip-validation", false, "skip desired state validation")
	flag.BoolVar(&c.keepUntrackedReleases, "keep-untracked-releases", false, "keep releases that are managed by Helmsman from the used DSFs in the command, and are no longer tracked in your desired state.")
	flag.BoolVar(&c.showDiff, "show-diff", false, "show helm diff results. Can expose sensitive information.")
	flag.BoolVar(&c.detailedExitCode, "detailed-exit-code", false, "returns a detailed exit code (0 - no changes, 1 - error, 2 - changes present)")
	flag.BoolVar(&c.noEnvSubst, "no-env-subst", false, "turn off environment substitution globally")
	flag.BoolVar(&c.substEnvValues, "subst-env-values", false, "turn on environment substitution in values files.")
	flag.BoolVar(&c.noRecursiveEnvExpand, "no-recursive-env-expand", false, "disable recursive environment values expansion")
	flag.BoolVar(&c.noSSMSubst, "no-ssm-subst", false, "turn off SSM parameter substitution globally")
	flag.BoolVar(&c.substSSMValues, "subst-ssm-values", false, "turn on SSM parameter substitution in values files.")
	flag.BoolVar(&c.updateDeps, "update-deps", false, "run 'helm dep up' for local charts")
	flag.BoolVar(&c.forceUpgrades, "force-upgrades", false, "use --force when upgrading helm releases. May cause resources to be recreated.")
	flag.BoolVar(&c.renameReplace, "replace-on-rename", false, "uninstall the existing release when a chart with a different name is used.")
	flag.BoolVar(&c.noCleanup, "no-cleanup", false, "keeps any credentials files that has been downloaded on the host where helmsman runs.")
	flag.BoolVar(&c.migrateContext, "migrate-context", false, "updates the context name for all apps defined in the DSF and applies Helmsman labels. Using this flag is required if you want to change context name after it has been set.")
	flag.BoolVar(&c.alwaysUpgrade, "always-upgrade", false, "upgrade release even if no changes are found")
	flag.BoolVar(&c.noUpdate, "no-update", false, "skip updating helm repos")
	flag.BoolVar(&c.kubectlDiff, "kubectl-diff", false, "use kubectl diff instead of helm diff. Defalts to false if the helm diff plugin is installed.")
	flag.BoolVar(&c.checkForChartUpdates, "check-for-chart-updates", false, "compares the chart versions in the state file to the latest versions in the chart repositories and shows available updates")
	flag.BoolVar(&c.downloadCharts, "download-charts", false, "download charts referenced by URLs in the state file")
	flag.BoolVar(&c.skipIgnoredApps, "skip-ignored", false, "skip ignored apps")
	flag.BoolVar(&c.skipPendingApps, "skip-pending", false, "skip pending helm releases")
	flag.IntVar(&c.pendingAppRetries, "pending-max-retries", 0, "max number of retries for pending helm releases")
	flag.BoolVar(&c.showSecrets, "show-secrets", false, "show helm diff results with secrets.")
	flag.BoolVar(&c.verify, "verify", false, "verify if charts are signed and valid before using them")
	flag.Usage = printUsage
	flag.Parse()
}

// Cli parses cmd flags, validates them and performs some initializations
func (c *cli) parse() {
	if c.version {
		fmt.Println("Helmsman version: " + appVersion)
		os.Exit(0)
	}

	if c.noFancy {
		c.noColors = true
		c.noBanner = true
	}
	verbose := c.verbose || c.debug
	initLogs(verbose, c.noColors)

	if !c.noBanner {
		fmt.Printf("%s version: %s\n%s", banner, appVersion, slogan)
	}

	if c.dryRun && c.apply {
		log.Fatal("--apply and --dry-run can't be used together.")
	}

	if c.destroy && c.apply {
		log.Fatal("--destroy and --apply can't be used together.")
	}

	if len(c.target) > 0 && len(c.group) > 0 {
		log.Fatal("--target and --group can't be used together.")
	}

	if len(flags.files) > 0 && len(flags.spec) > 0 {
		log.Fatal("-f and -spec can't be used together.")
	}

	if c.parallel < 1 {
		c.parallel = 1
	}

	log.Verbose("Helm client version: " + strings.TrimSpace(getHelmVersion()))
	if checkHelmVersion("<3.0.0") {
		log.Fatal("this version of Helmsman does not work with helm releases older than 3.0.0")
	}

	kubectlVersion := getKubectlVersion()
	log.Verbose("kubectl client version: " + kubectlVersion)

	if len(c.files) == 0 && len(c.spec) == 0 {
		log.Info("No desired state files provided.")
		os.Exit(0)
	}

	if c.kubeconfig != "" {
		os.Setenv("KUBECONFIG", c.kubeconfig)
	}

	if !ToolExists(kubectlBin) {
		log.Fatal("kubectl is not installed/configured correctly. Aborting!")
	}

	if !ToolExists(helmBin) {
		log.Fatal("" + helmBin + " is not installed/configured correctly. Aborting!")
	}

	if !c.kubectlDiff && !helmPluginExists("diff") {
		c.kubectlDiff = true
		log.Warning("helm diff not found, using kubectl diff")
	}

	if !c.noEnvSubst {
		log.Verbose("Substitution of env variables enabled")
		if c.substEnvValues {
			log.Verbose("Substitution of env variables in values enabled")
		}
	}

	if !c.noRecursiveEnvExpand {
		log.Verbose("Recursive environment variables expansion is enabled")
	}

	if !c.noSSMSubst {
		log.Verbose("Substitution of SSM variables enabled")
		if c.substSSMValues {
			log.Verbose("Substitution of SSM variables in values enabled")
		}
	}
}

// readState gets the desired state from files
func (c *cli) readState(s *State) error {
	// read the env file if it exists
	if _, err := os.Stat(".env"); err == nil {
		if !stringInSlice(".env", c.envFiles) {
			c.envFiles = append([]string{".env"}, c.envFiles...)
		}
	}

	if err := prepareEnv(c.envFiles); err != nil {
		return err
	}

	// wipe & create a temporary directory
	os.RemoveAll(tempFilesDir)
	_ = os.MkdirAll(tempFilesDir, 0o755)

	if len(c.spec) > 0 {
		sp := new(StateFiles)
		if err := sp.specFromYAML(c.spec); err != nil {
			return fmt.Errorf("error parsing spec file: %w", err)
		}

		for _, val := range sp.StateFiles {
			fo := fileOption{}
			fo.name = val.Path
			if err := isValidFile(fo.name, validManifestFiles); err != nil {
				return fmt.Errorf("invalid -spec file: %w", err)
			}
			c.files = append(c.files, fo)
		}
	}

	// read the TOML/YAML desired state file
	if err := s.build(c.files); err != nil {
		return fmt.Errorf("error building the state from files: %w", err)
	}

	s.disableApps(c.group, c.target, c.groupExcluded, c.targetExcluded)

	if c.skipIgnoredApps {
		s.Settings.SkipIgnoredApps = true
	}
	if c.skipPendingApps {
		s.Settings.SkipPendingApps = true
	}

	if !c.skipValidation {
		// validate the desired state content
		if len(c.files) > 0 {
			log.Info("Validating desired state definition")
			if err := s.validate(); err != nil { // syntax validation
				return err
			}
		}
	} else {
		log.Info("Desired state validation is skipped.")
	}

	if c.debug {
		s.print()
	}
	return nil
}

// getRunFlags returns dry-run and debug flags
func (c *cli) getRunFlags() []string {
	if c.dryRun {
		return []string{"--dry-run", "--debug"}
	}
	if c.debug {
		return []string{"--debug"}
	}
	return []string{}
}
