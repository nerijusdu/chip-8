package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	hub := newHub()
	messages := make(chan string)
	go hub.run(func() {
		Init()
		go GameLoop(func() {
			pixelsJson, _ := json.Marshal(data.Pixels)
			hub.Send([]byte(pixelsJson))
		}, messages)
	}, messages)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "home.html")
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	fmt.Println("App started")
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
