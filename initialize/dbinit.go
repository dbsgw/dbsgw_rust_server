package initialize

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func dbinit() {
	runmode, _ := beego.AppConfig.String("runmode")
	isDev := (runmode == "dev")
	registDatabase()
	if isDev {
		orm.Debug = isDev
	}
}

func registDatabase() {
	//初始化数据库
	dbUser, _ := beego.AppConfig.String("mysqluser")
	dbPass, _ := beego.AppConfig.String("mysqlpass")
	dbName, _ := beego.AppConfig.String("mysqldb")
	dbHost, _ := beego.AppConfig.String("mysqlhost")
	dbPort, _ := beego.AppConfig.String("mysqlport")
	maxIdleConn, _ := beego.AppConfig.Int("db_max_idle_conn")
	maxOpenConn, _ := beego.AppConfig.Int("db_max_open_open")
	err := orm.RegisterDriver("mysql", orm.DRMySQL)
	if err != nil {
		logs.Error("数据库链接错误", err)
		return
	}
	err = orm.RegisterDataBase("default", "mysql",
		dbUser+":"+dbPass+"@tcp("+dbHost+":"+dbPort+")/"+dbName+"?charset=utf8&parseTime=true&loc=Asia%2FShanghai",
	)
	if err != nil {
		logs.Error("数据库链接错误", err)
		return
	}
	orm.MaxIdleConnections(maxIdleConn)
	orm.MaxOpenConnections(maxOpenConn)
	orm.DefaultTimeLoc = time.UTC

}
