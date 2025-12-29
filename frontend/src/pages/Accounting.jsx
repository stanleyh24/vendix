import { useState, useEffect } from 'react'
import api from '../lib/api'
import Alert from '../components/Alert'
import { BookOpen, FileText, BarChart, Plus, Eye, Edit, Trash2, X, Settings } from 'lucide-react'

export default function Accounting() {
  const [activeTab, setActiveTab] = useState('accounts')
  const [alert, setAlert] = useState(null)
  const [loading, setLoading] = useState(false)

  // Accounts state
  const [accounts, setAccounts] = useState([])
  const [accountFilter, setAccountFilter] = useState('')
  const [showAccountModal, setShowAccountModal] = useState(false)
  const [editingAccount, setEditingAccount] = useState(null)
  const [accountForm, setAccountForm] = useState({
    account_code: '',
    account_name: '',
    account_type: 'ASSET',
    parent_account_code: '',
    description: '',
    is_active: true
  })

  // Journal entries state
  const [entries, setEntries] = useState([])
  const [showEntryModal, setShowEntryModal] = useState(false)
  const [expandedEntryId, setExpandedEntryId] = useState(null)

  // Reports state
  const [trialBalance, setTrialBalance] = useState([])
  const [balanceSheet, setBalanceSheet] = useState(null)
  const [incomeStatement, setIncomeStatement] = useState(null)

  // Account Mappings state
  const [accountMappings, setAccountMappings] = useState([])
  const [showMappingModal, setShowMappingModal] = useState(false)
  const [editingMapping, setEditingMapping] = useState(null)
  const [mappingForm, setMappingForm] = useState({
    transaction_type: '',
    account_code: '',
    description: ''
  })

  const tabs = [
    { id: 'accounts', label: 'Plan de Cuentas', icon: BookOpen },
    { id: 'entries', label: 'Asientos Contables', icon: FileText },
    { id: 'mappings', label: 'Configuración de Cuentas', icon: Settings },
    { id: 'reports', label: 'Reportes', icon: BarChart },
  ]

  useEffect(() => {
    if (activeTab === 'accounts') {
      fetchAccounts()
    } else if (activeTab === 'entries') {
      fetchEntries()
    } else if (activeTab === 'mappings') {
      fetchAccountMappings()
    } else if (activeTab === 'reports') {
      fetchTrialBalance()
    }
  }, [activeTab])

  const fetchAccounts = async () => {
    try {
      setLoading(true)
      console.log('Fetching accounts from /accounting/accounts...')
      const response = await api.get('/accounting/accounts')
      console.log('Accounts response:', response.data)
      
      if (response.data && Array.isArray(response.data)) {
        setAccounts(response.data)
        console.log(`Loaded ${response.data.length} accounts`)
      } else {
        console.warn('Response data is not an array:', response.data)
        setAccounts([])
      }
    } catch (error) {
      console.error('Error fetching accounts:', error)
      const errorMessage = error.response?.data?.error || error.message || 'Error al cargar cuentas'
      setAlert({ type: 'error', message: errorMessage })
      setAccounts([])
    } finally {
      setLoading(false)
    }
  }

  const fetchEntries = async () => {
    try {
      setLoading(true)
      const response = await api.get('/accounting/journal-entries')
      setEntries(response.data || [])
    } catch (error) {
      console.error('Error fetching entries:', error)
      const errorMessage = error.response?.data?.error || error.message || 'Error al cargar asientos'
      setAlert({ type: 'error', message: errorMessage })
      setEntries([])
    } finally {
      setLoading(false)
    }
  }

  const fetchTrialBalance = async () => {
    try {
      setLoading(true)
      const response = await api.get('/accounting/trial-balance')
      setTrialBalance(response.data || [])
    } catch (error) {
      console.error('Error fetching trial balance:', error)
      const errorMessage = error.response?.data?.error || error.message || 'Error al cargar balance'
      setAlert({ type: 'error', message: errorMessage })
      setTrialBalance([])
    } finally {
      setLoading(false)
    }
  }

  const fetchBalanceSheet = async () => {
    try {
      setLoading(true)
      const response = await api.get('/accounting/balance-sheet')
      setBalanceSheet(response.data)
    } catch (error) {
      setAlert({ type: 'error', message: 'Error al cargar balance general' })
    } finally {
      setLoading(false)
    }
  }

  const fetchIncomeStatement = async () => {
    try {
      setLoading(true)
      const response = await api.get('/accounting/income-statement')
      setIncomeStatement(response.data)
    } catch (error) {
      setAlert({ type: 'error', message: 'Error al cargar estado de resultados' })
    } finally {
      setLoading(false)
    }
  }

  const accountTypes = {
    ASSET: 'Activo',
    LIABILITY: 'Pasivo',
    EQUITY: 'Capital',
    REVENUE: 'Ingreso',
    EXPENSE: 'Gasto'
  }

  const getAccountTypeColor = (type) => {
    const colors = {
      ASSET: 'text-blue-600 bg-blue-50',
      LIABILITY: 'text-red-600 bg-red-50',
      EQUITY: 'text-purple-600 bg-purple-50',
      REVENUE: 'text-green-600 bg-green-50',
      EXPENSE: 'text-orange-600 bg-orange-50'
    }
    return colors[type] || 'text-gray-600 bg-gray-50'
  }

  const filteredAccounts = accounts.filter(account =>
    account.account_code.toLowerCase().includes(accountFilter.toLowerCase()) ||
    account.account_name.toLowerCase().includes(accountFilter.toLowerCase()) ||
    accountTypes[account.account_type]?.toLowerCase().includes(accountFilter.toLowerCase())
  )

  const calculateTrialBalanceTotals = () => {
    const totals = trialBalance.reduce((acc, entry) => {
      acc.debit += entry.debit
      acc.credit += entry.credit
      return acc
    }, { debit: 0, credit: 0 })
    return totals
  }

  const calculateEntryTotals = (entry) => {
    if (!entry.lines || entry.lines.length === 0) {
      return { debit: 0, credit: 0 }
    }
    return entry.lines.reduce((acc, line) => {
      acc.debit += line.debit || 0
      acc.credit += line.credit || 0
      return acc
    }, { debit: 0, credit: 0 })
  }

  const toggleEntryExpansion = (entryId) => {
    setExpandedEntryId(expandedEntryId === entryId ? null : entryId)
  }

  const openAccountModal = (account = null) => {
    if (account) {
      setEditingAccount(account)
      setAccountForm({
        account_code: account.account_code,
        account_name: account.account_name,
        account_type: account.account_type,
        parent_account_code: account.parent_account_code || '',
        description: account.description || '',
        is_active: account.is_active
      })
    } else {
      setEditingAccount(null)
      setAccountForm({
        account_code: '',
        account_name: '',
        account_type: 'ASSET',
        parent_account_code: '',
        description: '',
        is_active: true
      })
    }
    setShowAccountModal(true)
  }

  const closeAccountModal = () => {
    setShowAccountModal(false)
    setEditingAccount(null)
    setAccountForm({
      account_code: '',
      account_name: '',
      account_type: 'ASSET',
      parent_account_code: '',
      description: '',
      is_active: true
    })
  }

  const handleAccountSubmit = async (e) => {
    e.preventDefault()
    try {
      setLoading(true)
      const payload = {
        account_code: accountForm.account_code,
        account_name: accountForm.account_name,
        account_type: accountForm.account_type,
        description: accountForm.description || null,
        parent_account_code: accountForm.parent_account_code || null
      }

      if (editingAccount) {
        // Update account
        const updatePayload = {}
        if (accountForm.account_name !== editingAccount.account_name) {
          updatePayload.account_name = accountForm.account_name
        }
        if (accountForm.account_type !== editingAccount.account_type) {
          updatePayload.account_type = accountForm.account_type
        }
        if (accountForm.description !== (editingAccount.description || '')) {
          updatePayload.description = accountForm.description || null
        }
        if (accountForm.is_active !== editingAccount.is_active) {
          updatePayload.is_active = accountForm.is_active
        }
        if (accountForm.parent_account_code !== (editingAccount.parent_account_code || '')) {
          updatePayload.parent_account_code = accountForm.parent_account_code || null
        }

        await api.put(`/accounting/accounts/${editingAccount.account_code}`, updatePayload)
        setAlert({ type: 'success', message: 'Cuenta actualizada exitosamente' })
      } else {
        // Create account
        await api.post('/accounting/accounts', payload)
        setAlert({ type: 'success', message: 'Cuenta creada exitosamente' })
      }
      
      closeAccountModal()
      fetchAccounts()
    } catch (error) {
      const errorMessage = error.response?.data?.error || error.message || 'Error al guardar cuenta'
      setAlert({ type: 'error', message: errorMessage })
    } finally {
      setLoading(false)
    }
  }

  const handleDeleteAccount = async (account) => {
    if (!window.confirm(`¿Está seguro de eliminar la cuenta ${account.account_code} - ${account.account_name}?`)) {
      return
    }

    try {
      setLoading(true)
      await api.delete(`/accounting/accounts/${account.account_code}`)
      setAlert({ type: 'success', message: 'Cuenta eliminada exitosamente' })
      fetchAccounts()
    } catch (error) {
      const errorMessage = error.response?.data?.error || error.message || 'Error al eliminar cuenta'
      setAlert({ type: 'error', message: errorMessage })
    } finally {
      setLoading(false)
    }
  }

  // Account Mappings functions
  const fetchAccountMappings = async () => {
    try {
      setLoading(true)
      const response = await api.get('/accounting/account-mappings')
      setAccountMappings(response.data || [])
    } catch (error) {
      console.error('Error fetching account mappings:', error)
      const errorMessage = error.response?.data?.error || error.message || 'Error al cargar configuración de cuentas'
      setAlert({ type: 'error', message: errorMessage })
      setAccountMappings([])
    } finally {
      setLoading(false)
    }
  }

  const openMappingModal = (mapping = null) => {
    if (mapping) {
      setEditingMapping(mapping)
      setMappingForm({
        transaction_type: mapping.transaction_type,
        account_code: mapping.account_code,
        description: mapping.description || ''
      })
    } else {
      setEditingMapping(null)
      setMappingForm({
        transaction_type: '',
        account_code: '',
        description: ''
      })
    }
    setShowMappingModal(true)
  }

  const closeMappingModal = () => {
    setShowMappingModal(false)
    setEditingMapping(null)
    setMappingForm({
      transaction_type: '',
      account_code: '',
      description: ''
    })
  }

  const handleMappingSubmit = async (e) => {
    e.preventDefault()
    try {
      setLoading(true)
      const payload = {
        transaction_type: mappingForm.transaction_type,
        account_code: mappingForm.account_code,
        description: mappingForm.description || null
      }

      if (editingMapping) {
        await api.put(`/accounting/account-mappings/${editingMapping.transaction_type}`, {
          account_code: mappingForm.account_code,
          description: mappingForm.description || null
        })
        setAlert({ type: 'success', message: 'Configuración actualizada exitosamente' })
      } else {
        await api.post('/accounting/account-mappings', payload)
        setAlert({ type: 'success', message: 'Configuración creada exitosamente' })
      }
      
      closeMappingModal()
      fetchAccountMappings()
    } catch (error) {
      const errorMessage = error.response?.data?.error || error.message || 'Error al guardar configuración'
      setAlert({ type: 'error', message: errorMessage })
    } finally {
      setLoading(false)
    }
  }

  const transactionTypeLabels = {
    sales_cash: 'Ventas en Efectivo',
    sales_card: 'Ventas con Tarjeta',
    sales_transfer: 'Ventas con Transferencia',
    sales_credit: 'Ventas a Crédito (CXC)',
    sales_revenue: 'Ingresos por Ventas',
    sales_discount: 'Descuentos en Ventas',
    tax_payable: 'ITBIS por Pagar',
    tax_selective: 'Impuesto Selectivo por Pagar',
    expense_general: 'Gastos Generales',
    expense_utilities: 'Servicios Públicos',
    expense_rent: 'Alquileres',
    expense_supplies: 'Suministros',
    payment_cash: 'Pagos en Efectivo',
    payment_bank: 'Pagos Bancarios',
    purchase_accounts_payable: 'Cuentas por Pagar (Proveedores)',
    purchase_inventory: 'Compras de Inventario',
    purchase_expense: 'Costo de Compras',
    supplier_payment_cash: 'Pagos a Proveedores (Efectivo)',
    supplier_payment_bank: 'Pagos a Proveedores (Banco)'
  }

  return (
    <div className="max-w-7xl">
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-gray-900">Contabilidad</h1>
        <p className="mt-2 text-gray-600">
          Gestiona el plan de cuentas, asientos contables y reportes financieros
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

      <div className="bg-white rounded-lg shadow">
        {/* Tabs */}
        <div className="border-b border-gray-200">
          <nav className="flex -mb-px">
            {tabs.map((tab) => {
              const Icon = tab.icon
              return (
                <button
                  key={tab.id}
                  onClick={() => setActiveTab(tab.id)}
                  className={`
                    flex items-center px-6 py-4 border-b-2 font-medium text-sm
                    ${activeTab === tab.id
                      ? 'border-blue-500 text-blue-600'
                      : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                    }
                  `}
                >
                  <Icon className="w-5 h-5 mr-2" />
                  {tab.label}
                </button>
              )
            })}
          </nav>
        </div>

        {/* Content */}
        <div className="p-6">
          {loading && (
            <div className="flex justify-center items-center py-12">
              <div className="text-gray-600">Cargando...</div>
            </div>
          )}

          {/* Plan de Cuentas */}
          {activeTab === 'accounts' && !loading && (
            <div>
              <div className="flex justify-between items-center mb-6">
                <input
                  type="text"
                  placeholder="Buscar cuenta..."
                  value={accountFilter}
                  onChange={(e) => setAccountFilter(e.target.value)}
                  className="px-4 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                />
                <button
                  onClick={() => openAccountModal()}
                  className="flex items-center px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
                >
                  <Plus className="w-5 h-5 mr-2" />
                  Nueva Cuenta
                </button>
              </div>

              <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Código</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Nombre</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Tipo</th>
                      <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Balance</th>
                      <th className="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase">Estado</th>
                      <th className="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase">Acciones</th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {filteredAccounts.map((account) => (
                      <tr key={account.id} className="hover:bg-gray-50">
                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                          {account.account_code}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                          {account.account_name}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm">
                          <span className={`px-2 py-1 rounded-full text-xs font-medium ${getAccountTypeColor(account.account_type)}`}>
                            {accountTypes[account.account_type]}
                          </span>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">
                          DOP ${account.balance?.toLocaleString('es-DO', { minimumFractionDigits: 2 }) || '0.00'}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-center">
                          {account.is_active ? (
                            <span className="px-2 py-1 bg-green-100 text-green-800 rounded-full text-xs">Activa</span>
                          ) : (
                            <span className="px-2 py-1 bg-gray-100 text-gray-800 rounded-full text-xs">Inactiva</span>
                          )}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-center">
                          <div className="flex items-center justify-center gap-2">
                            <button
                              onClick={() => openAccountModal(account)}
                              className="p-1 text-blue-600 hover:text-blue-800 hover:bg-blue-50 rounded"
                              title="Editar cuenta"
                            >
                              <Edit className="w-4 h-4" />
                            </button>
                            <button
                              onClick={() => handleDeleteAccount(account)}
                              className="p-1 text-red-600 hover:text-red-800 hover:bg-red-50 rounded"
                              title="Eliminar cuenta"
                            >
                              <Trash2 className="w-4 h-4" />
                            </button>
                          </div>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>

              {filteredAccounts.length === 0 && (
                <div className="text-center py-12 text-gray-500">
                  No hay cuentas para mostrar
                </div>
              )}
            </div>
          )}

          {/* Asientos Contables */}
          {activeTab === 'entries' && !loading && (
            <div>
              <div className="flex justify-between items-center mb-6">
                <h2 className="text-lg font-semibold text-gray-900">Asientos Contables</h2>
                <button
                  onClick={() => setShowEntryModal(true)}
                  className="flex items-center px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
                >
                  <Plus className="w-5 h-5 mr-2" />
                  Nuevo Asiento
                </button>
              </div>

              <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase w-12"></th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Número</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Fecha</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Descripción</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Referencia</th>
                      <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Total Débito</th>
                      <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Total Crédito</th>
                      <th className="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase">Estado</th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {entries.map((entry) => {
                      const totals = calculateEntryTotals(entry)
                      const isExpanded = expandedEntryId === entry.id
                      return (
                        <>
                          <tr 
                            key={entry.id} 
                            className="hover:bg-gray-50 cursor-pointer"
                            onClick={() => toggleEntryExpansion(entry.id)}
                          >
                            <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                              {isExpanded ? '▼' : '▶'}
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                              {entry.entry_number}
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                              {new Date(entry.entry_date).toLocaleDateString('es-DO')}
                            </td>
                            <td className="px-6 py-4 text-sm text-gray-900">
                              {entry.description}
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                              {entry.reference || '-'}
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">
                              DOP ${totals.debit.toLocaleString('es-DO', { minimumFractionDigits: 2 })}
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">
                              DOP ${totals.credit.toLocaleString('es-DO', { minimumFractionDigits: 2 })}
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm text-center">
                              <span className={`px-2 py-1 rounded-full text-xs font-medium ${
                                entry.status === 'posted' ? 'bg-green-100 text-green-800' : 'bg-yellow-100 text-yellow-800'
                              }`}>
                                {entry.status === 'posted' ? 'Contabilizado' : 'Borrador'}
                              </span>
                            </td>
                          </tr>
                          {isExpanded && entry.lines && entry.lines.length > 0 && (
                            <tr>
                              <td colSpan="8" className="px-6 py-4 bg-gray-50">
                                <div className="ml-8">
                                  <h4 className="text-sm font-semibold text-gray-700 mb-3">Líneas del Asiento</h4>
                                  <div className="overflow-x-auto">
                                    <table className="min-w-full divide-y divide-gray-200 bg-white rounded-lg shadow-sm">
                                      <thead className="bg-gray-100">
                                        <tr>
                                          <th className="px-4 py-2 text-left text-xs font-medium text-gray-600 uppercase">Código</th>
                                          <th className="px-4 py-2 text-left text-xs font-medium text-gray-600 uppercase">Cuenta</th>
                                          <th className="px-4 py-2 text-left text-xs font-medium text-gray-600 uppercase">Descripción</th>
                                          <th className="px-4 py-2 text-right text-xs font-medium text-gray-600 uppercase">Débito</th>
                                          <th className="px-4 py-2 text-right text-xs font-medium text-gray-600 uppercase">Crédito</th>
                                        </tr>
                                      </thead>
                                      <tbody className="divide-y divide-gray-200">
                                        {entry.lines.map((line, idx) => (
                                          <tr key={idx} className="hover:bg-gray-50">
                                            <td className="px-4 py-2 whitespace-nowrap text-sm font-medium text-gray-900">
                                              {line.account_code}
                                            </td>
                                            <td className="px-4 py-2 text-sm text-gray-900">
                                              {line.account_name}
                                            </td>
                                            <td className="px-4 py-2 text-sm text-gray-500">
                                              {line.description || '-'}
                                            </td>
                                            <td className="px-4 py-2 whitespace-nowrap text-sm text-right text-gray-900">
                                              {line.debit > 0 ? `DOP $${line.debit.toLocaleString('es-DO', { minimumFractionDigits: 2 })}` : '-'}
                                            </td>
                                            <td className="px-4 py-2 whitespace-nowrap text-sm text-right text-gray-900">
                                              {line.credit > 0 ? `DOP $${line.credit.toLocaleString('es-DO', { minimumFractionDigits: 2 })}` : '-'}
                                            </td>
                                          </tr>
                                        ))}
                                        <tr className="bg-gray-100 font-semibold">
                                          <td colSpan="3" className="px-4 py-2 text-sm text-gray-900">TOTALES</td>
                                          <td className="px-4 py-2 whitespace-nowrap text-sm text-right text-gray-900">
                                            DOP ${totals.debit.toLocaleString('es-DO', { minimumFractionDigits: 2 })}
                                          </td>
                                          <td className="px-4 py-2 whitespace-nowrap text-sm text-right text-gray-900">
                                            DOP ${totals.credit.toLocaleString('es-DO', { minimumFractionDigits: 2 })}
                                          </td>
                                        </tr>
                                      </tbody>
                                    </table>
                                  </div>
                                </div>
                              </td>
                            </tr>
                          )}
                        </>
                      )
                    })}
                  </tbody>
                </table>
              </div>

              {entries.length === 0 && (
                <div className="text-center py-12 text-gray-500">
                  No hay asientos contables registrados
                </div>
              )}
            </div>
          )}

          {/* Configuración de Cuentas (Account Mappings) */}
          {activeTab === 'mappings' && !loading && (
            <div>
              <div className="flex justify-between items-center mb-6">
                <div>
                  <h2 className="text-lg font-semibold text-gray-900">Configuración de Cuentas</h2>
                  <p className="text-sm text-gray-600 mt-1">
                    Configure qué cuenta contable usar para cada tipo de transacción del sistema
                  </p>
                </div>
                <button
                  onClick={() => openMappingModal()}
                  className="flex items-center px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
                >
                  <Plus className="w-5 h-5 mr-2" />
                  Nueva Configuración
                </button>
              </div>

              <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Tipo de Transacción</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Código de Cuenta</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Nombre de Cuenta</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Descripción</th>
                      <th className="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase">Estado</th>
                      <th className="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase">Acciones</th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {accountMappings.map((mapping) => (
                      <tr key={mapping.id} className="hover:bg-gray-50">
                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                          {transactionTypeLabels[mapping.transaction_type] || mapping.transaction_type}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                          {mapping.account_code}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                          {mapping.account_name}
                        </td>
                        <td className="px-6 py-4 text-sm text-gray-500">
                          {mapping.description || '-'}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-center">
                          {mapping.is_active ? (
                            <span className="px-2 py-1 bg-green-100 text-green-800 rounded-full text-xs">Activa</span>
                          ) : (
                            <span className="px-2 py-1 bg-gray-100 text-gray-800 rounded-full text-xs">Inactiva</span>
                          )}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-center">
                          <button
                            onClick={() => openMappingModal(mapping)}
                            className="p-1 text-blue-600 hover:text-blue-800 hover:bg-blue-50 rounded"
                            title="Editar configuración"
                          >
                            <Edit className="w-4 h-4" />
                          </button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>

              {accountMappings.length === 0 && (
                <div className="text-center py-12 text-gray-500">
                  No hay configuraciones de cuentas. Crea una nueva configuración para comenzar.
                </div>
              )}
            </div>
          )}

          {/* Reportes */}
          {activeTab === 'reports' && !loading && (
            <div>
              <div className="mb-6 flex gap-4">
                <button
                  onClick={fetchTrialBalance}
                  className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
                >
                  Balance de Comprobación
                </button>
                <button
                  onClick={fetchBalanceSheet}
                  className="px-4 py-2 bg-purple-600 text-white rounded-md hover:bg-purple-700"
                >
                  Balance General
                </button>
                <button
                  onClick={fetchIncomeStatement}
                  className="px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700"
                >
                  Estado de Resultados
                </button>
              </div>

              {/* Trial Balance */}
              {trialBalance.length > 0 && !balanceSheet && !incomeStatement && (
                <div>
                  <h2 className="text-lg font-semibold text-gray-900 mb-4">Balance de Comprobación</h2>
                  <div className="overflow-x-auto">
                    <table className="min-w-full divide-y divide-gray-200">
                      <thead className="bg-gray-50">
                        <tr>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Código</th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Cuenta</th>
                          <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Débito</th>
                          <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Crédito</th>
                        </tr>
                      </thead>
                      <tbody className="bg-white divide-y divide-gray-200">
                        {trialBalance.map((entry, idx) => (
                          <tr key={idx} className="hover:bg-gray-50">
                            <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                              {entry.account_code}
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                              {entry.account_name}
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">
                              DOP ${entry.debit.toLocaleString('es-DO', { minimumFractionDigits: 2 })}
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">
                              DOP ${entry.credit.toLocaleString('es-DO', { minimumFractionDigits: 2 })}
                            </td>
                          </tr>
                        ))}
                        <tr className="bg-gray-100 font-semibold">
                          <td colSpan="2" className="px-6 py-4 text-sm text-gray-900">TOTALES</td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">
                            DOP ${calculateTrialBalanceTotals().debit.toLocaleString('es-DO', { minimumFractionDigits: 2 })}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">
                            DOP ${calculateTrialBalanceTotals().credit.toLocaleString('es-DO', { minimumFractionDigits: 2 })}
                          </td>
                        </tr>
                      </tbody>
                    </table>
                  </div>
                </div>
              )}

              {/* Balance Sheet */}
              {balanceSheet && (
                <div>
                  <div className="flex justify-between items-center mb-4">
                    <h2 className="text-lg font-semibold text-gray-900">{balanceSheet.statement_type}</h2>
                    <button
                      onClick={() => setBalanceSheet(null)}
                      className="text-sm text-blue-600 hover:underline"
                    >
                      Cerrar
                    </button>
                  </div>
                  <p className="text-sm text-gray-600 mb-6">Al {balanceSheet.period_end}</p>

                  {balanceSheet.sections?.map((section, idx) => (
                    <div key={idx} className="mb-6">
                      <h3 className="text-md font-semibold text-gray-800 mb-3">{section.title}</h3>
                      <table className="min-w-full">
                        <tbody>
                          {section.accounts.map((account, aidx) => (
                            <tr key={aidx}>
                              <td className="px-4 py-2 text-sm text-gray-600">{account.account_code}</td>
                              <td className="px-4 py-2 text-sm text-gray-900">{account.account_name}</td>
                              <td className="px-4 py-2 text-sm text-right text-gray-900">
                                DOP ${account.amount.toLocaleString('es-DO', { minimumFractionDigits: 2 })}
                              </td>
                            </tr>
                          ))}
                          <tr className="border-t border-gray-300 font-semibold">
                            <td colSpan="2" className="px-4 py-2 text-sm text-gray-900">Total {section.title}</td>
                            <td className="px-4 py-2 text-sm text-right text-gray-900">
                              DOP ${section.subtotal.toLocaleString('es-DO', { minimumFractionDigits: 2 })}
                            </td>
                          </tr>
                        </tbody>
                      </table>
                    </div>
                  ))}
                </div>
              )}

              {/* Income Statement */}
              {incomeStatement && (
                <div>
                  <div className="flex justify-between items-center mb-4">
                    <h2 className="text-lg font-semibold text-gray-900">{incomeStatement.statement_type}</h2>
                    <button
                      onClick={() => setIncomeStatement(null)}
                      className="text-sm text-blue-600 hover:underline"
                    >
                      Cerrar
                    </button>
                  </div>
                  <p className="text-sm text-gray-600 mb-6">
                    Del {incomeStatement.period_start} al {incomeStatement.period_end}
                  </p>

                  {incomeStatement.sections?.map((section, idx) => (
                    <div key={idx} className="mb-6">
                      <h3 className="text-md font-semibold text-gray-800 mb-3">{section.title}</h3>
                      <table className="min-w-full">
                        <tbody>
                          {section.accounts.map((account, aidx) => (
                            <tr key={aidx}>
                              <td className="px-4 py-2 text-sm text-gray-600">{account.account_code}</td>
                              <td className="px-4 py-2 text-sm text-gray-900">{account.account_name}</td>
                              <td className="px-4 py-2 text-sm text-right text-gray-900">
                                DOP ${account.amount.toLocaleString('es-DO', { minimumFractionDigits: 2 })}
                              </td>
                            </tr>
                          ))}
                          <tr className="border-t border-gray-300 font-semibold">
                            <td colSpan="2" className="px-4 py-2 text-sm text-gray-900">Total {section.title}</td>
                            <td className="px-4 py-2 text-sm text-right text-gray-900">
                              DOP ${section.subtotal.toLocaleString('es-DO', { minimumFractionDigits: 2 })}
                            </td>
                          </tr>
                        </tbody>
                      </table>
                    </div>
                  ))}

                  <div className="mt-6 pt-6 border-t-2 border-gray-400">
                    <div className="flex justify-between items-center">
                      <span className="text-lg font-bold text-gray-900">
                        {incomeStatement.total >= 0 ? 'Utilidad Neta' : 'Pérdida Neta'}
                      </span>
                      <span className={`text-lg font-bold ${incomeStatement.total >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                        DOP ${Math.abs(incomeStatement.total).toLocaleString('es-DO', { minimumFractionDigits: 2 })}
                      </span>
                    </div>
                  </div>
                </div>
              )}

              {trialBalance.length === 0 && !balanceSheet && !incomeStatement && (
                <div className="text-center py-12 text-gray-500">
                  Selecciona un reporte para visualizar
                </div>
              )}
            </div>
          )}
        </div>
      </div>

      {/* Account Modal */}
      {showAccountModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg shadow-xl max-w-2xl w-full mx-4 max-h-[90vh] overflow-y-auto">
            <div className="flex justify-between items-center p-6 border-b border-gray-200">
              <h2 className="text-xl font-semibold text-gray-900">
                {editingAccount ? 'Editar Cuenta' : 'Nueva Cuenta'}
              </h2>
              <button
                onClick={closeAccountModal}
                className="text-gray-400 hover:text-gray-600"
              >
                <X className="w-6 h-6" />
              </button>
            </div>

            <form onSubmit={handleAccountSubmit} className="p-6">
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Código de Cuenta *
                  </label>
                  <input
                    type="text"
                    value={accountForm.account_code}
                    onChange={(e) => setAccountForm({ ...accountForm, account_code: e.target.value })}
                    required
                    disabled={!!editingAccount}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500 disabled:bg-gray-100"
                    placeholder="Ej: 1111"
                  />
                  {editingAccount && (
                    <p className="mt-1 text-xs text-gray-500">El código de cuenta no se puede modificar</p>
                  )}
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Nombre de Cuenta *
                  </label>
                  <input
                    type="text"
                    value={accountForm.account_name}
                    onChange={(e) => setAccountForm({ ...accountForm, account_name: e.target.value })}
                    required
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                    placeholder="Ej: Caja General"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Tipo de Cuenta *
                  </label>
                  <select
                    value={accountForm.account_type}
                    onChange={(e) => setAccountForm({ ...accountForm, account_type: e.target.value })}
                    required
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                  >
                    <option value="ASSET">Activo</option>
                    <option value="LIABILITY">Pasivo</option>
                    <option value="EQUITY">Capital</option>
                    <option value="REVENUE">Ingreso</option>
                    <option value="EXPENSE">Gasto</option>
                  </select>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Cuenta Padre (Opcional)
                  </label>
                  <input
                    type="text"
                    value={accountForm.parent_account_code}
                    onChange={(e) => setAccountForm({ ...accountForm, parent_account_code: e.target.value })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                    placeholder="Ej: 1110 (dejar vacío si es cuenta principal)"
                  />
                  <p className="mt-1 text-xs text-gray-500">Código de la cuenta padre si es una subcuenta</p>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Descripción (Opcional)
                  </label>
                  <textarea
                    value={accountForm.description}
                    onChange={(e) => setAccountForm({ ...accountForm, description: e.target.value })}
                    rows={3}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                    placeholder="Descripción de la cuenta..."
                  />
                </div>

                {editingAccount && (
                  <div>
                    <label className="flex items-center">
                      <input
                        type="checkbox"
                        checked={accountForm.is_active}
                        onChange={(e) => setAccountForm({ ...accountForm, is_active: e.target.checked })}
                        className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                      />
                      <span className="ml-2 text-sm text-gray-700">Cuenta activa</span>
                    </label>
                  </div>
                )}
              </div>

              <div className="mt-6 flex justify-end gap-3">
                <button
                  type="button"
                  onClick={closeAccountModal}
                  className="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50"
                >
                  Cancelar
                </button>
                <button
                  type="submit"
                  disabled={loading}
                  className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50"
                >
                  {loading ? 'Guardando...' : editingAccount ? 'Actualizar' : 'Crear'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Account Mapping Modal */}
      {showMappingModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg shadow-xl max-w-2xl w-full mx-4 max-h-[90vh] overflow-y-auto">
            <div className="flex justify-between items-center p-6 border-b border-gray-200">
              <h2 className="text-xl font-semibold text-gray-900">
                {editingMapping ? 'Editar Configuración' : 'Nueva Configuración de Cuenta'}
              </h2>
              <button
                onClick={closeMappingModal}
                className="text-gray-400 hover:text-gray-600"
              >
                <X className="w-6 h-6" />
              </button>
            </div>

            <form onSubmit={handleMappingSubmit} className="p-6">
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Tipo de Transacción *
                  </label>
                  <select
                    value={mappingForm.transaction_type}
                    onChange={(e) => setMappingForm({ ...mappingForm, transaction_type: e.target.value })}
                    required
                    disabled={!!editingMapping}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500 disabled:bg-gray-100"
                  >
                    <option value="">Seleccione un tipo...</option>
                    <optgroup label="Ventas">
                      <option value="sales_cash">Ventas en Efectivo</option>
                      <option value="sales_card">Ventas con Tarjeta</option>
                      <option value="sales_transfer">Ventas con Transferencia</option>
                      <option value="sales_credit">Ventas a Crédito (CXC)</option>
                      <option value="sales_revenue">Ingresos por Ventas</option>
                      <option value="sales_discount">Descuentos en Ventas</option>
                    </optgroup>
                    <optgroup label="Impuestos">
                      <option value="tax_payable">ITBIS por Pagar</option>
                      <option value="tax_selective">Impuesto Selectivo por Pagar</option>
                    </optgroup>
                    <optgroup label="Gastos">
                      <option value="expense_general">Gastos Generales</option>
                      <option value="expense_utilities">Servicios Públicos</option>
                      <option value="expense_rent">Alquileres</option>
                      <option value="expense_supplies">Suministros</option>
                    </optgroup>
                    <optgroup label="Pagos">
                      <option value="payment_cash">Pagos en Efectivo</option>
                      <option value="payment_bank">Pagos Bancarios</option>
                    </optgroup>
                    <optgroup label="Compras">
                      <option value="purchase_accounts_payable">Cuentas por Pagar (Proveedores)</option>
                      <option value="purchase_inventory">Compras de Inventario</option>
                      <option value="purchase_expense">Costo de Compras</option>
                    </optgroup>
                    <optgroup label="Pagos a Proveedores">
                      <option value="supplier_payment_cash">Pagos a Proveedores (Efectivo)</option>
                      <option value="supplier_payment_bank">Pagos a Proveedores (Banco)</option>
                    </optgroup>
                  </select>
                  {editingMapping && (
                    <p className="mt-1 text-xs text-gray-500">El tipo de transacción no se puede modificar</p>
                  )}
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Código de Cuenta *
                  </label>
                  <input
                    type="text"
                    value={mappingForm.account_code}
                    onChange={(e) => setMappingForm({ ...mappingForm, account_code: e.target.value })}
                    required
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                    placeholder="Ej: 1111"
                    list="account-codes"
                  />
                  <datalist id="account-codes">
                    {accounts.map((acc) => (
                      <option key={acc.id} value={acc.account_code}>
                        {acc.account_name}
                      </option>
                    ))}
                  </datalist>
                  <p className="mt-1 text-xs text-gray-500">Ingrese el código de una cuenta existente en el plan de cuentas</p>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Descripción (Opcional)
                  </label>
                  <textarea
                    value={mappingForm.description}
                    onChange={(e) => setMappingForm({ ...mappingForm, description: e.target.value })}
                    rows={3}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                    placeholder="Descripción de esta configuración..."
                  />
                </div>
              </div>

              <div className="mt-6 flex justify-end gap-3">
                <button
                  type="button"
                  onClick={closeMappingModal}
                  className="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50"
                >
                  Cancelar
                </button>
                <button
                  type="submit"
                  disabled={loading}
                  className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50"
                >
                  {loading ? 'Guardando...' : editingMapping ? 'Actualizar' : 'Crear'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}

