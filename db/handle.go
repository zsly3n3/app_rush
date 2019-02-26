package db

import (
	"app/datastruct"
	"app/log"
	"app/osstool"
	"app/thirdParty"
	"app/tools"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-xorm/xorm"
)

func (handle *DBHandler) GetUserDataWithUnionId(unionId string) (*datastruct.UserInfo, bool) {
	isExist := false
	var u_data *datastruct.UserInfo
	engine := handle.mysqlEngine
	wxp := new(datastruct.WXPlatform)
	has, _ := engine.Where("w_x_u_u_i_d=?", unionId).Get(wxp)
	if has {
		isExist = true
		u_data = handle.GetUserDataFromDataBase(wxp.UserId)
	}
	return u_data, isExist
}

func (handle *DBHandler) GetUserDataWithToken(token string) (*datastruct.UserInfo, bool) {
	engine := handle.mysqlEngine
	u_data := new(datastruct.UserInfo)
	has, err := engine.Where("token=?", token).Get(u_data)
	if err != nil || !has {
		return nil, false
	}
	sql := "update user_info set login_time = ? where token = ?"
	engine.Exec(sql, time.Now().Unix(), token)
	return u_data, true
}

func (handle *DBHandler) GetUserDataFromDataBase(userId int) *datastruct.UserInfo {
	engine := handle.mysqlEngine
	u_info := new(datastruct.UserInfo)
	engine.Where("id=?", userId).Get(u_info)
	sql := "update user_info set login_time = ? , token = ?  where id = ?"
	engine.Exec(sql, time.Now().Unix(), tools.UniqueId(), userId)
	return u_info
}

func (handle *DBHandler) CreateUser(user *datastruct.UserInfo, wxp *datastruct.WXPlatform, referrer int) int {
	engine := handle.mysqlEngine
	session := engine.NewSession()
	defer session.Close()
	session.Begin()

	_, err := session.Insert(user)
	if err != nil {
		str := fmt.Sprintf("DBHandler->CreateUser Insert UserInfo :%s", err.Error())
		rollback(str, session)
		return -1
	}
	wxp.UserId = user.Id

	_, err = session.Insert(wxp)
	if err != nil {
		str := fmt.Sprintf("DBHandler->CreateUser Insert WXPlatform :%s", err.Error())
		rollback(str, session)
		return -1
	}

	if referrer > 0 && referrer < user.Id {
		referrer_user := new(datastruct.UserInfo)
		has, _ := session.Where("id=?", referrer).Get(referrer_user)
		if has {
			sender := referrer_user.Id
			receiver := user.Id
			invite := new(datastruct.InviteInfo)
			has, _ = session.Where("sender=? and receiver=?", referrer_user.Id, receiver).Get(invite)
			if !has {
				invite.Receiver = receiver
				invite.Sender = sender
				invite.CreatedAt = user.CreatedAt
				_, err = session.Insert(invite)
				if err != nil {
					str := fmt.Sprintf("DBHandler->CreateUser Insert InviteInfo :%s", err.Error())
					rollback(str, session)
				}
			}
		}
	}

	err = session.Commit()
	if err != nil {
		str := fmt.Sprintf("DBHandler->CreateUser Commit :%s", err.Error())
		rollback(str, session)
		return -1
	}

	return user.Id
}

func rollback(err_str string, session *xorm.Session) {
	log.Debug("will rollback,err_str:%v", err_str)
	session.Rollback()
}

func rollbackError(err_str string, session *xorm.Session) {
	log.Error("will rollback,err_str:%v", err_str)
	session.Rollback()
}

//data
func (handle *DBHandler) GetHomeData(platform int) *datastruct.ResponseHomeData {
	engine := handle.mysqlEngine
	ad := make([]*datastruct.AdInfo, 0)
	err := engine.Where("(platform = ? or platform = ?) and is_hidden = 0", datastruct.H5+1, platform).Desc("sort_id").Find(&ad)
	if err != nil {
		log.Debug("GetHomeData err:%s", err.Error())
		return nil
	}
	resp_ads := make([]*datastruct.ResponseAdData, 0)
	for _, v := range ad {
		resp_ad := new(datastruct.ResponseAdData)
		resp_ad.ImgUrl = osstool.CreateOSSURL(v.ImgName)
		if v.Location == datastruct.DownLoadAd {
			resp_ad.IsJump = 1
			resp_ad.JumpTo = handle.GetDirectDownloadApp()
		}
		resp_ads = append(resp_ads, resp_ad)
	}
	resp := new(datastruct.ResponseHomeData)
	resp.Ad = resp_ads
	return resp
}

type goodsData struct {
	datastruct.Goods            `xorm:"extends"`
	datastruct.GoodsClass       `xorm:"extends"`
	datastruct.RecommendedClass `xorm:"extends"`
}

func (handle *DBHandler) GetGoods(pageIndex int, pageSize int, classid int, user_id int) []*datastruct.ResponseGoodsData {
	engine := handle.mysqlEngine
	goods := make([]*goodsData, 0)
	start := (pageIndex - 1) * pageSize
	limit := pageSize
	engine.Table("goods").Join("INNER", "goods_class", "goods_class.id = goods.goods_class_id").Join("Left", "recommended_class", "recommended_class.id = goods.re_classid").Where("goods_class.id=? and goods.is_hidden=?", classid, 0).Desc("goods.sort_id").Desc("goods.id").Limit(limit, start).Find(&goods)
	resp_goods := make([]*datastruct.ResponseGoodsData, 0, len(goods))
	avatarCount := 2
	sql := "select avatar from tmp_data td inner join tmp_data_for_goods tdfg on td.id = tdfg.tmp_user_id where tdfg.goods_id=? order by td.id desc limit ?"
	for _, v := range goods {
		resp_good := new(datastruct.ResponseGoodsData)
		resp_good.Id = v.Goods.Id
		results, _ := engine.Query(sql, v.Goods.Id, avatarCount)
		avatar := make([]string, 0, len(results))
		for _, v := range results {
			avatar = append(avatar, string(v["avatar"][:]))
		}
		resp_good.Avatar = avatar
		resp_good.ImgUrl = tools.CreateGoodsImgUrl(v.Goods.ImgName)
		resp_good.Name = v.Goods.Name
		if v.RecommendedClass.Icon != "" {
			resp_good.ReImgUrl = tools.CreateGoodsImgUrl(v.RecommendedClass.Icon)
		}
		resp_good.RushPrice = v.Goods.RushPrice
		resp_good.SendedOut = v.Goods.SendedOut
		resp_goods = append(resp_goods, resp_good)
	}
	if pageIndex == 1 {
		uggt := new(datastruct.UserGetHomeGoodsDataTime)
		has, err := engine.Where("user_id=? and class_id=? ", user_id, classid).Get(uggt)
		if err == nil {
			now_time := time.Now().Unix()
			if has {
				update_uggt := new(datastruct.UserGetHomeGoodsDataTime)
				update_uggt.GetDataTime = now_time
				engine.Where("id=?", uggt.Id).Cols("get_data_time").Update(update_uggt)
			} else {
				new_uggt := new(datastruct.UserGetHomeGoodsDataTime)
				new_uggt.ClassId = classid
				new_uggt.UserId = user_id
				new_uggt.GetDataTime = now_time
				engine.Insert(new_uggt)
			}
		}
	}
	return resp_goods
}

func (handle *DBHandler) GetGoodsClass(userId int) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	classinfo := make([]*datastruct.GoodsClass, 0)
	err := engine.Where("is_hidden=?", 0).Desc("sort_id").Asc("id").Find(&classinfo)
	if err != nil {
		return nil, datastruct.GetDataFailed
	}
	resp := make([]*datastruct.ResponseGoodsClass, 0)
	classids := make([]interface{}, 0, len(classinfo))
	for _, v := range classinfo {
		gclass := new(datastruct.ResponseGoodsClass)
		gclass.ImgUrl = tools.CreateGoodsImgUrl(v.Icon)
		gclass.Id = v.Id
		gclass.Name = v.Name
		resp = append(resp, gclass)
		classids = append(classids, v.Id)
	}
	uggdt := new(datastruct.UserGetHomeGoodsDataTime)
	engine.Where("user_id=?", userId).NotIn("class_id", classids...).Delete(uggdt)
	return resp, datastruct.NULLError
}

func (handle *DBHandler) WorldRewardInfo(pageIndex int, pageSize int) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	start := (pageIndex - 1) * pageSize
	limit := pageSize
	sql := "select p.nick_name,p.avatar,g.name,g.id from pay_mode_rouge_game_succeed_history p inner join goods g on p.goods_id = g.id LIMIT ?,?"
	rs, err := engine.Query(sql, start, limit)
	if err != nil {
		return nil, datastruct.GetDataFailed
	}

	resp := make([]*datastruct.ResponseRewardHistory, 0, len(rs))
	for _, v := range rs {
		his := new(datastruct.ResponseRewardHistory)
		his.Desc = tools.CreateUserDescInfo(string(v["nick_name"][:]), string(v["name"][:]))
		his.Avatar = string(v["avatar"][:])
		his.GoodsId = tools.StringToInt(string(v["id"][:]))
		resp = append(resp, his)
	}
	return resp, datastruct.NULLError
}

func (handle *DBHandler) FreeRougeGame() (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	game := make([]*datastruct.FreeModeRougeGame, 0)
	err := engine.Asc("level").Find(&game)
	if err != nil {
		return nil, datastruct.GetDataFailed
	}
	arr := make([]*datastruct.ResponseFreeRougeGame, 0, len(game))
	for _, v := range game {
		free := new(datastruct.ResponseFreeRougeGame)
		free.Level = v.Level
		free.GameTime = v.GameTime
		free.Difficulty = v.Difficulty
		free.RougeCount = v.RougeCount
		arr = append(arr, free)
	}
	return arr, datastruct.NULLError
}

func (handle *DBHandler) PayRougeGame(userId int, goodid int) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	user := new(datastruct.UserInfo)
	_, err := engine.Id(userId).Get(user)
	if err != nil {
		str := fmt.Sprintf("DBHandler->PayRougeGame GetUser :%s", err.Error())
		log.Debug(str)
		return nil, datastruct.GetDataFailed
	}

	rushSetting := new(datastruct.RushLimitSetting)
	var has bool
	has, err = engine.Where("id=?", datastruct.DefaultId).Get(rushSetting)
	if err != nil || !has {
		log.Debug("DBHandler->PayRougeGame Get RushLimitSetting err")
		return nil, datastruct.GetDataFailed
	}
	if user.LotterySucceed >= rushSetting.CheatCount {
		return rushSetting.CheatTips, datastruct.CheatUser
	}

	session := engine.NewSession()
	defer session.Close()
	session.Begin()

	savegame := new(datastruct.SaveGameInfo)
	var continueGame bool
	continueGame, err = session.Where("user_id = ?", userId).Get(savegame)
	if err != nil {
		str := fmt.Sprintf("DBHandler->PayRougeGame Get SaveGameInfo err:%s", err.Error())
		rollback(str, session)
		return nil, datastruct.NothasgameRecord
	}
	var gameRecord *datastruct.ResponsSaveGame
	if continueGame {
		goodid = savegame.GoodsId
	}
	goods := new(datastruct.Goods)
	_, err = session.Id(goodid).Get(goods)
	if err != nil {
		str := fmt.Sprintf("DBHandler->PayRougeGame GetGoods :%s", err.Error())
		rollback(str, session)
		return nil, datastruct.GetDataFailed
	}
	if !continueGame {
		_, err = session.Id(goodid).Get(goods)
		if err != nil {
			str := fmt.Sprintf("DBHandler->PayRougeGame GetGoods :%s", err.Error())
			rollback(str, session)
			return nil, datastruct.GetDataFailed
		}
		if user.GoldCount < goods.RushPrice {
			str := fmt.Sprintf("DBHandler->PayRougeGame GoldLess")
			rollback(str, session)
			return nil, datastruct.GoldLess
		}
		user.GoldCount -= goods.RushPrice
		user.PayRushTotal += goods.RushPrice
		var affected int64
		affected, err = session.Id(userId).Cols("gold_count", "pay_rush_total").Update(user)
		if err != nil || affected <= 0 {
			str := fmt.Sprintf("DBHandler->PayRougeGame UpdateUser err")
			rollback(str, session)
			return nil, datastruct.GetDataFailed
		}

		goldChangeInfo := new(datastruct.GoldChangeInfo)
		goldChangeInfo.UserId = userId
		goldChangeInfo.CreatedAt = time.Now().Unix()
		goldChangeInfo.VarGold = goods.RushPrice
		goldChangeInfo.ChangeType = datastruct.RushConsumeType
		affected, err = session.Insert(goldChangeInfo)
		if err != nil || affected <= 0 {
			str := fmt.Sprintf("DBHandler->PayRougeGame Insert GoldChangeInfo err")
			rollback(str, session)
			return nil, datastruct.GetDataFailed
		}
	}

	game := make([]*datastruct.PayModeRougeGame, 0)
	err = session.Where("goods_id=?", goodid).Asc("level").Find(&game)
	if err != nil || len(game) < datastruct.MaxLevel {
		str := fmt.Sprintf("DBHandler->PayRougeGame GetPayModeRougeGame err")
		rollback(str, session)
		return nil, datastruct.GetDataFailed
	}
	if user.LotterySucceed >= rushSetting.LotteryCount {
		game[1].Difficulty = rushSetting.Diff2
		game[1].RougeCount = rushSetting.Diff2r
		game[1].GameTime = rushSetting.Diff2t
		game[2].Difficulty = rushSetting.Diff3
		game[2].RougeCount = rushSetting.Diff3r
		game[2].GameTime = rushSetting.Diff3t
	}

	var tmp_price int64
	if !continueGame {
		savegame = new(datastruct.SaveGameInfo)
		savegame.UserId = userId
		savegame.LevelId = game[0].Id
		savegame.GoodsId = goodid
		savegame.CreatedAt = time.Now().Unix()
		_, err = session.Insert(savegame)
		if err != nil {
			str := fmt.Sprintf("DBHandler->PayRougeGame Insert SaveGameInfo err:%s", err.Error())
			rollback(str, session)
			return nil, datastruct.GetDataFailed
		}
		tmp_price = goods.RushPrice
	} else {
		tmp_price = 0
		gameRecord = new(datastruct.ResponsSaveGame)
		gameRecord.CurrentId = savegame.LevelId
		gameRecord.Goods.Name = goods.Name
		gameRecord.Goods.Id = goods.Id
		// gameRecord.Goods.Brand = goods.Brand
		// gameRecord.Goods.Desc = goods.GoodsDesc
		gameRecord.Goods.ImgUrl = tools.CreateGoodsImgUrl(goods.ImgName)
		// gameRecord.Goods.PriceDesc = goods.PriceDesc
		// gameRecord.Goods.RushPriceDesc = goods.RushPriceDesc
		gameRecord.Goods.RushPrice = goods.RushPrice
		//gameRecord.Goods.SellerRecommendedUrl = tools.CreateDZTJImgUrl()
	}

	rewardPool := new(datastruct.GoodsRewardPool)
	has, err = session.Where("goods_id = ?", goodid).Get(rewardPool)
	if err != nil || !has {
		rollback("DBHandler->PayRougeGame Get GoodsRewardPool err", session)
		return nil, datastruct.GetDataFailed
	}
	var isDisturb bool
	if rewardPool.Current+tmp_price >= rewardPool.LimitAmount {
		isDisturb = false
	} else {
		isDisturb = true
	}

	if tmp_price > 0 {
		var sql string
		sql = "update goods_reward_pool set current = current + ? where goods_id = ?"
		_, err = session.Exec(sql, tmp_price, goodid)
		if err != nil {
			rollback("DBHandler->PayRougeGame Update GoodsRewardPool err", session)
			return nil, datastruct.UpdateDataFailed
		}
	}
	err = session.Commit()
	if err != nil {
		str := fmt.Sprintf("DBHandler->PayRougeGame Commit :%s", err.Error())
		rollback(str, session)
		return nil, datastruct.GetDataFailed
	}
	resp := new(datastruct.ResponseGameContinue)
	resp.List = game
	resp.GameRecord = gameRecord
	resp.Label = tools.BoolToInt(isDisturb)
	return resp, datastruct.NULLError
}

func (handle *DBHandler) LevelPass(userId int, Id int) datastruct.CodeType {
	engine := handle.mysqlEngine
	session := engine.NewSession()
	defer session.Close()
	session.Begin()

	paygame, err := handle.checkGameInfo(session, userId, Id)
	if err != nil || paygame.Level >= datastruct.MaxLevel {
		rollback("LevelPass checkGameInfo err", session)
		return datastruct.UpdateDataFailed
	}

	nextLevel := paygame.Level + 1
	new_paygame := new(datastruct.PayModeRougeGame)
	var has bool
	has, err = session.Where("level=? and goods_id=?", nextLevel, paygame.GoodsId).Get(new_paygame)
	if err != nil || !has {
		str := fmt.Sprintf("DBHandler->LevelPass Get nextLevel PayModeRougeGame err")
		rollback(str, session)
		return datastruct.UpdateDataFailed
	}
	now_time := time.Now().Unix()
	sql := fmt.Sprintf("REPLACE INTO save_game_info (user_id,level_id,goods_id,created_at)VALUES(%d,%d,%d,%d)", userId, new_paygame.Id, paygame.GoodsId, now_time)
	_, err = session.Exec(sql)
	if err != nil {
		str := fmt.Sprintf("DBHandler->LevelPass REPLACE INTO save_game_info :%s", err.Error())
		rollback(str, session)
		return datastruct.UpdateDataFailed
	}

	err = session.Commit()
	if err != nil {
		str := fmt.Sprintf("DBHandler->LevelPass Commit :%s", err.Error())
		rollback(str, session)
		return datastruct.UpdateDataFailed
	}

	return datastruct.NULLError
}

func (handle *DBHandler) LevelPassFailed(userId int, Id int, Number int) datastruct.CodeType {
	engine := handle.mysqlEngine
	savegame := new(datastruct.SaveGameInfo)

	has, err := engine.Where("user_id=?", userId).Get(savegame)
	if err != nil || !has || savegame.LevelId != Id {
		log.Debug("LevelPassFailed Get SaveGameInfo err")
		return datastruct.GetDataFailed
	}

	_, err = engine.Where("user_id=?", userId).Delete(new(datastruct.SaveGameInfo))
	if err != nil {
		log.Debug("LevelPassFailed Delete SaveGameInfo err")
		return datastruct.UpdateDataFailed
	}

	failed := new(datastruct.PayModeRougeGameFailed)
	failed.UserId = userId
	failed.CreatedAt = time.Now().Unix()
	failed.PayModeRougeGameId = Id
	failed.RougeNumber = Number
	_, err = engine.Insert(failed)
	if err != nil {
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

func (handle *DBHandler) LevelPassSucceed(userId int, Id int, platform datastruct.Platform) datastruct.CodeType {
	engine := handle.mysqlEngine
	session := engine.NewSession()
	defer session.Close()
	session.Begin()

	now_time := time.Now().Unix()
	checkCheat := new(datastruct.PayModeRougeGameSucceed)
	has, err := session.Where("user_id=?", userId).Desc("created_at").Limit(1, 0).Get(checkCheat)
	if err != nil {
		str := fmt.Sprintf("DBHandler->LevelPassSucceed Get checkCheat :%s", err.Error())
		rollback(str, session)
		return datastruct.UpdateDataFailed
	}
	if has && now_time-checkCheat.CreatedAt <= datastruct.MinInterval {
		rollback("LevelPassSucceed checkCheat failed", session)
		return datastruct.FastPayModeSucceed
	}

	var paygame *datastruct.PayModeRougeGame
	paygame, err = handle.checkGameInfo(session, userId, Id)
	if err != nil || paygame.Level != datastruct.MaxLevel {
		rollback("LevelPassSucceed checkCheat failed", session)
		return datastruct.UpdateDataFailed
	}

	sql := "update user_info set lottery_succeed = lottery_succeed + 1 where id = ?"
	res, err1 := session.Exec(sql, userId)
	affected, err2 := res.RowsAffected()
	if err1 != nil || err2 != nil || affected <= 0 {
		rollback("DBHandler->LevelPassSucceed UpdateUser", session)
		return datastruct.UpdateDataFailed
	}

	_, err = session.Where("user_id=?", userId).Delete(new(datastruct.SaveGameInfo))
	if err != nil {
		str := fmt.Sprintf("DBHandler->LevelPassSucceed Delete SaveGameInfo :%s", err.Error())
		rollback(str, session)
		return datastruct.UpdateDataFailed
	}

	succeedRecord := new(datastruct.PayModeRougeGameSucceed)
	succeedRecord.UserId = userId
	succeedRecord.CreatedAt = now_time
	succeedRecord.PayModeRougeGameId = Id

	user := new(datastruct.UserInfo)
	_, err = session.Id(userId).Get(user)
	if err != nil {
		str := fmt.Sprintf("DBHandler->LevelPassSucceed GetUser :%s", err.Error())
		rollback(str, session)
		return datastruct.UpdateDataFailed
	}

	goods := new(datastruct.Goods)
	_, err = session.Id(paygame.GoodsId).Get(goods)
	if err != nil {
		str := fmt.Sprintf("DBHandler->LevelPassSucceed GetGoods :%s", err.Error())
		rollback(str, session)
		return datastruct.UpdateDataFailed
	}

	sql = "update goods_reward_pool set current = 0 where goods_id = ?"
	_, err = session.Exec(sql, paygame.GoodsId)
	if err != nil {
		rollback("DBHandler->LevelPassSucceed Update GoodsRewardPool err", session)
		return datastruct.UpdateDataFailed
	}

	history := new(datastruct.PayModeRougeGameSucceedHistory)
	history.NickName = user.NickName
	history.Avatar = user.Avatar
	history.GoodsId = goods.Id
	_, err = session.Insert(succeedRecord, history)
	if err != nil {
		str := fmt.Sprintf("DBHandler->LevelPassSucceed insert PayModeRougeGameSucceed and PayModeRougeGameSucceedHistory:%s", err.Error())
		rollback(str, session)
		return datastruct.UpdateDataFailed
	}

	//创建订单
	err = createOrderData(session, userId, paygame.GoodsId, now_time, false, platform)
	if err != nil {
		rollback(err.Error(), session)
		return datastruct.UpdateDataFailed
	}

	err = session.Commit()
	if err != nil {
		str := fmt.Sprintf("DBHandler->LevelPassSucceed Commit :%s", err.Error())
		rollback(str, session)
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

func (handle *DBHandler) checkGameInfo(session *xorm.Session, userId int, Id int) (*datastruct.PayModeRougeGame, error) {
	paygame := new(datastruct.PayModeRougeGame)
	has, err := session.Id(Id).Get(paygame)
	if err != nil || !has {
		err_str := fmt.Sprintf("DBHandler->LevelPass GetPayModeRougeGame :%s", err.Error())
		rs_err := errors.New(err_str)
		return nil, rs_err
	}

	savegame := new(datastruct.SaveGameInfo)
	has, err = session.Where("user_id = ? and level_id = ?", userId, Id).Get(savegame)
	if err != nil || !has {
		err_str := fmt.Sprintf("DBHandler->LevelPass GetSaveGameInfo err")
		rs_err := errors.New(err_str)
		return nil, rs_err
	}
	return paygame, nil
}

func createOrderData(session *xorm.Session, userId int, GoodsId int, nowTime int64, isPurchase bool, platform datastruct.Platform) error {
	orderInfo := new(datastruct.OrderInfo)
	orderInfo.CreatedAt = nowTime
	orderInfo.UserId = userId
	orderInfo.Number = tools.UniqueId()
	orderInfo.GoodsId = GoodsId
	orderInfo.OrderState = datastruct.NotApply
	orderInfo.IsPurchase = tools.BoolToInt(isPurchase)
	orderInfo.Platform = platform
	_, err := session.Insert(orderInfo)
	if err != nil {
		err_str := fmt.Sprintf("createOrderData err:%s", err.Error())
		rs_err := errors.New(err_str)
		return rs_err
	}
	//库存减1
	return nil
}

//购买成功
func (handle *DBHandler) PurchaseSuccess(userId int, goodsid int, createTime int64, platform datastruct.Platform) datastruct.CodeType {
	engine := handle.mysqlEngine
	session := engine.NewSession()
	defer session.Close()
	session.Begin()

	goods := new(datastruct.Goods)
	_, err := session.Id(goodsid).Get(goods)
	if err != nil {
		str := fmt.Sprintf("DBHandler->Purchase GetGoods :%s", err.Error())
		rollbackError(str, session)
		return datastruct.GetDataFailed
	}

	sql := "update user_info set purchase_total = purchase_total + ? where id = ?"
	res, err1 := session.Exec(sql, float64(goods.Price), userId)
	affected, err2 := res.RowsAffected()
	if err1 != nil || err2 != nil || affected <= 0 {
		rollbackError("DBHandler->PurchaseSuccess UpdateUser", session)
		return datastruct.UpdateDataFailed
	}

	//创建订单
	err = createOrderData(session, userId, goodsid, createTime, true, platform)
	if err != nil {
		rollbackError(err.Error(), session)
		return datastruct.GetDataFailed
	}

	err = session.Commit()
	if err != nil {
		str := fmt.Sprintf("DBHandler->Purchase Commit :%s", err.Error())
		rollbackError(str, session)
		return datastruct.GetDataFailed
	}

	return datastruct.NULLError
}

func (handle *DBHandler) GetGoodsPrice(goodsid int) (*datastruct.Goods, datastruct.CodeType) {
	engine := handle.mysqlEngine
	goods := new(datastruct.Goods)
	has, err := engine.Id(goodsid).Get(goods)
	if err != nil || !has {
		log.Error("DBHandler->GetGoodsPrice err")
		return nil, datastruct.GetDataFailed
	}
	return goods, datastruct.NULLError
}

type orderData struct {
	datastruct.Goods     `xorm:"extends"`
	datastruct.OrderInfo `xorm:"extends"`
}

func (handle *DBHandler) GetNotAppliedOrderInfo(userId int, pageIndex int, pageSize int) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	start := (pageIndex - 1) * pageSize
	limit := pageSize
	orderInfo := make([]*orderData, 0)
	err := engine.Table("goods").Join("INNER", "order_info", "order_info.goods_id = goods.id").Where("user_id = ? and order_info.order_state = 0", userId).Desc("order_info.created_at").Limit(limit, start).Find(&orderInfo)
	if err != nil {
		return nil, datastruct.GetDataFailed
	}
	resp := make([]*datastruct.ResponseOrderInfo, 0)
	for _, v := range orderInfo {
		rs := new(datastruct.ResponseOrderInfo)
		rs.Id = v.GoodsId
		rs.ImgUrl = tools.CreateGoodsImgUrl(v.Goods.ImgName)
		rs.Name = v.Goods.Name
		if v.OrderInfo.IsPurchase == 0 {
			rs.Price = v.Goods.RushPrice
			rs.Remark = "已挑战成功"
		} else {
			rs.Price = v.Goods.Price
			rs.Remark = "买家已付款"
		}
		rs.PriceDesc = v.Goods.PriceDesc
		rs.Count = 1
		rs.Number = v.OrderInfo.Number
		resp = append(resp, rs)
	}
	return resp, datastruct.NULLError
}

func (handle *DBHandler) GetNotSendGoods(userId int, pageIndex int, pageSize int) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	start := (pageIndex - 1) * pageSize
	limit := pageSize
	sql := "select o.is_purchase,g.id as gid,g.name as gname,g.img_name,g.price,g.rush_price,g.price_desc,o.number,o.is_remind from send_goods s inner join order_info o on o.id = s.order_id inner join goods g on g.id = o.goods_id where s.send_goods_state = 0 and o.user_id = ? ORDER BY s.created_at desc LIMIT ?,?"
	results, err := engine.Query(sql, userId, start, limit)
	if err != nil {
		log.Debug("GetNotSendGoods err:%v", err.Error())
		return nil, datastruct.GetDataFailed
	}
	resp := make([]*datastruct.ResponseNotSendGoods, 0)
	for _, v := range results {
		rs := new(datastruct.ResponseNotSendGoods)
		rs.Count = 1
		rs.Id = tools.StringToInt(string(v["gid"][:]))
		rs.ImgUrl = tools.CreateGoodsImgUrl(string(v["img_name"][:]))
		rs.IsRemind = tools.StringToInt(string(v["is_remind"][:]))
		rs.Name = string(v["gname"][:])
		rs.Number = string(v["number"][:])
		is_purchase := tools.StringToInt(string(v["is_purchase"][:]))
		if is_purchase == 0 {
			rs.Price = tools.StringToInt64(string(v["rush_price"][:]))
		} else {
			rs.Price = tools.StringToInt64(string(v["price"][:]))
		}
		rs.PriceDesc = string(v["price_desc"][:])
		rs.Remark = "等待发货"
		resp = append(resp, rs)
	}
	return resp, datastruct.NULLError
}

func (handle *DBHandler) GetHasSendedGoods(userId int, pageIndex int, pageSize int) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	start := (pageIndex - 1) * pageSize
	limit := pageSize
	sql := "select o.is_purchase,g.id as gid,g.name as gname,g.img_name,g.price,g.rush_price,g.price_desc,o.number,s.express_number,s.express_agency from send_goods s inner join order_info o on o.id = s.order_id inner join goods g on g.id = o.goods_id where s.send_goods_state = 1 and s.sign_for_state = 0 and o.user_id = ? ORDER BY s.created_at desc LIMIT ?,?"
	results, err := engine.Query(sql, userId, start, limit)
	if err != nil {
		log.Debug("GetNotSendGoods err:%v", err.Error())
		return nil, datastruct.GetDataFailed
	}
	resp := make([]*datastruct.ResponseHasSendedGoods, 0)
	for _, v := range results {
		rs := new(datastruct.ResponseHasSendedGoods)
		rs.Count = 1
		rs.Id = tools.StringToInt(string(v["gid"][:]))
		rs.ImgUrl = tools.CreateGoodsImgUrl(string(v["img_name"][:]))
		rs.Name = string(v["gname"][:])
		rs.Number = string(v["number"][:])
		is_purchase := tools.StringToInt(string(v["is_purchase"][:]))
		if is_purchase == 0 {
			rs.Price = tools.StringToInt64(string(v["rush_price"][:]))
		} else {
			rs.Price = tools.StringToInt64(string(v["price"][:]))
		}
		rs.PriceDesc = string(v["price_desc"][:])
		rs.Remark = "已发货"
		rs.ExpressNumber = string(v["express_number"][:])
		rs.ExpressAgency = string(v["express_agency"][:])
		resp = append(resp, rs)
	}
	return resp, datastruct.NULLError
}

func (handle *DBHandler) GetAppraiseOrder(userId int, pageIndex int, pageSize int) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	start := (pageIndex - 1) * pageSize
	limit := pageSize
	sql := "select o.is_purchase,g.id as gid,g.name as gname,g.img_name,g.price,g.rush_price,g.price_desc,o.number,s.is_appraised from send_goods s inner join order_info o on o.id = s.order_id inner join goods g on g.id = o.goods_id where s.send_goods_state = 1 and s.sign_for_state = 1 and o.user_id = ? ORDER BY s.created_at desc LIMIT ?,?"
	results, err := engine.Query(sql, userId, start, limit)
	if err != nil {
		log.Debug("GetNotSendGoods err:%v", err.Error())
		return nil, datastruct.GetDataFailed
	}
	resp := make([]*datastruct.ResponseAppraiseOrder, 0)
	for _, v := range results {
		rs := new(datastruct.ResponseAppraiseOrder)
		rs.Count = 1
		rs.Id = tools.StringToInt(string(v["gid"][:]))
		rs.ImgUrl = tools.CreateGoodsImgUrl(string(v["img_name"][:]))
		rs.Name = string(v["gname"][:])
		rs.Number = string(v["number"][:])
		is_purchase := tools.StringToInt(string(v["is_purchase"][:]))
		if is_purchase == 0 {
			rs.Price = tools.StringToInt64(string(v["rush_price"][:]))
		} else {
			rs.Price = tools.StringToInt64(string(v["price"][:]))
		}
		rs.PriceDesc = string(v["price_desc"][:])
		rs.Remark = "已签收"
		rs.IsAppraised = tools.StringToInt(string(v["is_appraised"][:]))
		resp = append(resp, rs)
	}
	return resp, datastruct.NULLError
}

func (handle *DBHandler) ApplySend(userId int, body *datastruct.ApplySendBody) datastruct.CodeType {
	engine := handle.mysqlEngine
	session := engine.NewSession()
	defer session.Close()
	session.Begin()
	orderInfo := new(datastruct.OrderInfo)
	has, err := session.Where("number=?", body.OrderNumber).Get(orderInfo)
	if err != nil || !has || orderInfo.OrderState == datastruct.Apply {
		rollback("ApplySend get order err", session)
		return datastruct.UpdateDataFailed
	}
	sendGoods := new(datastruct.SendGoods)
	has, _ = session.Where("order_id=?", orderInfo.Id).Get(sendGoods)
	if has {
		rollback("ApplySend err:has sendGoods ", session)
		return datastruct.UpdateDataFailed
	}
	var affected int64
	orderInfo.OrderState = datastruct.Apply
	affected, err = session.Where("id=?", orderInfo.Id).Cols("order_state").Update(orderInfo)
	if err != nil || affected <= 0 {
		rollback("ApplySend update orderInfo err", session)
		return datastruct.UpdateDataFailed
	}
	sendGoods.Address = body.Address
	sendGoods.CreatedAt = time.Now().Unix()
	sendGoods.LinkMan = body.LinkMan
	sendGoods.OrderId = int64(orderInfo.Id)
	sendGoods.PhoneNumber = body.PhoneNumber
	sendGoods.Remark = body.Remark
	sendGoods.SendGoodsState = datastruct.NotSend

	affected, err = session.Insert(sendGoods)
	if err != nil || affected <= 0 {
		rollback("ApplySend Insert sendGoods err", session)
		return datastruct.UpdateDataFailed
	}

	sql := fmt.Sprintf("REPLACE INTO user_shipping_address (user_id,linkman,phone,addr,remark)VALUES(%d,'%s','%s','%s','%s')", userId, body.LinkMan, body.PhoneNumber, body.Address, body.Remark)
	_, err = session.Exec(sql)
	if err != nil {
		str := fmt.Sprintf("DBHandler->ApplySend REPLACE INTO user_shipping_address :%s", err.Error())
		rollback(str, session)
		return datastruct.UpdateDataFailed
	}

	err = session.Commit()
	if err != nil {
		str := fmt.Sprintf("DBHandler->ApplySend Commit :%s", err.Error())
		rollback(str, session)
		return datastruct.GetDataFailed
	}
	return datastruct.NULLError
}

type updateSendData struct {
	datastruct.SendGoods `xorm:"extends"`
	datastruct.OrderInfo `xorm:"extends"`
}

//充值成功
func (handle *DBHandler) DepositSucceed(userId int, money int64, platform datastruct.Platform) (int64, datastruct.CodeType) {
	engine := handle.mysqlEngine
	session := engine.NewSession()
	defer session.Close()
	session.Begin()
	user := new(datastruct.UserInfo)
	has, err := session.Where("id = ?", userId).Get(user)
	if err != nil || !has {
		str := fmt.Sprintf("DBHandler->DepositSucceed Get UserInfo err")
		rollbackError(str, session)
		return -1, datastruct.DepositFailed
	}
	user.GoldCount += money
	user.DepositTotal += money
	affected, err := session.Id(userId).Cols("gold_count", "deposit_total").Update(user)
	if err != nil || affected <= 0 {
		str := fmt.Sprintf("DBHandler->DepositSucceed Update Gold err")
		rollbackError(str, session)
		return -1, datastruct.DepositFailed
	}

	now_time := time.Now().Unix()
	goldChangeInfo := new(datastruct.GoldChangeInfo)
	goldChangeInfo.UserId = userId
	goldChangeInfo.CreatedAt = now_time
	goldChangeInfo.VarGold = money
	goldChangeInfo.ChangeType = datastruct.DepositType
	_, err = session.Insert(goldChangeInfo)
	if err != nil {
		str := fmt.Sprintf("DBHandler->DepositSucceed Insert GoldChangeInfo :%s", err.Error())
		rollbackError(str, session)
		return -1, datastruct.DepositFailed
	}

	userDepositInfo := new(datastruct.UserDepositInfo)
	userDepositInfo.CreatedAt = now_time
	userDepositInfo.UserId = userId
	userDepositInfo.Money = float64(money)
	userDepositInfo.PayPlatform = datastruct.WXPay
	userDepositInfo.Platform = platform
	_, err = session.Insert(userDepositInfo)
	if err != nil {
		str := fmt.Sprintf("DBHandler->DepositSucceed Insert UserDepositInfo :%s", err.Error())
		rollbackError(str, session)
		return -1, datastruct.DepositFailed
	}

	var receiver int
	receiver, err = agencyEarn(userDepositInfo.Id, userId, datastruct.AgentLevel1, money, userId, now_time, session)
	if err != nil {
		str := fmt.Sprintf("DBHandler->DepositSucceed Insert DepositInfo level_1 :%s", err.Error())
		rollbackError(str, session)
		return -1, datastruct.DepositFailed
	}
	if receiver > 0 {
		receiver, err = agencyEarn(userDepositInfo.Id, receiver, datastruct.AgentLevel2, money, userId, now_time, session)
		if err != nil {
			str := fmt.Sprintf("DBHandler->DepositSucceed Insert DepositInfo level_2:%s", err.Error())
			rollbackError(str, session)
			return -1, datastruct.DepositFailed
		}
		if receiver > 0 {
			receiver, err = agencyEarn(userDepositInfo.Id, receiver, datastruct.AgentLevel3, money, userId, now_time, session)
			if err != nil {
				str := fmt.Sprintf("DBHandler->DepositSucceed Insert DepositInfo level_3:%s", err.Error())
				rollbackError(str, session)
				return -1, datastruct.DepositFailed
			}
		}
	}
	err = session.Commit()
	if err != nil {
		str := fmt.Sprintf("DBHandler->DepositSucceed Commit :%s", err.Error())
		rollbackError(str, session)
		return -1, datastruct.DepositFailed
	}
	return user.GoldCount, datastruct.NULLError
}

func agencyEarn(depositId int, receiver int, agencyLevel datastruct.AgentLevelType, money int64, fromUserId int, nowTime int64, session *xorm.Session) (int, error) {
	inviteInfo := new(datastruct.InviteInfo)
	has, err := session.Where("receiver = ?", receiver).Get(inviteInfo)
	if err != nil {
		err_str := fmt.Sprintf("DBHandler->DepositSucceed Get InviteInfo:%s", err.Error())
		rs_err := errors.New(err_str)
		return -1, rs_err
	}
	if has {
		sender := new(datastruct.UserInfo)
		has, err = session.Where("id = ?", inviteInfo.Sender).Get(sender)
		if err != nil {
			err_str := fmt.Sprintf("DBHandler->agencyEarn Get UserInfo err:%s", err.Error())
			rs_err := errors.New(err_str)
			return -1, rs_err
		}
		if !has {
			log.Debug("DBHandler->agencyEarn not has userid:%v", inviteInfo.Sender)
			return 0, nil
		}
		agencyIdentifier := sender.MemberIdentifier
		agencyParams := new(datastruct.AgencyParams)
		has, err = session.Where("identifier = ?", agencyIdentifier).Get(agencyParams)
		if err != nil {
			err_str := fmt.Sprintf("DBHandler->DepositSucceed Get AgencyParams identifier:%d, err:%s", inviteInfo.Sender, err.Error())
			rs_err := errors.New(err_str)
			return -1, rs_err
		}
		if !has {
			if agencyIdentifier == datastruct.AgencyIdentifier {
				err_str := fmt.Sprintf("DBHandler->DepositSucceed Get AgencyParams err: not has identifier :%s", datastruct.AgencyIdentifier)
				rs_err := errors.New(err_str)
				return -1, rs_err
			} else {
				has, err = session.Where("identifier = ?", datastruct.AgencyIdentifier).Get(agencyParams)
				if err != nil {
					err_str := fmt.Sprintf("DBHandler->DepositSucceed Common AgencyParams sender:%d, err:%s", inviteInfo.Sender, err.Error())
					rs_err := errors.New(err_str)
					return -1, rs_err
				}
			}
		}

		var perMoney float64
		var perGold float64
		var addGold_int64 int64
		addGold_int64 = 0
		switch agencyLevel {
		case 1:
			perMoney = float64(agencyParams.Agency1MoneyPercent)
			perGold = float64(agencyParams.Agency1GoldPercent)
		case 2:
			perMoney = float64(agencyParams.Agency2MoneyPercent)
			perGold = float64(agencyParams.Agency2GoldPercent)
		case 3:
			perMoney = float64(agencyParams.Agency3MoneyPercent)
			perGold = float64(agencyParams.Agency3GoldPercent)
		}
		if perGold != 0 {
			addGold := tools.Decimal2(perGold * float64(money) / 100.0)
			addGold_int64 = int64(addGold)
			if addGold > float64(addGold_int64) {
				addGold_int64 = addGold_int64 + 1
			}
		}
		if perMoney != 0 {
			agency := new(datastruct.BalanceInfo)
			agency.DepositId = depositId
			agency.FromUserId = fromUserId
			agency.ToUserId = inviteInfo.Sender
			agency.AgencyLevel = int8(agencyLevel)
			agency.EarnBalance = tools.Decimal2(perMoney * float64(money) / 100.0)
			agency.EarnGold = addGold_int64
			agency.CreatedAt = nowTime
			_, err = session.Insert(agency)
			if err != nil {
				err_str := fmt.Sprintf("DBHandler->DepositSucceed Insert BalanceInfo :%s", err.Error())
				rs_err := errors.New(err_str)
				return -1, rs_err
			}
			sql := "update user_info set balance = balance + ? , balance_total = balance_total + ? where id = ?"
			_, err1 := session.Exec(sql, agency.EarnBalance, agency.EarnBalance, inviteInfo.Sender)
			if err1 != nil {
				rs_err := errors.New("DBHandler->DepositSucceed Add balance err:" + err1.Error())
				return -1, rs_err
			}
		}
		if addGold_int64 != 0 {
			goldChangeInfo := new(datastruct.GoldChangeInfo)
			goldChangeInfo.ChangeType = datastruct.ProxyRewardType
			goldChangeInfo.CreatedAt = nowTime
			goldChangeInfo.UserId = inviteInfo.Sender
			goldChangeInfo.VarGold = addGold_int64
			_, err = session.Insert(goldChangeInfo)
			if err != nil {
				err_str := fmt.Sprintf("DBHandler->DepositSucceed Insert GoldChangeInfo :%s", err.Error())
				rs_err := errors.New(err_str)
				return -1, rs_err
			}
			sql := "update user_info set gold_count = gold_count + ? where id = ?"
			_, err1 := session.Exec(sql, addGold_int64, inviteInfo.Sender)
			if err1 != nil {
				rs_err := errors.New("DBHandler->DepositSucceed Add Agency Gold err:" + err1.Error())
				return -1, rs_err
			}
		}
	} else {
		inviteInfo.Sender = 0
	}
	return inviteInfo.Sender, nil
}

func (handle *DBHandler) CommissionRank(userId int, pageIndex int, pageSize int) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	users := make([]*datastruct.UserInfo, 0)
	start := (pageIndex - 1) * pageSize
	limit := pageSize
	columnName := "balance_total"
	err := engine.Desc("balance_total").Asc("id").Limit(limit, start).Find(&users)
	if err != nil {
		return nil, datastruct.GetDataFailed
	}
	resp := new(datastruct.ResponseCommissionRank)
	list := make([]*datastruct.ResponseUserCommission, 0)
	for _, v := range users {
		resp_comm := new(datastruct.ResponseUserCommission)
		resp_comm.Avatar = v.Avatar
		resp_comm.NickName = v.NickName
		resp_comm.Total = v.BalanceTotal
		list = append(list, resp_comm)
	}

	self := new(datastruct.ResponseSelfCommission)
	table := "user_info"
	sql_str := fmt.Sprintf("SELECT b.* FROM ( SELECT t.*, @rownum := @rownum + 1 AS rownum FROM (SELECT @rownum := 0) r, (SELECT * FROM %s ORDER BY %s DESC) AS t ) AS b WHERE b.id = %d", table, columnName, userId)
	var results []map[string][]byte
	results, err = engine.Query(sql_str)
	if err != nil || len(results) <= 0 {
		return nil, datastruct.GetDataFailed
	}
	total_byte := results[0][columnName]
	self.Base.Total = tools.StringToFloat64(string(total_byte[:]))
	self.Base.Avatar = string(results[0]["avatar"][:])
	self.Base.NickName = string(results[0]["nick_name"][:])
	self.Rank = tools.StringToInt(string(results[0]["rownum"][:]))

	resp.List = list
	resp.Self = self
	return resp, datastruct.NULLError
}

func (handle *DBHandler) CommissionInfo(userId int, pageIndex int, pageSize int) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	arr := make([]*datastruct.BalanceInfo, 0)
	start := (pageIndex - 1) * pageSize
	limit := pageSize
	err := engine.Where("to_user_id=?", userId).Desc("created_at").Limit(limit, start).Find(&arr)
	if err != nil {
		return nil, datastruct.GetDataFailed
	}
	resp := make([]*datastruct.ResponseCommissionInfo, 0)
	for _, v := range arr {
		info := new(datastruct.ResponseCommissionInfo)
		info.AgencyLevel = v.AgencyLevel
		info.CreatedAt = v.CreatedAt
		info.EarnBalance = v.EarnBalance
		resp = append(resp, info)
	}
	return resp, datastruct.NULLError
}

func (handle *DBHandler) GetAgentCount(userId int) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	max_level := 3
	total_arr := make([]int, 0, max_level)
	for i := 1; i <= max_level; i++ {
		var sql string
		switch i {
		case 1:
			sql = "select count(*) from invite_info i join user_info u on i.receiver=u.id where i.sender = ?"
		case 2:
			sql = "select count(*) from invite_info i join user_info u on i.receiver=u.id where i.sender in (select u.id from invite_info i join user_info u on i.receiver=u.id where i.sender = ?)"
		case 3:
			sql = "select count(*) from invite_info i join user_info u on i.receiver=u.id where i.sender in (select u.id from invite_info i join user_info u on i.receiver=u.id where i.sender in (select u.id from invite_info i join user_info u on i.receiver=u.id where i.sender = ?))"
		}
		results, err := engine.Query(sql, userId)
		if err != nil {
			return nil, datastruct.GetDataFailed
		}
		total_str := string(results[0]["count(*)"][:])
		total_arr = append(total_arr, tools.StringToInt(total_str))
	}
	return total_arr, datastruct.NULLError
}

func (handle *DBHandler) GetAgentlevelN(userId int, level int, pageIndex int, pageSize int) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	start := (pageIndex - 1) * pageSize
	limit := pageSize

	var sql string
	switch level {
	case 1:
		sql = "select u.id,u.avatar,u.nick_name,u.created_at from invite_info i join user_info u on i.receiver=u.id where i.sender = ? order by i.created_at desc LIMIT ?,?"
	case 2:
		sql = "select u.id,u.avatar,u.nick_name,u.created_at from invite_info i join user_info u on i.receiver=u.id where i.sender in (select u.id from invite_info i join user_info u on i.receiver=u.id where i.sender = ?) order by i.created_at desc LIMIT ?,?"
	case 3:
		sql = "select u.id,u.avatar,u.nick_name,u.created_at from invite_info i join user_info u on i.receiver=u.id where i.sender in (select u.id from invite_info i join user_info u on i.receiver=u.id where i.sender in (select u.id from invite_info i join user_info u on i.receiver=u.id where i.sender = ?)) order by i.created_at desc LIMIT ?,?"
	}
	results, err := engine.Query(sql, userId, start, limit)

	if err != nil {
		log.Error("getAgentlevel_%d join :%s", level, err.Error())
		return nil, datastruct.GetDataFailed
	}
	resp := new(datastruct.ResponseAgentInfo)
	list := make([]*datastruct.ResponseAgent, 0)
	for _, v := range results {
		agent := new(datastruct.ResponseAgent)
		agent.Avatar = string(v["avatar"][:])
		agent.CreatedAt = tools.StringToInt64(string(v["created_at"][:]))
		agent.NickName = string(v["nick_name"][:])
		receiver_id := tools.StringToInt(string(v["id"][:]))
		balanceInfo := new(datastruct.BalanceInfo)
		total, err := engine.Where("from_user_id = ? and to_user_id = ?", receiver_id, userId).Sum(balanceInfo, "earn_balance")
		if err != nil {
			log.Error("getAgentlevel_%d Sum :%s", level, err.Error())
			return nil, datastruct.GetDataFailed
		}
		agent.EarnBalance = total
		list = append(list, agent)
	}
	user := new(datastruct.UserInfo)
	var has bool
	has, err = engine.Where("id=?", userId).Get(user)
	if err != nil || !has {
		return nil, datastruct.GetDataFailed
	}
	resp.Total = user.BalanceTotal
	resp.Balance = user.Balance
	resp.Agent = list
	return resp, datastruct.NULLError
}

func (handle *DBHandler) GetDrawCash(userId int, pageIndex int, pageSize int) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	arr := make([]*datastruct.DrawCashInfo, 0)
	start := (pageIndex - 1) * pageSize
	limit := pageSize
	err := engine.Where("user_id=?", userId).Desc("created_at").Limit(limit, start).Find(&arr)
	if err != nil {
		return nil, datastruct.GetDataFailed
	}
	resp := make([]*datastruct.ResponseDrawCashInfo, 0)
	for _, v := range arr {
		info := new(datastruct.ResponseDrawCashInfo)
		info.Charge = v.Charge
		if v.PaymentTime == "" {
			info.PaymentTime = tools.UnixToString(v.CreatedAt, "2006-01-02 15:04:05")
		} else {
			info.PaymentTime = v.PaymentTime
		}
		info.Poundage = v.Poundage
		info.State = v.State
		switch v.ArrivalType {
		case datastruct.DrawCashArrivalWX:
			info.ArrivalType = "微信钱包"
		case datastruct.DrawCashArrivalZFB:
			info.ArrivalType = "支付宝"
		}
		resp = append(resp, info)
	}
	return resp, datastruct.NULLError
}

func (handle *DBHandler) GetDrawCashRule(userId int) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	user := new(datastruct.UserInfo)
	has, err := engine.Where("id=?", userId).Get(user)
	if err != nil || !has {
		return nil, datastruct.GetDataFailed
	}
	params := new(datastruct.DrawCashParams)
	has, err = engine.Where("id=?", datastruct.DefaultId).Get(params)
	if err != nil || !has {
		return nil, datastruct.GetDataFailed
	}
	resp := new(datastruct.ResponseDrawCashRule)
	resp.Balance = user.Balance
	resp.MaxDrawCount = params.MaxDrawCount
	resp.MinCharge = params.MinCharge
	resp.MinPoundage = params.MinPoundage
	resp.PoundagePer = params.PoundagePer
	resp.RequireVerify = params.RequireVerify
	return resp, datastruct.NULLError
}

func (handle *DBHandler) GetDepositParams(userId int) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	user := new(datastruct.UserInfo)
	has, err := engine.Where("id=?", userId).Get(user)
	if err != nil || !has {
		return nil, datastruct.GetDataFailed
	}
	params := make([]*datastruct.DepositParams, 0)
	err = engine.Find(&params)
	if err != nil || len(params) <= 0 {
		return nil, datastruct.GetDataFailed
	}
	resp := new(datastruct.ResponseDepositParams)
	resp.GoldCount = user.GoldCount
	resp.List = params
	return resp, datastruct.NULLError
}

func (handle *DBHandler) GetUserInfo(userId int, platform datastruct.Platform) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	user := new(datastruct.UserInfo)
	has, err := engine.Where("id=?", userId).Get(user)
	if err != nil || !has {
		log.Error("GetUserInfo GetUser err")
		return nil, datastruct.GetDataFailed
	}
	resp := new(datastruct.ResponseUserInfo)
	resp.Avatar = user.Avatar
	resp.GoldCount = user.GoldCount
	resp.NickName = user.NickName
	resp.UserId = user.Id

	var DLQCount, DFHCount, YFHCount, DPJCount int64

	orderinfo := new(datastruct.OrderInfo)
	DLQCount, err = engine.Where("user_id = ? and order_state = 0", userId).Count(orderinfo)
	if err != nil {
		log.Error("GetUserInfo get DLQCount err:%v", err.Error())
		return nil, datastruct.GetDataFailed
	}
	sendgoods := new(datastruct.SendGoods)
	DFHCount, err = engine.Join("INNER", "order_info", "order_info.id = send_goods.order_id").Where("user_id=? and send_goods_state = 0", userId).Count(sendgoods)
	if err != nil {
		log.Error("GetUserInfo get DFHCount err:%v", err.Error())
		return nil, datastruct.GetDataFailed
	}

	YFHCount, err = engine.Join("INNER", "order_info", "order_info.id = send_goods.order_id").Where("user_id=? and send_goods_state = 1 and sign_for_state = 0", userId).Count(sendgoods)
	if err != nil {
		log.Error("GetUserInfo get YFHCount err:%v", err.Error())
		return nil, datastruct.GetDataFailed
	}

	DPJCount, err = engine.Join("INNER", "order_info", "order_info.id = send_goods.order_id").Where("user_id=? and send_goods_state = 1 and sign_for_state = 1 and is_appraised = 0", userId).Count(sendgoods)
	if err != nil {
		log.Error("GetUserInfo get DPJCount err:%v", err.Error())
		return nil, datastruct.GetDataFailed
	}

	adInfo := new(datastruct.AdInfo)
	has, err = engine.Where("location=?", datastruct.AppraiseAd).Get(adInfo)
	if err != nil {
		log.Error("GetUserInfo get AdInfo err:%v", err.Error())
		return nil, datastruct.GetDataFailed
	}

	gcg := new(datastruct.GoldCoinGift)
	has, err = engine.Where("id=?", datastruct.DefaultId).Get(gcg)
	if err != nil || !has {
		log.Error("GetUserInfo get GoldCoinGift err")
		return nil, datastruct.GetDataFailed
	}

	appraiseAd := new(datastruct.AppraiseAdInfo)
	appraiseAd.AppraiseAd = osstool.CreateOSSURL(adInfo.ImgName)
	appraiseAd.GoldCount = gcg.AppraisedGoldGift

	registerGift := new(datastruct.RegisterGift)
	if gcg.IsEnableRegisterGift == 0 {
		registerGift.IsGotRegisterGift = 1
	} else {
		registerGift.IsGotRegisterGift = user.IsGotRegisterGift
	}
	registerGift.GoldCount = gcg.RegisterGoldGift

	resp.DLQCount = DLQCount
	resp.DFHCount = DFHCount
	resp.YFHCount = YFHCount
	resp.DPJCount = DPJCount

	resp.AppraiseAd = appraiseAd
	resp.RegisterGift = registerGift
	return resp, datastruct.NULLError
}

func (handle *DBHandler) GetGoldInfo(userId int, pageIndex int, pageSize int) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	arr := make([]*datastruct.GoldChangeInfo, 0)
	start := (pageIndex - 1) * pageSize
	limit := pageSize
	err := engine.Where("user_id=?", userId).Desc("created_at").Limit(limit, start).Find(&arr)
	if err != nil {
		return nil, datastruct.NULLError
	}
	user := new(datastruct.UserInfo)
	var has bool
	has, err = engine.Where("id=?", userId).Get(user)
	if err != nil || !has {
		log.Debug("GetGoldInfo GetUser err")
		return nil, datastruct.GetDataFailed
	}
	resp := new(datastruct.ResponsGoldData)
	list := make([]*datastruct.ResponsGoldChange, 0, len(arr))
	for _, v := range arr {
		resp := new(datastruct.ResponsGoldChange)
		resp.ChangeType = v.ChangeType
		resp.CreatedAt = v.CreatedAt
		resp.GoldChange = v.VarGold
		list = append(list, resp)
	}
	resp.GoldCount = user.GoldCount
	resp.List = list
	return resp, datastruct.NULLError
}
func (handle *DBHandler) GetOpenId(userId int) (*datastruct.WXPlatform, datastruct.CodeType) {
	wx_user := new(datastruct.WXPlatform)
	engine := handle.mysqlEngine
	has, err := engine.Where("user_id=?", userId).Get(wx_user)
	if err != nil || !has {
		log.Error("GetOpenId GetWXPlatform err")
		return nil, datastruct.GetDataFailed
	}
	return wx_user, datastruct.NULLError
}

func (handle *DBHandler) ComputePoundage(userId int, amount float64, trade_no string, ip_addr string) (float64, datastruct.CodeType) {
	engine := handle.mysqlEngine
	session := engine.NewSession()
	defer session.Close()
	session.Begin()

	user := new(datastruct.UserInfo)

	has, err := session.Where("id=?", userId).Get(user)
	if err != nil || !has {
		rollback("DBHandler->ComputePoundage Get UserInfo err", session)
		return -1, datastruct.GetDataFailed
	}
	if user.Balance < amount {
		rollback("DBHandler->ComputePoundage  NotEnoughBalance", session)
		return -1, datastruct.NotEnoughBalance
	}
	user.Balance -= amount
	var affected int64
	affected, err = session.Id(userId).Cols("balance").Update(user)
	if err != nil || affected <= 0 {
		rollback("DBHandler->UserPayeeSuccess Update balance err", session)
		return -1, datastruct.UpdateDataFailed
	}

	params := new(datastruct.DrawCashParams)
	has, err = session.Where("id=?", datastruct.DefaultId).Get(params)
	if err != nil || !has {
		rollback("DBHandler->ComputePoundage  Get DrawCashParams err", session)
		return -1, datastruct.GetDataFailed
	}
	if amount < params.MinCharge {
		rollback("DBHandler->ComputePoundage  NotEnoughMinCharge", session)
		return -1, datastruct.NotEnoughMinCharge
	}

	today_unix, tomorrow_unix := tools.GetTodayTomorrowTime()

	arr := make([]*datastruct.DrawCashInfo, 0)
	err = session.Where("created_at >= ? and created_at < ? and user_id = ?", today_unix, tomorrow_unix, userId).Find(&arr)
	if err != nil {
		rollback("DBHandler->ComputePoundage  Get DrawCashInfo arr err:"+err.Error(), session)
		return -1, datastruct.GetDataFailed
	}
	if len(arr) >= params.MaxDrawCount {
		rollback("DBHandler->ComputePoundage  OverMaxDrawCount", session)
		return -1, datastruct.OverMaxDrawCount
	}
	rs_Poundage := amount * float64(params.PoundagePer) / 100.0
	if rs_Poundage < params.MinPoundage {
		rs_Poundage = params.MinPoundage
	}
	drawCashInfo := new(datastruct.DrawCashInfo)
	drawCashInfo.Charge = amount - rs_Poundage
	drawCashInfo.Poundage = rs_Poundage
	drawCashInfo.CreatedAt = time.Now().Unix()
	drawCashInfo.TradeNo = trade_no
	drawCashInfo.State = datastruct.DrawCashReview
	drawCashInfo.UserId = userId
	drawCashInfo.IpAddr = ip_addr
	drawCashInfo.Origin = amount
	affected, err = session.Insert(drawCashInfo)
	if err != nil || affected <= 0 {
		rollback("DBHandler->ComputePoundage Insert DrawCashInfo err", session)
		return -1, datastruct.UpdateDataFailed
	}

	err = session.Commit()
	if err != nil {
		str := fmt.Sprintf("DBHandler->ComputePoundage Commit :%s", err.Error())
		rollback(str, session)
		return -1, datastruct.UpdateDataFailed
	}

	if amount >= params.RequireVerify {
		return rs_Poundage, datastruct.PayeeReview
	}

	return rs_Poundage, datastruct.NULLError
}

func (handle *DBHandler) UpdateAppOpenId(userId int, openid string) {
	engine := handle.mysqlEngine
	session := engine.NewSession()
	defer session.Close()
	session.Begin()
	wx_user := new(datastruct.WXPlatform)
	_, err := session.Where("user_id=?", userId).Get(wx_user)
	if err != nil {
		str := "UpdateOpenId Get WXPlatform err:" + err.Error()
		rollback(str, session)
		return
	}

	//首次app登录送金币
	if wx_user.PayeeOpenid == "" {
		user := new(datastruct.UserInfo)
		user.IsGotDownLoadAppGift = 0
		var affected int64
		affected, err = engine.Cols("is_got_down_load_app_gift").Update(user)
		if err != nil || affected <= 0 {
			rollback("UpdateOpenId update IsGotDownLoadAppGift err", session)
			return
		}
	}
	wx_user.PayeeOpenid = openid
	wx_user.PayOpenidForKFPT = openid
	_, err = session.Where("user_id=?", userId).Cols("payee_openid", "pay_openid_for_k_f_p_t").Update(wx_user)
	if err != nil {
		rollback("UpdateOpenId payee_openid and pay_openid_for_k_f_p_t err err:"+err.Error(), session)
		return
	}
	err = session.Commit()
	if err != nil {
		str := fmt.Sprintf("DBHandler->UpdateOpenId Commit :%s", err.Error())
		rollback(str, session)
		return
	}
}

func (handle *DBHandler) UpdateGZHOpenId(userId int, openid string) {
	engine := handle.mysqlEngine
	session := engine.NewSession()
	defer session.Close()
	session.Begin()
	wx_user := new(datastruct.WXPlatform)
	_, err := session.Where("user_id=?", userId).Get(wx_user)
	if err != nil {
		str := "UpdateOpenId Get WXPlatform err:" + err.Error()
		rollback(str, session)
		return
	}
	wx_user.PayOpenidForGZH = openid
	_, err = session.Where("user_id=?", userId).Cols("pay_openid_for_g_z_h").Update(wx_user)
	if err != nil {
		rollback("UpdateOpenId pay_openid_for_g_z_h err err:"+err.Error(), session)
		return
	}
	err = session.Commit()
	if err != nil {
		str := fmt.Sprintf("DBHandler->UpdateOpenId Commit :%s", err.Error())
		rollback(str, session)
		return
	}
}

func (handle *DBHandler) UserPayeeSuccess(userId int, payeeData *thirdParty.ResponsWXPayeeXML) (float64, datastruct.CodeType) {
	engine := handle.mysqlEngine
	drawCashInfo := new(datastruct.DrawCashInfo)
	drawCashInfo.State = datastruct.DrawCashSucceed
	drawCashInfo.PaymentNo = payeeData.Payment_no
	drawCashInfo.PaymentTime = payeeData.Payment_time
	affected, err := engine.Where("user_id = ? and trade_no = ?", userId, payeeData.Partner_trade_no).Cols("state", "payment_no", "payment_time").Update(drawCashInfo)
	if err != nil || affected <= 0 {
		log.Debug("UserPayeeSuccess Update DrawCashInfo err")
		return -1, datastruct.UpdateDataFailed
	}
	user := new(datastruct.UserInfo)
	has, err := engine.Id(userId).Get(user)
	if err != nil || !has {
		return -1, datastruct.UpdateDataFailed
	}
	return user.Balance, datastruct.NULLError
}
func (handle *DBHandler) UserPayeefailed(userId int, amount float64, payeeData *thirdParty.ResponsWXPayeeXML) (float64, datastruct.CodeType) {
	engine := handle.mysqlEngine
	session := engine.NewSession()
	defer session.Close()
	session.Begin()

	user := new(datastruct.UserInfo)
	sql := "update user_info set balance = balance + ? where id = ?"
	_, err := session.Exec(sql, amount, userId)
	if err != nil {
		str := fmt.Sprintf("UserPayeefailed restore user's balance err :%s", err.Error())
		rollback(str, session)
		return -1, datastruct.UpdateDataFailed
	}
	drawCashInfo := new(datastruct.DrawCashInfo)
	drawCashInfo.State = datastruct.DrawCashFailed
	_, err = session.Where("user_id = ? and trade_no = ?", userId, payeeData.Partner_trade_no).Cols("state").Update(drawCashInfo)
	if err != nil {
		rollback("UserPayeefailed Update drawCashInfo err:"+err.Error(), session)
		return -1, datastruct.UpdateDataFailed
	}
	var has bool
	has, err = session.Id(userId).Get(user)
	if err != nil || !has {
		rollback("UserPayeefailed Get Userinfo err", session)
		return -1, datastruct.GetDataFailed
	}
	err = session.Commit()
	if err != nil {
		str := fmt.Sprintf("DBHandler->UserPayeefailed Commit :%s", err.Error())
		rollback(str, session)
		return -1, datastruct.UpdateDataFailed
	}
	return user.Balance, datastruct.NULLError
}

func (handle *DBHandler) GetDownLoadAppAddr() (string, datastruct.CodeType) {
	engine := handle.mysqlEngine
	addr := new(datastruct.AppDownloadAddr)
	_, err := engine.Desc("created_at").Limit(1, 0).Get(addr)
	if err != nil {
		return "", datastruct.GetDataFailed
	}
	return addr.DownLoadUrl, datastruct.NULLError
}

func (handle *DBHandler) GetAppAddr() (*datastruct.AppAddr, datastruct.CodeType) {
	engine := handle.mysqlEngine
	addr := new(datastruct.AppAddr)
	_, err := engine.Desc("created_at").Limit(1, 0).Get(addr)
	if err != nil {
		return nil, datastruct.GetDataFailed
	}
	return addr, datastruct.NULLError
}

func (handle *DBHandler) GetKfInfo() (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	kfInfo := new(datastruct.KfInfo)
	has, err := engine.Desc("w_x").Limit(1, 0).Get(kfInfo)
	if err != nil || !has {
		return nil, datastruct.GetDataFailed
	}
	mp := make(map[string]string)
	mp["qq"] = kfInfo.QQ
	mp["wx"] = kfInfo.WX
	mp["qrcode"] = tools.GetKfQRcode(kfInfo.QRcode)
	return mp, datastruct.NULLError
}

func (handle *DBHandler) GetEntryUrl() string {
	engine := handle.mysqlEngine
	entry := new(datastruct.EntryAddr)
	has, err := engine.Desc("created_at").Limit(1, 0).Get(entry)
	if err != nil || !has {
		return ""
	}
	return "https://" + entry.Url
}

func (handle *DBHandler) GetEntryPageUrl() string {
	engine := handle.mysqlEngine
	entry := new(datastruct.EntryAddr)
	has, err := engine.Desc("created_at").Limit(1, 0).Get(entry)
	if err != nil || !has {
		return ""
	}
	return "http://" + entry.PageUrl + ":9100"
}

func (handle *DBHandler) GetAppDownLoadShareUrl() string {
	engine := handle.mysqlEngine
	download := new(datastruct.AppDownloadAddr)
	has, err := engine.Desc("created_at").Limit(1, 0).Get(download)
	if err != nil || !has {
		return ""
	}
	return download.DownLoadUrl
}

func (handle *DBHandler) GetDirectDownloadApp() string {
	engine := handle.mysqlEngine
	direct := new(datastruct.AppDownloadAddr)
	has, err := engine.Desc("created_at").Limit(1, 0).Get(direct)
	if err != nil || !has {
		return ""
	}
	return direct.DirectDownLoadUrl
}

func (handle *DBHandler) GetAuthUrl() string {
	engine := handle.mysqlEngine
	auth := new(datastruct.AuthAddr)
	has, err := engine.Desc("created_at").Limit(1, 0).Get(auth)
	if err != nil || !has {
		return ""
	}
	return "http://" + auth.Url
}

func (handle *DBHandler) GetRedirect() (string, string) {
	engine := handle.mysqlEngine
	jump := new(datastruct.BlackListJump)
	has, err := engine.Desc("created_at").Limit(1, 0).Get(jump)
	if err != nil || !has {
		return "", "http://www.baidu.com"
	}
	return jump.PCJumpTo, jump.BLJumpTo
}

func (handle *DBHandler) CustomShareForApp(userid int) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	share := new(datastruct.AppCustomShare)
	engine.Desc("created_at").Limit(1, 0).Get(share)
	user := new(datastruct.UserInfo)
	engine.Where("id=?", userid).Get(user)
	resp := new(datastruct.ResponseShareData)
	resp.Desc = share.Desc
	resp.Title = strings.Replace(share.Title, "#nickname#", user.NickName, -1)
	link := handle.GetEntryUrl()
	resp.Link = tools.CreateInviteLink(tools.IntToString(userid), link)
	resp.ImgUrl = tools.GetShareImgUrl(share.ImgName)
	return resp, datastruct.NULLError
}

func (handle *DBHandler) CustomShareForGZH(userid int) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	share := new(datastruct.GZHCustomShare)
	engine.Desc("created_at").Limit(1, 0).Get(share)
	user := new(datastruct.UserInfo)
	engine.Where("id=?", userid).Get(user)
	resp := new(datastruct.ResponseShareData)
	resp.Desc = share.Desc
	resp.Title = strings.Replace(share.Title, "#nickname#", user.NickName, -1)
	link := handle.GetEntryUrl()
	resp.Link = tools.CreateInviteLink(tools.IntToString(userid), link)
	resp.ImgUrl = tools.GetShareImgUrl(share.ImgName)
	return resp, datastruct.NULLError
}

func (handle *DBHandler) GetCheckInData(userId int) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	checkIn := new(datastruct.CheckInInfo)
	has, err := engine.Where("user_id=?", userId).Get(checkIn)
	if err != nil {
		return nil, datastruct.GetDataFailed
	}
	continuousCount := 0
	isCheckedInToday := 0
	if has {
		yesterday_unix, today1_unix := tools.GetYesterdayTodayTime()
		today2_unix, tomorrow_unix := tools.GetTodayTomorrowTime()
		if checkIn.LastCheckIn >= yesterday_unix && checkIn.LastCheckIn < today1_unix {
			//昨天签到
			continuousCount = checkIn.ContinuousCount
		} else if checkIn.LastCheckIn >= today2_unix && checkIn.LastCheckIn < tomorrow_unix {
			//今天已签到
			continuousCount = checkIn.ContinuousCount
			isCheckedInToday = 1
		} else {
			//删除记录
			affected, err := engine.Where("user_id=?", userId).Delete(new(datastruct.CheckInInfo))
			if err != nil || affected <= 0 {
				return nil, datastruct.UpdateDataFailed
			}
		}
	}

	//get gold from CommonData  CheckInReward {
	list := make([]*datastruct.CheckInReward, 0)
	err = engine.Asc("day_index").Find(&list)
	if err != nil {
		return nil, datastruct.GetDataFailed
	}
	// }

	resp := make([]*datastruct.ResponseCheckInData, 0)
	for i := 0; i < continuousCount; i++ {
		data := new(datastruct.ResponseCheckInData)
		resp = append(resp, data)
		data.Gold = list[i].RewardGold
		data.IsCheckedIn = 1
	}
	for i := continuousCount; i < len(list); i++ {
		data := new(datastruct.ResponseCheckInData)
		resp = append(resp, data)
		data.Gold = list[i].RewardGold
		data.IsCheckedIn = 0
	}
	resp_map := make(map[string]interface{})
	resp_map["list"] = resp
	resp_map["ischeckedintoday"] = isCheckedInToday
	return resp_map, datastruct.NULLError
}

func (handle *DBHandler) CheckIn(userId int) (int, datastruct.CodeType) {
	engine := handle.mysqlEngine
	checkIn := new(datastruct.CheckInInfo)
	has, err := engine.Where("user_id=?", userId).Get(checkIn)
	if err != nil {
		return -1, datastruct.GetDataFailed
	}
	continuousCount := 0
	list := make([]*datastruct.CheckInReward, 0)
	err = engine.Asc("day_index").Find(&list)
	if err != nil {
		return -1, datastruct.GetDataFailed
	}
	var gold int64
	if has {
		yesterday_unix, today1_unix := tools.GetYesterdayTodayTime()
		today2_unix, tomorrow_unix := tools.GetTodayTomorrowTime()
		if checkIn.LastCheckIn >= yesterday_unix && checkIn.LastCheckIn < today1_unix {
			//昨天签到
			continuousCount = checkIn.ContinuousCount + 1
			checkIn.ContinuousCount = continuousCount
			checkIn.LastCheckIn = time.Now().Unix()
			affected, err := engine.Where("user_id=?", userId).Cols("continuous_count", "last_check_in").Update(checkIn)
			if err != nil || affected <= 0 {
				return -1, datastruct.UpdateDataFailed
			}
			gold = list[continuousCount-1].RewardGold
		} else if checkIn.LastCheckIn >= today2_unix && checkIn.LastCheckIn < tomorrow_unix {
			//今天已签到
			return checkIn.ContinuousCount, datastruct.TodayCheckedIn
		} else {
			//删除记录
			gold = list[0].RewardGold
			continuousCount := 1
			checkIn.ContinuousCount = continuousCount
			checkIn.LastCheckIn = time.Now().Unix()
			affected, err := engine.Where("user_id=?", userId).Cols("continuous_count", "last_check_in").Update(checkIn)
			if err != nil || affected <= 0 {
				return -1, datastruct.UpdateDataFailed
			}
		}
	} else {
		checkIn.ContinuousCount = 1
		checkIn.LastCheckIn = time.Now().Unix()
		checkIn.UserId = userId
		affected, err := engine.Insert(checkIn)
		if err != nil || affected <= 0 {
			return -1, datastruct.UpdateDataFailed
		}
		continuousCount = checkIn.ContinuousCount
		gold = list[0].RewardGold
	}
	session := engine.NewSession()
	defer session.Close()
	session.Begin()

	now_time := time.Now().Unix()
	goldChangeInfo := new(datastruct.GoldChangeInfo)
	goldChangeInfo.UserId = userId
	goldChangeInfo.CreatedAt = now_time
	goldChangeInfo.VarGold = gold
	goldChangeInfo.ChangeType = datastruct.CheckInType
	_, err = session.Insert(goldChangeInfo)
	if err != nil {
		str := fmt.Sprintf("DBHandler->CheckIn Insert GoldChangeInfo :%s", err.Error())
		rollback(str, session)
		return -1, datastruct.UpdateDataFailed
	}

	sql := "update user_info set gold_count = gold_count + ? where id = ?"
	res, err1 := session.Exec(sql, gold, userId)
	affected, err2 := res.RowsAffected()
	if err1 != nil || err2 != nil || affected <= 0 {
		rollback("DBHandler->CheckIn update user_info", session)
		return -1, datastruct.UpdateDataFailed
	}

	err = session.Commit()
	if err != nil {
		str := fmt.Sprintf("DBHandler->CheckIn Commit :%s", err.Error())
		rollback(str, session)
		return -1, datastruct.UpdateDataFailed
	}
	return continuousCount, datastruct.NULLError
}

// type AppMemberData struct {
// 	datastruct.MemberLevelData `xorm:"extends"`
// 	datastruct.AgencyParams    `xorm:"extends"`
// }

func (handle *DBHandler) AppGetMemberList(userId int) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	user := new(datastruct.UserInfo)
	has, err := engine.Where("id=?", userId).Get(user)
	if err != nil || !has {
		return nil, datastruct.GetDataFailed
	}
	max_level := 3
	members := make([]*datastruct.MemberLevelData, 0)
	err = engine.Where("level<=? and is_hidden = 0", max_level).Asc("level").Find(&members)
	if err != nil {
		return nil, datastruct.GetDataFailed
	}
	resp := make([]*datastruct.ResponseMemberData, 0)
	common := datastruct.AgencyIdentifier
	ap, code := getAgencyParams(common, engine)
	if code != datastruct.NULLError {
		return nil, code
	}
	m_data := createMemberData(ap)
	m_data.Price = 0
	m_data.Name = datastruct.DefaultMemberIdentifier
	m_data.Enable = 0
	var currentLevel int
	var currentIdentifier int
	if user.MemberIdentifier == common {
		m_data.Owned = 1
		currentLevel = 0
		currentIdentifier = 0
		m_data.VipId = 0
	} else {
		m_data.Owned = 0
		currentIdentifier = tools.StringToInt(user.MemberIdentifier)
		currentLevel, code = getMemberLevel(user.MemberIdentifier, engine)
		if code != datastruct.NULLError {
			return nil, code
		}
	}
	resp = append(resp, m_data)
	for _, v := range members {
		ap, code := getAgencyParams(tools.IntToString(v.Id), engine)
		if code != datastruct.NULLError {
			return nil, code
		}
		m_data = createMemberData(ap)
		m_data.Name = v.Name
		m_data.Price = v.Price
		m_data.VipId = v.Id
		if currentIdentifier == v.Id {
			m_data.Owned = 1
			m_data.Enable = 0
		} else {
			m_data.Owned = 0
			if currentLevel >= v.Level {
				m_data.Enable = 0
			} else {
				m_data.Enable = 1
			}
		}
		resp = append(resp, m_data)
	}
	return resp, datastruct.NULLError
}

func getMemberLevel(identifier string, engine *xorm.Engine) (int, datastruct.CodeType) {
	ml := new(datastruct.MemberLevelData)
	has, err := engine.Where("id=?", tools.StringToInt(identifier)).Get(ml)
	if err != nil || !has {
		return -1, datastruct.GetDataFailed
	}
	return ml.Level, datastruct.NULLError
}

func getAgencyParams(identifier string, engine *xorm.Engine) (*datastruct.AgencyParams, datastruct.CodeType) {
	ap := new(datastruct.AgencyParams)
	has, err := engine.Where("identifier=?", identifier).Get(ap)
	if err != nil || !has {
		return nil, datastruct.GetDataFailed
	}
	return ap, datastruct.NULLError
}
func createMemberData(ap *datastruct.AgencyParams) *datastruct.ResponseMemberData {
	m_data := new(datastruct.ResponseMemberData)
	m_data.A1Gold = ap.Agency1GoldPercent
	m_data.A2Gold = ap.Agency2GoldPercent
	m_data.A3Gold = ap.Agency3GoldPercent
	m_data.A1Money = ap.Agency1MoneyPercent
	m_data.A2Money = ap.Agency2MoneyPercent
	m_data.A3Money = ap.Agency3MoneyPercent
	return m_data
}

func (handle *DBHandler) QueryMemberLevelData(m_id int) (*datastruct.MemberLevelData, datastruct.CodeType) {
	engine := handle.mysqlEngine
	m_data := new(datastruct.MemberLevelData)
	has, err := engine.Where("id=?", m_id).Get(m_data)
	if err != nil || !has {
		log.Error("-----QueryMemberLevelData err")
		return nil, datastruct.GetDataFailed
	}
	return m_data, datastruct.NULLError
}

func (handle *DBHandler) IsRefreshMemberList(userId int, level int) datastruct.CodeType {
	engine := handle.mysqlEngine
	user := new(datastruct.UserInfo)
	has, err := engine.Where("id=?", userId).Get(user)
	if err != nil || !has {
		log.Error("-----IsRefreshMemberList err")
		return datastruct.GetDataFailed
	}
	if user.MemberIdentifier == datastruct.AgencyIdentifier {
		return datastruct.NULLError
	}
	m_data := new(datastruct.MemberLevelData)
	has, err = engine.Where("id=?", tools.StringToInt(user.MemberIdentifier)).Get(m_data)
	if err != nil || !has {
		return datastruct.GetDataFailed
	}
	if m_data.Level >= level {
		return datastruct.RefreshMemberList
	}
	return datastruct.NULLError
}

//购买Vip成功
func (handle *DBHandler) PurchaseVipSuccess(userId int, vipId int, createTime int64) datastruct.CodeType {
	engine := handle.mysqlEngine
	m_data := new(datastruct.MemberLevelData)
	has, err := engine.Where("id=?", vipId).Get(m_data)
	if err != nil || !has {
		log.Error("PurchaseVipSuccess Get MemberLevelData err")
		return datastruct.GetDataFailed
	}

	session := engine.NewSession()
	defer session.Close()
	session.Begin()
	user := new(datastruct.UserInfo)
	user.MemberIdentifier = tools.IntToString(m_data.Id)
	_, err = session.Cols("member_identifier").Where("id=?", userId).Update(user)
	if err != nil {
		str := fmt.Sprintf("DBHandler->PurchaseVipSuccess update UserInfo :%s", err.Error())
		rollbackError(str, session)
		return datastruct.UpdateDataFailed
	}
	m_order := new(datastruct.MemberLevelOrder)
	m_order.CreatedAt = createTime
	m_order.MemberLevelId = m_data.Id
	m_order.UserId = userId
	var affected int64
	affected, err = session.Insert(m_order)
	if err != nil || affected <= 0 {
		rollbackError("DBHandler->PurchaseVipSuccess Insert MemberLevelOrder err", session)
		return datastruct.UpdateDataFailed
	}
	err = session.Commit()
	if err != nil {
		str := fmt.Sprintf("DBHandler->PurchaseVipSuccess Commit :%s", err.Error())
		rollbackError(str, session)
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

// func (handle *DBHandler) UserActivate(userId int) datastruct.CodeType {
// 	engine := handle.mysqlEngine
// 	activity := new(datastruct.TodayUserActivityInfo)
// 	today_unix, tomorrow_unix := tools.GetTodayTomorrowTime()
// 	has, err := engine.Where("user_id=? and start_time >= ? and start_time < ?", userId, today_unix, tomorrow_unix).Get(activity)
// 	if err != nil {
// 		return datastruct.GetDataFailed
// 	}
// 	if !has {
// 		activity.UserId = userId
// 		activity.StartTime = time.Now().Unix()
// 		activity.EndTime = activity.StartTime
// 		var affected int64
// 		affected, err = engine.Insert(activity)
// 		if err != nil || affected <= 0 {
// 			return datastruct.UpdateDataFailed
// 		}
// 	}
// 	return datastruct.NULLError
// }

func (handle *DBHandler) StartLottery(userId int, rushprice int64) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	resp := new(datastruct.ResponseLotteryGoodsInfo)
	var goods_1 *datastruct.RandomLotteryGoods
	user := new(datastruct.UserInfo)
	has, err := engine.Where("id=?", userId).Get(user)
	if err != nil || !has {
		return nil, datastruct.GetDataFailed
	}
	allGoods := make([]*datastruct.RandomLotteryGoods, 0)
	if has && user.DepositTotal <= 0 && user.PurchaseTotal <= 0 { //没有充值和购买行为的用户,中奖率为0
		resp.IsGotReward = 0
	} else {
		goodsPool := new(datastruct.RandomLotteryGoodsPool)
		has, err = engine.Where("id=?", datastruct.DefaultId).Get(goodsPool)
		if err != nil || !has {
			return nil, datastruct.GetDataFailed
		}
		//根据当前池水值,查询出等值的商品
		err = engine.Where("is_hidden = 0 and price <= ?", goodsPool.Current).Find(&allGoods)
		length := len(allGoods)
		if err != nil {
			return nil, datastruct.GetDataFailed
		}
		if length <= 0 {
			//没有查出等值的商品,中奖率为0
			resp.IsGotReward = 0
		} else {
			goods_1 = getRandomGood(allGoods) //作为中奖商品
			random_rs := tools.RandInt(1, 101)
			var probability int
			if user.RandomLotterySucceed >= goodsPool.RandomLotteryCount {
				//当用户超过指定中奖次数,概率恒定为后端指定概率
				probability = goodsPool.Probability
			} else {
				probability = goods_1.Probability
			}
			if random_rs > 100-probability {
				//中奖
				resp.IsGotReward = 1
				session := engine.NewSession()
				defer session.Close()
				session.Begin()
				//增加用户中奖次数
				sql := "update user_info set random_lottery_succeed = random_lottery_succeed + 1 where id = ?"
				res, err1 := session.Exec(sql, userId)
				affected, err2 := res.RowsAffected()
				if err1 != nil || err2 != nil || affected <= 0 {
					rollback("DBHandler->StartLottery UpdateUser", session)
					return nil, datastruct.UpdateDataFailed
				}
				//更新中奖池水
				sql = "update random_lottery_goods_pool set current = case when current - ? > 0 then current - ? else 0 end"
				_, err = session.Exec(sql, goods_1.Price, goods_1.Price)
				if err != nil {
					rollback("DBHandler->StartLottery Update GoodsRewardPool err", session)
					return nil, datastruct.UpdateDataFailed
				}
				now_time := time.Now().Unix()
				orderId := tools.StringToInt64(tools.UniqueId())
				//添加中奖成功订单
				rlgs := new(datastruct.RandomLotteryGoodsSucceed)
				rlgs.UserId = userId
				rlgs.OrderId = orderId
				rlgs.LotteryGoodsId = goods_1.Id
				rlgs.CreatedAt = now_time
				affected, err = session.Insert(rlgs)
				if err != nil || affected <= 0 {
					rollback("DBHandler->StartLottery Insert RandomLotteryGoodsSucceed err", session)
					return nil, datastruct.UpdateDataFailed
				}
				//添加中奖成功记录
				rlgsh := new(datastruct.RandomLotteryGoodsSucceedHistory)
				rlgsh.DescInfo = fmt.Sprintf("抽中%s*1", goods_1.Name)
				rlgsh.NickName = user.NickName
				rlgsh.Avatar = user.Avatar
				affected, err = session.Insert(rlgsh)
				if err != nil || affected <= 0 {
					rollback("DBHandler->StartLottery Insert RandomLotteryGoodsSucceedHistory err", session)
					return nil, datastruct.UpdateDataFailed
				}
				//添加对应的发货订单
				sendGoods := new(datastruct.SendGoods)
				sendGoods.OrderId = orderId
				sendGoods.SendGoodsState = 0
				sendGoods.CreatedAt = now_time
				sendGoods.IsLotteryGoods = 1
				affected, err = session.Insert(sendGoods)
				if err != nil || affected <= 0 {
					rollback("DBHandler->StartLottery Insert SendGoods err", session)
					return nil, datastruct.UpdateDataFailed
				}
				err = session.Commit()
				if err != nil {
					str := fmt.Sprintf("DBHandler->StartLottery Commit :%s", err.Error())
					rollback(str, session)
					return -1, datastruct.UpdateDataFailed
				}
			} else {
				//没中奖
				resp.IsGotReward = 0
			}
		}
	}
	if resp.IsGotReward == 0 {
		sql := "update random_lottery_goods_pool set current = current + ?"
		_, err = engine.Exec(sql, rushprice)
		if err != nil {
			log.Debug("IsGotReward == 0,set current err:%s", err.Error())
			// rollback("DBHandler->PayRougeGame Update GoodsRewardPool err", session)
			return nil, datastruct.UpdateDataFailed
		}
	}
	if resp.IsGotReward == 1 {
		resp.GoodsImg = tools.CreateGoodsImgUrl(goods_1.ImgName)
	}
	return resp, datastruct.NULLError
}

func (handle *DBHandler) GetLotteryGoodsSucceedHistory() (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	rgsh := make([]*datastruct.RandomLotteryGoodsSucceedHistory, 0)
	err := engine.Desc("id").Find(&rgsh)
	if err != nil {
		return nil, datastruct.GetDataFailed
	}
	resp := make([]*datastruct.ResponseLotteryGoodsSucceedHistory, 0)
	for _, v := range rgsh {
		tmp := new(datastruct.ResponseLotteryGoodsSucceedHistory)
		tmp.Avatar = v.Avatar
		tmp.DescInfo = v.DescInfo
		tmp.NickName = v.NickName
		resp = append(resp, tmp)
	}
	return resp, datastruct.NULLError
}

func (handle *DBHandler) GetLotteryGoodsOrderInfo(userId int, pageIndex int, pageSize int) (interface{}, datastruct.CodeType) {

	start := (pageIndex - 1) * pageSize
	limit := pageSize
	engine := handle.mysqlEngine

	sql := "select rlg.name,rlg.img_name,rlg.price,rlgs.order_id,s.send_goods_state,s.is_appraised,s.express_number,s.express_agency from random_lottery_goods_succeed rlgs inner join random_lottery_goods rlg on rlgs.lottery_goods_id = rlg.id inner join send_goods s on rlgs.order_id = s.order_id where rlgs.user_id = ? ORDER BY rlgs.created_at desc LIMIT ?,?"

	result, err := engine.Query(sql, userId, start, limit)

	if err != nil {
		return nil, datastruct.GetDataFailed
	}

	resp := make([]interface{}, 0)
	for _, v := range result {
		send_goods_state := tools.StringToInt(string(v["send_goods_state"][:]))
		name := string(v["name"][:])
		imgurl := tools.CreateGoodsImgUrl(string(v["img_name"][:]))
		pricedesc := fmt.Sprintf("价值: ¥ %s.00", string(v["price"][:]))
		count := 1
		number := string(v["order_id"][:])
		var rs interface{}
		if send_goods_state == int(datastruct.NotSend) {
			tmp := new(datastruct.ResponseLotteryOrderNotSend)
			tmp.Name = name
			tmp.ImgUrl = imgurl
			tmp.PriceDesc = pricedesc
			tmp.Count = count
			tmp.Number = number
			tmp.Remark = "此奖品需要联系客服发货"
			tmp.PostAge = "邮费8元，晒单免邮费"
			tmp.IsSended = send_goods_state
			rs = tmp
		} else {
			tmp := new(datastruct.ResponseLotteryOrderHasSended)
			tmp.Name = name
			tmp.ImgUrl = imgurl
			tmp.PriceDesc = pricedesc
			tmp.Count = count
			tmp.Number = number
			tmp.Remark = "评价商品可免费得金币"
			tmp.ExpressNumber = string(v["express_number"][:])
			tmp.ExpressAgency = string(v["express_agency"][:])
			tmp.IsAppraised = tools.StringToInt(string(v["is_appraised"][:]))
			tmp.IsSended = send_goods_state
			rs = tmp
		}
		resp = append(resp, rs)
	}
	return resp, datastruct.NULLError
}

func getRandomGood(allGoods []*datastruct.RandomLotteryGoods) *datastruct.RandomLotteryGoods {
	length := len(allGoods)
	allGoods_index := make([]int, 0, length)
	for i := 0; i < length; i++ {
		allGoods_index = append(allGoods_index, i)
	}
	index_1 := tools.RandInt(0, len(allGoods_index))
	goods_1 := allGoods[allGoods_index[index_1]]
	return goods_1
}

func (handle *DBHandler) GetSharePosters(qrcode string) (interface{}, datastruct.CodeType) {
	sp := make([]*datastruct.SharePosters, 0)
	engine := handle.mysqlEngine
	start := 0
	limit := 3
	err := engine.Where("is_hidden = ?", 0).Desc("sort_id").Limit(limit, start).Find(&sp)
	if err != nil {
		return nil, datastruct.GetDataFailed
	}

	resp := new(datastruct.ResponseSharePosters)
	posters := make([]*datastruct.ResponsePosters, 0, len(sp))
	for _, v := range sp {
		poster := new(datastruct.ResponsePosters)
		poster.ImgUrl = osstool.CreateOSSURL(v.ImgName)
		poster.Icon = osstool.CreateOSSURL(v.IconName)
		poster.Location = v.Location
		posters = append(posters, poster)
	}
	resp.Posters = posters
	resp.QRcode = qrcode
	return resp, datastruct.NULLError
}

func (handle *DBHandler) GetUserAppraiseForApp(pageIndex int, pageSize int) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	start := (pageIndex - 1) * pageSize
	limit := pageSize

	limitStr := fmt.Sprintf(" LIMIT %d,%d", start, limit)
	orderby := " ORDER BY uap.created_at desc"

	sql := "select uap.goods_type,uap.id,uap.show_type,uap.desc,u.nick_name,u.avatar,g.name as gname,rlg.name as rlgname,uap.created_at from user_appraise uap inner join user_info u on u.id = uap.user_id left join goods g on uap.goods_id = g.id left join random_lottery_goods rlg on uap.goods_id = rlg.id where is_passed = 1" + orderby + limitStr

	results, _ := engine.Query(sql)
	list := make([]*datastruct.ResponseUserAppraise, 0, len(results))

	for _, v := range results {
		r_uap := new(datastruct.ResponseUserAppraise)
		r_uap.Avatar = string(v["avatar"][:])
		r_uap.CreatedAt = tools.StringToInt64(string(v["created_at"][:]))
		r_uap.Desc = string(v["desc"][:])

		goods_type := tools.StringToInt(string(v["goods_type"][:]))
		switch goods_type {
		case int(datastruct.RushGoods):
			r_uap.GoodsName = string(v["gname"][:])
		case int(datastruct.LotteryGoods):
			r_uap.GoodsName = string(v["rlgname"][:])
		}

		uap_id := tools.StringToInt(string(v["id"][:]))
		r_uap.NickName = string(v["nick_name"][:])

		showType := tools.StringToInt(string(v["show_type"][:]))
		if datastruct.UserAppraiseType(showType) != datastruct.OnlyT {
			tmp_uapps := make([]*datastruct.UserAppraisePic, 0)
			engine.Where("user_appraise_id = ?", uap_id).Asc("img_index").Find(&tmp_uapps)
			imgurls := make([]string, 0, len(tmp_uapps))
			for _, v := range tmp_uapps {
				url := osstool.CreateOSSURL(v.ImgName)
				imgurls = append(imgurls, url)
			}
			r_uap.ImgUrls = imgurls
		}
		list = append(list, r_uap)
	}
	return list, datastruct.NULLError
}

func (handle *DBHandler) UserAppraise(userId int, body *datastruct.UserAppraiseBody) datastruct.CodeType {
	engine := handle.mysqlEngine
	var sql string
	var err error
	var results []map[string][]byte
	switch body.GoodsType {
	case datastruct.RushGoods:
		sql = "select s.is_appraised,o.goods_id,s.order_id,s.sign_for_state from send_goods s inner join order_info o on s.order_id = o.id where o.number = ? and o.user_id = ?"
		results, err = engine.Query(sql, body.Number, userId)
	case datastruct.LotteryGoods:
		sql = "select s.is_appraised,r.lottery_goods_id as goods_id,s.order_id,s.sign_for_state from send_goods s inner join random_lottery_goods_succeed r on s.order_id = r.order_id where r.order_id = ? and r.user_id = ?"
		results, err = engine.Query(sql, tools.StringToInt64(body.Number), userId)
	}
	if err != nil || len(results) <= 0 {
		return datastruct.GetDataFailed
	}
	v := results[0]
	is_appraised := tools.StringToInt(string(v["is_appraised"][:]))
	if is_appraised == 1 {
		return datastruct.IsAppraised
	}
	if body.GoodsType == datastruct.RushGoods {
		sign_for_state := tools.StringToInt(string(v["sign_for_state"][:]))
		if sign_for_state == int(datastruct.NotSignGoods) {
			return datastruct.NotHasSignedGoods
		}
	}
	goods_id := tools.StringToInt(string(v["goods_id"][:]))
	order_id := tools.StringToInt(string(v["order_id"][:]))
	session := engine.NewSession()
	defer session.Close()
	session.Begin()

	now_time := time.Now().Unix()
	uap := new(datastruct.UserAppraise)
	uap.Desc = body.Desc
	uap.GoodsId = goods_id
	uap.IsPassed = 0
	uap.UserId = userId
	uap.CreatedAt = now_time
	uap.GoodsType = body.GoodsType
	length := len(body.ImgNames)
	if length < 0 {
		uap.ShowType = datastruct.OnlyT
	} else if body.Desc == "" {
		uap.ShowType = datastruct.OnlyP
	} else {
		uap.ShowType = datastruct.TAndP
	}
	var affected int64
	affected, err = session.Insert(uap)
	if err != nil || affected <= 0 {
		rollback("DBHandler->UserAppraise Insert err", session)
		return datastruct.UpdateDataFailed
	}
	for i := 0; i < length; i++ {
		v := body.ImgNames[i]
		uapp := new(datastruct.UserAppraisePic)
		uapp.ImgIndex = i
		uapp.ImgName = v
		uapp.UserAppraiseId = uap.Id
		affected, err = session.Insert(uapp)
		if err != nil || affected <= 0 {
			rollback("DBHandler->UserAppraise Insert GoodsImgs err", session)
			return datastruct.UpdateDataFailed
		}
	}
	sendgoods := new(datastruct.SendGoods)
	sendgoods.IsAppraised = 1
	affected, err = session.Where("order_id=?", order_id).Cols("is_appraised").Update(sendgoods)
	if err != nil || affected <= 0 {
		rollback("DBHandler->UserAppraise Update SendGoods err", session)
		return datastruct.UpdateDataFailed
	}
	gcg := new(datastruct.GoldCoinGift)
	var has bool
	has, err = session.Where("id=?", datastruct.DefaultId).Get(gcg)
	if err != nil || !has {
		rollback("UpdateOpenId Get GoldCoinGift err", session)
		return datastruct.GetDataFailed
	}
	add_goldcount := gcg.AppraisedGoldGift
	sql = "update user_info set gold_count = gold_count + ? where id = ?"
	res, err1 := session.Exec(sql, add_goldcount, userId)
	affected, err2 := res.RowsAffected()
	if err1 != nil || err2 != nil || affected <= 0 {
		rollback("UserAppraise update AddGoldCount err", session)
		return datastruct.UpdateDataFailed
	}
	goldChangeInfo := new(datastruct.GoldChangeInfo)
	goldChangeInfo.UserId = userId
	goldChangeInfo.CreatedAt = now_time
	goldChangeInfo.VarGold = add_goldcount
	goldChangeInfo.ChangeType = datastruct.AppraisedType
	_, err = session.Insert(goldChangeInfo)
	if err != nil {
		str := fmt.Sprintf("UserAppraise Insert GoldChangeInfo :%s", err.Error())
		rollback(str, session)
		return datastruct.UpdateDataFailed
	}

	err = session.Commit()
	if err != nil {
		str := fmt.Sprintf("DBHandler->UserAppraise Commit :%s", err.Error())
		rollback(str, session)
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

type goodsDetailData struct {
	datastruct.Goods            `xorm:"extends"`
	datastruct.RecommendedClass `xorm:"extends"`
}

func (handle *DBHandler) GetGoodsDetailForApp(goodsid int) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine

	sql := "select img_name from goods_imgs where goods_id = ? ORDER BY img_index asc"
	results, err := engine.Query(sql, goodsid)
	if err != nil {
		log.Debug("GetGoodsDetailForApp Query0 err:%v", err.Error())
		return nil, datastruct.GetDataFailed
	}
	goods_cover := make([]string, 0, len(results))
	for _, v := range results {
		goods_cover = append(goods_cover, tools.CreateGoodsImgUrl(string(v["img_name"][:])))
	}

	sql = "select img_name from goods_detail where goods_id = ? ORDER BY img_index asc"
	results, err = engine.Query(sql, goodsid)
	if err != nil {
		log.Debug("GetGoodsDetailForApp Query1 err:%v", err.Error())
		return nil, datastruct.GetDataFailed
	}
	goods_detail := make([]string, 0, len(results))
	for _, v := range results {
		url := osstool.CreateOSSURL(string(v["img_name"][:]))
		goods_detail = append(goods_detail, url)
	}

	queryGoods := make([]*goodsDetailData, 0)
	err = engine.Table("goods").Join("Left", "recommended_class", "recommended_class.id = goods.re_classid").Where("goods.id=?", goodsid).Limit(1, 0).Find(&queryGoods)
	if err != nil || len(queryGoods) <= 0 {
		log.Debug("GetGoodsDetailForApp Query2 err:%v", err.Error())
		return nil, datastruct.GetDataFailed
	}
	rs_query := queryGoods[0]

	total_avatar := 20
	avatarCount := 18
	tmpUsers := make([]*datastruct.TmpUser, 0, total_avatar)
	sql = "select avatar,nick_name from tmp_data td inner join tmp_data_for_goods tdfg on td.id = tdfg.tmp_user_id where tdfg.goods_id=? order by td.id desc limit ?"
	results, err = engine.Query(sql, goodsid, total_avatar-avatarCount)
	if err != nil {
		log.Debug("GetGoodsDetailForApp Query3.1 err:%v", err.Error())
		return nil, datastruct.GetDataFailed
	}
	for _, v := range results {
		tmpUser := new(datastruct.TmpUser)
		tmpUser.Avatar = string(v["avatar"][:])
		tmpUser.NickName = string(v["nick_name"][:])
		tmpUser.Desc = fmt.Sprintf("已获得%s*1", rs_query.Goods.Name)
		tmpUsers = append(tmpUsers, tmpUser)
	}

	sql = "select avatar,nick_name from pay_mode_rouge_game_succeed_history where goods_id=? order by id desc limit ?"
	results, err = engine.Query(sql, goodsid, avatarCount)
	if err != nil {
		log.Debug("GetGoodsDetailForApp Query3 err:%v", err.Error())
		return nil, datastruct.GetDataFailed
	}

	for _, v := range results {
		tmpUser := new(datastruct.TmpUser)
		tmpUser.Avatar = string(v["avatar"][:])
		tmpUser.NickName = string(v["nick_name"][:])
		tmpUser.Desc = fmt.Sprintf("已获得%s*1", rs_query.Goods.Name)
		tmpUsers = append(tmpUsers, tmpUser)
	}

	goodsDetailAppraise := new(datastruct.GoodsDetailAppraise)

	var AppraiseTotal int64
	AppraiseTotal, err = engine.Where("is_passed = 1").Count(new(datastruct.UserAppraise))
	if err != nil {
		log.Debug("GetGoodsDetailForApp Query4 err:%v", err.Error())
		return nil, datastruct.GetDataFailed
	}
	goodsDetailAppraise.Total = AppraiseTotal
	userAppraise := new(datastruct.ResponseUserAppraise)

	appraiseLimit := 1

	sql = "select uap.goods_type,uap.id,avatar,nick_name,uap.desc,show_type,uap.created_at,g.name as gname,rlg.name as rlgname from user_appraise uap inner join user_info u on uap.user_id=u.id left join goods g on g.id = uap.goods_id left join random_lottery_goods rlg on uap.goods_id = rlg.id where is_passed = 1 order by uap.created_at desc limit ?"
	results, err = engine.Query(sql, appraiseLimit)
	if err != nil || len(results) <= 0 {
		log.Debug("GetGoodsDetailForApp Query5 err")
		return nil, datastruct.GetDataFailed
	}
	tmp_v := results[0]
	userAppraise.Avatar = string(tmp_v["avatar"][:])
	userAppraise.NickName = string(tmp_v["nick_name"][:])
	userAppraise.CreatedAt = tools.StringToInt64(string(tmp_v["created_at"][:]))
	userAppraise.Desc = string(tmp_v["desc"][:])

	goods_type := tools.StringToInt(string(tmp_v["goods_type"][:]))
	switch goods_type {
	case int(datastruct.RushGoods):
		userAppraise.GoodsName = string(tmp_v["gname"][:])
	case int(datastruct.LotteryGoods):
		userAppraise.GoodsName = string(tmp_v["rlgname"][:])
	}

	userAppraise_id := tools.StringToInt(string(tmp_v["id"][:]))
	showType := tools.StringToInt(string(tmp_v["show_type"][:]))

	if datastruct.UserAppraiseType(showType) != datastruct.OnlyT {
		tmp_uapps := make([]*datastruct.UserAppraisePic, 0)
		engine.Where("user_appraise_id = ?", userAppraise_id).Asc("img_index").Find(&tmp_uapps)
		imgurls := make([]string, 0, len(tmp_uapps))
		for _, v := range tmp_uapps {
			url := osstool.CreateOSSURL(v.ImgName)
			imgurls = append(imgurls, url)
		}
		userAppraise.ImgUrls = imgurls
	}
	goodsDetailAppraise.UserAppraise = userAppraise
	resp := new(datastruct.ResponseGoodsDetail)
	resp.GoodsCover = goods_cover
	resp.GoodsName = rs_query.Goods.Name
	resp.Percent = rs_query.Percent
	resp.Price = rs_query.Price
	resp.PriceDesc = rs_query.PriceDesc
	resp.RushPrice = rs_query.RushPrice
	resp.SendedOut = rs_query.SendedOut
	resp.Words16 = rs_query.Words16
	resp.ReIcon = tools.CreateGoodsImgUrl(rs_query.RecommendedClass.Icon)
	resp.DetailImgs = goods_detail
	resp.TmpUsers = tmpUsers
	resp.Appraise = goodsDetailAppraise
	return resp, datastruct.NULLError
}

func (handle *DBHandler) UpdateUserAddress(userId int, body *datastruct.ReceiverForSendGoods) datastruct.CodeType {
	engine := handle.mysqlEngine
	sql := fmt.Sprintf("REPLACE INTO user_shipping_address (user_id,linkman,phone,addr,remark)VALUES(%d,'%s','%s','%s','%s')", userId, body.LinkMan, body.PhoneNumber, body.Address, body.Remark)
	_, err := engine.Exec(sql)
	if err != nil {
		log.Debug("UpdateUserAddress err:%v", err.Error())
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

func (handle *DBHandler) GetUserAddress(userId int) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	u_addr := new(datastruct.UserShippingAddress)
	has, err := engine.Where("user_id=?", userId).Get(u_addr)
	if err != nil {
		log.Debug("GetUserAddress err:%v", err.Error())
		return nil, datastruct.GetDataFailed
	}
	resp := new(datastruct.ReceiverForSendGoods)
	if has {
		resp.Address = u_addr.Addr
		resp.LinkMan = u_addr.Linkman
		resp.PhoneNumber = u_addr.Phone
		resp.Remark = u_addr.Remark
	}
	return resp, datastruct.NULLError
}

func (handle *DBHandler) GetAgencyPage(userId int) (*datastruct.ResponseAgencyPage, datastruct.CodeType) {
	engine := handle.mysqlEngine
	user := new(datastruct.UserInfo)
	has, err := engine.Where("id=?", userId).Get(user)
	if err != nil || !has {
		log.Debug("GetAgencyPage GetUser err")
		return nil, datastruct.GetDataFailed
	}
	member := new(datastruct.MemberInfo)
	if user.MemberIdentifier == datastruct.AgencyIdentifier {
		member.Name = datastruct.DefaultMemberIdentifier
	} else {
		memberData := new(datastruct.MemberLevelData)
		has, err = engine.Where("id=?", tools.StringToInt(user.MemberIdentifier)).Get(memberData)
		if err != nil || !has {
			log.Debug("GetAgencyPage MemberName err")
			return nil, datastruct.GetDataFailed
		}
		member.Name = memberData.Name
	}
	params := new(datastruct.AgencyParams)
	has, err = engine.Where("identifier = ?", user.MemberIdentifier).Get(params)
	if err != nil || !has {
		log.Debug("GetAgencyPage get AgencyParams err")
		return nil, datastruct.GetDataFailed
	}
	member.A1Money = params.Agency1MoneyPercent
	member.A2Money = params.Agency2MoneyPercent
	member.A3Money = params.Agency3MoneyPercent

	count := 10
	sql := "select avatar from user_info order by balance_total desc,id asc limit ?"
	results, err := engine.Query(sql, count)
	avatars := make([]string, 0)
	for _, v := range results {
		avatars = append(avatars, string(v["avatar"][:]))
	}
	resp := new(datastruct.ResponseAgencyPage)
	resp.Member = member
	resp.Avatar = avatars
	return resp, datastruct.NULLError
}

func (handle *DBHandler) RemindSendGoods(userId int, number string) datastruct.CodeType {
	engine := handle.mysqlEngine
	order := new(datastruct.OrderInfo)
	order.IsRemind = 1
	_, err := engine.Where("user_id = ? and number = ?", userId, number).Cols("is_remind").Update(order)
	if err != nil {
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

func (handle *DBHandler) AddSuggest(userId int, desc string) datastruct.CodeType {
	engine := handle.mysqlEngine
	scp := new(datastruct.SuggestionComplaintParams)
	has, err := engine.Where("id=?", datastruct.DefaultId).Get(scp)
	if err != nil || !has {
		if err != nil {
			log.Debug("AddSuggest err0:%v", err.Error())
		}
		return datastruct.UpdateDataFailed
	}
	count := scp.SuggestCountForDay
	today, tomorrow := tools.GetTodayTomorrowTime()
	var total int64
	total, err = engine.Where("user_id=? and created_at >= ? and created_at < ?", userId, today, tomorrow).Count(new(datastruct.Suggestion))
	if int(total) >= count {
		return datastruct.UpperLimit
	}
	sgt := new(datastruct.Suggestion)
	sgt.UserId = userId
	sgt.Desc = desc
	sgt.CreatedAt = time.Now().Unix()
	_, err = engine.Insert(sgt)
	if err != nil {
		log.Debug("AddSuggest Insert err:%v", err.Error())
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

func (handle *DBHandler) AddComplaint(userId int, complaintType string, desc string) datastruct.CodeType {
	engine := handle.mysqlEngine
	scp := new(datastruct.SuggestionComplaintParams)
	has, err := engine.Where("id=?", datastruct.DefaultId).Get(scp)
	if err != nil || !has {
		return datastruct.UpdateDataFailed
	}
	count := scp.ComplaintCountForDay
	bl := scp.ComplaintCountForBL
	today, tomorrow := tools.GetTodayTomorrowTime()
	var today_count, total int64
	total, err = engine.Where("user_id=?", userId).Count(new(datastruct.Complaint))

	if int(total) >= bl {
		return datastruct.BlackList
	}
	today_count, err = engine.Where("user_id=? and created_at >= ? and created_at < ?", userId, today, tomorrow).Count(new(datastruct.Complaint))
	if int(today_count) >= count {
		return datastruct.UpperLimit
	}
	cplt := new(datastruct.Complaint)
	cplt.UserId = userId
	cplt.Desc = desc
	cplt.CreatedAt = time.Now().Unix()
	cplt.ComplaintType = complaintType
	_, err = engine.Insert(cplt)
	if err != nil {
		return datastruct.UpdateDataFailed
	}
	if int(total)+1 >= bl {
		return datastruct.BlackList
	}
	return datastruct.NULLError
}

func (handle *DBHandler) IsDrawCashOnApp() bool {
	tf := false
	engine := handle.mysqlEngine
	gcg := new(datastruct.GoldCoinGift)
	has, err := engine.Where("id=?", datastruct.DefaultId).Get(gcg)
	if err != nil || !has {
		return tf
	}
	tf = tools.IntToBool(gcg.IsDrawCashOnlyApp)
	return tf
}

func (handle *DBHandler) GetDownLoadAppGift(userId int) datastruct.CodeType {
	engine := handle.mysqlEngine
	return addGoldGift(true, userId, engine)
}

func (handle *DBHandler) GetRegisterGift(userId int) datastruct.CodeType {
	engine := handle.mysqlEngine
	return addGoldGift(false, userId, engine)
}

func addGoldGift(isDownLoad bool, userId int, engine *xorm.Engine) datastruct.CodeType {
	gcg := new(datastruct.GoldCoinGift)
	has, err := engine.Where("id=?", datastruct.DefaultId).Get(gcg)
	if err != nil || !has {
		return datastruct.GetDataFailed
	}
	user := new(datastruct.UserInfo)
	has, err = engine.Where("id=?", userId).Get(user)
	if err != nil || !has {
		return datastruct.GetDataFailed
	}
	if isDownLoad {
		if user.IsGotDownLoadAppGift == 1 {
			return datastruct.NULLError
		}
	} else {
		if user.IsGotRegisterGift == 1 {
			return datastruct.NULLError
		}
	}
	session := engine.NewSession()
	defer session.Close()
	session.Begin()
	var add_goldcount int64
	var sql string
	var goldType datastruct.GoldChangeType
	if isDownLoad {
		add_goldcount = gcg.DownLoadAppGoldGift
		sql = "update user_info set gold_count = gold_count + ? , is_got_down_load_app_gift = 1  where id = ?"
		goldType = datastruct.DownLoadAppType
	} else {
		add_goldcount = gcg.RegisterGoldGift
		sql = "update user_info set gold_count = gold_count + ? , is_got_register_gift = 1  where id = ?"
		goldType = datastruct.RegisterType
	}
	res, err1 := session.Exec(sql, add_goldcount, userId)
	affected, err2 := res.RowsAffected()
	if err1 != nil || err2 != nil || affected <= 0 {
		rollback("addGoldGift update AddGoldCount err", session)
		return datastruct.UpdateDataFailed
	}
	goldChangeInfo := new(datastruct.GoldChangeInfo)
	goldChangeInfo.UserId = userId
	goldChangeInfo.CreatedAt = time.Now().Unix()
	goldChangeInfo.VarGold = add_goldcount
	goldChangeInfo.ChangeType = goldType
	_, err = session.Insert(goldChangeInfo)
	if err != nil {
		str := fmt.Sprintf("addGoldGift Insert GoldChangeInfo :%s", err.Error())
		rollback(str, session)
		return datastruct.UpdateDataFailed
	}
	err = session.Commit()
	if err != nil {
		str := fmt.Sprintf("addGoldGift Commit :%s", err.Error())
		rollback(str, session)
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

func (handle *DBHandler) UpdateTmpData(goods_id int, now_time int64) {
	engine := handle.mysqlEngine
	num := tools.RandInt(5, 11)
	sql := "update goods set sended_out = sended_out + ? where id = ?"
	_, err := engine.Exec(sql, num, goods_id)
	if err != nil {
		log.Debug("UpdateTmpData  update goods err:%s", err.Error())
		return
	}

	sql = "select id from tmp_data where id not in (select tmp_data.id from tmp_data INNER join tmp_data_for_goods on tmp_data.id = tmp_data_for_goods.tmp_user_id)"
	var results []map[string][]byte
	results, err = engine.Query(sql)
	if err != nil {
		log.Debug("UpdateTmpData query tmp_data err:%s", err.Error())
		return
	}
	id1, id2 := getRandomTmpId(results)
	count := 2
	tmp_arr := make([]*datastruct.TmpDataForGoods, 0, count)
	err = engine.Where("goods_id=?", goods_id).Limit(count, 0).Find(&tmp_arr)
	if err != nil {
		log.Debug("UpdateTmpData Find TmpDataForGoods err:%s", err.Error())
		return
	}
	var affected int64
	if len(tmp_arr) == 0 {
		tmpDataForGoods_1 := new(datastruct.TmpDataForGoods)
		tmpDataForGoods_2 := new(datastruct.TmpDataForGoods)
		tmpDataForGoods_1.GoodsId = goods_id
		tmpDataForGoods_1.TmpUserId = id1
		tmpDataForGoods_1.UpdateAt = now_time
		tmpDataForGoods_2.GoodsId = goods_id
		tmpDataForGoods_2.TmpUserId = id2
		tmpDataForGoods_2.UpdateAt = now_time
		affected, err = engine.Insert(tmpDataForGoods_1, tmpDataForGoods_2)
		if err != nil || affected <= 0 {
			log.Debug("UpdateTmpData Insert TmpDataForGoods err")
			return
		}
	} else if len(tmp_arr) == count {
		tmp_arr[0].TmpUserId = id1
		tmp_arr[0].UpdateAt = now_time
		affected, err = engine.Where("id=?", tmp_arr[0].Id).Cols("tmp_user_id", "update_at").Update(tmp_arr[0])
		if err != nil || affected <= 0 {
			log.Debug("UpdateTmpData update TmpDataForGoods1 err")
			return
		}
		tmp_arr[1].TmpUserId = id2
		tmp_arr[1].UpdateAt = now_time
		affected, err = engine.Where("id=?", tmp_arr[1].Id).Cols("tmp_user_id", "update_at").Update(tmp_arr[1])
		if err != nil || affected <= 0 {
			log.Debug("UpdateTmpData update TmpDataForGoods2 err")
			return
		}
	}
}

func getRandomTmpId(allIds []map[string][]byte) (int, int) {
	length := len(allIds)
	all_index := make([]int, 0, length)
	for i := 0; i < length; i++ {
		all_index = append(all_index, i)
	}
	index_1 := tools.RandInt(0, len(all_index))
	tmp_map := allIds[all_index[index_1]]
	id_1 := tools.StringToInt(string(tmp_map["id"][:]))

	all_index = tools.Remove(all_index, index_1)

	index_2 := tools.RandInt(0, len(all_index))
	tmp_map = allIds[all_index[index_2]]
	id_2 := tools.StringToInt(string(tmp_map["id"][:]))
	return id_1, id_2
}

func (handle *DBHandler) GetGoodsForTmpData() []int {
	engine := handle.mysqlEngine
	goods := make([]*datastruct.Goods, 0)
	err := engine.Where("is_hidden=0").Find(&goods)
	if err != nil {
		log.Debug("GetGoodsForTmpData err:%v", err.Error())
		return nil
	}
	goods_ids := make([]int, 0, len(goods))
	for _, v := range goods {
		goods_ids = append(goods_ids, v.Id)
	}
	return goods_ids
}

func (handle *DBHandler) TruncateTmpDataForGoods() {
	engine := handle.mysqlEngine
	_, err := engine.Exec("truncate table tmp_data_for_goods")
	if err != nil {
		log.Debug("--------------TruncateTmpDataForGoods err:%v", err.Error())
	}
}

func (handle *DBHandler) IsRefreshHomeGoodsData(userId int, classId int) (interface{}, datastruct.CodeType) {
	isRefresh := 0
	engine := handle.mysqlEngine
	uggt := new(datastruct.UserGetHomeGoodsDataTime)
	has, err := engine.Where("user_id=? and class_id=?", userId, classId).Get(uggt)
	if err != nil {
		return nil, datastruct.GetDataFailed
	}
	if has {
		sql := "select update_at from goods_class where id = ?"
		results, err := engine.Query(sql, classId)
		if err != nil {
			return nil, datastruct.GetDataFailed
		}
		update_at := tools.StringToInt64(string(results[0]["update_at"][:]))
		if update_at > uggt.GetDataTime {
			isRefresh = 1
		}
	} else {
		isRefresh = 1
	}
	return isRefresh, datastruct.NULLError
}

func (handle *DBHandler) UpdateGoodsClassTime() {
	engine := handle.mysqlEngine
	sql := "update goods_class set update_at = ? where id > 0"
	engine.Exec(sql, time.Now().Unix())
}

func (handle *DBHandler) GetGoldFromPoster(userId int, gpid int, addr string) (interface{}, datastruct.CodeType) {
	now_time := time.Now().Unix()
	engine := handle.mysqlEngine
	gp := new(datastruct.GoldPoster)
	invalid_tips := "该二维码已经失效，可联系客服扫最新的二维码领取哦"
	err_tips := "系统出错，我们已经在加急处理中了"
	resp := new(datastruct.ResponseGoldFromPoster)
	has, err := engine.Where("id=?", gpid).Get(gp)
	if err != nil {
		resp.Tips = err_tips
		resp.Succeed = 0
		return resp, datastruct.GetDataFailed
	}
	resp.Addr = addr
	if has {
		if now_time < gp.StartTime {
			resp.Succeed = 0
			resp.Tips = "活动还未开始，请在活动开始后来领取"
		} else if now_time > gp.EndTime {
			resp.Succeed = 0
			resp.Tips = invalid_tips
		} else {
			has, err = engine.Where("user_id=? and gold_poster_id=?", userId, gpid).Get(new(datastruct.SaveUserGetGoldPoster))
			if has {
				resp.Succeed = 0
				resp.Tips = "亲，你今天已经领取过了，明天再来吧"
			} else {
				resp.Succeed = 1
				resp.Tips = fmt.Sprintf("恭喜你获得%d个游戏币", gp.GoldCount)
				session := engine.NewSession()
				defer session.Close()
				session.Begin()
				add_goldcount := gp.GoldCount
				sql := "update user_info set gold_count = gold_count + ? where id = ?"
				goldType := datastruct.GoldPosterType
				res, err1 := session.Exec(sql, add_goldcount, userId)
				affected, err2 := res.RowsAffected()
				if err1 != nil || err2 != nil || affected <= 0 {
					rollbackError("GetGoldFromPoster update AddGoldCount err", session)
					resp.Succeed = 0
					resp.Tips = err_tips
					return resp, datastruct.UpdateDataFailed
				}
				goldChangeInfo := new(datastruct.GoldChangeInfo)
				goldChangeInfo.UserId = userId
				goldChangeInfo.CreatedAt = now_time
				goldChangeInfo.VarGold = int64(add_goldcount)
				goldChangeInfo.ChangeType = goldType
				_, err = session.Insert(goldChangeInfo)
				if err != nil {
					str := fmt.Sprintf("GetGoldFromPoster Insert GoldChangeInfo :%s", err.Error())
					rollbackError(str, session)
					resp.Succeed = 0
					resp.Tips = err_tips
					return resp, datastruct.UpdateDataFailed
				}
				sgp := new(datastruct.SaveUserGetGoldPoster)
				sgp.GoldPosterId = gpid
				sgp.UserId = userId
				_, err = session.Insert(sgp)
				if err != nil {
					str := fmt.Sprintf("GetGoldFromPoster Insert SaveUserGetGoldPoster :%s", err.Error())
					rollbackError(str, session)
					resp.Succeed = 0
					resp.Tips = err_tips
					return resp, datastruct.UpdateDataFailed
				}
				err = session.Commit()
				if err != nil {
					str := fmt.Sprintf("DBHandler->GetGoldFromPoster Commit :%s", err.Error())
					rollbackError(str, session)
					resp.Tips = err_tips
					resp.Succeed = 0
					return resp, datastruct.UpdateDataFailed
				}
			}
		}
	} else {
		resp.Succeed = 0
		resp.Tips = invalid_tips
	}
	return resp, datastruct.NULLError
}

func (handle *DBHandler) GetRandomLotteryList() (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	sql := "select img_name from random_lottery_goods order by created_at desc"
	results, err := engine.Query(sql)
	if err != nil {
		return nil, datastruct.GetDataFailed
	}
	list := make([]string, 0, len(results))
	for _, v := range results {
		list = append(list, tools.CreateGoodsImgUrl(string(v["img_name"][:])))
	}
	return list, datastruct.NULLError
}

// var valuesSlice = make([]interface{}, len(cols))
// has, err := engine.Where("id = ?", id).Cols(cols...).Get(&valuesSlice)
