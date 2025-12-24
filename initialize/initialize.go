package initialize

import (
	"crypto/md5"
	"fmt"
	"log"
	"project-manager/model"
	"project-manager/pkg"

	"gorm.io/gorm"
)

// 初始化数据库
func InitDB() error {
	db, err := pkg.OpenDB()
	if err != nil {
		return err
	}
	//自动迁移
	db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Project{},
		&model.Team{},
		&model.TeamProject{},
	)

	//初始化默认权限和角色
	err = initRoleAndPermission(db)
	if err != nil {
		return err
	}
	//打印日志
	log.Println("数据库初始化成功")
	return nil
}

func initRoleAndPermission(db *gorm.DB) error {
	//在这里初始化默认的权限和角色
	//user表中写入一个默认的管理员用户admin 初始密码为adminadmin使用md5加密
	var password = md5.Sum([]byte("adminadmin"))
	//检查用户是否存在
	var count int64
	db.Model(&model.User{}).Where("username = ?", "admin").Count(&count)
	if count > 0 {
		//用户已存在，直接返回
		return nil
	}
	//初始化admin用户
	pkg.Insert("users", &model.User{
		Username: "admin",
		Password: fmt.Sprintf("%x", password),
	})
	//Role表中插入一个admin角色
	pkg.Insert("roles", &model.Role{
		Name: "admin",
		Type: "System",
	})
	//user_roles表中插入关联
	var adminUser model.User
	var adminRole model.Role
	db.Model(&model.User{}).Where("username = ?", "admin").First(&adminUser)
	db.Model(&model.Role{}).Where("name = ?", "admin").First(&adminRole)
	db.Model(&adminUser).Association("Roles").Append(&adminRole)

	//Role表中插入team_lead角色和normal_user角色
	pkg.Insert("roles", &model.Role{
		Name: "normal_user",
		Type: "System",
	})
	pkg.Insert("roles", &model.Role{
		Name: "team_lead",
		Type: "System",
	})
	return nil
}

// 初始化系统
func InitSystem() error {
	return nil
}
