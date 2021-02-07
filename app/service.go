package app

import (
	"errors"
	"fmt"
	"github.com/RedisLabs/redis-recommend/redrec"
	"github.com/garyburd/redigo/redis"
)

var ErrInvalidArgument = errors.New("invalid/missing argument")

type IRecommenderService interface {
	Rate(item string, user string, score float64) error
	GetRecommendedItems(user string, count int) ([]string, error)
	GetUserItems(user string, max int)([]string, error)
	BatchUpdate(max int) error
	UpdateSuggestedItems(user string, max int) error
	GetProbability(user string, item string) (float64, error)
	//SuggestedItems(user string, max int)([]string, error)
	//SimilarItems(item string, max int)([]string, error)
	//TopItems(user string, max int)([]string, error)
}

type RecommenderService struct {
	conn redis.Conn
	recommender *redrec.Redrec
}

func (svc RecommenderService) UpdateSuggestedItems(user string, max int) error {
	return svc.recommender.UpdateSuggestedItems(user, max);
}

func NewRecommenderService(repo IRepository) IRecommenderService {
	recommender, _ := redrec.New(repo.ConnUrl())
	return RecommenderService{
		conn: *repo.Conn(),
		recommender: recommender,
	}
}

func (svc RecommenderService) GetProbability(user string, item string) (float64, error){
	return svc.recommender.CalcItemProbability(item, user)
}

func (svc RecommenderService) Rate(item string, user string, score float64) error  {
	if user == "" || item == "" || score == 0 {
		return ErrInvalidArgument
	}
	return svc.recommender.Rate(item, user, score)
}

func (svc RecommenderService) GetRecommendedItems(user string, count int)([]string, error)  {
	if user == "" || count == 0 {
		return nil, ErrInvalidArgument
	}

	svc.UpdateSuggestedItems(user, count);

	response, err := svc.recommender.GetUserSuggestions(user, count)
	if err != nil {
		return nil, err
	}
	var  items []string
	for i := 0; i < len(response); i += 2 {
		items = append(items, response[i])
	}
	return items, nil
}

func (svc RecommenderService) GetUserItems(user string, max int) ([]string, error)  {
	if user == "" || max == 0 {
		return nil, ErrInvalidArgument
	}
	items, err := redis.Strings(svc.conn.Do("ZREVRANGE", fmt.Sprintf("user:%s:items", user), 0, max))
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (svc RecommenderService) BatchUpdate(max int) error {
	if max == 0 {
		return ErrInvalidArgument
	}
	return svc.recommender.BatchUpdateSimilarUsers(max)
}