package lead

import (
	"errors"
	"lead-bitrix/entities"
	"lead-bitrix/internal/http-server/handlers"
	"lead-bitrix/internal/storage/pgx"
	"log/slog"
	"net/http"
)

type UpdateLeadRequest struct {
	Phone  string  `json:"phone" validate:"required"`
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

		// 1. Проверка метода
		if r.Method != http.MethodPatch {
			handlers.RespondError(w, "Only PATCH method is allowed", http.StatusMethodNotAllowed)
			return
		}

		// 2. Декодирование запроса
		var req UpdateLeadRequest
		if err := handlers.DecodeAndValidate(w, r, &req); err != nil {
			log.Error("Failed to decode request", "error", err)
			return
		}

		// 3. Валидация телефона
		if req.Phone == "" {
			handlers.RespondError(w, "Phone is required", http.StatusBadRequest)
			return
		}

		// 5. Подготовка данных для обновления
		updateData := entities.LeadBitrix{
			Phone:  req.Phone,
			Name:   *req.Name,
			Email:  *req.Email,
			Source: *req.Source,
		}

		// 6. Обновление в БД
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

		// 7. Успешный ответ
		handlers.RespondJSON(w, http.StatusOK, map[string]interface{}{
			"status": "success",
			"data":   updatedLead,
		})
	}
}
