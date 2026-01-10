import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { 
  Bell, 
  X, 
  Check, 
  CheckCheck, 
  AlertCircle, 
  Clock, 
  ExternalLink,
  Filter,
  Search
} from 'lucide-react'
import api from '../lib/api'

export default function NotificationsPage() {
  const navigate = useNavigate()
  const [notifications, setNotifications] = useState([])
  const [unreadCount, setUnreadCount] = useState(0)
  const [loading, setLoading] = useState(true)
  const [filterType, setFilterType] = useState('all') // all, purchase_due_soon, purchase_overdue
  const [filterRead, setFilterRead] = useState('all') // all, read, unread
  const [searchTerm, setSearchTerm] = useState('')
  const [page, setPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const limit = 20

  useEffect(() => {
    loadNotifications()
    loadUnreadCount()
  }, [page, filterType, filterRead, searchTerm])

  const loadUnreadCount = async () => {
    try {
      const response = await api.get('/notifications/unread-count')
      setUnreadCount(response.data.count || 0)
    } catch (error) {
      console.error('Error loading unread count:', error)
    }
  }

  const loadNotifications = async () => {
    try {
      setLoading(true)
      const params = new URLSearchParams({
        limit: '1000', // Cargar más para poder filtrar en frontend
        offset: '0',
      })

      if (filterRead === 'unread') {
        params.append('unread_only', 'true')
      }

      if (filterType !== 'all') {
        params.append('type', filterType)
      }

      const response = await api.get(`/notifications?${params.toString()}`)
      let data = response.data || []

      // Filtrar por leídas si es necesario
      if (filterRead === 'read') {
        data = data.filter(n => n.is_read)
      }

      // Filtrar por búsqueda
      if (searchTerm) {
        data = data.filter(n => 
          n.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
          n.message.toLowerCase().includes(searchTerm.toLowerCase())
        )
      }

      // Calcular paginación
      const total = data.length
      setTotalPages(Math.ceil(total / limit))
      
      // Aplicar paginación
      const start = (page - 1) * limit
      const end = start + limit
      const paginatedData = data.slice(start, end)
      
      setNotifications(paginatedData)
    } catch (error) {
      console.error('Error loading notifications:', error)
    } finally {
      setLoading(false)
    }
  }

  const markAsRead = async (id) => {
    try {
      await api.post(`/notifications/${id}/read`)
      setNotifications(prev => 
        prev.map(n => n.id === id ? { ...n, is_read: true } : n)
      )
      setUnreadCount(prev => Math.max(0, prev - 1))
    } catch (error) {
      console.error('Error marking notification as read:', error)
    }
  }

  const markAllAsRead = async () => {
    try {
      await api.post('/notifications/read-all')
      setNotifications(prev => 
        prev.map(n => ({ ...n, is_read: true }))
      )
      setUnreadCount(0)
    } catch (error) {
      console.error('Error marking all as read:', error)
    }
  }

  const deleteNotification = async (id) => {
    try {
      await api.delete(`/notifications/${id}`)
      setNotifications(prev => prev.filter(n => n.id !== id))
      if (!notifications.find(n => n.id === id)?.is_read) {
        setUnreadCount(prev => Math.max(0, prev - 1))
      }
    } catch (error) {
      console.error('Error deleting notification:', error)
    }
  }

  const getNotificationIcon = (type, priority) => {
    if (type === 'purchase_overdue' || priority === 'urgent') {
      return <AlertCircle className="w-5 h-5 text-red-500" />
    }
    if (type === 'purchase_due_soon' || priority === 'high') {
      return <Clock className="w-5 h-5 text-orange-500" />
    }
    return <Bell className="w-5 h-5 text-gray-500" />
  }

  const getPriorityColor = (priority) => {
    switch (priority) {
      case 'urgent':
        return 'border-l-red-500 bg-red-50'
      case 'high':
        return 'border-l-orange-500 bg-orange-50'
      case 'normal':
        return 'border-l-blue-500 bg-blue-50'
      default:
        return 'border-l-gray-500 bg-gray-50'
    }
  }

  const getPriorityBadge = (priority) => {
    const badges = {
      urgent: { label: 'Urgente', color: 'bg-red-100 text-red-800' },
      high: { label: 'Alta', color: 'bg-orange-100 text-orange-800' },
      normal: { label: 'Normal', color: 'bg-blue-100 text-blue-800' },
      low: { label: 'Baja', color: 'bg-gray-100 text-gray-800' },
    }
    const badge = badges[priority] || badges.normal
    return (
      <span className={`px-2 py-1 rounded-full text-xs font-medium ${badge.color}`}>
        {badge.label}
      </span>
    )
  }

  const formatDate = (dateString) => {
    const date = new Date(dateString)
    const now = new Date()
    const diffMs = now - date
    const diffMins = Math.floor(diffMs / 60000)
    const diffHours = Math.floor(diffMs / 3600000)
    const diffDays = Math.floor(diffMs / 86400000)

    if (diffMins < 1) return 'Hace un momento'
    if (diffMins < 60) return `Hace ${diffMins} min`
    if (diffHours < 24) return `Hace ${diffHours} h`
    if (diffDays < 7) return `Hace ${diffDays} días`
    return date.toLocaleDateString('es-DO', { 
      day: 'numeric', 
      month: 'short',
      year: date.getFullYear() !== now.getFullYear() ? 'numeric' : undefined
    })
  }

  const handleNotificationClick = (notification) => {
    if (notification.entity_type === 'purchase' && notification.entity_id) {
      navigate(`/purchases?highlight=${notification.entity_id}`)
    }
  }

  const filteredNotifications = notifications.filter(n => {
    if (searchTerm) {
      const search = searchTerm.toLowerCase()
      return n.title.toLowerCase().includes(search) || 
             n.message.toLowerCase().includes(search)
    }
    return true
  })

  return (
    <div className="p-6">
      {/* Header */}
      <div className="mb-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-[#212121] flex items-center gap-3">
              <Bell className="w-8 h-8 text-[#FF6B00]" />
              Notificaciones
            </h1>
            <p className="text-gray-600 mt-1">
              {unreadCount > 0 
                ? `${unreadCount} notificación${unreadCount !== 1 ? 'es' : ''} sin leer`
                : 'Todas las notificaciones están leídas'
              }
            </p>
          </div>
          {unreadCount > 0 && (
            <button
              onClick={markAllAsRead}
              className="px-4 py-2 bg-[#FF6B00] text-white rounded-lg hover:bg-[#E55A00] transition-colors flex items-center gap-2"
            >
              <CheckCheck className="w-5 h-5" />
              Marcar todas como leídas
            </button>
          )}
        </div>
      </div>

      {/* Filtros y búsqueda */}
      <div className="bg-white rounded-lg shadow-sm p-4 mb-6">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
            <input
              type="text"
              placeholder="Buscar notificaciones..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#FF6B00] focus:border-transparent"
            />
          </div>
          <select
            value={filterType}
            onChange={(e) => {
              setFilterType(e.target.value)
              setPage(1)
            }}
            className="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#FF6B00] focus:border-transparent"
          >
            <option value="all">Todos los tipos</option>
            <option value="purchase_due_soon">Compras próximas a vencer</option>
            <option value="purchase_overdue">Compras vencidas</option>
          </select>
          <select
            value={filterRead}
            onChange={(e) => {
              setFilterRead(e.target.value)
              setPage(1)
            }}
            className="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#FF6B00] focus:border-transparent"
          >
            <option value="all">Todas</option>
            <option value="unread">Sin leer</option>
            <option value="read">Leídas</option>
          </select>
          <div className="flex items-center gap-2 text-sm text-gray-600">
            <Filter className="w-5 h-5" />
            <span>{filteredNotifications.length} resultado{filteredNotifications.length !== 1 ? 's' : ''}</span>
          </div>
        </div>
      </div>

      {/* Lista de notificaciones */}
      <div className="bg-white rounded-lg shadow-sm overflow-hidden">
        {loading ? (
          <div className="p-8 text-center text-gray-500">
            Cargando notificaciones...
          </div>
        ) : filteredNotifications.length === 0 ? (
          <div className="p-12 text-center">
            <Bell className="w-16 h-16 mx-auto mb-4 text-gray-300" />
            <p className="text-gray-500 text-lg font-medium">No hay notificaciones</p>
            <p className="text-gray-400 text-sm mt-2">
              {searchTerm || filterType !== 'all' || filterRead !== 'all'
                ? 'Intenta ajustar los filtros'
                : 'No tienes notificaciones en este momento'}
            </p>
          </div>
        ) : (
          <>
            <div className="divide-y divide-gray-200">
              {filteredNotifications.map((notification) => {
                const isPurchaseNotification = notification.entity_type === 'purchase' && notification.entity_id
                
                return (
                  <div
                    key={notification.id}
                    className={`p-6 border-l-4 transition-all ${
                      notification.is_read 
                        ? 'bg-white' 
                        : getPriorityColor(notification.priority)
                    } ${isPurchaseNotification ? 'cursor-pointer hover:bg-gray-50' : ''}`}
                    onClick={isPurchaseNotification ? () => handleNotificationClick(notification) : undefined}
                  >
                    <div className="flex items-start gap-4">
                      <div className="flex-shrink-0 mt-0.5">
                        {getNotificationIcon(notification.type, notification.priority)}
                      </div>
                      <div className="flex-1 min-w-0">
                        <div className="flex items-start justify-between gap-4">
                          <div className="flex-1">
                            <div className="flex items-center gap-3 mb-2">
                              <h3 className={`text-base font-semibold ${
                                notification.is_read ? 'text-gray-700' : 'text-gray-900'
                              }`}>
                                {notification.title}
                              </h3>
                              {!notification.is_read && (
                                <span className="w-2 h-2 bg-[#FF6B00] rounded-full"></span>
                              )}
                              {getPriorityBadge(notification.priority)}
                              {isPurchaseNotification && (
                                <ExternalLink className="w-4 h-4 text-[#FF6B00] flex-shrink-0" />
                              )}
                            </div>
                            <p className="text-sm text-gray-600 mb-2">
                              {notification.message}
                            </p>
                            <div className="flex items-center gap-4 text-xs text-gray-500">
                              {notification.days_until_due !== undefined && (
                                <span className="text-orange-600 font-medium">
                                  Vence en {notification.days_until_due} días
                                </span>
                              )}
                              {notification.days_overdue !== undefined && (
                                <span className="text-red-600 font-medium">
                                  Vencida hace {notification.days_overdue} días
                                </span>
                              )}
                              <span>{formatDate(notification.created_at)}</span>
                              {notification.entity_type && (
                                <span className="px-2 py-1 bg-gray-100 rounded text-gray-600">
                                  {notification.entity_type === 'purchase' ? 'Compra' : notification.entity_type}
                                </span>
                              )}
                            </div>
                          </div>
                          <div className="flex items-center gap-2">
                            {!notification.is_read && (
                              <button
                                onClick={(e) => {
                                  e.stopPropagation()
                                  markAsRead(notification.id)
                                }}
                                className="p-2 hover:bg-gray-200 rounded-lg transition-colors"
                                title="Marcar como leída"
                              >
                                <Check className="w-5 h-5 text-gray-500" />
                              </button>
                            )}
                            <button
                              onClick={(e) => {
                                e.stopPropagation()
                                deleteNotification(notification.id)
                              }}
                              className="p-2 hover:bg-gray-200 rounded-lg transition-colors"
                              title="Eliminar"
                            >
                              <X className="w-5 h-5 text-gray-400" />
                            </button>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                )
              })}
            </div>

            {/* Paginación */}
            {totalPages > 1 && (
              <div className="px-6 py-4 border-t border-gray-200 bg-gray-50 flex items-center justify-between">
                <div className="text-sm text-gray-600">
                  Página {page} de {totalPages}
                </div>
                <div className="flex items-center gap-2">
                  <button
                    onClick={() => setPage(prev => Math.max(1, prev - 1))}
                    disabled={page === 1}
                    className="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    Anterior
                  </button>
                  <button
                    onClick={() => setPage(prev => Math.min(totalPages, prev + 1))}
                    disabled={page === totalPages}
                    className="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    Siguiente
                  </button>
                </div>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  )
}

