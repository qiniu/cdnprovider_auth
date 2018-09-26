package cdnprovider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"define"
	"errcode"
	"filelog"
	"util"
)

// 百度云鉴权模块
type BaiduAuthService struct {
	BdyAuthConf BaiduyunAuthConfiguration
}

type BaiduyunAuthConfiguration struct {
	AccessKeyId string `json:"accessKeyId"`
	SecretKey   string `json:"secretKey"`
}

var (
	baiduAuthVersion        = "bce-auth-v1"
	baiduPreffix            = "x-bce-"
	defaultExpiresInSeconds = 1800
)

func (s *BaiduAuthService) GenerateSignatureMethod(log *filelog.ReqLogger, authReq define.AuthReq) (authRes define.AuthRes, err error) {
	if s.BdyAuthConf.AccessKeyId == "" || s.BdyAuthConf.SecretKey == "" {
		log.Errorf("invalid auth conf for baiduyun, accessKeyId: %v, secretKey: %v", s.BdyAuthConf.AccessKeyId, s.BdyAuthConf.SecretKey)
		err = errcode.InvalidAuthConfErr
		return
	}
	paramsStr := make(map[string]string, len(authReq.Params))
	paramsData, _ := json.Marshal(authReq.Params)
	err = json.Unmarshal(paramsData, &paramsStr)
	if err != nil {
		log.Errorf("baidu transfer params err: ", err)
		err = errcode.InvalidParamsErr
		return
	}
	auth := s.generateSign(authReq.Method, authReq.Path, authReq.Headers, paramsStr)

	authRes.Auth = []define.AuthInfo{
		define.AuthInfo{
			Name:     "Authorization",
			Value:    auth,
			Location: define.LocationTypeHeader,
		},
	}
	return
}

func (s *BaiduAuthService) generateSign(method string, path string, headers map[string]string, params map[string]string) string {

	authString := baiduAuthVersion + "/" + s.BdyAuthConf.AccessKeyId + "/" +
		getTimeStamp() + "/" + strconv.Itoa(defaultExpiresInSeconds)
	//使用sk和authString生成signKey
	signingKey := util.HmacSha256Hex(s.BdyAuthConf.SecretKey, authString)

	//生成标准化URI
	canonicalURI := GetCanonicalURIPath(path)

	//生成标准化QueryString
	canonicalQueryString := GetCanonicalQueryString(params)

	//生成标准化header
	canonicalHeader, signedHeaders := s.getCanonicalHeaders(headers)

	//组成标准请求串
	canonicalRequest := method + "\n" + canonicalURI + "\n" + canonicalQueryString + "\n" + canonicalHeader

	//使用signKey和标准请求串完成签名
	signature := util.HmacSha256Hex(signingKey, canonicalRequest)

	authorizationHeader := authString + "/" + signedHeaders + "/" + signature
	return authorizationHeader
}

func getTimeStamp() string {
	return time.Now().UTC().Format(define.ISO8601Format)
}

/*
对HTTP请求中的Header部分进行选择性编码的结果
次序：
    将Header的名字变成全小写。
    将Header的值去掉开头和结尾的空白字符。
    经过上一步之后值为空字符串的Header忽略，其余的转换为 UriEncode(name) + ":" + UriEncode(value) 的形式。
    把上面转换后的所有字符串按照字典序进行排序。
    将排序后的字符串按顺序用\n符号连接起来得到最终的CanonicalQueryHeaders。

*/
func (s *BaiduAuthService) getCanonicalHeaders(headers map[string]string) (canonHdrs string, signedHdrs string) {
	headersToSign := map[string]string{"host": "", "content-md5": "", "content-length": "", "content-type": "", "x-bce-date": ""}

	canonHdrs = ""
	signedHdrs = ""
	//如果没有headers，则返回空串
	if headers == nil || len(headers) == 0 {
		return
	}

	headerStrings := make([]string, 0)
	signHeaders := make([]string, 0, 0)
	var item string = ""
	for k, v := range headers {
		k = UriEncode(strings.ToLower(strings.TrimSpace(k)), true)
		strValue := UriEncode(strings.TrimSpace(v), true)
		_, ok := headersToSign[k]
		if ok || strings.HasPrefix(k, baiduPreffix) {
			//如果value为nil，则赋值为空串
			if len(strValue) == 0 {
				// 值为空字符串的忽略
				continue
			}
			item = k + ":" + strValue
			headerStrings = append(headerStrings, item)
			signHeaders = append(signHeaders, k)
		}
	}
	//字典序排序
	sort.Strings(headerStrings)
	sort.Strings(signHeaders)

	canonHdrs = strings.Join(headerStrings, "\n")
	signedHdrs = strings.Join(signHeaders, ";")
	return
}

/*
RFC 3986规定，"URI非保留字符"包括以下字符：字母（A-Z，a-z）、数字（0-9）、连字号（-）、点号（.）、下划线（_)、波浪线（~），算法实现如下：
1. 将字符串转换成UTF-8编码的字节流
2. 保留所有“URI非保留字符”原样不变
3. 对其余字节做一次RFC 3986中规定的百分号编码（Percent-encoding），即一个“%”后面跟着两个表示该字节值的十六进制字母，字母一律采用大写形式。
}*/
func UriEncode(uri string, encodeSlash bool) string {
	var byteBuf bytes.Buffer
	for _, b := range []byte(uri) {
		if (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') || (b >= '0' && b <= '9') ||
			b == '-' || b == '_' || b == '.' || b == '~' || (b == '/' && !encodeSlash) {
			byteBuf.WriteByte(b)
		} else {
			byteBuf.WriteString(fmt.Sprintf("%%%02X", b))
		}
	}
	return byteBuf.String()
}

/*
对URL中的绝对路径进行编码后的结果。要求绝对路径必须以“/”开头，不以“/”开头的需要补充上，空路径为“/”
*/
func GetCanonicalURIPath(path string) string {
	if len(path) == 0 {
		return "/"
	}
	canonicalPath := path
	if strings.HasPrefix(path, "/") {
		canonicalPath = path[1:]
	}
	canonicalPath = UriEncode(canonicalPath, false)
	return "/" + canonicalPath
}

/*
 生成标准化QueryString,对于URL中的Query String（Query String即URL中“？”后面的“key1 = valve1 & key2 = valve2 ”字符串）进行编码后的结果。
 编码方法为：
 1. 将Query String根据&拆开成若干项，每一项是key=value或者只有key的形式。
 2. 对拆开后的每一项进行如下处理：
       对于key是authorization，直接忽略。
       对于只有key的项，转换为UriEncode(key) + "="的形式。
       对于key=value的项，转换为 UriEncode(key) + "=" + UriEncode(value) 的形式。这里value可以是空字符串。
 3. 将上面转换后的所有字符串按照字典顺序排序。
 4. 将排序后的字符串按顺序用 & 符号链接起来。
*/
func GetCanonicalQueryString(qmap map[string]string) string {
	if len(qmap) == 0 {
		return ""
	}
	var item string
	array := make([]string, 0)
	for k, v := range qmap {
		if strings.ToLower(k) != "authorization" {
			if len(v) > 0 {
				item = UriEncode(k, true) + "=" + UriEncode(v, true)
			} else {
				item = UriEncode(k, true) + "="
			}
			array = append(array, item)
		} else {
			/* 忽略 */
		}
	}
	sort.Strings(array)
	return strings.Join(array, "&")
}
