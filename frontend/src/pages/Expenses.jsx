import { useState, useEffect } from 'react';
import { Receipt, Plus, Search, Edit2, Trash2, Calendar, DollarSign, FileText, Tag } from 'lucide-react';
import Alert from '../components/Alert';
import api from '../lib/api';

export default function Expenses() {
  const [expenses, setExpenses] = useState([]);
  const [loading, setLoading] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const [filterCategory, setFilterCategory] = useState('all');
  const [filterPeriod, setFilterPeriod] = useState('all'); // all, month, week
  const [showModal, setShowModal] = useState(false);
  const [editingExpense, setEditingExpense] = useState(null);
  const [alert, setAlert] = useState(null);
  const [deleteConfirm, setDeleteConfirm] = useState(null);

  // Formulario
  const [formData, setFormData] = useState({
    category: 'supplies',
    description: '',
    amount: '',
    payment_date: new Date().toISOString().split('T')[0],
    supplier: '',
    payment_method: 'cash',
    reference: '',
    notes: '',
  });

  // Categorías de gastos
  const categories = [
    { value: 'supplies', label: 'Suministros', color: 'bg-blue-500' },
    { value: 'rent', label: 'Alquiler', color: 'bg-purple-500' },
    { value: 'utilities', label: 'Servicios', color: 'bg-green-500' },
    { value: 'salaries', label: 'Salarios', color: 'bg-orange-500' },
    { value: 'marketing', label: 'Marketing', color: 'bg-pink-500' },
    { value: 'taxes', label: 'Impuestos', color: 'bg-red-500' },
    { value: 'other', label: 'Otros', color: 'bg-gray-500' },
  ];

  // Cargar gastos del backend
  useEffect(() => {
    loadExpenses();
  }, []);

  const loadExpenses = async () => {
    setLoading(true);
    try {
      const response = await api.get('/expenses');
      const data = response.data || [];
      
      // Map backend data to frontend format
      const mappedExpenses = data.map(exp => ({
        id: exp.id,
        category: exp.payment_method || 'other',
        description: exp.notes || 'Sin descripción',
        amount: exp.amount,
        date: exp.payment_date?.split('T')[0] || exp.payment_date,
        supplier: exp.reference || 'Sin proveedor',
        invoice_number: exp.payment_number,
        created_at: exp.created_at
      }));
      
      setExpenses(mappedExpenses);
    } catch (error) {
      console.error('Error loading expenses:', error);
      showAlert('error', 'Error', 'No se pudieron cargar los gastos');
      setExpenses([]);
    } finally {
      setLoading(false);
    }
  };

  const showAlert = (type, title, message) => {
    setAlert({ type, title, message });
    setTimeout(() => setAlert(null), 5000);
  };

  // Filtrar gastos
  const filteredExpenses = expenses.filter(expense => {
    const matchesSearch = 
      expense.description.toLowerCase().includes(searchTerm.toLowerCase()) ||
      expense.supplier.toLowerCase().includes(searchTerm.toLowerCase()) ||
      expense.invoice_number.toLowerCase().includes(searchTerm.toLowerCase());
    
    const matchesCategory = filterCategory === 'all' || expense.category === filterCategory;

    return matchesSearch && matchesCategory;
  });

  // Abrir modal para crear
  const handleCreate = () => {
    setEditingExpense(null);
    setFormData({
      category: 'supplies',
      description: '',
      amount: '',
      payment_date: new Date().toISOString().split('T')[0],
      supplier: '',
      payment_method: 'cash',
      reference: '',
      notes: '',
    });
    setShowModal(true);
  };

  // Abrir modal para editar
  const handleEdit = (expense) => {
    setEditingExpense(expense);
    setFormData({
      category: expense.category,
      description: expense.description,
      amount: expense.amount.toString(),
      payment_date: expense.date,
      supplier: expense.supplier,
      payment_method: 'cash',
      reference: expense.supplier,
      notes: expense.description,
    });
    setShowModal(true);
  };

  // Guardar gasto
  const handleSubmit = async (e) => {
    e.preventDefault();

    try {
      const payload = {
        description: formData.description,
        amount: parseFloat(formData.amount),
        payment_date: formData.payment_date,
        supplier: formData.supplier,
        payment_method: formData.payment_method,
        reference: formData.reference,
        notes: formData.notes,
        category: formData.category,
        currency: 'DOP'
      };

      if (editingExpense) {
        await api.put(`/expenses/${editingExpense.id}`, payload);
        showAlert('success', 'Gasto actualizado', 'El gasto se actualizó correctamente');
      } else {
        await api.post('/expenses', payload);
        showAlert('success', 'Gasto registrado', 'El gasto se registró correctamente');
      }
      
      setShowModal(false);
      loadExpenses(); // Reload data
    } catch (error) {
      showAlert('error', 'Error', error.response?.data?.error || 'No se pudo guardar el gasto');
    }
  };

  // Eliminar gasto
  const handleDelete = (expense) => {
    setDeleteConfirm(expense);
  };

  const confirmDelete = async () => {
    try {
      await api.delete(`/expenses/${deleteConfirm.id}`);
      showAlert('success', 'Gasto eliminado', 'El gasto se eliminó correctamente');
      setDeleteConfirm(null);
      loadExpenses(); // Reload data
    } catch (error) {
      showAlert('error', 'Error', error.response?.data?.error || 'No se pudo eliminar el gasto');
      setDeleteConfirm(null);
    }
  };

  // Calcular total
  const calculateTotal = () => {
    return filteredExpenses.reduce((sum, expense) => sum + expense.amount, 0);
  };

  const getCategoryColor = (category) => {
    return categories.find(c => c.value === category)?.color || 'bg-gray-500';
  };

  const getCategoryLabel = (category) => {
    return categories.find(c => c.value === category)?.label || category;
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-[#212121] flex items-center gap-3">
            <Receipt className="w-8 h-8 text-[#FF6B00]" />
            Gastos
          </h1>
          <p className="text-gray-600 mt-1">Registra y controla tus gastos operativos</p>
        </div>
        <button
          onClick={handleCreate}
          className="btn-primary flex items-center gap-2"
        >
          <Plus className="w-5 h-5" />
          Nuevo Gasto
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

      {/* Resumen de gastos */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="card bg-gradient-to-br from-red-50 to-white">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600">Total de Gastos</p>
              <p className="text-3xl font-bold text-[#D32F2F] mt-1">
                ${calculateTotal().toFixed(2)}
              </p>
            </div>
            <div className="w-12 h-12 bg-red-100 rounded-full flex items-center justify-center">
              <DollarSign className="w-6 h-6 text-red-600" />
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600">Gastos Registrados</p>
              <p className="text-3xl font-bold text-[#212121] mt-1">
                {filteredExpenses.length}
              </p>
            </div>
            <div className="w-12 h-12 bg-orange-100 rounded-full flex items-center justify-center">
              <FileText className="w-6 h-6 text-[#FF6B00]" />
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600">Promedio por Gasto</p>
              <p className="text-3xl font-bold text-[#212121] mt-1">
                ${filteredExpenses.length > 0 ? (calculateTotal() / filteredExpenses.length).toFixed(2) : '0.00'}
              </p>
            </div>
            <div className="w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center">
              <Receipt className="w-6 h-6 text-blue-600" />
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
                placeholder="Buscar por descripción, proveedor o número de factura..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="input-field pl-10"
              />
            </div>
          </div>

          {/* Filtros por categoría */}
          <div className="flex flex-wrap gap-2">
            <button
              onClick={() => setFilterCategory('all')}
              className={`px-4 py-2 rounded-lg font-medium transition-colors text-sm ${
                filterCategory === 'all'
                  ? 'bg-[#FF6B00] text-white'
                  : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
              }`}
            >
              Todos ({expenses.length})
            </button>
            {categories.map(cat => (
              <button
                key={cat.value}
                onClick={() => setFilterCategory(cat.value)}
                className={`px-4 py-2 rounded-lg font-medium transition-colors text-sm ${
                  filterCategory === cat.value
                    ? `${cat.color} text-white`
                    : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                }`}
              >
                {cat.label} ({expenses.filter(e => e.category === cat.value).length})
              </button>
            ))}
          </div>
        </div>
      </div>

      {/* Lista de gastos */}
      {loading ? (
        <div className="card text-center py-12">
          <div className="inline-block w-8 h-8 border-4 border-[#FF6B00] border-t-transparent rounded-full animate-spin"></div>
          <p className="mt-4 text-gray-600">Cargando gastos...</p>
        </div>
      ) : filteredExpenses.length === 0 ? (
        <div className="card text-center py-12">
          <Receipt className="w-16 h-16 text-gray-300 mx-auto mb-4" />
          <h3 className="text-lg font-semibold text-gray-600 mb-2">
            {searchTerm ? 'No se encontraron gastos' : 'No hay gastos registrados'}
          </h3>
          <p className="text-gray-500 mb-6">
            {searchTerm ? 'Intenta con otro término de búsqueda' : 'Comienza registrando tu primer gasto'}
          </p>
          {!searchTerm && (
            <button onClick={handleCreate} className="btn-primary">
              <Plus className="w-5 h-5 inline mr-2" />
              Registrar Gasto
            </button>
          )}
        </div>
      ) : (
        <div className="space-y-3">
          {filteredExpenses.map((expense) => (
            <div
              key={expense.id}
              className="card hover:shadow-card-hover transition-all duration-200"
            >
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <div className="flex items-center gap-3 mb-2">
                    <span className={`badge ${getCategoryColor(expense.category)} text-white`}>
                      {getCategoryLabel(expense.category)}
                    </span>
                    <div className="flex items-center gap-1 text-sm text-gray-500">
                      <Calendar className="w-4 h-4" />
                      {new Date(expense.date).toLocaleDateString('es-DO', {
                        year: 'numeric',
                        month: 'long',
                        day: 'numeric'
                      })}
                    </div>
                  </div>

                  <h3 className="text-lg font-semibold text-[#212121] mb-2">
                    {expense.description}
                  </h3>

                  <div className="grid grid-cols-1 md:grid-cols-3 gap-3 text-sm">
                    <div className="flex items-center gap-2">
                      <Tag className="w-4 h-4 text-gray-400" />
                      <span className="text-gray-600">
                        Proveedor: <span className="font-medium">{expense.supplier}</span>
                      </span>
                    </div>
                    {expense.invoice_number && (
                      <div className="flex items-center gap-2">
                        <FileText className="w-4 h-4 text-gray-400" />
                        <span className="text-gray-600">
                          Factura: <span className="font-medium">{expense.invoice_number}</span>
                        </span>
                      </div>
                    )}
                  </div>
                </div>

                <div className="flex items-start gap-3">
                  <div className="text-right">
                    <p className="text-2xl font-bold text-[#D32F2F]">
                      ${expense.amount.toFixed(2)}
                    </p>
                  </div>
                  <div className="flex gap-1">
                    <button
                      onClick={() => handleEdit(expense)}
                      className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
                      title="Editar"
                    >
                      <Edit2 className="w-4 h-4 text-gray-600" />
                    </button>
                    <button
                      onClick={() => handleDelete(expense)}
                      className="p-2 hover:bg-red-50 rounded-lg transition-colors"
                      title="Eliminar"
                    >
                      <Trash2 className="w-4 h-4 text-red-600" />
                    </button>
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Modal de crear/editar */}
      {showModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-2xl shadow-card-hover max-w-2xl w-full max-h-[90vh] overflow-y-auto">
            <div className="sticky top-0 bg-white border-b border-gray-200 px-6 py-4 rounded-t-2xl">
              <h2 className="text-2xl font-bold text-[#212121]">
                {editingExpense ? 'Editar Gasto' : 'Nuevo Gasto'}
              </h2>
            </div>

            <form onSubmit={handleSubmit} className="p-6 space-y-4">
              {/* Categoría */}
              <div>
                <label className="block text-sm font-medium text-[#212121] mb-2">
                  Categoría *
                </label>
                <select
                  required
                  value={formData.category}
                  onChange={(e) => setFormData({ ...formData, category: e.target.value })}
                  className="input-field"
                >
                  {categories.map(cat => (
                    <option key={cat.value} value={cat.value}>{cat.label}</option>
                  ))}
                </select>
              </div>

              {/* Descripción */}
              <div>
                <label className="block text-sm font-medium text-[#212121] mb-2">
                  Descripción *
                </label>
                <textarea
                  required
                  value={formData.description}
                  onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                  className="input-field"
                  rows="3"
                  placeholder="Describe el gasto..."
                />
              </div>

              {/* Monto y Fecha */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-[#212121] mb-2">
                    Monto * (RD$)
                  </label>
                  <input
                    type="number"
                    required
                    step="0.01"
                    min="0"
                    value={formData.amount}
                    onChange={(e) => setFormData({ ...formData, amount: e.target.value })}
                    className="input-field"
                    placeholder="5000.00"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-[#212121] mb-2">
                    Fecha *
                  </label>
                  <input
                    type="date"
                    required
                    value={formData.payment_date}
                    onChange={(e) => setFormData({ ...formData, payment_date: e.target.value })}
                    className="input-field"
                  />
                </div>
              </div>

              {/* Proveedor y Método de Pago */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-[#212121] mb-2">
                    Proveedor *
                  </label>
                  <input
                    type="text"
                    required
                    value={formData.supplier}
                    onChange={(e) => setFormData({ ...formData, supplier: e.target.value })}
                    className="input-field"
                    placeholder="Nombre del proveedor"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-[#212121] mb-2">
                    Método de Pago *
                  </label>
                  <select
                    required
                    value={formData.payment_method}
                    onChange={(e) => setFormData({ ...formData, payment_method: e.target.value })}
                    className="input-field"
                  >
                    <option value="cash">Efectivo</option>
                    <option value="card">Tarjeta</option>
                    <option value="transfer">Transferencia</option>
                    <option value="check">Cheque</option>
                  </select>
                </div>
              </div>

              {/* Referencia */}
              <div>
                <label className="block text-sm font-medium text-[#212121] mb-2">
                  Referencia (Número de Factura)
                </label>
                <input
                  type="text"
                  value={formData.reference}
                  onChange={(e) => setFormData({ ...formData, reference: e.target.value })}
                  className="input-field"
                  placeholder="INV-2025-001"
                />
              </div>

              {/* Notas */}
              <div>
                <label className="block text-sm font-medium text-[#212121] mb-2">
                  Notas
                </label>
                <textarea
                  value={formData.notes}
                  onChange={(e) => setFormData({ ...formData, notes: e.target.value })}
                  className="input-field"
                  rows="2"
                  placeholder="Notas adicionales (opcional)"
                />
              </div>

              {/* Botones */}
              <div className="flex gap-3 pt-4 border-t border-gray-200">
                <button
                  type="submit"
                  className="btn-primary flex-1"
                >
                  {editingExpense ? 'Actualizar Gasto' : 'Registrar Gasto'}
                </button>
                <button
                  type="button"
                  onClick={() => setShowModal(false)}
                  className="btn-outline flex-1"
                >
                  Cancelar
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Modal de confirmación de eliminación */}
      {deleteConfirm && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-2xl shadow-card-hover max-w-md w-full p-6">
            <div className="text-center">
              <div className="w-16 h-16 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-4">
                <Trash2 className="w-8 h-8 text-red-600" />
              </div>
              <h3 className="text-xl font-bold text-[#212121] mb-2">
                ¿Eliminar Gasto?
              </h3>
              <p className="text-gray-600 mb-1">
                Estás a punto de eliminar:
              </p>
              <p className="font-semibold text-[#212121] mb-2">
                {deleteConfirm.description}
              </p>
              <p className="text-lg font-bold text-[#D32F2F] mb-4">
                ${deleteConfirm.amount.toFixed(2)}
              </p>
              <p className="text-sm text-gray-500 mb-6">
                Esta acción no se puede deshacer.
              </p>
              <div className="flex gap-3">
                <button
                  onClick={() => setDeleteConfirm(null)}
                  className="btn-outline flex-1"
                >
                  Cancelar
                </button>
                <button
                  onClick={confirmDelete}
                  className="btn-error flex-1"
                >
                  Eliminar
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

