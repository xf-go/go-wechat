package mp

import "time"

// AccessTokenResp .
type AccessTokenResp struct {
	AccessToken string        `json:"access_token"`
	ExpiresIn   time.Duration `json:"expires_in"`
}
