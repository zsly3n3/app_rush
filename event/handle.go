package event

import (
	"app/conf"
	"app/datastruct"
	"app/db"
	"app/log"
	"app/osstool"
	"app/thirdParty"
	"app/tools"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func (handle *EventHandler) GetDBHandler() *db.DBHandler {
	return handle.dbHandler
}

func (handle *EventHandler) IsExistUser(token string) (int, bool, bool) {
	conn := handle.cacheHandler.GetConn()
	defer conn.Close()
	var userId int
	var tf bool
	var isBlackList bool
	userId, tf, isBlackList = handle.cacheHandler.IsExistUserWithConn(conn, token)
	if !tf {
		var user *datastruct.UserInfo
		user, tf = handle.dbHandler.GetUserDataWithToken(token)
		if tf {
			userId = user.Id
			isBlackList = tools.IntToBool(user.IsBlackList)
			handle.cacheHandler.SetUserAllData(conn, user)
			handle.cacheHandler.AddExpire(conn, token)
		}
	} else {
		handle.cacheHandler.AddExpire(conn, token)
	}
	return userId, tf, isBlackList
}

func (handle *EventHandler) AppLogin(c *gin.Context) {
	var body datastruct.AppLoginBody
	err := c.BindJSON(&body)
	code := datastruct.NULLError
	if err == nil {
		if body.Code == "" {
			c.JSON(200, gin.H{
				"code": datastruct.WXCodeInvalid,
			})
			return
		}
		openid, access_token, isError := tools.GetOpenidAndAccessToken(body.Code, datastruct.WX_KFPT_AppID, datastruct.WX_KFPT_AppSecret)
		if isError {
			c.JSON(200, gin.H{
				"code": datastruct.WXCodeInvalid,
			})
			return
		}
		handle.login(openid, access_token, datastruct.APP, 0, c)
	} else {
		code = datastruct.JsonParseFailedFromPostBody
		c.JSON(200, gin.H{
			"code": code,
		})
	}
}

func (handle *EventHandler) H5Login(c *gin.Context) {
	var body datastruct.H5LoginBody
	err := c.BindJSON(&body)
	code := datastruct.NULLError
	if err == nil {
		if body.Code == "" {
			c.JSON(200, gin.H{
				"code": datastruct.WXCodeInvalid,
			})
			return
		}
		openid, access_token, isError := tools.GetOpenidAndAccessToken(body.Code, datastruct.WX_GZH_AppID, datastruct.WX_GZH_AppSecret)
		if isError {
			c.JSON(200, gin.H{
				"code": datastruct.WXCodeInvalid,
			})
			return
		}
		handle.login(openid, access_token, datastruct.H5, body.Referrer, c)
	} else {
		code = datastruct.JsonParseFailedFromPostBody
		c.JSON(200, gin.H{
			"code": code,
		})
	}
}

func (handle *EventHandler) login(openid string, access_token string, platform datastruct.Platform, referrer int, c *gin.Context) {
	wxp_user, isError := tools.GetWXUserData(openid, access_token, platform)
	if isError {
		c.JSON(200, gin.H{
			"code": datastruct.WXCodeInvalid,
		})
		return
	}
	conn := handle.cacheHandler.GetConn()
	defer conn.Close()

	u_data, isExistMysql := handle.dbHandler.GetUserDataWithUnionId(wxp_user.UnionId) //find in mysql
	if !isExistMysql {
		u_data = handle.createUser(wxp_user, platform, referrer)
	}
	handle.UpdateUserInfo(u_data.Id, openid, platform)
	handle.cacheHandler.SetUserAllData(conn, u_data)
	handle.cacheHandler.AddExpire(conn, u_data.Token)
	mp := datastruct.ResponseLoginData(u_data)
	c.JSON(200, gin.H{
		"code": datastruct.NULLError,
		"data": mp,
	})
}
func (handle *EventHandler) UpdateUserInfo(userId int, openid string, platform datastruct.Platform) {
	if platform == datastruct.APP { //修改支付openid
		handle.dbHandler.UpdateAppOpenId(userId, openid)
	} else {
		handle.dbHandler.UpdateGZHOpenId(userId, openid)
	}
}

func (handle *EventHandler) GetHomeData(platform int) *datastruct.ResponseHomeData {
	return handle.dbHandler.GetHomeData(platform)
}

func (handle *EventHandler) GetGoods(pageIndex int, pageSize int, classid int, user_id int) []*datastruct.ResponseGoodsData {
	return handle.dbHandler.GetGoods(pageIndex, pageSize, classid, user_id)
}

func (handle *EventHandler) GetGoodsClass(userId int) (interface{}, datastruct.CodeType) {
	return handle.dbHandler.GetGoodsClass(userId)
}

func (handle *EventHandler) WorldRewardInfo(pageIndex int, pageSize int) (interface{}, datastruct.CodeType) {
	return handle.dbHandler.WorldRewardInfo(pageIndex, pageSize)
}

func (handle *EventHandler) FreeRougeGame() (interface{}, datastruct.CodeType) {
	return handle.dbHandler.FreeRougeGame()
}

func (handle *EventHandler) PayRougeGame(userId int, goodsid int, c *gin.Context) {
	rs, code := handle.dbHandler.PayRougeGame(userId, goodsid)
	if code == datastruct.NULLError || code == datastruct.CheatUser {
		c.JSON(200, gin.H{
			"code": code,
			"data": rs,
		})
	} else {
		c.JSON(200, gin.H{
			"code": code,
		})
	}
}

func (handle *EventHandler) LevelPass(userId int, c *gin.Context) {
	var body datastruct.LevelPassBody
	err := c.BindJSON(&body)
	var code datastruct.CodeType
	if err == nil {
		code = handle.dbHandler.LevelPass(userId, body.Id)
	} else {
		code = datastruct.JsonParseFailedFromPostBody
	}
	c.JSON(200, gin.H{
		"code": code,
	})
}

func (handle *EventHandler) LevelPassFailed(userId int, c *gin.Context) {
	var body datastruct.LevelPassFailedBody
	err := c.BindJSON(&body)
	var code datastruct.CodeType
	if err == nil {
		code = handle.dbHandler.LevelPassFailed(userId, body.Id, body.Number)
	} else {
		code = datastruct.JsonParseFailedFromPostBody
	}
	c.JSON(200, gin.H{
		"code": code,
	})
}

func (handle *EventHandler) LevelPassSucceed(userId int, c *gin.Context) {
	var body datastruct.LevelPassBody
	err := c.BindJSON(&body)
	var code datastruct.CodeType
	if err == nil {
		code = handle.dbHandler.LevelPassSucceed(userId, body.Id, body.Platform)
	} else {
		code = datastruct.JsonParseFailedFromPostBody
	}
	c.JSON(200, gin.H{
		"code": code,
	})
}

func (handle *EventHandler) GetNotAppliedOrderInfo(userId int, pageIndex int, pageSize int) (interface{}, datastruct.CodeType) {
	return handle.dbHandler.GetNotAppliedOrderInfo(userId, pageIndex, pageSize)
}

func (handle *EventHandler) GetNotSendGoods(userId int, pageIndex int, pageSize int) (interface{}, datastruct.CodeType) {
	return handle.dbHandler.GetNotSendGoods(userId, pageIndex, pageSize)
}

func (handle *EventHandler) GetHasSendedGoods(userId int, pageIndex int, pageSize int) (interface{}, datastruct.CodeType) {
	return handle.dbHandler.GetHasSendedGoods(userId, pageIndex, pageSize)
}

func (handle *EventHandler) GetAppraiseOrder(userId int, pageIndex int, pageSize int) (interface{}, datastruct.CodeType) {
	return handle.dbHandler.GetAppraiseOrder(userId, pageIndex, pageSize)
}

func (handle *EventHandler) ApplySend(userId int, c *gin.Context) {
	var body datastruct.ApplySendBody
	err := c.BindJSON(&body)
	var code datastruct.CodeType
	if err == nil && body.Address != "" && body.LinkMan != "" && body.PhoneNumber != "" {
		code = handle.dbHandler.ApplySend(userId, &body)
	} else {
		code = datastruct.ParamError
	}
	c.JSON(200, gin.H{
		"code": code,
	})
}

func (handle *EventHandler) CommissionRank(userId int, pageIndex int, pageSize int) (interface{}, datastruct.CodeType) {
	return handle.dbHandler.CommissionRank(userId, pageIndex, pageSize)
}

func (handle *EventHandler) CommissionInfo(userId int, pageIndex int, pageSize int) (interface{}, datastruct.CodeType) {
	return handle.dbHandler.CommissionInfo(userId, pageIndex, pageSize)
}

func (handle *EventHandler) GetAgentlevelN(userId int, level int, pageIndex int, pageSize int) (interface{}, datastruct.CodeType) {
	return handle.dbHandler.GetAgentlevelN(userId, level, pageIndex, pageSize)
}

func (handle *EventHandler) GetDrawCash(userId int, pageIndex int, pageSize int) (interface{}, datastruct.CodeType) {
	return handle.dbHandler.GetDrawCash(userId, pageIndex, pageSize)
}

func (handle *EventHandler) GetDrawCashRule(userId int) (interface{}, datastruct.CodeType) {
	return handle.dbHandler.GetDrawCashRule(userId)
}

func (handle *EventHandler) GetDepositParams(userId int) (interface{}, datastruct.CodeType) {
	return handle.dbHandler.GetDepositParams(userId)
}

func (handle *EventHandler) GetInviteUrl(userId int) (interface{}, datastruct.CodeType) {
	entryUrl := handle.dbHandler.GetEntryUrl()
	url_str := tools.CreateInviteLink(tools.IntToString(userId), entryUrl)
	return url_str, datastruct.NULLError
}

func (handle *EventHandler) GetInviteUrlNoToken() (interface{}, datastruct.CodeType) {
	entryUrl := handle.dbHandler.GetEntryUrl()
	return entryUrl, datastruct.NULLError
}

func (handle *EventHandler) GetAppDownLoadShareUrl(userId int) (interface{}, datastruct.CodeType) {
	shareDownLoad := handle.dbHandler.GetAppDownLoadShareUrl()
	url_str := tools.CreateDownLoadLink(tools.IntToString(userId), shareDownLoad)
	return url_str, datastruct.NULLError
}

func (handle *EventHandler) GetDirectDownloadApp() string {
	return handle.dbHandler.GetDirectDownloadApp()
}

func (handle *EventHandler) DepositSucceed(userId int, money int64) (interface{}, datastruct.CodeType) {
	return handle.dbHandler.DepositSucceed(userId, money, datastruct.H5)
}

func (handle *EventHandler) GetUserInfo(userId int, platform datastruct.Platform) (interface{}, datastruct.CodeType) {
	return handle.dbHandler.GetUserInfo(userId, platform)
}

func (handle *EventHandler) GetGoldInfo(userId int, pageIndex int, pageSize int) (interface{}, datastruct.CodeType) {
	return handle.dbHandler.GetGoldInfo(userId, pageIndex, pageSize)
}

func (handle *EventHandler) Purchase(userId int, purchaseType datastruct.PurchaseType, purchase_id int, ip_addr string, platform datastruct.Platform) (interface{}, datastruct.CodeType) {
	var code datastruct.CodeType
	var m_data *datastruct.MemberLevelData
	if purchaseType == datastruct.VipType {
		m_data, code = handle.dbHandler.QueryMemberLevelData(purchase_id)
		if code != datastruct.NULLError {
			return nil, code
		}
		code = handle.dbHandler.IsRefreshMemberList(userId, m_data.Level)
		if code != datastruct.NULLError {
			return nil, code
		}
	}

	var wx_user *datastruct.WXPlatform
	wx_user, code = handle.dbHandler.GetOpenId(userId)
	if code != datastruct.NULLError {
		return nil, datastruct.PurchaseFailed
	}
	orderForm := new(datastruct.OrderForm)
	orderForm.PurchaseType = purchaseType
	orderForm.PurchaseId = purchase_id

	var money float64
	switch purchaseType {
	case datastruct.GoodsType:
		var goods *datastruct.Goods
		goods, code = handle.dbHandler.GetGoodsPrice(purchase_id)
		if code != datastruct.NULLError {
			return nil, code
		}
		money = float64(goods.Price)
		orderForm.Desc = fmt.Sprintf("%s 价格:%f", goods.Name, money)

	case datastruct.GoldType:
		money = commondata.DepositParams[purchase_id]
		orderForm.Desc = fmt.Sprintf("%s%d个", datastruct.ProductDesc, int(money))

	case datastruct.VipType:
		money = float64(m_data.Price)
		orderForm.Desc = fmt.Sprintf("%s 价格:%f", m_data.Name, money)
	}

	orderForm.Id = tools.UniqueId()
	orderForm.CreatedAt = time.Now().Unix()
	orderForm.Money = int64(money)
	orderForm.UserId = userId
	orderForm.Platform = platform
	wx_pay := new(datastruct.WX_PayInfo)
	wx_pay.OrderForm = orderForm

	if platform == datastruct.H5 {
		wx_pay.Appid = datastruct.WX_GZH_AppID
		wx_pay.Mch_id = datastruct.WX_GZH_Mch_Id
		wx_pay.PaySecret = datastruct.WX_GZH_PaySecret
		orderForm.OpenId = wx_user.PayOpenidForGZH
		wx_pay.Trade_type = "JSAPI"
	} else {
		wx_pay.Appid = datastruct.WX_KFPT_AppID
		wx_pay.Mch_id = datastruct.WX_KFPT_Mch_Id
		wx_pay.PaySecret = datastruct.WX_KFPT_PaySecret
		orderForm.OpenId = wx_user.PayOpenidForKFPT
		wx_pay.Trade_type = "APP"
	}
	resp_wx_pay, err := thirdParty.WXunifyChargeReq(wx_pay, ip_addr)
	if err != nil {
		log.Error("WXunifyChargeReq err:%v", err.Error())
		return nil, datastruct.DepositFailed
	}
	var rs_interface interface{}
	if platform == datastruct.H5 {
		resp := new(datastruct.ResponsH5Pay)
		resp.AppId = resp_wx_pay.AppId
		resp.NonceStr = resp_wx_pay.NonceStr
		resp.TimeStamp = fmt.Sprintf("%d", orderForm.CreatedAt)
		resp.Package = fmt.Sprintf("prepay_id=%s", resp_wx_pay.Prepay_id)
		resp.SignType = "MD5"

		var m map[string]interface{}
		m = make(map[string]interface{}, 0)
		m["appId"] = resp.AppId
		m["nonceStr"] = resp.NonceStr
		m["package"] = resp.Package
		m["timeStamp"] = resp.TimeStamp
		m["signType"] = resp.SignType

		resp.PaySign = thirdParty.WXpayCalcSign(m, wx_pay.PaySecret)
		rs_interface = resp
	} else {

		resp := new(datastruct.ResponsAppPay)
		resp.PartnerId = resp_wx_pay.Mch_id
		resp.NonceStr = resp_wx_pay.NonceStr
		resp.TimeStamp = fmt.Sprintf("%d", orderForm.CreatedAt)
		resp.Package = "Sign=WXPay"
		resp.PrepayId = resp_wx_pay.Prepay_id

		var m map[string]interface{}
		m = make(map[string]interface{}, 0)
		m["partnerid"] = resp.PartnerId
		m["noncestr"] = resp.NonceStr
		m["package"] = resp.Package
		m["timestamp"] = resp.TimeStamp
		m["prepayid"] = resp.PrepayId
		m["appid"] = resp_wx_pay.AppId
		resp.Sign = thirdParty.WXpayCalcSign(m, wx_pay.PaySecret)
		rs_interface = resp
	}
	handle.cacheHandler.CreateOrderForm(orderForm)
	return rs_interface, datastruct.NULLError
}

func (handle *EventHandler) WxPayResultCall(c *gin.Context) {
	var body datastruct.WXPayResultNoticeBody
	err := c.BindXML(&body)
	if err == nil {
		if body.Result_code != "SUCCESS" {
			c.XML(200, gin.H{
				"return_code": "FAIL",
				"return_msg":  "OK",
			})
			return
		}
		mp := make(map[string]interface{})
		mp["appid"] = body.AppId
		mp["bank_type"] = body.Bank_type
		mp["cash_fee"] = body.Cash_fee
		mp["fee_type"] = body.Fee_type
		mp["is_subscribe"] = body.Is_subscribe
		mp["mch_id"] = body.Mch_id
		mp["nonce_str"] = body.Nonce_str
		mp["openid"] = body.Openid
		mp["out_trade_no"] = body.Out_trade_no
		mp["result_code"] = body.Result_code
		mp["return_code"] = body.Return_code
		mp["time_end"] = body.Time_end
		mp["total_fee"] = body.Total_fee
		mp["trade_type"] = body.Trade_type
		mp["transaction_id"] = body.Transaction_id
		paySecret := datastruct.WX_KFPT_PaySecret
		if body.Trade_type == "JSAPI" {
			paySecret = datastruct.WX_GZH_PaySecret
		}
		if thirdParty.WXpayCalcSign(mp, paySecret) == body.Sign {
			handle.payMutex.Lock()
			defer handle.payMutex.Unlock()
			orderForm, code := handle.cacheHandler.GetPayData(body.Out_trade_no)
			if code != datastruct.NULLError {
				c.XML(200, gin.H{
					"return_code": "FAIL",
					"return_msg":  "OK",
				})
				return
			}
			switch orderForm.PurchaseType {
			case datastruct.GoldType:
				_, code = handle.dbHandler.DepositSucceed(orderForm.UserId, orderForm.Money, orderForm.Platform)
			case datastruct.GoodsType:
				code = handle.dbHandler.PurchaseSuccess(orderForm.UserId, orderForm.PurchaseId, orderForm.CreatedAt, orderForm.Platform)
			case datastruct.VipType:
				code = handle.dbHandler.PurchaseVipSuccess(orderForm.UserId, orderForm.PurchaseId, orderForm.CreatedAt)
			}
			if code != datastruct.NULLError {
				c.XML(200, gin.H{
					"return_code": "FAIL",
					"return_msg":  "OK",
				})
				return
			}
			handle.cacheHandler.DeletedKeys([]interface{}{body.Out_trade_no})
			c.XML(200, gin.H{
				"return_code": "SUCCESS",
				"return_msg":  "OK",
			})
		} else {
			log.Error("WxPayResultCall WXpayCalcSign err")
			c.XML(200, gin.H{
				"return_code": "FAIL",
				"return_msg":  "OK",
			})
		}
	} else {
		log.Error("WxPayResultCall BindXML ERR:%s", err.Error())
		c.XML(200, gin.H{
			"return_code": "FAIL",
			"return_msg":  "OK",
		})
	}
}

func (handle *EventHandler) UserPayee(userId int, ip_addr string, c *gin.Context) (interface{}, datastruct.CodeType) {
	var body datastruct.PayeeBody
	err := c.BindJSON(&body)
	var code datastruct.CodeType
	var rs interface{}
	if err == nil {
		var wxp_user *datastruct.WXPlatform
		wxp_user, code = handle.dbHandler.GetOpenId(userId)
		if code != datastruct.NULLError {
			return nil, datastruct.GetDataFailed
		}
		isOnlyApp := handle.dbHandler.IsDrawCashOnApp()
		payeeOpenid := ""
		if isOnlyApp {
			if wxp_user.PayeeOpenid == "" {
				return nil, datastruct.PayeeOnlyInApp
			}
			payeeOpenid = wxp_user.PayeeOpenid
		} else {
			payeeOpenid = wxp_user.PayOpenidForGZH
		}

		trade_no := tools.UniqueId()
		var poundage float64
		poundage, code = handle.dbHandler.ComputePoundage(userId, body.Amount, trade_no, ip_addr)
		rs_amount := body.Amount - poundage
		if code == datastruct.PayeeReview {
			return rs_amount, code
		}
		if code != datastruct.NULLError {
			return nil, code
		}
		rs_payee, err := thirdParty.WXunifyPayeeReq(rs_amount, payeeOpenid, ip_addr, trade_no, isOnlyApp)
		if err == nil && rs_payee.Result_code == "SUCCESS" && rs_payee.Partner_trade_no == trade_no {
			rs, code = handle.dbHandler.UserPayeeSuccess(userId, rs_payee)
		} else {
			if rs_payee.Err_code_des != "" {
				rs, code = handle.dbHandler.UserPayeefailed(userId, body.Amount, rs_payee)
				log.Debug("UserPayee err:%s", rs_payee.Err_code_des)
			}
			code = datastruct.WeChatPayeeError
		}
	} else {
		code = datastruct.JsonParseFailedFromPostBody
		rs = nil
	}
	return rs, code
}

func (handle *EventHandler) CustomShareForApp(userId int) (interface{}, datastruct.CodeType) {
	return handle.dbHandler.CustomShareForApp(userId)
}

func (handle *EventHandler) CustomShareForGZH(userId int) (interface{}, datastruct.CodeType) {
	return handle.dbHandler.CustomShareForGZH(userId)
}

func (handle *EventHandler) GetAuthAddr(userId string) string {
	authUrl := handle.dbHandler.GetAuthUrl()
	link := tools.CreateAuthLink(userId, datastruct.WX_GZH_AppID, authUrl)
	return link
}

func (handle *EventHandler) GetKfInfo() (interface{}, datastruct.CodeType) {
	return handle.dbHandler.GetKfInfo()
}

func (handle *EventHandler) GetAppAddr() (interface{}, datastruct.CodeType) {
	addr_info, code := handle.dbHandler.GetAppAddr()
	var port_str string
	if conf.Common.Mode == conf.Debug {
		port_str = ":7006"
	} else {
		port_str = ":8080"
	}
	return "http://" + addr_info.Url + port_str, code
}

func (handle *EventHandler) GetDownLoadAppAddr() (interface{}, datastruct.CodeType) {
	return handle.dbHandler.GetDownLoadAppAddr()
}

func (handle *EventHandler) GetBlackListRedirect() string {
	return commondata.BlacklistRedirect
}

func (handle *EventHandler) GetPCRedirect() string {
	return commondata.PcRedirect
}

func (handle *EventHandler) GetEntryPageUrl() string {
	return handle.dbHandler.GetEntryPageUrl()
}

func (handle *EventHandler) GetServerInfoFromMemory() (string, int) {
	commondata.ServerInfo.RWMutex.RLock()
	defer commondata.ServerInfo.RWMutex.RUnlock()
	return commondata.ServerInfo.Version, commondata.ServerInfo.IsMaintain
}
func (handle *EventHandler) GetCheckInData(userId int) (interface{}, datastruct.CodeType) {
	return handle.dbHandler.GetCheckInData(userId)
}
func (handle *EventHandler) CheckIn(userId int) (interface{}, datastruct.CodeType) {
	return handle.dbHandler.CheckIn(userId)
}

func (handle *EventHandler) AppGetMemberList(userId int) (interface{}, datastruct.CodeType) {
	return handle.dbHandler.AppGetMemberList(userId)
}

func (handle *EventHandler) StartLottery(userId int, rushprice int64, c *gin.Context) {
	now_time := time.Now().Unix()
	commondata.LotteryQueue.RWMutex.Lock()
	defer commondata.LotteryQueue.RWMutex.Unlock()
	handle_nowTime := time.Now().Unix()
	if handle_nowTime-now_time >= datastruct.MaxTimeOutForLotteryQueue {
		c.JSON(200, gin.H{
			"code": datastruct.LotteryTimeOut,
		})
	} else {
		rs, code := handle.dbHandler.StartLottery(userId, rushprice)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": rs,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	}
}
func (handle *EventHandler) GetLotteryGoodsSucceedHistory() (interface{}, datastruct.CodeType) {
	return handle.dbHandler.GetLotteryGoodsSucceedHistory()
}

func (handle *EventHandler) GetLotteryGoodsOrderInfo(userId int, pageIndex int, pageSize int) (interface{}, datastruct.CodeType) {
	return handle.dbHandler.GetLotteryGoodsOrderInfo(userId, pageIndex, pageSize)
}

func (handle *EventHandler) GetSharePosters(userId int) (interface{}, datastruct.CodeType) {
	qrcode, _ := handle.GetInviteUrl(userId)
	return handle.dbHandler.GetSharePosters(qrcode.(string))
}

func (handle *EventHandler) GetUserAppraiseForApp(pageIndex int, pageSize int) (interface{}, datastruct.CodeType) {
	return handle.dbHandler.GetUserAppraiseForApp(pageIndex, pageSize)
}

func (handle *EventHandler) UserAppraise(userId int, c *gin.Context) datastruct.CodeType {
	var body datastruct.UserAppraiseBody
	err := c.BindJSON(&body)
	if err != nil || body.Number == "" || (body.Desc == "" && len(body.ImgNames) == 0) {
		return datastruct.ParamError
	}
	code := handle.dbHandler.UserAppraise(userId, &body)
	if code != datastruct.NULLError {
		for _, v := range body.ImgNames {
			go osstool.DeleteFile(commondata.OSSBucket, v)
		}
	}
	return code
}

func (handle *EventHandler) GetGoodsDetailForApp(goodsid int) (interface{}, datastruct.CodeType) {
	return handle.dbHandler.GetGoodsDetailForApp(goodsid)
}

func (handle *EventHandler) UpdateUserAddress(userId int, c *gin.Context) datastruct.CodeType {
	var body datastruct.ReceiverForSendGoods
	err := c.BindJSON(&body)
	if err != nil || body.PhoneNumber == "" || body.Address == "" || body.LinkMan == "" {
		return datastruct.ParamError
	}
	return handle.dbHandler.UpdateUserAddress(userId, &body)
}
func (handle *EventHandler) GetUserAddress(userId int) (interface{}, datastruct.CodeType) {
	return handle.dbHandler.GetUserAddress(userId)
}

func (handle *EventHandler) GetAgencyPage(userId int) (*datastruct.ResponseAgencyPage, datastruct.CodeType) {
	resp, code := handle.dbHandler.GetAgencyPage(userId)
	if code == datastruct.NULLError {
		entryUrl := handle.dbHandler.GetEntryUrl()
		url_str := tools.CreateInviteLink(tools.IntToString(userId), entryUrl)
		resp.InviteUrl = url_str
	}
	return resp, code
}

func (handle *EventHandler) RemindSendGoods(userId int, number string) datastruct.CodeType {
	return handle.dbHandler.RemindSendGoods(userId, number)
}

func (handle *EventHandler) AddSuggest(userId int, c *gin.Context) datastruct.CodeType {
	var body datastruct.SuggestBody
	err := c.BindJSON(&body)
	if err != nil || body.Desc == "" {
		return datastruct.ParamError
	}
	return handle.dbHandler.AddSuggest(userId, body.Desc)
}

func (handle *EventHandler) AddComplaint(userId int, c *gin.Context) datastruct.CodeType {
	var body datastruct.ComplaintBody
	err := c.BindJSON(&body)
	if err != nil || body.Desc == "" {
		return datastruct.ParamError
	}
	code := handle.dbHandler.AddComplaint(userId, body.ComplaintType, body.Desc)
	if code == datastruct.BlackList {
		code = handle.UpdateUserBlackList(1, userId)
	}
	return code
}

func (handle *EventHandler) GetDownLoadAppGift(userId int) datastruct.CodeType {
	return handle.dbHandler.GetDownLoadAppGift(userId)
}

func (handle *EventHandler) GetRegisterGift(userId int) datastruct.CodeType {
	return handle.dbHandler.GetRegisterGift(userId)
}

func (handle *EventHandler) IsRefreshHomeGoodsData(userId int, classId int) (interface{}, datastruct.CodeType) {
	return handle.dbHandler.IsRefreshHomeGoodsData(userId, classId)
}

func (handle *EventHandler) GetGoldFromPoster(userId int, gpid int) (interface{}, datastruct.CodeType) {
	addr, _ := handle.GetAppAddr()
	return handle.dbHandler.GetGoldFromPoster(userId, gpid, addr.(string))
}

func (handle *EventHandler) GetRandomLotteryList() (interface{}, datastruct.CodeType) {
	return handle.dbHandler.GetRandomLotteryList()
}
func (handle *EventHandler) GetAgentCount(userId int) (interface{}, datastruct.CodeType) {
	return handle.dbHandler.GetAgentCount(userId)
}

// func (handle *EventHandler) UserActivate(userId int) datastruct.CodeType {
// 	return handle.dbHandler.UserActivate(userId)
// }
