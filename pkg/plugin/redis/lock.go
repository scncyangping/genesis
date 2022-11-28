package redis

import (
	"github.com/go-redis/redis"
	"time"
)

type RLock struct {
	conn    redis.Cmdable
	timeout time.Duration
	key     string
	val     string
}

func NewRedisLock(conn redis.Cmdable, key, val string, timeout time.Duration) *RLock {
	return &RLock{conn: conn, timeout: timeout, key: key, val: val}
}

func (lock *RLock) TryLock() (bool, error) {
	return lock.conn.SetNX(lock.key, lock.val, lock.timeout).Result()
}

func (lock *RLock) UnLock() error {
	luaDel := redis.NewScript("if redis.call('get',KEYS[1]) == ARGV[1] then " +
		"return redis.call('del',KEYS[1]) else return 0 end")
	return luaDel.Run(lock.conn, []string{lock.key}, lock.val).Err()
}

func (lock *RLock) GetLockKey() string {
	return lock.key
}

func (lock *RLock) GetLockVal() string {
	return lock.val
}
