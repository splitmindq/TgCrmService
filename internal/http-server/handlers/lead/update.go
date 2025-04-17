package lead

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"lead-bitrix/entities"
	"lead-bitrix/internal/http-server/handlers"
	"lead-bitrix/internal/storage/pgx"
	"log/slog"
	"net/http"
)

type UpdateLeadRequest struct {
	Name   *string `json:"name,omitempty"`
	Email  *string `json:"email,omitempty"`
	Source *string `json:"source,omitempty"`
}

func UpdateLead(log *slog.Logger, storage *pgx.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.lead.UpdateLead"
		log = log.With(
			slog.String("op", op),
			slog.String("method", r.Method),
		)

		if r.Method != http.MethodPatch {
			handlers.RespondError(w, "Only PATCH method is allowed", http.StatusMethodNotAllowed)
			return
		}

		phone := chi.URLParam(r, "phone")
		if phone == "" {
			handlers.RespondError(w, "phone is required", http.StatusBadRequest)
			return
		}

		var req UpdateLeadRequest
		if err := handlers.DecodeAndValidate(w, r, &req); err != nil {
			log.Error("Failed to decode request", "error", err)
			return
		}

		if req.Name == nil && req.Email == nil && req.Source == nil {
			handlers.RespondError(w, "At least one field must be provided to update", http.StatusBadRequest)
			return
		}

		updateData := entities.LeadBitrix{
			Phone:  phone,
			Name:   *req.Name,
			Email:  *req.Email,
			Source: *req.Source,
		}

		updatedLead, err := storage.UpdateLead(r.Context(), updateData)
		if err != nil {
			log.Error("Failed to update lead", "error", err)

			switch {
			case errors.Is(err, pgx.ErrNotFound):
				handlers.RespondError(w, "Lead not found", http.StatusNotFound)
			case errors.Is(err, pgx.ErrEmailExists):
				handlers.RespondError(w, "Email already exists", http.StatusConflict)
			default:
				handlers.RespondError(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		handlers.RespondJSON(w, http.StatusOK, map[string]interface{}{
			"status": "success",
			"data":   updatedLead,
		})
	}
}
