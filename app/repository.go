package app

import (
	"github.com/garyburd/redigo/redis"
)

type IRepository interface {
	Conn() *redis.Conn
	ConnUrl() string
}
