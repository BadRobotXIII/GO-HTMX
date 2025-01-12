package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/badrobotxiii/Go-HTMX/internal/hardware"
	"github.com/coder/websocket"
)

type server struct {
	subscribeBuffer  int
	mux              http.ServeMux
	subscribersMutex sync.Mutex
	subscribers      map[*subscriber]struct{}
}

http.ServeMux

type subscriber struct {
	msgs chan []byte
}

func NewServer() *server {
	s := &server{
		subscribeBuffer: 10,
		subscribers:     make(map[*subscriber]struct{}),
	}

	s.mux.Handle("/", http.FileServer(http.Dir("/htmx")))
	s.mux.HandleFunc("/ws", s.subscribeHandler)
	return s
}

func (s *server) subscribeHandler(writer http.ResponseWriter, req *http.Request) {
	err := s.subscribe(req.Context(), writer, req)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (s *server) addSubscriber(subscriber *subscriber) {
	s.subscribersMutex.Lock()
	s.subscribers[subscriber] = struct{}{}
	s.subscribersMutex.Unlock()
	fmt.Println("Added subscriber")
}

func (s *server) subscribe(ctx context.Context, writer http.ResponseWriter, req *http.Request) error {
	var c *websocket.Conn
	subscriber := &subscriber{
		msgs: make(chan []byte, s.subscribeBuffer),
	}
	s.addSubscriber(subscriber)

	c, err := websocket.Accept(writer, req, nil)
	if err != nil {
		return err
	}
	defer c.CloseNow()

	ctx = c.CloseRead(ctx)

	for {
		select {
		case msg := <-subscriber.msgs:
			ctx, cancel := context.WithTimeout(ctx, time.Second)
			defer cancel()
			err := c.Write(ctx, websocket.MessageText, msg)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}

}

func main() {
	fmt.Println("Starting system monitor...")
	go func() {
		for {
			systemMonitor, err := hardware.GetSystem()
			if err != nil {
				fmt.Println(err)
			}

			diskMonitor, err := hardware.GetDisk()
			if err != nil {
				fmt.Println(err)
			}

			cpuMonitor, err := hardware.GetCPU()
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(systemMonitor)
			fmt.Println(diskMonitor)
			fmt.Println(cpuMonitor)

			time.Sleep(3 * time.Second)
		}
	}()

	fmt.Println("Starting server...")
	srvr := NewServer()
	err := http.ListenAndServe(":3000", &srvr.mux)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
