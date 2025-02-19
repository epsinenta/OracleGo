package main

import (
	"OracleGo/internal/db"
	"OracleGo/internal/handlers"
	"OracleGo/internal/log"
	"OracleGo/internal/middlewares"
	"OracleGo/internal/redis"
	"OracleGo/internal/repository"
	"OracleGo/internal/services"

	"net/http"
)

func main() {
	logger := log.New()

	database, err := db.NewDatabaseManager()
	if err != nil {
		logger.LogFatal(err)
	}
	defer func() {
		if err := database.GratefulStop(); err != nil {
			logger.Log(err)
		}
	}()

	cache, err := redis.NewRedisManager()
	if err != nil {
		logger.LogFatal(err)
	}
	defer func() {
		if err := cache.GratefulStop(); err != nil {
			logger.Log(err)
		}
	}()

	repo := repository.NewRepositoryManager(database, cache)

	svc, err := services.NewServiceManager(repo, cache)
	if err != nil {
		logger.LogFatal(err)
	}

	mm := middlewares.NewMiddlewareManager(svc, logger)

	hm, err := handlers.NewHandlerManager(svc, logger)
	if err != nil {
		logger.LogFatal(err)
	}

	defaultChain := mm.BuildChain(mm.Recoverer, mm.Logger, mm.AuthenticationChecker)
	chainWithRegCheck := mm.BuildChain(mm.Recoverer, mm.Logger, mm.AuthenticationChecker, mm.RegistrationChecker)
	// chainWithAuth := middlewares.BuildChain(logger, defaultChain, middlewares.Authenticator)

	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.Handle("/", defaultChain(http.HandlerFunc(hm.HomeHandler)))

	http.Handle("/login", chainWithRegCheck(http.HandlerFunc(hm.LoginHandler)))
	http.Handle("/register", chainWithRegCheck(http.HandlerFunc(hm.RegisterHandler)))
	http.Handle("/profile", defaultChain(http.HandlerFunc(hm.ProfileHandler)))

	http.Handle("/prediction", defaultChain(http.HandlerFunc(hm.PredictionHandler)))
	http.Handle("/recommendations", defaultChain(http.HandlerFunc(hm.RecommendationsHandler)))
	http.Handle("/statistics", defaultChain(http.HandlerFunc(hm.StatisticsHandler)))
	http.Handle("/team-analysis", defaultChain(http.HandlerFunc(hm.TeamsHandler)))

	http.Handle("/profile-complete", chainWithRegCheck(http.HandlerFunc(hm.ProfileCompleteHandler)))
	http.Handle("/confirm", defaultChain(http.HandlerFunc(hm.ConfirmationHandler)))

	logger.Info().Msg(("Server started on :8080"))
	logger.Fatal().Err(http.ListenAndServe(":8080", nil)).Msg("")
}
