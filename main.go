package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)
var (
	table map[string]string
	tableLock sync.RWMutex
)
func add(request, destination string){
	tableLock.Lock()
	defer tableLock.Unlock()
	if _, ok := table[request]; ok {
		return
	}
	table[request] = destination
	log.Printf("%s pointing to %s added!\n", request, destination)
}
func remove(request string){
	tableLock.Lock()
	defer tableLock.Unlock()
	if _, ok := table[request]; ok {
		delete(table, request)
		log.Printf("%s deleted\n", request)
	}
}
func handler(w http.ResponseWriter, r *http.Request){
	domain := r.Host
	pos := strings.Index(r.Host, ":")
	if pos != -1 {
		domain = r.Host[0:pos]
	}

	if _, ok := table[domain]; !ok {
		return
	}

	url, err := url.Parse(table[domain])
	if err != nil{
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(url)

	r.URL.Host = url.Host
	r.URL.Scheme = url.Scheme
	r.Host = url.Host

	log.Printf("Serving content from %s\n", url)

	proxy.ServeHTTP(w, r)
}
func main() {
	var(
		tls bool
		certificate, key string
	)
	table = make(map[string]string)
	flag.BoolVar(&tls, "tls", false, "TLS enabled")
	flag.StringVar(&certificate, "certificate", "", "Path to TLS certificate")
	flag.StringVar(&key, "key", "", "Path to TLS key")
	flag.Parse()

	if tls{
		if _, err := ioutil.ReadFile(certificate); err != nil {
			panic(err)
		}
		if _, err := ioutil.ReadFile(key); err != nil {
			panic(err)
		}
	}

	http.HandleFunc("/", handler)

	if tls{
		log.Println("Listening on 443 (TLS)")
		go func() {
			if err := http.ListenAndServeTLS(":443", certificate, key,nil); err != nil {
				log.Fatal(err)
			}
		}()
	}

	log.Println("Listening on 80")
	go func() {
		if err := http.ListenAndServe(":80",nil); err != nil {
			log.Fatal(err)
		}
	}()

	go signalHandler()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Pretender> ")
		userInput, _ := reader.ReadString('\n')
		userInput = strings.Replace(userInput, "\n", "", -1)
		command := strings.Fields(userInput)
		if len(command) < 1 {
			continue
		}
		switch command[0] {
		case "add":
			if len(command) < 3{
				continue
			}
			request, destination := command[1], command[2]
			add(request, destination)
		case "remove":
			if len(command) < 2{
				continue
			}
			request := command[1]
			remove(request)
		default:
			fmt.Println("Invalid command.\nUse 'add domain.com https://www.google.com' to add a new domain.\nUse 'remove new.domain' to remove a domain.")
		}
	}
}
func signalHandler() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	s := <-signals
	log.Printf("Received signal: %s\n", s)
	os.Exit(1)
}