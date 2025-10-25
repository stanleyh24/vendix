/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        // Paleta de colores Vendix
        vendix: {
          orange: '#FF6B00',      // Color principal (energía, ventas, dinamismo)
          'orange-dark': '#E66000', // Variante más oscura para hover
          'orange-light': '#FF8533', // Variante más clara
        },
        carbon: '#212121',         // Negro carbón para tipografía y contraste
        neutral: {
          light: '#E6E6E6',        // Gris claro para bordes y elementos neutros
          lighter: '#F5F5F5',      // Gris más claro para fondos secundarios
        },
        saas: {
          blue: '#1877F2',         // Azul SaaS secundario (confianza tecnológica)
          'blue-dark': '#1565D8',  // Variante más oscura para hover
        },
        // Colores de feedback visual
        feedback: {
          success: '#00C853',      // Verde para acciones exitosas
          'success-dark': '#00A344',
          error: '#D32F2F',        // Rojo para errores
          'error-dark': '#B71C1C',
          info: '#1877F2',         // Azul para información
          'info-dark': '#1565D8',
          warning: '#FFA000',      // Amarillo/Naranja para advertencias
          'warning-dark': '#FF8F00',
        },
      },
      // Bordes redondeados personalizados
      borderRadius: {
        'xl': '1rem',
        '2xl': '1.5rem',
        '3xl': '2rem',
      },
      // Sombras personalizadas
      boxShadow: {
        'soft': '0 2px 8px rgba(0, 0, 0, 0.08)',
        'card': '0 4px 12px rgba(0, 0, 0, 0.1)',
        'card-hover': '0 8px 24px rgba(0, 0, 0, 0.15)',
      },
      // Asegurar que los colores personalizados estén disponibles
      backgroundColor: (theme) => theme('colors'),
      textColor: (theme) => theme('colors'),
      borderColor: (theme) => theme('colors'),
    },
  },
  plugins: [],
}

