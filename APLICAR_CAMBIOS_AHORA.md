# 🚀 Aplicar Cambios AHORA - Guía Paso a Paso

**¡Todo está listo!** Sigue estos pasos para aplicar las mejoras.

---

## ⚡ Opción Rápida (Recomendada)

```bash
cd /home/scorpion/desa/Projects/vendix

# Configurar credenciales (ajusta si es necesario)
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=vendix
export DB_USER=postgres
export DB_PASSWORD=postgres

# Aplicar TODAS las migraciones (v17, v18, v19)
./scripts/migrate-tenants.sh

# Reiniciar backend
docker-compose restart api
```

**Refrescar navegador**: `Ctrl+F5`

**¡Listo!** ✅

---

## 📋 Opción Detallada (Paso a Paso)

### Paso 1: Ir al Directorio
```bash
cd /home/scorpion/desa/Projects/vendix
```

### Paso 2: Configurar Variables de Entorno
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=vendix
export DB_USER=postgres
export DB_PASSWORD=postgres
```

### Paso 3: Aplicar Migración v17 (Inventario)
```bash
./scripts/add-stock-to-products.sh
```
**Resultado esperado**:
```
========================================
  Migración: Sistema de Inventario
  Versión: 17
========================================
...
✓ Migración completada exitosamente en todos los tenants
```

### Paso 4: Aplicar Migración v18 (Cliente Genérico)
```bash
./scripts/create-generic-customer.sh
```
**Resultado esperado**:
```
========================================
  Migración: Cliente Genérico
  Versión: 18
========================================
...
✓ Migración completada exitosamente
Cliente Genérico:
  ID: 00000000-0000-0000-0000-000000000001
  Nombre: Cliente Genérico
```

### Paso 5: Aplicar Migración v19 (Tipo NCF)
```bash
./scripts/add-ncf-type.sh
```
**Resultado esperado**:
```
========================================
  Migración: Tipo de Comprobante NCF
  Versión: 19
========================================
...
✓ Migración completada exitosamente
Tipo de Comprobante Fiscal (NCF):
  • NCF 01 = Crédito Fiscal (empresas con RNC)
  • NCF 02 = Consumidor Final (default)
```

### Paso 6: Reiniciar Backend
```bash
docker-compose restart api
```

**Esperar** ~5-10 segundos a que el servicio reinicie.

### Paso 7: Refrescar Frontend
```
1. Abrir navegador
2. Presionar Ctrl+F5 (Windows/Linux) o Cmd+Shift+R (Mac)
3. Limpiar caché si es necesario
```

---

## ✅ Verificar que Todo Funciona

### Test 1: Productos con Stock
```
1. Ir a "Productos"
2. Click "Nuevo Producto"
3. Verificar que hay campo "Cantidad en Stock"
4. Crear producto con stock 100
5. Verificar indicador verde 🟢 aparece
✅ Si ves el campo de stock → Migración v17 OK
```

### Test 2: Cliente Genérico
```
1. Ir a "Ventas"
2. Agregar productos al carrito
3. Click "Completar Venta"
4. Verificar botón "Cliente Genérico" aparece
✅ Si ves el botón → Migración v18 OK
```

### Test 3: Tipo de NCF
```
1. En modal de checkout
2. Verificar sección "Tipo de Comprobante Fiscal"
3. Ver botones NCF 02 y NCF 01
✅ Si ves los botones → Migración v19 OK
```

### Test 4: Validación NCF + Cliente
```
1. Click "Cliente Genérico"
2. Verificar que NCF 01 se pone gris (disabled)
3. Intentar click en NCF 01
4. Verificar que no hace nada
5. Ver mensaje: "⚠️ Crédito Fiscal requiere cliente con RNC"
✅ Si funciona → Validación OK
```

### Test 5: Venta Completa
```
1. Agregar productos
2. Click "Completar Venta"
3. Seleccionar "Cliente Genérico"
4. Dejar NCF 02 (default)
5. Dejar Efectivo (default)
6. Click "Confirmar Venta"
7. Ver mensaje: "Venta procesada"
8. Verificar stock del producto se redujo
✅ Si todo funciona → Sistema 100% OK
```

---

## 🎯 Lo que Deberías Ver

### En Productos
```
┌─────────────────────────────┐
│ 📦 Laptop Dell XPS 15       │
│ LAP-001                     │
│ Stock: 48.00 unidad 🟢      │ ← NUEVO
│ Precio: RD$ 75,000.00       │
│ [Activo]                    │
└─────────────────────────────┘
```

### En Ventas (Modal Checkout)
```
┌─────────────────────────────────────┐
│ 📋 Completar Venta           [X]   │
├─────────────────────────────────────┤
│ 👤 Cliente                          │
│ [Cliente Genérico] [Buscar]         │ ← NUEVO
│                                     │
│ 📄 Tipo de Comprobante              │
│ [NCF 02 🟠] [NCF 01 (bloq.)]       │ ← NUEVO
│                                     │
│ 💵 Método de Pago                   │
│ [💵 Efectivo] [💳 Tarjeta]          │ ← NUEVO
│ [🏦 Transfer.] [💰 Mixto]           │
│                                     │
│ 📊 Resumen                          │
│ Total: RD$ XXX,XXX.XX               │
│                                     │
│ [Cancelar] [✓ Confirmar Venta]     │
└─────────────────────────────────────┘
```

---

## ❌ Errores Comunes y Soluciones

### Error 1: Script no ejecuta
```
bash: ./scripts/migrate-tenants.sh: Permission denied
```
**Solución**:
```bash
chmod +x ./scripts/*.sh
./scripts/migrate-tenants.sh
```

### Error 2: No conecta a PostgreSQL
```
Error: No se puede conectar a PostgreSQL
```
**Solución**:
```bash
# Verificar que PostgreSQL está corriendo
docker-compose ps

# Verificar credenciales
echo $DB_PASSWORD

# Probar conexión manual
psql -h localhost -U postgres -d vendix -c '\dt'
```

### Error 3: Tenant no encontrado
```
No se encontraron tenants activos
```
**Solución**:
```bash
# Crear tenant demo si no existe
./scripts/create-test-tenant.sh

# O verificar tenants existentes
psql -h localhost -U postgres -d vendix -c "SELECT schema_name FROM tenants;"
```

### Error 4: Migración ya aplicada
```
⊙ Ya tiene la columna stock_quantity - Saltando
```
**Esto es normal**: La migración detecta que ya está aplicada y la salta.

---

## 📊 Checklist Final

Marca cada item cuando lo completes:

- [ ] Variables de entorno configuradas
- [ ] Migración v17 aplicada (stock)
- [ ] Migración v18 aplicada (cliente genérico)
- [ ] Migración v19 aplicada (NCF type)
- [ ] Backend reiniciado
- [ ] Frontend refrescado (Ctrl+F5)
- [ ] Test: Campo stock en productos
- [ ] Test: Botón cliente genérico
- [ ] Test: Selector NCF en modal
- [ ] Test: Métodos de pago visibles
- [ ] Test: Venta completa exitosa
- [ ] Test: Stock se reduce automáticamente
- [ ] Test: NCF 01 bloqueado con genérico

---

## 🎯 ¿Qué Debes Hacer Ahora?

### 1. Aplicar Migraciones (5 minutos)
```bash
export DB_PASSWORD=postgres
./scripts/migrate-tenants.sh
docker-compose restart api
```

### 2. Probar el Sistema (10 minutos)
```
- Crear producto con stock
- Realizar venta de prueba
- Verificar modal de checkout
- Probar todas las opciones
```

### 3. Capacitar Usuarios (30 minutos)
```
- Mostrar nueva interfaz
- Explicar cliente genérico
- Explicar tipos de NCF
- Practicar ventas rápidas
```

### 4. ¡Comenzar a Vender! 🚀
```
- Sistema listo para producción
- Todas las validaciones activas
- Documentación disponible
```

---

## 📞 Soporte

Si algo no funciona:

1. **Revisar logs del backend**:
   ```bash
   docker-compose logs -f api
   ```

2. **Revisar consola del navegador**:
   - F12 → Console
   - Ver errores JavaScript

3. **Consultar documentación**:
   - `GUIA_RAPIDA_VENTAS.md` - Guía de uso
   - `SESION_COMPLETA_MEJORAS_VENTAS.md` - Resumen técnico

4. **Re-aplicar migraciones**:
   ```bash
   ./scripts/migrate-tenants.sh
   ```

---

## 🎉 ¡Estás a UN COMANDO de Distancia!

```bash
export DB_PASSWORD=postgres && ./scripts/migrate-tenants.sh && docker-compose restart api
```

**¡Ejecuta este comando y todo estará listo!** 🚀

---

**Fecha**: 19 de Octubre, 2025  
**Sistema**: Vendix - Facturación Electrónica RD  
**Estado**: ✅ COMPLETADO  
**Archivos Modificados**: 23  
**Migraciones**: 3 (v17, v18, v19)  
**Tiempo de Aplicación**: ~5 minutos

---

## 🎊 ¡A Vender!

Tu sistema ahora está al nivel de los mejores POS del mercado.

**¡Disfrútalo!** 🎉

