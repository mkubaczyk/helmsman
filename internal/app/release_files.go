package app

import (
	"strings"
)

// substituteVarsInStaticFiles loops through the values/secrets files and substitutes variables into them.
func (r *Release) substituteVarsInStaticFiles() {
	if r.ValuesFile != "" {
		r.ValuesFile = substituteVarsInYaml(r.ValuesFile)
	}
	if r.SecretsFile != "" {
		r.SecretsFile = substituteVarsInYaml(r.SecretsFile)
	}

	for i := range r.ValuesFiles {
		r.ValuesFiles[i] = substituteVarsInYaml(r.ValuesFiles[i])
	}
	for i := range r.SecretsFiles {
		r.SecretsFiles[i] = substituteVarsInYaml(r.SecretsFiles[i])
	}

	for key, val := range r.Hooks {
		if key != "deleteOnSuccess" && key != "successTimeout" && key != "successCondition" {
			hook := val.(string)
			if isOfType(hook, []string{".yaml", ".yml"}) {
				r.Hooks[key] = substituteVarsInYaml(hook)
			}
		}
	}
}

// resolvePaths resolves relative paths of certs/keys/chart/value file/secret files/etc and replace them with a absolute paths
func (r *Release) resolvePaths(dir, downloadDest string) {
	if r.ValuesFile != "" {
		r.ValuesFile, _ = resolveOnePath(r.ValuesFile, dir, downloadDest)
	}
	if r.SecretsFile != "" {
		r.SecretsFile, _ = resolveOnePath(r.SecretsFile, dir, downloadDest)
	}

	for i, file := range r.ValuesFiles {
		r.ValuesFiles[i], _ = resolveOnePath(file, dir, downloadDest)
	}
	for i, file := range r.SecretsFiles {
		r.SecretsFiles[i], _ = resolveOnePath(file, dir, downloadDest)
	}

	for key, val := range r.Hooks {
		if key != "deleteOnSuccess" && key != "successTimeout" && key != "successCondition" {
			hook := strings.Fields(val.(string))
			if isOfType(hook[0], validHookFiles) && !ToolExists(hook[0]) {
				hook[0], _ = resolveOnePath(hook[0], dir, downloadDest)
				r.Hooks[key] = strings.Join(hook, " ")
			}
		}
	}
}

// getValuesFiles return partial install/upgrade release command to substitute the -f flag in Helm.
func (r *Release) getValuesFiles() []string {
	var fileList []string

	if r.ValuesFile != "" {
		fileList = append(fileList, r.ValuesFile)
	} else if len(r.ValuesFiles) > 0 {
		fileList = append(fileList, r.ValuesFiles...)
	}

	if r.SecretsFile != "" || len(r.SecretsFiles) > 0 {
		if settings.EyamlEnabled {
			if !ToolExists("eyaml") {
				log.Fatal("hiera-eyaml is not installed/configured correctly. Aborting!")
			}
		} else if settings.VaultEnabled {
			if !helmPluginExists("vault") {
				log.Fatal("helm vault plugin is not installed/configured correctly. Aborting!")
			}
		} else {
			if !helmPluginExists("secrets") {
				log.Fatal("helm secrets plugin is not installed/configured correctly. Aborting!")
			}
		}
	}
	if r.SecretsFile != "" {
		if !isOfType(r.SecretsFile, []string{".dec"}) {
			if err := decryptSecret(r.SecretsFile); err != nil {
				log.Fatal(err.Error())
			}
			r.SecretsFile += ".dec"
		}
		fileList = append(fileList, r.SecretsFile)
	} else if len(r.SecretsFiles) > 0 {
		for i := 0; i < len(r.SecretsFiles); i++ {
			if isOfType(r.SecretsFiles[i], []string{".dec"}) {
				// if .dec extension is added before to the secret filename, don't add it again.
				// This happens at upgrade time (where diff and upgrade both call this function)
				// and we don't need to decrypt the file again
				continue
			}

			if err := decryptSecret(r.SecretsFiles[i]); err != nil {
				log.Fatal(err.Error())
			}
			r.SecretsFiles[i] += ".dec"
		}
		fileList = append(fileList, r.SecretsFiles...)
	}

	fileListArgs := []string{}
	for _, file := range fileList {
		fileListArgs = append(fileListArgs, "-f", file)
	}
	return fileListArgs
}
