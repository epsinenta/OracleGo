package ml

import (
	"OracleGo/internal/net"
	_ "fmt"
	"net/http"
)

func RecommendationsHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}
	net.RenderTemplate(w, r, "prediction.html", data)
}
