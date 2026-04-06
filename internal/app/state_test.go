package app

import (
	"os"
	"testing"
)

func setupTestCase(t *testing.T) (func(t *testing.T), error) {
	t.Log("setup test case")
	if err := os.MkdirAll(tempFilesDir, 0o755); err != nil {
		return nil, err
	}
	flags.tempDir = tempFilesDir
	if err := os.MkdirAll(os.TempDir()+"/helmsman-tests/myapp", os.ModePerm); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(os.TempDir()+"/helmsman-tests/dir-with space/myapp", os.ModePerm); err != nil {
		return nil, err
	}
	cmd := helmCmd([]string{"create", os.TempDir() + "/helmsman-tests/dir-with space/myapp"}, "creating an empty local chart directory")
	if _, err := cmd.Exec(); err != nil {
		log.Fatalf("Command failed: %v", err)
	}

	return func(t *testing.T) {
		t.Log("teardown test case")
		os.RemoveAll(tempFilesDir)
		os.RemoveAll(os.TempDir() + "/helmsman-tests/")
		flags.tempDir = ""
	}, nil
}

func Test_state_validate(t *testing.T) {
	type fields struct {
		Metadata     map[string]string
		Certificates map[string]string
		Settings     Config
		Namespaces   map[string]*Namespace
		HelmRepos    map[string]string
		Apps         map[string]*Release
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "test case 1",
			fields: fields{
				Metadata: make(map[string]string),
				Certificates: map[string]string{
					"caCrt": "s3://some-bucket/12345.crt",
					"caKey": "s3://some-bucket/12345.key",
				},
				Settings: Config{
					KubeContext: "minikube",
					Username:    "admin",
					Password:    "$K8S_PASSWORD",
					ClusterURI:  "https://192.168.99.100:8443",
				},
				Namespaces: map[string]*Namespace{
					"staging": {false, Limits{}, make(map[string]string), make(map[string]string), &Quotas{}, false},
				},
				HelmRepos: map[string]string{
					"deprecated-stable": "https://kubernetes-charts.storage.googleapis.com",
					"myrepo":            "s3://my-repo/charts",
				},
				Apps: make(map[string]*Release),
			},
			want: true,
		}, {
			name: "test case 2 -- settings/empty_context",
			fields: fields{
				Metadata: make(map[string]string),
				Certificates: map[string]string{
					"caCrt": "s3://some-bucket/12345.crt",
					"caKey": "s3://some-bucket/12345.key",
				},
				Settings: Config{
					KubeContext: "",
					Username:    "admin",
					Password:    "$K8S_PASSWORD",
					ClusterURI:  "https://192.168.99.100:8443",
				},
				Namespaces: map[string]*Namespace{
					"staging": {false, Limits{}, make(map[string]string), make(map[string]string), &Quotas{}, false},
				},
				HelmRepos: map[string]string{
					"deprecated-stable": "https://kubernetes-charts.storage.googleapis.com",
					"myrepo":            "s3://my-repo/charts",
				},
				Apps: make(map[string]*Release),
			},
			want: false,
		}, {
			name: "test case 3 -- settings/optional_params",
			fields: fields{
				Metadata: make(map[string]string),
				Certificates: map[string]string{
					"caCrt": "s3://some-bucket/12345.crt",
					"caKey": "s3://some-bucket/12345.key",
				},
				Settings: Config{
					KubeContext: "minikube",
				},
				Namespaces: map[string]*Namespace{
					"staging": {false, Limits{}, make(map[string]string), make(map[string]string), &Quotas{}, false},
				},
				HelmRepos: map[string]string{
					"deprecated-stable": "https://kubernetes-charts.storage.googleapis.com",
					"myrepo":            "s3://my-repo/charts",
				},
				Apps: make(map[string]*Release),
			},
			want: true,
		}, {
			name: "test case 4 -- settings/password-passed-directly",
			fields: fields{
				Metadata: make(map[string]string),
				Certificates: map[string]string{
					"caCrt": "s3://some-bucket/12345.crt",
					"caKey": "s3://some-bucket/12345.key",
				},
				Settings: Config{
					KubeContext: "minikube",
					Username:    "admin",
					Password:    "K8S_PASSWORD",
					ClusterURI:  "https://192.168.99.100:8443",
				},
				Namespaces: map[string]*Namespace{
					"staging": {false, Limits{}, make(map[string]string), make(map[string]string), &Quotas{}, false},
				},
				HelmRepos: map[string]string{
					"deprecated-stable": "https://kubernetes-charts.storage.googleapis.com",
					"myrepo":            "s3://my-repo/charts",
				},
				Apps: make(map[string]*Release),
			},
			want: true,
		}, {
			name: "test case 5 -- settings/clusterURI-empty-env-var",
			fields: fields{
				Metadata: make(map[string]string),
				Certificates: map[string]string{
					"caCrt": "s3://some-bucket/12345.crt",
					"caKey": "s3://some-bucket/12345.key",
				},
				Settings: Config{
					KubeContext: "minikube",
					Username:    "admin",
					Password:    "K8S_PASSWORD",
					ClusterURI:  "$URI", // unset env
				},
				Namespaces: map[string]*Namespace{
					"staging": {false, Limits{}, make(map[string]string), make(map[string]string), &Quotas{}, false},
				},
				HelmRepos: map[string]string{
					"deprecated-stable": "https://kubernetes-charts.storage.googleapis.com",
					"myrepo":            "s3://my-repo/charts",
				},
				Apps: make(map[string]*Release),
			},
			want: false,
		}, {
			name: "test case 6 -- settings/clusterURI-invalid",
			fields: fields{
				Metadata: make(map[string]string),
				Certificates: map[string]string{
					"caCrt": "s3://some-bucket/12345.crt",
					"caKey": "s3://some-bucket/12345.key",
				},
				Settings: Config{
					KubeContext: "minikube",
					Username:    "admin",
					Password:    "K8S_PASSWORD",
					ClusterURI:  "https//192.168.99.100:8443", // invalid url
				},
				Namespaces: map[string]*Namespace{
					"staging": {false, Limits{}, make(map[string]string), make(map[string]string), &Quotas{}, false},
				},
				HelmRepos: map[string]string{
					"deprecated-stable": "https://kubernetes-charts.storage.googleapis.com",
					"myrepo":            "s3://my-repo/charts",
				},
				Apps: make(map[string]*Release),
			},
			want: false,
		}, {
			name: "test case 7 -- certifications/missing key",
			fields: fields{
				Metadata: make(map[string]string),
				Certificates: map[string]string{
					"caCrt": "s3://some-bucket/12345.crt",
				},
				Settings: Config{
					KubeContext: "minikube",
					Username:    "admin",
					Password:    "$K8S_PASSWORD",
					ClusterURI:  "https://192.168.99.100:8443",
				},
				Namespaces: map[string]*Namespace{
					"staging": {false, Limits{}, make(map[string]string), make(map[string]string), &Quotas{}, false},
				},
				HelmRepos: map[string]string{
					"deprecated-stable": "https://kubernetes-charts.storage.googleapis.com",
					"myrepo":            "s3://my-repo/charts",
				},
				Apps: make(map[string]*Release),
			},
			want: false,
		}, {
			name: "test case 8 -- certifications/nil_value",
			fields: fields{
				Metadata:     make(map[string]string),
				Certificates: nil,
				Settings: Config{
					KubeContext: "minikube",
					Username:    "admin",
					Password:    "$K8S_PASSWORD",
					ClusterURI:  "https://192.168.99.100:8443",
				},
				Namespaces: map[string]*Namespace{
					"staging": {false, Limits{}, make(map[string]string), make(map[string]string), &Quotas{}, false},
				},
				HelmRepos: map[string]string{
					"deprecated-stable": "https://kubernetes-charts.storage.googleapis.com",
					"myrepo":            "s3://my-repo/charts",
				},
				Apps: make(map[string]*Release),
			},
			want: false,
		}, {
			name: "test case 9 -- certifications/invalid_s3",
			fields: fields{
				Metadata: make(map[string]string),
				Certificates: map[string]string{
					"caCrt": "s3://some-bucket/12345.crt",
					"caKey": "http://someurl.com/",
				},
				Settings: Config{
					KubeContext: "minikube",
					Username:    "admin",
					Password:    "$K8S_PASSWORD",
					ClusterURI:  "https://192.168.99.100:8443",
				},
				Namespaces: map[string]*Namespace{
					"staging": {false, Limits{}, make(map[string]string), make(map[string]string), &Quotas{}, false},
				},
				HelmRepos: map[string]string{
					"deprecated-stable": "https://kubernetes-charts.storage.googleapis.com",
					"myrepo":            "s3://my-repo/charts",
				},
				Apps: make(map[string]*Release),
			},
			want: false,
		}, {
			name: "test case 10 -- certifications/nil_value_pass",
			fields: fields{
				Metadata:     make(map[string]string),
				Certificates: nil,
				Settings: Config{
					KubeContext: "minikube",
				},
				Namespaces: map[string]*Namespace{
					"staging": {false, Limits{}, make(map[string]string), make(map[string]string), &Quotas{}, false},
				},
				HelmRepos: map[string]string{
					"deprecated-stable": "https://kubernetes-charts.storage.googleapis.com",
					"myrepo":            "s3://my-repo/charts",
				},
				Apps: make(map[string]*Release),
			},
			want: true,
		}, {
			name: "test case 11 -- namespaces/nil_value",
			fields: fields{
				Metadata:     make(map[string]string),
				Certificates: nil,
				Settings: Config{
					KubeContext: "minikube",
				},
				Namespaces: nil,
				HelmRepos: map[string]string{
					"deprecated-stable": "https://kubernetes-charts.storage.googleapis.com",
					"myrepo":            "s3://my-repo/charts",
				},
				Apps: make(map[string]*Release),
			},
			want: false,
		}, {
			name: "test case 12 -- namespaces/empty",
			fields: fields{
				Metadata:     make(map[string]string),
				Certificates: nil,
				Settings: Config{
					KubeContext: "minikube",
				},
				Namespaces: map[string]*Namespace{},
				HelmRepos: map[string]string{
					"deprecated-stable": "https://kubernetes-charts.storage.googleapis.com",
					"myrepo":            "s3://my-repo/charts",
				},
				Apps: make(map[string]*Release),
			},
			want: false,
		}, {
			name: "test case 17 -- helmRepos/nil_value",
			fields: fields{
				Metadata:     make(map[string]string),
				Certificates: nil,
				Settings: Config{
					KubeContext: "minikube",
				},
				Namespaces: map[string]*Namespace{
					"staging": {false, Limits{}, make(map[string]string), make(map[string]string), &Quotas{}, false},
				},
				HelmRepos: nil,
				Apps:      make(map[string]*Release),
			},
			want: true,
		}, {
			name: "test case 18 -- helmRepos/empty",
			fields: fields{
				Metadata:     make(map[string]string),
				Certificates: nil,
				Settings: Config{
					KubeContext: "minikube",
				},
				Namespaces: map[string]*Namespace{
					"staging": {false, Limits{}, make(map[string]string), make(map[string]string), &Quotas{}, false},
				},
				HelmRepos: map[string]string{},
				Apps:      make(map[string]*Release),
			},
			want: true,
		}, {
			name: "test case 19 -- helmRepos/empty_repo_value",
			fields: fields{
				Metadata:     make(map[string]string),
				Certificates: nil,
				Settings: Config{
					KubeContext: "minikube",
				},
				Namespaces: map[string]*Namespace{
					"staging": {false, Limits{}, make(map[string]string), make(map[string]string), &Quotas{}, false},
				},
				HelmRepos: map[string]string{
					"deprecated-stable": "https://kubernetes-charts.storage.googleapis.com",
					"myrepo":            "",
				},
				Apps: make(map[string]*Release),
			},
			want: false,
		}, {
			name: "test case 20 -- helmRepos/invalid_repo_value",
			fields: fields{
				Metadata:     make(map[string]string),
				Certificates: nil,
				Settings: Config{
					KubeContext: "minikube",
				},
				Namespaces: map[string]*Namespace{
					"staging": {false, Limits{}, make(map[string]string), make(map[string]string), &Quotas{}, false},
				},
				HelmRepos: map[string]string{
					"deprecated-stable": "https://kubernetes-charts.storage.googleapis.com",
					"myrepo":            "s3//my-repo/charts",
				},
				Apps: make(map[string]*Release),
			},
			want: false,
		},
	}
	os.Setenv("K8S_PASSWORD", "my-fake-password")
	os.Setenv("SET_URI", "https://192.168.99.100:8443")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := State{
				Metadata:     tt.fields.Metadata,
				Certificates: tt.fields.Certificates,
				Settings:     tt.fields.Settings,
				Namespaces:   tt.fields.Namespaces,
				HelmRepos:    tt.fields.HelmRepos,
				Apps:         tt.fields.Apps,
			}
			err := s.validate()
			switch err.(type) {
			case nil:
				if tt.want != true {
					t.Errorf("state.validate() = %v, want error", err)
				}
			case error:
				if tt.want != false {
					t.Errorf("state.validate() = %v, want nil", err)
				}
			}
		})
	}
}

func createFullReleasePointer(name, chart, version string) *Release {
	return &Release{
		Name:          name,
		Description:   "",
		Namespace:     "",
		Enabled:       True,
		Chart:         chart,
		Version:       version,
		ValuesFile:    "",
		ValuesFiles:   []string{},
		SecretsFile:   "",
		SecretsFiles:  []string{},
		Test:          False,
		Protected:     False,
		Wait:          False,
		Priority:      0,
		Set:           make(map[string]string),
		SetString:     make(map[string]string),
		HelmFlags:     []string{},
		HelmDiffFlags: []string{},
		NoHooks:       False,
		Timeout:       0,
		PostRenderer:  "",
	}
}

func Test_state_getReleaseChartsInfo(t *testing.T) {
	type args struct {
		apps map[string]*Release
	}

	tests := []struct {
		name       string
		targetFlag []string
		groupFlag  []string
		args       args
		want       bool
	}{
		{
			name: "test case 1: valid local path with no chart",
			args: args{
				apps: map[string]*Release{
					"app": createFullReleasePointer("app", os.TempDir()+"/helmsman-tests/myapp", ""),
				},
			},
			want: false,
		}, {
			name: "test case 2: invalid local path",
			args: args{
				apps: map[string]*Release{
					"app": createFullReleasePointer("app", os.TempDir()+"/does-not-exist/myapp", ""),
				},
			},
			want: false,
		}, {
			name: "test case 3: valid chart local path with whitespace",
			args: args{
				apps: map[string]*Release{
					"app": createFullReleasePointer("app", os.TempDir()+"/helmsman-tests/dir-with space/myapp", "0.1.0"),
				},
			},
			want: true,
		}, {
			name: "test case 4: valid chart from repo",
			args: args{
				apps: map[string]*Release{
					"app": createFullReleasePointer("app", "prometheus-community/prometheus", "11.16.5"),
				},
			},
			want: true,
		}, {
			name:       "test case 5: invalid local path for chart ignored with -target flag, while other app was targeted",
			targetFlag: []string{"notThisOne"},
			args: args{
				apps: map[string]*Release{
					"app": createFullReleasePointer("app", os.TempDir()+"/does-not-exist/myapp", ""),
				},
			},
			want: true,
		}, {
			name:       "test case 6: invalid local path for chart included with -target flag",
			targetFlag: []string{"app"},
			args: args{
				apps: map[string]*Release{
					"app": createFullReleasePointer("app", os.TempDir()+"/does-not-exist/myapp", ""),
				},
			},
			want: false,
		}, {
			name:       "test case 7: multiple valid local apps with the same chart version",
			targetFlag: []string{"app"},
			args: args{
				apps: map[string]*Release{
					"app1": createFullReleasePointer("app2", os.TempDir()+"/helmsman-tests/dir-with space/myapp", "0.1.0"),
					"app2": createFullReleasePointer("app3", os.TempDir()+"/helmsman-tests/dir-with space/myapp", "0.1.0"),
					"app3": createFullReleasePointer("app4", os.TempDir()+"/helmsman-tests/dir-with space/myapp", "0.1.0"),
					"app4": createFullReleasePointer("app5", os.TempDir()+"/helmsman-tests/dir-with space/myapp", "0.1.0"),
					"app5": createFullReleasePointer("app5", os.TempDir()+"/helmsman-tests/dir-with space/myapp", "0.1.0"),
				},
			},
			want: true,
		},
	}

	teardownTestCase, err := setupTestCase(t)
	if err != nil {
		t.Errorf("setupTestCase() = %v", err)
		return
	}
	defer teardownTestCase(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stt := &State{Apps: tt.args.apps}
			stt.disableApps(tt.groupFlag, tt.targetFlag, []string{}, []string{})
			err := stt.getReleaseChartsInfo()
			switch err.(type) {
			case nil:
				if tt.want != true {
					t.Errorf("getReleaseChartsInfo() = %v, want error", err)
				}
			case error:
				if tt.want != false {
					t.Errorf("getReleaseChartsInfo() = %v, want nil", err)
				}
			}
		})
	}
}
