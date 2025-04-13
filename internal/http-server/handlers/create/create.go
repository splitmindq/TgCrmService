package create

import (
	"fmt"
	"lead-bitrix/entities"
	"lead-bitrix/internal/http-server/handlers"
	"log/slog"
	"net/http"
)

func NewLead(log *slog.Logger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.NewLead"

		log = log.With(
			slog.String("op:", op),
		)

		if r.Method != "POST" {
			log.Error("Invalid method:", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var lead entities.LeadBitrix
		err := handlers.DecodeAndValidate(w, r, &lead, log)
		if err != nil {
			log.Error("Ошибка разбора JSON: "+err.Error(), log)
		}

		fmt.Println(lead)
		
	}

}
