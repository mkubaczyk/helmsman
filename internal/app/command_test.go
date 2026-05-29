package app

import (
	"fmt"
	"testing"
)

func TestToolExists(t *testing.T) {
	type args struct {
		tool string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test case 1 -- checking helm exists.",
			args: args{
				tool: helmBin,
			},
			want: true,
		}, {
			name: "test case 2 -- checking kubectl exists.",
			args: args{
				tool: kubectlBin,
			},
			want: true,
		}, {
			name: "test case 3 -- checking nonExistingTool exists.",
			args: args{
				tool: "nonExistingTool",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToolExists(tt.args.tool); got != tt.want {
				t.Errorf("toolExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommandExec(t *testing.T) {
	type input struct {
		cmd  string
		args []string
		desc string
	}
	type expected struct {
		code   int
		err    error
		output string
	}
	tests := []struct {
		name  string
		input input
		want  expected
	}{
		{
			name: "echo",
			input: input{
				cmd:  "bash",
				args: []string{"-c", "echo this is fun"},
				desc: "A bash command execution test with echo.",
			},
			want: expected{
				code:   0,
				output: "this is fun",
				err:    nil,
			},
		}, {
			name: "exitCode",
			input: input{
				cmd:  "bash",
				args: []string{"-c", "echo $?"},
				desc: "A bash command execution test with exitCode.",
			},
			want: expected{
				code:   0,
				output: "0",
				err:    nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Command{
				Cmd:         tt.input.cmd,
				Args:        tt.input.args,
				Description: tt.input.desc,
			}
			got, err := c.Exec()
			if err != nil && tt.want.err == nil {
				t.Errorf("command.exec() unexpected error got:\n%v want:\n%v", err, tt.want.err)
			}
			if err != nil && tt.want.err != nil {
				if err.Error() != tt.want.err.Error() {
					t.Errorf("command.exec() unexpected error got:\n%v want:\n%v", err, tt.want.err)
				}
			}
			if got.code != tt.want.code {
				t.Errorf("command.exec() unexpected code got = %v, want = %v", got.code, tt.want.code)
			}
			if got.output != tt.want.output {
				t.Errorf("command.exec() unexpected output got:\n%v want:\n%v", got.output, tt.want.output)
			}
		})
	}
}

func TestPipeExec(t *testing.T) {
	type expected struct {
		code   int
		err    error
		output string
	}
	tests := []struct {
		name  string
		input CmdPipe
		want  expected
	}{
		{
			name: "echo",
			input: CmdPipe{
				Command{
					Cmd:         "printf",
					Args:        []string{"first string\nsecond string\nthird string"},
					Description: "multiline printf",
				},
			},
			want: expected{
				code:   0,
				output: "first string\nsecond string\nthird string",
				err:    nil,
			},
		}, {
			name: "line count",
			input: CmdPipe{
				Command{
					Cmd:         "printf",
					Args:        []string{"first string\nsecond string\nthird string\n"},
					Description: "multiline printf",
				},
				Command{
					Cmd:         "wc",
					Args:        []string{"-l"},
					Description: "line count",
				},
			},
			want: expected{
				code:   0,
				output: "3",
				err:    nil,
			},
		}, {
			name: "grep",
			input: CmdPipe{
				Command{
					Cmd:         "printf",
					Args:        []string{"first string\nsecond string\nthird string\n"},
					Description: "multiline printf",
				},
				Command{
					Cmd:         "grep",
					Args:        []string{"second"},
					Description: "grep",
				},
			},
			want: expected{
				code:   0,
				output: "second string",
				err:    nil,
			},
		}, {
			name: "grep no matches",
			input: CmdPipe{
				Command{
					Cmd:         "printf",
					Args:        []string{"first string\nsecond string\nthird string\n"},
					Description: "multiline printf",
				},
				Command{
					Cmd:         "grep",
					Args:        []string{"fourth"},
					Description: "grep",
				},
			},
			want: expected{
				code:   1,
				output: "",
				err:    newExitError("grep", 1, "", "", fmt.Errorf("exit status 1")),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.input.Exec()
			if err != nil && tt.want.err == nil {
				t.Errorf("command.exec() unexpected error got:\n%v want:\n%v", err, tt.want.err)
			}
			if err != nil && tt.want.err != nil {
				if err.Error() != tt.want.err.Error() {
					t.Errorf("command.exec() unexpected error got:\n%v want:\n%v", err, tt.want.err)
				}
			}
			if got.code != tt.want.code {
				t.Errorf("command.exec() unexpected code got = %v, want = %v", got.code, tt.want.code)
			}
			if got.output != tt.want.output {
				t.Errorf("command.exec() unexpected output got:\n%v want:\n%v", got.output, tt.want.output)
			}
		})
	}
}

func TestCommand_String(t *testing.T) {
	tests := []struct {
		name     string
		cmd      Command
		expected string
	}{
		{
			"regular",
			kubectl([]string{"config", "set-cluster", "CONTEXT", "--server=http://localhost:8080", "--certificate-authority=cacert.crt"}, ""),
			"kubectl config set-cluster CONTEXT --server=http://localhost:8080 --certificate-authority=cacert.crt",
		},
		{
			"cert-key",
			kubectl([]string{"config", "set-credentials", "USER", "--client-key=client.key", "--client-certificate=client.crt"}, ""),
			"kubectl config set-credentials USER --client-key=client.key --client-certificate=client.crt",
		},
		{
			"password",
			kubectl([]string{"config", "set-credentials", "USER", "--username=foo", "--password=secret"}, ""),
			"kubectl config set-credentials USER --username=foo --password=******",
		},
		{
			"password2",
			kubectl([]string{"config", "set-credentials", "USER", "--username", "foo", "--password", "secret"}, ""),
			"kubectl config set-credentials USER --username foo --password=******",
		},
		{
			"token",
			kubectl([]string{"config", "set-credentials", "USER", "--token=secret"}, ""),
			"kubectl config set-credentials USER --token=******",
		},
		{
			"token2",
			kubectl([]string{"config", "set-credentials", "USER", "--token", "secret"}, ""),
			"kubectl config set-credentials USER --token=******",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := test.cmd.String()
			if actual != test.expected {
				t.Errorf("command.String() unexpected got = %s, want = %s\n", actual, test.expected)
			}
		})
	}
}
