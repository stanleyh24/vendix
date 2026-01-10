package reports

// InventoryValuationItem represents a single item in the inventory valuation report
type InventoryValuationItem struct {
	ProductID      string   `json:"product_id"`
	ProductCode    string   `json:"product_code"`
	ProductName    string   `json:"product_name"`
	ProductType    string   `json:"product_type"`
	Unit           string   `json:"unit"`
	StockQuantity  float64  `json:"stock_quantity"`
	UnitCost       *float64 `json:"unit_cost,omitempty"`
	UnitPrice      float64  `json:"unit_price"` // Precio de venta unitario
	TotalValue     float64  `json:"total_value"` // stock_quantity * unit_cost (o 0 si no hay costo)
	TotalSaleValue float64  `json:"total_sale_value"` // stock_quantity * unit_price (valor de venta total)
	SupplierID     *string  `json:"supplier_id,omitempty"`
	SupplierName   *string  `json:"supplier_name,omitempty"`
	LastUpdated    string   `json:"last_updated"`
}

// InventoryValuationReport represents the complete inventory valuation report
type InventoryValuationReport struct {
	ReportDate        string                    `json:"report_date"`
	TotalItems        int                       `json:"total_items"`
	TotalValue        float64                   `json:"total_value"`         // Suma de todos los valores (costo)
	TotalValueWithCost float64                  `json:"total_value_with_cost"` // Solo productos con costo definido
	TotalSaleValue    float64                   `json:"total_sale_value"`    // Valor de venta total del inventario (precio * cantidad)
	Items             []InventoryValuationItem  `json:"items"`
	SummaryByType     map[string]TypeSummary    `json:"summary_by_type,omitempty"` // Agrupado por tipo de producto
	SummaryBySupplier map[string]SupplierSummary `json:"summary_by_supplier,omitempty"` // Agrupado por proveedor
}

// TypeSummary represents summary grouped by product type
type TypeSummary struct {
	ProductType    string  `json:"product_type"`
	ItemCount      int     `json:"item_count"`
	TotalQuantity  float64 `json:"total_quantity"`
	TotalValue     float64 `json:"total_value"`
}

// SupplierSummary represents summary grouped by supplier
type SupplierSummary struct {
	SupplierID     string  `json:"supplier_id"`
	SupplierName   string  `json:"supplier_name"`
	ItemCount      int     `json:"item_count"`
	TotalQuantity  float64 `json:"total_quantity"`
	TotalValue     float64 `json:"total_value"`
}

// SalesReportItem represents a single sale in the sales report
type SalesReportItem struct {
	SaleID         string  `json:"sale_id"`
	SaleNumber     string  `json:"sale_number"`
	SaleDate       string  `json:"sale_date"`
	CustomerID     *string `json:"customer_id,omitempty"`
	CustomerName   string  `json:"customer_name"`
	PaymentType    string  `json:"payment_type"`
	Subtotal       float64 `json:"subtotal"`
	TaxAmount      float64 `json:"tax_amount"`
	Total          float64 `json:"total"`
	Cost           float64 `json:"cost"` // Costo de productos vendidos (COGS)
	Profit         float64 `json:"profit"` // Ganancia (total - cost)
	Status         string  `json:"status"`
}

// SalesReportSummary represents summary statistics for the sales report
type SalesReportSummary struct {
	TotalSales      int     `json:"total_sales"`
	TotalSubtotal   float64 `json:"total_subtotal"`
	TotalTax        float64 `json:"total_tax"`
	TotalAmount     float64 `json:"total_amount"`
	TotalCost       float64 `json:"total_cost"`
	TotalProfit     float64 `json:"total_profit"`
	AverageSale     float64 `json:"average_sale"`
	ProfitMargin    float64 `json:"profit_margin"` // Porcentaje de ganancia
}

// PaymentTypeSummary represents sales grouped by payment type
type PaymentTypeSummary struct {
	PaymentType string  `json:"payment_type"`
	Count       int     `json:"count"`
	Subtotal    float64 `json:"subtotal"`
	Tax         float64 `json:"tax"`
	Total       float64 `json:"total"`
}

// CustomerSummary represents sales grouped by customer
type CustomerSummary struct {
	CustomerID   string  `json:"customer_id"`
	CustomerName string  `json:"customer_name"`
	Count        int     `json:"count"`
	Subtotal     float64 `json:"subtotal"`
	Tax          float64 `json:"tax"`
	Total        float64 `json:"total"`
}

// ProductSalesSummary represents top products sold
type ProductSalesSummary struct {
	ProductID    string  `json:"product_id"`
	ProductCode  string  `json:"product_code"`
	ProductName  string  `json:"product_name"`
	Quantity     float64 `json:"quantity"`
	UnitPrice    float64 `json:"unit_price"`
	Subtotal     float64 `json:"subtotal"`
	Tax          float64 `json:"tax"`
	Total        float64 `json:"total"`
	Cost         float64 `json:"cost"`
	Profit       float64 `json:"profit"`
}

// DailySalesSummary represents sales grouped by day
type DailySalesSummary struct {
	Date      string  `json:"date"`
	Count     int     `json:"count"`
	Subtotal  float64 `json:"subtotal"`
	Tax       float64 `json:"tax"`
	Total     float64 `json:"total"`
	Cost      float64 `json:"cost"`
	Profit    float64 `json:"profit"`
}

// SalesReport represents the complete sales report
type SalesReport struct {
	ReportDate        string                `json:"report_date"`
	StartDate         string                `json:"start_date"`
	EndDate           string                `json:"end_date"`
	Summary           SalesReportSummary    `json:"summary"`
	Items             []SalesReportItem     `json:"items,omitempty"`
	ByPaymentType     []PaymentTypeSummary  `json:"by_payment_type,omitempty"`
	ByCustomer        []CustomerSummary     `json:"by_customer,omitempty"`
	TopProducts       []ProductSalesSummary `json:"top_products,omitempty"`
	DailyBreakdown    []DailySalesSummary   `json:"daily_breakdown,omitempty"`
}

// ReturnsReportItem represents a single return in the returns report
type ReturnsReportItem struct {
	ReturnID       string  `json:"return_id"`
	ReturnNumber   string  `json:"return_number"`
	ReturnDate     string  `json:"return_date"`
	ReturnType     string  `json:"return_type"` // "sale" or "invoice"
	OriginalNumber *string `json:"original_number,omitempty"` // sale_number or invoice_number
	CustomerID     *string `json:"customer_id,omitempty"`
	CustomerName   string  `json:"customer_name"`
	ReturnReason   *string `json:"return_reason,omitempty"`
	Status         string  `json:"status"`
	RefundMethod   *string `json:"refund_method,omitempty"`
	RefundStatus   *string `json:"refund_status,omitempty"`
	Subtotal       float64 `json:"subtotal"`
	TaxAmount      float64 `json:"tax_amount"`
	Total          float64 `json:"total"`
	RefundAmount   float64 `json:"refund_amount"`
}

// ReturnsReportSummary represents summary statistics for the returns report
type ReturnsReportSummary struct {
	TotalReturns      int     `json:"total_returns"`
	TotalSubtotal     float64 `json:"total_subtotal"`
	TotalTax          float64 `json:"total_tax"`
	TotalAmount       float64 `json:"total_amount"`
	TotalRefunded     float64 `json:"total_refunded"`
	AverageReturn     float64 `json:"average_return"`
	ByStatus          map[string]StatusSummary `json:"by_status,omitempty"`
}

// StatusSummary represents summary grouped by status
type StatusSummary struct {
	Status      string  `json:"status"`
	Count       int     `json:"count"`
	Subtotal    float64 `json:"subtotal"`
	Tax         float64 `json:"tax"`
	Total       float64 `json:"total"`
}

// ReturnReasonSummary represents returns grouped by reason
type ReturnReasonSummary struct {
	Reason    string  `json:"reason"`
	Count     int     `json:"count"`
	Subtotal  float64 `json:"subtotal"`
	Tax       float64 `json:"tax"`
	Total     float64 `json:"total"`
}

// RefundMethodSummary represents returns grouped by refund method
type RefundMethodSummary struct {
	RefundMethod string  `json:"refund_method"`
	Count        int     `json:"count"`
	Total        float64 `json:"total"`
}

// ReturnCustomerSummary represents returns grouped by customer
type ReturnCustomerSummary struct {
	CustomerID   string  `json:"customer_id"`
	CustomerName string  `json:"customer_name"`
	Count        int     `json:"count"`
	Subtotal     float64 `json:"subtotal"`
	Tax          float64 `json:"tax"`
	Total        float64 `json:"total"`
}

// ProductReturnsSummary represents products most frequently returned
type ProductReturnsSummary struct {
	ProductID    string  `json:"product_id"`
	ProductCode  string  `json:"product_code"`
	ProductName  string  `json:"product_name"`
	Quantity     float64 `json:"quantity"`
	Count        int     `json:"count"` // Number of returns containing this product
	Subtotal     float64 `json:"subtotal"`
	Tax          float64 `json:"tax"`
	Total        float64 `json:"total"`
}

// DailyReturnsSummary represents returns grouped by day
type DailyReturnsSummary struct {
	Date     string  `json:"date"`
	Count    int     `json:"count"`
	Subtotal float64 `json:"subtotal"`
	Tax      float64 `json:"tax"`
	Total    float64 `json:"total"`
}

// ReturnsReport represents the complete returns report
type ReturnsReport struct {
	ReportDate      string                  `json:"report_date"`
	StartDate       string                  `json:"start_date"`
	EndDate         string                  `json:"end_date"`
	Summary         ReturnsReportSummary    `json:"summary"`
	Items           []ReturnsReportItem     `json:"items,omitempty"`
	ByReason        []ReturnReasonSummary   `json:"by_reason,omitempty"`
	ByRefundMethod  []RefundMethodSummary   `json:"by_refund_method,omitempty"`
	ByCustomer      []ReturnCustomerSummary `json:"by_customer,omitempty"`
	TopProducts     []ProductReturnsSummary `json:"top_products,omitempty"`
	DailyBreakdown  []DailyReturnsSummary   `json:"daily_breakdown,omitempty"`
}

// DashboardSummary represents comprehensive dashboard KPIs and metrics
type DashboardSummary struct {
	ReportDate string `json:"report_date"`
	
	// Sales KPIs
	SalesToday      SalesKPIs `json:"sales_today"`
	SalesThisMonth  SalesKPIs `json:"sales_this_month"`
	SalesLastMonth  SalesKPIs `json:"sales_last_month"`
	
	// Invoices KPIs
	InvoicesToday      InvoiceKPIs `json:"invoices_today"`
	InvoicesThisMonth  InvoiceKPIs `json:"invoices_this_month"`
	InvoicesPending    int          `json:"invoices_pending"`
	InvoicesOverdue    int          `json:"invoices_overdue"`
	
	// Returns KPIs (NUEVO)
	ReturnsPending    int     `json:"returns_pending"`
	ReturnsThisMonth  int     `json:"returns_this_month"`
	ReturnsTotalValue float64 `json:"returns_total_value"`
	
	// Credit Notes KPIs (NUEVO)
	CreditNotesPending int     `json:"credit_notes_pending"`
	CreditNotesThisMonth int   `json:"credit_notes_this_month"`
	CreditNotesTotalValue float64 `json:"credit_notes_total_value"`
	
	// Purchases KPIs
	PurchasesThisMonth int     `json:"purchases_this_month"`
	PurchasesPending   int     `json:"purchases_pending"`
	PurchasesTotalValue float64 `json:"purchases_total_value"`
	
	// Expenses KPIs
	ExpensesToday     float64 `json:"expenses_today"`
	ExpensesThisMonth float64 `json:"expenses_this_month"`
	
	// Inventory KPIs
	InventoryTotalItems    int     `json:"inventory_total_items"`
	InventoryTotalValue    float64 `json:"inventory_total_value"`
	InventoryLowStockItems int     `json:"inventory_low_stock_items"`
	
	// Accounts Receivable/Payable
	AccountsReceivableTotal float64 `json:"accounts_receivable_total"`
	AccountsPayableTotal    float64 `json:"accounts_payable_total"`
	
	// Top Data
	TopProducts []ProductSalesSummary `json:"top_products,omitempty"`
	TopCustomers []CustomerSummary    `json:"top_customers,omitempty"`
	
	// Recent Activity
	RecentSales     []SalesReportItem     `json:"recent_sales,omitempty"`
	RecentInvoices  []InvoiceSummaryItem `json:"recent_invoices,omitempty"`
	RecentReturns   []ReturnSummaryItem  `json:"recent_returns,omitempty"`
	RecentCreditNotes []CreditNoteSummaryItem `json:"recent_credit_notes,omitempty"`
}

// SalesKPIs represents sales key performance indicators
type SalesKPIs struct {
	Count      int     `json:"count"`
	Subtotal   float64 `json:"subtotal"`
	Tax        float64 `json:"tax"`
	Total      float64 `json:"total"`
	Cost       float64 `json:"cost"`
	Profit     float64 `json:"profit"`
	ProfitMargin float64 `json:"profit_margin"`
}

// InvoiceKPIs represents invoice key performance indicators
type InvoiceKPIs struct {
	Count      int     `json:"count"`
	Subtotal   float64 `json:"subtotal"`
	Tax        float64 `json:"tax"`
	Total      float64 `json:"total"`
	PaidAmount float64 `json:"paid_amount"`
	Outstanding float64 `json:"outstanding"`
}

// InvoiceSummaryItem represents a summary invoice item for dashboard
type InvoiceSummaryItem struct {
	ID            string  `json:"id"`
	InvoiceNumber string  `json:"invoice_number"`
	CustomerName  string  `json:"customer_name"`
	Total         float64 `json:"total"`
	Status        string  `json:"status"`
	IssueDate     string  `json:"issue_date"`
}

// ReturnSummaryItem represents a summary return item for dashboard
type ReturnSummaryItem struct {
	ID            string  `json:"id"`
	ReturnNumber  string  `json:"return_number"`
	CustomerName  string  `json:"customer_name"`
	Total         float64 `json:"total"`
	Status        string  `json:"status"`
	ReturnDate    string  `json:"return_date"`
}

// CreditNoteSummaryItem represents a summary credit note item for dashboard
type CreditNoteSummaryItem struct {
	ID                string  `json:"id"`
	CreditNoteNumber  string  `json:"credit_note_number"`
	CustomerName      string  `json:"customer_name"`
	Total             float64 `json:"total"`
	Status            string  `json:"status"`
	IssueDate         string  `json:"issue_date"`
}

// TaxReport represents a comprehensive tax report (ITBIS)
type TaxReport struct {
	ReportDate        string                `json:"report_date"`
	StartDate         string                `json:"start_date"`
	EndDate           string                `json:"end_date"`
	Summary           TaxReportSummary      `json:"summary"`
	ByDocumentType    []TaxByDocumentType   `json:"by_document_type,omitempty"`
	ByNCFType         []TaxByNCFType        `json:"by_ncf_type,omitempty"`
	DailyBreakdown    []TaxDailyBreakdown   `json:"daily_breakdown,omitempty"`
	SalesTax          TaxBreakdown          `json:"sales_tax"`
	InvoicesTax       TaxBreakdown          `json:"invoices_tax"`
	PurchasesTax      TaxBreakdown          `json:"purchases_tax"`
	CreditNotesTax    TaxBreakdown          `json:"credit_notes_tax"`
	ReturnsTax        TaxBreakdown          `json:"returns_tax"`
	ByWithholding     []WithholdingSummary  `json:"by_withholding,omitempty"`
}

// TaxReportSummary represents summary statistics for the tax report
type TaxReportSummary struct {
	TotalITBISCollected  float64 `json:"total_itbis_collected"` // ITBIS cobrado (ventas + facturas)
	TotalITBISPaid       float64 `json:"total_itbis_paid"`     // ITBIS pagado (compras)
	TotalITBISNet        float64 `json:"total_itbis_net"`      // ITBIS neto a pagar (cobrado - pagado)
	TotalWithholdingTax  float64 `json:"total_withholding_tax"` // Total retenciones aplicadas (ISR, ITBIS)
	TotalNetAmount       float64 `json:"total_net_amount"`      // Total neto recibido después de retenciones
	TotalSales           int     `json:"total_sales"`
	TotalInvoices        int     `json:"total_invoices"`
	TotalPurchases       int     `json:"total_purchases"`
	TotalCreditNotes     int     `json:"total_credit_notes"`
	TotalReturns         int     `json:"total_returns"`
}

// TaxByDocumentType represents tax breakdown by document type
type TaxByDocumentType struct {
	DocumentType string  `json:"document_type"` // "sales", "invoices", "purchases", "credit_notes", "returns"
	Count        int     `json:"count"`
	Subtotal     float64 `json:"subtotal"`
	TaxAmount    float64 `json:"tax_amount"`
	Total        float64 `json:"total"`
}

// TaxByNCFType represents tax breakdown by NCF type
type TaxByNCFType struct {
	NCFType     string  `json:"ncf_type"`     // "01", "02", "03", "04", etc.
	Description string  `json:"description"`   // "Factura Fiscal", "Factura Consumidor Final", etc.
	Count       int     `json:"count"`
	Subtotal    float64 `json:"subtotal"`
	TaxAmount   float64 `json:"tax_amount"`
	Total       float64 `json:"total"`
}

// TaxDailyBreakdown represents daily tax breakdown
type TaxDailyBreakdown struct {
	Date         string  `json:"date"`
	SalesTax     float64 `json:"sales_tax"`
	InvoicesTax  float64 `json:"invoices_tax"`
	PurchasesTax float64 `json:"purchases_tax"`
	TotalTax     float64 `json:"total_tax"`
}

// TaxBreakdown represents tax breakdown for a specific document type
type TaxBreakdown struct {
	Count      int     `json:"count"`
	Subtotal   float64 `json:"subtotal"`
	TaxAmount  float64 `json:"tax_amount"`
	Total      float64 `json:"total"`
}

// WithholdingSummary represents summary of withholding taxes
type WithholdingSummary struct {
	WithholdingType string  `json:"withholding_type"` // 'isr', 'itbis'
	Count           int     `json:"count"`
	TotalAmount     float64 `json:"total_amount"`     // Total de documentos con retención
	TotalWithheld   float64 `json:"total_withheld"`   // Total retenido
}

