package filelog

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	reqidKey = "X-Reqid"
)

var ServerLogger *Logger

type Logger struct {
	w *Writer
}

func NewLogger(dir, prefix string, timeMode int64, chunkBits uint) (r *Logger, err error) {
	_, err = os.Stat(dir)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	w, err := NewWriter(dir, prefix, timeMode, chunkBits)
	if err != nil {
		return
	}
	r = &Logger{w}
	return
}

func (r *Logger) Close() (err error) {
	return r.w.Close()
}

func (r *Logger) Log(msg []byte) (err error) {
	msg = append(msg, '\n')
	_, err = r.w.Write(msg)
	return
}

func (r *Logger) logLevel(reqId, level string, v ...interface{}) (err error) {
	var formatMsg string
	if reqId != "" {
		formatMsg = fmt.Sprintf("%s [%s] [%s] %s", time.Now().In(time.Local).Format(time.RFC3339), reqId, level, fmt.Sprintln(v...))
	} else {
		formatMsg = fmt.Sprintf("%s [%s] %s", time.Now().In(time.Local).Format(time.RFC3339), level, fmt.Sprintln(v...))
	}
	_, err = r.w.Write([]byte(formatMsg))
	return
}

type ReqLogger struct {
	ReqId string
	*Logger
}

func NewReqLogger(w http.ResponseWriter, req *http.Request) (reqLog *ReqLogger) {
	reqId := req.Header.Get(reqidKey)
	if reqId == "" {
		reqId = genReqId()
		req.Header.Set(reqidKey, reqId)
	}
	h := w.Header()
	h.Set(reqidKey, reqId)
	return &ReqLogger{ReqId: reqId, Logger: ServerLogger}
}

func (r *ReqLogger) Info(v ...interface{}) (err error) {
	return r.logLevel(r.ReqId, "INFO", v...)
}

func (r *ReqLogger) Infof(format string, v ...interface{}) (err error) {
	return r.Info(fmt.Sprintf(format, v...))
}

func (r *ReqLogger) Error(v ...interface{}) (err error) {
	return r.logLevel(r.ReqId, "ERROR", v...)
}

func (r *ReqLogger) Errorf(format string, v ...interface{}) (err error) {
	return r.Error(fmt.Sprintf(format, v...))
}

func (r *ReqLogger) Warn(v ...interface{}) (err error) {
	return r.logLevel(r.ReqId, "WARN", v...)
}

func (r *ReqLogger) Warnf(format string, v ...interface{}) (err error) {
	return r.Warn(fmt.Sprintf(format, v...))
}

func genReqId() string {
	var b [12]byte
	binary.LittleEndian.PutUint32(b[:], uint32(os.Getpid()))
	binary.LittleEndian.PutUint64(b[4:], uint64(time.Now().UnixNano()))
	return base64.URLEncoding.EncodeToString(b[:])
}

func InitServerLogger(dir, prefix string, timeMode int64, chunkBits uint) (err error) {
	ServerLogger, err = NewLogger(dir, prefix, timeMode, chunkBits)
	return
}
