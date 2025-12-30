package setting

import (
	"log"
	"os"
	"time"

	"github.com/go-ini/ini"
)

var (
	Cfg *ini.File

	RunMode      string
	HTTPPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	JWTSecret    string
)

func init() {
	var err error
	Cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'conf/app.ini': %v", err)
	}
	LoadBase()
	LoadServer()
	LoadApp()
}

func LoadBase() {
	// 优先使用环境变量 RUN_MODE（便于容器或运行时覆盖），否则从配置文件读取
	if m := os.Getenv("RUN_MODE"); m != "" {
		RunMode = m
		return
	}
	sec, err := Cfg.GetSection("server")
	if err != nil {
		RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")
		return
	}
	RunMode = sec.Key("RUN_MODE").MustString("debug")
}

func LoadServer() {
	sec, err := Cfg.GetSection("server")
	if err != nil {
		log.Fatalf("Fail to get section 'server'")
	}
	HTTPPort = sec.Key("HTTP_PORT").MustInt(8080)
	ReadTimeout = time.Duration(sec.Key("READ_TIMEOUT").MustInt(60)) * time.Second
	WriteTimeout = time.Duration(sec.Key("WRITE_TIMEOUT").MustInt(60)) * time.Second
}

func LoadApp() {
	sec, err := Cfg.GetSection("app")
	if err != nil {
		log.Fatalf("Fail to get section 'app': %v", err)
	}
	JWTSecret = sec.Key("JWT_SECRET").MustString("!@)*#)!@U#@*!@!)")
}
