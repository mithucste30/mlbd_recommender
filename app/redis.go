package app

import (
	"fmt"
	"github.com/RedisLabs/redis-recommend/redrec"
	"github.com/garyburd/redigo/redis"
)

func NewRedis(url string) (Repository, error)  {
	conn, err := redis.DialURL(url)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	redisRec, _ := redrec.New(url)

	rr := &RedisRepository{
		conn: conn,
		redisRec: redisRec,
	}

	return rr, nil
}

type RedisRepository struct {
	conn redis.Conn
	redisRec *redrec.Redrec
}

func (repo *RedisRepository) Conn() (redis.Conn)  {
	return repo.conn
}

func (repo *RedisRepository) Recommender() *redrec.Redrec {
	return repo.redisRec
}
