import { useState, useEffect } from 'react';
import { Users, Plus, Search, Edit2, Trash2, Mail, Phone, MapPin, CreditCard, Building2, User } from 'lucide-react';
import api from '../lib/api';
import Alert from '../components/Alert';

export default function Customers() {
  const [customers, setCustomers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState('');
  const [filterActive, setFilterActive] = useState('all'); // all, active, inactive
  const [filterType, setFilterType] = useState('all'); // all, individual, business
  const [showModal, setShowModal] = useState(false);
  const [editingCustomer, setEditingCustomer] = useState(null);
  const [alert, setAlert] = useState(null);
  const [deleteConfirm, setDeleteConfirm] = useState(null);

  // Formulario
  const [formData, setFormData] = useState({
    customer_type: 'individual',
    tax_id: '',
    name: '',
    email: '',
    phone: '',
    address: '',
    city: '',
    state: '',
    postal_code: '',
    country: 'República Dominicana',
    is_government_entity: false,
    default_withholding_rate: '',
  });

  // Cargar clientes
  useEffect(() => {
    loadCustomers();
  }, []);

  const loadCustomers = async () => {
    try {
      setLoading(true);
      const response = await api.get('/customers');
      setCustomers(response.data || []);
    } catch (error) {
      showAlert('error', 'Error al cargar clientes', error.response?.data?.error || 'Error desconocido');
    } finally {
      setLoading(false);
    }
  };

  const showAlert = (type, title, message) => {
    setAlert({ type, title, message });
    setTimeout(() => setAlert(null), 5000);
  };

  // Filtrar clientes
  const filteredCustomers = customers.filter(customer => {
    const matchesSearch = 
      customer.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      (customer.email && customer.email.toLowerCase().includes(searchTerm.toLowerCase())) ||
      (customer.tax_id && customer.tax_id.toLowerCase().includes(searchTerm.toLowerCase()));
    
    const matchesFilter = filterActive === 'all' ? true :
                         filterActive === 'active' ? customer.is_active :
                         !customer.is_active;

    const matchesType = filterType === 'all' ? true : customer.customer_type === filterType;

    return matchesSearch && matchesFilter && matchesType;
  });

  // Abrir modal para crear
  const handleCreate = () => {
    setEditingCustomer(null);
    setFormData({
      customer_type: 'individual',
      tax_id: '',
      name: '',
      email: '',
      phone: '',
      address: '',
      city: '',
      state: '',
      postal_code: '',
      country: 'República Dominicana',
      is_government_entity: false,
      default_withholding_rate: '',
    });
    setShowModal(true);
  };

  // Abrir modal para editar
  const handleEdit = (customer) => {
    setEditingCustomer(customer);
    setFormData({
      customer_type: customer.customer_type,
      tax_id: customer.tax_id || '',
      name: customer.name,
      email: customer.email || '',
      phone: customer.phone || '',
      address: customer.address || '',
      city: customer.city || '',
      state: customer.state || '',
      postal_code: customer.postal_code || '',
      country: customer.country || 'República Dominicana',
      is_government_entity: customer.is_government_entity || false,
      default_withholding_rate: customer.default_withholding_rate || '',
    });
    setShowModal(true);
  };

  // Guardar cliente
  const handleSubmit = async (e) => {
    e.preventDefault();

    // Validar que cliente gubernamental tenga RNC
    if (formData.is_government_entity && (!formData.tax_id || formData.tax_id.trim() === '')) {
      showAlert('error', 'Error de validación', 'Los clientes gubernamentales deben tener RNC');
      return;
    }

    const payload = {
      customer_type: formData.customer_type,
      tax_id: formData.tax_id || null,
      name: formData.name,
      email: formData.email || null,
      phone: formData.phone || null,
      address: formData.address || null,
      city: formData.city || null,
      state: formData.state || null,
      postal_code: formData.postal_code || null,
      country: formData.country || null,
      is_government_entity: formData.is_government_entity || false,
      default_withholding_rate: formData.default_withholding_rate ? parseFloat(formData.default_withholding_rate) : null,
    };

    try {
      if (editingCustomer) {
        // Actualizar
        await api.put(`/customers/${editingCustomer.id}`, payload);
        showAlert('success', 'Cliente actualizado', 'El cliente se actualizó correctamente');
      } else {
        // Crear
        await api.post('/customers', payload);
        showAlert('success', 'Cliente creado', 'El cliente se creó correctamente');
      }
      
      setShowModal(false);
      loadCustomers();
    } catch (error) {
      showAlert('error', 'Error al guardar', error.response?.data?.error || 'Error desconocido');
    }
  };

  // Eliminar cliente
  const handleDelete = async (customer) => {
    setDeleteConfirm(customer);
  };

  const confirmDelete = async () => {
    try {
      await api.delete(`/customers/${deleteConfirm.id}`);
      showAlert('success', 'Cliente eliminado', 'El cliente se eliminó correctamente');
      setDeleteConfirm(null);
      loadCustomers();
    } catch (error) {
      showAlert('error', 'Error al eliminar', error.response?.data?.error || 'Error desconocido');
    }
  };

  // Cambiar estado activo/inactivo
  const toggleActive = async (customer) => {
    try {
      await api.put(`/customers/${customer.id}`, {
        is_active: !customer.is_active
      });
      showAlert('success', 'Estado actualizado', `Cliente ${!customer.is_active ? 'activado' : 'desactivado'}`);
      loadCustomers();
    } catch (error) {
      showAlert('error', 'Error al actualizar', error.response?.data?.error || 'Error desconocido');
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-[#212121] flex items-center gap-3">
            <Users className="w-8 h-8 text-[#FF6B00]" />
            Clientes
          </h1>
          <p className="text-gray-600 mt-1">Gestiona tu cartera de clientes</p>
        </div>
        <button
          onClick={handleCreate}
          className="btn-primary flex items-center gap-2"
        >
          <Plus className="w-5 h-5" />
          Nuevo Cliente
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

      {/* Filtros y búsqueda */}
      <div className="card">
        <div className="flex flex-col gap-4">
          {/* Búsqueda */}
          <div className="flex-1">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
              <input
                type="text"
                placeholder="Buscar por nombre, email o RNC/Cédula..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="input-field pl-10"
              />
            </div>
          </div>

          {/* Filtros */}
          <div className="flex flex-wrap gap-2">
            {/* Filtro por tipo */}
            <div className="flex gap-2">
              <button
                onClick={() => setFilterType('all')}
                className={`px-4 py-2 rounded-lg font-medium transition-colors text-sm ${
                  filterType === 'all'
                    ? 'bg-[#FF6B00] text-white'
                    : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                }`}
              >
                Todos
              </button>
              <button
                onClick={() => setFilterType('individual')}
                className={`px-4 py-2 rounded-lg font-medium transition-colors text-sm flex items-center gap-1 ${
                  filterType === 'individual'
                    ? 'bg-[#1877F2] text-white'
                    : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                }`}
              >
                <User className="w-4 h-4" />
                Individual
              </button>
              <button
                onClick={() => setFilterType('business')}
                className={`px-4 py-2 rounded-lg font-medium transition-colors text-sm flex items-center gap-1 ${
                  filterType === 'business'
                    ? 'bg-[#1877F2] text-white'
                    : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                }`}
              >
                <Building2 className="w-4 h-4" />
                Empresa
              </button>
            </div>

            {/* Filtro de estado */}
            <div className="flex gap-2 ml-auto">
              <button
                onClick={() => setFilterActive('active')}
                className={`px-4 py-2 rounded-lg font-medium transition-colors text-sm ${
                  filterActive === 'active'
                    ? 'bg-[#00C853] text-white'
                    : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                }`}
              >
                Activos ({customers.filter(c => c.is_active).length})
              </button>
              <button
                onClick={() => setFilterActive('inactive')}
                className={`px-4 py-2 rounded-lg font-medium transition-colors text-sm ${
                  filterActive === 'inactive'
                    ? 'bg-gray-500 text-white'
                    : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                }`}
              >
                Inactivos ({customers.filter(c => !c.is_active).length})
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* Lista de clientes */}
      {loading ? (
        <div className="card text-center py-12">
          <div className="inline-block w-8 h-8 border-4 border-[#FF6B00] border-t-transparent rounded-full animate-spin"></div>
          <p className="mt-4 text-gray-600">Cargando clientes...</p>
        </div>
      ) : filteredCustomers.length === 0 ? (
        <div className="card text-center py-12">
          <Users className="w-16 h-16 text-gray-300 mx-auto mb-4" />
          <h3 className="text-lg font-semibold text-gray-600 mb-2">
            {searchTerm ? 'No se encontraron clientes' : 'No hay clientes'}
          </h3>
          <p className="text-gray-500 mb-6">
            {searchTerm ? 'Intenta con otro término de búsqueda' : 'Comienza agregando tu primer cliente'}
          </p>
          {!searchTerm && (
            <button onClick={handleCreate} className="btn-primary">
              <Plus className="w-5 h-5 inline mr-2" />
              Crear Cliente
            </button>
          )}
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {filteredCustomers.map((customer) => (
            <div
              key={customer.id}
              className="card hover:shadow-card-hover transition-all duration-200"
            >
              {/* Header del card */}
              <div className="flex items-start justify-between mb-4">
                <div className="flex-1">
                  <div className="flex items-center gap-2 mb-2 flex-wrap">
                    {customer.customer_type === 'individual' ? (
                      <User className="w-5 h-5 text-[#FF6B00]" />
                    ) : (
                      <Building2 className="w-5 h-5 text-[#FF6B00]" />
                    )}
                    <span className="text-xs font-medium text-gray-500 capitalize">
                      {customer.customer_type === 'individual' ? 'Persona' : 'Empresa'}
                    </span>
                    {customer.is_government_entity && (
                      <span className="text-xs font-medium px-2 py-0.5 bg-blue-100 text-blue-700 rounded-full">
                        Gubernamental
                      </span>
                    )}
                  </div>
                  <h3 className="text-lg font-semibold text-[#212121]">
                    {customer.name}
                  </h3>
                  {customer.tax_id && (
                    <div className="flex items-center gap-1 mt-1">
                      <CreditCard className="w-3 h-3 text-gray-400" />
                      <span className="text-xs text-gray-600">{customer.tax_id}</span>
                    </div>
                  )}
                </div>
                <div className="flex gap-1">
                  <button
                    onClick={() => handleEdit(customer)}
                    className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
                    title="Editar"
                  >
                    <Edit2 className="w-4 h-4 text-gray-600" />
                  </button>
                  <button
                    onClick={() => handleDelete(customer)}
                    className="p-2 hover:bg-red-50 rounded-lg transition-colors"
                    title="Eliminar"
                  >
                    <Trash2 className="w-4 h-4 text-red-600" />
                  </button>
                </div>
              </div>

              {/* Info del cliente */}
              <div className="space-y-2 mb-4">
                {customer.email && (
                  <div className="flex items-center gap-2 text-sm">
                    <Mail className="w-4 h-4 text-gray-400 flex-shrink-0" />
                    <span className="text-gray-600 truncate">{customer.email}</span>
                  </div>
                )}
                {customer.phone && (
                  <div className="flex items-center gap-2 text-sm">
                    <Phone className="w-4 h-4 text-gray-400 flex-shrink-0" />
                    <span className="text-gray-600">{customer.phone}</span>
                  </div>
                )}
                {customer.address && (
                  <div className="flex items-start gap-2 text-sm">
                    <MapPin className="w-4 h-4 text-gray-400 flex-shrink-0 mt-0.5" />
                    <span className="text-gray-600 line-clamp-2">
                      {customer.address}
                      {customer.city && `, ${customer.city}`}
                    </span>
                  </div>
                )}
              </div>

              {/* Footer del card */}
              <div className="flex items-center justify-between pt-4 border-t border-gray-100">
                <button
                  onClick={() => toggleActive(customer)}
                  className={`badge ${
                    customer.is_active ? 'badge-success' : 'badge bg-gray-500 text-white'
                  } cursor-pointer hover:opacity-80 transition-opacity`}
                >
                  {customer.is_active ? 'Activo' : 'Inactivo'}
                </button>
                <span className="text-xs text-gray-400">
                  {new Date(customer.created_at).toLocaleDateString('es-DO')}
                </span>
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
                {editingCustomer ? 'Editar Cliente' : 'Nuevo Cliente'}
              </h2>
            </div>

            <form onSubmit={handleSubmit} className="p-6 space-y-4">
              {/* Tipo de cliente */}
              <div>
                <label className="block text-sm font-medium text-[#212121] mb-2">
                  Tipo de Cliente *
                </label>
                <div className="grid grid-cols-2 gap-3">
                  <button
                    type="button"
                    onClick={() => setFormData({ ...formData, customer_type: 'individual' })}
                    className={`p-3 rounded-lg border-2 transition-colors flex items-center justify-center gap-2 ${
                      formData.customer_type === 'individual'
                        ? 'border-[#FF6B00] bg-orange-50 text-[#FF6B00]'
                        : 'border-gray-200 hover:border-gray-300'
                    }`}
                  >
                    <User className="w-5 h-5" />
                    <span className="font-medium">Persona</span>
                  </button>
                  <button
                    type="button"
                    onClick={() => setFormData({ ...formData, customer_type: 'business' })}
                    className={`p-3 rounded-lg border-2 transition-colors flex items-center justify-center gap-2 ${
                      formData.customer_type === 'business'
                        ? 'border-[#FF6B00] bg-orange-50 text-[#FF6B00]'
                        : 'border-gray-200 hover:border-gray-300'
                    }`}
                  >
                    <Building2 className="w-5 h-5" />
                    <span className="font-medium">Empresa</span>
                  </button>
                </div>
              </div>

              {/* Entidad Gubernamental */}
              <div>
                <label className="flex items-center gap-2 cursor-pointer">
                  <input
                    type="checkbox"
                    checked={formData.is_government_entity}
                    onChange={(e) => setFormData({ ...formData, is_government_entity: e.target.checked })}
                    className="w-5 h-5 text-[#FF6B00] border-gray-300 rounded focus:ring-[#FF6B00]"
                  />
                  <span className="text-sm font-medium text-[#212121]">
                    Entidad Gubernamental
                  </span>
                </label>
                <p className="text-xs text-gray-500 mt-1 ml-7">
                  Marca esta opción si el cliente es una entidad del Estado (requiere RNC válido)
                </p>
              </div>

              {/* RNC/Cédula y Nombre */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-[#212121] mb-2">
                    {formData.customer_type === 'individual' ? 'Cédula' : 'RNC'} {formData.is_government_entity && '*'}
                  </label>
                  <input
                    type="text"
                    required={formData.is_government_entity}
                    value={formData.tax_id}
                    onChange={(e) => setFormData({ ...formData, tax_id: e.target.value })}
                    className="input-field"
                    placeholder={formData.customer_type === 'individual' ? '000-0000000-0' : '000-00000-0'}
                  />
                  {formData.is_government_entity && (
                    <p className="text-xs text-amber-600 mt-1">
                      RNC obligatorio para entidades gubernamentales
                    </p>
                  )}
                </div>

                <div>
                  <label className="block text-sm font-medium text-[#212121] mb-2">
                    Nombre Completo *
                  </label>
                  <input
                    type="text"
                    required
                    value={formData.name}
                    onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                    className="input-field"
                    placeholder={formData.customer_type === 'individual' ? 'Juan Pérez' : 'Empresa S.A.'}
                  />
                </div>
              </div>

              {/* Tasa de Retención por Defecto (solo si es gubernamental) */}
              {formData.is_government_entity && (
                <div>
                  <label className="block text-sm font-medium text-[#212121] mb-2">
                    Tasa de Retención por Defecto (opcional)
                  </label>
                  <div className="relative">
                    <input
                      type="number"
                      step="0.01"
                      min="0"
                      max="1"
                      value={formData.default_withholding_rate}
                      onChange={(e) => setFormData({ ...formData, default_withholding_rate: e.target.value })}
                      className="input-field pr-12"
                      placeholder="0.05"
                    />
                    <span className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 text-sm">
                      (5% = 0.05)
                    </span>
                  </div>
                  <p className="text-xs text-gray-500 mt-1">
                    Tasa de retención ISR por defecto (ej: 0.05 para 5%). Si no se especifica, se usará 5% por defecto.
                  </p>
                </div>
              )}

              {/* Email y Teléfono */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-[#212121] mb-2">
                    Email
                  </label>
                  <input
                    type="email"
                    value={formData.email}
                    onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                    className="input-field"
                    placeholder="cliente@example.com"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-[#212121] mb-2">
                    Teléfono
                  </label>
                  <input
                    type="tel"
                    value={formData.phone}
                    onChange={(e) => setFormData({ ...formData, phone: e.target.value })}
                    className="input-field"
                    placeholder="809-555-1234"
                  />
                </div>
              </div>

              {/* Dirección */}
              <div>
                <label className="block text-sm font-medium text-[#212121] mb-2">
                  Dirección
                </label>
                <input
                  type="text"
                  value={formData.address}
                  onChange={(e) => setFormData({ ...formData, address: e.target.value })}
                  className="input-field"
                  placeholder="Calle Principal #123"
                />
              </div>

              {/* Ciudad, Provincia y Código Postal */}
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div>
                  <label className="block text-sm font-medium text-[#212121] mb-2">
                    Ciudad
                  </label>
                  <input
                    type="text"
                    value={formData.city}
                    onChange={(e) => setFormData({ ...formData, city: e.target.value })}
                    className="input-field"
                    placeholder="Santo Domingo"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-[#212121] mb-2">
                    Provincia
                  </label>
                  <input
                    type="text"
                    value={formData.state}
                    onChange={(e) => setFormData({ ...formData, state: e.target.value })}
                    className="input-field"
                    placeholder="Distrito Nacional"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-[#212121] mb-2">
                    Código Postal
                  </label>
                  <input
                    type="text"
                    value={formData.postal_code}
                    onChange={(e) => setFormData({ ...formData, postal_code: e.target.value })}
                    className="input-field"
                    placeholder="10101"
                  />
                </div>
              </div>

              {/* País */}
              <div>
                <label className="block text-sm font-medium text-[#212121] mb-2">
                  País
                </label>
                <input
                  type="text"
                  value={formData.country}
                  onChange={(e) => setFormData({ ...formData, country: e.target.value })}
                  className="input-field"
                  placeholder="República Dominicana"
                />
              </div>

              {/* Botones */}
              <div className="flex gap-3 pt-4 border-t border-gray-200">
                <button
                  type="submit"
                  className="btn-primary flex-1"
                >
                  {editingCustomer ? 'Actualizar Cliente' : 'Crear Cliente'}
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
                ¿Eliminar Cliente?
              </h3>
              <p className="text-gray-600 mb-1">
                Estás a punto de eliminar a:
              </p>
              <p className="font-semibold text-[#212121] mb-4">
                {deleteConfirm.name}
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
