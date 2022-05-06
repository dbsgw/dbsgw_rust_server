package main

import (
	_ "dbsgw_rust_server/initialize"
	"dbsgw_rust_server/models"
	_ "dbsgw_rust_server/routers"
	"encoding/gob"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/beego/beego/v2/server/web/session/redis"
)

func init() {
	gob.Register(models.UserBase{})
	gob.Register(models.UserAuth{})
}
func main() {

	beego.BConfig.WebConfig.Session.SessionProvider = "redis"
	beego.BConfig.WebConfig.Session.SessionProviderConfig = "106.12.125.65:6379,,72812E30873455DCEE2CE2D1EE26E4AB"
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()

}
