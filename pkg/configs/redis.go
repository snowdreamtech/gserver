package configs

import "github.com/gin-gonic/gin"

// RedisConfig Redis Configg
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
}

var defaultRedisConfig = RedisConfig{
	Host:     "localhost",
	Port:     6379,
	Password: "",
}

// GetRedisConfigWithContext Get RedisConfig from context
func GetRedisConfigWithContext(c *gin.Context) (redisConfig *RedisConfig) {
	value, exists := c.Get(ConfigKey)

	if !exists {
		return nil
	}

	confs, ok := value.(*Configs)

	if !ok {
		return nil
	}

	return &confs.Redis
}

// GetRedisConfig Get RedisConfig from context
func GetRedisConfig() (redisConfig *RedisConfig) {
	if c == nil {
		return &defaultRedisConfig
	}

	return &c.Redis
}
