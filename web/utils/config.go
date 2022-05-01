package utils

import (
	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
)

var (
	G_server_name  string // 项目名称
	G_server_addr  string // 服务器ip地址
	G_server_port  string // 服务器端口
	G_mysql_addr   string // mysql地址
	G_mysql_port   string // mysql端口
	G_mysql_dbname string // mysql数据库名
	G_mysql_user   string // mysql用户名
	G_mysql_pass   string // mysql密码
	G_redis_addr   string // redis地址
	G_redis_port   string // redis端口
	G_redis_dbnum  string // redis数据库编号
)

func InitConfig() {
	if appconf, err := config.NewConfig("ini", "../conf/app.conf"); err != nil {
		logs.Debug(err)
		return
	} else {
		G_server_name, _ = appconf.String("appname")
		G_server_addr, _ = appconf.String("httpaddr")
		G_server_port, _ = appconf.String("httpport")
		G_redis_addr, _ = appconf.String("redisaddr")
		G_redis_port, _ = appconf.String("redisport")
		G_redis_dbnum, _ = appconf.String("redisdbnum")
		G_mysql_addr, _ = appconf.String("mysqladdr")
		G_mysql_port, _ = appconf.String("mysqlport")
		G_mysql_dbname, _ = appconf.String("mysqldbname")
		G_mysql_user, _ = appconf.String("mysqluser")
		G_mysql_pass, _ = appconf.String("mysqlpass")
		return
	}
}

func init() {
	InitConfig()
}
