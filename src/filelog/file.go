package filelog

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"
)

type File struct {
	sync.Mutex

	f    *os.File
	name string
	t    int64
	off  int64
	ref  int32
}

func (r *File) Write(p []byte) (n int, err error) {
	r.Lock()
	n, err = r.f.Write(p)
	r.Unlock()

	atomic.AddInt64(&r.off, int64(len(p)))
	return
}

func (r *File) Release() {
	ref := atomic.AddInt32(&r.ref, -1)
	if ref == 0 {
		r.f.Close()
	}
}

func (r *File) Acquire() {
	atomic.AddInt32(&r.ref, 1)
}

const (
	DefaultTimeMode  = 3600 // 默认1小时
	DefaultChunkBits = 32   // 默认大小4G
	MinChunkBits     = 26   // min大小64M
)

type Writer struct {
	prefix    string
	timeMode  int64
	sizeLimit int64
	file      *File
	sync.RWMutex
}

// timeMode 单位sec。保证可以平均切分；不能过小，大于1s; 不能过大，小于1day
func NewWriter(dir, prefix string, timeMode int64, chunkBits uint) (writer *Writer, err error) {
	if timeMode <= 0 {
		timeMode = DefaultTimeMode
	}
	if chunkBits == 0 {
		chunkBits = DefaultChunkBits
	} else if chunkBits < MinChunkBits {
		chunkBits = MinChunkBits
	}

	if 86400%timeMode != 0 {
		return nil, errors.New("wrong timeMode")
	}

	writer = &Writer{filepath.Join(dir, prefix), timeMode, 1 << chunkBits, nil, sync.RWMutex{}}
	if prefix == "" {
		writer.prefix += string(filepath.Separator)
	}
	return
}

func genFileName(now int64) string { // TODO
	t := time.Unix(now, 0)
	return t.Format("20060102150405")
}

func now(t time.Time, timeMode int64) int64 {
	if timeMode <= 3600 {
		return t.Truncate(time.Duration(timeMode * 1e9)).Unix()
	}
	today := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return today.Unix() + int64(t.Sub(today))/1e9/timeMode*timeMode
}

func (r *Writer) getFile() (f *File, err error) {
	t := time.Now()
	nowT := now(t, r.timeMode)
	name := ""

	r.RLock()
	f = r.file
	if f != nil && f.t == nowT {
		if atomic.LoadInt64(&f.off) < r.sizeLimit {
			f.Acquire()
			r.RUnlock()
			return
		}
		name = r.prefix + t.Format("20060102150405") // size rotate 时按秒来分割
	} else {
		name = r.prefix + genFileName(nowT)
	}
	r.RUnlock()

	fp, err := os.OpenFile(name, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return
	}

	info, err := fp.Stat()
	if err != nil {
		return
	}
	f = &File{f: fp, name: name, t: nowT, off: info.Size(), ref: 2}

	r.Lock()
	of := r.file
	if of != nil && of.name == name {
		f.ref = 1
		f, of = of, f
		f.Acquire()
	} else {
		r.file = f
	}
	r.Unlock()

	if of != nil {
		of.Release()
	}

	return
}

func (r *Writer) Write(data []byte) (n int, err error) {
	f, err := r.getFile()
	if err != nil {
		return
	}
	defer f.Release()

	return f.Write(data)
}

func (r *Writer) Close() error {
	r.Lock()
	defer r.Unlock()

	if r.file != nil {
		r.file.Release()
	}
	r.file = nil
	return nil
}
