package configs

import (
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"snowdream.tech/gserver/pkg/env"
	"snowdream.tech/gserver/pkg/os"
	"snowdream.tech/gserver/pkg/tools"
)

// default config names
const (
	// DebugConfigName debug config name
	DebugConfigName = "development"
	// ReleaseConfigName release config name
	ReleaseConfigName = "production"
	// TestConfigName test config name
	TestConfigName = "test"
)

var (
	// Used for flags.
	configFile string

	// Config Types
	configTypes = []string{"json", "env", "ini", "yaml", "toml", "hcl", "properties"}

	// Config Paths
	configPaths = []string{"./", "./configs/", "/etc/" + env.ProjectName, "$HOME/." + env.ProjectName}
)

const (
	// ConfigKey ConfigKey
	ConfigKey = "snowdream.tech/gserver/pkg/configs/configkey"
)

// Configs Configs
type Configs struct {
	Version  string         `mapstructure:"version"`
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
}

var c *Configs = &Configs{
	App:      defaultAppConfig,
	Database: defaultDatabaseConfig,
	Redis:    defaultRedisConfig,
}

// InitConfig init config
func InitConfig() (conf *Configs) {
	if configFile != "" {
		if os.IsExistFile(configFile) {
			// tools.DebugPrintF("[WARNING] %s does not exist or is Not a os.", configFile)
			return nil
		}

		// Use config file from the flag.
		viper.SetConfigFile(configFile)
	} else {
		configName := ""
		switch gin.Mode() {
		case gin.DebugMode:
			configName = DebugConfigName
		case gin.ReleaseMode:
			configName = ReleaseConfigName
		case gin.TestMode:
			configName = TestConfigName
		}

		configFile := ""
		configPath := ""
		configType := ""
		isFinded := false

		for i := 0; i < len(configPaths); i++ {
			configPath = configPaths[i]

			for j := 0; j < len(configTypes); j++ {
				configType = configTypes[j]
				configFile = configPath + configName + "." + configType

				if os.IsExistFile(configFile) {
					isFinded = true
					break
				}
			}

			if isFinded {
				viper.SetConfigFile(configFile)
				break
			}
		}

		if configFile == "" || !os.IsExistFile(configFile) {
			// tools.DebugPrintF("[WARNING] %s.(json/env/ini/yaml/toml/hcl/properties) does not exist or is Not a os.", configName)

			return nil
		}
	}

	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	if err != nil {
		tools.DebugPrintF("[WARNING] Failed to read the config file %s,\n Error:\n %s .", viper.ConfigFileUsed(), err)
		return c
	}

	tools.DebugPrintF("[INFO] The config file %s has been used.", viper.ConfigFileUsed())

	err = viper.Unmarshal(&c)

	if err != nil {
		tools.DebugPrintF("[WARNING] Failed to unmarshal the config file %s,\n Error:\n %s .", viper.ConfigFileUsed(), err)
		return c
	}

	tools.DebugPrintF("[INFO] The config file %s has been Unmarshalled.", viper.ConfigFileUsed())

	viper.WatchConfig()

	viper.OnConfigChange(func(e fsnotify.Event) {
		tools.DebugPrintF("[INFO] The config file %s has been changed.", e.Name)

		if err := viper.Unmarshal(&c); err != nil {
			tools.DebugPrintF("[WARNING] Unmarshal conf failed, err:%s ", err)
		}
	})

	return c
}

// ConfigFile config file path
func ConfigFile() *string {
	return &configFile
}

// GetConfigs Get Configs
func GetConfigs() *Configs {
	return c
}
