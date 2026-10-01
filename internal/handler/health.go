package handler

import (
	"net/http"

	"github.com/garrychrstn/top-go/internal/util"
)

// Health reports service liveness.
func Health(w http.ResponseWriter, _ *http.Request) {
	util.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
