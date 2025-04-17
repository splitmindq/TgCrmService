package handlers

import (
	"encoding/json"
	"lead-bitrix/internal/validator"
	"log/slog"
	"net/http"
)

var validate = validator.Validate

func RespondError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}

func RespondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("Failed to encode JSON response", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func DecodeAndValidate(w http.ResponseWriter, r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		RespondError(w, "Ошибка разбора JSON: "+err.Error(), http.StatusBadRequest)
		return err
	}

	if err := validate.Struct(v); err != nil {
		RespondError(w, "Ошибка валидации: "+err.Error(), http.StatusUnprocessableEntity)
		return err
	}

	return nil
}

func DecodeLead(w http.ResponseWriter, r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		RespondError(w, "Ошибка разбора JSON: "+err.Error(), http.StatusBadRequest)
		return err
	}
	return nil
}
