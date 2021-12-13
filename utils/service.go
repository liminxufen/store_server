package utils

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/store_server/logger"
	log "github.com/store_server/logger"
)

//instance class
type Instance interface {
	Start() error
	Stop() error
}

var (
	defaultListenSignal = []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP}
)

func contains(list []os.Signal, item os.Signal) (has bool) {
	for _, ele := range list {
		if ele == item {
			return true
		}
	}
	return
}

//service define
type Service struct {
	sf       map[os.Signal]func(os.Signal)
	instance Instance
}

func NewService(instance Instance) (s *Service) {
	return &Service{instance: instance}
}

func (ser *Service) Register(signal2 os.Signal, f func(signal2 os.Signal)) {
	if ser.sf == nil {
		ser.sf = make(map[os.Signal]func(os.Signal))
	}
	ser.sf[signal2] = f
}

func (ser *Service) Forever() {
	sigChan := make(chan os.Signal, 1)

	go func() {
		if err := ser.instance.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	sigs := defaultListenSignal[0:]
	for sig := range ser.sf {
		if !contains(sigs, sig) {
			sigs = append(sigs, sig)
		}
	}

	signal.Notify(sigChan, sigs...)
	for {
		s := <-sigChan
		logger.Entry().Infof("recv signal:%v", s)
		f, ok := ser.sf[s]
		if !ok {
			logger.Entry().Infof("unregistered signal:%v", s)
		} else {
			f(s)
			continue
		}
		switch s {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			if err := ser.instance.Stop(); err != nil {
				logger.Entry().Errorf("stop instance err:%v", err)
			}
			os.Exit(0)
		case syscall.SIGHUP:
		}
	}

}
