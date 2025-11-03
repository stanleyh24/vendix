import { useState, useEffect } from 'react';
import { ShoppingCart, Plus, Trash2, Search, User, X, Calculator, DollarSign, Printer, FileText } from 'lucide-react';
import api from '../lib/api';
import Alert from '../components/Alert';
import InvoicePrint from '../components/InvoicePrint';

// Cliente genérico constante
const GENERIC_CUSTOMER = {
  id: '00000000-0000-0000-0000-000000000001',
  name: 'Cliente Genérico',
  tax_id: '000000000',
  customer_type: 'individual',
  is_active: true,
  is_generic: true
};

export default function Pos() {
  const [customers, setCustomers] = useState([]);
  const [products, setProducts] = useState([]);
  const [selectedCustomer, setSelectedCustomer] = useState(null);
  const [ncfType, setNcfType] = useState('02'); // 01=Crédito Fiscal, 02=Consumidor Final (default)
  const [paymentMethod, setPaymentMethod] = useState('cash'); // cash, card, transfer, mixed
  const [cart, setCart] = useState([]);
  const [searchCustomer, setSearchCustomer] = useState('');
  const [searchProduct, setSearchProduct] = useState('');
  const [showCustomerModal, setShowCustomerModal] = useState(false);
  const [showCheckoutModal, setShowCheckoutModal] = useState(false); // Modal de completar venta
  const [showSuccessModal, setShowSuccessModal] = useState(false); // Modal de éxito después de la venta
  const [createdInvoiceId, setCreatedInvoiceId] = useState(null); // ID de la factura creada
  const [showInvoicePrint, setShowInvoicePrint] = useState(false); // Modal de impresión
  const [alert, setAlert] = useState(null);
  const [loading, setLoading] = useState(false);

  // Cargar datos iniciales
  useEffect(() => {
    loadCustomers();
    loadProducts();
  }, []);

  const loadCustomers = async () => {
    try {
      const response = await api.get('/customers');
      setCustomers(response.data || []);
    } catch (error) {
      console.error('Error loading customers:', error);
    }
  };

  const loadProducts = async () => {
    try {
      const response = await api.get('/products');
      setProducts(response.data?.filter(p => p.is_active) || []);
    } catch (error) {
      console.error('Error loading products:', error);
    }
  };

  const showAlert = (type, title, message) => {
    setAlert({ type, title, message });
    setTimeout(() => setAlert(null), 5000);
  };

  // Filtrar clientes y productos
  const filteredCustomers = customers.filter(c => 
    c.name.toLowerCase().includes(searchCustomer.toLowerCase()) ||
    (c.email && c.email.toLowerCase().includes(searchCustomer.toLowerCase()))
  );

  const filteredProducts = products.filter(p =>
    p.name.toLowerCase().includes(searchProduct.toLowerCase()) ||
    p.code.toLowerCase().includes(searchProduct.toLowerCase())
  );

  // Seleccionar cliente
  const selectCustomer = (customer) => {
    setSelectedCustomer(customer);
    setShowCustomerModal(false);
    setSearchCustomer('');
    
    // Si se selecciona cliente genérico, forzar NCF 02 (Consumidor Final)
    if (customer.is_generic) {
      setNcfType('02');
    }
  };

  // Agregar producto al carrito
  const addToCart = (product) => {
    const existingItem = cart.find(item => item.product.id === product.id);
    
    if (existingItem) {
      setCart(cart.map(item =>
        item.product.id === product.id
          ? { ...item, quantity: item.quantity + 1 }
          : item
      ));
    } else {
      setCart([...cart, { product, quantity: 1 }]);
    }
    
    setSearchProduct('');
  };

  // Actualizar cantidad
  const updateQuantity = (productId, newQuantity) => {
    if (newQuantity <= 0) {
      removeFromCart(productId);
      return;
    }
    
    setCart(cart.map(item =>
      item.product.id === productId
        ? { ...item, quantity: newQuantity }
        : item
    ));
  };

  // Remover del carrito
  const removeFromCart = (productId) => {
    setCart(cart.filter(item => item.product.id !== productId));
  };

  // Calcular totales
  const calculateSubtotal = () => {
    return cart.reduce((sum, item) => sum + (item.product.price * item.quantity), 0);
  };

  const calculateTax = () => {
    return cart.reduce((sum, item) => {
      const itemSubtotal = item.product.price * item.quantity;
      return sum + (itemSubtotal * item.product.tax_rate);
    }, 0);
  };

  const calculateTotal = () => {
    return calculateSubtotal() + calculateTax();
  };

  // Abrir modal de checkout
  const openCheckout = () => {
    if (cart.length === 0) {
      showAlert('warning', 'Carrito vacío', 'Agrega al menos un producto para procesar la venta');
      return;
    }
    setShowCheckoutModal(true);
  };

  // Procesar venta
  const processSale = async () => {
    setLoading(true);
    try {
      const saleData = {
        customer_id: selectedCustomer ? selectedCustomer.id : null,
        payment_type: paymentMethod,
        notes: `Venta procesada desde POS - ${paymentMethod}`,
        lines: cart.map(item => ({
          product_id: item.product.id,
          description: item.product.name,
          quantity: item.quantity,
          unit_price: item.product.price,
          tax_rate: item.product.tax_rate,
        })),
      };

      const response = await api.post('/sales', saleData);
      setCreatedInvoiceId(response.data.invoice_id);
      setShowCheckoutModal(false);
      setShowSuccessModal(true);
    } catch (error) {
      showAlert('error', 'Error al procesar venta', error.response?.data?.error || 'Error desconocido');
    } finally {
      setLoading(false);
    }
  };

  // Crear nueva venta (limpiar todo)
  const startNewSale = () => {
    setCart([]);
    setSelectedCustomer(null);
    setNcfType('02');
    setPaymentMethod('cash');
    setShowSuccessModal(false);
    setCreatedInvoiceId(null);
    showAlert('success', 'Listo', 'Puedes comenzar una nueva venta');
  };

  const printInvoice = () => {
    setShowSuccessModal(false);
    setShowInvoicePrint(true);
  };

  const clearSale = () => {
    setCart([]);
    setSelectedCustomer(null);
    setNcfType('02');
    setPaymentMethod('cash');
  };

  return (
    <div className="flex flex-col" style={{ height: 'calc(100vh - 64px)' }}>
      {/* Header */}
      <div className="bg-white border-b border-gray-200 px-6 py-4 flex-shrink-0">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-[#212121] flex items-center gap-3">
              <ShoppingCart className="w-8 h-8 text-[#FF6B00]" />
              Punto de Venta
            </h1>
            <p className="text-gray-600 mt-1">Punto de venta rápido</p>
          </div>
          {selectedCustomer && (
            <div className="bg-[#F5F5F5] rounded-lg px-4 py-2">
              <p className="text-xs text-gray-600">Cliente:</p>
              <p className="font-semibold text-[#212121]">{selectedCustomer.name}</p>
            </div>
          )}
        </div>
      </div>

      {/* Alert */}
      {alert && (
        <div className="px-6 pt-4 flex-shrink-0">
          <Alert
            type={alert.type}
            title={alert.title}
            message={alert.message}
            onClose={() => setAlert(null)}
          />
        </div>
      )}

      {/* Main Content */}
      <div className="flex-1 flex overflow-hidden">
        {/* Panel de productos */}
        <div className="flex-1 overflow-y-auto p-6 bg-gray-50">
          <div className="mb-4">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
              <input
                type="text"
                placeholder="Buscar productos..."
                value={searchProduct}
                onChange={(e) => setSearchProduct(e.target.value)}
                className="input-field pl-10"
              />
            </div>
          </div>

          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
            {filteredProducts.map((product) => (
              <button
                key={product.id}
                onClick={() => addToCart(product)}
                className="bg-white rounded-lg border border-gray-200 p-4 hover:shadow-md hover:border-[#FF6B00] transition-all duration-200 text-left"
              >
                <h3 className="font-semibold text-[#212121] mb-1 truncate">
                  {product.name}
                </h3>
                <p className="text-xs text-gray-500 mb-2">{product.code}</p>
                <div className="flex items-center justify-between">
                  <span className="text-lg font-bold text-[#FF6B00]">
                    ${product.price.toFixed(2)}
                  </span>
                  <Plus className="w-5 h-5 text-gray-400" />
                </div>
              </button>
            ))}
          </div>

          {filteredProducts.length === 0 && (
            <div className="text-center py-12">
              <p className="text-gray-500">No se encontraron productos</p>
            </div>
          )}
        </div>

        {/* Panel del carrito */}
        <div className="w-96 bg-white border-l border-gray-200 flex flex-col h-full overflow-hidden">

          {/* Items del carrito */}
          <div className="flex-1 overflow-y-auto p-4 min-h-0">
            {cart.length === 0 ? (
              <div className="text-center py-12">
                <ShoppingCart className="w-16 h-16 text-gray-300 mx-auto mb-4" />
                <p className="text-gray-500 text-sm">Carrito vacío</p>
                <p className="text-gray-400 text-xs mt-1">Agrega productos para comenzar</p>
              </div>
            ) : (
              <div className="space-y-2">
                {cart.map((item) => (
                  <div
                    key={item.product.id}
                    className="bg-[#F5F5F5] rounded-lg p-3"
                  >
                    <div className="flex items-start justify-between mb-2">
                      <div className="flex-1">
                        <h4 className="font-semibold text-sm text-[#212121]">
                          {item.product.name}
                        </h4>
                        <p className="text-xs text-gray-600">{item.product.code}</p>
                      </div>
                      <button
                        onClick={() => removeFromCart(item.product.id)}
                        className="p-1 hover:bg-red-100 rounded text-red-600"
                      >
                        <Trash2 className="w-4 h-4" />
                      </button>
                    </div>

                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-2">
                        <button
                          onClick={() => updateQuantity(item.product.id, item.quantity - 1)}
                          className="w-7 h-7 bg-white rounded border border-gray-300 hover:bg-gray-100 flex items-center justify-center"
                        >
                          -
                        </button>
                        <span className="w-8 text-center font-semibold">
                          {item.quantity}
                        </span>
                        <button
                          onClick={() => updateQuantity(item.product.id, item.quantity + 1)}
                          className="w-7 h-7 bg-white rounded border border-gray-300 hover:bg-gray-100 flex items-center justify-center"
                        >
                          +
                        </button>
                      </div>
                      <div className="text-right">
                        <p className="text-xs text-gray-600">
                          ${item.product.price.toFixed(2)} c/u
                        </p>
                        <p className="font-bold text-[#FF6B00]">
                          ${(item.product.price * item.quantity).toFixed(2)}
                        </p>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>

          {/* Resumen de totales */}
          <div className="border-t border-gray-200 p-4 space-y-3 flex-shrink-0 bg-white">
            <div className="flex items-center justify-between text-sm">
              <span className="text-gray-600">Subtotal:</span>
              <span className="font-semibold text-[#212121]">
                ${calculateSubtotal().toFixed(2)}
              </span>
            </div>
            <div className="flex items-center justify-between text-sm">
              <span className="text-gray-600">ITBIS:</span>
              <span className="font-semibold text-[#212121]">
                ${calculateTax().toFixed(2)}
              </span>
            </div>
            <div className="flex items-center justify-between pt-3 border-t border-gray-200">
              <span className="font-bold text-lg text-[#212121]">Total:</span>
              <span className="font-bold text-2xl text-[#FF6B00]">
                ${calculateTotal().toFixed(2)}
              </span>
            </div>
          </div>

          {/* Botones de acción */}
          <div className="p-4 space-y-2 border-t border-gray-200 flex-shrink-0 bg-white">
            <button
              onClick={openCheckout}
              disabled={cart.length === 0}
              className="btn-primary w-full flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <Calculator className="w-5 h-5" />
              Completar Venta
            </button>
            
            <button
              onClick={clearSale}
              disabled={loading}
              className="btn-outline w-full disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Limpiar Todo
            </button>
          </div>
        </div>
      </div>

      {/* Modal de Completar Venta (Checkout) */}
      {showCheckoutModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-2xl shadow-card-hover max-w-3xl w-full max-h-[90vh] overflow-y-auto">
            <div className="sticky top-0 bg-white border-b border-gray-200 px-6 py-4 rounded-t-2xl flex items-center justify-between">
              <h2 className="text-2xl font-bold text-[#212121]">Completar Venta</h2>
              <button
                onClick={() => setShowCheckoutModal(false)}
                className="p-2 hover:bg-gray-100 rounded-lg"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            <div className="p-6 space-y-6">
              {/* Sección 1: Cliente */}
              <div>
                <h3 className="text-lg font-semibold text-[#212121] mb-3 flex items-center gap-2">
                  <User className="w-5 h-5 text-[#FF6B00]" />
                  Cliente
                </h3>
                
                {selectedCustomer ? (
                  <div className="flex items-center justify-between bg-[#F5F5F5] rounded-lg p-3">
                    <div className="flex items-center gap-2">
                      <User className="w-5 h-5 text-[#FF6B00]" />
                      <div>
                        <p className="font-semibold text-sm text-[#212121]">
                          {selectedCustomer.name}
                        </p>
                        {selectedCustomer.tax_id && !selectedCustomer.is_generic && (
                          <p className="text-xs text-gray-600">{selectedCustomer.tax_id}</p>
                        )}
                        {selectedCustomer.is_generic && (
                          <p className="text-xs text-gray-500">Consumidor Final</p>
                        )}
                      </div>
                    </div>
                    <button
                      onClick={() => setSelectedCustomer(null)}
                      className="p-1 hover:bg-gray-200 rounded"
                    >
                      <X className="w-4 h-4 text-gray-600" />
                    </button>
                  </div>
                ) : (
                  <div className="grid grid-cols-2 gap-3">
                    <button
                      onClick={() => selectCustomer(GENERIC_CUSTOMER)}
                      className="bg-[#FF6B00] text-white py-3 px-4 rounded-lg hover:bg-[#E55D00] transition-colors flex items-center justify-center gap-2"
                    >
                      <User className="w-5 h-5" />
                      Cliente Genérico
                    </button>
                    <button
                      onClick={() => setShowCustomerModal(true)}
                      className="btn-outline flex items-center justify-center gap-2"
                    >
                      <Search className="w-5 h-5" />
                      Buscar Cliente
                    </button>
                  </div>
                )}
              </div>

              {/* Sección 2: Tipo de Comprobante Fiscal */}
              <div>
                <h3 className="text-lg font-semibold text-[#212121] mb-3">
                  Tipo de Comprobante Fiscal
                </h3>
                <div className="grid grid-cols-2 gap-3">
                  <button
                    onClick={() => setNcfType('02')}
                    className={`py-4 px-4 rounded-lg font-medium transition-all ${
                      ncfType === '02'
                        ? 'bg-[#FF6B00] text-white shadow-md ring-2 ring-[#FF6B00] ring-offset-2'
                        : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                    }`}
                  >
                    <div className="text-sm font-bold mb-1">NCF 02</div>
                    <div className="text-xs opacity-90">Consumidor Final</div>
                    <div className="text-[10px] mt-1 opacity-75">Para personas sin RNC</div>
                  </button>
                  <button
                    onClick={() => setNcfType('01')}
                    disabled={selectedCustomer?.is_generic}
                    title={selectedCustomer?.is_generic ? 'Requiere cliente con RNC' : ''}
                    className={`py-4 px-4 rounded-lg font-medium transition-all ${
                      ncfType === '01'
                        ? 'bg-[#00C853] text-white shadow-md ring-2 ring-[#00C853] ring-offset-2'
                        : selectedCustomer?.is_generic
                        ? 'bg-gray-200 text-gray-400 cursor-not-allowed opacity-60'
                        : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                    }`}
                  >
                    <div className="text-sm font-bold mb-1">NCF 01</div>
                    <div className="text-xs opacity-90">Crédito Fiscal</div>
                    <div className="text-[10px] mt-1 opacity-75">Para empresas con RNC</div>
                  </button>
                </div>
                {selectedCustomer?.is_generic ? (
                  <p className="text-xs text-yellow-600 mt-2 flex items-center gap-1">
                    ⚠️ Crédito Fiscal requiere cliente con RNC válido
                  </p>
                ) : (
                  <p className="text-xs text-gray-500 mt-2">
                    {ncfType === '01' ? 'Permite deducir ITBIS' : 'Para ventas generales'}
                  </p>
                )}
              </div>

              {/* Sección 3: Método de Pago */}
              <div>
                <h3 className="text-lg font-semibold text-[#212121] mb-3 flex items-center gap-2">
                  <DollarSign className="w-5 h-5 text-[#FF6B00]" />
                  Método de Pago
                </h3>
                <div className="grid grid-cols-2 gap-3">
                  <button
                    onClick={() => setPaymentMethod('cash')}
                    className={`py-3 px-4 rounded-lg font-medium transition-all ${
                      paymentMethod === 'cash'
                        ? 'bg-[#FF6B00] text-white shadow-md'
                        : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                    }`}
                  >
                    💵 Efectivo
                  </button>
                  <button
                    onClick={() => setPaymentMethod('card')}
                    className={`py-3 px-4 rounded-lg font-medium transition-all ${
                      paymentMethod === 'card'
                        ? 'bg-[#FF6B00] text-white shadow-md'
                        : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                    }`}
                  >
                    💳 Tarjeta
                  </button>
                  <button
                    onClick={() => setPaymentMethod('transfer')}
                    className={`py-3 px-4 rounded-lg font-medium transition-all ${
                      paymentMethod === 'transfer'
                        ? 'bg-[#FF6B00] text-white shadow-md'
                        : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                    }`}
                  >
                    🏦 Transferencia
                  </button>
                  <button
                    onClick={() => setPaymentMethod('mixed')}
                    className={`py-3 px-4 rounded-lg font-medium transition-all ${
                      paymentMethod === 'mixed'
                        ? 'bg-[#FF6B00] text-white shadow-md'
                        : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                    }`}
                  >
                    💰 Mixto
                  </button>
                </div>
              </div>

              {/* Resumen de la Venta */}
              <div className="bg-[#F5F5F5] rounded-lg p-4">
                <h3 className="text-lg font-semibold text-[#212121] mb-3">Resumen</h3>
                <div className="space-y-2">
                  <div className="flex justify-between text-sm">
                    <span className="text-gray-600">Productos:</span>
                    <span className="font-medium">{cart.length} item(s)</span>
                  </div>
                  <div className="flex justify-between text-sm">
                    <span className="text-gray-600">Subtotal:</span>
                    <span className="font-semibold">RD$ {calculateSubtotal().toFixed(2)}</span>
                  </div>
                  <div className="flex justify-between text-sm">
                    <span className="text-gray-600">ITBIS:</span>
                    <span className="font-semibold">RD$ {calculateTax().toFixed(2)}</span>
                  </div>
                  <div className="border-t border-gray-300 pt-2 flex justify-between">
                    <span className="font-bold text-lg">Total:</span>
                    <span className="font-bold text-2xl text-[#FF6B00]">
                      RD$ {calculateTotal().toFixed(2)}
                    </span>
                  </div>
                </div>
              </div>

              {/* Botones de acción */}
              <div className="flex gap-3 pt-4">
                <button
                  onClick={() => setShowCheckoutModal(false)}
                  className="btn-outline flex-1"
                  disabled={loading}
                >
                  Cancelar
                </button>
                <button
                  onClick={processSale}
                  disabled={loading || !selectedCustomer}
                  className="btn-primary flex-1 flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {loading ? (
                    <>
                      <div className="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
                      Procesando...
                    </>
                  ) : !selectedCustomer ? (
                    <>
                      <Calculator className="w-5 h-5" />
                      Selecciona un Cliente
                    </>
                  ) : (
                    <>
                      <Calculator className="w-5 h-5" />
                      Confirmar Venta
                    </>
                  )}
                </button>
              </div>
              
              {/* Mensaje de ayuda si no hay cliente */}
              {!selectedCustomer && !loading && (
                <div className="text-center mt-3">
                  <p className="text-xs text-yellow-600 flex items-center justify-center gap-1">
                    ⚠️ Debes seleccionar un cliente para continuar
                  </p>
                </div>
              )}
            </div>
          </div>
        </div>
      )}

      {/* Modal de selección de cliente */}
      {showCustomerModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-2xl shadow-card-hover max-w-2xl w-full max-h-[80vh] overflow-y-auto">
            <div className="sticky top-0 bg-white border-b border-gray-200 px-6 py-4 rounded-t-2xl flex items-center justify-between">
              <h2 className="text-2xl font-bold text-[#212121]">Seleccionar Cliente</h2>
              <button
                onClick={() => setShowCustomerModal(false)}
                className="p-2 hover:bg-gray-100 rounded-lg"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            <div className="p-6">
              {/* Opción de cliente genérico */}
              <button
                onClick={() => selectCustomer(GENERIC_CUSTOMER)}
                className="w-full mb-4 p-4 bg-[#FF6B00] bg-opacity-10 border-2 border-[#FF6B00] rounded-lg hover:bg-opacity-20 transition-colors"
              >
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 bg-[#FF6B00] rounded-full flex items-center justify-center">
                    <User className="w-6 h-6 text-white" />
                  </div>
                  <div className="text-left">
                    <p className="font-bold text-[#212121]">Cliente Genérico</p>
                    <p className="text-sm text-gray-600">Consumidor Final - Venta rápida</p>
                  </div>
                </div>
              </button>

              <div className="relative mb-4">
                <div className="absolute inset-0 flex items-center">
                  <div className="w-full border-t border-gray-300"></div>
                </div>
                <div className="relative flex justify-center text-sm">
                  <span className="px-2 bg-white text-gray-500">O buscar cliente específico</span>
                </div>
              </div>

              {/* Búsqueda */}
              <div className="mb-4">
                <div className="relative">
                  <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
                  <input
                    type="text"
                    placeholder="Buscar cliente por nombre o email..."
                    value={searchCustomer}
                    onChange={(e) => setSearchCustomer(e.target.value)}
                    className="input-field pl-10"
                    autoFocus
                  />
                </div>
              </div>

              {/* Lista de clientes */}
              <div className="space-y-2 max-h-96 overflow-y-auto">
                {filteredCustomers.map((customer) => (
                  <button
                    key={customer.id}
                    onClick={() => selectCustomer(customer)}
                    className="w-full text-left p-4 bg-[#F5F5F5] hover:bg-[#FF6B00] hover:bg-opacity-10 rounded-lg transition-colors"
                  >
                    <div className="flex items-center justify-between">
                      <div>
                        <p className="font-semibold text-[#212121]">{customer.name}</p>
                        {customer.tax_id && (
                          <p className="text-sm text-gray-600">{customer.tax_id}</p>
                        )}
                        {customer.email && (
                          <p className="text-xs text-gray-500">{customer.email}</p>
                        )}
                      </div>
                      <span className={`badge ${customer.is_active ? 'badge-success' : 'badge bg-gray-500 text-white'}`}>
                        {customer.is_active ? 'Activo' : 'Inactivo'}
                      </span>
                    </div>
                  </button>
                ))}

                {filteredCustomers.length === 0 && (
                  <div className="text-center py-8">
                    <p className="text-gray-500">No se encontraron clientes</p>
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Modal de Éxito después de procesar la venta */}
      {showSuccessModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-2xl shadow-card-hover max-w-md w-full">
            <div className="p-6 text-center">
              {/* Icono de éxito */}
              <div className="w-20 h-20 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-4">
                <svg className="w-10 h-10 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
              </div>

              {/* Mensaje */}
              <h2 className="text-2xl font-bold text-[#212121] mb-2">¡Venta Completada!</h2>
              <p className="text-gray-600 mb-6">
                La factura se ha generado exitosamente por un total de{' '}
                <span className="font-bold text-[#FF6B00]">
                  RD$ {calculateTotal().toFixed(2)}
                </span>
              </p>

              {/* Botones de acción */}
              <div className="space-y-3">
                <button
                  onClick={printInvoice}
                  className="btn-primary w-full flex items-center justify-center gap-2"
                >
                  <Printer className="w-5 h-5" />
                  Imprimir Factura
                </button>
                
                <button
                  onClick={startNewSale}
                  className="btn-outline w-full flex items-center justify-center gap-2"
                >
                  <ShoppingCart className="w-5 h-5" />
                  Nueva Venta
                </button>
              </div>

              {/* Información adicional */}
              <div className="mt-6 pt-6 border-t border-gray-200">
                <p className="text-xs text-gray-500">
                  Puedes encontrar todas tus facturas en la sección de <strong>Facturas</strong>
                </p>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Modal de Impresión de Factura */}
      {showInvoicePrint && createdInvoiceId && (
        <InvoicePrint
          invoiceId={createdInvoiceId}
          onClose={() => {
            setShowInvoicePrint(false);
            startNewSale();
          }}
        />
      )}
    </div>
  );
}
