import { useState, useEffect } from 'react';
import { UserPlus, Plus, Search, Edit2, Trash2, Mail, Phone, MapPin, Briefcase, DollarSign, Calendar, Building2 } from 'lucide-react';
import api from '../lib/api';
import Alert from '../components/Alert';

export default function Employees() {
  const [employees, setEmployees] = useState([]);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState('');
  const [filterActive, setFilterActive] = useState('all');
  const [showModal, setShowModal] = useState(false);
  const [editingEmployee, setEditingEmployee] = useState(null);
  const [alert, setAlert] = useState(null);
  const [deleteConfirm, setDeleteConfirm] = useState(null);

  const [formData, setFormData] = useState({
    employee_code: '',
    first_name: '',
    last_name: '',
    email: '',
    phone: '',
    tax_id: '',
    address: '',
    city: '',
    state: '',
    postal_code: '',
    country: 'DO',
    department: '',
    position: '',
    hire_date: '',
    salary: '',
    salary_type: 'monthly',
  });

  useEffect(() => {
    loadEmployees();
  }, []);

  const loadEmployees = async () => {
    try {
      setLoading(true);
      const response = await api.get('/employees');
      setEmployees(response.data || []);
    } catch (error) {
      showAlert('error', 'Error al cargar empleados', error.response?.data?.error || 'Error desconocido');
    } finally {
      setLoading(false);
    }
  };

  const showAlert = (type, title, message) => {
    setAlert({ type, title, message });
    setTimeout(() => setAlert(null), 5000);
  };

  const filteredEmployees = employees.filter(emp => {
    const matchesSearch = 
      `${emp.first_name} ${emp.last_name}`.toLowerCase().includes(searchTerm.toLowerCase()) ||
      (emp.employee_code && emp.employee_code.toLowerCase().includes(searchTerm.toLowerCase())) ||
      (emp.email && emp.email.toLowerCase().includes(searchTerm.toLowerCase()));
    
    const matchesFilter = filterActive === 'all' ? true :
                         filterActive === 'active' ? emp.is_active :
                         !emp.is_active;

    return matchesSearch && matchesFilter;
  });

  const handleCreate = () => {
    setEditingEmployee(null);
    setFormData({
      employee_code: '',
      first_name: '',
      last_name: '',
      email: '',
      phone: '',
      tax_id: '',
      address: '',
      city: '',
      state: '',
      postal_code: '',
      country: 'DO',
      department: '',
      position: '',
      hire_date: '',
      salary: '',
      salary_type: 'monthly',
    });
    setShowModal(true);
  };

  const handleEdit = (employee) => {
    setEditingEmployee(employee);
    setFormData({
      employee_code: employee.employee_code,
      first_name: employee.first_name,
      last_name: employee.last_name,
      email: employee.email || '',
      phone: employee.phone || '',
      tax_id: employee.tax_id || '',
      address: employee.address || '',
      city: employee.city || '',
      state: employee.state || '',
      postal_code: employee.postal_code || '',
      country: employee.country || 'DO',
      department: employee.department || '',
      position: employee.position || '',
      hire_date: employee.hire_date ? employee.hire_date.split('T')[0] : '',
      salary: employee.salary ? employee.salary.toString() : '',
      salary_type: employee.salary_type || 'monthly',
    });
    setShowModal(true);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    const payload = {
      employee_code: formData.employee_code,
      first_name: formData.first_name,
      last_name: formData.last_name,
      email: formData.email || null,
      phone: formData.phone || null,
      tax_id: formData.tax_id || null,
      address: formData.address || null,
      city: formData.city || null,
      state: formData.state || null,
      postal_code: formData.postal_code || null,
      country: formData.country || null,
      department: formData.department || null,
      position: formData.position || null,
      hire_date: formData.hire_date || null,
      salary: formData.salary ? parseFloat(formData.salary) : null,
      salary_type: formData.salary_type,
    };

    try {
      if (editingEmployee) {
        await api.put(`/employees/${editingEmployee.id}`, payload);
        showAlert('success', 'Empleado actualizado', 'El empleado se actualizó correctamente');
      } else {
        await api.post('/employees', payload);
        showAlert('success', 'Empleado creado', 'El empleado se creó correctamente');
      }
      
      setShowModal(false);
      loadEmployees();
    } catch (error) {
      showAlert('error', 'Error al guardar', error.response?.data?.error || 'Error desconocido');
    }
  };

  const handleDelete = async (employee) => {
    setDeleteConfirm(employee);
  };

  const confirmDelete = async () => {
    try {
      await api.delete(`/employees/${deleteConfirm.id}`);
      showAlert('success', 'Empleado eliminado', 'El empleado se eliminó correctamente');
      setDeleteConfirm(null);
      loadEmployees();
    } catch (error) {
      showAlert('error', 'Error al eliminar', error.response?.data?.error || 'Error desconocido');
    }
  };

  const toggleActive = async (employee) => {
    try {
      await api.put(`/employees/${employee.id}`, {
        is_active: !employee.is_active
      });
      showAlert('success', 'Estado actualizado', `Empleado ${!employee.is_active ? 'activado' : 'desactivado'}`);
      loadEmployees();
    } catch (error) {
      showAlert('error', 'Error al actualizar', error.response?.data?.error || 'Error desconocido');
    }
  };

  const formatSalary = (salary, salaryType) => {
    if (!salary) return 'No asignado';
    const formatted = new Intl.NumberFormat('es-DO', {
      style: 'currency',
      currency: 'DOP',
      minimumFractionDigits: 0,
    }).format(salary);
    
    const typeLabels = {
      monthly: 'mes',
      hourly: 'hora',
      daily: 'día'
    };
    
    return `${formatted}/${typeLabels[salaryType] || salaryType}`;
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-[#212121] flex items-center gap-3">
            <UserPlus className="w-8 h-8 text-[#FF6B00]" />
            Empleados
          </h1>
          <p className="text-gray-600 mt-1">Gestiona tu equipo de trabajo</p>
        </div>
        <button
          onClick={handleCreate}
          className="btn-primary flex items-center gap-2"
        >
          <Plus className="w-5 h-5" />
          Nuevo Empleado
        </button>
      </div>

      {alert && (
        <Alert
          type={alert.type}
          title={alert.title}
          message={alert.message}
          onClose={() => setAlert(null)}
        />
      )}

      <div className="card">
        <div className="flex flex-col gap-4">
          <div className="flex-1">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
              <input
                type="text"
                placeholder="Buscar por nombre, código o email..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="input-field pl-10"
              />
            </div>
          </div>

          <div className="flex gap-2">
            <button
              onClick={() => setFilterActive('all')}
              className={`px-4 py-2 rounded-lg font-medium transition-colors text-sm ${
                filterActive === 'all'
                  ? 'bg-[#FF6B00] text-white'
                  : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
              }`}
            >
              Todos
            </button>
            <button
              onClick={() => setFilterActive('active')}
              className={`px-4 py-2 rounded-lg font-medium transition-colors text-sm ${
                filterActive === 'active'
                  ? 'bg-[#00C853] text-white'
                  : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
              }`}
            >
              Activos ({employees.filter(e => e.is_active).length})
            </button>
            <button
              onClick={() => setFilterActive('inactive')}
              className={`px-4 py-2 rounded-lg font-medium transition-colors text-sm ${
                filterActive === 'inactive'
                  ? 'bg-gray-500 text-white'
                  : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
              }`}
            >
              Inactivos ({employees.filter(e => !e.is_active).length})
            </button>
          </div>
        </div>
      </div>

      {loading ? (
        <div className="card text-center py-12">
          <div className="inline-block w-8 h-8 border-4 border-[#FF6B00] border-t-transparent rounded-full animate-spin"></div>
          <p className="mt-4 text-gray-600">Cargando empleados...</p>
        </div>
      ) : filteredEmployees.length === 0 ? (
        <div className="card text-center py-12">
          <UserPlus className="w-16 h-16 text-gray-300 mx-auto mb-4" />
          <h3 className="text-lg font-semibold text-gray-600 mb-2">
            {searchTerm ? 'No se encontraron empleados' : 'No hay empleados'}
          </h3>
          <p className="text-gray-500 mb-6">
            {searchTerm ? 'Intenta con otro término de búsqueda' : 'Comienza agregando tu primer empleado'}
          </p>
          {!searchTerm && (
            <button onClick={handleCreate} className="btn-primary">
              <Plus className="w-5 h-5 inline mr-2" />
              Crear Empleado
            </button>
          )}
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {filteredEmployees.map((employee) => (
            <div
              key={employee.id}
              className="card hover:shadow-card-hover transition-all duration-200"
            >
              <div className="flex items-start justify-between mb-4">
                <div className="flex-1">
                  <div className="flex items-center gap-2 mb-2">
                    <Briefcase className="w-5 h-5 text-[#FF6B00]" />
                    <span className="text-xs font-medium text-gray-500">
                      {employee.employee_code}
                    </span>
                  </div>
                  <h3 className="text-lg font-semibold text-[#212121]">
                    {employee.first_name} {employee.last_name}
                  </h3>
                  {employee.position && (
                    <p className="text-sm text-gray-600 mt-1">{employee.position}</p>
                  )}
                </div>
                <div className="flex gap-1">
                  <button
                    onClick={() => handleEdit(employee)}
                    className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
                    title="Editar"
                  >
                    <Edit2 className="w-4 h-4 text-gray-600" />
                  </button>
                  <button
                    onClick={() => handleDelete(employee)}
                    className="p-2 hover:bg-red-50 rounded-lg transition-colors"
                    title="Eliminar"
                  >
                    <Trash2 className="w-4 h-4 text-red-600" />
                  </button>
                </div>
              </div>

              <div className="space-y-2 mb-4">
                {employee.email && (
                  <div className="flex items-center gap-2 text-sm">
                    <Mail className="w-4 h-4 text-gray-400 flex-shrink-0" />
                    <span className="text-gray-600 truncate">{employee.email}</span>
                  </div>
                )}
                {employee.phone && (
                  <div className="flex items-center gap-2 text-sm">
                    <Phone className="w-4 h-4 text-gray-400 flex-shrink-0" />
                    <span className="text-gray-600">{employee.phone}</span>
                  </div>
                )}
                {employee.department && (
                  <div className="flex items-center gap-2 text-sm">
                    <Building2 className="w-4 h-4 text-gray-400 flex-shrink-0" />
                    <span className="text-gray-600">{employee.department}</span>
                  </div>
                )}
                {employee.salary && (
                  <div className="flex items-center gap-2 text-sm">
                    <DollarSign className="w-4 h-4 text-gray-400 flex-shrink-0" />
                    <span className="text-gray-600">{formatSalary(employee.salary, employee.salary_type)}</span>
                  </div>
                )}
                {employee.hire_date && (
                  <div className="flex items-center gap-2 text-sm">
                    <Calendar className="w-4 h-4 text-gray-400 flex-shrink-0" />
                    <span className="text-gray-600">
                      Ingreso: {new Date(employee.hire_date).toLocaleDateString('es-DO')}
                    </span>
                  </div>
                )}
              </div>

              <div className="flex items-center justify-between pt-4 border-t border-gray-100">
                <button
                  onClick={() => toggleActive(employee)}
                  className={`badge ${
                    employee.is_active ? 'badge-success' : 'badge bg-gray-500 text-white'
                  } cursor-pointer hover:opacity-80 transition-opacity`}
                >
                  {employee.is_active ? 'Activo' : 'Inactivo'}
                </button>
                <span className="text-xs text-gray-400">
                  {new Date(employee.created_at).toLocaleDateString('es-DO')}
                </span>
              </div>
            </div>
          ))}
        </div>
      )}

      {showModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-2xl shadow-card-hover max-w-3xl w-full max-h-[90vh] overflow-y-auto">
            <div className="sticky top-0 bg-white border-b border-gray-200 px-6 py-4 rounded-t-2xl">
              <h2 className="text-2xl font-bold text-[#212121]">
                {editingEmployee ? 'Editar Empleado' : 'Nuevo Empleado'}
              </h2>
            </div>

            <form onSubmit={handleSubmit} className="p-6 space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-[#212121] mb-2">
                    Código de Empleado *
                  </label>
                  <input
                    type="text"
                    required
                    value={formData.employee_code}
                    onChange={(e) => setFormData({ ...formData, employee_code: e.target.value })}
                    className="input-field"
                    placeholder="EMP001"
                    disabled={!!editingEmployee}
                  />
                </div>
                <div className="grid grid-cols-2 gap-2">
                  <div>
                    <label className="block text-sm font-medium text-[#212121] mb-2">
                      Nombre *
                    </label>
                    <input
                      type="text"
                      required
                      value={formData.first_name}
                      onChange={(e) => setFormData({ ...formData, first_name: e.target.value })}
                      className="input-field"
                      placeholder="Juan"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-[#212121] mb-2">
                      Apellido *
                    </label>
                    <input
                      type="text"
                      required
                      value={formData.last_name}
                      onChange={(e) => setFormData({ ...formData, last_name: e.target.value })}
                      className="input-field"
                      placeholder="Pérez"
                    />
                  </div>
                </div>
              </div>

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
                    placeholder="empleado@example.com"
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

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-[#212121] mb-2">
                    Cédula/RNC
                  </label>
                  <input
                    type="text"
                    value={formData.tax_id}
                    onChange={(e) => setFormData({ ...formData, tax_id: e.target.value })}
                    className="input-field"
                    placeholder="000-0000000-0"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-[#212121] mb-2">
                    Fecha de Ingreso
                  </label>
                  <input
                    type="date"
                    value={formData.hire_date}
                    onChange={(e) => setFormData({ ...formData, hire_date: e.target.value })}
                    className="input-field"
                  />
                </div>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-[#212121] mb-2">
                    Departamento
                  </label>
                  <input
                    type="text"
                    value={formData.department}
                    onChange={(e) => setFormData({ ...formData, department: e.target.value })}
                    className="input-field"
                    placeholder="Ventas"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-[#212121] mb-2">
                    Cargo
                  </label>
                  <input
                    type="text"
                    value={formData.position}
                    onChange={(e) => setFormData({ ...formData, position: e.target.value })}
                    className="input-field"
                    placeholder="Gerente"
                  />
                </div>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div>
                  <label className="block text-sm font-medium text-[#212121] mb-2">
                    Salario
                  </label>
                  <input
                    type="number"
                    step="0.01"
                    value={formData.salary}
                    onChange={(e) => setFormData({ ...formData, salary: e.target.value })}
                    className="input-field"
                    placeholder="50000"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-[#212121] mb-2">
                    Tipo de Salario *
                  </label>
                  <select
                    required
                    value={formData.salary_type}
                    onChange={(e) => setFormData({ ...formData, salary_type: e.target.value })}
                    className="input-field"
                  >
                    <option value="monthly">Mensual</option>
                    <option value="hourly">Por Hora</option>
                    <option value="daily">Por Día</option>
                  </select>
                </div>
              </div>

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

              <div className="flex gap-3 pt-4 border-t border-gray-200">
                <button
                  type="submit"
                  className="btn-primary flex-1"
                >
                  {editingEmployee ? 'Actualizar Empleado' : 'Crear Empleado'}
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

      {deleteConfirm && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-2xl shadow-card-hover max-w-md w-full p-6">
            <div className="text-center">
              <div className="w-16 h-16 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-4">
                <Trash2 className="w-8 h-8 text-red-600" />
              </div>
              <h3 className="text-xl font-bold text-[#212121] mb-2">
                ¿Eliminar Empleado?
              </h3>
              <p className="text-gray-600 mb-1">
                Estás a punto de eliminar a:
              </p>
              <p className="font-semibold text-[#212121] mb-4">
                {deleteConfirm.first_name} {deleteConfirm.last_name}
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

