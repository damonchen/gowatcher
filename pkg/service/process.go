package service

import (
	"context"
	"os"
	"os/exec"
)

func process(cmd string, args []string, envs []string, pid int) (int, context.CancelFunc, error) {
	if pid != 0 {
		log.Debugf("process pid %d", pid)
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
		log.Errorf("find process error %s", err)
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
	cmd.Args = args
	cmd.Env = append(os.Environ(), envs...)

	err := cmd.Start()
	if err != nil {
		log.Errorf("start command %s %+v error", execPath, args)
		return 0, cancel, err
	}
	go func() {
		_ = cmd.Wait()
	}()

	return cmd.Process.Pid, cancel, nil
}
