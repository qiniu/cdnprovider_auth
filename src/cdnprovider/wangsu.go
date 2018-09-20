package cdnprovider

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"io"

	"define"
	"errcode"
	"filelog"
)

// 网宿鉴权模块
type WanggsuService struct {
	WangsuAuthConf WangsuAuthConfiguration
}

type WangsuAuthConfiguration struct {
	ApiKey   string `json:"apiKey"`
	UserName string `json:"userName"`
	Password string `json:"password"`
}

func (s *WanggsuService) GenerateSignatureMethod(log *filelog.ReqLogger, authReq define.AuthReq) (authRes define.AuthRes, err error) {
	if s.WangsuAuthConf.ApiKey == "" || s.WangsuAuthConf.UserName == "" || s.WangsuAuthConf.Password == "" {
		log.Errorf("invalid auth conf for wangsu, apiKey: %v, userName: %v, password: %v",
			s.WangsuAuthConf.ApiKey, s.WangsuAuthConf.UserName, s.WangsuAuthConf.Password)
		err = errcode.InvalidAuthConfErr
		return
	}

	switch authReq.Operation {
	case define.OperationPrefresh:
		return s.getPrefreshPwd(log, authReq)
	case define.OperationRefresh:
		return s.getRefreshPwd(log, authReq)
	case define.OperationDomainConf:
		return s.getDomainConfPwd(log, authReq)
	}

	return
}

func (s *WanggsuService) getRefreshPwd(log *filelog.ReqLogger, authReq define.AuthReq) (authRes define.AuthRes, err error) {
	param, ok := authReq.Params["ParamString"].(string)
	if !ok {
		log.Error("ParamString is essential")
		err = errcode.InvalidParamsErr
		return
	}
	m := md5.New()
	io.WriteString(m, s.WangsuAuthConf.UserName+s.WangsuAuthConf.Password+param)
	pwd := hex.EncodeToString(m.Sum(nil))

	authRes.Auth = []define.AuthInfo{
		define.AuthInfo{
			Name:     "username",
			Value:    s.WangsuAuthConf.UserName,
			Location: define.LocationTypeBody,
		},
		define.AuthInfo{
			Name:     "password",
			Value:    pwd,
			Location: define.LocationTypeBody,
		},
	}
	return
}

func (s *WanggsuService) getPrefreshPwd(log *filelog.ReqLogger, authReq define.AuthReq) (authRes define.AuthRes, err error) {

	param, ok := authReq.Params["ParamString"].(string)
	if !ok {
		log.Error("ParamString is essential")
		err = errcode.InvalidParamsErr
		return
	}

	m := md5.New()
	io.WriteString(m, param+s.WangsuAuthConf.UserName+"chinanetcenter"+s.WangsuAuthConf.Password)
	pwd := hex.EncodeToString(m.Sum(nil))

	authRes.Auth = []define.AuthInfo{
		define.AuthInfo{
			Name:     "username",
			Value:    s.WangsuAuthConf.UserName,
			Location: define.LocationTypeBody,
		},
		define.AuthInfo{
			Name:     "password",
			Value:    pwd,
			Location: define.LocationTypeBody,
		},
	}
	return
}

func (s *WanggsuService) getDomainConfPwd(log *filelog.ReqLogger, authReq define.AuthReq) (authRes define.AuthRes, err error) {
	date, ok := authReq.Params["Date"].(string)
	if !ok {
		log.Error("Date is essential")
		err = errcode.InvalidParamsErr
		return
	}

	h := hmac.New(sha1.New, []byte(s.WangsuAuthConf.ApiKey))
	h.Write([]byte(date))
	pwd := base64.StdEncoding.EncodeToString(h.Sum(nil))
	auth := base64.StdEncoding.EncodeToString([]byte(s.WangsuAuthConf.UserName + ":" + pwd))

	authRes.Auth = []define.AuthInfo{
		define.AuthInfo{
			Name:     "Authorization",
			Value:    "Basic " + auth,
			Location: define.LocationTypeHeader,
		},
	}
	return
}
