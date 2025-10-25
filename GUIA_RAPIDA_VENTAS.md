# 🚀 Guía Rápida - Sistema de Ventas Mejorado

**¡Todo listo para usar!** Esta guía te muestra cómo usar las nuevas funcionalidades.

---

## 📦 1. Gestionar Productos con Stock

### Crear Producto
1. Ir a **"Productos"**
2. Click **"Nuevo Producto"**
3. Llenar formulario:
   - Código: `LAP-001`
   - Nombre: `Laptop Dell XPS 15`
   - Tipo: **Producto** (importante!)
   - Precio: `75000`
   - **Stock Inicial: `50`** ← Campo nuevo
4. **Guardar**

### Indicadores de Stock
- 🟢 **Verde** (> 10): Stock normal
- 🟡 **Amarillo** (≤ 10): Stock bajo
- 🔴 **Rojo** (≤ 0): Sin stock

---

## 🛒 2. Realizar una Venta

### Paso 1: Agregar Productos
```
Panel izquierdo → Buscar productos → Click en producto
```

### Paso 2: Completar Venta
```
Panel derecho → Click "Completar Venta" (botón naranja)
```

### Paso 3: Modal de Checkout

Se abre un modal con **4 secciones**:

#### A. Seleccionar Cliente
```
[Cliente Genérico]  → Venta rápida (consumidor final)
[Buscar Cliente]    → Cliente específico (empresa/persona)
```

#### B. Tipo de Comprobante
```
[NCF 02] 🟠 → Consumidor Final (default)
[NCF 01] 🟢 → Crédito Fiscal (empresas)
```

#### C. Método de Pago
```
💵 Efectivo       💳 Tarjeta
🏦 Transferencia  💰 Mixto
```

#### D. Resumen
```
Productos: X items
Subtotal:  RD$ XX,XXX.XX
ITBIS:     RD$ X,XXX.XX
──────────────────────
Total:     RD$ XX,XXX.XX
```

### Paso 4: Confirmar
```
[Confirmar Venta] → Procesa la venta
```

**Resultado**:
- ✅ Factura creada
- ✅ Stock reducido automáticamente
- ✅ Venta registrada
- ✅ Carrito limpiado

---

## 💨 Venta Rápida (10 segundos)

Para ventas al mostrador sin identificar cliente:

```
1. Agregar productos al carrito
2. Click "Completar Venta"
3. Click "Cliente Genérico"
4. NCF 02 (ya seleccionado) ✓
5. Efectivo (ya seleccionado) ✓
6. Click "Confirmar Venta"
```

**¡Listo!** Venta procesada en segundos.

---

## 🏢 Venta a Empresa (30 segundos)

Para ventas a empresas con RNC:

```
1. Agregar productos al carrito
2. Click "Completar Venta"
3. Click "Buscar Cliente"
   → Buscar empresa
   → Seleccionar
4. Click "NCF 01" (Crédito Fiscal) 🟢
5. Seleccionar método (ej: Transferencia)
6. Revisar resumen
7. Click "Confirmar Venta"
```

**¡Listo!** Venta con crédito fiscal registrada.

---

## 🎯 Casos Comunes

### Caso 1: Tienda de Retail
```
Mayoría de ventas:
→ Cliente Genérico
→ NCF 02 (Consumidor Final)
→ Efectivo o Tarjeta
```

### Caso 2: Distribuidora
```
Todas las ventas:
→ Clientes específicos (empresas)
→ NCF 01 (Crédito Fiscal)
→ Transferencia bancaria
```

### Caso 3: Restaurante
```
Todas las ventas:
→ Cliente Genérico
→ NCF 02 (Consumidor Final)
→ Efectivo o Tarjeta
```

### Caso 4: Mixto (Retail + B2B)
```
Venta B2C:
→ Cliente Genérico + NCF 02 + Efectivo

Venta B2B:
→ Cliente Específico + NCF 01 + Transferencia
```

---

## 🔍 Detalles Importantes

### Stock Automático
- ✅ Al vender, stock se reduce solo
- ✅ Si no hay stock suficiente → venta rechazada
- ✅ Solo productos físicos (servicios NO)

### Cliente Genérico
- ✅ Siempre disponible (creado automáticamente)
- ✅ ID fijo en todos los tenants
- ✅ Para ventas sin RNC

### Tipo de NCF
- ✅ NCF 02 es el default (más común)
- ✅ NCF 01 para empresas que deducen impuestos
- ✅ Cambio con un click

### Método de Pago
- ✅ Efectivo es el default
- ✅ 4 opciones disponibles
- ✅ Fácil de cambiar

---

## 🛠️ Aplicar Migraciones

### Primera Vez (Sistema Nuevo)
```bash
docker-compose up -d
```
¡Listo! Migraciones automáticas.

### Sistema Existente
```bash
# 1. Configurar variables
export DB_PASSWORD=postgres

# 2. Ejecutar migraciones
cd /home/scorpion/desa/Projects/vendix
./scripts/migrate-tenants.sh

# 3. Reiniciar backend
docker-compose restart api

# 4. Refrescar frontend
# Ctrl+F5 en el navegador
```

**¡Listo para usar!**

---

## ❓ Preguntas Frecuentes

### ¿Puedo vender sin seleccionar cliente?
✅ Sí, el sistema usa "Cliente Genérico" automáticamente.

### ¿Qué pasa si vendo más stock del disponible?
❌ La venta se rechaza con mensaje: "stock insuficiente para [producto]"

### ¿Puedo cambiar el tipo de NCF después?
❌ No directamente. Deberías cancelar la factura y crear una nueva.

### ¿Se guarda el método de pago?
⚠️ Actualmente solo en frontend. Próxima versión guardará en BD.

### ¿Los servicios también tienen stock?
❌ No. Solo productos físicos (`product_type: "product"`).

### ¿Puedo agregar más métodos de pago?
✅ Sí, es fácil agregar botones en el código.

---

## 📊 Atajos de Teclado (Futuros)

Sugerencias para implementar:

- `F2` → Abrir modal de checkout
- `F3` → Cliente genérico
- `F4` → Buscar cliente
- `ESC` → Cerrar modal
- `Enter` → Confirmar venta

---

## 🎁 Extras Implementados

### Feedback Visual
- ✅ Loading spinners al procesar
- ✅ Alertas de éxito/error
- ✅ Estados disabled claros
- ✅ Hover effects en botones

### Validaciones
- ✅ Carrito vacío → no permite checkout
- ✅ Stock insuficiente → rechaza venta
- ✅ Tipo de NCF inválido → error 400

### Reset Automático
- ✅ Al completar venta → todo se limpia
- ✅ Al limpiar todo → vuelve a defaults
- ✅ Al cancelar modal → mantiene carrito

---

## 🎊 ¡Todo Implementado y Funcionando!

### ✅ Lo que tienes ahora:

1. **Inventario Inteligente** 📦
   - Stock en tiempo real
   - Reducción automática
   - Validaciones

2. **Ventas Flexibles** 👤
   - Cliente genérico
   - Clientes registrados
   - Automático

3. **Cumplimiento Fiscal** 📄
   - NCF 01 (Crédito Fiscal)
   - NCF 02 (Consumidor Final)
   - DGII RD

4. **Checkout Profesional** 🎯
   - Modal organizado
   - 4 métodos de pago
   - Resumen claro

### 🚀 Próximos Pasos:

1. **Aplicar migraciones** (si no lo has hecho):
   ```bash
   ./scripts/migrate-tenants.sh
   docker-compose restart api
   ```

2. **Probar el sistema**:
   - Crear productos con stock
   - Realizar ventas de prueba
   - Verificar stock se reduce

3. **Capacitar usuarios**:
   - Mostrar nuevo modal de checkout
   - Explicar tipos de NCF
   - Practicar ventas rápidas

---

**¡Disfruta tu nuevo sistema de ventas mejorado!** 🎉

Para más detalles técnicos, consulta:
- `RESUMEN_FINAL_SESION_VENTAS.md` - Resumen técnico completo
- `MODAL_COMPLETAR_VENTA_IMPLEMENTADO.md` - Detalles del modal

