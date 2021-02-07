package app

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)

type RedisRepository struct {
	conn *redis.Conn
	url string
}

func NewRedisRepository(connUrl string) (IRepository, error)  {
	url := "redis://"+connUrl+"/0"
	conn, err := redis.DialURL(url)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &RedisRepository{
		conn: &conn,
		url: url,
	}, nil
}

func (repo *RedisRepository) Conn() *redis.Conn {
	return repo.conn
}

func (repo *RedisRepository) ConnUrl() string  {
	return repo.url
}
