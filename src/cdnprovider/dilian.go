package cdnprovider

import (
	"time"

	"define"
	"errcode"
	"filelog"
	"util"
)

// 帝联鉴权模块
type DilianAuthService struct {
	DilianAuthConf DilianAuthConfiguration
}

type DilianAuthConfiguration struct {
	AccessKeyId string `json:"accessKeyId"`
	AccessKey   string `json:"accessKey"`
}

func (s *DilianAuthService) GenerateSignatureMethod(log *filelog.ReqLogger, authReq define.AuthReq) (authRes define.AuthRes, err error) {
	if s.DilianAuthConf.AccessKeyId == "" || s.DilianAuthConf.AccessKey == "" {
		log.Errorf("invalid auth conf for dilian, accessKeyId: %v, accessKey: %v", s.DilianAuthConf.AccessKeyId, s.DilianAuthConf.AccessKey)
		err = errcode.InvalidAuthConfErr
		return
	}
	bodyData := []byte{}
	v, ok := authReq.Params["ParamString"].(string)
	if ok {
		bodyData = []byte(v)
	}

	md5BodyData := util.DoMd5(string(bodyData))
	credential := s.DilianAuthConf.AccessKeyId + "/" + time.Now().In(define.Local).Format("20060102150405") + "/dnioncloud"
	var signature string
	signature = util.DoMd5(authReq.Method + "\n" + s.DilianAuthConf.AccessKey + "\n" + authReq.Path + "\n" + md5BodyData + "\n" + credential)
	authorization := "Algorithm=md5," + "Credential=" + credential + ",Signature=" + signature

	authRes.Auth = []define.AuthInfo{
		define.AuthInfo{
			Name:     "Authorization",
			Value:    authorization,
			Location: define.LocationTypeHeader,
		},
	}
	return
}
