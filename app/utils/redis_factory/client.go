package redis_factory

import (
	"github.com/gomodule/redigo/redis"
	"goskeleton/app/core/event"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/global/variable"
	"goskeleton/app/utils/config"
	"time"
)

func createRedisClientPool() *redis.Pool {

	configFac := config.CreateYamlFactory()
	redisPool := &redis.Pool{
		MaxIdle:     configFac.GetInt("Redis.MaxIdle"),                        //最大空闲数
		MaxActive:   configFac.GetInt("Redis.MaxActive"),                      //最大活跃数
		IdleTimeout: configFac.GetDuration("Redis.IdleTimeout") * time.Second, //最大的空闲连接等待时间，超过此时间后，空闲连接将被关闭
		Dial: func() (redis.Conn, error) {
			//此处对应redis ip及端口号
			conn, err := redis.Dial("tcp", configFac.GetString("Redis.Host")+":"+configFac.GetString("Redis.Port"))
			if err != nil {
				variable.ZapLog.Error(my_errors.ErrorsRedisInitConnFail + err.Error())
			}
			auth := configFac.GetString("Redis.Auth") //通过配置项设置redis密码
			if len(auth) >= 1 {
				if _, err := conn.Do("AUTH", auth); err != nil {
					conn.Close()
					variable.ZapLog.Error(my_errors.ErrorsRedisAuhtFail + err.Error())
				}
			}
			conn.Do("select", configFac.GetInt("Redis.IndexDb"))
			return conn, err
		},
	}
	// 将redis的关闭事件，注册在全局事件统一管理器，由程序退出时统一销毁
	event.CreateEventManageFactory().Set(variable.EventDestroyPrefix+"Redis", func(args ...interface{}) {
		redisPool.Close()
	})
	return redisPool
}

//  从连接池获取一个redis连接
func GetOneRedisClient() *RedisClient {
	poolConn := createRedisClientPool()
	//defer clientConn.Close()
	return &RedisClient{poolConn.Get()}
}

// 定义一个redis客户端结构体
type RedisClient struct {
	client redis.Conn
}

// 为redis-go 客户端封装统一操作函数入口
func (r *RedisClient) Execute(cmd string, args ...interface{}) (interface{}, error) {
	return r.client.Do(cmd, args...)
}

// 释放连接池
func (r *RedisClient) RelaseOneRedisClientPool() {
	r.client.Close()
}

//  封装几个数据类型转换的函数

//bool 类型转换
func (r *RedisClient) Bool(reply interface{}, err error) (bool, error) {
	return redis.Bool(reply, err)
}

//string 类型转换
func (r *RedisClient) String(reply interface{}, err error) (string, error) {
	return redis.String(reply, err)
}

//strings 类型转换
func (r *RedisClient) Strings(reply interface{}, err error) ([]string, error) {
	return redis.Strings(reply, err)
}

//Float64 类型转换
func (r *RedisClient) Float64(reply interface{}, err error) (float64, error) {
	return redis.Float64(reply, err)
}

//int 类型转换
func (r *RedisClient) Int(reply interface{}, err error) (int, error) {
	return redis.Int(reply, err)
}

//int64 类型转换
func (r *RedisClient) Int64(reply interface{}, err error) (int64, error) {
	return redis.Int64(reply, err)
}

//uint64 类型转换
func (r *RedisClient) Uint64(reply interface{}, err error) (uint64, error) {
	return redis.Uint64(reply, err)
}

//Bytes 类型转换
func (r *RedisClient) Bytes(reply interface{}, err error) ([]byte, error) {
	return redis.Bytes(reply, err)
}
