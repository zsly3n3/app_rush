package event

import (
	"app/datastruct"
	"app/log"
	"app/osstool"
	"app/thirdParty"

	"github.com/gin-gonic/gin"
)

func (handle *EventHandler) WebLogin(body *datastruct.WebLoginBody) (interface{}, datastruct.CodeType) {
	return handle.dbHandler.WebLogin(body)
}

func (handle *EventHandler) UpdateSendInfo(c *gin.Context) {
	var body datastruct.UpdateSendInfoBody
	err := c.BindJSON(&body)
	var code datastruct.CodeType
	if err == nil && body.OrderNumber != "" && body.ExpressAgency != "" && body.ExpressNumber != "" {
		code = handle.dbHandler.UpdateSendInfo(&body)
	} else {
		code = datastruct.JsonParseFailedFromPostBody
	}
	c.JSON(200, gin.H{
		"code": code,
	})
}
func (handle *EventHandler) EditGoods(body *datastruct.EditGoodsBody) datastruct.CodeType {
	return handle.dbHandler.EditGoods(body)
}

func (handle *EventHandler) WebGetGoods(body *datastruct.WebGetGoodsBody) (interface{}, datastruct.CodeType) {
	return handle.dbHandler.WebGetGoods(body)
}

func (handle *EventHandler) EditDomain(c *gin.Context) datastruct.CodeType {
	var body datastruct.WebEditDomainBody
	err := c.BindJSON(&body)
	var code datastruct.CodeType
	if err == nil && body.DownLoadUrl != "" && body.AppDomain != "" && body.AuthDomain != "" && body.EntryDomain != "" && body.IOSApp != "" && body.AndroidApp != "" {
		code = handle.dbHandler.EditDomain(&body)
	} else {
		code = datastruct.JsonParseFailedFromPostBody
	}
	return code
}

func (handle *EventHandler) GetDomain() (interface{}, datastruct.CodeType) {
	resp, code := handle.dbHandler.GetDomain()
	resp.Version, _ = handle.GetServerInfoFromMemory()
	return resp, code
}

func (handle *EventHandler) GetBlackListJump() interface{} {
	return handle.dbHandler.GetBlackListJump()
}

func (handle *EventHandler) EditBlackListJump(body *datastruct.BlackListJumpBody) interface{} {
	return handle.dbHandler.EditBlackListJump(body)
}

func (handle *EventHandler) GetRushOrder(body *datastruct.GetRushOrderBody) (interface{}, datastruct.CodeType) {
	return handle.dbHandler.GetRushOrder(body)
}

func (handle *EventHandler) GetPurchaseOrder(body *datastruct.GetPurchaseBody) (interface{}, datastruct.CodeType) {
	return handle.dbHandler.GetPurchaseOrder(body)
}
func (handle *EventHandler) GetSendGoodsOrder(body *datastruct.GetSendGoodsBody) (interface{}, datastruct.CodeType) {
	return handle.dbHandler.GetSendGoodsOrder(body)
}

func (handle *EventHandler) GetDefaultAgency() (interface{}, datastruct.CodeType) {
	return handle.dbHandler.GetDefaultAgency()
}

func (handle *EventHandler) UpdateDefaultAgency(c *gin.Context) {
	var body datastruct.DefaultAgencyBody
	err := c.BindJSON(&body)
	var code datastruct.CodeType
	if err == nil && checkAgencyNumber(body.Agent1Gold) && checkAgencyNumber(body.Agent2Gold) && checkAgencyNumber(body.Agent3Gold) && checkAgencyNumber(body.Agent1Money) && checkAgencyNumber(body.Agent2Money) && checkAgencyNumber(body.Agent3Money) {
		code = handle.dbHandler.UpdateDefaultAgency(&body)
	} else {
		code = datastruct.JsonParseFailedFromPostBody
	}
	c.JSON(200, gin.H{
		"code": code,
	})
}

func checkAgencyNumber(value int) bool {
	tf := true
	if value < 0 || value > 100 {
		tf = false
	}
	return tf
}

func (handle *EventHandler) EditMemberLevel(c *gin.Context) {
	var body datastruct.EditLevelDataBody
	err := c.BindJSON(&body)
	var code datastruct.CodeType
	if err == nil && (body.IsHidden == 0 || body.IsHidden == 1) && body.Price > 0 && (body.Level >= 1 && body.Level <= 8) && checkAgencyNumber(body.Agent1Gold) && checkAgencyNumber(body.Agent2Gold) && checkAgencyNumber(body.Agent3Gold) && checkAgencyNumber(body.Agent1Money) && checkAgencyNumber(body.Agent2Money) && checkAgencyNumber(body.Agent3Money) {
		code = handle.dbHandler.EditMemberLevel(&body)
	} else {
		code = datastruct.JsonParseFailedFromPostBody
	}
	c.JSON(200, gin.H{
		"code": code,
	})
}

func (handle *EventHandler) GetMemberLevel(name string) (interface{}, datastruct.CodeType) {
	return handle.dbHandler.GetMemberLevel(name)
}

func (handle *EventHandler) GetMembers(body *datastruct.WebGetMembersBody) (interface{}, datastruct.CodeType) {
	return handle.dbHandler.GetMembers(body)
}

func (handle *EventHandler) UpdateUserBlackList(state int, userId int) datastruct.CodeType {
	token, code := handle.dbHandler.UpdateUserBlackList(userId, state)
	if code == datastruct.NULLError {
		handle.cacheHandler.UpdateBlackList(token, state)
	}
	return code
}

func (handle *EventHandler) UpdateUserLevel(userId int, state int) datastruct.CodeType {
	return handle.dbHandler.UpdateUserLevel(userId, state)
}

func (handle *EventHandler) WebChangeGold(userId int, gold int64) (interface{}, datastruct.CodeType) {
	return handle.dbHandler.WebChangeGold(userId, gold)
}

func (handle *EventHandler) MyPrentices(body *datastruct.WebGetAgencyInfoBody) (interface{}, datastruct.CodeType) {
	return handle.dbHandler.MyPrentices(body)
}
func (handle *EventHandler) GetServerInfoFromDB() (*datastruct.WebServerInfoBody, datastruct.CodeType) {
	return handle.dbHandler.GetServerInfo()
}

func (handle *EventHandler) EditServerInfo(version string, isMaintain int) datastruct.CodeType {
	code := handle.dbHandler.EditServerInfo(version, isMaintain)
	if code != datastruct.NULLError {
		return code
	}
	commondata.ServerInfo.RWMutex.Lock()
	defer commondata.ServerInfo.RWMutex.Unlock()
	commondata.ServerInfo.Version = version
	commondata.ServerInfo.IsMaintain = isMaintain
	return datastruct.NULLError
}

func (handle *EventHandler) UpdateGoodsClassState(classid int, isHidden int) datastruct.CodeType {
	if classid <= 0 || isHidden < 0 || isHidden > 1 {
		return datastruct.ParamError
	}
	return handle.dbHandler.UpdateGoodsClassState(classid, isHidden)
}

func (handle *EventHandler) EditGoodsClass(body *datastruct.WebEditGoodsClassBody) datastruct.CodeType {
	return handle.dbHandler.EditGoodsClass(body)
}
func (handle *EventHandler) GetAllGoodsClasses(body *datastruct.WebQueryGoodsClassBody) (interface{}, datastruct.CodeType) {
	return handle.dbHandler.GetAllGoodsClasses(body)
}

func (handle *EventHandler) GetAllDepositInfo(c *gin.Context) (interface{}, datastruct.CodeType) {
	var body datastruct.WebQueryDepositInfoBody
	err := c.BindJSON(&body)
	if err != nil || body.PageSize <= 0 || body.PageIndex <= 0 || body.Platform < 0 || body.Platform > 2 {
		return nil, datastruct.ParamError
	}
	return handle.dbHandler.GetAllDepositInfo(&body)
}

func (handle *EventHandler) GetAllDrawInfo(c *gin.Context) (interface{}, datastruct.CodeType) {
	var body datastruct.WebQueryDrawInfoBody
	err := c.BindJSON(&body)
	if err != nil || body.PageSize <= 0 || body.PageIndex <= 0 || body.State < 0 || body.State > 3 {
		return nil, datastruct.ParamError
	}
	return handle.dbHandler.GetAllDrawInfo(&body)
}

func (handle *EventHandler) GetAllMembers() (interface{}, datastruct.CodeType) {
	return handle.dbHandler.GetAllMembers()
}

func (handle *EventHandler) UpdateMemberLevelState(c *gin.Context) datastruct.CodeType {
	var body datastruct.WebUpdateMemberLevelBody
	err := c.BindJSON(&body)
	if err != nil || body.Id <= 0 || body.IsHidden < 0 || body.IsHidden > 1 {
		return datastruct.ParamError
	}
	return handle.dbHandler.UpdateMemberLevelState(&body)
}

/*
func (handle *EventHandler) GetNewUsers(c *gin.Context) (interface{}, datastruct.CodeType) {
	var body datastruct.WebNewsUserBody
	err := c.BindJSON(&body)
	if err != nil || body.StartTime < 0 || body.EndTime < 0 || body.StartTime >= body.EndTime {
		return nil, datastruct.ParamError
	}
	return handle.dbHandler.GetNewUsers(&body)
}

func (handle *EventHandler) GetTotalEarn(c *gin.Context) (interface{}, datastruct.CodeType) {
	var body datastruct.WebNewsUserBody
	err := c.BindJSON(&body)
	if err != nil || body.StartTime < 0 || body.EndTime < 0 || body.StartTime >= body.EndTime {
		return nil, datastruct.ParamError
	}
	return handle.dbHandler.GetTotalEarn(&body)
}

func (handle *EventHandler) GetDepositUsers(c *gin.Context) (interface{}, datastruct.CodeType) {
	var body datastruct.WebNewsUserBody
	err := c.BindJSON(&body)
	if err != nil || body.StartTime < 0 || body.EndTime < 0 || body.StartTime >= body.EndTime {
		return nil, datastruct.ParamError
	}
	return handle.dbHandler.GetDepositUsers(&body)
}

func (handle *EventHandler) GetActivityUsers(c *gin.Context) (interface{}, datastruct.CodeType) {
	var body datastruct.WebNewsUserBody
	err := c.BindJSON(&body)
	if err != nil || body.StartTime < 0 || body.EndTime < 0 || body.StartTime >= body.EndTime {
		return nil, datastruct.ParamError
	}
	return handle.dbHandler.GetActivityUsers(&body)
}*/

func (handle *EventHandler) GetMemberOrder(c *gin.Context) (interface{}, datastruct.CodeType) {
	var body datastruct.WebMemberOrderBody
	err := c.BindJSON(&body)
	if err != nil || body.PageIndex <= 0 || body.PageSize <= 0 || body.StartTime < 0 || body.EndTime < 0 || body.StartTime > body.EndTime {
		return nil, datastruct.ParamError
	}
	return handle.dbHandler.GetMemberOrder(&body)
}

func (handle *EventHandler) DeleteMemberOrder(c *gin.Context) datastruct.CodeType {
	var body datastruct.WebDeleteMemberBody
	err := c.BindJSON(&body)
	if err != nil || body.Id <= 0 {
		return datastruct.ParamError
	}
	return handle.dbHandler.DeleteMemberOrder(body.Id)
}

func (handle *EventHandler) EditRandomLotteryGoods(c *gin.Context) datastruct.CodeType {
	var body datastruct.EditRandomLotteryGoodsBody
	err := c.BindJSON(&body)
	if err != nil || !checkRandomLotteryGoods(&body) {
		return datastruct.ParamError
	}
	return handle.dbHandler.EditRandomLotteryGoods(&body)
}

func checkRandomLotteryGoods(body *datastruct.EditRandomLotteryGoodsBody) bool {
	tf := true
	if body.IsHidden < 0 || body.IsHidden > 1 || body.Base64str == "" || body.Classid <= 0 || body.Name == "" || body.Price <= 0 || body.Probability > 100 || body.Probability < 0 {
		tf = false
	}
	return tf
}

func (handle *EventHandler) GetRandomLotteryGoods(c *gin.Context) (interface{}, datastruct.CodeType) {
	var body datastruct.WebGetGoodsBody
	err := c.BindJSON(&body)
	if err != nil || body.PageIndex <= 0 || body.PageSize <= 0 || body.IsHidden < 0 {
		return nil, datastruct.ParamError
	}
	return handle.dbHandler.WebGetRandomLotteryGoods(&body)
}

func (handle *EventHandler) EditRandomLotteryGoodsPool(c *gin.Context) datastruct.CodeType {
	var body datastruct.WebResponseRandomLotteryPool
	err := c.BindJSON(&body)
	if err != nil || body.Current < 0 || body.Probability <= 0 || body.RandomLotteryCount <= 0 {
		return datastruct.ParamError
	}
	return handle.dbHandler.EditRandomLotteryGoodsPool(&body)
}

func (handle *EventHandler) GetRandomLotteryGoodsPool() (interface{}, datastruct.CodeType) {
	return handle.dbHandler.WebGetRandomLotteryGoodsPool()
}

func (handle *EventHandler) GetRandomLotteryOrder(body *datastruct.GetSendGoodsBody) (interface{}, datastruct.CodeType) {
	return handle.dbHandler.GetRandomLotteryOrder(body)
}

func (handle *EventHandler) UpdateLotteryGoodsSendState(c *gin.Context) datastruct.CodeType {
	var body datastruct.WebResponseLotteryGoodsSendStateBody
	err := c.BindJSON(&body)
	if err != nil || body.OrderNumber == "" || body.ExpressNumber == "" || body.ExpressAgency == "" || body.LinkMan == "" || body.PhoneNumber == "" || body.Address == "" {
		return datastruct.ParamError
	}
	return handle.dbHandler.UpdateLotteryGoodsSendState(&body)
}

func (handle *EventHandler) GetRushLimitSetting() (interface{}, datastruct.CodeType) {
	return handle.dbHandler.GetRushLimitSetting()
}

func (handle *EventHandler) EditRushLimitSetting(c *gin.Context) datastruct.CodeType {
	var body datastruct.WebRushLimitSettingBody
	err := c.BindJSON(&body)
	if err != nil || !checkRushLimitSettingBody(&body) {
		return datastruct.ParamError
	}
	return handle.dbHandler.EditRushLimitSetting(&body)
}

func checkRushLimitSettingBody(body *datastruct.WebRushLimitSettingBody) bool {
	tf := true
	if body.Diff2 < 0 || body.Diff3 < 0 || body.Diff2r <= 0 || body.Diff2t <= 0 || body.Diff3r <= 0 || body.Diff3t <= 0 || body.CheatCount <= 0 || body.LotteryCount <= 0 || body.CheatTips == "" {
		tf = false
	}
	return tf
}

func (handle *EventHandler) GetWebStatistics(c *gin.Context) (interface{}, datastruct.CodeType) {
	var body datastruct.WebNewsUserBody
	err := c.BindJSON(&body)
	if err != nil || body.StartTime < 0 || body.EndTime < 0 || body.StartTime >= body.EndTime {
		return nil, datastruct.ParamError
	}
	if body.EndTime-body.StartTime > 3600*24*31*2 {
		return nil, datastruct.DateTooLong
	}
	var day_sec int64
	day_sec = 3600 * 24
	list := make([]interface{}, 0)
	startTime := body.StartTime
	for {
		new_body := new(datastruct.WebNewsUserBody)
		new_body.StartTime = startTime
		new_body.EndTime = startTime + day_sec
		new_body.RPlatform = body.RPlatform
		new_body.PayPlatform = body.PayPlatform
		rs, _ := handle.dbHandler.WebStatistics(new_body)
		list = append(list, rs)
		if new_body.EndTime >= body.EndTime {
			break
		}
		startTime = new_body.EndTime
	}
	var rs interface{}
	rs = list
	return rs, datastruct.NULLError
}

func (handle *EventHandler) GetActiveUsers(c *gin.Context) (interface{}, datastruct.CodeType) {
	var body datastruct.WebActiveUserBody
	err := c.BindJSON(&body)
	if err != nil || body.StartTime < 0 || body.EndTime < 0 || body.StartTime >= body.EndTime {
		return nil, datastruct.ParamError
	}
	if body.EndTime-body.StartTime > 3600*24*31*2 {
		return nil, datastruct.DateTooLong
	}
	var day_sec int64
	day_sec = 3600 * 24
	list := make([]interface{}, 0)
	startTime := body.StartTime
	for {
		new_body := new(datastruct.WebActiveUserBody)
		new_body.StartTime = startTime
		new_body.EndTime = startTime + day_sec
		new_body.RPlatform = body.RPlatform
		rs, _ := handle.dbHandler.GetActiveUsers(new_body)
		list = append(list, rs)
		if new_body.EndTime >= body.EndTime {
			break
		}
		startTime = new_body.EndTime
	}
	var rs interface{}
	rs = list
	return rs, datastruct.NULLError
}

func (handle *EventHandler) EditReClass(c *gin.Context) datastruct.CodeType {
	var body datastruct.WebEditReClassBody
	err := c.BindJSON(&body)
	if err != nil || body.Id < 0 || body.Name == "" || body.Icon == "" {
		return datastruct.ParamError
	}
	return handle.dbHandler.EditReClass(&body)
}

func (handle *EventHandler) GetAllReClass(c *gin.Context) (interface{}, datastruct.CodeType) {
	var body datastruct.WebQueryReClassBody
	err := c.BindJSON(&body)
	if err != nil {
		return nil, datastruct.ParamError
	}
	return handle.dbHandler.GetAllReClass(&body)
}

func (handle *EventHandler) EditSharePoster(c *gin.Context) datastruct.CodeType {
	var body datastruct.WebEditSharePosterBody
	err := c.BindJSON(&body)
	if err != nil || body.Location < 0 || body.Id < 0 || body.IsHidden < 0 || body.IsHidden > 1 || body.SortId < 0 || body.ImgUrl == "" || body.Icon == "" {
		return datastruct.ParamError
	}
	return handle.dbHandler.EditSharePoster(&body, commondata.OSSBucket)
}

func (handle *EventHandler) GetAllSharePosters(c *gin.Context) (interface{}, datastruct.CodeType) {
	var body datastruct.WebQuerySharePostersBody
	err := c.BindJSON(&body)
	if err != nil || body.IsHidden < 0 || body.IsHidden > 2 {
		return nil, datastruct.ParamError
	}
	return handle.dbHandler.GetAllSharePosters(&body)
}
func (handle *EventHandler) UpdateSharePosterState(c *gin.Context) datastruct.CodeType {
	var body datastruct.WebHiddenPostersBody
	err := c.BindJSON(&body)
	if err != nil || body.IsHidden < 0 || body.IsHidden > 1 {
		return datastruct.ParamError
	}
	return handle.dbHandler.UpdateSharePosterState(&body)
}

func (handle *EventHandler) EditUserAppraise(c *gin.Context) datastruct.CodeType {
	var body datastruct.WebEditUserAppraiseBody
	err := c.BindJSON(&body)
	if err != nil || body.GoodsType < 0 || (body.Desc == "" && len(body.ImgNames) == 0) || body.Id < 0 || body.UserId < 0 || body.GoodsId < 0 || body.IsPassed < 0 || body.IsPassed > 1 || body.TimeStamp <= 0 {
		return datastruct.ParamError
	}
	return handle.dbHandler.EditUserAppraise(&body, commondata.OSSBucket)
}

func (handle *EventHandler) DeleteUserAppraise(c *gin.Context) datastruct.CodeType {
	var body datastruct.WebDeleteUserAppraiseBody
	err := c.BindJSON(&body)
	if err != nil || body.Id < 0 {
		return datastruct.ParamError
	}
	return handle.dbHandler.DeleteUserAppraise(&body, commondata.OSSBucket)
}

func (handle *EventHandler) GetUserAppraise(c *gin.Context) (interface{}, datastruct.CodeType) {
	var body datastruct.WebQueryUserAppraiseBody
	err := c.BindJSON(&body)
	if err != nil || body.IsPassed < 0 || body.IsPassed > 1 || body.PageIndex < 0 || body.PageSize < 0 {
		return nil, datastruct.ParamError
	}
	return handle.dbHandler.GetUserAppraise(&body)
}

func (handle *EventHandler) UpdateSignForState(c *gin.Context) datastruct.CodeType {
	var body datastruct.UpdateSignForBody
	err := c.BindJSON(&body)
	if err != nil || body.OrderNumber == "" {
		return datastruct.ParamError
	}
	return handle.dbHandler.UpdateSignForState(body.OrderNumber)
}

func (handle *EventHandler) EditGoodsDetail(c *gin.Context) datastruct.CodeType {
	var body datastruct.WebEditGoodsDetailBody
	err := c.BindJSON(&body)
	if err != nil || len(body.ImgNames) <= 0 || body.GoodsId < 0 {
		return datastruct.ParamError
	}
	code := handle.dbHandler.EditGoodsDetail(&body, commondata.OSSBucket)
	if code != datastruct.NULLError {
		for _, v := range body.ImgNames {
			go osstool.DeleteFile(commondata.OSSBucket, v)
		}
	}
	return code
}

func (handle *EventHandler) GetGoodsDetail(goodsid int) (interface{}, datastruct.CodeType) {
	return handle.dbHandler.GetGoodsDetail(goodsid)
}

func (handle *EventHandler) GetSCParams() (interface{}, datastruct.CodeType) {
	return handle.dbHandler.GetSCParams()
}

func (handle *EventHandler) UpdateSCParams(c *gin.Context) datastruct.CodeType {
	var body datastruct.WebResponseSCP
	err := c.BindJSON(&body)
	if err != nil || body.SCFD < 0 || body.CCFBL < 0 || body.CCFD < 0 {
		return datastruct.ParamError
	}
	return handle.dbHandler.UpdateSCParams(&body)
}

func (handle *EventHandler) GetSuggestion(c *gin.Context) (interface{}, datastruct.CodeType) {
	var body datastruct.WebQuerySuggestBody
	err := c.BindJSON(&body)
	if err != nil || body.PageIndex <= 0 || body.PageSize <= 0 {
		return nil, datastruct.ParamError
	}
	return handle.dbHandler.GetSuggestion(&body)
}

func (handle *EventHandler) DeleteSuggestion(c *gin.Context) datastruct.CodeType {
	var body datastruct.WebDeleteUserAppraiseBody
	err := c.BindJSON(&body)
	if err != nil || body.Id <= 0 {
		return datastruct.ParamError
	}
	return handle.dbHandler.DeleteSuggestion(body.Id)
}

func (handle *EventHandler) DeleteComplaint(c *gin.Context) datastruct.CodeType {
	var body datastruct.WebDeleteUserAppraiseBody
	err := c.BindJSON(&body)
	if err != nil || body.Id <= 0 {
		return datastruct.ParamError
	}
	return handle.dbHandler.DeleteComplaint(body.Id)
}

func (handle *EventHandler) GetComplaint(c *gin.Context) (interface{}, datastruct.CodeType) {
	var body datastruct.WebQueryComplaintBody
	err := c.BindJSON(&body)
	if err != nil || body.PageIndex <= 0 || body.PageSize <= 0 {
		return nil, datastruct.ParamError
	}
	return handle.dbHandler.GetComplaint(&body)
}

func (handle *EventHandler) GetDrawCashParams() (interface{}, datastruct.CodeType) {
	return handle.dbHandler.GetDrawCashParams()
}

func (handle *EventHandler) UpdateDrawCashParams(c *gin.Context) datastruct.CodeType {
	var body datastruct.WebResponseDrawCashParams
	err := c.BindJSON(&body)
	if err != nil || body.MaxDrawCount <= 0 || body.MinCharge <= 0 || body.MinPoundage < 0 || body.PoundagePer < 0 || body.RequireVerify <= 0 {
		return datastruct.ParamError
	}
	return handle.dbHandler.UpdateDrawCashParams(&body)
}

func (handle *EventHandler) DrawCashPass(c *gin.Context) datastruct.CodeType {
	var body datastruct.WebDeleteMemberBody
	err := c.BindJSON(&body)
	if err != nil || body.Id <= 0 {
		return datastruct.ParamError
	}
	payeeOpenid, isOnlyApp, data, code := handle.dbHandler.GetDrawCashInfo(body.Id)
	if code != datastruct.NULLError {
		return code
	}
	rs_payee, err := thirdParty.WXunifyPayeeReq(data.Charge, payeeOpenid, data.IpAddr, data.TradeNo, isOnlyApp)
	if err == nil && rs_payee.Result_code == "SUCCESS" && rs_payee.Partner_trade_no == data.TradeNo {
		_, code = handle.dbHandler.UserPayeeSuccess(data.UserId, rs_payee)
	} else {
		if rs_payee.Err_code_des != "" {
			_, code = handle.dbHandler.UserPayeefailed(data.UserId, data.Origin, rs_payee)
			log.Debug("UserPayee err:%s", rs_payee.Err_code_des)
		}
		code = datastruct.WeChatPayeeError
	}
	return code
}

func (handle *EventHandler) DeleteSharePosters(c *gin.Context) datastruct.CodeType {
	var body datastruct.WebDeleteMemberBody
	err := c.BindJSON(&body)
	if err != nil || body.Id <= 0 {
		return datastruct.ParamError
	}
	return handle.dbHandler.DeleteSharePosters(body.Id)
}

func (handle *EventHandler) DeleteAd(c *gin.Context) datastruct.CodeType {
	var body datastruct.WebDeleteMemberBody
	err := c.BindJSON(&body)
	if err != nil || body.Id <= 0 {
		return datastruct.ParamError
	}
	return handle.dbHandler.DeleteAd(body.Id)
}

func (handle *EventHandler) GetAllAd(c *gin.Context) (interface{}, datastruct.CodeType) {
	var body datastruct.WebQueryAdBody
	err := c.BindJSON(&body)
	if err != nil || body.IsHidden < 0 || body.IsHidden > 2 || body.Location < 0 || body.Location > 3 || body.Platform < 0 || body.Platform > 2 {
		return nil, datastruct.ParamError
	}
	return handle.dbHandler.GetAllAd(&body)
}

func (handle *EventHandler) EditAd(c *gin.Context) datastruct.CodeType {
	var body datastruct.WebResponseAdInfo
	err := c.BindJSON(&body)
	if err != nil || body.Id < 0 || body.SortId < 0 || body.IsHidden < 0 || body.IsHidden > 1 || body.Location < 0 || body.Location > 2 || body.Platform < 0 || body.Platform > 2 || body.IsJump < 0 || body.IsJump > 1 || body.ImgUrl == "" {
		return datastruct.ParamError
	}
	return handle.dbHandler.EditAd(&body, commondata.OSSBucket)
}
func (handle *EventHandler) GetGoldCoinGift() (interface{}, datastruct.CodeType) {
	return handle.dbHandler.GetGoldCoinGift()
}

func (handle *EventHandler) EditGoldCoinGift(c *gin.Context) datastruct.CodeType {
	var body datastruct.WebResponseGoldCoinGift
	err := c.BindJSON(&body)
	if err != nil || body.AppraisedGoldGift <= 0 || body.DownLoadAppGoldGift <= 0 || body.RegisterGoldGift <= 0 || body.IsEnableRegisterGift < 0 || body.IsEnableRegisterGift > 1 || body.IsDrawCashOnlyApp < 0 || body.IsDrawCashOnlyApp > 1 {
		return datastruct.ParamError
	}
	return handle.dbHandler.EditGoldCoinGift(&body)
}

func (handle *EventHandler) GetWebUsers() (interface{}, datastruct.CodeType) {
	return handle.dbHandler.GetWebUsers()
}

func (handle *EventHandler) DeleteWebUser(c *gin.Context) datastruct.CodeType {
	var body datastruct.WebDeleteMemberBody
	err := c.BindJSON(&body)
	if err != nil || body.Id <= datastruct.AdminLevelID {
		return datastruct.ParamError
	}
	return handle.dbHandler.DeleteWebUser(body.Id)
}

func (handle *EventHandler) EditWebUser(c *gin.Context) datastruct.CodeType {
	var body datastruct.WebEditPermissionUserBody
	err := c.BindJSON(&body)
	if err != nil || body.Id == datastruct.AdminLevelID || body.Id < 0 || body.LoginName == "" || body.Name == "" || len(body.PermissionIds) <= 0 || (body.Pwd == "" && body.Id == 0) {
		return datastruct.ParamError
	}
	return handle.dbHandler.EditWebUser(&body)
}

func (handle *EventHandler) GetAllMenuInfo() (interface{}, datastruct.CodeType) {
	return handle.dbHandler.GetAllMenuInfo()
}

func (handle *EventHandler) CheckPermission(token string, method string, url string) bool {
	return handle.dbHandler.CheckPermission(token, method, url)
}

func (handle *EventHandler) UpdateWebUserPwd(c *gin.Context) datastruct.CodeType {
	tokens, isExist := c.Request.Header["Webtoken"]
	rs_token := ""
	if isExist {
		token := tokens[0]
		if token != "" {
			rs_token = token
		}
	}
	var body datastruct.WebUserPwdBody
	err := c.BindJSON(&body)
	if err != nil || body.NewPwd == "" || body.OldPwd == "" || rs_token == "" {
		return datastruct.ParamError
	}
	return handle.dbHandler.UpdateWebUserPwd(&body, rs_token)
}
