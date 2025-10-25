import React from 'react';
import Logo from '../components/Logo';
import Alert from '../components/Alert';

/**
 * Guía de Estilos de Vendix
 * Componente para visualizar y probar la paleta de colores y componentes UI
 */
export default function StyleGuide() {
  return (
    <div className="min-h-screen p-8" style={{ backgroundColor: '#F5F5F5' }}>
      <div className="max-w-7xl mx-auto space-y-8">
        {/* Header */}
        <div className="card">
          <Logo variant="full" size="xl" className="justify-center mb-4" />
          <h1 className="text-4xl font-bold mb-2 text-center" style={{ color: '#FF6B00' }}>
            Vendix Style Guide
          </h1>
          <p className="text-center" style={{ color: '#212121' }}>
            Paleta de colores y componentes de la plataforma de facturación electrónica
          </p>
        </div>

        {/* Logo Variants */}
        <div className="card">
          <h2 className="text-2xl font-bold mb-6" style={{ color: '#212121' }}>🎨 Logo Vendix</h2>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
            {/* Logo completo */}
            <div className="space-y-4">
              <h3 className="font-semibold text-lg">Logo Completo</h3>
              <div className="p-6 bg-white rounded-lg border border-gray-200 flex items-center justify-center">
                <Logo variant="full" size="sm" />
              </div>
              <div className="p-6 bg-white rounded-lg border border-gray-200 flex items-center justify-center">
                <Logo variant="full" size="md" />
              </div>
              <div className="p-6 bg-white rounded-lg border border-gray-200 flex items-center justify-center">
                <Logo variant="full" size="lg" />
              </div>
            </div>

            {/* Logo icono */}
            <div className="space-y-4">
              <h3 className="font-semibold text-lg">Logo Icono (Favicon)</h3>
              <div className="p-6 bg-white rounded-lg border border-gray-200 flex items-center justify-center gap-4">
                <Logo variant="icon" size="sm" />
                <Logo variant="icon" size="md" />
                <Logo variant="icon" size="lg" />
              </div>
              
              <h3 className="font-semibold text-lg mt-8">Logo Sidebar</h3>
              <div className="p-6 bg-[#212121] rounded-lg flex items-center justify-center">
                <Logo variant="sidebar" />
              </div>
            </div>
          </div>
        </div>

        {/* Paleta de Colores */}
        <div className="card">
          <h2 className="text-2xl font-bold mb-6" style={{ color: '#212121' }}>🎨 Paleta de Colores Principal</h2>
          
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {/* Naranja Vendix */}
            <div>
              <div className="h-32 rounded-lg mb-3 flex items-center justify-center" style={{ backgroundColor: '#FF6B00' }}>
                <span className="text-white font-semibold">Naranja Vendix</span>
              </div>
              <p className="text-sm font-mono" style={{ color: '#212121' }}>#FF6B00</p>
              <p className="text-xs text-gray-600">Color principal</p>
            </div>

            {/* Naranja Oscuro */}
            <div>
              <div className="h-32 rounded-lg mb-3 flex items-center justify-center" style={{ backgroundColor: '#E66000' }}>
                <span className="text-white font-semibold">Naranja Oscuro</span>
              </div>
              <p className="text-sm font-mono" style={{ color: '#212121' }}>#E66000</p>
              <p className="text-xs text-gray-600">Hover / Active</p>
            </div>

            {/* Naranja Claro */}
            <div>
              <div className="h-32 rounded-lg mb-3 flex items-center justify-center" style={{ backgroundColor: '#FF8533' }}>
                <span className="text-white font-semibold">Naranja Claro</span>
              </div>
              <p className="text-sm font-mono" style={{ color: '#212121' }}>#FF8533</p>
              <p className="text-xs text-gray-600">Badges / Destacados</p>
            </div>

            {/* Azul SaaS */}
            <div>
              <div className="h-32 rounded-lg mb-3 flex items-center justify-center" style={{ backgroundColor: '#1877F2' }}>
                <span className="text-white font-semibold">Azul SaaS</span>
              </div>
              <p className="text-sm font-mono" style={{ color: '#212121' }}>#1877F2</p>
              <p className="text-xs text-gray-600">Color secundario</p>
            </div>

            {/* Negro Carbón */}
            <div>
              <div className="h-32 rounded-lg mb-3 flex items-center justify-center" style={{ backgroundColor: '#212121' }}>
                <span className="text-white font-semibold">Negro Carbón</span>
              </div>
              <p className="text-sm font-mono" style={{ color: '#212121' }}>#212121</p>
              <p className="text-xs text-gray-600">Tipografía / Contraste</p>
            </div>

            {/* Gris Claro */}
            <div>
              <div className="h-32 rounded-lg mb-3 flex items-center justify-center border border-gray-300" style={{ backgroundColor: '#E6E6E6' }}>
                <span className="font-semibold" style={{ color: '#212121' }}>Gris Claro</span>
              </div>
              <p className="text-sm font-mono" style={{ color: '#212121' }}>#E6E6E6</p>
              <p className="text-xs text-gray-600">Bordes / Neutros</p>
            </div>
          </div>

          <h3 className="text-xl font-bold mb-4 mt-8" style={{ color: '#212121' }}>Colores de Feedback Visual</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            {/* Verde Éxito */}
            <div>
              <div className="h-24 rounded-lg mb-3 flex items-center justify-center" style={{ backgroundColor: '#00C853' }}>
                <span className="text-white font-semibold">Éxito</span>
              </div>
              <p className="text-sm font-mono" style={{ color: '#212121' }}>#00C853</p>
              <p className="text-xs text-gray-600">Acciones exitosas</p>
            </div>

            {/* Rojo Error */}
            <div>
              <div className="h-24 rounded-lg mb-3 flex items-center justify-center" style={{ backgroundColor: '#D32F2F' }}>
                <span className="text-white font-semibold">Error</span>
              </div>
              <p className="text-sm font-mono" style={{ color: '#212121' }}>#D32F2F</p>
              <p className="text-xs text-gray-600">Errores y alertas</p>
            </div>

            {/* Azul Info */}
            <div>
              <div className="h-24 rounded-lg mb-3 flex items-center justify-center" style={{ backgroundColor: '#1877F2' }}>
                <span className="text-white font-semibold">Info</span>
              </div>
              <p className="text-sm font-mono" style={{ color: '#212121' }}>#1877F2</p>
              <p className="text-xs text-gray-600">Información</p>
            </div>

            {/* Amarillo Advertencia */}
            <div>
              <div className="h-24 rounded-lg mb-3 flex items-center justify-center" style={{ backgroundColor: '#FFA000' }}>
                <span className="text-white font-semibold">Advertencia</span>
              </div>
              <p className="text-sm font-mono" style={{ color: '#212121' }}>#FFA000</p>
              <p className="text-xs text-gray-600">Advertencias</p>
            </div>
          </div>
        </div>

        {/* Alertas y Feedback */}
        <div className="card">
          <h2 className="text-2xl font-bold mb-6" style={{ color: '#212121' }}>✨ Alertas y Feedback Visual</h2>
          
          <div className="space-y-4">
            <Alert 
              type="success" 
              title="Éxito" 
              message="La factura #12345 se ha generado correctamente y enviado a DGII."
            />
            
            <Alert 
              type="error" 
              title="Error" 
              message="No se pudo procesar el pago. Por favor, verifica los datos de la tarjeta."
            />
            
            <Alert 
              type="info" 
              title="Información" 
              message="Recuerda que tienes 15 facturas pendientes de envío a DGII."
            />
            
            <Alert 
              type="warning" 
              title="Advertencia" 
              message="Tu suscripción vence en 3 días. Renueva ahora para evitar la interrupción del servicio."
            />
          </div>
        </div>

        {/* Botones */}
        <div className="card">
          <h2 className="text-2xl font-bold mb-6" style={{ color: '#212121' }}>🔘 Botones</h2>
          
          <div className="space-y-4">
            <div className="flex flex-wrap gap-4 items-center">
              <button className="btn-primary">Botón Principal</button>
              <button className="btn-primary" disabled>
                Botón Deshabilitado
              </button>
              <button className="btn-primary text-sm px-3 py-1">
                Botón Pequeño
              </button>
              <button className="btn-primary text-lg px-6 py-3">
                Botón Grande
              </button>
            </div>

            <div className="flex flex-wrap gap-4 items-center">
              <button className="btn-secondary">Botón Secundario</button>
              <button className="btn-outline">Botón Outline</button>
            </div>

            <h3 className="text-lg font-semibold mt-6 mb-3">Botones de Feedback</h3>
            <div className="flex flex-wrap gap-4 items-center">
              <button className="btn-success">Guardar Cambios</button>
              <button className="btn-error">Eliminar</button>
              <button className="btn-warning">Advertencia</button>
            </div>
          </div>
        </div>

        {/* Inputs */}
        <div className="card">
          <h2 className="text-2xl font-bold text-carbon mb-6">📝 Inputs</h2>
          
          <div className="max-w-md space-y-4">
            <div>
              <label className="block text-sm font-medium text-carbon mb-2">
                Nombre del cliente
              </label>
              <input 
                type="text" 
                className="input-field" 
                placeholder="Ingrese el nombre"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-carbon mb-2">
                Email
              </label>
              <input 
                type="email" 
                className="input-field" 
                placeholder="cliente@example.com"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-carbon mb-2">
                Descripción
              </label>
              <textarea 
                className="input-field" 
                rows="3"
                placeholder="Ingrese una descripción"
              />
            </div>
          </div>
        </div>

        {/* Badges */}
        <div className="card">
          <h2 className="text-2xl font-bold mb-6" style={{ color: '#212121' }}>🏷️ Badges</h2>
          
          <h3 className="text-lg font-semibold mb-3">Badges de Color Vendix</h3>
          <div className="flex flex-wrap gap-3 mb-6">
            <span className="badge-orange">Pendiente</span>
            <span className="badge-blue">Enviada</span>
          </div>

          <h3 className="text-lg font-semibold mb-3">Badges de Estado</h3>
          <div className="flex flex-wrap gap-3">
            <span className="badge-success">Completada</span>
            <span className="badge-error">Rechazada</span>
            <span className="badge-warning">En Proceso</span>
            <span className="badge-info">Información</span>
            <span className="badge bg-purple-500 text-white">Borrador</span>
            <span className="badge bg-gray-500 text-white">Cancelada</span>
          </div>
        </div>

        {/* Cards */}
        <div className="card">
          <h2 className="text-2xl font-bold text-carbon mb-6">📄 Cards</h2>
          
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            {/* Card de estadística */}
            <div className="card bg-white">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-500">Facturas del mes</p>
                  <p className="text-3xl font-bold text-carbon">247</p>
                </div>
                <div className="w-12 h-12 bg-vendix-orange-light rounded-full flex items-center justify-center">
                  <span className="text-white text-xl">📊</span>
                </div>
              </div>
            </div>

            <div className="card bg-white">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-500">Ingresos totales</p>
                  <p className="text-3xl font-bold text-carbon">$125K</p>
                </div>
                <div className="w-12 h-12 bg-saas-blue rounded-full flex items-center justify-center">
                  <span className="text-white text-xl">💰</span>
                </div>
              </div>
            </div>

            <div className="card bg-white">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-500">Clientes activos</p>
                  <p className="text-3xl font-bold text-carbon">1,234</p>
                </div>
                <div className="w-12 h-12 bg-green-500 rounded-full flex items-center justify-center">
                  <span className="text-white text-xl">👥</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Links */}
        <div className="card">
          <h2 className="text-2xl font-bold text-carbon mb-6">🔗 Links</h2>
          
          <div className="space-y-2">
            <div>
              <a href="#" className="link">Enlace normal</a>
            </div>
            <div>
              <a href="#" className="text-saas-blue hover:text-saas-blue-dark underline">
                Enlace azul
              </a>
            </div>
            <div>
              <a href="#" className="text-carbon hover:text-vendix-orange transition-colors">
                Enlace sin subrayado
              </a>
            </div>
          </div>
        </div>

        {/* Tipografía */}
        <div className="card">
          <h2 className="text-2xl font-bold text-carbon mb-6">📝 Tipografía</h2>
          
          <div className="space-y-4">
            <h1 className="text-4xl font-bold text-carbon">Heading 1</h1>
            <h2 className="text-3xl font-bold text-carbon">Heading 2</h2>
            <h3 className="text-2xl font-bold text-carbon">Heading 3</h3>
            <h4 className="text-xl font-bold text-carbon">Heading 4</h4>
            <h5 className="text-lg font-bold text-carbon">Heading 5</h5>
            <h6 className="text-base font-bold text-carbon">Heading 6</h6>
            
            <p className="text-base text-carbon mt-4">
              Párrafo normal. Lorem ipsum dolor sit amet, consectetur adipiscing elit. 
              Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
            </p>
            
            <p className="text-sm text-gray-600">
              Texto secundario más pequeño. Ut enim ad minim veniam, quis nostrud 
              exercitation ullamco laboris.
            </p>
            
            <p className="text-xs text-gray-500">
              Texto muy pequeño para notas o información complementaria.
            </p>
          </div>
        </div>

        {/* Ejemplo de Formulario Completo */}
        <div className="card">
          <h2 className="text-2xl font-bold text-carbon mb-6">📋 Formulario Completo</h2>
          
          <form className="max-w-2xl space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-carbon mb-2">
                  Nombre *
                </label>
                <input 
                  type="text" 
                  className="input-field" 
                  placeholder="Juan Pérez"
                  required
                />
              </div>
              
              <div>
                <label className="block text-sm font-medium text-carbon mb-2">
                  RNC/Cédula *
                </label>
                <input 
                  type="text" 
                  className="input-field" 
                  placeholder="000-0000000-0"
                  required
                />
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-carbon mb-2">
                Email *
              </label>
              <input 
                type="email" 
                className="input-field" 
                placeholder="cliente@example.com"
                required
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-carbon mb-2">
                Dirección
              </label>
              <textarea 
                className="input-field" 
                rows="3"
                placeholder="Calle Principal #123, Santo Domingo"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-carbon mb-2">
                Plan
              </label>
              <select className="input-field">
                <option>Seleccione un plan</option>
                <option>Básico</option>
                <option>Profesional</option>
                <option>Empresarial</option>
              </select>
            </div>

            <div className="flex items-center gap-2">
              <input 
                type="checkbox" 
                id="terms" 
                className="w-4 h-4 text-vendix-orange border-neutral-light rounded focus:ring-vendix-orange"
              />
              <label htmlFor="terms" className="text-sm text-carbon">
                Acepto los términos y condiciones
              </label>
            </div>

            <div className="flex gap-3 pt-4">
              <button type="submit" className="btn-primary">
                Guardar Cliente
              </button>
              <button type="button" className="btn-outline">
                Cancelar
              </button>
              <button type="button" className="btn-secondary">
                Vista Previa
              </button>
            </div>
          </form>
        </div>

      </div>
    </div>
  );
}

