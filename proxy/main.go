package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	freeling "../client"
)

var (
	addr   = flag.String("addr", ":8080", "")
	debug  = flag.Bool("debug", false, "")
	langs = flag.String("langs", "de=:10001;es=:10002", "")

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
	config := map[string]string{}
	for _, pair := range strings.Split(*langs, ";") {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			logExit("--config %q was invalid; must take form '<lang>=<addr>[;...]'", *langs) 
		}
		config[parts[0]] = parts[1]
	}


	clients = make(map[string]*freeling.Client)
	for lang, addr := range config {
		client, err := freeling.New(addr)
		if err != nil {
			logExit("Failed to create Freeling client for %s (%q): %v", lang, addr, err)
		}
		client.Debug = *debug
		defer client.Close()
		clients[lang] = client
	}

	for lang := range clients {
		http.HandleFunc(fmt.Sprintf("/freeling-%s-json", lang), handlerForLang(lang))
	}
	
	log.Printf("Listening for HTTP connections on %s...", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
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
