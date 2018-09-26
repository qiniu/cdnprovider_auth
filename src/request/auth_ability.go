package request

import (
	"net/http"

	"filelog"
)

type AuthAbilityService struct{}

/*
	仅用于检测服务的可用性，不做其他处理
*/
func (s *AuthAbilityService) HandleRequest(log *filelog.ReqLogger, w http.ResponseWriter, r *http.Request) (authRes interface{}, err error) {
	log.Info("check server ability success!")
	return nil, nil
}
