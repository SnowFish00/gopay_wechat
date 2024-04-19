package mysql

import (
	"log"
	"pay/global"
	model_srv "pay/model/service_model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func DSN() string {
	sql_cfg := global.ReturnCfg().Mysql
	return sql_cfg.User + ":" + sql_cfg.Password + "@tcp(" + sql_cfg.Host + ":" + sql_cfg.Port + ")/" + sql_cfg.Database + "?" + sql_cfg.Options
}

func Mysql() *gorm.DB {
	mcfg := DSN()
	if db, err := gorm.Open(mysql.Open(mcfg)); err != nil {
		log.Fatal("Connect mysql failed: ", err)
		return nil
	} else {
		//自动建表
		db.AutoMigrate(&model_srv.ChargeMessage{})
		db.AutoMigrate(&model_srv.HttpChargeBlance{})
		db.AutoMigrate(&model_srv.HttpReduceBlance{})
		log.Print("Connect mysql success: ", mcfg)
		return db
	}
}
