package app

import (
	"errors"
	"fmt"
	"github.com/go-kit/kit/log"
	kitHttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"time"
)

var errBadRoute = errors.New("bad route")

func Serve(port int, redisHost string, doc bool) {
	logger := log.NewLogfmtLogger(os.Stdout)
	logger = log.With(logger, "listen", port, "caller", log.DefaultCaller)
	repo, err := NewRedisRepository(redisHost)
	if err != nil {
		log.With(logger, "Failed to initialize repository, reason:", err.Error())
		return
	}
	svc, err := NewRecommenderService(repo)
	if err != nil {
		log.With(logger, "Failed to initialize service, reason: ", err.Error())
		return
	}
	svc = loggingMiddleware(logger)(svc)
	rateHandler := kitHttp.NewServer(makeRateEndpoint(svc), decodeRateRequest, encodeResponse)
	suggestedItemsHandler := kitHttp.NewServer(makeSuggestedItemsEndpoint(svc), decodeSuggestedItemsRequest, encodeResponse)
	userItemsHandler := kitHttp.NewServer(makeUserItemsEndpoint(svc), decodeUserItemsRequest, encodeResponse)
	r := mux.NewRouter()
	r.Handle("/api/rate", rateHandler).Methods("POST")
	r.Handle("/api/users/{user}/suggestions", suggestedItemsHandler).Queries("count", "{count:[0-9]+}").Methods("GET")
	r.Handle("/api/users/{user}/suggestions", suggestedItemsHandler).Methods("GET")
	r.Handle("/api/users/{user}/items", userItemsHandler).Methods("GET")
	if doc == true {
		r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("doc/dist"))))
	}
	http := &http.Server{
		Handler: r,
		Addr: fmt.Sprintf("0.0.0.0:%v", port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Printf("Starting server. Listening at 0.0.0.0:%v\n", port)
	log.With(logger, "error", http.ListenAndServe())
}
