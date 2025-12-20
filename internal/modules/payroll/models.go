package payroll

import (
	"time"

	"github.com/google/uuid"
)

type PayrollPeriod struct {
	ID          uuid.UUID `db:"id" json:"id"`
	PeriodCode  string    `db:"period_code" json:"period_code"` // Format: YYYY-MM
	PeriodStart time.Time  `db:"period_start" json:"period_start"`
	PeriodEnd   time.Time  `db:"period_end" json:"period_end"`
	Status      string    `db:"status" json:"status"` // 'draft', 'processing', 'completed', 'paid'
	TotalGross  float64   `db:"total_gross" json:"total_gross"`
	TotalDeductions float64 `db:"total_deductions" json:"total_deductions"`
	TotalNet    float64   `db:"total_net" json:"total_net"`
	Notes       *string   `db:"notes" json:"notes,omitempty"`
	CreatedBy   uuid.UUID `db:"created_by" json:"created_by"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type PayrollEntry struct {
	ID                uuid.UUID `db:"id" json:"id"`
	PayrollPeriodID   uuid.UUID `db:"payroll_period_id" json:"payroll_period_id"`
	EmployeeID        uuid.UUID `db:"employee_id" json:"employee_id"`
	EmployeeCode      string    `db:"employee_code" json:"employee_code"`
	EmployeeName      string    `db:"employee_name" json:"employee_name"`
	Department        *string   `db:"department" json:"department,omitempty"`
	Position          *string   `db:"position" json:"position,omitempty"`
	Salary            float64   `db:"salary" json:"salary"`
	SalaryType        string    `db:"salary_type" json:"salary_type"`
	HoursWorked       *float64  `db:"hours_worked" json:"hours_worked,omitempty"`
	DaysWorked        *float64  `db:"days_worked" json:"days_worked,omitempty"`
	GrossSalary       float64   `db:"gross_salary" json:"gross_salary"`
	OvertimeHours     *float64  `db:"overtime_hours" json:"overtime_hours,omitempty"`
	OvertimePay       float64   `db:"overtime_pay" json:"overtime_pay"`
	Bonuses           float64   `db:"bonuses" json:"bonuses"`
	Commissions       float64   `db:"commissions" json:"commissions"`
	TotalGross        float64   `db:"total_gross" json:"total_gross"`
	TaxDeduction      float64   `db:"tax_deduction" json:"tax_deduction"` // ISR
	SocialSecurity    float64   `db:"social_security" json:"social_security"` // TSS
	OtherDeductions    float64   `db:"other_deductions" json:"other_deductions"`
	TotalDeductions   float64   `db:"total_deductions" json:"total_deductions"`
	NetSalary         float64   `db:"net_salary" json:"net_salary"`
	Notes             *string   `db:"notes" json:"notes,omitempty"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`
}

type CreatePayrollPeriodRequest struct {
	PeriodCode  string    `json:"period_code" validate:"required"` // Format: YYYY-MM
	PeriodStart time.Time `json:"period_start" validate:"required"`
	PeriodEnd   time.Time `json:"period_end" validate:"required"`
	Notes       *string   `json:"notes,omitempty"`
}

type UpdatePayrollPeriodRequest struct {
	Status *string `json:"status,omitempty" validate:"omitempty,oneof=draft processing completed paid"`
	Notes  *string `json:"notes,omitempty"`
}

type CreatePayrollEntryRequest struct {
	EmployeeID      string   `json:"employee_id" validate:"required"`
	HoursWorked     *float64 `json:"hours_worked,omitempty"`
	DaysWorked      *float64 `json:"days_worked,omitempty"`
	OvertimeHours   *float64 `json:"overtime_hours,omitempty"`
	Bonuses         *float64 `json:"bonuses,omitempty"`
	Commissions     *float64 `json:"commissions,omitempty"`
	OtherDeductions *float64 `json:"other_deductions,omitempty"`
	Notes           *string  `json:"notes,omitempty"`
}

type UpdatePayrollEntryRequest struct {
	HoursWorked     *float64 `json:"hours_worked,omitempty"`
	DaysWorked      *float64 `json:"days_worked,omitempty"`
	OvertimeHours   *float64 `json:"overtime_hours,omitempty"`
	Bonuses         *float64 `json:"bonuses,omitempty"`
	Commissions     *float64 `json:"commissions,omitempty"`
	OtherDeductions *float64 `json:"other_deductions,omitempty"`
	Notes           *string  `json:"notes,omitempty"`
}

type PayrollPeriodResponse struct {
	ID             string    `json:"id"`
	PeriodCode     string    `json:"period_code"`
	PeriodStart    time.Time `json:"period_start"`
	PeriodEnd      time.Time `json:"period_end"`
	Status         string    `json:"status"`
	TotalGross     float64   `json:"total_gross"`
	TotalDeductions float64  `json:"total_deductions"`
	TotalNet       float64   `json:"total_net"`
	Notes          *string   `json:"notes,omitempty"`
	CreatedBy      string    `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Entries        []*PayrollEntryResponse `json:"entries,omitempty"`
}

type PayrollEntryResponse struct {
	ID              string    `json:"id"`
	PayrollPeriodID string    `json:"payroll_period_id"`
	EmployeeID      string    `json:"employee_id"`
	EmployeeCode    string    `json:"employee_code"`
	EmployeeName    string    `json:"employee_name"`
	Department      *string   `json:"department,omitempty"`
	Position        *string   `json:"position,omitempty"`
	Salary          float64   `json:"salary"`
	SalaryType      string    `json:"salary_type"`
	HoursWorked     *float64  `json:"hours_worked,omitempty"`
	DaysWorked      *float64  `json:"days_worked,omitempty"`
	GrossSalary     float64   `json:"gross_salary"`
	OvertimeHours   *float64  `json:"overtime_hours,omitempty"`
	OvertimePay     float64   `json:"overtime_pay"`
	Bonuses         float64  `json:"bonuses"`
	Commissions     float64  `json:"commissions"`
	TotalGross      float64  `json:"total_gross"`
	TaxDeduction    float64  `json:"tax_deduction"`
	SocialSecurity  float64  `json:"social_security"`
	OtherDeductions float64  `json:"other_deductions"`
	TotalDeductions float64  `json:"total_deductions"`
	NetSalary       float64  `json:"net_salary"`
	Notes           *string  `json:"notes,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (p *PayrollPeriod) ToResponse() *PayrollPeriodResponse {
	return &PayrollPeriodResponse{
		ID:              p.ID.String(),
		PeriodCode:      p.PeriodCode,
		PeriodStart:     p.PeriodStart,
		PeriodEnd:       p.PeriodEnd,
		Status:          p.Status,
		TotalGross:      p.TotalGross,
		TotalDeductions: p.TotalDeductions,
		TotalNet:        p.TotalNet,
		Notes:           p.Notes,
		CreatedBy:       p.CreatedBy.String(),
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
	}
}

func (e *PayrollEntry) ToResponse() *PayrollEntryResponse {
	return &PayrollEntryResponse{
		ID:              e.ID.String(),
		PayrollPeriodID: e.PayrollPeriodID.String(),
		EmployeeID:      e.EmployeeID.String(),
		EmployeeCode:    e.EmployeeCode,
		EmployeeName:    e.EmployeeName,
		Department:      e.Department,
		Position:        e.Position,
		Salary:          e.Salary,
		SalaryType:      e.SalaryType,
		HoursWorked:     e.HoursWorked,
		DaysWorked:      e.DaysWorked,
		GrossSalary:     e.GrossSalary,
		OvertimeHours:   e.OvertimeHours,
		OvertimePay:     e.OvertimePay,
		Bonuses:         e.Bonuses,
		Commissions:     e.Commissions,
		TotalGross:      e.TotalGross,
		TaxDeduction:    e.TaxDeduction,
		SocialSecurity:  e.SocialSecurity,
		OtherDeductions: e.OtherDeductions,
		TotalDeductions: e.TotalDeductions,
		NetSalary:       e.NetSalary,
		Notes:           e.Notes,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
	}
}

