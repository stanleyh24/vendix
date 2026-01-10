import { useState, useEffect } from 'react';
import { FileText, Search, Eye, Send, XCircle, Download, Calendar, User, DollarSign, Edit2, X, Save, Plus } from 'lucide-react';
import Alert from '../components/Alert';
import CreateInvoiceModal from '../components/CreateInvoiceModal';
import api from '../lib/api';

export default function Invoices() {
  const [invoices, setInvoices] = useState([]);
  const [loading, setLoading] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const [filterStatus, setFilterStatus] = useState('all');
  const [alert, setAlert] = useState(null);
  const [selectedInvoice, setSelectedInvoice] = useState(null);
  const [showDetailModal, setShowDetailModal] = useState(false);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [isEditing, setIsEditing] = useState(false);
  const [editForm, setEditForm] = useState({ notes: '', terms: '' });

  // Estados de factura
  const statuses = [
    { value: 'draft', label: 'Borrador', color: 'bg-gray-500', badge: 'badge bg-gray-500 text-white' },
    { value: 'pending', label: 'Pendiente', color: 'bg-yellow-500', badge: 'badge-warning' },
    { value: 'paid', label: 'Pagada', color: 'bg-green-500', badge: 'badge-success' },
    { value: 'sent', label: 'Enviada', color: 'bg-blue-500', badge: 'badge-blue' },
    { value: 'cancelled', label: 'Anulada', color: 'bg-red-500', badge: 'badge-error' },
    { value: 'overdue', label: 'Vencida', color: 'bg-red-600', badge: 'badge bg-red-600 text-white' },
  ];

  // Cargar facturas del backend
  useEffect(() => {
    loadInvoices();
  }, []);

  const loadInvoices = async () => {
    setLoading(true);
    try {
      // Include customer information in the request
      const response = await api.get('/invoices?include=customer');
      const data = response.data || [];
      
      // Map backend data to frontend format
      const mappedInvoices = data.map(inv => {
        // Use customer object if available, otherwise fallback to customer_name
        let customerData = {
          id: inv.customer_id,
          name: inv.customer_name || 'Cliente',
          tax_id: ''
        };
        
        // If customer object is included, use it
        if (inv.customer && typeof inv.customer === 'object') {
          customerData = {
            id: inv.customer.id || inv.customer_id,
            name: inv.customer.name || inv.customer_name || 'Cliente',
            tax_id: inv.customer.tax_id || '',
            email: inv.customer.email,
            phone: inv.customer.phone,
            address: inv.customer.address,
            city: inv.customer.city,
            state: inv.customer.state,
            postal_code: inv.customer.postal_code,
            country: inv.customer.country
          };
        }
        
        return {
          id: inv.id,
          invoice_number: inv.invoice_number,
          ncf: inv.ncf,
          ncf_type: inv.ncf_type,
          customer: customerData,
          date: inv.issue_date?.split('T')[0] || inv.issue_date,
          due_date: inv.due_date?.split('T')[0] || inv.due_date,
          status: inv.status,
          subtotal: inv.subtotal,
          tax: inv.tax_amount,
          total: inv.total,
          paid_amount: inv.paid_amount || 0,
          withholding_tax_amount: inv.withholding_tax_amount,
          withholding_tax_type: inv.withholding_tax_type,
          withholding_rate: inv.withholding_rate,
          withholding_exempt: inv.withholding_exempt,
          net_amount: inv.net_amount,
          items: inv.lines || [],
          created_at: inv.created_at
        };
      });
      
      setInvoices(mappedInvoices);
    } catch (error) {
      console.error('Error loading invoices:', error);
      showAlert('error', 'Error', 'No se pudieron cargar las facturas');
      setInvoices([]);
    } finally {
      setLoading(false);
    }
  };

  const showAlert = (type, title, message) => {
    setAlert({ type, title, message });
    setTimeout(() => setAlert(null), 5000);
  };

  // Filtrar facturas
  const filteredInvoices = invoices.filter(invoice => {
    const matchesSearch = 
      invoice.invoice_number.toLowerCase().includes(searchTerm.toLowerCase()) ||
      invoice.customer.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      (invoice.customer.tax_id && invoice.customer.tax_id.toLowerCase().includes(searchTerm.toLowerCase()));
    
    const matchesStatus = filterStatus === 'all' || invoice.status === filterStatus;

    return matchesSearch && matchesStatus;
  });

  // Acciones de factura
  const handleSendInvoice = async (invoice) => {
    try {
      await api.post(`/invoices/${invoice.id}/send`);
      showAlert('success', 'Factura enviada', `Factura ${invoice.invoice_number} enviada al cliente`);
      loadInvoices();  // Reload to update status
    } catch (error) {
      showAlert('error', 'Error', error.response?.data?.error || 'No se pudo enviar la factura');
    }
  };

  const handleCancelInvoice = async (invoice) => {
    try {
      await api.post(`/invoices/${invoice.id}/cancel`);
      showAlert('warning', 'Factura anulada', `Factura ${invoice.invoice_number} ha sido anulada`);
      loadInvoices();  // Reload to update status
    } catch (error) {
      showAlert('error', 'Error', error.response?.data?.error || 'No se pudo anular la factura');
    }
  };

  const handleDownloadInvoice = async (invoice) => {
    try {
      // Usa el cliente axios ya configurado con Authorization y X-Tenant-ID
      const response = await api.get(`/invoices/${invoice.id}/pdf`, { responseType: 'blob' });
      // response.data is already a Blob when responseType is 'blob'
      const blob = response.data;
      const link = document.createElement('a');
      link.href = URL.createObjectURL(blob);
      link.download = `invoice-${invoice.invoice_number}.pdf`;
      document.body.appendChild(link);
      link.click();
      link.remove();
    } catch (err) {
      const message = err.response?.data?.error || 'No se pudo descargar el PDF';
      showAlert('error', 'Error', message);
    }
  };

  // Ver detalles de factura
  const handleViewDetails = async (invoice) => {
    try {
      // Include customer information in the request
      const response = await api.get(`/invoices/${invoice.id}?include=customer`);
      const invoiceData = response.data;
      
      // Use customer object if available, otherwise fallback to customer_name
      let customerData = {
        id: invoiceData.customer_id,
        name: invoiceData.customer_name || 'Cliente',
        tax_id: ''
      };
      
      // If customer object is included, use it
      if (invoiceData.customer && typeof invoiceData.customer === 'object') {
        customerData = {
          id: invoiceData.customer.id || invoiceData.customer_id,
          name: invoiceData.customer.name || invoiceData.customer_name || 'Cliente',
          tax_id: invoiceData.customer.tax_id || '',
          email: invoiceData.customer.email,
          phone: invoiceData.customer.phone,
          address: invoiceData.customer.address,
          city: invoiceData.customer.city,
          state: invoiceData.customer.state,
          postal_code: invoiceData.customer.postal_code,
          country: invoiceData.customer.country
        };
      }
      
      // Map backend data
      const mappedInvoice = {
        ...invoiceData,
        customer: customerData,
        date: invoiceData.issue_date?.split('T')[0] || invoiceData.issue_date,
        due_date: invoiceData.due_date?.split('T')[0] || invoiceData.due_date,
        tax: invoiceData.tax_amount,
        items: invoiceData.lines || []
      };
      
      setSelectedInvoice(mappedInvoice);
      setEditForm({
        notes: invoiceData.notes || '',
        terms: invoiceData.terms || ''
      });
      setIsEditing(false);
      setShowDetailModal(true);
    } catch (error) {
      showAlert('error', 'Error', 'No se pudieron cargar los detalles de la factura');
    }
  };

  // Guardar cambios
  const handleSaveChanges = async () => {
    try {
      await api.put(`/invoices/${selectedInvoice.id}`, {
        notes: editForm.notes || null,
        terms: editForm.terms || null
      });
      
      showAlert('success', 'Cambios guardados', 'La factura se actualizó correctamente');
      setIsEditing(false);
      loadInvoices();
      
      // Recargar detalles
      handleViewDetails(selectedInvoice);
    } catch (error) {
      showAlert('error', 'Error', error.response?.data?.error || 'No se pudieron guardar los cambios');
    }
  };

  // Cerrar modal
  const handleCloseModal = () => {
    setShowDetailModal(false);
    setSelectedInvoice(null);
    setIsEditing(false);
  };

  // Calcular totales
  const calculateTotals = () => {
    const total = filteredInvoices.reduce((sum, inv) => sum + inv.total, 0);
    const paid = filteredInvoices.filter(inv => inv.status === 'paid').reduce((sum, inv) => sum + inv.total, 0);
    const pending = filteredInvoices.filter(inv => inv.status === 'pending').reduce((sum, inv) => sum + inv.total, 0);
    
    return { total, paid, pending };
  };

  const totals = calculateTotals();

  const getStatusBadge = (status) => {
    return statuses.find(s => s.value === status)?.badge || 'badge';
  };

  const getStatusLabel = (status) => {
    return statuses.find(s => s.value === status)?.label || status;
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-[#212121] flex items-center gap-3">
            <FileText className="w-8 h-8 text-[#FF6B00]" />
            Facturación
          </h1>
          <p className="text-gray-600 mt-1">Gestiona tus facturas electrónicas</p>
        </div>
        <button
          onClick={() => setShowCreateModal(true)}
          className="btn-primary flex items-center gap-2"
        >
          <Plus className="w-5 h-5" />
          Nueva Factura
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
              <p className="text-sm text-gray-600">Total Facturado</p>
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
              <p className="text-sm text-gray-600">Pagadas</p>
              <p className="text-3xl font-bold text-[#00C853] mt-1">
                ${totals.paid.toFixed(2)}
              </p>
            </div>
            <div className="w-12 h-12 bg-green-100 rounded-full flex items-center justify-center">
              <FileText className="w-6 h-6 text-green-600" />
            </div>
          </div>
        </div>

        <div className="card bg-gradient-to-br from-yellow-50 to-white">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600">Pendientes</p>
              <p className="text-3xl font-bold text-[#FFA000] mt-1">
                ${totals.pending.toFixed(2)}
              </p>
            </div>
            <div className="w-12 h-12 bg-yellow-100 rounded-full flex items-center justify-center">
              <Calendar className="w-6 h-6 text-yellow-600" />
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
                placeholder="Buscar por número, cliente o RNC..."
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
              Todas ({invoices.length})
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
                {status.label} ({invoices.filter(i => i.status === status.value).length})
              </button>
            ))}
          </div>
        </div>
      </div>

      {/* Lista de facturas */}
      {loading ? (
        <div className="card text-center py-12">
          <div className="inline-block w-8 h-8 border-4 border-[#FF6B00] border-t-transparent rounded-full animate-spin"></div>
          <p className="mt-4 text-gray-600">Cargando facturas...</p>
        </div>
      ) : filteredInvoices.length === 0 ? (
        <div className="card text-center py-12">
          <FileText className="w-16 h-16 text-gray-300 mx-auto mb-4" />
          <h3 className="text-lg font-semibold text-gray-600 mb-2">
            {searchTerm ? 'No se encontraron facturas' : 'No hay facturas'}
          </h3>
          <p className="text-gray-500 mb-6">
            {searchTerm ? 'Intenta con otro término de búsqueda' : 'Las facturas aparecerán cuando proceses ventas'}
          </p>
        </div>
      ) : (
        <div className="card p-0 overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Factura</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Cliente</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Emisión</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Vence</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Estado</th>
                <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Subtotal</th>
                <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">ITBIS</th>
                <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Total</th>
                <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Retención</th>
                <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Monto Neto</th>
                <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Acciones</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {filteredInvoices.map((invoice) => (
                <tr key={invoice.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-mono font-semibold text-[#FF6B00]">{invoice.invoice_number}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-[#212121]">{invoice.customer.name}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-700">{new Date(invoice.date).toLocaleDateString('es-DO')}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-700">{new Date(invoice.due_date).toLocaleDateString('es-DO')}</td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <span className={getStatusBadge(invoice.status)}>{getStatusLabel(invoice.status)}</span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">{invoice.subtotal.toFixed(2)}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">{invoice.tax.toFixed(2)}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-right font-semibold text-[#FF6B00]">{invoice.total.toFixed(2)}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-right">
                    {invoice.withholding_tax_amount && invoice.withholding_tax_amount > 0 ? (
                      <span className="text-red-600 font-medium">-{invoice.withholding_tax_amount.toFixed(2)}</span>
                    ) : (
                      <span className="text-gray-400">-</span>
                    )}
                    {invoice.ncf_type === '15' && (
                      <span className="ml-1 text-xs text-blue-600" title="Factura Gubernamental">🏛️</span>
                    )}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-right font-semibold">
                    {invoice.net_amount !== undefined && invoice.net_amount !== invoice.total ? (
                      <span className="text-green-600">{invoice.net_amount.toFixed(2)}</span>
                    ) : (
                      <span className="text-gray-400">-</span>
                    )}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-right text-sm">
                    <button
                      onClick={() => handleViewDetails(invoice)}
                      className="p-2 hover:bg-blue-50 rounded-lg mr-1"
                      title="Ver detalles"
                    >
                      <Eye className="w-4 h-4 text-blue-600" />
                    </button>
                    <button
                      onClick={() => handleDownloadInvoice(invoice)}
                      className="p-2 hover:bg-gray-100 rounded-lg mr-1"
                      title="Descargar PDF"
                    >
                      <Download className="w-4 h-4 text-gray-600" />
                    </button>
                    {invoice.status !== 'sent' && invoice.status !== 'cancelled' && invoice.status !== 'paid' && (
                      <button
                        onClick={() => handleSendInvoice(invoice)}
                        className="p-2 hover:bg-blue-50 rounded-lg mr-1"
                        title="Enviar a DGII"
                      >
                        <Send className="w-4 h-4 text-blue-600" />
                      </button>
                    )}
                    {invoice.status !== 'cancelled' && invoice.status !== 'paid' && (
                      <button
                        onClick={() => handleCancelInvoice(invoice)}
                        className="p-2 hover:bg-red-50 rounded-lg"
                        title="Anular factura"
                      >
                        <XCircle className="w-4 h-4 text-red-600" />
                      </button>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
            <tfoot className="bg-gray-50">
              <tr>
                <td className="px-6 py-3 text-sm font-medium text-gray-700" colSpan={9}>Total listado</td>
                <td className="px-6 py-3 text-sm font-bold text-right text-[#FF6B00]">{totals.total.toFixed(2)}</td>
                <td className="px-6 py-3"></td>
              </tr>
            </tfoot>
          </table>
        </div>
      )}

      {/* Modal de detalles */}
      {showDetailModal && selectedInvoice && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-2xl shadow-card-hover max-w-4xl w-full max-h-[90vh] overflow-y-auto">
            {/* Header del modal */}
            <div className="sticky top-0 bg-white border-b border-gray-200 px-6 py-4 rounded-t-2xl flex items-center justify-between">
              <div className="flex items-center gap-4">
                <h2 className="text-2xl font-bold text-[#212121]">
                  {selectedInvoice.invoice_number}
                </h2>
                <span className={getStatusBadge(selectedInvoice.status)}>
                  {getStatusLabel(selectedInvoice.status)}
                </span>
              </div>
              <div className="flex items-center gap-2">
                {selectedInvoice.status !== 'paid' && selectedInvoice.status !== 'cancelled' && (
                  <>
                    {isEditing ? (
                      <>
                        <button
                          onClick={handleSaveChanges}
                          className="btn-primary flex items-center gap-2 text-sm"
                        >
                          <Save className="w-4 h-4" />
                          Guardar
                        </button>
                        <button
                          onClick={() => {
                            setIsEditing(false);
                            setEditForm({
                              notes: selectedInvoice.notes || '',
                              terms: selectedInvoice.terms || ''
                            });
                          }}
                          className="btn-outline text-sm"
                        >
                          Cancelar
                        </button>
                      </>
                    ) : (
                      <button
                        onClick={() => setIsEditing(true)}
                        className="btn-primary flex items-center gap-2 text-sm"
                      >
                        <Edit2 className="w-4 h-4" />
                        Editar
                      </button>
                    )}
                  </>
                )}
                <button
                  onClick={handleCloseModal}
                  className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
                >
                  <X className="w-5 h-5 text-gray-600" />
                </button>
              </div>
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
                      <p className="font-semibold text-[#212121]">{selectedInvoice.customer.name}</p>
                    </div>
                    {selectedInvoice.customer.tax_id && (
                      <p className="text-sm text-gray-600 ml-6">RNC: {selectedInvoice.customer.tax_id}</p>
                    )}
                  </div>
                </div>

                <div className="card">
                  <h3 className="text-sm font-semibold text-gray-600 mb-3">FECHAS</h3>
                  <div className="space-y-2">
                    <div className="flex items-center justify-between">
                      <span className="text-sm text-gray-600">Fecha de emisión:</span>
                      <span className="font-medium">{new Date(selectedInvoice.date).toLocaleDateString('es-DO')}</span>
                    </div>
                    <div className="flex items-center justify-between">
                      <span className="text-sm text-gray-600">Fecha de vencimiento:</span>
                      <span className="font-medium">{new Date(selectedInvoice.due_date).toLocaleDateString('es-DO')}</span>
                    </div>
                    {selectedInvoice.sent_at && (
                      <div className="flex items-center justify-between">
                        <span className="text-sm text-gray-600">Enviada:</span>
                        <span className="font-medium">{new Date(selectedInvoice.sent_at).toLocaleDateString('es-DO')}</span>
                      </div>
                    )}
                  </div>
                </div>
              </div>

              {/* Items de la factura */}
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
                      {selectedInvoice.items.map((item, idx) => (
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
                    <span className="text-lg font-semibold">${selectedInvoice.subtotal?.toFixed(2) || '0.00'}</span>
                  </div>
                  <div className="flex justify-between items-center">
                    <span className="text-gray-600">ITBIS (18%):</span>
                    <span className="text-lg font-semibold">${selectedInvoice.tax?.toFixed(2) || '0.00'}</span>
                  </div>
                  <div className="border-t border-gray-300 pt-3 flex justify-between items-center">
                    <span className="text-xl font-bold text-[#212121]">TOTAL:</span>
                    <span className="text-xl font-bold text-[#FF6B00]">${selectedInvoice.total?.toFixed(2) || '0.00'}</span>
                  </div>
                  {/* Mostrar retención si existe */}
                  {selectedInvoice.withholding_tax_amount && selectedInvoice.withholding_tax_amount > 0 && (
                    <>
                      <div className="flex justify-between items-center pt-2 border-t border-gray-300">
                        <div className="flex items-center gap-2">
                          <span className="text-gray-600">Retención ISR</span>
                          {selectedInvoice.withholding_rate && (
                            <span className="text-xs text-gray-500">
                              ({(selectedInvoice.withholding_rate * 100).toFixed(2)}%)
                            </span>
                          )}
                        </div>
                        <span className="text-lg font-semibold text-red-600">
                          -${selectedInvoice.withholding_tax_amount.toFixed(2)}
                        </span>
                      </div>
                      <div className="border-t-2 border-gray-400 pt-3 flex justify-between items-center">
                        <span className="text-xl font-bold text-[#212121]">Monto Neto a Recibir:</span>
                        <span className="text-2xl font-bold text-green-600">
                          ${(selectedInvoice.net_amount || selectedInvoice.total - selectedInvoice.withholding_tax_amount).toFixed(2)}
                        </span>
                      </div>
                    </>
                  )}
                  {selectedInvoice.ncf_type === '15' && (
                    <div className="mt-2 p-2 bg-blue-50 border border-blue-200 rounded text-xs text-blue-700">
                      🏛️ Factura Gubernamental (NCF tipo 15)
                    </div>
                  )}
                  {selectedInvoice.paid_amount > 0 && (
                    <div className="flex justify-between items-center text-green-600">
                      <span className="font-medium">Pagado:</span>
                      <span className="text-lg font-semibold">${selectedInvoice.paid_amount?.toFixed(2) || '0.00'}</span>
                    </div>
                  )}
                </div>
              </div>

              {/* Notas y términos */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <label className="block text-sm font-semibold text-gray-600 mb-2">
                    Notas
                    {selectedInvoice.status === 'paid' && (
                      <span className="ml-2 text-xs text-gray-500">(solo lectura)</span>
                    )}
                  </label>
                  {isEditing ? (
                    <textarea
                      value={editForm.notes}
                      onChange={(e) => setEditForm({ ...editForm, notes: e.target.value })}
                      className="input-field"
                      rows="4"
                      placeholder="Notas adicionales..."
                    />
                  ) : (
                    <div className="card min-h-[100px] text-sm text-gray-700">
                      {selectedInvoice.notes || 'Sin notas'}
                    </div>
                  )}
                </div>

                <div>
                  <label className="block text-sm font-semibold text-gray-600 mb-2">
                    Términos y Condiciones
                    {selectedInvoice.status === 'paid' && (
                      <span className="ml-2 text-xs text-gray-500">(solo lectura)</span>
                    )}
                  </label>
                  {isEditing ? (
                    <textarea
                      value={editForm.terms}
                      onChange={(e) => setEditForm({ ...editForm, terms: e.target.value })}
                      className="input-field"
                      rows="4"
                      placeholder="Términos y condiciones..."
                    />
                  ) : (
                    <div className="card min-h-[100px] text-sm text-gray-700">
                      {selectedInvoice.terms || 'Sin términos especificados'}
                    </div>
                  )}
                </div>
              </div>

              {/* Acciones adicionales */}
              {selectedInvoice.status !== 'cancelled' && (
                <div className="flex gap-3 pt-4 border-t border-gray-200">
                  <button
                    onClick={() => handleDownloadInvoice(selectedInvoice)}
                    className="btn-outline flex-1 flex items-center justify-center gap-2"
                  >
                    <Download className="w-4 h-4" />
                    Descargar PDF
                  </button>
                  
                  {selectedInvoice.status !== 'sent' && selectedInvoice.status !== 'paid' && (
                    <button
                      onClick={() => {
                        handleSendInvoice(selectedInvoice);
                        handleCloseModal();
                      }}
                      className="btn-primary flex-1 flex items-center justify-center gap-2"
                    >
                      <Send className="w-4 h-4" />
                      Enviar a DGII
                    </button>
                  )}
                  
                  {selectedInvoice.status !== 'paid' && (
                    <button
                      onClick={() => {
                        if (confirm('¿Estás seguro de que deseas anular esta factura?')) {
                          handleCancelInvoice(selectedInvoice);
                          handleCloseModal();
                        }
                      }}
                      className="btn-error flex-1 flex items-center justify-center gap-2"
                    >
                      <XCircle className="w-4 h-4" />
                      Anular Factura
                    </button>
                  )}
                </div>
              )}
            </div>
          </div>
        </div>
      )}

      {/* Modal de crear factura */}
      <CreateInvoiceModal
        isOpen={showCreateModal}
        onClose={() => setShowCreateModal(false)}
        onSuccess={(invoice) => {
          showAlert('success', 'Factura creada', `Factura ${invoice.invoice_number} creada exitosamente`);
          loadInvoices();
        }}
      />
    </div>
  );
}
