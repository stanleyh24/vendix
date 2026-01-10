import { useState, useEffect } from 'react';
import { 
  DollarSign, TrendingUp, FileText, ShoppingCart, Package, 
  Users, ArrowUpRight, ArrowDownRight, AlertCircle, CheckCircle,
  RotateCcw, FileX, Receipt, CreditCard, TrendingDown, Activity,
  Calendar, Clock, BarChart3
} from 'lucide-react';
import api from '../lib/api';
import Alert from '../components/Alert';

export default function Dashboard() {
  const [dashboardData, setDashboardData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [alert, setAlert] = useState(null);

  useEffect(() => {
    loadDashboardData();
  }, []);

  const loadDashboardData = async () => {
    setLoading(true);
    try {
      const response = await api.get('/reports/dashboard');
      setDashboardData(response.data);
    } catch (error) {
      console.error('Error loading dashboard:', error);
      showAlert('error', 'Error', 'No se pudieron cargar los datos del dashboard');
    } finally {
      setLoading(false);
    }
  };

  const showAlert = (type, title, message) => {
    setAlert({ type, title, message });
    setTimeout(() => setAlert(null), 5000);
  };

  const formatCurrency = (value) => {
    return new Intl.NumberFormat('es-DO', {
      style: 'currency',
      currency: 'DOP',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(value);
  };

  const formatNumber = (value) => {
    return new Intl.NumberFormat('es-DO').format(value);
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-[#FF6B00]"></div>
      </div>
    );
  }

  if (!dashboardData) {
    return (
      <div className="p-6">
        <Alert type="error" title="Error" message="No se pudieron cargar los datos del dashboard" />
      </div>
    );
  }

  const data = dashboardData;

  // Calcular porcentaje de cambio para ventas
  const salesChange = data.sales_last_month.total > 0
    ? ((data.sales_this_month.total - data.sales_last_month.total) / data.sales_last_month.total) * 100
    : 0;

  // Calcular porcentaje de cambio para ingresos
  const revenueChange = data.sales_last_month.total > 0
    ? ((data.sales_this_month.total - data.sales_last_month.total) / data.sales_last_month.total) * 100
    : 0;

  return (
    <div className="min-h-screen bg-gray-50">
      {alert && (
        <div className="p-4">
          <Alert type={alert.type} title={alert.title} message={alert.message} />
        </div>
      )}

      {/* Header */}
      <div className="bg-white border-b border-gray-200 px-6 py-4">
        <h1 className="text-2xl font-bold text-gray-900">Dashboard</h1>
        <p className="text-sm text-gray-600 mt-1">Resumen general del negocio</p>
      </div>

      <div className="p-6 space-y-6">
        {/* KPIs Principales - Ventas e Ingresos */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {/* Ventas Hoy */}
          <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-sm font-medium text-gray-600">Ventas Hoy</h3>
              <ShoppingCart className="w-5 h-5 text-blue-500" />
            </div>
            <div className="text-2xl font-bold text-gray-900 mb-2">
              {formatCurrency(data.sales_today.total)}
            </div>
            <div className="flex items-center text-sm">
              <span className="text-gray-600">{data.sales_today.count} transacciones</span>
            </div>
          </div>

          {/* Ventas del Mes */}
          <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-sm font-medium text-gray-600">Ventas del Mes</h3>
              <TrendingUp className="w-5 h-5 text-green-500" />
            </div>
            <div className="text-2xl font-bold text-gray-900 mb-2">
              {formatCurrency(data.sales_this_month.total)}
            </div>
            <div className="flex items-center text-sm">
              {salesChange >= 0 ? (
                <ArrowUpRight className="w-4 h-4 text-green-500 mr-1" />
              ) : (
                <ArrowDownRight className="w-4 h-4 text-red-500 mr-1" />
              )}
              <span className={salesChange >= 0 ? 'text-green-600' : 'text-red-600'}>
                {Math.abs(salesChange).toFixed(1)}% vs mes anterior
              </span>
            </div>
          </div>

          {/* Facturación del Mes */}
          <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-sm font-medium text-gray-600">Facturación del Mes</h3>
              <FileText className="w-5 h-5 text-purple-500" />
            </div>
            <div className="text-2xl font-bold text-gray-900 mb-2">
              {formatCurrency(data.invoices_this_month.total)}
            </div>
            <div className="flex items-center text-sm">
              <span className="text-gray-600">{data.invoices_this_month.count} facturas</span>
            </div>
          </div>

          {/* Ganancia del Mes */}
          <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-sm font-medium text-gray-600">Ganancia del Mes</h3>
              <DollarSign className="w-5 h-5 text-green-500" />
            </div>
            <div className="text-2xl font-bold text-gray-900 mb-2">
              {formatCurrency(data.sales_this_month.profit)}
            </div>
            <div className="flex items-center text-sm">
              <span className="text-gray-600">
                {data.sales_this_month.profit_margin.toFixed(1)}% margen
              </span>
            </div>
          </div>
        </div>

        {/* KPIs Secundarios - Módulos Nuevos */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {/* Devoluciones Pendientes */}
          <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-sm font-medium text-gray-600">Devoluciones</h3>
              <RotateCcw className="w-5 h-5 text-orange-500" />
            </div>
            <div className="text-2xl font-bold text-gray-900 mb-2">
              {data.returns_pending}
            </div>
            <div className="flex items-center text-sm">
              <span className="text-gray-600">
                {data.returns_this_month} este mes
              </span>
            </div>
            {data.returns_total_value > 0 && (
              <div className="mt-2 text-xs text-gray-500">
                {formatCurrency(data.returns_total_value)} en devoluciones
              </div>
            )}
          </div>

          {/* Notas de Crédito */}
          <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-sm font-medium text-gray-600">Notas de Crédito</h3>
              <FileX className="w-5 h-5 text-red-500" />
            </div>
            <div className="text-2xl font-bold text-gray-900 mb-2">
              {data.credit_notes_pending}
            </div>
            <div className="flex items-center text-sm">
              <span className="text-gray-600">
                {data.credit_notes_this_month} este mes
              </span>
            </div>
            {data.credit_notes_total_value > 0 && (
              <div className="mt-2 text-xs text-gray-500">
                {formatCurrency(data.credit_notes_total_value)} en notas
              </div>
            )}
          </div>

          {/* Compras Pendientes */}
          <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-sm font-medium text-gray-600">Compras</h3>
              <Package className="w-5 h-5 text-indigo-500" />
            </div>
            <div className="text-2xl font-bold text-gray-900 mb-2">
              {data.purchases_pending}
            </div>
            <div className="flex items-center text-sm">
              <span className="text-gray-600">
                {data.purchases_this_month} este mes
              </span>
            </div>
            {data.purchases_total_value > 0 && (
              <div className="mt-2 text-xs text-gray-500">
                {formatCurrency(data.purchases_total_value)} en compras
              </div>
            )}
          </div>

          {/* Gastos del Mes */}
          <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-sm font-medium text-gray-600">Gastos del Mes</h3>
              <TrendingDown className="w-5 h-5 text-red-500" />
            </div>
            <div className="text-2xl font-bold text-gray-900 mb-2">
              {formatCurrency(data.expenses_this_month)}
            </div>
            <div className="flex items-center text-sm">
              <span className="text-gray-600">
                {formatCurrency(data.expenses_today)} hoy
              </span>
            </div>
          </div>
        </div>

        {/* KPIs Financieros */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {/* Cuentas por Cobrar */}
          <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-sm font-medium text-gray-600">Cuentas por Cobrar</h3>
              <CreditCard className="w-5 h-5 text-blue-500" />
            </div>
            <div className="text-2xl font-bold text-gray-900 mb-2">
              {formatCurrency(data.accounts_receivable_total)}
            </div>
            <div className="flex items-center text-sm text-gray-600">
              <AlertCircle className="w-4 h-4 mr-1" />
              {data.invoices_pending} facturas pendientes
            </div>
            {data.invoices_overdue > 0 && (
              <div className="mt-2 text-xs text-red-600">
                {data.invoices_overdue} facturas vencidas
              </div>
            )}
          </div>

          {/* Cuentas por Pagar */}
          <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-sm font-medium text-gray-600">Cuentas por Pagar</h3>
              <Receipt className="w-5 h-5 text-orange-500" />
            </div>
            <div className="text-2xl font-bold text-gray-900 mb-2">
              {formatCurrency(data.accounts_payable_total)}
            </div>
            <div className="flex items-center text-sm text-gray-600">
              <AlertCircle className="w-4 h-4 mr-1" />
              {data.purchases_pending} compras pendientes
            </div>
          </div>

          {/* Inventario */}
          <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-sm font-medium text-gray-600">Inventario</h3>
              <Package className="w-5 h-5 text-green-500" />
            </div>
            <div className="text-2xl font-bold text-gray-900 mb-2">
              {formatCurrency(data.inventory_total_value)}
            </div>
            <div className="flex items-center text-sm text-gray-600">
              <span>{data.inventory_total_items} productos</span>
            </div>
            {data.inventory_low_stock_items > 0 && (
              <div className="mt-2 text-xs text-orange-600 flex items-center">
                <AlertCircle className="w-3 h-3 mr-1" />
                {data.inventory_low_stock_items} productos con stock bajo
              </div>
            )}
          </div>
        </div>

        {/* Actividades Recientes y Top Data */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Ventas Recientes */}
          <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-semibold text-gray-900">Ventas Recientes</h2>
              <Activity className="w-5 h-5 text-gray-400" />
            </div>
            <div className="space-y-3">
              {data.recent_sales && data.recent_sales.length > 0 ? (
                data.recent_sales.slice(0, 5).map((sale) => (
                  <div key={sale.sale_id} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                    <div>
                      <div className="font-medium text-gray-900">{sale.sale_number}</div>
                      <div className="text-sm text-gray-600">{sale.customer_name}</div>
                    </div>
                    <div className="text-right">
                      <div className="font-semibold text-gray-900">{formatCurrency(sale.total)}</div>
                      <div className="text-xs text-gray-500">{sale.sale_date}</div>
                    </div>
                  </div>
                ))
              ) : (
                <p className="text-sm text-gray-500 text-center py-4">No hay ventas recientes</p>
              )}
            </div>
          </div>

          {/* Facturas Recientes */}
          <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-semibold text-gray-900">Facturas Recientes</h2>
              <FileText className="w-5 h-5 text-gray-400" />
            </div>
            <div className="space-y-3">
              {data.recent_invoices && data.recent_invoices.length > 0 ? (
                data.recent_invoices.slice(0, 5).map((invoice) => (
                  <div key={invoice.id} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                    <div>
                      <div className="font-medium text-gray-900">{invoice.invoice_number}</div>
                      <div className="text-sm text-gray-600">{invoice.customer_name}</div>
                    </div>
                    <div className="text-right">
                      <div className="font-semibold text-gray-900">{formatCurrency(invoice.total)}</div>
                      <div className="text-xs text-gray-500">{invoice.issue_date}</div>
                    </div>
                  </div>
                ))
              ) : (
                <p className="text-sm text-gray-500 text-center py-4">No hay facturas recientes</p>
              )}
            </div>
          </div>
        </div>

        {/* Devoluciones y Notas de Crédito Recientes */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Devoluciones Recientes */}
          <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-semibold text-gray-900">Devoluciones Recientes</h2>
              <RotateCcw className="w-5 h-5 text-gray-400" />
            </div>
            <div className="space-y-3">
              {data.recent_returns && data.recent_returns.length > 0 ? (
                data.recent_returns.slice(0, 5).map((ret) => (
                  <div key={ret.id} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                    <div>
                      <div className="font-medium text-gray-900">{ret.return_number}</div>
                      <div className="text-sm text-gray-600">{ret.customer_name}</div>
                    </div>
                    <div className="text-right">
                      <div className="font-semibold text-gray-900">{formatCurrency(ret.total)}</div>
                      <div className="text-xs text-gray-500">{ret.return_date}</div>
                    </div>
                  </div>
                ))
              ) : (
                <p className="text-sm text-gray-500 text-center py-4">No hay devoluciones recientes</p>
              )}
            </div>
          </div>

          {/* Notas de Crédito Recientes */}
          <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-semibold text-gray-900">Notas de Crédito Recientes</h2>
              <FileX className="w-5 h-5 text-gray-400" />
            </div>
            <div className="space-y-3">
              {data.recent_credit_notes && data.recent_credit_notes.length > 0 ? (
                data.recent_credit_notes.slice(0, 5).map((cn) => (
                  <div key={cn.id} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                    <div>
                      <div className="font-medium text-gray-900">{cn.credit_note_number}</div>
                      <div className="text-sm text-gray-600">{cn.customer_name}</div>
                    </div>
                    <div className="text-right">
                      <div className="font-semibold text-gray-900">{formatCurrency(cn.total)}</div>
                      <div className="text-xs text-gray-500">{cn.issue_date}</div>
                    </div>
                  </div>
                ))
              ) : (
                <p className="text-sm text-gray-500 text-center py-4">No hay notas de crédito recientes</p>
              )}
            </div>
          </div>
        </div>

        {/* Top Productos y Top Clientes */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Top Productos */}
          <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-semibold text-gray-900">Top Productos (Últimos 30 días)</h2>
              <BarChart3 className="w-5 h-5 text-gray-400" />
            </div>
            <div className="space-y-3">
              {data.top_products && data.top_products.length > 0 ? (
                data.top_products.map((product, index) => (
                  <div key={product.product_id} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                    <div className="flex items-center">
                      <div className="w-8 h-8 bg-blue-100 text-blue-600 rounded-full flex items-center justify-center font-semibold text-sm mr-3">
                        {index + 1}
                      </div>
                      <div>
                        <div className="font-medium text-gray-900">{product.product_name}</div>
                        <div className="text-sm text-gray-600">{product.product_code}</div>
                      </div>
                    </div>
                    <div className="text-right">
                      <div className="font-semibold text-gray-900">{formatCurrency(product.total)}</div>
                      <div className="text-xs text-gray-500">{formatNumber(product.quantity)} unidades</div>
                    </div>
                  </div>
                ))
              ) : (
                <p className="text-sm text-gray-500 text-center py-4">No hay datos de productos</p>
              )}
            </div>
          </div>

          {/* Top Clientes */}
          <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-semibold text-gray-900">Top Clientes (Últimos 30 días)</h2>
              <Users className="w-5 h-5 text-gray-400" />
            </div>
            <div className="space-y-3">
              {data.top_customers && data.top_customers.length > 0 ? (
                data.top_customers.map((customer, index) => (
                  <div key={customer.customer_id} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                    <div className="flex items-center">
                      <div className="w-8 h-8 bg-green-100 text-green-600 rounded-full flex items-center justify-center font-semibold text-sm mr-3">
                        {index + 1}
                      </div>
                      <div>
                        <div className="font-medium text-gray-900">{customer.customer_name}</div>
                        <div className="text-sm text-gray-600">{customer.count} compras</div>
                      </div>
                    </div>
                    <div className="text-right">
                      <div className="font-semibold text-gray-900">{formatCurrency(customer.total)}</div>
                      <div className="text-xs text-gray-500">{formatCurrency(customer.subtotal)} + ITBIS</div>
                    </div>
                  </div>
                ))
              ) : (
                <p className="text-sm text-gray-500 text-center py-4">No hay datos de clientes</p>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
