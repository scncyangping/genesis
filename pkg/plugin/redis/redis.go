// @Author: YangPing
// @Create: 2023/10/23
// @Description: redis插件配置

package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"time"
)

var Client *RClient

type Config interface {
	GetPoolSize() int
	GetAddr() []string
	GetPwd() string
	GetDialTimeout() time.Duration
	GetReadTimeout() time.Duration
	GetWriteTimeout() time.Duration
}
type RClient struct {
	redis.Cmdable
}

func NewRedisConn(o Config) (client *RClient, err error) {
	var redisCli redis.Cmdable
	if len(o.GetAddr()) > 1 {
		redisCli = redis.NewClusterClient(
			&redis.ClusterOptions{
				Addrs:        o.GetAddr(),
				PoolSize:     o.GetPoolSize(),
				DialTimeout:  o.GetDialTimeout(),
				ReadTimeout:  o.GetReadTimeout(),
				WriteTimeout: o.GetWriteTimeout(),
				Password:     o.GetPwd(),
			},
		)
	} else {
		redisCli = redis.NewClient(
			&redis.Options{
				Addr:         o.GetAddr()[0],
				DialTimeout:  o.GetDialTimeout(),
				ReadTimeout:  o.GetReadTimeout(),
				WriteTimeout: o.GetWriteTimeout(),
				Password:     o.GetPwd(),
				PoolSize:     o.GetPoolSize(),
				DB:           0,
			},
		)
	}

	err = redisCli.Ping(context.Background()).Err()
	if nil != err {
		return nil, errors.Wrapf(err, "Redis Init Error: Host: %v, Error:%v ", o.GetAddr, err)
	}

	client = new(RClient)
	client.Cmdable = redisCli
	Client = client
	return client, nil
}

func (c *RClient) Process(ctx context.Context, cmd redis.Cmder) error {
	switch redisCli := c.Cmdable.(type) {
	case *redis.ClusterClient:
		return redisCli.Process(ctx, cmd)
	case *redis.Client:
		return redisCli.Process(ctx, cmd)
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
