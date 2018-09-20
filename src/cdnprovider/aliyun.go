package cdnprovider

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"net/url"
	"sort"
	"strings"

	"define"
	"errcode"
	"filelog"
)

// 阿里云鉴权模块
type AliAuthService struct {
	AliAuthConf AliAuthConfiguration
}

type AliAuthConfiguration struct {
	AccessKeyId string `json:"accessKeyId"`
	SecretKey   string `json:"secretKey"`
}

func (s *AliAuthService) GenerateSignatureMethod(log *filelog.ReqLogger, authReq define.AuthReq) (authRes define.AuthRes, err error) {
	if s.AliAuthConf.AccessKeyId == "" || s.AliAuthConf.SecretKey == "" {
		log.Errorf("invalid auth for aliyun conf accessKeyId: %v, secretKey: %v", s.AliAuthConf.AccessKeyId, s.AliAuthConf.SecretKey)
		err = errcode.InvalidAuthConfErr
		return
	}
	params := map[string]string{}
	keys := make([]string, 0, len(params))

	authReq.Params["AccessKeyId"] = s.AliAuthConf.AccessKeyId
	for k, v := range authReq.Params {
		ss, ok := v.(string)
		if ok {
			params[k] = ss
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	urlParams := ""
	for _, key := range keys {
		urlParams += "&" + key + "=" + strings.Replace(url.QueryEscape(params[key]), "+", "%20", -1)
	}
	urlParams = urlParams[1:]

	encodedUrlParams := url.QueryEscape(urlParams)
	StringToSign := authReq.Method + "&" + url.QueryEscape("/") + "&" + encodedUrlParams
	hmacObj := hmac.New(sha1.New, []byte(s.AliAuthConf.SecretKey+"&"))
	hmacObj.Write([]byte(StringToSign))
	signature := url.QueryEscape(base64.StdEncoding.EncodeToString(hmacObj.Sum(nil)))

	authRes.Auth = []define.AuthInfo{
		define.AuthInfo{
			Name:     "AccessKeyId",
			Value:    s.AliAuthConf.AccessKeyId,
			Location: define.LocationTypeUrlQuery,
		},
		define.AuthInfo{
			Name:     "Signature",
			Value:    signature,
			Location: define.LocationTypeUrlQuery,
		},
	}
	return
}
