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
      const response = await api.get(`/invoices/${invoiceId}`);
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
      month: 'long',
      day: 'numeric'
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
      {/* Estilos para impresión */}
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
            width: 100%;
          }
          .no-print {
            display: none !important;
          }
          .print-break {
            page-break-after: always;
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

          {/* Contenido de la factura */}
          <div id="invoice-print-content" className="p-8 bg-white">
            {/* Header de la factura */}
            <div className="border-b-4 border-[#FF6B00] pb-6 mb-6">
              <div className="flex items-start justify-between">
                <div>
                  <h1 className="text-4xl font-bold text-[#212121] mb-2">FACTURA</h1>
                  <p className="text-gray-600 text-lg">{invoice.invoice_number}</p>
                  {invoice.ncf && (
                    <p className="text-sm text-gray-500 mt-1">NCF: {invoice.ncf}</p>
                  )}
                  <p className="text-sm text-gray-500">
                    Tipo: {invoice.ncf_type === '01' ? 'Crédito Fiscal' : 'Consumidor Final'}
                  </p>
                </div>
                <div className="text-right">
                  <div className="bg-[#FF6B00] text-white px-4 py-2 rounded-lg inline-block mb-4">
                    <span className="text-sm font-medium">VENDIX</span>
                  </div>
                  <p className="text-sm text-gray-600">Sistema de Facturación</p>
                  <p className="text-sm text-gray-600">República Dominicana</p>
                </div>
              </div>
            </div>

            {/* Información del cliente y fechas */}
            <div className="grid grid-cols-2 gap-8 mb-8">
              <div>
                <h3 className="text-sm font-semibold text-gray-500 uppercase mb-2">
                  Información del Cliente
                </h3>
                <div className="bg-gray-50 rounded-lg p-4">
                  <p className="font-bold text-lg text-[#212121] mb-1">
                    {invoice.customer_name || 'Cliente Genérico'}
                  </p>
                  {invoice.customer_id && (
                    <p className="text-sm text-gray-600">
                      ID: {invoice.customer_id}
                    </p>
                  )}
                </div>
              </div>
              <div>
                <h3 className="text-sm font-semibold text-gray-500 uppercase mb-2">
                  Información de la Factura
                </h3>
                <div className="bg-gray-50 rounded-lg p-4 space-y-2">
                  <div className="flex justify-between">
                    <span className="text-sm text-gray-600">Fecha de Emisión:</span>
                    <span className="text-sm font-semibold text-[#212121]">
                      {formatDate(invoice.issue_date)}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-sm text-gray-600">Fecha de Vencimiento:</span>
                    <span className="text-sm font-semibold text-[#212121]">
                      {formatDate(invoice.due_date)}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-sm text-gray-600">Estado:</span>
                    <span className={`text-sm font-semibold ${
                      invoice.status === 'draft' ? 'text-yellow-600' :
                      invoice.status === 'sent' ? 'text-blue-600' :
                      invoice.status === 'paid' ? 'text-green-600' :
                      'text-gray-600'
                    }`}>
                      {invoice.status === 'draft' ? 'Borrador' :
                       invoice.status === 'sent' ? 'Enviada' :
                       invoice.status === 'paid' ? 'Pagada' :
                       invoice.status}
                    </span>
                  </div>
                </div>
              </div>
            </div>

            {/* Tabla de productos/servicios */}
            <div className="mb-8">
              <table className="w-full border-collapse">
                <thead>
                  <tr className="bg-[#212121] text-white">
                    <th className="text-left py-3 px-4 text-sm font-semibold">#</th>
                    <th className="text-left py-3 px-4 text-sm font-semibold">Descripción</th>
                    <th className="text-right py-3 px-4 text-sm font-semibold">Cantidad</th>
                    <th className="text-right py-3 px-4 text-sm font-semibold">Precio Unit.</th>
                    <th className="text-right py-3 px-4 text-sm font-semibold">ITBIS</th>
                    <th className="text-right py-3 px-4 text-sm font-semibold">Total</th>
                  </tr>
                </thead>
                <tbody>
                  {invoice.lines && invoice.lines.map((line, index) => (
                    <tr key={line.id} className="border-b border-gray-200">
                      <td className="py-3 px-4 text-sm text-gray-600">
                        {line.line_number}
                      </td>
                      <td className="py-3 px-4">
                        <p className="font-medium text-[#212121]">{line.description}</p>
                      </td>
                      <td className="py-3 px-4 text-right text-sm text-gray-600">
                        {line.quantity}
                      </td>
                      <td className="py-3 px-4 text-right text-sm text-gray-600">
                        {formatCurrency(line.unit_price)}
                      </td>
                      <td className="py-3 px-4 text-right text-sm text-gray-600">
                        {formatCurrency(line.tax_amount)}
                      </td>
                      <td className="py-3 px-4 text-right font-semibold text-[#212121]">
                        {formatCurrency(line.line_total)}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            {/* Totales */}
            <div className="flex justify-end mb-8">
              <div className="w-80">
                <div className="bg-gray-50 rounded-lg p-4 space-y-3">
                  <div className="flex justify-between text-sm">
                    <span className="text-gray-600">Subtotal:</span>
                    <span className="font-semibold text-[#212121]">
                      {formatCurrency(invoice.subtotal)}
                    </span>
                  </div>
                  <div className="flex justify-between text-sm">
                    <span className="text-gray-600">ITBIS (18%):</span>
                    <span className="font-semibold text-[#212121]">
                      {formatCurrency(invoice.tax_amount)}
                    </span>
                  </div>
                  <div className="border-t-2 border-gray-300 pt-3 flex justify-between">
                    <span className="font-bold text-lg text-[#212121]">TOTAL:</span>
                    <span className="font-bold text-2xl text-[#FF6B00]">
                      {formatCurrency(invoice.total)}
                    </span>
                  </div>
                  {invoice.paid_amount > 0 && (
                    <>
                      <div className="flex justify-between text-sm">
                        <span className="text-gray-600">Pagado:</span>
                        <span className="font-semibold text-green-600">
                          {formatCurrency(invoice.paid_amount)}
                        </span>
                      </div>
                      <div className="flex justify-between text-sm">
                        <span className="text-gray-600">Balance Pendiente:</span>
                        <span className="font-semibold text-red-600">
                          {formatCurrency(invoice.total - invoice.paid_amount)}
                        </span>
                      </div>
                    </>
                  )}
                </div>
              </div>
            </div>

            {/* Notas y términos */}
            {(invoice.notes || invoice.terms) && (
              <div className="border-t border-gray-200 pt-6 space-y-4">
                {invoice.notes && (
                  <div>
                    <h3 className="text-sm font-semibold text-gray-500 uppercase mb-2">
                      Notas
                    </h3>
                    <p className="text-sm text-gray-600">{invoice.notes}</p>
                  </div>
                )}
                {invoice.terms && (
                  <div>
                    <h3 className="text-sm font-semibold text-gray-500 uppercase mb-2">
                      Términos y Condiciones
                    </h3>
                    <p className="text-sm text-gray-600">{invoice.terms}</p>
                  </div>
                )}
              </div>
            )}

            {/* Footer */}
            <div className="border-t border-gray-200 pt-6 mt-8">
              <p className="text-center text-xs text-gray-500">
                Este documento fue generado electrónicamente por VENDIX
              </p>
              <p className="text-center text-xs text-gray-500 mt-1">
                Factura generada el {new Date().toLocaleString('es-DO')}
              </p>
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

