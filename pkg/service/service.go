package service

import (
	"sync"

	logger "github.com/phachon/go-logger"

	"github.com/damonchen/gowatcher/pkg/config"
)

type Service struct {
	*Watcher
}

var (
	log  *logger.Logger
	once sync.Once
)

func NewService(cfg *config.Config) *Service {
	return &Service{
		NewWatcher(cfg),
	}
}

func (svc *Service) Run() error {
	once.Do(func() {
		log = logger.NewLogger()
	})

	defer svc.Close()

	done := make(chan bool)
	err := svc.watch()
	if err != nil {
		log.Errorf("watch error %s", err)
		return err
	}

	<-done
	return nil
}

func (svc *Service) Close() {
	svc.close()
}