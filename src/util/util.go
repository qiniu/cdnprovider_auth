package util

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
)

func GetResponseBody(body io.ReadCloser, contentLength int64) ([]byte, error) {
	defer func() {
		io.Copy(ioutil.Discard, body)
		body.Close()
	}()
	if contentLength == 0 || body == nil {
		return []byte{}, nil
	}

	var b bytes.Buffer
	_, err := b.ReadFrom(body)
	if err != nil {
		return []byte{}, err
	}
	return b.Bytes(), nil
}

func ParseReqArgs(body io.ReadCloser, contentLength int64, req interface{}) (err error) {
	reqBodyData, err := GetResponseBody(body, contentLength)
	if err != nil {
		return
	}
	err = json.Unmarshal(reqBodyData, req)
	if err != nil {
		return
	}
	return
}

func DoMd5(param string) (result string) {
	md5byte := []byte(param)
	hasher := md5.New()
	hasher.Write(md5byte)
	return hex.EncodeToString(hasher.Sum(nil))
}

func HmacSha256Hex(key, signStr string) string {
	hasher := hmac.New(sha256.New, []byte(key))
	hasher.Write([]byte(signStr))
	return hex.EncodeToString(hasher.Sum(nil))
}

func InStringSlice(s string, sslice []string) bool {
	for _, v := range sslice {
		if s == v {
			return true
		}
	}
	return false
}
