package core

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"sync/atomic"
	"time"
	"unsafe"

	"go-wechat/tools/curl"
)

// access_token 中控服务器接口.
type AccessTokenServer interface {
	Token() (token string, err error)                           // 请求中控服务器返回缓存的 access_token
	RefreshToken(currentToken string) (token string, err error) // 请求中控服务器刷新 access_token                     // 接口标识, 没有实际意义
}

var _ AccessTokenServer = (*DefaultAccessTokenServer)(nil)

// DefaultAccessTokenServer 微信公众平台
type DefaultAccessTokenServer struct {
	appID      string //wxb64d76135de4360c
	appSecret  string //c91638621846716ae81e385e98b95f15
	httpClient *http.Client

	refreshTokenRequestChan  chan string             // chan currentToken
	refreshTokenResponseChan chan refreshTokenResult // chan {token, err}

	tokenCache unsafe.Pointer // *accessToken
}

// NewDefaultAccessTokenServer 创建一个新的 DefaultAccessTokenServer, 如果 httpClient == nil 则默认使用 util.DefaultHttpClient.
func NewDefaultAccessTokenServer(appID, appSecret string, httpClient *http.Client) (srv *DefaultAccessTokenServer) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	srv = &DefaultAccessTokenServer{
		appID:                    url.QueryEscape(appID),
		appSecret:                url.QueryEscape(appSecret),
		httpClient:               httpClient,
		refreshTokenRequestChan:  make(chan string),
		refreshTokenResponseChan: make(chan refreshTokenResult),
	}

	go srv.tokenUpdateDaemon(time.Hour * 24 * time.Duration(100+rand.Int63n(200)))
	return
}

func (srv *DefaultAccessTokenServer) tokenUpdateDaemon(initTickDuration time.Duration) {
	tickDuration := initTickDuration

NEW_TICK_DURATION:
	ticker := time.NewTicker(tickDuration)
	for {
		select {
		case currentToken := <-srv.refreshTokenRequestChan:
			accessToken, cached, err := srv.updateToken(currentToken)
			if err != nil {
				srv.refreshTokenResponseChan <- refreshTokenResult{err: err}
				break
			}
			srv.refreshTokenResponseChan <- refreshTokenResult{token: accessToken.Token}
			if !cached {
				tickDuration = time.Duration(accessToken.ExpiresIn) * time.Second
				ticker.Stop()
				goto NEW_TICK_DURATION
			}

		case <-ticker.C:
			accessToken, _, err := srv.updateToken("")
			if err != nil {
				break
			}
			newTickDuration := time.Duration(accessToken.ExpiresIn) * time.Second
			if abs(tickDuration-newTickDuration) > time.Second*5 {
				tickDuration = newTickDuration
				ticker.Stop()
				goto NEW_TICK_DURATION
			}
		}
	}
}

func abs(x time.Duration) time.Duration {
	if x >= 0 {
		return x
	}
	return -x
}

type refreshTokenResult struct {
	token string
	err   error
}

type AccessToken struct {
	AccessToken string
	ExpiresIn   time.Duration
	ExpiresAt   time.Time
}

// NewWechatMp .
// func NewWechatMp() *WechatMp {
// 	return &WechatMp{
// 		appID:     "wxb64d76135de4360c",
// 		appSecret: "c91638621846716ae81e385e98b95f15",
// 	}
// }

func (srv *DefaultAccessTokenServer) Token() (token string, err error) {
	if p := (*accessToken)(atomic.LoadPointer(&srv.tokenCache)); p != nil {
		return p.Token, nil
	}
	return srv.RefreshToken("")
}

func (srv *DefaultAccessTokenServer) RefreshToken(currentToken string) (token string, err error) {
	srv.refreshTokenRequestChan <- currentToken
	rslt := <-srv.refreshTokenResponseChan
	return rslt.token, rslt.err
}

type accessToken struct {
	Token     string `json:"access_token"`
	ExpiresIn int64  `json:"expires_in"`
}

// 从微信服务器获取新的 access_token 并存入缓存, 同时返回该 access_token.
func (srv *DefaultAccessTokenServer) updateToken(currentToken string) (token *accessToken, cached bool, err error) {
	reqURL := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
	reqURL = fmt.Sprintf(reqURL, srv.appID, srv.appSecret)
	res, err := curl.Get(reqURL)
	if err != nil {
		atomic.StorePointer(&srv.tokenCache, nil)
		return
	}
	fmt.Println("res: ", string(res))

	var result struct {
		Error
		accessToken
	}
	err = json.Unmarshal(res, &result)
	if err != nil {
		atomic.StorePointer(&srv.tokenCache, nil)
		return
	}
	if result.ErrCode != ErrCodeOK {
		atomic.StorePointer(&srv.tokenCache, nil)
		err = &result.Error
		return
	}

	tokenCopy := result.accessToken
	atomic.StorePointer(&srv.tokenCache, unsafe.Pointer(&tokenCopy))
	token = &tokenCopy
	return
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("r: ", r)
	fmt.Println("1")
}

func DecodeJSONHttpResponse(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}
