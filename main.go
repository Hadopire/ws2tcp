package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Hadopire/ws2tcp/proxy"
	"github.com/gorilla/websocket"
)

const (
	defaultPort = "8000"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func Connection(w http.ResponseWriter, r *http.Request) {
	source, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	target, err := proxy.Init(source)
	if err != nil {
		source.Close()
		return
	}

	go proxy.TargetToSource(source, target)
	go proxy.SourceToTarget(source, target)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	http.HandleFunc("/", Connection)

	var port string
	flag.StringVar(&port, "p", defaultPort, "specifies the port")
	flag.Parse()

	fmt.Println("listening on port " + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
