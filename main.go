package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/opencodeco/sidekick/proxy"
)

var (
	port    string
	appPort string
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	flag.StringVar(&port, "port", "9601", "sidekick port")
	flag.StringVar(&appPort, "app-port", "8888", "application port")
	flag.Parse()

	appCmd := exec.Command(flag.Args()[0], flag.Args()[1:]...)
	http.HandleFunc("/", proxy.Proxy(appPort))

	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for sig := range c {
			fmt.Printf("captured %v, terminating the application\n", sig)
			err := appCmd.Process.Signal(syscall.SIGTERM)
			if err != nil {
				fmt.Println("error terminating application", err)
			}
			fmt.Println("exiting sidekick, bye bye")
			os.Exit(0)
		}
	}()

	err := appCmd.Start()
	if err != nil {
		fmt.Println("error running application", err)
	}

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("error running sidekick", err)
	}
}
