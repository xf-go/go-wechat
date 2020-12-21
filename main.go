package main

import (
	"net/http"

	"go-wechat/mp"
)

func main() {
	http.HandleFunc("/message", mp.Message)
	http.ListenAndServe(":8011", nil)
}
