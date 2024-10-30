package home

import (
	"OracleGo/internal/net"
	_ "fmt"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}

	net.RenderTemplate(w, r, "home.html", data)
}
