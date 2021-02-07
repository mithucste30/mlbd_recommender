package app

import (
	"errors"
	kitHttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

var errBadRoute = errors.New("bad route")

func Serve(port int, redisHost string) {
	var err error
	var repo IRepository
	repo, err = NewRedisRepository(redisHost)
	var svc IRecommenderService
	svc = NewRecommenderService(repo)
	if err != nil {
		log.Fatalf("Failed to initialize service, reason: %v", err.Error())
		return
	}
	rateHandler := kitHttp.NewServer(makeRateEndpoint(svc), decodeRateRequest, encodeResponse)
	suggestedItemsHandler := kitHttp.NewServer(makeSuggestedItemsEndpoint(svc), decodeSuggestedItemsRequest, encodeResponse)
	r := mux.NewRouter()
	r.Handle("/api/rate", rateHandler).Methods("POST")
	r.Handle("/api/users/{id}/suggestions", suggestedItemsHandler).Methods("GET")
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("app/doc/dist"))))

	srv := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
