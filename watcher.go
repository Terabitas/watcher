package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/nildev/watcher/config"
	"github.com/nildev/watcher/server"
	"github.com/nildev/watcher/version"
	"github.com/rakyll/globalconf"
)

var (
	GitHash        = ""
	BuiltTimestamp = ""
	Version        = ""
	ctxLog         *log.Entry
)

const (
	DefaultConfigFile = "/etc/watcher/watcher.conf"
)

func init() {
	version.Version = Version
	version.GitHash = GitHash
	version.BuiltTimestamp = BuiltTimestamp

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
}

func main() {
	ctxLog = log.WithField("version", version.Version).WithField("git-hash", version.GitHash).WithField("build-time", version.BuiltTimestamp)
	userset := flag.NewFlagSet("watcherd", flag.ExitOnError)

	printVersion := userset.Bool("version", false, "Print the version and exit")
	cfgPath := userset.String("config", DefaultConfigFile, fmt.Sprintf("Path to config file. watcherd will look for a config at %s by default.", DefaultConfigFile))

	err := userset.Parse(os.Args[1:])
	if err == flag.ErrHelp {
		userset.Usage()
		os.Exit(1)
	}

	args := userset.Args()
	if len(args) == 1 && args[0] == "version" {
		*printVersion = true
	} else if len(args) != 0 {
		userset.Usage()
		os.Exit(1)
	}

	if *printVersion {
		fmt.Printf("Version: %s \n", version.Version)
		fmt.Printf("Git hash: %s \n", version.GitHash)
		fmt.Printf("Build timestamp: %s \n", version.BuiltTimestamp)
		os.Exit(0)
	}

	cfgset := flag.NewFlagSet("watcherd", flag.ExitOnError)
	// Generic
	cfgset.Int("verbosity", 0, "Logging level")
	cfgset.Int("report_interval", 1, "Interval which defines how often to push metrics")
	cfgset.String("push_endpoint", "http://nil.services/register", "Endpoint where to push metrics")
	cfgset.String("metrics_endpoint", "http://127.0.0.1:19999", "Endpoint where to query for metrics")
	cfgset.String("metrics", "", "Metrics")

	globalconf.Register("", cfgset)
	cfg, err := getConfig(cfgset, *cfgPath)
	if err != nil {
		ctxLog.Fatalf(err.Error())
	}

	srv, err := server.New(*cfg)
	if err != nil {
		ctxLog.Fatalf("Failed creating Server: %v", err.Error())
	}
	srv.Run()

	reconfigure := func() {
		ctxLog.Infof("Reloading configuration from %s", *cfgPath)

		cfg, err := getConfig(cfgset, *cfgPath)
		if err != nil {
			ctxLog.Fatalf(err.Error())
		}

		ctxLog.Infof("Restarting server components")
		srv.Stop()

		srv, err = server.New(*cfg)
		if err != nil {
			ctxLog.Fatalf(err.Error())
		}
		srv.Run()
	}

	shutdown := func() {
		ctxLog.Infof("Gracefully shutting down")
		srv.Stop()
		srv.Purge()
		os.Exit(0)
	}

	writeState := func() {
		ctxLog.Infof("Dumping server state")

		encoded, err := json.Marshal(srv)
		if err != nil {
			ctxLog.Errorf("Failed to dump server state: %v", err)
			return
		}

		if _, err := os.Stdout.Write(encoded); err != nil {
			ctxLog.Errorf("Failed to dump server state: %v", err)
			return
		}

		os.Stdout.Write([]byte("\n"))
	}

	signals := map[os.Signal]func(){
		syscall.SIGHUP:  reconfigure,
		syscall.SIGTERM: shutdown,
		syscall.SIGINT:  shutdown,
		syscall.SIGUSR1: writeState,
		syscall.SIGABRT: shutdown,
	}

	listenForSignals(signals)
}

func getConfig(flagset *flag.FlagSet, userCfgFile string) (*config.Config, error) {
	opts := globalconf.Options{EnvPrefix: "WATCHERD_"}

	if userCfgFile != "" {
		// Fail hard if a user-provided config is not usable
		fi, err := os.Stat(userCfgFile)
		if err != nil {
			ctxLog.Fatalf("Unable to use config file %s: %v", userCfgFile, err)
		}
		if fi.IsDir() {
			ctxLog.Fatalf("Provided config %s is a directory, not a file", userCfgFile)
		}
		opts.Filename = userCfgFile
	} else if _, err := os.Stat(DefaultConfigFile); err == nil {
		opts.Filename = DefaultConfigFile
	}

	gconf, err := globalconf.NewWithOptions(&opts)
	if err != nil {
		return nil, err
	}

	gconf.ParseSet("", flagset)

	cfg := config.Config{
		Verbosity:       (*flagset.Lookup("verbosity")).Value.(flag.Getter).Get().(int),
		ReportInterval:  (*flagset.Lookup("report_interval")).Value.(flag.Getter).Get().(int),
		PushEndpoint:    (*flagset.Lookup("push_endpoint")).Value.(flag.Getter).Get().(string),
		MetricsEndpoint: (*flagset.Lookup("metrics_endpoint")).Value.(flag.Getter).Get().(string),
		Metrics:         (*flagset.Lookup("metrics")).Value.(flag.Getter).Get().(string),
	}

	log.SetLevel(log.Level(cfg.Verbosity))

	ctxLog.Infof("Loaded config: [%+v]", cfg)

	return &cfg, nil
}

func listenForSignals(sigmap map[os.Signal]func()) {
	sigchan := make(chan os.Signal, 1)

	for k := range sigmap {
		signal.Notify(sigchan, k)
	}

	for true {
		sig := <-sigchan
		handler, ok := sigmap[sig]
		if ok {
			handler()
		}
	}
}
