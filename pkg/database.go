package pkg

import (
	"fmt"
	"log"
	"project-manager/pkg/setting"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB 全局数据库连接
var DB *gorm.DB

// 数据库连接
func OpenDB() (*gorm.DB, error) {
	var (
		err                          error
		dbName, user, password, host string
	)
	sec, err := setting.Cfg.GetSection("database")
	if err != nil {
		log.Fatal(2, "Fail to get section 'database': %v", err)
	}
	dbName = sec.Key("NAME").String()
	user = sec.Key("USER").String()
	password = sec.Key("PASSWORD").String()
	host = sec.Key("HOST").String()

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		user,
		password,
		host,
		dbName)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//设置超时时间
	if err != nil {
		log.Println(err)
	}

	return DB, nil
}

// 关闭数据库连接
func CloseDB() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Println(err)
		return
	}
	sqlDB.Close()
}

func Insert(tableName string, value interface{}) *gorm.DB {
	return DB.Table(tableName).Create(value)
}

func Query(tableName string, query string) *gorm.DB {
	return DB.Table(tableName).Where(query)
}

func Delete(tableName string, query string) *gorm.DB {
	return DB.Table(tableName).Where(query).Delete(nil)
}
