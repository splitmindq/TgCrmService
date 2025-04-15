package pgx

import (
	"context"
	"fmt"
	"lead-bitrix/entities"
)

func (s *Storage) SaveLead(ctx context.Context, lead entities.LeadBitrix) error {

	const query = `INSERT INTO leads(name,email,phone,source)
					VALUES ($1,$2,$3,$4)`

	_, err := s.db.Exec(ctx, query,
		lead.Name,
		lead.Email,
		lead.Phone,
		lead.Source)

	if err != nil {
		return fmt.Errorf("failed to save lead: %w", err)
	}
	
	return nil
}
