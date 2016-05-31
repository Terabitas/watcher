package server

import (
	"net/http"

	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/nildev/watcher/config"
	"github.com/nildev/watcher/reporter"
	"github.com/nildev/watcher/version"
)

// Server type
type Server struct {
	stop    chan bool
	cfg     config.Config
	handler http.Handler
}

// New type
func New(cfg config.Config) (*Server, error) {
	srv := Server{
		cfg:  cfg,
		stop: nil,
	}
	return &srv, nil
}

// Run server
func (s *Server) Run() {
	ctxLog := log.WithField("version", version.Version).WithField("git-hash", version.GitHash).WithField("build-time", version.BuiltTimestamp)

	ctxLog.Infof("Starting watcher service push to [%s] from [%s] every %d s", s.cfg.PushEndpoint, s.cfg.MetricsEndpoint, s.cfg.ReportInterval)
	s.stop = make(chan bool)

	go func() {
		var r reporter.Reporter
	created:
		for {
			r = reporter.NewRemoteReporter(s.cfg.PushEndpoint, reporter.NewLocalMetadataFetcher())
			if r == nil {
				ctxLog.Errorf("Could not create reporter")
				time.Sleep(time.Second * 3)
				continue
			}
			break created
		}

		for {
			ctxLog.Infof(".")
			err := r.Report([]byte(`{}`))
			if err != nil {
				ctxLog.Errorf("%s", err)
			}
			time.Sleep(time.Second * time.Duration(s.cfg.ReportInterval))
		}
	}()
}

// Stop server
func (s *Server) Stop() {
	close(s.stop)
}

// Purge server
func (s *Server) Purge() {
}
