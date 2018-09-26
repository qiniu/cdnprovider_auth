package errcode

import (
	"define"
)

const (
	ErrBadRequest            = 400000
	ErrInvalidParams         = 400001
	ErrInvalidCdnProvider    = 400002
	ErrNullAuthInfo          = 400003
	ErrUnRecognizedCacheType = 400004
	ErrInvalidAuthConf       = 400005

	ErrForbiddenRequest = 401000

	ErrNoSuchAuthInfoCache = 404000

	ErrQueryRedisData   = 500000
	ErrDeleteRedisData  = 500001
	ErrUpsertRedisData  = 500002
	ErrQueryMemoryData  = 500003
	ErrDeleteMemoryData = 500004
	ErrUpsertMemoryData = 500005
	ErrLoadConfFile     = 500006
)

var (
	InvalidParamsErr         = define.NewError(ErrInvalidParams, "invalid params")
	ForbiddenRequestErr      = define.NewError(ErrForbiddenRequest, "forbidden request")
	InvalidCdnProviderErr    = define.NewError(ErrInvalidCdnProvider, "invalid cdn provider")
	NullAuthInfoErr          = define.NewError(ErrNullAuthInfo, "null auth info")
	NoSuchAuthInfoCacheErr   = define.NewError(ErrNoSuchAuthInfoCache, "no such auth info cache")
	QueryRedisDataErr        = define.NewError(ErrQueryRedisData, "query data from redis failed")
	UnRecognizedCacheTypeErr = define.NewError(ErrUnRecognizedCacheType, "unrecognized cache type")
	InvalidAuthConfErr       = define.NewError(ErrInvalidAuthConf, "invalid auth conf")
)
