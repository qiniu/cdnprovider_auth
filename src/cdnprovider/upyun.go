package cdnprovider

import (
	"define"
	"errcode"
	"filelog"
)

// 又拍鉴权模块
type UpyunAuthService struct {
	UpyunAuthConf UpyunAuthConfiguration
}

type UpyunAuthConfiguration struct {
	Token string `json:"token"`
}

func (s *UpyunAuthService) GenerateSignatureMethod(log *filelog.ReqLogger, authReq define.AuthReq) (authRes define.AuthRes, err error) {
	if s.UpyunAuthConf.Token == "" {
		log.Errorf("invalid auth conf for upyun, token: %v", s.UpyunAuthConf.Token)
		err = errcode.InvalidAuthConfErr
		return
	}
	authRes.Auth = []define.AuthInfo{
		define.AuthInfo{
			Name:     "Authorization",
			Value:    "Bearer " + s.UpyunAuthConf.Token,
			Location: define.LocationTypeHeader,
		},
	}
	return
}
