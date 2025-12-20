import { useState, useEffect } from 'react';
import { DollarSign, Calendar, Plus, Edit2, Trash2, Eye, Users, FileText, CheckCircle } from 'lucide-react';
import api from '../lib/api';
import Alert from '../components/Alert';

export default function Payroll() {
  const [activeTab, setActiveTab] = useState('periods');
  const [alert, setAlert] = useState(null);
  const [loading, setLoading] = useState(false);

  // Periods state
  const [periods, setPeriods] = useState([]);
  const [showPeriodModal, setShowPeriodModal] = useState(false);
  const [editingPeriod, setEditingPeriod] = useState(null);
  const [periodFormData, setPeriodFormData] = useState({
    period_code: '',
    period_start: '',
    period_end: '',
    notes: '',
  });

  // Entries state
  const [selectedPeriod, setSelectedPeriod] = useState(null);
  const [entries, setEntries] = useState([]);
  const [employees, setEmployees] = useState([]);
  const [showEntryModal, setShowEntryModal] = useState(false);
  const [editingEntry, setEditingEntry] = useState(null);
  const [entryFormData, setEntryFormData] = useState({
    employee_id: '',
    hours_worked: '',
    days_worked: '',
    overtime_hours: '',
    bonuses: '',
    commissions: '',
    other_deductions: '',
    notes: '',
  });

  const tabs = [
    { id: 'periods', label: 'Períodos', icon: Calendar },
    { id: 'entries', label: 'Entradas', icon: FileText },
  ];

  useEffect(() => {
    if (activeTab === 'periods') {
      fetchPeriods();
    } else if (activeTab === 'entries') {
      fetchEmployees();
      if (selectedPeriod) {
        fetchEntries(selectedPeriod.id);
      }
    }
  }, [activeTab, selectedPeriod]);

  const fetchPeriods = async () => {
    try {
      setLoading(true);
      const response = await api.get('/payroll/periods');
      setPeriods(response.data || []);
    } catch (error) {
      showAlert('error', 'Error al cargar períodos', error.response?.data?.error || 'Error desconocido');
    } finally {
      setLoading(false);
    }
  };

  const fetchEntries = async (periodId) => {
    try {
      setLoading(true);
      const response = await api.get(`/payroll/periods/${periodId}/entries`);
      setEntries(response.data || []);
    } catch (error) {
      showAlert('error', 'Error al cargar entradas', error.response?.data?.error || 'Error desconocido');
    } finally {
      setLoading(false);
    }
  };

  const fetchEmployees = async () => {
    try {
      const response = await api.get('/employees/active');
      setEmployees(response.data || []);
    } catch (error) {
      showAlert('error', 'Error al cargar empleados', error.response?.data?.error || 'Error desconocido');
    }
  };

  const showAlert = (type, title, message) => {
    setAlert({ type, title, message });
    setTimeout(() => setAlert(null), 5000);
  };

  const handleCreatePeriod = () => {
    setEditingPeriod(null);
    const now = new Date();
    const year = now.getFullYear();
    const month = String(now.getMonth() + 1).padStart(2, '0');
    const periodCode = `${year}-${month}`;
    
    const firstDay = new Date(year, now.getMonth(), 1);
    const lastDay = new Date(year, now.getMonth() + 1, 0);
    
    setPeriodFormData({
      period_code: periodCode,
      period_start: firstDay.toISOString().split('T')[0],
      period_end: lastDay.toISOString().split('T')[0],
      notes: '',
    });
    setShowPeriodModal(true);
  };

  const handleSubmitPeriod = async (e) => {
    e.preventDefault();
    try {
      const payload = {
        period_code: periodFormData.period_code,
        period_start: periodFormData.period_start,
        period_end: periodFormData.period_end,
        notes: periodFormData.notes || null,
      };

      if (editingPeriod) {
        await api.put(`/payroll/periods/${editingPeriod.id}`, payload);
        showAlert('success', 'Período actualizado', 'El período se actualizó correctamente');
      } else {
        await api.post('/payroll/periods', payload);
        showAlert('success', 'Período creado', 'El período se creó correctamente');
      }
      
      setShowPeriodModal(false);
      fetchPeriods();
    } catch (error) {
      showAlert('error', 'Error al guardar', error.response?.data?.error || 'Error desconocido');
    }
  };

  const handleCreateEntry = () => {
    if (!selectedPeriod) {
      showAlert('error', 'Selecciona un período', 'Primero debes seleccionar un período de nómina');
      return;
    }
    setEditingEntry(null);
    setEntryFormData({
      employee_id: '',
      hours_worked: '',
      days_worked: '',
      overtime_hours: '',
      bonuses: '',
      commissions: '',
      other_deductions: '',
      notes: '',
    });
    setShowEntryModal(true);
  };

  const handleSubmitEntry = async (e) => {
    e.preventDefault();
    try {
      const payload = {
        employee_id: entryFormData.employee_id,
        hours_worked: entryFormData.hours_worked ? parseFloat(entryFormData.hours_worked) : null,
        days_worked: entryFormData.days_worked ? parseFloat(entryFormData.days_worked) : null,
        overtime_hours: entryFormData.overtime_hours ? parseFloat(entryFormData.overtime_hours) : null,
        bonuses: entryFormData.bonuses ? parseFloat(entryFormData.bonuses) : null,
        commissions: entryFormData.commissions ? parseFloat(entryFormData.commissions) : null,
        other_deductions: entryFormData.other_deductions ? parseFloat(entryFormData.other_deductions) : null,
        notes: entryFormData.notes || null,
      };

      if (editingEntry) {
        await api.put(`/payroll/entries/${editingEntry.id}`, payload);
        showAlert('success', 'Entrada actualizada', 'La entrada se actualizó correctamente');
      } else {
        await api.post(`/payroll/periods/${selectedPeriod.id}/entries`, payload);
        showAlert('success', 'Entrada creada', 'La entrada se creó correctamente');
      }
      
      setShowEntryModal(false);
      if (selectedPeriod) {
        fetchEntries(selectedPeriod.id);
      }
    } catch (error) {
      showAlert('error', 'Error al guardar', error.response?.data?.error || 'Error desconocido');
    }
  };

  const handleDeletePeriod = async (period) => {
    if (!confirm(`¿Estás seguro de eliminar el período ${period.period_code}?`)) {
      return;
    }
    try {
      await api.delete(`/payroll/periods/${period.id}`);
      showAlert('success', 'Período eliminado', 'El período se eliminó correctamente');
      fetchPeriods();
    } catch (error) {
      showAlert('error', 'Error al eliminar', error.response?.data?.error || 'Error desconocido');
    }
  };

  const handleDeleteEntry = async (entry) => {
    if (!confirm(`¿Estás seguro de eliminar la entrada de ${entry.employee_name}?`)) {
      return;
    }
    try {
      await api.delete(`/payroll/entries/${entry.id}`);
      showAlert('success', 'Entrada eliminada', 'La entrada se eliminó correctamente');
      if (selectedPeriod) {
        fetchEntries(selectedPeriod.id);
      }
    } catch (error) {
      showAlert('error', 'Error al eliminar', error.response?.data?.error || 'Error desconocido');
    }
  };

  const formatCurrency = (amount) => {
    return new Intl.NumberFormat('es-DO', {
      style: 'currency',
      currency: 'DOP',
      minimumFractionDigits: 2,
    }).format(amount);
  };

  const getStatusColor = (status) => {
    const colors = {
      draft: 'bg-gray-500',
      processing: 'bg-blue-500',
      completed: 'bg-green-500',
      paid: 'bg-purple-500',
    };
    return colors[status] || 'bg-gray-500';
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-[#212121] flex items-center gap-3">
            <DollarSign className="w-8 h-8 text-[#FF6B00]" />
            Nómina
          </h1>
          <p className="text-gray-600 mt-1">Gestiona períodos y pagos de nómina</p>
        </div>
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
      <div className="card">
        <div className="flex gap-2 border-b border-gray-200">
          {tabs.map((tab) => {
            const Icon = tab.icon;
            return (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id)}
                className={`flex items-center gap-2 px-6 py-3 font-medium transition-colors border-b-2 ${
                  activeTab === tab.id
                    ? 'border-[#FF6B00] text-[#FF6B00]'
                    : 'border-transparent text-gray-600 hover:text-[#212121]'
                }`}
              >
                <Icon className="w-5 h-5" />
                {tab.label}
              </button>
            );
          })}
        </div>

        {/* Periods Tab */}
        {activeTab === 'periods' && (
          <div className="p-6">
            <div className="flex justify-between items-center mb-6">
              <h2 className="text-xl font-semibold text-[#212121]">Períodos de Nómina</h2>
              <button
                onClick={handleCreatePeriod}
                className="btn-primary flex items-center gap-2"
              >
                <Plus className="w-5 h-5" />
                Nuevo Período
              </button>
            </div>

            {loading ? (
              <div className="text-center py-12">
                <div className="inline-block w-8 h-8 border-4 border-[#FF6B00] border-t-transparent rounded-full animate-spin"></div>
                <p className="mt-4 text-gray-600">Cargando períodos...</p>
              </div>
            ) : periods.length === 0 ? (
              <div className="text-center py-12">
                <Calendar className="w-16 h-16 text-gray-300 mx-auto mb-4" />
                <h3 className="text-lg font-semibold text-gray-600 mb-2">No hay períodos</h3>
                <p className="text-gray-500 mb-6">Comienza creando tu primer período de nómina</p>
                <button onClick={handleCreatePeriod} className="btn-primary">
                  <Plus className="w-5 h-5 inline mr-2" />
                  Crear Período
                </button>
              </div>
            ) : (
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead>
                    <tr className="border-b border-gray-200">
                      <th className="text-left py-3 px-4 font-semibold text-[#212121]">Período</th>
                      <th className="text-left py-3 px-4 font-semibold text-[#212121]">Fecha Inicio</th>
                      <th className="text-left py-3 px-4 font-semibold text-[#212121]">Fecha Fin</th>
                      <th className="text-left py-3 px-4 font-semibold text-[#212121]">Estado</th>
                      <th className="text-right py-3 px-4 font-semibold text-[#212121]">Total Bruto</th>
                      <th className="text-right py-3 px-4 font-semibold text-[#212121]">Total Neto</th>
                      <th className="text-right py-3 px-4 font-semibold text-[#212121]">Acciones</th>
                    </tr>
                  </thead>
                  <tbody>
                    {periods.map((period) => (
                      <tr key={period.id} className="border-b border-gray-100 hover:bg-gray-50">
                        <td className="py-3 px-4 font-medium">{period.period_code}</td>
                        <td className="py-3 px-4">{new Date(period.period_start).toLocaleDateString('es-DO')}</td>
                        <td className="py-3 px-4">{new Date(period.period_end).toLocaleDateString('es-DO')}</td>
                        <td className="py-3 px-4">
                          <span className={`badge ${getStatusColor(period.status)} text-white capitalize`}>
                            {period.status}
                          </span>
                        </td>
                        <td className="py-3 px-4 text-right">{formatCurrency(period.total_gross)}</td>
                        <td className="py-3 px-4 text-right font-semibold">{formatCurrency(period.total_net)}</td>
                        <td className="py-3 px-4">
                          <div className="flex gap-2 justify-end">
                            <button
                              onClick={() => {
                                setSelectedPeriod(period);
                                setActiveTab('entries');
                              }}
                              className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
                              title="Ver entradas"
                            >
                              <Eye className="w-4 h-4 text-gray-600" />
                            </button>
                            <button
                              onClick={() => handleDeletePeriod(period)}
                              className="p-2 hover:bg-red-50 rounded-lg transition-colors"
                              title="Eliminar"
                            >
                              <Trash2 className="w-4 h-4 text-red-600" />
                            </button>
                          </div>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        )}

        {/* Entries Tab */}
        {activeTab === 'entries' && (
          <div className="p-6">
            <div className="flex justify-between items-center mb-6">
              <div>
                <h2 className="text-xl font-semibold text-[#212121]">Entradas de Nómina</h2>
                {selectedPeriod && (
                  <p className="text-sm text-gray-600 mt-1">
                    Período: {selectedPeriod.period_code} - {new Date(selectedPeriod.period_start).toLocaleDateString('es-DO')} a {new Date(selectedPeriod.period_end).toLocaleDateString('es-DO')}
                  </p>
                )}
              </div>
              {selectedPeriod && (
                <button
                  onClick={handleCreateEntry}
                  className="btn-primary flex items-center gap-2"
                >
                  <Plus className="w-5 h-5" />
                  Nueva Entrada
                </button>
              )}
            </div>

            {!selectedPeriod ? (
              <div className="text-center py-12">
                <FileText className="w-16 h-16 text-gray-300 mx-auto mb-4" />
                <h3 className="text-lg font-semibold text-gray-600 mb-2">Selecciona un período</h3>
                <p className="text-gray-500 mb-6">Ve a la pestaña de Períodos y selecciona uno para ver sus entradas</p>
                <button onClick={() => setActiveTab('periods')} className="btn-primary">
                  Ver Períodos
                </button>
              </div>
            ) : loading ? (
              <div className="text-center py-12">
                <div className="inline-block w-8 h-8 border-4 border-[#FF6B00] border-t-transparent rounded-full animate-spin"></div>
                <p className="mt-4 text-gray-600">Cargando entradas...</p>
              </div>
            ) : entries.length === 0 ? (
              <div className="text-center py-12">
                <Users className="w-16 h-16 text-gray-300 mx-auto mb-4" />
                <h3 className="text-lg font-semibold text-gray-600 mb-2">No hay entradas</h3>
                <p className="text-gray-500 mb-6">Comienza agregando entradas para este período</p>
                <button onClick={handleCreateEntry} className="btn-primary">
                  <Plus className="w-5 h-5 inline mr-2" />
                  Crear Entrada
                </button>
              </div>
            ) : (
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead>
                    <tr className="border-b border-gray-200">
                      <th className="text-left py-3 px-4 font-semibold text-[#212121]">Empleado</th>
                      <th className="text-left py-3 px-4 font-semibold text-[#212121]">Cargo</th>
                      <th className="text-right py-3 px-4 font-semibold text-[#212121]">Salario Bruto</th>
                      <th className="text-right py-3 px-4 font-semibold text-[#212121]">Extras</th>
                      <th className="text-right py-3 px-4 font-semibold text-[#212121]">Total Bruto</th>
                      <th className="text-right py-3 px-4 font-semibold text-[#212121]">Deducciones</th>
                      <th className="text-right py-3 px-4 font-semibold text-[#212121]">Neto</th>
                      <th className="text-right py-3 px-4 font-semibold text-[#212121]">Acciones</th>
                    </tr>
                  </thead>
                  <tbody>
                    {entries.map((entry) => (
                      <tr key={entry.id} className="border-b border-gray-100 hover:bg-gray-50">
                        <td className="py-3 px-4">
                          <div>
                            <div className="font-medium">{entry.employee_name}</div>
                            <div className="text-sm text-gray-500">{entry.employee_code}</div>
                          </div>
                        </td>
                        <td className="py-3 px-4">{entry.position || '-'}</td>
                        <td className="py-3 px-4 text-right">{formatCurrency(entry.gross_salary)}</td>
                        <td className="py-3 px-4 text-right">
                          {formatCurrency(entry.overtime_pay + entry.bonuses + entry.commissions)}
                        </td>
                        <td className="py-3 px-4 text-right font-medium">{formatCurrency(entry.total_gross)}</td>
                        <td className="py-3 px-4 text-right text-red-600">{formatCurrency(entry.total_deductions)}</td>
                        <td className="py-3 px-4 text-right font-semibold text-green-600">{formatCurrency(entry.net_salary)}</td>
                        <td className="py-3 px-4">
                          <div className="flex gap-2 justify-end">
                            <button
                              onClick={() => {
                                setEditingEntry(entry);
                                setEntryFormData({
                                  employee_id: entry.employee_id,
                                  hours_worked: entry.hours_worked?.toString() || '',
                                  days_worked: entry.days_worked?.toString() || '',
                                  overtime_hours: entry.overtime_hours?.toString() || '',
                                  bonuses: entry.bonuses?.toString() || '',
                                  commissions: entry.commissions?.toString() || '',
                                  other_deductions: entry.other_deductions?.toString() || '',
                                  notes: entry.notes || '',
                                });
                                setShowEntryModal(true);
                              }}
                              className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
                              title="Editar"
                            >
                              <Edit2 className="w-4 h-4 text-gray-600" />
                            </button>
                            <button
                              onClick={() => handleDeleteEntry(entry)}
                              className="p-2 hover:bg-red-50 rounded-lg transition-colors"
                              title="Eliminar"
                            >
                              <Trash2 className="w-4 h-4 text-red-600" />
                            </button>
                          </div>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        )}
      </div>

      {/* Period Modal */}
      {showPeriodModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-2xl shadow-card-hover max-w-md w-full">
            <div className="sticky top-0 bg-white border-b border-gray-200 px-6 py-4 rounded-t-2xl">
              <h2 className="text-2xl font-bold text-[#212121]">
                {editingPeriod ? 'Editar Período' : 'Nuevo Período'}
              </h2>
            </div>
            <form onSubmit={handleSubmitPeriod} className="p-6 space-y-4">
              <div>
                <label className="block text-sm font-medium text-[#212121] mb-2">
                  Código de Período *
                </label>
                <input
                  type="text"
                  required
                  value={periodFormData.period_code}
                  onChange={(e) => setPeriodFormData({ ...periodFormData, period_code: e.target.value })}
                  className="input-field"
                  placeholder="2024-01"
                />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-[#212121] mb-2">
                    Fecha Inicio *
                  </label>
                  <input
                    type="date"
                    required
                    value={periodFormData.period_start}
                    onChange={(e) => setPeriodFormData({ ...periodFormData, period_start: e.target.value })}
                    className="input-field"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-[#212121] mb-2">
                    Fecha Fin *
                  </label>
                  <input
                    type="date"
                    required
                    value={periodFormData.period_end}
                    onChange={(e) => setPeriodFormData({ ...periodFormData, period_end: e.target.value })}
                    className="input-field"
                  />
                </div>
              </div>
              <div>
                <label className="block text-sm font-medium text-[#212121] mb-2">
                  Notas
                </label>
                <textarea
                  value={periodFormData.notes}
                  onChange={(e) => setPeriodFormData({ ...periodFormData, notes: e.target.value })}
                  className="input-field"
                  rows="3"
                />
              </div>
              <div className="flex gap-3 pt-4 border-t border-gray-200">
                <button type="submit" className="btn-primary flex-1">
                  {editingPeriod ? 'Actualizar' : 'Crear'}
                </button>
                <button
                  type="button"
                  onClick={() => setShowPeriodModal(false)}
                  className="btn-outline flex-1"
                >
                  Cancelar
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Entry Modal */}
      {showEntryModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-2xl shadow-card-hover max-w-2xl w-full max-h-[90vh] overflow-y-auto">
            <div className="sticky top-0 bg-white border-b border-gray-200 px-6 py-4 rounded-t-2xl">
              <h2 className="text-2xl font-bold text-[#212121]">
                {editingEntry ? 'Editar Entrada' : 'Nueva Entrada'}
              </h2>
            </div>
            <form onSubmit={handleSubmitEntry} className="p-6 space-y-4">
              <div>
                <label className="block text-sm font-medium text-[#212121] mb-2">
                  Empleado *
                </label>
                <select
                  required
                  value={entryFormData.employee_id}
                  onChange={(e) => setEntryFormData({ ...entryFormData, employee_id: e.target.value })}
                  className="input-field"
                  disabled={!!editingEntry}
                >
                  <option value="">Selecciona un empleado</option>
                  {employees.map((emp) => (
                    <option key={emp.id} value={emp.id}>
                      {emp.first_name} {emp.last_name} ({emp.employee_code})
                    </option>
                  ))}
                </select>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-[#212121] mb-2">
                    Horas Trabajadas
                  </label>
                  <input
                    type="number"
                    step="0.01"
                    value={entryFormData.hours_worked}
                    onChange={(e) => setEntryFormData({ ...entryFormData, hours_worked: e.target.value })}
                    className="input-field"
                    placeholder="0"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-[#212121] mb-2">
                    Días Trabajados
                  </label>
                  <input
                    type="number"
                    step="0.01"
                    value={entryFormData.days_worked}
                    onChange={(e) => setEntryFormData({ ...entryFormData, days_worked: e.target.value })}
                    className="input-field"
                    placeholder="0"
                  />
                </div>
              </div>
              <div>
                <label className="block text-sm font-medium text-[#212121] mb-2">
                  Horas Extras
                </label>
                <input
                  type="number"
                  step="0.01"
                  value={entryFormData.overtime_hours}
                  onChange={(e) => setEntryFormData({ ...entryFormData, overtime_hours: e.target.value })}
                  className="input-field"
                  placeholder="0"
                />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-[#212121] mb-2">
                    Bonificaciones
                  </label>
                  <input
                    type="number"
                    step="0.01"
                    value={entryFormData.bonuses}
                    onChange={(e) => setEntryFormData({ ...entryFormData, bonuses: e.target.value })}
                    className="input-field"
                    placeholder="0.00"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-[#212121] mb-2">
                    Comisiones
                  </label>
                  <input
                    type="number"
                    step="0.01"
                    value={entryFormData.commissions}
                    onChange={(e) => setEntryFormData({ ...entryFormData, commissions: e.target.value })}
                    className="input-field"
                    placeholder="0.00"
                  />
                </div>
              </div>
              <div>
                <label className="block text-sm font-medium text-[#212121] mb-2">
                  Otras Deducciones
                </label>
                <input
                  type="number"
                  step="0.01"
                  value={entryFormData.other_deductions}
                  onChange={(e) => setEntryFormData({ ...entryFormData, other_deductions: e.target.value })}
                  className="input-field"
                  placeholder="0.00"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-[#212121] mb-2">
                  Notas
                </label>
                <textarea
                  value={entryFormData.notes}
                  onChange={(e) => setEntryFormData({ ...entryFormData, notes: e.target.value })}
                  className="input-field"
                  rows="3"
                />
              </div>
              <div className="flex gap-3 pt-4 border-t border-gray-200">
                <button type="submit" className="btn-primary flex-1">
                  {editingEntry ? 'Actualizar' : 'Crear'}
                </button>
                <button
                  type="button"
                  onClick={() => setShowEntryModal(false)}
                  className="btn-outline flex-1"
                >
                  Cancelar
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}

