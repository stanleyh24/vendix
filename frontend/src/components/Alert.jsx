import React from 'react';
import { CheckCircle, AlertCircle, Info, AlertTriangle, X } from 'lucide-react';

/**
 * Componente de Alerta para feedback visual
 * @param {string} type - Tipo de alerta: 'success', 'error', 'info', 'warning'
 * @param {string} title - Título de la alerta (opcional)
 * @param {string} message - Mensaje de la alerta
 * @param {function} onClose - Función para cerrar la alerta (opcional)
 */
export default function Alert({ type = 'info', title, message, onClose, className = '' }) {
  const icons = {
    success: CheckCircle,
    error: AlertCircle,
    info: Info,
    warning: AlertTriangle,
  };

  const Icon = icons[type] || Info;

  const classNames = {
    success: 'alert-success',
    error: 'alert-error',
    info: 'alert-info',
    warning: 'alert-warning',
  };

  return (
    <div className={`${classNames[type]} ${className}`}>
      <Icon className="w-5 h-5 flex-shrink-0 mt-0.5" />
      <div className="flex-1">
        {title && <h4 className="font-semibold mb-1">{title}</h4>}
        <p className="text-sm">{message}</p>
      </div>
      {onClose && (
        <button
          onClick={onClose}
          className="flex-shrink-0 hover:opacity-70 transition-opacity"
          aria-label="Cerrar alerta"
        >
          <X className="w-5 h-5" />
        </button>
      )}
    </div>
  );
}

/**
 * Componente de Toast para notificaciones temporales
 */
export function Toast({ type = 'info', message, duration = 3000, onClose }) {
  React.useEffect(() => {
    if (duration && onClose) {
      const timer = setTimeout(onClose, duration);
      return () => clearTimeout(timer);
    }
  }, [duration, onClose]);

  return (
    <div className="fixed top-4 right-4 z-50 animate-slide-in">
      <Alert type={type} message={message} onClose={onClose} />
    </div>
  );
}

