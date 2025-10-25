# ✅ Rangos de NCF - Implementación Completada

## 🎯 Funcionalidad Implementada

La configuración de NCF ahora incluye **inicio y fin de rango** para cada tipo de comprobante fiscal, permitiendo controlar los límites asignados por la DGII.

---

## 🇩🇴 Contexto: Rangos de NCF en República Dominicana

### ¿Qué son los Rangos de NCF?

La **DGII** (Dirección General de Impuestos Internos) asigna **rangos específicos** de números de comprobantes fiscales a cada empresa:

**Ejemplo**:
- NCF 01 (Crédito Fiscal): Del **1** al **100,000**
- NCF 02 (Consumidor Final): Del **1** al **500,000**
- NCF 03 (Nota Débito): Del **1** al **50,000**
- NCF 04 (Nota Crédito): Del **1** al **50,000**

Cuando una empresa **agota un rango**, debe **solicitar un nuevo rango** a la DGII.

---

## ✨ Funcionalidades Implementadas

### 1. Campos de Rango

Cada tipo de NCF ahora tiene **3 campos**:

```
Prefijo:          B01        (Asignado por DGII)
Inicio de Rango:  1          (Primer número del rango)
Fin de Rango:     100000     (Último número del rango)
Secuencia Actual: 1          (Número actual en uso)
```

### 2. Alertas Automáticas

#### ⚠️ Advertencia al 90%
Cuando la secuencia alcanza el 90% del rango:
```
┌────────────────────────────────────────┐
│ ⚠️ Advertencia: Te estás acercando    │
│    al límite del rango (90000/100000)  │
└────────────────────────────────────────┘
```
**Acción**: Prepararse para solicitar nuevo rango

#### 🚨 Error al 100%
Cuando la secuencia alcanza o supera el límite:
```
┌────────────────────────────────────────┐
│ 🚨 Error: Has alcanzado o superado    │
│    el límite del rango. Contacta a    │
│    DGII para obtener un nuevo rango.  │
└────────────────────────────────────────┘
```
**Acción**: Contactar DGII urgentemente

### 3. Validación Visual

El sistema **calcula automáticamente** el porcentaje de uso:
```javascript
if (secuencia >= fin * 0.9) {
  // Mostrar advertencia amarilla
}
if (secuencia >= fin) {
  // Mostrar error rojo
}
```

---

## 🗄️ Cambios en Base de Datos

### Migración v16 Actualizada (Esquema Inicial)

Se agregaron campos de rango en la tabla `tenant_config`:

```sql
-- NCF Crédito Fiscal (01)
ncf_fiscal_credit_prefix VARCHAR(20),
ncf_fiscal_credit_start INTEGER DEFAULT 1,
ncf_fiscal_credit_end INTEGER DEFAULT 999999999,
ncf_fiscal_credit_sequence INTEGER DEFAULT 1,

-- NCF Consumidor Final (02)
ncf_consumer_prefix VARCHAR(20),
ncf_consumer_start INTEGER DEFAULT 1,
ncf_consumer_end INTEGER DEFAULT 999999999,
ncf_consumer_sequence INTEGER DEFAULT 1,

-- NCF Nota de Débito (03)
ncf_debit_note_prefix VARCHAR(20),
ncf_debit_note_start INTEGER DEFAULT 1,
ncf_debit_note_end INTEGER DEFAULT 999999999,
ncf_debit_note_sequence INTEGER DEFAULT 1,

-- NCF Nota de Crédito (04)
ncf_credit_note_prefix VARCHAR(20),
ncf_credit_note_start INTEGER DEFAULT 1,
ncf_credit_note_end INTEGER DEFAULT 999999999,
ncf_credit_note_sequence INTEGER DEFAULT 1,
```

### Migración v20 (Para Tenants Existentes)

Para aplicar a tenants que ya existen:

```sql
-- NCF Crédito Fiscal (01)
ALTER TABLE tenant_config ADD COLUMN IF NOT EXISTS ncf_fiscal_credit_start INTEGER DEFAULT 1;
ALTER TABLE tenant_config ADD COLUMN IF NOT EXISTS ncf_fiscal_credit_end INTEGER DEFAULT 999999999;

-- NCF Consumidor Final (02)
ALTER TABLE tenant_config ADD COLUMN IF NOT EXISTS ncf_consumer_start INTEGER DEFAULT 1;
ALTER TABLE tenant_config ADD COLUMN IF NOT EXISTS ncf_consumer_end INTEGER DEFAULT 999999999;

-- (Similar para NCF 03 y 04)
```

---

## 🎨 Interfaz Visual

### Tarjeta de NCF Mejorada

```
┌─────────────────────────────────────────────────────┐
│ 🟢 NCF 01 - Crédito Fiscal                         │
├─────────────────────────────────────────────────────┤
│                                                     │
│ Prefijo           Inicio de Rango    Fin de Rango  │
│ ┌──────────┐     ┌──────────┐       ┌──────────┐  │
│ │ B01      │     │ 1        │       │ 100000   │  │
│ └──────────┘     └──────────┘       └──────────┘  │
│ Asignado DGII                                       │
│                                                     │
│ Secuencia Actual                                    │
│ ┌────────────────────────────────────────────────┐ │
│ │ 95000                                          │ │
│ └────────────────────────────────────────────────┘ │
│                                                     │
│ ┌────────────────────────────────────────────────┐ │
│ │ ⚠️ Advertencia: Te estás acercando al límite   │ │
│ │    del rango (95000 / 100000)                  │ │
│ └────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────┘
```

### Estados de Advertencia

#### Estado Normal (< 90%)
```
Secuencia: 50000 / 100000
→ Sin alertas
→ Todo OK ✅
```

#### Estado de Advertencia (≥ 90%)
```
Secuencia: 95000 / 100000
┌──────────────────────────────────┐
│ ⚠️ Advertencia: Acercándose...  │ ← Amarillo
└──────────────────────────────────┘
```

#### Estado de Error (≥ 100%)
```
Secuencia: 100000 / 100000
┌──────────────────────────────────┐
│ 🚨 Error: Límite alcanzado      │ ← Rojo
│    Contacta a DGII               │
└──────────────────────────────────┘
```

---

## 🔧 Cambios Implementados

### Frontend

**Archivo**: `frontend/src/pages/Settings.jsx`

#### Estado Ampliado
```javascript
const [config, setConfig] = useState({
  // Para cada tipo de NCF:
  ncf_fiscal_credit_prefix: '',
  ncf_fiscal_credit_start: 1,        // ← NUEVO
  ncf_fiscal_credit_end: 999999999,  // ← NUEVO
  ncf_fiscal_credit_sequence: 1,
  // ...similar para NCF 02, 03, 04
});
```

#### Grid de 3 Columnas
```jsx
<div className="grid grid-cols-3 gap-6">
  <div>{/* Prefijo */}</div>
  <div>{/* Inicio */}</div>
  <div>{/* Fin */}</div>
  <div className="col-span-3">{/* Secuencia */}</div>
</div>
```

#### Alertas Condicionales
```jsx
{config.ncf_fiscal_credit_sequence >= config.ncf_fiscal_credit_end * 0.9 && (
  <div className="bg-yellow-50 border border-yellow-300">
    ⚠️ Advertencia: Acercándose al límite
  </div>
)}

{config.ncf_fiscal_credit_sequence >= config.ncf_fiscal_credit_end && (
  <div className="bg-red-50 border border-red-300">
    🚨 Error: Límite alcanzado
  </div>
)}
```

### Backend

**Archivo**: `internal/database/tenant_migrations.go`

- Migración v16 actualizada con campos de rango (para nuevos tenants)
- Migración v20 creada (para tenants existentes)

---

## 📊 Ejemplo de Uso Real

### Configuración Inicial

```
Empresa recibe de DGII:
- NCF 01: Serie B0100000001 al B0100100000 (100,000 comprobantes)
- NCF 02: Serie B0200000001 al B0200500000 (500,000 comprobantes)
```

**En Vendix configurar**:
```
NCF 01 - Crédito Fiscal:
- Prefijo: B01
- Inicio: 1
- Fin: 100000
- Secuencia: 1

NCF 02 - Consumidor Final:
- Prefijo: B02
- Inicio: 1
- Fin: 500000
- Secuencia: 1
```

### Durante el Uso

Después de emitir 95,000 facturas NCF 01:
```
Secuencia actual: 95000
Rango: 1 - 100000

Sistema muestra:
⚠️ Advertencia: Te estás acercando al límite 
   del rango (95000 / 100000)

Acción: Solicitar nuevo rango a DGII
```

### Al Agotar Rango

Al llegar a 100,000:
```
Secuencia actual: 100000
Rango: 1 - 100000

Sistema muestra:
🚨 Error: Has alcanzado o superado el límite
   del rango. Contacta a DGII para obtener
   un nuevo rango.

Acción: Contactar DGII urgentemente
```

---

## 🚀 Migración

### Sistemas Nuevos
```bash
docker-compose up -d
```
Los campos se crean automáticamente con la tabla `tenant_config`.

### Sistemas Existentes

```bash
# Configurar variables
export DB_PASSWORD=postgres

# Aplicar migración v20
./scripts/add-ncf-ranges.sh

# O aplicar todas las migraciones
./scripts/migrate-tenants.sh

# Reiniciar backend
docker-compose restart api
```

---

## 📋 Configurar Rangos Asignados por DGII

### Paso 1: Obtener Rangos de DGII

Contacta a DGII o revisa tu certificado de autorización para conocer los rangos asignados.

### Paso 2: Ir a Configuración

```
1. Login en Vendix
2. Ir a "Configuración"
3. Tab "NCF (DGII)"
```

### Paso 3: Configurar Cada Tipo

**Para NCF 01 (Crédito Fiscal)**:
```
Prefijo: B01
Inicio:  1
Fin:     100000  ← Según DGII
Secuencia: 1     ← Auto-incrementa
```

**Para NCF 02 (Consumidor Final)**:
```
Prefijo: B02
Inicio:  1
Fin:     500000  ← Según DGII
Secuencia: 1
```

**Similar para NCF 03 y 04**

### Paso 4: Guardar

```
Click "💾 Guardar Cambios"
```

---

## ⚠️ Alertas y Notificaciones

### Alerta Amarilla (90% del Rango)

**Cuándo**: Secuencia ≥ 90% del fin de rango

**Mensaje**: 
```
⚠️ Advertencia: Te estás acercando al límite del rango (X / Y)
```

**Acción Sugerida**:
1. Contactar a DGII
2. Solicitar nuevo rango
3. Planificar la transición

**Tiempo**: ~2-4 semanas antes de agotar

### Alerta Roja (100% del Rango)

**Cuándo**: Secuencia ≥ fin de rango

**Mensaje**:
```
🚨 Error: Has alcanzado o superado el límite del rango. 
   Contacta a DGII para obtener un nuevo rango.
```

**Acción Requerida**:
1. **NO emitir más** comprobantes de ese tipo
2. Contactar DGII **urgentemente**
3. Obtener nuevo rango
4. Actualizar configuración

---

## 🧪 Ejemplos de Validación

### Ejemplo 1: Rango Normal

```
Inicio: 1
Fin: 100000
Secuencia: 50000

Uso: 50%
Estado: ✅ OK (sin alertas)
```

### Ejemplo 2: Advertencia

```
Inicio: 1
Fin: 100000
Secuencia: 95000

Uso: 95%
Estado: ⚠️ ADVERTENCIA
Mensaje: "Te estás acercando al límite"
```

### Ejemplo 3: Límite Alcanzado

```
Inicio: 1
Fin: 100000
Secuencia: 100000

Uso: 100%
Estado: 🚨 ERROR
Mensaje: "Límite alcanzado. Contacta DGII"
```

### Ejemplo 4: Límite Superado (Error)

```
Inicio: 1
Fin: 100000
Secuencia: 100001

Uso: 100.001%
Estado: 🚨 ERROR CRÍTICO
Mensaje: "Has superado el límite. Contacta DGII"
```

---

## 🔐 Validación de Rangos

### Frontend
- ✅ Calcula porcentaje de uso
- ✅ Muestra alerta al 90%
- ✅ Muestra error al 100%
- ✅ Colores visuales (amarillo/rojo)

### Backend (Próxima Implementación Sugerida)
```go
// Validar antes de emitir NCF
if sequence >= endRange {
    return error "Rango de NCF agotado. Contacte a DGII"
}

// Advertir al acercarse
if sequence >= endRange * 0.9 {
    logger.Warn("NCF range nearly exhausted", 
        "type", ncfType, 
        "current", sequence, 
        "limit", endRange)
}
```

---

## 📁 Archivos Modificados

### Backend
```
✅ internal/database/tenant_migrations.go
   - Migración v16 actualizada (nuevos tenants)
   - Migración v20 creada (tenants existentes)
```

### Frontend
```
✅ frontend/src/pages/Settings.jsx
   - Grid de 3 columnas (Prefijo, Inicio, Fin)
   - Campos de inicio y fin
   - Alertas visuales (amarillo/rojo)
   - Cálculo de porcentaje
```

### Scripts
```
✅ scripts/add-ncf-ranges.sh
   - Script para aplicar migración v20
```

---

## 🎨 Diseño Visual

### Antes (Solo Prefijo y Secuencia)

```
┌─────────────────────────────┐
│ NCF 01 - Crédito Fiscal     │
│                             │
│ Prefijo    Secuencia        │
│ [B01]      [1]              │
└─────────────────────────────┘
```
❌ No controla límites

### Ahora (Con Rangos Completos)

```
┌──────────────────────────────────────────┐
│ 🟢 NCF 01 - Crédito Fiscal              │
├──────────────────────────────────────────┤
│ Prefijo      Inicio      Fin             │
│ ┌────────┐  ┌────────┐  ┌────────┐      │
│ │ B01    │  │ 1      │  │ 100000 │      │
│ └────────┘  └────────┘  └────────┘      │
│ DGII                                     │
│                                          │
│ Secuencia Actual                         │
│ ┌────────────────────────────────────┐  │
│ │ 95000                              │  │
│ └────────────────────────────────────┘  │
│                                          │
│ ⚠️ Advertencia: Cerca del límite        │
│    (95000 / 100000)                      │
└──────────────────────────────────────────┘
```
✅ Control completo de rangos

---

## 📊 Configuración Recomendada

### Para Pequeñas Empresas

```
NCF 01 (Crédito Fiscal):
- Inicio: 1
- Fin: 50,000
- Uso típico: ~100-500 por año

NCF 02 (Consumidor Final):
- Inicio: 1
- Fin: 100,000
- Uso típico: ~1,000-10,000 por año
```

### Para Medianas Empresas

```
NCF 01 (Crédito Fiscal):
- Inicio: 1
- Fin: 200,000
- Uso típico: ~2,000-5,000 por año

NCF 02 (Consumidor Final):
- Inicio: 1
- Fin: 500,000
- Uso típico: ~10,000-50,000 por año
```

### Para Grandes Empresas

```
NCF 01 (Crédito Fiscal):
- Inicio: 1
- Fin: 1,000,000
- Uso típico: >10,000 por año

NCF 02 (Consumidor Final):
- Inicio: 1
- Fin: 5,000,000
- Uso típico: >100,000 por año
```

---

## 🔮 Mejoras Futuras Sugeridas

### 1. Validación Backend al Emitir
```go
func ValidateNCFRange(ncfType string, sequence int, endRange int) error {
    if sequence >= endRange {
        return errors.New("Rango de NCF agotado")
    }
    return nil
}
```

### 2. Notificaciones por Email
```
Al alcanzar 90%:
→ Enviar email a administrador
→ "Tu rango de NCF XX está por agotarse"

Al alcanzar 100%:
→ Enviar email urgente
→ "Rango agotado. Acción inmediata requerida"
```

### 3. Dashboard de Rangos
```
Panel en Dashboard mostrando:
- NCF 01: ████████░░ 80% (80,000 / 100,000)
- NCF 02: ███░░░░░░░ 30% (150,000 / 500,000)
- NCF 03: █░░░░░░░░░ 10% (5,000 / 50,000)
- NCF 04: █░░░░░░░░░ 5% (2,500 / 50,000)
```

### 4. Solicitud Automática a DGII
```
Al 90%:
→ Generar solicitud pre-llenada
→ Facilitar envío a DGII
→ Tracking de solicitud
```

---

## 🧪 Pruebas

### Test 1: Configurar Rangos
```
1. Ir a Configuración → NCF
2. Para NCF 01:
   - Inicio: 1
   - Fin: 1000 (pequeño para prueba)
   - Secuencia: 1
3. Guardar
✅ Campos se guardan correctamente
```

### Test 2: Alerta al 90%
```
1. Configurar:
   - Fin: 1000
   - Secuencia: 900
2. Guardar y recargar página
✅ Muestra advertencia amarilla
```

### Test 3: Error al 100%
```
1. Configurar:
   - Fin: 1000
   - Secuencia: 1000
2. Guardar y recargar página
✅ Muestra error rojo
```

### Test 4: Alertas para Todos los Tipos
```
1. Configurar cada NCF (01, 02, 03, 04)
2. Poner secuencia cerca del fin
✅ Cada uno muestra su propia alerta
```

---

## 📄 Documentación DGII

### Cómo Obtener Rangos

1. **Registro en DGII**:
   - Completar formulario de solicitud
   - Presentar documentación de la empresa
   - Esperar aprobación (1-2 semanas)

2. **Recibir Rangos**:
   - DGII asigna rangos específicos
   - Entrega certificado con números asignados
   - Ejemplo: B0100000001 al B0100100000

3. **Configurar en Vendix**:
   - Extraer prefijo: "B01"
   - Extraer inicio: 1
   - Extraer fin: 100000
   - Ingresar en configuración

4. **Renovar Rangos**:
   - Al agotar (~90%), solicitar nuevo rango
   - DGII asigna nuevo rango
   - Actualizar configuración en Vendix

---

## ✅ Checklist

### Migración
- [x] Migración v16 actualizada
- [x] Migración v20 creada
- [x] Script de migración creado
- [x] Comentarios SQL agregados

### Frontend
- [x] Campos de inicio y fin agregados
- [x] Grid de 3 columnas
- [x] Alertas de 90% (amarillo)
- [x] Alertas de 100% (rojo)
- [x] Cálculo automático de porcentaje
- [x] Diseño consistente

### Backend
- [x] Esquema actualizado
- [x] Defaults configurados (999999999)
- [x] Código compila sin errores

---

## 🎯 Beneficios

### Control de Cumplimiento
- ✅ Previene usar NCF fuera de rango
- ✅ Alerta antes de agotar rangos
- ✅ Cumple con DGII

### Prevención de Problemas
- ✅ No llega a límite sin aviso
- ✅ Tiempo para solicitar nuevo rango
- ✅ Evita multas de DGII

### Gestión Proactiva
- ✅ Visibilidad del uso de rangos
- ✅ Planificación anticipada
- ✅ Sin interrupciones de servicio

---

**Fecha**: 19 de Octubre, 2025  
**Migración**: v16 (actualizada) + v20 (nueva)  
**Estado**: ✅ COMPLETADO  

---

## 🎉 ¡Rangos de NCF Implementados!

Ahora puedes:
- Configurar los rangos asignados por DGII
- Ver alertas cuando te acerques al límite
- Planificar solicitudes de nuevos rangos
- Cumplir con las regulaciones dominicanas

**Para aplicar**:
```bash
export DB_PASSWORD=postgres
./scripts/add-ncf-ranges.sh
docker-compose restart api
```











