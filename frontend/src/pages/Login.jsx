import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { useForm } from 'react-hook-form'
import { useAuthStore } from '../stores/authStore'
import api from '../lib/api'
import Logo from '../components/Logo'

export default function Login() {
  const navigate = useNavigate()
  const { setAuth, setTenant } = useAuthStore()
  const [error, setError] = useState('')
  const { register, handleSubmit, formState: { errors } } = useForm()

  const onSubmit = async (data) => {
    try {
      setError('')
      const response = await api.post('/auth/login', {
        email: data.email,
        password: data.password,
      })

      const { access_token, refresh_token, user } = response.data
      setAuth(user, access_token, refresh_token)
      navigate('/')
    } catch (err) {
      setError(err.response?.data?.error || 'Login failed')
    }
  }

  return (
    <div className="flex items-center justify-center min-h-screen bg-[#F5F5F5]">
      <div className="w-full max-w-md p-8 bg-white rounded-2xl shadow-card">
        {/* Logo y título */}
        <div className="text-center mb-8">
          <Logo variant="full" size="lg" className="justify-center mb-3" />
          <p className="text-[#212121] text-sm">Facturación Electrónica RD</p>
        </div>

        <h2 className="mb-6 text-2xl font-bold text-center text-[#212121]">
          Iniciar Sesión
        </h2>

        {error && (
          <div className="p-3 mb-4 text-sm text-red-700 bg-red-100 rounded-lg border border-red-300">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
          <div>
            <label className="block mb-2 text-sm font-medium text-[#212121]">
              Tenant ID
            </label>
            <input
              {...register('tenantId', { required: true })}
              defaultValue="demo"
              className="input-field"
            />
          </div>

          <div>
            <label className="block mb-2 text-sm font-medium text-[#212121]">
              Email
            </label>
            <input
              type="email"
              {...register('email', { required: 'El email es requerido' })}
              className="input-field"
              placeholder="tu@email.com"
            />
            {errors.email && (
              <p className="mt-1 text-sm text-red-600">{errors.email.message}</p>
            )}
          </div>

          <div>
            <label className="block mb-2 text-sm font-medium text-[#212121]">
              Contraseña
            </label>
            <input
              type="password"
              {...register('password', { required: 'La contraseña es requerida' })}
              className="input-field"
              placeholder="••••••••"
            />
            {errors.password && (
              <p className="mt-1 text-sm text-red-600">{errors.password.message}</p>
            )}
          </div>

          <button
            type="submit"
            className="btn-primary w-full"
          >
            Iniciar Sesión
          </button>
        </form>

        <p className="mt-6 text-sm text-center text-gray-600">
          ¿No tienes una cuenta?{' '}
          <Link to="/register" className="link">
            Registrarse
          </Link>
        </p>
      </div>
    </div>
  )
}

