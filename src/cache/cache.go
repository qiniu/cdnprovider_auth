package cache

import (
	"time"

	"define"
	"errcode"
)

var Cache CacheClient

type CacheClient interface {
	Upsert(cdnProvider define.CDNProvider, authRes define.AuthRes, expiration time.Duration) (err error)
	Delete(cdnProvider define.CDNProvider) (err error)
	Get(cdnProvider define.CDNProvider) (authRes define.AuthRes, err error)
	Close()
}

type CacheConf struct {
	CacheType      string         `json:"cacheType"`      // 缓存方式
	RedisCacheConf RedisCacheConf `json:"redisCacheConf"` // redis 缓存配置
}

type RedisCacheConf struct {
	Host     string `json:"host"`     // redis host
	Password string `json:"password"` // redis 密码
}

func NewCacheClient(cacheConf CacheConf) (cacheClient CacheClient, err error) {
	switch cacheConf.CacheType {
	case "redis":
		return NewRedisCacheClient(cacheConf.RedisCacheConf.Host, cacheConf.RedisCacheConf.Password)
	case "memory":
		return NewMemoryCache()
	case "": // 无缓存配置
		return nil, nil
	default:
		return nil, errcode.UnRecognizedCacheTypeErr
	}
	return
}

func InitCacheClient(cacheConf CacheConf) (err error) {
	Cache, err = NewCacheClient(cacheConf)
	return err
}
