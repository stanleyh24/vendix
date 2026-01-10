import { useState, useEffect } from 'react';
import { FileText, Search, Eye, Send, XCircle, Download, Calendar, User, DollarSign, X, Receipt } from 'lucide-react';
import Alert from '../components/Alert';
import api from '../lib/api';

export default function CreditNotes() {
  const [creditNotes, setCreditNotes] = useState([]);
  const [loading, setLoading] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const [filterStatus, setFilterStatus] = useState('all');
  const [alert, setAlert] = useState(null);
  const [selectedCreditNote, setSelectedCreditNote] = useState(null);
  const [showDetailModal, setShowDetailModal] = useState(false);

  // Estados de nota de crédito
  const statuses = [
    { value: 'draft', label: 'Borrador', color: 'bg-gray-500', badge: 'badge bg-gray-500 text-white' },
    { value: 'pending', label: 'Pendiente', color: 'bg-yellow-500', badge: 'badge-warning' },
    { value: 'sent', label: 'Enviada a DGII', color: 'bg-blue-500', badge: 'badge-blue' },
    { value: 'cancelled', label: 'Anulada', color: 'bg-red-500', badge: 'badge-error' },
  ];

  // Cargar notas de crédito del backend
  useEffect(() => {
    loadCreditNotes();
  }, []);

  const loadCreditNotes = async () => {
    setLoading(true);
    try {
      const response = await api.get('/credit-notes');
      const data = response.data || [];
      
      // Map backend data to frontend format
      const mappedCreditNotes = data.map(cn => ({
        id: cn.id,
        credit_note_number: cn.credit_note_number,
        ncf: cn.ncf,
        original_invoice_id: cn.original_invoice_id,
        original_invoice_number: cn.original_invoice_number || 'N/A',
        return_id: cn.return_id,
        customer_id: cn.customer_id,
        customer_name: cn.customer_name || 'Cliente',
        date: cn.issue_date?.split('T')[0] || cn.issue_date,
        reason: cn.reason || 'Devolución de productos',
        status: cn.status,
        subtotal: cn.subtotal,
        tax: cn.tax_amount,
        total: cn.total,
        currency: cn.currency || 'DOP',
        notes: cn.notes,
        dgii_status: cn.dgii_status,
        tracking_code: cn.tracking_code,
        signed_at: cn.signed_at,
        sent_at: cn.sent_at,
        items: cn.lines || [],
        created_at: cn.created_at
      }));
      
      setCreditNotes(mappedCreditNotes);
    } catch (error) {
      console.error('Error loading credit notes:', error);
      showAlert('error', 'Error', 'No se pudieron cargar las notas de crédito');
      setCreditNotes([]);
    } finally {
      setLoading(false);
    }
  };

  const showAlert = (type, title, message) => {
    setAlert({ type, title, message });
    setTimeout(() => setAlert(null), 5000);
  };

  // Filtrar notas de crédito
  const filteredCreditNotes = creditNotes.filter(creditNote => {
    const matchesSearch = 
      creditNote.credit_note_number.toLowerCase().includes(searchTerm.toLowerCase()) ||
      creditNote.ncf?.toLowerCase().includes(searchTerm.toLowerCase()) ||
      creditNote.original_invoice_number.toLowerCase().includes(searchTerm.toLowerCase()) ||
      creditNote.customer_name.toLowerCase().includes(searchTerm.toLowerCase());
    
    const matchesStatus = filterStatus === 'all' || creditNote.status === filterStatus;

    return matchesSearch && matchesStatus;
  });

  // Acciones de nota de crédito
  const handleSendCreditNote = async (creditNote) => {
    try {
      await api.post(`/credit-notes/${creditNote.id}/send`);
      showAlert('success', 'Nota de Crédito enviada', `Nota de Crédito ${creditNote.credit_note_number} enviada a DGII`);
      loadCreditNotes();
    } catch (error) {
      showAlert('error', 'Error', error.response?.data?.error || 'No se pudo enviar la nota de crédito');
    }
  };

  const handleCancelCreditNote = async (creditNote) => {
    try {
      await api.post(`/credit-notes/${creditNote.id}/cancel`);
      showAlert('warning', 'Nota de Crédito anulada', `Nota de Crédito ${creditNote.credit_note_number} ha sido anulada`);
      loadCreditNotes();
    } catch (error) {
      showAlert('error', 'Error', error.response?.data?.error || 'No se pudo anular la nota de crédito');
    }
  };

  const handleDownloadCreditNote = async (creditNote) => {
    try {
      const response = await api.get(`/credit-notes/${creditNote.id}/pdf`, { responseType: 'blob' });
      const blob = response.data;
      const link = document.createElement('a');
      link.href = URL.createObjectURL(blob);
      link.download = `credit-note-${creditNote.credit_note_number}.pdf`;
      document.body.appendChild(link);
      link.click();
      link.remove();
    } catch (err) {
      const message = err.response?.data?.error || 'No se pudo descargar el PDF';
      showAlert('error', 'Error', message);
    }
  };

  // Ver detalles de nota de crédito
  const handleViewDetails = async (creditNote) => {
    try {
      const response = await api.get(`/credit-notes/${creditNote.id}`);
      const creditNoteData = response.data;
      
      const mappedCreditNote = {
        ...creditNoteData,
        customer_name: creditNoteData.customer_name || 'Cliente',
        date: creditNoteData.issue_date?.split('T')[0] || creditNoteData.issue_date,
        tax: creditNoteData.tax_amount,
        items: creditNoteData.lines || []
      };
      
      setSelectedCreditNote(mappedCreditNote);
      setShowDetailModal(true);
    } catch (error) {
      showAlert('error', 'Error', 'No se pudieron cargar los detalles de la nota de crédito');
    }
  };

  // Cerrar modal
  const handleCloseModal = () => {
    setShowDetailModal(false);
    setSelectedCreditNote(null);
  };

  // Calcular totales
  const calculateTotals = () => {
    const total = filteredCreditNotes.reduce((sum, cn) => sum + cn.total, 0);
    const sent = filteredCreditNotes.filter(cn => cn.status === 'sent').reduce((sum, cn) => sum + cn.total, 0);
    const pending = filteredCreditNotes.filter(cn => cn.status === 'pending' || cn.status === 'draft').reduce((sum, cn) => sum + cn.total, 0);
    
    return { total, sent, pending };
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
            <Receipt className="w-8 h-8 text-[#FF6B00]" />
            Notas de Crédito
          </h1>
          <p className="text-gray-600 mt-1">Gestiona tus notas de crédito fiscales (NCF tipo 04)</p>
        </div>
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
              <p className="text-sm text-gray-600">Total Notas de Crédito</p>
              <p className="text-3xl font-bold text-[#FF6B00] mt-1">
                ${totals.total.toFixed(2)}
              </p>
            </div>
            <div className="w-12 h-12 bg-orange-100 rounded-full flex items-center justify-center">
              <DollarSign className="w-6 h-6 text-[#FF6B00]" />
            </div>
          </div>
        </div>

        <div className="card bg-gradient-to-br from-blue-50 to-white">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600">Enviadas a DGII</p>
              <p className="text-3xl font-bold text-[#2196F3] mt-1">
                ${totals.sent.toFixed(2)}
              </p>
            </div>
            <div className="w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center">
              <Send className="w-6 h-6 text-blue-600" />
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
                placeholder="Buscar por número, NCF, factura original o cliente..."
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
              Todas ({creditNotes.length})
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
                {status.label} ({creditNotes.filter(cn => cn.status === status.value).length})
              </button>
            ))}
          </div>
        </div>
      </div>

      {/* Lista de notas de crédito */}
      {loading ? (
        <div className="card text-center py-12">
          <div className="inline-block w-8 h-8 border-4 border-[#FF6B00] border-t-transparent rounded-full animate-spin"></div>
          <p className="mt-4 text-gray-600">Cargando notas de crédito...</p>
        </div>
      ) : filteredCreditNotes.length === 0 ? (
        <div className="card text-center py-12">
          <Receipt className="w-16 h-16 text-gray-300 mx-auto mb-4" />
          <h3 className="text-lg font-semibold text-gray-600 mb-2">
            {searchTerm ? 'No se encontraron notas de crédito' : 'No hay notas de crédito'}
          </h3>
          <p className="text-gray-500 mb-6">
            {searchTerm ? 'Intenta con otro término de búsqueda' : 'Las notas de crédito se generan automáticamente desde devoluciones o manualmente desde facturas'}
          </p>
        </div>
      ) : (
        <div className="card p-0 overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Nota de Crédito</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">NCF</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Factura Original</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Cliente</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Fecha</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Motivo</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Estado</th>
                <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Subtotal</th>
                <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">ITBIS</th>
                <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Total</th>
                <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Acciones</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {filteredCreditNotes.map((creditNote) => (
                <tr key={creditNote.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-mono font-semibold text-[#FF6B00]">{creditNote.credit_note_number}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-mono text-gray-700">{creditNote.ncf || 'N/A'}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-700">{creditNote.original_invoice_number}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-[#212121]">{creditNote.customer_name}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-700">{new Date(creditNote.date).toLocaleDateString('es-DO')}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-700 max-w-xs truncate">{creditNote.reason}</td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <span className={getStatusBadge(creditNote.status)}>{getStatusLabel(creditNote.status)}</span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">{creditNote.subtotal.toFixed(2)}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">{creditNote.tax.toFixed(2)}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-right font-semibold text-[#FF6B00]">{creditNote.total.toFixed(2)}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-right text-sm">
                    <button
                      onClick={() => handleViewDetails(creditNote)}
                      className="p-2 hover:bg-blue-50 rounded-lg mr-1"
                      title="Ver detalles"
                    >
                      <Eye className="w-4 h-4 text-blue-600" />
                    </button>
                    <button
                      onClick={() => handleDownloadCreditNote(creditNote)}
                      className="p-2 hover:bg-gray-100 rounded-lg mr-1"
                      title="Descargar PDF"
                    >
                      <Download className="w-4 h-4 text-gray-600" />
                    </button>
                    {creditNote.status !== 'sent' && creditNote.status !== 'cancelled' && (
                      <button
                        onClick={() => handleSendCreditNote(creditNote)}
                        className="p-2 hover:bg-blue-50 rounded-lg mr-1"
                        title="Enviar a DGII"
                      >
                        <Send className="w-4 h-4 text-blue-600" />
                      </button>
                    )}
                    {creditNote.status !== 'cancelled' && (
                      <button
                        onClick={() => handleCancelCreditNote(creditNote)}
                        className="p-2 hover:bg-red-50 rounded-lg"
                        title="Anular nota de crédito"
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
      {showDetailModal && selectedCreditNote && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-2xl shadow-card-hover max-w-4xl w-full max-h-[90vh] overflow-y-auto">
            {/* Header del modal */}
            <div className="sticky top-0 bg-white border-b border-gray-200 px-6 py-4 rounded-t-2xl flex items-center justify-between">
              <div className="flex items-center gap-4">
                <h2 className="text-2xl font-bold text-[#212121]">
                  {selectedCreditNote.credit_note_number}
                </h2>
                <span className={getStatusBadge(selectedCreditNote.status)}>
                  {getStatusLabel(selectedCreditNote.status)}
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
              {/* Información general */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="card">
                  <h3 className="text-sm font-semibold text-gray-600 mb-3">INFORMACIÓN GENERAL</h3>
                  <div className="space-y-2">
                    <div className="flex items-center justify-between">
                      <span className="text-sm text-gray-600">NCF:</span>
                      <span className="font-mono font-semibold text-[#212121]">{selectedCreditNote.ncf || 'N/A'}</span>
                    </div>
                    <div className="flex items-center justify-between">
                      <span className="text-sm text-gray-600">Factura Original:</span>
                      <span className="font-medium">{selectedCreditNote.original_invoice_number}</span>
                    </div>
                    {selectedCreditNote.return_id && (
                      <div className="flex items-center justify-between">
                        <span className="text-sm text-gray-600">Devolución:</span>
                        <span className="font-medium">#{selectedCreditNote.return_id.substring(0, 8)}...</span>
                      </div>
                    )}
                  </div>
                </div>

                <div className="card">
                  <h3 className="text-sm font-semibold text-gray-600 mb-3">FECHA Y CLIENTE</h3>
                  <div className="space-y-2">
                    <div className="flex items-center justify-between">
                      <span className="text-sm text-gray-600">Fecha de emisión:</span>
                      <span className="font-medium">{new Date(selectedCreditNote.date).toLocaleDateString('es-DO')}</span>
                    </div>
                    <div className="flex items-center gap-2">
                      <User className="w-4 h-4 text-gray-400" />
                      <p className="font-semibold text-[#212121]">{selectedCreditNote.customer_name}</p>
                    </div>
                    {selectedCreditNote.sent_at && (
                      <div className="flex items-center justify-between">
                        <span className="text-sm text-gray-600">Enviada:</span>
                        <span className="font-medium">{new Date(selectedCreditNote.sent_at).toLocaleDateString('es-DO')}</span>
                      </div>
                    )}
                  </div>
                </div>
              </div>

              {/* Motivo */}
              {selectedCreditNote.reason && (
                <div className="card">
                  <h3 className="text-sm font-semibold text-gray-600 mb-2">MOTIVO</h3>
                  <p className="text-sm text-gray-700">{selectedCreditNote.reason}</p>
                </div>
              )}

              {/* Items de la nota de crédito */}
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
                      {selectedCreditNote.items.map((item, idx) => (
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
                    <span className="text-lg font-semibold">${selectedCreditNote.subtotal?.toFixed(2) || '0.00'}</span>
                  </div>
                  <div className="flex justify-between items-center">
                    <span className="text-gray-600">ITBIS (18%):</span>
                    <span className="text-lg font-semibold">${selectedCreditNote.tax?.toFixed(2) || '0.00'}</span>
                  </div>
                  <div className="border-t border-gray-300 pt-3 flex justify-between items-center">
                    <span className="text-xl font-bold text-[#212121]">TOTAL:</span>
                    <span className="text-2xl font-bold text-[#FF6B00]">${selectedCreditNote.total?.toFixed(2) || '0.00'}</span>
                  </div>
                </div>
              </div>

              {/* Estado DGII */}
              {selectedCreditNote.dgii_status && (
                <div className="card bg-blue-50 border border-blue-200">
                  <h3 className="text-sm font-semibold text-blue-900 mb-2">ESTADO DGII</h3>
                  <div className="space-y-2">
                    <div className="flex items-center justify-between">
                      <span className="text-sm text-blue-700">Estado:</span>
                      <span className="font-medium text-blue-900">{selectedCreditNote.dgii_status}</span>
                    </div>
                    {selectedCreditNote.tracking_code && (
                      <div className="flex items-center justify-between">
                        <span className="text-sm text-blue-700">Tracking Code:</span>
                        <span className="font-mono text-blue-900">{selectedCreditNote.tracking_code}</span>
                      </div>
                    )}
                  </div>
                </div>
              )}

              {/* Notas */}
              {selectedCreditNote.notes && (
                <div className="card">
                  <h3 className="text-sm font-semibold text-gray-600 mb-2">NOTAS</h3>
                  <p className="text-sm text-gray-700">{selectedCreditNote.notes}</p>
                </div>
              )}

              {/* Acciones */}
              {selectedCreditNote.status !== 'cancelled' && (
                <div className="flex gap-3 pt-4 border-t border-gray-200">
                  <button
                    onClick={() => {
                      handleDownloadCreditNote(selectedCreditNote);
                      handleCloseModal();
                    }}
                    className="btn-outline flex-1 flex items-center justify-center gap-2"
                  >
                    <Download className="w-4 h-4" />
                    Descargar PDF
                  </button>
                  
                  {selectedCreditNote.status !== 'sent' && (
                    <button
                      onClick={() => {
                        handleSendCreditNote(selectedCreditNote);
                        handleCloseModal();
                      }}
                      className="btn-primary flex-1 flex items-center justify-center gap-2"
                    >
                      <Send className="w-4 h-4" />
                      Enviar a DGII
                    </button>
                  )}
                  
                  <button
                    onClick={() => {
                      if (confirm('¿Estás seguro de que deseas anular esta nota de crédito?')) {
                        handleCancelCreditNote(selectedCreditNote);
                        handleCloseModal();
                      }
                    }}
                    className="btn-error flex-1 flex items-center justify-center gap-2"
                  >
                    <XCircle className="w-4 h-4" />
                    Anular Nota de Crédito
                  </button>
                </div>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

