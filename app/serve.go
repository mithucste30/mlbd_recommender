package app

import (
	"errors"
	"flag"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

var errBadRoute = errors.New("bad route")

func Serve(port int, redisHost string)  {
	////var repo Repository
	////repo, _ = NewRedis(redisHost)
	//
	////var svc RecommenderService
	////svc = RecommenderServiceImpl{}.New(repo)
	//
	//r := mux.NewRouter()
	////_, err := http.Dir("app/swaggerui/dist").Open("index.html")
	////if err != nil {
	////	fmt.Println("error opening file")
	////}
	//ui := http.FileServer(http.Dir("app/swaggerui/dist"))
	//r.Handle("/", ui)
	////rateHandler := httpTransport.NewServer(makeRateEndpoint(svc), decodeRateRequest, encodeResponse)
	////r.Handle("/api/rate", rateHandler).Methods("POST")
	////suggestedItemsHandler := httpTransport.NewServer(makeSuggestedItemsEndpoint(svc), decodeSuggestedItemsRequest, encodeResponse)
	////r.Handle("/api/users/{id}/suggestions", suggestedItemsHandler).Methods("GET")
	////serverMux := http.NewServeMux()
	////serverMux.Handle("/", r)
	//http.Handle("/", ui)
	//serverPort := fmt.Sprintf(":%v", port)
	//log.Fatal(http.ListenAndServe(serverPort, nil))

	var dir string

	flag.StringVar(&dir, "dir", "/usr/src/app/app/swaggerui/dist", "the directory to serve files from. Defaults to the current dir")
	flag.Parse()
	r := mux.NewRouter()

	// This will serve files under http://localhost:8000/static/<filename>
	r.PathPrefix("/docs").Handler(http.StripPrefix("/docs", http.FileServer(http.Dir(dir))))

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

