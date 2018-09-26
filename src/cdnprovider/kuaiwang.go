package cdnprovider

import (
	"encoding/json"
	"net/http"
	"time"

	"cache"
	"define"
	"errcode"
	"filelog"
	"util"
)

// 快网鉴权模块
type KuaiwangAuthService struct {
	KuaiwangAuthConf KaiuwangAuthConfiguration
	CacheClient      cache.CacheClient
}

type KaiuwangAuthConfiguration struct {
	AppId     string `json:"appId"`
	AppSecret string `json:"appSecret"`
}

var (
	kwHost    = "https://cdncs-api.fastweb.com.cn"
	tokenURI  = "/oauth/access_token"
	grantType = "client_credentials"
)

type kwAuthParam struct {
	GrantType string `json:"grant_type"`
	AppId     string `json:"appid"`
	AppSecret string `json:"appsecret"`
}

type kwAuthResult struct {
	Status int64  `json:"status"`
	Info   string `json:"info"`
	Result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	} `json:"result"`
}

func (s *KuaiwangAuthService) GenerateSignatureMethod(log *filelog.ReqLogger, authReq define.AuthReq) (authRes define.AuthRes, err error) {
	if s.KuaiwangAuthConf.AppId == "" || s.KuaiwangAuthConf.AppSecret == "" {
		log.Errorf("invalid auth conf for kuaiwang, appId: %v, appSecret: %v", s.KuaiwangAuthConf.AppId, s.KuaiwangAuthConf.AppSecret)
		err = errcode.InvalidAuthConfErr
		return
	}
	// 从缓存中获取鉴权信息
	if s.CacheClient != nil {
		if authRes, err = s.CacheClient.Get(define.CDNProviderKuaiwang); err == nil {
			log.Info("kuaiwang get auth info from cache")
			return
		}
		log.Warn("kuaiwang get auth info from cache err: ", err)
	}

	log.Info("kuaiwang get auth info from cdn provider service")
	// 没有缓存或缓存过期，请求快网鉴权服务，获取鉴权信息
	token, expires, err := s.getAccessToken(log)
	if err != nil {
		log.Error("kuaiwang request new auth info err:", err)
		return
	}
	authRes.Auth = []define.AuthInfo{
		define.AuthInfo{
			Name:     "access_token",
			Value:    token,
			Location: define.LocationTypeBody,
		},
	}
	// 更新缓存
	if s.CacheClient != nil {
		if err_ := s.CacheClient.Upsert(define.CDNProviderKuaiwang, authRes, time.Second*time.Duration(expires)); err_ != nil {
			log.Error("update kuaiwang auth info cache failed, err:", err_)
		}
	}
	return
}

func (s *KuaiwangAuthService) getAccessToken(log *filelog.ReqLogger) (token string, expires int, err error) {
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	url := kwHost + tokenURI
	authReqParam := kwAuthParam{
		GrantType: grantType,
		AppId:     s.KuaiwangAuthConf.AppId,
		AppSecret: s.KuaiwangAuthConf.AppSecret,
	}

	authData, _ := json.Marshal(authReqParam)
	res, err := util.DoHTTPRequest(log, url, http.MethodPost, string(authData), headers)
	if err != nil {
		log.Error("kuaiwang request new auth info failed, err:", err)
		err = define.NewError(errcode.ErrBadRequest, err.Error())
		return
	}

	var authResult kwAuthResult
	if err = util.ParseReqArgs(res.Body, res.ContentLength, &authResult); err != nil {
		log.Error("kuaiwang parse request body failed, err:", err)
		err = define.NewError(errcode.ErrBadRequest, err.Error())
		return
	}

	token = authResult.Result.AccessToken
	if token == "" {
		log.Error("kuaiwang token is null")
		err = errcode.NullAuthInfoErr
		return
	}
	// 默认时间为 86400s，提前 1 分钟过期，以防在过期临界点的请求获取过期的鉴权信息
	expires = authResult.Result.ExpiresIn - 1*60
	return
}
