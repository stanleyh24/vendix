import { User } from 'lucide-react'
import SalesChart from '../components/SalesChart'

export default function Dashboard() {
  // Datos de ejemplo que coinciden con la imagen
  const kpiCards = [
    {
      title: 'Ventas',
      value: '$18.000',
      subtitle: 'Nueva venta',
      subtitleValue: '$2.700'
    },
    {
      title: 'Ingresos',
      value: '$24,500',
      subtitle: '↗ 6%',
      subtitleValue: '',
      isPositive: true
    },
    {
      title: 'Facturación',
      value: '8',
      subtitle: 'Facturas',
      subtitleValue: 'Ver',
      hasButton: true
    }
  ]


  return (
    <div className="min-h-screen bg-white">
      {/* Header con icono de usuario */}
      <div className="flex justify-end p-6">
        <button className="p-2 text-gray-600 hover:bg-gray-100 rounded-lg transition-colors">
          <User className="w-6 h-6" />
        </button>
      </div>

      {/* KPI Cards */}
      <div className="px-6 pb-6">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {kpiCards.map((card, index) => (
            <div key={index} className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
              <h3 className="text-sm font-medium text-gray-600 mb-2">{card.title}</h3>
              <div className="text-3xl font-bold text-gray-900 mb-2">{card.value}</div>
              <div className="flex items-center justify-between">
                <div className="flex items-center">
                  <span className="text-sm text-gray-600">{card.subtitle}</span>
                  {card.subtitleValue && (
                    <span className={`text-sm ml-2 ${card.isPositive ? 'text-green-600' : 'text-gray-900'}`}>
                      {card.subtitleValue}
                    </span>
                  )}
                </div>
                {card.hasButton && (
                  <button className="bg-[#FF6B00] text-white text-xs px-3 py-1 rounded-lg hover:bg-[#E66000] transition-colors">
                    {card.subtitleValue}
                  </button>
                )}
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Gráfico de ventas */}
      <div className="px-6">
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-6">Últimas ventas</h2>
          <SalesChart />
        </div>
      </div>
    </div>
  )
}

