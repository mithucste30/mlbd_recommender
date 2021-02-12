package app

import "github.com/go-kit/kit/log"

func loggingMiddleware(logger log.Logger) ServiceMiddleware {
	return func(next IRecommenderService) IRecommenderService {
		return logmw{logger, next}
	}
}

type logmw struct {
	logger log.Logger
	IRecommenderService
}

func (mw logmw) Rate(item string, user string, score float64) error  {
	mw.logger.Log("method", "rate", "user", user, "item", item, "score", score);
	return mw.IRecommenderService.Rate(item, user, score)
}

func (mw logmw) GetRecommendedItems(user string, count int) ([]string, error)  {
	mw.logger.Log("method", "GetRecommendedItems", "user", user, "count", count);
	return mw.IRecommenderService.GetRecommendedItems(user, count);
}

func (mw logmw) GetUserItems(user string, count int) ([]string, error)  {
	mw.logger.Log("method", "GetRecommendedItems", "user", user, "count", count);
	return mw.IRecommenderService.GetUserItems(user, count);
}


