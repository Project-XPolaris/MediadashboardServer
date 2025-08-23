package plugins

import (
	"github.com/allentom/harukap/plugins/nacos"
	"github.com/allentom/harukap/plugins/youlog"
)

var DefaultYouLogPlugin = &youlog.Plugin{}

// DefaultNacosPlugin 对外暴露的全局 nacos 插件实例，供其他包使用
var DefaultNacosPlugin *nacos.NacosPlugin
