package utils

import (
	"log"

	"gitlab.com/lilh/go-demo/internal/config"
	"gitlab.mvalley.com/datapack/cain/pkg/client"
	"gorm.io/gorm"
)

var BatchSize = 20000
var RecruitmentDatapackDB *gorm.DB

func InitRecruitmentDatapackDB() {
	if RecruitmentDatapackDB != nil {
		return
	}
	db, err := client.InitGormV2(config.Config.MysqlConfig.RecruitmentDatapackMySQLConfig)
	if err != nil {
		panic(err)
	}
	log.Println("init RecruitmentDatapackDB success")
	RecruitmentDatapackDB = db
}
