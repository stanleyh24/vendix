import { useState, useEffect } from 'react';
import { ShoppingCart, Search, Eye, Download, Calendar, User, DollarSign, Filter, RefreshCw } from 'lucide-react';
import Alert from '../components/Alert';
import api from '../lib/api';

export default function SalesManagement() {
  const [sales, setSales] = useState([]);
  const [loading, setLoading] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const [filterStatus, setFilterStatus] = useState('all');
  const [alert, setAlert] = useState(null);
  const [selectedSale, setSelectedSale] = useState(null);
  const [showDetailModal, setShowDetailModal] = useState(false);

  // Estados de venta
  const statuses = [
    { value: 'completed', label: 'Completada', color: 'bg-green-500', badge: 'badge-success' },
    { value: 'cancelled', label: 'Cancelada', color: 'bg-red-500', badge: 'badge-error' },
    { value: 'refunded', label: 'Reembolsada', color: 'bg-orange-500', badge: 'badge-warning' },
  ];

  // Tipos de pago
  const paymentTypes = {
    cash: { label: 'Efectivo', icon: '💵' },
    card: { label: 'Tarjeta', icon: '💳' },
    transfer: { label: 'Transferencia', icon: '🏦' },
    check: { label: 'Cheque', icon: '📝' },
  };

  // Cargar ventas del backend
  useEffect(() => {
    loadSales();
  }, []);

  const loadSales = async () => {
    setLoading(true);
    try {
      const response = await api.get('/sales');
      const data = response.data || [];
      
      // Map backend data to frontend format
      const mappedSales = data.map(sale => ({
        id: sale.id,
        sale_number: sale.sale_number,
        customer: { 
          id: sale.customer_id, 
          name: sale.customer_name || 'Cliente Genérico',
          tax_id: '' 
        },
        invoice_id: sale.invoice_id,
        payment_type: sale.payment_type,
        status: sale.status,
        subtotal: sale.subtotal,
        tax: sale.tax_amount,
        total: sale.total,
        items: sale.lines || [],
        created_at: sale.created_at
      }));
      
      setSales(mappedSales);
    } catch (error) {
      console.error('Error loading sales:', error);
      showAlert('error', 'Error', 'No se pudieron cargar las ventas');
      setSales([]);
    } finally {
      setLoading(false);
    }
  };

  const showAlert = (type, title, message) => {
    setAlert({ type, title, message });
    setTimeout(() => setAlert(null), 5000);
  };

  // Filtrar ventas
  const filteredSales = sales.filter(sale => {
    const matchesSearch = 
      sale.sale_number.toLowerCase().includes(searchTerm.toLowerCase()) ||
      sale.customer.name.toLowerCase().includes(searchTerm.toLowerCase());
    
    const matchesStatus = filterStatus === 'all' || sale.status === filterStatus;

    return matchesSearch && matchesStatus;
  });

  // Ver detalles de venta
  const handleViewDetails = async (sale) => {
    try {
      const response = await api.get(`/sales/${sale.id}`);
      const saleData = response.data;
      
      // Map backend data
      const mappedSale = {
        ...saleData,
        customer: { 
          id: saleData.customer_id, 
          name: saleData.customer_name || 'Cliente Genérico',
          tax_id: '' 
        },
        tax: saleData.tax_amount,
        items: saleData.lines || []
      };
      
      setSelectedSale(mappedSale);
      setShowDetailModal(true);
    } catch (error) {
      showAlert('error', 'Error', 'No se pudieron cargar los detalles de la venta');
    }
  };

  // Cerrar modal
  const handleCloseModal = () => {
    setShowDetailModal(false);
    setSelectedSale(null);
  };

  // Calcular totales
  const calculateTotals = () => {
    const total = filteredSales.reduce((sum, sale) => sum + sale.total, 0);
    const completed = filteredSales.filter(sale => sale.status === 'completed').reduce((sum, sale) => sum + sale.total, 0);
    const cancelled = filteredSales.filter(sale => sale.status === 'cancelled').reduce((sum, sale) => sum + sale.total, 0);
    
    return { total, completed, cancelled };
  };

  const totals = calculateTotals();

  const getStatusBadge = (status) => {
    return statuses.find(s => s.value === status)?.badge || 'badge';
  };

  const getStatusLabel = (status) => {
    return statuses.find(s => s.value === status)?.label || status;
  };

  const getPaymentTypeInfo = (paymentType) => {
    return paymentTypes[paymentType] || { label: paymentType, icon: '💰' };
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-[#212121] flex items-center gap-3">
            <ShoppingCart className="w-8 h-8 text-[#FF6B00]" />
            Gestión de Ventas
          </h1>
          <p className="text-gray-600 mt-1">Historial y gestión de ventas</p>
        </div>
        <button
          onClick={loadSales}
          disabled={loading}
          className="btn-outline flex items-center gap-2 disabled:opacity-50"
        >
          <RefreshCw className={`w-5 h-5 ${loading ? 'animate-spin' : ''}`} />
          Actualizar
        </button>
      </div>

      {/* Alert */}
      {alert && (
        <Alert
          type={alert.type}
          title={alert.title}
          message={alert.message}
          onClose={() => setAlert(null)}
        />
      )}

      {/* KPI Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="card">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600">Total Ventas</p>
              <p className="text-3xl font-bold text-[#FF6B00] mt-1">
                ${totals.total.toFixed(2)}
              </p>
            </div>
            <div className="w-12 h-12 bg-orange-100 rounded-full flex items-center justify-center">
              <DollarSign className="w-6 h-6 text-[#FF6B00]" />
            </div>
          </div>
        </div>

        <div className="card bg-gradient-to-br from-green-50 to-white">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600">Completadas</p>
              <p className="text-3xl font-bold text-[#00C853] mt-1">
                ${totals.completed.toFixed(2)}
              </p>
            </div>
            <div className="w-12 h-12 bg-green-100 rounded-full flex items-center justify-center">
              <ShoppingCart className="w-6 h-6 text-green-600" />
            </div>
          </div>
        </div>

        <div className="card bg-gradient-to-br from-red-50 to-white">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600">Canceladas</p>
              <p className="text-3xl font-bold text-[#F44336] mt-1">
                ${totals.cancelled.toFixed(2)}
              </p>
            </div>
            <div className="w-12 h-12 bg-red-100 rounded-full flex items-center justify-center">
              <Calendar className="w-6 h-6 text-red-600" />
            </div>
          </div>
        </div>
      </div>

      {/* Filtros y búsqueda */}
      <div className="card">
        <div className="flex flex-col gap-4">
          {/* Búsqueda */}
          <div className="flex-1">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
              <input
                type="text"
                placeholder="Buscar por número de venta o cliente..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="input-field pl-10"
              />
            </div>
          </div>

          {/* Filtros por estado */}
          <div className="flex flex-wrap gap-2">
            <button
              onClick={() => setFilterStatus('all')}
              className={`px-4 py-2 rounded-lg font-medium transition-colors text-sm ${
                filterStatus === 'all'
                  ? 'bg-[#FF6B00] text-white'
                  : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
              }`}
            >
              Todas ({sales.length})
            </button>
            {statuses.map(status => (
              <button
                key={status.value}
                onClick={() => setFilterStatus(status.value)}
                className={`px-4 py-2 rounded-lg font-medium transition-colors text-sm ${
                  filterStatus === status.value
                    ? `${status.color} text-white`
                    : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                }`}
              >
                {status.label} ({sales.filter(s => s.status === status.value).length})
              </button>
            ))}
          </div>
        </div>
      </div>

      {/* Lista de ventas */}
      {loading ? (
        <div className="card text-center py-12">
          <div className="inline-block w-8 h-8 border-4 border-[#FF6B00] border-t-transparent rounded-full animate-spin"></div>
          <p className="mt-4 text-gray-600">Cargando ventas...</p>
        </div>
      ) : filteredSales.length === 0 ? (
        <div className="card text-center py-12">
          <ShoppingCart className="w-16 h-16 text-gray-300 mx-auto mb-4" />
          <h3 className="text-lg font-semibold text-gray-600 mb-2">
            {searchTerm ? 'No se encontraron ventas' : 'No hay ventas'}
          </h3>
          <p className="text-gray-500 mb-6">
            {searchTerm ? 'Intenta con otro término de búsqueda' : 'Las ventas aparecerán cuando proceses transacciones'}
          </p>
        </div>
      ) : (
        <div className="space-y-3">
          {filteredSales.map((sale) => (
            <div
              key={sale.id}
              className="card hover:shadow-card-hover transition-all duration-200"
            >
              <div className="flex items-start justify-between">
                {/* Info principal */}
                <div className="flex-1">
                  <div className="flex items-center gap-3 mb-3">
                    <span className="font-mono font-bold text-lg text-[#FF6B00]">
                      {sale.sale_number}
                    </span>
                    <span className={getStatusBadge(sale.status)}>
                      {getStatusLabel(sale.status)}
                    </span>
                    <span className="badge bg-blue-100 text-blue-800">
                      {getPaymentTypeInfo(sale.payment_type).icon} {getPaymentTypeInfo(sale.payment_type).label}
                    </span>
                  </div>

                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
                    <div>
                      <div className="flex items-center gap-2 mb-2">
                        <User className="w-4 h-4 text-gray-400" />
                        <div>
                          <p className="font-semibold text-[#212121]">{sale.customer.name}</p>
                          <p className="text-sm text-gray-600">{sale.customer.tax_id}</p>
                        </div>
                      </div>
                    </div>

                    <div>
                      <div className="flex items-center gap-2">
                        <Calendar className="w-4 h-4 text-gray-400" />
                        <div>
                          <p className="text-sm text-gray-600">
                            Fecha: <span className="font-medium">{new Date(sale.created_at).toLocaleDateString('es-DO')}</span>
                          </p>
                          <p className="text-sm text-gray-600">
                            Hora: <span className="font-medium">{new Date(sale.created_at).toLocaleTimeString('es-DO')}</span>
                          </p>
                        </div>
                      </div>
                    </div>
                  </div>

                  {/* Items de la venta */}
                  <div className="bg-[#F5F5F5] rounded-lg p-3 mb-3">
                    <p className="text-xs font-semibold text-gray-600 mb-2">ITEMS:</p>
                    <div className="space-y-1">
                      {sale.items.map((item, idx) => (
                        <div key={idx} className="flex items-center justify-between text-sm">
                          <span className="text-gray-700">
                            {item.quantity}x {item.description}
                          </span>
                          <span className="font-medium text-gray-900">
                            ${(item.quantity * item.unit_price).toFixed(2)}
                          </span>
                        </div>
                      ))}
                    </div>
                  </div>

                  {/* Totales */}
                  <div className="grid grid-cols-3 gap-3 text-sm">
                    <div>
                      <p className="text-xs text-gray-600">Subtotal</p>
                      <p className="font-semibold text-gray-900">${sale.subtotal.toFixed(2)}</p>
                    </div>
                    <div>
                      <p className="text-xs text-gray-600">ITBIS</p>
                      <p className="font-semibold text-gray-900">${sale.tax.toFixed(2)}</p>
                    </div>
                    <div>
                      <p className="text-xs text-gray-600">Total</p>
                      <p className="font-bold text-lg text-[#FF6B00]">${sale.total.toFixed(2)}</p>
                    </div>
                  </div>
                </div>

                {/* Acciones */}
                <div className="flex flex-col gap-2 ml-4">
                  <button
                    onClick={() => handleViewDetails(sale)}
                    className="p-2 hover:bg-blue-50 rounded-lg transition-colors"
                    title="Ver detalles"
                  >
                    <Eye className="w-5 h-5 text-blue-600" />
                  </button>
                  
                  <button
                    onClick={() => showAlert('info', 'Descargando', `Generando PDF de venta ${sale.sale_number}`)}
                    className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
                    title="Descargar PDF"
                  >
                    <Download className="w-5 h-5 text-gray-600" />
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Modal de detalles */}
      {showDetailModal && selectedSale && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-2xl shadow-card-hover max-w-4xl w-full max-h-[90vh] overflow-y-auto">
            {/* Header del modal */}
            <div className="sticky top-0 bg-white border-b border-gray-200 px-6 py-4 rounded-t-2xl flex items-center justify-between">
              <div className="flex items-center gap-4">
                <h2 className="text-2xl font-bold text-[#212121]">
                  {selectedSale.sale_number}
                </h2>
                <span className={getStatusBadge(selectedSale.status)}>
                  {getStatusLabel(selectedSale.status)}
                </span>
                <span className="badge bg-blue-100 text-blue-800">
                  {getPaymentTypeInfo(selectedSale.payment_type).icon} {getPaymentTypeInfo(selectedSale.payment_type).label}
                </span>
              </div>
              <button
                onClick={handleCloseModal}
                className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
              >
                <X className="w-5 h-5 text-gray-600" />
              </button>
            </div>

            {/* Contenido del modal */}
            <div className="p-6 space-y-6">
              {/* Información del cliente y fechas */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="card">
                  <h3 className="text-sm font-semibold text-gray-600 mb-3">INFORMACIÓN DEL CLIENTE</h3>
                  <div className="space-y-2">
                    <div className="flex items-center gap-2">
                      <User className="w-4 h-4 text-gray-400" />
                      <p className="font-semibold text-[#212121]">{selectedSale.customer.name}</p>
                    </div>
                    {selectedSale.customer.tax_id && (
                      <p className="text-sm text-gray-600 ml-6">RNC: {selectedSale.customer.tax_id}</p>
                    )}
                  </div>
                </div>

                <div className="card">
                  <h3 className="text-sm font-semibold text-gray-600 mb-3">INFORMACIÓN DE LA VENTA</h3>
                  <div className="space-y-2">
                    <div className="flex items-center justify-between">
                      <span className="text-sm text-gray-600">Fecha:</span>
                      <span className="font-medium">{new Date(selectedSale.created_at).toLocaleDateString('es-DO')}</span>
                    </div>
                    <div className="flex items-center justify-between">
                      <span className="text-sm text-gray-600">Hora:</span>
                      <span className="font-medium">{new Date(selectedSale.created_at).toLocaleTimeString('es-DO')}</span>
                    </div>
                    <div className="flex items-center justify-between">
                      <span className="text-sm text-gray-600">Método de pago:</span>
                      <span className="font-medium">
                        {getPaymentTypeInfo(selectedSale.payment_type).icon} {getPaymentTypeInfo(selectedSale.payment_type).label}
                      </span>
                    </div>
                    {selectedSale.invoice_id && (
                      <div className="flex items-center justify-between">
                        <span className="text-sm text-gray-600">Factura:</span>
                        <span className="font-medium text-[#FF6B00]">#{selectedSale.invoice_id}</span>
                      </div>
                    )}
                  </div>
                </div>
              </div>

              {/* Items de la venta */}
              <div className="card">
                <h3 className="text-sm font-semibold text-gray-600 mb-4">ARTÍCULOS</h3>
                <div className="overflow-x-auto">
                  <table className="w-full">
                    <thead className="bg-[#F5F5F5]">
                      <tr>
                        <th className="text-left p-3 text-sm font-semibold text-gray-600">Descripción</th>
                        <th className="text-center p-3 text-sm font-semibold text-gray-600">Cantidad</th>
                        <th className="text-right p-3 text-sm font-semibold text-gray-600">Precio Unit.</th>
                        <th className="text-right p-3 text-sm font-semibold text-gray-600">ITBIS</th>
                        <th className="text-right p-3 text-sm font-semibold text-gray-600">Total</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-200">
                      {selectedSale.items.map((item, idx) => (
                        <tr key={idx} className="hover:bg-gray-50">
                          <td className="p-3 text-sm">{item.description}</td>
                          <td className="p-3 text-sm text-center">{item.quantity}</td>
                          <td className="p-3 text-sm text-right">${item.unit_price?.toFixed(2) || '0.00'}</td>
                          <td className="p-3 text-sm text-right">${item.tax_amount?.toFixed(2) || '0.00'}</td>
                          <td className="p-3 text-sm text-right font-medium">${item.line_total?.toFixed(2) || '0.00'}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>

              {/* Totales */}
              <div className="card bg-[#F5F5F5]">
                <div className="space-y-3">
                  <div className="flex justify-between items-center">
                    <span className="text-gray-600">Subtotal:</span>
                    <span className="text-lg font-semibold">${selectedSale.subtotal?.toFixed(2) || '0.00'}</span>
                  </div>
                  <div className="flex justify-between items-center">
                    <span className="text-gray-600">ITBIS (18%):</span>
                    <span className="text-lg font-semibold">${selectedSale.tax?.toFixed(2) || '0.00'}</span>
                  </div>
                  <div className="border-t border-gray-300 pt-3 flex justify-between items-center">
                    <span className="text-xl font-bold text-[#212121]">TOTAL:</span>
                    <span className="text-2xl font-bold text-[#FF6B00]">${selectedSale.total?.toFixed(2) || '0.00'}</span>
                  </div>
                </div>
              </div>

              {/* Notas */}
              {selectedSale.notes && (
                <div className="card">
                  <h3 className="text-sm font-semibold text-gray-600 mb-2">NOTAS</h3>
                  <p className="text-sm text-gray-700">{selectedSale.notes}</p>
                </div>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
