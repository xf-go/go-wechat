package ctr

import (
	"fmt"
	"go-wechat/mp/core"

	"net/http"
)

// GetAccessToken .
func GetAccessToken(w http.ResponseWriter, r *http.Request) {
	srv := core.NewDefaultAccessTokenServer("wxb64d76135de4360c", "c91638621846716ae81e385e98b95f15", nil)
	res, err := srv.Token()
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	fmt.Println("res: ", res)
}
