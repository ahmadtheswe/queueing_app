package routers

import (
	"net/http"
	"time"

	"github.com/ahmadtheswe/queueing_app/common/utils/http_utils"
)

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"version":     "1.0.0",
		"server_time": time.Now().Format(time.RFC3339),
		"server":      "backend",
	}

	http_utils.SendSuccessResponse(w, data)
}

func SetupRoutes() {
	http.Handle("/health", http_utils.NewMethodHandler().Get(healthCheckHandler))
}
