package config

import (
	"errors"
	"os"
	"strings"
	"sync"
)

var (
	globalSetting *GlobalSetting
	configFile    string
	initOnce      sync.Once
)

func Setting() *GlobalSetting {
	return globalSetting
}

func MustInitialize() {
	initOnce.Do(func() {
		initSetting()
	})
}

func initSetting() {
	if strings.ToLower(os.Getenv("APP_ENV")) == "test" {
		configFile = "kir.test.json"
	} else {
		configFile = "app.json"
	}

	setting := &GlobalSetting{}
	if err := Initialize(configFile, setting); err != nil {
		panic(err)
	}

	if setting == nil {
		panic(errors.New("failed load config"))
	}

	globalSetting = setting
}
