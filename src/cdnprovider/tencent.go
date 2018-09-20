package cdnprovider

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"sort"
	"strings"

	"define"
	"errcode"
	"filelog"
)

// 腾讯云鉴权模块
type TencentAuthService struct {
	TencentAuthConf TencentAuthConfiguration
}

type TencentAuthConfiguration struct {
	SecretId  string `json:"secretId"`
	SecretKey string `json:"secretKey"`
}

func (s *TencentAuthService) GenerateSignatureMethod(log *filelog.ReqLogger, authReq define.AuthReq) (authRes define.AuthRes, err error) {
	if s.TencentAuthConf.SecretId == "" || s.TencentAuthConf.SecretKey == "" {
		log.Errorf("invalid auth conf for tencent, secretId: %v, secretKey: %v", s.TencentAuthConf.SecretId, s.TencentAuthConf.SecretKey)
		err = errcode.InvalidAuthConfErr
		return
	}

	authReq.Params["SecretId"] = s.TencentAuthConf.SecretId

	keys := make([]string, 0, len(authReq.Params))
	for k, _ := range authReq.Params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var plainParms string
	for i := range keys {
		k := keys[i]
		plainParms += "&" + fmt.Sprintf("%v", k) + "=" + fmt.Sprintf("%v", authReq.Params[k])
	}
	plainText := strings.ToUpper(authReq.Method) + authReq.Host + authReq.Path + "?" + plainParms[1:]

	hmacObj := hmac.New(sha1.New, []byte(s.TencentAuthConf.SecretKey))
	hmacObj.Write([]byte(plainText))
	sign := base64.StdEncoding.EncodeToString(hmacObj.Sum(nil))

	authRes.Auth = []define.AuthInfo{
		define.AuthInfo{
			Name:     "SecretId",
			Value:    s.TencentAuthConf.SecretId,
			Location: define.LocationTypeUrlQuery,
		},
		define.AuthInfo{
			Name:     "Signature",
			Value:    sign,
			Location: define.LocationTypeUrlQuery,
		},
	}
	return
}
