package main

import (
	"net/http"

	"github.com/ut-code/Raxcel/server/api"
)

func main() {
	http.HandleFunc("/", api.Handler)
	http.ListenAndServe(":8080", nil)
}
