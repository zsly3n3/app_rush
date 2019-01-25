package event

import (
	"app/log"
	"time"
)

func (data *CommonData) createTicker(times time.Duration) {
	if !data.isExistTicker {
		data.isExistTicker = true
		data.ticker = time.NewTicker(times)
		go data.selectTicker()
	}
}

// func (handle *EventHandler)stopTicker(){
//     if handle.ticker != nil{
// 	   handle.ticker.Stop()
// 	   handle.isExistTicker = false
//     }
// }

func (data *CommonData) selectTicker() {
	for {
		select {
		case <-data.ticker.C:
			data.UpdateTmpData()
		}
	}
}

func (data *CommonData) UpdateTmpData() {
	dbHanle := data.eventHandler.GetDBHandler()
	goods_ids := dbHanle.GetGoodsForTmpData()
	if goods_ids == nil || len(goods_ids) == 0 {
		log.Debug("Ticker UpdateTmpData has no goods")
		return
	}
	now_time := time.Now().Unix()
	for _, v := range goods_ids {
		dbHanle.UpdateTmpData(v, now_time)
	}
}
