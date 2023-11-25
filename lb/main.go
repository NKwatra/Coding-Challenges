package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type LBState struct {
	lastIndex     int
	mu            sync.Mutex
	activeServers []string
}

func main() {
	defer func() {
		resp := recover()
		if val, ok := resp.(string); ok {
			fmt.Fprintf(os.Stdout, "lb: %s\n", val)
		}
	}()
	domainsPtr := flag.String("d", "", "List of domains for web servers")
	checkPtr := flag.Int("c", 10, "Interval for health check(seconds)")
	checkPathPtr := flag.String("p", "/", "Path on servers for health check")
	flag.Parse()
	if len(*domainsPtr) == 0 {
		panic("At least one backend server is required")
	}
	http.HandleFunc("/", generateHandler(*domainsPtr, *checkPtr, *checkPathPtr))
	http.ListenAndServe(":8080", nil)
}

func generateHandler(domains string, checkInterval int, checkPath string) func(w http.ResponseWriter, r *http.Request) {
	domainList := strings.Split(domains, ",")
	state := LBState{activeServers: make([]string, len(domainList))}
	copy(state.activeServers, domainList)
	client := &http.Client{}

	go setUpHealthCheck(domainList, &state, checkInterval, checkPath)
	return func(w http.ResponseWriter, r *http.Request) {
		completed := make(chan bool)
		state.mu.Lock()
		index := (state.lastIndex + 1) % len(state.activeServers)
		state.lastIndex = index
		server := state.activeServers[index]
		state.mu.Unlock()
		go processRequest(r, w, completed, client, server)
		<-completed
	}
}

func processRequest(r *http.Request, w http.ResponseWriter, completed chan bool, client *http.Client, server string) {
	headers := r.Header
	fmt.Fprintf(os.Stdout, "%v :Received request from %s\n", time.Now().Format(time.RFC1123), r.RemoteAddr)
	fmt.Fprintf(os.Stdout, "%s %s %s\n", r.Method, r.URL.Path, r.Proto)
	fmt.Fprintf(os.Stdout, "Host: %s\n", r.Host)
	fmt.Fprintf(os.Stdout, "User-Agent: %s\n", r.UserAgent())
	fmt.Fprintf(os.Stdout, "Accept: %s\n", headers.Get("Accept"))
	req, err := http.NewRequest(r.Method, server+r.URL.Path, r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Internal Server Error")
		return
	}
	for key, header := range headers {
		req.Header.Add(key, strings.Join(header, ","))
	}
	reqStartTime := time.Now().UnixMilli()
	res, err := client.Do(req)
	reqEndTime := time.Now().UnixMilli()
	if err != nil {
		fmt.Fprintln(os.Stdout)
		fmt.Fprintf(os.Stdout, "Response from server: %s %d (%d)ms\n", "Unknown", 502, reqEndTime-reqStartTime)
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprint(w, "<html><head><h1>502<br/>Bad Gateway<h1></head></html>\n")
	} else {
		fmt.Fprintln(os.Stdout)
		fmt.Fprintf(os.Stdout, "Response from server: %s %d (%d)ms\n", res.Proto, res.StatusCode, reqEndTime-reqStartTime)
		w.WriteHeader(res.StatusCode)
		w.Header().Set("Content-Type", res.Header.Get("Content-Type"))
		w.Header().Set("Content-Length", res.Header.Get("Content-Length"))
		io.Copy(w, res.Body)
		res.Body.Close()
	}
	completed <- true
}

func setUpHealthCheck(domains []string, state *LBState, interval int, path string) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	for {
		<-ticker.C
		performCheck(domains, state, path)
	}
}

func performCheck(domains []string, state *LBState, path string) {
	channels := make([]chan string, len(domains))
	for i := range channels {
		channels[i] = createRequestChannel(domains[i])
	}
	c := fanIn(channels)
	timeout := time.After(time.Duration(2) * time.Second)
	activeServers := make([]string, 0)
loop:
	for {
		select {
		case domain := <-c:
			activeServers = append(activeServers, domain)
		case <-timeout:
			break loop
		}
	}
	state.mu.Lock()
	state.activeServers = activeServers
	state.mu.Unlock()
}

func createRequestChannel(path string) chan string {
	c := make(chan string)
	go func() {
		res, err := http.Get(path)
		if err == nil && res.StatusCode == 200 {
			c <- path
		}
	}()
	return c
}

func fanIn(channels []chan string) chan string {
	c := make(chan string)
	for i := range channels {
		channel := channels[i]
		go func() {
			for {
				c <- <-channel
			}
		}()
	}
	return c
}
