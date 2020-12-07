package xredis

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

var (
	redisMgr *RedisMgr
)

// RedisMgrConfig struct config for redis mgr
type RedisMgrConfig map[string]Config

// RedisMgr struct for redis mgr
type RedisMgr struct {
	instances map[string]*IRedis
}

// Get get redis interface
func (mgr *RedisMgr) Get(tag string, args ...interface{}) *IRedis {
	if redisMgr == nil {
		return nil
	}
	rdb, ok := mgr.instances[tag]
	if !ok {
		return nil
	}
	return rdb
}

// Close get redis interface
func (mgr *RedisMgr) Close() {
	for _, rds := range mgr.instances {
		rds.Close()
	}
}

// ForEach for every rds
func (mgr *RedisMgr) ForEach(cb func(string, redis.Conn)) {
	for tag, rds := range mgr.instances {
		cb(tag, rds.GetConn())
	}
}

//NewRedisMgr new mgr
func NewRedisMgr(mgrConf RedisMgrConfig) *RedisMgr {
	mgr := &RedisMgr{
		instances: map[string]*IRedis{},
	}

	for tag, conf := range mgrConf {
		mgr.instances[tag] = NewIRedis(conf)
	}

	return mgr
}

// RedisInit init redis mgr
func RedisInit(mgrConf RedisMgrConfig) *RedisMgr {
	redisMgr = NewRedisMgr(mgrConf)
	return redisMgr
}

// RedisFinal final redis mgr
func RedisFinal() {
	if redisMgr != nil {
		redisMgr.Close()
	}
}

// Redis redis by tag
func Redis(tag string) *IRedis {
	return redisMgr.Get(tag)
}

// ForEach for each
func ForEach(cb func(string, redis.Conn)) error {
	if redisMgr != nil {
		redisMgr.ForEach(cb)
		return nil
	}
	return fmt.Errorf("no redis ")
}
