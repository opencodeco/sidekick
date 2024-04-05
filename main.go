package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/opencodeco/sidekick/proxy"
)

var (
	port      string
	appPort   string
	logFormat string
	logLevel  string
)

func main() {
	flag.StringVar(&port, "port", "9601", "sidekick port")
	flag.StringVar(&appPort, "app-port", "8888", "application port")
	flag.StringVar(&logFormat, "log-format", "text", "log format")
	flag.StringVar(&logLevel, "log-level", "info", "log level")
	flag.Parse()

	level := &slog.LevelVar{}
	level.UnmarshalText([]byte(logLevel))
	logOpts := slog.HandlerOptions{
		Level: level,
	}

	var logHandler slog.Handler
	logHandler = slog.NewTextHandler(os.Stdout, &logOpts)
	if logFormat == "json" {
		logHandler = slog.NewJSONHandler(os.Stdout, &logOpts)
	}
	log := slog.New(logHandler)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	appCmd := exec.Command(flag.Args()[0], flag.Args()[1:]...)
	http.HandleFunc("/", proxy.Proxy(appPort))

	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for sig := range c {
			log.Info(fmt.Sprintf("captured %v, terminating the application", sig))

			err := appCmd.Process.Signal(syscall.SIGTERM)
			if err != nil {
				log.Error("error terminating application", err)
			}

			log.Debug("exiting sidekick, bye bye")
			os.Exit(0)
		}
	}()

	err := appCmd.Start()
	log.Info(fmt.Sprintf("application started at %s", appPort))
	if err != nil {
		log.Error("error running application", err)
	}

	log.Info(fmt.Sprintf("starting sidekick at %s", port))
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("error running sidekick", err)
	}
}
