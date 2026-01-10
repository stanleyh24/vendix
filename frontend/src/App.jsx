import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom'
import { useAuthStore } from './stores/authStore'
import Login from './pages/Login'
import Register from './pages/Register'
import Dashboard from './pages/Dashboard'
import Pos from './pages/Pos'
import SalesManagement from './pages/SalesManagement'
import Customers from './pages/Customers'
import Products from './pages/Products'
import Expenses from './pages/Expenses'
import Invoices from './pages/Invoices'
import CreditNotes from './pages/CreditNotes'
import Suppliers from './pages/Suppliers'
import Purchases from './pages/Purchases'
import Accounting from './pages/Accounting'
import Reports from './pages/Reports'
import Returns from './pages/Returns'
import Settings from './pages/Settings'
import StyleGuide from './pages/StyleGuide'
import Employees from './pages/Employees'
import Payroll from './pages/Payroll'
import NotificationsPage from './pages/NotificationsPage'
import Layout from './components/Layout'

function PrivateRoute({ children }) {
  const { isAuthenticated } = useAuthStore()
  return isAuthenticated ? children : <Navigate to="/login" />
}

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        
        <Route path="/" element={<PrivateRoute><Layout /></PrivateRoute>}>
          <Route index element={<Dashboard />} />
          <Route path="pos" element={<Pos />} />
          {/* Legacy route: redirect /sales to /pos */}
          <Route path="sales" element={<Navigate to="/pos" replace />} />
          <Route path="sales-management" element={<SalesManagement />} />
          <Route path="customers" element={<Customers />} />
          <Route path="products" element={<Products />} />
          <Route path="suppliers" element={<Suppliers />} />
          <Route path="purchases" element={<Purchases />} />
          <Route path="expenses" element={<Expenses />} />
          <Route path="invoices" element={<Invoices />} />
          <Route path="credit-notes" element={<CreditNotes />} />
          <Route path="accounting" element={<Accounting />} />
          <Route path="reports" element={<Reports />} />
          <Route path="returns" element={<Returns />} />
          <Route path="employees" element={<Employees />} />
          <Route path="payroll" element={<Payroll />} />
          <Route path="notifications" element={<NotificationsPage />} />
          <Route path="settings" element={<Settings />} />
          <Route path="style-guide" element={<StyleGuide />} />
        </Route>
      </Routes>
    </Router>
  )
}

export default App

