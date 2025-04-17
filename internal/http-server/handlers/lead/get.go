package lead

import (
	"lead-bitrix/internal/http-server/handlers"
	"lead-bitrix/internal/storage/pgx"
	"log/slog"
	"net/http"
)

func GetLeads(log *slog.Logger, storage *pgx.Storage) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		const op = "Handlers.Get.GetLeads"

		log = log.With(
			slog.String("op", op),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		)

		if r.Method != http.MethodGet {
			log.Warn("invalid method attempt")
			handlers.RespondError(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		leads, err := storage.GetAllLeads(r.Context())
		if err != nil {
			log.Error(op, err)
			handlers.RespondError(w, "Failed to GET leads", http.StatusInternalServerError)
			return
		}

		if len(leads) == 0 {
			log.Info("no leads found")
			handlers.RespondJSON(w, http.StatusOK, map[string]interface{}{
				"message": "No leads found",
				"data":    []string{},
			})
			return
		}

		log.Info("leads retrieved successfully", slog.Int("count", len(leads)))
		handlers.RespondJSON(w, http.StatusOK, map[string]interface{}{
			"count": len(leads),
			"data":  leads,
		})
	}
}
