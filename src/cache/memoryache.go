package cache

import (
	"encoding/json"
	"time"

	"github.com/coocood/freecache"

	"define"
	"errcode"
)

type MemoryCacheClient struct {
	*freecache.Cache
}

func NewMemoryCache() (client *MemoryCacheClient, err error) {
	cacheSize := 100 * 1024 * 1024
	cache := freecache.NewCache(cacheSize)
	return &MemoryCacheClient{cache}, nil
}

func (c *MemoryCacheClient) Upsert(cdnProvider define.CDNProvider, authRes define.AuthRes, expiration time.Duration) (err error) {
	data, _ := json.Marshal(authRes)
	err = c.Cache.Set([]byte(cdnProvider), data, int(expiration.Seconds()))
	if err != nil {
		return define.NewError(errcode.ErrUpsertMemoryData, err.Error())
	}
	return
}
func (c *MemoryCacheClient) Delete(cdnProvider define.CDNProvider) (err error) {
	c.Cache.Del([]byte(cdnProvider))
	return
}
func (c *MemoryCacheClient) Get(cdnProvider define.CDNProvider) (authRes define.AuthRes, err error) {
	data, err := c.Cache.Get([]byte(cdnProvider))
	if err != nil {
		if err == freecache.ErrNotFound {
			err = errcode.NoSuchAuthInfoCacheErr
			return
		}
		err = define.NewError(errcode.ErrQueryMemoryData, err.Error())
		return
	}
	err = json.Unmarshal(data, &authRes)
	if err != nil {
		err = define.NewError(errcode.ErrQueryMemoryData, err.Error())
		return
	}
	return
}

func (c *MemoryCacheClient) Close() {
	c.Cache.Clear()
	return
}
