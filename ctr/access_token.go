package ctr

import (
	"fmt"
	"go-wechat/mp"

	"net/http"
)

// GetAccessToken .
func GetAccessToken(w http.ResponseWriter, r *http.Request) {
	wechatMp := mp.NewWechatMp()
	res, err := wechatMp.GetAccessToken()
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	fmt.Println("res: ", res)
}
