package event

import (
	"app/cache"
	"app/datastruct"
	"app/db"
	"app/osstool"
	"log"
	"sync"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var commondata *CommonData

type EventHandler struct {
	dbHandler    *db.DBHandler
	cacheHandler *cache.CACHEHandler
	payMutex     *sync.Mutex //读写互斥量
}

func CreateEventHandler() *EventHandler {
	eventHandler := new(EventHandler)
	eventHandler.dbHandler = db.CreateDBHandler()
	eventHandler.cacheHandler = cache.CreateCACHEHandler()
	eventHandler.payMutex = new(sync.Mutex)
	commondata = createCommonData(eventHandler)
	return eventHandler
}

type CommonData struct {
	ServerInfo        *ServerInfo
	DepositParams     map[int]float64 //key为充值id, value为充值价格
	PcRedirect        string
	BlacklistRedirect string
	OSSBucket         *oss.Bucket
	LotteryQueue      *LotteryQueue
	ticker            *time.Ticker
	isExistTicker     bool
	eventHandler      *EventHandler
}

type LotteryQueue struct {
	RWMutex *sync.RWMutex //读写互斥量
}

type ServerInfo struct {
	RWMutex    *sync.RWMutex //读写互斥量
	Version    string        //当前服务端版本号
	IsMaintain int           //0为不维护,1为维护
}

func createServerInfo(dbHandler *db.DBHandler) *ServerInfo {
	server := new(ServerInfo)
	server.RWMutex = new(sync.RWMutex)
	db_data, code := dbHandler.GetServerInfo()
	if code != datastruct.NULLError {
		log.Fatal("GetServerInfo error from db")
		return nil
	}
	server.IsMaintain = db_data.IsMaintain
	server.Version = db_data.Version
	return server
}

func createCommonData(eventHandler *EventHandler) *CommonData {
	commondata := new(CommonData)
	dbHandler := eventHandler.GetDBHandler()
	commondata.eventHandler = eventHandler
	commondata.ServerInfo = createServerInfo(dbHandler)
	commondata.PcRedirect, commondata.BlacklistRedirect = dbHandler.GetRedirect()
	commondata.DepositParams = dbHandler.DepositParams()
	commondata.LotteryQueue = createLotteryQueue()
	commondata.OSSBucket = osstool.CreateOSSBucket()
	commondata.isExistTicker = false
	commondata.createTicker(datastruct.TmpDataTimes)
	dbHandler.TruncateTmpDataForGoods()
	commondata.UpdateTmpData()
	return commondata
}

func createLotteryQueue() *LotteryQueue {
	lotteryQueue := new(LotteryQueue)
	lotteryQueue.RWMutex = new(sync.RWMutex)
	return lotteryQueue
}
