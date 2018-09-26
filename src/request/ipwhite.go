package request

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"conf"
	"define"
	"errcode"
	"filelog"
	"util"
)

type IpWhiteService struct {
	ConfPath string
}

/*
	更新服务配置的 IP 白名单列表
*/
func (s *IpWhiteService) HandleRequest(log *filelog.ReqLogger, w http.ResponseWriter, r *http.Request) (authRes interface{}, err error) {
	ipWhiteReq := define.IpWhiteListReq{}
	if err = util.ParseReqArgs(r.Body, r.ContentLength, &ipWhiteReq); err != nil {
		log.Error("unmarshal request body to target struct failed, err: ", err)
		err = errcode.InvalidParamsErr
		return
	}

	serverConf := &conf.CdnProviderAuthConf{}
	cfgData, err := ioutil.ReadFile(s.ConfPath)
	if err != nil {
		log.Errorf("read conf file from %s failed, err: %v", s.ConfPath, err)
		err = define.NewError(errcode.ErrLoadConfFile, err.Error())
		return
	}
	if err = json.Unmarshal(cfgData, serverConf); err != nil {
		log.Errorf("parse conf file from %s failed, err: %v", s.ConfPath, err)
		err = define.NewError(errcode.ErrLoadConfFile, err.Error())
		return
	}

	// 更新白名单
	serverConf.IpWhiteList = ipWhiteReq.IpWhiteList
	newData, _ := json.MarshalIndent(serverConf, "", "\t")

	// 更新进程的配置
	conf.ServerConf.IpWhiteList = ipWhiteReq.IpWhiteList

	// 更新当前配置文件
	err = ioutil.WriteFile(s.ConfPath, []byte(newData), 0644)
	if err != nil {
		log.Errorf("write conf file from %s failed, err: %v", s.ConfPath, err)
		err = define.NewError(errcode.ErrLoadConfFile, err.Error())
		return
	}
	return
}
