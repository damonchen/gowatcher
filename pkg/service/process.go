package service

import (
	"context"
	"os"
	"os/exec"
)

func process(cmd string, args []string, envs []string, pid int) (int, context.CancelFunc, error) {
	if pid != 0 {
		err := kill(pid)
		if err != nil {
			return 0, nil, err
		}
	}

	return startProcess(cmd, args, envs)
}

func kill(pid int) error {
	ps, err := os.FindProcess(pid)

	if err != nil {
		return err
	}

	err = ps.Kill()
	if err != nil {
		return err
	}

	_, err = ps.Wait()
	return err
}

func startProcess(execPath string, args []string, envs []string) (int, context.CancelFunc, error) {
	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, execPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Args = append([]string{execPath}, args...)
	cmd.Env = append(os.Environ(), envs...)

	err := cmd.Start()
	if err != nil {
		return 0, cancel, err
	}
	go func() {
		_ = cmd.Wait()
	}()

	return cmd.Process.Pid, cancel, nil
}
