package event

import (
	"app/datastruct"
	"app/tools"
	"time"
)

func (handle *EventHandler) createUser(wx_user *datastruct.WXUserData, platform datastruct.Platform, referrer int) *datastruct.UserInfo {
	user := new(datastruct.UserInfo)
	user.MemberIdentifier = datastruct.AgencyIdentifier
	user.Token = tools.UniqueId()
	user.NickName = wx_user.NickName
	user.Sex = wx_user.Sex
	user.Avatar = wx_user.HeadImgUrl
	user.IsCheat = 0
	user.LotterySucceed = 0
	user.DepositTotal = 0
	user.BalanceTotal = 0
	user.GoldCount = 0
	user.CreatedAt = time.Now().Unix()
	user.LoginTime = user.CreatedAt
	user.Platform = platform
	user.IsGotRegisterGift = 0
	user.IsGotDownLoadAppGift = 1
	wxp := new(datastruct.WXPlatform)
	wxp.WXUUID = wx_user.UnionId

	switch platform {
	case datastruct.H5:
		wxp.PayOpenidForGZH = wx_user.OpenId
	case datastruct.APP:
		wxp.PayOpenidForKFPT = wx_user.OpenId
	}

	user.Id = handle.dbHandler.CreateUser(user, wxp, referrer)
	return user
}
