package app

import "testing"

var _ = func() bool {
	testing.Init()
	return true
}()

func init() {
	flags.parse()
}

func Test_readState(t *testing.T) {
	type result struct {
		numApps        int
		numNSs         int
		numEnabledApps int
		numEnabledNSs  int
	}
	tests := []struct {
		name  string
		flags cli
		want  result
	}{
		{
			name: "yaml minimal example; no validation",
			flags: cli{
				files:          fileOptionArray([]fileOption{{"../../examples/minimal-example.yaml", 0}}),
				skipValidation: true,
			},
			want: result{
				numApps:        2,
				numNSs:         1,
				numEnabledApps: 2,
				numEnabledNSs:  1,
			},
		},
		{
			name: "toml minimal example; no validation",
			flags: cli{
				files:          fileOptionArray([]fileOption{{"../../examples/minimal-example.toml", 0}}),
				skipValidation: true,
			},
			want: result{
				numApps:        2,
				numNSs:         1,
				numEnabledApps: 2,
				numEnabledNSs:  1,
			},
		},
		{
			name: "yaml minimal example; no validation with bad target",
			flags: cli{
				target:         stringArray([]string{"foo"}),
				files:          fileOptionArray([]fileOption{{"../../examples/minimal-example.yaml", 0}}),
				skipValidation: true,
			},
			want: result{
				numApps:        2,
				numNSs:         1,
				numEnabledApps: 0,
				numEnabledNSs:  0,
			},
		},
		{
			name: "yaml minimal example; no validation; target jenkins",
			flags: cli{
				target:         stringArray([]string{"jenkins"}),
				files:          fileOptionArray([]fileOption{{"../../examples/minimal-example.yaml", 0}}),
				skipValidation: true,
			},
			want: result{
				numApps:        2,
				numNSs:         1,
				numEnabledApps: 1,
				numEnabledNSs:  1,
			},
		},
		{
			name: "yaml and toml minimal examples merged; no validation",
			flags: cli{
				files:          fileOptionArray([]fileOption{{"../../examples/minimal-example.yaml", 0}, {"../../examples/minimal-example.toml", 0}}),
				skipValidation: true,
			},
			want: result{
				numApps:        2,
				numNSs:         1,
				numEnabledApps: 2,
				numEnabledNSs:  1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := State{}
			if err := tt.flags.readState(&s); err != nil {
				t.Errorf("readState() = Unexpected error: %v", err)
			}
			if len(s.Apps) != tt.want.numApps {
				t.Errorf("readState() = app count mismatch: want: %d, got: %d", tt.want.numApps, len(s.Apps))
			}
			if len(s.Namespaces) != tt.want.numNSs {
				t.Errorf("readState() = NS count mismatch: want: %d, got: %d", tt.want.numNSs, len(s.Namespaces))
			}

			var enabledApps, enabledNSs int
			for _, a := range s.Apps {
				if !a.disabled {
					enabledApps++
				}
			}
			if enabledApps != tt.want.numEnabledApps {
				t.Errorf("readState() = app count mismatch: want: %d, got: %d", tt.want.numEnabledApps, enabledApps)
			}
			for _, n := range s.Namespaces {
				if !n.disabled {
					enabledNSs++
				}
			}
			if enabledNSs != tt.want.numEnabledNSs {
				t.Errorf("readState() = app count mismatch: want: %d, got: %d", tt.want.numEnabledNSs, enabledNSs)
			}
		})
	}
}
