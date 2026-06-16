package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/krizad/go-gin-api/internal/core/domain"
)

type ExampleRepository interface {
	Create(ctx context.Context, name, email string) (*domain.Example, error)
	GetByID(ctx context.Context, id int64) (*domain.Example, error)
	GetByEmail(ctx context.Context, email string) (*domain.Example, error)
	List(ctx context.Context) ([]domain.Example, error)
	Update(ctx context.Context, id int64, name, email string) (*domain.Example, error)
	Delete(ctx context.Context, id int64) (bool, error)
	Count(ctx context.Context) (int, error)
}

type exampleRepo struct {
	db *sql.DB
}

func NewExampleRepository(db *sql.DB) ExampleRepository {
	return &exampleRepo{db: db}
}

func (r *exampleRepo) Create(ctx context.Context, name, email string) (*domain.Example, error) {
	example := &domain.Example{}
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO examples (name, email) VALUES ($1, $2) RETURNING id, name, email, created_at, updated_at`,
		name, email,
	).Scan(&example.ID, &example.Name, &example.Email, &example.CreatedAt, &example.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create example: %w", err)
	}
	return example, nil
}

func (r *exampleRepo) GetByID(ctx context.Context, id int64) (*domain.Example, error) {
	example := &domain.Example{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, email, created_at, updated_at FROM examples WHERE id = $1`, id,
	).Scan(&example.ID, &example.Name, &example.Email, &example.CreatedAt, &example.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get example: %w", err)
	}
	return example, nil
}

func (r *exampleRepo) GetByEmail(ctx context.Context, email string) (*domain.Example, error) {
	example := &domain.Example{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, email, created_at, updated_at FROM examples WHERE email = $1`, email,
	).Scan(&example.ID, &example.Name, &example.Email, &example.CreatedAt, &example.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get example by email: %w", err)
	}
	return example, nil
}

func (r *exampleRepo) List(ctx context.Context) ([]domain.Example, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, email, created_at, updated_at FROM examples ORDER BY id`,
	)
	if err != nil {
		return nil, fmt.Errorf("list examples: %w", err)
	}
	defer rows.Close()

	var examples []domain.Example
	for rows.Next() {
		var e domain.Example
		if err := rows.Scan(&e.ID, &e.Name, &e.Email, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan example: %w", err)
		}
		examples = append(examples, e)
	}
	if examples == nil {
		examples = []domain.Example{}
	}
	return examples, rows.Err()
}

func (r *exampleRepo) Update(ctx context.Context, id int64, name, email string) (*domain.Example, error) {
	example := &domain.Example{}
	err := r.db.QueryRowContext(ctx,
		`UPDATE examples SET name = $1, email = $2, updated_at = NOW() WHERE id = $3 RETURNING id, name, email, created_at, updated_at`,
		name, email, id,
	).Scan(&example.ID, &example.Name, &example.Email, &example.CreatedAt, &example.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("update example: %w", err)
	}
	return example, nil
}

func (r *exampleRepo) Delete(ctx context.Context, id int64) (bool, error) {
	result, err := r.db.ExecContext(ctx, `DELETE FROM examples WHERE id = $1`, id)
	if err != nil {
		return false, fmt.Errorf("delete example: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("rows affected: %w", err)
	}
	return affected > 0, nil
}

func (r *exampleRepo) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM examples`).Scan(&count)
	return count, err
}
