package app

import (
	"bytes"
	"fmt"
	"math"
	"os/exec"
	"strings"
	"time"
)

// Command type representing all executable commands Helmsman needs
// to execute in order to inspect the environment|releases|charts etc.
type Command struct {
	Cmd         string
	Args        []string
	Description string
}

// CmdPipe is a os/exec.Commnad wrapper for UNIX pipe
type CmdPipe []Command

type ExitStatus struct {
	code   int
	errors string
	output string
}

func (e ExitStatus) String() string {
	str := strings.TrimSpace(e.output)
	if errs := strings.TrimSpace(e.errors); errs != "" {
		str = fmt.Sprintf("%s\n--- stderr ---\n%s", str, errs)
	}
	return str
}

func (c *Command) String() string {
	var sb strings.Builder
	sb.WriteString(c.Cmd)
	for i := 0; i < len(c.Args); i++ {
		arg := c.Args[i]
		sb.WriteRune(' ')
		if strings.HasPrefix(arg, "--token=") {
			sb.WriteString("--token=******")
			continue
		}
		if strings.HasPrefix(arg, "--password=") {
			sb.WriteString("--password=******")
			continue
		}
		if arg == "--token" {
			sb.WriteString(arg)
			sb.WriteString("=******")
			i++
			continue
		}
		if arg == "--password" {
			sb.WriteString(arg)
			sb.WriteString("=******")
			i++
			continue
		}
		sb.WriteString(arg)
	}
	return sb.String()
}

// RetryExec runs exec command with retry
func (c *Command) RetryExec(attempts int) (ExitStatus, error) {
	return c.RetryExecWithThreshold(attempts, 0)
}

// RetryExecWithThreshold runs exec command with retry and allows specifying the threshold for the exit code to be considered erroneous
func (c *Command) RetryExecWithThreshold(attempts, exitCodeThreshold int) (ExitStatus, error) {
	var (
		result ExitStatus
		err    error
	)

	for i := 0; i < attempts; i++ {
		result, err = c.Exec()
		if err == nil || (result.code >= 0 && result.code <= exitCodeThreshold) {
			return result, nil
		}
		if i < (attempts - 1) {
			time.Sleep(time.Duration(math.Pow(2, float64(2+i))) * time.Second)
			log.Infof("Retrying %s due to error: %v", c.Description, err)
		}
	}

	return result, fmt.Errorf("%s, failed after %d attempts with: %w", c.Description, attempts, err)
}

func (c *Command) command() *exec.Cmd {
	// Only use non-empty string args
	var args []string

	for _, str := range c.Args {
		if str != "" {
			args = append(args, str)
		}
	}

	log.Verbose(c.Description)
	log.Debug(c.String())

	return exec.Command(c.Cmd, args...)
}

// Exec executes the executable command and returns the exit code and execution result
func (c *Command) Exec() (ExitStatus, error) {
	var stdout, stderr bytes.Buffer
	cmd := c.command()
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		log.Info("cmd.Start: " + err.Error())
		return ExitStatus{
			code:   1,
			errors: err.Error(),
		}, err
	}

	err := cmd.Wait()
	res := ExitStatus{
		output: strings.TrimSpace(stdout.String()),
		errors: strings.TrimSpace(stderr.String()),
	}
	if err != nil {
		res.code = 126
		if exiterr, ok := err.(*exec.ExitError); ok {
			res.code = exiterr.ExitCode()
		}
		err = newExitError(c.Description, res.code, res.output, res.errors, err)
	}
	return res, err
}

// Exec pipes the executable commands and returns the exit code and execution result
func (p CmdPipe) Exec() (ExitStatus, error) {
	var (
		stdout, stderr bytes.Buffer
		stack          []*exec.Cmd
	)

	l := len(p) - 1
	if l < 0 {
		// nonthing to do here
		return ExitStatus{}, nil
	}
	if l == 0 {
		// it's just one command we can just run it
		return p[0].Exec()
	}

	for i, c := range p {
		stack = append(stack, c.command())
		stack[i].Stderr = &stderr
		if i > 0 {
			stack[i].Stdin, _ = stack[i-1].StdoutPipe()
		}
	}
	stack[l].Stdout = &stdout

	err := call(stack)
	res := ExitStatus{
		output: strings.TrimSpace(stdout.String()),
		errors: strings.TrimSpace(stderr.String()),
	}
	if err != nil {
		res.code = 126
		if exiterr, ok := err.(*exec.ExitError); ok {
			res.code = exiterr.ExitCode()
		}
		err = newExitError(p[l].Description, res.code, res.output, res.errors, err)
	}
	return res, err
}

// RetryExec runs piped commands with retry
func (p CmdPipe) RetryExec(attempts int) (ExitStatus, error) {
	return p.RetryExecWithThreshold(attempts, 0)
}

// RetryExecWithThreshold runs piped commands with retry and allows specifying the threshold for the exit code to be considered erroneous
func (p CmdPipe) RetryExecWithThreshold(attempts, exitCodeThreshold int) (ExitStatus, error) {
	var (
		result ExitStatus
		err    error
	)

	l := len(p) - 1
	for i := 0; i < attempts; i++ {
		result, err = p.Exec()
		if err == nil || (result.code >= 0 && result.code <= exitCodeThreshold) {
			return result, nil
		}
		if i < (attempts - 1) {
			time.Sleep(time.Duration(math.Pow(2, float64(2+i))) * time.Second)
			log.Infof("Retrying %s due to error: %v", p[l].Description, err)
		}
	}

	return result, fmt.Errorf("%s, failed after %d attempts with: %w", p[l].Description, attempts, err)
}

func call(stack []*exec.Cmd) (err error) {
	if stack[0].Process == nil {
		if err = stack[0].Start(); err != nil {
			return err
		}
	}
	if len(stack) > 1 {
		if err = stack[1].Start(); err != nil {
			return err
		}
		defer func() {
			if err == nil {
				err = call(stack[1:])
			} else {
				err = stack[1].Wait()
			}
			if err != nil {
				log.Infof("call: %v", err)
			}
		}()
	}
	return stack[0].Wait()
}

func newExitError(cmd string, code int, stdout, stderr string, cause error) error {
	return fmt.Errorf(
		"%s failed with non-zero exit code %d: %w\noutput: %s",
		cmd, code, cause,
		fmt.Sprintf(
			"\n--- stdout ---\n%s\n--- stderr ---\n%s",
			strings.TrimSpace(stdout),
			strings.TrimSpace(stderr),
		),
	)
}

// ToolExists returns true if the tool is present in the environment and false otherwise.
// It takes as input the tool's command to check if it is recognizable or not. e.g. helm or kubectl
func ToolExists(tool string) bool {
	_, err := exec.LookPath(tool)
	return err == nil
}
