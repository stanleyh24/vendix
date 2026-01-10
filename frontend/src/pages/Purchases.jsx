import { useState, useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';
import { ShoppingBag, Plus, Search, Edit2, Trash2, Calendar, DollarSign, Package, CheckCircle, XCircle } from 'lucide-react';
import Alert from '../components/Alert';
import api from '../lib/api';

export default function Purchases() {
  const [searchParams, setSearchParams] = useSearchParams();
  const [purchases, setPurchases] = useState([]);
  const [suppliers, setSuppliers] = useState([]);
  const [products, setProducts] = useState([]);
  const [loading, setLoading] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const [filterSupplier, setFilterSupplier] = useState('all');
  const [filterStatus, setFilterStatus] = useState('all');
  const [showModal, setShowModal] = useState(false);
  const [showReceiveModal, setShowReceiveModal] = useState(false);
  const [editingPurchase, setEditingPurchase] = useState(null);
  const [receivingPurchase, setReceivingPurchase] = useState(null);
  const [alert, setAlert] = useState(null);
  const [deleteConfirm, setDeleteConfirm] = useState(null);
  const [highlightedPurchaseId, setHighlightedPurchaseId] = useState(null);

  // Formulario
  const [formData, setFormData] = useState({
    supplier_id: '',
    purchase_date: new Date().toISOString().split('T')[0],
    expected_delivery_date: '',
    payment_method: 'credit',
    payment_terms: 30, // Días de crédito por defecto
    reference: '',
    notes: '',
    lines: [{ product_id: '', description: '', quantity: 1, unit_price: 0, tax_rate: 18 }],
  });

  // Cargar datos iniciales
  useEffect(() => {
    const loadData = async () => {
      await Promise.all([loadSuppliers(), loadProducts(), loadPurchases()]);
    };
    loadData();
  }, []);

  // Manejar highlight de compra desde notificaciones
  useEffect(() => {
    const highlightId = searchParams.get('highlight');
    if (highlightId) {
      setHighlightedPurchaseId(highlightId);
      // Remover el parámetro de la URL después de un tiempo
      setTimeout(() => {
        setSearchParams({}, { replace: true });
        setHighlightedPurchaseId(null);
      }, 5000); // Remover después de 5 segundos
      
      // Scroll a la compra destacada después de cargar
      setTimeout(() => {
        const element = document.getElementById(`purchase-${highlightId}`);
        if (element) {
          element.scrollIntoView({ behavior: 'smooth', block: 'center' });
        }
      }, 500);
    }
  }, [searchParams, setSearchParams]);

  const loadPurchases = async () => {
    setLoading(true);
    try {
      const response = await api.get('/purchases');
      setPurchases(response.data || []);
    } catch (error) {
      console.error('Error loading purchases:', error);
      showAlert('error', 'Error', 'No se pudieron cargar las compras');
      setPurchases([]);
    } finally {
      setLoading(false);
    }
  };

  const loadSuppliers = async () => {
    try {
      const response = await api.get('/suppliers?active=true');
      setSuppliers(response.data || []);
    } catch (error) {
      console.error('Error loading suppliers:', error);
      setSuppliers([]);
    }
  };

  const loadProducts = async () => {
    try {
      const response = await api.get('/products');
      setProducts(response.data || []);
    } catch (error) {
      console.error('Error loading products:', error);
      setProducts([]);
    }
  };

  const showAlert = (type, title, message) => {
    setAlert({ type, title, message });
    setTimeout(() => setAlert(null), 5000);
  };

  // Filtrar compras
  const filteredPurchases = purchases.filter(purchase => {
    const matchesSearch = purchase.purchase_number?.toLowerCase().includes(searchTerm.toLowerCase()) ||
      purchase.supplier_name?.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesSupplier = filterSupplier === 'all' || purchase.supplier_id === filterSupplier;
    const matchesStatus = filterStatus === 'all' || purchase.status === filterStatus;
    return matchesSearch && matchesSupplier && matchesStatus;
  });

  // Calcular totales del formulario
  const calculateTotals = () => {
    let subtotal = 0;
    let taxAmount = 0;
    formData.lines.forEach(line => {
      const lineSubtotal = (line.quantity || 0) * (line.unit_price || 0);
      const lineTax = lineSubtotal * ((line.tax_rate || 0) / 100);
      subtotal += lineSubtotal;
      taxAmount += lineTax;
    });
    return { subtotal, taxAmount, total: subtotal + taxAmount };
  };

  const totals = calculateTotals();

  // Abrir modal para crear/editar
  const openModal = (purchase = null) => {
    if (purchase) {
      setEditingPurchase(purchase);
      setFormData({
        supplier_id: purchase.supplier_id,
        purchase_date: purchase.purchase_date?.split('T')[0] || new Date().toISOString().split('T')[0],
        expected_delivery_date: purchase.expected_delivery_date?.split('T')[0] || '',
        payment_method: purchase.payment_method || 'credit',
        reference: purchase.reference || '',
        notes: purchase.notes || '',
        lines: purchase.lines?.map(line => ({
          product_id: line.product_id || '',
          description: line.description || '',
          quantity: line.quantity || 1,
          unit_price: line.unit_price || 0,
          tax_rate: line.tax_rate || 18,
        })) || [{ product_id: '', description: '', quantity: 1, unit_price: 0, tax_rate: 18 }],
      });
    } else {
      setEditingPurchase(null);
      setFormData({
        supplier_id: '',
        purchase_date: new Date().toISOString().split('T')[0],
        expected_delivery_date: '',
        payment_method: 'credit',
        reference: '',
        notes: '',
        lines: [{ product_id: '', description: '', quantity: 1, unit_price: 0, tax_rate: 18 }],
      });
    }
    setShowModal(true);
  };

  // Cerrar modal
  const closeModal = () => {
    setShowModal(false);
    setEditingPurchase(null);
  };

  // Agregar línea al formulario
  const addLine = () => {
    setFormData({
      ...formData,
      lines: [...formData.lines, { product_id: '', description: '', quantity: 1, unit_price: 0, tax_rate: 18 }],
    });
  };

  // Eliminar línea del formulario
  const removeLine = (index) => {
    if (formData.lines.length > 1) {
      setFormData({
        ...formData,
        lines: formData.lines.filter((_, i) => i !== index),
      });
    }
  };

  // Actualizar línea del formulario
  const updateLine = (index, field, value) => {
    const newLines = [...formData.lines];
    newLines[index][field] = value;
    
    // Si se selecciona un producto, actualizar descripción y precio
    if (field === 'product_id' && value) {
      const product = products.find(p => p.id === value);
      if (product) {
        newLines[index].description = product.name;
        newLines[index].unit_price = product.cost || product.price || 0;
        newLines[index].tax_rate = product.tax_rate || 18;
      }
    }
    
    setFormData({ ...formData, lines: newLines });
  };

  // Guardar compra
  const handleSave = async () => {
    if (!formData.supplier_id) {
      showAlert('error', 'Error', 'Debe seleccionar un proveedor');
      return;
    }

    if (formData.lines.length === 0 || formData.lines.some(line => !line.description || line.quantity <= 0 || line.unit_price <= 0)) {
      showAlert('error', 'Error', 'Debe agregar al menos una línea válida');
      return;
    }

    try {
      const payload = {
        ...formData,
        payment_terms: formData.payment_method === 'credit' ? formData.payment_terms : undefined,
        lines: formData.lines.map(line => ({
          product_id: line.product_id || null,
          description: line.description,
          quantity: parseFloat(line.quantity),
          unit_price: parseFloat(line.unit_price),
          tax_rate: parseFloat(line.tax_rate),
        })),
      };

      if (editingPurchase) {
        await api.put(`/purchases/${editingPurchase.id}`, payload);
        showAlert('success', 'Éxito', 'Compra actualizada correctamente');
      } else {
        await api.post('/purchases', payload);
        showAlert('success', 'Éxito', 'Compra creada correctamente');
      }

      await loadPurchases();
      closeModal();
    } catch (error) {
      console.error('Error saving purchase:', error);
      showAlert('error', 'Error', error.response?.data?.error || 'Error al guardar la compra');
    }
  };

  // Eliminar compra
  const handleDelete = async (id) => {
    try {
      await api.delete(`/purchases/${id}`);
      showAlert('success', 'Éxito', 'Compra eliminada correctamente');
      await loadPurchases();
      setDeleteConfirm(null);
    } catch (error) {
      console.error('Error deleting purchase:', error);
      showAlert('error', 'Error', error.response?.data?.error || 'Error al eliminar la compra');
    }
  };

  // Abrir modal de recepción
  const openReceiveModal = (purchase) => {
    setReceivingPurchase(purchase);
    setShowReceiveModal(true);
  };

  // Cerrar modal de recepción
  const closeReceiveModal = () => {
    setShowReceiveModal(false);
    setReceivingPurchase(null);
  };

  // Recibir compra
  const handleReceive = async (receivedLines) => {
    try {
      await api.post(`/purchases/${receivingPurchase.id}/receive`, {
        received_date: new Date().toISOString().split('T')[0],
        lines: receivedLines.map(line => ({
          line_id: line.id,
          received_quantity: parseFloat(line.received_quantity),
        })),
      });
      showAlert('success', 'Éxito', 'Compra recibida correctamente');
      await loadPurchases();
      closeReceiveModal();
    } catch (error) {
      console.error('Error receiving purchase:', error);
      showAlert('error', 'Error', error.response?.data?.error || 'Error al recibir la compra');
    }
  };

  const getStatusBadge = (status) => {
    const badges = {
      draft: { label: 'Borrador', color: 'bg-gray-100 text-gray-800' },
      pending: { label: 'Pendiente', color: 'bg-yellow-100 text-yellow-800' },
      received: { label: 'Recibida', color: 'bg-green-100 text-green-800' },
      cancelled: { label: 'Cancelada', color: 'bg-red-100 text-red-800' },
    };
    const badge = badges[status] || badges.draft;
    return <span className={`px-2 py-1 rounded-full text-xs font-medium ${badge.color}`}>{badge.label}</span>;
  };

  const getPaymentStatusBadge = (status) => {
    const badges = {
      pending: { label: 'Pendiente', color: 'bg-yellow-100 text-yellow-800' },
      partial: { label: 'Parcial', color: 'bg-blue-100 text-blue-800' },
      paid: { label: 'Pagado', color: 'bg-green-100 text-green-800' },
      cancelled: { label: 'Cancelado', color: 'bg-red-100 text-red-800' },
    };
    const badge = badges[status] || badges.pending;
    return <span className={`px-2 py-1 rounded-full text-xs font-medium ${badge.color}`}>{badge.label}</span>;
  };

  return (
    <div className="p-6">
      {alert && <Alert type={alert.type} title={alert.title} message={alert.message} onClose={() => setAlert(null)} />}

      {/* Header */}
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-[#212121] flex items-center gap-3">
          <ShoppingBag className="w-8 h-8 text-[#FF6B00]" />
          Compras
        </h1>
        <p className="text-gray-600 mt-1">Gestiona las compras de productos a proveedores</p>
      </div>

      {/* Filtros y búsqueda */}
      <div className="bg-white rounded-lg shadow-sm p-4 mb-6">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
            <input
              type="text"
              placeholder="Buscar por número o proveedor..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#FF6B00] focus:border-transparent"
            />
          </div>
          <select
            value={filterSupplier}
            onChange={(e) => setFilterSupplier(e.target.value)}
            className="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#FF6B00] focus:border-transparent"
          >
            <option value="all">Todos los proveedores</option>
            {suppliers.map(supplier => (
              <option key={supplier.id} value={supplier.id}>{supplier.name}</option>
            ))}
          </select>
          <select
            value={filterStatus}
            onChange={(e) => setFilterStatus(e.target.value)}
            className="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#FF6B00] focus:border-transparent"
          >
            <option value="all">Todos los estados</option>
            <option value="draft">Borrador</option>
            <option value="pending">Pendiente</option>
            <option value="received">Recibida</option>
            <option value="cancelled">Cancelada</option>
          </select>
          <button
            onClick={() => openModal()}
            className="bg-[#FF6B00] text-white px-4 py-2 rounded-lg hover:bg-[#E55A00] transition-colors flex items-center justify-center gap-2"
          >
            <Plus className="w-5 h-5" />
            Nueva Compra
          </button>
        </div>
      </div>

      {/* Tabla de compras */}
      <div className="bg-white rounded-lg shadow-sm overflow-hidden">
        {loading ? (
          <div className="p-8 text-center text-gray-500">Cargando...</div>
        ) : filteredPurchases.length === 0 ? (
          <div className="p-8 text-center text-gray-500">No hay compras registradas</div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="bg-gray-50 border-b border-gray-200">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Número</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Proveedor</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Fecha</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Vencimiento</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Total</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Estado</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Pago</th>
                  <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Acciones</th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {filteredPurchases.map((purchase) => {
                  // Calcular si está próxima a vencer o vencida
                  let isOverdue = false;
                  let isDueSoon = false;
                  if (purchase.payment_due_date && purchase.payment_method === 'credit' && purchase.payment_status !== 'paid') {
                    const dueDate = new Date(purchase.payment_due_date);
                    const today = new Date();
                    const diffTime = dueDate - today;
                    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
                    isOverdue = diffDays < 0;
                    isDueSoon = diffDays >= 0 && diffDays <= 7;
                  }
                  
                  const isHighlighted = highlightedPurchaseId === purchase.id;
                  
                  return (
                  <tr 
                    key={purchase.id}
                    id={`purchase-${purchase.id}`}
                    className={`hover:bg-gray-50 transition-all ${
                      isHighlighted 
                        ? 'bg-yellow-100 border-l-4 border-l-yellow-500 shadow-md' 
                        : isOverdue 
                          ? 'bg-red-50 border-l-4 border-l-red-500' 
                          : isDueSoon 
                            ? 'bg-orange-50 border-l-4 border-l-orange-500' 
                            : ''
                    }`}
                  >
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="text-sm font-medium text-gray-900">{purchase.purchase_number}</div>
                      {purchase.reference && (
                        <div className="text-sm text-gray-500">Ref: {purchase.reference}</div>
                      )}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="text-sm text-gray-900">{purchase.supplier_name || 'N/A'}</div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="text-sm text-gray-900">
                        {new Date(purchase.purchase_date).toLocaleDateString('es-DO')}
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      {purchase.payment_due_date ? (() => {
                        const dueDate = new Date(purchase.payment_due_date);
                        const today = new Date();
                        const diffTime = dueDate - today;
                        const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
                        const isOverdue = diffDays < 0;
                        const isDueSoon = diffDays >= 0 && diffDays <= 7;
                        
                        return (
                          <div className="text-sm">
                            <div className={`font-medium ${
                              isOverdue ? 'text-red-600' : isDueSoon ? 'text-orange-600' : 'text-gray-900'
                            }`}>
                              {dueDate.toLocaleDateString('es-DO')}
                            </div>
                            {isOverdue && (
                              <div className="text-xs text-red-600 font-medium">
                                Vencida hace {Math.abs(diffDays)} días
                              </div>
                            )}
                            {isDueSoon && !isOverdue && (
                              <div className="text-xs text-orange-600 font-medium">
                                Vence en {diffDays} días
                              </div>
                            )}
                          </div>
                        );
                      })() : (
                        <div className="text-sm text-gray-400">-</div>
                      )}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="text-sm font-medium text-gray-900">
                        ${purchase.total?.toFixed(2) || '0.00'}
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      {getStatusBadge(purchase.status)}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      {getPaymentStatusBadge(purchase.payment_status)}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                      <div className="flex items-center justify-end gap-2">
                        {purchase.status !== 'received' && purchase.status !== 'cancelled' && (
                          <button
                            onClick={() => openReceiveModal(purchase)}
                            className="text-green-600 hover:text-green-900"
                            title="Recibir compra"
                          >
                            <CheckCircle className="w-5 h-5" />
                          </button>
                        )}
                        {purchase.status === 'draft' && (
                          <>
                            <button
                              onClick={() => openModal(purchase)}
                              className="text-blue-600 hover:text-blue-900"
                              title="Editar"
                            >
                              <Edit2 className="w-5 h-5" />
                            </button>
                            <button
                              onClick={() => setDeleteConfirm(purchase.id)}
                              className="text-red-600 hover:text-red-900"
                              title="Eliminar"
                            >
                              <Trash2 className="w-5 h-5" />
                            </button>
                          </>
                        )}
                      </div>
                    </td>
                  </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* Modal de crear/editar compra */}
      {showModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-lg shadow-xl max-w-4xl w-full max-h-[90vh] overflow-y-auto">
            <div className="p-6 border-b border-gray-200">
              <h2 className="text-2xl font-bold text-gray-900">
                {editingPurchase ? 'Editar Compra' : 'Nueva Compra'}
              </h2>
            </div>
            <div className="p-6 space-y-4">
              {/* Información básica */}
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Proveedor *</label>
                  <select
                    value={formData.supplier_id}
                    onChange={(e) => setFormData({ ...formData, supplier_id: e.target.value })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#FF6B00] focus:border-transparent"
                    required
                  >
                    <option value="">Seleccionar proveedor</option>
                    {suppliers.map(supplier => (
                      <option key={supplier.id} value={supplier.id}>{supplier.name}</option>
                    ))}
                  </select>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Fecha de Compra *</label>
                  <input
                    type="date"
                    value={formData.purchase_date}
                    onChange={(e) => setFormData({ ...formData, purchase_date: e.target.value })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#FF6B00] focus:border-transparent"
                    required
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Fecha Esperada de Entrega</label>
                  <input
                    type="date"
                    value={formData.expected_delivery_date}
                    onChange={(e) => setFormData({ ...formData, expected_delivery_date: e.target.value })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#FF6B00] focus:border-transparent"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Método de Pago *</label>
                  <select
                    value={formData.payment_method}
                    onChange={(e) => setFormData({ ...formData, payment_method: e.target.value })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#FF6B00] focus:border-transparent"
                    required
                  >
                    <option value="credit">Crédito</option>
                    <option value="cash">Efectivo</option>
                    <option value="bank">Banco</option>
                    <option value="check">Cheque</option>
                  </select>
                </div>
                {formData.payment_method === 'credit' && (
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">Términos de Crédito (días)</label>
                    <input
                      type="number"
                      min="1"
                      value={formData.payment_terms}
                      onChange={(e) => setFormData({ ...formData, payment_terms: parseInt(e.target.value) || 30 })}
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#FF6B00] focus:border-transparent"
                      placeholder="30"
                    />
                    <p className="text-xs text-gray-500 mt-1">Días de crédito para el pago</p>
                  </div>
                )}
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Referencia</label>
                  <input
                    type="text"
                    value={formData.reference}
                    onChange={(e) => setFormData({ ...formData, reference: e.target.value })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#FF6B00] focus:border-transparent"
                    placeholder="Número de factura del proveedor"
                  />
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Notas</label>
                <textarea
                  value={formData.notes}
                  onChange={(e) => setFormData({ ...formData, notes: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#FF6B00] focus:border-transparent"
                  rows="2"
                  placeholder="Notas adicionales..."
                />
              </div>

              {/* Líneas de compra */}
              <div>
                <div className="flex items-center justify-between mb-3">
                  <label className="block text-sm font-medium text-gray-700">Productos / Líneas de Compra *</label>
                  <button
                    type="button"
                    onClick={addLine}
                    className="text-sm text-[#FF6B00] hover:text-[#E55A00] flex items-center gap-1"
                  >
                    <Plus className="w-4 h-4" />
                    Agregar Línea
                  </button>
                </div>
                <div className="space-y-3">
                  {formData.lines.map((line, index) => (
                    <div key={index} className="border border-gray-200 rounded-lg p-4">
                      <div className="grid grid-cols-12 gap-3">
                        <div className="col-span-12 md:col-span-4">
                          <label className="block text-xs font-medium text-gray-700 mb-1">Producto</label>
                          <select
                            value={line.product_id}
                            onChange={(e) => updateLine(index, 'product_id', e.target.value)}
                            className="w-full px-2 py-1 text-sm border border-gray-300 rounded focus:ring-2 focus:ring-[#FF6B00] focus:border-transparent"
                          >
                            <option value="">Seleccionar producto</option>
                            {products.map(product => (
                              <option key={product.id} value={product.id}>{product.name}</option>
                            ))}
                          </select>
                        </div>
                        <div className="col-span-12 md:col-span-4">
                          <label className="block text-xs font-medium text-gray-700 mb-1">Descripción *</label>
                          <input
                            type="text"
                            value={line.description}
                            onChange={(e) => updateLine(index, 'description', e.target.value)}
                            className="w-full px-2 py-1 text-sm border border-gray-300 rounded focus:ring-2 focus:ring-[#FF6B00] focus:border-transparent"
                            required
                          />
                        </div>
                        <div className="col-span-4 md:col-span-1">
                          <label className="block text-xs font-medium text-gray-700 mb-1">Cantidad *</label>
                          <input
                            type="number"
                            value={line.quantity}
                            onChange={(e) => updateLine(index, 'quantity', parseFloat(e.target.value) || 0)}
                            className="w-full px-2 py-1 text-sm border border-gray-300 rounded focus:ring-2 focus:ring-[#FF6B00] focus:border-transparent"
                            min="0.01"
                            step="0.01"
                            required
                          />
                        </div>
                        <div className="col-span-4 md:col-span-1">
                          <label className="block text-xs font-medium text-gray-700 mb-1">Precio *</label>
                          <input
                            type="number"
                            value={line.unit_price}
                            onChange={(e) => updateLine(index, 'unit_price', parseFloat(e.target.value) || 0)}
                            className="w-full px-2 py-1 text-sm border border-gray-300 rounded focus:ring-2 focus:ring-[#FF6B00] focus:border-transparent"
                            min="0"
                            step="0.01"
                            required
                          />
                        </div>
                        <div className="col-span-4 md:col-span-1">
                          <label className="block text-xs font-medium text-gray-700 mb-1">ITBIS %</label>
                          <input
                            type="number"
                            value={line.tax_rate}
                            onChange={(e) => updateLine(index, 'tax_rate', parseFloat(e.target.value) || 0)}
                            className="w-full px-2 py-1 text-sm border border-gray-300 rounded focus:ring-2 focus:ring-[#FF6B00] focus:border-transparent"
                            min="0"
                            step="0.01"
                          />
                        </div>
                        <div className="col-span-12 md:col-span-1 flex items-end">
                          <button
                            type="button"
                            onClick={() => removeLine(index)}
                            className="w-full px-2 py-1 text-red-600 hover:bg-red-50 rounded text-sm"
                            disabled={formData.lines.length === 1}
                          >
                            <Trash2 className="w-4 h-4 mx-auto" />
                          </button>
                        </div>
                      </div>
                      <div className="mt-2 text-xs text-gray-600">
                        Subtotal: ${((line.quantity || 0) * (line.unit_price || 0)).toFixed(2)} | 
                        ITBIS: ${(((line.quantity || 0) * (line.unit_price || 0)) * ((line.tax_rate || 0) / 100)).toFixed(2)} | 
                        Total: ${(((line.quantity || 0) * (line.unit_price || 0)) * (1 + (line.tax_rate || 0) / 100)).toFixed(2)}
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              {/* Totales */}
              <div className="border-t border-gray-200 pt-4">
                <div className="flex justify-end">
                  <div className="w-64 space-y-2">
                    <div className="flex justify-between text-sm">
                      <span className="text-gray-600">Subtotal:</span>
                      <span className="font-medium">${totals.subtotal.toFixed(2)}</span>
                    </div>
                    <div className="flex justify-between text-sm">
                      <span className="text-gray-600">ITBIS:</span>
                      <span className="font-medium">${totals.taxAmount.toFixed(2)}</span>
                    </div>
                    <div className="flex justify-between text-lg font-bold border-t border-gray-200 pt-2">
                      <span>Total:</span>
                      <span className="text-[#FF6B00]">${totals.total.toFixed(2)}</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <div className="p-6 border-t border-gray-200 flex justify-end gap-3">
              <button
                onClick={closeModal}
                className="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50"
              >
                Cancelar
              </button>
              <button
                onClick={handleSave}
                className="px-4 py-2 bg-[#FF6B00] text-white rounded-lg hover:bg-[#E55A00]"
              >
                {editingPurchase ? 'Actualizar' : 'Crear'} Compra
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Modal de recepción */}
      {showReceiveModal && receivingPurchase && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-lg shadow-xl max-w-3xl w-full">
            <div className="p-6 border-b border-gray-200">
              <h2 className="text-2xl font-bold text-gray-900">Recibir Compra</h2>
              <p className="text-sm text-gray-600 mt-1">Compra: {receivingPurchase.purchase_number}</p>
            </div>
            <div className="p-6">
              <ReceivePurchaseForm
                purchase={receivingPurchase}
                onReceive={handleReceive}
                onCancel={closeReceiveModal}
              />
            </div>
          </div>
        </div>
      )}

      {/* Confirmación de eliminación */}
      {deleteConfirm && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg shadow-xl p-6 max-w-md w-full">
            <h3 className="text-lg font-bold text-gray-900 mb-4">Confirmar Eliminación</h3>
            <p className="text-gray-600 mb-6">¿Está seguro de que desea eliminar esta compra? Esta acción no se puede deshacer.</p>
            <div className="flex justify-end gap-3">
              <button
                onClick={() => setDeleteConfirm(null)}
                className="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50"
              >
                Cancelar
              </button>
              <button
                onClick={() => handleDelete(deleteConfirm)}
                className="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700"
              >
                Eliminar
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

// Componente para recibir compra
function ReceivePurchaseForm({ purchase, onReceive, onCancel }) {
  const [receivedLines, setReceivedLines] = useState(
    purchase.lines?.map(line => ({
      id: line.id,
      description: line.description,
      quantity: line.quantity,
      received_quantity: line.received_quantity || 0,
    })) || []
  );

  const updateReceivedQuantity = (index, quantity) => {
    const newLines = [...receivedLines];
    newLines[index].received_quantity = parseFloat(quantity) || 0;
    if (newLines[index].received_quantity > newLines[index].quantity) {
      newLines[index].received_quantity = newLines[index].quantity;
    }
    setReceivedLines(newLines);
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    onReceive(receivedLines);
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="space-y-3">
        {receivedLines.map((line, index) => (
          <div key={line.id} className="border border-gray-200 rounded-lg p-4">
            <div className="flex items-center justify-between mb-2">
              <div className="flex-1">
                <div className="font-medium text-gray-900">{line.description}</div>
                <div className="text-sm text-gray-600">
                  Cantidad ordenada: {line.quantity}
                </div>
              </div>
              <div className="flex items-center gap-3">
                <label className="text-sm font-medium text-gray-700">Recibida:</label>
                <input
                  type="number"
                  value={line.received_quantity}
                  onChange={(e) => updateReceivedQuantity(index, e.target.value)}
                  className="w-24 px-2 py-1 border border-gray-300 rounded focus:ring-2 focus:ring-[#FF6B00] focus:border-transparent"
                  min="0"
                  max={line.quantity}
                  step="0.01"
                  required
                />
              </div>
            </div>
          </div>
        ))}
      </div>
      <div className="flex justify-end gap-3 pt-4 border-t border-gray-200">
        <button
          type="button"
          onClick={onCancel}
          className="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50"
        >
          Cancelar
        </button>
        <button
          type="submit"
          className="px-4 py-2 bg-[#FF6B00] text-white rounded-lg hover:bg-[#E55A00]"
        >
          Recibir Compra
        </button>
      </div>
    </form>
  );
}

