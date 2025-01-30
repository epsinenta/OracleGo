package handlers

import (
	"OracleGo/internal/net"
	_ "fmt"
	"net/http"
)

func (hm *HandlersManager) RecommendationsHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}
	net.RenderTemplate(w, r, "prediction.html", data)
}
