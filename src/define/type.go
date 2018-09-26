package define

import (
	"encoding/json"
	"time"

	"filelog"
)

type AuthReq struct {
	CdnProvider CDNProvider            `json:"cdnProvider"`
	Host        string                 `json:"host"`
	Path        string                 `json:"path"`
	Method      string                 `json:"method"`
	Params      map[string]interface{} `json:"params"`
	Headers     map[string]string      `json:"headers"`
	Operation   Operation              `json:"operation"`
}

type Operation string

const (
	OperationRefresh    Operation = "refresh"
	OperationPrefresh   Operation = "prefresh"
	OperationDomainConf Operation = "domainConf"
)

type IpWhiteListReq struct {
	IpWhiteList []string `json:"ipWhiteList"`
}

type AuthRes struct {
	Auth []AuthInfo `json:"auth"`
}

type AuthInfo struct {
	Name     string       `json:"name"`
	Value    string       `json:"value"`
	Location LocationType `json:"location"`
}

type LocationType string

const (
	LocationTypeBody     LocationType = "body"
	LocationTypeUrlQuery LocationType = "urlquery"
	LocationTypeUrlPath  LocationType = "urlpath"
	LocationTypeHeader   LocationType = "header"
)

type CDNProvider string

const (
	CDNProviderAliyun     CDNProvider = "aliyun"
	CDNProviderBaidu      CDNProvider = "baidu"
	CDNProviderBaishanyun CDNProvider = "baishanyun"
	CDNProviderDilian     CDNProvider = "dilian"
	CDNProviderKuaiwang   CDNProvider = "kuaiwang"
	CDNProviderTencent    CDNProvider = "tencent"
	CDNProviderUPYun      CDNProvider = "upyun"
	CDNProviderWangsu     CDNProvider = "wangsu"
	CDNProviderYunfan     CDNProvider = "yunfan"
)

var ValidCDNProviderMap = map[CDNProvider]bool{
	CDNProviderAliyun:     true,
	CDNProviderBaidu:      true,
	CDNProviderBaishanyun: true,
	CDNProviderDilian:     true,
	CDNProviderKuaiwang:   true,
	CDNProviderTencent:    true,
	CDNProviderUPYun:      true,
	CDNProviderWangsu:     true,
	CDNProviderYunfan:     true,
}

func (c CDNProvider) String() string {
	return string(c)
}

type ErrorInfo struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}

func (e *ErrorInfo) Error() string {
	errData, _ := json.Marshal(e)
	return string(errData)
}

func (e *ErrorInfo) ErrorCode() int {
	return e.Code
}

func (e *ErrorInfo) HttpCode() int {
	return e.Code / 1000
}

func NewError(code int, err string) *ErrorInfo {
	return &ErrorInfo{code, err}
}

type GenerateSignatureMethod interface {
	GenerateSignatureMethod(*filelog.ReqLogger, AuthReq) (AuthRes, error)
}

const (
	ISO8601Format = "2006-01-02T15:04:05Z"
)

var (
	Local = time.FixedZone("CST", 8*3600)
)
