package employees

import (
	"context"
	"fmt"

	"vendix/internal/database"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, schema string, employee *Employee) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.employees (
			id, employee_code, first_name, last_name, email, phone, tax_id, 
			address, city, state, postal_code, country, department, position, 
			hire_date, salary, salary_type, is_active, metadata, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
	`, schema)

	_, err := r.db.ExecContext(ctx, query,
		employee.ID,
		employee.EmployeeCode,
		employee.FirstName,
		employee.LastName,
		employee.Email,
		employee.Phone,
		employee.TaxID,
		employee.Address,
		employee.City,
		employee.State,
		employee.PostalCode,
		employee.Country,
		employee.Department,
		employee.Position,
		employee.HireDate,
		employee.Salary,
		employee.SalaryType,
		employee.IsActive,
		employee.Metadata,
		employee.CreatedAt,
		employee.UpdatedAt,
	)

	return err
}

func (r *Repository) GetByID(ctx context.Context, schema string, id uuid.UUID) (*Employee, error) {
	var employee Employee
	query := fmt.Sprintf(`
		SELECT id, employee_code, first_name, last_name, email, phone, tax_id, 
			address, city, state, postal_code, country, department, position, 
			hire_date, salary, salary_type, is_active, metadata, created_at, updated_at
		FROM %s.employees
		WHERE id = $1
	`, schema)

	err := r.db.GetContext(ctx, &employee, query, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("employee not found")
		}
		return nil, err
	}
	return &employee, nil
}

func (r *Repository) GetByCode(ctx context.Context, schema, code string) (*Employee, error) {
	var employee Employee
	query := fmt.Sprintf(`
		SELECT id, employee_code, first_name, last_name, email, phone, tax_id, 
			address, city, state, postal_code, country, department, position, 
			hire_date, salary, salary_type, is_active, metadata, created_at, updated_at
		FROM %s.employees
		WHERE employee_code = $1
	`, schema)

	err := r.db.GetContext(ctx, &employee, query, code)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("employee not found")
		}
		return nil, err
	}
	return &employee, nil
}

func (r *Repository) List(ctx context.Context, schema string) ([]*Employee, error) {
	var employees []*Employee
	query := fmt.Sprintf(`
		SELECT id, employee_code, first_name, last_name, email, phone, tax_id, 
			address, city, state, postal_code, country, department, position, 
			hire_date, salary, salary_type, is_active, metadata, created_at, updated_at
		FROM %s.employees
		ORDER BY last_name ASC, first_name ASC
	`, schema)

	err := r.db.SelectContext(ctx, &employees, query)
	return employees, err
}

func (r *Repository) ListActive(ctx context.Context, schema string) ([]*Employee, error) {
	var employees []*Employee
	query := fmt.Sprintf(`
		SELECT id, employee_code, first_name, last_name, email, phone, tax_id, 
			address, city, state, postal_code, country, department, position, 
			hire_date, salary, salary_type, is_active, metadata, created_at, updated_at
		FROM %s.employees
		WHERE is_active = true
		ORDER BY last_name ASC, first_name ASC
	`, schema)

	err := r.db.SelectContext(ctx, &employees, query)
	return employees, err
}

func (r *Repository) Update(ctx context.Context, schema string, employee *Employee) error {
	query := fmt.Sprintf(`
		UPDATE %s.employees
		SET first_name = $1, last_name = $2, email = $3, phone = $4, 
			address = $5, city = $6, state = $7, postal_code = $8, 
			department = $9, position = $10, salary = $11, salary_type = $12, 
			is_active = $13, updated_at = $14
		WHERE id = $15
	`, schema)

	_, err := r.db.ExecContext(ctx, query,
		employee.FirstName,
		employee.LastName,
		employee.Email,
		employee.Phone,
		employee.Address,
		employee.City,
		employee.State,
		employee.PostalCode,
		employee.Department,
		employee.Position,
		employee.Salary,
		employee.SalaryType,
		employee.IsActive,
		employee.UpdatedAt,
		employee.ID,
	)

	return err
}

func (r *Repository) Delete(ctx context.Context, schema string, id uuid.UUID) error {
	query := fmt.Sprintf(`
		DELETE FROM %s.employees
		WHERE id = $1
	`, schema)

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

