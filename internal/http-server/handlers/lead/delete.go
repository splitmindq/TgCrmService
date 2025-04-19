package lead

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"lead-bitrix/internal/http-server/bitrix"
	"lead-bitrix/internal/http-server/handlers"
	"lead-bitrix/internal/storage/pgx"
	"lead-bitrix/internal/telegram"
	"log/slog"
	"net/http"
)

//todo Log-request-info, middleware
//todo https://chat.deepseek.com/a/chat/s/bfa70afb-8ba3-492a-836b-9341a62f0d50 check this
//todo check delete method in 0.0.0.0 interface

func DelLead(log *slog.Logger, bot *telegram.Bot, storage *pgx.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "Handlers.DelLead"

		log = log.With(
			slog.String("op", op),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		)

		if r.Method != "DELETE" {
			log.Error(op, "Invalid method")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		email := chi.URLParam(r, "email")
		if email == "" {
			handlers.RespondError(w, "Email is required", 400)
			return
		}

		lead, err := storage.GetLead(r.Context(), email)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				handlers.RespondError(w, "Lead not found", http.StatusNotFound)
			} else {
				handlers.RespondError(w, "Failed to get lead", http.StatusInternalServerError)
			}
			return
		}

		err = bitrix.DeleteLead(log, lead)
		if err != nil {
			handlers.RespondError(w, "Failed to delete lead", http.StatusInternalServerError)
			return
		}

		err = storage.DeleteLead(r.Context(), email)
		if err != nil {
			handlers.RespondError(w, "Failed to delete lead", http.StatusInternalServerError)
			return
		}

		leadInfo := fmt.Sprintf("\nLead name: %s\nLead Phone: %s\n"+
			"Lead Email: %s\nLead Source: %s\n", lead.Form.Name, lead.Form.Phone, lead.Form.Email, lead.Form.Source)

		err = bot.SendNotification(fmt.Sprintf("Deleted lead: %s", leadInfo))
		if err != nil {
			log.Error("Failed to send notification", err)
		}

		log.Info("lead deleted successfully")
		w.WriteHeader(http.StatusNoContent)

	}

}
