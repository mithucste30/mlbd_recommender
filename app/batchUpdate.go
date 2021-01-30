package app

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"os"
)

func BatchUpdate(port int, redisHost string)  {
	var repo Repository
	repo, _ = NewRedis(redisHost)

	var svc RecommenderService
	svc = RecommenderServiceImpl{}.New(repo)
	err := svc.BatchUpdate(-1) // lets batchUpdate all users.
	chekErrorAndExit(repo.Conn(), err)
}

func chekErrorAndExit(conn redis.Conn, err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		conn.Close()
		os.Exit(1)
	}
}
