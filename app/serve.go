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
	var repo Repository
	repo, _ = NewRedis(redisHost)
	var svc RecommenderService
	svc = RecommenderServiceImpl{}.New(repo)
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
