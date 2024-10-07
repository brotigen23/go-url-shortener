package handlers

import "net/http"

func StatusHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	// намеренно добавлена ошибка в JSON
	rw.Write([]byte(`{"status":"ok}`))
}
