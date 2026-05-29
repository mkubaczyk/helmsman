package app

import (
	"reflect"
	"testing"
	"time"
)

func Test_createPlan(t *testing.T) {
	tests := []struct {
		name string
		want *plan
	}{
		{
			name: "test creating a plan",
			want: &plan{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createPlan(); reflect.DeepEqual(got, tt.want) {
				t.Errorf("createPlan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_plan_addCommand(t *testing.T) {
	type fields struct {
		Commands  []orderedCommand
		Decisions []orderedDecision
		Created   time.Time
	}
	type args struct {
		c Command
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "testing command 1",
			fields: fields{
				Commands:  []orderedCommand{},
				Decisions: []orderedDecision{},
				Created:   time.Now(),
			},
			args: args{
				c: Command{
					Cmd:         "bash",
					Args:        []string{"-c", "echo this is fun"},
					Description: "A bash command execution test with echo.",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &plan{
				Commands:  tt.fields.Commands,
				Decisions: tt.fields.Decisions,
				Created:   tt.fields.Created,
			}
			r := &Release{}
			p.addCommand(tt.args.c, 0, r, []hookCmd{}, []hookCmd{})
			if got := len(p.Commands); got != 1 {
				t.Errorf("addCommand(): got  %v, want 1", got)
			}
		})
	}
}

func Test_plan_addDecision(t *testing.T) {
	type fields struct {
		Commands  []orderedCommand
		Decisions []orderedDecision
		Created   time.Time
	}
	type args struct {
		decision string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "testing decision adding",
			fields: fields{
				Commands:  []orderedCommand{},
				Decisions: []orderedDecision{},
				Created:   time.Now(),
			},
			args: args{
				decision: "This is a test decision.",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &plan{
				Commands:  tt.fields.Commands,
				Decisions: tt.fields.Decisions,
				Created:   tt.fields.Created,
			}
			p.addDecision(tt.args.decision, 0, noop)
			if got := len(p.Decisions); got != 1 {
				t.Errorf("addDecision(): got  %v, want 1", got)
			}
		})
	}
}

// func Test_plan_execPlan(t *testing.T) {
// 	type fields struct {
// 		Commands  []command
// 		Decisions []string
// 		Created   time.Time
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 	}{
// 		{
// 			name: "testing executing a plan",
// 			fields: fields{
// 				Commands: []command{
// 					{
// 						Cmd:         "bash",
// 						Args:        []string{"-c", "touch hello.world"},
// 						Description: "Creating hello.world file.",
// 					}, {
// 						Cmd:         "bash",
// 						Args:        []string{"-c", "touch hello.world1"},
// 						Description: "Creating hello.world1 file.",
// 					},
// 				},
// 				Decisions: []string{"Create hello.world.", "Create hello.world1."},
// 				Created:   time.Now(),
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			p := plan{
// 				Commands:  tt.fields.Commands,
// 				Decisions: tt.fields.Decisions,
// 				Created:   tt.fields.Created,
// 			}
// 			p.execPlan()
// 			c := command{
// 				Cmd:         "bash",
// 				Args:        []string{"-c", "ls | grep hello.world | wc -l"},
// 				Description: "",
// 			}
// 			if _, got := c.exec(); strings.TrimSpace(got) != "2" {
// 				t.Errorf("execPlan(): got  %v, want hello world, again!", got)
// 			}
// 		})
// 	}
// }

// func Test_plan_printPlanCmds(t *testing.T) {
// 	type fields struct {
// 		Commands  []command
// 		Decisions []string
// 		Created   time.Time
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 	}{
// 	// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			p := plan{
// 				Commands:  tt.fields.Commands,
// 				Decisions: tt.fields.Decisions,
// 				Created:   tt.fields.Created,
// 			}
// 			p.printPlanCmds()
// 		})
// 	}
// }

// func Test_plan_printPlan(t *testing.T) {
// 	type fields struct {
// 		Commands  []command
// 		Decisions []string
// 		Created   time.Time
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 	}{
// 	// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			p := plan{
// 				Commands:  tt.fields.Commands,
// 				Decisions: tt.fields.Decisions,
// 				Created:   tt.fields.Created,
// 			}
// 			p.printPlan()
// 		})
// 	}
// }
