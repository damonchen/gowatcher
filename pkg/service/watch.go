package service

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/damonchen/gowatcher/pkg/config"
	"github.com/fsnotify/fsnotify"
)

type ProcessInfo struct {
	pid    int
	cancel context.CancelFunc
}

type Watcher struct {
	cfg          *config.Config
	watcher      *fsnotify.Watcher
	processInfos map[string]ProcessInfo
}

func NewWatcher(cfg *config.Config) *Watcher {
	return &Watcher{
		cfg:          cfg,
		processInfos: make(map[string]ProcessInfo),
	}
}

func (w *Watcher) isIgnoreFile(filename string) bool {
	excludePaths := w.cfg.ExcludedPaths
	for _, excludePath := range excludePaths {
		if strings.HasPrefix(filename, excludePath) {
			return true
		}
	}
	return false
}

func (w *Watcher) isWatchFile(filename string) bool {
	includePaths := w.cfg.IncludePaths
	for _, includePath := range includePaths {
		if strings.HasPrefix(filename, includePath) {
			return true
		}
	}
	return false
}

func (w *Watcher) watch() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	go w.watchEvent(watcher)

	for _, includePath := range w.cfg.IncludePaths {
		log.Infof("start watch %+v\n", includePath)
		// if path, should list all files then add
		files, err := filepath.Glob(includePath)
		if err != nil {
			return err
		}

		for _, file := range files {
			watcher.Add(file)
		}
	}
	w.watcher = watcher
	return nil
}

func (w *Watcher) close() {
	_ = w.watcher.Close()
}

func (w *Watcher) watchEvent(watcher *fsnotify.Watcher) {
	for {
		select {
		case e, ok := <-watcher.Events:
			if !ok {
				continue
			}
			log.Infof("watch file event %+v", e)
			if w.isIgnoreFile(e.Name) {
				log.Debugf("file %s changed, but in ignore file list", e.Name)
				continue
			}

			//if !w.isWatchFile(e.Name) {
			//	log.Debugf("file %s changed, but in not in process file list", e.Name)
			//	continue
			//}

			log.Infof("watch file %s changed", e.Name)

			err := w.runCmd(e.Name)
			if err != nil {
				log.Errorf("run command error %s, env: %+v\n", err, os.Environ())
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				continue
			}
			log.Errorf("watch error occur %s", err)
		}
	}
}

func prepareCmd(cmd string, filename string) string {
	basename := filepath.Base(filename)
	ext := filepath.Ext(filename)
	filenameWithoutExt := filename[:len(filename)-len(ext)]
	basenameWithoutExt := filepath.Base(filenameWithoutExt)
	fileDir := filepath.Dir(filename)

	cmd = strings.ReplaceAll(cmd, "${basename}", basename)
	cmd = strings.ReplaceAll(cmd, "${basenameWithoutExt}", basenameWithoutExt)
	cmd = strings.ReplaceAll(cmd, "${filenameWithoutExt}", filenameWithoutExt)
	cmd = strings.ReplaceAll(cmd, "${filename}", filename)
	cmd = strings.ReplaceAll(cmd, "${fileDir}", fileDir)
	return cmd
}

func (w *Watcher) runCmd(name string) error {
	for _, cmd := range w.cfg.Command {
		info, exists := w.processInfos[cmd.Cmd]
		var pid int
		if exists {
			pid = info.pid
		}
		execPath := prepareCmd(cmd.Cmd, name)

		log.Infof("will restart cmd %s", execPath)

		pid, cancel, err := process(execPath, cmd.Args, cmd.Envs, pid)
		if err != nil {
			return err
		}

		w.processInfos[cmd.Cmd] = ProcessInfo{
			pid:    pid,
			cancel: cancel,
		}
	}
	return nil
}
