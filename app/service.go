package app

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

var ErrInvalidArgument = errors.New("invalid/missing argument")

type RecommenderService interface {
	New(repo Repository) RecommenderService
	Rate(item string, user string, score float64) error
	GetRecommendedItems(user string, count int) ([]string, error)
	GetUserItems(user string, max int)([]string, error)
	BatchUpdate(max int) error
	UpdateSuggestedItems(user string, max int) error
	GetProbability(user string, item string) (float64, error)
	//SuggestedItems(user string, max int)([]string, error)
	//SimilarItems(item string, max int)([]string, error)
	//TopItems(user string, max int)([]string, error)
	//AddNewUser(user string) error
}

type RecommenderServiceImpl struct {
	repo Repository
}

func (svc RecommenderServiceImpl) UpdateSuggestedItems(user string, max int) error {
	return svc.repo.Recommender().UpdateSuggestedItems(user, max);
}

func (RecommenderServiceImpl) New(repo Repository) RecommenderService {
	return RecommenderServiceImpl{
		repo: repo,
	}
}

func (svc RecommenderServiceImpl) GetProbability(user string, item string) (float64, error){
	return svc.repo.Recommender().CalcItemProbability(item, user)
}

func (svc RecommenderServiceImpl) Rate(item string, user string, score float64) error  {
	if user == "" || item == "" || score == 0 {
		return ErrInvalidArgument
	}
	return svc.repo.Recommender().Rate(item, user, score)
}

func (svc RecommenderServiceImpl) GetRecommendedItems(user string, count int)([]string, error)  {
	if user == "" || count == 0 {
		return nil, ErrInvalidArgument
	}

	svc.UpdateSuggestedItems(user, count);

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
	if user == "" || max == 0 {
		return nil, ErrInvalidArgument
	}
	items, err := redis.Strings(svc.repo.Conn().Do("ZREVRANGE", fmt.Sprintf("user:%s:items", user), 0, max))
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (svc RecommenderServiceImpl) BatchUpdate(max int) error {
	if max == 0 {
		return ErrInvalidArgument
	}
	return svc.repo.Recommender().BatchUpdateSimilarUsers(max)
}