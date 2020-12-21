package mp

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("r: ", r)
	fmt.Println("1")
}

func Message(w http.ResponseWriter, r *http.Request) {
	file, err := os.OpenFile("1.log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	file.Write(bs)
	fmt.Println("message")
}
