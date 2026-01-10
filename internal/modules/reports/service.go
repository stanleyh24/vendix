package reports

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"vendix/internal/config"
	"vendix/internal/database"
	"vendix/internal/logger"
	"vendix/internal/modules/creditnotes"
	"vendix/internal/modules/invoices"
	"vendix/internal/modules/payments"
	"vendix/internal/modules/products"
	"vendix/internal/modules/purchases"
	"vendix/internal/modules/returns"
	"vendix/internal/modules/sales"
)

type Service struct {
	db             *database.DB
	invoicesRepo   *invoices.Repository
	purchasesRepo  *purchases.Repository
	salesRepo      *sales.Repository
	returnsRepo    *returns.Repository
	creditNotesRepo *creditnotes.Repository
	paymentsRepo   *payments.Repository
	productsRepo   *products.Repository
	cfg            *config.Config
}

func NewService(db *database.DB, cfg *config.Config) *Service {
	return &Service{
		db:              db,
		invoicesRepo:    invoices.NewRepository(db),
		purchasesRepo:   purchases.NewRepository(db),
		salesRepo:       sales.NewRepository(db),
		returnsRepo:     returns.NewRepository(db),
		creditNotesRepo: creditnotes.NewRepository(db),
		paymentsRepo:    payments.NewRepository(db),
		productsRepo:     products.NewRepository(db),
		cfg:             cfg,
	}
}

// GetInventoryValuation generates an inventory valuation report
func (s *Service) GetInventoryValuation(ctx context.Context, schema string, includeZeroStock bool) (*InventoryValuationReport, error) {
	// Query to get products with their supplier information
	var query string
	if includeZeroStock {
		query = fmt.Sprintf(`
			SELECT 
				p.id::text as product_id,
				p.code as product_code,
				p.name as product_name,
				p.product_type,
				p.unit,
				p.stock_quantity,
				p.cost as unit_cost,
				p.price as unit_price,
				p.updated_at,
				s.id::text as supplier_id,
				s.name as supplier_name
			FROM %s.products p
			LEFT JOIN %s.suppliers s ON p.supplier_id = s.id
			WHERE p.is_active = true
			ORDER BY p.name ASC
		`, schema, schema)
	} else {
		query = fmt.Sprintf(`
			SELECT 
				p.id::text as product_id,
				p.code as product_code,
				p.name as product_name,
				p.product_type,
				p.unit,
				p.stock_quantity,
				p.cost as unit_cost,
				p.price as unit_price,
				p.updated_at,
				s.id::text as supplier_id,
				s.name as supplier_name
			FROM %s.products p
			LEFT JOIN %s.suppliers s ON p.supplier_id = s.id
			WHERE p.is_active = true
			  AND p.stock_quantity > 0
			ORDER BY p.name ASC
		`, schema, schema)
	}

	type row struct {
		ProductID     string     `db:"product_id"`
		ProductCode   string     `db:"product_code"`
		ProductName   string     `db:"product_name"`
		ProductType   string     `db:"product_type"`
		Unit          string     `db:"unit"`
		StockQuantity float64    `db:"stock_quantity"`
		UnitCost      *float64   `db:"unit_cost"`
		UnitPrice     float64    `db:"unit_price"`
		UpdatedAt     time.Time  `db:"updated_at"`
		SupplierID    *string    `db:"supplier_id"`
		SupplierName  *string    `db:"supplier_name"`
	}

	var rows []row
	if err := s.db.SelectContext(ctx, &rows, query); err != nil {
		logger.Error("Failed to get inventory valuation data", "error", err)
		return nil, fmt.Errorf("failed to get inventory valuation data: %w", err)
	}

	// Build report
	report := &InventoryValuationReport{
		ReportDate:        time.Now().Format("2006-01-02"),
		Items:             []InventoryValuationItem{},
		SummaryByType:     make(map[string]TypeSummary),
		SummaryBySupplier: make(map[string]SupplierSummary),
	}

	totalValue := 0.0
	totalValueWithCost := 0.0
	totalSaleValue := 0.0

	for _, r := range rows {
		// Calculate total value (cost)
		var totalValueItem float64
		if r.UnitCost != nil {
			totalValueItem = r.StockQuantity * *r.UnitCost
			totalValueWithCost += totalValueItem
		}
		// Note: totalValue includes all items, even without cost (for reporting purposes)
		totalValue += totalValueItem

		// Calculate total sale value (price * quantity)
		totalSaleValueItem := r.StockQuantity * r.UnitPrice
		totalSaleValue += totalSaleValueItem

		item := InventoryValuationItem{
			ProductID:      r.ProductID,
			ProductCode:    r.ProductCode,
			ProductName:    r.ProductName,
			ProductType:    r.ProductType,
			Unit:           r.Unit,
			StockQuantity:  r.StockQuantity,
			UnitCost:       r.UnitCost,
			UnitPrice:      r.UnitPrice,
			TotalValue:     totalValueItem,
			TotalSaleValue: totalSaleValueItem,
			SupplierID:     r.SupplierID,
			SupplierName:   r.SupplierName,
			LastUpdated:    r.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		report.Items = append(report.Items, item)

		// Update summary by type
		typeSummary := report.SummaryByType[r.ProductType]
		typeSummary.ProductType = r.ProductType
		typeSummary.ItemCount++
		typeSummary.TotalQuantity += r.StockQuantity
		typeSummary.TotalValue += totalValueItem
		report.SummaryByType[r.ProductType] = typeSummary

		// Update summary by supplier
		supplierKey := "Sin Proveedor"
		if r.SupplierID != nil && *r.SupplierID != "" {
			supplierKey = *r.SupplierID
		}

		supplierSummary := report.SummaryBySupplier[supplierKey]
		if r.SupplierID != nil && *r.SupplierID != "" {
			supplierSummary.SupplierID = *r.SupplierID
			if r.SupplierName != nil {
				supplierSummary.SupplierName = *r.SupplierName
			}
		} else {
			supplierSummary.SupplierID = ""
			supplierSummary.SupplierName = "Sin Proveedor"
		}
		supplierSummary.ItemCount++
		supplierSummary.TotalQuantity += r.StockQuantity
		supplierSummary.TotalValue += totalValueItem
		report.SummaryBySupplier[supplierKey] = supplierSummary
	}

	report.TotalItems = len(report.Items)
	report.TotalValue = totalValue
	report.TotalValueWithCost = totalValueWithCost
	report.TotalSaleValue = totalSaleValue

	logger.Info("Inventory valuation report generated",
		"schema", schema,
		"total_items", report.TotalItems,
		"total_value", report.TotalValue,
		"total_sale_value", report.TotalSaleValue,
	)

	return report, nil
}

// GetAccountsReceivableReport generates an accounts receivable aging report
func (s *Service) GetAccountsReceivableReport(ctx context.Context, schema string, asOfDate *string) (*invoices.AccountsReceivableReport, error) {
	return s.invoicesRepo.GetAccountsReceivableReport(ctx, schema, func() time.Time {
		if asOfDate != nil && *asOfDate != "" {
			date, err := time.Parse("2006-01-02", *asOfDate)
			if err == nil {
				return date
			}
		}
		return time.Now()
	}())
}

// GetAccountsPayableReport generates an accounts payable aging report
func (s *Service) GetAccountsPayableReport(ctx context.Context, schema string, asOfDate *string) (*purchases.AccountsPayableReport, error) {
	return s.purchasesRepo.GetAccountsPayableReport(ctx, schema, func() time.Time {
		if asOfDate != nil && *asOfDate != "" {
			date, err := time.Parse("2006-01-02", *asOfDate)
			if err == nil {
				return date
			}
		}
		return time.Now()
	}())
}

// GetSalesReport generates a comprehensive sales report
func (s *Service) GetSalesReport(ctx context.Context, schema string, startDate, endDate *string) (*SalesReport, error) {
	// Parse dates
	now := time.Now()
	var start, end time.Time
	
	if startDate != nil && *startDate != "" {
		parsed, err := time.Parse("2006-01-02", *startDate)
		if err == nil {
			start = parsed
		} else {
			start = now.AddDate(0, 0, -30) // Default: last 30 days
		}
	} else {
		start = now.AddDate(0, 0, -30) // Default: last 30 days
	}
	
	if endDate != nil && *endDate != "" {
		parsed, err := time.Parse("2006-01-02", *endDate)
		if err == nil {
			end = parsed
		} else {
			end = now
		}
	} else {
		end = now
	}

	// Get sales data
	rows, err := s.salesRepo.GetSalesReportData(ctx, schema, start, end)
	if err != nil {
		logger.Error("Failed to get sales report data", "error", err)
		return nil, fmt.Errorf("failed to get sales report data: %w", err)
	}

	// Build report
	report := &SalesReport{
		ReportDate:     now.Format("2006-01-02"),
		StartDate:      start.Format("2006-01-02"),
		EndDate:        end.Format("2006-01-02"),
		Summary:        SalesReportSummary{},
		Items:           []SalesReportItem{},
		ByPaymentType:   []PaymentTypeSummary{},
		ByCustomer:      []CustomerSummary{},
		TopProducts:     []ProductSalesSummary{},
		DailyBreakdown: []DailySalesSummary{},
	}

	// Maps for aggregations
	paymentTypeMap := make(map[string]*PaymentTypeSummary)
	customerMap := make(map[string]*CustomerSummary)
	dailyMap := make(map[string]*DailySalesSummary)

	var totalSubtotal, totalTax, totalAmount, totalCost float64
	totalSales := len(rows)

	// Process each sale
	for _, row := range rows {
		profit := row.Total - row.Cost
		
		totalSubtotal += row.Subtotal
		totalTax += row.TaxAmount
		totalAmount += row.Total
		totalCost += row.Cost

		// Add to items
		item := SalesReportItem{
			SaleID:       row.SaleID,
			SaleNumber:   row.SaleNumber,
			SaleDate:     row.SaleDate.Format("2006-01-02"),
			CustomerID:   row.CustomerID,
			CustomerName: row.CustomerName,
			PaymentType:  row.PaymentType,
			Subtotal:     row.Subtotal,
			TaxAmount:    row.TaxAmount,
			Total:        row.Total,
			Cost:         row.Cost,
			Profit:       profit,
			Status:       row.Status,
		}
		report.Items = append(report.Items, item)

		// Aggregate by payment type
		if paymentTypeMap[row.PaymentType] == nil {
			paymentTypeMap[row.PaymentType] = &PaymentTypeSummary{
				PaymentType: row.PaymentType,
			}
		}
		pt := paymentTypeMap[row.PaymentType]
		pt.Count++
		pt.Subtotal += row.Subtotal
		pt.Tax += row.TaxAmount
		pt.Total += row.Total

		// Aggregate by customer
		customerKey := "generic"
		if row.CustomerID != nil && *row.CustomerID != "" {
			customerKey = *row.CustomerID
		}
		if customerMap[customerKey] == nil {
			customerMap[customerKey] = &CustomerSummary{
				CustomerID:   customerKey,
				CustomerName: row.CustomerName,
			}
		}
		cs := customerMap[customerKey]
		cs.Count++
		cs.Subtotal += row.Subtotal
		cs.Tax += row.TaxAmount
		cs.Total += row.Total

		// Aggregate by day
		dayKey := row.SaleDate.Format("2006-01-02")
		if dailyMap[dayKey] == nil {
			dailyMap[dayKey] = &DailySalesSummary{
				Date: dayKey,
			}
		}
		ds := dailyMap[dayKey]
		ds.Count++
		ds.Subtotal += row.Subtotal
		ds.Tax += row.TaxAmount
		ds.Total += row.Total
		ds.Cost += row.Cost
		ds.Profit += profit
	}

	// Convert maps to slices
	for _, pt := range paymentTypeMap {
		report.ByPaymentType = append(report.ByPaymentType, *pt)
	}
	for _, cs := range customerMap {
		report.ByCustomer = append(report.ByCustomer, *cs)
	}
	for _, ds := range dailyMap {
		report.DailyBreakdown = append(report.DailyBreakdown, *ds)
	}

	// Calculate summary
	totalProfit := totalAmount - totalCost
	averageSale := 0.0
	if totalSales > 0 {
		averageSale = totalAmount / float64(totalSales)
	}
	profitMargin := 0.0
	if totalAmount > 0 {
		profitMargin = (totalProfit / totalAmount) * 100
	}

	report.Summary = SalesReportSummary{
		TotalSales:    totalSales,
		TotalSubtotal: totalSubtotal,
		TotalTax:      totalTax,
		TotalAmount:   totalAmount,
		TotalCost:     totalCost,
		TotalProfit:   totalProfit,
		AverageSale:   averageSale,
		ProfitMargin:  profitMargin,
	}

	// Get top products (need separate query)
	topProducts, err := s.getTopProducts(ctx, schema, start, end)
	if err == nil {
		report.TopProducts = topProducts
	} else {
		logger.Warn("Failed to get top products", "error", err)
	}

	logger.Info("Sales report generated",
		"schema", schema,
		"start_date", report.StartDate,
		"end_date", report.EndDate,
		"total_sales", totalSales,
		"total_amount", totalAmount,
	)

	return report, nil
}

// getTopProducts gets the top selling products for the period
func (s *Service) getTopProducts(ctx context.Context, schema string, startDate, endDate time.Time) ([]ProductSalesSummary, error) {
	query := fmt.Sprintf(`
		SELECT 
			p.id::text as product_id,
			p.code as product_code,
			p.name as product_name,
			SUM(sl.quantity) as quantity,
			AVG(sl.unit_price) as unit_price,
			SUM(sl.line_total - sl.tax_amount) as subtotal,
			SUM(sl.tax_amount) as tax,
			SUM(sl.line_total) as total,
			SUM(sl.quantity * COALESCE(p.cost, 0)) as cost
		FROM %s.sale_lines sl
		INNER JOIN %s.sales s ON sl.sale_id = s.id
		LEFT JOIN %s.products p ON sl.product_id = p.id
		WHERE DATE(s.created_at) >= $1
		  AND DATE(s.created_at) <= $2
		  AND s.status = 'completed'
		  AND sl.product_id IS NOT NULL
		GROUP BY p.id, p.code, p.name
		ORDER BY SUM(sl.line_total) DESC
		LIMIT 20
	`, schema, schema, schema)

	type row struct {
		ProductID   string  `db:"product_id"`
		ProductCode string  `db:"product_code"`
		ProductName string  `db:"product_name"`
		Quantity    float64 `db:"quantity"`
		UnitPrice   float64 `db:"unit_price"`
		Subtotal    float64 `db:"subtotal"`
		Tax         float64 `db:"tax"`
		Total       float64 `db:"total"`
		Cost        float64 `db:"cost"`
	}

	var rows []row
	if err := s.db.SelectContext(ctx, &rows, query, startDate, endDate); err != nil {
		return nil, err
	}

	products := make([]ProductSalesSummary, len(rows))
	for i, r := range rows {
		products[i] = ProductSalesSummary{
			ProductID:   r.ProductID,
			ProductCode: r.ProductCode,
			ProductName: r.ProductName,
			Quantity:    r.Quantity,
			UnitPrice:   r.UnitPrice,
			Subtotal:    r.Subtotal,
			Tax:         r.Tax,
			Total:       r.Total,
			Cost:        r.Cost,
			Profit:      r.Total - r.Cost,
		}
	}

	return products, nil
}

// GetReturnsReport generates a comprehensive returns report
func (s *Service) GetReturnsReport(ctx context.Context, schema string, startDate, endDate *string) (*ReturnsReport, error) {
	// Parse dates
	now := time.Now()
	var start, end time.Time
	
	if startDate != nil && *startDate != "" {
		parsed, err := time.Parse("2006-01-02", *startDate)
		if err == nil {
			start = parsed
		} else {
			start = now.AddDate(0, 0, -30) // Default: last 30 days
		}
	} else {
		start = now.AddDate(0, 0, -30) // Default: last 30 days
	}
	
	if endDate != nil && *endDate != "" {
		parsed, err := time.Parse("2006-01-02", *endDate)
		if err == nil {
			end = parsed
		} else {
			end = now
		}
	} else {
		end = now
	}

	// Get returns data
	rows, err := s.returnsRepo.GetReturnsReportData(ctx, schema, start, end)
	if err != nil {
		logger.Error("Failed to get returns report data", "error", err)
		return nil, fmt.Errorf("failed to get returns report data: %w", err)
	}

	// Get product returns data
	productRows, err := s.returnsRepo.GetReturnsReportProductsData(ctx, schema, start, end)
	if err != nil {
		logger.Error("Failed to get product returns data", "error", err)
		// Don't fail the report if product data fails
		productRows = []returns.ReturnsReportProductRow{}
	}

	// Build report
	report := &ReturnsReport{
		ReportDate:      now.Format("2006-01-02"),
		StartDate:       start.Format("2006-01-02"),
		EndDate:         end.Format("2006-01-02"),
		Summary:         ReturnsReportSummary{ByStatus: make(map[string]StatusSummary)},
		Items:           []ReturnsReportItem{},
		ByReason:        []ReturnReasonSummary{},
		ByRefundMethod:  []RefundMethodSummary{},
		ByCustomer:      []ReturnCustomerSummary{},
		TopProducts:     []ProductReturnsSummary{},
		DailyBreakdown:  []DailyReturnsSummary{},
	}

	// Maps for aggregations
	reasonMap := make(map[string]*ReturnReasonSummary)
	refundMethodMap := make(map[string]*RefundMethodSummary)
	customerMap := make(map[string]*ReturnCustomerSummary)
	dailyMap := make(map[string]*DailyReturnsSummary)

	totalReturns := 0
	totalSubtotal := 0.0
	totalTax := 0.0
	totalAmount := 0.0
	totalRefunded := 0.0

	// Process returns
	for _, row := range rows {
		totalReturns++
		totalSubtotal += row.Subtotal
		totalTax += row.TaxAmount
		totalAmount += row.Total
		totalRefunded += row.RefundAmount

		// Build item
		item := ReturnsReportItem{
			ReturnID:       row.ReturnID,
			ReturnNumber:   row.ReturnNumber,
			ReturnDate:     row.ReturnDate.Format("2006-01-02"),
			ReturnType:     row.ReturnType,
			OriginalNumber: row.OriginalNumber,
			CustomerID:     row.CustomerID,
			CustomerName:   row.CustomerName,
			ReturnReason:   row.ReturnReason,
			Status:         row.Status,
			RefundMethod:   row.RefundMethod,
			RefundStatus:   row.RefundStatus,
			Subtotal:       row.Subtotal,
			TaxAmount:      row.TaxAmount,
			Total:          row.Total,
			RefundAmount:   row.RefundAmount,
		}
		report.Items = append(report.Items, item)

		// Update status summary
		statusSummary := report.Summary.ByStatus[row.Status]
		statusSummary.Status = row.Status
		statusSummary.Count++
		statusSummary.Subtotal += row.Subtotal
		statusSummary.Tax += row.TaxAmount
		statusSummary.Total += row.Total
		report.Summary.ByStatus[row.Status] = statusSummary

		// Update reason summary
		reasonKey := "Sin motivo"
		if row.ReturnReason != nil && *row.ReturnReason != "" {
			reasonKey = *row.ReturnReason
		}
		if _, exists := reasonMap[reasonKey]; !exists {
			reasonMap[reasonKey] = &ReturnReasonSummary{Reason: reasonKey}
		}
		reasonMap[reasonKey].Count++
		reasonMap[reasonKey].Subtotal += row.Subtotal
		reasonMap[reasonKey].Tax += row.TaxAmount
		reasonMap[reasonKey].Total += row.Total

		// Update refund method summary
		if row.RefundMethod != nil && *row.RefundMethod != "" {
			if _, exists := refundMethodMap[*row.RefundMethod]; !exists {
				refundMethodMap[*row.RefundMethod] = &RefundMethodSummary{RefundMethod: *row.RefundMethod}
			}
			refundMethodMap[*row.RefundMethod].Count++
			refundMethodMap[*row.RefundMethod].Total += row.RefundAmount
		}

		// Update customer summary
		customerKey := "Sin cliente"
		if row.CustomerID != nil && *row.CustomerID != "" {
			customerKey = *row.CustomerID
		}
		if _, exists := customerMap[customerKey]; !exists {
			customerMap[customerKey] = &ReturnCustomerSummary{
				CustomerID:   customerKey,
				CustomerName: row.CustomerName,
			}
		}
		customerMap[customerKey].Count++
		customerMap[customerKey].Subtotal += row.Subtotal
		customerMap[customerKey].Tax += row.TaxAmount
		customerMap[customerKey].Total += row.Total

		// Update daily breakdown
		dateKey := row.ReturnDate.Format("2006-01-02")
		if _, exists := dailyMap[dateKey]; !exists {
			dailyMap[dateKey] = &DailyReturnsSummary{Date: dateKey}
		}
		dailyMap[dateKey].Count++
		dailyMap[dateKey].Subtotal += row.Subtotal
		dailyMap[dateKey].Tax += row.TaxAmount
		dailyMap[dateKey].Total += row.Total
	}

	// Convert maps to slices
	for _, summary := range reasonMap {
		report.ByReason = append(report.ByReason, *summary)
	}
	for _, summary := range refundMethodMap {
		report.ByRefundMethod = append(report.ByRefundMethod, *summary)
	}
	for _, summary := range customerMap {
		report.ByCustomer = append(report.ByCustomer, *summary)
	}
	for _, summary := range dailyMap {
		report.DailyBreakdown = append(report.DailyBreakdown, *summary)
	}

	// Sort daily breakdown by date
	for i := 0; i < len(report.DailyBreakdown); i++ {
		for j := i + 1; j < len(report.DailyBreakdown); j++ {
			if report.DailyBreakdown[i].Date > report.DailyBreakdown[j].Date {
				report.DailyBreakdown[i], report.DailyBreakdown[j] = report.DailyBreakdown[j], report.DailyBreakdown[i]
			}
		}
	}

	// Process product returns
	products := make([]ProductReturnsSummary, len(productRows))
	for i, r := range productRows {
		products[i] = ProductReturnsSummary{
			ProductID:   r.ProductID,
			ProductCode: r.ProductCode,
			ProductName: r.ProductName,
			Quantity:    r.Quantity,
			Count:       r.ReturnCount,
			Subtotal:    r.Subtotal,
			Tax:         r.Tax,
			Total:       r.Total,
		}
	}
	report.TopProducts = products

	// Calculate summary
	averageReturn := 0.0
	if totalReturns > 0 {
		averageReturn = totalAmount / float64(totalReturns)
	}

	report.Summary = ReturnsReportSummary{
		TotalReturns:  totalReturns,
		TotalSubtotal: totalSubtotal,
		TotalTax:      totalTax,
		TotalAmount:   totalAmount,
		TotalRefunded: totalRefunded,
		AverageReturn: averageReturn,
		ByStatus:      report.Summary.ByStatus,
	}

	logger.Info("Returns report generated",
		"schema", schema,
		"start_date", start.Format("2006-01-02"),
		"end_date", end.Format("2006-01-02"),
		"total_returns", totalReturns,
		"total_amount", totalAmount,
	)

	return report, nil
}

// GetDashboardSummary generates comprehensive dashboard KPIs and metrics
func (s *Service) GetDashboardSummary(ctx context.Context, schema string) (*DashboardSummary, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	startOfLastMonth := startOfMonth.AddDate(0, -1, 0)

	summary := &DashboardSummary{
		ReportDate: now.Format("2006-01-02"),
	}

	// Helper function to calculate SalesKPIs
	calculateSalesKPIs := func(start, end time.Time) (SalesKPIs, error) {
		query := fmt.Sprintf(`
			SELECT 
				COUNT(DISTINCT s.id) as count,
				COALESCE(SUM(s.subtotal), 0) as subtotal,
				COALESCE(SUM(s.tax_amount), 0) as tax,
				COALESCE(SUM(s.total), 0) as total,
				COALESCE(SUM(sl.quantity * COALESCE(p.cost, 0)), 0) as cost
			FROM %s.sales s
			LEFT JOIN %s.sale_lines sl ON s.id = sl.sale_id
			LEFT JOIN %s.products p ON sl.product_id = p.id
			WHERE s.created_at >= $1 AND s.created_at < $2
			  AND s.status = 'completed'
		`, schema, schema, schema)

		type row struct {
			Count    int     `db:"count"`
			Subtotal float64 `db:"subtotal"`
			Tax      float64 `db:"tax"`
			Total    float64 `db:"total"`
			Cost     float64 `db:"cost"`
		}

		var r row
		err := s.db.GetContext(ctx, &r, query, start, end)
		if err != nil && err != sql.ErrNoRows {
			return SalesKPIs{}, err
		}

		profit := r.Total - r.Cost
		profitMargin := 0.0
		if r.Total > 0 {
			profitMargin = (profit / r.Total) * 100
		}

		return SalesKPIs{
			Count:        r.Count,
			Subtotal:     r.Subtotal,
			Tax:          r.Tax,
			Total:        r.Total,
			Cost:         r.Cost,
			Profit:       profit,
			ProfitMargin: profitMargin,
		}, nil
	}

	// Helper function to calculate InvoiceKPIs
	calculateInvoiceKPIs := func(start, end time.Time) (InvoiceKPIs, error) {
		query := fmt.Sprintf(`
			SELECT 
				COUNT(*) as count,
				COALESCE(SUM(subtotal), 0) as subtotal,
				COALESCE(SUM(tax_amount), 0) as tax,
				COALESCE(SUM(total), 0) as total,
				COALESCE(SUM(paid_amount), 0) as paid_amount
			FROM %s.invoices
			WHERE issue_date >= $1 AND issue_date < $2
		`, schema)

		type row struct {
			Count       int     `db:"count"`
			Subtotal    float64 `db:"subtotal"`
			Tax         float64 `db:"tax"`
			Total       float64 `db:"total"`
			PaidAmount  float64 `db:"paid_amount"`
		}

		var r row
		err := s.db.GetContext(ctx, &r, query, start, end)
		if err != nil && err != sql.ErrNoRows {
			return InvoiceKPIs{}, err
		}

		return InvoiceKPIs{
			Count:       r.Count,
			Subtotal:   r.Subtotal,
			Tax:        r.Tax,
			Total:      r.Total,
			PaidAmount: r.PaidAmount,
			Outstanding: r.Total - r.PaidAmount,
		}, nil
	}

	// Get Sales KPIs
	var err error
	summary.SalesToday, err = calculateSalesKPIs(today, today.AddDate(0, 0, 1))
	if err != nil {
		logger.Error("Failed to get sales today", "error", err)
	}

	summary.SalesThisMonth, err = calculateSalesKPIs(startOfMonth, now)
	if err != nil {
		logger.Error("Failed to get sales this month", "error", err)
	}

	summary.SalesLastMonth, err = calculateSalesKPIs(startOfLastMonth, startOfMonth)
	if err != nil {
		logger.Error("Failed to get sales last month", "error", err)
	}

	// Get Invoice KPIs
	summary.InvoicesToday, err = calculateInvoiceKPIs(today, today.AddDate(0, 0, 1))
	if err != nil {
		logger.Error("Failed to get invoices today", "error", err)
	}

	summary.InvoicesThisMonth, err = calculateInvoiceKPIs(startOfMonth, now)
	if err != nil {
		logger.Error("Failed to get invoices this month", "error", err)
	}

	// Get pending and overdue invoices
	pendingQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %s.invoices WHERE status = 'pending'`, schema)
	overdueQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %s.invoices WHERE status = 'overdue'`, schema)
	s.db.GetContext(ctx, &summary.InvoicesPending, pendingQuery)
	s.db.GetContext(ctx, &summary.InvoicesOverdue, overdueQuery)

	// Get Returns KPIs
	returnsQuery := fmt.Sprintf(`
		SELECT 
			COUNT(*) FILTER (WHERE status = 'pending') as pending,
			COUNT(*) FILTER (WHERE return_date >= $1) as this_month,
			COALESCE(SUM(total) FILTER (WHERE return_date >= $1), 0) as total_value
		FROM %s.returns
	`, schema)
	type returnsRow struct {
		Pending    int     `db:"pending"`
		ThisMonth  int     `db:"this_month"`
		TotalValue float64 `db:"total_value"`
	}
	var returnsData returnsRow
	if err := s.db.GetContext(ctx, &returnsData, returnsQuery, startOfMonth); err != nil && err != sql.ErrNoRows {
		logger.Error("Failed to get returns data", "error", err)
	} else {
		summary.ReturnsPending = returnsData.Pending
		summary.ReturnsThisMonth = returnsData.ThisMonth
		summary.ReturnsTotalValue = returnsData.TotalValue
	}

	// Get Credit Notes KPIs
	creditNotesQuery := fmt.Sprintf(`
		SELECT 
			COUNT(*) FILTER (WHERE status = 'pending') as pending,
			COUNT(*) FILTER (WHERE issue_date >= $1) as this_month,
			COALESCE(SUM(total) FILTER (WHERE issue_date >= $1), 0) as total_value
		FROM %s.credit_notes
	`, schema)
	type creditNotesRow struct {
		Pending    int     `db:"pending"`
		ThisMonth  int     `db:"this_month"`
		TotalValue float64 `db:"total_value"`
	}
	var creditNotesData creditNotesRow
	if err := s.db.GetContext(ctx, &creditNotesData, creditNotesQuery, startOfMonth); err != nil && err != sql.ErrNoRows {
		logger.Error("Failed to get credit notes data", "error", err)
	} else {
		summary.CreditNotesPending = creditNotesData.Pending
		summary.CreditNotesThisMonth = creditNotesData.ThisMonth
		summary.CreditNotesTotalValue = creditNotesData.TotalValue
	}

	// Get Purchases KPIs
	purchasesQuery := fmt.Sprintf(`
		SELECT 
			COUNT(*) FILTER (WHERE purchase_date >= $1) as this_month,
			COUNT(*) FILTER (WHERE payment_status IN ('pending', 'partial')) as pending,
			COALESCE(SUM(total) FILTER (WHERE purchase_date >= $1), 0) as total_value
		FROM %s.purchases
	`, schema)
	type purchasesRow struct {
		ThisMonth  int     `db:"this_month"`
		Pending    int     `db:"pending"`
		TotalValue float64 `db:"total_value"`
	}
	var purchasesData purchasesRow
	if err := s.db.GetContext(ctx, &purchasesData, purchasesQuery, startOfMonth); err != nil && err != sql.ErrNoRows {
		logger.Error("Failed to get purchases data", "error", err)
	} else {
		summary.PurchasesThisMonth = purchasesData.ThisMonth
		summary.PurchasesPending = purchasesData.Pending
		summary.PurchasesTotalValue = purchasesData.TotalValue
	}

	// Get Expenses KPIs
	expensesTodayQuery := fmt.Sprintf(`
		SELECT COALESCE(SUM(amount), 0) 
		FROM %s.payments 
		WHERE payment_date >= $1 AND payment_date < $2 AND status = 'completed'
	`, schema)
	expensesMonthQuery := fmt.Sprintf(`
		SELECT COALESCE(SUM(amount), 0) 
		FROM %s.payments 
		WHERE payment_date >= $1 AND status = 'completed'
	`, schema)
	s.db.GetContext(ctx, &summary.ExpensesToday, expensesTodayQuery, today, today.AddDate(0, 0, 1))
	s.db.GetContext(ctx, &summary.ExpensesThisMonth, expensesMonthQuery, startOfMonth)

	// Get Inventory KPIs
	inventoryQuery := fmt.Sprintf(`
		SELECT 
			COUNT(*) FILTER (WHERE is_active = true) as total_items,
			COUNT(*) FILTER (WHERE is_active = true AND stock_quantity <= 10 AND stock_quantity > 0) as low_stock,
			COALESCE(SUM(stock_quantity * COALESCE(cost, 0)) FILTER (WHERE is_active = true), 0) as total_value
		FROM %s.products
	`, schema)
	type inventoryRow struct {
		TotalItems    int     `db:"total_items"`
		LowStockItems int     `db:"low_stock"`
		TotalValue    float64 `db:"total_value"`
	}
	var inventoryData inventoryRow
	if err := s.db.GetContext(ctx, &inventoryData, inventoryQuery); err != nil && err != sql.ErrNoRows {
		logger.Error("Failed to get inventory data", "error", err)
	} else {
		summary.InventoryTotalItems = inventoryData.TotalItems
		summary.InventoryLowStockItems = inventoryData.LowStockItems
		summary.InventoryTotalValue = inventoryData.TotalValue
	}

	// Get Accounts Receivable and Payable
	arQuery := fmt.Sprintf(`
		SELECT COALESCE(SUM(total - paid_amount), 0) 
		FROM %s.invoices 
		WHERE status IN ('sent', 'pending', 'overdue') AND (total - paid_amount) > 0
	`, schema)
	apQuery := fmt.Sprintf(`
		SELECT COALESCE(SUM(p.total - COALESCE(pp.paid_amount, 0)), 0)
		FROM %s.purchases p
		LEFT JOIN (
			SELECT purchase_id, SUM(amount) as paid_amount
			FROM %s.supplier_payment_allocations
			GROUP BY purchase_id
		) pp ON p.id = pp.purchase_id
		WHERE p.payment_status IN ('pending', 'partial')
		  AND (p.total - COALESCE(pp.paid_amount, 0)) > 0
	`, schema, schema)
	s.db.GetContext(ctx, &summary.AccountsReceivableTotal, arQuery)
	s.db.GetContext(ctx, &summary.AccountsPayableTotal, apQuery)

	// Get Top Products (last 30 days)
	topProducts, err := s.getTopProducts(ctx, schema, now.AddDate(0, 0, -30), now)
	if err == nil {
		// Limit to top 5
		if len(topProducts) > 5 {
			summary.TopProducts = topProducts[:5]
		} else {
			summary.TopProducts = topProducts
		}
	}

	// Get Top Customers (last 30 days)
	topCustomers, err := s.getTopCustomers(ctx, schema, now.AddDate(0, 0, -30), now, 5)
	if err == nil {
		summary.TopCustomers = topCustomers
	}

	// Get Recent Sales (last 10)
	recentSales, err := s.getRecentSales(ctx, schema, 10)
	if err == nil {
		summary.RecentSales = recentSales
	}

	// Get Recent Invoices (last 10)
	recentInvoices, err := s.getRecentInvoices(ctx, schema, 10)
	if err == nil {
		summary.RecentInvoices = recentInvoices
	}

	// Get Recent Returns (last 10)
	recentReturns, err := s.getRecentReturns(ctx, schema, 10)
	if err == nil {
		summary.RecentReturns = recentReturns
	}

	// Get Recent Credit Notes (last 10)
	recentCreditNotes, err := s.getRecentCreditNotes(ctx, schema, 10)
	if err == nil {
		summary.RecentCreditNotes = recentCreditNotes
	}

	return summary, nil
}

// Helper methods for dashboard data
func (s *Service) getTopCustomers(ctx context.Context, schema string, start, end time.Time, limit int) ([]CustomerSummary, error) {
	query := fmt.Sprintf(`
		SELECT 
			c.id::text as customer_id,
			COALESCE(c.name, 'Cliente Genérico') as customer_name,
			COUNT(DISTINCT s.id) as count,
			COALESCE(SUM(s.subtotal), 0) as subtotal,
			COALESCE(SUM(s.tax_amount), 0) as tax,
			COALESCE(SUM(s.total), 0) as total
		FROM %s.sales s
		LEFT JOIN %s.customers c ON s.customer_id = c.id
		WHERE s.created_at >= $1 AND s.created_at <= $2
		  AND s.status = 'completed'
		GROUP BY c.id, c.name
		ORDER BY total DESC
		LIMIT $3
	`, schema, schema)

	var customers []CustomerSummary
	err := s.db.SelectContext(ctx, &customers, query, start, end, limit)
	return customers, err
}

func (s *Service) getRecentSales(ctx context.Context, schema string, limit int) ([]SalesReportItem, error) {
	query := fmt.Sprintf(`
		SELECT 
			s.id::text as sale_id,
			s.sale_number,
			DATE(s.created_at)::text as sale_date,
			s.customer_id::text as customer_id,
			COALESCE(c.name, 'Cliente Genérico') as customer_name,
			s.payment_type,
			s.subtotal,
			s.tax_amount,
			s.total,
			COALESCE(SUM(sl.quantity * COALESCE(p.cost, 0)), 0) as cost,
			s.status
		FROM %s.sales s
		LEFT JOIN %s.customers c ON s.customer_id = c.id
		LEFT JOIN %s.sale_lines sl ON s.id = sl.sale_id
		LEFT JOIN %s.products p ON sl.product_id = p.id
		WHERE s.status = 'completed'
		GROUP BY s.id, s.sale_number, s.created_at, s.customer_id, c.name, s.payment_type, s.subtotal, s.tax_amount, s.total, s.status
		ORDER BY s.created_at DESC
		LIMIT $1
	`, schema, schema, schema, schema)

	var sales []SalesReportItem
	err := s.db.SelectContext(ctx, &sales, query, limit)
	if err != nil {
		return nil, err
	}

	// Calculate profit for each sale
	for i := range sales {
		sales[i].Profit = sales[i].Total - sales[i].Cost
	}

	return sales, nil
}

func (s *Service) getRecentInvoices(ctx context.Context, schema string, limit int) ([]InvoiceSummaryItem, error) {
	query := fmt.Sprintf(`
		SELECT 
			i.id::text as id,
			i.invoice_number,
			COALESCE(c.name, 'Cliente Genérico') as customer_name,
			i.total,
			i.status,
			i.issue_date::text as issue_date
		FROM %s.invoices i
		LEFT JOIN %s.customers c ON i.customer_id = c.id
		ORDER BY i.created_at DESC
		LIMIT $1
	`, schema, schema)

	var invoices []InvoiceSummaryItem
	err := s.db.SelectContext(ctx, &invoices, query, limit)
	return invoices, err
}

func (s *Service) getRecentReturns(ctx context.Context, schema string, limit int) ([]ReturnSummaryItem, error) {
	query := fmt.Sprintf(`
		SELECT 
			r.id::text as id,
			r.return_number,
			COALESCE(c.name, 'Cliente Genérico') as customer_name,
			r.total,
			r.status,
			r.return_date::text as return_date
		FROM %s.returns r
		LEFT JOIN %s.customers c ON r.customer_id = c.id
		ORDER BY r.created_at DESC
		LIMIT $1
	`, schema, schema)

	var returns []ReturnSummaryItem
	err := s.db.SelectContext(ctx, &returns, query, limit)
	return returns, err
}

func (s *Service) getRecentCreditNotes(ctx context.Context, schema string, limit int) ([]CreditNoteSummaryItem, error) {
	query := fmt.Sprintf(`
		SELECT 
			cn.id::text as id,
			cn.credit_note_number,
			COALESCE(c.name, 'Cliente Genérico') as customer_name,
			cn.total,
			cn.status,
			cn.issue_date::text as issue_date
		FROM %s.credit_notes cn
		LEFT JOIN %s.customers c ON cn.customer_id = c.id
		ORDER BY cn.created_at DESC
		LIMIT $1
	`, schema, schema)

	var creditNotes []CreditNoteSummaryItem
	err := s.db.SelectContext(ctx, &creditNotes, query, limit)
	return creditNotes, err
}

// GetTaxReport generates a comprehensive tax report (ITBIS)
func (s *Service) GetTaxReport(ctx context.Context, schema string, startDate, endDate *string) (*TaxReport, error) {
	// Parse dates
	now := time.Now()
	var start, end time.Time
	
	if startDate != nil && *startDate != "" {
		parsed, err := time.Parse("2006-01-02", *startDate)
		if err == nil {
			start = parsed
		} else {
			start = now.AddDate(0, 0, -30) // Default: last 30 days
		}
	} else {
		start = now.AddDate(0, 0, -30) // Default: last 30 days
	}
	
	if endDate != nil && *endDate != "" {
		parsed, err := time.Parse("2006-01-02", *endDate)
		if err == nil {
			end = parsed
		} else {
			end = now
		}
	} else {
		end = now
	}

	report := &TaxReport{
		ReportDate:      now.Format("2006-01-02"),
		StartDate:       start.Format("2006-01-02"),
		EndDate:         end.Format("2006-01-02"),
		Summary:         TaxReportSummary{},
		ByDocumentType:  []TaxByDocumentType{},
		ByNCFType:       []TaxByNCFType{},
		DailyBreakdown:  []TaxDailyBreakdown{},
		SalesTax:        TaxBreakdown{},
		InvoicesTax:     TaxBreakdown{},
		PurchasesTax:    TaxBreakdown{},
		CreditNotesTax:  TaxBreakdown{},
		ReturnsTax:      TaxBreakdown{},
	}

	// Get Sales Tax (ITBIS cobrado)
	salesQuery := fmt.Sprintf(`
		SELECT 
			COUNT(*) as count,
			COALESCE(SUM(subtotal), 0) as subtotal,
			COALESCE(SUM(tax_amount), 0) as tax_amount,
			COALESCE(SUM(total), 0) as total
		FROM %s.sales
		WHERE DATE(created_at) >= $1 AND DATE(created_at) <= $2
		  AND status = 'completed'
	`, schema)
	type salesRow struct {
		Count     int     `db:"count"`
		Subtotal  float64 `db:"subtotal"`
		TaxAmount float64 `db:"tax_amount"`
		Total     float64 `db:"total"`
	}
	var salesData salesRow
	if err := s.db.GetContext(ctx, &salesData, salesQuery, start, end); err != nil && err != sql.ErrNoRows {
		logger.Error("Failed to get sales tax data", "error", err)
	} else {
		report.SalesTax = TaxBreakdown{
			Count:     salesData.Count,
			Subtotal:  salesData.Subtotal,
			TaxAmount: salesData.TaxAmount,
			Total:     salesData.Total,
		}
		report.Summary.TotalSales = salesData.Count
	}

	// Get Invoices Tax (ITBIS cobrado)
	invoicesQuery := fmt.Sprintf(`
		SELECT 
			COUNT(*) as count,
			COALESCE(SUM(subtotal), 0) as subtotal,
			COALESCE(SUM(tax_amount), 0) as tax_amount,
			COALESCE(SUM(total), 0) as total
		FROM %s.invoices
		WHERE DATE(issue_date) >= $1 AND DATE(issue_date) <= $2
	`, schema)
	type invoicesRow struct {
		Count     int     `db:"count"`
		Subtotal  float64 `db:"subtotal"`
		TaxAmount float64 `db:"tax_amount"`
		Total     float64 `db:"total"`
	}
	var invoicesData invoicesRow
	if err := s.db.GetContext(ctx, &invoicesData, invoicesQuery, start, end); err != nil && err != sql.ErrNoRows {
		logger.Error("Failed to get invoices tax data", "error", err)
	} else {
		report.InvoicesTax = TaxBreakdown{
			Count:     invoicesData.Count,
			Subtotal:  invoicesData.Subtotal,
			TaxAmount: invoicesData.TaxAmount,
			Total:     invoicesData.Total,
		}
		report.Summary.TotalInvoices = invoicesData.Count
	}

	// Get Purchases Tax (ITBIS pagado - crédito fiscal)
	purchasesQuery := fmt.Sprintf(`
		SELECT 
			COUNT(*) as count,
			COALESCE(SUM(subtotal), 0) as subtotal,
			COALESCE(SUM(tax_amount), 0) as tax_amount,
			COALESCE(SUM(total), 0) as total
		FROM %s.purchases
		WHERE DATE(purchase_date) >= $1 AND DATE(purchase_date) <= $2
	`, schema)
	type purchasesRow struct {
		Count     int     `db:"count"`
		Subtotal  float64 `db:"subtotal"`
		TaxAmount float64 `db:"tax_amount"`
		Total     float64 `db:"total"`
	}
	var purchasesData purchasesRow
	if err := s.db.GetContext(ctx, &purchasesData, purchasesQuery, start, end); err != nil && err != sql.ErrNoRows {
		logger.Error("Failed to get purchases tax data", "error", err)
	} else {
		report.PurchasesTax = TaxBreakdown{
			Count:     purchasesData.Count,
			Subtotal:  purchasesData.Subtotal,
			TaxAmount: purchasesData.TaxAmount,
			Total:     purchasesData.Total,
		}
		report.Summary.TotalPurchases = purchasesData.Count
	}

	// Get Credit Notes Tax (reduce ITBIS cobrado)
	creditNotesQuery := fmt.Sprintf(`
		SELECT 
			COUNT(*) as count,
			COALESCE(SUM(subtotal), 0) as subtotal,
			COALESCE(SUM(tax_amount), 0) as tax_amount,
			COALESCE(SUM(total), 0) as total
		FROM %s.credit_notes
		WHERE DATE(issue_date) >= $1 AND DATE(issue_date) <= $2
	`, schema)
	type creditNotesRow struct {
		Count     int     `db:"count"`
		Subtotal  float64 `db:"subtotal"`
		TaxAmount float64 `db:"tax_amount"`
		Total     float64 `db:"total"`
	}
	var creditNotesData creditNotesRow
	if err := s.db.GetContext(ctx, &creditNotesData, creditNotesQuery, start, end); err != nil && err != sql.ErrNoRows {
		logger.Error("Failed to get credit notes tax data", "error", err)
	} else {
		report.CreditNotesTax = TaxBreakdown{
			Count:     creditNotesData.Count,
			Subtotal:  creditNotesData.Subtotal,
			TaxAmount: creditNotesData.TaxAmount,
			Total:     creditNotesData.Total,
		}
		report.Summary.TotalCreditNotes = creditNotesData.Count
	}

	// Get Returns Tax (reduce ITBIS cobrado)
	returnsQuery := fmt.Sprintf(`
		SELECT 
			COUNT(*) as count,
			COALESCE(SUM(subtotal), 0) as subtotal,
			COALESCE(SUM(tax_amount), 0) as tax_amount,
			COALESCE(SUM(total), 0) as total
		FROM %s.returns
		WHERE DATE(return_date) >= $1 AND DATE(return_date) <= $2
	`, schema)
	type returnsRow struct {
		Count     int     `db:"count"`
		Subtotal  float64 `db:"subtotal"`
		TaxAmount float64 `db:"tax_amount"`
		Total     float64 `db:"total"`
	}
	var returnsData returnsRow
	if err := s.db.GetContext(ctx, &returnsData, returnsQuery, start, end); err != nil && err != sql.ErrNoRows {
		logger.Error("Failed to get returns tax data", "error", err)
	} else {
		report.ReturnsTax = TaxBreakdown{
			Count:     returnsData.Count,
			Subtotal:  returnsData.Subtotal,
			TaxAmount: returnsData.TaxAmount,
			Total:     returnsData.Total,
		}
		report.Summary.TotalReturns = returnsData.Count
	}

	// Calculate summary
	// ITBIS cobrado = ventas + facturas
	itbisCollected := report.SalesTax.TaxAmount + report.InvoicesTax.TaxAmount
	// ITBIS pagado = compras
	itbisPaid := report.PurchasesTax.TaxAmount
	// Reducciones = notas de crédito + devoluciones
	itbisReductions := report.CreditNotesTax.TaxAmount + report.ReturnsTax.TaxAmount
	// ITBIS neto = cobrado - pagado - reducciones
	itbisNet := itbisCollected - itbisPaid - itbisReductions

	report.Summary.TotalITBISCollected = itbisCollected
	report.Summary.TotalITBISPaid = itbisPaid
	report.Summary.TotalITBISNet = itbisNet

	// Get Withholding Tax (Retenciones)
	withholdingQuery := fmt.Sprintf(`
		SELECT 
			COALESCE(withholding_tax_type, 'unknown') as withholding_tax_type,
			COUNT(*) as count,
			COALESCE(SUM(total), 0) as total_amount,
			COALESCE(SUM(withholding_tax_amount), 0) as total_withheld
		FROM %s.invoices
		WHERE DATE(issue_date) >= $1 AND DATE(issue_date) <= $2
		  AND withholding_tax_amount IS NOT NULL
		  AND withholding_tax_amount > 0
		GROUP BY withholding_tax_type
	`, schema)
	type withholdingRow struct {
		WithholdingType string  `db:"withholding_tax_type"`
		Count           int     `db:"count"`
		TotalAmount     float64 `db:"total_amount"`
		TotalWithheld   float64 `db:"total_withheld"`
	}
	var withholdingRows []withholdingRow
	totalWithholding := 0.0
	totalNetAmount := 0.0
	if err := s.db.SelectContext(ctx, &withholdingRows, withholdingQuery, start, end); err == nil {
		report.ByWithholding = make([]WithholdingSummary, len(withholdingRows))
		for i, row := range withholdingRows {
			report.ByWithholding[i] = WithholdingSummary{
				WithholdingType: row.WithholdingType,
				Count:           row.Count,
				TotalAmount:     row.TotalAmount,
				TotalWithheld:   row.TotalWithheld,
			}
			totalWithholding += row.TotalWithheld
		}
	}

	// Calcular total neto recibido (después de retenciones)
	netAmountQuery := fmt.Sprintf(`
		SELECT COALESCE(SUM(net_amount), 0)
		FROM %s.invoices
		WHERE DATE(issue_date) >= $1 AND DATE(issue_date) <= $2
	`, schema)
	if err := s.db.GetContext(ctx, &totalNetAmount, netAmountQuery, start, end); err != nil && err != sql.ErrNoRows {
		logger.Error("Failed to get total net amount", "error", err)
	}

	report.Summary.TotalWithholdingTax = totalWithholding
	report.Summary.TotalNetAmount = totalNetAmount

	// Build by document type
	report.ByDocumentType = []TaxByDocumentType{
		{
			DocumentType: "sales",
			Count:        report.SalesTax.Count,
			Subtotal:     report.SalesTax.Subtotal,
			TaxAmount:    report.SalesTax.TaxAmount,
			Total:        report.SalesTax.Total,
		},
		{
			DocumentType: "invoices",
			Count:        report.InvoicesTax.Count,
			Subtotal:     report.InvoicesTax.Subtotal,
			TaxAmount:    report.InvoicesTax.TaxAmount,
			Total:        report.InvoicesTax.Total,
		},
		{
			DocumentType: "purchases",
			Count:        report.PurchasesTax.Count,
			Subtotal:     report.PurchasesTax.Subtotal,
			TaxAmount:    report.PurchasesTax.TaxAmount,
			Total:        report.PurchasesTax.Total,
		},
		{
			DocumentType: "credit_notes",
			Count:        report.CreditNotesTax.Count,
			Subtotal:     report.CreditNotesTax.Subtotal,
			TaxAmount:    report.CreditNotesTax.TaxAmount,
			Total:        report.CreditNotesTax.Total,
		},
		{
			DocumentType: "returns",
			Count:        report.ReturnsTax.Count,
			Subtotal:     report.ReturnsTax.Subtotal,
			TaxAmount:    report.ReturnsTax.TaxAmount,
			Total:        report.ReturnsTax.Total,
		},
	}

	// Get tax by NCF type
	ncfQuery := fmt.Sprintf(`
		SELECT 
			COALESCE(ncf_type, 'N/A') as ncf_type,
			COUNT(*) as count,
			COALESCE(SUM(subtotal), 0) as subtotal,
			COALESCE(SUM(tax_amount), 0) as tax_amount,
			COALESCE(SUM(total), 0) as total
		FROM %s.invoices
		WHERE DATE(issue_date) >= $1 AND DATE(issue_date) <= $2
		GROUP BY ncf_type
		ORDER BY ncf_type
	`, schema)
	type ncfRow struct {
		NCFType   string  `db:"ncf_type"`
		Count     int     `db:"count"`
		Subtotal  float64 `db:"subtotal"`
		TaxAmount float64 `db:"tax_amount"`
		Total     float64 `db:"total"`
	}
	var ncfRows []ncfRow
	if err := s.db.SelectContext(ctx, &ncfRows, ncfQuery, start, end); err == nil {
		ncfDescriptions := map[string]string{
			"01": "Factura Fiscal",
			"02": "Factura Consumidor Final",
			"03": "Nota de Débito",
			"04": "Nota de Crédito",
			"11": "Factura de Exportación",
			"12": "Factura Gubernamental",
			"13": "Factura Especial",
			"14": "Factura Regímenes Especiales",
			"15": "Factura Gubernamental",
		}
		report.ByNCFType = make([]TaxByNCFType, len(ncfRows))
		for i, row := range ncfRows {
			report.ByNCFType[i] = TaxByNCFType{
				NCFType:     row.NCFType,
				Description: ncfDescriptions[row.NCFType],
				Count:       row.Count,
				Subtotal:    row.Subtotal,
				TaxAmount:   row.TaxAmount,
				Total:       row.Total,
			}
		}
	}

	// Get daily breakdown
	dailyQuery := fmt.Sprintf(`
		SELECT 
			DATE(created_at)::text as date,
			COALESCE(SUM(CASE WHEN source = 'sales' THEN tax_amount ELSE 0 END), 0) as sales_tax,
			COALESCE(SUM(CASE WHEN source = 'invoices' THEN tax_amount ELSE 0 END), 0) as invoices_tax,
			COALESCE(SUM(CASE WHEN source = 'purchases' THEN tax_amount ELSE 0 END), 0) as purchases_tax
		FROM (
			SELECT created_at, tax_amount, 'sales' as source
			FROM %s.sales
			WHERE DATE(created_at) >= $1 AND DATE(created_at) <= $2 AND status = 'completed'
			UNION ALL
			SELECT issue_date as created_at, tax_amount, 'invoices' as source
			FROM %s.invoices
			WHERE DATE(issue_date) >= $1 AND DATE(issue_date) <= $2
			UNION ALL
			SELECT purchase_date as created_at, tax_amount, 'purchases' as source
			FROM %s.purchases
			WHERE DATE(purchase_date) >= $1 AND DATE(purchase_date) <= $2
		) combined
		GROUP BY DATE(created_at)
		ORDER BY DATE(created_at)
	`, schema, schema, schema)
	type dailyRow struct {
		Date         string  `db:"date"`
		SalesTax     float64 `db:"sales_tax"`
		InvoicesTax  float64 `db:"invoices_tax"`
		PurchasesTax float64 `db:"purchases_tax"`
	}
	var dailyRows []dailyRow
	if err := s.db.SelectContext(ctx, &dailyRows, dailyQuery, start, end); err == nil {
		report.DailyBreakdown = make([]TaxDailyBreakdown, len(dailyRows))
		for i, row := range dailyRows {
			report.DailyBreakdown[i] = TaxDailyBreakdown{
				Date:         row.Date,
				SalesTax:     row.SalesTax,
				InvoicesTax:  row.InvoicesTax,
				PurchasesTax: row.PurchasesTax,
				TotalTax:     row.SalesTax + row.InvoicesTax - row.PurchasesTax, // Neto del día
			}
		}
	}

	logger.Info("Tax report generated",
		"schema", schema,
		"start_date", report.StartDate,
		"end_date", report.EndDate,
		"itbis_collected", itbisCollected,
		"itbis_paid", itbisPaid,
		"itbis_net", itbisNet,
	)

	return report, nil
}
