package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	freeling "github.com/ewaters/go-freeling/client"
)

var (
	addr       = flag.String("addr", ":8080", "")
	debug      = flag.Bool("debug", false, "")
	langs      = flag.String("langs", "de=:10001;es=:10002", "")
	retryDelay = flag.Duration("retry_delay", 5*time.Second, "")

	config  map[string]string
	clients map[string]*freeling.Client
)

func logExit(msg string, args ...interface{}) {
	fmt.Printf(msg+"\n", args...)
	os.Exit(1)
}

func main() {
	flag.Parse()

	if *langs == "" {
		logExit("--langs=... is required")
	}
	config = map[string]string{}
	for _, pair := range strings.Split(*langs, ";") {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			logExit("--config %q was invalid; must take form '<lang>=<addr>[;...]'", *langs)
		}
		config[parts[0]] = parts[1]
	}

	clients = make(map[string]*freeling.Client)
	connectAll()

	var endpoints []string
	for lang := range clients {
		url := fmt.Sprintf("/freeling-%s-json", lang)
		http.HandleFunc(url, handlerForLang(lang))
		endpoints = append(endpoints, url)
	}
	http.HandleFunc("/", rootHandler(endpoints))

	log.Printf("Listening for HTTP connections on %s...", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func connectAll() {
	for lang := range config {
		if _, ok := clients[lang]; ok {
			continue
		}
		connectOneRetry(lang)
	}
}

func connectOneRetry(lang string) {
	for {
		err := connectOne(lang)
		if err == nil {
			return
		}
		log.Printf("Failed connectOne(%q): %v; will retry in %v", lang, err, *retryDelay)
		time.Sleep(*retryDelay)
	}
}

func connectOne(lang string) error {
	addr, ok := config[lang]
	if !ok {
		return fmt.Errorf("Can't connectOne(%q): no addr", lang)
	}
	if client, ok := clients[lang]; ok {
		client.Close()
	}
	log.Printf("Dialing %q on %q", lang, addr)
	client, err := freeling.New(addr)
	if err != nil {
		return err
	}
	client.Debug = *debug
	clients[lang] = client
	return nil
}

func rootHandler(endpoints []string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		for _, endpoint := range endpoints {
			fmt.Fprintf(w, "Endpoint <a href=%q>%s</a></br>", endpoint, endpoint)
		}
	}
}

func handlerForLang(lang string) func(http.ResponseWriter, *http.Request) {
	client := clients[lang]
	if client == nil {
		logExit("No client found for language %s", lang)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		failed := func(format string, args ...interface{}) {
			fmt.Fprintf(w, "Failed: "+format, args...)
		}

		msg := r.FormValue("text")
		log.Printf("[%s] Request from %s for %q", lang, r.RemoteAddr, msg)
		if msg == "" {
			failed("'text' required")
			return
		}

		strs, err := client.Send(msg)
		if err != nil && strings.Contains(err.Error(), "broken pipe") {
			log.Printf("Broken pipe %q; reconnecting", lang)
			connectOneRetry(lang)
			strs, err = client.Send(msg)
		}
		if err != nil {
			failed("freeling Send(%q) failed: %v", msg, err)
			return
		}
		res := strings.Join(strs, ", ")

		w.Header().Add("Content-Type", "application/json")
		if *debug {
			fmt.Printf("\n\n[%s]\n\n", res)
		}
		fmt.Fprintf(w, "[%s]", res)
	}
}
