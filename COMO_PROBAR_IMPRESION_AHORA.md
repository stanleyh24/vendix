# 🚀 Cómo Probar la Impresión de Facturas AHORA

## Pre-requisitos

✅ El backend debe estar corriendo (puerto 8080)  
✅ El frontend debe estar corriendo (puerto 5173)  
✅ Debe haber al menos un producto y un cliente en el sistema  

## Pasos Rápidos para Probar

### 1. Iniciar el Sistema (si no está corriendo)

```bash
# Terminal 1 - Backend
cd /home/scorpion/desa/Projects/vendix
make dev

# Terminal 2 - Frontend
cd /home/scorpion/desa/Projects/vendix/frontend
npm run dev
```

### 2. Acceder al Punto de Venta

1. Abre tu navegador en: `http://localhost:5173`
2. Inicia sesión con tus credenciales
3. Ve a la sección **"Ventas"** o **"Punto de Venta"**

### 3. Crear una Venta de Prueba

#### Paso 1: Agregar Productos
- Busca un producto en el panel izquierdo
- Haz clic sobre el producto para agregarlo al carrito
- Verás el producto aparecer en el panel derecho
- Ajusta la cantidad si es necesario (botones + / -)

#### Paso 2: Completar Venta
- Haz clic en el botón **"Completar Venta"** (botón naranja grande)
- Se abrirá un modal con 3 secciones

#### Paso 3: Seleccionar Cliente
Tienes dos opciones:
- **Opción A**: Haz clic en **"Cliente Genérico"** (más rápido para pruebas)
- **Opción B**: Haz clic en **"Buscar Cliente"** y selecciona un cliente específico

#### Paso 4: Tipo de Comprobante Fiscal
- **NCF 02 (Consumidor Final)**: Para cliente genérico o público general
- **NCF 01 (Crédito Fiscal)**: Solo disponible con cliente específico que tenga RNC

#### Paso 5: Método de Pago
Selecciona cualquiera (para pruebas no importa):
- 💵 Efectivo
- 💳 Tarjeta
- 🏦 Transferencia
- 💰 Mixto

#### Paso 6: Confirmar
- Revisa el resumen (productos, subtotal, ITBIS, total)
- Haz clic en **"Confirmar Venta"**

### 4. ⭐ Modal de Éxito (NUEVO)

Después de confirmar, verás un modal de éxito con:
- ✅ Mensaje de confirmación
- 💰 Total de la venta
- 🖨️ Botón "Imprimir Factura" (botón naranja)
- 🛒 Botón "Nueva Venta"

### 5. ⭐ Imprimir la Factura (NUEVO)

#### Opción 1: Imprimir Inmediatamente
1. Haz clic en **"Imprimir Factura"**
2. Se abrirá la vista previa de la factura
3. Revisa todos los datos:
   - ✅ Número de factura (Ej: INV-00000001)
   - ✅ Tipo de comprobante (NCF 01 o 02)
   - ✅ Cliente
   - ✅ Fechas
   - ✅ Productos con cantidades y precios
   - ✅ Subtotal, ITBIS (18%), Total
4. Haz clic en el botón **"Imprimir"** (botón naranja con icono de impresora)
5. Se abrirá el diálogo de impresión de tu navegador

#### Opción 2: Guardar como PDF
1. Desde la vista previa, haz clic en **"Imprimir"**
2. En el diálogo de impresión, selecciona **"Guardar como PDF"** como destino
3. Elige ubicación y nombre del archivo
4. Haz clic en **"Guardar"**
5. ✅ ¡Tendrás tu factura en PDF!

#### Opción 3: Continuar sin Imprimir
1. En el modal de éxito, haz clic en **"Nueva Venta"**
2. El sistema se limpia y queda listo para la siguiente venta
3. Puedes imprimir la factura más tarde desde el listado de facturas

### 6. Diálogo de Impresión del Navegador

Cuando se abra el diálogo de impresión:

**Para Imprimir en Impresora:**
- Selecciona tu impresora
- Ajusta orientación: **Vertical (Portrait)**
- Ajusta tamaño: **Carta** o **A4**
- Haz clic en **"Imprimir"**

**Para Guardar como PDF:**
- Destino: **"Guardar como PDF"** o **"Microsoft Print to PDF"**
- Haz clic en **"Guardar"** o **"Imprimir"**
- Elige ubicación y nombre
- ✅ Listo!

## Vista Previa - Elementos a Verificar

### Encabezado:
- [ ] Logo/Nombre "VENDIX"
- [ ] Número de factura (INV-00000XXX)
- [ ] NCF (si aplica)
- [ ] Tipo de comprobante

### Información del Cliente:
- [ ] Nombre del cliente
- [ ] RNC (si aplica)
- [ ] Tipo de cliente

### Detalles de la Factura:
- [ ] Fecha de emisión (hoy)
- [ ] Fecha de vencimiento (30 días después)
- [ ] Estado (normalmente "draft")

### Tabla de Productos:
- [ ] Número de línea
- [ ] Descripción del producto
- [ ] Cantidad
- [ ] Precio unitario
- [ ] ITBIS calculado
- [ ] Total de la línea

### Totales:
- [ ] Subtotal (suma sin impuestos)
- [ ] ITBIS 18% (calculado correctamente)
- [ ] Total (subtotal + ITBIS)

### Pie de Página:
- [ ] "Este documento fue generado electrónicamente por VENDIX"
- [ ] Fecha y hora de generación

## Pruebas Recomendadas

### Prueba 1: Cliente Genérico + NCF 02
1. Agrega 1 producto al carrito
2. Selecciona "Cliente Genérico"
3. NCF 02 (Consumidor Final) - debe estar preseleccionado
4. Confirma y verifica la factura

### Prueba 2: Cliente Específico + NCF 01
1. Agrega varios productos al carrito
2. Busca y selecciona un cliente con RNC
3. Selecciona NCF 01 (Crédito Fiscal)
4. Confirma y verifica la factura

### Prueba 3: Múltiples Productos
1. Agrega 5+ productos diferentes
2. Ajusta cantidades variadas
3. Verifica que la tabla se vea bien en la impresión

### Prueba 4: Guardar como PDF
1. Completa una venta
2. Imprime como PDF
3. Abre el PDF y verifica que todo se vea correcto

### Prueba 5: Nueva Venta sin Imprimir
1. Completa una venta
2. Haz clic en "Nueva Venta" sin imprimir
3. Verifica que el carrito se limpió
4. Completa otra venta

## Atajos de Teclado

| Atajo | Acción |
|-------|--------|
| `Ctrl + P` (Windows/Linux) | Imprimir (desde vista previa) |
| `Cmd + P` (Mac) | Imprimir (desde vista previa) |
| `Esc` | Cerrar modal |

## Solución Rápida de Problemas

### "No puedo ver el botón de imprimir"
- ✅ Verifica que completaste la venta correctamente
- ✅ Debe aparecer un modal de éxito primero
- ✅ Haz clic en "Imprimir Factura" en ese modal

### "La factura no carga"
- ✅ Verifica que el backend esté corriendo
- ✅ Abre la consola del navegador (F12) y busca errores
- ✅ Verifica tu conexión a internet

### "Los totales están mal"
- ✅ Verifica la configuración de ITBIS de los productos (debe ser 0.18 = 18%)
- ✅ Revisa los precios de los productos
- ✅ El cálculo es: Subtotal + (Subtotal × 0.18) = Total

### "El diseño se ve mal al imprimir"
- ✅ Usa un navegador moderno (Chrome, Firefox, Edge)
- ✅ Ajusta la escala al 100% en el diálogo de impresión
- ✅ Selecciona orientación vertical
- ✅ Desactiva "Gráficos de fondo" si está activado

## Verificación Final

Después de probar, verifica:
- [ ] Puedes completar una venta
- [ ] El modal de éxito aparece correctamente
- [ ] La vista previa muestra todos los datos
- [ ] Los cálculos son correctos
- [ ] Puedes imprimir o guardar como PDF
- [ ] Puedes iniciar una nueva venta
- [ ] El diseño se ve profesional

## Próximos Pasos

Una vez que hayas probado exitosamente:

1. **Prueba con datos reales** de tu negocio
2. **Imprime en papel** para ver el resultado físico
3. **Comparte** la factura PDF con un cliente de prueba
4. **Revisa** el listado de facturas (próxima funcionalidad)
5. **Personaliza** si es necesario (logo, colores, etc.)

## Soporte

Si encuentras algún problema:
1. Revisa la consola del navegador (F12)
2. Verifica los logs del backend
3. Consulta la documentación completa en:
   - `IMPRESION_FACTURAS_IMPLEMENTADO.md` (técnico)
   - `GUIA_IMPRESION_FACTURAS.md` (usuario)

---

## Resumen del Flujo Completo

```
📦 Agregar productos → 🛒 Completar venta → 👤 Seleccionar cliente 
  ↓
📄 Elegir NCF → 💰 Método de pago → ✅ Confirmar
  ↓
🎉 Modal de éxito → 🖨️ Imprimir Factura (NUEVO)
  ↓
👁️ Vista previa → 🖨️ Imprimir/PDF → ✅ ¡Listo!
```

**¡Todo listo para usar! 🎉**

---

**Fecha**: 19 de octubre de 2025  
**Versión**: 1.0  
**Estado**: ✅ FUNCIONANDO

