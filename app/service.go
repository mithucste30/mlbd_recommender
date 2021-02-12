package app

import (
	"errors"
	"fmt"
	"github.com/RedisLabs/redis-recommend/redrec"
	"github.com/garyburd/redigo/redis"
	"math"
	"strconv"
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

type ServiceMiddleware func(IRecommenderService) IRecommenderService

type RecommenderService struct {
	conn redis.Conn
	recommender *redrec.Redrec
}

func (svc RecommenderService) UpdateSuggestedItems(user string, max int) error {
	return svc.updateSuggestedItems(user, max);
}

func NewRecommenderService(repo IRepository) (IRecommenderService, error) {
	recommender, err := redrec.New(repo.ConnUrl())
	if err != nil {
		return nil, err
	}
	return RecommenderService{
		conn: *repo.Conn(),
		recommender: recommender,
	}, nil
}

func (svc RecommenderService) GetProbability(user string, item string) (float64, error){
	return svc.calcItemProbability(item, user)
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

	err := svc.BatchUpdate(MaxNumber)
	if err != nil {
		return nil, err
	}

	err = svc.UpdateSuggestedItems(user, count);
	if err != nil {
		return nil, err
	}

	response, err := svc.getUserSuggestions(user, count)
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
	return svc.batchUpdateSimilarUsers(max)
}

func (svc RecommenderService) getUserSuggestions(user string, max int) ([]string, error) {
	items, err := redis.Strings(svc.conn.Do("ZREVRANGE", fmt.Sprintf("user:%s:suggestions", user), 0, max, "WITHSCORES"))
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (svc RecommenderService) batchUpdateSimilarUsers(max int) error  {
	users, err := redis.Strings(svc.conn.Do("SMEMBERS", "users"))
	if err != nil {
		return err
	}
	for _, user := range users {
		candidates, err := svc.getSimilarityCandidates(user, max)
		args := []interface{}{}
		args = append(args, fmt.Sprintf("user:%s:similars", user))
		for _, candidate := range candidates {
			if candidate != user {
				score, _ := svc.calcSimilarity(user, candidate)
				args = append(args, score)
				args = append(args, candidate)
			}
		}
		if len(args) > 1 {
			_, err = svc.conn.Do("ZADD", args...)
			if err != nil {
				return errors.New(fmt.Sprintf("ZADD ERR0: %v", err))
			}
		}
	}
	return nil
}

func (svc RecommenderService) updateSuggestedItems(user string, max int) error {
	items, err := svc.getSuggestCandidates(user, max)
	if max > len(items) {
		max = len(items)
	}

	args := []interface{}{}
	args = append(args, fmt.Sprintf("user:%s:suggestions", user))
	for _, item := range items {
		probability, _ := svc.calcItemProbability(user, item)
		args = append(args, probability)
		args = append(args, item)
	}
	if len(args) > 1 {
		_, err = svc.conn.Do("ZADD", args...)
		if err != nil {
			return errors.New(fmt.Sprintf("ZADD ERR1: %v", err))
		}
	}
	return nil
}

func (svc RecommenderService) getSimilarityCandidates(user string, max int) ([]string, error)  {
	items, err := svc.GetUserItems(user, max)
	if max > len(items) {
		max = len(items)
	}

	args := []interface{}{}
	args = append(args, "ztmp", float64(max))
	for i := 0; i < max; i++ {
		args = append(args, fmt.Sprintf("item:%s:scores", items[i]))
	}

	_, err = svc.conn.Do("ZUNIONSTORE", args...)
	if err != nil {
		return nil, err
	}

	users, err := redis.Strings(svc.conn.Do("ZRANGE", "ztmp", 0, -1))
	if err != nil {
		return nil, err
	}

	_, err = svc.conn.Do("DEL", "ztmp")
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (svc RecommenderService) calcSimilarity(user string, simuser string) (float64, error) {
	_, err := svc.conn.Do("ZINTERSTORE",
		"ztmp", 2, fmt.Sprintf("user:%s:items", user), fmt.Sprintf("user:%s:items", simuser), "WEIGHTS", 1, -1)
	if err != nil {
		return 0, err
	}

	userDiffs, err := redis.Strings(svc.conn.Do("ZRANGE", "ztmp", 0, -1, "WITHSCORES"))
	svc.conn.Do("DEL", "ztmp")
	if err != nil {
		return 0, err
	}

	if len(userDiffs) == 0 {
		return 0, nil
	}

	var score float64
	for i := 1; i < len(userDiffs); i += 2 {
		diffVal, _ := strconv.ParseFloat(userDiffs[i], 64)
		score += diffVal * diffVal
	}
	score /= float64(len(userDiffs) / 2)
	score = math.Sqrt(score)

	return score, nil
}

func (svc RecommenderService) getSuggestCandidates(user string, max int) ([]string, error) {
	similarUsers, err := redis.Strings(svc.conn.Do("ZRANGE", fmt.Sprintf("user:%s:similars", user), 0, max))
	if err != nil {
		return nil, err
	}

	max = len(similarUsers)
	args := []interface{}{}
	args = append(args, "ztmp", float64(max+1), fmt.Sprintf("user:%s:items", user))
	weights := []interface{}{}
	weights = append(weights, "WEIGHTS", -1.0)
	for _, simuser := range similarUsers {
		args = append(args, fmt.Sprintf("user:%s:items", simuser))
		weights = append(weights, 1.0)
	}

	args = append(args, weights...)
	args = append(args, "AGGREGATE", "MIN")
	_, err = svc.conn.Do("ZUNIONSTORE", args...)
	if err != nil {
		return nil, err
	}

	candidates, err := redis.Strings(svc.conn.Do("ZRANGEBYSCORE", "ztmp", 0, "inf"))
	if err != nil {
		return nil, err
	}

	_, err = svc.conn.Do("DEL", "ztmp")
	if err != nil {
		return nil, err
	}

	return candidates, nil
}

func (svc RecommenderService) calcItemProbability(user string, item string) (float64, error) {
	_, err := svc.conn.Do("ZINTERSTORE",
		"ztmp", 2, fmt.Sprintf("user:%s:similars", user), fmt.Sprintf("item:%s:scores", item), "WEIGHTS", 0, 1)
	if err != nil {
		return 0, err
	}

	scores, err := redis.Strings(svc.conn.Do("ZRANGE", "ztmp", 0, -1, "WITHSCORES"))
	svc.conn.Do("DEL", "ztmp")
	if err != nil {
		return 0, err
	}

	if len(scores) == 0 {
		return 0, nil
	}

	var score float64
	for i := 1; i < len(scores); i += 2 {
		val, _ := strconv.ParseFloat(scores[i], 64)
		score += val
	}
	score /= float64(len(scores) / 2)

	return score, nil
}