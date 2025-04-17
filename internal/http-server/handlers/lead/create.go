package lead

import (
	"fmt"
	"lead-bitrix/entities"
	"lead-bitrix/internal/http-server/handlers"
	"lead-bitrix/internal/storage/pgx"
	"lead-bitrix/internal/telegram"
	"log/slog"
	"net/http"
)

//todo Log-request-info, middleware
//todo https://chat.deepseek.com/a/chat/s/bfa70afb-8ba3-492a-836b-9341a62f0d50 check this
//todo check delete method in 0.0.0.0 interface

func NewLead(log *slog.Logger, bot *telegram.Bot, storage *pgx.Storage) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		const op = "Handlers.create.NewLead"
		log = log.With(
			slog.String("op", op),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		)
		var lead entities.LeadBitrix

		if r.Method != "POST" {
			handlers.RespondError(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		err := handlers.DecodeAndValidate(w, r, &lead)
		if err != nil {
			log.Error("Failed to decode and validate request", err)
			handlers.RespondError(w, "Invalid request data", http.StatusBadRequest)
			return
		}
		if err = storage.SaveLead(r.Context(), lead); err != nil {
			log.Error("Failed to save lead", err)
			msg := err.Error()
			handlers.RespondError(w, msg, 500)
			return
		}

		leadInfo := fmt.Sprintf("Lead name: %s\nLead Phone: %s\n"+
			"Lead Email: %s\nLead Source: %s\n", lead.Name, lead.Phone, lead.Email, lead.Source)

		err = bot.SendNotification(leadInfo)
		if err != nil {
			log.Error("Failed to send notification", err)
		}

		log.Info("lead created successfully")
		handlers.RespondJSON(w, http.StatusCreated, lead)

	}

}
