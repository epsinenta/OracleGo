package main

import (
	"OracleGo/internal/db"
	"OracleGo/internal/handlers"
	"OracleGo/internal/log"
	"OracleGo/internal/net"
	"OracleGo/internal/redis"
	"OracleGo/internal/repository"
	"OracleGo/internal/services"

	"net/http"
)

func main() {
	logger := log.New()

	db, err := db.NewDatabaseManager()
	if err != nil {
		logger.LogFatal(err)
	}
	defer func() {
		if err := db.GratefulStop(); err != nil {
			logger.Log(err)
		}
	}()

	redis, err := redis.NewRedisManager()
	if err != nil {
		logger.LogFatal(err)
	}
	defer func() {
		if err := redis.GratefulStop(); err != nil {
			logger.Log(err)
		}
	}()

	repo := repository.NewRepositoryManager(db, redis)

	svc := services.NewServiceManager(repo)

	hm := handlers.NewHandlerManager(svc, logger)

	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", hm.HomeHandler)

	http.Handle("/login", net.RedirectIfAuthenticated(http.HandlerFunc(hm.LoginHandler)))
	http.Handle("/register", net.RedirectIfAuthenticated(http.HandlerFunc(hm.RegisterHandler)))
	http.Handle("/profile", net.SessionMiddleware(http.HandlerFunc(hm.ProfileHandler)))

	http.HandleFunc("/prediction", hm.PredictionHandler)
	http.HandleFunc("/recommendations", hm.RecommendationsHandler)
	http.HandleFunc("/statistics", hm.StatisticsHandler)
	http.HandleFunc("/team-analysis", hm.TeamsHandler)

	logger.Info().Msg(("Server started on :8080"))
	logger.Fatal().Err(http.ListenAndServe(":8080", nil)).Msg("")
}
