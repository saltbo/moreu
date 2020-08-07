package oauth2

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"gopkg.in/resty.v1"
)

type Config struct {
	AuthURL      string `yaml:"auth_url"`
	TokenURL     string `yaml:"token_url"`
	UserURL      string `yaml:"user_url"`
	ClientId     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	Scope        string `yaml:"scope"`
}

type JSONResponse struct {
	Code   int               `json:"code"`
	ErrMsg string            `json:"err_msg"`
	Data   map[string]string `json:"data"`
}

func NewJSONResponse() *JSONResponse {
	return &JSONResponse{}
}

type Oauth2 struct {
	cli  *resty.Client
	conf Config
}

func New(conf Config) (*Oauth2, error) {
	return &Oauth2{
		cli:  resty.New(),
		conf: conf,
	}, nil
}

func (c *Oauth2) Authorize(redirectUri string) string {
	v := make(url.Values)
	v.Set("client_id", c.conf.ClientId)
	v.Set("redirect_uri", redirectUri)
	v.Set("response_type", "code")
	v.Set("scope", c.conf.Scope)
	v.Set("state", "2333")
	return fmt.Sprintf("%s?%s", c.conf.AuthURL, v.Encode())
}

func (c *Oauth2) accessToken(code, redirectUri string) (string, error) {
	body := map[string]string{
		"client_id":     c.conf.ClientId,
		"client_secret": c.conf.ClientSecret,
		"code":          code,
		"grant_type":    "authorization_code",
		"redirect_uri":  redirectUri,
	}
	errResp := NewJSONResponse()
	req := c.cli.R().SetBody(body).SetHeader("Content-Type", "application/json").SetResult(&errResp)
	resp, err := req.Post(c.conf.TokenURL)
	if err != nil {
		return "", err
	} else if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("%s", resp.Status())
	}

	if errResp.Code != 0 {
		return "", fmt.Errorf(errResp.ErrMsg)
	}

	succResp := make(map[string]interface{})
	if err := json.Unmarshal(resp.Body(), &succResp); err != nil {
		return "", err
	}

	return succResp["access_token"].(string), nil
}

func (sso *Oauth2) GetUserInfo(code, redirectUri string) (map[string]string, error) {
	token, err := sso.accessToken(code, redirectUri)
	if err != nil {
		return nil, err
	}

	respBody := make(map[string]string)
	req := sso.cli.R().SetAuthToken(token).SetResult(&respBody)
	resp, err := req.Get(sso.conf.UserURL)
	if err != nil {
		return nil, err
	} else if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("%s", resp.Status())
	}

	return respBody, nil
}
