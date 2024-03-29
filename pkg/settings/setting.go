package settings

import (
	"fmt"
	"gopkg.in/ini.v1"
	"log"
	"time"
)

var (
	cfg *ini.File

	RunMode string

	ServerIp     string
	ServerPort   string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	EmailTitle     string
	EmailSender    string
	EmailSPassword string
	EmailSMTPAddr  string
	EmailSMTPPort  int

	MysqlUser     string
	MysqlPassword string
	MysqlIp       string
	MysqlPort     string
	MysqlDatabase string
	MysqlDSN      string

	RedisPassword string
	RedisIp       string
	RedisPort     string
	RedisDB       int
)

func init() {
	var err error
	cfg, err = ini.Load("conf/config.ini")
	if err != nil {
		panic(err)
	}
	LoadBase()
	LoadServer()
	LoadMysql()
	LoadRedis()
	LoadEmail()
}

func LoadBase() {
	RunMode = cfg.Section("").Key("RUN_MODE").String()
}

func LoadEmail() {
	EmailTitle = cfg.Section("email").Key("Title").String()
	EmailSender = cfg.Section("email").Key("Sender").String()
	EmailSPassword = cfg.Section("email").Key("SPassword").String()
	EmailSMTPAddr = cfg.Section("email").Key("SMTPAddr").String()
	EmailSMTPPort, _ = cfg.Section("email").Key("SMTPPort").Int()
}

func LoadMysql() {
	MysqlUser = cfg.Section("mysql").Key("USER").String()
	MysqlPassword = cfg.Section("mysql").Key("PASSWORD").String()
	MysqlIp = cfg.Section("mysql").Key("IP").String()
	MysqlPort = cfg.Section("mysql").Key("PORT").String()
	MysqlDatabase = cfg.Section("mysql").Key("DATABASE").String()

	MysqlDSN = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		MysqlUser, MysqlPassword, MysqlIp, MysqlPort, MysqlDatabase)
}

func LoadRedis() {
	RedisPassword = cfg.Section("redis").Key("PASSWORD").String()
	RedisIp = cfg.Section("redis").Key("IP").String()
	RedisPort = cfg.Section("redis").Key("PORT").String()
	var err error
	RedisDB, err = cfg.Section("redis").Key("DB").Int()
	if err != nil {
		panic("failed to load redis conf!")
	}
}

func LoadServer() {
	sec, err := cfg.GetSection("server")
	if err != nil {
		log.Fatalf("fail to get section 'server': %v", err)
	}
	ServerIp = sec.Key("SERVER_IP").String()
	ServerPort = sec.Key("SERVER_PORT").String()
	ReadTimeout = time.Duration(sec.Key("READ_TIMEOUT").MustInt(60)) * time.Second
	WriteTimeout = time.Duration(sec.Key("WRITE_TIMEOUT").MustInt(60)) * time.Second
}
