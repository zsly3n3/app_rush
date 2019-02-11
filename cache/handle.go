package cache

import (
	"app/datastruct"
	"app/log"
	"app/tools"

	"github.com/gomodule/redigo/redis"
)

func (handle *CACHEHandler) clearData() {
	conn := handle.GetConn()
	defer conn.Close()
	conn.Do("flushdb")
}

func (handle *CACHEHandler) IsExistUserWithConn(conn redis.Conn, key string) (int, bool, bool) {
	value, err := redis.Values(conn.Do("hmget", key, datastruct.UserIdField, datastruct.IsBlackListField))
	length := len(value)
	if err != nil || length == 0 {
		return -1, false, false
	}
	var userId int
	var isBlackList int
	for i := 0; i < length; i++ {
		x_value := value[i]
		if x_value == nil {
			return -1, false, false
		}
		tmp := x_value.([]byte)
		str := string(tmp[:])
		switch i {
		case 0:
			userId = tools.StringToInt(str)
		case 1:
			isBlackList = tools.StringToInt(str)
		}
	}
	log.Debug("CACHEHandler IsExistUserWithConn userid:%s", userId)
	return userId, true, tools.IntToBool(isBlackList)
}

func (handle *CACHEHandler) GetConn() redis.Conn {
	conn := handle.redisClient.Get()
	return conn
}

func (handle *CACHEHandler) AddExpire(conn redis.Conn, token string) {
	_, err := conn.Do("expire", token, datastruct.EXPIRETIME)
	if err != nil {
		log.Debug("expire token:%s err:%s", token, err.Error())
		return
	}
}

//创建充值订单
func (handle *CACHEHandler) CreateOrderForm(pay *datastruct.OrderForm) datastruct.CodeType {
	conn := handle.GetConn()
	defer conn.Close()
	value, isError := tools.InterfaceToString(pay)
	if isError {
		log.Error("CreateOrderForm InterfaceToString err")
		return datastruct.UpdateDataFailed
	}
	conn.Send("MULTI")
	conn.Send("SETNX", pay.Id, value)
	conn.Send("expire", pay.Id, datastruct.WXOrderMaxSec)
	_, err := conn.Do("EXEC")
	if err != nil {
		log.Error("CreateOrderForm payInfo:%v err:%s", pay, err.Error())
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

//获取充值订单
func (handle *CACHEHandler) GetPayData(payId string) (*datastruct.OrderForm, datastruct.CodeType) {
	conn := handle.GetConn()
	defer conn.Close()
	value, err := redis.String(conn.Do("get", payId))
	if err != nil || value == "" {
		log.Error("GetPayData payId:%v err", payId)
		return nil, datastruct.GetDataFailed
	}
	pay, isError := tools.BytesToOrderForm([]byte(value))
	if isError {
		log.Error("GetPayData BytesToOrderForm err")
		return nil, datastruct.GetDataFailed
	}
	return pay, datastruct.NULLError
}

/*
//创建付费游戏状态便于检测
func (handle *CACHEHandler) CreatePayGameState(state *datastruct.GameState) datastruct.CodeType {
	conn := handle.GetConn()
	defer conn.Close()
	value, isError := tools.InterfaceToString(state)
	if isError {
		return datastruct.UpdateDataFailed
	}
	key := fmt.Sprintf("%s%d", datastruct.GameStateKey, state.UserId)
	conn.Send("MULTI")
	conn.Send("SET", key, value)
	_, err := conn.Do("EXEC")
	if err != nil {
		log.Debug("CreatePayGameState state:%v err:%s", state, err.Error())
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

//获取付费游戏状态
func (handle *CACHEHandler) GetPayGameState(userId int) (*datastruct.GameState, datastruct.CodeType) {
	conn := handle.GetConn()
	defer conn.Close()
	key := fmt.Sprintf("%s%d", datastruct.GameStateKey, userId)
	value, err := redis.String(conn.Do("get", key))
	if err != nil || value == "" {
		return nil, datastruct.GetDataFailed
	}
	state, isError := tools.BytesToGameState([]byte(value))
	if isError {
		return nil, datastruct.GetDataFailed
	}
	return state, datastruct.NULLError
}*/

func (handle *CACHEHandler) SetUserAllData(conn redis.Conn, u_data *datastruct.UserInfo) {
	key := u_data.Token
	//add

	args := make([]interface{}, 0, 5)
	args = append(args, key)

	args = append(args, datastruct.UserIdField)
	args = append(args, u_data.Id)
	args = append(args, datastruct.IsBlackListField)
	args = append(args, u_data.IsBlackList)

	_, err := conn.Do("hmset", args...)

	if err != nil {
		log.Debug("CACHEHandler SetPlayerData err:%s", err.Error())
	}
}

func (handle *CACHEHandler) UpdateBlackList(token string, isBlackList int) {
	conn := handle.GetConn()
	defer conn.Close()

	args := make([]interface{}, 0, 3)
	args = append(args, token)

	args = append(args, datastruct.IsBlackListField)
	args = append(args, isBlackList)

	_, err := conn.Do("hset", args...)

	if err != nil {
		log.Debug("CACHEHandler UpdateBlackList err:%s", err.Error())
	}
}

func (handle *CACHEHandler) DeletedKeys(keys []interface{}) {
	conn := handle.GetConn()
	defer conn.Close()
	_, err := conn.Do("del", keys...)
	if err != nil {
		log.Debug("CACHEHandler DeletedKeys err:%s", err.Error())
	}
}

func (handle *CACHEHandler) GetTmpAvatars() []string {

	return nil
}

func (handle *CACHEHandler) GetTmpUsers() []*datastruct.TmpUser {

	return nil
}

// func (handle *CACHEHandler) IsExistVipOrder(orderId string) bool {
// 	isExist := false
// 	conn := handle.GetConn()
// 	defer conn.Close()
// 	value, err := redis.String(conn.Do("get", orderId))
// 	if err != nil {
// 		return isExist
// 	}
// 	if value != "" {
// 		isExist = true
// 	}
// 	return isExist
// }
