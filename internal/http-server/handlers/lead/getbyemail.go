package lead

import (
	"github.com/go-chi/chi/v5"
	"lead-bitrix/internal/http-server/handlers"
	"lead-bitrix/internal/storage/pgx"
	"log/slog"
	"net/http"
)

func LeadGetByEmail(log *slog.Logger, storage *pgx.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		const op = "handlers.LeadGetByEmail"

		log = log.With(
			slog.String("op", op),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		)

		email := chi.URLParam(r, "email")
		if email == "" {
			log.Error("no email address provided")
			handlers.RespondError(w, "No email address provided", http.StatusBadRequest)
			return
		}

		lead, err := storage.GetLead(r.Context(), email)
		if err != nil {
			log.Error(err.Error())
			handlers.RespondError(w, "Error getting lead", http.StatusNotFound)
			return
		}

		log.Info("lead retrieved successfully")
		handlers.RespondJSON(w, http.StatusOK, lead)
	}
}
