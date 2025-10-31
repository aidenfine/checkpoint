package checkpoint

import (
	"os"
	"sync"
)

type Config struct {
	ServiceUrl string
	Logs       LogConfig
}

var (
	cfg  *Config
	once sync.Once
)

type LogConfig struct {
	Level string
}

func LoadConfigs() {
	once.Do(func() {
		cfg = &Config{
			ServiceUrl: os.Getenv("SERVICE_URL"),
			Logs: LogConfig{
				Level: getEnvWithDefault("LOG_LEVEL", "ERROR"),
			},
		}
	})
}
func getEnvWithDefault(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val

}

func Get() *Config {
	if cfg == nil {
		panic("Config not loaded, call config.Load()") // may not need to panic here
	}
	return cfg
}
