import { useState, useEffect } from 'react';
import { Package, Plus, Search, Edit2, Trash2, DollarSign, Tag, Box } from 'lucide-react';
import api from '../lib/api';
import Alert from '../components/Alert';

export default function Products() {
  const [products, setProducts] = useState([]);
  const [suppliers, setSuppliers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState('');
  const [filterActive, setFilterActive] = useState('all'); // all, active, inactive
  const [showModal, setShowModal] = useState(false);
  const [editingProduct, setEditingProduct] = useState(null);
  const [alert, setAlert] = useState(null);
  const [deleteConfirm, setDeleteConfirm] = useState(null);

  // Formulario
  const [formData, setFormData] = useState({
    code: '',
    name: '',
    description: '',
    product_type: 'product',
    unit: 'unidad',
    price: '',
    cost: '',
    tax_rate: '0.18', // 18% ITBIS RD por defecto
    stock_quantity: '0',
    supplier_id: '',
  });

  // Cargar productos
  useEffect(() => {
    loadProducts();
    loadSuppliers();
  }, []);

  const loadProducts = async () => {
    try {
      setLoading(true);
      const response = await api.get('/products');
      setProducts(response.data || []);
    } catch (error) {
      showAlert('error', 'Error al cargar productos', error.response?.data?.error || 'Error desconocido');
    } finally {
      setLoading(false);
    }
  };

  const loadSuppliers = async () => {
    try {
      const response = await api.get('/suppliers');
      setSuppliers(response.data || []);
    } catch (error) {
      // opcional: no bloquear si falla
    }
  };

  const showAlert = (type, title, message) => {
    setAlert({ type, title, message });
    setTimeout(() => setAlert(null), 5000);
  };

  // Filtrar productos
  const filteredProducts = products.filter(product => {
    const matchesSearch = product.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         product.code.toLowerCase().includes(searchTerm.toLowerCase());
    
    const matchesFilter = filterActive === 'all' ? true :
                         filterActive === 'active' ? product.is_active :
                         !product.is_active;

    return matchesSearch && matchesFilter;
  });

  // Abrir modal para crear
  const handleCreate = () => {
    setEditingProduct(null);
    setFormData({
      code: '',
      name: '',
      description: '',
      product_type: 'product',
      unit: 'unidad',
      price: '',
      cost: '',
      tax_rate: '0.18',
      stock_quantity: '0',
      supplier_id: '',
    });
    setShowModal(true);
  };

  // Abrir modal para editar
  const handleEdit = (product) => {
    setEditingProduct(product);
    setFormData({
      code: product.code,
      name: product.name,
      description: product.description || '',
      product_type: product.product_type,
      unit: product.unit,
      price: product.price.toString(),
      cost: product.cost ? product.cost.toString() : '',
      tax_rate: product.tax_rate.toString(),
      stock_quantity: product.stock_quantity?.toString() || '0',
      supplier_id: product.supplier_id || '',
    });
    setShowModal(true);
  };

  // Guardar producto
  const handleSubmit = async (e) => {
    e.preventDefault();

    const payload = {
      code: formData.code,
      name: formData.name,
      description: formData.description || null,
      product_type: formData.product_type,
      unit: formData.unit,
      price: parseFloat(formData.price),
      cost: formData.cost ? parseFloat(formData.cost) : null,
      tax_rate: parseFloat(formData.tax_rate),
      stock_quantity: parseFloat(formData.stock_quantity) || 0,
      supplier_id: formData.supplier_id || null,
    };

    try {
      if (editingProduct) {
        // Actualizar
        await api.put(`/products/${editingProduct.id}`, payload);
        showAlert('success', 'Producto actualizado', 'El producto se actualizó correctamente');
      } else {
        // Crear
        await api.post('/products', payload);
        showAlert('success', 'Producto creado', 'El producto se creó correctamente');
      }
      
      setShowModal(false);
      loadProducts();
    } catch (error) {
      showAlert('error', 'Error al guardar', error.response?.data?.error || 'Error desconocido');
    }
  };

  // Eliminar producto
  const handleDelete = async (product) => {
    setDeleteConfirm(product);
  };

  const confirmDelete = async () => {
    try {
      await api.delete(`/products/${deleteConfirm.id}`);
      showAlert('success', 'Producto eliminado', 'El producto se eliminó correctamente');
      setDeleteConfirm(null);
      loadProducts();
    } catch (error) {
      showAlert('error', 'Error al eliminar', error.response?.data?.error || 'Error desconocido');
    }
  };

  // Cambiar estado activo/inactivo
  const toggleActive = async (product) => {
    try {
      await api.put(`/products/${product.id}`, {
        is_active: !product.is_active
      });
      showAlert('success', 'Estado actualizado', `Producto ${!product.is_active ? 'activado' : 'desactivado'}`);
      loadProducts();
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
            <Package className="w-8 h-8 text-[#FF6B00]" />
            Productos
          </h1>
          <p className="text-gray-600 mt-1">Gestiona tu catálogo de productos y servicios</p>
        </div>
        <button
          onClick={handleCreate}
          className="btn-primary flex items-center gap-2"
        >
          <Plus className="w-5 h-5" />
          Nuevo Producto
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
        <div className="flex flex-col md:flex-row gap-4">
          {/* Búsqueda */}
          <div className="flex-1">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
              <input
                type="text"
                placeholder="Buscar por nombre o código..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="input-field pl-10"
              />
            </div>
          </div>

          {/* Filtro de estado */}
          <div className="flex gap-2">
            <button
              onClick={() => setFilterActive('all')}
              className={`px-4 py-2 rounded-lg font-medium transition-colors ${
                filterActive === 'all'
                  ? 'bg-[#FF6B00] text-white'
                  : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
              }`}
            >
              Todos ({products.length})
            </button>
            <button
              onClick={() => setFilterActive('active')}
              className={`px-4 py-2 rounded-lg font-medium transition-colors ${
                filterActive === 'active'
                  ? 'bg-[#00C853] text-white'
                  : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
              }`}
            >
              Activos ({products.filter(p => p.is_active).length})
            </button>
            <button
              onClick={() => setFilterActive('inactive')}
              className={`px-4 py-2 rounded-lg font-medium transition-colors ${
                filterActive === 'inactive'
                  ? 'bg-gray-500 text-white'
                  : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
              }`}
            >
              Inactivos ({products.filter(p => !p.is_active).length})
            </button>
          </div>
        </div>
      </div>

      {/* Lista de productos */}
      {loading ? (
        <div className="card text-center py-12">
          <div className="inline-block w-8 h-8 border-4 border-[#FF6B00] border-t-transparent rounded-full animate-spin"></div>
          <p className="mt-4 text-gray-600">Cargando productos...</p>
        </div>
      ) : filteredProducts.length === 0 ? (
        <div className="card text-center py-12">
          <Package className="w-16 h-16 text-gray-300 mx-auto mb-4" />
          <h3 className="text-lg font-semibold text-gray-600 mb-2">
            {searchTerm ? 'No se encontraron productos' : 'No hay productos'}
          </h3>
          <p className="text-gray-500 mb-6">
            {searchTerm ? 'Intenta con otro término de búsqueda' : 'Comienza agregando tu primer producto'}
          </p>
          {!searchTerm && (
            <button onClick={handleCreate} className="btn-primary">
              <Plus className="w-5 h-5 inline mr-2" />
              Crear Producto
            </button>
          )}
        </div>
      ) : (
        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
          {filteredProducts.map((product) => (
            <div
              key={product.id}
              className="card p-4 hover:shadow-card-hover transition-all duration-200"
            >
              {/* Header del card */}
              <div className="flex items-start justify-between mb-3">
                <div className="flex-1">
                  <div className="flex items-center gap-2 mb-1">
                    <Box className="w-4 h-4 text-[#FF6B00]" />
                    <span className="text-[11px] font-mono text-gray-500">{product.code}</span>
                  </div>
                  <h3 className="text-base font-semibold text-[#212121]">
                    {product.name}
                  </h3>
                  {product.description && (
                    <p className="text-xs text-gray-600 mt-1 line-clamp-2">
                      {product.description}
                    </p>
                  )}
                </div>
                <div className="flex gap-1">
                  <button
                    onClick={() => handleEdit(product)}
                    className="p-1.5 hover:bg-gray-100 rounded-lg transition-colors"
                    title="Editar"
                  >
                    <Edit2 className="w-4 h-4 text-gray-600" />
                  </button>
                  <button
                    onClick={() => handleDelete(product)}
                    className="p-1.5 hover:bg-red-50 rounded-lg transition-colors"
                    title="Eliminar"
                  >
                    <Trash2 className="w-4 h-4 text-red-600" />
                  </button>
                </div>
              </div>

              {/* Info del producto */}
              <div className="space-y-1.5 mb-3">
                <div className="flex items-center justify-between text-sm">
                  <span className="text-gray-600">Tipo:</span>
                  <span className="font-medium capitalize">{product.product_type}</span>
                </div>
                <div className="flex items-center justify-between text-sm">
                  <span className="text-gray-600">Unidad:</span>
                  <span className="font-medium capitalize">{product.unit}</span>
                </div>
                {product.product_type === 'product' && (
                  <div className="flex items-center justify-between text-sm">
                    <span className="text-gray-600">Stock:</span>
                    <span className={`font-bold ${
                      product.stock_quantity <= 0 ? 'text-red-600' :
                      product.stock_quantity <= 10 ? 'text-yellow-600' :
                      'text-[#00C853]'
                    }`}>
                      {product.stock_quantity.toFixed(2)}
                    </span>
                  </div>
                )}
                <div className="flex items-center justify-between">
                  <span className="text-gray-600 text-sm">Precio:</span>
                  <span className="text-lg font-bold text-[#FF6B00]">
                    ${product.price.toFixed(2)}
                  </span>
                </div>
                {product.cost && (
                  <div className="flex items-center justify-between text-sm">
                    <span className="text-gray-600">Costo:</span>
                    <span className="font-medium">${product.cost.toFixed(2)}</span>
                  </div>
                )}
                <div className="flex items-center justify-between text-sm">
                  <span className="text-gray-600">Impuesto:</span>
                  <span className="font-medium">{(product.tax_rate * 100).toFixed(0)}%</span>
                </div>
              </div>

              {/* Footer del card */}
              <div className="flex items-center justify-between pt-3 border-t border-gray-100">
                <button
                  onClick={() => toggleActive(product)}
                  className={`badge ${
                    product.is_active ? 'badge-success' : 'badge bg-gray-500 text-white'
                  } cursor-pointer hover:opacity-80 transition-opacity`}
                >
                  {product.is_active ? 'Activo' : 'Inactivo'}
                </button>
                <span className="text-xs text-gray-400">
                  {new Date(product.created_at).toLocaleDateString('es-DO')}
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
                {editingProduct ? 'Editar Producto' : 'Nuevo Producto'}
              </h2>
            </div>

            <form onSubmit={handleSubmit} className="p-6 space-y-4">
              {/* Código y Nombre */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-[#212121] mb-2">
                    Código *
                  </label>
                  <input
                    type="text"
                    required
                    value={formData.code}
                    onChange={(e) => setFormData({ ...formData, code: e.target.value })}
                    className="input-field"
                    placeholder="PROD-001"
                    disabled={!!editingProduct}
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-[#212121] mb-2">
                    Nombre *
                  </label>
                  <input
                    type="text"
                    required
                    value={formData.name}
                    onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                    className="input-field"
                    placeholder="Laptop Dell XPS 15"
                  />
                </div>
              </div>

              {/* Descripción */}
              <div>
                <label className="block text-sm font-medium text-[#212121] mb-2">
                  Descripción
                </label>
                <textarea
                  value={formData.description}
                  onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                  className="input-field"
                  rows="3"
                  placeholder="Descripción detallada del producto..."
                />
              </div>

              {/* Tipo, Unidad y Proveedor */}
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div>
                  <label className="block text-sm font-medium text-[#212121] mb-2">
                    Tipo *
                  </label>
                  <select
                    required
                    value={formData.product_type}
                    onChange={(e) => setFormData({ ...formData, product_type: e.target.value })}
                    className="input-field"
                  >
                    <option value="product">Producto</option>
                    <option value="service">Servicio</option>
                  </select>
                </div>

                <div>
                  <label className="block text-sm font-medium text-[#212121] mb-2">
                    Unidad *
                  </label>
                  <select
                    required
                    value={formData.unit}
                    onChange={(e) => setFormData({ ...formData, unit: e.target.value })}
                    className="input-field"
                  >
                    <option value="unidad">Unidad</option>
                    <option value="caja">Caja</option>
                    <option value="paquete">Paquete</option>
                    <option value="metro">Metro</option>
                    <option value="kilogramo">Kilogramo</option>
                    <option value="litro">Litro</option>
                    <option value="hora">Hora</option>
                    <option value="dia">Día</option>
                    <option value="mes">Mes</option>
                  </select>
                </div>

                <div>
                  <label className="block text-sm font-medium text-[#212121] mb-2">
                    Proveedor
                  </label>
                  <select
                    value={formData.supplier_id}
                    onChange={(e) => setFormData({ ...formData, supplier_id: e.target.value })}
                    className="input-field"
                  >
                    <option value="">Sin proveedor</option>
                    {suppliers.map(s => (
                      <option key={s.id} value={s.id}>{s.name}</option>
                    ))}
                  </select>
                </div>
              </div>

              {/* Precio, Costo y Tax */}
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div>
                  <label className="block text-sm font-medium text-[#212121] mb-2">
                    Precio * (RD$)
                  </label>
                  <input
                    type="number"
                    required
                    step="0.01"
                    min="0"
                    value={formData.price}
                    onChange={(e) => setFormData({ ...formData, price: e.target.value })}
                    className="input-field"
                    placeholder="1500.00"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-[#212121] mb-2">
                    Costo (RD$)
                  </label>
                  <input
                    type="number"
                    step="0.01"
                    min="0"
                    value={formData.cost}
                    onChange={(e) => setFormData({ ...formData, cost: e.target.value })}
                    className="input-field"
                    placeholder="1000.00"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-[#212121] mb-2">
                    ITBIS * (%)
                  </label>
                  <select
                    required
                    value={formData.tax_rate}
                    onChange={(e) => setFormData({ ...formData, tax_rate: e.target.value })}
                    className="input-field"
                  >
                    <option value="0">0% (Exento)</option>
                    <option value="0.16">16%</option>
                    <option value="0.18">18% (Estándar)</option>
                  </select>
                </div>
              </div>

              {/* Stock Quantity (solo para productos, no para servicios) */}
              {formData.product_type === 'product' && (
                <div>
                  <label className="block text-sm font-medium text-[#212121] mb-2">
                    Cantidad en Stock *
                  </label>
                  <input
                    type="number"
                    required
                    step="0.01"
                    min="0"
                    value={formData.stock_quantity}
                    onChange={(e) => setFormData({ ...formData, stock_quantity: e.target.value })}
                    className="input-field"
                    placeholder="100.00"
                  />
                  <p className="text-xs text-gray-500 mt-1">
                    La cantidad en stock se reducirá automáticamente al realizar ventas
                  </p>
                </div>
              )}

              {/* Botones */}
              <div className="flex gap-3 pt-4 border-t border-gray-200">
                <button
                  type="submit"
                  className="btn-primary flex-1"
                >
                  {editingProduct ? 'Actualizar Producto' : 'Crear Producto'}
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
                ¿Eliminar Producto?
              </h3>
              <p className="text-gray-600 mb-1">
                Estás a punto de eliminar:
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
