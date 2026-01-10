import { useState, useEffect } from 'react'
import api from '../lib/api'
import Alert from '../components/Alert'
import { RotateCcw, Plus, Eye, CheckCircle, XCircle, Clock, Search, X, RefreshCw } from 'lucide-react'

export default function Returns() {
  const [returns, setReturns] = useState([])
  const [sales, setSales] = useState([])
  const [alert, setAlert] = useState(null)
  const [loading, setLoading] = useState(false)
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [showDetailsModal, setShowDetailsModal] = useState(false)
  const [selectedReturn, setSelectedReturn] = useState(null)
  const [selectedSale, setSelectedSale] = useState(null)
  const [returnLines, setReturnLines] = useState([])
  const [searchTerm, setSearchTerm] = useState('')
  const [statusFilter, setStatusFilter] = useState('')
  const [rejectReason, setRejectReason] = useState('')
  const [showRejectModal, setShowRejectModal] = useState(false)
  const [processReturn, setProcessReturn] = useState(null)
  const [showProcessModal, setShowProcessModal] = useState(false)
  const [refundMethod, setRefundMethod] = useState('cash')
  const [refundAmount, setRefundAmount] = useState('')

  useEffect(() => {
    loadReturns()
  }, [statusFilter])

  const loadReturns = async () => {
    try {
      setLoading(true)
      const params = new URLSearchParams()
      if (statusFilter) {
        params.append('status', statusFilter)
      }
      const url = `/returns${params.toString() ? '?' + params.toString() : ''}`
      const response = await api.get(url)
      setReturns(response.data || [])
    } catch (error) {
      console.error('Error loading returns:', error)
      showAlert('error', 'Error', error.response?.data?.error || 'Error al cargar devoluciones')
    } finally {
      setLoading(false)
    }
  }

  const loadSales = async () => {
    try {
      const response = await api.get('/sales')
      setSales(response.data || [])
    } catch (error) {
      console.error('Error loading sales:', error)
    }
  }

  const showAlert = (type, title, message) => {
    setAlert({ type, title, message })
    setTimeout(() => setAlert(null), 5000)
  }

  const formatCurrency = (amount) => {
    return new Intl.NumberFormat('es-DO', {
      style: 'currency',
      currency: 'DOP',
      minimumFractionDigits: 2,
    }).format(amount || 0)
  }

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString('es-DO', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
    })
  }

  const getStatusBadge = (status) => {
    const badges = {
      pending: { bg: 'bg-yellow-100', text: 'text-yellow-800', label: 'Pendiente', icon: Clock },
      approved: { bg: 'bg-blue-100', text: 'text-blue-800', label: 'Aprobada', icon: CheckCircle },
      rejected: { bg: 'bg-red-100', text: 'text-red-800', label: 'Rechazada', icon: XCircle },
      completed: { bg: 'bg-green-100', text: 'text-green-800', label: 'Completada', icon: CheckCircle },
      refunded: { bg: 'bg-purple-100', text: 'text-purple-800', label: 'Reembolsada', icon: RotateCcw },
    }
    const badge = badges[status] || badges.pending
    const Icon = badge.icon
    return (
      <span className={`inline-flex items-center gap-1 px-2 py-1 rounded-full text-xs font-medium ${badge.bg} ${badge.text}`}>
        <Icon className="w-3 h-3" />
        {badge.label}
      </span>
    )
  }

  const handleCreateReturn = () => {
    setSelectedSale(null)
    setReturnLines([])
    loadSales()
    setShowCreateModal(true)
  }

  const selectSale = (sale) => {
    setSelectedSale(sale)
    // Initialize return lines with all sale lines, setting quantity to 0
    const lines = (sale.lines || []).map((line, index) => ({
      ...line,
      returnQuantity: 0,
      itemCondition: null,
      restockQuantity: 0,
    }))
    setReturnLines(lines)
  }

  const updateReturnLine = (lineId, field, value) => {
    setReturnLines(returnLines.map(line => {
      if (line.id === lineId) {
        const updated = { ...line, [field]: value }
        // Auto-set restock quantity based on condition
        if (field === 'itemCondition') {
          if (value === 'new' || value === 'used') {
            updated.restockQuantity = updated.returnQuantity || 0
          } else if (value === 'defective' || value === 'damaged') {
            updated.restockQuantity = 0
          }
        }
        // Validate return quantity doesn't exceed available
        if (field === 'returnQuantity') {
          const available = line.quantity - (line.returned_quantity || 0)
          if (value > available) {
            showAlert('error', 'Error', `La cantidad a devolver (${value}) excede la cantidad disponible (${available})`)
            return line
          }
        }
        return updated
      }
      return line
    }))
  }

  const handleSubmitReturn = async () => {
    // Validate at least one line has quantity > 0
    const validLines = returnLines.filter(line => line.returnQuantity > 0)
    if (validLines.length === 0) {
      showAlert('error', 'Error', 'Debe seleccionar al menos una línea con cantidad a devolver')
      return
    }

    // Validate quantities
    for (const line of validLines) {
      const available = line.quantity - (line.returned_quantity || 0)
      if (line.returnQuantity > available) {
        showAlert('error', 'Error', `La cantidad a devolver excede la disponible para ${line.description}`)
        return
      }
    }

    try {
      setLoading(true)
      const payload = {
        sale_id: selectedSale.id,
        lines: validLines.map(line => ({
          sale_line_id: line.id,
          quantity: parseFloat(line.returnQuantity),
          item_condition: line.itemCondition || 'used',
          restock_quantity: line.restockQuantity || 0,
        })),
      }

      await api.post('/returns', payload)
      showAlert('success', 'Éxito', 'Devolución creada exitosamente')
      setShowCreateModal(false)
      setSelectedSale(null)
      setReturnLines([])
      loadReturns()
    } catch (error) {
      console.error('Error creating return:', error)
      showAlert('error', 'Error', error.response?.data?.error || 'Error al crear devolución')
    } finally {
      setLoading(false)
    }
  }

  const handleApprove = async (returnId) => {
    if (!window.confirm('¿Está seguro de aprobar esta devolución?')) {
      return
    }

    try {
      setLoading(true)
      await api.patch(`/returns/${returnId}/approve`)
      showAlert('success', 'Éxito', 'Devolución aprobada exitosamente')
      loadReturns()
    } catch (error) {
      console.error('Error approving return:', error)
      showAlert('error', 'Error', error.response?.data?.error || 'Error al aprobar devolución')
    } finally {
      setLoading(false)
    }
  }

  const handleReject = async () => {
    if (!rejectReason.trim()) {
      showAlert('error', 'Error', 'Debe proporcionar una razón para rechazar')
      return
    }

    try {
      setLoading(true)
      await api.patch(`/returns/${selectedReturn.id}/reject`, {
        reason: rejectReason,
      })
      showAlert('success', 'Éxito', 'Devolución rechazada exitosamente')
      setShowRejectModal(false)
      setRejectReason('')
      setSelectedReturn(null)
      loadReturns()
    } catch (error) {
      console.error('Error rejecting return:', error)
      showAlert('error', 'Error', error.response?.data?.error || 'Error al rechazar devolución')
    } finally {
      setLoading(false)
    }
  }

  const handleProcess = async () => {
    if (!refundMethod) {
      showAlert('error', 'Error', 'Debe seleccionar un método de reembolso')
      return
    }

    const amount = refundAmount ? parseFloat(refundAmount) : null
    if (amount !== null && amount > processReturn.total) {
      showAlert('error', 'Error', 'El monto de reembolso no puede exceder el total de la devolución')
      return
    }

    try {
      setLoading(true)
      await api.post(`/returns/${processReturn.id}/process`, {
        refund_method: refundMethod,
        refund_amount: amount,
      })
      showAlert('success', 'Éxito', 'Devolución procesada exitosamente')
      setShowProcessModal(false)
      setProcessReturn(null)
      setRefundMethod('cash')
      setRefundAmount('')
      loadReturns()
    } catch (error) {
      console.error('Error processing return:', error)
      showAlert('error', 'Error', error.response?.data?.error || 'Error al procesar devolución')
    } finally {
      setLoading(false)
    }
  }

  const viewDetails = async (returnId) => {
    try {
      const response = await api.get(`/returns/${returnId}`)
      setSelectedReturn(response.data)
      setShowDetailsModal(true)
    } catch (error) {
      console.error('Error loading return details:', error)
      showAlert('error', 'Error', 'Error al cargar detalles de la devolución')
    }
  }

  const openRejectModal = (ret) => {
    setSelectedReturn(ret)
    setRejectReason('')
    setShowRejectModal(true)
  }

  const openProcessModal = (ret) => {
    setProcessReturn(ret)
    setRefundMethod('cash')
    setRefundAmount('')
    setShowProcessModal(true)
  }

  const filteredReturns = returns.filter(ret => {
    const matchesSearch = 
      ret.return_number?.toLowerCase().includes(searchTerm.toLowerCase()) ||
      ret.sale_number?.toLowerCase().includes(searchTerm.toLowerCase()) ||
      ret.customer_name?.toLowerCase().includes(searchTerm.toLowerCase())
    return matchesSearch
  })

  return (
    <div className="max-w-7xl">
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-[#212121] flex items-center gap-3">
          <RotateCcw className="w-8 h-8 text-[#FF6B00]" />
          Devoluciones
        </h1>
        <p className="mt-2 text-gray-600">
          Gestión de devoluciones de productos vendidos
        </p>
      </div>

      {alert && (
        <div className="mb-6">
          <Alert
            type={alert.type}
            title={alert.title}
            message={alert.message}
            onClose={() => setAlert(null)}
          />
        </div>
      )}

      {/* Filters and Actions */}
      <div className="bg-white rounded-lg shadow p-6 mb-6">
        <div className="flex flex-col md:flex-row gap-4 items-center justify-between">
          <div className="flex flex-1 gap-4 w-full md:w-auto">
            <div className="relative flex-1 md:flex-initial">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
              <input
                type="text"
                placeholder="Buscar por número, venta, cliente..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="input-field pl-10 w-full md:w-80"
              />
            </div>
            <select
              value={statusFilter}
              onChange={(e) => setStatusFilter(e.target.value)}
              className="input-field"
            >
              <option value="">Todos los estados</option>
              <option value="pending">Pendiente</option>
              <option value="approved">Aprobada</option>
              <option value="rejected">Rechazada</option>
              <option value="completed">Completada</option>
              <option value="refunded">Reembolsada</option>
            </select>
          </div>
          <div className="flex gap-2">
            <button
              onClick={loadReturns}
              className="flex items-center gap-2 px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50"
            >
              <RefreshCw className="w-4 h-4" />
              Actualizar
            </button>
            <button
              onClick={handleCreateReturn}
              className="flex items-center gap-2 px-4 py-2 bg-[#FF6B00] text-white rounded-lg hover:bg-[#E55A00]"
            >
              <Plus className="w-5 h-5" />
              Nueva Devolución
            </button>
          </div>
        </div>
      </div>

      {/* Returns Table */}
      <div className="bg-white rounded-lg shadow overflow-hidden">
        {loading && !returns.length ? (
          <div className="flex justify-center items-center py-12">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-[#FF6B00]"></div>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Número</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Fecha</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Venta</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Cliente</th>
                  <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Total</th>
                  <th className="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase">Estado</th>
                  <th className="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase">Acciones</th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {filteredReturns.length === 0 ? (
                  <tr>
                    <td colSpan="7" className="px-6 py-8 text-center text-gray-500">
                      No hay devoluciones registradas
                    </td>
                  </tr>
                ) : (
                  filteredReturns.map((ret) => (
                    <tr key={ret.id} className="hover:bg-gray-50">
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                        {ret.return_number}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                        {formatDate(ret.return_date)}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                        {ret.sale_number || '-'}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                        {ret.customer_name || 'Cliente Genérico'}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">
                        {formatCurrency(ret.total)}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-center">
                        {getStatusBadge(ret.status)}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-center">
                        <div className="flex items-center justify-center gap-2">
                          <button
                            onClick={() => viewDetails(ret.id)}
                            className="p-1 text-blue-600 hover:text-blue-800 hover:bg-blue-50 rounded"
                            title="Ver detalles"
                          >
                            <Eye className="w-4 h-4" />
                          </button>
                          {ret.status === 'pending' && (
                            <>
                              <button
                                onClick={() => handleApprove(ret.id)}
                                className="p-1 text-green-600 hover:text-green-800 hover:bg-green-50 rounded"
                                title="Aprobar"
                              >
                                <CheckCircle className="w-4 h-4" />
                              </button>
                              <button
                                onClick={() => openRejectModal(ret)}
                                className="p-1 text-red-600 hover:text-red-800 hover:bg-red-50 rounded"
                                title="Rechazar"
                              >
                                <XCircle className="w-4 h-4" />
                              </button>
                            </>
                          )}
                          {ret.status === 'approved' && (
                            <button
                              onClick={() => openProcessModal(ret)}
                              className="p-1 text-purple-600 hover:text-purple-800 hover:bg-purple-50 rounded"
                              title="Procesar"
                            >
                              <RotateCcw className="w-4 h-4" />
                            </button>
                          )}
                        </div>
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* Create Return Modal */}
      {showCreateModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-lg shadow-xl max-w-4xl w-full max-h-[90vh] overflow-y-auto">
            <div className="flex justify-between items-center p-6 border-b border-gray-200">
              <h2 className="text-xl font-semibold text-gray-900">Nueva Devolución</h2>
              <button
                onClick={() => {
                  setShowCreateModal(false)
                  setSelectedSale(null)
                  setReturnLines([])
                }}
                className="text-gray-400 hover:text-gray-600"
              >
                <X className="w-6 h-6" />
              </button>
            </div>

            <div className="p-6">
              {!selectedSale ? (
                <div>
                  <h3 className="text-lg font-medium mb-4">Seleccionar Venta</h3>
                  <div className="space-y-2 max-h-96 overflow-y-auto">
                    {sales.map((sale) => (
                      <button
                        key={sale.id}
                        onClick={() => selectSale(sale)}
                        className="w-full text-left p-4 border border-gray-200 rounded-lg hover:bg-gray-50 hover:border-[#FF6B00] transition-colors"
                      >
                        <div className="flex justify-between items-center">
                          <div>
                            <div className="font-medium text-gray-900">{sale.sale_number}</div>
                            <div className="text-sm text-gray-500">
                              {formatDate(sale.created_at)} - {sale.customer_name || 'Cliente Genérico'}
                            </div>
                          </div>
                          <div className="text-right">
                            <div className="font-semibold text-gray-900">{formatCurrency(sale.total)}</div>
                            <div className="text-xs text-gray-500">{sale.lines?.length || 0} líneas</div>
                          </div>
                        </div>
                      </button>
                    ))}
                  </div>
                </div>
              ) : (
                <div>
                  <div className="mb-4 pb-4 border-b border-gray-200">
                    <div className="flex justify-between items-center">
                      <div>
                        <h3 className="text-lg font-medium">Venta: {selectedSale.sale_number}</h3>
                        <p className="text-sm text-gray-500">
                          {formatDate(selectedSale.created_at)} - {selectedSale.customer_name || 'Cliente Genérico'}
                        </p>
                      </div>
                      <button
                        onClick={() => {
                          setSelectedSale(null)
                          setReturnLines([])
                        }}
                        className="text-sm text-blue-600 hover:text-blue-800"
                      >
                        Cambiar venta
                      </button>
                    </div>
                  </div>

                  <div className="space-y-4">
                    <h4 className="font-medium">Líneas a Devolver</h4>
                    {returnLines.map((line, index) => {
                      const available = line.quantity - (line.returned_quantity || 0)
                      return (
                        <div key={line.id} className="border border-gray-200 rounded-lg p-4">
                          <div className="flex justify-between items-start mb-2">
                            <div className="flex-1">
                              <div className="font-medium">{line.description}</div>
                              <div className="text-sm text-gray-500">
                                Precio: {formatCurrency(line.unit_price)} | 
                                Vendido: {line.quantity} | 
                                Ya devuelto: {(line.returned_quantity || 0).toFixed(2)} | 
                                Disponible: {available.toFixed(2)}
                              </div>
                            </div>
                          </div>
                          <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mt-3">
                            <div>
                              <label className="block text-sm font-medium text-gray-700 mb-1">
                                Cantidad a Devolver *
                              </label>
                              <input
                                type="number"
                                min="0"
                                max={available}
                                step="0.01"
                                value={line.returnQuantity || ''}
                                onChange={(e) => updateReturnLine(line.id, 'returnQuantity', parseFloat(e.target.value) || 0)}
                                className="input-field"
                                placeholder="0.00"
                              />
                              <p className="text-xs text-gray-500 mt-1">Máximo: {available.toFixed(2)}</p>
                            </div>
                            <div>
                              <label className="block text-sm font-medium text-gray-700 mb-1">
                                Condición
                              </label>
                              <select
                                value={line.itemCondition || ''}
                                onChange={(e) => updateReturnLine(line.id, 'itemCondition', e.target.value)}
                                className="input-field"
                              >
                                <option value="">Seleccionar...</option>
                                <option value="new">Nuevo</option>
                                <option value="used">Usado</option>
                                <option value="defective">Defectuoso</option>
                                <option value="damaged">Dañado</option>
                              </select>
                            </div>
                            <div>
                              <label className="block text-sm font-medium text-gray-700 mb-1">
                                Cantidad a Reabastecer
                              </label>
                              <input
                                type="number"
                                min="0"
                                max={line.returnQuantity || 0}
                                step="0.01"
                                value={line.restockQuantity || ''}
                                onChange={(e) => updateReturnLine(line.id, 'restockQuantity', parseFloat(e.target.value) || 0)}
                                className="input-field"
                                disabled={!line.itemCondition}
                              />
                            </div>
                          </div>
                        </div>
                      )
                    })}
                  </div>

                  <div className="mt-6 flex justify-end gap-3">
                    <button
                      onClick={() => {
                        setShowCreateModal(false)
                        setSelectedSale(null)
                        setReturnLines([])
                      }}
                      className="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50"
                    >
                      Cancelar
                    </button>
                    <button
                      onClick={handleSubmitReturn}
                      disabled={loading}
                      className="px-4 py-2 bg-[#FF6B00] text-white rounded-lg hover:bg-[#E55A00] disabled:opacity-50"
                    >
                      {loading ? 'Guardando...' : 'Crear Devolución'}
                    </button>
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>
      )}

      {/* Details Modal */}
      {showDetailsModal && selectedReturn && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-lg shadow-xl max-w-3xl w-full max-h-[90vh] overflow-y-auto">
            <div className="flex justify-between items-center p-6 border-b border-gray-200">
              <h2 className="text-xl font-semibold text-gray-900">Detalles de Devolución</h2>
              <button
                onClick={() => {
                  setShowDetailsModal(false)
                  setSelectedReturn(null)
                }}
                className="text-gray-400 hover:text-gray-600"
              >
                <X className="w-6 h-6" />
              </button>
            </div>

            <div className="p-6 space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="text-sm font-medium text-gray-500">Número</label>
                  <div className="text-lg font-semibold">{selectedReturn.return_number}</div>
                </div>
                <div>
                  <label className="text-sm font-medium text-gray-500">Estado</label>
                  <div>{getStatusBadge(selectedReturn.status)}</div>
                </div>
                <div>
                  <label className="text-sm font-medium text-gray-500">Fecha</label>
                  <div>{formatDate(selectedReturn.return_date)}</div>
                </div>
                <div>
                  <label className="text-sm font-medium text-gray-500">Venta</label>
                  <div>{selectedReturn.sale_number || '-'}</div>
                </div>
                <div>
                  <label className="text-sm font-medium text-gray-500">Cliente</label>
                  <div>{selectedReturn.customer_name || 'Cliente Genérico'}</div>
                </div>
                <div>
                  <label className="text-sm font-medium text-gray-500">Total</label>
                  <div className="text-lg font-semibold">{formatCurrency(selectedReturn.total)}</div>
                </div>
              </div>

              {selectedReturn.lines && selectedReturn.lines.length > 0 && (
                <div>
                  <h3 className="font-medium mb-3">Líneas de Devolución</h3>
                  <div className="overflow-x-auto">
                    <table className="min-w-full divide-y divide-gray-200">
                      <thead className="bg-gray-50">
                        <tr>
                          <th className="px-4 py-2 text-left text-xs font-medium text-gray-500">Producto</th>
                          <th className="px-4 py-2 text-right text-xs font-medium text-gray-500">Cantidad</th>
                          <th className="px-4 py-2 text-right text-xs font-medium text-gray-500">Precio</th>
                          <th className="px-4 py-2 text-right text-xs font-medium text-gray-500">Total</th>
                          <th className="px-4 py-2 text-left text-xs font-medium text-gray-500">Condición</th>
                        </tr>
                      </thead>
                      <tbody className="bg-white divide-y divide-gray-200">
                        {selectedReturn.lines.map((line, idx) => (
                          <tr key={idx}>
                            <td className="px-4 py-2 text-sm">{line.description}</td>
                            <td className="px-4 py-2 text-sm text-right">{line.quantity}</td>
                            <td className="px-4 py-2 text-sm text-right">{formatCurrency(line.unit_price)}</td>
                            <td className="px-4 py-2 text-sm text-right font-medium">{formatCurrency(line.line_total)}</td>
                            <td className="px-4 py-2 text-sm">{line.item_condition || '-'}</td>
                          </tr>
                        ))}
                      </tbody>
                      <tfoot className="bg-gray-50">
                        <tr>
                          <td colSpan="3" className="px-4 py-2 text-sm font-medium text-right">Total:</td>
                          <td className="px-4 py-2 text-sm font-bold">{formatCurrency(selectedReturn.total)}</td>
                          <td></td>
                        </tr>
                      </tfoot>
                    </table>
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>
      )}

      {/* Reject Modal */}
      {showRejectModal && selectedReturn && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-lg shadow-xl max-w-md w-full">
            <div className="flex justify-between items-center p-6 border-b border-gray-200">
              <h2 className="text-xl font-semibold text-gray-900">Rechazar Devolución</h2>
              <button
                onClick={() => {
                  setShowRejectModal(false)
                  setSelectedReturn(null)
                  setRejectReason('')
                }}
                className="text-gray-400 hover:text-gray-600"
              >
                <X className="w-6 h-6" />
              </button>
            </div>

            <div className="p-6">
              <p className="mb-4 text-gray-600">
                Devolución: <strong>{selectedReturn.return_number}</strong>
              </p>
              <div className="mb-4">
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Razón del Rechazo *
                </label>
                <textarea
                  value={rejectReason}
                  onChange={(e) => setRejectReason(e.target.value)}
                  rows={4}
                  className="input-field"
                  placeholder="Explique por qué se rechaza esta devolución..."
                />
              </div>
              <div className="flex justify-end gap-3">
                <button
                  onClick={() => {
                    setShowRejectModal(false)
                    setSelectedReturn(null)
                    setRejectReason('')
                  }}
                  className="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50"
                >
                  Cancelar
                </button>
                <button
                  onClick={handleReject}
                  disabled={loading || !rejectReason.trim()}
                  className="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 disabled:opacity-50"
                >
                  {loading ? 'Rechazando...' : 'Rechazar'}
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Process Modal */}
      {showProcessModal && processReturn && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-lg shadow-xl max-w-md w-full">
            <div className="flex justify-between items-center p-6 border-b border-gray-200">
              <h2 className="text-xl font-semibold text-gray-900">Procesar Devolución</h2>
              <button
                onClick={() => {
                  setShowProcessModal(false)
                  setProcessReturn(null)
                  setRefundMethod('cash')
                  setRefundAmount('')
                }}
                className="text-gray-400 hover:text-gray-600"
              >
                <X className="w-6 h-6" />
              </button>
            </div>

            <div className="p-6 space-y-4">
              <div>
                <p className="text-sm text-gray-600 mb-2">
                  Devolución: <strong>{processReturn.return_number}</strong>
                </p>
                <p className="text-sm text-gray-600">
                  Total a reembolsar: <strong className="text-lg">{formatCurrency(processReturn.total)}</strong>
                </p>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Método de Reembolso *
                </label>
                <select
                  value={refundMethod}
                  onChange={(e) => setRefundMethod(e.target.value)}
                  className="input-field"
                >
                  <option value="cash">Efectivo</option>
                  <option value="card">Tarjeta</option>
                  <option value="transfer">Transferencia</option>
                  <option value="credit_note">Nota de Crédito</option>
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Monto a Reembolsar (opcional, por defecto: total)
                </label>
                <input
                  type="number"
                  step="0.01"
                  min="0"
                  max={processReturn.total}
                  value={refundAmount}
                  onChange={(e) => setRefundAmount(e.target.value)}
                  className="input-field"
                  placeholder={formatCurrency(processReturn.total)}
                />
                <p className="text-xs text-gray-500 mt-1">
                  Si se deja vacío, se reembolsará el total de la devolución
                </p>
              </div>

              <div className="flex justify-end gap-3 pt-4">
                <button
                  onClick={() => {
                    setShowProcessModal(false)
                    setProcessReturn(null)
                    setRefundMethod('cash')
                    setRefundAmount('')
                  }}
                  className="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50"
                >
                  Cancelar
                </button>
                <button
                  onClick={handleProcess}
                  disabled={loading}
                  className="px-4 py-2 bg-[#FF6B00] text-white rounded-lg hover:bg-[#E55A00] disabled:opacity-50"
                >
                  {loading ? 'Procesando...' : 'Procesar Devolución'}
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

