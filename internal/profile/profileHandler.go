package profile

import (
	"OracleGo/internal/net"
	_ "fmt"
	"net/http"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}

	net.RenderTemplate(w, r, "profile.html", data)
}
