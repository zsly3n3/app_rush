package routes

import (
	"app/datastruct"
	"app/event"
	"app/tools"

	"github.com/gin-gonic/gin"
	//"app/log"
)

func editDomain(r *gin.Engine, eventHandler *event.EventHandler) {
	r.POST("/web/domain", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": eventHandler.EditDomain(c),
		})
	})
}

func updateSendInfo(r *gin.Engine, eventHandler *event.EventHandler) {
	r.POST("/web/updatesendinfo", func(c *gin.Context) {
		eventHandler.UpdateSendInfo(c)
	})
}

func updateDefaultAgency(r *gin.Engine, eventHandler *event.EventHandler) {
	r.POST("/web/defaultagency", func(c *gin.Context) {
		eventHandler.UpdateDefaultAgency(c)
	})
}

func editMemberLevel(r *gin.Engine, eventHandler *event.EventHandler) {
	r.POST("/web/editmemberlevel", func(c *gin.Context) {
		eventHandler.EditMemberLevel(c)
	})
}

func getDefaultAgency(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/web/defaultagency", func(c *gin.Context) {
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
	r.POST("/web/editgoods", func(c *gin.Context) {
		// data1, _ := ioutil.ReadAll(c.Request.Body)
		// log.Debug("---body/---%v", string(data1))
		// return
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
	r.POST("/web/getgoods", func(c *gin.Context) {
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
	r.GET("/web/domain", func(c *gin.Context) {
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
	r.GET("/web/blacklist", func(c *gin.Context) {
		data := eventHandler.GetBlackListJump()
		c.JSON(200, gin.H{
			"code": datastruct.NULLError,
			"data": data,
		})
	})
}

func editBlackListJump(r *gin.Engine, eventHandler *event.EventHandler) {
	r.POST("/web/blacklist", func(c *gin.Context) {
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
	r.POST("/web/getpurchaseorder", func(c *gin.Context) {
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
	r.POST("/web/getrushorder", func(c *gin.Context) {
		// data1, _ := ioutil.ReadAll(c.Request.Body)
		// log.Debug("---body/---%v", string(data1))
		// return
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
	r.POST("/web/getsendgoods", func(c *gin.Context) {
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
	r.POST("/web/getmemberlevel", func(c *gin.Context) {
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
	r.POST("/web/getmembers", func(c *gin.Context) {
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
	r.POST("/web/updateuserbl", func(c *gin.Context) {
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
	r.POST("/web/updateuserlevel", func(c *gin.Context) {
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
	r.POST("/web/addgoldcount", func(c *gin.Context) {
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
	r.POST("/web/myprentices", func(c *gin.Context) {
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
	r.GET("/web/serverinfo", func(c *gin.Context) {
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
	r.POST("/web/serverinfo", func(c *gin.Context) {
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
	r.GET("/web/updateGoodsClassState/:classid/:ishidden", func(c *gin.Context) {
		classid := tools.StringToInt(c.Param("classid"))
		ishidden := tools.StringToInt(c.Param("ishidden"))
		code := eventHandler.UpdateGoodsClassState(classid, ishidden)
		c.JSON(200, gin.H{
			"code": code,
		})
	})
}

func editGoodsClass(r *gin.Engine, eventHandler *event.EventHandler) {
	r.POST("/web/editgoodsclass", func(c *gin.Context) {
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
	r.POST("/web/allgoodsclasses", func(c *gin.Context) {
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
	r.POST("/web/querydepositinfo", func(c *gin.Context) {
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
	r.POST("/web/querydrawinfo", func(c *gin.Context) {
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
	r.GET("/web/availablemembers", func(c *gin.Context) {
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
	r.POST("/web/memberstate", func(c *gin.Context) {
		code := eventHandler.UpdateMemberLevelState(c)
		c.JSON(200, gin.H{
			"code": code,
		})
	})
}

/*
func getNewUsers(r *gin.Engine, eventHandler *event.EventHandler) {
	r.POST("/web/newusers", func(c *gin.Context) {
		data, code := eventHandler.GetNewUsers(c)
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

func getTotalEarn(r *gin.Engine, eventHandler *event.EventHandler) {
	r.POST("/web/totaleran", func(c *gin.Context) {
		data, code := eventHandler.GetTotalEarn(c)
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

func getDepositUsers(r *gin.Engine, eventHandler *event.EventHandler) {
	r.POST("/web/depositusers", func(c *gin.Context) {
		data, code := eventHandler.GetDepositUsers(c)
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

func getActivityUsers(r *gin.Engine, eventHandler *event.EventHandler) {
	r.POST("/web/activityusers", func(c *gin.Context) {
		data, code := eventHandler.GetActivityUsers(c)
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
*/

func getMemberOrder(r *gin.Engine, eventHandler *event.EventHandler) {
	r.POST("/web/getmemberorder", func(c *gin.Context) {
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
	r.POST("/web/deletememberorder", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": eventHandler.DeleteMemberOrder(c),
		})
	})
}

func editRandomLotteryGoods(r *gin.Engine, eventHandler *event.EventHandler) {
	r.POST("/web/editlotterygoods", func(c *gin.Context) {
		// data1, _ := ioutil.ReadAll(c.Request.Body)
		// log.Debug("---body/---%v", string(data1))
		// return
		code := eventHandler.EditRandomLotteryGoods(c)
		c.JSON(200, gin.H{
			"code": code,
		})
	})
}

func getRandomLotteryGoods(r *gin.Engine, eventHandler *event.EventHandler) {
	r.POST("/web/getlotterygoods", func(c *gin.Context) {
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
	r.POST("/web/editlotteryPool", func(c *gin.Context) {
		// data1, _ := ioutil.ReadAll(c.Request.Body)
		// log.Debug("---body/---%v", string(data1))
		// return
		code := eventHandler.EditRandomLotteryGoodsPool(c)
		c.JSON(200, gin.H{
			"code": code,
		})
	})
}

func getRandomLotteryGoodsPool(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/web/getlotterypool", func(c *gin.Context) {
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
	r.POST("/web/getlotteryorder", func(c *gin.Context) {
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
	r.POST("/web/updatelotterysendstate", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": eventHandler.UpdateLotteryGoodsSendState(c),
		})
	})
}

func getRushLimitSetting(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/web/getrushlimitsetting", func(c *gin.Context) {
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
	r.POST("/web/editrushlimitsetting", func(c *gin.Context) {
		// data1, _ := ioutil.ReadAll(c.Request.Body)
		// log.Debug("---body/---%v", string(data1))
		// return
		code := eventHandler.EditRushLimitSetting(c)
		c.JSON(200, gin.H{
			"code": code,
		})
	})
}

func getWebStatistics(r *gin.Engine, eventHandler *event.EventHandler) {
	r.POST("/web/statistics", func(c *gin.Context) {
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
	r.POST("/web/activeUsers", func(c *gin.Context) {
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
	r.POST("/web/editreclass", func(c *gin.Context) {
		// data1, _ := ioutil.ReadAll(c.Request.Body)
		// log.Debug("---body/---%v", string(data1))
		// return
		code := eventHandler.EditReClass(c)
		c.JSON(200, gin.H{
			"code": code,
		})
	})
}

func getAllReClass(r *gin.Engine, eventHandler *event.EventHandler) {
	r.POST("/web/getallreclasses", func(c *gin.Context) {
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
	r.POST("/web/editshareposter", func(c *gin.Context) {
		// data1, _ := ioutil.ReadAll(c.Request.Body)
		// log.Debug("---body/---%v", string(data1))
		// return
		code := eventHandler.EditSharePoster(c)
		c.JSON(200, gin.H{
			"code": code,
		})
	})
}

func getAllSharePosters(r *gin.Engine, eventHandler *event.EventHandler) {
	r.POST("/web/allshareposters", func(c *gin.Context) {
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
	r.POST("/web/updateposterstate", func(c *gin.Context) {
		code := eventHandler.UpdateSharePosterState(c)
		c.JSON(200, gin.H{
			"code": code,
		})
	})
}

func editUserAppraise(r *gin.Engine, eventHandler *event.EventHandler) {
	r.POST("/web/edituserappraise", func(c *gin.Context) {
		// data1, _ := ioutil.ReadAll(c.Request.Body)
		// log.Debug("---body/---%v", string(data1))
		// return
		code := eventHandler.EditUserAppraise(c)
		c.JSON(200, gin.H{
			"code": code,
		})
	})
}
func deleteUserAppraise(r *gin.Engine, eventHandler *event.EventHandler) {
	r.POST("/web/deleteuserappraise", func(c *gin.Context) {
		// data1, _ := ioutil.ReadAll(c.Request.Body)
		// log.Debug("---body/---%v", string(data1))
		// return
		code := eventHandler.DeleteUserAppraise(c)
		c.JSON(200, gin.H{
			"code": code,
		})
	})
}

func getUserAppraise(r *gin.Engine, eventHandler *event.EventHandler) {
	r.POST("/web/getuserappraise", func(c *gin.Context) {
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
	r.POST("/web/updatesignforstate", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": eventHandler.UpdateSignForState(c),
		})

	})
}

func editGoodsDetail(r *gin.Engine, eventHandler *event.EventHandler) {
	r.POST("/web/editgoodsdetail", func(c *gin.Context) {
		// data1, _ := ioutil.ReadAll(c.Request.Body)
		// log.Debug("---body/---%v", string(data1))
		// return
		code := eventHandler.EditGoodsDetail(c)
		c.JSON(200, gin.H{
			"code": code,
		})
	})
}
func getGoodsDetail(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/web/goodsdetail/:goodsid", func(c *gin.Context) {
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
	r.GET("/web/scparams", func(c *gin.Context) {
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
	r.POST("/web/scparams", func(c *gin.Context) {
		code := eventHandler.UpdateSCParams(c)
		c.JSON(200, gin.H{
			"code": code,
		})
	})
}

func getSuggestion(r *gin.Engine, eventHandler *event.EventHandler) {
	r.POST("/web/getsuggest", func(c *gin.Context) {
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
	r.POST("/web/delsuggestion", func(c *gin.Context) {
		code := eventHandler.DeleteSuggestion(c)
		c.JSON(200, gin.H{
			"code": code,
		})
	})
}

func deleteComplaint(r *gin.Engine, eventHandler *event.EventHandler) {
	r.POST("/web/delcomplaint", func(c *gin.Context) {
		code := eventHandler.DeleteComplaint(c)
		c.JSON(200, gin.H{
			"code": code,
		})
	})
}

func getComplaint(r *gin.Engine, eventHandler *event.EventHandler) {
	r.POST("/web/getcomplaint", func(c *gin.Context) {
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
	r.GET("/web/drawcashparams", func(c *gin.Context) {
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
	r.POST("/web/drawcashparams", func(c *gin.Context) {
		code := eventHandler.UpdateDrawCashParams(c)
		c.JSON(200, gin.H{
			"code": code,
		})
	})
}

func drawCashPass(r *gin.Engine, eventHandler *event.EventHandler) {
	r.POST("/web/drawcashpass", func(c *gin.Context) {
		code := eventHandler.DrawCashPass(c)
		c.JSON(200, gin.H{
			"code": code,
		})
	})
}

func deleteSharePosters(r *gin.Engine, eventHandler *event.EventHandler) {
	r.POST("/web/delposters", func(c *gin.Context) {
		code := eventHandler.DeleteSharePosters(c)
		c.JSON(200, gin.H{
			"code": code,
		})
	})
}

func deleteAd(r *gin.Engine, eventHandler *event.EventHandler) {
	r.POST("/web/deletead", func(c *gin.Context) {
		code := eventHandler.DeleteAd(c)
		c.JSON(200, gin.H{
			"code": code,
		})
	})
}

func getAllAd(r *gin.Engine, eventHandler *event.EventHandler) {
	r.POST("/web/adinfo", func(c *gin.Context) {
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
	r.POST("/web/editadinfo", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": eventHandler.EditAd(c),
		})
	})
}

func getGoldCoinGift(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/web/goldcoingift", func(c *gin.Context) {
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
	r.POST("/web/goldcoingift", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": eventHandler.EditGoldCoinGift(c),
		})
	})
}

func getWebUsers(r *gin.Engine, eventHandler *event.EventHandler) {
	r.GET("/web/webusers", func(c *gin.Context) {
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
	r.POST("/web/deletewebuser", func(c *gin.Context) {
		code := eventHandler.DeleteWebUser(c)
		c.JSON(200, gin.H{
			"code": code,
		})
	})
}

func editWebUser(r *gin.Engine, eventHandler *event.EventHandler) {
	r.POST("/web/editwebuser", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": eventHandler.EditWebUser(c),
		})
	})
}

func WebRegister(r *gin.Engine, eventHandler *event.EventHandler) {
	editDomain(r, eventHandler)
	updateSendInfo(r, eventHandler)
	webLogin(r, eventHandler)
	editGoods(r, eventHandler)
	webGetGoods(r, eventHandler)
	getDomain(r, eventHandler)
	getBlackListJump(r, eventHandler)
	editBlackListJump(r, eventHandler)
	getRushOrder(r, eventHandler)
	getPurchaseOrder(r, eventHandler)
	getSendGoodsOrder(r, eventHandler)
	updateDefaultAgency(r, eventHandler)
	getDefaultAgency(r, eventHandler)
	editMemberLevel(r, eventHandler)
	getMemberLevel(r, eventHandler)
	webGetMembers(r, eventHandler)
	updateUserBlackList(r, eventHandler)
	updateUserLevel(r, eventHandler)
	webChangeGold(r, eventHandler)
	myPrentices(r, eventHandler)
	getServerInfo(r, eventHandler)
	editServerInfo(r, eventHandler)
	updateGoodsClassState(r, eventHandler)
	editGoodsClass(r, eventHandler)
	getAllGoodsClasses(r, eventHandler)
	getAllDepositInfo(r, eventHandler)
	getAllDrawInfo(r, eventHandler)
	getAllMembers(r, eventHandler)
	updateMemberLevelState(r, eventHandler)
	// getNewUsers(r, eventHandler)
	// getTotalEarn(r, eventHandler)
	// getDepositUsers(r, eventHandler)
	// getActivityUsers(r, eventHandler)
	getMemberOrder(r, eventHandler)
	deleteMemberOrder(r, eventHandler)
	editRandomLotteryGoods(r, eventHandler)
	getRandomLotteryGoods(r, eventHandler)
	editRandomLotteryGoodsPool(r, eventHandler)
	getRandomLotteryGoodsPool(r, eventHandler)
	getRandomLotteryOrder(r, eventHandler)
	updateLotteryGoodsSendState(r, eventHandler)
	getRushLimitSetting(r, eventHandler)
	editRushLimitSetting(r, eventHandler)
	getWebStatistics(r, eventHandler)
	editReClass(r, eventHandler)
	getAllReClass(r, eventHandler)
	editSharePoster(r, eventHandler)
	getAllSharePosters(r, eventHandler)
	updateSharePosterState(r, eventHandler)
	editUserAppraise(r, eventHandler)
	deleteUserAppraise(r, eventHandler)
	getUserAppraise(r, eventHandler)
	updateSignForState(r, eventHandler)
	editGoodsDetail(r, eventHandler)
	getGoodsDetail(r, eventHandler)
	getSCParams(r, eventHandler)
	updateSCParams(r, eventHandler)
	getSuggestion(r, eventHandler)
	deleteSuggestion(r, eventHandler)
	getComplaint(r, eventHandler)
	deleteComplaint(r, eventHandler)
	getDrawCashParams(r, eventHandler)
	updateDrawCashParams(r, eventHandler)
	drawCashPass(r, eventHandler)
	deleteSharePosters(r, eventHandler)
	deleteAd(r, eventHandler)
	getAllAd(r, eventHandler)
	editAd(r, eventHandler)
	getGoldCoinGift(r, eventHandler)
	editGoldCoinGift(r, eventHandler)
	getActiveUsers(r, eventHandler)
	getWebUsers(r, eventHandler)
	deleteWebUser(r, eventHandler)
	editWebUser(r, eventHandler)
}
