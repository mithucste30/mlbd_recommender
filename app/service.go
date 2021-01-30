package app

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)

type RecommenderService interface {
	New(repo Repository) RecommenderService
	Rate(user string, item string, score float64) error
	GetRecommendedItems(user string, count int) ([]string, error)
	GetUserItems(user string, max int)([]string, error)
	BatchUpdate(max int) error
	//SuggestedItems(user string, max int)([]string, error)
	//SimilarItems(item string, max int)([]string, error)
	//TopItems(user string, max int)([]string, error)
	//AddNewUser(user string) error
}

type RecommenderServiceImpl struct {
	repo Repository
}

func (RecommenderServiceImpl) New(repo Repository) RecommenderService {
	return RecommenderServiceImpl{
		repo: repo,
	}
}

func (svc RecommenderServiceImpl) Rate(user string, item string, score float64) error  {
	return svc.repo.Recommender().Rate(user, item, score)
}

func (svc RecommenderServiceImpl) GetRecommendedItems(user string, count int)([]string, error)  {
	response, err := svc.repo.Recommender().GetUserSuggestions(user, count)
	if err != nil {
		return nil, err
	}
	var  items []string
	for i := 0; i < len(response); i += 2 {
		items = append(items, response[i])
	}
	return items, nil
}

func (svc RecommenderServiceImpl) GetUserItems(user string, max int) ([]string, error)  {
	items, err := redis.Strings(svc.repo.Conn().Do("ZREVRANGE", fmt.Sprintf("user:%s:items", user), 0, max))
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (svc RecommenderServiceImpl) BatchUpdate(max int) error {
	return svc.repo.Recommender().BatchUpdateSimilarUsers(max)
}