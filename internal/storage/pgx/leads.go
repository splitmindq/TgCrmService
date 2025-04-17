package pgx

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

func (s *Storage) DeleteLead(ctx context.Context, email string) error {
	const query = `DELETE FROM leads WHERE email=$1`

	_, err := s.db.Exec(ctx, query, email)
	if err != nil {
		return fmt.Errorf("failed to delete lead: %w", err)
	}

	return nil
}

func (s *Storage) GetLead(ctx context.Context, email string) (*entities.LeadBitrix, error) {
	var lead entities.LeadBitrix

	const query = `SELECT name, email, phone, source FROM leads WHERE email=$1`
	err := s.db.QueryRow(ctx, query, email).Scan(
		&lead.Name,
		&lead.Email,
		&lead.Phone,
		&lead.Source,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch lead: %w", err)
	}
	return &lead, nil
}

func (s *Storage) GetAllLeads(ctx context.Context) ([]entities.LeadBitrix, error) {

	var leads []entities.LeadBitrix
	const query = `SELECT name, email, phone, source FROM leads`
	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch leads: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var lead entities.LeadBitrix
		err := rows.Scan(&lead.Name, &lead.Email, &lead.Phone, &lead.Source)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows: %w", err)
		}
		leads = append(leads, lead)
	}
	return leads, nil
}

func (s *Storage) UpdateLead(ctx context.Context, lead entities.LeadBitrix) (*entities.LeadBitrix, error) {
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
        RETURNING name, email, phone, source`

	var result entities.LeadBitrix
	err := s.db.QueryRow(ctx, query,
		lead.Name,
		lead.Email,
		lead.Source,
		lead.Phone,
	).Scan(
		&result.Name,
		&result.Email,
		&result.Phone,
		&result.Source,
	)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrEmailExists
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &result, nil
}
