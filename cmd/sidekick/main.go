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

	"sidekick/internal/components"
	"sidekick/internal/utils"
)

var (
	port      string
	appPort   string
	logFormat string
	logLevel  string
)

func main() {
	parseFlags()
	utils.SetupLogger(logLevel, logFormat)
	setupComponentHandlers()
	app := createApp()
	gracefulShutdown(app)
	startAndListen(app, appPort, port)	
}

func parseFlags() {
	flag.StringVar(&port, "port", "9601", "sidekick port")
	flag.StringVar(&appPort, "app-port", "8888", "application port")
	flag.StringVar(&logFormat, "log-format", "text", "log format")
	flag.StringVar(&logLevel, "log-level", "info", "log level")
	flag.Parse()
}

func setupComponentHandlers() {
	http.HandleFunc("/health", components.Health)
	http.HandleFunc("/", components.Proxy(appPort))
}

func createApp() *exec.Cmd {
	return exec.Command(flag.Args()[0], flag.Args()[1:]...)
}

func gracefulShutdown(app *exec.Cmd) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for sig := range c {
			slog.Info(fmt.Sprintf("captured %v, terminating the application", sig))

			err := app.Process.Signal(syscall.SIGTERM)
			if err != nil {
				slog.Error("error terminating application", err)
			}

			slog.Debug("exiting sidekick, bye bye")
			os.Exit(0)
		}
	}()
}

func startAndListen(app *exec.Cmd, appPort string, port string) {
	err := app.Start()
	slog.Info(fmt.Sprintf("application started at %s", appPort))
	if err != nil {
		slog.Error("error running application", "err", err)
		os.Exit(1)
	}

	slog.Info(fmt.Sprintf("starting sidekick at %s", port))
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		slog.Error("error running sidekick", "err", err)
		os.Exit(1)
	}
}
