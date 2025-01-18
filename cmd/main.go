package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
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

type subscriber struct {
	msgs chan []byte
}

func NewServer() *server {
	s := &server{
		subscribeBuffer: 10,
		subscribers:     make(map[*subscriber]struct{}),
	}
	s.mux.Handle("/", http.FileServer(http.Dir("C:\\Users\\kzulf\\Dropbox\\Coding\\GoLang\\Go-htmx\\htmx")))
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
			ctx, cancel := context.WithTimeout(ctx, time.Second*10)
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

func (s *server) broadcast(msg []byte) {
	s.subscribersMutex.Lock()
	defer s.subscribersMutex.Unlock()
	for subscriber := range s.subscribers {
		subscriber.msgs <- msg
	}
}

func main() {
	fmt.Println("Starting system monitor...")
	srvr := NewServer()
	go func(s *server) {
		for {
			systemMonitor, err := hardware.GetSystem()
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(systemMonitor)
			}

			sysMon := systemMonitor

			fmt.Printf("Sys Mon: %v\n", sysMon)
			fmt.Printf("Runtime OS: %v\n", sysMon.RunTimeOS)

			diskMonitor, err := hardware.GetDisk()
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(diskMonitor)
			}

			cpuMonitor, err := hardware.GetCPU()
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(cpuMonitor)
			}

			timeStamp := time.Now().Format("2006-01-02 15:04:05")

			msg := []byte(`
			<div hx-swap-oob="innerHTML:#update-timestamp"> ` + timeStamp + ` </div>
			<div hx-swap-oob="innerHTML:#operating-system"> ` + systemMonitor.RunTimeOS + ` </div>
			<div hx-swap-oob="innerHTML:#disk-total"> ` + strconv.FormatUint(diskMonitor.DiscTotal, 10) + ` </div>
			<div hx-swap-oob="innerHTML:#cpu-type"> ` + cpuMonitor.CpuType + ` </div>`)

			s.broadcast(msg)

			time.Sleep(1 * time.Second)
		}
	}(srvr)

	fmt.Println("Starting server...")

	err := http.ListenAndServe(":8081", &srvr.mux)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
