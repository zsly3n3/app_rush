package routes

import (
	"app/datastruct"
	"app/event"
	"app/log"
	"app/tools"

	"github.com/gin-gonic/gin"
	//"app/log"
)

func DepositSucceed(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/depositSucceed/:userid/:money", func(c *gin.Context) {
		userid := tools.StringToInt(c.Param("userid"))
		money := tools.StringToInt64(c.Param("money"))
		data, code := eventHandler.DepositSucceed(userid, money)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func getAuthAddr(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/auth/:userid", func(c *gin.Context) {
		userid := c.Param("userid")
		data := eventHandler.GetAuthAddr(userid)
		c.JSON(200, gin.H{
			"code": datastruct.NULLError,
			"data": data,
		})
	})
}

func getDownLoadAppAddr(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/downloadapp", func(c *gin.Context) {
		data, code := eventHandler.GetDownLoadAppAddr()
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}
func getDirectDownloadApp(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/directdownloadapp", func(c *gin.Context) {
		data := eventHandler.GetDirectDownloadApp()
		c.JSON(200, gin.H{
			"code": datastruct.NULLError,
			"data": data,
		})
	})
}
func getAppDownLoadShareUrl(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/downloadappshareurl", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		userId, _, tf := checkToken(c, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.GetAppDownLoadShareUrl(userId)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func getAppAddr(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/addr", func(c *gin.Context) {
		data, code := eventHandler.GetAppAddr()
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func getKfInfo(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/kfinfo", func(c *gin.Context) {
		data, code := eventHandler.GetKfInfo()
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func checkToken(c *gin.Context, eventHandler *event.EventHandler) (int, string, bool) {
	tokens, isExist := c.Request.Header["Apptoken"]
	tf := false
	var token string
	var userId int
	var isBlackList bool
	if isExist {
		token = tokens[0]
		if token != "" {
			userId, tf, isBlackList = eventHandler.IsExistUser(token)
			if tf && isBlackList {
				url := eventHandler.GetBlackListRedirect()
				c.JSON(200, gin.H{
					"code": datastruct.AppRedirect,
					"data": url,
				})
				return -1, "", false
			}
		}
	}
	if !tf {
		c.JSON(200, gin.H{
			"code": datastruct.TokenError,
		})
	}
	return userId, token, tf
}

func checkVersion(c *gin.Context, eventHandler *event.EventHandler) bool {
	serverVersion, isMaintain := eventHandler.GetServerInfoFromMemory()
	if isMaintain == 1 {
		c.JSON(int(datastruct.Maintenance), gin.H{
			"code": datastruct.NULLError,
		})
		return false
	}
	// isPcInfo, tf := c.Request.Header["Ispc"]
	// if tf && isPcInfo[0] == "1" {
	// 	c.JSON(200, gin.H{
	// 		"code": datastruct.AppRedirect,
	// 		"data": eventHandler.GetPCRedirect(),
	// 	})
	// 	return false
	// }
	version, isExist := c.Request.Header["Appversion"]
	if isExist && version[0] == serverVersion {
		return true
	} else {
		c.JSON(200, gin.H{
			"code": datastruct.VersionError,
			"data": eventHandler.GetDirectDownloadApp(),
		})
		return false
	}
}

func appLogin(r *gin.Engine, eventHandler *event.EventHandler) {
	r.POST("/app/user/login", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		eventHandler.AppLogin(c)
	})
}

func h5Login(r *gin.Engine, eventHandler *event.EventHandler) {
	r.POST("/app/user/loginh5", func(c *gin.Context) {
		// if !checkVersion(c, eventHandler) {
		// 	return
		// }
		eventHandler.H5Login(c)
	})
}

func getHomeData(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/homedata", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		platforms, tf := c.Request.Header["Platform"]
		if !tf {
			c.JSON(200, gin.H{
				"code": datastruct.HeaderParamError,
			})
			return
		}
		if !checkPlatform(platforms) {
			c.JSON(200, gin.H{
				"code": datastruct.HeaderParamError,
			})
			return
		}
		platform := tools.StringToInt(platforms[0])
		c.JSON(200, gin.H{
			"code": datastruct.NULLError,
			"data": eventHandler.GetHomeData(platform),
		})
	})
}

func checkPlatform(platforms []string) bool {
	tf := true
	tmp := tools.StringToInt(platforms[0])
	if tmp < int(datastruct.APP) || tmp > int(datastruct.H5) {
		tf = false
	}
	return tf
}

func getGoods(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/goods/:pageindex/:pagesize/:classid", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		userId, _, tf := checkToken(c, eventHandler)
		if !tf {
			return
		}
		pageIndex := tools.StringToInt(c.Param("pageindex"))
		pageSize := tools.StringToInt(c.Param("pagesize"))
		classid := tools.StringToInt(c.Param("classid"))
		if pageIndex <= 0 || pageSize <= 0 || classid <= 0 {
			c.JSON(200, gin.H{
				"code": datastruct.ParamError,
			})
			return
		}
		c.JSON(200, gin.H{
			"code": datastruct.NULLError,
			"data": eventHandler.GetGoods(pageIndex, pageSize, classid, userId),
		})
	})
}

func getGoodsClass(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/goodsclass", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		userId, _, tf := checkToken(c, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.GetGoodsClass(userId)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func worldRewardInfo(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/rewardinfo/:pageindex/:pagesize", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		pageIndex := tools.StringToInt(c.Param("pageindex"))
		pageSize := tools.StringToInt(c.Param("pagesize"))
		if pageIndex <= 0 || pageSize <= 0 {
			c.JSON(200, gin.H{
				"code": datastruct.ParamError,
			})
			return
		}
		data, code := eventHandler.WorldRewardInfo(pageIndex, pageSize)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func freeRougeGame(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/freerougegame", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		data, code := eventHandler.FreeRougeGame()
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func payRougeGame(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/user/payrougegame/:goodsid", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		userId, _, tf := checkToken(c, eventHandler)
		if !tf {
			return
		}
		goodsid := c.Param("goodsid")
		eventHandler.PayRougeGame(userId, tools.StringToInt(goodsid), c)
	})
}

func levelPass(r *gin.Engine, eventHandler *event.EventHandler) {
	r.POST("/app/user/levelpass", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		userId, _, tf := checkToken(c, eventHandler)
		if !tf {
			return
		}
		eventHandler.LevelPass(userId, c)
	})
}

func levelPassFailed(r *gin.Engine, eventHandler *event.EventHandler) {
	r.POST("/app/user/levelpassfailed", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		userId, _, tf := checkToken(c, eventHandler)
		if !tf {
			return
		}
		eventHandler.LevelPassFailed(userId, c)
	})
}

func levelPassSucceed(r *gin.Engine, eventHandler *event.EventHandler) {
	r.POST("/app/user/levelpasssucceed", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		userId, _, tf := checkToken(c, eventHandler)
		if !tf {
			return
		}
		eventHandler.LevelPassSucceed(userId, c)
	})
}

func getNotAppliedOrderInfo(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/user/notappliedorder/:pageindex/:pagesize", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		pageIndex := tools.StringToInt(c.Param("pageindex"))
		pageSize := tools.StringToInt(c.Param("pagesize"))
		if pageIndex <= 0 || pageSize <= 0 {
			c.JSON(200, gin.H{
				"code": datastruct.ParamError,
			})
			return
		}
		userId, _, tf := checkToken(c, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.GetNotAppliedOrderInfo(userId, pageIndex, pageSize)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func getNotSendGoods(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/user/notsendgoods/:pageindex/:pagesize", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		pageIndex := tools.StringToInt(c.Param("pageindex"))
		pageSize := tools.StringToInt(c.Param("pagesize"))
		if pageIndex <= 0 || pageSize <= 0 {
			c.JSON(200, gin.H{
				"code": datastruct.ParamError,
			})
			return
		}
		userId, _, tf := checkToken(c, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.GetNotSendGoods(userId, pageIndex, pageSize)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func getHasSendedGoods(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/user/hassendedgoods/:pageindex/:pagesize", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		pageIndex := tools.StringToInt(c.Param("pageindex"))
		pageSize := tools.StringToInt(c.Param("pagesize"))
		if pageIndex <= 0 || pageSize <= 0 {
			c.JSON(200, gin.H{
				"code": datastruct.ParamError,
			})
			return
		}
		userId, _, tf := checkToken(c, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.GetHasSendedGoods(userId, pageIndex, pageSize)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func getAppraiseOrder(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/user/appraiseorder/:pageindex/:pagesize", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		pageIndex := tools.StringToInt(c.Param("pageindex"))
		pageSize := tools.StringToInt(c.Param("pagesize"))
		if pageIndex <= 0 || pageSize <= 0 {
			c.JSON(200, gin.H{
				"code": datastruct.ParamError,
			})
			return
		}
		userId, _, tf := checkToken(c, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.GetAppraiseOrder(userId, pageIndex, pageSize)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func applySend(r *gin.Engine, eventHandler *event.EventHandler) {
	r.POST("/app/user/applysend", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		userId, _, tf := checkToken(c, eventHandler)
		if !tf {
			return
		}
		eventHandler.ApplySend(userId, c)
	})
}

func commissionRank(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/user/commissionrank/:pageindex/:pagesize", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		pageIndex := tools.StringToInt(c.Param("pageindex"))
		pageSize := tools.StringToInt(c.Param("pagesize"))
		if pageIndex <= 0 || pageSize <= 0 {
			c.JSON(200, gin.H{
				"code": datastruct.ParamError,
			})
			return
		}
		userId, _, tf := checkToken(c, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.CommissionRank(userId, pageIndex, pageSize)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func commissionInfo(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/user/commissioninfo/:pageindex/:pagesize", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		pageIndex := tools.StringToInt(c.Param("pageindex"))
		pageSize := tools.StringToInt(c.Param("pagesize"))
		if pageIndex <= 0 || pageSize <= 0 {
			c.JSON(200, gin.H{
				"code": datastruct.ParamError,
			})
			return
		}
		userId, _, tf := checkToken(c, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.CommissionInfo(userId, pageIndex, pageSize)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func getAgentlevelN(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/user/agentleveln/:level/:pageindex/:pagesize", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		level := tools.StringToInt(c.Param("level"))
		pageIndex := tools.StringToInt(c.Param("pageindex"))
		pageSize := tools.StringToInt(c.Param("pagesize"))
		if level <= 0 || level > 3 || pageIndex <= 0 || pageSize <= 0 {
			c.JSON(200, gin.H{
				"code": datastruct.ParamError,
			})
			return
		}
		userId, _, tf := checkToken(c, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.GetAgentlevelN(userId, level, pageIndex, pageSize)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func getDrawCash(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/user/drawcash/:pageindex/:pagesize", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		pageIndex := tools.StringToInt(c.Param("pageindex"))
		pageSize := tools.StringToInt(c.Param("pagesize"))
		if pageIndex <= 0 || pageSize <= 0 {
			c.JSON(200, gin.H{
				"code": datastruct.ParamError,
			})
			return
		}
		userId, _, tf := checkToken(c, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.GetDrawCash(userId, pageIndex, pageSize)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func getDrawCashRule(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/user/drawcashrule", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		userId, _, tf := checkToken(c, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.GetDrawCashRule(userId)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func getDepositParams(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/user/depositparams", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		userId, _, tf := checkToken(c, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.GetDepositParams(userId)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func getInviteUrl(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/user/inviteurl", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		userId, _, tf := checkToken(c, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.GetInviteUrl(userId)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func getInviteUrlNoToken(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/user/inviteurlnotoken", func(c *gin.Context) {
		data, code := eventHandler.GetInviteUrlNoToken()
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func getGoldInfo(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/user/goldchange/:pageindex/:pagesize", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		userId, _, tf := checkToken(c, eventHandler)
		if !tf {
			return
		}
		pageIndex := tools.StringToInt(c.Param("pageindex"))
		pageSize := tools.StringToInt(c.Param("pagesize"))
		if pageIndex <= 0 || pageSize <= 0 {
			c.JSON(200, gin.H{
				"code": datastruct.ParamError,
			})
			return
		}
		data, code := eventHandler.GetGoldInfo(userId, pageIndex, pageSize)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func getUserInfo(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/user/info", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}

		platforms, tf := c.Request.Header["Platform"]
		if !tf {
			c.JSON(200, gin.H{
				"code": datastruct.HeaderParamError,
			})
			return
		}
		if !checkPlatform(platforms) {
			c.JSON(200, gin.H{
				"code": datastruct.HeaderParamError,
			})
			return
		}
		platform := tools.StringToInt(platforms[0])

		userId, _, tf := checkToken(c, eventHandler)
		if !tf {
			return
		}

		data, code := eventHandler.GetUserInfo(userId, datastruct.Platform(platform))
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func h5purchase(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/user/h5purchase/:class/:id", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		ip_addr, isExist := c.Request.Header["Remote_addr"]
		if !isExist || ip_addr[0] == "" {
			c.JSON(200, gin.H{
				"code": datastruct.GetIpAdrrError,
			})
			return
		}
		userId, _, tf := checkToken(c, eventHandler)
		if !tf {
			return
		}
		class := c.Param("class")
		purchase_id := tools.StringToInt(c.Param("id"))
		if !checkPurchaseClass(class) || purchase_id <= 0 {
			c.JSON(200, gin.H{
				"code": datastruct.ParamError,
			})
			return
		}
		data, code := eventHandler.Purchase(userId, datastruct.PurchaseType(class), purchase_id, ip_addr[0], datastruct.H5)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func checkPurchaseClass(class string) bool {
	tf := true
	if class != string(datastruct.GoodsType) && class != string(datastruct.GoldType) && class != string(datastruct.VipType) {
		tf = false
	}
	return tf
}

func appPurchase(r *gin.Engine, eventHandler *event.EventHandler) {
	r.POST("/app/user/purchase", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		ip_addr, isExist := c.Request.Header["Remote_addr"]
		if !isExist || ip_addr[0] == "" {
			c.JSON(200, gin.H{
				"code": datastruct.GetIpAdrrError,
			})
			return
		}
		userId, _, tf := checkToken(c, eventHandler)
		if !tf {
			return
		}
		var body datastruct.AppPurchaseBody
		err := c.BindJSON(&body)
		purchase_id := tools.StringToInt(body.Id)
		if err != nil || purchase_id <= 0 || !checkPurchaseClass(body.Class) {
			if err != nil {
				log.Debug("err:%v", err.Error())
			}
			c.JSON(200, gin.H{
				"code": datastruct.ParamError,
			})
			return
		}
		data, code := eventHandler.Purchase(userId, datastruct.PurchaseType(body.Class), purchase_id, ip_addr[0], datastruct.APP)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func wxPayResultCall(r *gin.Engine, eventHandler *event.EventHandler) {
	r.POST(datastruct.WXPayCallRoute, func(c *gin.Context) {
		eventHandler.WxPayResultCall(c)
	})
}

func userPayee(r *gin.Engine, eventHandler *event.EventHandler) {
	r.POST("/app/user/payee", func(c *gin.Context) {
		ip_addr, isExist := c.Request.Header["Remote_addr"]
		if !isExist || ip_addr[0] == "" {
			c.JSON(200, gin.H{
				"code": datastruct.GetIpAdrrError,
			})
			return
		}
		if !checkVersion(c, eventHandler) {
			return
		}
		userId, _, tf := checkToken(c, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.UserPayee(userId, ip_addr[0], c)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}
func userShare(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/user/share", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		userId, _, tf := checkToken(c, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.CustomShareForApp(userId)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func customShareForGZH(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/user/customshareforgzh", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		userId, _, tf := checkToken(c, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.CustomShareForGZH(userId)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func getEntryPageUrl(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/entrypageurl", func(c *gin.Context) {
		data := eventHandler.GetEntryPageUrl()
		c.JSON(200, gin.H{
			"code": datastruct.NULLError,
			"data": data,
		})
	})
}
func getCheckInData(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/user/checkindata", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		userId, _, tf := checkToken(c, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.GetCheckInData(userId)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func checkIn(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/user/checkin", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		userId, _, tf := checkToken(c, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.CheckIn(userId)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func appGetMemberList(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/user/memberlist", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		userId, _, tf := checkToken(c, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.AppGetMemberList(userId)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func startLottery(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/user/startlottery/:rushprice", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		rushprice := tools.StringToInt64(c.Param("rushprice"))
		if rushprice <= 0 {
			c.JSON(200, gin.H{
				"code": datastruct.ParamError,
			})
			return
		}
		userId, _, tf := checkToken(c, eventHandler)
		if !tf {
			return
		}
		eventHandler.StartLottery(userId, rushprice, c)
	})
}

func getLotteryGoodsSucceedHistory(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/getlotteryhistory", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		data, code := eventHandler.GetLotteryGoodsSucceedHistory()
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func getLotteryGoodsOrderInfo(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/user/getlotteryorder/:pageindex/:pagesize", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		pageIndex := tools.StringToInt(c.Param("pageindex"))
		pageSize := tools.StringToInt(c.Param("pagesize"))
		if pageIndex <= 0 || pageSize <= 0 {
			c.JSON(200, gin.H{
				"code": datastruct.ParamError,
			})
			return
		}
		userId, _, tf := checkToken(c, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.GetLotteryGoodsOrderInfo(userId, pageIndex, pageSize)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func getSharePosters(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/user/shareposters", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		userId, _, tf := checkToken(c, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.GetSharePosters(userId)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func getUserAppraiseForApp(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/user/appraise/:pageindex/:pagesize", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		pageIndex := tools.StringToInt(c.Param("pageindex"))
		pageSize := tools.StringToInt(c.Param("pagesize"))
		if pageIndex <= 0 || pageSize <= 0 {
			c.JSON(200, gin.H{
				"code": datastruct.ParamError,
			})
			return
		}
		data, code := eventHandler.GetUserAppraiseForApp(pageIndex, pageSize)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

// func userActivate(r *gin.Engine, eventHandler *event.EventHandler) {
// 	r.GET("/app/user/activate", func(c *gin.Context) {
// 		if !checkVersion(c, eventHandler) {
// 			return
// 		}
// 		userId, _, tf := checkToken(c, eventHandler)
// 		if !tf {
// 			return
// 		}
// 		c.JSON(200, gin.H{
// 			"code": eventHandler.UserActivate(userId),
// 		})
// 	})
// }

func userAppraise(r *gin.Engine, eventHandler *event.EventHandler) {
	r.POST("/app/user/appraise", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		userId, _, tf := checkToken(c, eventHandler)
		if !tf {
			return
		}
		c.JSON(200, gin.H{
			"code": eventHandler.UserAppraise(userId, c),
		})
	})
}

func getGoodsDetailForApp(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/goodsDetail/:goodsid", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		goodsid := tools.StringToInt(c.Param("goodsid"))
		if goodsid <= 0 {
			c.JSON(200, gin.H{
				"code": datastruct.ParamError,
			})
			return
		}
		data, code := eventHandler.GetGoodsDetailForApp(goodsid)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func updateUserAddress(r *gin.Engine, eventHandler *event.EventHandler) {
	r.POST("/app/user/address", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		userId, _, tf := checkToken(c, eventHandler)
		if !tf {
			return
		}
		c.JSON(200, gin.H{
			"code": eventHandler.UpdateUserAddress(userId, c),
		})
	})
}
func getUserAddress(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/user/address", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		userId, _, tf := checkToken(c, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.GetUserAddress(userId)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func getAgencyPage(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/user/agencypage", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		userId, _, tf := checkToken(c, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.GetAgencyPage(userId)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func remindSendGoods(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/user/remindsendgoods/:number", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		userId, _, tf := checkToken(c, eventHandler)
		if !tf {
			return
		}
		number := c.Param("number")
		if number == "" {
			c.JSON(200, gin.H{
				"code": datastruct.ParamError,
			})
			return
		}
		c.JSON(200, gin.H{
			"code": eventHandler.RemindSendGoods(userId, number),
		})
	})
}

func addSuggest(r *gin.Engine, eventHandler *event.EventHandler) {
	r.POST("/app/user/suggest", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		userId, _, tf := checkToken(c, eventHandler)
		if !tf {
			return
		}
		c.JSON(200, gin.H{
			"code": eventHandler.AddSuggest(userId, c),
		})
	})
}

func addComplaint(r *gin.Engine, eventHandler *event.EventHandler) {
	r.POST("/app/user/complaint", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		userId, _, tf := checkToken(c, eventHandler)
		if !tf {
			return
		}
		c.JSON(200, gin.H{
			"code": eventHandler.AddComplaint(userId, c),
		})
	})
}

func getDownLoadAppGift(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/user/downloadappgift", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		userId, _, tf := checkToken(c, eventHandler)
		if !tf {
			return
		}
		c.JSON(200, gin.H{
			"code": eventHandler.GetDownLoadAppGift(userId),
		})
	})
}

func getRegisterGift(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/user/registergift", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		userId, _, tf := checkToken(c, eventHandler)
		if !tf {
			return
		}
		c.JSON(200, gin.H{
			"code": eventHandler.GetRegisterGift(userId),
		})
	})
}

func appGetDefaultAgency(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/defaultagency", func(c *gin.Context) {
		data, code := eventHandler.GetDefaultAgency()
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func IsRefreshHomeGoodsData(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/isrefreshgoods/:classid", func(c *gin.Context) {
		if !checkVersion(c, eventHandler) {
			return
		}
		userId, _, tf := checkToken(c, eventHandler)
		if !tf {
			return
		}
		classid := tools.StringToInt(c.Param("classid"))
		if classid <= 0 {
			c.JSON(200, gin.H{
				"code": datastruct.ParamError,
			})
			return
		}
		data, code := eventHandler.IsRefreshHomeGoodsData(userId, classid)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func getGoldFromPoster(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/app/goldposter/:userid/:gpid", func(c *gin.Context) {
		userid := tools.StringToInt(c.Param("userid"))
		gpid := tools.StringToInt(c.Param("gpid"))
		if userid <= 0 || gpid <= 0 {
			c.JSON(200, gin.H{
				"code": datastruct.ParamError,
			})
			return
		}
		data, code := eventHandler.GetGoldFromPoster(userid, gpid)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func Register(r *gin.Engine, eventHandler *event.EventHandler) {
	getAuthAddr(r, eventHandler)
	getKfInfo(r, eventHandler)
	getHomeData(r, eventHandler)
	getGoods(r, eventHandler)
	getGoodsClass(r, eventHandler)
	worldRewardInfo(r, eventHandler)
	freeRougeGame(r, eventHandler)
	payRougeGame(r, eventHandler)
	levelPass(r, eventHandler)
	levelPassFailed(r, eventHandler)
	levelPassSucceed(r, eventHandler)
	getNotAppliedOrderInfo(r, eventHandler)
	getNotSendGoods(r, eventHandler)
	applySend(r, eventHandler)
	commissionRank(r, eventHandler)
	commissionInfo(r, eventHandler)
	getAgentlevelN(r, eventHandler)
	getDrawCash(r, eventHandler)
	getDrawCashRule(r, eventHandler)
	getDepositParams(r, eventHandler)
	getInviteUrl(r, eventHandler)
	getUserInfo(r, eventHandler)
	appLogin(r, eventHandler)
	h5Login(r, eventHandler)
	getGoldInfo(r, eventHandler)
	h5purchase(r, eventHandler)
	appPurchase(r, eventHandler)
	userPayee(r, eventHandler)
	wxPayResultCall(r, eventHandler)
	userShare(r, eventHandler)
	getAppAddr(r, eventHandler)
	getDownLoadAppAddr(r, eventHandler)
	getDirectDownloadApp(r, eventHandler)
	DepositSucceed(r, eventHandler)
	getAppDownLoadShareUrl(r, eventHandler)
	getInviteUrlNoToken(r, eventHandler)
	customShareForGZH(r, eventHandler)
	getEntryPageUrl(r, eventHandler)
	getCheckInData(r, eventHandler)
	checkIn(r, eventHandler)
	appGetMemberList(r, eventHandler)
	startLottery(r, eventHandler)
	getLotteryGoodsSucceedHistory(r, eventHandler)
	getLotteryGoodsOrderInfo(r, eventHandler)
	getSharePosters(r, eventHandler)
	getUserAppraiseForApp(r, eventHandler)
	userAppraise(r, eventHandler)
	getGoodsDetailForApp(r, eventHandler)
	updateUserAddress(r, eventHandler)
	getUserAddress(r, eventHandler)
	getAgencyPage(r, eventHandler)
	getHasSendedGoods(r, eventHandler)
	getAppraiseOrder(r, eventHandler)
	remindSendGoods(r, eventHandler)
	addSuggest(r, eventHandler)
	addComplaint(r, eventHandler)
	getDownLoadAppGift(r, eventHandler)
	getRegisterGift(r, eventHandler)
	appGetDefaultAgency(r, eventHandler)
	IsRefreshHomeGoodsData(r, eventHandler)
	getGoldFromPoster(r, eventHandler) //金币海报领取
}
