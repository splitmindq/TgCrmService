package create

import (
	"fmt"
	"lead-bitrix/entities"
	"lead-bitrix/internal/http-server/handlers"
	"lead-bitrix/internal/telegram"
	"log/slog"
	"net/http"
)

func NewLead(log *slog.Logger, bot *telegram.Bot) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		const op = "Handlers.create.NewLead"
		log = log.With("handler", op)

		var lead entities.LeadBitrix

		err := handlers.DecodeAndValidate(w, r, &lead)
		if err != nil {
			log.Error("Failed to decode and validate request", err)
			return
		}

		leadInfo := fmt.Sprintf("Lead name: %s\nLead Phone: %s\n"+
			"Lead Email: %s\nLead Source: %s\n", lead.Name, lead.Phone, lead.Email, lead.Source)

		err = bot.SendNotification(leadInfo)
		if err != nil {
			log.Error("Failed to send notification", err)
			handlers.RespondError(w, "Failed to send notification", 500)
			return
		}

		//todo storage
		//todo bitrix service

	}

}
