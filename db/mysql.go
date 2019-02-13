package db

import (
	"app/conf"
	"app/datastruct"
	"app/log"
	"app/tools"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type DBHandler struct {
	mysqlEngine *xorm.Engine
}

func CreateDBHandler() *DBHandler {
	dbHandler := new(DBHandler)
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4", conf.Server.DB_UserName, conf.Server.DB_Pwd, conf.Server.DB_IP, conf.Server.DB_Name)
	engine, err := xorm.NewEngine("mysql", dsn)
	errhandle(err)
	err = engine.Ping()
	errhandle(err)
	//日志打印SQL
	engine.ShowSQL(true)
	//设置连接池的空闲数大小
	engine.SetMaxIdleConns(20)
	SyncDB(engine)
	initData(engine)
	dbHandler.mysqlEngine = engine
	//go timerTask(engine)
	return dbHandler
}

func SyncDB(engine *xorm.Engine) {
	arr := make([]interface{}, 0, 15)
	var err error

	arr = append(arr, new(datastruct.UserInfo))
	// arr = append(arr, new(datastruct.WXPlatform))
	arr = append(arr, new(datastruct.GoodsClass))
	arr = append(arr, new(datastruct.AdInfo))
	arr = append(arr, new(datastruct.Goods))
	// arr = append(arr, new(datastruct.GoodsRewardPool))
	// arr = append(arr, new(datastruct.FreeModeRougeGame))
	// arr = append(arr, new(datastruct.PayModeRougeGame))
	// arr = append(arr, new(datastruct.SaveGameInfo))
	// arr = append(arr, new(datastruct.InviteInfo))
	// arr = append(arr, new(datastruct.AgencyParams))
	// arr = append(arr, new(datastruct.PayModeRougeGameFailed))
	// arr = append(arr, new(datastruct.PayModeRougeGameSucceed))
	arr = append(arr, new(datastruct.PayModeRougeGameSucceedHistory))
	arr = append(arr, new(datastruct.OrderInfo))
	arr = append(arr, new(datastruct.SendGoods))

	//arr = append(arr, new(datastruct.BalanceInfo))
	arr = append(arr, new(datastruct.DrawCashParams))
	arr = append(arr, new(datastruct.DrawCashInfo))
	// arr = append(arr, new(datastruct.DepositParams))
	//arr = append(arr, new(datastruct.GoldChangeInfo))
	// arr = append(arr, new(datastruct.KfInfo))
	arr = append(arr, new(datastruct.WebUser))
	// arr = append(arr, new(datastruct.Role))
	// arr = append(arr, new(datastruct.EntryAddr))
	// arr = append(arr, new(datastruct.AuthAddr))
	// arr = append(arr, new(datastruct.AppAddr))

	// arr = append(arr, new(datastruct.BlackListJump))
	// arr = append(arr, new(datastruct.MemberLevelData))
	// arr = append(arr, new(datastruct.AppDownloadAddr))
	// arr = append(arr, new(datastruct.AppCustomShare))
	// arr = append(arr, new(datastruct.GZHCustomShare))

	// arr = append(arr, new(datastruct.ServerVersion))
	// arr = append(arr, new(datastruct.CheckInInfo))
	arr = append(arr, new(datastruct.CheckInReward))
	//arr = append(arr, new(datastruct.MemberLevelOrder))
	//arr = append(arr, new(datastruct.UserDepositInfo))

	//arr = append(arr, new(datastruct.RushLimitSetting))
	//arr = append(arr, new(datastruct.RandomLotteryGoods))
	//arr = append(arr, new(datastruct.RandomLotteryGoodsPool))
	//arr = append(arr, new(datastruct.RandomLotteryGoodsSucceed))
	//arr = append(arr, new(datastruct.RandomLotteryGoodsSucceedHistory))
	arr = append(arr, new(datastruct.RecommendedClass))
	arr = append(arr, new(datastruct.SharePosters))
	arr = append(arr, new(datastruct.GoodsImgs))
	arr = append(arr, new(datastruct.UserAppraise))
	arr = append(arr, new(datastruct.UserAppraisePic))
	arr = append(arr, new(datastruct.GoodsDetail))
	arr = append(arr, new(datastruct.UserShippingAddress))

	arr = append(arr, new(datastruct.SuggestionComplaintParams))
	arr = append(arr, new(datastruct.Suggestion))
	arr = append(arr, new(datastruct.Complaint))
	arr = append(arr, new(datastruct.GoldCoinGift))
	arr = append(arr, new(datastruct.TmpData))
	arr = append(arr, new(datastruct.TmpDataForGoods))
	arr = append(arr, new(datastruct.MasterMenu))
	arr = append(arr, new(datastruct.SecondaryMenu))

	// err = engine.DropTables(arr...)
	// errhandle(err)
	// err = engine.CreateTables(arr...)
	// errhandle(err)

	err = engine.Sync2(arr...)
	errhandle(err)

	// role := createRoleData()
	// _, err = engine.Insert(&role)
	// errhandle(err)
	webUser := createLoginData(datastruct.AdminLevelID)
	_, err = engine.Insert(&webUser)
	errhandle(err)
}

func createRoleData() []datastruct.Role {
	admin := datastruct.Role{
		Id:   datastruct.AdminLevelID,
		Desc: "admin",
	}
	guest := datastruct.Role{
		Id:   datastruct.NormalLevelID,
		Desc: "normal",
	}

	return []datastruct.Role{admin, guest}
}

func createLoginData(adminLevelID int) []datastruct.WebUser {
	now_time := time.Now().Unix()
	admin := datastruct.WebUser{
		Name:      "管理员",
		LoginName: "admin",
		Pwd:       "123@s678",
		RoleId:    adminLevelID,
		Token:     tools.UniqueId(),
		CreatedAt: now_time,
		UpdatedAt: now_time,
	}
	return []datastruct.WebUser{admin}
}

func initData(engine *xorm.Engine) {
	execStr := fmt.Sprintf("ALTER DATABASE %s CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci;", conf.Server.DB_Name)
	_, err := engine.Exec(execStr)
	errhandle(err)

	_, err = engine.Exec("ALTER TABLE user_info CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;")
	errhandle(err)

	execStr = fmt.Sprintf("ALTER TABLE user_info AUTO_INCREMENT = %d", datastruct.UserIdStart)
	_, err = engine.Exec(execStr)
	errhandle(err)

	_, err = engine.Exec("ALTER TABLE user_info CHANGE nick_name nick_name VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;")
	errhandle(err)
}

func errhandle(err error) {
	if err != nil {
		log.Fatal("db error is %v", err.Error())
	}
}

func (handle *DBHandler) DepositParams() map[int]float64 {
	engine := handle.mysqlEngine
	params := make([]*datastruct.DepositParams, 0)
	err := engine.Find(&params)
	errhandle(err)
	mp := make(map[int]float64)
	for _, v := range params {
		mp[v.Id] = v.Money
	}
	return mp
}

// /*定时任务*/
// func timerTask(engine *xorm.Engine) {
// 	c := cron.New()
// 	spec := "0 50 14 * * ?"
// 	c.AddFunc(spec, func() {
// 		yesterday, today := tools.GetYesterdayTodayTime()
// 		activity := make([]*datastruct.TodayUserActivityInfo, 0)
// 		err := engine.Where("start_time >= ? and start_time < ?", yesterday, today).Find(&activity)
// 		if err != nil {
// 			log.Debug("timerTask Get TodayUserActivityInfo err:%s", err)
// 			return
// 		}
// 		total := len(activity)
// 		var new_users int64
// 		users := new(datastruct.UserInfo)
// 		query_str := "user_info.created_at >= ? and user_info.created_at < ? and today_user_activity_info.start_time >= ? and today_user_activity_info.start_time < ?"
// 		args := make([]interface{}, 0, 4)
// 		args = append(args, yesterday)
// 		args = append(args, today)
// 		args = append(args, yesterday)
// 		args = append(args, today)
// 		new_users, err = engine.Table("today_user_activity_info").Join("INNER", "user_info", "user_info.id = today_user_activity_info.user_id").Where(query_str, args...).Count(users)
// 		if err != nil {
// 			log.Debug("timerTask Get TodayUserActivityInfo new_users err:%s", err)
// 			return
// 		}
// 		old_users := total - int(new_users)
// 		code := tools.SaveActivityData(activity, old_users, int(new_users), tools.Int64ToString(yesterday))
// 		if code == datastruct.NULLError {
// 			engine.Where("start_time >= ? and start_time < ?", yesterday, today).Delete(activity)
// 		}
// 	})
// 	c.Start()
// }
