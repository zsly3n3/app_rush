package db

import (
	"app/conf"
	"app/datastruct"
	"app/log"
	"app/osstool"
	"app/tools"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"

	"github.com/go-xorm/xorm"
)

func (handle *DBHandler) WebLogin(body *datastruct.WebLoginBody) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	user := new(datastruct.WebUser)
	has, err := engine.Where("login_name=? and pwd=?", body.Account, body.Pwd).Get(user)
	if err != nil || !has {
		return nil, datastruct.LoginFailed
	}
	p_user := new(datastruct.WebResponsePermissionUser)
	p_user.Name = user.Name
	p_user.Token = user.Token
	var permission []*datastruct.MasterInfo
	if user.RoleId == datastruct.AdminLevelID {
		permission = getAllMenu(engine)
	} else {
		permission = make([]*datastruct.MasterInfo, 0)
		sql := "select mm.id,mm.name from web_permission wp join secondary_menu sm on wp.secondary_id = sm.id join master_menu mm on mm.id = sm.master_id where user_id = ? GROUP BY mm.id order by mm.id asc"
		rs, _ := engine.Query(sql, user.Id)
		for _, v := range rs {
			master_id := tools.StringToInt(string(v["id"][:]))
			master_name := string(v["name"][:])
			m_info := new(datastruct.MasterInfo)
			m_info.MasterId = master_id
			m_info.Name = master_name
			secondary := make([]*datastruct.SecondaryInfo, 0)
			sql := "select sm.id,sm.name from web_permission wp join secondary_menu sm on wp.secondary_id = sm.id join master_menu mm on mm.id = sm.master_id where wp.user_id = ? and mm.id = ? order by sm.id asc"
			rs, _ = engine.Query(sql, user.Id, master_id)
			for _, v := range rs {
				secondaryInfo := new(datastruct.SecondaryInfo)
				secondaryInfo.Name = string(v["name"][:])
				secondaryInfo.SecondaryId = tools.StringToInt(string(v["id"][:]))
				secondary = append(secondary, secondaryInfo)
			}
			m_info.Secondary = secondary
			permission = append(permission, m_info)
		}
	}
	p_user.Permission = permission
	return p_user, datastruct.NULLError
}

func (handle *DBHandler) EditGoods(body *datastruct.EditGoodsBody) datastruct.CodeType {
	engine := handle.mysqlEngine
	session := engine.NewSession()
	defer session.Close()
	session.Begin()
	goods := new(datastruct.Goods)
	goods.Brand = body.Brand
	goods.GoodsClassId = body.Classid
	goods.GoodsDesc = body.Goodsdesc
	goods.Name = body.Name
	goods.Price = body.Price
	goods.PriceDesc = body.Pricedesc
	goods.RushPrice = body.Rushprice
	goods.RushPriceDesc = body.Rushpricedesc
	goods.SortId = body.Sortid
	goods.IsHidden = body.IsHidden
	goods.Count = body.Count
	goods.Type = body.Type
	goods.OriginalPrice = body.OriginalPrice
	goods.PostAge = body.PostAge
	goods.ImgName = ""
	goods.SendedOut = body.SendedOut
	goods.ReClassid = body.ReClassid //可为0
	goods.Words16 = body.Words16
	goods.Percent = body.Percent
	var has bool
	var err error
	var affected int64
	var isUpdate bool
	rewardPool := new(datastruct.GoodsRewardPool)
	rewardPool.LimitAmount = body.LimitAmount
	if body.Goodsid <= 0 {
		affected, err = session.Insert(goods)
		if err != nil || affected <= 0 {
			rollback("DBHandler->EditGoods InsertGoods err", session)
			return datastruct.UpdateDataFailed
		}
		rewardPool.Current = 0
		rewardPool.GoodsId = goods.Id
		affected, err = session.Insert(rewardPool)
		if err != nil || affected <= 0 {
			rollback("DBHandler->EditGoods Insert rewardPool err", session)
			return datastruct.UpdateDataFailed
		}
		isUpdate = false
	} else {
		affected, err = session.Where("id=?", body.Goodsid).Cols("percent", "words16", "re_classid", "sended_out", "name", "price", "price_desc", "rush_price", "rush_price_desc", "goods_desc", "brand", "sort_id", "goods_class_id", "is_hidden", "count", "type", "original_price", "post_age").Update(goods)
		if err != nil {
			rollback("DBHandler->EditGoods UpdateGoods err", session)
			return datastruct.UpdateDataFailed
		}
		affected, err = session.Where("goods_id=?", body.Goodsid).Cols("limit_amount").Update(rewardPool)
		if err != nil {
			rollback("DBHandler->EditGoods UpdateGoods err", session)
			return datastruct.UpdateDataFailed
		}
		goods.Id = body.Goodsid
		isUpdate = true
	}
	length := len(body.Base64str)
	if body.Goodsid > 0 {
		tmp := make([]*datastruct.GoodsImgs, 0)
		err = session.Where("goods_id = ?", body.Goodsid).Find(&tmp)
		if err != nil {
			rollback("DBHandler->EditGoods GetGoodsImgs err:"+err.Error(), session)
			return datastruct.UpdateDataFailed
		}
		var isDelete bool
		for _, v_imgs := range tmp {
			isDelete = true
			for _, v := range body.Base64str {
				if strings.Contains(v, conf.Server.Domain) {
					str_arr := strings.Split(v, "/")
					filename := str_arr[len(str_arr)-1]
					if v_imgs.ImgName == filename {
						isDelete = false
					}
				}
			}
			if isDelete {
				deleteFile(tools.GetImgPath() + v_imgs.ImgName)
			}
		}
		for i, v := range body.Base64str {
			if strings.Contains(v, conf.Server.Domain) {
				tmp := new(datastruct.GoodsImgs)
				has, err = session.Where("goods_id=? and img_index=?", goods.Id, i).Get(tmp)
				if err != nil {
					rollback("DBHandler->EditGoods 0 get err:"+err.Error(), session)
					return datastruct.GetDataFailed
				}
				str_arr := strings.Split(v, "/")
				filename := str_arr[len(str_arr)-1]
				if has {
					tmp.ImgName = filename
					_, err = session.Where("goods_id=? and img_index=?", goods.Id, i).Cols("img_name").Update(tmp)
					if err != nil {
						rollback("DBHandler->EditGoods 0 Update GoodsImgs err:"+err.Error(), session)
						return datastruct.UpdateDataFailed
					}
				} else {
					goodsimgs := new(datastruct.GoodsImgs)
					goodsimgs.GoodsId = goods.Id
					goodsimgs.ImgName = filename
					goodsimgs.ImgIndex = i
					affected, err = session.Insert(goodsimgs)
					if err != nil || affected <= 0 {
						rollback("DBHandler->EditGoods 0 Insert GoodsImgs err", session)
						return datastruct.UpdateDataFailed
					}
				}
			}
		}
		_, err = session.Where("goods_id=? and img_index >= ?", body.Goodsid, length).Delete(new(datastruct.GoodsImgs))
		if err != nil {
			rollback("DBHandler->EditGoods DeleteGoodsImgs err:"+err.Error(), session)
			return datastruct.UpdateDataFailed
		}
	}

	for i := 0; i < length; i++ {
		v := body.Base64str[i]
		if !strings.Contains(v, conf.Server.Domain) {
			imgName := fmt.Sprintf("%s.png", tools.UniqueId())
			path := tools.GetImgPath() + imgName
			arr_str := strings.Split(v, ",")
			var isError bool
			if len(arr_str) > 1 {
				isError = tools.CreateImgFromBase64(&arr_str[1], path)
			} else {
				isError = tools.CreateImgFromBase64(&arr_str[0], path)
			}
			if isError {
				rollback("DBHandler->EditGoods CreateImgFromBase64 err", session)
				return datastruct.UpdateDataFailed
			}
			if i == 0 {
				goods.ImgName = imgName
				affected, err = session.Where("id=?", goods.Id).Cols("img_name").Update(goods)
				if err != nil || affected <= 0 {
					rollback("DBHandler->EditGoods UpdateGoods img_name err", session)
					return datastruct.UpdateDataFailed
				}
			}
			tmp := new(datastruct.GoodsImgs)
			has, err = session.Where("goods_id=? and img_index=?", goods.Id, i).Get(tmp)
			if err != nil {
				rollback("DBHandler->EditGoods get err:"+err.Error(), session)
				return datastruct.GetDataFailed
			}
			if has {
				tmp.ImgName = imgName
				_, err = session.Where("goods_id=? and img_index=?", goods.Id, i).Cols("img_name").Update(tmp)
				if err != nil {
					rollback("DBHandler->EditGoods Update GoodsImgs err:"+err.Error(), session)
					return datastruct.UpdateDataFailed
				}
			} else {
				goodsimgs := new(datastruct.GoodsImgs)
				goodsimgs.GoodsId = goods.Id
				goodsimgs.ImgName = imgName
				goodsimgs.ImgIndex = i
				affected, err = session.Insert(goodsimgs)
				if err != nil || affected <= 0 {
					rollback("DBHandler->EditGoods Insert GoodsImgs err", session)
					return datastruct.UpdateDataFailed
				}
			}
		}
	}

	for k, v := range body.LevelData {
		paymode := new(datastruct.PayModeRougeGame)
		paymode.Level = k + 1
		paymode.GoodsId = goods.Id
		paymode.RougeCount = v.Count
		paymode.Difficulty = v.Difficulty
		paymode.GameTime = v.Time
		if paymode.RougeCount <= 0 || paymode.Difficulty < 0 || paymode.GameTime <= 0 {
			rollback("DBHandler->EditGoods PayModeRougeGame ParamError", session)
			return datastruct.ParamError
		}
		if !isUpdate {
			affected, err = session.Insert(paymode)
			if err != nil || affected <= 0 {
				rollback("DBHandler->EditGoods Insert PayModeRougeGame err", session)
				return datastruct.UpdateDataFailed
			}
		} else {
			affected, err = session.Where("level=? and goods_id=?", paymode.Level, paymode.GoodsId).Cols("rouge_count", "difficulty", "game_time").Update(paymode)
			if err != nil {
				rollback("DBHandler->EditGoods Update PayModeRougeGame err", session)
				return datastruct.UpdateDataFailed
			}
		}
	}

	err = session.Commit()
	if err != nil {
		str := fmt.Sprintf("DBHandler->EditGoods Commit :%s", err.Error())
		rollback(str, session)
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

type webGoodsData struct {
	datastruct.Goods            `xorm:"extends"`
	datastruct.GoodsClass       `xorm:"extends"`
	datastruct.GoodsRewardPool  `xorm:"extends"`
	datastruct.RecommendedClass `xorm:"extends"`
}

type webFailedData struct {
	datastruct.Goods                  `xorm:"extends"`
	datastruct.PayModeRougeGameFailed `xorm:"extends"`
	datastruct.PayModeRougeGame       `xorm:"extends"`
}

func (handle *DBHandler) WebGetGoods(body *datastruct.WebGetGoodsBody) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	goods := make([]*webGoodsData, 0)
	start := (body.PageIndex - 1) * body.PageSize
	limit := body.PageSize
	query := ""
	args := make([]interface{}, 0)
	if body.Classid <= 0 {
		query = "goods_class.id > ?"
		args = append(args, 0)
	} else {
		query = "goods_class.id = ?"
		args = append(args, body.Classid)
	}
	if body.IsHidden >= 2 {
		query += " and goods.is_hidden >= ?"
		args = append(args, 0)
	} else {
		query += " and goods.is_hidden = ?"
		args = append(args, body.IsHidden)
	}
	if body.Name != "" {
		query += " and goods.name like ?"
		args = append(args, "%"+body.Name+"%")
	}
	tmp_total := new(datastruct.Goods)

	count, _ := engine.Table("goods").Join("INNER", "goods_class", "goods_class.id = goods.goods_class_id").Join("INNER", "goods_reward_pool", "goods.id = goods_reward_pool.goods_id").Join("Left", "recommended_class", "recommended_class.id = goods.re_classid").Where(query, args...).Desc("goods_class.sort_id").Asc("goods.id").Count(tmp_total)
	engine.Table("goods").Join("INNER", "goods_class", "goods_class.id = goods.goods_class_id").Join("INNER", "goods_reward_pool", "goods.id = goods_reward_pool.goods_id").Join("Left", "recommended_class", "recommended_class.id = goods.re_classid").Where(query, args...).Desc("goods_class.sort_id").Asc("goods.id").Limit(limit, start).Find(&goods)
	resp_goods := make([]*datastruct.WebResponseGoodsData, 0, len(goods))
	for _, v := range goods {
		resp_good := new(datastruct.WebResponseGoodsData)
		goodsimgs := make([]*datastruct.GoodsImgs, 0)
		engine.Where("goods_id=?", v.Goods.Id).Asc("img_index").Find(&goodsimgs)
		imgurls := make([]string, 0, len(goodsimgs))
		for _, img_v := range goodsimgs {
			imgurls = append(imgurls, tools.CreateGoodsImgUrl(img_v.ImgName))
		}
		resp_good.ImgUrls = imgurls
		resp_good.Brand = v.Goods.Brand
		resp_good.Count = v.Goods.Count
		resp_good.GoodsClass = v.GoodsClass.Name
		resp_good.GoodsDesc = v.Goods.GoodsDesc
		resp_good.Id = v.Goods.Id
		resp_good.IsHidden = v.Goods.IsHidden
		resp_good.LimitAmount = v.GoodsRewardPool.LimitAmount
		resp_good.Name = v.Goods.Name
		resp_good.OriginalPrice = v.Goods.OriginalPrice
		resp_good.PostAge = v.Goods.PostAge
		resp_good.Price = v.Goods.Price
		resp_good.PriceDesc = v.Goods.PriceDesc
		resp_good.RushPrice = v.Goods.RushPrice
		resp_good.RushPriceDesc = v.Goods.RushPriceDesc
		resp_good.SortId = v.Goods.SortId
		resp_good.Type = v.Goods.Type
		resp_good.SendedOut = v.Goods.SendedOut
		resp_good.Words16 = v.Goods.Words16
		resp_good.Percent = v.Goods.Percent
		if v.RecommendedClass.Icon != "" {
			resp_good.ReClassIcon = tools.CreateGoodsImgUrl(v.RecommendedClass.Icon)
		}

		sql1 := "select count(*) from pay_mode_rouge_game_failed a INNER JOIN pay_mode_rouge_game b ON a.pay_mode_rouge_game_id=b.id"
		sql2 := " inner join goods c on b.goods_id = c.id where c.id=%d"
		sql := fmt.Sprintf(sql1+sql2, v.Goods.Id)
		results, _ := engine.Query(sql)
		strTotal := string(results[0]["count(*)"][:])
		resp_good.FailedTotal = tools.StringToInt64(strTotal) * v.Goods.RushPrice

		paymode_info := make([]*datastruct.PayModeRougeGame, 0, datastruct.MaxLevel)
		rs_level := make([]*datastruct.GoodsLevelInfo, 0, datastruct.MaxLevel)
		engine.Where("goods_id=?", v.Goods.Id).Asc("level").Find(&paymode_info)
		for _, v := range paymode_info {
			level := new(datastruct.GoodsLevelInfo)
			level.Count = v.RougeCount
			level.Difficulty = v.Difficulty
			level.Time = v.GameTime
			rs_level = append(rs_level, level)
		}
		resp_good.Level = rs_level
		resp_goods = append(resp_goods, resp_good)
	}
	resp := new(datastruct.WebResponseGoods)
	resp.List = resp_goods
	resp.Total = int(count)
	return resp, datastruct.NULLError
}

func (handle *DBHandler) EditDomain(body *datastruct.WebEditDomainBody) datastruct.CodeType {
	engine := handle.mysqlEngine
	session := engine.NewSession()
	defer session.Close()
	session.Begin()

	now_time := time.Now().Unix()
	sql := fmt.Sprintf("REPLACE INTO entry_addr (url,page_url,created_at)VALUES('%s','%s',%d)", body.EntryDomain, body.EntryPageUrl, now_time)
	_, err := session.Exec(sql)
	if err != nil {
		str := fmt.Sprintf("DBHandler->EditDomain REPLACE INTO entry_addr :%s", err.Error())
		rollback(str, session)
		return datastruct.UpdateDataFailed
	}

	sql = fmt.Sprintf("REPLACE INTO auth_addr (url,created_at)VALUES('%s',%d)", body.AuthDomain, now_time)
	_, err = session.Exec(sql)
	if err != nil {
		str := fmt.Sprintf("DBHandler->EditDomain REPLACE INTO auth_addr :%s", err.Error())
		rollback(str, session)
		return datastruct.UpdateDataFailed
	}

	sql = fmt.Sprintf("REPLACE INTO app_addr (url,created_at)VALUES('%s',%d)", body.AppDomain, now_time)
	_, err = session.Exec(sql)
	if err != nil {
		str := fmt.Sprintf("DBHandler->EditDomain REPLACE INTO app_addr :%s", err.Error())
		rollback(str, session)
		return datastruct.UpdateDataFailed
	}

	sql = fmt.Sprintf("REPLACE INTO app_download_addr (ios_url,android_url,down_load_url,direct_down_load_url,created_at)VALUES('%s','%s','%s','%s',%d)", body.IOSApp, body.AndroidApp, body.DownLoadUrl, body.DirectDownLoadUrl, now_time)
	_, err = session.Exec(sql)
	if err != nil {
		str := fmt.Sprintf("DBHandler->EditDomain REPLACE INTO app_download_addr :%s", err.Error())
		rollback(str, session)
		return datastruct.UpdateDataFailed
	}

	err = session.Commit()
	if err != nil {
		str := fmt.Sprintf("DBHandler->EditDomain Commit :%s", err.Error())
		rollback(str, session)
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

func (handle *DBHandler) GetDomain() (*datastruct.WebResponseEditDomain, datastruct.CodeType) {
	engine := handle.mysqlEngine
	resp := new(datastruct.WebResponseEditDomain)

	entry := new(datastruct.EntryAddr)
	has, err := engine.Desc("created_at").Limit(1, 0).Get(entry)
	if err == nil && has {
		resp.EntryDomain = entry.Url
		resp.EntryPageUrl = entry.PageUrl
	}

	auth := new(datastruct.AuthAddr)
	has, err = engine.Desc("created_at").Limit(1, 0).Get(auth)
	if err == nil && has {
		resp.AuthDomain = auth.Url
	}

	app := new(datastruct.AppAddr)
	has, err = engine.Desc("created_at").Limit(1, 0).Get(app)
	if err == nil && has {
		resp.AppDomain = app.Url
	}

	download := new(datastruct.AppDownloadAddr)
	has, err = engine.Desc("created_at").Limit(1, 0).Get(download)
	if err == nil && has {
		resp.IOSApp = download.IosUrl
		resp.AndroidApp = download.AndroidUrl
		resp.DownLoadUrl = download.DownLoadUrl
		resp.DirectDownLoadUrl = download.DirectDownLoadUrl
	}

	return resp, datastruct.NULLError
}

func (handle *DBHandler) GetBlackListJump() interface{} {
	resp := new(datastruct.BlackListJumpBody)
	engine := handle.mysqlEngine
	table := new(datastruct.BlackListJump)
	has, err := engine.Desc("created_at").Limit(1, 0).Get(table)
	if err == nil && has {
		resp.BLJumpTo = table.BLJumpTo
		resp.PCJumpTo = table.PCJumpTo
	}
	return resp
}

func (handle *DBHandler) EditBlackListJump(body *datastruct.BlackListJumpBody) datastruct.CodeType {
	engine := handle.mysqlEngine
	now_time := time.Now().Unix()
	sql := fmt.Sprintf("REPLACE INTO black_list_jump (b_l_jump_to,p_c_jump_to,created_at)VALUES('%s','%s',%d)", body.BLJumpTo, body.PCJumpTo, now_time)
	_, err := engine.Exec(sql)
	if err != nil {
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

func (handle *DBHandler) getOrderCount(sql string) int {
	engine := handle.mysqlEngine
	results, _ := engine.Query(sql)
	total := tools.StringToInt(string(results[0]["count(*)"][:]))
	return total
}

func (handle *DBHandler) getRushOrderAmount(sql string) int64 {
	engine := handle.mysqlEngine
	results, _ := engine.Query(sql)
	var amount int64
	amount = 0
	str_sum := string(results[0]["sum(g.rush_price)"][:])
	if str_sum != "" {
		amount += tools.StringToInt64(str_sum)
	}
	return amount
}

func (handle *DBHandler) getPurchaseOrderAmount(sql string) int64 {
	engine := handle.mysqlEngine
	results, _ := engine.Query(sql)
	var amount int64
	amount = 0
	str_sum := string(results[0]["sum(g.price)"][:])
	if str_sum != "" {
		amount += tools.StringToInt64(str_sum)
	}
	return amount
}

const count_sql1 = "select count(*) from save_game_info s inner join user_info u on s.user_id = u.id inner join goods g on g.id = s.goods_id"
const count_sql2 = "select count(*) from order_info o inner join user_info u on o.user_id = u.id inner join goods g on g.id = o.goods_id"
const count_sql3 = "select count(*) from pay_mode_rouge_game_failed pf inner join user_info u on pf.user_id = u.id inner join pay_mode_rouge_game p on pf.pay_mode_rouge_game_id = p.id inner join goods g on g.id = p.goods_id"
const price_sum_sql1 = "select sum(g.price) from save_game_info s inner join user_info u on s.user_id = u.id inner join goods g on g.id = s.goods_id"
const price_sum_sql2 = "select sum(g.price) from order_info o inner join user_info u on o.user_id = u.id inner join goods g on g.id = o.goods_id"
const price_sum_sql3 = "select sum(g.price) from pay_mode_rouge_game_failed pf inner join user_info u on pf.user_id = u.id inner join pay_mode_rouge_game p on pf.pay_mode_rouge_game_id = p.id inner join goods g on g.id = p.goods_id"

const rush_sum_sql1 = "select sum(g.rush_price) from save_game_info s inner join user_info u on s.user_id = u.id inner join goods g on g.id = s.goods_id"
const rush_sum_sql2 = "select sum(g.rush_price) from order_info o inner join user_info u on o.user_id = u.id inner join goods g on g.id = o.goods_id"
const rush_sum_sql3 = "select sum(g.rush_price) from pay_mode_rouge_game_failed pf inner join user_info u on pf.user_id = u.id inner join pay_mode_rouge_game p on pf.pay_mode_rouge_game_id = p.id inner join goods g on g.id = p.goods_id"

func (handle *DBHandler) getAllRushOrderValue() (int, int64, int, int64) {

	total := handle.getOrderCount(count_sql1)
	total += handle.getOrderCount(count_sql2 + " where o.is_purchase = 0")
	total += handle.getOrderCount(count_sql3)

	amount := handle.getRushOrderAmount(rush_sum_sql1)
	amount += handle.getRushOrderAmount(rush_sum_sql2 + " where o.is_purchase = 0")
	amount += handle.getRushOrderAmount(rush_sum_sql3)

	today_unix, tomorrow_unix := tools.GetTodayTomorrowTime()

	today_count_sql1 := fmt.Sprintf(count_sql1+" where s.created_at >= %d and s.created_at < %d", today_unix, tomorrow_unix)
	today_count_sql2 := fmt.Sprintf(count_sql2+" where o.is_purchase = 0 and o.created_at >= %d and o.created_at < %d", today_unix, tomorrow_unix)
	today_count_sql3 := fmt.Sprintf(count_sql3+" where pf.created_at >= %d and pf.created_at < %d", today_unix, tomorrow_unix)

	today_total := handle.getOrderCount(today_count_sql1)
	today_total += handle.getOrderCount(today_count_sql2)
	today_total += handle.getOrderCount(today_count_sql3)

	today_sum_sql1 := fmt.Sprintf(rush_sum_sql1+" where s.created_at >= %d and s.created_at < %d", today_unix, tomorrow_unix)
	today_sum_sql2 := fmt.Sprintf(rush_sum_sql2+" where o.is_purchase = 0 and o.created_at >= %d and o.created_at < %d", today_unix, tomorrow_unix)
	today_sum_sql3 := fmt.Sprintf(rush_sum_sql3+" where pf.created_at >= %d and pf.created_at < %d", today_unix, tomorrow_unix)

	today_amount := handle.getRushOrderAmount(today_sum_sql1)
	today_amount += handle.getRushOrderAmount(today_sum_sql2)
	today_amount += handle.getRushOrderAmount(today_sum_sql3)

	return total, amount, today_total, today_amount
}

func (handle *DBHandler) GetRushOrder(body *datastruct.GetRushOrderBody) (interface{}, datastruct.CodeType) {
	resp := new(datastruct.WebResponseRushOrderInfo)
	resp.Total, resp.TotalAmount, resp.TodayCount, resp.TodayAmount = handle.getAllRushOrderValue()

	query := " where 1=1"
	if body.UserName != "" {
		query += " and u.nick_name like " + "'%" + body.UserName + "%'"
	}
	if body.GoodsName != "" {
		query += " and g.name like " + "'%" + body.GoodsName + "%'"
	}
	if body.EndTime > 0 && body.StartTime > 0 && body.EndTime > body.StartTime {
		switch body.State {
		case datastruct.RushPayed:
			query += " and s.created_at >= " + tools.Int64ToString(body.StartTime) + " and s.created_at < " + tools.Int64ToString(body.EndTime)
		case datastruct.RushFinishedApply:
			query += " and o.created_at >= " + tools.Int64ToString(body.StartTime) + " and o.created_at < " + tools.Int64ToString(body.EndTime)
		case datastruct.RushFinishedNotApply:
			query += " and o.created_at >= " + tools.Int64ToString(body.StartTime) + " and o.created_at < " + tools.Int64ToString(body.EndTime)
		case datastruct.RushFinishedFailed:
			query += " and pf.created_at >= " + tools.Int64ToString(body.StartTime) + " and pf.created_at < " + tools.Int64ToString(body.EndTime)
		}
	}

	var current_count_sql string
	var current_amount_sql string
	switch body.State {
	case datastruct.RushPayed:
		current_count_sql = count_sql1 + query
		current_amount_sql = rush_sum_sql1 + query
	case datastruct.RushFinishedApply:
		query = query + " and o.order_state = 1 and o.is_purchase = 0"
		current_count_sql = count_sql2 + query
		current_amount_sql = rush_sum_sql2 + query
	case datastruct.RushFinishedNotApply:
		query = query + " and o.order_state = 0 and o.is_purchase = 0"
		current_count_sql = count_sql2 + query
		current_amount_sql = rush_sum_sql2 + query
	case datastruct.RushFinishedFailed:
		current_count_sql = count_sql3 + query
		current_amount_sql = rush_sum_sql3 + query
	}

	current_total := handle.getOrderCount(current_count_sql)
	current_amount := handle.getRushOrderAmount(current_amount_sql)
	resp.CurrentTotal = current_total
	resp.CurrentTotalAmount = current_amount

	engine := handle.mysqlEngine
	start := (body.PageIndex - 1) * body.PageSize
	limit := body.PageSize
	limitStr := fmt.Sprintf(" LIMIT %d,%d", start, limit)

	var sql string
	var orderby string
	switch body.State {
	case datastruct.RushPayed:
		orderby = " ORDER BY s.created_at desc"
		sql = "select u.nick_name,u.avatar,g.img_name,g.name,g.rush_price,s.created_at from save_game_info s inner join user_info u on s.user_id = u.id inner join goods g on g.id = s.goods_id"
	case datastruct.RushFinishedApply:
		fallthrough
	case datastruct.RushFinishedNotApply:
		orderby = " ORDER BY o.created_at desc"
		sql = "select u.nick_name,u.avatar,g.img_name,g.name,g.rush_price,o.created_at from order_info o inner join user_info u on o.user_id = u.id inner join goods g on g.id = o.goods_id"
	case datastruct.RushFinishedFailed:
		orderby = " ORDER BY pf.created_at desc"
		sql = "select u.nick_name,u.avatar,g.img_name,g.name,g.rush_price,pf.created_at,p.level,pf.rouge_number from pay_mode_rouge_game_failed pf inner join user_info u on pf.user_id = u.id inner join pay_mode_rouge_game p on pf.pay_mode_rouge_game_id = p.id inner join goods g on g.id = p.goods_id"
	}
	sql = sql + query + orderby + limitStr
	results, _ := engine.Query(sql)
	list := make([]*datastruct.RushOrderInfo, 0, len(results))
	for _, v := range results {
		orderInfo := new(datastruct.RushOrderInfo)
		orderInfo.UserName = string(v["nick_name"][:])
		orderInfo.Avatar = string(v["avatar"][:])
		orderInfo.GoodsName = string(v["name"][:])
		orderInfo.RushPrice = tools.StringToInt64(string(v["rush_price"][:]))
		orderInfo.Time = tools.StringToInt64(string(v["created_at"][:]))
		imgName := string(v["img_name"][:])
		orderInfo.GoodsImgUrl = tools.CreateGoodsImgUrl(imgName)
		orderInfo.State = body.State
		var mark string
		switch body.State {
		case datastruct.RushPayed:
			mark = "当前用户的存档,闯关ing"
		case datastruct.RushFinishedApply:
			mark = "闯关成功,用户已申请发货"
		case datastruct.RushFinishedNotApply:
			mark = "闯关成功,用户未申请发货"
		case datastruct.RushFinishedFailed:
			level := string(v["level"][:])
			rouge_count := string(v["rouge_number"][:])
			mark = "闯关失败,在关卡第" + level + "关,第" + rouge_count + "只时失败(有三种失败场景:时间结束,主动退出,没插入成功)"
		}
		orderInfo.ReMark = mark
		list = append(list, orderInfo)
	}
	resp.List = list
	return resp, datastruct.NULLError
}

func (handle *DBHandler) getAllPurchaseOrderValue() (int, int64, int, int64) {

	total := handle.getOrderCount(count_sql2 + " where o.is_purchase = 1")
	amount := handle.getPurchaseOrderAmount(price_sum_sql2 + " where o.is_purchase = 1")

	today_unix, tomorrow_unix := tools.GetTodayTomorrowTime()

	today_count_sql := fmt.Sprintf(count_sql2+" where o.is_purchase = 1 and o.created_at >= %d and o.created_at < %d", today_unix, tomorrow_unix)
	today_total := handle.getOrderCount(today_count_sql)

	today_sum_sql := fmt.Sprintf(price_sum_sql2+" where o.is_purchase = 1 and o.created_at >= %d and o.created_at < %d", today_unix, tomorrow_unix)
	today_amount := handle.getPurchaseOrderAmount(today_sum_sql)

	return total, amount, today_total, today_amount
}

func (handle *DBHandler) GetPurchaseOrder(body *datastruct.GetPurchaseBody) (interface{}, datastruct.CodeType) {
	resp := new(datastruct.WebResponsePurchaseOrderInfo)
	resp.Total, resp.TotalAmount, resp.TodayCount, resp.TodayAmount = handle.getAllPurchaseOrderValue()

	query := " where 1=1"
	if body.UserName != "" {
		query += " and u.nick_name like " + "'%" + body.UserName + "%'"
	}
	if body.GoodsName != "" {
		query += " and g.name like " + "'%" + body.GoodsName + "%'"
	}
	if body.EndTime > 0 && body.StartTime > 0 && body.EndTime > body.StartTime {
		query += " and o.created_at >= " + tools.Int64ToString(body.StartTime) + " and o.created_at < " + tools.Int64ToString(body.EndTime)
	}
	var current_count_sql string
	var current_amount_sql string

	switch body.State {
	case datastruct.Apply:
		query = query + " and o.order_state = 1 and o.is_purchase = 1"
		current_count_sql = count_sql2 + query
		current_amount_sql = price_sum_sql2 + query
	case datastruct.NotApply:
		query = query + " and o.order_state = 0 and o.is_purchase = 1"
		current_count_sql = count_sql2 + query
		current_amount_sql = price_sum_sql2 + query
	default:
		query = query + " and o.is_purchase = 1"
		current_count_sql = count_sql2 + query
		current_amount_sql = price_sum_sql2 + query
	}

	current_total := handle.getOrderCount(current_count_sql)
	current_amount := handle.getPurchaseOrderAmount(current_amount_sql)
	resp.CurrentTotal = current_total
	resp.CurrentTotalAmount = current_amount

	engine := handle.mysqlEngine
	start := (body.PageIndex - 1) * body.PageSize
	limit := body.PageSize
	limitStr := fmt.Sprintf(" LIMIT %d,%d", start, limit)

	sql := "select u.nick_name,u.avatar,g.img_name,g.name,g.price,o.created_at,o.order_state from order_info o inner join user_info u on o.user_id = u.id inner join goods g on g.id = o.goods_id"
	orderby := " ORDER BY o.created_at desc"
	sql = sql + query + orderby + limitStr
	results, _ := engine.Query(sql)
	list := make([]*datastruct.PurchaseOrderInfo, 0, len(results))
	for _, v := range results {
		orderInfo := new(datastruct.PurchaseOrderInfo)
		orderInfo.UserName = string(v["nick_name"][:])
		orderInfo.Avatar = string(v["avatar"][:])
		orderInfo.GoodsName = string(v["name"][:])
		orderInfo.Price = tools.StringToInt64(string(v["price"][:]))
		orderInfo.Time = tools.StringToInt64(string(v["created_at"][:]))
		imgName := string(v["img_name"][:])
		orderInfo.GoodsImgUrl = tools.CreateGoodsImgUrl(imgName)
		orderInfo.State = datastruct.OrderType(tools.StringToInt(string(v["order_state"][:])))
		var mark string
		switch orderInfo.State {
		case datastruct.Apply:
			mark = "购买成功,用户已申请发货"
		case datastruct.NotApply:
			mark = "购买成功,用户未申请发货"
		}
		orderInfo.ReMark = mark
		list = append(list, orderInfo)
	}
	resp.List = list
	return resp, datastruct.NULLError
}

const inner_str = "from send_goods s inner join order_info o on s.order_id = o.id inner join user_info u on o.user_id = u.id inner join goods g on g.id = o.goods_id"

func (handle *DBHandler) getAllSendGoodsValue() (int, int64, int, int64) {
	count_sql := "select count(*) " + inner_str
	total := handle.getOrderCount(count_sql)
	rush_sum_sql := "select sum(g.rush_price) " + inner_str + " where o.is_purchase = 0"
	price_sum_sql := "select sum(g.price) " + inner_str + " where o.is_purchase = 1"

	amount := handle.getRushOrderAmount(rush_sum_sql)
	amount += handle.getPurchaseOrderAmount(price_sum_sql)

	today_unix, tomorrow_unix := tools.GetTodayTomorrowTime()

	today_count_sql := fmt.Sprintf(count_sql+" where s.created_at >= %d and s.created_at < %d", today_unix, tomorrow_unix)
	today_total := handle.getOrderCount(today_count_sql)

	today_rush_sum_sql := fmt.Sprintf(rush_sum_sql+" and s.created_at >= %d and s.created_at < %d", today_unix, tomorrow_unix)
	today_price_sum_sql := fmt.Sprintf(price_sum_sql+" and s.created_at >= %d and s.created_at < %d", today_unix, tomorrow_unix)

	today_amount := handle.getPurchaseOrderAmount(today_rush_sum_sql)
	today_amount += handle.getPurchaseOrderAmount(today_price_sum_sql)
	return total, amount, today_total, today_amount
}

func (handle *DBHandler) GetSendGoodsOrder(body *datastruct.GetSendGoodsBody) (interface{}, datastruct.CodeType) {
	resp := new(datastruct.WebResponseSendGoodsInfo)
	resp.Total, resp.TotalAmount, resp.TodayCount, resp.TodayAmount = handle.getAllSendGoodsValue()

	query := " where 1=1"
	if body.OrderId != "" {
		query += " and o.number like " + "'%" + body.OrderId + "%'"
	}
	if body.UserName != "" {
		query += " and u.nick_name like " + "'%" + body.UserName + "%'"
	}
	if body.GoodsName != "" {
		query += " and g.name like " + "'%" + body.GoodsName + "%'"
	}
	if body.EndTime > 0 && body.StartTime > 0 && body.EndTime > body.StartTime {
		query += " and s.created_at >= " + tools.Int64ToString(body.StartTime) + " and s.created_at < " + tools.Int64ToString(body.EndTime)
	}

	var current_count_sql string
	var current_rush_sql string
	var current_price_sql string

	switch body.State {
	case datastruct.Sended:
		query = query + " and s.send_goods_state = 1"
		current_count_sql = "select count(*) " + inner_str + query
		current_rush_sql = "select sum(g.rush_price) " + inner_str + query + " and o.is_purchase = 0"
		current_price_sql = "select sum(g.price) " + inner_str + query + " and o.is_purchase = 1"
	case datastruct.NotSend:
		query = query + " and s.send_goods_state = 0"
		current_count_sql = "select count(*) " + inner_str + query
		current_rush_sql = "select sum(g.rush_price) " + inner_str + query + " and o.is_purchase = 0"
		current_price_sql = "select sum(g.price) " + inner_str + query + " and o.is_purchase = 1"
	default:
		current_count_sql = "select count(*) " + inner_str + query
		current_rush_sql = "select sum(g.rush_price) " + inner_str + query + " and o.is_purchase = 0"
		current_price_sql = "select sum(g.price) " + inner_str + query + " and o.is_purchase = 1"
	}

	current_total := handle.getOrderCount(current_count_sql)
	current_amount := handle.getRushOrderAmount(current_rush_sql)
	current_amount += handle.getPurchaseOrderAmount(current_price_sql)
	resp.CurrentTotal = current_total
	resp.CurrentTotalAmount = current_amount

	engine := handle.mysqlEngine
	start := (body.PageIndex - 1) * body.PageSize
	limit := body.PageSize
	limitStr := fmt.Sprintf(" LIMIT %d,%d", start, limit)

	sql := "select s.sign_for_state,s.express_agency,s.express_number,s.address,s.phone_number,s.link_man,o.is_purchase,o.number,u.nick_name,u.avatar,g.img_name,g.name,g.price,g.rush_price,s.send_goods_state,s.created_at " + inner_str
	orderby := " ORDER BY s.created_at desc"
	sql = sql + query + orderby + limitStr
	results, _ := engine.Query(sql)

	list := make([]*datastruct.SendGoodsInfo, 0, len(results))

	for _, v := range results {
		sendGoodsInfo := new(datastruct.SendGoodsInfo)
		sendGoodsInfo.Avatar = string(v["avatar"][:])
		imgName := string(v["img_name"][:])
		sendGoodsInfo.Count = 1
		sendGoodsInfo.GoodsImgUrl = tools.CreateGoodsImgUrl(imgName)
		sendGoodsInfo.GoodsName = string(v["name"][:])
		sendGoodsInfo.OrderId = string(v["number"][:])
		sendGoodsInfo.SignForState = datastruct.SignForType(tools.StringToInt(string(v["sign_for_state"][:])))
		is_purchase := tools.StringToInt(string(v["is_purchase"][:]))
		if is_purchase == 1 {
			sendGoodsInfo.ReMark = "直接购买"
			sendGoodsInfo.Price = tools.StringToInt64(string(v["price"][:]))
		} else {
			sendGoodsInfo.ReMark = "闯关成功"
			sendGoodsInfo.Price = tools.StringToInt64(string(v["rush_price"][:]))
		}
		sendGoodsInfo.State = datastruct.SendGoodsType(tools.StringToInt(string(v["send_goods_state"][:])))

		receiver := new(datastruct.ReceiverForSendGoods)
		receiver.Address = string(v["address"][:])
		receiver.LinkMan = string(v["link_man"][:])
		receiver.PhoneNumber = string(v["phone_number"][:])
		sendGoodsInfo.Receiver = receiver

		if sendGoodsInfo.State == datastruct.Sended {
			sender := new(datastruct.SenderForSendGoods)
			sender.ExpressAgency = string(v["express_agency"][:])
			sender.ExpressNumber = string(v["express_number"][:])
			sendGoodsInfo.Sender = sender
		}

		sendGoodsInfo.Time = tools.StringToInt64(string(v["created_at"][:]))
		sendGoodsInfo.UserName = string(v["nick_name"][:])
		list = append(list, sendGoodsInfo)
	}
	resp.List = list
	return resp, datastruct.NULLError
}

func (handle *DBHandler) UpdateSendInfo(body *datastruct.UpdateSendInfoBody) datastruct.CodeType {
	engine := handle.mysqlEngine
	updateSend := make([]*updateSendData, 0)
	err := engine.Table("send_goods").Join("INNER", "order_info", "order_info.id = send_goods.order_id").Where("order_info.number=?", body.OrderNumber).Find(&updateSend)
	if err != nil || len(updateSend) <= 0 || updateSend[0].SendGoodsState == datastruct.Sended {
		return datastruct.UpdateDataFailed
	}

	sendgoods := new(datastruct.SendGoods)
	sendgoods.SendGoodsState = datastruct.Sended
	sendgoods.ExpressNumber = body.ExpressNumber
	sendgoods.ExpressAgency = body.ExpressAgency
	var affected int64
	affected, err = engine.Where("order_id=?", updateSend[0].OrderId).Cols("send_goods_state", "express_number", "express_agency").Update(sendgoods)
	if affected <= 0 || err != nil {
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

func (handle *DBHandler) UpdateDefaultAgency(body *datastruct.DefaultAgencyBody) datastruct.CodeType {
	engine := handle.mysqlEngine
	sql := fmt.Sprintf("REPLACE INTO agency_params (identifier,agency1_gold_percent,agency2_gold_percent,agency3_gold_percent,agency1_money_percent,agency2_money_percent,agency3_money_percent)VALUES('%s',%d,%d,%d,%d,%d,%d)", datastruct.AgencyIdentifier, body.Agent1Gold, body.Agent2Gold, body.Agent3Gold, body.Agent1Money, body.Agent2Money, body.Agent3Money)
	_, err := engine.Exec(sql)
	if err != nil {
		log.Debug("DBHandler->UpdateDefaultAgency err:%v", err.Error())
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

func (handle *DBHandler) GetDefaultAgency() (interface{}, datastruct.CodeType) {
	params := new(datastruct.AgencyParams)
	engine := handle.mysqlEngine
	has, err := engine.Where("identifier=?", datastruct.AgencyIdentifier).Get(params)
	if err != nil {
		log.Debug("DBHandler->GetDefaultAgency err:%s", err)
		return nil, datastruct.GetDataFailed
	}
	resp := new(datastruct.DefaultAgencyBody)
	if has {
		resp.Agent1Gold = params.Agency1GoldPercent
		resp.Agent2Gold = params.Agency2GoldPercent
		resp.Agent3Gold = params.Agency3GoldPercent
		resp.Agent1Money = params.Agency1MoneyPercent
		resp.Agent2Money = params.Agency2MoneyPercent
		resp.Agent3Money = params.Agency3MoneyPercent
	}
	return resp, datastruct.NULLError
}

func (handle *DBHandler) EditMemberLevel(body *datastruct.EditLevelDataBody) datastruct.CodeType {
	engine := handle.mysqlEngine
	session := engine.NewSession()
	defer session.Close()
	session.Begin()

	member := new(datastruct.MemberLevelData)
	member.Name = body.Name
	member.IsHidden = body.IsHidden
	member.Price = body.Price
	member.Level = body.Level
	var identifier_id int
	if body.Id <= 0 {
		affected, err := session.Insert(member)
		if err != nil || affected <= 0 {
			rollback("DBHandler->EditMemberLevel insert member err", session)
			return datastruct.UpdateDataFailed
		}
		identifier_id = member.Id
	} else {
		_, err := session.Where("id=?", body.Id).Cols("name", "level", "price", "is_hidden").Update(member)
		if err != nil {
			rollback("DBHandler->EditMemberLevel Update err", session)
			return datastruct.UpdateDataFailed
		}
		identifier_id = body.Id
	}

	sql := fmt.Sprintf("REPLACE INTO agency_params (identifier,agency1_gold_percent,agency2_gold_percent,agency3_gold_percent,agency1_money_percent,agency2_money_percent,agency3_money_percent)VALUES('%s',%d,%d,%d,%d,%d,%d)", tools.IntToString(identifier_id), body.Agent1Gold, body.Agent2Gold, body.Agent3Gold, body.Agent1Money, body.Agent2Money, body.Agent3Money)
	_, err := session.Exec(sql)
	if err != nil {
		str := fmt.Sprintf("DBHandler->EditMemberLevel REPLACE INTO agency_params :%s", err.Error())
		rollback(str, session)
		return datastruct.UpdateDataFailed
	}

	err = session.Commit()
	if err != nil {
		str := fmt.Sprintf("DBHandler->EditMemberLevel Commit :%s", err.Error())
		rollback(str, session)
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

func (handle *DBHandler) GetMemberLevel(name string) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	list := make([]*datastruct.MemberLevelData, 0)
	var err error
	if name == "" {
		err = engine.Desc("level").Find(&list)
	} else {
		err = engine.Where("name like ?", "%"+name+"%").Desc("level").Find(&list)
	}

	if err != nil {
		log.Debug("--------GetMemberLevel err:%s", err)
		return nil, datastruct.GetDataFailed
	}

	resp := make([]*datastruct.EditLevelDataBody, 0)
	var has bool
	for _, v := range list {
		params := new(datastruct.AgencyParams)
		has, err = engine.Where("identifier = ?", tools.IntToString(v.Id)).Get(params)
		if err != nil {
			log.Debug("--------GetMemberLevel get AgencyParams err:%s", err)
			return nil, datastruct.GetDataFailed
		}
		if !has {
			log.Debug("--------GetMemberLevel get AgencyParams err: not has %s", params.Identifier)
			return nil, datastruct.GetDataFailed
		}
		levelData := new(datastruct.EditLevelDataBody)
		levelData.Agent1Gold = params.Agency1GoldPercent
		levelData.Agent2Gold = params.Agency2GoldPercent
		levelData.Agent3Gold = params.Agency3GoldPercent
		levelData.Agent1Money = params.Agency1MoneyPercent
		levelData.Agent2Money = params.Agency2MoneyPercent
		levelData.Agent3Money = params.Agency3MoneyPercent
		levelData.Id = v.Id
		levelData.IsHidden = v.IsHidden
		levelData.Level = v.Level
		levelData.Name = v.Name
		levelData.Price = v.Price
		resp = append(resp, levelData)
	}

	return resp, datastruct.NULLError
}

func (handle *DBHandler) GetMembers(body *datastruct.WebGetMembersBody) (interface{}, datastruct.CodeType) {

	query := ""
	args := make([]interface{}, 0)
	if body.IsBlacklist >= 2 {
		query = "user_info.is_black_list >= ?"
		args = append(args, 0)
	} else {
		query = "user_info.is_black_list = ?"
		args = append(args, body.IsBlacklist)
	}

	if body.NickName != "" {
		query += " and user_info.nick_name like ?"
		args = append(args, "%"+body.NickName+"%")
	}
	users := make([]*datastruct.UserInfo, 0)
	engine := handle.mysqlEngine
	start := (body.PageIndex - 1) * body.PageSize
	limit := body.PageSize
	resp := new(datastruct.WebResponseMembersInfo)
	var levelName string
	var userLevelId int
	if body.LevelId < 0 {
		total, err := engine.Where(query, args...).Desc("created_at").Count(new(datastruct.UserInfo))
		if err != nil {
			log.Debug("GetMembers1-------------------Err:%v", err.Error())
			return nil, datastruct.GetDataFailed
		}
		resp.CurrentTotal = int(total)
		engine.Where(query, args...).Desc("created_at").Limit(limit, start).Find(&users)

	} else if body.LevelId >= 0 {
		if body.LevelId == 0 {
			levelName = datastruct.DefaultMemberIdentifier
			userLevelId = 0
			levelDataArr := make([]*datastruct.MemberLevelData, 0)
			err := engine.Find(&levelDataArr)
			if err != nil {
				log.Debug("GetMembers2-------------------Err:%v", err.Error())
				return nil, datastruct.GetDataFailed
			}
			str_arr := make([]interface{}, 0, len(levelDataArr))
			for _, v := range levelDataArr {
				str_arr = append(str_arr, fmt.Sprintf("%d", v.Id))
			}
			total, err := engine.Where(query, args...).NotIn("member_identifier", str_arr...).Desc("created_at").Count(new(datastruct.UserInfo))
			if err != nil {
				log.Debug("GetMembers9-------------------Err:%v", err.Error())
				return nil, datastruct.GetDataFailed
			}
			resp.CurrentTotal = int(total)
			engine.Where(query, args...).NotIn("member_identifier", str_arr...).Desc("created_at").Limit(limit, start).Find(&users)
		} else {
			levelData := new(datastruct.MemberLevelData)
			has, err := engine.Where("id = ?", body.LevelId).Get(levelData)
			if err != nil || !has {
				log.Debug("GetMembers3-------------------Err:%v", err.Error())
				return nil, datastruct.GetDataFailed
			}
			levelName = levelData.Name
			userLevelId = levelData.Id
			query += " and user_info.member_identifier = ?"
			args = append(args, tools.IntToString(body.LevelId))
			total, err := engine.Where(query, args...).Desc("created_at").Count(new(datastruct.UserInfo))
			if err != nil {
				log.Debug("GetMembers4-------------------Err:%v", err.Error())
				return nil, datastruct.GetDataFailed
			}
			resp.CurrentTotal = int(total)
			engine.Where(query, args...).Desc("created_at").Limit(limit, start).Find(&users)
		}
	}
	members := make([]*datastruct.WebResponseMember, 0, len(users))
	for _, v := range users {
		member := new(datastruct.WebResponseMember)
		member.Avatar = v.Avatar
		member.Balance = v.Balance
		member.BalanceTotal = v.BalanceTotal
		member.CreateTime = v.CreatedAt
		member.GoldCount = v.GoldCount
		member.Id = v.Id
		member.IsBlacklist = v.IsBlackList
		member.NickName = v.NickName
		if body.LevelId < 0 {
			if v.MemberIdentifier == datastruct.AgencyIdentifier {
				levelName = datastruct.DefaultMemberIdentifier
				userLevelId = 0
			} else {
				levelid := tools.StringToInt(v.MemberIdentifier)
				levelData := new(datastruct.MemberLevelData)
				has, err := engine.Where("id = ?", levelid).Get(levelData)
				if err != nil {
					log.Debug("GetMembers5-------------------Err:%v", err.Error())
					return nil, datastruct.GetDataFailed
				}
				if has {
					levelName = levelData.Name
					userLevelId = levelData.Id
				} else {
					levelName = datastruct.DefaultMemberIdentifier
					userLevelId = 0
				}
			}
		}
		member.LevelName = levelName
		member.LevelId = userLevelId
		members = append(members, member)
	}
	resp.Members = members
	user := new(datastruct.UserInfo)
	total, err := engine.Count(user)
	if err != nil {
		log.Debug("GetMembers6-------------------Err:%v", err.Error())
		return nil, datastruct.GetDataFailed
	}
	resp.Total = int(total)
	today_unix, tomorrow_unix := tools.GetTodayTomorrowTime()
	total, err = engine.Where("created_at >= ? and created_at < ?", today_unix, tomorrow_unix).Count(user)
	if err != nil {
		log.Debug("GetMembers7-------------------Err:%v", err.Error())
		return nil, datastruct.GetDataFailed
	}
	resp.TodayCreated = int(total)
	return resp, datastruct.NULLError
}

func (handle *DBHandler) UpdateUserBlackList(userid int, state int) (string, datastruct.CodeType) {
	engine := handle.mysqlEngine
	user := new(datastruct.UserInfo)
	has, err := engine.Where("id = ?", userid).Get(user)
	if err != nil || !has {
		return "", datastruct.GetDataFailed
	}
	token := user.Token
	user.IsBlackList = state
	_, err = engine.Where("id = ?", userid).Cols("is_black_list").Update(user)
	if err != nil {
		return "", datastruct.UpdateDataFailed
	}
	return token, datastruct.NULLError
}

func (handle *DBHandler) UpdateUserLevel(userid int, levelid int) datastruct.CodeType {
	engine := handle.mysqlEngine
	user := new(datastruct.UserInfo)
	if levelid == 0 {
		user.MemberIdentifier = datastruct.AgencyIdentifier
	} else {
		user.MemberIdentifier = tools.IntToString(levelid)
	}
	_, err := engine.Where("id = ?", userid).Cols("member_identifier").Update(user)
	if err != nil {
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

func (handle *DBHandler) WebChangeGold(userid int, gold int64) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	session := engine.NewSession()
	defer session.Close()
	session.Begin()
	if gold != 0 {
		sql := "update user_info set gold_count = case when gold_count + ? > 0 then gold_count + ? else 0 end where id = ?"
		_, err := session.Exec(sql, gold, gold, userid)
		if err != nil {
			str := fmt.Sprintf("DBHandler->WebChangeGold Update user_info :%s", err.Error())
			rollback(str, session)
			return -1, datastruct.UpdateDataFailed
		}

		goldChange := new(datastruct.GoldChangeInfo)
		goldChange.CreatedAt = time.Now().Unix()
		goldChange.UserId = userid
		goldChange.VarGold = int64(math.Abs(float64(gold)))
		if gold > 0 {
			goldChange.ChangeType = datastruct.GrantType
		} else {
			goldChange.ChangeType = datastruct.DeductType
		}
		var affected int64
		affected, err = session.Insert(goldChange)
		if err != nil || affected <= 0 {
			str := fmt.Sprintf("DBHandler->WebChangeGold Update user_info :%s", err.Error())
			rollback(str, session)
			return -1, datastruct.UpdateDataFailed
		}
	}
	user := new(datastruct.UserInfo)
	has, err := session.Where("id=?", userid).Get(user)
	if err != nil || !has {
		str := fmt.Sprintf("DBHandler->WebChangeGold Get user_info err")
		rollback(str, session)
		return -1, datastruct.GetDataFailed
	}

	err = session.Commit()
	if err != nil {
		str := fmt.Sprintf("DBHandler->EditDomain Commit :%s", err.Error())
		rollback(str, session)
		return -1, datastruct.UpdateDataFailed
	}
	return user.GoldCount, datastruct.NULLError
}

func (handle *DBHandler) GetServerInfo() (*datastruct.WebServerInfoBody, datastruct.CodeType) {
	engine := handle.mysqlEngine
	serverInfo := new(datastruct.ServerVersion)
	has, err := engine.Desc("created_at").Limit(1, 0).Get(serverInfo)
	if err != nil || !has {
		return nil, datastruct.GetDataFailed
	}
	resp := new(datastruct.WebServerInfoBody)
	resp.IsMaintain = serverInfo.IsMaintain
	resp.Version = serverInfo.Version
	return resp, datastruct.NULLError
}

func (handle *DBHandler) EditServerInfo(version string, isMaintain int) datastruct.CodeType {
	engine := handle.mysqlEngine
	now_time := time.Now().Unix()
	sql := fmt.Sprintf("REPLACE INTO server_version (version,is_maintain,created_at)VALUES('%s',%d,%d)", version, isMaintain, now_time)
	_, err := engine.Exec(sql)
	if err != nil {
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

func (handle *DBHandler) getAgencyStatisticsWithUsers(users []int) (*datastruct.WebAgencyStatistics, []int) {
	if users == nil || len(users) <= 0 {
		statistics := new(datastruct.WebAgencyStatistics)
		statistics.Count = 0
		statistics.DepositTotal = 0
		statistics.PayRushTotal = 0
		statistics.PurchaseTotal = 0
		return statistics, nil
	}
	ids_str := getUsersStr(users)
	inner_str := "from invite_info i inner join user_info u on i.receiver = u.id where i.sender in (" + ids_str + ")"
	count_sql := "select u.id,u.deposit_total,u.pay_rush_total,u.purchase_total " + inner_str
	engine := handle.mysqlEngine
	results, _ := engine.Query(count_sql)

	var depositTotal int64
	depositTotal = 0

	var payRushTotal int64
	payRushTotal = 0

	var purchaseTotal float64
	purchaseTotal = 0

	count := len(results)
	ids := make([]int, 0, count)
	for _, v := range results {
		depositTotal += tools.StringToInt64(string(v["deposit_total"][:]))
		payRushTotal += tools.StringToInt64(string(v["pay_rush_total"][:]))
		purchaseTotal += tools.StringToFloat64(string(v["purchase_total"][:]))
		ids = append(ids, tools.StringToInt(string(v["id"][:])))
	}

	statistics := new(datastruct.WebAgencyStatistics)
	statistics.Count = count
	statistics.DepositTotal = depositTotal
	statistics.PayRushTotal = payRushTotal
	statistics.PurchaseTotal = purchaseTotal
	return statistics, ids
}

func getUsersStr(users []int) string {
	ids_str := ""
	num := len(users)
	for k, v := range users {
		if k == num-1 {
			ids_str += tools.IntToString(v)
		} else {
			ids_str += tools.IntToString(v) + ","
		}
	}
	return ids_str
}

func (handle *DBHandler) UpdateGoodsClassState(classid int, isHidden int) datastruct.CodeType {
	engine := handle.mysqlEngine
	gclass := new(datastruct.GoodsClass)
	gclass.IsHidden = isHidden
	_, err := engine.Where("id = ?", classid).Cols("is_hidden").Update(gclass)
	if err != nil {
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

func (handle *DBHandler) EditGoodsClass(body *datastruct.WebEditGoodsClassBody) datastruct.CodeType {
	engine := handle.mysqlEngine
	gclass := new(datastruct.GoodsClass)
	gclass.IsHidden = body.IsHidden
	gclass.Name = body.Name
	gclass.SortId = body.SortId
	session := engine.NewSession()
	defer session.Close()
	session.Begin()

	isUpdateIcon := false
	if !strings.Contains(body.Icon, conf.Server.Domain) && body.Icon != "" {
		isUpdateIcon = true
		if body.Id > 0 {
			tmp := new(datastruct.GoodsClass)
			has, err := session.Where("id=?", body.Id).Get(tmp)
			if err != nil || !has {
				rollback("DBHandler->EditGoodsClass Get GoodsClass err", session)
				return datastruct.UpdateDataFailed
			}
			deleteFile(tools.GetImgPath() + tmp.Icon)
		}
		imgName := fmt.Sprintf("%s.png", tools.UniqueId())
		path := tools.GetImgPath() + imgName
		arr_str := strings.Split(body.Icon, ",")
		var isError bool
		if len(arr_str) > 1 {
			isError = tools.CreateImgFromBase64(&arr_str[1], path)
		} else {
			isError = tools.CreateImgFromBase64(&arr_str[0], path)
		}
		if isError {
			rollback("DBHandler->EditGoodsClass CreateImgFromBase64 err", session)
			return datastruct.UpdateDataFailed
		}
		gclass.Icon = imgName
	}

	var err error
	if body.Id == 0 {
		_, err = session.Insert(gclass)
	} else {
		if isUpdateIcon {
			_, err = session.Where("id = ?", body.Id).Cols("name,is_hidden,sort_id,icon").Update(gclass)
		} else {
			_, err = session.Where("id = ?", body.Id).Cols("name,is_hidden,sort_id").Update(gclass)
		}
	}

	if err != nil {
		deleteFile(tools.GetImgPath() + gclass.Icon)
		rollback("DBHandler->EditGoodsClass err:%s"+err.Error(), session)
		return datastruct.UpdateDataFailed
	}

	err = session.Commit()
	if err != nil {
		deleteFile(tools.GetImgPath() + gclass.Icon)
		str := fmt.Sprintf("DBHandler->EditGoodsClass Commit :%s", err.Error())
		rollback(str, session)
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

func deleteFile(path string) {
	tools.DeleteFile(path)
}

func (handle *DBHandler) GetAllGoodsClasses(body *datastruct.WebQueryGoodsClassBody) (interface{}, datastruct.CodeType) {
	gclasses := make([]*datastruct.GoodsClass, 0)
	engine := handle.mysqlEngine
	query := ""
	args := make([]interface{}, 0)
	if body.IsHidden == 2 {
		query = "is_hidden >= ?"
		args = append(args, 0)
	} else {
		query = "is_hidden = ?"
		args = append(args, body.IsHidden)
	}
	if body.Name != "" {
		query += " and name like ?"
		args = append(args, "%"+body.Name+"%")
	}
	err := engine.Where(query, args...).Desc("sort_id").Find(&gclasses)
	if err != nil {
		return nil, datastruct.GetDataFailed
	}
	resp := make([]*datastruct.WebResponseGoodsClass, 0)
	for _, v := range gclasses {
		gclass := new(datastruct.WebResponseGoodsClass)
		gclass.Icon = tools.CreateGoodsImgUrl(v.Icon)
		gclass.Id = v.Id
		gclass.IsHidden = v.IsHidden
		gclass.Name = v.Name
		gclass.SortId = v.SortId
		resp = append(resp, gclass)
	}

	return resp, datastruct.NULLError
}

type webDepositInfo struct {
	datastruct.UserDepositInfo `xorm:"extends"`
	datastruct.UserInfo        `xorm:"extends"`
	datastruct.BalanceInfo     `xorm:"extends"`
}

func (handle *DBHandler) GetAllDepositInfo(body *datastruct.WebQueryDepositInfoBody) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	resp := new(datastruct.WebResponseDepositInfo)
	userDepositInfo := new(datastruct.UserDepositInfo)
	resp.Count, _ = engine.Count(userDepositInfo)
	resp.Amount, _ = engine.Sum(userDepositInfo, "money")

	today_unix, tomorrow_unix := tools.GetTodayTomorrowTime()
	resp.TodayCount, _ = engine.Where("created_at >= ? and created_at < ?", today_unix, tomorrow_unix).Count(userDepositInfo)
	resp.TodayAmount, _ = engine.Where("created_at >= ? and created_at < ?", today_unix, tomorrow_unix).Sum(userDepositInfo, "money")

	balanceInfo := new(datastruct.BalanceInfo)
	resp.EarnMoney, _ = engine.Sum(balanceInfo, "earn_balance")
	resp.EarnGold, _ = engine.Sum(balanceInfo, "earn_gold")

	args := make([]interface{}, 0)
	query := "1=1"
	if body.Platform != 2 {
		query += " and user_deposit_info.platform = ?"
		args = append(args, body.Platform)
	}
	if body.Name != "" {
		query += " and user_info.nick_name like ?"
		args = append(args, "%"+body.Name+"%")
	}
	if body.EndTime > 0 && body.StartTime > 0 && body.EndTime > body.StartTime {
		query += " and user_deposit_info.created_at >= " + tools.Int64ToString(body.StartTime) + " and user_deposit_info.created_at < " + tools.Int64ToString(body.EndTime)
	}

	start := (body.PageIndex - 1) * body.PageSize
	limit := body.PageSize
	depositInfo := make([]*webDepositInfo, 0)
	resp.CurrentTotal, _ = engine.Table("user_deposit_info").Join("INNER", "user_info", "user_info.id = user_deposit_info.user_id").Where(query, args...).Desc("user_deposit_info.id").Count(balanceInfo)

	engine.Table("user_deposit_info").Join("INNER", "user_info", "user_info.id = user_deposit_info.user_id").Where(query, args...).Desc("user_deposit_info.id").Limit(limit, start).Find(&depositInfo)

	depositUsers := make([]*datastruct.WebDepositUser, 0)
	for _, v := range depositInfo {
		depositUser := new(datastruct.WebDepositUser)
		depositUser.CreateTime = v.UserDepositInfo.CreatedAt
		depositUser.Pay = int64(v.UserDepositInfo.Money)
		depositUser.Platform = int(v.UserDepositInfo.Platform)
		depositUser.NickName = v.UserInfo.NickName
		depositUser.Avatar = v.UserInfo.Avatar
		depositUser.Agency = queryAgencyEarn(v.UserDepositInfo.UserId, v.UserDepositInfo.Id, engine)
		depositUsers = append(depositUsers, depositUser)
	}
	resp.Users = depositUsers
	return resp, datastruct.NULLError
}

func queryAgencyEarn(fromUserId int, depositId int, engine *xorm.Engine) []*datastruct.WebDepositAgency {
	agencys := make([]*datastruct.WebDepositAgency, 0, datastruct.MaxLevel)
	for i := 0; i < datastruct.MaxLevel; i++ {
		level := i + 1
		balanceInfo := new(datastruct.BalanceInfo)
		has, err := engine.Where("agency_level = ? and from_user_id = ? and deposit_id=?", level, fromUserId, depositId).Get(balanceInfo)
		if err != nil {
			log.Debug("queryAgencyEarn err:%s", err.Error())
			continue
		}
		agency := new(datastruct.WebDepositAgency)
		if has {
			agency.Gold = balanceInfo.EarnGold
			agency.Money = balanceInfo.EarnBalance
		} else {
			agency.Gold = 0
			agency.Money = 0
		}
		agencys = append(agencys, agency)
	}
	return agencys
}

type webDrawCashInfo struct {
	datastruct.DrawCashInfo `xorm:"extends"`
	datastruct.UserInfo     `xorm:"extends"`
}

func (handle *DBHandler) GetAllDrawInfo(body *datastruct.WebQueryDrawInfoBody) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	resp := new(datastruct.WebResponseDrawInfo)

	drawCashInfo := new(datastruct.DrawCashInfo)
	resp.Amount, _ = engine.Where("state = 1").Sum(drawCashInfo, "charge")
	resp.Poundage, _ = engine.Where("state = 1").Sum(drawCashInfo, "poundage")

	today_unix, tomorrow_unix := tools.GetTodayTomorrowTime()
	resp.TodayAmount, _ = engine.Where("created_at >= ? and created_at < ? and state = 1", today_unix, tomorrow_unix).Sum(drawCashInfo, "charge")

	args := make([]interface{}, 0)
	query := "1=1"
	if body.Name != "" {
		query += " and user_info.nick_name like ?"
		args = append(args, "%"+body.Name+"%")
	}
	if body.EndTime > 0 && body.StartTime > 0 && body.EndTime > body.StartTime {
		query += " and draw_cash_info.created_at >= " + tools.Int64ToString(body.StartTime) + " and draw_cash_info.created_at < " + tools.Int64ToString(body.EndTime)
	}
	if body.State != 3 {
		query += " and draw_cash_info.state = ?"
		args = append(args, body.State)
	}
	start := (body.PageIndex - 1) * body.PageSize
	limit := body.PageSize
	list := make([]*webDrawCashInfo, 0)

	resp.CurrentTotal, _ = engine.Table("draw_cash_info").Join("INNER", "user_info", "user_info.id = draw_cash_info.user_id").Where(query, args...).Desc("draw_cash_info.created_at").Count(drawCashInfo)
	engine.Table("draw_cash_info").Join("INNER", "user_info", "user_info.id = draw_cash_info.user_id").Where(query, args...).Desc("draw_cash_info.created_at").Limit(limit, start).Find(&list)

	users := make([]*datastruct.WebDrawUser, 0)
	for _, v := range list {
		user := new(datastruct.WebDrawUser)
		user.Id = v.DrawCashInfo.Id
		user.Avatar = v.UserInfo.Avatar
		user.NickName = v.UserInfo.NickName
		user.Charge = v.DrawCashInfo.Charge
		user.PaymentNo = v.DrawCashInfo.PaymentNo
		if v.DrawCashInfo.PaymentTime == "" {
			user.CreateTime = tools.UnixToString(v.DrawCashInfo.CreatedAt, "2006-01-02 15:04:05")
		} else {
			user.CreateTime = v.DrawCashInfo.PaymentTime
		}
		switch v.ArrivalType {
		case datastruct.DrawCashArrivalWX:
			user.ArrivalType = "微信钱包"
		case datastruct.DrawCashArrivalZFB:
			user.ArrivalType = "支付宝"
		}
		user.Poundage = v.DrawCashInfo.Poundage
		user.State = v.DrawCashInfo.State
		users = append(users, user)
	}
	resp.Users = users
	return resp, datastruct.NULLError
}

func (handle *DBHandler) GetAllMembers() (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	list := make([]*datastruct.MemberLevelData, 0)
	engine.Asc("level").Find(&list)
	resp := make([]*datastruct.WebResponseNotHiddenMember, 0, len(list))
	for _, v := range list {
		m_data := new(datastruct.WebResponseNotHiddenMember)
		m_data.Id = v.Id
		m_data.Name = v.Name
		resp = append(resp, m_data)
	}
	return resp, datastruct.NULLError
}

func (handle *DBHandler) UpdateMemberLevelState(body *datastruct.WebUpdateMemberLevelBody) datastruct.CodeType {
	engine := handle.mysqlEngine
	m_data := new(datastruct.MemberLevelData)
	m_data.IsHidden = body.IsHidden
	_, err := engine.Where("id=?", body.Id).Cols("is_hidden").Update(m_data)
	if err != nil {
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

func (handle *DBHandler) GetRushLimitSetting() (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	setting := new(datastruct.RushLimitSetting)
	has, err := engine.Where("id=?", datastruct.DefaultId).Get(setting)
	if err != nil {
		return nil, datastruct.GetDataFailed
	}
	resp := new(datastruct.WebRushLimitSettingBody)
	if has {
		resp.CheatCount = setting.CheatCount
		resp.CheatTips = setting.CheatTips
		resp.Diff2 = setting.Diff2
		resp.Diff3 = setting.Diff3
		resp.Diff2r = setting.Diff2r
		resp.Diff2t = setting.Diff2t
		resp.Diff3r = setting.Diff3r
		resp.Diff3t = setting.Diff3t
		resp.LotteryCount = setting.LotteryCount
	}
	return resp, datastruct.NULLError
}

func (handle *DBHandler) EditRushLimitSetting(body *datastruct.WebRushLimitSettingBody) datastruct.CodeType {
	engine := handle.mysqlEngine
	setting := new(datastruct.RushLimitSetting)
	has, err := engine.Where("id=?", datastruct.DefaultId).Get(setting)
	if err != nil {
		return datastruct.GetDataFailed
	}
	setting.CheatCount = body.CheatCount
	setting.CheatTips = body.CheatTips
	setting.Diff2 = body.Diff2
	setting.Diff3 = body.Diff3
	setting.Diff2r = body.Diff2r
	setting.Diff2t = body.Diff2t
	setting.Diff3r = body.Diff3r
	setting.Diff3t = body.Diff3t
	setting.LotteryCount = body.LotteryCount
	if has {
		_, err = engine.Where("id=?", datastruct.DefaultId).Update(setting)
		if err != nil {
			return datastruct.UpdateDataFailed
		}
	} else {
		setting.Id = datastruct.DefaultId
		_, err = engine.Insert(setting)
		if err != nil {
			return datastruct.UpdateDataFailed
		}
	}
	return datastruct.NULLError
}

/*
func (handle *DBHandler) GetNewUsers(body *datastruct.WebNewsUserBody) (interface{}, datastruct.CodeType) {
	user := new(datastruct.UserInfo)
	engine := handle.mysqlEngine
	var start, end int64
	if body.StartTime == 0 || body.EndTime == 0 {
		today_unix, tomorrow_unix := tools.GetTodayTomorrowTime()
		start = today_unix
		end = tomorrow_unix
	} else {
		start = body.StartTime
		end = body.EndTime
	}
	count, err := engine.Where("created_at >= ? and created_at < ?", start, end).Count(user)
	if err != nil {
		return nil, datastruct.GetDataFailed
	}
	return count, datastruct.NULLError
}

func (handle *DBHandler) GetTotalEarn(body *datastruct.WebNewsUserBody) (interface{}, datastruct.CodeType) {
	userDepositInfo := new(datastruct.UserDepositInfo)
	engine := handle.mysqlEngine
	var start, end int64
	if body.StartTime == 0 || body.EndTime == 0 {
		today_unix, tomorrow_unix := tools.GetTodayTomorrowTime()
		start = today_unix
		end = tomorrow_unix
	} else {
		start = body.StartTime
		end = body.EndTime
	}
	total, err := engine.Where("created_at >= ? and created_at < ?", start, end).Sum(userDepositInfo, "money")
	if err != nil {
		return nil, datastruct.GetDataFailed
	}
	return total, datastruct.NULLError
}

func (handle *DBHandler) GetDepositUsers(body *datastruct.WebNewsUserBody) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	var start, end int64
	if body.StartTime == 0 || body.EndTime == 0 {
		today_unix, tomorrow_unix := tools.GetTodayTomorrowTime()
		start = today_unix
		end = tomorrow_unix
	} else {
		start = body.StartTime
		end = body.EndTime
	}
	total_sql := "SELECT count(*) FROM user_deposit_info where created_at>= ? and created_at < ? GROUP BY user_id"
	rs, err := engine.Query(total_sql, start, end)
	if err != nil {
		return nil, datastruct.GetDataFailed
	}
	count := len(rs)
	new_users_sql := "SELECT count(*) FROM user_deposit_info udi join user_info u on udi.user_id = u.id where u.created_at>= ? and u.created_at < ? and udi.created_at >= ? and udi.created_at < ? GROUP BY user_id"
	rs, err = engine.Query(new_users_sql, start, end, start, end)
	if err != nil {
		return nil, datastruct.GetDataFailed
	}
	new_users := len(rs)
	old_users := count - new_users
	resp := new(datastruct.WebResponseOldAndNewUsers)
	resp.OldUsers = old_users
	resp.NewUsers = new_users
	return resp, datastruct.NULLError
}

func (handle *DBHandler) GetActivityUsers(body *datastruct.WebNewsUserBody) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	var new_users, old_users int
	var code datastruct.CodeType
	resp := new(datastruct.WebResponseOldAndNewUsers)
	if body.StartTime == 0 || body.EndTime == 0 {
		today_unix, tomorrow_unix := tools.GetTodayTomorrowTime()
		old_users, new_users, code = computeToday(today_unix, tomorrow_unix, engine)
		if code != datastruct.NULLError {
			return nil, code
		}
	} else {

		today_unix, tomorrow_unix := tools.GetTodayTomorrowTime()
		if body.EndTime >= today_unix {
			old_users, new_users, code = computeToday(today_unix, tomorrow_unix, engine)
			if code != datastruct.NULLError {
				return nil, code
			}
		} else {
			old_users = 0
			new_users = 0
		}
		current := body.StartTime
		for {
			fileName := tools.Int64ToString(current)
			tmp_old, tmp_new := tools.GetActivityData(fileName)
			old_users += tmp_old
			new_users += tmp_new
			current += 24 * 3600
			if current > body.EndTime {
				break
			}
		}
	}
	resp.NewUsers = new_users
	resp.OldUsers = old_users
	return resp, datastruct.NULLError
}

func computeToday(start int64, end int64, engine *xorm.Engine) (int, int, datastruct.CodeType) {
	var total int64
	var new_users int64
	var old_users int64
	var err error
	todayActivity := new(datastruct.TodayUserActivityInfo)
	total, err = engine.Where("start_time >= ? and start_time < ? ", start, end).Count(todayActivity)
	if err != nil {
		return -1, -1, datastruct.GetDataFailed
	}
	users := new(datastruct.UserInfo)
	query_str := "user_info.created_at >= ? and user_info.created_at < ? and today_user_activity_info.start_time >= ? and today_user_activity_info.start_time < ?"
	args := make([]interface{}, 0, 4)
	args = append(args, start)
	args = append(args, end)
	args = append(args, start)
	args = append(args, end)
	new_users, err = engine.Table("today_user_activity_info").Join("INNER", "user_info", "user_info.id = today_user_activity_info.user_id").Where(query_str, args...).Count(users)
	if err != nil {
		return -1, -1, datastruct.GetDataFailed
	}
	old_users = total - new_users
	return int(old_users), int(new_users), datastruct.NULLError
}*/

func (handle *DBHandler) MyPrentices(body *datastruct.WebGetAgencyInfoBody) (interface{}, datastruct.CodeType) {
	statistics1, users1 := handle.getAgencyStatisticsWithUsers([]int{body.UserId})
	statistics2, users2 := handle.getAgencyStatisticsWithUsers(users1)
	statistics3, _ := handle.getAgencyStatisticsWithUsers(users2)
	statistics_list := make([]*datastruct.WebAgencyStatistics, 0, 3)
	statistics_list = append(statistics_list, statistics1)
	statistics_list = append(statistics_list, statistics2)
	statistics_list = append(statistics_list, statistics3)
	resp := new(datastruct.WebResponseAgencyData)
	resp.Statistics = statistics_list

	engine := handle.mysqlEngine
	switch body.Level {
	case 1:
		resp.Users, _, resp.CurrentTotal = getAgencyUser(engine, []int{body.UserId}, false, body)
	case 2:
		_, ids, _ := getAgencyUser(engine, []int{body.UserId}, true, nil)
		resp.Users, _, resp.CurrentTotal = getAgencyUser(engine, ids, false, body)
	case 3:
		_, ids1, _ := getAgencyUser(engine, []int{body.UserId}, true, nil)
		_, ids2, _ := getAgencyUser(engine, ids1, true, nil)
		resp.Users, _, resp.CurrentTotal = getAgencyUser(engine, ids2, false, body)
	}
	return resp, datastruct.NULLError
}

func getAgencyUser(engine *xorm.Engine, users []int, isQueryAll bool, body *datastruct.WebGetAgencyInfoBody) ([]*datastruct.WebAgencyUser, []int, int) {
	resp := make([]*datastruct.WebAgencyUser, 0)
	if users == nil || len(users) <= 0 {
		return resp, nil, 0
	}
	ids_str := getUsersStr(users)
	inner_str := "from invite_info i inner join user_info u on i.receiver = u.id where i.sender in (" + ids_str + ")"

	var rs_sql string
	var current_count int
	if !isQueryAll {
		start := (body.PageIndex - 1) * body.PageSize
		limit := body.PageSize
		limitStr := fmt.Sprintf(" ORDER BY i.created_at desc LIMIT %d,%d", start, limit)
		query_str := ""
		if body.StartTime != 0 && body.EndTime != 0 {
			query_str += " and i.created_at >= " + tools.Int64ToString(body.StartTime) + " and i.created_at < " + tools.Int64ToString(body.EndTime)
		}
		if body.Name != "" {
			query_str += " and u.nick_name like " + "'%" + body.Name + "%'"
		}
		count_sql := "select count(*) " + inner_str + query_str
		count_results, _ := engine.Query(count_sql)
		current_count = tools.StringToInt(string(count_results[0]["count(*)"][:]))

		inner_str += query_str + limitStr

		rs_sql = "select u.balance,u.gold_count,i.created_at,u.avatar,u.nick_name,u.deposit_total,u.pay_rush_total,u.purchase_total " + inner_str

	} else {
		rs_sql = "select u.id " + inner_str
	}

	results, _ := engine.Query(rs_sql)

	ids := make([]int, 0, len(results))
	if !isQueryAll {
		for _, v := range results {
			agencyUser := new(datastruct.WebAgencyUser)
			agencyUser.Avatar = string(v["avatar"][:])
			agencyUser.Balance = tools.StringToFloat64(string(v["balance"][:]))
			agencyUser.CreatedAt = tools.StringToInt64(string(v["created_at"][:]))
			agencyUser.DepositTotal = tools.StringToInt64(string(v["deposit_total"][:]))
			agencyUser.GoldCount = tools.StringToInt64(string(v["gold_count"][:]))
			agencyUser.NickName = string(v["nick_name"][:])
			agencyUser.PayRushTotal = tools.StringToInt64(string(v["pay_rush_total"][:]))
			agencyUser.PurchaseTotal = tools.StringToFloat64(string(v["purchase_total"][:]))
			resp = append(resp, agencyUser)
		}
	} else {
		for _, v := range results {
			ids = append(ids, tools.StringToInt(string(v["id"][:])))
		}
	}

	return resp, ids, current_count
}

type webMemberOrderData struct {
	datastruct.MemberLevelOrder `xorm:"extends"`
	datastruct.MemberLevelData  `xorm:"extends"`
	datastruct.UserInfo         `xorm:"extends"`
}

func (handle *DBHandler) GetMemberOrder(body *datastruct.WebMemberOrderBody) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	start := (body.PageIndex - 1) * body.PageSize
	limit := body.PageSize
	query_str := "1=1"
	args := make([]interface{}, 0)
	if body.StartTime != 0 && body.EndTime != 0 {
		query_str += " and member_level_order.created_at >= ? and member_level_order.created_at < ?"
		args = append(args, body.StartTime)
		args = append(args, body.EndTime)
	}
	if body.Name != "" {
		query_str += " and user_info.nick_name like ?"
		args = append(args, "%"+body.Name+"%")
	}
	if body.LevelName != "" {
		query_str += " and member_level_data.name like ?"
		args = append(args, "%"+body.LevelName+"%")
	}
	resp := new(datastruct.WebResponseMemberOrder)
	arr := make([]*webMemberOrderData, 0)
	m_order := new(datastruct.MemberLevelData)
	today_unix, tomorrow_unix := tools.GetTodayTomorrowTime()
	count, _ := engine.Table("member_level_order").Join("INNER", "user_info", "user_info.id = member_level_order.user_id").Join("INNER", "member_level_data", "member_level_data.id = member_level_order.member_level_id").Count(m_order)
	todayCount, _ := engine.Table("member_level_order").Join("INNER", "user_info", "user_info.id = member_level_order.user_id").Join("INNER", "member_level_data", "member_level_data.id = member_level_order.member_level_id").Where("member_level_order.created_at >= ? and member_level_order.created_at < ?", today_unix, tomorrow_unix).Count(m_order)
	current, _ := engine.Table("member_level_order").Join("INNER", "user_info", "user_info.id = member_level_order.user_id").Join("INNER", "member_level_data", "member_level_data.id = member_level_order.member_level_id").Where(query_str, args...).Count(m_order)
	resp.Count = count
	resp.TodayCount = todayCount
	resp.CurrentTotal = current

	amount, _ := engine.Table("member_level_order").Join("INNER", "user_info", "user_info.id = member_level_order.user_id").Join("INNER", "member_level_data", "member_level_data.id = member_level_order.member_level_id").Sum(m_order, "member_level_data.price")

	todayAmount, _ := engine.Table("member_level_order").Join("INNER", "user_info", "user_info.id = member_level_order.user_id").Join("INNER", "member_level_data", "member_level_data.id = member_level_order.member_level_id").Where("member_level_order.created_at >= ? and member_level_order.created_at < ?", today_unix, tomorrow_unix).Sum(m_order, "member_level_data.price")
	resp.Amount = amount
	resp.TodayAmount = todayAmount

	engine.Table("member_level_order").Join("INNER", "user_info", "user_info.id = member_level_order.user_id").Join("INNER", "member_level_data", "member_level_data.id = member_level_order.member_level_id").Where(query_str, args...).Desc("member_level_order.created_at").Limit(limit, start).Find(&arr)

	list := make([]*datastruct.WebMemberOrder, 0)
	for _, v := range arr {
		memberorder := new(datastruct.WebMemberOrder)
		memberorder.Avatar = v.UserInfo.Avatar
		memberorder.CreatedAt = v.MemberLevelOrder.CreatedAt
		memberorder.Id = v.MemberLevelOrder.Id
		memberorder.LevelName = v.MemberLevelData.Name
		memberorder.NickName = v.UserInfo.NickName
		memberorder.Price = v.MemberLevelData.Price
		list = append(list, memberorder)
	}
	resp.List = list
	return resp, datastruct.NULLError
}

func (handle *DBHandler) DeleteMemberOrder(orderId int) datastruct.CodeType {
	engine := handle.mysqlEngine
	affected, err := engine.Where("id=?", orderId).Delete(new(datastruct.MemberLevelOrder))
	if err != nil || affected <= 0 {
		log.Debug("DeleteMemberOrder err")
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

func (handle *DBHandler) EditRandomLotteryGoods(body *datastruct.EditRandomLotteryGoodsBody) datastruct.CodeType {
	engine := handle.mysqlEngine
	goods := new(datastruct.RandomLotteryGoods)
	goods.Probability = body.Probability
	goods.CreatedAt = time.Now().Unix()
	goods.Brand = body.Brand
	goods.GoodsClassId = body.Classid
	goods.GoodsDesc = body.Goodsdesc
	goods.Name = body.Name
	goods.Price = body.Price
	goods.IsHidden = body.IsHidden
	goods.ImgName = ""
	if !strings.Contains(body.Base64str, conf.Server.Domain) {
		if body.Goodsid > 0 {
			tmp := new(datastruct.Goods)
			has, err := engine.Where("id=?", body.Goodsid).Get(tmp)
			if err != nil || !has {
				log.Debug("DBHandler->EditGoods GetGoods err")
				return datastruct.UpdateDataFailed
			}
			filePath := tools.GetImgPath() + tmp.ImgName
			tools.DeleteFile(filePath)
		}
		imgName := fmt.Sprintf("%s.png", tools.UniqueId())
		path := tools.GetImgPath() + imgName
		arr_str := strings.Split(body.Base64str, ",")
		var isError bool
		if len(arr_str) > 1 {
			isError = tools.CreateImgFromBase64(&arr_str[1], path)
		} else {
			isError = tools.CreateImgFromBase64(&arr_str[0], path)
		}
		if isError {
			log.Debug("DBHandler->EditGoods CreateImgFromBase64 err")
			return datastruct.UpdateDataFailed
		}
		goods.ImgName = imgName
	}

	var err error
	var affected int64
	if body.Goodsid <= 0 {
		affected, err = engine.Insert(goods)
		if err != nil || affected <= 0 {
			log.Debug("DBHandler->EditGoods InsertGoods err")
			return datastruct.UpdateDataFailed
		}
	} else {
		if goods.ImgName == "" {
			affected, err = engine.Where("id=?", body.Goodsid).Cols("name", "price", "goods_desc", "brand", "goods_class_id", "is_hidden", "created_at", "probability").Update(goods)
		} else {
			affected, err = engine.Where("id=?", body.Goodsid).Cols("name", "price", "goods_desc", "brand", "goods_class_id", "is_hidden", "created_at", "probability", "img_name").Update(goods)
		}
		if err != nil {
			log.Debug("DBHandler->EditGoods UpdateGoods err")
			return datastruct.UpdateDataFailed
		}
	}
	return datastruct.NULLError
}

type webRandomLotteryGoodsData struct {
	datastruct.RandomLotteryGoods `xorm:"extends"`
	datastruct.GoodsClass         `xorm:"extends"`
}

func (handle *DBHandler) WebGetRandomLotteryGoods(body *datastruct.WebGetGoodsBody) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	goods := make([]*webRandomLotteryGoodsData, 0)
	start := (body.PageIndex - 1) * body.PageSize
	limit := body.PageSize
	query := ""
	args := make([]interface{}, 0)
	if body.Classid <= 0 {
		query = "goods_class.id > ?"
		args = append(args, 0)
	} else {
		query = "goods_class.id = ?"
		args = append(args, body.Classid)
	}
	if body.IsHidden >= 2 {
		query += " and random_lottery_goods.is_hidden >= ?"
		args = append(args, 0)
	} else {
		query += " and random_lottery_goods.is_hidden = ?"
		args = append(args, body.IsHidden)
	}
	if body.Name != "" {
		query += " and random_lottery_goods.name like ?"
		args = append(args, "%"+body.Name+"%")
	}
	tmp_total := new(datastruct.RandomLotteryGoods)
	count, _ := engine.Table("random_lottery_goods").Join("INNER", "goods_class", "goods_class.id = random_lottery_goods.goods_class_id").Where(query, args...).Desc("random_lottery_goods.created_at").Count(tmp_total)
	engine.Table("random_lottery_goods").Join("INNER", "goods_class", "goods_class.id = random_lottery_goods.goods_class_id").Where(query, args...).Desc("random_lottery_goods.created_at").Limit(limit, start).Find(&goods)
	resp_goods := make([]*datastruct.WebResponseRandomLotteryGoodsData, 0, len(goods))
	for _, v := range goods {
		resp_good := new(datastruct.WebResponseRandomLotteryGoodsData)
		resp_good.ImgUrl = tools.CreateGoodsImgUrl(v.RandomLotteryGoods.ImgName)
		resp_good.Brand = v.RandomLotteryGoods.Brand
		resp_good.GoodsDesc = v.RandomLotteryGoods.GoodsDesc
		resp_good.Id = v.RandomLotteryGoods.Id
		resp_good.IsHidden = v.RandomLotteryGoods.IsHidden
		resp_good.Name = v.RandomLotteryGoods.Name
		resp_good.Price = v.RandomLotteryGoods.Price
		resp_good.Probability = v.RandomLotteryGoods.Probability
		resp_good.GoodsClass = v.GoodsClass.Name
		resp_goods = append(resp_goods, resp_good)
	}
	resp := new(datastruct.WebResponseRandomLotteryGoods)
	resp.List = resp_goods
	resp.Total = int(count)
	return resp, datastruct.NULLError
}

func (handle *DBHandler) WebGetRandomLotteryGoodsPool() (interface{}, datastruct.CodeType) {
	goodsPool := new(datastruct.RandomLotteryGoodsPool)
	engine := handle.mysqlEngine
	has, err := engine.Where("id=?", datastruct.DefaultId).Get(goodsPool)
	if err != nil {
		return nil, datastruct.GetDataFailed
	}
	resp := new(datastruct.WebResponseRandomLotteryPool)
	if has {
		resp.Current = goodsPool.Current
		resp.Probability = goodsPool.Probability
		resp.RandomLotteryCount = goodsPool.RandomLotteryCount
	}
	return resp, datastruct.NULLError
}

func (handle *DBHandler) EditRandomLotteryGoodsPool(body *datastruct.WebResponseRandomLotteryPool) datastruct.CodeType {
	goodsPool := new(datastruct.RandomLotteryGoodsPool)
	engine := handle.mysqlEngine
	has, err := engine.Where("id=?", datastruct.DefaultId).Get(goodsPool)
	if err != nil {
		return datastruct.GetDataFailed
	}
	goodsPool.Current = body.Current
	goodsPool.Probability = body.Probability
	goodsPool.RandomLotteryCount = body.RandomLotteryCount
	if has {
		_, err = engine.Where("id=?", datastruct.DefaultId).Update(goodsPool)
		if err != nil {
			return datastruct.UpdateDataFailed
		}
	} else {
		goodsPool.Id = datastruct.DefaultId
		affected, err := engine.Insert(goodsPool)
		if err != nil || affected <= 0 {
			return datastruct.UpdateDataFailed
		}
	}
	return datastruct.NULLError
}

const randomLottery_inner_str = "from send_goods s inner join random_lottery_goods_succeed r on s.order_id = r.order_id inner join user_info u on r.user_id = u.id inner join random_lottery_goods g on g.id = r.lottery_goods_id"

func (handle *DBHandler) getAllLotteryOrderValue() (int, int) {
	count_sql := "select count(*) " + randomLottery_inner_str
	total := handle.getOrderCount(count_sql)
	today_unix, tomorrow_unix := tools.GetTodayTomorrowTime()
	today_count_sql := fmt.Sprintf(count_sql+" where s.created_at >= %d and s.created_at < %d", today_unix, tomorrow_unix)
	today_total := handle.getOrderCount(today_count_sql)
	return total, today_total
}

func (handle *DBHandler) GetRandomLotteryOrder(body *datastruct.GetSendGoodsBody) (interface{}, datastruct.CodeType) {
	resp := new(datastruct.WebResponseLotteryOrderInfo)
	resp.Total, resp.TodayCount = handle.getAllLotteryOrderValue()

	query := " where is_lottery_goods = 1"
	if body.OrderId != "" {
		query += " and r.order_id like " + "'%" + body.OrderId + "%'"
	}
	if body.UserName != "" {
		query += " and u.nick_name like " + "'%" + body.UserName + "%'"
	}
	if body.GoodsName != "" {
		query += " and g.name like " + "'%" + body.GoodsName + "%'"
	}
	if body.EndTime > 0 && body.StartTime > 0 && body.EndTime > body.StartTime {
		query += " and r.created_at >= " + tools.Int64ToString(body.StartTime) + " and r.created_at < " + tools.Int64ToString(body.EndTime)
	}

	var current_count_sql string

	switch body.State {
	case datastruct.Sended:
		query = query + " and s.send_goods_state = 1"
		current_count_sql = "select count(*) " + randomLottery_inner_str + query

	case datastruct.NotSend:
		query = query + " and s.send_goods_state = 0"
		current_count_sql = "select count(*) " + randomLottery_inner_str + query
	default:
		current_count_sql = "select count(*) " + randomLottery_inner_str + query
	}

	current_total := handle.getOrderCount(current_count_sql)
	resp.CurrentTotal = current_total

	engine := handle.mysqlEngine
	start := (body.PageIndex - 1) * body.PageSize
	limit := body.PageSize
	limitStr := fmt.Sprintf(" LIMIT %d,%d", start, limit)

	sql := "select s.express_agency,s.express_number,s.address,s.phone_number,s.link_man,r.order_id,u.nick_name,u.avatar,g.img_name,g.name,g.price,s.send_goods_state,r.created_at " + randomLottery_inner_str
	orderby := " ORDER BY s.created_at desc"
	sql = sql + query + orderby + limitStr
	results, _ := engine.Query(sql)
	list := make([]*datastruct.SendGoodsInfo, 0, len(results))

	for _, v := range results {
		sendGoodsInfo := new(datastruct.SendGoodsInfo)
		sendGoodsInfo.Avatar = string(v["avatar"][:])
		imgName := string(v["img_name"][:])
		sendGoodsInfo.Count = 1
		sendGoodsInfo.GoodsImgUrl = tools.CreateGoodsImgUrl(imgName)
		sendGoodsInfo.GoodsName = string(v["name"][:])
		sendGoodsInfo.OrderId = string(v["order_id"][:])
		sendGoodsInfo.ReMark = "抽奖成功"
		sendGoodsInfo.Price = tools.StringToInt64(string(v["price"][:]))
		sendGoodsInfo.State = datastruct.SendGoodsType(tools.StringToInt(string(v["send_goods_state"][:])))

		receiver := new(datastruct.ReceiverForSendGoods)
		receiver.Address = string(v["address"][:])
		receiver.LinkMan = string(v["link_man"][:])
		receiver.PhoneNumber = string(v["phone_number"][:])
		sendGoodsInfo.Receiver = receiver

		if sendGoodsInfo.State == datastruct.Sended {
			sender := new(datastruct.SenderForSendGoods)
			sender.ExpressAgency = string(v["express_agency"][:])
			sender.ExpressNumber = string(v["express_number"][:])
			sendGoodsInfo.Sender = sender
		}

		sendGoodsInfo.Time = tools.StringToInt64(string(v["created_at"][:]))
		sendGoodsInfo.UserName = string(v["nick_name"][:])
		list = append(list, sendGoodsInfo)
	}
	resp.List = list

	return resp, datastruct.NULLError
}

func (handle *DBHandler) UpdateLotteryGoodsSendState(body *datastruct.WebResponseLotteryGoodsSendStateBody) datastruct.CodeType {
	engine := handle.mysqlEngine
	updateSend := make([]*updateSendData, 0)
	err := engine.Table("send_goods").Join("INNER", "random_lottery_goods_succeed", "random_lottery_goods_succeed.order_id = send_goods.order_id").Where("random_lottery_goods_succeed.order_id=?", tools.StringToInt64(body.OrderNumber)).Find(&updateSend)
	if err != nil || len(updateSend) <= 0 || updateSend[0].SendGoodsState == datastruct.Sended {
		return datastruct.UpdateDataFailed
	}

	sendgoods := new(datastruct.SendGoods)
	sendgoods.SendGoodsState = datastruct.Sended
	sendgoods.ExpressNumber = body.ExpressNumber
	sendgoods.ExpressAgency = body.ExpressAgency
	sendgoods.LinkMan = body.LinkMan
	sendgoods.PhoneNumber = body.PhoneNumber
	sendgoods.Address = body.Address

	_, err = engine.Where("order_id=?", updateSend[0].OrderId).Cols("send_goods_state", "express_number", "express_agency", "link_man", "phone_number", "address").Update(sendgoods)
	if err != nil {
		log.Debug("UpdateLotteryGoodsSendState err:%s", err.Error())
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

type webStatisticsData struct {
	datastruct.RandomLotteryGoodsSucceed `xorm:"extends"`
	datastruct.UserInfo                  `xorm:"extends"`
}

func (handle *DBHandler) WebStatistics(body *datastruct.WebNewsUserBody) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	users := make([]*webStatisticsData, 0)
	engine.Table("random_lottery_goods_succeed").Join("INNER", "user_info", "user_info.id = random_lottery_goods_succeed.user_id").Where("random_lottery_goods_succeed.created_at >= ? and random_lottery_goods_succeed.created_at < ?", body.StartTime, body.EndTime).Find(&users)
	var total, paytotal int
	total = 0
	paytotal = 0
	for _, v := range users {
		ids1, total1, paytotal1 := getAgencyUserCount(engine, []int{v.UserId}, false, body)
		ids2, total2, paytotal2 := getAgencyUserCount(engine, ids1, false, body)
		_, total3, paytotal3 := getAgencyUserCount(engine, ids2, true, body)
		total += total1 + total2 + total3
		paytotal += paytotal1 + paytotal2 + paytotal3
	}
	var rs float64
	if total != 0 {
		rs = float64(paytotal) / float64(total)
	} else {
		rs = 0
	}

	queryRegisterUser := ""
	if body.RPlatform <= datastruct.H5 {
		queryRegisterUser = " and u.platform = " + tools.IntToString(int(body.RPlatform))
	}
	queryPayUser := ""
	if body.PayPlatform <= datastruct.H5 {
		queryPayUser = " and o.platform = " + tools.IntToString(int(body.PayPlatform))
	}

	usersCountPurchase := "select count(*) from order_info o inner join user_info u on o.user_id = u.id inner join goods g on g.id = o.goods_id where o.is_purchase = 1 and o.created_at>= ? and o.created_at < ?" + queryRegisterUser + queryPayUser + " GROUP BY o.user_id"

	purchasePriceSum := "select sum(g.price) from order_info o inner join user_info u on o.user_id = u.id inner join goods g on g.id = o.goods_id where o.is_purchase = 1 and o.created_at>= ? and o.created_at < ?" + queryRegisterUser + queryPayUser

	results, _ := engine.Query(usersCountPurchase, body.StartTime, body.EndTime)
	purchase_users := len(results)

	results, _ = engine.Query(purchasePriceSum, body.StartTime, body.EndTime)
	str_sum := string(results[0]["sum(g.price)"][:])
	purchase_amount := tools.StringToInt64(str_sum)

	if body.PayPlatform <= datastruct.H5 {
		queryPayUser = " and udi.platform = " + tools.IntToString(int(body.PayPlatform))
	}
	usersCountDeposit := "select count(*) from user_deposit_info udi join user_info u on udi.user_id = u.id where udi.created_at >= ? and udi.created_at < ?" + queryRegisterUser + queryPayUser + " GROUP BY udi.user_id"
	depositMoneySum := "select sum(udi.money) from user_deposit_info udi join user_info u on udi.user_id = u.id where udi.created_at >= ? and udi.created_at < ?" + queryRegisterUser + queryPayUser

	results, _ = engine.Query(usersCountDeposit, body.StartTime, body.EndTime)
	deposit_users := len(results)

	results, _ = engine.Query(depositMoneySum, body.StartTime, body.EndTime)
	str_sum = string(results[0]["sum(udi.money)"][:])
	deposit_amount := tools.StringToFloat64(str_sum)

	statistics := new(datastruct.WebResponseStatistics)
	statistics.AgentPayRate = rs
	statistics.GrowthCount = total
	statistics.UsersForPurchase = purchase_users
	statistics.PurchaseAmount = purchase_amount
	statistics.UsersForDeposit = deposit_users
	statistics.DepositAmount = deposit_amount
	statistics.UsersForPay = purchase_users + deposit_users
	statistics.UserPayAmount = float64(purchase_amount) + deposit_amount
	statistics.Date = tools.UnixToString(body.StartTime, "2006-01-02")
	return statistics, datastruct.NULLError
}

func (handle *DBHandler) GetActiveUsers(body *datastruct.WebActiveUserBody) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	queryRegisterUser := ""
	if body.RPlatform <= datastruct.H5 {
		queryRegisterUser = " and platform = " + tools.IntToString(int(body.RPlatform))
	}
	statistics := new(datastruct.WebResponseActiveUsers)
	statistics.Date = tools.UnixToString(body.StartTime, "2006-01-02")

	newUserSql := "select count(*) from user_info where created_at >= ? and created_at < ?" + queryRegisterUser
	results, _ := engine.Query(newUserSql, body.StartTime, body.EndTime)
	strTotal := string(results[0]["count(*)"][:])
	statistics.NewUsers = tools.StringToInt64(strTotal)

	activeUserSql := "select count(*) from user_info where login_time >= ? and login_time < ?" + queryRegisterUser
	results, _ = engine.Query(activeUserSql, body.StartTime, body.EndTime)
	strTotal = string(results[0]["count(*)"][:])
	statistics.ActiveUsers = tools.StringToInt64(strTotal)

	return statistics, datastruct.NULLError
}

func (handle *DBHandler) GetCommissionStatistics(body *datastruct.WebActiveUserBody) (interface{}, datastruct.CodeType) {
	/*
		engine := handle.mysqlEngine
		queryRegisterUser := ""
		if body.RPlatform <= datastruct.H5 {
			queryRegisterUser = " and platform = " + tools.IntToString(int(body.RPlatform))
		}
		statistics := new(datastruct.WebResponseActiveUsers)
		statistics.Date = tools.UnixToString(body.StartTime, "2006-01-02")

		newUserSql := "select count(*) from user_info where created_at >= ? and created_at < ?" + queryRegisterUser
		results, _ := engine.Query(newUserSql, body.StartTime, body.EndTime)
		strTotal := string(results[0]["count(*)"][:])
		statistics.NewUsers = tools.StringToInt64(strTotal)

		activeUserSql := "select count(*) from user_info where login_time >= ? and login_time < ?" + queryRegisterUser
		results, _ = engine.Query(activeUserSql, body.StartTime, body.EndTime)
		strTotal = string(results[0]["count(*)"][:])
		statistics.ActiveUsers = tools.StringToInt64(strTotal)
		return statistics, datastruct.NULLError
	*/
	return nil, datastruct.NULLError
}

func getAgencyUserCount(engine *xorm.Engine, users []int, isEnd bool, body *datastruct.WebNewsUserBody) ([]int, int, int) {
	if users == nil || len(users) <= 0 {
		return nil, 0, 0
	}
	ids_str := getUsersStr(users)
	var sql string
	if body.StartTime != 0 && body.EndTime != 0 {
		sql = "select u.id,u.deposit_total,u.purchase_total from invite_info i inner join user_info u on i.receiver = u.id where u.created_at >= " + tools.Int64ToString(body.StartTime) + " and u.created_at < " + tools.Int64ToString(body.EndTime) + " and i.sender in (" + ids_str + ")"
	} else {
		sql = "select u.id,u.deposit_total,u.purchase_total from invite_info i inner join user_info u on i.receiver = u.id where i.sender in (" + ids_str + ")"
	}
	results, _ := engine.Query(sql)
	count := len(results)
	paycount := 0
	ids := make([]int, 0, count)
	if !isEnd {
		for _, v := range results {
			ids = append(ids, tools.StringToInt(string(v["id"][:])))
			depositTotal := tools.StringToInt64(string(v["deposit_total"][:]))
			purchaseTotal := tools.StringToFloat64(string(v["purchase_total"][:]))
			if depositTotal > 0 || purchaseTotal > 0 {
				paycount += 1
			}
		}
	}
	return ids, count, paycount
}

func (handle *DBHandler) EditReClass(body *datastruct.WebEditReClassBody) datastruct.CodeType {
	engine := handle.mysqlEngine
	reclass := new(datastruct.RecommendedClass)
	reclass.Name = body.Name
	session := engine.NewSession()
	defer session.Close()
	session.Begin()

	isUpdateIcon := false
	if !strings.Contains(body.Icon, conf.Server.Domain) && body.Icon != "" {
		isUpdateIcon = true
		if body.Id > 0 {
			tmp := new(datastruct.RecommendedClass)
			has, err := session.Where("id=?", body.Id).Get(tmp)
			if err != nil || !has {
				rollback("DBHandler->EditReClass Get ReClass err", session)
				return datastruct.UpdateDataFailed
			}
			deleteFile(tools.GetImgPath() + tmp.Icon)
		}
		imgName := fmt.Sprintf("%s.png", tools.UniqueId())
		path := tools.GetImgPath() + imgName
		arr_str := strings.Split(body.Icon, ",")
		var isError bool
		if len(arr_str) > 1 {
			isError = tools.CreateImgFromBase64(&arr_str[1], path)
		} else {
			isError = tools.CreateImgFromBase64(&arr_str[0], path)
		}
		if isError {
			rollback("DBHandler->EditReClass CreateImgFromBase64 err", session)
			return datastruct.UpdateDataFailed
		}
		reclass.Icon = imgName
	}

	var err error
	if body.Id == 0 {
		_, err = session.Insert(reclass)
	} else {
		if isUpdateIcon {
			_, err = session.Where("id = ?", body.Id).Cols("name,icon").Update(reclass)
		} else {
			_, err = session.Where("id = ?", body.Id).Cols("name").Update(reclass)
		}
	}

	if err != nil {
		deleteFile(tools.GetImgPath() + reclass.Icon)
		rollback("DBHandler->EditReClass err:%s"+err.Error(), session)
		return datastruct.UpdateDataFailed
	}

	err = session.Commit()
	if err != nil {
		deleteFile(tools.GetImgPath() + reclass.Icon)
		str := fmt.Sprintf("DBHandler->EditReClass Commit :%s", err.Error())
		rollback(str, session)
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

func (handle *DBHandler) GetAllReClass(body *datastruct.WebQueryReClassBody) (interface{}, datastruct.CodeType) {
	reclasses := make([]*datastruct.RecommendedClass, 0)
	engine := handle.mysqlEngine
	query := "1=1"
	args := make([]interface{}, 0)
	if body.Name != "" {
		query += " and name like ?"
		args = append(args, "%"+body.Name+"%")
	}
	err := engine.Where(query, args...).Find(&reclasses)
	if err != nil {
		return nil, datastruct.GetDataFailed
	}

	resp := make([]*datastruct.WebResponseReClass, 0)
	for _, v := range reclasses {
		reclass := new(datastruct.WebResponseReClass)
		reclass.Icon = tools.CreateGoodsImgUrl(v.Icon)
		reclass.Id = v.Id
		reclass.Name = v.Name
		resp = append(resp, reclass)
	}

	return resp, datastruct.NULLError
}

func (handle *DBHandler) EditSharePoster(body *datastruct.WebEditSharePosterBody, bucket *oss.Bucket) datastruct.CodeType {
	engine := handle.mysqlEngine
	sp := new(datastruct.SharePosters)
	sp.IsHidden = body.IsHidden
	sp.SortId = body.SortId
	sp.Location = body.Location

	isUpdateImg := false
	isUpdateIcon := false
	isContains := strings.Contains(body.ImgUrl, datastruct.OSSEndpoint)
	tmp := new(datastruct.SharePosters)
	if body.Id > 0 {
		has, err := engine.Where("id=?", body.Id).Get(tmp)
		if err != nil || !has {
			log.Debug("DBHandler->EditSharePosters Get ReClass err")
			return datastruct.UpdateDataFailed
		}
	}
	if !isContains {
		isUpdateImg = true
		if body.Id > 0 {
			osstool.DeleteFile(bucket, tmp.ImgName)
		}
		sp.ImgName = body.ImgUrl
	}
	isContains = strings.Contains(body.Icon, datastruct.OSSEndpoint)
	if !isContains {
		isUpdateIcon = true
		if body.Id > 0 {
			osstool.DeleteFile(bucket, tmp.IconName)
		}
		sp.IconName = body.Icon
	}

	var err1, err2, err3 error
	if body.Id == 0 {
		_, err1 = engine.Insert(sp)
	} else {
		_, err1 = engine.Where("id = ?", body.Id).Cols("sort_id", "is_hidden", "location").Update(sp)
		if isUpdateImg {
			_, err2 = engine.Where("id = ?", body.Id).Cols("img_name").Update(sp)
		}
		if isUpdateIcon {
			_, err3 = engine.Where("id = ?", body.Id).Cols("icon_name").Update(sp)
		}
	}

	if err1 != nil || err2 != nil || err3 != nil {
		if err2 != nil {
			osstool.DeleteFile(bucket, sp.ImgName)
		}
		if err3 != nil {
			osstool.DeleteFile(bucket, sp.IconName)
		}
		log.Debug("DBHandler->EditSharePosters eidt err")
		return datastruct.UpdateDataFailed
	}

	return datastruct.NULLError
}

func (handle *DBHandler) GetAllSharePosters(body *datastruct.WebQuerySharePostersBody) (interface{}, datastruct.CodeType) {
	sp := make([]*datastruct.SharePosters, 0)
	engine := handle.mysqlEngine
	query := ""
	args := make([]interface{}, 0, 1)
	if body.IsHidden == 2 {
		query = "is_hidden >= ?"
		args = append(args, 0)
	} else {
		query = "is_hidden = ?"
		args = append(args, body.IsHidden)
	}
	err := engine.Where(query, args...).Find(&sp)
	if err != nil {
		return nil, datastruct.GetDataFailed
	}

	resp := make([]*datastruct.WebEditSharePosterBody, 0)
	for _, v := range sp {
		wsp := new(datastruct.WebEditSharePosterBody)
		wsp.ImgUrl = osstool.CreateOSSURL(v.ImgName)
		wsp.Icon = osstool.CreateOSSURL(v.IconName)
		wsp.Id = v.Id
		wsp.IsHidden = v.IsHidden
		wsp.SortId = v.SortId
		wsp.Location = v.Location
		resp = append(resp, wsp)
	}
	return resp, datastruct.NULLError
}

func (handle *DBHandler) UpdateSharePosterState(body *datastruct.WebHiddenPostersBody) datastruct.CodeType {
	engine := handle.mysqlEngine
	sp := new(datastruct.SharePosters)
	sp.IsHidden = body.IsHidden
	_, err := engine.Where("id = ?", body.Id).Cols("is_hidden").Update(sp)
	if err != nil {
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

func (handle *DBHandler) EditUserAppraise(body *datastruct.WebEditUserAppraiseBody, bucket *oss.Bucket) datastruct.CodeType {
	engine := handle.mysqlEngine

	has, err := engine.Where("id=?", body.UserId).Get(new(datastruct.UserInfo))
	if err != nil {
		log.Debug("EditUserAppraise GetUser err:%v", err.Error())
		return datastruct.GetDataFailed
	}
	if !has {
		return datastruct.NotExistUser
	}
	switch body.GoodsType {
	case datastruct.RushGoods:
		has, err = engine.Where("id=?", body.GoodsId).Get(new(datastruct.Goods))
	case datastruct.LotteryGoods:
		has, err = engine.Where("id=?", body.GoodsId).Get(new(datastruct.RandomLotteryGoods))
	}

	if err != nil {
		log.Debug("EditUserAppraise GetGoods err:%v", err.Error())
		return datastruct.GetDataFailed
	}
	if !has {
		return datastruct.NotExistGoods
	}

	uap := new(datastruct.UserAppraise)
	uap.Desc = body.Desc
	uap.GoodsId = body.GoodsId
	uap.IsPassed = body.IsPassed
	uap.UserId = body.UserId
	uap.GoodsType = body.GoodsType
	uap.CreatedAt = body.TimeStamp
	session := engine.NewSession()
	defer session.Close()
	session.Begin()
	length := len(body.ImgNames)
	if length < 0 {
		uap.ShowType = datastruct.OnlyT
	} else if body.Desc == "" {
		uap.ShowType = datastruct.OnlyP
	} else {
		uap.ShowType = datastruct.TAndP
	}
	var affected int64
	if body.Id > 0 {
		isExist := new(datastruct.UserAppraise)
		has, err = session.Where("id=?", body.Id).Get(isExist)
		if err != nil || !has {
			rollback("DBHandler->EditUserAppraise not exist", session)
			return datastruct.NotExist
		}
		tmp := make([]*datastruct.UserAppraisePic, 0)
		err = session.Where("user_appraise_id = ?", body.Id).Find(&tmp)
		if err != nil {
			rollback("DBHandler->EditUserAppraise Imgs err:"+err.Error(), session)
			return datastruct.UpdateDataFailed
		}
		var isDelete bool
		for _, v_imgs := range tmp {
			isDelete = true
			for _, v := range body.ImgNames {
				if strings.Contains(v, datastruct.OSSEndpoint) {
					str_arr := strings.Split(v, datastruct.OSSEndpoint+"/")
					filename := str_arr[len(str_arr)-1]
					if v_imgs.ImgName == filename {
						isDelete = false
					}
				}
			}
			if isDelete {
				osstool.DeleteFile(bucket, v_imgs.ImgName)
			}
		}
		for i, v := range body.ImgNames {
			if strings.Contains(v, datastruct.OSSEndpoint) {
				str_arr := strings.Split(v, datastruct.OSSEndpoint+"/")
				filename := str_arr[len(str_arr)-1]
				tmp := new(datastruct.UserAppraisePic)
				has, err = session.Where("user_appraise_id=? and img_index=?", body.Id, i).Get(tmp)
				if err != nil {
					rollback("DBHandler->EditUserAppraise 0 get err:"+err.Error(), session)
					return datastruct.GetDataFailed
				}
				if has {
					tmp.ImgName = filename
					_, err = session.Where("user_appraise_id=? and img_index=?", body.Id, i).Cols("img_name").Update(tmp)
					if err != nil {
						rollback("DBHandler->EditUserAppraise 0 Update Imgs err:"+err.Error(), session)
						return datastruct.UpdateDataFailed
					}
				} else {
					userAppraisePic := new(datastruct.UserAppraisePic)
					userAppraisePic.UserAppraiseId = body.Id
					userAppraisePic.ImgName = filename
					userAppraisePic.ImgIndex = i
					affected, err = session.Insert(userAppraisePic)
					if err != nil || affected <= 0 {
						rollback("DBHandler->EditUserAppraise 0 Insert Imgs err", session)
						return datastruct.UpdateDataFailed
					}
				}
			}
		}

		_, err = session.Where("user_appraise_id = ? and img_index >= ?", body.Id, length).Delete(new(datastruct.UserAppraisePic))
		if err != nil {
			rollback("DBHandler->EditUserAppraise DeleteGoodsImgs err:"+err.Error(), session)
			return datastruct.UpdateDataFailed
		}
		_, err = session.Where("id=?", body.Id).Cols("goods_type", "user_id", "desc", "show_type", "goods_id", "is_passed", "created_at").Update(uap)
		if err != nil {
			rollback("DBHandler->EditUserAppraise Update err:"+err.Error(), session)
			return datastruct.UpdateDataFailed
		}
	} else {
		_, err = session.Insert(uap)
		if err != nil {
			rollback("DBHandler->EditUserAppraise Insert err:"+err.Error(), session)
			return datastruct.UpdateDataFailed
		}
	}

	var userAppraiseId int
	for i := 0; i < length; i++ {
		v := body.ImgNames[i]
		if !strings.Contains(v, datastruct.OSSEndpoint) {
			if body.Id > 0 {
				userAppraiseId = body.Id
			} else {
				userAppraiseId = uap.Id
			}
			tmp := new(datastruct.UserAppraisePic)
			has, err = session.Where("user_appraise_id = ? and img_index = ?", userAppraiseId, i).Get(tmp)
			if err != nil {
				rollback("DBHandler->EditUserAppraise get err:"+err.Error(), session)
				return datastruct.GetDataFailed
			}
			if has {
				tmp.ImgName = v
				_, err = session.Where("user_appraise_id = ? and img_index = ?", userAppraiseId, i).Cols("img_name").Update(tmp)
				if err != nil {
					rollback("DBHandler->EditUserAppraise Delete err:"+err.Error(), session)
					return datastruct.UpdateDataFailed
				}
			} else {
				uapp := new(datastruct.UserAppraisePic)
				uapp.ImgIndex = i
				uapp.ImgName = v
				uapp.UserAppraiseId = userAppraiseId
				affected, err = session.Insert(uapp)
				if err != nil || affected <= 0 {
					rollback("DBHandler->EditUserAppraise Insert Imgs err", session)
					return datastruct.UpdateDataFailed
				}
			}
		}
	}

	err = session.Commit()
	if err != nil {
		str := fmt.Sprintf("DBHandler->EditUserAppraise Commit :%s", err.Error())
		rollback(str, session)
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

func (handle *DBHandler) DeleteUserAppraise(body *datastruct.WebDeleteUserAppraiseBody, bucket *oss.Bucket) datastruct.CodeType {
	engine := handle.mysqlEngine
	session := engine.NewSession()
	defer session.Close()
	session.Begin()

	_, err := session.Where("id = ?", body.Id).Delete(new(datastruct.UserAppraise))
	if err != nil {
		str := fmt.Sprintf("DBHandler->DeleteUserAppraise delete UserAppraise :%s", err.Error())
		rollback(str, session)
		return datastruct.UpdateDataFailed
	}
	tmp_arr := make([]*datastruct.UserAppraisePic, 0)
	session.Where("user_appraise_id = ?", body.Id).Find(&tmp_arr)
	for _, v := range tmp_arr {
		osstool.DeleteFile(bucket, v.ImgName)
	}

	_, err = session.Where("user_appraise_id = ?", body.Id).Delete(new(datastruct.UserAppraisePic))
	if err != nil {
		str := fmt.Sprintf("DBHandler->DeleteUserAppraise delete UserAppraisePic :%s", err.Error())
		rollback(str, session)
		return datastruct.UpdateDataFailed
	}

	err = session.Commit()
	if err != nil {
		str := fmt.Sprintf("DBHandler->DeleteUserAppraise Commit :%s", err.Error())
		rollback(str, session)
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

func (handle *DBHandler) GetUserAppraise(body *datastruct.WebQueryUserAppraiseBody) (interface{}, datastruct.CodeType) {
	query := fmt.Sprintf(" is_passed = %d and uap.goods_type = %d", body.IsPassed, body.GoodsType)

	if body.Desc != "" {
		query += " and uap.desc like " + "'%" + body.Desc + "%'"
	}
	if body.GoodsName != "" {
		switch body.GoodsType {
		case datastruct.RushGoods:
			query += " and g.name like " + "'%" + body.GoodsName + "%'"
		case datastruct.LotteryGoods:
			query += " and rlg.name like " + "'%" + body.GoodsName + "%'"
		}
	}
	if body.UserName != "" {
		query += " and u.nick_name like " + "'%" + body.UserName + "%'"
	}

	start := (body.PageIndex - 1) * body.PageSize
	limit := body.PageSize
	engine := handle.mysqlEngine
	limitStr := fmt.Sprintf(" LIMIT %d,%d", start, limit)

	var sql_count string
	var sql string
	switch body.GoodsType {
	case datastruct.RushGoods:
		sql_count = "select count(*) from user_appraise uap inner join user_info u on u.id = uap.user_id inner join goods g on uap.goods_id = g.id where "
		sql = "select g.id as gid,uap.show_type,uap.id,uap.desc,uap.user_id,u.nick_name,u.avatar,g.name as gname,uap.created_at from user_appraise uap inner join user_info u on u.id = uap.user_id inner join goods g on uap.goods_id = g.id where "
	case datastruct.LotteryGoods:
		sql_count = "select count(*) from user_appraise uap inner join user_info u on u.id = uap.user_id inner join random_lottery_goods rlg on uap.goods_id = rlg.id where "
		sql = "select rlg.id as gid,uap.show_type,uap.id,uap.desc,uap.user_id,u.nick_name,u.avatar,rlg.name as gname,uap.created_at from user_appraise uap inner join user_info u on u.id = uap.user_id inner join random_lottery_goods rlg on uap.goods_id = rlg.id where "
	}
	results, _ := engine.Query(sql_count + query)
	strTotal := string(results[0]["count(*)"][:])
	currentTotal := tools.StringToInt64(strTotal)

	orderby := " ORDER BY uap.created_at desc"
	results, _ = engine.Query(sql + query + orderby + limitStr)
	list := make([]*datastruct.WebUserAppraise, 0, len(results))
	for _, v := range results {
		web_uap := new(datastruct.WebUserAppraise)
		web_uap.Avatar = string(v["avatar"][:])
		web_uap.CreatedAt = tools.StringToInt64(string(v["created_at"][:]))
		web_uap.Desc = string(v["desc"][:])
		web_uap.GoodsType = body.GoodsType
		web_uap.GoodsName = string(v["gname"][:])
		web_uap.GoodsId = tools.StringToInt(string(v["gid"][:]))
		web_uap.Id = tools.StringToInt(string(v["id"][:]))
		web_uap.NickName = string(v["nick_name"][:])
		web_uap.UserId = tools.StringToInt(string(v["user_id"][:]))
		showType := tools.StringToInt(string(v["show_type"][:]))
		if datastruct.UserAppraiseType(showType) != datastruct.OnlyT {
			tmp_uapps := make([]*datastruct.UserAppraisePic, 0)
			engine.Where("user_appraise_id = ?", web_uap.Id).Asc("img_index").Find(&tmp_uapps)
			imgurls := make([]string, 0, len(tmp_uapps))
			for _, v := range tmp_uapps {
				url := osstool.CreateOSSURL(v.ImgName)
				imgurls = append(imgurls, url)
			}
			web_uap.ImgUrls = imgurls
		}
		web_uap.IsPassed = body.IsPassed
		list = append(list, web_uap)
	}
	resp := new(datastruct.WebResponseUserAppraise)
	resp.CurrentTotal = currentTotal
	resp.List = list
	return resp, datastruct.NULLError
}

func (handle *DBHandler) UpdateSignForState(number string) datastruct.CodeType {
	engine := handle.mysqlEngine
	updateSend := make([]*updateSendData, 0)
	err := engine.Table("send_goods").Join("INNER", "order_info", "order_info.id = send_goods.order_id").Where("order_info.number=?", number).Find(&updateSend)
	if err != nil || len(updateSend) <= 0 || updateSend[0].SignForState == datastruct.GoodsSigned {
		return datastruct.UpdateDataFailed
	}
	if updateSend[0].SendGoodsState != datastruct.Sended {
		return datastruct.NotHasSendGoods
	}
	sendgoods := new(datastruct.SendGoods)
	sendgoods.SignForState = datastruct.GoodsSigned
	var affected int64
	affected, err = engine.Where("order_id=?", updateSend[0].OrderId).Cols("sign_for_state").Update(sendgoods)
	if affected <= 0 || err != nil {
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

func (handle *DBHandler) EditGoodsDetail(body *datastruct.WebEditGoodsDetailBody, bucket *oss.Bucket) datastruct.CodeType {
	engine := handle.mysqlEngine
	session := engine.NewSession()
	defer session.Close()
	session.Begin()
	length := len(body.ImgNames)

	tmp := make([]*datastruct.GoodsDetail, 0)
	err := session.Where("goods_id = ?", body.GoodsId).Find(&tmp)
	if err != nil {
		rollback("DBHandler->EditGoodsDetail GetImgs err:"+err.Error(), session)
		return datastruct.UpdateDataFailed
	}
	var has bool
	var isDelete bool
	var affected int64
	for _, v_imgs := range tmp {
		isDelete = true
		for _, v := range body.ImgNames {
			if strings.Contains(v, datastruct.OSSEndpoint) {
				str_arr := strings.Split(v, datastruct.OSSEndpoint+"/")
				filename := str_arr[len(str_arr)-1]
				if v_imgs.ImgName == filename {
					isDelete = false
				}
			}
		}
		if isDelete {
			osstool.DeleteFile(bucket, v_imgs.ImgName)
		}
	}

	for i, v := range body.ImgNames {
		if strings.Contains(v, datastruct.OSSEndpoint) {
			str_arr := strings.Split(v, datastruct.OSSEndpoint+"/")
			filename := str_arr[len(str_arr)-1]
			tmp := new(datastruct.GoodsDetail)
			has, err = session.Where("goods_id=? and img_index=?", body.GoodsId, i).Get(tmp)
			if err != nil {
				rollback("DBHandler->EditGoodsDetail 0 get err:"+err.Error(), session)
				return datastruct.GetDataFailed
			}
			if has {
				tmp.ImgName = filename
				_, err = session.Where("goods_id=? and img_index=?", body.GoodsId, i).Cols("img_name").Update(tmp)
				if err != nil {
					rollback("DBHandler->EditGoodsDetail 0 Update GoodsImgs err:"+err.Error(), session)
					return datastruct.UpdateDataFailed
				}
			} else {
				goodsDetail := new(datastruct.GoodsDetail)
				goodsDetail.GoodsId = body.GoodsId
				goodsDetail.ImgName = filename
				goodsDetail.ImgIndex = i
				affected, err = session.Insert(goodsDetail)
				if err != nil || affected <= 0 {
					rollback("DBHandler->EditGoodsDetail 0 Insert GoodsImgs err", session)
					return datastruct.UpdateDataFailed
				}
			}
		}
	}

	_, err = session.Where("goods_id=? and img_index >= ?", body.GoodsId, length).Delete(new(datastruct.GoodsDetail))
	if err != nil {
		rollback("DBHandler->EditGoodsDetail DeleteImgs err:"+err.Error(), session)
		return datastruct.UpdateDataFailed
	}

	for i := 0; i < length; i++ {
		v := body.ImgNames[i]
		if !strings.Contains(v, datastruct.OSSEndpoint) {
			tmp := new(datastruct.GoodsDetail)
			has, err = session.Where("goods_id = ? and img_index = ?", body.GoodsId, i).Get(tmp)
			if err != nil {
				rollback("DBHandler->EditGoodsDetail get err:"+err.Error(), session)
				return datastruct.GetDataFailed
			}
			if has {
				tmp.ImgName = v
				_, err = session.Where("goods_id = ? and img_index = ?", body.GoodsId, i).Cols("img_name").Update(tmp)
				if err != nil {
					rollback("DBHandler->EditGoodsDetail Delete err:"+err.Error(), session)
					return datastruct.UpdateDataFailed
				}
			} else {
				gd := new(datastruct.GoodsDetail)
				gd.ImgIndex = i
				gd.ImgName = v
				gd.GoodsId = body.GoodsId
				affected, err = session.Insert(gd)
				if err != nil || affected <= 0 {
					rollback("DBHandler->EditGoodsDetail Insert GoodsImgs err", session)
					return datastruct.UpdateDataFailed
				}
			}
		}
	}

	err = session.Commit()
	if err != nil {
		str := fmt.Sprintf("DBHandler->EditGoodsDetail Commit :%s", err.Error())
		rollback(str, session)
		return datastruct.UpdateDataFailed
	}

	return datastruct.NULLError
}

func (handle *DBHandler) GetGoodsDetail(goodsid int) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	tmp_arr := make([]*datastruct.GoodsDetail, 0)
	engine.Where("goods_id=?", goodsid).Asc("img_index").Find(&tmp_arr)
	resp := make([]string, 0, len(tmp_arr))
	for _, v := range tmp_arr {
		resp = append(resp, osstool.CreateOSSURL(v.ImgName))
	}
	return resp, datastruct.NULLError
}

func (handle *DBHandler) GetSCParams() (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	scp := new(datastruct.SuggestionComplaintParams)
	has, err := engine.Where("id=?", datastruct.DefaultId).Get(scp)
	if err != nil || !has {
		return nil, datastruct.GetDataFailed
	}
	resp := new(datastruct.WebResponseSCP)
	resp.CCFBL = scp.ComplaintCountForBL
	resp.CCFD = scp.ComplaintCountForDay
	resp.SCFD = scp.SuggestCountForDay
	return resp, datastruct.NULLError
}
func (handle *DBHandler) UpdateSCParams(body *datastruct.WebResponseSCP) datastruct.CodeType {
	engine := handle.mysqlEngine
	scp := new(datastruct.SuggestionComplaintParams)
	scp.ComplaintCountForBL = body.CCFBL
	scp.ComplaintCountForDay = body.CCFD
	scp.SuggestCountForDay = body.SCFD
	_, err := engine.Where("id=?", datastruct.DefaultId).Cols("suggest_count_for_day", "complaint_count_for_day", "complaint_count_for_b_l").Update(scp)
	if err != nil {
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

func (handle *DBHandler) GetSuggestion(body *datastruct.WebQuerySuggestBody) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	query := ""
	start := (body.PageIndex - 1) * body.PageSize
	limit := body.PageSize
	limitStr := fmt.Sprintf(" LIMIT %d,%d", start, limit)
	orderby := " ORDER BY s.id desc"
	if body.Desc != "" {
		query += " and s.desc like " + "'%" + body.Desc + "%'"
	}
	if body.UserName != "" {
		query += " and u.nick_name like " + "'%" + body.UserName + "%'"
	}
	sql := "select u.avatar,u.nick_name,s.desc,s.created_at,s.id from suggestion s inner join user_info u on u.id = s.user_id where 1=1"
	sql_count := "select count(*) from suggestion s inner join user_info u on u.id = s.user_id where 1=1"
	result_count, _ := engine.Query(sql_count + query)
	strTotal := string(result_count[0]["count(*)"][:])

	results, _ := engine.Query(sql + query + orderby + limitStr)
	resp := new(datastruct.WebResponseSuggestInfo)
	list := make([]*datastruct.WebResponseSuggest, 0)
	for _, v := range results {
		tmp := new(datastruct.WebResponseSuggest)
		tmp.Avatar = string(v["avatar"][:])
		tmp.NickName = string(v["nick_name"][:])
		tmp.CreatedAt = tools.StringToInt64(string(v["created_at"][:]))
		tmp.Desc = string(v["desc"][:])
		tmp.Id = tools.StringToInt(string(v["id"][:]))
		list = append(list, tmp)
	}
	resp.List = list
	resp.CurrentTotal = tools.StringToInt64(strTotal)
	return resp, datastruct.NULLError
}

func (handle *DBHandler) GetComplaint(body *datastruct.WebQueryComplaintBody) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	query := ""
	start := (body.PageIndex - 1) * body.PageSize
	limit := body.PageSize
	limitStr := fmt.Sprintf(" LIMIT %d,%d", start, limit)
	orderby := " ORDER BY c.id desc"
	if body.ComplaintType != "" {
		query += " and c.complaint_type like " + "'%" + body.ComplaintType + "%'"
	}
	if body.Desc != "" {
		query += " and c.desc like " + "'%" + body.Desc + "%'"
	}
	if body.UserName != "" {
		query += " and u.nick_name like " + "'%" + body.UserName + "%'"
	}
	sql := "select u.avatar,u.nick_name,c.desc,c.created_at,c.id,c.complaint_type from complaint c inner join user_info u on u.id = c.user_id where 1=1"
	sql_count := "select count(*) from complaint c inner join user_info u on u.id = c.user_id where 1=1"
	result_count, _ := engine.Query(sql_count + query)
	strTotal := string(result_count[0]["count(*)"][:])

	results, _ := engine.Query(sql + query + orderby + limitStr)
	resp := new(datastruct.WebResponseComplaintInfo)
	list := make([]*datastruct.WebResponseComplaint, 0)
	for _, v := range results {
		tmp := new(datastruct.WebResponseComplaint)
		tmp.Avatar = string(v["avatar"][:])
		tmp.NickName = string(v["nick_name"][:])
		tmp.CreatedAt = tools.StringToInt64(string(v["created_at"][:]))
		tmp.Desc = string(v["desc"][:])
		tmp.Id = tools.StringToInt(string(v["id"][:]))
		tmp.ComplaintType = string(v["complaint_type"][:])
		list = append(list, tmp)
	}
	resp.List = list
	resp.CurrentTotal = tools.StringToInt64(strTotal)
	return resp, datastruct.NULLError
}

func (handle *DBHandler) DeleteSuggestion(id int) datastruct.CodeType {
	engine := handle.mysqlEngine
	_, err := engine.Where("id = ?", id).Delete(new(datastruct.Suggestion))
	if err != nil {
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

func (handle *DBHandler) DeleteComplaint(id int) datastruct.CodeType {
	engine := handle.mysqlEngine
	_, err := engine.Where("id = ?", id).Delete(new(datastruct.Complaint))
	if err != nil {
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

func (handle *DBHandler) GetDrawCashParams() (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	params := new(datastruct.DrawCashParams)
	has, err := engine.Where("id=?", datastruct.DefaultId).Get(params)
	if err != nil || !has {
		return nil, datastruct.GetDataFailed
	}
	resp := new(datastruct.WebResponseDrawCashParams)
	resp.MaxDrawCount = params.MaxDrawCount
	resp.MinCharge = params.MinCharge
	resp.MinPoundage = params.MinPoundage
	resp.PoundagePer = params.PoundagePer
	resp.RequireVerify = params.RequireVerify
	return resp, datastruct.NULLError
}
func (handle *DBHandler) UpdateDrawCashParams(body *datastruct.WebResponseDrawCashParams) datastruct.CodeType {
	engine := handle.mysqlEngine
	params := new(datastruct.DrawCashParams)
	params.MaxDrawCount = body.MaxDrawCount
	params.MinCharge = body.MinCharge
	params.MinPoundage = body.MinPoundage
	params.PoundagePer = body.PoundagePer
	params.RequireVerify = body.RequireVerify
	_, err := engine.Where("id = ?", datastruct.DefaultId).Cols("min_charge", "min_poundage", "max_draw_count", "poundage_per", "require_verify").Update(params)
	if err != nil {
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

func (handle *DBHandler) GetDrawCashInfo(Id int) (string, bool, *datastruct.DrawCashInfo, datastruct.CodeType) {
	engine := handle.mysqlEngine
	drawCashInfo := new(datastruct.DrawCashInfo)
	has, err := engine.Where("id = ?", Id).Get(drawCashInfo)
	if err != nil || !has {
		return "", false, nil, datastruct.GetDataFailed
	}
	wx_user := new(datastruct.WXPlatform)
	has, err = engine.Where("user_id = ?", drawCashInfo.UserId).Get(wx_user)
	if err != nil || !has {
		return "", false, nil, datastruct.GetDataFailed
	}
	if drawCashInfo.State != datastruct.DrawCashReview {
		return "", false, nil, datastruct.UpdateDataFailed
	}
	isOnlyApp := handle.IsDrawCashOnApp()
	payeeOpenid := ""
	if isOnlyApp {
		if wx_user.PayOpenidForKFPT == "" {
			return "", false, nil, datastruct.PayeeOnlyInApp
		}
		payeeOpenid = wx_user.PayOpenidForKFPT

	} else {
		payeeOpenid = wx_user.PayOpenidForGZH
	}
	return payeeOpenid, isOnlyApp, drawCashInfo, datastruct.NULLError
}

func (handle *DBHandler) DeleteSharePosters(Id int) datastruct.CodeType {
	engine := handle.mysqlEngine
	rs, err := engine.Where("id=?", Id).Delete(new(datastruct.SharePosters))
	if err != nil || rs <= 0 {
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

func (handle *DBHandler) DeleteAd(Id int) datastruct.CodeType {
	engine := handle.mysqlEngine
	rs, err := engine.Where("id=?", Id).Delete(new(datastruct.AdInfo))
	if err != nil || rs <= 0 {
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

func (handle *DBHandler) GetAllAd(body *datastruct.WebQueryAdBody) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	query := "1=1"
	args := make([]interface{}, 0)
	if body.Location != 3 {
		query += " and location = ?"
		args = append(args, body.Location)
	}
	if body.Platform != 2 {
		query += " and platform = ?"
		args = append(args, body.Platform)
	}
	if body.IsHidden != 2 {
		query += " and is_hidden = ?"
		args = append(args, body.IsHidden)
	}
	adInfo := make([]*datastruct.AdInfo, 0)
	engine.Where(query, args...).Desc("sort_id", "id").Find(&adInfo)
	resp := make([]*datastruct.WebResponseAdInfo, 0, len(adInfo))
	for _, v := range adInfo {
		tmp := new(datastruct.WebResponseAdInfo)
		tmp.Id = v.Id
		tmp.ImgUrl = osstool.CreateOSSURL(v.ImgName)
		tmp.IsHidden = v.IsHidden
		tmp.IsJump = v.IsJump
		tmp.JumpTo = v.JumpTo
		tmp.Location = v.Location
		tmp.Platform = v.Platform
		tmp.SortId = v.SortId
		resp = append(resp, tmp)
	}
	return resp, datastruct.NULLError
}

func (handle *DBHandler) EditAd(body *datastruct.WebResponseAdInfo, bucket *oss.Bucket) datastruct.CodeType {
	engine := handle.mysqlEngine
	adInfo := new(datastruct.AdInfo)
	adInfo.IsHidden = body.IsHidden
	adInfo.IsJump = body.IsJump
	adInfo.JumpTo = body.JumpTo
	adInfo.Location = body.Location
	adInfo.Platform = body.Platform
	adInfo.SortId = body.SortId

	var err error
	var affected int64

	if body.Id > 0 {
		isUpdateImg := false
		if !strings.Contains(body.ImgUrl, datastruct.OSSEndpoint) {
			adInfo.ImgName = body.ImgUrl
			isUpdateImg = true
			var has bool
			tmp := new(datastruct.AdInfo)
			has, err = engine.Where("id=?", body.Id).Get(tmp)
			if err != nil || !has {
				log.Debug("DBHandler->EditAd Get AdInfo err")
				return datastruct.GetDataFailed
			}
			osstool.DeleteFile(bucket, tmp.ImgName)
		}
		if isUpdateImg {
			affected, err = engine.Cols("img_name", "is_jump", "jump_to", "sort_id", "location", "platform", "is_hidden").Where("id=?", body.Id).Update(adInfo)
		} else {
			affected, err = engine.Cols("is_jump", "jump_to", "sort_id", "location", "platform", "is_hidden").Where("id=?", body.Id).Update(adInfo)
		}
	} else {
		adInfo.ImgName = body.ImgUrl
		affected, err = engine.Insert(adInfo)
	}
	if err != nil || affected <= 0 {
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

func (handle *DBHandler) GetGoldCoinGift() (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	gcg := new(datastruct.GoldCoinGift)
	has, err := engine.Where("id=?", datastruct.DefaultId).Get(gcg)
	if err != nil || !has {
		return nil, datastruct.GetDataFailed
	}
	resp := new(datastruct.WebResponseGoldCoinGift)
	resp.AppraisedGoldGift = gcg.AppraisedGoldGift
	resp.DownLoadAppGoldGift = gcg.DownLoadAppGoldGift
	resp.IsEnableRegisterGift = gcg.IsEnableRegisterGift
	resp.RegisterGoldGift = gcg.RegisterGoldGift
	resp.IsDrawCashOnlyApp = gcg.IsDrawCashOnlyApp
	return resp, datastruct.NULLError
}

func (handle *DBHandler) EditGoldCoinGift(body *datastruct.WebResponseGoldCoinGift) datastruct.CodeType {
	engine := handle.mysqlEngine
	gcg := new(datastruct.GoldCoinGift)
	gcg.AppraisedGoldGift = body.AppraisedGoldGift
	gcg.DownLoadAppGoldGift = body.DownLoadAppGoldGift
	gcg.IsEnableRegisterGift = body.IsEnableRegisterGift
	gcg.RegisterGoldGift = body.RegisterGoldGift
	gcg.IsDrawCashOnlyApp = body.IsDrawCashOnlyApp
	_, err := engine.Where("id=?", datastruct.DefaultId).Cols("down_load_app_gold_gift", "appraised_gold_gift", "register_gold_gift", "is_enable_register_gift", "is_draw_cash_only_app").Update(gcg)
	if err != nil {
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

func (handle *DBHandler) EditWebUser(body *datastruct.WebEditPermissionUserBody, token string) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	isUpdate := false
	if body.Id > 0 {
		isUpdate = true
	}
	sql := "select id from web_user where login_name = ?"
	results, err := engine.Query(sql, body.LoginName)
	isUpdateMine := 0 //是否修改自己的数据
	if err != nil {
		log.Debug("EditWebUser Query sql err:%v", err.Error())
		return isUpdateMine, datastruct.UpdateDataFailed
	}
	count := len(results)
	if count > 0 {
		if isUpdate {
			query_id := string(results[0]["id"][:])
			if count >= 2 || tools.StringToInt(query_id) != body.Id {
				return isUpdateMine, datastruct.LoginNameAlreadyExisted
			}
		} else {
			return isUpdateMine, datastruct.LoginNameAlreadyExisted
		}
	}
	session := engine.NewSession()
	defer session.Close()
	session.Begin()

	now_time := time.Now().Unix()
	web_user := new(datastruct.WebUser)
	web_user.LoginName = body.LoginName
	web_user.Name = body.Name
	web_user.UpdatedAt = now_time
	var user_id int
	if !isUpdate {
		web_user.Pwd = body.Pwd
		web_user.CreatedAt = now_time
		web_user.RoleId = datastruct.NormalLevelID
		web_user.Token = tools.UniqueId()
		_, err = session.Insert(web_user)
		user_id = web_user.Id
	} else {
		user_id = body.Id
		isUpdatePwd := false
		if body.Pwd != "" {
			isUpdatePwd = true
			web_user.Pwd = body.Pwd
		}
		if isUpdatePwd {
			_, err = session.Where("id=?", user_id).Cols("name", "login_name", "pwd", "updated_at").Update(web_user)
		} else {
			_, err = session.Where("id=?", user_id).Cols("name", "login_name", "updated_at").Update(web_user)
		}
	}
	if err != nil {
		str := fmt.Sprintf("EditPermissionUser err0:%s", err.Error())
		rollbackError(str, session)
		return isUpdateMine, datastruct.UpdateDataFailed
	}
	permission := new(datastruct.WebPermission)
	_, err = session.Where("user_id=?", user_id).Delete(permission)
	if err != nil {
		str := fmt.Sprintf("EditPermissionUser err1:%s", err.Error())
		rollbackError(str, session)
		return isUpdateMine, datastruct.UpdateDataFailed
	}
	for _, v := range body.PermissionIds {
		permission := new(datastruct.WebPermission)
		permission.UserId = user_id
		permission.SecondaryId = v
		_, err = session.Insert(permission)
		if err != nil {
			str := fmt.Sprintf("EditPermissionUser err2:%s", err.Error())
			rollbackError(str, session)
			return isUpdateMine, datastruct.UpdateDataFailed
		}
	}

	err = session.Commit()
	if err != nil {
		str := fmt.Sprintf("DBHandler->EditPermissionUser Commit :%s", err.Error())
		rollbackError(str, session)
		return isUpdateMine, datastruct.UpdateDataFailed
	}

	webUser := new(datastruct.WebUser)
	engine.Where("token=?", token).Get(webUser)
	if webUser.Id == body.Id {
		isUpdateMine = 1
	}
	return isUpdateMine, datastruct.NULLError
}

func (handle *DBHandler) DeleteWebUser(webUserId int) datastruct.CodeType {
	engine := handle.mysqlEngine
	session := engine.NewSession()
	defer session.Close()
	session.Begin()

	var err error
	web_user := new(datastruct.WebUser)
	_, err = session.Where("id=?", webUserId).Delete(web_user)
	if err != nil {
		str := fmt.Sprintf("DeletePermissionUser err0:%s", err.Error())
		rollbackError(str, session)
		return datastruct.UpdateDataFailed
	}

	permission := new(datastruct.WebPermission)
	_, err = session.Where("user_id=?", webUserId).Delete(permission)
	if err != nil {
		str := fmt.Sprintf("DeletePermissionUser err1:%s", err.Error())
		rollbackError(str, session)
		return datastruct.UpdateDataFailed
	}

	err = session.Commit()
	if err != nil {
		str := fmt.Sprintf("DBHandler->DeletePermissionUser Commit :%s", err.Error())
		rollbackError(str, session)
		return datastruct.UpdateDataFailed
	}

	return datastruct.NULLError
}

func (handle *DBHandler) GetWebUsers() (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	users := make([]datastruct.WebUser, 0)
	engine.Where("role_id <> ?", datastruct.AdminLevelID).Asc("created_at").Find(&users)
	web_users := make([]*datastruct.WebResponseAllWebUser, 0)
	for _, v := range users {
		user := new(datastruct.WebResponseAllWebUser)
		user.Id = v.Id
		user.LoginName = v.LoginName
		user.Name = v.Name
		permission := make([]datastruct.WebPermission, 0)
		engine.Where("user_id=?", v.Id).Asc("secondary_id").Find(&permission)
		permissionIds := make([]int, 0, len(permission))
		for _, v := range permission {
			permissionIds = append(permissionIds, v.SecondaryId)
		}
		user.PermissionIds = permissionIds
		web_users = append(web_users, user)
	}
	return web_users, datastruct.NULLError
}

func (handle *DBHandler) GetAllMenuInfo() (interface{}, datastruct.CodeType) {
	return getAllMenu(handle.mysqlEngine), datastruct.NULLError
}

func getAllMenu(engine *xorm.Engine) []*datastruct.MasterInfo {
	master_menu := make([]*datastruct.MasterMenu, 0, 40)
	permission := make([]*datastruct.MasterInfo, 0)
	engine.Asc("id").Find(&master_menu)
	for _, v := range master_menu {
		m_info := new(datastruct.MasterInfo)
		m_info.MasterId = v.Id
		m_info.Name = v.Name
		secondary := make([]*datastruct.SecondaryInfo, 0)
		secondary_menu := make([]*datastruct.SecondaryMenu, 0, 40)
		engine.Where("master_id=?", m_info.MasterId).Asc("id").Find(&secondary_menu)
		for _, v := range secondary_menu {
			secondaryInfo := new(datastruct.SecondaryInfo)
			secondaryInfo.Name = v.Name
			secondaryInfo.SecondaryId = v.Id
			secondary = append(secondary, secondaryInfo)
		}
		m_info.Secondary = secondary
		permission = append(permission, m_info)
	}
	return permission
}

func (handle *DBHandler) CheckPermission(token string, method string, url string) bool {
	engine := handle.mysqlEngine
	user := new(datastruct.WebUser)
	has, err := engine.Where("token=?", token).Get(user)
	if err != nil || !has {
		return false
	}
	if user.RoleId == datastruct.AdminLevelID {
		return true
	}
	sql := "select count(*) from web_secondary_menu_api wsma join web_permission wp on wsma.secondary_id = wp.secondary_id where wp.user_id = ? and method = ? and url = ?"
	results, err := engine.Query(sql, user.Id, method, url)
	if err != nil {
		return false
	}
	count_str := string(results[0]["count(*)"][:])
	if tools.StringToInt(count_str) <= 0 {
		return false
	}
	return true
}

func (handle *DBHandler) UpdateWebUserPwd(body *datastruct.WebUserPwdBody, token string) datastruct.CodeType {
	engine := handle.mysqlEngine
	user := new(datastruct.WebUser)
	user.Pwd = body.NewPwd
	user.Token = tools.UniqueId()
	rs, err := engine.Where("token = ? and pwd = ?", token, body.OldPwd).Cols("pwd", "token").Update(user)
	if err != nil {
		log.Error("UpdateWebUserPwd err:%v", err.Error())
		return datastruct.UpdateDataFailed
	}
	if rs <= 0 {
		log.Error("UpdateWebUserPwd affected row:%v", rs)
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

// var valuesSlice = make([]interface{}, len(cols))
// has, err := engine.Where("id = ?", id).Cols(cols...).Get(&valuesSlice)
