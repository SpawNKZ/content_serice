package main

import (
	"context"
	"flag"
	"fmt"
	contentRepo "github.com/SpawNKZ/content_service/content/repo"
	contentService "github.com/SpawNKZ/content_service/content/service"
	content "github.com/SpawNKZ/content_service/content/transports"
	contentHistoryRepo "github.com/SpawNKZ/content_service/content_history/repo"
	contentHistoryService "github.com/SpawNKZ/content_service/content_history/service"
	contentStatusRepo "github.com/SpawNKZ/content_service/content_status/repo"
	contentStatusService "github.com/SpawNKZ/content_service/content_status/service"
	contentStatus "github.com/SpawNKZ/content_service/content_status/transports"
	"github.com/SpawNKZ/content_service/db"
	"github.com/SpawNKZ/content_service/mb"
	postRepo "github.com/SpawNKZ/content_service/post/repo"
	postService "github.com/SpawNKZ/content_service/post/service"
	post "github.com/SpawNKZ/content_service/post/transports"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"

	"github.com/SpawNKZ/content_service/utils"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		fmt.Println(err)
	}
	var (
		addr     = config.HTTP_PORT
		httpAddr = flag.String("http.addr", ":"+addr, "HTTP listen address")
	)

	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	dbConn, err := db.NewDBConnection(context.Background(), config.DBSource, logger)
	if err != nil {
		logger.Log("cannot connect to DB", err)
	}

	nc, closeNats, err := mb.NewNatsConnection(config.NATS, logger)
	if err != nil {
		logger.Log("cannot connect to NATS", err)
	}

	defer func() {
		closeNats()
	}()

	var contentStatusSvc contentStatusService.Service
	contentStatusRepository := contentStatusRepo.NewRepository(dbConn)
	contentStatusSvc = contentStatusService.New(contentStatusRepository)
	contentStatusSvc = contentStatusService.NewLoggingService(log.With(logger, "component", "contentStatus"), contentStatusSvc)

	var contentHistorySvc contentHistoryService.Service
	contentHistoryRepository := contentHistoryRepo.NewRepository(dbConn)
	contentHistorySvc = contentHistoryService.New(contentHistoryRepository)
	contentHistorySvc = contentHistoryService.NewLoggingService(log.With(logger, "component", "contentHistory"), contentHistorySvc)

	var contentSvc contentService.Service
	contentRepository := contentRepo.NewRepository(dbConn)
	subjectRepository := contentRepo.NewSubjectRepository(nc)
	microtopicRepository := contentRepo.NewMicrotopicRepository(nc)
	contentSvc = contentService.New(contentRepository, subjectRepository, microtopicRepository, contentHistorySvc, contentStatusSvc)
	contentSvc = contentService.NewLoggingService(log.With(logger, "component", "content"), contentSvc)

	var postSvc postService.Service
	postRepository := postRepo.NewRepository(dbConn)
	postSvc = postService.New(postRepository, contentSvc)
	postSvc = postService.NewLoggingService(log.With(logger, "component", "post"), postSvc)
	//org_svc.Seed(context.Background())

	httpLogger := log.With(logger, "component", "http")

	r := mux.NewRouter()
	r.Use(accessControl)

	contentStatusHandler := contentStatus.MakeHTTPHandler(contentStatusSvc, httpLogger)
	contentHandler := content.MakeHTTPHandler(contentSvc, httpLogger)
	postHandler := post.MakeHTTPHandler(postSvc, httpLogger)

	r.PathPrefix("/api/v1/content_status").Handler(contentStatusHandler)
	r.PathPrefix("/api/v1/content").Handler(contentHandler)
	r.PathPrefix("/api/v1/post").Handler(postHandler)

	http.Handle("/", r)

	errs := make(chan error, 2)
	go func() {
		logger.Log("transport", "http", "address", *httpAddr, "msg", "listening")
		errs <- http.ListenAndServe(*httpAddr, nil)
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("exit", <-errs)
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}
