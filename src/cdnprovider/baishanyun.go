package cdnprovider

import (
	"define"
	"errcode"
	"filelog"
)

// 白山云鉴权模块
type BaishanyunAuthService struct {
	BaishanAuthConf BaishanAuthConfiguration
}

type BaishanAuthConfiguration struct {
	Token string `json:"token"`
}

func (s *BaishanyunAuthService) GenerateSignatureMethod(log *filelog.ReqLogger, authReq define.AuthReq) (authRes define.AuthRes, err error) {
	if s.BaishanAuthConf.Token == "" {
		log.Errorf("invalid auth conf for baishanyun token: %v", s.BaishanAuthConf.Token)
		err = errcode.InvalidAuthConfErr
		return
	}
	authRes.Auth = []define.AuthInfo{
		define.AuthInfo{
			Name:     "token",
			Value:    s.BaishanAuthConf.Token,
			Location: define.LocationTypeUrlQuery,
		},
	}
	return
}
