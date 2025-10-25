import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { useForm } from 'react-hook-form'
import api from '../lib/api'

export default function Register() {
  const navigate = useNavigate()
  const [error, setError] = useState('')
  const [success, setSuccess] = useState(false)
  const { register, handleSubmit, formState: { errors }, watch } = useForm()

  const onSubmit = async (data) => {
    try {
      setError('')
      await api.post('/auth/register', {
        email: data.email,
        password: data.password,
        first_name: data.firstName,
        last_name: data.lastName,
        role: 'admin',
      })

      setSuccess(true)
      setTimeout(() => navigate('/login'), 2000)
    } catch (err) {
      setError(err.response?.data?.error || 'Registration failed')
    }
  }

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-100">
      <div className="w-full max-w-md p-8 bg-white rounded-lg shadow-md">
        <h2 className="mb-6 text-2xl font-bold text-center text-gray-900">
          Create Account
        </h2>

        {error && (
          <div className="p-3 mb-4 text-sm text-red-700 bg-red-100 rounded">
            {error}
          </div>
        )}

        {success && (
          <div className="p-3 mb-4 text-sm text-green-700 bg-green-100 rounded">
            Account created! Redirecting to login...
          </div>
        )}

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block mb-2 text-sm font-medium text-gray-700">
                First Name
              </label>
              <input
                {...register('firstName')}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
            <div>
              <label className="block mb-2 text-sm font-medium text-gray-700">
                Last Name
              </label>
              <input
                {...register('lastName')}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
          </div>

          <div>
            <label className="block mb-2 text-sm font-medium text-gray-700">
              Email
            </label>
            <input
              type="email"
              {...register('email', { required: 'Email is required' })}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <div>
            <label className="block mb-2 text-sm font-medium text-gray-700">
              Password
            </label>
            <input
              type="password"
              {...register('password', { 
                required: 'Password is required',
                minLength: { value: 8, message: 'Password must be at least 8 characters' }
              })}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            {errors.password && (
              <p className="mt-1 text-sm text-red-600">{errors.password.message}</p>
            )}
          </div>

          <button
            type="submit"
            className="w-full px-4 py-2 text-white bg-blue-600 rounded-md hover:bg-blue-700"
          >
            Create Account
          </button>
        </form>

        <p className="mt-4 text-sm text-center text-gray-600">
          Already have an account?{' '}
          <Link to="/login" className="text-blue-600 hover:underline">
            Sign in
          </Link>
        </p>
      </div>
    </div>
  )
}

