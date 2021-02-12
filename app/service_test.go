package app

import (
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"testing"
)

func initService() (IRecommenderService, *miniredis.Miniredis) {
	MRedis, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	var repo IRepository
	repo, _ = NewRedisRepository("redis://"+MRedis.Addr()+"/0")
	svc, err := NewRecommenderService(repo)
	if err != nil {
		panic(err)
	}
	return svc, MRedis
}

func TestRecommenderServiceImpl_Rate(t *testing.T) {
	svc, conn := initService()
	defer conn.Close()
	err := svc.Rate("1", "1", 4.5)
	err = svc.Rate("2", "1", 3.5)
	err = svc.Rate("3", "1", 4.2)
	err = svc.Rate("1", "2", 4.2)

	user1Items, err := svc.GetUserItems("1", 100)
	user2Items, err := svc.GetUserItems("2", 100)
	if err != nil {
		panic(err)
	}
	if len(user1Items) != 3 {
		t.Error(fmt.Sprintf("User1 Expected To have 3 items, got %v", len(user1Items)))
	}

	if len(user2Items) != 1 {
		t.Error(fmt.Sprintf("User2 Expected To have 1 items, got %v", len(user2Items)))
	}
}

func TestRecommenderServiceImpl_GetUserItems(t *testing.T) {
	TestRecommenderServiceImpl_Rate(t)
}

func TestRecommenderServiceImpl_BatchUpdate(t *testing.T) {
	svc, conn := initService()
	defer conn.Close()
	err := svc.Rate("1", "1", 4.2)
	err = svc.Rate("2", "1", 4.2)
	err = svc.Rate("3", "1", 4.2)
	err = svc.Rate("1", "2", 3.2)
	err = svc.Rate("2", "2", 3.5)
	err = svc.Rate("3", "2", 4.2)
	err = svc.Rate("1", "3", 3.2)
	err = svc.Rate("2", "3", 3.5)
	err = svc.Rate("3", "3", 4.2)
	err = svc.Rate("1", "4", 3.2)
	err = svc.Rate("2", "4", 3.5)
	err = svc.Rate("7", "4", 4.2)
	err = svc.Rate("1", "5", 3.2)
	err = svc.Rate("2", "5", 3.5)
	err = svc.Rate("5", "5", 4.2)
	err = svc.Rate("1", "6", 3.2)
	err = svc.Rate("2", "6", 3.5)
	err = svc.Rate("6", "6", 4.2)

	if err != nil {
		panic(err)
	}

	err = svc.BatchUpdate(MaxNumber)

	if err != nil {
		t.Error("Batch update failed, reason: ", err.Error())
		return
	}

	user1Items, err := svc.GetRecommendedItems("1", 100)

	if err != nil {
		t.Errorf("Get recommended items failed: %v", err.Error())
		return
	}

	if len(user1Items) != 3 {
		t.Errorf("User 1 expected to have 3 recommended items got: %v", len(user1Items))
	}
}

func TestRecommenderServiceImpl_GetRecommendedItems(t *testing.T) {
	TestRecommenderServiceImpl_BatchUpdate(t)
}

func TestRecommenderServiceImpl_GetProbability(t *testing.T) {
	svc, conn := initService()
	defer conn.Close()
	err := svc.Rate("1", "1", 4.2)
	err = svc.Rate("2", "1", 4.2)
	err = svc.Rate("3", "1", 4.2)
	err = svc.Rate("1", "2", 3.2)
	err = svc.Rate("2", "2", 3.5)
	err = svc.Rate("3", "2", 4.2)
	err = svc.Rate("1", "3", 3.2)
	err = svc.Rate("2", "3", 3.5)
	err = svc.Rate("3", "3", 4.2)
	err = svc.Rate("1", "4", 3.2)
	err = svc.Rate("2", "4", 3.5)
	err = svc.Rate("7", "4", 4.2)
	err = svc.Rate("1", "5", 3.2)
	err = svc.Rate("2", "5", 3.5)
	err = svc.Rate("5", "5", 4.2)
	err = svc.Rate("1", "6", 3.2)
	err = svc.Rate("2", "6", 3.5)
	err = svc.Rate("6", "6", 4.2)

	if err != nil {
		panic(err)
	}

	err = svc.BatchUpdate(9223372036854775806)

	if err != nil {
		t.Error("Batch update failed, reason: ", err.Error())
		return
	}

	ZeroProbability, ZeroProbabilityErr := svc.GetProbability("1", "100")
	SomeProbability, SomeProbabilityErr := svc.GetProbability("1", "6")

	if ZeroProbabilityErr != nil {
		t.Errorf("Failed to calculate probability, reason: %v", ZeroProbabilityErr.Error())
	}

	if SomeProbabilityErr != nil {
		t.Errorf("Failed to calculate probability, reason: %v", SomeProbabilityErr.Error())
	}

	if ZeroProbability != 0 {
		t.Errorf("User 1 expected to have 0 probability on item 100 but got: %v", ZeroProbability)
	}

	if SomeProbability == 0 {
		t.Errorf("User 1 expected to have some probability on item 6 but got: %v", SomeProbability)
	}
}


