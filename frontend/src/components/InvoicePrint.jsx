import { useState, useEffect, useRef } from 'react';
import { X, Printer, Download, FileText } from 'lucide-react';
import api from '../lib/api';
import Alert from './Alert';

/**
 * Componente para visualizar e imprimir facturas
 * @param {string} invoiceId - ID de la factura a imprimir
 * @param {function} onClose - Función para cerrar el modal
 */
export default function InvoicePrint({ invoiceId, onClose }) {
  const [invoice, setInvoice] = useState(null);
  const [loading, setLoading] = useState(true);
  const [alert, setAlert] = useState(null);
  const printRef = useRef();

  useEffect(() => {
    if (invoiceId) {
      loadInvoice();
    }
  }, [invoiceId]);

  const loadInvoice = async () => {
    try {
      setLoading(true);
      // Include customer information in the request
      const response = await api.get(`/invoices/${invoiceId}?include=customer`);
      setInvoice(response.data);
    } catch (error) {
      showAlert('error', 'Error', 'No se pudo cargar la factura');
      console.error('Error loading invoice:', error);
    } finally {
      setLoading(false);
    }
  };

  const showAlert = (type, title, message) => {
    setAlert({ type, title, message });
    setTimeout(() => setAlert(null), 5000);
  };

  const handlePrint = () => {
    window.print();
  };

  const formatDate = (dateString) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('es-DO', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit'
    });
  };

  const formatCurrency = (amount) => {
    return new Intl.NumberFormat('es-DO', {
      style: 'currency',
      currency: 'DOP'
    }).format(amount);
  };

  if (loading) {
    return (
      <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div className="bg-white rounded-2xl p-8 max-w-4xl w-full mx-4">
          <div className="flex items-center justify-center">
            <div className="w-8 h-8 border-4 border-[#FF6B00] border-t-transparent rounded-full animate-spin"></div>
            <span className="ml-3 text-gray-600">Cargando factura...</span>
          </div>
        </div>
      </div>
    );
  }

  if (!invoice) {
    return (
      <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div className="bg-white rounded-2xl p-8 max-w-4xl w-full mx-4">
          <Alert type="error" title="Error" message="No se pudo cargar la factura" />
          <button onClick={onClose} className="btn-primary mt-4">
            Cerrar
          </button>
        </div>
      </div>
    );
  }

  return (
    <>
      {/* Estilos para impresión POS (80mm) */}
      <style>{`
        @media print {
          body * {
            visibility: hidden;
          }
          #invoice-print-content,
          #invoice-print-content * {
            visibility: visible;
          }
          #invoice-print-content {
            position: absolute;
            left: 0;
            top: 0;
            width: 80mm !important;
            max-width: 80mm !important;
            padding: 5mm !important;
            font-size: 10pt !important;
          }
          .no-print {
            display: none !important;
          }
          .print-break {
            page-break-after: always;
          }
          @page {
            size: 80mm auto;
            margin: 0;
          }
        }
      `}</style>

      {/* Modal */}
      <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4 overflow-y-auto">
        <div className="bg-white rounded-2xl shadow-2xl max-w-4xl w-full my-8">
          {/* Header - No se imprime */}
          <div className="no-print sticky top-0 bg-white border-b border-gray-200 px-6 py-4 rounded-t-2xl flex items-center justify-between">
            <h2 className="text-2xl font-bold text-[#212121] flex items-center gap-2">
              <FileText className="w-6 h-6 text-[#FF6B00]" />
              Vista Previa de Factura
            </h2>
            <div className="flex items-center gap-2">
              <button
                onClick={handlePrint}
                className="btn-primary flex items-center gap-2"
              >
                <Printer className="w-4 h-4" />
                Imprimir
              </button>
              <button
                onClick={onClose}
                className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
              >
                <X className="w-5 h-5" />
              </button>
            </div>
          </div>

          {/* Alert */}
          {alert && (
            <div className="no-print px-6 pt-4">
              <Alert
                type={alert.type}
                title={alert.title}
                message={alert.message}
                onClose={() => setAlert(null)}
              />
            </div>
          )}

          {/* Contenido de la factura - Formato POS (80mm) */}
          <div id="invoice-print-content" className="max-w-[80mm] mx-auto bg-white">
            {/* Header del recibo POS */}
            <div className="text-center border-b-2 border-dashed border-gray-300 pb-3 mb-3">
              <div className="text-xl font-bold text-[#212121] mb-1">VENDIX</div>
              <div className="text-xs text-gray-600">Sistema de Facturación</div>
              <div className="text-xs text-gray-600">República Dominicana</div>
            </div>

            {/* Información básica */}
            <div className="text-center border-b-2 border-dashed border-gray-300 pb-3 mb-3">
              <div className="text-sm font-bold uppercase mb-1">Recibo de Venta</div>
              <div className="text-xs">#{invoice.invoice_number}</div>
              {invoice.ncf && (
                <div className="text-xs">NCF: {invoice.ncf}</div>
              )}
              <div className="text-xs mt-1">
                {invoice.ncf_type === '01' ? 'Crédito Fiscal' : 
                 invoice.ncf_type === '15' ? 'Gubernamental' : 
                 'Consumidor Final'}
              </div>
              {invoice.ncf_type === '15' && (
                <div className="text-xs text-blue-600 font-semibold mt-1">🏛️ Factura Gubernamental</div>
              )}
            </div>

            {/* Cliente compacto */}
            <div className="border-b-2 border-dashed border-gray-300 pb-2 mb-2">
              <div className="text-xs font-bold uppercase">Cliente:</div>
              <div className="text-xs">
                {invoice.customer && typeof invoice.customer === 'object' 
                  ? (invoice.customer.name || invoice.customer_name || 'Cliente Genérico')
                  : (invoice.customer_name || 'Cliente Genérico')}
              </div>
              {invoice.customer && typeof invoice.customer === 'object' && invoice.customer.tax_id && (
                <div className="text-xs text-gray-600">RNC: {invoice.customer.tax_id}</div>
              )}
            </div>

            {/* Fecha compacta */}
            <div className="border-b-2 border-dashed border-gray-300 pb-2 mb-2 text-xs">
              <div>Fecha: {formatDate(invoice.issue_date)}</div>
              {invoice.due_date !== invoice.issue_date && (
                <div>Vence: {formatDate(invoice.due_date)}</div>
              )}
            </div>

            {/* Items - Formato compacto para POS */}
            <div className="border-b-2 border-dashed border-gray-300 pb-2 mb-2">
              {invoice.lines && invoice.lines.map((line) => (
                <div key={line.id} className="mb-2 pb-2 border-b border-gray-200 last:border-0">
                  <div className="flex justify-between items-start text-xs mb-1">
                    <div className="flex-1">
                      <div className="font-semibold">{line.description}</div>
                      <div className="text-gray-600">
                        {line.quantity} x {formatCurrency(line.unit_price)}
                      </div>
                    </div>
                    <div className="text-right font-bold">
                      {formatCurrency(line.line_total)}
                    </div>
                  </div>
                </div>
              ))}
            </div>

            {/* Totales compactos */}
            <div className="border-b-2 border-dashed border-gray-300 pb-2 mb-2 space-y-1">
              <div className="flex justify-between text-xs">
                <span>Subtotal:</span>
                <span>{formatCurrency(invoice.subtotal)}</span>
              </div>
              <div className="flex justify-between text-xs">
                <span>ITBIS:</span>
                <span>{formatCurrency(invoice.tax_amount)}</span>
              </div>
              <div className="flex justify-between font-bold text-sm border-t border-gray-300 pt-1">
                <span>TOTAL:</span>
                <span>{formatCurrency(invoice.total)}</span>
              </div>
              {/* Mostrar retención si existe */}
              {invoice.withholding_tax_amount && invoice.withholding_tax_amount > 0 && (
                <>
                  <div className="flex justify-between text-xs pt-1 border-t border-gray-300">
                    <span className="flex items-center gap-1">
                      Retención ISR
                      {invoice.withholding_rate && (
                        <span className="text-gray-500">
                          ({(invoice.withholding_rate * 100).toFixed(2)}%)
                        </span>
                      )}
                    </span>
                    <span className="text-red-600 font-semibold">
                      -{formatCurrency(invoice.withholding_tax_amount)}
                    </span>
                  </div>
                  <div className="flex justify-between font-bold text-sm border-t-2 border-gray-400 pt-1">
                    <span>MONTO NETO A RECIBIR:</span>
                    <span className="text-green-600">
                      {formatCurrency(invoice.net_amount || (invoice.total - invoice.withholding_tax_amount))}
                    </span>
                  </div>
                </>
              )}
            </div>

            {/* Notas compactas si existen */}
            {invoice.notes && (
              <div className="border-b-2 border-dashed border-gray-300 pb-2 mb-2 text-xs text-gray-600">
                <div className="font-bold mb-1">NOTA:</div>
                <div>{invoice.notes}</div>
              </div>
            )}

            {/* Footer compacto */}
            <div className="text-center border-t-2 border-dashed border-gray-300 pt-3 mt-3">
              <div className="text-xs text-gray-500">
                ¡Gracias por su compra!
              </div>
              <div className="text-[10px] text-gray-400 mt-1">
                {new Date().toLocaleString('es-DO')}
              </div>
              <div className="text-[10px] text-gray-400 mt-2">
                Documento generado electrónicamente
              </div>
            </div>
          </div>

          {/* Footer - No se imprime */}
          <div className="no-print border-t border-gray-200 px-6 py-4 bg-gray-50 rounded-b-2xl flex items-center justify-between">
            <p className="text-sm text-gray-600">
              Puedes imprimir esta factura o guardarla como PDF
            </p>
            <button
              onClick={onClose}
              className="btn-outline"
            >
              Cerrar
            </button>
          </div>
        </div>
      </div>
    </>
  );
}

