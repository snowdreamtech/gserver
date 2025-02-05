package configs

import "github.com/gin-gonic/gin"

// DatabaseConfig Database Config
type DatabaseConfig struct {
	Type     string `mapstructure:"type"`
	DSN      string `mapstructure:"dsn"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Dbname   string `mapstructure:"dbname"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	TimeZone string `mapstructure:"timezone"`
}

var defaultDatabaseConfig = DatabaseConfig{
	Type:     "pg",
	Host:     "localhost",
	Port:     5432,
	Dbname:   "postgres",
	User:     "postgres",
	Password: "postgres",
	TimeZone: "Asia/Shanghai",
	DSN:      "",
}

// GetDatabaseConfigWithContext Get DatabaseConfig from context
func GetDatabaseConfigWithContext(c *gin.Context) (databaseConfig *DatabaseConfig) {
	value, exists := c.Get(ConfigKey)

	if !exists {
		return nil
	}

	confs, ok := value.(*Configs)

	if !ok {
		return nil
	}

	return &confs.Database
}

// GetDatabaseConfig Get DatabaseConfig from context
func GetDatabaseConfig() (databaseConfig *DatabaseConfig) {
	if c == nil {
		return &defaultDatabaseConfig
	}

	return &c.Database
}
