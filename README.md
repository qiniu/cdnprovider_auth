# CDNHub 统一授权服务使用说明 [![Build Status](https://api.travis-ci.org/qiniu/logkit.svg)](http://travis-ci.org/qiniu/cdnprovider_auth)

## 服务简述

CDNHub 统一授权服务是七牛 CDNHub 平台为用户提供的针对第三方 CDN 厂商的通用授权服务。基于第三方厂商提供的账号信息为 CdnHub 平台提供授权信息，允许用户在该平台上统一管理多厂商 CDN 资源，实现刷新预取，配置修改，用量统计等功能。现已支持的第三方 CDN 厂商：阿里云、百度云、白山云、帝联科技、快网、腾讯云、网宿、云帆加速、又拍云。

## 服务部署

### 源码下载编译
```
git clone https://github.com/qiniu/cdnprovider_auth
cd cdnprovider_auth
make & cd bin
```

### 二进制文件下载
* [Mac](http://pebc2c9b2.bkt.clouddn.com/mac/cdnprovider_auth)
* [Linux](http://pebc2c9b2.bkt.clouddn.com/linux/cdnprovider_auth)
* [Windows](http://pebc2c9b2.bkt.clouddn.com/windows/cdnprovider_auth.exe)


### 服务配置
下载 cdnprovider_auth.conf 或自建配置文件, 配置其中将使用CDN厂商的授权信息。**请注意：** 配置文件不支持 "//" 注释,自建配置文件需删除注释。
```
{
	// 服务绑定地址
	"bindHost": "0.0.0.0:8000",
	
	// 如部署 https 服务，需配置证书公钥和私钥文件的地址，否则不需要填写 certPath 和 keyPath 两项
	"certPath": "/home/cert/cert.crt",
	"keyPath": "/home/cert/cert.key",
	
	// 服务日志配置，必填
	// logdir: 日志文件存储目录，必填
	// chunkbits: 日志切割大小，取值范围：26(2^26:64M) - 32(2^32:4G)，默认值为32
	"logConfig": {
		"logdir": "/home/qiniu/cdn/log",
		"chunkbits": 29
	},
	
	// 部分厂商的授权信息需要缓存，如选择不缓存，可不配置此项
	// cacheType: 授权信息缓存方式, 取值范围：redis(redis 缓存)，memory(内存缓存)
	// redisCacheConf: 选择使用 redis 作为缓存方式时，需填写 redis 的地址和密码，否则不需要配置此项
	"cacheConf": {
		"cacheType": "redis",
		"redisCacheConf": {
			"host": "127.0.0.1:6379",
			"password": "root"
		}
	},
	
	// 服务白名单，服务仅接收来自白名单中 ip 的请求，其余请求将拒绝。如果 ipWhiteList 为空，则不会检验客户端 ip
	"ipWhiteList": [
		""
	],
	
	// 使用阿里云 CDN 业务的需配置此项
	// 在 [阿里云官网]-[管理控制台]-[accesskeys] 中获取鉴权信息
	"aliyunConf": {
		"accessKeyId": "aliyunAccessKeyId",
		"secretKey": "aliyunSecretKey"
	},
	
	// 使用百度云 CDN 业务的需配置此项
	// 在 [百度云官网]-[安全认证]-[Access Key] 中获取鉴权信息
	"baiduyunConf": {
		"accessKeyId": "baiduyunAccessKey",
		"secretKey": "baiduSecretKey"
	},
	
	// 使用白山云 CDN 业务的需配置此项
	// 白山鉴权信息请联系其官方人员获取
	"baishanyunConf": {
		"token": "baishanyunToken"
	},
	
	// 使用帝联科技 CDN 业务的需配置此项
	// 在 [帝联科技官网]-[个人面板]-[密钥管理] 中获取鉴权信息
	"dilianConf": {
		"accessKeyId": "dilianAccessKeyId",
		"accessKey": "dilianAccessKey"
	},
	
	// 使用快网 CDN 业务的需配置此项，且推荐配置缓存策略，如部署多进程必须配置 redis 缓存，具体请参考快网 API 手册
	// 快网鉴权信息请联系其官方人员获取
	"kuaiwangConf": {
		"appId": "kuaiwangAppId",
		"appSecret": "kuaiwangAppSecret"
	},
	
	// 使用腾讯 CDN 业务的需配置此项
	// 在 [腾讯云官网]-[项目管理]-[云API密钥]-[API密钥管理] 中获取鉴权信息
	"tencentConf": {
		"secretId": "tencentSeretId",
		"secretKey": "tencentSecretKey"
	},
	
	// 使用网宿 CDN 业务的需配置此项
	// 网宿鉴权信息请联系其官方人员获取
	"wangsuConf": {
		"userName": "wangsuUserName",
		"apiKey": "wangsuApiKey",
		"password": "wangsuPassword"
	},
	
	// 使用云帆加速 CDN 业务的需配置此项
	// 在 [云帆加速官网]-[安全认证]-[Access Key] 中获取鉴权信息
	"yunfanConf": {
		"accessKey": "yunfanAccessKey",
		"secretKey": "yunfanSecretKey"
	},
    
	// 使用又拍云 CDN 业务的需配置此项
	// 又拍云鉴权信息请联系其官方人员获取
	"upyunConf": {
		"token": "upyunToken"
	},
}
```
### 服务启动
使用二进制文件启动服务（指定配置文件路径）：

```
./cdnprovider_auth cdnprovider_auth.conf
```

使用 docker 镜像启动服务：
```
docker pull reg.qiniu.com/cdnproviderauth/cdnprovider_auth:v1
docker run -d  -p 8080:8080  -v /local/cdnprovider_auth.conf:/app/auth.conf reg.qiniu.com/cdnproviderauth/cdnprovider_auth:v1
```

镜像中，cdnprovider\_auth 读取 /app/ 目录下的配置文件 auth.conf，需要把本地的配置文件挂载到镜像中去才能启动，比如本地的配置文件为 /local/cdnprovider_auth.conf，挂载到镜像中的 /app/auth.conf 文件。

### 服务更新

下载源码生成二进制文件或直接下载二进制文件，将已经在运行的 cdnprovider_auth 二进制文件替换成新的二进制文件，重启服务方式：

1. kill cdnprovider\_auth 服务进程，并重新启动
2. 热重启服务：kill -HUP cdnprovider\_auth 服务进程
