package payroll

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

// PayrollPeriod methods

func (r *Repository) CreatePeriod(ctx context.Context, schema string, period *PayrollPeriod) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.payroll_periods (
			id, period_code, period_start, period_end, status, 
			total_gross, total_deductions, total_net, notes, 
			created_by, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`, schema)

	_, err := r.db.ExecContext(ctx, query,
		period.ID,
		period.PeriodCode,
		period.PeriodStart,
		period.PeriodEnd,
		period.Status,
		period.TotalGross,
		period.TotalDeductions,
		period.TotalNet,
		period.Notes,
		period.CreatedBy,
		period.CreatedAt,
		period.UpdatedAt,
	)

	return err
}

func (r *Repository) GetPeriodByID(ctx context.Context, schema string, id uuid.UUID) (*PayrollPeriod, error) {
	var period PayrollPeriod
	query := fmt.Sprintf(`
		SELECT id, period_code, period_start, period_end, status, 
			total_gross, total_deductions, total_net, notes, 
			created_by, created_at, updated_at
		FROM %s.payroll_periods
		WHERE id = $1
	`, schema)

	err := r.db.GetContext(ctx, &period, query, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("payroll period not found")
		}
		return nil, err
	}
	return &period, nil
}

func (r *Repository) GetPeriodByCode(ctx context.Context, schema, code string) (*PayrollPeriod, error) {
	var period PayrollPeriod
	query := fmt.Sprintf(`
		SELECT id, period_code, period_start, period_end, status, 
			total_gross, total_deductions, total_net, notes, 
			created_by, created_at, updated_at
		FROM %s.payroll_periods
		WHERE period_code = $1
	`, schema)

	err := r.db.GetContext(ctx, &period, query, code)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("payroll period not found")
		}
		return nil, err
	}
	return &period, nil
}

func (r *Repository) ListPeriods(ctx context.Context, schema string) ([]*PayrollPeriod, error) {
	var periods []*PayrollPeriod
	query := fmt.Sprintf(`
		SELECT id, period_code, period_start, period_end, status, 
			total_gross, total_deductions, total_net, notes, 
			created_by, created_at, updated_at
		FROM %s.payroll_periods
		ORDER BY period_code DESC
	`, schema)

	err := r.db.SelectContext(ctx, &periods, query)
	return periods, err
}

func (r *Repository) UpdatePeriod(ctx context.Context, schema string, period *PayrollPeriod) error {
	query := fmt.Sprintf(`
		UPDATE %s.payroll_periods
		SET status = $1, total_gross = $2, total_deductions = $3, 
			total_net = $4, notes = $5, updated_at = $6
		WHERE id = $7
	`, schema)

	_, err := r.db.ExecContext(ctx, query,
		period.Status,
		period.TotalGross,
		period.TotalDeductions,
		period.TotalNet,
		period.Notes,
		period.UpdatedAt,
		period.ID,
	)

	return err
}

func (r *Repository) DeletePeriod(ctx context.Context, schema string, id uuid.UUID) error {
	query := fmt.Sprintf(`
		DELETE FROM %s.payroll_periods
		WHERE id = $1
	`, schema)

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// PayrollEntry methods

func (r *Repository) CreateEntry(ctx context.Context, schema string, entry *PayrollEntry) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.payroll_entries (
			id, payroll_period_id, employee_id, employee_code, employee_name,
			department, position, salary, salary_type, hours_worked, days_worked,
			gross_salary, overtime_hours, overtime_pay, bonuses, commissions,
			total_gross, tax_deduction, social_security, other_deductions,
			total_deductions, net_salary, notes, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25)
	`, schema)

	_, err := r.db.ExecContext(ctx, query,
		entry.ID,
		entry.PayrollPeriodID,
		entry.EmployeeID,
		entry.EmployeeCode,
		entry.EmployeeName,
		entry.Department,
		entry.Position,
		entry.Salary,
		entry.SalaryType,
		entry.HoursWorked,
		entry.DaysWorked,
		entry.GrossSalary,
		entry.OvertimeHours,
		entry.OvertimePay,
		entry.Bonuses,
		entry.Commissions,
		entry.TotalGross,
		entry.TaxDeduction,
		entry.SocialSecurity,
		entry.OtherDeductions,
		entry.TotalDeductions,
		entry.NetSalary,
		entry.Notes,
		entry.CreatedAt,
		entry.UpdatedAt,
	)

	return err
}

func (r *Repository) GetEntryByID(ctx context.Context, schema string, id uuid.UUID) (*PayrollEntry, error) {
	var entry PayrollEntry
	query := fmt.Sprintf(`
		SELECT id, payroll_period_id, employee_id, employee_code, employee_name,
			department, position, salary, salary_type, hours_worked, days_worked,
			gross_salary, overtime_hours, overtime_pay, bonuses, commissions,
			total_gross, tax_deduction, social_security, other_deductions,
			total_deductions, net_salary, notes, created_at, updated_at
		FROM %s.payroll_entries
		WHERE id = $1
	`, schema)

	err := r.db.GetContext(ctx, &entry, query, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("payroll entry not found")
		}
		return nil, err
	}
	return &entry, nil
}

func (r *Repository) ListEntriesByPeriod(ctx context.Context, schema string, periodID uuid.UUID) ([]*PayrollEntry, error) {
	var entries []*PayrollEntry
	query := fmt.Sprintf(`
		SELECT id, payroll_period_id, employee_id, employee_code, employee_name,
			department, position, salary, salary_type, hours_worked, days_worked,
			gross_salary, overtime_hours, overtime_pay, bonuses, commissions,
			total_gross, tax_deduction, social_security, other_deductions,
			total_deductions, net_salary, notes, created_at, updated_at
		FROM %s.payroll_entries
		WHERE payroll_period_id = $1
		ORDER BY employee_name ASC
	`, schema)

	err := r.db.SelectContext(ctx, &entries, query, periodID)
	return entries, err
}

func (r *Repository) UpdateEntry(ctx context.Context, schema string, entry *PayrollEntry) error {
	query := fmt.Sprintf(`
		UPDATE %s.payroll_entries
		SET hours_worked = $1, days_worked = $2, overtime_hours = $3,
			overtime_pay = $4, bonuses = $5, commissions = $6,
			total_gross = $7, tax_deduction = $8, social_security = $9,
			other_deductions = $10, total_deductions = $11, net_salary = $12,
			notes = $13, updated_at = $14
		WHERE id = $15
	`, schema)

	_, err := r.db.ExecContext(ctx, query,
		entry.HoursWorked,
		entry.DaysWorked,
		entry.OvertimeHours,
		entry.OvertimePay,
		entry.Bonuses,
		entry.Commissions,
		entry.TotalGross,
		entry.TaxDeduction,
		entry.SocialSecurity,
		entry.OtherDeductions,
		entry.TotalDeductions,
		entry.NetSalary,
		entry.Notes,
		entry.UpdatedAt,
		entry.ID,
	)

	return err
}

func (r *Repository) DeleteEntry(ctx context.Context, schema string, id uuid.UUID) error {
	query := fmt.Sprintf(`
		DELETE FROM %s.payroll_entries
		WHERE id = $1
	`, schema)

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *Repository) DeleteEntriesByPeriod(ctx context.Context, schema string, periodID uuid.UUID) error {
	query := fmt.Sprintf(`
		DELETE FROM %s.payroll_entries
		WHERE payroll_period_id = $1
	`, schema)

	_, err := r.db.ExecContext(ctx, query, periodID)
	return err
}

