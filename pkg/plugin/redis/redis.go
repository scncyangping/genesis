// @Author: YangPing
// @Create: 2023/10/23
// @Description: redis插件配置

package redis

import (
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"time"
)

var Client *RClient

type Config interface {
	PoolSize() int
	Addr() []string
	Pwd() string
	DialTimeout() time.Duration
	ReadTimeout() time.Duration
	WriteTimeout() time.Duration
}
type RClient struct {
	redis.Cmdable
}

func NewRedisConn(o Config) (client *RClient, err error) {
	var redisCli redis.Cmdable
	if len(o.Addr()) > 1 {
		redisCli = redis.NewClusterClient(
			&redis.ClusterOptions{
				Addrs:        o.Addr(),
				PoolSize:     o.PoolSize(),
				DialTimeout:  o.DialTimeout(),
				ReadTimeout:  o.ReadTimeout(),
				WriteTimeout: o.WriteTimeout(),
				Password:     o.Pwd(),
			},
		)
	} else {
		redisCli = redis.NewClient(
			&redis.Options{
				Addr:         o.Addr()[0],
				DialTimeout:  o.DialTimeout(),
				ReadTimeout:  o.ReadTimeout(),
				WriteTimeout: o.WriteTimeout(),
				Password:     o.Pwd(),
				PoolSize:     o.PoolSize(),
				DB:           0,
			},
		)
	}
	err = redisCli.Ping().Err()
	if nil != err {
		return nil, errors.Wrapf(err, "Redis Init Error: Host: %v, Error:%v ", o.Addr, err)
	}

	client = new(RClient)
	client.Cmdable = redisCli
	Client = client
	return client, nil
}

func (c *RClient) Process(cmd redis.Cmder) error {
	switch redisCli := c.Cmdable.(type) {
	case *redis.ClusterClient:
		return redisCli.Process(cmd)
	case *redis.Client:
		return redisCli.Process(cmd)
	default:
		return nil
	}
}

func (c *RClient) Close() error {
	switch redisCli := c.Cmdable.(type) {
	case *redis.ClusterClient:
		return redisCli.Close()
	case *redis.Client:
		return redisCli.Close()
	default:
		return nil
	}
}
