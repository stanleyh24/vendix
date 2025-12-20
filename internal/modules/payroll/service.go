package payroll

import (
	"context"
	"fmt"
	"time"

	"vendix/internal/config"
	"vendix/internal/database"
	"vendix/internal/logger"
	"vendix/internal/modules/employees"

	"github.com/google/uuid"
)

type Service struct {
	repo         *Repository
	employeeSvc  *employees.Service
	cfg          *config.Config
}

func NewService(db *database.DB, cfg *config.Config) *Service {
	return &Service{
		repo:        NewRepository(db),
		employeeSvc: employees.NewService(db, cfg),
		cfg:         cfg,
	}
}

// Constants for Dominican Republic payroll calculations
const (
	ISR_RATE_LOW    = 0.15  // 15% for first bracket
	ISR_RATE_HIGH   = 0.20  // 20% for second bracket
	ISR_THRESHOLD   = 416220.00 // Monthly threshold in DOP
	TSS_RATE        = 0.0304   // 3.04% (employee contribution)
	OVERTIME_RATE   = 1.5      // 1.5x hourly rate for overtime
)

// calculateISR calculates income tax (ISR) for Dominican Republic
func calculateISR(grossSalary float64) float64 {
	if grossSalary <= ISR_THRESHOLD {
		return grossSalary * ISR_RATE_LOW
	}
	return (ISR_THRESHOLD * ISR_RATE_LOW) + ((grossSalary - ISR_THRESHOLD) * ISR_RATE_HIGH)
}

// calculateTSS calculates Social Security (TSS) contribution
func calculateTSS(grossSalary float64) float64 {
	return grossSalary * TSS_RATE
}

// calculateGrossSalary calculates gross salary based on type
func (s *Service) calculateGrossSalary(employee *employees.Employee, hoursWorked, daysWorked *float64) (float64, error) {
	if employee.Salary == nil {
		return 0, fmt.Errorf("employee has no salary configured")
	}

	baseSalary := *employee.Salary

	switch employee.SalaryType {
	case "monthly":
		// For monthly, return full salary
		return baseSalary, nil
	case "hourly":
		if hoursWorked == nil {
			return 0, fmt.Errorf("hours worked required for hourly employees")
		}
		return baseSalary * (*hoursWorked), nil
	case "daily":
		if daysWorked == nil {
			return 0, fmt.Errorf("days worked required for daily employees")
		}
		// Assuming 22 working days per month standard
		dailyRate := baseSalary / 22.0
		return dailyRate * (*daysWorked), nil
	default:
		return 0, fmt.Errorf("invalid salary type: %s", employee.SalaryType)
	}
}

// calculateOvertimePay calculates overtime pay
func (s *Service) calculateOvertimePay(employee *employees.Employee, overtimeHours *float64) (float64, error) {
	if overtimeHours == nil || *overtimeHours == 0 {
		return 0, nil
	}

	if employee.Salary == nil {
		return 0, fmt.Errorf("employee has no salary configured")
	}

	var hourlyRate float64
	switch employee.SalaryType {
	case "monthly":
		// Assuming 160 hours per month (8 hours * 20 days)
		hourlyRate = *employee.Salary / 160.0
	case "hourly":
		hourlyRate = *employee.Salary
	case "daily":
		// Assuming 8 hours per day
		hourlyRate = (*employee.Salary / 22.0) / 8.0
	default:
		return 0, fmt.Errorf("invalid salary type: %s", employee.SalaryType)
	}

	return hourlyRate * OVERTIME_RATE * (*overtimeHours), nil
}

func (s *Service) CreatePeriod(ctx context.Context, schema, userID string, req *CreatePayrollPeriodRequest) (*PayrollPeriodResponse, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Check if period code already exists
	existing, _ := s.repo.GetPeriodByCode(ctx, schema, req.PeriodCode)
	if existing != nil {
		return nil, fmt.Errorf("payroll period code already exists")
	}

	period := &PayrollPeriod{
		ID:             uuid.New(),
		PeriodCode:     req.PeriodCode,
		PeriodStart:    req.PeriodStart,
		PeriodEnd:       req.PeriodEnd,
		Status:          "draft",
		TotalGross:     0,
		TotalDeductions: 0,
		TotalNet:        0,
		Notes:           req.Notes,
		CreatedBy:       userUUID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.repo.CreatePeriod(ctx, schema, period); err != nil {
		return nil, err
	}

	logger.Info("Payroll period created", "id", period.ID, "code", period.PeriodCode, "schema", schema)

	return period.ToResponse(), nil
}

func (s *Service) ListPeriods(ctx context.Context, schema string) ([]*PayrollPeriodResponse, error) {
	periods, err := s.repo.ListPeriods(ctx, schema)
	if err != nil {
		return nil, err
	}

	responses := make([]*PayrollPeriodResponse, len(periods))
	for i, p := range periods {
		responses[i] = p.ToResponse()
	}

	return responses, nil
}

func (s *Service) GetPeriodByID(ctx context.Context, schema, id string) (*PayrollPeriodResponse, error) {
	periodID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid period ID: %w", err)
	}

	period, err := s.repo.GetPeriodByID(ctx, schema, periodID)
	if err != nil {
		return nil, err
	}

	// Load entries
	entries, err := s.repo.ListEntriesByPeriod(ctx, schema, periodID)
	if err != nil {
		return nil, err
	}

	response := period.ToResponse()
	response.Entries = make([]*PayrollEntryResponse, len(entries))
	for i, e := range entries {
		response.Entries[i] = e.ToResponse()
	}

	return response, nil
}

func (s *Service) UpdatePeriod(ctx context.Context, schema, id string, req *UpdatePayrollPeriodRequest) (*PayrollPeriodResponse, error) {
	periodID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid period ID: %w", err)
	}

	period, err := s.repo.GetPeriodByID(ctx, schema, periodID)
	if err != nil {
		return nil, err
	}

	if req.Status != nil {
		period.Status = *req.Status
	}
	if req.Notes != nil {
		period.Notes = req.Notes
	}

	period.UpdatedAt = time.Now()

	if err := s.repo.UpdatePeriod(ctx, schema, period); err != nil {
		return nil, err
	}

	logger.Info("Payroll period updated", "id", period.ID, "schema", schema)

	return period.ToResponse(), nil
}

func (s *Service) DeletePeriod(ctx context.Context, schema, id string) error {
	periodID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid period ID: %w", err)
	}

	if err := s.repo.DeletePeriod(ctx, schema, periodID); err != nil {
		return err
	}

	logger.Info("Payroll period deleted", "id", periodID, "schema", schema)

	return nil
}

func (s *Service) CreateEntry(ctx context.Context, schema, periodID string, req *CreatePayrollEntryRequest) (*PayrollEntryResponse, error) {
	periodUUID, err := uuid.Parse(periodID)
	if err != nil {
		return nil, fmt.Errorf("invalid period ID: %w", err)
	}

	employeeUUID, err := uuid.Parse(req.EmployeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee ID: %w", err)
	}

	// Get employee via service
	employeeResp, err := s.employeeSvc.GetByID(ctx, schema, req.EmployeeID)
	if err != nil {
		return nil, fmt.Errorf("employee not found: %w", err)
	}

	// Calculate gross salary using response data
	if employeeResp.Salary == nil {
		return nil, fmt.Errorf("employee has no salary configured")
	}

	baseSalary := *employeeResp.Salary
	var grossSalary float64

	switch employeeResp.SalaryType {
	case "monthly":
		grossSalary = baseSalary
	case "hourly":
		if req.HoursWorked == nil {
			return nil, fmt.Errorf("hours worked required for hourly employees")
		}
		grossSalary = baseSalary * (*req.HoursWorked)
	case "daily":
		if req.DaysWorked == nil {
			return nil, fmt.Errorf("days worked required for daily employees")
		}
		dailyRate := baseSalary / 22.0
		grossSalary = dailyRate * (*req.DaysWorked)
	default:
		return nil, fmt.Errorf("invalid salary type: %s", employeeResp.SalaryType)
	}

	// Calculate overtime pay
	var overtimePay float64
	if req.OvertimeHours != nil && *req.OvertimeHours > 0 {
		var hourlyRate float64
		switch employeeResp.SalaryType {
		case "monthly":
			hourlyRate = baseSalary / 160.0
		case "hourly":
			hourlyRate = baseSalary
		case "daily":
			hourlyRate = (baseSalary / 22.0) / 8.0
		default:
			return nil, fmt.Errorf("invalid salary type: %s", employeeResp.SalaryType)
		}
		overtimePay = hourlyRate * OVERTIME_RATE * (*req.OvertimeHours)
	}

	bonuses := 0.0
	if req.Bonuses != nil {
		bonuses = *req.Bonuses
	}

	commissions := 0.0
	if req.Commissions != nil {
		commissions = *req.Commissions
	}

	totalGross := grossSalary + overtimePay + bonuses + commissions

	// Calculate deductions
	taxDeduction := calculateISR(totalGross)
	socialSecurity := calculateTSS(totalGross)

	otherDeductions := 0.0
	if req.OtherDeductions != nil {
		otherDeductions = *req.OtherDeductions
	}

	totalDeductions := taxDeduction + socialSecurity + otherDeductions
	netSalary := totalGross - totalDeductions

	entry := &PayrollEntry{
		ID:              uuid.New(),
		PayrollPeriodID: periodUUID,
		EmployeeID:      employeeUUID,
		EmployeeCode:    employeeResp.EmployeeCode,
		EmployeeName:    fmt.Sprintf("%s %s", employeeResp.FirstName, employeeResp.LastName),
		Department:      employeeResp.Department,
		Position:        employeeResp.Position,
		Salary:          grossSalary,
		SalaryType:      employeeResp.SalaryType,
		HoursWorked:     req.HoursWorked,
		DaysWorked:      req.DaysWorked,
		GrossSalary:     grossSalary,
		OvertimeHours:   req.OvertimeHours,
		OvertimePay:     overtimePay,
		Bonuses:         bonuses,
		Commissions:     commissions,
		TotalGross:      totalGross,
		TaxDeduction:    taxDeduction,
		SocialSecurity:  socialSecurity,
		OtherDeductions: otherDeductions,
		TotalDeductions: totalDeductions,
		NetSalary:       netSalary,
		Notes:           req.Notes,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.repo.CreateEntry(ctx, schema, entry); err != nil {
		return nil, err
	}

	// Update period totals
	if err := s.updatePeriodTotals(ctx, schema, periodUUID); err != nil {
		logger.Error("Failed to update period totals", "error", err)
	}

	logger.Info("Payroll entry created", "id", entry.ID, "employee_id", employeeUUID, "schema", schema)

	return entry.ToResponse(), nil
}

func (s *Service) updatePeriodTotals(ctx context.Context, schema string, periodID uuid.UUID) error {
	entries, err := s.repo.ListEntriesByPeriod(ctx, schema, periodID)
	if err != nil {
		return err
	}

	var totalGross, totalDeductions, totalNet float64
	for _, entry := range entries {
		totalGross += entry.TotalGross
		totalDeductions += entry.TotalDeductions
		totalNet += entry.NetSalary
	}

	period, err := s.repo.GetPeriodByID(ctx, schema, periodID)
	if err != nil {
		return err
	}

	period.TotalGross = totalGross
	period.TotalDeductions = totalDeductions
	period.TotalNet = totalNet
	period.UpdatedAt = time.Now()

	return s.repo.UpdatePeriod(ctx, schema, period)
}

func (s *Service) ListEntriesByPeriod(ctx context.Context, schema, periodID string) ([]*PayrollEntryResponse, error) {
	periodUUID, err := uuid.Parse(periodID)
	if err != nil {
		return nil, fmt.Errorf("invalid period ID: %w", err)
	}

	entries, err := s.repo.ListEntriesByPeriod(ctx, schema, periodUUID)
	if err != nil {
		return nil, err
	}

	responses := make([]*PayrollEntryResponse, len(entries))
	for i, e := range entries {
		responses[i] = e.ToResponse()
	}

	return responses, nil
}

func (s *Service) UpdateEntry(ctx context.Context, schema, id string, req *UpdatePayrollEntryRequest) (*PayrollEntryResponse, error) {
	entryID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid entry ID: %w", err)
	}

	entry, err := s.repo.GetEntryByID(ctx, schema, entryID)
	if err != nil {
		return nil, err
	}

	// Get employee to recalculate
	employeeResp, err := s.employeeSvc.GetByID(ctx, schema, entry.EmployeeID.String())
	if err != nil {
		return nil, fmt.Errorf("employee not found: %w", err)
	}

	// Update fields
	if req.HoursWorked != nil {
		entry.HoursWorked = req.HoursWorked
	}
	if req.DaysWorked != nil {
		entry.DaysWorked = req.DaysWorked
	}
	if req.OvertimeHours != nil {
		entry.OvertimeHours = req.OvertimeHours
	}
	if req.Bonuses != nil {
		entry.Bonuses = *req.Bonuses
	}
	if req.Commissions != nil {
		entry.Commissions = *req.Commissions
	}
	if req.OtherDeductions != nil {
		entry.OtherDeductions = *req.OtherDeductions
	}
	if req.Notes != nil {
		entry.Notes = req.Notes
	}

	// Recalculate gross salary
	if employeeResp.Salary == nil {
		return nil, fmt.Errorf("employee has no salary configured")
	}
	baseSalary := *employeeResp.Salary
	var grossSalary float64
	
	switch employeeResp.SalaryType {
	case "monthly":
		grossSalary = baseSalary
	case "hourly":
		if entry.HoursWorked == nil {
			return nil, fmt.Errorf("hours worked required for hourly employees")
		}
		grossSalary = baseSalary * (*entry.HoursWorked)
	case "daily":
		if entry.DaysWorked == nil {
			return nil, fmt.Errorf("days worked required for daily employees")
		}
		dailyRate := baseSalary / 22.0
		grossSalary = dailyRate * (*entry.DaysWorked)
	default:
		return nil, fmt.Errorf("invalid salary type: %s", employeeResp.SalaryType)
	}
	
	entry.GrossSalary = grossSalary
	entry.Salary = grossSalary

	// Recalculate overtime pay
	var overtimePay float64
	if entry.OvertimeHours != nil && *entry.OvertimeHours > 0 {
		var hourlyRate float64
		switch employeeResp.SalaryType {
		case "monthly":
			hourlyRate = baseSalary / 160.0
		case "hourly":
			hourlyRate = baseSalary
		case "daily":
			hourlyRate = (baseSalary / 22.0) / 8.0
		default:
			return nil, fmt.Errorf("invalid salary type: %s", employeeResp.SalaryType)
		}
		overtimePay = hourlyRate * OVERTIME_RATE * (*entry.OvertimeHours)
	}
	entry.OvertimePay = overtimePay

	entry.TotalGross = entry.GrossSalary + entry.OvertimePay + entry.Bonuses + entry.Commissions
	entry.TaxDeduction = calculateISR(entry.TotalGross)
	entry.SocialSecurity = calculateTSS(entry.TotalGross)
	entry.TotalDeductions = entry.TaxDeduction + entry.SocialSecurity + entry.OtherDeductions
	entry.NetSalary = entry.TotalGross - entry.TotalDeductions

	entry.UpdatedAt = time.Now()

	if err := s.repo.UpdateEntry(ctx, schema, entry); err != nil {
		return nil, err
	}

	// Update period totals
	if err := s.updatePeriodTotals(ctx, schema, entry.PayrollPeriodID); err != nil {
		logger.Error("Failed to update period totals", "error", err)
	}

	logger.Info("Payroll entry updated", "id", entry.ID, "schema", schema)

	return entry.ToResponse(), nil
}

func (s *Service) DeleteEntry(ctx context.Context, schema, id string) error {
	entryID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid entry ID: %w", err)
	}

	entry, err := s.repo.GetEntryByID(ctx, schema, entryID)
	if err != nil {
		return err
	}

	periodID := entry.PayrollPeriodID

	if err := s.repo.DeleteEntry(ctx, schema, entryID); err != nil {
		return err
	}

	// Update period totals
	if err := s.updatePeriodTotals(ctx, schema, periodID); err != nil {
		logger.Error("Failed to update period totals", "error", err)
	}

	logger.Info("Payroll entry deleted", "id", entryID, "schema", schema)

	return nil
}

