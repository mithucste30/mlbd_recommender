package app

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
	"net/http"
)

func makeRateEndpoint(svc RecommenderService) endpoint.Endpoint  {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(rateRequest)
		err = svc.Rate(req.User, req.Item, req.Score)
		if err != nil {
			return rateResponse{Error: err}, nil
		}
		return rateResponse{ Message: "successful" }, nil
	}
}

func makeSuggestedItemsEndpoint(svc RecommenderService) endpoint.Endpoint  {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(suggestedItemsRequest)

		suggestedItems, err := svc.GetRecommendedItems(req.User, -1)
		if err != nil {
			return nil, err
		}
		results :=  make([]string, len(suggestedItems))
		start := req.Page * req.Per
		end := start + req.Per
		for i := start; i < end; i++ {
			results = append(results, suggestedItems[i])
		}
		return suggestedItemsResponse{Items: results}, nil
	}
}

type rateRequest struct {
	User string `json:"user"`
	Item string `json:"item"`
	Score float64 `json:"score"`
}
type rateResponse struct {
	Error error `json:"err,omitempty"`
	Message string `json:"message"`
}

type suggestedItemsRequest struct {
	User string `json:"user"`
	Page int `json:"page"`
	Per  int `json:"per"`
}

type suggestedItemsResponse struct {
	Error error `json:"err,omitempty"`
	Items []string `json:"items"`
}

func decodeRateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request rateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeSuggestedItemsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errBadRoute
	}
	return suggestedItemsRequest{User: id}, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
