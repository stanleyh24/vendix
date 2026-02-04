# ✅ Configuración de Tenant Implementada

## 📋 Resumen

Se ha implementado un módulo completo de configuración para cada tenant, permitiendo que cada empresa configure sus datos, preferencias fiscales, facturación y personalización.

---

## 🗄️ Base de Datos

### Nueva Tabla: `tenant_config`

Se agregó la migración v15 que crea la tabla `tenant_config` con los siguientes campos:

#### **Información Básica de la Empresa**
- `company_legal_name` - Nombre legal de la empresa (requerido)
- `company_trade_name` - Nombre comercial
- `company_tax_id` - RNC/Cédula
- `company_email` - Email de contacto
- `company_phone` - Teléfono de contacto
- `company_website` - Sitio web

#### **Dirección**
- `address_line1` - Dirección línea 1
- `address_line2` - Dirección línea 2
- `city` - Ciudad
- `state_province` - Provincia/Estado
- `postal_code` - Código postal
- `country` - País (default: 'DO')

#### **Configuración Fiscal**
- `tax_regime` - Régimen fiscal
- `dgii_rnc` - RNC para DGII
- `dgii_user` - Usuario DGII
- `dgii_api_key` - API Key DGII
- `dgii_environment` - Ambiente (sandbox/production)
- `ncf_enabled` - NCF habilitado (boolean)

#### **Configuración de Facturación**
- `invoice_prefix` - Prefijo de facturas (default: 'INV')
- `invoice_next_number` - Siguiente número de factura
- `invoice_terms` - Términos de factura
- `invoice_footer` - Pie de página de factura
- `payment_terms_days` - Días de pago (default: 30)

#### **Numeración NCF**
- `ncf_fiscal_credit_prefix` - Prefijo crédito fiscal
- `ncf_fiscal_credit_sequence` - Secuencia crédito fiscal
- `ncf_consumer_prefix` - Prefijo consumidor final
- `ncf_consumer_sequence` - Secuencia consumidor final
- `ncf_debit_note_prefix` - Prefijo nota de débito
- `ncf_debit_note_sequence` - Secuencia nota de débito
- `ncf_credit_note_prefix` - Prefijo nota de crédito
- `ncf_credit_note_sequence` - Secuencia nota de crédito

#### **Configuración Monetaria**
- `default_currency` - Moneda por defecto (default: 'DOP')
- `default_tax_rate` - Tasa de impuesto (default: 18.00)

#### **Notificaciones**
- `notification_email` - Email para notificaciones
- `send_invoice_emails` - Enviar facturas por email (boolean)
- `send_payment_reminders` - Enviar recordatorios de pago (boolean)

#### **Branding**
- `logo_url` - URL del logo
- `primary_color` - Color primario (default: '#3B82F6')
- `secondary_color` - Color secundario (default: '#10B981')

#### **Localización**
- `timezone` - Zona horaria (default: 'America/Santo_Domingo')
- `date_format` - Formato de fecha (default: 'DD/MM/YYYY')
- `time_format` - Formato de hora (default: 'HH:mm')
- `locale` - Idioma/localización (default: 'es_DO')

---

## 🔧 Backend (Go)

### Módulo: `internal/modules/tenantconfig`

#### Archivos Creados

1. **models.go**
   - `TenantConfig` - Modelo principal
   - `UpdateTenantConfigRequest` - DTO para actualizaciones

2. **repository.go**
   - `Get(ctx)` - Obtiene la configuración del tenant
   - `Update(ctx, config)` - Actualiza la configuración

3. **service.go**
   - `Get(ctx)` - Lógica de negocio para obtener config
   - `Update(ctx, request)` - Lógica de negocio para actualizar

4. **handler.go**
   - `GET /api/v1/tenant/config` - Obtiene configuración
   - `PUT /api/v1/tenant/config` - Actualiza configuración

### Integración

Las rutas fueron registradas en `internal/server/server.go`:

```go
import "github.com/stanleyh24/vendix/internal/modules/tenantconfig"

// En setupRoutes(), después de protected := tenant.Group("")
tenantconfig.RegisterRoutes(protected, s.DB)
```

---

## 🎨 Frontend (React)

### Página: `frontend/src/pages/Settings.jsx`

#### Características

✅ **Interfaz con pestañas** para organizar configuración:
- **Empresa** - Datos básicos y dirección
- **Configuración Fiscal** - DGII, NCF, régimen fiscal
- **Facturación** - Prefijos, numeración, términos
- **NCF (DGII)** - Configuración de secuencias NCF
- **Notificaciones** - Preferencias de email
- **Identidad Visual** - Logo y colores

✅ **Formulario completo** con todos los campos
✅ **Validación** de campos requeridos
✅ **Feedback visual** con alertas de éxito/error
✅ **Estado de carga** mientras guarda

#### Integración con API

```javascript
// Obtener configuración
const response = await api.get('/config')

// Actualizar configuración
await api.put('/config', config)
```

### Navegación

Se agregó el enlace de **Configuración** al menú lateral en `Layout.jsx`:

```javascript
{ name: 'Configuración', href: '/settings', icon: Settings }
```

---

## 🚀 Uso

### 1. Ejecutar Migraciones

Al inicializar el schema de un tenant nuevo, la tabla `tenant_config` se creará automáticamente con valores por defecto.

```bash
# Inicializar schema de tenant
curl -X POST http://localhost:8080/api/v1/public/tenants/{tenant_id}/init
```

### 2. Acceder a la Configuración

Una vez autenticado, navega a:
- **Frontend**: http://localhost:5173/settings
- **API**: GET http://localhost:8080/api/v1/tenant/config

### 3. Actualizar Configuración

```bash
curl -X PUT http://localhost:8080/api/v1/tenant/config \
  -H "Authorization: Bearer {token}" \
  -H "X-Tenant-ID: {tenant_id}" \
  -H "Content-Type: application/json" \
  -d '{
    "company_legal_name": "Mi Empresa SRL",
    "company_tax_id": "123456789",
    "invoice_prefix": "FACT",
    "default_tax_rate": 18.00
  }'
```

---

## 🎯 Casos de Uso

### 1. Configuración Inicial del Tenant

Al crear un nuevo tenant, el administrador puede:
1. Configurar datos de la empresa
2. Establecer información fiscal
3. Configurar numeración de facturas
4. Personalizar colores y logo

### 2. Actualización de Datos Fiscales

Cuando cambian los datos de DGII:
1. Actualizar RNC, usuario y API key
2. Cambiar ambiente (sandbox → production)
3. Configurar secuencias de NCF

### 3. Personalización de Marca

Para adaptar el sistema a la marca:
1. Subir logo de la empresa
2. Configurar colores corporativos
3. Establecer formato de fecha/hora

---

## 🔒 Seguridad

- ✅ Rutas protegidas con `AuthMiddleware`
- ✅ Aislamiento por tenant con `TenantMiddleware`
- ✅ Campos sensibles (API keys) en variables de entorno
- ✅ Validación de datos en backend

---

## 📝 Notas

### Tabla Única por Tenant

Cada tenant tiene **una sola fila** en `tenant_config`. La configuración se actualiza mediante UPDATE, no INSERT.

### Valores por Defecto

Al crear el schema de un tenant, se inserta una fila con valores por defecto:

```sql
INSERT INTO tenant_config (company_legal_name) VALUES ('Mi Empresa');
```

### Campos Opcionales

La mayoría de campos son opcionales (`NULL`), excepto:
- `company_legal_name`
- `country`
- `invoice_prefix`
- `default_currency`
- `dgii_environment`
- Colores primario y secundario
- Configuraciones de localización

---

## ✅ Próximos Pasos

### Mejoras Sugeridas

1. **Upload de Logo**
   - Implementar endpoint para subir archivos
   - Integrar con S3 o almacenamiento local

2. **Validación de RNC**
   - Validar formato de RNC dominicano
   - Integrar con API de DGII para verificación

3. **Preview de Colores**
   - Mostrar preview del logo y colores
   - Aplicar colores dinámicamente en el UI

4. **Historial de Cambios**
   - Auditar cambios en configuración
   - Mostrar quién y cuándo modificó

5. **Export/Import**
   - Exportar configuración como JSON
   - Importar configuración de otro tenant

---

## 🧪 Testing

### Backend

```bash
# Test endpoint de configuración
go test ./internal/modules/tenantconfig/...
```

### Frontend

```bash
cd frontend
npm test -- Settings.test.jsx
```

---

## 📚 Documentación API

### GET /api/v1/tenant/config

**Headers:**
```
Authorization: Bearer {access_token}
X-Tenant-ID: {tenant_id}
```

**Response:** 200 OK
```json
{
  "id": "uuid",
  "company_legal_name": "Mi Empresa SRL",
  "company_tax_id": "123456789",
  "invoice_prefix": "FACT",
  "invoice_next_number": 1001,
  "default_currency": "DOP",
  "default_tax_rate": 18.00,
  ...
}
```

### PUT /api/v1/tenant/config

**Headers:**
```
Authorization: Bearer {access_token}
X-Tenant-ID: {tenant_id}
Content-Type: application/json
```

**Body:** (todos los campos opcionales)
```json
{
  "company_legal_name": "Mi Empresa SRL",
  "company_email": "contacto@miempresa.com",
  "invoice_prefix": "FACT",
  "default_tax_rate": 18.00
}
```

**Response:** 200 OK
```json
{
  "id": "uuid",
  "company_legal_name": "Mi Empresa SRL",
  ...
}
```

---

## ✨ Conclusión

Se ha implementado un módulo completo de configuración de tenant que permite a cada empresa:

✅ Gestionar sus datos básicos y fiscales  
✅ Configurar numeración de documentos  
✅ Personalizar la marca (logo y colores)  
✅ Establecer preferencias de notificación  
✅ Configurar opciones de localización  

El sistema está listo para ser utilizado y puede extenderse fácilmente con nuevas funcionalidades.

