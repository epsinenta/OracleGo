package main

import (
	"OracleGo/internal/auth"
	"OracleGo/internal/home"
	"OracleGo/internal/ml"
	"OracleGo/internal/net"
	"OracleGo/internal/profile"
	"OracleGo/internal/statistics"
	_ "fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", home.HomeHandler)

	http.Handle("/login", net.RedirectIfAuthenticated(http.HandlerFunc(auth.LoginHandler)))
	http.Handle("/register", net.RedirectIfAuthenticated(http.HandlerFunc(auth.RegisterHandler)))
	http.Handle("/profile", net.SessionMiddleware(http.HandlerFunc(profile.ProfileHandler)))

	http.HandleFunc("/prediction", ml.PredictionHandler)
	http.HandleFunc("/recommendations", ml.RecommendationsHandler)
	http.HandleFunc("/statistics", statistics.StatisticsHandler)
	http.HandleFunc("/team-analysis", statistics.TeamsHandler)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
