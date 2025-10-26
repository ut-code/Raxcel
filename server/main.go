package main

import (
	"net/http"

	handler "github.com/ut-code/Raxcel/server/api"
)

func main() {
	http.HandleFunc("/", handler.Handler)
	http.ListenAndServe(":8080", nil)
}
