package app

import (
	"os"
	"reflect"
	"testing"
)

func setupStateFileTestCase(t *testing.T) (func(t *testing.T), error) {
	t.Log("setup test case")
	if err := os.MkdirAll(tempFilesDir, 0o755); err != nil {
		t.Errorf("setupStateFileTestCase(), failed to create temp files dir: %v", err)
		return nil, err
	}
	flags.tempDir = tempFilesDir

	return func(t *testing.T) {
		t.Log("teardown test case")
		os.RemoveAll(tempFilesDir)
		flags.tempDir = ""
	}, nil
}

func Test_fromTOML(t *testing.T) {
	type args struct {
		file string
		s    *State
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test case 1 -- invalid TOML",
			args: args{
				file: "../../../tests/invalid_example.toml",
				s:    new(State),
			},
			want: false,
		}, {
			name: "test case 2 -- valid TOML",
			args: args{
				file: "../../examples/example.toml",
				s:    new(State),
			},
			want: true,
		},
	}
	os.Setenv("ORG_PATH", "sample")
	os.Setenv("VALUE", "sample")

	teardownTestCase, err := setupStateFileTestCase(t)
	if err != nil {
		t.Errorf("setupStateFileTestCase(), got: %v", err)
	}
	defer teardownTestCase(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.s.fromTOML(tt.args.file)
			got := err == nil
			if got != tt.want {
				t.Errorf("fromToml() = %v, want %v", got, tt.want)
			}
		})
	}
	os.Unsetenv("ORG_PATH")
	os.Unsetenv("VALUE")
}

func Test_fromTOML_Expand(t *testing.T) {
	type args struct {
		file string
		s    *State
	}
	tests := []struct {
		name    string
		args    args
		section string
		field   string
		want    string
	}{
		{
			name: "test case 1 -- valid TOML expand ClusterURI",
			args: args{
				file: "../../examples/example.toml",
				s:    new(State),
			},
			section: "Settings",
			field:   "ClusterURI",
			want:    "https://192.168.99.100:8443",
		},
		{
			name: "test case 2 -- valid TOML expand org",
			args: args{
				file: "../../examples/example.toml",
				s:    new(State),
			},
			section: "Metadata",
			field:   "org",
			want:    "example.com/sample/",
		},
	}
	os.Setenv("SET_URI", "https://192.168.99.100:8443")
	os.Setenv("ORG_PATH", "sample")
	os.Setenv("VALUE", "sample")

	teardownTestCase, err := setupStateFileTestCase(t)
	if err != nil {
		t.Errorf("setupStateFileTestCase(), got: %v", err)
	}
	defer teardownTestCase(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.s.fromTOML(tt.args.file)
			if err != nil {
				t.Errorf("fromToml(), got: %v", err)
			}

			tomlVal := reflect.ValueOf(tt.args.s).Elem()
			tomlType := reflect.TypeOf(tt.args.s)

			if tomlType.Kind() != reflect.Struct {

				section := tomlVal.FieldByName(tt.section)
				sectionType := reflect.TypeOf(section)

				if section.IsValid() && section.Kind() == reflect.Struct {
					field := section.FieldByName(tt.field)
					if sectionType.Kind() == reflect.String {
						if field.String() != tt.want {
							t.Errorf("fromToml().section.field = %v, got: %v", tt.want, field.String())
						}
					}
				} else if section.IsValid() && section.Kind() == reflect.Map {
					found := false
					value := ""
					for _, key := range section.MapKeys() {
						if key.String() == tt.field {
							found = true
							value = section.MapIndex(key).String()
						}
					}
					if !found {
						t.Errorf("fromToml().section.field = '%v' not found", tt.field)
					} else if value != tt.want {
						t.Errorf("fromToml().section.field = %v, got: %v", tt.want, value)
					}

				} else {
					t.Errorf("fromToml().section = struct, got: %v", sectionType.Kind())
				}

			} else {
				t.Errorf("fromToml() = struct, got: %v", tomlType.Kind())
			}
		})
	}
	os.Unsetenv("ORG_PATH")
	os.Unsetenv("SET_URI")
	os.Unsetenv("VALUE")
}

func Test_fromYAML(t *testing.T) {
	type args struct {
		file string
		s    *State
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test case 1 -- invalid YAML",
			args: args{
				file: "../../tests/invalid_example.yaml",
				s:    new(State),
			},
			want: false,
		}, {
			name: "test case 2 -- valid TOML",
			args: args{
				file: "../../examples/example.yaml",
				s:    new(State),
			},
			want: true,
		},
	}
	os.Setenv("VALUE", "sample")
	os.Setenv("ORG_PATH", "sample")

	teardownTestCase, err := setupStateFileTestCase(t)
	if err != nil {
		t.Errorf("setupStateFileTestCase(), got: %v", err)
	}
	defer teardownTestCase(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.s.fromYAML(tt.args.file)
			got := err == nil
			if got != tt.want {
				t.Errorf("fromYaml() = %v, want %v", got, tt.want)
			}
		})
	}
	os.Unsetenv("ORG_PATH")
	os.Unsetenv("VALUE")
}

func Test_fromYAML_UnsetVars(t *testing.T) {
	type args struct {
		file string
		s    *State
	}
	tests := []struct {
		name      string
		args      args
		targetVar string
		want      bool
	}{
		{
			name: "test case 1 -- unset ORG_PATH env var",
			args: args{
				file: "../../examples/example.yaml",
				s:    new(State),
			},
			targetVar: "ORG_PATH",
			want:      false,
		},
		{
			name: "test case 2 -- unset VALUE var",
			args: args{
				file: "../../examples/example.yaml",
				s:    new(State),
			},
			targetVar: "VALUE",
			want:      false,
		},
	}

	teardownTestCase, err := setupStateFileTestCase(t)
	if err != nil {
		t.Errorf("setupStateFileTestCase(), got: %v", err)
	}
	defer teardownTestCase(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.targetVar {
			case "ORG_PATH":
				os.Setenv("VALUE", "sample")
			case "VALUE":
				os.Setenv("ORG_PATH", "sample")
			}
			err := tt.args.s.fromYAML(tt.args.file)
			got := err == nil
			if got != tt.want {
				t.Errorf("fromYaml() = %v, want %v", got, tt.want)
			}
		})
		os.Unsetenv("ORG_PATH")
		os.Unsetenv("VALUE")
	}
}

func Test_fromYAML_Expand(t *testing.T) {
	type args struct {
		file string
		s    *State
	}
	tests := []struct {
		name    string
		args    args
		section string
		field   string
		want    string
	}{
		{
			name: "test case 1 -- valid YAML expand ClusterURI",
			args: args{
				file: "../../examples/example.yaml",
				s:    new(State),
			},
			section: "Settings",
			field:   "ClusterURI",
			want:    "https://192.168.99.100:8443",
		},
		{
			name: "test case 2 -- valid YAML expand org",
			args: args{
				file: "../../examples/example.yaml",
				s:    new(State),
			},
			section: "Metadata",
			field:   "org",
			want:    "example.com/sample/",
		},
	}
	os.Setenv("SET_URI", "https://192.168.99.100:8443")
	os.Setenv("ORG_PATH", "sample")
	os.Setenv("VALUE", "sample")

	teardownTestCase, err := setupStateFileTestCase(t)
	if err != nil {
		t.Errorf("setupStateFileTestCase(), got: %v", err)
	}
	defer teardownTestCase(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.s.fromYAML(tt.args.file)
			if err != nil {
				t.Errorf("fromYaml(), got: %v", err)
			}

			yamlVal := reflect.ValueOf(tt.args.s).Elem()
			yamlType := reflect.TypeOf(tt.args.s)

			if yamlType.Kind() != reflect.Struct {

				section := yamlVal.FieldByName(tt.section)
				sectionType := reflect.TypeOf(section)

				if section.IsValid() && section.Kind() == reflect.Struct {
					field := section.FieldByName(tt.field)
					if sectionType.Kind() == reflect.String {
						if field.String() != tt.want {
							t.Errorf("fromYaml().section.field = %v, got: %v", tt.want, field.String())
						}
					}
				} else if section.IsValid() && section.Kind() == reflect.Map {
					found := false
					value := ""
					for _, key := range section.MapKeys() {
						if key.String() == tt.field {
							found = true
							value = section.MapIndex(key).String()
						}
					}
					if !found {
						t.Errorf("fromYaml().section.field = '%v' not found", tt.field)
					} else if value != tt.want {
						t.Errorf("fromYaml().section.field = %v, got: %v", tt.want, value)
					}

				} else {
					t.Errorf("fromYaml().section = struct, got: %v", sectionType.Kind())
				}

			} else {
				t.Errorf("fromYaml() = struct, got: %v", yamlType.Kind())
			}
		})
	}
	os.Unsetenv("ORG_PATH")
	os.Unsetenv("SET_URI")
	os.Unsetenv("VALUE")
}

func Test_build(t *testing.T) {
	teardownTestCase, err := setupStateFileTestCase(t)
	if err != nil {
		t.Errorf("setupStateFileTestCase(), got: %v", err)
	}
	defer teardownTestCase(t)
	s := new(State)
	files := fileOptionArray{
		fileOption{name: "../../examples/composition/main.yaml"},
		fileOption{name: "../../examples/composition/kyverno.yaml"},
		fileOption{name: "../../examples/composition/argo.yaml"},
		fileOption{name: "../../examples/composition/artifactory.yaml"},
	}
	err = s.build(files)
	if err != nil {
		t.Errorf("build() - unexpected error: %v", err)
	}
	if len(s.Apps) != 5 {
		t.Errorf("build() - unexpected number of apps, wanted 5 got %d", len(s.Apps))
	}
	if len(s.HelmRepos) != 3 {
		t.Errorf("build() - unexpected number of repos, wanted 3 got %d", len(s.Apps))
	}
}

func Test_DSFMergeWithOverwrite(t *testing.T) {
	teardownTestCase, err := setupStateFileTestCase(t)
	if err != nil {
		t.Errorf("setupStateFileTestCase(), got: %v", err)
	}
	defer teardownTestCase(t)
	s := new(State)
	files := fileOptionArray{
		fileOption{name: "../../examples/minimal-example.yaml"},
		fileOption{name: "../../examples/minimal-example-overwrite.yaml"},
	}
	err = s.build(files)
	if err != nil {
		t.Errorf("build() - unexpected error: %v", err)
	}
	if len(s.Apps) != 2 {
		t.Errorf("build() - unexpected number of apps, wanted 5 got %d", len(s.Apps))
	}
	if len(s.HelmRepos) != 2 {
		t.Errorf("build() - unexpected number of repos, wanted 3 got %d", len(s.Apps))
	}
	if s.Apps["jenkins"].Enabled.Value != false {
		t.Errorf("build() - unexpected status of a release, wanted 'enabled'=false got %v", s.Apps["jenkins"].Enabled.Value)
	}
}
