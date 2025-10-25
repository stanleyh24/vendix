import { useState, useEffect } from 'react'
import api from '../lib/api'
import Alert from '../components/Alert'
import { BookOpen, FileText, BarChart, Plus, Eye } from 'lucide-react'

export default function Accounting() {
  const [activeTab, setActiveTab] = useState('accounts')
  const [alert, setAlert] = useState(null)
  const [loading, setLoading] = useState(false)

  // Accounts state
  const [accounts, setAccounts] = useState([])
  const [accountFilter, setAccountFilter] = useState('')
  const [showAccountModal, setShowAccountModal] = useState(false)

  // Journal entries state
  const [entries, setEntries] = useState([])
  const [showEntryModal, setShowEntryModal] = useState(false)

  // Reports state
  const [trialBalance, setTrialBalance] = useState([])
  const [balanceSheet, setBalanceSheet] = useState(null)
  const [incomeStatement, setIncomeStatement] = useState(null)

  const tabs = [
    { id: 'accounts', label: 'Plan de Cuentas', icon: BookOpen },
    { id: 'entries', label: 'Asientos Contables', icon: FileText },
    { id: 'reports', label: 'Reportes', icon: BarChart },
  ]

  useEffect(() => {
    if (activeTab === 'accounts') {
      fetchAccounts()
    } else if (activeTab === 'entries') {
      fetchEntries()
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
                  onClick={() => setShowAccountModal(true)}
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

              <div className="space-y-4">
                {entries.map((entry) => (
                  <div key={entry.id} className="border border-gray-200 rounded-lg p-4">
                    <div className="flex justify-between items-start mb-3">
                      <div>
                        <h3 className="font-semibold text-gray-900">{entry.entry_number}</h3>
                        <p className="text-sm text-gray-600">{new Date(entry.entry_date).toLocaleDateString('es-DO')}</p>
                        <p className="text-sm text-gray-700 mt-1">{entry.description}</p>
                      </div>
                      <span className={`px-3 py-1 rounded-full text-xs font-medium ${
                        entry.status === 'posted' ? 'bg-green-100 text-green-800' : 'bg-yellow-100 text-yellow-800'
                      }`}>
                        {entry.status === 'posted' ? 'Contabilizado' : 'Borrador'}
                      </span>
                    </div>

                    <div className="overflow-x-auto">
                      <table className="min-w-full text-sm">
                        <thead className="bg-gray-50">
                          <tr>
                            <th className="px-3 py-2 text-left text-xs font-medium text-gray-500">Cuenta</th>
                            <th className="px-3 py-2 text-right text-xs font-medium text-gray-500">Débito</th>
                            <th className="px-3 py-2 text-right text-xs font-medium text-gray-500">Crédito</th>
                          </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-100">
                          {entry.lines?.map((line, idx) => (
                            <tr key={idx}>
                              <td className="px-3 py-2">
                                <div className="font-medium text-gray-900">{line.account_code}</div>
                                <div className="text-gray-600">{line.account_name}</div>
                              </td>
                              <td className="px-3 py-2 text-right text-gray-900">
                                {line.debit > 0 ? `DOP $${line.debit.toLocaleString('es-DO', { minimumFractionDigits: 2 })}` : '-'}
                              </td>
                              <td className="px-3 py-2 text-right text-gray-900">
                                {line.credit > 0 ? `DOP $${line.credit.toLocaleString('es-DO', { minimumFractionDigits: 2 })}` : '-'}
                              </td>
                            </tr>
                          ))}
                        </tbody>
                      </table>
                    </div>
                  </div>
                ))}
              </div>

              {entries.length === 0 && (
                <div className="text-center py-12 text-gray-500">
                  No hay asientos contables registrados
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
    </div>
  )
}

