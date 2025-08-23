package config

import (
	"fmt"
	"os"

	"github.com/allentom/harukap/config"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

var DefaultConfigProvider *config.Provider

func InitConfigProvider() error {
	var err error
	customConfigPath := os.Getenv("CONFIG_PATH")
	if customConfigPath != "" {
		logrus.Info("using custom config path: ", customConfigPath)
	}
	DefaultConfigProvider, err = config.NewProvider(func(provider *config.Provider) {
		ReadConfig(provider)
	}, customConfigPath)
	return err
}

var Instance Config

type AuthConfig struct {
	Name   string
	Enable bool
	AppId  string
	Secret string
	Url    string
	Type   string
}

type ProxyConfig struct {
	Target      string `json:"target"`
	Prefix      string `json:"prefix"`
	Name        string `json:"name"`
	UseNacos    bool   `json:"useNacos" mapstructure:"useNacos"`
	ServiceName string `json:"serviceName" mapstructure:"serviceName"`
	Group       string `json:"group" mapstructure:"group"`
	Scheme      string `json:"scheme" mapstructure:"scheme"`
}

type Config struct {
	YouAuthConfig       *AuthConfig
	Auths               []*AuthConfig
	YouAuthConfigPrefix string
	EnableAnonymous     bool
	ServiceProxy        []*ProxyConfig
}

func ReadConfig(provider *config.Provider) {
	configer := provider.Manager
	configer.SetDefault("addr", ":8000")
	configer.SetDefault("application", "My Service")
	configer.SetDefault("instance", "main")

	Instance = Config{}
	Instance.ServiceProxy = make([]*ProxyConfig, 0)

	rawAuth := configer.GetStringMap("auth")
	for key := range rawAuth {
		authConfig := &AuthConfig{}
		err := mapstructure.Decode(rawAuth[key], authConfig)
		if err != nil {
			panic(err)
		}
		Instance.Auths = append(Instance.Auths, authConfig)
		if authConfig.Type == "youauth" {
			Instance.YouAuthConfig = authConfig
			Instance.YouAuthConfigPrefix = fmt.Sprintf("auth.%s", key)
		}
		if authConfig.Type == "anonymous" {
			Instance.EnableAnonymous = configer.GetBool(fmt.Sprintf("auth.%s.enable", key))
		}
	}

	rawProxy := configer.GetStringMap("proxy")
	for key := range rawProxy {
		proxyConfig := &ProxyConfig{}
		err := mapstructure.Decode(rawProxy[key], proxyConfig)
		if err != nil {
			panic(err)
		}
		Instance.ServiceProxy = append(Instance.ServiceProxy, proxyConfig)
	}
}
