package sub

import (
	"encoding/base64"
	"strings"
	"x-ui/logger"
	"x-ui/web/service"

	"github.com/gin-gonic/gin"
)

type SUBController struct {
	subService     SubService
	settingService service.SettingService
}

func NewSUBController(g *gin.RouterGroup) *SUBController {
	a := &SUBController{}
	a.initRouter(g)
	return a
}

func (a *SUBController) initRouter(g *gin.RouterGroup) {
	g = g.Group("/")

	g.GET("/:subid", a.subs)
}

func (a *SUBController) subs(c *gin.Context) {
	subEncrypt, _ := a.settingService.GetSubEncrypt()
	subShowInfo, _ := a.settingService.GetSubShowInfo()
	subRemoteEnabled, _ := a.settingService.GetSubRemoteEnable()
	subId := c.Param("subid")
	host := strings.Split(c.Request.Host, ":")[0]
	subs, headers, err := a.subService.GetSubs(subId, host, subShowInfo)
	rSubs := ""
	if subRemoteEnabled {
		rSubs = a.subService.getRemoteSubsBySubId(subId)
	}
	if err != nil || len(subs) == 0 {
		c.String(400, "Error!")
	} else {
		result := ""
		for _, sub := range subs {
			result += sub + "\n"
		}

		if rSubs != "" && subRemoteEnabled {
			result += rSubs
		}

		if subEncrypt {
			logger.Debug(result)
			result = base64.StdEncoding.EncodeToString([]byte(result))
		}
		// Add headers
		c.Writer.Header().Set("Subscription-Userinfo", headers[0])
		c.Writer.Header().Set("Profile-Update-Interval", headers[1])
		c.Writer.Header().Set("Profile-Title", headers[2])

		c.String(200, result)
	}
}
