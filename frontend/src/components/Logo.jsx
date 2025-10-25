import React from 'react';

/**
 * Logo de Vendix
 * "V" en naranja, resto en negro, con ícono de venta en la "i"
 */
export default function Logo({ variant = 'full', size = 'md', className = '' }) {
  const sizes = {
    sm: { text: 'text-xl', icon: '12' },
    md: { text: 'text-2xl', icon: '16' },
    lg: { text: 'text-4xl', icon: '24' },
    xl: { text: 'text-6xl', icon: '32' },
  };

  const currentSize = sizes[size] || sizes.md;

  // Logo completo con texto
  if (variant === 'full') {
    return (
      <div className={`flex items-center ${className}`}>
        <span className={`font-bold ${currentSize.text}`}>
          <span style={{ color: '#FF6B00' }}>V</span>
          <span style={{ color: '#212121' }}>end</span>
          <span style={{ color: '#212121' }} className="relative inline-block">
            i
            <svg
              width={currentSize.icon}
              height={currentSize.icon}
              viewBox="0 0 24 24"
              fill="none"
              className="absolute -top-1 left-1/2 -translate-x-1/2"
              style={{ marginTop: '-2px' }}
            >
              <path
                d="M9 2L7 9H3L6 16H7L8 22H16L17 16H18L21 9H17L15 2H9Z"
                fill="#FF6B00"
                opacity="0.7"
              />
            </svg>
          </span>
          <span style={{ color: '#212121' }}>x</span>
        </span>
      </div>
    );
  }

  // Logo compacto (solo "V" en círculo naranja)
  if (variant === 'icon') {
    const iconSizes = {
      sm: 'w-8 h-8 text-base',
      md: 'w-10 h-10 text-lg',
      lg: 'w-14 h-14 text-2xl',
      xl: 'w-20 h-20 text-4xl',
    };

    return (
      <div
        className={`${iconSizes[size]} rounded-full flex items-center justify-center font-bold text-white ${className}`}
        style={{ backgroundColor: '#FF6B00' }}
      >
        V
      </div>
    );
  }

  // Logo para sidebar (V + texto compacto)
  if (variant === 'sidebar') {
    return (
      <div className={`flex items-center gap-2 ${className}`}>
        <div
          className="w-8 h-8 rounded-full flex items-center justify-center font-bold text-white text-sm"
          style={{ backgroundColor: '#FF6B00' }}
        >
          V
        </div>
        <span className="font-bold text-lg text-white">Vendix</span>
      </div>
    );
  }

  return null;
}

/**
 * SVG del logo para exportar o usar en otros contextos
 */
export function LogoSVG({ width = 200, height = 60 }) {
  return (
    <svg width={width} height={height} viewBox="0 0 200 60" fill="none" xmlns="http://www.w3.org/2000/svg">
      <text x="10" y="45" fontFamily="system-ui, -apple-system, sans-serif" fontWeight="bold" fontSize="42">
        <tspan fill="#FF6B00">V</tspan>
        <tspan fill="#212121">end</tspan>
        <tspan fill="#212121">i</tspan>
        <tspan fill="#212121">x</tspan>
      </text>
      {/* Ícono de venta sobre la "i" */}
      <path
        d="M108 12L107 18H104L106 23H106.5L107 28H111L111.5 23H112L115 18H112L111 12H108Z"
        fill="#FF6B00"
        opacity="0.7"
      />
    </svg>
  );
}

