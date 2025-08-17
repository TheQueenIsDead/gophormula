package main

import (
	"gophormula/pkg/dash"
	"net/http"
)

func main() {
	h := http.NewServeMux()
	h.HandleFunc("/", dash.Index)
	err := http.ListenAndServe(":8080", h)
	if err != nil {
		panic(err)
	}
}
