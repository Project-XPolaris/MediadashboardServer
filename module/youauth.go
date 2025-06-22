package module

import (
	"bytes"
	"encoding/gob"
	"mediadashboard/config"
	"mediadashboard/database"

	"github.com/allentom/harukap"
	"github.com/allentom/harukap/commons"
	"github.com/allentom/harukap/module/auth"
)

// Auth 是认证模块的实例，包含认证插件列表
var Auth = &auth.AuthModule{
	Plugins: []harukap.AuthPlugin{},
}

// CreateAuthModule 初始化认证模块
// 设置配置提供者并添加用户序列化器
func CreateAuthModule() {
	Auth.ConfigProvider = config.DefaultConfigProvider
	Auth.AddCacheStore(&UserSerializer{})
	Auth.InitModule()
}

// UserSerializer 实现了用户数据的序列化和反序列化
type UserSerializer struct {
}

// Serialize 将用户数据序列化为字节数组
// 使用 gob 编码器将用户对象编码为二进制格式
func (s *UserSerializer) Serialize(data interface{}) ([]byte, error) {
	user := data.(*database.User)
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(user)
	if err != nil {
		return nil, err
	}
	return result.Bytes(), nil
}

// Deserialize 将字节数组反序列化为用户对象
// 使用 gob 解码器将二进制数据转换回用户对象
func (s *UserSerializer) Deserialize(raw []byte) (commons.AuthUser, error) {
	var user database.User
	decoder := gob.NewDecoder(bytes.NewReader(raw))
	err := decoder.Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
