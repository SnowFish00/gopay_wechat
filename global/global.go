package global

import (
	model_cfg "pay/model/config_model"

	"github.com/go-pay/gopay/wechat/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	cfg      model_cfg.Config
	clientV3 *wechat.ClientV3
	logger   *zap.SugaredLogger
	db       *gorm.DB
)

func SetCfg(cfgv model_cfg.Config) {
	cfg = cfgv
}

func SetClient(clientv *wechat.ClientV3) {
	clientV3 = clientv
}

func SetLogger(loggerv *zap.SugaredLogger) {
	logger = loggerv
}

func SetDB(dbv *gorm.DB) {
	db = dbv
}

func ReturnCfg() model_cfg.Config {
	return cfg
}

func ReturnClient() *wechat.ClientV3 {
	return clientV3
}

func ReturnLogger() *zap.SugaredLogger {
	return logger
}

func ReturnDB() *gorm.DB {
	return db
}
