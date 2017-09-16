package main

import (
	"neo/config"
	"neo/core"
	"neo/dal"
	"neo/util"
)

var (
	conf *config.Config
	db   interface{}
	err  error
)

func init() {
	conf = config.GetConfig()
	db, err = dal.Open(conf)
	if err != nil {
		panic(err)
	}
}

func main() {
	// get rate
	accessor := &core.RateAccessor{}
	rate, err := accessor.GetLatestRateInfo(db)
	if err != nil {
		panic(err)
	}

	util.SendMail(rate, conf)
}
