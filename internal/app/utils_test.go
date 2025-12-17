package app

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name     string
		expected func(key string) string
		key      string
	}{
		{
			name: "direct",
			key:  "BAR",
			expected: func(key string) string {
				expected := "myValue"
				os.Setenv(key, expected)
				return expected
			},
		},
		{
			name: "string_with_var",
			key:  "BAR",
			expected: func(key string) string {
				expected := "contains myValue"
				os.Setenv("FOO", "myValue")
				os.Setenv(key, "contains ${FOO}")
				return expected
			},
		},
		{
			name: "nested_one_level",
			key:  "BAR",
			expected: func(key string) string {
				expected := "myValue"
				os.Setenv("FOO", expected)
				os.Setenv(key, "${FOO}")
				return expected
			},
		},
		{
			name: "nested_two_levels",
			key:  "BAR",
			expected: func(key string) string {
				expected := "myValue"
				os.Setenv("FOZ", expected)
				os.Setenv("FOO", "$FOZ")
				os.Setenv(key, "${FOO}")
				return expected
			},
		},
	}
	for _, tt := range tests {
		expected := tt.expected(tt.key)
		value := getEnv(tt.key)
		if value != expected {
			t.Errorf("getEnv() - unexpected value: wanted: %s got: %s", expected, value)
		}
	}
}

func TestOciRefToFilename(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "no_repo",
			in:   "my-chart:1.2.3",
			want: "my-chart-1.2.3.tgz",
		},
		{
			name: "two_colons",
			in:   "my:chart:1.2.3",
			want: "my:chart-1.2.3.tgz",
		},
		{
			name: "with_Host",
			in:   "my-repo.example.com/charts/my-chart:1.2.3",
			want: "my-chart-1.2.3.tgz",
		},
		{
			name: "full_url",
			in:   "oci://my-repo.example.com/charts/my-chart:1.2.3",
			want: "my-chart-1.2.3.tgz",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ociRefToFilename(tt.in)
			if err != nil {
				t.Errorf("ociRefToFilename() unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("ociRefToFilename() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isOfType(t *testing.T) {
	type args struct {
		filename  string
		filetypes []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test case 1 -- valid xml check",
			args: args{
				filename:  "name.xml",
				filetypes: []string{".xml"},
			},
			want: true,
		}, {
			name: "test case 2 -- valid yaml check",
			args: args{
				filename:  "another_name.yaml",
				filetypes: []string{".yaml", ".yml"},
			},
			want: true,
		}, {
			name: "test case 3 -- valid (short) yaml check",
			args: args{
				filename:  "another_name.yml",
				filetypes: []string{".yaml", ".yml"},
			},
			want: true,
		}, {
			name: "test case 4 -- invalid yaml check",
			args: args{
				filename:  "name.xml",
				filetypes: []string{".yaml", ".yml"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isOfType(tt.args.filename, tt.args.filetypes); got != tt.want {
				t.Errorf("isOfType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readFile(t *testing.T) {
	type args struct {
		filepath string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test case 1 -- successful reading.",
			args: args{
				filepath: "../../tests/values.yaml",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := readFile(tt.args.filepath); got != tt.want {
				t.Errorf("readFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_eyamlSecrets(t *testing.T) {
	type args struct {
		r *Release
		s *Config
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "decryptSecrets - valid eyaml-based secrets decryption",
			args: args{
				s: &Config{
					EyamlEnabled:        true,
					EyamlPublicKeyPath:  "./../../tests/keys/public_key.pkcs7.pem",
					EyamlPrivateKeyPath: "./../../tests/keys/private_key.pkcs7.pem",
				},
				r: &Release{
					Name:        "release1",
					Namespace:   "namespace",
					Version:     "1.0.0",
					Enabled:     True,
					SecretsFile: "./../../tests/secrets/valid_eyaml_secrets.yaml",
				},
			},
			want: true,
		},
		{
			name: "decryptSecrets - not existing eyaml-based secrets file",
			args: args{
				s: &Config{
					EyamlEnabled:        true,
					EyamlPublicKeyPath:  "./../../tests/keys/public_key.pkcs7.pem",
					EyamlPrivateKeyPath: "./../../tests/keys/private_key.pkcs7.pem",
				},
				r: &Release{
					Name:        "release1",
					Namespace:   "namespace",
					Version:     "1.0.0",
					Enabled:     True,
					SecretsFile: "./../../tests/secrets/invalid_eyaml_secrets.yaml",
				},
			},
			want: false,
		},
		{
			name: "decryptSecrets - not existing private eyaml key",
			args: args{
				s: &Config{
					EyamlEnabled:        true,
					EyamlPublicKeyPath:  "./../../tests/keys/public_key.pkcs7.pem",
					EyamlPrivateKeyPath: "./../../tests/keys/invalid_private_key.pkcs7.pem",
				},
				r: &Release{
					Name:        "release1",
					Namespace:   "namespace",
					Version:     "1.0.0",
					Enabled:     True,
					SecretsFile: "./../../tests/secrets/valid_eyaml_secrets.yaml",
				},
			},
			want: false,
		},
		{
			name: "decryptSecrets - not existing public eyaml key",
			args: args{
				s: &Config{
					EyamlEnabled: true,
					// https://github.com/voxpupuli/hiera-eyaml/commit/760fd05e7cbb34b5380f87a87290deb790ae0aaf
					// Public key is only needed for encryption, not decryption
					EyamlPublicKeyPath:  "./../../tests/keys/invalid_public_key.pkcs7.pem",
					EyamlPrivateKeyPath: "./../../tests/keys/private_key.pkcs7.pem",
				},
				r: &Release{
					Name:        "release1",
					Namespace:   "namespace",
					Version:     "1.0.0",
					Enabled:     True,
					SecretsFile: "./../../tests/secrets/valid_eyaml_secrets.yaml",
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(tt.want)
			settings = &Config{}
			settings.EyamlEnabled = tt.args.s.EyamlEnabled
			settings.EyamlPublicKeyPath = tt.args.s.EyamlPublicKeyPath
			settings.EyamlPrivateKeyPath = tt.args.s.EyamlPrivateKeyPath
			err := decryptSecret(tt.args.r.SecretsFile)
			switch err.(type) {
			case nil:
				if tt.want != true {
					t.Errorf("decryptSecret() = %v, want error", err)
				}
			case error:
				if tt.want != false {
					t.Errorf("decryptSecret() = %v, want nil", err)
				}
			}
			if _, err := os.Stat(tt.args.r.SecretsFile + ".dec"); err == nil {
				defer deleteFile(tt.args.r.SecretsFile + ".dec")
			}
		})
	}
}
