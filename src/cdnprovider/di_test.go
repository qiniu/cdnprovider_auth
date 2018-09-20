package cdnprovider

import (
	"define"
	"fmt"
	"testing"
)

func TestDl(t *testing.T) {
	s := DilianAuthService{
		DilianAuthConf: DilianAuthConfiguration{AccessKeyId: "accesskeyidexample8", AccessKey: "1234567890abcdef"},
	}
	authReq := define.AuthReq{
		Path:   "/v3/api/push",
		Method: "POST",
		Params: map[string]interface{}{
			"ParamString": "http://dl-u0rjq.coolplayer.net/a.txt",
		},
	}
	res, err := s.GenerateSignatureMethod(nil, authReq)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
}
