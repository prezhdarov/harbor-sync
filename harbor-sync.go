package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-kit/log/level"
	"github.com/prezhdarov/prometheus-exporter/config"
	"github.com/prometheus/common/promlog"

	harbor "github.com/prezhdarov/harbor-sync/pkg/harbor"
)

var (
	hbSrc     = flag.String("src.harbor", "", "Source harbor host")
	hbSrcUser = flag.String("src.username", "", "Source username")
	hbSrcPass = flag.String("src.password", "", "Source password")
	hbDst     = flag.String("dst.harbor", "", "Source harbor host")
	hbDstUser = flag.String("dst.username", "", "Source username")
	hbDstPass = flag.String("dst.password", "", "Source password")

	logLevel  = flag.String("log.level", "debug", "Log Level minimums. Available options are: debug,info,warn and error")
	logFormat = flag.String("log.format", "logfmt", "Log output format. Available options are: logfmt and json")
)

func usage() {
	const s = `
fortim-exporter collects metrics data from Forti Managers. 
`
	config.Usage(s)
}

func setLogger(lf, ll *string) *promlog.Config {
	promlogFormat := &promlog.AllowedFormat{}
	promlogFormat.Set(*lf)

	promlogLevel := &promlog.AllowedLevel{}
	promlogLevel.Set(*ll)

	promlogConfig := &promlog.Config{}
	promlogConfig.Format = promlogFormat
	promlogConfig.Level = promlogLevel

	return promlogConfig
}

func main() {

	flag.CommandLine.SetOutput(os.Stdout)
	flag.Usage = usage
	config.Parse()

	logger := promlog.New(setLogger(logFormat, logLevel))

	if *hbSrc == "" || *hbSrcUser == "" || *hbSrcPass == "" {
		level.Error(logger).Log("msg", "Source host, username and password cannot be empty")
		os.Exit(0)
	}

	if *hbDst == "" || *hbDstUser == "" || *hbDstPass == "" {
		level.Error(logger).Log("msg", "Destination host, username and password cannot be empty")
		os.Exit(0)
	}

	err := harbor.SyncLabels(&harbor.Harbor{Host: hbSrc, User: hbSrcUser, Pass: hbSrcPass}, &harbor.Harbor{Host: hbDst, User: hbDstUser, Pass: hbDstPass}, logger)
	if err != nil {
		fmt.Printf("can't do this: %s\n", err)
	}

}
