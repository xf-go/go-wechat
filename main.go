package main

import (
	"net/http"

	"go-wechat/ctr"
)

func main() {
	http.HandleFunc("/getAccessToken", ctr.GetAccessToken)
	http.HandleFunc("/message", ctr.Message)
	http.ListenAndServe(":8011", nil)
}
