package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func respondError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}

func DecodeAndValidate(w http.ResponseWriter, r *http.Request, v interface{}, log *slog.Logger) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		respondError(w, "Ошибка разбора JSON: "+err.Error(), http.StatusBadRequest)
		return err
	}

	return nil
}
