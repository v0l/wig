package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	_version string = "v0.1"
)

var _settings = LoadSettings()

type Server struct {
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Upgrade") == "websocket" {
		fmt.Printf("New client: %s\n", r.RemoteAddr)
		NewClient(w, r)
	} else {
		body := "Hello World\n"
		w.Header().Set("Server", "go-web-irc")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", fmt.Sprint(len(body)))
		fmt.Fprint(w, body)
	}
}

func main() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	signal.Notify(sigchan, syscall.SIGTERM)

	server := Server{}
	go func() {
		http.HandleFunc("/ws", server.ServeHTTP)
		hter := http.ListenAndServe(_settings.WsAddress, nil)
		if hter != nil {
			fmt.Errorf(hter.Error())
			os.Exit(11)
		}
	}()

	fmt.Println("WS port started on:", _settings.WsAddress)
	<-sigchan
}
