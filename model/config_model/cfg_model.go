package model_cfg

type Config struct {
	WxClient   WxClient
	HttpServer HttpServer
	Mysql      Mysql
	Log        Log
	SynUrl     SynUrl
	Test       Test
}

type WxClient struct {
	MchId      string `mapstructure:"mchId"`
	AppID      string `mapstructure:"appId"`
	AppSecret  string `mapstructure:"appSecret"`
	SerialNo   string `mapstructure:"serialNo"`
	ApiV3Key   string `mapstructure:"apiV3Key"`
	PrivateKey string `mapstructure:"privateKey"`
	NotifyUrl  string `mapstructure:"notifyUrl"`
}

type HttpServer struct {
	Origin     string `mapstructure:"origin"`
	Port       string `mapstructure:"port"`
	AdminToken string `mapstructure:"adminToken"`
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

type SynUrl struct {
	AddURL    string `mapstructure:"add_url"`
	ReduceURL string `mapstructure:"reduce_url"`
}

type Test struct {
	OpenID string `mapstructure:"wx_openid"`
}
