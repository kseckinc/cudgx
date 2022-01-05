package logger

import "go.uber.org/zap"

//Config 是日志库的配置类
type Config struct {
	Level         zap.AtomicLevel        `json:"level" yaml:"level"`
	Encoding      string                 `json:"encoding" yaml:"encoding"`
	InitialFields map[string]interface{} `json:"initialFields" yaml:"initialFields"`
}
