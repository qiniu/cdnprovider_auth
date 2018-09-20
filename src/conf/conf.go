package conf

import (
	"encoding/json"
	"io/ioutil"

	"cache"
	"cdnprovider"
)

var ServerConf *CdnProviderAuthConf

type CdnProviderAuthConf struct {
	BindHost    string          `json:"bindHost"`    // 绑定的服务地址
	CertPath    string          `json:"certPath"`    // 服务证书
	KeyPath     string          `json:"keyPath"`     // 服务私钥
	LogConfig   LogConfig       `json:"logConfig"`   // 日志配置
	CacheConf   cache.CacheConf `json:"cacheConf"`   // 缓存信息配置
	IpWhiteList []string        `json:"ipWhiteList"` // IP白名单列表

	AliyunConf     cdnprovider.AliAuthConfiguration      `json:"aliyunConf"`     // 阿里云鉴权信息配置
	BaiduyunConf   cdnprovider.BaiduyunAuthConfiguration `json:"baiduyunConf"`   // 百度云鉴权信息配置
	BaishanyunConf cdnprovider.BaishanAuthConfiguration  `json:"baishanyunConf"` // 白山云鉴权信息配置
	DilianConf     cdnprovider.DilianAuthConfiguration   `json:"dilianConf"`     // 帝联鉴权信息配置
	KuaiwangConf   cdnprovider.KaiuwangAuthConfiguration `json:"kuaiwangConf"`   // 快网鉴权信息配置
	UpyunConf      cdnprovider.UpyunAuthConfiguration    `json:"upyunConf"`      // 又拍鉴权信息配置
	WangsuConf     cdnprovider.WangsuAuthConfiguration   `json:"wangsuConf"`     // 网宿鉴权信息配置
	YunfanConf     cdnprovider.YunfanAuthConfiguration   `json:"yunfanConf"`     // 云帆鉴权信息配置
	TencentConf    cdnprovider.TencentAuthConfiguration  `json:"tencentConf"`    // 腾讯云鉴权信息配置
}

type LogConfig struct {
	LogDir    string `json:"logdir"`    // 日志目录
	Chunkbits uint   `json:"chunkbits"` // 日志大小
}

func InitServiceConf(path string) (err error) {
	ServerConf = &CdnProviderAuthConf{}
	cfgData, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	if err = json.Unmarshal(cfgData, ServerConf); err != nil {
		return
	}
	return
}
