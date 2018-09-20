package cache

import (
	"encoding/json"
	"time"

	redis "gopkg.in/redis.v5"

	"define"
	"errcode"
)

type RedisCacheClient struct {
	client *redis.Client
}

func NewRedisCacheClient(url, passwd string) (client *RedisCacheClient, err error) {
	client = &RedisCacheClient{}
	client.client = redis.NewClient(&redis.Options{
		Addr:     url,
		Password: passwd,
		DB:       0,
	})
	_, err = client.client.Ping().Result()
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *RedisCacheClient) Upsert(cdnProvider define.CDNProvider, authRes define.AuthRes, expiration time.Duration) (err error) {
	data, _ := json.Marshal(authRes)
	err = c.client.Set(cdnProvider.String(), data, expiration).Err()
	if err != nil {
		return define.NewError(errcode.ErrUpsertRedisData, err.Error())
	}
	return nil
}

func (c *RedisCacheClient) Delete(cdnProvider define.CDNProvider) (err error) {
	err = c.client.Del(cdnProvider.String()).Err()
	if err != nil {
		return define.NewError(errcode.ErrDeleteRedisData, err.Error())
	}
	return
}

func (c *RedisCacheClient) Get(cdnProvider define.CDNProvider) (authRes define.AuthRes, err error) {
	value, err := c.client.Get(cdnProvider.String()).Bytes()
	if err != nil {
		if err == redis.Nil {
			err = errcode.NoSuchAuthInfoCacheErr
			return
		}
		err = define.NewError(errcode.ErrQueryRedisData, err.Error())
		return
	}
	err = json.Unmarshal(value, &authRes)
	if err != nil {
		err = define.NewError(errcode.ErrQueryRedisData, err.Error())
		return
	}
	return
}

func (c *RedisCacheClient) Close() {
	c.client.Close()
}
