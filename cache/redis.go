package cache

import (
	"app/conf"
	"app/log"
	"time"

	"github.com/gomodule/redigo/redis"
)

type CACHEHandler struct {
	redisClient *redis.Pool
}

func getRedisPool() *redis.Pool {
	RedisClient := &redis.Pool{
		// 从配置文件获取maxidle以及maxactive，取不到则用后面的默认值
		MaxIdle:     30,  //最大的空闲连接数，表示即使没有redis连接时依然可以保持N个空闲的连接，而不被清除，随时处于待命状态。
		MaxActive:   100, //最大的激活连接数，表示同时最多有N个连接
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", conf.Server.Redis_IP)
			errhandle(err)
			if _, err := conn.Do("AUTH", conf.Server.Redis_Pwd); err != nil {
				conn.Close()
				errhandle(err)
			}
			conn.Do("SELECT", conf.Server.Redis_Name)
			return conn, nil
		},
	}
	return RedisClient
}

func CreateCACHEHandler() *CACHEHandler {
	cacheHandler := new(CACHEHandler)
	cacheHandler.redisClient = getRedisPool()
	if conf.Common.Mode == conf.Debug {
		cacheHandler.clearData()
	}
	return cacheHandler
}

func errhandle(err error) {
	if err != nil {
		log.Fatal("cache error is %v", err.Error())
	}
}
