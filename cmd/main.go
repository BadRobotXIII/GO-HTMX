package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/badrobotxiii/Go-HTMX/internal/hardware"
	"github.com/coder/websocket"
	"github.com/pkg/browser"
)

type server struct { //Server structure
	subscribeBuffer  int
	mux              http.ServeMux
	subscribersMutex sync.Mutex
	subscribers      map[*subscriber]struct{}
}

type browserTab struct {
	Description          string `json:"description"`
	DevToolsFrontendUrl  string `json:"devtoolsFrontendUrl"`
	FaviconUrl           string `json:"faviconurl"`
	Id                   string `json:"id"`
	ThumbnailUrl         string `json:"thumbnailUrl"`
	Title                string `json:"title"`
	Url                  string `json:"url"`
	WebSocketDebuggerUrl string `json:"webSocketDebuggerUrl"`
}

type subscriber struct {
	msgs chan []byte
}

func GetTabs(tabUrl string) ([]browserTab, error) {
	u := "url"
	fmt.Printf("%s\n", u)
	resp, err := http.Get("http://localhost:8081/json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tabs []browserTab
	err = json.NewDecoder(resp.Body).Decode(&tabs)
	if err != nil {
		return nil, err
	}

	return tabs, nil
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

	tabs, e := GetTabs("http://localhost:8081")
	if e != nil {
		fmt.Printf("Error retrieving tab data %s\n", e)
	}
	fmt.Print(tabs)

	if !true {
		browser.OpenURL("http://localhost:8081")
	}

	go func(s *server) {
		for {
			//Get system information from hardware monitor
			systemMonitor, err := hardware.GetSystem()
			if err != nil {
				fmt.Println(err)
			}

			//Get disc information from hardware monitor
			diskMonitor, err := hardware.GetDisc()
			if err != nil {
				fmt.Println(err)
			}

			//Get CPU information from hardware monitor
			cpuMonitor, err := hardware.GetCPU()
			if err != nil {
				fmt.Println(err)
			}

			//Acquire current time
			timeStamp := time.Now().Format("2006-01-02 15:04:05")

			//Format message buffer
			msg := []byte(`
			<div hx-swap-oob="innerHTML:#update-timestamp"> ` + timeStamp + ` </div> 
			<div hx-swap-oob="innerHTML:#operating-system"> ` + systemMonitor.RunTimeOS + ` </div>
			<div hx-swap-oob="innerHTML:#host-name"> ` + systemMonitor.HostName + ` </div>
			<div hx-swap-oob="innerHTML:#mem-total"> ` + strconv.FormatUint(systemMonitor.VmTotal, 10) + ` </div>
			<div hx-swap-oob="innerHTML:#mem-used"> ` + strconv.FormatUint(systemMonitor.VmUsed, 10) + ` </div>
			<div hx-swap-oob="innerHTML:#disc-total"> ` + strconv.FormatUint(diskMonitor.DiscTotal, 10) + ` </div>
			<div hx-swap-oob="innerHTML:#disc-used"> ` + strconv.FormatUint(diskMonitor.DiscUsed, 10) + ` </div>
			<div hx-swap-oob="innerHTML:#disc-free"> ` + strconv.FormatUint(diskMonitor.DiskFree, 10) + ` </div>
			<div hx-swap-oob="innerHTML:#cpu-type"> ` + cpuMonitor.CpuType + ` </div>
			<div hx-swap-oob="innerHTML:#cpu-cores"> ` + strconv.FormatInt(int64(cpuMonitor.CpuCores), 10) + ` </div>
			<div hx-swap-oob="innerHTML:#cpu-speed"> ` + fmt.Sprintf("%s Mhz", strconv.FormatFloat(cpuMonitor.CpuSpeed, 'f', 1, 64)) + ` </div>`)

			//Broadcast message
			s.broadcast(msg)

			//Delay loop
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
