import { useState, useEffect } from 'react'
import api from '../lib/api'
import Alert from '../components/Alert'

export default function Settings() {
  const [activeTab, setActiveTab] = useState('company')
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [alert, setAlert] = useState(null)
  const [config, setConfig] = useState({
    // Company info
    company_legal_name: '',
    company_trade_name: '',
    company_tax_id: '',
    company_email: '',
    company_phone: '',
    company_website: '',

    // Address
    address_line1: '',
    address_line2: '',
    city: '',
    state_province: '',
    postal_code: '',
    country: 'DO',

    // Fiscal
    tax_regime: '',
    dgii_rnc: '',
    dgii_user: '',
    dgii_api_key: '',
    dgii_environment: 'sandbox',
    ncf_enabled: true,

    // Invoicing
    invoice_prefix: 'INV',
    invoice_next_number: 1,
    invoice_terms: '',
    invoice_footer: '',
    payment_terms_days: 30,

    // NCF sequences
    ncf_fiscal_credit_prefix: '',
    ncf_fiscal_credit_start: 1,
    ncf_fiscal_credit_end: 999999999,
    ncf_fiscal_credit_sequence: 1,
    ncf_consumer_prefix: '',
    ncf_consumer_start: 1,
    ncf_consumer_end: 999999999,
    ncf_consumer_sequence: 1,
    ncf_debit_note_prefix: '',
    ncf_debit_note_start: 1,
    ncf_debit_note_end: 999999999,
    ncf_debit_note_sequence: 1,
    ncf_credit_note_prefix: '',
    ncf_credit_note_start: 1,
    ncf_credit_note_end: 999999999,
    ncf_credit_note_sequence: 1,

    // Currency
    default_currency: 'DOP',
    default_tax_rate: 18.00,

    // Notifications
    notification_email: '',
    send_invoice_emails: true,
    send_payment_reminders: true,

    // Branding
    logo_url: '',
    primary_color: '#3B82F6',
    secondary_color: '#10B981',

    // Localization
    timezone: 'America/Santo_Domingo',
    date_format: 'DD/MM/YYYY',
    time_format: 'HH:mm',
    locale: 'es_DO',
  })

  useEffect(() => {
    fetchConfig()
  }, [])

  const fetchConfig = async () => {
    try {
      setLoading(true)
      const response = await api.get('/config')
      // Ensure all fields have proper values
      const data = response.data
      setConfig({
        ...config,
        ...data,
        // Ensure non-null values for required fields
        company_legal_name: data.company_legal_name || '',
        country: data.country || 'DO',
        dgii_environment: data.dgii_environment || 'sandbox',
        ncf_enabled: data.ncf_enabled !== undefined ? data.ncf_enabled : true,
        invoice_prefix: data.invoice_prefix || 'INV',
        invoice_next_number: data.invoice_next_number || 1,
        payment_terms_days: data.payment_terms_days || 30,
        ncf_fiscal_credit_prefix: data.ncf_fiscal_credit_prefix || '',
        ncf_fiscal_credit_start: data.ncf_fiscal_credit_start || 1,
        ncf_fiscal_credit_end: data.ncf_fiscal_credit_end || 999999999,
        ncf_fiscal_credit_sequence: data.ncf_fiscal_credit_sequence || 1,
        ncf_consumer_prefix: data.ncf_consumer_prefix || '',
        ncf_consumer_start: data.ncf_consumer_start || 1,
        ncf_consumer_end: data.ncf_consumer_end || 999999999,
        ncf_consumer_sequence: data.ncf_consumer_sequence || 1,
        ncf_debit_note_prefix: data.ncf_debit_note_prefix || '',
        ncf_debit_note_start: data.ncf_debit_note_start || 1,
        ncf_debit_note_end: data.ncf_debit_note_end || 999999999,
        ncf_debit_note_sequence: data.ncf_debit_note_sequence || 1,
        ncf_credit_note_prefix: data.ncf_credit_note_prefix || '',
        ncf_credit_note_start: data.ncf_credit_note_start || 1,
        ncf_credit_note_end: data.ncf_credit_note_end || 999999999,
        ncf_credit_note_sequence: data.ncf_credit_note_sequence || 1,
        default_currency: data.default_currency || 'DOP',
        default_tax_rate: data.default_tax_rate || 18.00,
        send_invoice_emails: data.send_invoice_emails !== undefined ? data.send_invoice_emails : true,
        send_payment_reminders: data.send_payment_reminders !== undefined ? data.send_payment_reminders : true,
        primary_color: data.primary_color || '#3B82F6',
        secondary_color: data.secondary_color || '#10B981',
        timezone: data.timezone || 'America/Santo_Domingo',
        date_format: data.date_format || 'DD/MM/YYYY',
        time_format: data.time_format || 'HH:mm',
        locale: data.locale || 'es_DO',
      })
    } catch (error) {
      console.error('Error loading config:', error)
      setAlert({ 
        type: 'error', 
        message: error.response?.data?.error || 'Error al cargar la configuración'
      })
    } finally {
      setLoading(false)
    }
  }

  const handleChange = (field, value) => {
    setConfig(prev => ({ ...prev, [field]: value }))
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    try {
      setSaving(true)
      await api.put('/config', config)
      setAlert({ type: 'success', message: 'Configuración guardada exitosamente' })
    } catch (error) {
      setAlert({ 
        type: 'error', 
        message: error.response?.data?.error || 'Error al guardar la configuración' 
      })
    } finally {
      setSaving(false)
    }
  }

  const tabs = [
    { id: 'company', label: 'Empresa' },
    { id: 'fiscal', label: 'Configuración Fiscal' },
    { id: 'invoicing', label: 'Facturación' },
    { id: 'ncf', label: 'NCF (DGII)' },
    { id: 'notifications', label: 'Notificaciones' },
    { id: 'branding', label: 'Identidad Visual' },
  ]

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-lg text-gray-600">Cargando configuración...</div>
      </div>
    )
  }

  return (
    <div className="max-w-6xl">
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-[#212121]">⚙️ Configuración</h1>
        <p className="mt-2 text-gray-600">
          Administra la configuración de tu empresa y preferencias del sistema
        </p>
      </div>

      {alert && (
        <div className="mb-6">
          <Alert
            type={alert.type}
            message={alert.message}
            onClose={() => setAlert(null)}
          />
        </div>
      )}

      <div className="bg-white rounded-2xl shadow-card">
        {/* Tabs */}
        <div className="border-b border-gray-200">
          <nav className="flex -mb-px space-x-8 overflow-x-auto px-6">
            {tabs.map((tab) => (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id)}
                className={`
                  whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm transition-colors
                  ${activeTab === tab.id
                    ? 'border-[#FF6B00] text-[#FF6B00]'
                    : 'border-transparent text-gray-500 hover:text-[#FF6B00] hover:border-gray-300'
                  }
                `}
              >
                {tab.label}
              </button>
            ))}
          </nav>
        </div>

        {/* Content */}
        <form onSubmit={handleSubmit} className="p-8">
          {activeTab === 'company' && (
            <div className="space-y-8">
              <div>
                <h2 className="text-xl font-bold text-[#212121] mb-6 flex items-center gap-2">
                  🏢 Información de la Empresa
                </h2>
                <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
                  <div>
                    <label className="block text-sm font-semibold text-[#212121] mb-2">
                      Nombre Legal de la Empresa *
                    </label>
                    <input
                      type="text"
                      value={config.company_legal_name}
                      onChange={(e) => handleChange('company_legal_name', e.target.value)}
                      className="input-field"
                      required
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-semibold text-[#212121] mb-2">
                      Nombre Comercial
                    </label>
                    <input
                      type="text"
                      value={config.company_trade_name || ''}
                      onChange={(e) => handleChange('company_trade_name', e.target.value)}
                      className="input-field"
                      placeholder="Nombre como aparece al público"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-semibold text-[#212121] mb-2">
                      RNC / Cédula
                    </label>
                    <input
                      type="text"
                      value={config.company_tax_id || ''}
                      onChange={(e) => handleChange('company_tax_id', e.target.value)}
                      className="input-field"
                      placeholder="131234567"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-semibold text-[#212121] mb-2">
                      Email Corporativo
                    </label>
                    <input
                      type="email"
                      value={config.company_email || ''}
                      onChange={(e) => handleChange('company_email', e.target.value)}
                      className="input-field"
                      placeholder="contacto@miempresa.com"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-semibold text-[#212121] mb-2">
                      Teléfono
                    </label>
                    <input
                      type="tel"
                      value={config.company_phone || ''}
                      onChange={(e) => handleChange('company_phone', e.target.value)}
                      className="input-field"
                      placeholder="+1 809-555-1234"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-semibold text-[#212121] mb-2">
                      Sitio Web
                    </label>
                    <input
                      type="url"
                      value={config.company_website || ''}
                      onChange={(e) => handleChange('company_website', e.target.value)}
                      className="input-field"
                      placeholder="https://www.miempresa.com"
                    />
                  </div>
                </div>
              </div>

              <div className="pt-6 border-t border-gray-200">
                <h2 className="text-xl font-bold text-[#212121] mb-6 flex items-center gap-2">
                  📍 Dirección Fiscal
                </h2>
                <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
                  <div className="md:col-span-2">
                    <label className="block text-sm font-semibold text-[#212121] mb-2">
                      Dirección (Línea 1)
                    </label>
                    <input
                      type="text"
                      value={config.address_line1 || ''}
                      onChange={(e) => handleChange('address_line1', e.target.value)}
                      className="input-field"
                      placeholder="Av. 27 de Febrero #123"
                    />
                  </div>

                  <div className="md:col-span-2">
                    <label className="block text-sm font-semibold text-[#212121] mb-2">
                      Dirección (Línea 2)
                    </label>
                    <input
                      type="text"
                      value={config.address_line2 || ''}
                      onChange={(e) => handleChange('address_line2', e.target.value)}
                      className="input-field"
                      placeholder="Edificio, piso, suite (opcional)"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-semibold text-[#212121] mb-2">
                      Ciudad
                    </label>
                    <input
                      type="text"
                      value={config.city || ''}
                      onChange={(e) => handleChange('city', e.target.value)}
                      className="input-field"
                      placeholder="Santo Domingo"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-semibold text-[#212121] mb-2">
                      Provincia / Estado
                    </label>
                    <input
                      type="text"
                      value={config.state_province || ''}
                      onChange={(e) => handleChange('state_province', e.target.value)}
                      className="input-field"
                      placeholder="Distrito Nacional"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-semibold text-[#212121] mb-2">
                      Código Postal
                    </label>
                    <input
                      type="text"
                      value={config.postal_code || ''}
                      onChange={(e) => handleChange('postal_code', e.target.value)}
                      className="input-field"
                      placeholder="10101"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-semibold text-[#212121] mb-2">
                      País
                    </label>
                    <select
                      value={config.country}
                      onChange={(e) => handleChange('country', e.target.value)}
                      className="input-field"
                    >
                      <option value="DO">🇩🇴 República Dominicana</option>
                      <option value="US">🇺🇸 Estados Unidos</option>
                      <option value="ES">🇪🇸 España</option>
                      <option value="MX">🇲🇽 México</option>
                      <option value="CO">🇨🇴 Colombia</option>
                    </select>
                  </div>
                </div>
              </div>
            </div>
          )}

          {activeTab === 'fiscal' && (
            <div className="space-y-8">
              <div>
                <h2 className="text-xl font-bold text-[#212121] mb-6 flex items-center gap-2">
                  🏛️ Configuración DGII
                </h2>
                <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
                  <div>
                    <label className="block text-sm font-semibold text-[#212121] mb-2">
                      Régimen Fiscal
                    </label>
                    <input
                      type="text"
                      value={config.tax_regime || ''}
                      onChange={(e) => handleChange('tax_regime', e.target.value)}
                      className="input-field"
                      placeholder="Ej: Régimen Simplificado"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-semibold text-[#212121] mb-2">
                      RNC (DGII)
                    </label>
                    <input
                      type="text"
                      value={config.dgii_rnc || ''}
                      onChange={(e) => handleChange('dgii_rnc', e.target.value)}
                      className="input-field"
                      placeholder="131234567"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-semibold text-[#212121] mb-2">
                      Usuario DGII
                    </label>
                    <input
                      type="text"
                      value={config.dgii_user || ''}
                      onChange={(e) => handleChange('dgii_user', e.target.value)}
                      className="input-field"
                      placeholder="usuario_dgii"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-semibold text-[#212121] mb-2">
                      API Key DGII
                    </label>
                    <input
                      type="password"
                      value={config.dgii_api_key || ''}
                      onChange={(e) => handleChange('dgii_api_key', e.target.value)}
                      className="input-field"
                      placeholder="•••••••••••••"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-semibold text-[#212121] mb-2">
                      Ambiente DGII
                    </label>
                    <select
                      value={config.dgii_environment}
                      onChange={(e) => handleChange('dgii_environment', e.target.value)}
                      className="input-field"
                    >
                      <option value="sandbox">🧪 Sandbox (Pruebas)</option>
                      <option value="production">🚀 Producción</option>
                    </select>
                  </div>

                  <div className="flex items-center mt-6">
                    <input
                      type="checkbox"
                      checked={config.ncf_enabled}
                      onChange={(e) => handleChange('ncf_enabled', e.target.checked)}
                      className="w-5 h-5 text-[#FF6B00] border-gray-300 rounded focus:ring-[#FF6B00] focus:ring-2"
                    />
                    <label className="block ml-3 text-sm font-medium text-[#212121]">
                      Habilitar NCF (Comprobante Fiscal)
                    </label>
                  </div>
                </div>
              </div>

              <div className="pt-6 border-t border-gray-200">
                <h2 className="text-xl font-bold text-[#212121] mb-6 flex items-center gap-2">
                  💰 Configuración Monetaria
                </h2>
                <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
                  <div>
                    <label className="block text-sm font-semibold text-[#212121] mb-2">
                      Moneda por Defecto
                    </label>
                    <select
                      value={config.default_currency}
                      onChange={(e) => handleChange('default_currency', e.target.value)}
                      className="input-field"
                    >
                      <option value="DOP">💵 DOP - Peso Dominicano</option>
                      <option value="USD">💵 USD - Dólar Estadounidense</option>
                      <option value="EUR">💶 EUR - Euro</option>
                    </select>
                  </div>

                  <div>
                    <label className="block text-sm font-semibold text-[#212121] mb-2">
                      Tasa de Impuesto por Defecto (%)
                    </label>
                    <input
                      type="number"
                      step="0.01"
                      value={config.default_tax_rate}
                      onChange={(e) => handleChange('default_tax_rate', parseFloat(e.target.value))}
                      className="input-field"
                      placeholder="18.00"
                    />
                    <p className="mt-2 text-xs text-gray-500">💡 ITBIS en República Dominicana es 18%</p>
                  </div>
                </div>
              </div>
            </div>
          )}

          {activeTab === 'invoicing' && (
            <div className="space-y-8">
              <div>
                <h2 className="text-xl font-bold text-[#212121] mb-6 flex items-center gap-2">
                  📄 Configuración de Facturación
                </h2>
                <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
                  <div>
                    <label className="block text-sm font-semibold text-[#212121] mb-2">
                      Prefijo de Factura
                    </label>
                    <input
                      type="text"
                      value={config.invoice_prefix}
                      onChange={(e) => handleChange('invoice_prefix', e.target.value)}
                      className="input-field"
                      placeholder="INV"
                    />
                    <p className="mt-2 text-xs text-gray-500">💡 Ej: INV, FACT, FCT</p>
                  </div>

                  <div>
                    <label className="block text-sm font-semibold text-[#212121] mb-2">
                      Siguiente Número de Factura
                    </label>
                    <input
                      type="number"
                      value={config.invoice_next_number}
                      onChange={(e) => handleChange('invoice_next_number', parseInt(e.target.value))}
                      className="input-field"
                      placeholder="1"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-semibold text-[#212121] mb-2">
                      Días de Pago por Defecto
                    </label>
                    <input
                      type="number"
                      value={config.payment_terms_days}
                      onChange={(e) => handleChange('payment_terms_days', parseInt(e.target.value))}
                      className="input-field"
                      placeholder="30"
                    />
                    <p className="mt-2 text-xs text-gray-500">💡 Plazo de pago en días</p>
                  </div>
                </div>
              </div>

              <div>
                <label className="block text-sm font-semibold text-[#212121] mb-2">
                  Términos y Condiciones de Factura
                </label>
                <textarea
                  value={config.invoice_terms || ''}
                  onChange={(e) => handleChange('invoice_terms', e.target.value)}
                  rows={4}
                  className="input-field resize-none"
                  placeholder="Términos y condiciones que aparecerán en las facturas..."
                />
              </div>

              <div>
                <label className="block text-sm font-semibold text-[#212121] mb-2">
                  Pie de Página de Factura
                </label>
                <textarea
                  value={config.invoice_footer || ''}
                  onChange={(e) => handleChange('invoice_footer', e.target.value)}
                  rows={3}
                  className="input-field resize-none"
                  placeholder="Gracias por su preferencia. ¡Vuelva pronto!"
                />
              </div>
            </div>
          )}

          {activeTab === 'ncf' && (
            <div className="space-y-8">
              <div className="p-5 bg-gradient-to-r from-[#FF6B00] to-[#E55D00] bg-opacity-10 border-l-4 border-[#FF6B00] rounded-lg">
                <p className="text-sm text-[#212121] font-medium">
                  📋 <strong>Números de Comprobante Fiscal (NCF)</strong>
                </p>
                <p className="text-xs text-gray-600 mt-1">
                  Los NCF son requeridos por la DGII en República Dominicana. Configure aquí los prefijos y secuencias para cada tipo de comprobante.
                </p>
              </div>

              <div className="grid grid-cols-1 gap-8">
                {/* NCF 01 - Crédito Fiscal */}
                <div className="bg-[#F5F5F5] rounded-lg p-6">
                  <h3 className="mb-4 text-lg font-bold text-[#212121] flex items-center gap-2">
                    🟢 NCF 01 - Crédito Fiscal
                  </h3>
                  <div className="grid grid-cols-1 gap-6 md:grid-cols-3">
                    <div>
                      <label className="block text-sm font-semibold text-[#212121] mb-2">
                        Prefijo
                      </label>
                      <input
                        type="text"
                        value={config.ncf_fiscal_credit_prefix || ''}
                        onChange={(e) => handleChange('ncf_fiscal_credit_prefix', e.target.value)}
                        className="input-field bg-white"
                        placeholder="B01"
                      />
                      <p className="mt-1 text-xs text-gray-500">Asignado por DGII</p>
                    </div>
                    <div>
                      <label className="block text-sm font-semibold text-[#212121] mb-2">
                        Inicio de Rango
                      </label>
                      <input
                        type="number"
                        value={config.ncf_fiscal_credit_start}
                        onChange={(e) => handleChange('ncf_fiscal_credit_start', parseInt(e.target.value))}
                        className="input-field bg-white"
                        placeholder="1"
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-semibold text-[#212121] mb-2">
                        Fin de Rango
                      </label>
                      <input
                        type="number"
                        value={config.ncf_fiscal_credit_end}
                        onChange={(e) => handleChange('ncf_fiscal_credit_end', parseInt(e.target.value))}
                        className="input-field bg-white"
                        placeholder="100000"
                      />
                    </div>
                    <div className="md:col-span-3">
                      <label className="block text-sm font-semibold text-[#212121] mb-2">
                        Secuencia Actual
                      </label>
                      <input
                        type="number"
                        value={config.ncf_fiscal_credit_sequence}
                        onChange={(e) => handleChange('ncf_fiscal_credit_sequence', parseInt(e.target.value))}
                        className="input-field bg-white"
                        placeholder="1"
                      />
                      {config.ncf_fiscal_credit_sequence >= config.ncf_fiscal_credit_end * 0.9 && (
                        <div className="mt-2 p-3 bg-yellow-50 border border-yellow-300 rounded-lg">
                          <p className="text-xs text-yellow-800 font-medium flex items-center gap-1">
                            ⚠️ Advertencia: Te estás acercando al límite del rango ({config.ncf_fiscal_credit_sequence} / {config.ncf_fiscal_credit_end})
                          </p>
                        </div>
                      )}
                      {config.ncf_fiscal_credit_sequence >= config.ncf_fiscal_credit_end && (
                        <div className="mt-2 p-3 bg-red-50 border border-red-300 rounded-lg">
                          <p className="text-xs text-red-800 font-medium flex items-center gap-1">
                            🚨 Error: Has alcanzado o superado el límite del rango. Contacta a DGII para obtener un nuevo rango.
                          </p>
                        </div>
                      )}
                    </div>
                  </div>
                </div>

                {/* NCF 02 - Consumidor Final */}
                <div className="bg-[#F5F5F5] rounded-lg p-6">
                  <h3 className="mb-4 text-lg font-bold text-[#212121] flex items-center gap-2">
                    🟠 NCF 02 - Consumidor Final
                  </h3>
                  <div className="grid grid-cols-1 gap-6 md:grid-cols-3">
                    <div>
                      <label className="block text-sm font-semibold text-[#212121] mb-2">
                        Prefijo
                      </label>
                      <input
                        type="text"
                        value={config.ncf_consumer_prefix || ''}
                        onChange={(e) => handleChange('ncf_consumer_prefix', e.target.value)}
                        className="input-field bg-white"
                        placeholder="B02"
                      />
                      <p className="mt-1 text-xs text-gray-500">Asignado por DGII</p>
                    </div>
                    <div>
                      <label className="block text-sm font-semibold text-[#212121] mb-2">
                        Inicio de Rango
                      </label>
                      <input
                        type="number"
                        value={config.ncf_consumer_start}
                        onChange={(e) => handleChange('ncf_consumer_start', parseInt(e.target.value))}
                        className="input-field bg-white"
                        placeholder="1"
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-semibold text-[#212121] mb-2">
                        Fin de Rango
                      </label>
                      <input
                        type="number"
                        value={config.ncf_consumer_end}
                        onChange={(e) => handleChange('ncf_consumer_end', parseInt(e.target.value))}
                        className="input-field bg-white"
                        placeholder="100000"
                      />
                    </div>
                    <div className="md:col-span-3">
                      <label className="block text-sm font-semibold text-[#212121] mb-2">
                        Secuencia Actual
                      </label>
                      <input
                        type="number"
                        value={config.ncf_consumer_sequence}
                        onChange={(e) => handleChange('ncf_consumer_sequence', parseInt(e.target.value))}
                        className="input-field bg-white"
                        placeholder="1"
                      />
                      {config.ncf_consumer_sequence >= config.ncf_consumer_end * 0.9 && (
                        <div className="mt-2 p-3 bg-yellow-50 border border-yellow-300 rounded-lg">
                          <p className="text-xs text-yellow-800 font-medium flex items-center gap-1">
                            ⚠️ Advertencia: Te estás acercando al límite del rango ({config.ncf_consumer_sequence} / {config.ncf_consumer_end})
                          </p>
                        </div>
                      )}
                      {config.ncf_consumer_sequence >= config.ncf_consumer_end && (
                        <div className="mt-2 p-3 bg-red-50 border border-red-300 rounded-lg">
                          <p className="text-xs text-red-800 font-medium flex items-center gap-1">
                            🚨 Error: Has alcanzado o superado el límite del rango. Contacta a DGII para obtener un nuevo rango.
                          </p>
                        </div>
                      )}
                    </div>
                  </div>
                </div>

                {/* NCF 03 - Nota de Débito */}
                <div className="bg-[#F5F5F5] rounded-lg p-6">
                  <h3 className="mb-4 text-lg font-bold text-[#212121] flex items-center gap-2">
                    📝 NCF 03 - Nota de Débito
                  </h3>
                  <div className="grid grid-cols-1 gap-6 md:grid-cols-3">
                    <div>
                      <label className="block text-sm font-semibold text-[#212121] mb-2">
                        Prefijo
                      </label>
                      <input
                        type="text"
                        value={config.ncf_debit_note_prefix || ''}
                        onChange={(e) => handleChange('ncf_debit_note_prefix', e.target.value)}
                        className="input-field bg-white"
                        placeholder="B03"
                      />
                      <p className="mt-1 text-xs text-gray-500">Asignado por DGII</p>
                    </div>
                    <div>
                      <label className="block text-sm font-semibold text-[#212121] mb-2">
                        Inicio de Rango
                      </label>
                      <input
                        type="number"
                        value={config.ncf_debit_note_start}
                        onChange={(e) => handleChange('ncf_debit_note_start', parseInt(e.target.value))}
                        className="input-field bg-white"
                        placeholder="1"
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-semibold text-[#212121] mb-2">
                        Fin de Rango
                      </label>
                      <input
                        type="number"
                        value={config.ncf_debit_note_end}
                        onChange={(e) => handleChange('ncf_debit_note_end', parseInt(e.target.value))}
                        className="input-field bg-white"
                        placeholder="100000"
                      />
                    </div>
                    <div className="md:col-span-3">
                      <label className="block text-sm font-semibold text-[#212121] mb-2">
                        Secuencia Actual
                      </label>
                      <input
                        type="number"
                        value={config.ncf_debit_note_sequence}
                        onChange={(e) => handleChange('ncf_debit_note_sequence', parseInt(e.target.value))}
                        className="input-field bg-white"
                        placeholder="1"
                      />
                      {config.ncf_debit_note_sequence >= config.ncf_debit_note_end * 0.9 && config.ncf_debit_note_sequence < config.ncf_debit_note_end && (
                        <div className="mt-2 p-3 bg-yellow-50 border border-yellow-300 rounded-lg">
                          <p className="text-xs text-yellow-800 font-medium flex items-center gap-1">
                            ⚠️ Advertencia: Te estás acercando al límite del rango
                          </p>
                        </div>
                      )}
                    </div>
                  </div>
                </div>

                {/* NCF 04 - Nota de Crédito */}
                <div className="bg-[#F5F5F5] rounded-lg p-6">
                  <h3 className="mb-4 text-lg font-bold text-[#212121] flex items-center gap-2">
                    📝 NCF 04 - Nota de Crédito
                  </h3>
                  <div className="grid grid-cols-1 gap-6 md:grid-cols-3">
                    <div>
                      <label className="block text-sm font-semibold text-[#212121] mb-2">
                        Prefijo
                      </label>
                      <input
                        type="text"
                        value={config.ncf_credit_note_prefix || ''}
                        onChange={(e) => handleChange('ncf_credit_note_prefix', e.target.value)}
                        className="input-field bg-white"
                        placeholder="B04"
                      />
                      <p className="mt-1 text-xs text-gray-500">Asignado por DGII</p>
                    </div>
                    <div>
                      <label className="block text-sm font-semibold text-[#212121] mb-2">
                        Inicio de Rango
                      </label>
                      <input
                        type="number"
                        value={config.ncf_credit_note_start}
                        onChange={(e) => handleChange('ncf_credit_note_start', parseInt(e.target.value))}
                        className="input-field bg-white"
                        placeholder="1"
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-semibold text-[#212121] mb-2">
                        Fin de Rango
                      </label>
                      <input
                        type="number"
                        value={config.ncf_credit_note_end}
                        onChange={(e) => handleChange('ncf_credit_note_end', parseInt(e.target.value))}
                        className="input-field bg-white"
                        placeholder="100000"
                      />
                    </div>
                    <div className="md:col-span-3">
                      <label className="block text-sm font-semibold text-[#212121] mb-2">
                        Secuencia Actual
                      </label>
                      <input
                        type="number"
                        value={config.ncf_credit_note_sequence}
                        onChange={(e) => handleChange('ncf_credit_note_sequence', parseInt(e.target.value))}
                        className="input-field bg-white"
                        placeholder="1"
                      />
                      {config.ncf_credit_note_sequence >= config.ncf_credit_note_end * 0.9 && config.ncf_credit_note_sequence < config.ncf_credit_note_end && (
                        <div className="mt-2 p-3 bg-yellow-50 border border-yellow-300 rounded-lg">
                          <p className="text-xs text-yellow-800 font-medium flex items-center gap-1">
                            ⚠️ Advertencia: Te estás acercando al límite del rango
                          </p>
                        </div>
                      )}
                    </div>
                  </div>
                </div>
              </div>
            </div>
          )}

          {activeTab === 'notifications' && (
            <div className="space-y-8">
              <div>
                <h2 className="text-xl font-bold text-[#212121] mb-6 flex items-center gap-2">
                  📧 Notificaciones por Email
                </h2>
                <div>
                  <label className="block text-sm font-semibold text-[#212121] mb-2">
                    Email de Notificaciones
                  </label>
                  <input
                    type="email"
                    value={config.notification_email || ''}
                    onChange={(e) => handleChange('notification_email', e.target.value)}
                    className="input-field"
                    placeholder="notificaciones@miempresa.com"
                  />
                  <p className="mt-2 text-xs text-gray-500">
                    💡 Email donde se recibirán las notificaciones del sistema
                  </p>
                </div>
              </div>

              <div className="pt-6 border-t border-gray-200">
                <h2 className="text-xl font-bold text-[#212121] mb-6 flex items-center gap-2">
                  🔔 Preferencias de Notificación
                </h2>
                <div className="space-y-5">
                  <div className="flex items-start bg-[#F5F5F5] rounded-lg p-4 hover:bg-gray-100 transition-colors">
                    <div className="flex items-center h-5 mt-1">
                      <input
                        type="checkbox"
                        checked={config.send_invoice_emails}
                        onChange={(e) => handleChange('send_invoice_emails', e.target.checked)}
                        className="w-5 h-5 text-[#FF6B00] border-gray-300 rounded focus:ring-[#FF6B00] focus:ring-2"
                      />
                    </div>
                    <div className="ml-4">
                      <label className="text-sm font-semibold text-[#212121] cursor-pointer">
                        📨 Enviar facturas por email
                      </label>
                      <p className="text-xs text-gray-600 mt-1">
                        Enviar automáticamente las facturas a los clientes por correo electrónico
                      </p>
                    </div>
                  </div>

                  <div className="flex items-start bg-[#F5F5F5] rounded-lg p-4 hover:bg-gray-100 transition-colors">
                    <div className="flex items-center h-5 mt-1">
                      <input
                        type="checkbox"
                        checked={config.send_payment_reminders}
                        onChange={(e) => handleChange('send_payment_reminders', e.target.checked)}
                        className="w-5 h-5 text-[#FF6B00] border-gray-300 rounded focus:ring-[#FF6B00] focus:ring-2"
                      />
                    </div>
                    <div className="ml-4">
                      <label className="text-sm font-semibold text-[#212121] cursor-pointer">
                        ⏰ Recordatorios de pago
                      </label>
                      <p className="text-xs text-gray-600 mt-1">
                        Enviar recordatorios automáticos para facturas vencidas
                      </p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          )}

          {activeTab === 'branding' && (
            <div className="space-y-8">
              <div>
                <h2 className="text-xl font-bold text-[#212121] mb-6 flex items-center gap-2">
                  🎨 Identidad Visual
                </h2>
                <div>
                  <label className="block text-sm font-semibold text-[#212121] mb-2">
                    URL del Logo
                  </label>
                  <input
                    type="url"
                    value={config.logo_url || ''}
                    onChange={(e) => handleChange('logo_url', e.target.value)}
                    className="input-field"
                    placeholder="https://ejemplo.com/logo.png"
                  />
                  <p className="mt-2 text-xs text-gray-500">
                    💡 URL pública del logo de tu empresa (aparecerá en facturas)
                  </p>
                </div>
              </div>

              <div className="pt-6 border-t border-gray-200">
                <h2 className="text-xl font-bold text-[#212121] mb-6 flex items-center gap-2">
                  🎨 Colores de Marca
                </h2>
                <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
                  <div>
                    <label className="block text-sm font-semibold text-[#212121] mb-2">
                      Color Primario
                    </label>
                    <div className="flex gap-3">
                      <input
                        type="color"
                        value={config.primary_color}
                        onChange={(e) => handleChange('primary_color', e.target.value)}
                        className="h-12 w-16 rounded-lg border-2 border-gray-300 cursor-pointer shadow-sm hover:border-[#FF6B00] transition-colors"
                      />
                      <input
                        type="text"
                        value={config.primary_color}
                        onChange={(e) => handleChange('primary_color', e.target.value)}
                        className="input-field flex-1 font-mono"
                        placeholder="#3B82F6"
                      />
                    </div>
                    <div 
                      className="mt-3 h-10 rounded-lg border-2 border-gray-200" 
                      style={{ backgroundColor: config.primary_color }}
                    ></div>
                  </div>

                  <div>
                    <label className="block text-sm font-semibold text-[#212121] mb-2">
                      Color Secundario
                    </label>
                    <div className="flex gap-3">
                      <input
                        type="color"
                        value={config.secondary_color}
                        onChange={(e) => handleChange('secondary_color', e.target.value)}
                        className="h-12 w-16 rounded-lg border-2 border-gray-300 cursor-pointer shadow-sm hover:border-[#FF6B00] transition-colors"
                      />
                      <input
                        type="text"
                        value={config.secondary_color}
                        onChange={(e) => handleChange('secondary_color', e.target.value)}
                        className="input-field flex-1 font-mono"
                        placeholder="#10B981"
                      />
                    </div>
                    <div 
                      className="mt-3 h-10 rounded-lg border-2 border-gray-200" 
                      style={{ backgroundColor: config.secondary_color }}
                    ></div>
                  </div>
                </div>
              </div>

              <div className="pt-6 border-t border-gray-200">
                <h2 className="text-xl font-bold text-[#212121] mb-6 flex items-center gap-2">
                  🌍 Localización
                </h2>
                <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
                  <div>
                    <label className="block text-sm font-semibold text-[#212121] mb-2">
                      Zona Horaria
                    </label>
                    <select
                      value={config.timezone}
                      onChange={(e) => handleChange('timezone', e.target.value)}
                      className="input-field"
                    >
                      <option value="America/Santo_Domingo">🇩🇴 América/Santo Domingo</option>
                      <option value="America/New_York">🇺🇸 América/Nueva York</option>
                      <option value="America/Los_Angeles">🇺🇸 América/Los Ángeles</option>
                      <option value="Europe/Madrid">🇪🇸 Europa/Madrid</option>
                    </select>
                  </div>

                  <div>
                    <label className="block text-sm font-semibold text-[#212121] mb-2">
                      Formato de Fecha
                    </label>
                    <select
                      value={config.date_format}
                      onChange={(e) => handleChange('date_format', e.target.value)}
                      className="input-field"
                    >
                      <option value="DD/MM/YYYY">📅 DD/MM/YYYY (Ej: 19/10/2025)</option>
                      <option value="MM/DD/YYYY">📅 MM/DD/YYYY (Ej: 10/19/2025)</option>
                      <option value="YYYY-MM-DD">📅 YYYY-MM-DD (Ej: 2025-10-19)</option>
                    </select>
                  </div>

                  <div>
                    <label className="block text-sm font-semibold text-[#212121] mb-2">
                      Formato de Hora
                    </label>
                    <select
                      value={config.time_format}
                      onChange={(e) => handleChange('time_format', e.target.value)}
                      className="input-field"
                    >
                      <option value="HH:mm">🕐 24 horas (HH:mm)</option>
                      <option value="hh:mm A">🕐 12 horas (hh:mm AM/PM)</option>
                    </select>
                  </div>

                  <div>
                    <label className="block text-sm font-semibold text-[#212121] mb-2">
                      Idioma
                    </label>
                    <select
                      value={config.locale}
                      onChange={(e) => handleChange('locale', e.target.value)}
                      className="input-field"
                    >
                      <option value="es_DO">🇩🇴 Español (República Dominicana)</option>
                      <option value="es_ES">🇪🇸 Español (España)</option>
                      <option value="en_US">🇺🇸 English (United States)</option>
                    </select>
                  </div>
                </div>
              </div>
            </div>
          )}

          <div className="flex justify-end gap-3 pt-8 mt-8 border-t border-gray-200">
            <button
              type="button"
              onClick={fetchConfig}
              disabled={saving}
              className="btn-outline disabled:opacity-50"
            >
              🔄 Recargar
            </button>
            <button
              type="submit"
              disabled={saving}
              className="btn-primary px-8 disabled:opacity-50 flex items-center gap-2"
            >
              {saving ? (
                <>
                  <div className="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
                  Guardando...
                </>
              ) : (
                <>
                  💾 Guardar Cambios
                </>
              )}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
