import { useState, useEffect } from 'react'
import api from '../lib/api'
import Alert from '../components/Alert'
import { BarChart, Package, DollarSign, TrendingUp, Download, RefreshCw, CheckSquare, XSquare, FileText, RotateCcw } from 'lucide-react'

export default function Reports() {
  const [activeTab, setActiveTab] = useState('inventory')
  const [alert, setAlert] = useState(null)
  const [loading, setLoading] = useState(false)

  // Inventory Valuation state
  const [inventoryReport, setInventoryReport] = useState(null)
  const [includeZeroStock, setIncludeZeroStock] = useState(false)

  // Accounts Receivable state
  const [arReport, setArReport] = useState(null)
  const [arAsOfDate, setArAsOfDate] = useState(new Date().toISOString().split('T')[0])

  // Returns Report state
  const [returnsReport, setReturnsReport] = useState(null)
  const [returnsStartDate, setReturnsStartDate] = useState(() => {
    const date = new Date()
    date.setDate(date.getDate() - 30)
    return date.toISOString().split('T')[0]
  })
  const [returnsEndDate, setReturnsEndDate] = useState(new Date().toISOString().split('T')[0])

  // Sales Report state
  const [salesReport, setSalesReport] = useState(null)
  const [salesStartDate, setSalesStartDate] = useState(() => {
    const date = new Date()
    date.setDate(date.getDate() - 30)
    return date.toISOString().split('T')[0]
  })
  const [salesEndDate, setSalesEndDate] = useState(new Date().toISOString().split('T')[0])

  // Tax Report state
  const [taxReport, setTaxReport] = useState(null)
  const [taxStartDate, setTaxStartDate] = useState(() => {
    const date = new Date()
    date.setDate(date.getDate() - 30)
    return date.toISOString().split('T')[0]
  })
  const [taxEndDate, setTaxEndDate] = useState(new Date().toISOString().split('T')[0])

  const tabs = [
    { id: 'inventory', label: 'Valoración de Inventario', icon: Package },
    { id: 'receivables', label: 'Cuentas por Cobrar', icon: FileText },
    { id: 'returns', label: 'Devoluciones', icon: RotateCcw },
    { id: 'sales', label: 'Ventas', icon: TrendingUp },
    { id: 'taxes', label: 'Impuestos', icon: DollarSign },
  ]

  useEffect(() => {
    if (activeTab === 'inventory') {
      fetchInventoryValuation()
    } else if (activeTab === 'receivables') {
      fetchAccountsReceivable()
    } else if (activeTab === 'returns') {
      fetchReturnsReport()
    } else if (activeTab === 'sales') {
      fetchSalesReport()
    } else if (activeTab === 'taxes') {
      fetchTaxReport()
    }
  }, [activeTab, includeZeroStock, arAsOfDate, returnsStartDate, returnsEndDate, salesStartDate, salesEndDate, taxStartDate, taxEndDate])

  const fetchInventoryValuation = async () => {
    try {
      setLoading(true)
      const url = `/reports/inventory-valuation${includeZeroStock ? '?include_zero_stock=true' : ''}`
      const response = await api.get(url)
      setInventoryReport(response.data)
    } catch (error) {
      console.error('Error fetching inventory valuation:', error)
      const errorMessage = error.response?.data?.error || error.message || 'Error al cargar el informe de valoración'
      setAlert({ type: 'error', message: errorMessage })
      setInventoryReport(null)
    } finally {
      setLoading(false)
    }
  }

  const fetchAccountsReceivable = async () => {
    try {
      setLoading(true)
      const url = `/reports/accounts-receivable${arAsOfDate ? `?as_of_date=${arAsOfDate}` : ''}`
      const response = await api.get(url)
      setArReport(response.data)
    } catch (error) {
      console.error('Error fetching accounts receivable:', error)
      const errorMessage = error.response?.data?.error || error.message || 'Error al cargar el informe de cuentas por cobrar'
      setAlert({ type: 'error', message: errorMessage })
      setArReport(null)
    } finally {
      setLoading(false)
    }
  }

  const fetchReturnsReport = async () => {
    try {
      setLoading(true)
      const params = new URLSearchParams()
      if (returnsStartDate) params.append('start_date', returnsStartDate)
      if (returnsEndDate) params.append('end_date', returnsEndDate)
      const url = `/reports/returns?${params.toString()}`
      const response = await api.get(url)
      setReturnsReport(response.data)
    } catch (error) {
      console.error('Error fetching returns report:', error)
      const errorMessage = error.response?.data?.error || error.message || 'Error al cargar el informe de devoluciones'
      setAlert({ type: 'error', title: 'Error', message: errorMessage })
      setReturnsReport(null)
    } finally {
      setLoading(false)
    }
  }

  const fetchSalesReport = async () => {
    try {
      setLoading(true)
      const params = new URLSearchParams()
      if (salesStartDate) params.append('start_date', salesStartDate)
      if (salesEndDate) params.append('end_date', salesEndDate)
      const url = `/reports/sales?${params.toString()}`
      const response = await api.get(url)
      setSalesReport(response.data)
    } catch (error) {
      console.error('Error fetching sales report:', error)
      const errorMessage = error.response?.data?.error || error.message || 'Error al cargar el informe de ventas'
      setAlert({ type: 'error', title: 'Error', message: errorMessage })
      setSalesReport(null)
    } finally {
      setLoading(false)
    }
  }

  const fetchTaxReport = async () => {
    try {
      setLoading(true)
      const params = new URLSearchParams()
      if (taxStartDate) params.append('start_date', taxStartDate)
      if (taxEndDate) params.append('end_date', taxEndDate)
      const url = `/reports/taxes?${params.toString()}`
      const response = await api.get(url)
      setTaxReport(response.data)
    } catch (error) {
      console.error('Error fetching tax report:', error)
      const errorMessage = error.response?.data?.error || error.message || 'Error al cargar el informe de impuestos'
      setAlert({ type: 'error', title: 'Error', message: errorMessage })
      setTaxReport(null)
    } finally {
      setLoading(false)
    }
  }

  const formatCurrency = (amount) => {
    return new Intl.NumberFormat('es-DO', {
      style: 'currency',
      currency: 'DOP',
      minimumFractionDigits: 2,
    }).format(amount || 0)
  }

  const formatNumber = (num) => {
    return new Intl.NumberFormat('es-DO', {
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    }).format(num || 0)
  }

  const showAlert = (type, title, message) => {
    setAlert({ type, title, message })
    setTimeout(() => setAlert(null), 5000)
  }

  return (
    <div className="p-6">
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-[#212121] flex items-center gap-3">
          <BarChart className="w-8 h-8 text-[#FF6B00]" />
          Reportes
        </h1>
        <p className="text-gray-600 mt-1">Informes y análisis del negocio</p>
      </div>

      {alert && (
        <Alert
          type={alert.type}
          title={alert.title}
          message={alert.message}
          onClose={() => setAlert(null)}
        />
      )}

      {/* Tabs */}
      <div className="mb-6 border-b border-gray-200">
        <nav className="flex space-x-8">
          {tabs.map((tab) => {
            const Icon = tab.icon
            return (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id)}
                className={`flex items-center gap-2 py-4 px-1 border-b-2 font-medium text-sm transition-colors ${
                  activeTab === tab.id
                    ? 'border-[#FF6B00] text-[#FF6B00]'
                    : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                }`}
              >
                <Icon className="w-5 h-5" />
                {tab.label}
              </button>
            )
          })}
        </nav>
      </div>

      {/* Tab Content */}
      <div className="bg-white rounded-lg shadow">
        {activeTab === 'inventory' && (
          <div className="p-6">
            <div className="flex justify-between items-center mb-6">
              <div>
                <h2 className="text-2xl font-bold text-gray-900">Valoración de Inventario</h2>
                <p className="text-gray-600 mt-1">
                  Informe detallado del valor del inventario actual
                </p>
              </div>
              <div className="flex gap-3">
                <label className="flex items-center gap-2 cursor-pointer">
                  <input
                    type="checkbox"
                    checked={includeZeroStock}
                    onChange={(e) => setIncludeZeroStock(e.target.checked)}
                    className="w-4 h-4 text-[#FF6B00] border-gray-300 rounded focus:ring-[#FF6B00]"
                  />
                  <span className="text-sm text-gray-700">Incluir productos sin stock</span>
                </label>
                <button
                  onClick={fetchInventoryValuation}
                  disabled={loading}
                  className="flex items-center gap-2 px-4 py-2 bg-[#FF6B00] text-white rounded-lg hover:bg-[#E55A00] transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  <RefreshCw className={`w-4 h-4 ${loading ? 'animate-spin' : ''}`} />
                  Actualizar
                </button>
              </div>
            </div>

            {loading ? (
              <div className="flex justify-center items-center py-12">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-[#FF6B00]"></div>
              </div>
            ) : inventoryReport ? (
              <>
                {/* Summary Cards */}
                <div className="grid grid-cols-1 md:grid-cols-5 gap-4 mb-6">
                  <div className="bg-blue-50 p-4 rounded-lg border border-blue-200">
                    <div className="text-sm text-blue-600 font-medium">Total de Productos</div>
                    <div className="text-2xl font-bold text-blue-900 mt-1">
                      {inventoryReport.total_items || 0}
                    </div>
                  </div>
                  <div className="bg-green-50 p-4 rounded-lg border border-green-200">
                    <div className="text-sm text-green-600 font-medium">Valor Total (Costo)</div>
                    <div className="text-2xl font-bold text-green-900 mt-1">
                      {formatCurrency(inventoryReport.total_value || 0)}
                    </div>
                  </div>
                  <div className="bg-purple-50 p-4 rounded-lg border border-purple-200">
                    <div className="text-sm text-purple-600 font-medium">Con Costo Definido</div>
                    <div className="text-2xl font-bold text-purple-900 mt-1">
                      {formatCurrency(inventoryReport.total_value_with_cost || 0)}
                    </div>
                  </div>
                  <div className="bg-amber-50 p-4 rounded-lg border border-amber-200">
                    <div className="text-sm text-amber-600 font-medium">Valor de Venta Total</div>
                    <div className="text-2xl font-bold text-amber-900 mt-1">
                      {formatCurrency(inventoryReport.total_sale_value || 0)}
                    </div>
                  </div>
                  <div className="bg-orange-50 p-4 rounded-lg border border-orange-200">
                    <div className="text-sm text-orange-600 font-medium">Fecha del Informe</div>
                    <div className="text-lg font-semibold text-orange-900 mt-1">
                      {inventoryReport.report_date || 'N/A'}
                    </div>
                  </div>
                </div>

                {/* Items Table */}
                <div className="mb-6">
                  <h3 className="text-lg font-semibold text-gray-900 mb-4">Detalle de Productos</h3>
                  <div className="overflow-x-auto">
                    <table className="min-w-full divide-y divide-gray-200">
                      <thead className="bg-gray-50">
                        <tr>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Código
                          </th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Producto
                          </th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Tipo
                          </th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Cantidad
                          </th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Costo Unitario
                          </th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Valor Total
                          </th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Proveedor
                          </th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Última Actualización
                          </th>
                        </tr>
                      </thead>
                      <tbody className="bg-white divide-y divide-gray-200">
                        {inventoryReport.items && inventoryReport.items.length > 0 ? (
                          inventoryReport.items.map((item, index) => (
                            <tr key={item.product_id || index} className="hover:bg-gray-50">
                              <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                                {item.product_code}
                              </td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                                {item.product_name}
                              </td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                {item.product_type}
                              </td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                                {formatNumber(item.stock_quantity)} {item.unit}
                              </td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                                {item.unit_cost ? formatCurrency(item.unit_cost) : 'N/A'}
                              </td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm font-semibold text-gray-900">
                                {formatCurrency(item.total_value)}
                              </td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                {item.supplier_name || 'Sin proveedor'}
                              </td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                {item.last_updated || 'N/A'}
                              </td>
                            </tr>
                          ))
                        ) : (
                          <tr>
                            <td colSpan="8" className="px-6 py-8 text-center text-gray-500">
                              No hay productos en inventario
                            </td>
                          </tr>
                        )}
                      </tbody>
                      {inventoryReport.items && inventoryReport.items.length > 0 && (
                        <tfoot className="bg-gray-50">
                          <tr>
                            <td colSpan="4" className="px-6 py-4 text-right text-sm font-semibold text-gray-900">
                              Total:
                            </td>
                            <td colSpan="1" className="px-6 py-4 text-sm text-gray-500"></td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm font-bold text-gray-900">
                              {formatCurrency(inventoryReport.total_value || 0)}
                            </td>
                            <td colSpan="2"></td>
                          </tr>
                        </tfoot>
                      )}
                    </table>
                  </div>
                </div>

                {/* Summary by Type */}
                {inventoryReport.summary_by_type && Object.keys(inventoryReport.summary_by_type).length > 0 && (
                  <div className="mb-6">
                    <h3 className="text-lg font-semibold text-gray-900 mb-4">Resumen por Tipo de Producto</h3>
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                      {Object.values(inventoryReport.summary_by_type).map((summary, index) => (
                        <div key={index} className="bg-gray-50 p-4 rounded-lg border border-gray-200">
                          <div className="font-semibold text-gray-900 mb-2">{summary.product_type || 'Sin tipo'}</div>
                          <div className="text-sm text-gray-600 space-y-1">
                            <div>Productos: {summary.item_count}</div>
                            <div>Cantidad: {formatNumber(summary.total_quantity)}</div>
                            <div className="font-semibold text-gray-900">
                              Valor: {formatCurrency(summary.total_value)}
                            </div>
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                )}

                {/* Summary by Supplier */}
                {inventoryReport.summary_by_supplier && Object.keys(inventoryReport.summary_by_supplier).length > 0 && (
                  <div>
                    <h3 className="text-lg font-semibold text-gray-900 mb-4">Resumen por Proveedor</h3>
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                      {Object.values(inventoryReport.summary_by_supplier).map((summary, index) => (
                        <div key={index} className="bg-gray-50 p-4 rounded-lg border border-gray-200">
                          <div className="font-semibold text-gray-900 mb-2">
                            {summary.supplier_name || 'Sin proveedor'}
                          </div>
                          <div className="text-sm text-gray-600 space-y-1">
                            <div>Productos: {summary.item_count}</div>
                            <div>Cantidad: {formatNumber(summary.total_quantity)}</div>
                            <div className="font-semibold text-gray-900">
                              Valor: {formatCurrency(summary.total_value)}
                            </div>
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                )}
              </>
            ) : (
              <div className="text-center py-12 text-gray-500">
                No se pudo cargar el informe de valoración
              </div>
            )}
          </div>
        )}

        {activeTab === 'receivables' && (
          <div className="p-6">
            <div className="flex justify-between items-center mb-6">
              <div>
                <h2 className="text-2xl font-bold text-gray-900">Cuentas por Cobrar</h2>
                <p className="text-gray-600 mt-1">
                  Informe de envejecimiento de cuentas por cobrar
                </p>
              </div>
              <div className="flex gap-3">
                <input
                  type="date"
                  value={arAsOfDate}
                  onChange={(e) => setArAsOfDate(e.target.value)}
                  className="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#FF6B00] focus:border-transparent"
                />
                <button
                  onClick={fetchAccountsReceivable}
                  disabled={loading}
                  className="flex items-center gap-2 px-4 py-2 bg-[#FF6B00] text-white rounded-lg hover:bg-[#E55A00] transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  <RefreshCw className={`w-4 h-4 ${loading ? 'animate-spin' : ''}`} />
                  Actualizar
                </button>
              </div>
            </div>

            {loading ? (
              <div className="flex justify-center items-center py-12">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-[#FF6B00]"></div>
              </div>
            ) : arReport ? (
              <>
                {/* Summary Cards */}
                <div className="grid grid-cols-1 md:grid-cols-5 gap-4 mb-6">
                  <div className="bg-blue-50 p-4 rounded-lg border border-blue-200">
                    <div className="text-sm text-blue-600 font-medium">Total Facturas</div>
                    <div className="text-2xl font-bold text-blue-900 mt-1">
                      {arReport.total_invoices || 0}
                    </div>
                  </div>
                  <div className="bg-green-50 p-4 rounded-lg border border-green-200">
                    <div className="text-sm text-green-600 font-medium">Total Pendiente</div>
                    <div className="text-2xl font-bold text-green-900 mt-1">
                      {formatCurrency(arReport.total_outstanding || 0)}
                    </div>
                  </div>
                  <div className="bg-yellow-50 p-4 rounded-lg border border-yellow-200">
                    <div className="text-sm text-yellow-600 font-medium">0-30 Días</div>
                    <div className="text-2xl font-bold text-yellow-900 mt-1">
                      {formatCurrency(arReport.total_0_30 || 0)}
                    </div>
                  </div>
                  <div className="bg-orange-50 p-4 rounded-lg border border-orange-200">
                    <div className="text-sm text-orange-600 font-medium">31-60 Días</div>
                    <div className="text-2xl font-bold text-orange-900 mt-1">
                      {formatCurrency(arReport.total_31_60 || 0)}
                    </div>
                  </div>
                  <div className="bg-red-50 p-4 rounded-lg border border-red-200">
                    <div className="text-sm text-red-600 font-medium">+60 Días</div>
                    <div className="text-2xl font-bold text-red-900 mt-1">
                      {formatCurrency((arReport.total_61_90 || 0) + (arReport.total_over_90 || 0))}
                    </div>
                  </div>
                </div>

                {/* Aging Breakdown */}
                <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
                  <div className="bg-gray-50 p-4 rounded-lg border border-gray-200">
                    <div className="text-sm text-gray-600 font-medium">0-30 Días</div>
                    <div className="text-xl font-bold text-gray-900 mt-1">
                      {formatCurrency(arReport.total_0_30 || 0)}
                    </div>
                  </div>
                  <div className="bg-gray-50 p-4 rounded-lg border border-gray-200">
                    <div className="text-sm text-gray-600 font-medium">31-60 Días</div>
                    <div className="text-xl font-bold text-gray-900 mt-1">
                      {formatCurrency(arReport.total_31_60 || 0)}
                    </div>
                  </div>
                  <div className="bg-gray-50 p-4 rounded-lg border border-gray-200">
                    <div className="text-sm text-gray-600 font-medium">61-90 Días</div>
                    <div className="text-xl font-bold text-gray-900 mt-1">
                      {formatCurrency(arReport.total_61_90 || 0)}
                    </div>
                  </div>
                  <div className="bg-gray-50 p-4 rounded-lg border border-gray-200">
                    <div className="text-sm text-gray-600 font-medium">+90 Días</div>
                    <div className="text-xl font-bold text-gray-900 mt-1">
                      {formatCurrency(arReport.total_over_90 || 0)}
                    </div>
                  </div>
                </div>

                {/* Items Table */}
                <div className="mb-6">
                  <h3 className="text-lg font-semibold text-gray-900 mb-4">Detalle de Facturas</h3>
                  <div className="overflow-x-auto">
                    <table className="min-w-full divide-y divide-gray-200">
                      <thead className="bg-gray-50">
                        <tr>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Factura
                          </th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Cliente
                          </th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Fecha Emisión
                          </th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Fecha Vencimiento
                          </th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Total
                          </th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Pagado
                          </th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Saldo
                          </th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Días Vencido
                          </th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Estado
                          </th>
                        </tr>
                      </thead>
                      <tbody className="bg-white divide-y divide-gray-200">
                        {arReport.items && arReport.items.length > 0 ? (
                          arReport.items.map((item, index) => {
                            const getStatusColor = (status) => {
                              switch (status) {
                                case 'overdue':
                                  return 'bg-red-100 text-red-800'
                                case 'pending':
                                  return 'bg-yellow-100 text-yellow-800'
                                case 'paid':
                                  return 'bg-green-100 text-green-800'
                                default:
                                  return 'bg-gray-100 text-gray-800'
                              }
                            }
                            return (
                              <tr key={item.invoice_id || index} className="hover:bg-gray-50">
                                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                                  {item.invoice_number}
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                                  {item.customer_name}
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                  {new Date(item.issue_date).toLocaleDateString('es-DO')}
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                  {new Date(item.due_date).toLocaleDateString('es-DO')}
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                                  {formatCurrency(item.total)}
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                  {formatCurrency(item.paid_amount)}
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap text-sm font-semibold text-gray-900">
                                  {formatCurrency(item.balance)}
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                  {item.days_overdue !== null && item.days_overdue !== undefined
                                    ? `${item.days_overdue} días`
                                    : '-'}
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap">
                                  <span className={`px-2 py-1 text-xs font-medium rounded-full ${getStatusColor(item.status)}`}>
                                    {item.status === 'overdue' ? 'Vencida' : 
                                     item.status === 'pending' ? 'Pendiente' : 
                                     item.status === 'paid' ? 'Pagada' : item.status}
                                  </span>
                                </td>
                              </tr>
                            )
                          })
                        ) : (
                          <tr>
                            <td colSpan="9" className="px-6 py-8 text-center text-gray-500">
                              No hay facturas pendientes
                            </td>
                          </tr>
                        )}
                      </tbody>
                      {arReport.items && arReport.items.length > 0 && (
                        <tfoot className="bg-gray-50">
                          <tr>
                            <td colSpan="6" className="px-6 py-4 text-right text-sm font-semibold text-gray-900">
                              Total Pendiente:
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm font-bold text-gray-900">
                              {formatCurrency(arReport.total_outstanding || 0)}
                            </td>
                            <td colSpan="2"></td>
                          </tr>
                        </tfoot>
                      )}
                    </table>
                  </div>
                </div>

                {/* Summary by Customer */}
                {arReport.summary_by_customer && Object.keys(arReport.summary_by_customer).length > 0 && (
                  <div>
                    <h3 className="text-lg font-semibold text-gray-900 mb-4">Resumen por Cliente</h3>
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                      {Object.values(arReport.summary_by_customer).map((summary, index) => (
                        <div key={index} className="bg-gray-50 p-4 rounded-lg border border-gray-200">
                          <div className="font-semibold text-gray-900 mb-2">{summary.customer_name}</div>
                          <div className="text-sm text-gray-600 space-y-1">
                            <div>Facturas: {summary.invoice_count}</div>
                            <div className="font-semibold text-gray-900">
                              Total: {formatCurrency(summary.total_outstanding)}
                            </div>
                            <div className="text-xs text-gray-500 mt-2">
                              <div>0-30: {formatCurrency(summary.total_0_30)}</div>
                              <div>31-60: {formatCurrency(summary.total_31_60)}</div>
                              <div>61-90: {formatCurrency(summary.total_61_90)}</div>
                              <div>+90: {formatCurrency(summary.total_over_90)}</div>
                            </div>
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                )}
              </>
            ) : (
              <div className="text-center py-12 text-gray-500">
                No se pudo cargar el informe de cuentas por cobrar
              </div>
            )}
          </div>
        )}

        {activeTab === 'returns' && (
          <div className="p-6">
            <div className="flex justify-between items-center mb-6">
              <div>
                <h2 className="text-2xl font-bold text-gray-900">Reporte de Devoluciones</h2>
                <p className="text-gray-600 mt-1">
                  Informe detallado de devoluciones y reembolsos
                </p>
              </div>
              <div className="flex gap-3">
                <input
                  type="date"
                  value={returnsStartDate}
                  onChange={(e) => setReturnsStartDate(e.target.value)}
                  className="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#FF6B00] focus:border-transparent"
                />
                <input
                  type="date"
                  value={returnsEndDate}
                  onChange={(e) => setReturnsEndDate(e.target.value)}
                  className="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#FF6B00] focus:border-transparent"
                />
                <button
                  onClick={fetchReturnsReport}
                  disabled={loading}
                  className="flex items-center gap-2 px-4 py-2 bg-[#FF6B00] text-white rounded-lg hover:bg-[#E55A00] transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  <RefreshCw className={`w-4 h-4 ${loading ? 'animate-spin' : ''}`} />
                  Actualizar
                </button>
              </div>
            </div>

            {loading ? (
              <div className="flex justify-center items-center py-12">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-[#FF6B00]"></div>
              </div>
            ) : returnsReport ? (
              <>
                {/* Summary Cards */}
                <div className="grid grid-cols-1 md:grid-cols-5 gap-4 mb-6">
                  <div className="bg-red-50 p-4 rounded-lg border border-red-200">
                    <div className="text-sm text-red-600 font-medium">Total Devoluciones</div>
                    <div className="text-2xl font-bold text-red-900 mt-1">
                      {returnsReport.summary?.total_returns || 0}
                    </div>
                  </div>
                  <div className="bg-orange-50 p-4 rounded-lg border border-orange-200">
                    <div className="text-sm text-orange-600 font-medium">Total Subtotal</div>
                    <div className="text-2xl font-bold text-orange-900 mt-1">
                      {formatCurrency(returnsReport.summary?.total_subtotal || 0)}
                    </div>
                  </div>
                  <div className="bg-yellow-50 p-4 rounded-lg border border-yellow-200">
                    <div className="text-sm text-yellow-600 font-medium">Total ITBIS</div>
                    <div className="text-2xl font-bold text-yellow-900 mt-1">
                      {formatCurrency(returnsReport.summary?.total_tax || 0)}
                    </div>
                  </div>
                  <div className="bg-red-100 p-4 rounded-lg border border-red-300">
                    <div className="text-sm text-red-700 font-medium">Total Devolución</div>
                    <div className="text-2xl font-bold text-red-900 mt-1">
                      {formatCurrency(returnsReport.summary?.total_amount || 0)}
                    </div>
                  </div>
                  <div className="bg-green-50 p-4 rounded-lg border border-green-200">
                    <div className="text-sm text-green-600 font-medium">Total Reembolsado</div>
                    <div className="text-2xl font-bold text-green-900 mt-1">
                      {formatCurrency(returnsReport.summary?.total_refunded || 0)}
                    </div>
                  </div>
                </div>

                {/* Date Range Info */}
                <div className="mb-6 text-sm text-gray-600">
                  Período: {returnsReport.start_date} al {returnsReport.end_date} | 
                  Promedio por devolución: {formatCurrency(returnsReport.summary?.average_return || 0)}
                </div>

                {/* Summary by Status */}
                {returnsReport.summary?.by_status && Object.keys(returnsReport.summary.by_status).length > 0 && (
                  <div className="mb-6">
                    <h3 className="text-lg font-semibold text-gray-900 mb-4">Resumen por Estado</h3>
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                      {Object.values(returnsReport.summary.by_status).map((summary, index) => (
                        <div key={index} className="bg-gray-50 p-4 rounded-lg border border-gray-200">
                          <div className="font-semibold text-gray-900 mb-2">{summary.status}</div>
                          <div className="text-sm text-gray-600 space-y-1">
                            <div>Cantidad: {summary.count}</div>
                            <div className="font-semibold text-gray-900">
                              Total: {formatCurrency(summary.total)}
                            </div>
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                )}

                {/* Summary by Reason */}
                {returnsReport.by_reason && returnsReport.by_reason.length > 0 && (
                  <div className="mb-6">
                    <h3 className="text-lg font-semibold text-gray-900 mb-4">Resumen por Motivo</h3>
                    <div className="overflow-x-auto">
                      <table className="min-w-full divide-y divide-gray-200">
                        <thead className="bg-gray-50">
                          <tr>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Motivo</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Cantidad</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Subtotal</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">ITBIS</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Total</th>
                          </tr>
                        </thead>
                        <tbody className="bg-white divide-y divide-gray-200">
                          {returnsReport.by_reason.map((item, index) => (
                            <tr key={index} className="hover:bg-gray-50">
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{item.reason}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">{item.count}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">{formatCurrency(item.subtotal)}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">{formatCurrency(item.tax)}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm font-semibold text-right text-gray-900">{formatCurrency(item.total)}</td>
                            </tr>
                          ))}
                        </tbody>
                      </table>
                    </div>
                  </div>
                )}

                {/* Summary by Refund Method */}
                {returnsReport.by_refund_method && returnsReport.by_refund_method.length > 0 && (
                  <div className="mb-6">
                    <h3 className="text-lg font-semibold text-gray-900 mb-4">Resumen por Método de Reembolso</h3>
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                      {returnsReport.by_refund_method.map((item, index) => (
                        <div key={index} className="bg-gray-50 p-4 rounded-lg border border-gray-200">
                          <div className="font-semibold text-gray-900 mb-2">{item.refund_method}</div>
                          <div className="text-sm text-gray-600 space-y-1">
                            <div>Cantidad: {item.count}</div>
                            <div className="font-semibold text-gray-900">
                              Total: {formatCurrency(item.total)}
                            </div>
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                )}

                {/* Top Products */}
                {returnsReport.top_products && returnsReport.top_products.length > 0 && (
                  <div className="mb-6">
                    <h3 className="text-lg font-semibold text-gray-900 mb-4">Productos Más Devueltos</h3>
                    <div className="overflow-x-auto">
                      <table className="min-w-full divide-y divide-gray-200">
                        <thead className="bg-gray-50">
                          <tr>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Producto</th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Código</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Cantidad Devuelta</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase"># Devoluciones</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Total</th>
                          </tr>
                        </thead>
                        <tbody className="bg-white divide-y divide-gray-200">
                          {returnsReport.top_products.map((product, index) => (
                            <tr key={index} className="hover:bg-gray-50">
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{product.product_name}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{product.product_code}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">{formatNumber(product.quantity)}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">{product.count}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm font-semibold text-right text-gray-900">{formatCurrency(product.total)}</td>
                            </tr>
                          ))}
                        </tbody>
                      </table>
                    </div>
                  </div>
                )}

                {/* Daily Breakdown */}
                {returnsReport.daily_breakdown && returnsReport.daily_breakdown.length > 0 && (
                  <div className="mb-6">
                    <h3 className="text-lg font-semibold text-gray-900 mb-4">Desglose Diario</h3>
                    <div className="overflow-x-auto">
                      <table className="min-w-full divide-y divide-gray-200">
                        <thead className="bg-gray-50">
                          <tr>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Fecha</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Cantidad</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Subtotal</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">ITBIS</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Total</th>
                          </tr>
                        </thead>
                        <tbody className="bg-white divide-y divide-gray-200">
                          {returnsReport.daily_breakdown.map((day, index) => (
                            <tr key={index} className="hover:bg-gray-50">
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{day.date}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">{day.count}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">{formatCurrency(day.subtotal)}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">{formatCurrency(day.tax)}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm font-semibold text-right text-gray-900">{formatCurrency(day.total)}</td>
                            </tr>
                          ))}
                        </tbody>
                      </table>
                    </div>
                  </div>
                )}

                {/* Items Table */}
                <div className="mb-6">
                  <h3 className="text-lg font-semibold text-gray-900 mb-4">Detalle de Devoluciones</h3>
                  <div className="overflow-x-auto">
                    <table className="min-w-full divide-y divide-gray-200">
                      <thead className="bg-gray-50">
                        <tr>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Número</th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Fecha</th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Tipo</th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Original</th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Cliente</th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Motivo</th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Estado</th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Método Reembolso</th>
                          <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Total</th>
                          <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Reembolsado</th>
                        </tr>
                      </thead>
                      <tbody className="bg-white divide-y divide-gray-200">
                        {returnsReport.items && returnsReport.items.length > 0 ? (
                          returnsReport.items.map((item, index) => (
                            <tr key={item.return_id || index} className="hover:bg-gray-50">
                              <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{item.return_number}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{item.return_date}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{item.return_type === 'sale' ? 'Venta' : 'Factura'}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{item.original_number || 'N/A'}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{item.customer_name}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 max-w-xs truncate">{item.return_reason || 'N/A'}</td>
                              <td className="px-6 py-4 whitespace-nowrap">
                                <span className={`px-2 py-1 text-xs rounded-full ${
                                  item.status === 'completed' ? 'bg-green-100 text-green-800' :
                                  item.status === 'approved' ? 'bg-blue-100 text-blue-800' :
                                  item.status === 'pending' ? 'bg-yellow-100 text-yellow-800' :
                                  item.status === 'rejected' ? 'bg-red-100 text-red-800' :
                                  'bg-gray-100 text-gray-800'
                                }`}>
                                  {item.status}
                                </span>
                              </td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{item.refund_method || 'N/A'}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm font-semibold text-right text-gray-900">{formatCurrency(item.total)}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">{formatCurrency(item.refund_amount)}</td>
                            </tr>
                          ))
                        ) : (
                          <tr>
                            <td colSpan="10" className="px-6 py-8 text-center text-gray-500">
                              No hay devoluciones en el período seleccionado
                            </td>
                          </tr>
                        )}
                      </tbody>
                    </table>
                  </div>
                </div>
              </>
            ) : (
              <div className="text-center py-12 text-gray-500">
                No se pudo cargar el informe de devoluciones
              </div>
            )}
          </div>
        )}

        {activeTab === 'sales' && (
          <div className="p-6">
            <div className="flex justify-between items-center mb-6">
              <div>
                <h2 className="text-2xl font-bold text-gray-900">Reporte de Ventas</h2>
                <p className="text-gray-600 mt-1">
                  Informe detallado de ventas, ganancias y análisis por cliente y producto
                </p>
              </div>
              <div className="flex gap-3">
                <input
                  type="date"
                  value={salesStartDate}
                  onChange={(e) => setSalesStartDate(e.target.value)}
                  className="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#FF6B00] focus:border-transparent"
                />
                <input
                  type="date"
                  value={salesEndDate}
                  onChange={(e) => setSalesEndDate(e.target.value)}
                  className="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#FF6B00] focus:border-transparent"
                />
                <button
                  onClick={fetchSalesReport}
                  disabled={loading}
                  className="flex items-center gap-2 px-4 py-2 bg-[#FF6B00] text-white rounded-lg hover:bg-[#E55A00] transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  <RefreshCw className={`w-4 h-4 ${loading ? 'animate-spin' : ''}`} />
                  Actualizar
                </button>
              </div>
            </div>

            {loading ? (
              <div className="flex justify-center items-center py-12">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-[#FF6B00]"></div>
              </div>
            ) : salesReport ? (
              <>
                {/* Summary Cards */}
                <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
                  <div className="bg-blue-50 p-4 rounded-lg border border-blue-200">
                    <div className="text-sm text-blue-600 font-medium">Total Ventas</div>
                    <div className="text-2xl font-bold text-blue-900 mt-1">
                      {salesReport.summary?.total_sales || 0}
                    </div>
                  </div>
                  <div className="bg-green-50 p-4 rounded-lg border border-green-200">
                    <div className="text-sm text-green-600 font-medium">Total Ingresos</div>
                    <div className="text-2xl font-bold text-green-900 mt-1">
                      {formatCurrency(salesReport.summary?.total_amount || 0)}
                    </div>
                  </div>
                  <div className="bg-purple-50 p-4 rounded-lg border border-purple-200">
                    <div className="text-sm text-purple-600 font-medium">Total Ganancia</div>
                    <div className="text-2xl font-bold text-purple-900 mt-1">
                      {formatCurrency(salesReport.summary?.total_profit || 0)}
                    </div>
                  </div>
                  <div className="bg-amber-50 p-4 rounded-lg border border-amber-200">
                    <div className="text-sm text-amber-600 font-medium">Margen de Ganancia</div>
                    <div className="text-2xl font-bold text-amber-900 mt-1">
                      {salesReport.summary?.profit_margin?.toFixed(1) || 0}%
                    </div>
                  </div>
                </div>

                {/* Additional Summary Cards */}
                <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
                  <div className="bg-gray-50 p-4 rounded-lg border border-gray-200">
                    <div className="text-sm text-gray-600 font-medium">Subtotal</div>
                    <div className="text-xl font-bold text-gray-900 mt-1">
                      {formatCurrency(salesReport.summary?.total_subtotal || 0)}
                    </div>
                  </div>
                  <div className="bg-yellow-50 p-4 rounded-lg border border-yellow-200">
                    <div className="text-sm text-yellow-600 font-medium">Total ITBIS</div>
                    <div className="text-xl font-bold text-yellow-900 mt-1">
                      {formatCurrency(salesReport.summary?.total_tax || 0)}
                    </div>
                  </div>
                  <div className="bg-red-50 p-4 rounded-lg border border-red-200">
                    <div className="text-sm text-red-600 font-medium">Total Costo</div>
                    <div className="text-xl font-bold text-red-900 mt-1">
                      {formatCurrency(salesReport.summary?.total_cost || 0)}
                    </div>
                  </div>
                  <div className="bg-indigo-50 p-4 rounded-lg border border-indigo-200">
                    <div className="text-sm text-indigo-600 font-medium">Venta Promedio</div>
                    <div className="text-xl font-bold text-indigo-900 mt-1">
                      {formatCurrency(salesReport.summary?.average_sale || 0)}
                    </div>
                  </div>
                </div>

                {/* Date Range Info */}
                <div className="mb-6 text-sm text-gray-600">
                  Período: {salesReport.start_date} al {salesReport.end_date} | 
                  Fecha del reporte: {salesReport.report_date}
                </div>

                {/* Sales by Payment Type */}
                {salesReport.by_payment_type && salesReport.by_payment_type.length > 0 && (
                  <div className="mb-6">
                    <h3 className="text-lg font-semibold text-gray-900 mb-4">Ventas por Tipo de Pago</h3>
                    <div className="overflow-x-auto">
                      <table className="min-w-full divide-y divide-gray-200">
                        <thead className="bg-gray-50">
                          <tr>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Tipo de Pago</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Cantidad</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Subtotal</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">ITBIS</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Total</th>
                          </tr>
                        </thead>
                        <tbody className="bg-white divide-y divide-gray-200">
                          {salesReport.by_payment_type.map((item, index) => (
                            <tr key={index} className="hover:bg-gray-50">
                              <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                                {item.payment_type === 'cash' ? 'Efectivo' :
                                 item.payment_type === 'card' ? 'Tarjeta' :
                                 item.payment_type === 'transfer' ? 'Transferencia' :
                                 item.payment_type === 'check' ? 'Cheque' :
                                 item.payment_type === 'mixed' ? 'Mixto' :
                                 item.payment_type}
                              </td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">{item.count}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">{formatCurrency(item.subtotal)}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">{formatCurrency(item.tax)}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm font-semibold text-right text-gray-900">{formatCurrency(item.total)}</td>
                            </tr>
                          ))}
                        </tbody>
                      </table>
                    </div>
                  </div>
                )}

                {/* Top Customers */}
                {salesReport.by_customer && salesReport.by_customer.length > 0 && (
                  <div className="mb-6">
                    <h3 className="text-lg font-semibold text-gray-900 mb-4">Top Clientes</h3>
                    <div className="overflow-x-auto">
                      <table className="min-w-full divide-y divide-gray-200">
                        <thead className="bg-gray-50">
                          <tr>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Cliente</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Compras</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Subtotal</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">ITBIS</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Total</th>
                          </tr>
                        </thead>
                        <tbody className="bg-white divide-y divide-gray-200">
                          {salesReport.by_customer.map((customer, index) => (
                            <tr key={customer.customer_id || index} className="hover:bg-gray-50">
                              <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{customer.customer_name}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">{customer.count}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">{formatCurrency(customer.subtotal)}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">{formatCurrency(customer.tax)}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm font-semibold text-right text-gray-900">{formatCurrency(customer.total)}</td>
                            </tr>
                          ))}
                        </tbody>
                      </table>
                    </div>
                  </div>
                )}

                {/* Top Products */}
                {salesReport.top_products && salesReport.top_products.length > 0 && (
                  <div className="mb-6">
                    <h3 className="text-lg font-semibold text-gray-900 mb-4">Top Productos Vendidos</h3>
                    <div className="overflow-x-auto">
                      <table className="min-w-full divide-y divide-gray-200">
                        <thead className="bg-gray-50">
                          <tr>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Producto</th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Código</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Cantidad</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Precio Unit.</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Subtotal</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">ITBIS</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Total</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Costo</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Ganancia</th>
                          </tr>
                        </thead>
                        <tbody className="bg-white divide-y divide-gray-200">
                          {salesReport.top_products.map((product, index) => (
                            <tr key={product.product_id || index} className="hover:bg-gray-50">
                              <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{product.product_name}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{product.product_code}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">{formatNumber(product.quantity)}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">{formatCurrency(product.unit_price)}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">{formatCurrency(product.subtotal)}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">{formatCurrency(product.tax)}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm font-semibold text-right text-gray-900">{formatCurrency(product.total)}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-red-600">{formatCurrency(product.cost)}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm font-semibold text-right text-green-600">{formatCurrency(product.profit)}</td>
                            </tr>
                          ))}
                        </tbody>
                      </table>
                    </div>
                  </div>
                )}

                {/* Daily Breakdown */}
                {salesReport.daily_breakdown && salesReport.daily_breakdown.length > 0 && (
                  <div className="mb-6">
                    <h3 className="text-lg font-semibold text-gray-900 mb-4">Desglose Diario</h3>
                    <div className="overflow-x-auto">
                      <table className="min-w-full divide-y divide-gray-200">
                        <thead className="bg-gray-50">
                          <tr>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Fecha</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Ventas</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Subtotal</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">ITBIS</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Total</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Costo</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Ganancia</th>
                          </tr>
                        </thead>
                        <tbody className="bg-white divide-y divide-gray-200">
                          {salesReport.daily_breakdown.map((day, index) => (
                            <tr key={index} className="hover:bg-gray-50">
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{day.date}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">{day.count}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">{formatCurrency(day.subtotal)}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">{formatCurrency(day.tax)}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm font-semibold text-right text-gray-900">{formatCurrency(day.total)}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-red-600">{formatCurrency(day.cost)}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm font-semibold text-right text-green-600">{formatCurrency(day.profit)}</td>
                            </tr>
                          ))}
                        </tbody>
                      </table>
                    </div>
                  </div>
                )}

                {/* Sales Items Table */}
                {salesReport.items && salesReport.items.length > 0 && (
                  <div className="mb-6">
                    <h3 className="text-lg font-semibold text-gray-900 mb-4">Detalle de Ventas</h3>
                    <div className="overflow-x-auto">
                      <table className="min-w-full divide-y divide-gray-200">
                        <thead className="bg-gray-50">
                          <tr>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Número</th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Fecha</th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Cliente</th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Tipo de Pago</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Subtotal</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">ITBIS</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Total</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Costo</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Ganancia</th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Estado</th>
                          </tr>
                        </thead>
                        <tbody className="bg-white divide-y divide-gray-200">
                          {salesReport.items.map((item, index) => (
                            <tr key={item.sale_id || index} className="hover:bg-gray-50">
                              <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{item.sale_number}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{item.sale_date}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{item.customer_name}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                {item.payment_type === 'cash' ? 'Efectivo' :
                                 item.payment_type === 'card' ? 'Tarjeta' :
                                 item.payment_type === 'transfer' ? 'Transferencia' :
                                 item.payment_type === 'check' ? 'Cheque' :
                                 item.payment_type === 'mixed' ? 'Mixto' :
                                 item.payment_type}
                              </td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">{formatCurrency(item.subtotal)}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">{formatCurrency(item.tax_amount)}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm font-semibold text-right text-gray-900">{formatCurrency(item.total)}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-red-600">{formatCurrency(item.cost)}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm font-semibold text-right text-green-600">{formatCurrency(item.profit)}</td>
                              <td className="px-6 py-4 whitespace-nowrap">
                                <span className={`px-2 py-1 text-xs rounded-full ${
                                  item.status === 'completed' ? 'bg-green-100 text-green-800' :
                                  'bg-gray-100 text-gray-800'
                                }`}>
                                  {item.status}
                                </span>
                              </td>
                            </tr>
                          ))}
                        </tbody>
                      </table>
                    </div>
                  </div>
                )}

                {(!salesReport.items || salesReport.items.length === 0) && 
                 (!salesReport.by_payment_type || salesReport.by_payment_type.length === 0) && 
                 (!salesReport.by_customer || salesReport.by_customer.length === 0) && 
                 (!salesReport.top_products || salesReport.top_products.length === 0) && 
                 (!salesReport.daily_breakdown || salesReport.daily_breakdown.length === 0) && (
                  <div className="text-center py-12 text-gray-500">
                    No hay datos de ventas en el período seleccionado
                  </div>
                )}
              </>
            ) : (
              <div className="text-center py-12 text-gray-500">
                No se pudo cargar el informe de ventas
              </div>
            )}
          </div>
        )}

        {activeTab === 'taxes' && (
          <div className="p-6">
            <div className="flex justify-between items-center mb-6">
              <div>
                <h2 className="text-2xl font-bold text-gray-900">Reporte de Impuestos (ITBIS)</h2>
                <p className="text-gray-600 mt-1">
                  Informe detallado de ITBIS cobrado, pagado y neto a pagar
                </p>
              </div>
              <div className="flex gap-3">
                <input
                  type="date"
                  value={taxStartDate}
                  onChange={(e) => setTaxStartDate(e.target.value)}
                  className="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#FF6B00] focus:border-transparent"
                />
                <input
                  type="date"
                  value={taxEndDate}
                  onChange={(e) => setTaxEndDate(e.target.value)}
                  className="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#FF6B00] focus:border-transparent"
                />
                <button
                  onClick={fetchTaxReport}
                  disabled={loading}
                  className="flex items-center gap-2 px-4 py-2 bg-[#FF6B00] text-white rounded-lg hover:bg-[#E55A00] transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  <RefreshCw className={`w-4 h-4 ${loading ? 'animate-spin' : ''}`} />
                  Actualizar
                </button>
              </div>
            </div>

            {loading ? (
              <div className="flex justify-center items-center py-12">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-[#FF6B00]"></div>
              </div>
            ) : taxReport ? (
              <>
                {/* Summary Cards - ITBIS Principal */}
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
                  <div className="bg-green-50 p-4 rounded-lg border border-green-200">
                    <div className="text-sm text-green-600 font-medium">ITBIS Cobrado</div>
                    <div className="text-2xl font-bold text-green-900 mt-1">
                      {formatCurrency(taxReport.summary?.total_itbis_collected || 0)}
                    </div>
                    <div className="text-xs text-green-700 mt-2">
                      Ventas + Facturas
                    </div>
                  </div>
                  <div className="bg-red-50 p-4 rounded-lg border border-red-200">
                    <div className="text-sm text-red-600 font-medium">ITBIS Pagado</div>
                    <div className="text-2xl font-bold text-red-900 mt-1">
                      {formatCurrency(taxReport.summary?.total_itbis_paid || 0)}
                    </div>
                    <div className="text-xs text-red-700 mt-2">
                      Compras (Crédito Fiscal)
                    </div>
                  </div>
                  <div className="bg-blue-50 p-4 rounded-lg border border-blue-200">
                    <div className="text-sm text-blue-600 font-medium">ITBIS Neto a Pagar</div>
                    <div className="text-2xl font-bold text-blue-900 mt-1">
                      {formatCurrency(taxReport.summary?.total_itbis_net || 0)}
                    </div>
                    <div className="text-xs text-blue-700 mt-2">
                      Cobrado - Pagado - Reducciones
                    </div>
                  </div>
                </div>

                {/* Retenciones Summary Cards */}
                {(taxReport.summary?.total_withholding_tax > 0 || taxReport.by_withholding?.length > 0) && (
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
                    <div className="bg-amber-50 p-4 rounded-lg border border-amber-200">
                      <div className="text-sm text-amber-600 font-medium">Retenciones Aplicadas</div>
                      <div className="text-2xl font-bold text-amber-900 mt-1">
                        {formatCurrency(taxReport.summary?.total_withholding_tax || 0)}
                      </div>
                      <div className="text-xs text-amber-700 mt-2">
                        ISR retenido en facturas gubernamentales
                      </div>
                    </div>
                    <div className="bg-emerald-50 p-4 rounded-lg border border-emerald-200">
                      <div className="text-sm text-emerald-600 font-medium">Monto Neto Recibido</div>
                      <div className="text-2xl font-bold text-emerald-900 mt-1">
                        {formatCurrency(taxReport.summary?.total_net_amount || 0)}
                      </div>
                      <div className="text-xs text-emerald-700 mt-2">
                        Total después de retenciones
                      </div>
                    </div>
                  </div>
                )}

                {/* Document Counts */}
                <div className="grid grid-cols-1 md:grid-cols-5 gap-4 mb-6">
                  <div className="bg-gray-50 p-4 rounded-lg border border-gray-200">
                    <div className="text-sm text-gray-600 font-medium">Ventas</div>
                    <div className="text-xl font-bold text-gray-900 mt-1">
                      {taxReport.summary?.total_sales || 0}
                    </div>
                  </div>
                  <div className="bg-gray-50 p-4 rounded-lg border border-gray-200">
                    <div className="text-sm text-gray-600 font-medium">Facturas</div>
                    <div className="text-xl font-bold text-gray-900 mt-1">
                      {taxReport.summary?.total_invoices || 0}
                    </div>
                  </div>
                  <div className="bg-gray-50 p-4 rounded-lg border border-gray-200">
                    <div className="text-sm text-gray-600 font-medium">Compras</div>
                    <div className="text-xl font-bold text-gray-900 mt-1">
                      {taxReport.summary?.total_purchases || 0}
                    </div>
                  </div>
                  <div className="bg-gray-50 p-4 rounded-lg border border-gray-200">
                    <div className="text-sm text-gray-600 font-medium">Notas de Crédito</div>
                    <div className="text-xl font-bold text-gray-900 mt-1">
                      {taxReport.summary?.total_credit_notes || 0}
                    </div>
                  </div>
                  <div className="bg-gray-50 p-4 rounded-lg border border-gray-200">
                    <div className="text-sm text-gray-600 font-medium">Devoluciones</div>
                    <div className="text-xl font-bold text-gray-900 mt-1">
                      {taxReport.summary?.total_returns || 0}
                    </div>
                  </div>
                </div>

                {/* Date Range Info */}
                <div className="mb-6 text-sm text-gray-600">
                  Período: {taxReport.start_date} al {taxReport.end_date} | 
                  Fecha del reporte: {taxReport.report_date}
                </div>

                {/* Tax by Document Type */}
                {taxReport.by_document_type && taxReport.by_document_type.length > 0 && (
                  <div className="mb-6">
                    <h3 className="text-lg font-semibold text-gray-900 mb-4">ITBIS por Tipo de Documento</h3>
                    <div className="overflow-x-auto">
                      <table className="min-w-full divide-y divide-gray-200">
                        <thead className="bg-gray-50">
                          <tr>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Tipo de Documento</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Cantidad</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Subtotal</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">ITBIS</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Total</th>
                          </tr>
                        </thead>
                        <tbody className="bg-white divide-y divide-gray-200">
                          {taxReport.by_document_type.map((item, index) => (
                            <tr key={index} className="hover:bg-gray-50">
                              <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                                {item.document_type === 'sales' ? 'Ventas' :
                                 item.document_type === 'invoices' ? 'Facturas' :
                                 item.document_type === 'purchases' ? 'Compras' :
                                 item.document_type === 'credit_notes' ? 'Notas de Crédito' :
                                 item.document_type === 'returns' ? 'Devoluciones' :
                                 item.document_type}
                              </td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">{item.count}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">{formatCurrency(item.subtotal)}</td>
                              <td className={`px-6 py-4 whitespace-nowrap text-sm font-semibold text-right ${
                                item.document_type === 'purchases' ? 'text-red-600' :
                                item.document_type === 'credit_notes' || item.document_type === 'returns' ? 'text-orange-600' :
                                'text-green-600'
                              }`}>
                                {formatCurrency(item.tax_amount)}
                              </td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">{formatCurrency(item.total)}</td>
                            </tr>
                          ))}
                        </tbody>
                        <tfoot className="bg-gray-50">
                          <tr>
                            <td className="px-6 py-4 whitespace-nowrap text-sm font-semibold text-gray-900">Total</td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm font-semibold text-right text-gray-900">
                              {taxReport.by_document_type.reduce((sum, item) => sum + item.count, 0)}
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm font-semibold text-right text-gray-900">
                              {formatCurrency(taxReport.by_document_type.reduce((sum, item) => sum + item.subtotal, 0))}
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm font-semibold text-right text-blue-600">
                              {formatCurrency(taxReport.summary?.total_itbis_net || 0)}
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm font-semibold text-right text-gray-900">
                              {formatCurrency(taxReport.by_document_type.reduce((sum, item) => sum + item.total, 0))}
                            </td>
                          </tr>
                        </tfoot>
                      </table>
                    </div>
                  </div>
                )}

                {/* Retenciones por Tipo */}
                {taxReport.by_withholding && taxReport.by_withholding.length > 0 && (
                  <div className="mb-6">
                    <h3 className="text-lg font-semibold text-gray-900 mb-4">Desglose de Retenciones</h3>
                    <div className="overflow-x-auto">
                      <table className="min-w-full divide-y divide-gray-200">
                        <thead className="bg-gray-50">
                          <tr>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Tipo de Retención</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Documentos</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Monto Total</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Total Retenido</th>
                          </tr>
                        </thead>
                        <tbody className="bg-white divide-y divide-gray-200">
                          {taxReport.by_withholding.map((item, index) => (
                            <tr key={index} className="hover:bg-gray-50">
                              <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                                {item.withholding_type === 'isr' ? 'ISR (Impuesto sobre la Renta)' :
                                 item.withholding_type === 'itbis' ? 'ITBIS' :
                                 item.withholding_type}
                              </td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">{item.count}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">{formatCurrency(item.total_amount)}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm font-semibold text-right text-amber-600">
                                {formatCurrency(item.total_withheld)}
                              </td>
                            </tr>
                          ))}
                        </tbody>
                        <tfoot className="bg-gray-50">
                          <tr>
                            <td className="px-6 py-4 whitespace-nowrap text-sm font-semibold text-gray-900">Total</td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm font-semibold text-right text-gray-900">
                              {taxReport.by_withholding.reduce((sum, item) => sum + item.count, 0)}
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm font-semibold text-right text-gray-900">
                              {formatCurrency(taxReport.by_withholding.reduce((sum, item) => sum + item.total_amount, 0))}
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm font-semibold text-right text-amber-600">
                              {formatCurrency(taxReport.summary?.total_withholding_tax || 0)}
                            </td>
                          </tr>
                        </tfoot>
                      </table>
                    </div>
                  </div>
                )}

                {/* Tax by NCF Type */}
                {taxReport.by_ncf_type && taxReport.by_ncf_type.length > 0 && (
                  <div className="mb-6">
                    <h3 className="text-lg font-semibold text-gray-900 mb-4">ITBIS por Tipo de NCF</h3>
                    <div className="overflow-x-auto">
                      <table className="min-w-full divide-y divide-gray-200">
                        <thead className="bg-gray-50">
                          <tr>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Tipo NCF</th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Descripción</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Cantidad</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Subtotal</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">ITBIS</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Total</th>
                          </tr>
                        </thead>
                        <tbody className="bg-white divide-y divide-gray-200">
                          {taxReport.by_ncf_type.map((item, index) => (
                            <tr key={index} className="hover:bg-gray-50">
                              <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{item.ncf_type}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">{item.description || 'N/A'}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">{item.count}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">{formatCurrency(item.subtotal)}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm font-semibold text-right text-green-600">{formatCurrency(item.tax_amount)}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">{formatCurrency(item.total)}</td>
                            </tr>
                          ))}
                        </tbody>
                      </table>
                    </div>
                  </div>
                )}

                {/* Detailed Breakdown by Document Type */}
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-6">
                  {/* Sales Tax */}
                  <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-4">
                    <h4 className="text-sm font-semibold text-gray-900 mb-3">Ventas</h4>
                    <div className="space-y-2 text-sm">
                      <div className="flex justify-between">
                        <span className="text-gray-600">Documentos:</span>
                        <span className="font-medium">{taxReport.sales_tax?.count || 0}</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-600">Subtotal:</span>
                        <span>{formatCurrency(taxReport.sales_tax?.subtotal || 0)}</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-600">ITBIS:</span>
                        <span className="font-semibold text-green-600">{formatCurrency(taxReport.sales_tax?.tax_amount || 0)}</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-600">Total:</span>
                        <span>{formatCurrency(taxReport.sales_tax?.total || 0)}</span>
                      </div>
                    </div>
                  </div>

                  {/* Invoices Tax */}
                  <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-4">
                    <h4 className="text-sm font-semibold text-gray-900 mb-3">Facturas</h4>
                    <div className="space-y-2 text-sm">
                      <div className="flex justify-between">
                        <span className="text-gray-600">Documentos:</span>
                        <span className="font-medium">{taxReport.invoices_tax?.count || 0}</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-600">Subtotal:</span>
                        <span>{formatCurrency(taxReport.invoices_tax?.subtotal || 0)}</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-600">ITBIS:</span>
                        <span className="font-semibold text-green-600">{formatCurrency(taxReport.invoices_tax?.tax_amount || 0)}</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-600">Total:</span>
                        <span>{formatCurrency(taxReport.invoices_tax?.total || 0)}</span>
                      </div>
                    </div>
                  </div>

                  {/* Purchases Tax */}
                  <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-4">
                    <h4 className="text-sm font-semibold text-gray-900 mb-3">Compras (Crédito Fiscal)</h4>
                    <div className="space-y-2 text-sm">
                      <div className="flex justify-between">
                        <span className="text-gray-600">Documentos:</span>
                        <span className="font-medium">{taxReport.purchases_tax?.count || 0}</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-600">Subtotal:</span>
                        <span>{formatCurrency(taxReport.purchases_tax?.subtotal || 0)}</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-600">ITBIS:</span>
                        <span className="font-semibold text-red-600">{formatCurrency(taxReport.purchases_tax?.tax_amount || 0)}</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-600">Total:</span>
                        <span>{formatCurrency(taxReport.purchases_tax?.total || 0)}</span>
                      </div>
                    </div>
                  </div>

                  {/* Credit Notes Tax */}
                  <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-4">
                    <h4 className="text-sm font-semibold text-gray-900 mb-3">Notas de Crédito</h4>
                    <div className="space-y-2 text-sm">
                      <div className="flex justify-between">
                        <span className="text-gray-600">Documentos:</span>
                        <span className="font-medium">{taxReport.credit_notes_tax?.count || 0}</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-600">Subtotal:</span>
                        <span>{formatCurrency(taxReport.credit_notes_tax?.subtotal || 0)}</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-600">ITBIS:</span>
                        <span className="font-semibold text-orange-600">{formatCurrency(taxReport.credit_notes_tax?.tax_amount || 0)}</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-600">Total:</span>
                        <span>{formatCurrency(taxReport.credit_notes_tax?.total || 0)}</span>
                      </div>
                    </div>
                  </div>

                  {/* Returns Tax */}
                  <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-4">
                    <h4 className="text-sm font-semibold text-gray-900 mb-3">Devoluciones</h4>
                    <div className="space-y-2 text-sm">
                      <div className="flex justify-between">
                        <span className="text-gray-600">Documentos:</span>
                        <span className="font-medium">{taxReport.returns_tax?.count || 0}</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-600">Subtotal:</span>
                        <span>{formatCurrency(taxReport.returns_tax?.subtotal || 0)}</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-600">ITBIS:</span>
                        <span className="font-semibold text-orange-600">{formatCurrency(taxReport.returns_tax?.tax_amount || 0)}</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-600">Total:</span>
                        <span>{formatCurrency(taxReport.returns_tax?.total || 0)}</span>
                      </div>
                    </div>
                  </div>
                </div>

                {/* Daily Breakdown */}
                {taxReport.daily_breakdown && taxReport.daily_breakdown.length > 0 && (
                  <div className="mb-6">
                    <h3 className="text-lg font-semibold text-gray-900 mb-4">Desglose Diario de ITBIS</h3>
                    <div className="overflow-x-auto">
                      <table className="min-w-full divide-y divide-gray-200">
                        <thead className="bg-gray-50">
                          <tr>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Fecha</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">ITBIS Ventas</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">ITBIS Facturas</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">ITBIS Compras</th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">ITBIS Neto del Día</th>
                          </tr>
                        </thead>
                        <tbody className="bg-white divide-y divide-gray-200">
                          {taxReport.daily_breakdown.map((day, index) => (
                            <tr key={index} className="hover:bg-gray-50">
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{day.date}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-green-600">{formatCurrency(day.sales_tax)}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-green-600">{formatCurrency(day.invoices_tax)}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-red-600">{formatCurrency(day.purchases_tax)}</td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm font-semibold text-right text-blue-600">{formatCurrency(day.total_tax)}</td>
                            </tr>
                          ))}
                        </tbody>
                      </table>
                    </div>
                  </div>
                )}

                {(!taxReport.by_document_type || taxReport.by_document_type.length === 0) && 
                 (!taxReport.by_ncf_type || taxReport.by_ncf_type.length === 0) && 
                 (!taxReport.by_withholding || taxReport.by_withholding.length === 0) &&
                 (!taxReport.daily_breakdown || taxReport.daily_breakdown.length === 0) && (
                  <div className="text-center py-12 text-gray-500">
                    No hay datos de impuestos en el período seleccionado
                  </div>
                )}
              </>
            ) : (
              <div className="text-center py-12 text-gray-500">
                No se pudo cargar el informe de impuestos
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  )
}
