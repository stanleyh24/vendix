import { useEffect, useState } from 'react'
import { Users, Plus, Search, Edit2, Trash2 } from 'lucide-react'
import Alert from '../components/Alert'
import api from '../lib/api'

export default function Suppliers() {
  const [suppliers, setSuppliers] = useState([])
  const [loading, setLoading] = useState(false)
  const [search, setSearch] = useState('')
  const [showModal, setShowModal] = useState(false)
  const [editing, setEditing] = useState(null)
  const [alert, setAlert] = useState(null)
  const [confirmDelete, setConfirmDelete] = useState(null)

  const [form, setForm] = useState({
    name: '',
    tax_id: '',
    email: '',
    phone: '',
    address: '',
    city: '',
    state: '',
    postal_code: '',
    country: 'DO',
    is_active: true,
  })

  useEffect(() => { load() }, [])

  const showAlert = (type, title, message) => {
    setAlert({ type, title, message })
    setTimeout(() => setAlert(null), 5000)
  }

  const load = async () => {
    setLoading(true)
    try {
      const { data } = await api.get('/suppliers')
      setSuppliers(data || [])
    } catch (e) {
      showAlert('error', 'Error', 'No se pudieron cargar los proveedores')
    } finally {
      setLoading(false)
    }
  }

  const filtered = suppliers.filter(s =>
    s.name.toLowerCase().includes(search.toLowerCase()) ||
    (s.tax_id || '').toLowerCase().includes(search.toLowerCase()) ||
    (s.email || '').toLowerCase().includes(search.toLowerCase())
  )

  const openCreate = () => {
    setEditing(null)
    setForm({ name: '', tax_id: '', email: '', phone: '', address: '', city: '', state: '', postal_code: '', country: 'DO', is_active: true })
    setShowModal(true)
  }

  const openEdit = (supplier) => {
    setEditing(supplier)
    setForm({
      name: supplier.name || '',
      tax_id: supplier.tax_id || '',
      email: supplier.email || '',
      phone: supplier.phone || '',
      address: supplier.address || '',
      city: supplier.city || '',
      state: supplier.state || '',
      postal_code: supplier.postal_code || '',
      country: supplier.country || 'DO',
      is_active: supplier.is_active,
    })
    setShowModal(true)
  }

  const save = async (e) => {
    e.preventDefault()
    try {
      const payload = {
        name: form.name,
        tax_id: form.tax_id || null,
        email: form.email || null,
        phone: form.phone || null,
        address: form.address || null,
        city: form.city || null,
        state: form.state || null,
        postal_code: form.postal_code || null,
        country: form.country || 'DO',
        is_active: form.is_active,
      }
      if (editing) {
        await api.put(`/suppliers/${editing.id}`, payload)
        showAlert('success', 'Actualizado', 'Proveedor actualizado correctamente')
      } else {
        await api.post('/suppliers', payload)
        showAlert('success', 'Creado', 'Proveedor creado correctamente')
      }
      setShowModal(false)
      load()
    } catch (e) {
      showAlert('error', 'Error', e.response?.data?.error || 'No se pudo guardar')
    }
  }

  const doDelete = async () => {
    try {
      await api.delete(`/suppliers/${confirmDelete.id}`)
      showAlert('success', 'Eliminado', 'Proveedor eliminado')
      setConfirmDelete(null)
      load()
    } catch (e) {
      showAlert('error', 'Error', e.response?.data?.error || 'No se pudo eliminar')
      setConfirmDelete(null)
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-[#212121] flex items-center gap-3">
            <Users className="w-8 h-8 text-[#FF6B00]" />
            Proveedores
          </h1>
          <p className="text-gray-600 mt-1">Gestiona tus proveedores</p>
        </div>
        <button onClick={openCreate} className="btn-primary flex items-center gap-2">
          <Plus className="w-5 h-5" />
          Nuevo Proveedor
        </button>
      </div>

      {alert && (
        <Alert type={alert.type} title={alert.title} message={alert.message} onClose={() => setAlert(null)} />
      )}

      <div className="card">
        <div className="flex items-center gap-4">
          <div className="flex-1">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
              <input
                type="text"
                placeholder="Buscar por nombre, RNC o email..."
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                className="input-field pl-10"
              />
            </div>
          </div>
        </div>
      </div>

      {loading ? (
        <div className="card text-center py-12">
          <div className="inline-block w-8 h-8 border-4 border-[#FF6B00] border-t-transparent rounded-full animate-spin"></div>
          <p className="mt-4 text-gray-600">Cargando proveedores...</p>
        </div>
      ) : (
        <div className="card p-0 overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Nombre</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">RNC</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Email</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Teléfono</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Ciudad</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Estado</th>
                <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Acciones</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {filtered.map(s => (
                <tr key={s.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-[#212121]">{s.name}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-700">{s.tax_id || '-'}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-700">{s.email || '-'}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-700">{s.phone || '-'}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-700">{s.city || '-'}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm">
                    <span className={`badge ${s.is_active ? 'badge-success' : 'badge bg-gray-500 text-white'}`}>
                      {s.is_active ? 'Activo' : 'Inactivo'}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-right text-sm">
                    <button onClick={() => openEdit(s)} className="p-2 hover:bg-gray-100 rounded-lg mr-1" title="Editar">
                      <Edit2 className="w-4 h-4 text-gray-600" />
                    </button>
                    <button onClick={() => setConfirmDelete(s)} className="p-2 hover:bg-red-50 rounded-lg" title="Eliminar">
                      <Trash2 className="w-4 h-4 text-red-600" />
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {showModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-2xl shadow-card-hover max-w-2xl w-full max-h-[90vh] overflow-y-auto">
            <div className="sticky top-0 bg-white border-b border-gray-200 px-6 py-4 rounded-t-2xl">
              <h2 className="text-2xl font-bold text-[#212121]">{editing ? 'Editar Proveedor' : 'Nuevo Proveedor'}</h2>
            </div>
            <form onSubmit={save} className="p-6 space-y-4">
              <div>
                <label className="block text-sm font-medium text-[#212121] mb-2">Nombre *</label>
                <input required value={form.name} onChange={e => setForm({ ...form, name: e.target.value })} className="input-field" />
              </div>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-[#212121] mb-2">RNC</label>
                  <input value={form.tax_id} onChange={e => setForm({ ...form, tax_id: e.target.value })} className="input-field" />
                </div>
                <div>
                  <label className="block text-sm font-medium text-[#212121] mb-2">Email</label>
                  <input type="email" value={form.email} onChange={e => setForm({ ...form, email: e.target.value })} className="input-field" />
                </div>
              </div>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-[#212121] mb-2">Teléfono</label>
                  <input value={form.phone} onChange={e => setForm({ ...form, phone: e.target.value })} className="input-field" />
                </div>
                <div>
                  <label className="block text-sm font-medium text-[#212121] mb-2">País</label>
                  <input value={form.country} onChange={e => setForm({ ...form, country: e.target.value })} className="input-field" />
                </div>
              </div>
              <div>
                <label className="block text-sm font-medium text-[#212121] mb-2">Dirección</label>
                <input value={form.address} onChange={e => setForm({ ...form, address: e.target.value })} className="input-field" />
              </div>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div>
                  <label className="block text-sm font-medium text-[#212121] mb-2">Ciudad</label>
                  <input value={form.city} onChange={e => setForm({ ...form, city: e.target.value })} className="input-field" />
                </div>
                <div>
                  <label className="block text-sm font-medium text-[#212121] mb-2">Provincia</label>
                  <input value={form.state} onChange={e => setForm({ ...form, state: e.target.value })} className="input-field" />
                </div>
                <div>
                  <label className="block text-sm font-medium text-[#212121] mb-2">Código Postal</label>
                  <input value={form.postal_code} onChange={e => setForm({ ...form, postal_code: e.target.value })} className="input-field" />
                </div>
              </div>
              <div className="flex gap-3 pt-4 border-t border-gray-200">
                <button type="submit" className="btn-primary flex-1">{editing ? 'Actualizar' : 'Crear'}</button>
                <button type="button" onClick={() => setShowModal(false)} className="btn-outline flex-1">Cancelar</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {confirmDelete && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-2xl shadow-card-hover max-w-md w-full p-6">
            <h3 className="text-xl font-bold text-[#212121] mb-2">Eliminar proveedor</h3>
            <p className="text-gray-600 mb-6">Esta acción no se puede deshacer.</p>
            <div className="flex gap-3">
              <button onClick={() => setConfirmDelete(null)} className="btn-outline flex-1">Cancelar</button>
              <button onClick={doDelete} className="btn-error flex-1">Eliminar</button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
