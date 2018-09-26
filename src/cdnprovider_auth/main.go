package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cache"
	"conf"
	"filelog"
	"gracehttp"
	"request"
)

func usage() {
	fmt.Fprintf(os.Stderr, "USAGE\n")
	fmt.Fprintf(os.Stderr, "  %s [conf loaded from a file] \n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\n")
}

func main() {
	serveMux := http.NewServeMux()
	serveMux.Handle("/auth", &request.RequestHandle{&request.SignatureService{}})
	serveMux.Handle("/auth/ability", &request.RequestHandle{&request.AuthAbilityService{}})
	serveMux.Handle("/auth/ipwhite", &request.RequestHandle{&request.IpWhiteService{os.Args[1]}})

	internalSrv := &http.Server{Addr: conf.ServerConf.BindHost, Handler: serveMux}
	graceInternalSrv := &gracehttp.Server{
		Name:   "qiniu-cdnprovider-auth-server",
		Server: internalSrv,
	}

	// set https config
	if conf.ServerConf.CertPath != "" && conf.ServerConf.KeyPath != "" {
		graceInternalSrv.CertFile = conf.ServerConf.CertPath
		graceInternalSrv.KeyFile = conf.ServerConf.KeyPath
	}

	defer filelog.ServerLogger.Close()

	gracehttpServers := []*gracehttp.Server{graceInternalSrv}
	authApp := gracehttp.NewApp(gracehttpServers)
	err := authApp.Run()
	if err != nil {
		log.Fatalf("Failed to run auth app, err: %v", err)
	}
}

func init() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}
	var err error
	// init server conf
	confPath := os.Args[1]
	if err = conf.InitServiceConf(confPath); err != nil {
		log.Fatal(err)
	}

	// init server log
	if err = filelog.InitServerLogger(conf.ServerConf.LogConfig.LogDir, "", int64(time.Hour.Seconds())*24, conf.ServerConf.LogConfig.Chunkbits); err != nil {
		log.Fatal(err)
	}

	// init cache client
	err = cache.InitCacheClient(conf.ServerConf.CacheConf)
	if err != nil {
		log.Fatal(err)
	}
}
