import { Outlet, Link, useNavigate } from 'react-router-dom'
import { useAuthStore } from '../stores/authStore'
import { 
  LayoutDashboard, 
  ShoppingCart, 
  BarChart3,
  Users, 
  Package, 
  Receipt, 
  FileText, 
  BookOpen,
  Settings,
  User
} from 'lucide-react'

export default function Layout() {
  const navigate = useNavigate()
  const { user, logout } = useAuthStore()

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  const navigation = [
    { name: 'Dashboard', href: '/', icon: LayoutDashboard },
    { name: 'Ventas', href: '/sales', icon: ShoppingCart },
    { name: 'Historial Ventas', href: '/sales-management', icon: BarChart3 },
    { name: 'Clientes', href: '/customers', icon: Users },
    { name: 'Inventario', href: '/products', icon: Package },
    { name: 'Gastos', href: '/expenses', icon: Receipt },
    { name: 'Facturación', href: '/invoices', icon: FileText },
    { name: 'Contabilidad', href: '/accounting', icon: BookOpen },
    { name: 'Configuración', href: '/settings', icon: Settings },
  ]

  return (
    <div className="min-h-screen bg-white">
      {/* Sidebar */}
      <div className="fixed inset-y-0 left-0 w-64 bg-[#FF6B00] shadow-lg">
        <div className="flex flex-col h-full">
          {/* Logo */}
          <div className="flex items-center px-6 py-6">
            <div className="w-8 h-8 bg-white rounded-lg flex items-center justify-center mr-3">
              <span className="text-[#FF6B00] font-bold text-lg">V</span>
            </div>
            <span className="text-carbon font-bold text-2xl">Vendix</span>
          </div>
          
          {/* Navigation */}
          <nav className="flex-1 px-4 space-y-2">
            {navigation.map((item) => {
              const Icon = item.icon
              const isActive = window.location.pathname === item.href
              return (
                <Link
                  key={item.name}
                  to={item.href}
                  className={`flex items-center px-4 py-3 text-sm font-medium rounded-lg transition-colors duration-200 ${
                    isActive 
                      ? 'bg-white bg-opacity-20 text-white' 
                      : 'text-white hover:bg-white hover:bg-opacity-10'
                  }`}
                >
                  <Icon className="w-5 h-5 mr-3" />
                  {item.name}
                </Link>
              )
            })}
          </nav>

          {/* User profile icon at bottom */}
          <div className="p-4">
            <div className="flex items-center justify-center">
              <button className="p-2 text-white hover:bg-white hover:bg-opacity-10 rounded-lg transition-colors">
                <User className="w-6 h-6" />
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* Main content */}
      <div className="pl-64">
        <main className="p-8">
          <Outlet />
        </main>
      </div>
    </div>
  )
}

