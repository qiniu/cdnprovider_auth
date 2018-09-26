package util

import (
	"crypto/tls"
	"io"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"filelog"
)

func DoHTTPRequest(log *filelog.ReqLogger, url string, method string, body string, headers map[string]string) (res *http.Response, err error) {
	var reqBody io.Reader
	if body != "" {
		reqBody = strings.NewReader(body)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		log.Error("new request failed, err: ", err)
		return
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}
	reqBytes, _ := httputil.DumpRequest(req, true)
	log.Info("http request:\n", string(reqBytes))

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	cl := &http.Client{Transport: tr, Timeout: 5 * time.Second}

	if res, err = cl.Do(req); err != nil {
		log.Error("do request err: ", err)
		return
	}
	resBytes, _ := httputil.DumpResponse(res, true)
	log.Info("http response:\n", string(resBytes))
	return
}
