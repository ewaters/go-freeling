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
	flAddr = flag.String("fl_addr", "", "")
	addr   = flag.String("addr", ":8080", "")
	debug  = flag.Bool("debug", false, "")
)

func logExit(msg string, args ...interface{}) {
	fmt.Printf(msg+"\n", args...)
	os.Exit(1)
}

func main() {
	flag.Parse()

	if *flAddr == "" {
		logExit("--fl_addr=... is required")
	}

	client, err := freeling.New(*flAddr)
	if err != nil {
		log.Fatal(err)
	}
	client.Debug = *debug
	defer client.Close()

	http.HandleFunc("/freeling-es-json", func(w http.ResponseWriter, r *http.Request) {
		failed := func(format string, args ...interface{}) {
			fmt.Fprintf(w, "Failed: "+format, args...)
		}

		msg := r.FormValue("text")
		log.Printf("Request from %s for %q", r.RemoteAddr, msg)
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
	})

	log.Printf("Listening for HTTP connections on %s...", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
