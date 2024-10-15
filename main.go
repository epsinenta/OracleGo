package main

import (
	"OracleGo/db"
	"fmt"
	_ "fmt"
	"html/template"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/prediction", PredictionHandler)
	http.HandleFunc("/statistics", StatisticsHandler)
	http.HandleFunc("/profile", ProfileHandler)
	http.HandleFunc("/teams", TeamAnalysisHandler)
	http.HandleFunc("/recommendations", RecommendationsHandler)

	// Статические файлы
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "home.html")
}

type Hero struct {
	Name string
}
type Team struct {
	Name string
}

type PageData struct {
	Hero []Hero
	Team []Team
}

func PredictionHandler(w http.ResponseWriter, r *http.Request) {
	dbManager, err := db.NewDatabaseManager()
	if err != nil {
		log.Fatalf("Не удалось создать DataBaseManager: %v", err)
	}
	heroesRows, err := dbManager.GetRows("heroes_list", []string{"hero_name"}, map[string][]string{"patch": {"7.35c"}})
	if err != nil {
		log.Fatalf("Не удалось провести запрос %v", err)
	}
	var heroes []Hero
	for _, row := range heroesRows {
		heroes = append(heroes, Hero{row[0]})
	}
	teamsRows, err := dbManager.GetRows("teams_roasters", []string{"team_name"}, map[string][]string{})
	if err != nil {
		log.Fatalf("Не удалось провести запрос %v", err)
	}
	fmt.Print(teamsRows)
	var teams []Team
	for _, row := range teamsRows {
		teams = append(teams, Team{row[0]})
	}
	parsedTemplate, err := template.ParseFiles("templates/prediction.html")
	if err != nil {
		log.Fatalf("Ошибка при парсинге шаблона: %v", err)
	}
	parsedTemplate.Execute(w, PageData{Hero: heroes, Team: teams})
}

type HeroPerPatch struct {
	Patch   string
	Name    string
	Winrate string
}

func StatisticsHandler(w http.ResponseWriter, r *http.Request) {
	dbManager, err := db.NewDatabaseManager()
	if err != nil {
		log.Fatalf("Не удалось создать DataBaseManager: %v", err)
	}
	resultRows, err := dbManager.GetRows("heroes_list", []string{"*"}, map[string][]string{})
	if err != nil {
		log.Fatalf("Не удалось провести запрос %v", err)
	}
	var heroes []HeroPerPatch
	for _, row := range resultRows {
		heroes = append(heroes, HeroPerPatch{row[0], row[1], row[2]})
	}
	parsedTemplate, err := template.ParseFiles("templates/statistics.html")
	if err != nil {
		log.Fatalf("Ошибка при парсинге шаблона: %v", err)
	}
	parsedTemplate.Execute(w, struct{ HeroPerPatch []HeroPerPatch }{heroes})
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "profile.html")
}

func TeamAnalysisHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "teams.html")
}

func RecommendationsHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "recommendations.html")
}

func renderTemplate(w http.ResponseWriter, tmpl string) {
	parsedTemplate, _ := template.ParseFiles("templates/" + tmpl)
	parsedTemplate.Execute(w, nil)
}
