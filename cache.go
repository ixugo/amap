package amap

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ixugo/amap/conc"
)

// Cache 缓存接口，用户可以自定义实现（如Redis等）
type Cache interface {
	// Get 获取缓存值
	Get(key string) ([]byte, bool)
	// Set 设置缓存值，带过期时间
	Set(key string, value []byte)
}

// TTLMapCache 基于TTL Map的缓存实现
type TTLMapCache struct {
	ttl    time.Duration
	ttlMap *conc.TTLMap[string, []byte]
}

// NewTTLMapCache 创建新的TTL Map缓存
func NewTTLMapCache(ttl time.Duration) *TTLMapCache {
	return &TTLMapCache{
		ttl:    ttl,
		ttlMap: conc.NewTTLMap[string, []byte](),
	}
}

// Get 获取缓存值
func (c *TTLMapCache) Get(key string) ([]byte, bool) {
	return c.ttlMap.Load(key)
}

// Set 设置缓存值
func (c *TTLMapCache) Set(key string, value []byte) {
	c.ttlMap.Store(key, value, c.ttl)
}

// Delete 删除缓存
func (c *TTLMapCache) Delete(key string) {
	c.ttlMap.Delete(key)
}

// generateCacheKey 生成缓存键
// 为防止碰撞，请使用唯一标识作为 prefix
func generateCacheKey(prefix string, params interface{}) string {
	data, _ := json.Marshal(params)
	hash := md5.Sum(data)
	return fmt.Sprintf("%s:%x", prefix, hash)
}
