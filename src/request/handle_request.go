package request

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"

	"conf"
	"define"
	"errcode"
	"filelog"
	"util"
)

type RequestHandle struct {
	HandleRequestFunc handleRequestFunc
}

type handleRequestFunc interface {
	HandleRequest(*filelog.ReqLogger, http.ResponseWriter, *http.Request) (interface{}, error)
}

func (rh *RequestHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log := filelog.NewReqLogger(w, r)

	remoteIP := getRemoteIp(r)
	log.Infof("==> ServeHTTP request from ip: %s", remoteIP)

	var err error
	var resData interface{}

	defer func() {
		if err != nil {
			log.Info("==> ServeHTTP failed err: ", err)
			defErr, ok := err.(*define.ErrorInfo)
			if !ok {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(fmt.Sprintf(`{"code":40000,"errMsg":"%s"}`, err.Error())))
				return
			}
			w.WriteHeader(defErr.HttpCode())
			w.Write([]byte(defErr.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		if resData != nil {
			resDataByte, _ := json.Marshal(resData)
			w.Write(resDataByte)
		}
	}()

	// 仅支持 POST 请求
	if r.Method != http.MethodPost {
		log.Error("forbidden request method: ", r.Method)
		err = errcode.ForbiddenRequestErr
		return
	}

	// 校验客户端的 IP 是否在白名单内
	whiteIpCount := len(conf.ServerConf.IpWhiteList)
	if whiteIpCount > 0 && !util.InStringSlice(remoteIP, conf.ServerConf.IpWhiteList) {
		err = errcode.ForbiddenRequestErr
		return
	}

	resData, err = rh.HandleRequestFunc.HandleRequest(log, w, r)
	return
}

func getRemoteIp(r *http.Request) (remoteIp string) {
	if remoteIp = r.Header.Get("X-Forwarded-For"); remoteIp != "" {
		remoteIp = strings.Split(remoteIp, ",")[0]
	} else {
		remoteIp, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	return
}
