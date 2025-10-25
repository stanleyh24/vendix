import React from 'react';

export default function SalesChart() {
  // Datos del gráfico que coinciden con la imagen
  const salesData = [
    { month: 'Jan', value: 12 },
    { month: 'Feb', value: 19 },
    { month: 'Mar', value: 15 },
    { month: 'Apr', value: 25 },
    { month: 'May', value: 18 },
    { month: 'Jun', value: 22 },
    { month: 'Jul', value: 20 },
    { month: 'Ago', value: 16 },
    { month: 'Sep', value: 14 },
    { month: 'Nov', value: 18 },
    { month: 'Dec', value: 21 }
  ];

  const maxValue = Math.max(...salesData.map(d => d.value));
  const chartHeight = 200;

  return (
    <div className="relative">
      {/* Líneas de la grilla de fondo */}
      <div className="absolute inset-0 pointer-events-none">
        {[0, 25, 50, 75, 100].map((percent, index) => (
          <div
            key={index}
            className="absolute w-full border-t border-gray-100"
            style={{ bottom: `${percent}%` }}
          ></div>
        ))}
      </div>

      {/* Gráfico de barras */}
      <div className="h-64 flex items-end justify-between space-x-1 relative z-10">
        {salesData.map((data, index) => {
          const barHeight = (data.value / maxValue) * chartHeight;
          return (
            <div key={index} className="flex flex-col items-center flex-1 group">
              {/* Barra */}
              <div 
                className="bg-[#FF6B00] rounded-t w-full mb-2 transition-all duration-300 hover:bg-[#E66000] cursor-pointer relative group"
                style={{ 
                  height: `${barHeight}px`,
                  minHeight: '4px'
                }}
              >
                {/* Tooltip en hover */}
                <div className="absolute -top-8 left-1/2 transform -translate-x-1/2 bg-gray-800 text-white text-xs px-2 py-1 rounded opacity-0 group-hover:opacity-100 transition-opacity duration-200 whitespace-nowrap">
                  {data.value} ventas
                </div>
              </div>
              
              {/* Etiqueta del mes */}
              <span className="text-xs text-gray-600 font-medium">{data.month}</span>
            </div>
          );
        })}
      </div>
    </div>
  );
}
