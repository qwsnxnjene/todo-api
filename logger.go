package main

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var Logger *zap.Logger

func initLogger() {
	var err error
	if viper.GetString("ENV") == "production" {
		Logger, err = zap.NewProduction()
	} else {
		Logger, err = zap.NewDevelopment()
	}
	if err != nil {
		panic(err)
	}
}
