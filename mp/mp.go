package mp

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go-wechat/tools/curl"
)

// WechatMp 微信公众平台
type WechatMp struct {
	appID     string //wxb64d76135de4360c
	appSecret string //c91638621846716ae81e385e98b95f15
}

// NewWechatMp .
func NewWechatMp() *WechatMp {
	return &WechatMp{
		appID:     "wxb64d76135de4360c",
		appSecret: "c91638621846716ae81e385e98b95f15",
	}
}

// GetAccessToken 获取Access token
func (wm *WechatMp) GetAccessToken() (*AccessTokenResp, error) {
	reqURL := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
	reqURL = fmt.Sprintf(reqURL, wm.appID, wm.appSecret)
	res, err := curl.Get(reqURL)
	if err != nil {
		return nil, err
	}
	fmt.Println("res: ", string(res))

	resp := &AccessTokenResp{}
	err = json.Unmarshal(res, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("r: ", r)
	fmt.Println("1")
}
