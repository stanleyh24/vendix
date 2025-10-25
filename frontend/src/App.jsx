import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom'
import { useAuthStore } from './stores/authStore'
import Login from './pages/Login'
import Register from './pages/Register'
import Dashboard from './pages/Dashboard'
import Sales from './pages/Sales'
import SalesManagement from './pages/SalesManagement'
import Customers from './pages/Customers'
import Products from './pages/Products'
import Expenses from './pages/Expenses'
import Invoices from './pages/Invoices'
import Accounting from './pages/Accounting'
import Reports from './pages/Reports'
import Settings from './pages/Settings'
import StyleGuide from './pages/StyleGuide'
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
          <Route path="sales" element={<Sales />} />
          <Route path="sales-management" element={<SalesManagement />} />
          <Route path="customers" element={<Customers />} />
          <Route path="products" element={<Products />} />
          <Route path="expenses" element={<Expenses />} />
          <Route path="invoices" element={<Invoices />} />
          <Route path="accounting" element={<Accounting />} />
          <Route path="settings" element={<Settings />} />
          <Route path="style-guide" element={<StyleGuide />} />
        </Route>
      </Routes>
    </Router>
  )
}

export default App

