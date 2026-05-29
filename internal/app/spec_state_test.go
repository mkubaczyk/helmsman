package app

import (
	"testing"
)

func Test_specFromYAML(t *testing.T) {
	type args struct {
		file string
		s    *StateFiles
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test case 1 -- Valid YAML",
			args: args{
				file: "../../examples/example-spec.yaml",
				s:    new(StateFiles),
			},
			want: true,
		}, {
			name: "test case 2 -- Invalid Yaml",
			args: args{
				file: "../../tests/Invalid_example_spec.yaml",
				s:    new(StateFiles),
			},
			want: false,
		}, {
			name: "test case 3 -- Commposition example",
			args: args{
				file: "../../examples/composition/spec.yaml",
				s:    new(StateFiles),
			},
			want: true,
		},
	}

	teardownTestCase, err := setupStateFileTestCase(t)
	if err != nil {
		t.Fatal(err)
	}
	defer teardownTestCase(t)
	for _, tt := range tests {
		// os.Args = append(os.Args, "-f ../../examples/example.yaml")
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.s.specFromYAML(tt.args.file)
			if err != nil {
				t.Log(err)
			}

			got := err == nil
			if got != tt.want {
				t.Errorf("specFromYaml() = %v, want %v", got, tt.want)
			}
		})
	}
}
