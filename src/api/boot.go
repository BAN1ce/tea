package api

import (
	"flag"
	"fmt"
	"net/http"
)

var (
	port = flag.Int("port", 4001, "http port")
)

func HttpServerBoot() {


	http.HandleFunc("/members", GetHandle)
	fmt.Printf("Listening on :%d\n", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		fmt.Println(err)
	}
}
