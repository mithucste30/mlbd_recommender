package app

import (
	"github.com/RedisLabs/redis-recommend/redrec"
	"github.com/garyburd/redigo/redis"
)

type Repository interface {
	Conn() redis.Conn
	Recommender() *redrec.Redrec
}
