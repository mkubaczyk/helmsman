---
version: v3.0.0
---

# CMD reference

This lists available CMD options in Helmsman:

> you can find the CMD options for the version you are using by typing: `helmsman -h` or `helmsman --help`

  `--always-upgrade`
        upgrade release even if no changes are found.

  `--apply`
        apply the plan directly.

  `--context-override string`
        override releases context defined in release state with this one.

  `--debug`
        show the debug execution logs and actual helm/kubectl commands. This can log secrets and should only be used for debugging purposes.

  `--verbose`
        show verbose execution logs.

  `--destroy`
        delete all deployed releases.

  `-detailed-exit-code`
        returns a detailed exit code (0 - no changes, 1 - error, 2 - changes present)

  `--diff-context num`
        number of lines of context to show around changes in helm diff output.

  `-p`
        max number of concurrent helm releases to run

  `--dry-run`
        apply the dry-run (do not update) option for helm commands.

  `-e value`
        additional file(s) to load environment variables from, may be supplied more than once, it extends default .env file lookup, every next file takes precedence over previous ones in case of having the same environment variables defined.
        If a `.env` file exists, it will be loaded by default, if additional env files are specified using the `-e` flag, the environment file will be loaded in order where the last file will take precedence.

  `-f value`
        desired state file name(s), may be supplied more than once to merge state files.

  `--force-upgrades`
        use --force when upgrading helm releases. May cause resources to be recreated.

  `--keep-untracked-releases`
        keep releases that are managed by Helmsman from the used DSFs in the command, and are no longer tracked in your desired state.

  `--kubeconfig`
        path to the kubeconfig file to use for CLI requests. Defaults to false if the helm diff plugin is installed.

   `--kubectl-diff`
        Use kubectl diff instead of helm diff

  `--migrate-context`
        Updates the context name for all apps defined in the DSF and applies Helmsman labels. Using this flag is required if you want to change context name after it has been set.

  `--no-banner`
        don't show the banner.

  `--no-color`
        don't use colors.

  `--no-env-subst`
        turn off environment substitution globally.

  `--no-recursive-env-expand`
        disable recursive environment variables expansion.

  `--subst-env-values`
        turn on environment substitution in values files.

  `--no-fancy`
        don't display the banner and don't use colors.

  `--no-ns`
        don't create namespaces.

  `--no-ssm-subst`
        turn off SSM parameter substitution globally.

  `--replace-on-rename`
        uninstall the existing release when a chart with a different name is used.

  `--spec string`
        specification file name, contains locations of desired state files to be merged

  `--subst-ssm-values`
        turn on SSM parameter substitution in values files.

  `--ns-override string`
        override defined namespaces with this one.

  `--show-diff`
        show helm diff results. Can expose sensitive information.

  `--skip-validation`
        skip desired state validation.

  `--target`
        limit execution to specific app.

  `--exclude-target`
        exclude specific app from execution.

  `--group`
        limit execution to specific group of apps.

  `--exclude-group`
        exclude specific group of apps from execution.

  `--update-deps`
        run 'helm dep up' for local chart

  `--check-for-chart-updates`
        compares the chart versions in the state file to the latest versions in the chart repositories and shows available updates

  `--v`    
        show the version.

  `--verify`
        verify the downloaded charts.