package model_cfg

type Config struct {
	WxClient   WxClient
	HttpServer HttpServer
	Mysql      Mysql
	Log        Log
}

type WxClient struct {
	MchId      string `mapstructure:"mchId"`
	SerialNo   string `mapstructure:"serialNo"`
	ApiV3Key   string `mapstructure:"apiV3Key"`
	PrivateKey string `mapstructure:"privateKey"`
}

type HttpServer struct {
	Origin string `mapstructure:"origin"`
	Port   string `mapstructure:"port"`
}

type Mysql struct {
	Options  string `mapstructure:"options"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Database string `mapstructure:"database"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type Log struct {
	Path string `mapstructure:"path"`
}
