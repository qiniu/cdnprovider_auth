package request

import (
	"net/http"

	"cache"
	"cdnprovider"
	"conf"
	"define"
	"errcode"
	"filelog"
	"util"
)

type SignatureService struct {
	GenerateSignatureMethodMap map[define.CDNProvider]define.GenerateSignatureMethod
}

/*
	根据请求信息生成相应供应商的鉴权信息
*/
func (s *SignatureService) HandleRequest(log *filelog.ReqLogger, w http.ResponseWriter, r *http.Request) (authRes interface{}, err error) {
	if s.GenerateSignatureMethodMap == nil {
		s.initSignatureMethodMap()
	}

	authReq := define.AuthReq{}
	if err = util.ParseReqArgs(r.Body, r.ContentLength, &authReq); err != nil {
		log.Error("unmarshal request body to target struct failed, err: ", err)
		err = errcode.InvalidParamsErr
		return
	}

	log.Infof("generate signature args: %#v", authReq)
	// 判断是否为适配的CDN供应商
	if _, ok := define.ValidCDNProviderMap[authReq.CdnProvider]; !ok {
		log.Error("unrecognized cdn provider: ", authReq.CdnProvider)
		err = errcode.InvalidCdnProviderErr
		return
	}
	if authReq.Params == nil {
		authReq.Params = make(map[string]interface{}, 0)
	}

	// 获取对应的CDN供应商的鉴权信息
	authRes, err = s.GenerateSignatureMethodMap[authReq.CdnProvider].GenerateSignatureMethod(log, authReq)
	if err != nil {
		log.Errorf("generate signature for %s failed, err: %v", authReq.CdnProvider, err)
		return
	}
	return
}

func (s *SignatureService) initSignatureMethodMap() {
	s.GenerateSignatureMethodMap = make(map[define.CDNProvider]define.GenerateSignatureMethod, len(define.ValidCDNProviderMap))
	for cdn, _ := range define.ValidCDNProviderMap {
		switch cdn {
		case define.CDNProviderAliyun:
			aliAuthService := &cdnprovider.AliAuthService{conf.ServerConf.AliyunConf}
			s.GenerateSignatureMethodMap[cdn] = aliAuthService
		case define.CDNProviderBaidu:
			baiduAuthService := &cdnprovider.BaiduAuthService{conf.ServerConf.BaiduyunConf}
			s.GenerateSignatureMethodMap[cdn] = baiduAuthService
		case define.CDNProviderBaishanyun:
			baishanyunAuthService := &cdnprovider.BaishanyunAuthService{conf.ServerConf.BaishanyunConf}
			s.GenerateSignatureMethodMap[cdn] = baishanyunAuthService
		case define.CDNProviderDilian:
			dilianAuthService := &cdnprovider.DilianAuthService{conf.ServerConf.DilianConf}
			s.GenerateSignatureMethodMap[cdn] = dilianAuthService
		case define.CDNProviderKuaiwang:
			kuaiwangService := &cdnprovider.KuaiwangAuthService{conf.ServerConf.KuaiwangConf, cache.Cache}
			s.GenerateSignatureMethodMap[cdn] = kuaiwangService
		case define.CDNProviderTencent:
			tencentAuthService := &cdnprovider.TencentAuthService{conf.ServerConf.TencentConf}
			s.GenerateSignatureMethodMap[cdn] = tencentAuthService
		case define.CDNProviderUPYun:
			upyunAuthService := &cdnprovider.UpyunAuthService{conf.ServerConf.UpyunConf}
			s.GenerateSignatureMethodMap[cdn] = upyunAuthService
		case define.CDNProviderWangsu:
			wangsuAuthService := &cdnprovider.WanggsuService{conf.ServerConf.WangsuConf}
			s.GenerateSignatureMethodMap[cdn] = wangsuAuthService
		case define.CDNProviderYunfan:
			yunfanService := &cdnprovider.YunfanAuthService{conf.ServerConf.YunfanConf}
			s.GenerateSignatureMethodMap[cdn] = yunfanService
		}
	}
}
