package xredis

import (
	"errors"
	"fmt"
	"time"

	jsoniter "github.com/json-iterator/go"

	"github.com/garyburd/redigo/redis"
)

var (
	json                = jsoniter.ConfigCompatibleWithStandardLibrary
	addr                = "redis://*****:6379"
	redisMaxIdle        = 3   //最大空闲连接数
	redisIdleTimeoutSec = 240 //最大空闲连接时间
	password            = "*****"
)

type Config struct {
	Addr           string `json:"Addr"`
	Password       string `json:"Password"`
	Database       int    `json:"Database"`
	MaxActive      int    `json:"MaxActive"`
	MaxIdle        int    `json:"MaxIdle"`
	IdleTimeoutSec int    `json:"IdleTimeoutSec"`
}

type IRedis struct {
	redisPool *redis.Pool
}

//NewIRedis new
func NewIRedis(config Config) *IRedis {
	return &IRedis{
		redisPool: NewRedisPool(config),
	}
}

// NewRedisPool 返回redis连接池
func NewRedisPool(config Config) *redis.Pool {
	idleTimeoutSec := time.Duration(redisIdleTimeoutSec)
	if config.IdleTimeoutSec > 0 {
		idleTimeoutSec = time.Duration(config.IdleTimeoutSec) * time.Second
	}
	return &redis.Pool{
		MaxIdle:     config.MaxIdle,
		MaxActive:   config.MaxActive,
		IdleTimeout: time.Duration(idleTimeoutSec),
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(config.Addr)
			if err != nil {
				return nil, fmt.Errorf("redis connection error: %s", err)
			}
			//验证redis密码
			if config.Password != "" {
				if _, authErr := c.Do("AUTH", config.Password); authErr != nil {
					return nil, fmt.Errorf("redis auth password error: %s", authErr)
				}
			}
			c.Do("SELECT", config.Database)
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				return fmt.Errorf("ping redis error: %s", err)
			}
			return nil
		},
	}
}

//Close end redis conn
func (i *IRedis) Close() error {
	if i.redisPool != nil {
		return i.redisPool.Close()
	}
	return nil
}

//GetConn get redis conn
func (i *IRedis) GetConn() redis.Conn {
	if i.redisPool == nil {
		return nil
	}
	return i.redisPool.Get()
}

// RdbSet redis  SET
func (i *IRedis) RdbSet(key, v string) (bool, error) {
	rds := i.GetConn()
	if rds == nil {
		return false, fmt.Errorf("conn error")
	}
	b, err := redis.Bool(rds.Do("SET", key, v))
	if err != nil {
		return false, err
	}
	return b, nil
}

// RdbGet redis  GET
func (i *IRedis) RdbGet(key string) (string, error) {
	rds := i.GetConn()
	if rds == nil {
		return "", fmt.Errorf("conn error")
	}
	val, err := redis.String(rds.Do("GET", key))
	if err != nil {
		return "", err
	}

	return val, nil
}

// RdbSetKeyExp redis key EXPIRE
func (i *IRedis) RdbSetKeyExp(key string, ex int) error {
	rds := i.GetConn()
	if rds == nil {
		return fmt.Errorf("conn error")
	}
	_, err := rds.Do("EXPIRE", key, ex)
	if err != nil {
		return err
	}
	return nil
}

// RdbSetExp redis EXPIRE
func (i *IRedis) RdbSetExp(key, v string, ex int) error {
	rds := i.GetConn()
	if rds == nil {
		return fmt.Errorf("conn error")
	}
	_, err := rds.Do("SET", key, v, "EX", ex)
	if err != nil {
		return err
	}
	return nil
}

//RdbExists redis EXISTS
func (i *IRedis) RdbExists(key string) bool {
	rds := i.GetConn()
	if rds == nil {
		return false
	}
	b, err := redis.Bool(rds.Do("EXISTS", key))
	if err != nil {
		return false
	}
	return b
}

// RdbDel edis DEL
func (i *IRedis) RdbDel(key string) error {
	rds := i.GetConn()
	if rds == nil {
		return fmt.Errorf("conn error")
	}
	_, err := rds.Do("DEL", key)
	if err != nil {
		return err
	}
	return nil
}

// RdbSetJson redis SETNX
func (i *IRedis) RdbSetJson(key string, data interface{}) error {
	rds := i.GetConn()
	if rds == nil {
		return fmt.Errorf("conn error")
	}
	value, _ := json.Marshal(data)
	n, _ := rds.Do("SETNX", key, value)
	if n != int64(1) {
		return errors.New("set failed")
	}
	return nil
}

// RdbGetJson redis GET return map
func (i *IRedis) RdbGetJson(key string) (map[string]string, error) {
	rds := i.GetConn()
	if rds == nil {
		return nil, fmt.Errorf("conn error")
	}
	var jsonData map[string]string
	bv, err := redis.Bytes(rds.Do("GET", key))
	if err != nil {
		return nil, err
	}
	errJson := json.Unmarshal(bv, &jsonData)
	if errJson != nil {
		return nil, err
	}
	return jsonData, nil
}

//RdbHSet redis hSet 注意 设置什么类型 取的时候需要获取对应类型
func (i *IRedis) RdbHSet(key string, field string, data interface{}) error {
	rds := i.GetConn()
	if rds == nil {
		return fmt.Errorf("conn error")
	}
	_, err := rds.Do("HSET", key, field, data)
	if err != nil {
		return err
	}
	return nil
}

//RdbHGet redis hGet 注意 设置什么类型 取的时候需要获取对应类型
func (i *IRedis) RdbHGet(key, field string) (interface{}, error) {
	rds := i.GetConn()
	if rds == nil {
		return false, fmt.Errorf("conn error")
	}
	data, err := rds.Do("HGET", key, field)
	if err != nil {
		return nil, err
	}
	return data, nil
}

//RdbHGetJson redis hGet json
func (i *IRedis) RdbHGetJson(key, field string) (map[string]interface{}, error) {
	rds := i.GetConn()
	if rds == nil {
		return nil, fmt.Errorf("conn error")
	}
	var jsonData map[string]interface{}
	bv, err := redis.Bytes(rds.Do("GET", key))
	if err != nil {
		return nil, err
	}
	errJson := json.Unmarshal(bv, &jsonData)
	if errJson != nil {
		return nil, err
	}
	return jsonData, nil
}

//RdbHGetAll redis hGetAll return map
func (i *IRedis) RdbHGetAll(key string) (map[string]string, error) {
	rds := i.GetConn()
	if rds == nil {
		return nil, fmt.Errorf("conn error")
	}
	data, err2 := redis.StringMap(rds.Do("HGETALL", key))
	_, err := data, err2
	if err != nil {
		return nil, err
	}
	return data, nil
}

//RdbIncr redis INCR 将 key 中储存的数字值增一
func (i *IRedis) RdbIncr(key string) error {
	rds := i.GetConn()
	if rds == nil {
		return fmt.Errorf("conn error")
	}
	_, err := rds.Do("INCR", key)
	if err != nil {
		return err
	}
	return nil
}

//RdbIncrBy redis INCRBY 将 key 所储存的值加上增量 n
func (i *IRedis) RdbIncrBy(key string, n int) error {
	rds := i.GetConn()
	if rds == nil {
		return fmt.Errorf("conn error")
	}
	_, err := rds.Do("INCRBY", key, n)
	if err != nil {
		return err
	}
	return nil
}

//RdbDecr redis DECR 将 key 中储存的数字值减一。
func (i *IRedis) RdbDecr(key string) error {
	rds := i.GetConn()
	if rds == nil {
		return fmt.Errorf("conn error")
	}
	_, err := rds.Do("DECR", key)
	if err != nil {
		return err
	}
	return nil
}

//RdbDecrBy redis DECRBY 将 key 所储存的值减去减量 n
func (i *IRedis) RdbDecrBy(key string, n int) error {
	rds := i.GetConn()
	if rds == nil {
		return fmt.Errorf("conn error")
	}
	_, err := rds.Do("DECRBY", key, n)
	if err != nil {
		return err
	}
	return nil
}

//RdbSAdd redis SADD 将一个或多个 member 元素加入到集合 key 当中，已经存在于集合的 member 元素将被忽略。
func (i *IRedis) RdbSAdd(key, v string) error {
	rds := i.GetConn()
	if rds == nil {
		return fmt.Errorf("conn error")
	}
	_, err := rds.Do("SADD", key, v)
	if err != nil {
		return err
	}
	return nil
}

//RdbSMembers redis SMEMBERS 返回集合 key 中的所有成员; return map
func (i *IRedis) RdbSMembers(key string) (interface{}, error) {
	rds := i.GetConn()
	if rds == nil {
		return nil, fmt.Errorf("conn error")
	}
	data, err := redis.Strings(rds.Do("SMEMBERS", key))
	if err != nil {
		return nil, err
	}
	return data, nil
}

//RdbSISMembers redis SISMEMBER 判断 member 元素是否集合 key 的成员。 return bool
func (i *IRedis) RdbSISMembers(key, v string) bool {
	rds := i.GetConn()
	if rds == nil {
		return false
	}
	b, err := redis.Bool(rds.Do("SISMEMBER", key, v))
	if err != nil {
		return false
	}
	return b
}
