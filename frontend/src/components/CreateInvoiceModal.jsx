import { useState, useEffect } from 'react';
import { X, Save, Plus, Trash2, User, Calendar, FileText, DollarSign } from 'lucide-react';
import api from '../lib/api';
import Alert from './Alert';

export default function CreateInvoiceModal({ isOpen, onClose, onSuccess }) {
  const [formData, setFormData] = useState({
    customer_id: '',
    ncf_type: '02',
    issue_date: new Date().toISOString().split('T')[0],
    due_date: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
    status: 'draft',
    notes: '',
    terms: '',
    lines: [{ description: '', quantity: 1, unit_price: 0, tax_rate: 0.18 }],
    withholding_tax_type: null,
    withholding_rate: null,
    withholding_exempt: false,
  });
  
  const [customers, setCustomers] = useState([]);
  const [loading, setLoading] = useState(false);
  const [alert, setAlert] = useState(null);

  // Estados de factura disponibles
  const statusOptions = [
    { value: 'draft', label: 'Borrador', description: 'Factura en proceso de creación' },
    { value: 'pending', label: 'Pendiente', description: 'Factura creada, esperando pago' },
    { value: 'paid', label: 'Pagada', description: 'Factura completamente pagada' },
    { value: 'sent', label: 'Enviada', description: 'Factura enviada a DGII' },
    { value: 'cancelled', label: 'Cancelada', description: 'Factura cancelada' },
    { value: 'overdue', label: 'Vencida', description: 'Factura vencida sin pago' }
  ];

  // Tipos de NCF
  const ncfTypes = [
    { value: '01', label: 'Crédito Fiscal', description: 'Para empresas con RNC' },
    { value: '02', label: 'Consumidor Final', description: 'Para personas sin RNC' },
    { value: '15', label: 'Gubernamental', description: 'Para entidades del Estado (requiere cliente gubernamental)' }
  ];

  useEffect(() => {
    if (isOpen) {
      loadCustomers();
    }
  }, [isOpen]);

  const loadCustomers = async () => {
    try {
      const response = await api.get('/customers');
      setCustomers(response.data || []);
    } catch (error) {
      console.error('Error loading customers:', error);
    }
  };

  const showAlert = (type, title, message) => {
    setAlert({ type, title, message });
    setTimeout(() => setAlert(null), 5000);
  };

  const handleInputChange = (field, value) => {
    setFormData(prev => ({
      ...prev,
      [field]: value
    }));
  };

  const handleLineChange = (index, field, value) => {
    const newLines = [...formData.lines];
    newLines[index] = {
      ...newLines[index],
      [field]: value
    };
    setFormData(prev => ({
      ...prev,
      lines: newLines
    }));
  };

  const addLine = () => {
    setFormData(prev => ({
      ...prev,
      lines: [...prev.lines, { description: '', quantity: 1, unit_price: 0, tax_rate: 0.18 }]
    }));
  };

  const removeLine = (index) => {
    if (formData.lines.length > 1) {
      const newLines = formData.lines.filter((_, i) => i !== index);
      setFormData(prev => ({
        ...prev,
        lines: newLines
      }));
    }
  };

  const calculateLineTotal = (line) => {
    const subtotal = line.quantity * line.unit_price;
    const tax = subtotal * line.tax_rate;
    return subtotal + tax;
  };

  const calculateTotals = () => {
    const subtotal = formData.lines.reduce((sum, line) => 
      sum + (line.quantity * line.unit_price), 0
    );
    const tax = formData.lines.reduce((sum, line) => 
      sum + (line.quantity * line.unit_price * line.tax_rate), 0
    );
    const total = subtotal + tax;
    
    // Calcular retención si aplica
    const selectedCustomer = customers.find(c => c.id === formData.customer_id);
    const isGovInvoice = formData.ncf_type === '15' && selectedCustomer?.is_government_entity && !formData.withholding_exempt;
    
    let withholdingAmount = 0;
    let withholdingRate = 0.05; // 5% por defecto
    
    if (isGovInvoice) {
      if (formData.withholding_rate !== null && formData.withholding_rate !== undefined) {
        withholdingRate = parseFloat(formData.withholding_rate) || 0.05;
      } else if (selectedCustomer?.default_withholding_rate) {
        withholdingRate = selectedCustomer.default_withholding_rate;
      }
      withholdingAmount = total * withholdingRate;
    }
    
    const netAmount = total - withholdingAmount;
    
    return { subtotal, tax, total, withholdingAmount, withholdingRate, netAmount };
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    // Validaciones
    if (!formData.customer_id) {
      showAlert('error', 'Error', 'Debes seleccionar un cliente');
      return;
    }

    // Validar NCF tipo 15 requiere cliente gubernamental
    const selectedCustomer = customers.find(c => c.id === formData.customer_id);
    if (formData.ncf_type === '15') {
      if (!selectedCustomer) {
        showAlert('error', 'Error', 'Debes seleccionar un cliente');
        return;
      }
      if (!selectedCustomer.is_government_entity) {
        showAlert('error', 'Error', 'NCF tipo 15 (Gubernamental) solo se puede usar con clientes que sean entidades gubernamentales');
        return;
      }
    }

    if (formData.lines.some(line => !line.description.trim())) {
      showAlert('error', 'Error', 'Todas las líneas deben tener una descripción');
      return;
    }

    if (formData.lines.some(line => line.quantity <= 0)) {
      showAlert('error', 'Error', 'Todas las cantidades deben ser mayores a 0');
      return;
    }

    if (formData.lines.some(line => line.unit_price < 0)) {
      showAlert('error', 'Error', 'Todos los precios deben ser mayores o iguales a 0');
      return;
    }

    setLoading(true);
    try {
      // Preparar payload con retenciones
      const totals = calculateTotals();
      const payload = {
        ...formData,
        withholding_tax_type: totals.withholdingAmount > 0 ? 'isr' : null,
        withholding_rate: totals.withholdingAmount > 0 ? totals.withholdingRate : null,
        withholding_exempt: formData.withholding_exempt,
      };
      
      const response = await api.post('/invoices', payload);
      showAlert('success', 'Factura creada', `Factura ${response.data.invoice_number} creada exitosamente`);
      onSuccess && onSuccess(response.data);
      onClose();
    } catch (error) {
      showAlert('error', 'Error', error.response?.data?.error || 'Error al crear la factura');
    } finally {
      setLoading(false);
    }
  };

  const handleClose = () => {
    setFormData({
      customer_id: '',
      ncf_type: '02',
      issue_date: new Date().toISOString().split('T')[0],
      due_date: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
      status: 'draft',
      notes: '',
      terms: '',
      lines: [{ description: '', quantity: 1, unit_price: 0, tax_rate: 0.18 }],
      withholding_tax_type: null,
      withholding_rate: null,
      withholding_exempt: false,
    });
    setAlert(null);
    onClose();
  };

  if (!isOpen) return null;

  const totals = calculateTotals();

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-2xl shadow-card-hover max-w-4xl w-full max-h-[90vh] overflow-y-auto">
        {/* Header */}
        <div className="sticky top-0 bg-white border-b border-gray-200 px-6 py-4 rounded-t-2xl flex items-center justify-between">
          <h2 className="text-2xl font-bold text-[#212121] flex items-center gap-2">
            <FileText className="w-6 h-6 text-[#FF6B00]" />
            Crear Nueva Factura
          </h2>
          <button
            onClick={handleClose}
            className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Alert */}
        {alert && (
          <div className="px-6 pt-4">
            <Alert
              type={alert.type}
              title={alert.title}
              message={alert.message}
              onClose={() => setAlert(null)}
            />
          </div>
        )}

        {/* Form */}
        <form onSubmit={handleSubmit} className="p-6 space-y-6">
          {/* Información básica */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {/* Cliente */}
            <div>
              <label className="block text-sm font-semibold text-gray-600 mb-2">
                Cliente *
              </label>
              <select
                value={formData.customer_id}
                onChange={(e) => handleInputChange('customer_id', e.target.value)}
                className="input-field"
                required
              >
                <option value="">Seleccionar cliente</option>
                {customers.map(customer => (
                  <option key={customer.id} value={customer.id}>
                    {customer.name} {customer.tax_id && `(${customer.tax_id})`}
                  </option>
                ))}
              </select>
            </div>

            {/* Tipo de NCF */}
            <div>
              <label className="block text-sm font-semibold text-gray-600 mb-2">
                Tipo de NCF
              </label>
              <select
                value={formData.ncf_type}
                onChange={(e) => {
                  handleInputChange('ncf_type', e.target.value);
                  // Reset retenciones si cambia el tipo de NCF
                  if (e.target.value !== '15') {
                    handleInputChange('withholding_exempt', false);
                    handleInputChange('withholding_rate', null);
                  }
                }}
                className="input-field"
              >
                {ncfTypes.map(type => (
                  <option key={type.value} value={type.value}>
                    {type.label} - {type.description}
                  </option>
                ))}
              </select>
              {formData.ncf_type === '15' && (
                <p className="text-xs text-amber-600 mt-1">
                  Requiere cliente gubernamental
                </p>
              )}
            </div>

            {/* Fecha de emisión */}
            <div>
              <label className="block text-sm font-semibold text-gray-600 mb-2">
                Fecha de Emisión *
              </label>
              <input
                type="date"
                value={formData.issue_date}
                onChange={(e) => handleInputChange('issue_date', e.target.value)}
                className="input-field"
                required
              />
            </div>

            {/* Fecha de vencimiento */}
            <div>
              <label className="block text-sm font-semibold text-gray-600 mb-2">
                Fecha de Vencimiento *
              </label>
              <input
                type="date"
                value={formData.due_date}
                onChange={(e) => handleInputChange('due_date', e.target.value)}
                className="input-field"
                required
              />
            </div>

            {/* Status */}
            <div>
              <label className="block text-sm font-semibold text-gray-600 mb-2">
                Estado
              </label>
              <select
                value={formData.status}
                onChange={(e) => handleInputChange('status', e.target.value)}
                className="input-field"
              >
                {statusOptions.map(status => (
                  <option key={status.value} value={status.value}>
                    {status.label} - {status.description}
                  </option>
                ))}
              </select>
            </div>
          </div>

          {/* Líneas de factura */}
          <div>
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-semibold text-[#212121]">Artículos</h3>
              <button
                type="button"
                onClick={addLine}
                className="btn-primary flex items-center gap-2 text-sm"
              >
                <Plus className="w-4 h-4" />
                Agregar Línea
              </button>
            </div>

            <div className="space-y-4">
              {formData.lines.map((line, index) => (
                <div key={index} className="grid grid-cols-12 gap-4 items-end">
                  <div className="col-span-5">
                    <label className="block text-xs font-semibold text-gray-600 mb-1">
                      Descripción *
                    </label>
                    <input
                      type="text"
                      value={line.description}
                      onChange={(e) => handleLineChange(index, 'description', e.target.value)}
                      className="input-field"
                      placeholder="Descripción del artículo"
                      required
                    />
                  </div>
                  
                  <div className="col-span-2">
                    <label className="block text-xs font-semibold text-gray-600 mb-1">
                      Cantidad *
                    </label>
                    <input
                      type="number"
                      min="0.01"
                      step="0.01"
                      value={line.quantity}
                      onChange={(e) => handleLineChange(index, 'quantity', parseFloat(e.target.value) || 0)}
                      className="input-field"
                      required
                    />
                  </div>
                  
                  <div className="col-span-2">
                    <label className="block text-xs font-semibold text-gray-600 mb-1">
                      Precio Unit. *
                    </label>
                    <input
                      type="number"
                      min="0"
                      step="0.01"
                      value={line.unit_price}
                      onChange={(e) => handleLineChange(index, 'unit_price', parseFloat(e.target.value) || 0)}
                      className="input-field"
                      required
                    />
                  </div>
                  
                  <div className="col-span-2">
                    <label className="block text-xs font-semibold text-gray-600 mb-1">
                      ITBIS %
                    </label>
                    <input
                      type="number"
                      min="0"
                      max="100"
                      step="0.01"
                      value={line.tax_rate * 100}
                      onChange={(e) => handleLineChange(index, 'tax_rate', (parseFloat(e.target.value) || 0) / 100)}
                      className="input-field"
                    />
                  </div>
                  
                  <div className="col-span-1">
                    <button
                      type="button"
                      onClick={() => removeLine(index)}
                      disabled={formData.lines.length === 1}
                      className="p-2 text-red-600 hover:bg-red-50 rounded-lg disabled:opacity-50 disabled:cursor-not-allowed"
                      title="Eliminar línea"
                    >
                      <Trash2 className="w-4 h-4" />
                    </button>
                  </div>
                  
                  <div className="col-span-12 mt-2">
                    <div className="text-right text-sm">
                      <span className="text-gray-600">Total línea: </span>
                      <span className="font-semibold text-[#FF6B00]">
                        ${calculateLineTotal(line).toFixed(2)}
                      </span>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* Retenciones (solo para facturas gubernamentales) */}
          {formData.ncf_type === '15' && (() => {
            const selectedCustomer = customers.find(c => c.id === formData.customer_id);
            const showWithholding = selectedCustomer?.is_government_entity;
            
            if (!showWithholding) return null;
            
            return (
              <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
                <h3 className="text-lg font-semibold text-[#212121] mb-3">Retenciones</h3>
                <div className="space-y-3">
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input
                      type="checkbox"
                      checked={formData.withholding_exempt}
                      onChange={(e) => handleInputChange('withholding_exempt', e.target.checked)}
                      className="w-4 h-4 text-[#FF6B00] border-gray-300 rounded focus:ring-[#FF6B00]"
                    />
                    <span className="text-sm text-gray-700">
                      Exento de Retención (Emisor Electrónico)
                    </span>
                  </label>
                  
                  {!formData.withholding_exempt && (
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        Tasa de Retención (%)
                      </label>
                      <div className="relative">
                        <input
                          type="number"
                          step="0.01"
                          min="0"
                          max="100"
                          value={formData.withholding_rate !== null && formData.withholding_rate !== undefined 
                            ? (formData.withholding_rate * 100) 
                            : (selectedCustomer?.default_withholding_rate 
                              ? (selectedCustomer.default_withholding_rate * 100) 
                              : '5')}
                          onChange={(e) => handleInputChange('withholding_rate', (parseFloat(e.target.value) || 5) / 100)}
                          className="input-field pr-12"
                          placeholder="5"
                        />
                        <span className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 text-sm">%</span>
                      </div>
                      <p className="text-xs text-gray-500 mt-1">
                        Por defecto: {selectedCustomer?.default_withholding_rate ? (selectedCustomer.default_withholding_rate * 100).toFixed(2) : '5'}%
                      </p>
                    </div>
                  )}
                </div>
              </div>
            );
          })()}

          {/* Totales */}
          <div className="bg-[#F5F5F5] rounded-lg p-4">
            <h3 className="text-lg font-semibold text-[#212121] mb-3">Resumen</h3>
            <div className="space-y-2">
              <div className="flex justify-between text-sm">
                <span className="text-gray-600">Subtotal:</span>
                <span className="font-semibold">${totals.subtotal.toFixed(2)}</span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-gray-600">ITBIS:</span>
                <span className="font-semibold">${totals.tax.toFixed(2)}</span>
              </div>
              <div className="border-t border-gray-300 pt-2 flex justify-between">
                <span className="font-bold text-lg">Total:</span>
                <span className="font-bold text-xl text-[#FF6B00]">${totals.total.toFixed(2)}</span>
              </div>
              
              {/* Mostrar retención si aplica */}
              {totals.withholdingAmount > 0 && (
                <>
                  <div className="flex justify-between text-sm pt-2 border-t border-gray-300">
                    <span className="text-gray-600">Retención ISR ({(totals.withholdingRate * 100).toFixed(2)}%):</span>
                    <span className="font-semibold text-red-600">-${totals.withholdingAmount.toFixed(2)}</span>
                  </div>
                  <div className="flex justify-between pt-2 border-t-2 border-gray-400">
                    <span className="font-bold text-lg">Monto Neto a Recibir:</span>
                    <span className="font-bold text-2xl text-green-600">${totals.netAmount.toFixed(2)}</span>
                  </div>
                </>
              )}
            </div>
          </div>

          {/* Notas y términos */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <label className="block text-sm font-semibold text-gray-600 mb-2">
                Notas
              </label>
              <textarea
                value={formData.notes}
                onChange={(e) => handleInputChange('notes', e.target.value)}
                className="input-field"
                rows="4"
                placeholder="Notas adicionales..."
              />
            </div>

            <div>
              <label className="block text-sm font-semibold text-gray-600 mb-2">
                Términos y Condiciones
              </label>
              <textarea
                value={formData.terms}
                onChange={(e) => handleInputChange('terms', e.target.value)}
                className="input-field"
                rows="4"
                placeholder="Términos y condiciones..."
              />
            </div>
          </div>

          {/* Botones */}
          <div className="flex gap-3 pt-6 border-t border-gray-200">
            <button
              type="button"
              onClick={handleClose}
              className="btn-outline flex-1"
              disabled={loading}
            >
              Cancelar
            </button>
            <button
              type="submit"
              disabled={loading}
              className="btn-primary flex-1 flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {loading ? (
                <>
                  <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
                  Creando...
                </>
              ) : (
                <>
                  <Save className="w-4 h-4" />
                  Crear Factura
                </>
              )}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
