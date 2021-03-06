// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"redPacket/controllers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
)

func init() {
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))
	beego.Router("/redpacket", &controllers.RedPacketController{})
	beego.Router("/redpacket/sendpacket", &controllers.RedPacketController{}, "post:SendPacket")
	beego.Router("/redpacket/getlocal", &controllers.RedPacketController{}, "post:GetLocal")
	beego.Router("/redpacket/grad", &controllers.RedPacketController{}, "post:Grad")
	beego.Router("/redpacket/getdetail", &controllers.RedPacketController{}, "post:GetDetail")
	beego.Router("/redpacket/test", &controllers.RedPacketController{}, "get:Test")
	beego.Router("/redpacket/autograd", &controllers.RedPacketController{}, "post:AutoGrad")
}
