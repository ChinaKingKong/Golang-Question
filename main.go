package main

import (
	"fmt"
	"golang-question/config"
	"golang-question/errorx"
)

const (
	ErrCodeSecretTooShort = 1001
)

type Secret string

// Validate 验证Secret是否有效
// 如果Secret的长度小于8，则返回ErrCodeSecretTooShort错误
// 否则返回nil表示验证通过
func (s Secret) Validate() errorx.Error {
	if len(s) < 8 {
		return errorx.Cf(ErrCodeSecretTooShort, "invalid secret %s", s)
	}
	return nil
}

type Config struct {
	Secret Secret `yaml:"secret" json:"secret"`
}

var conf = config.Local[Config]().Watch().InitData(Config{
	Secret: "hello world",
})

func main() {
	// 获取配置中的Secret
	s := conf.Get().Secret

	// 验证Secret是否有效
	if err := s.Validate(); err != nil {
		// 如果验证失败，打印错误信息
		fmt.Printf("validate error: %+v\n", err)
	}

	// 更新配置中的Secret
	if err := conf.Update(Config{Secret: Secret("updated secret")}); err != nil {
		// 如果更新失败，打印错误信息
		fmt.Printf("update error: %+v\n", err)
	} else {
		// 如果更新成功，打印新的Secret
		fmt.Printf("Secret updated to: %s\n", conf.Get().Secret)
	}
}
