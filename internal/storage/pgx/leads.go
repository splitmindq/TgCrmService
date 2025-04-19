package pgx

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"lead-bitrix/entities"
	"strings"
)

func (s *Storage) SaveLead(ctx context.Context, lead entities.Lead) error {

	if lead.Id == 0 {
		const query = `INSERT INTO leads(name,email,phone,source)
				VALUES ($1,$2,$3,$4)`

		_, err := s.db.Exec(ctx, query,
			lead.Form.Name,
			lead.Form.Email,
			lead.Form.Phone,
			lead.Form.Source)

		if err != nil {
			return fmt.Errorf("failed to save lead: %w", err)
		}
	} else {
		const query = `INSERT INTO leads(id,name,email,phone,source)
    							VALUES ($1,$2,$3,$4,$5)
    							ON CONFLICT(id) DO UPDATE
    							SET name = EXCLUDED.name,
    							    email = EXCLUDED.email,
    							    phone = EXCLUDED.phone,
    							    source = EXCLUDED.source`
		_, err := s.db.Exec(ctx, query,
			lead.Id,
			lead.Form.Name,
			lead.Form.Email,
			lead.Form.Phone,
			lead.Form.Source)

		if err != nil {
			return fmt.Errorf("failed to save lead: %w", err)
		}
	}
	return nil
}

func (s *Storage) DeleteLead(ctx context.Context, email string) error {
	const query = `DELETE FROM leads WHERE email=$1`

	_, err := s.db.Exec(ctx, query, email)
	if err != nil {
		return fmt.Errorf("failed to delete lead: %w", err)
	}

	return nil
}

func (s *Storage) GetLead(ctx context.Context, email string) (*entities.Lead, error) {
	var lead entities.Lead

	const query = `SELECT id,name, email, phone, source FROM leads WHERE email=$1`
	err := s.db.QueryRow(ctx, query, email).Scan(
		&lead.Id,
		&lead.Form.Name,
		&lead.Form.Email,
		&lead.Form.Phone,
		&lead.Form.Source,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch lead: %w", err)
	}
	return &lead, nil
}

func (s *Storage) GetAllLeads(ctx context.Context) ([]entities.Lead, error) {

	var leads []entities.Lead
	const query = `SELECT id,name, email, phone, source FROM leads`
	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch leads: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var lead entities.Lead
		err := rows.Scan(&lead.Id, &lead.Form.Name, &lead.Form.Email, &lead.Form.Phone, &lead.Form.Source)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows: %w", err)
		}
		leads = append(leads, lead)
	}
	return leads, nil
}

func (s *Storage) UpdateLead(ctx context.Context, lead entities.JsonForm) (*entities.Lead, error) {
	const op = "storage.pgx.UpdateLead"

	if lead.Phone == "" {
		return nil, fmt.Errorf("%s: phone is required", op)
	}

	query := `
        UPDATE leads
        SET 
            name = COALESCE($1, name),
            email = COALESCE(NULLIF($2, ''), email),
            source = COALESCE($3, source)
        WHERE phone = $4
        RETURNING id, name, email, phone, source`

	var result entities.Lead
	err := s.db.QueryRow(ctx, query,
		lead.Name,
		lead.Email,
		lead.Source,
		lead.Phone,
	).Scan(
		&result.Id,
		&result.Form.Name,
		&result.Form.Email,
		&result.Form.Phone,
		&result.Form.Source,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				if strings.Contains(pgErr.Message, "email") {
					return nil, ErrEmailExists
				} else if strings.Contains(pgErr.Message, "phone") {
					return nil, ErrPhoneExists
				}
			}
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &result, nil
}
