package routes

import (
	"app/datastruct"
	"app/event"
	"app/tools"

	"github.com/gin-gonic/gin"
)

func editDomain(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/domain"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		c.JSON(200, gin.H{
			"code": eventHandler.EditDomain(c),
		})
	})
}

func updateSendInfo(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/updatesendinfo"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		eventHandler.UpdateSendInfo(c)
	})
}

func updateDefaultAgency(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/defaultagency"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		eventHandler.UpdateDefaultAgency(c)
	})
}

func editMemberLevel(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/editmemberlevel"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		eventHandler.EditMemberLevel(c)
	})
}

func getDefaultAgency(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/defaultagency"
	r.GET(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
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

func webLogin(r *gin.Engine, eventHandler *event.EventHandler) {
	r.POST("/web/login", func(c *gin.Context) {
		var body datastruct.WebLoginBody
		err := c.BindJSON(&body)
		if err != nil {
			c.JSON(200, gin.H{
				"code": datastruct.ParamError,
			})
			return
		}
		data, code := eventHandler.WebLogin(&body)
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

func editGoods(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/editgoods"
	r.POST(url, func(c *gin.Context) {
		// data1, _ := ioutil.ReadAll(c.Request.Body)
		// log.Debug("---body/---%v", string(data1))
		// return
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		var body datastruct.EditGoodsBody
		err := c.BindJSON(&body)
		if err != nil || !checkEditGoods(&body) {
			c.JSON(200, gin.H{
				"code": datastruct.ParamError,
			})
			return
		}
		code := eventHandler.EditGoods(&body)
		c.JSON(200, gin.H{
			"code": code,
		})
	})
}
func webGetGoods(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/getgoods"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		var body datastruct.WebGetGoodsBody
		err := c.BindJSON(&body)
		if err != nil || body.PageIndex <= 0 || body.PageSize <= 0 || body.IsHidden < 0 {
			c.JSON(200, gin.H{
				"code": datastruct.ParamError,
			})
			return
		}
		data, code := eventHandler.WebGetGoods(&body)
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

func getDomain(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/domain"
	r.GET(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.GetDomain()
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

func getBlackListJump(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/blacklist"
	r.GET(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		data := eventHandler.GetBlackListJump()
		c.JSON(200, gin.H{
			"code": datastruct.NULLError,
			"data": data,
		})
	})
}

func editBlackListJump(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/blacklist"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		var body datastruct.BlackListJumpBody
		err := c.BindJSON(&body)
		if err != nil {
			c.JSON(200, gin.H{
				"code": datastruct.ParamError,
			})
			return
		}
		code := eventHandler.EditBlackListJump(&body)
		c.JSON(200, gin.H{
			"code": code,
		})
	})
}

func checkEditGoods(body *datastruct.EditGoodsBody) bool {
	tf := true
	if body.LimitAmount < 0 || body.PostAge < 0 || body.OriginalPrice < 0 || body.Count < 0 || body.Type < 0 || body.Type > 1 || body.IsHidden < 0 || body.IsHidden > 1 || len(body.LevelData) < datastruct.MaxLevel || len(body.Base64str) <= 0 || body.Brand == "" || body.Classid <= 0 || body.Goodsdesc == "" || body.Name == "" || body.Price <= 0 || body.Pricedesc == "" || body.Rushprice <= 0 || body.Rushpricedesc == "" || body.Sortid < 0 {
		tf = false
	}
	return tf
}

func getPurchaseOrder(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/getpurchaseorder"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		var body datastruct.GetPurchaseBody
		err := c.BindJSON(&body)
		if err != nil || body.State < 0 || body.State > 2 || body.PageIndex <= 0 || body.PageSize <= 0 {
			c.JSON(200, gin.H{
				"code": datastruct.ParamError,
			})
			return
		}
		data, code := eventHandler.GetPurchaseOrder(&body)
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

func getRushOrder(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/getrushorder"
	r.POST(url, func(c *gin.Context) {
		// data1, _ := ioutil.ReadAll(c.Request.Body)
		// log.Debug("---body/---%v", string(data1))
		// return
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		var body datastruct.GetRushOrderBody
		err := c.BindJSON(&body)
		if err != nil || body.State < -1 || body.State > 3 || body.PageIndex <= 0 || body.PageSize <= 0 {
			c.JSON(200, gin.H{
				"code": datastruct.ParamError,
			})
			return
		}
		data, code := eventHandler.GetRushOrder(&body)
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

func getSendGoodsOrder(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/getsendgoods"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		var body datastruct.GetSendGoodsBody
		err := c.BindJSON(&body)
		if err != nil || body.State < 0 || body.State > 2 || body.PageIndex <= 0 || body.PageSize <= 0 {
			c.JSON(200, gin.H{
				"code": datastruct.ParamError,
			})
			return
		}
		data, code := eventHandler.GetSendGoodsOrder(&body)
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

func getMemberLevel(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/getmemberlevel"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		var body datastruct.GetMemberLevelDataBody
		err := c.BindJSON(&body)
		if err != nil {
			c.JSON(200, gin.H{
				"code": datastruct.ParamError,
			})
			return
		}
		data, code := eventHandler.GetMemberLevel(body.Name)
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

func webGetMembers(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/getmembers"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		var body datastruct.WebGetMembersBody
		err := c.BindJSON(&body)
		if err != nil || body.PageIndex <= 0 || body.PageSize <= 0 || body.IsBlacklist < 0 || body.IsBlacklist > 2 {
			c.JSON(200, gin.H{
				"code": datastruct.ParamError,
			})
			return
		}
		data, code := eventHandler.GetMembers(&body)
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

func updateUserBlackList(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/updateuserbl"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		var body datastruct.WebUpdateUserBlBody
		err := c.BindJSON(&body)
		if err != nil || body.UserId <= 0 || body.IsBlacklist < 0 || body.IsBlacklist > 1 {
			c.JSON(200, gin.H{
				"code": datastruct.ParamError,
			})
			return
		}
		c.JSON(200, gin.H{
			"code": eventHandler.UpdateUserBlackList(body.IsBlacklist, body.UserId),
		})
	})
}

func updateUserLevel(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/updateuserlevel"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		var body datastruct.WebUpdateUserLevelBody
		err := c.BindJSON(&body)
		if err != nil || body.UserId <= 0 || body.LevelId < 0 {
			c.JSON(200, gin.H{
				"code": datastruct.ParamError,
			})
			return
		}
		c.JSON(200, gin.H{
			"code": eventHandler.UpdateUserLevel(body.UserId, body.LevelId),
		})
	})
}
func webChangeGold(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/addgoldcount"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		var body datastruct.WebAddGoldBody
		err := c.BindJSON(&body)
		if err != nil || body.UserId <= 0 {
			c.JSON(200, gin.H{
				"code": datastruct.ParamError,
			})
			return
		}
		data, code := eventHandler.WebChangeGold(body.UserId, body.Gold)
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
func myPrentices(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/myprentices"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		var body datastruct.WebGetAgencyInfoBody
		err := c.BindJSON(&body)
		if err != nil || body.PageIndex <= 0 || body.PageSize <= 0 || body.UserId <= 0 || body.Level < 1 || body.Level > 3 || body.StartTime > body.EndTime || body.StartTime < 0 || body.EndTime < 0 {
			c.JSON(200, gin.H{
				"code": datastruct.ParamError,
			})
			return
		}
		data, code := eventHandler.MyPrentices(&body)
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

func getServerInfo(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/serverinfo"
	r.GET(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.GetServerInfoFromDB()
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

func editServerInfo(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/serverinfo"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		var body datastruct.WebServerInfoBody
		err := c.BindJSON(&body)
		if err != nil || body.IsMaintain < 0 || body.IsMaintain > 1 {
			c.JSON(200, gin.H{
				"code": datastruct.ParamError,
			})
			return
		}
		c.JSON(200, gin.H{
			"code": eventHandler.EditServerInfo(body.Version, body.IsMaintain),
		})
	})
}

func updateGoodsClassState(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/updateGoodsClassState/:classid/:ishidden"
	r.GET(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		classid := tools.StringToInt(c.Param("classid"))
		ishidden := tools.StringToInt(c.Param("ishidden"))
		code := eventHandler.UpdateGoodsClassState(classid, ishidden)
		c.JSON(200, gin.H{
			"code": code,
		})
	})
}

func editGoodsClass(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/editgoodsclass"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		var body datastruct.WebEditGoodsClassBody
		err := c.BindJSON(&body)
		if err != nil || body.IsHidden < 0 || body.IsHidden > 1 || body.SortId < 0 || body.Id < 0 {
			c.JSON(200, gin.H{
				"code": datastruct.ParamError,
			})
			return
		}
		c.JSON(200, gin.H{
			"code": eventHandler.EditGoodsClass(&body),
		})
	})
}

func getAllGoodsClasses(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/allgoodsclasses"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		var body datastruct.WebQueryGoodsClassBody
		err := c.BindJSON(&body)
		if err != nil || body.IsHidden < 0 || body.IsHidden > 2 {
			c.JSON(200, gin.H{
				"code": datastruct.ParamError,
			})
			return
		}
		data, code := eventHandler.GetAllGoodsClasses(&body)
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

func getAllDepositInfo(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/querydepositinfo"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.GetAllDepositInfo(c)
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

func getAllDrawInfo(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/querydrawinfo"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.GetAllDrawInfo(c)
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

func getAllMembers(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/availablemembers"
	r.GET(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.GetAllMembers()
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

func updateMemberLevelState(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/memberstate"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		code := eventHandler.UpdateMemberLevelState(c)
		c.JSON(200, gin.H{
			"code": code,
		})
	})
}

func getMemberOrder(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/getmemberorder"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.GetMemberOrder(c)
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

func deleteMemberOrder(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/deletememberorder"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		c.JSON(200, gin.H{
			"code": eventHandler.DeleteMemberOrder(c),
		})
	})
}

func editRandomLotteryGoods(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/editlotterygoods"
	r.POST(url, func(c *gin.Context) {
		// data1, _ := ioutil.ReadAll(c.Request.Body)
		// log.Debug("---body/---%v", string(data1))
		// return
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		code := eventHandler.EditRandomLotteryGoods(c)
		c.JSON(200, gin.H{
			"code": code,
		})
	})
}

func getRandomLotteryGoods(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/getlotterygoods"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.GetRandomLotteryGoods(c)
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

func editRandomLotteryGoodsPool(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/editlotterypool"
	r.POST(url, func(c *gin.Context) {
		// data1, _ := ioutil.ReadAll(c.Request.Body)
		// log.Debug("---body/---%v", string(data1))
		// return
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		code := eventHandler.EditRandomLotteryGoodsPool(c)
		c.JSON(200, gin.H{
			"code": code,
		})
	})
}

func getRandomLotteryGoodsPool(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/getlotterypool"
	r.GET(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.GetRandomLotteryGoodsPool()
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

func getRandomLotteryOrder(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/getlotteryorder"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		var body datastruct.GetSendGoodsBody
		err := c.BindJSON(&body)
		if err != nil || body.State < 0 || body.State > 2 || body.PageIndex <= 0 || body.PageSize <= 0 {
			c.JSON(200, gin.H{
				"code": datastruct.ParamError,
			})
			return
		}
		data, code := eventHandler.GetRandomLotteryOrder(&body)
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

func updateLotteryGoodsSendState(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/updatelotterysendstate"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		c.JSON(200, gin.H{
			"code": eventHandler.UpdateLotteryGoodsSendState(c),
		})
	})
}

func getRushLimitSetting(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/getrushlimitsetting"
	r.GET(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.GetRushLimitSetting()
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

func editRushLimitSetting(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/editrushlimitsetting"
	r.POST(url, func(c *gin.Context) {
		// data1, _ := ioutil.ReadAll(c.Request.Body)
		// log.Debug("---body/---%v", string(data1))
		// return
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		code := eventHandler.EditRushLimitSetting(c)
		c.JSON(200, gin.H{
			"code": code,
		})
	})
}

func getWebStatistics(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/statistics"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.GetWebStatistics(c)
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

func getActiveUsers(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/activeusers"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.GetActiveUsers(c)
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

func editReClass(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/editreclass"
	r.POST(url, func(c *gin.Context) {
		// data1, _ := ioutil.ReadAll(c.Request.Body)
		// log.Debug("---body/---%v", string(data1))
		// return
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		code := eventHandler.EditReClass(c)
		c.JSON(200, gin.H{
			"code": code,
		})
	})
}

func getAllReClass(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/getallreclasses"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.GetAllReClass(c)
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

func editSharePoster(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/editshareposter"
	r.POST(url, func(c *gin.Context) {
		// data1, _ := ioutil.ReadAll(c.Request.Body)
		// log.Debug("---body/---%v", string(data1))
		// return
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		code := eventHandler.EditSharePoster(c)
		c.JSON(200, gin.H{
			"code": code,
		})
	})
}

func getAllSharePosters(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/allshareposters"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.GetAllSharePosters(c)
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

func updateSharePosterState(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/updateposterstate"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		code := eventHandler.UpdateSharePosterState(c)
		c.JSON(200, gin.H{
			"code": code,
		})
	})
}

func editUserAppraise(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/edituserappraise"
	r.POST(url, func(c *gin.Context) {
		// data1, _ := ioutil.ReadAll(c.Request.Body)
		// log.Debug("---body/---%v", string(data1))
		// return
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		code := eventHandler.EditUserAppraise(c)
		c.JSON(200, gin.H{
			"code": code,
		})
	})
}
func deleteUserAppraise(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/deleteuserappraise"
	r.POST(url, func(c *gin.Context) {
		// data1, _ := ioutil.ReadAll(c.Request.Body)
		// log.Debug("---body/---%v", string(data1))
		// return
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		code := eventHandler.DeleteUserAppraise(c)
		c.JSON(200, gin.H{
			"code": code,
		})
	})
}

func getUserAppraise(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/getuserappraise"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.GetUserAppraise(c)
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
func updateSignForState(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/updatesignforstate"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		c.JSON(200, gin.H{
			"code": eventHandler.UpdateSignForState(c),
		})
	})
}

func editGoodsDetail(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/editgoodsdetail"
	r.POST(url, func(c *gin.Context) {
		// data1, _ := ioutil.ReadAll(c.Request.Body)
		// log.Debug("---body/---%v", string(data1))
		// return
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		code := eventHandler.EditGoodsDetail(c)
		c.JSON(200, gin.H{
			"code": code,
		})
	})
}
func getGoodsDetail(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/goodsdetail/:goodsid"
	r.GET(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		goodsid := tools.StringToInt(c.Param("goodsid"))
		if goodsid <= 0 {
			c.JSON(200, gin.H{
				"code": datastruct.ParamError,
			})
			return
		}
		data, code := eventHandler.GetGoodsDetail(goodsid)
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

func getSCParams(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/scparams"
	r.GET(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.GetSCParams()
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

func updateSCParams(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/scparams"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		code := eventHandler.UpdateSCParams(c)
		c.JSON(200, gin.H{
			"code": code,
		})
	})
}

func getSuggestion(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/getsuggest"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.GetSuggestion(c)
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

func deleteSuggestion(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/delsuggestion"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		code := eventHandler.DeleteSuggestion(c)
		c.JSON(200, gin.H{
			"code": code,
		})
	})
}

func deleteComplaint(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/delcomplaint"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		code := eventHandler.DeleteComplaint(c)
		c.JSON(200, gin.H{
			"code": code,
		})
	})
}

func getComplaint(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/getcomplaint"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.GetComplaint(c)
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

func getDrawCashParams(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/drawcashparams"
	r.GET(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.GetDrawCashParams()
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

func updateDrawCashParams(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/drawcashparams"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		code := eventHandler.UpdateDrawCashParams(c)
		c.JSON(200, gin.H{
			"code": code,
		})
	})
}

func drawCashPass(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/drawcashpass"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		code := eventHandler.DrawCashPass(c)
		c.JSON(200, gin.H{
			"code": code,
		})
	})
}

func deleteSharePosters(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/delposters"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		code := eventHandler.DeleteSharePosters(c)
		c.JSON(200, gin.H{
			"code": code,
		})
	})
}

func deleteAd(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/deletead"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		code := eventHandler.DeleteAd(c)
		c.JSON(200, gin.H{
			"code": code,
		})
	})
}

func getAllAd(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/adinfo"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.GetAllAd(c)
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

func editAd(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/editadinfo"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		c.JSON(200, gin.H{
			"code": eventHandler.EditAd(c),
		})
	})
}

func getGoldCoinGift(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/goldcoingift"
	r.GET(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.GetGoldCoinGift()
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

func editGoldCoinGift(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/goldcoingift"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		c.JSON(200, gin.H{
			"code": eventHandler.EditGoldCoinGift(c),
		})
	})
}

func getWebUsers(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/webusers"
	r.GET(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.GetWebUsers()
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

func deleteWebUser(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/deletewebuser"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		code := eventHandler.DeleteWebUser(c)
		c.JSON(200, gin.H{
			"code": code,
		})
	})
}

func editWebUser(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/editwebuser"
	r.POST(url, func(c *gin.Context) {
		token, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.EditWebUser(c, token)
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

func getAllMenuInfo(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/allmenu"
	r.GET(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.GetAllMenuInfo()
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
func checkPermission(c *gin.Context, url string, eventHandler *event.EventHandler) (string, bool) {
	rs := false
	tokens, isExist := c.Request.Header["Webtoken"]
	var rs_token string
	if isExist {
		token := tokens[0]
		if token != "" {
			rs_token = token
			method := c.Request.Method
			rs = eventHandler.CheckPermission(token, method, url)
		}
	}
	if !rs {
		c.JSON(200, gin.H{
			"code": datastruct.WebPermissionDenied,
		})
	}
	return rs_token, rs
}

func updateWebUserPwd(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/updatepwd"
	r.POST(url, func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": eventHandler.UpdateWebUserPwd(c),
		})
	})
}

func getCommissionStatistics(r *gin.Engine, eventHandler *event.EventHandler) {
	url := "/web/commissionstatistics"
	r.POST(url, func(c *gin.Context) {
		_, tf := checkPermission(c, url, eventHandler)
		if !tf {
			return
		}
		data, code := eventHandler.GetCommissionStatistics(c)
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

func WebRegister(r *gin.Engine, eventHandler *event.EventHandler) {
	editDomain(r, eventHandler)                  //添加或修改域名
	updateSendInfo(r, eventHandler)              //商品已发货
	webLogin(r, eventHandler)                    //web登录
	editGoods(r, eventHandler)                   //添加或修改商品信息
	webGetGoods(r, eventHandler)                 //商品查询
	getDomain(r, eventHandler)                   //获取域名信息
	getBlackListJump(r, eventHandler)            //获取黑名单跳转信息
	editBlackListJump(r, eventHandler)           //添加或修改黑名单跳转信息
	getRushOrder(r, eventHandler)                //闯关订单查询
	getPurchaseOrder(r, eventHandler)            //直接购买订单查询
	getSendGoodsOrder(r, eventHandler)           //发货订单查询
	updateDefaultAgency(r, eventHandler)         //修改默认佣金设置
	getDefaultAgency(r, eventHandler)            //获取默认佣金设置
	editMemberLevel(r, eventHandler)             //添加或修改会员等级信息
	getMemberLevel(r, eventHandler)              //获取所有会员等级
	webGetMembers(r, eventHandler)               //会员查询
	updateUserBlackList(r, eventHandler)         //修改会员黑明单状态
	updateUserLevel(r, eventHandler)             //修改会员等级
	webChangeGold(r, eventHandler)               //充值金币
	myPrentices(r, eventHandler)                 //我的下线
	getServerInfo(r, eventHandler)               //获取服务器状态信息
	editServerInfo(r, eventHandler)              //添加或修改服务器状态信息
	updateGoodsClassState(r, eventHandler)       //修改商品类型状态
	editGoodsClass(r, eventHandler)              //添加或修改商品类型
	getAllGoodsClasses(r, eventHandler)          //获取所有商品类型
	getAllDepositInfo(r, eventHandler)           //查询充值信息
	getAllDrawInfo(r, eventHandler)              //查询提现信息
	getAllMembers(r, eventHandler)               //查询所有会员等级信息
	updateMemberLevelState(r, eventHandler)      //会员等级状态修改
	getMemberOrder(r, eventHandler)              //获取等级订单
	deleteMemberOrder(r, eventHandler)           //删除等级订单
	editRandomLotteryGoods(r, eventHandler)      //添加或修改抽奖商品
	getRandomLotteryGoods(r, eventHandler)       //抽奖商品查询
	editRandomLotteryGoodsPool(r, eventHandler)  //编辑抽奖商品池水
	getRandomLotteryGoodsPool(r, eventHandler)   //获取抽奖商品池水
	getRandomLotteryOrder(r, eventHandler)       //抽奖订单发货查询
	updateLotteryGoodsSendState(r, eventHandler) //抽奖商品已发货
	getRushLimitSetting(r, eventHandler)         //获取闯关限制设置
	editRushLimitSetting(r, eventHandler)        //编辑闯关限制设置
	getWebStatistics(r, eventHandler)            //首页统计数据
	editReClass(r, eventHandler)                 //添加或修改推荐类型
	getAllReClass(r, eventHandler)               //获取所有推荐类型
	editSharePoster(r, eventHandler)             //添加或修改分享海报
	getAllSharePosters(r, eventHandler)          //获取所有分享海报
	updateSharePosterState(r, eventHandler)      //修改分享海报状态
	editUserAppraise(r, eventHandler)            //添加或修改用户评价
	deleteUserAppraise(r, eventHandler)          //删除用户评价
	getUserAppraise(r, eventHandler)             //获取用户评价
	updateSignForState(r, eventHandler)          //商品已签收
	editGoodsDetail(r, eventHandler)             //修改商品详情
	getGoodsDetail(r, eventHandler)              //获取商品详情
	getSCParams(r, eventHandler)                 //获取建议与投诉参数设置
	updateSCParams(r, eventHandler)              //修改建议与投诉参数设置
	getSuggestion(r, eventHandler)               //查询用户反馈
	deleteSuggestion(r, eventHandler)            //删除用户反馈
	getComplaint(r, eventHandler)                //查询用户投诉
	deleteComplaint(r, eventHandler)             //删除用户投诉
	getDrawCashParams(r, eventHandler)           //获取提现参数
	updateDrawCashParams(r, eventHandler)        //修改提现参数
	drawCashPass(r, eventHandler)                //提现审核通过
	deleteSharePosters(r, eventHandler)          //删除分享海报
	deleteAd(r, eventHandler)                    //删除广告
	getAllAd(r, eventHandler)                    //查询广告
	editAd(r, eventHandler)                      //添加或修改广告
	getGoldCoinGift(r, eventHandler)             //获取金币赠送设置
	editGoldCoinGift(r, eventHandler)            //修改金币赠送设置
	getActiveUsers(r, eventHandler)              //活跃用户数
	getWebUsers(r, eventHandler)                 //获取web用户
	deleteWebUser(r, eventHandler)               //删除web用户
	editWebUser(r, eventHandler)                 //添加或修改web用户
	getAllMenuInfo(r, eventHandler)              //获取所有菜单信息
	updateWebUserPwd(r, eventHandler)            //修改web用户密码,不需要权限检测
	getCommissionStatistics(r, eventHandler)
}
