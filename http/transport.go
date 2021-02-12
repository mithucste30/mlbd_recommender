package http

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
	"github.com/mithucste30/mlbd_recommender/app"
	"net/http"
	"strconv"
)

func makeRateEndpoint(svc app.IRecommenderService) endpoint.Endpoint  {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(rateRequest)
		err = svc.Rate(req.Item, req.User, req.Score)
		if err != nil {
			return err, nil
		}
		return rateResponse{ Message: "successful" }, nil
	}
}

func makeSuggestedItemsEndpoint(svc app.IRecommenderService) endpoint.Endpoint  {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(suggestedItemsRequest)
		count := app.MaxNumber
		if req.Count != 0 {
			count = req.Count
		}
		suggestedItems, err := svc.GetRecommendedItems(req.User, count)
		if err != nil {
			return err, nil
		}
		return suggestedItemsResponse{Items: suggestedItems}, nil
	}
}

func makeUserItemsEndpoint(svc app.IRecommenderService) endpoint.Endpoint  {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(userItemsRequest)
		fetchCount := app.MaxNumber
		if req.Count != 0 {
			fetchCount = req.Count
		}
		userItems, err := svc.GetUserItems(req.User, fetchCount)
		if err != nil {
			return err, nil
		}
		return userItemsResponse{Items: userItems}, nil
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
	User string
	Count int
}

type userItemsRequest struct {
	User string
	Count int
}

type suggestedItemsResponse struct {
	Error error `json:"err,omitempty"`
	Items []string `json:"items"`
}

type userItemsResponse struct {
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
	user, ok := vars["user"]
	if !ok {
		return nil, errBadRoute
	}
	if countParam, ok := vars["count"]; ok {
		if count, err := strconv.Atoi(countParam); err == nil {
			return suggestedItemsRequest{User: user, Count: count}, nil
		}
	}
	return suggestedItemsRequest{User: user}, nil
}

func decodeUserItemsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	user, ok := vars["user"]
	if !ok {
		return nil, errBadRoute
	}
	if countParam, ok := vars["count"]; ok {
		if count, err := strconv.Atoi(countParam); err == nil {
			return userItemsRequest{User: user, Count: count}, nil
		}
	}
	return userItemsRequest{User: user}, nil
}


func encodeError(_ context.Context, err error, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case app.ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	encodeErr := json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
	if encodeErr != nil{
		return encodeErr
	}
	return nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(error); ok {
		if err := encodeError(ctx, e, w); err != nil {
			return err
		}
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}