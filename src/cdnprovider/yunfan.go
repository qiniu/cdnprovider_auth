package cdnprovider

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"

	"define"
	"errcode"
	"filelog"
)

// 云帆鉴权模块
type YunfanAuthService struct {
	YunfanAuthConf YunfanAuthConfiguration
}

type YunfanAuthConfiguration struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}

func (s *YunfanAuthService) GenerateSignatureMethod(log *filelog.ReqLogger, authReq define.AuthReq) (authRes define.AuthRes, err error) {
	if s.YunfanAuthConf.AccessKey == "" || s.YunfanAuthConf.SecretKey == "" {
		log.Errorf("invalid auth conf for yunfan, accessKey: %v, secretKey: %v", s.YunfanAuthConf.AccessKey, s.YunfanAuthConf.SecretKey)
		err = errcode.InvalidAuthConfErr
		return
	}
	bodyData := []byte{}
	v, ok := authReq.Params["ParamString"].(string)
	if ok {
		bodyData = []byte(v)
	}

	URI := authReq.Path + "\n" + string(bodyData)
	sign := hmac.New(sha1.New, []byte(s.YunfanAuthConf.SecretKey))
	if _, err = sign.Write([]byte(URI)); err != nil {
		log.Error("hmac-sha1 auth info err:", err)
		err = errcode.InvalidParamsErr
		return
	}
	encodedSign := base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%x", sign.Sum(nil))))
	authorization := s.YunfanAuthConf.AccessKey + ":" + encodedSign

	authRes.Auth = []define.AuthInfo{
		define.AuthInfo{
			Name:     "Authorization",
			Value:    authorization,
			Location: define.LocationTypeHeader,
		},
	}
	return
}
