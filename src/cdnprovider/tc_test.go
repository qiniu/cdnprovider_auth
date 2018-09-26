package cdnprovider

import (
	"testing"

	"define"
	"fmt"
)

func TestTc(t *testing.T) {
	s := TencentAuthService{TencentAuthConf: TencentAuthConfiguration{SecretId: "AKIDT8G5AsY1D3MChWooNq1rFSw1fyBVCX9D", SecretKey: "pxPgRWDbCy86ZYyqBTDk7WmeRZSmPco0"}}

	args := define.AuthReq{
		Host:   "cdn.api.qcloud.com",
		Path:   "/v2/index.php",
		Method: "GET",
		Params: map[string]interface{}{
			"offset":    0,
			"Action":    "DescribeCdnHosts",
			"Nonce":     "13029",
			"SecretId":  "AKIDT8G5AsY1D3MChWooNq1rFSw1fyBVCX9D",
			"Timestamp": "1463122059",
			"limit":     "10",
		},
	}
	res, _ := s.GenerateSignatureMethod(nil, args)
	fmt.Println(res)
}
