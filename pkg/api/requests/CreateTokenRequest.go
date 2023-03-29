package requests

type CreateTokenRequest struct {
	AppId string `json:"appId"`
	Token string `json:"token"`
}
