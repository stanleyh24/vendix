# 🖨️ Guía Rápida: Impresión de Facturas

## ¿Cómo imprimir una factura después de una venta?

### Paso 1: Completar la Venta
1. Agrega productos al carrito en el **Punto de Venta**
2. Haz clic en **"Completar Venta"**
3. Selecciona el cliente (o usa Cliente Genérico)
4. Elige el tipo de comprobante fiscal:
   - **NCF 01**: Crédito Fiscal (para empresas con RNC)
   - **NCF 02**: Consumidor Final (para público general)
5. Selecciona el método de pago
6. Haz clic en **"Confirmar Venta"**

### Paso 2: Modal de Éxito
Después de confirmar la venta, verás un modal de éxito con dos opciones:

#### 🖨️ Imprimir Factura
- Abre una vista previa completa de la factura
- Puedes revisarla antes de imprimir
- Haz clic en el botón **"Imprimir"** para abrir el diálogo de impresión

#### 🛒 Nueva Venta
- Limpia el carrito y prepara el sistema para la siguiente venta
- Úsala si no necesitas imprimir en este momento

### Paso 3: Vista Previa e Impresión
En la vista previa podrás ver:
- ✅ Número de factura
- ✅ NCF (si aplica)
- ✅ Datos del cliente
- ✅ Lista detallada de productos
- ✅ Subtotal, ITBIS y Total
- ✅ Fechas de emisión y vencimiento

**Opciones disponibles:**
1. **Imprimir**: Abre el diálogo de impresión del navegador
2. **Cerrar**: Vuelve al punto de venta

### Paso 4: Diálogo de Impresión del Navegador
Desde aquí puedes:
- 🖨️ Imprimir en una impresora física
- 📄 Guardar como PDF
- ⚙️ Ajustar opciones de impresión (orientación, márgenes, etc.)

## Consejos para una Mejor Impresión

### Configuración Recomendada:
- **Orientación**: Vertical (Portrait)
- **Tamaño de papel**: Carta (8.5" x 11") o A4
- **Márgenes**: Predeterminados o mínimos
- **Escala**: 100% (sin ajustar)
- **Fondo**: Deshabilitado (opcional)

### Para Guardar como PDF:
1. En el diálogo de impresión, selecciona **"Guardar como PDF"** como destino
2. Elige la ubicación y nombre del archivo
3. Haz clic en **"Guardar"**

## Atajos de Teclado

| Atajo | Acción |
|-------|--------|
| `Ctrl + P` (Windows/Linux) | Imprimir |
| `Cmd + P` (Mac) | Imprimir |
| `Esc` | Cerrar vista previa |

## Elementos de la Factura Impresa

### Encabezado:
- Logo/Nombre del sistema (VENDIX)
- Número de factura
- NCF y tipo de comprobante

### Información del Cliente:
- Nombre del cliente
- RNC (si aplica)
- Tipo de cliente

### Detalles de la Factura:
- Fecha de emisión
- Fecha de vencimiento
- Estado de la factura

### Tabla de Productos:
| # | Descripción | Cantidad | Precio Unit. | ITBIS | Total |
|---|-------------|----------|--------------|-------|-------|
| 1 | Producto X  | 2        | $50.00       | $18.00| $118.00|

### Totales:
- **Subtotal**: Suma de productos sin impuestos
- **ITBIS (18%)**: Impuesto sobre el subtotal
- **TOTAL**: Monto final a pagar

### Pie de Página:
- Notas adicionales (si existen)
- Términos y condiciones (si aplican)
- Fecha y hora de generación

## Solución de Problemas

### La factura no se carga:
1. Verifica tu conexión a internet
2. Actualiza la página (F5)
3. Intenta de nuevo

### El formato se ve mal al imprimir:
1. Verifica que estés usando un navegador moderno
2. Ajusta la escala al 100%
3. Desactiva los gráficos de fondo si es necesario

### No puedo imprimir:
1. Verifica que tu impresora esté conectada y encendida
2. Intenta guardar como PDF primero
3. Verifica los permisos de impresión del navegador

### Los totales no se ven correctos:
1. Verifica los datos de la venta
2. Comprueba la configuración de impuestos
3. Contacta al administrador si persiste

## Preguntas Frecuentes

### ¿Puedo reimprimir una factura?
Sí, puedes acceder al listado de facturas y seleccionar la que deseas reimprimir.

### ¿Se puede personalizar el diseño de la factura?
Actualmente el diseño es estándar, pero futuras versiones permitirán personalización.

### ¿Puedo enviar la factura por email?
Esta función estará disponible en una próxima actualización.

### ¿Las facturas se guardan automáticamente?
Sí, todas las facturas se guardan en la base de datos automáticamente al confirmar la venta.

### ¿Qué pasa si cierro la ventana sin imprimir?
No hay problema, la factura queda guardada y puedes imprimirla desde el listado de facturas.

### ¿Puedo editar una factura después de crearla?
Las facturas no se pueden editar una vez creadas para mantener integridad fiscal. Puedes cancelarla y crear una nueva si es necesario.

## Mejores Prácticas

### 📋 Antes de Imprimir:
- Revisa que todos los datos sean correctos
- Verifica el cliente y tipo de NCF
- Confirma los productos y cantidades
- Revisa los totales

### 🖨️ Al Imprimir:
- Usa papel de buena calidad para facturas oficiales
- Mantén copias de respaldo (PDF)
- Organiza las facturas por fecha
- Guarda un registro digital

### 📦 Después de Imprimir:
- Entrega la factura al cliente
- Archiva una copia para tu contabilidad
- Verifica que se registró correctamente en el sistema
- Actualiza tu inventario si es necesario

## Soporte

Si experimentas problemas con la impresión de facturas:

1. **Revisa esta guía** para soluciones comunes
2. **Verifica tu configuración** de navegador y impresora
3. **Contacta al soporte técnico** si el problema persiste

---

**¿Necesitas ayuda adicional?**  
Consulta la documentación completa en `IMPRESION_FACTURAS_IMPLEMENTADO.md`

