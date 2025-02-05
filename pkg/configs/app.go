package configs

import "github.com/gin-gonic/gin"

// AppConfig App Config
type AppConfig struct {
	Host                string   `mapstructure:"host"`
	Port                string   `mapstructure:"port"`
	Basic               bool     `mapstructure:"basic"`
	Gzip                bool     `mapstructure:"gzip"`
	User                string   `mapstructure:"user"`
	LogDir              string   `mapstructure:"logdir"`
	RateLimiter         string   `mapstructure:"ratelimiter"`
	ReadTimeout         int64    `mapstructure:"readtimeout"`
	WriteTimeout        int64    `mapstructure:"writetimeout"`
	WwwRoot             string   `mapstructure:"wwwroot"`
	AutoIndexTimeFormat string   `mapstructure:"autoindextimeformat"`
	AutoIndexExactSize  bool     `mapstructure:"autoindexexactsize"`
	PreviewHTML         bool     `mapstructure:"previewhtml"`
	EnableHTTPS         bool     `mapstructure:"enablehttps"`
	HTTPSPort           string   `mapstructure:"httpsport"`
	HTTPSCertFile       string   `mapstructure:"httpscertfile"`
	HTTPSKeyFile        string   `mapstructure:"httpskeyfile"`
	HTTPSCertsDir       string   `mapstructure:"httpscertsdir"`
	HTTPSDomains        []string `mapstructure:"httpsdomains"`
	ContactEmail        string   `mapstructure:"contactemail"`
	SpeedLimiter        int64    `mapstructure:"speedlimiter"`
	RefererLimiter      bool     `mapstructure:"refererlimiter"`
}

var defaultAppConfig = AppConfig{
	Host:                "",
	Port:                "",
	Basic:               false,
	Gzip:                true,
	User:                "admin:admin",
	LogDir:              ".",
	RateLimiter:         "",
	ReadTimeout:         10,
	WriteTimeout:        10,
	WwwRoot:             "",
	AutoIndexTimeFormat: "2006-01-02 15:04:05",
	AutoIndexExactSize:  false,
	PreviewHTML:         true,
	EnableHTTPS:         false,
	HTTPSPort:           "",
	HTTPSCertFile:       "",
	HTTPSKeyFile:        "",
	HTTPSCertsDir:       "certs",
	HTTPSDomains:        nil,
	ContactEmail:        "",
	SpeedLimiter:        0,
	RefererLimiter:      false,
}

// GetAppConfigWithContext Get AppConfig from context
func GetAppConfigWithContext(c *gin.Context) (appconfig *AppConfig) {
	value, exists := c.Get(ConfigKey)

	if !exists {
		return nil
	}

	confs, ok := value.(*Configs)

	if !ok {
		return nil
	}

	return &confs.App
}

// GetAppConfig Get AppConfig from context
func GetAppConfig() (appconfig *AppConfig) {
	if c == nil {
		return &defaultAppConfig
	}

	return &c.App
}
