Eres un asistente de desarrollo senior. Crea un SaaS de facturación electrónica para República Dominicana con el siguiente alcance técnico y entregables. El objetivo es entregar un monolito modular con backend en Go, base de datos PostgreSQL (multi-tenant con un schema por tenant), y frontend en React. No uses GORM; escribe SQL explícito y/o usa pgx/sqlx. Provee código completo, scripts de migración, ejemplos de llamadas, pruebas y documentación.

Requisitos arquitectónicos y operativos

Arquitectura: Monolito modular (capas: api, services, repositories, modules, jobs). Cada módulo debe encapsular rutas, lógica de negocio, repositorios y migraciones. Diseño preparado para separación futura a microservicios.

Multi-tenant: Cada tenant tendrá su propio schema en Postgres. Debe incluir:

Gestión de tenant (crear, listar, actualizar, eliminar).

Automatización de creación de schema y ejecución de migraciones por tenant.

Middleware que selecciona el schema correcto según subdominio o header (X-Tenant-ID) y ejecuta queries dentro del schema.

Autenticación y permisos:

JWT + refresh tokens.

Roles y permisos (Admin tenant, Accountant, Billing Clerk, Viewer).

RBAC que permite habilitar/inhabilitar acceso a módulos por plan.

Base de datos y migraciones:

Repositorio db/migrations con SQL puro.

Scripts para crear esquema base, tablas compartidas (si corresponde) y tablas por tenant.

Ejemplo de migración para crear tablas: customers, products, invoices, invoice_lines, tax_rates, payments, credit_notes, journal_entries, audit_logs, settings.

Comunicación con autoridades fiscales (DGII):

Módulo de compliance que prepare los documentos electrónicos en el formato requerido (JSON/XML), gestione certificados digitales y firme digitalmente las facturas si es requerido.

Nota: incluye hooks/adapter para integrar con APIs externas (DGII o proveedores autorizados). Implementa la interfaz y un mock adapter para pruebas (no requiere integración real a DGII).

Facturación y módulos esenciales

Tenant Admin (configuración fiscal, certificado, moneda, serie de facturación).

Catálogo: clientes, productos/servicios, unidades, impuestos.

Facturación: crear factura, enviar factura electrónica, anular factura, notas de crédito/débito.

Pagos: registrar cobros, conciliación, integraciones (Stripe).

Punto de venta (POS) básico (opcional como módulo).

Reportes fiscales (libro de ventas, libro de compras) y export CSV/PDF.

Contabilidad básica: asientos, cierre mensual, export a contabilidad (CSV/OFX).

Auditoría y logs inmutables.

Planes y suscripciones (billing): control de módulos por plan, límites (facturas/mes), facturación recurrente (Stripe).

Administración multi-tenant: panel global (superadmin) para gestionar tenants y ver usage.

Middleware y políticas

Middleware que valida que el tenant tiene habilitado el módulo solicitado según su plan (plan_modules).

Rate limiting por tenant y por usuario.

Integraciones

Stripe: crear cargos, suscripciones, webhooks de pago, facturación recurrente y ejemplo de reconciliación.

Mail (envío de comprobantes).

S3 compatible para almacenamiento de PDFs/XML.

Entregables devops

Dockerfile y docker-compose.yml para entorno local (Postgres, Redis, backend, frontend).

Scripts para crear tenant de prueba y ejecutar migraciones.

CI (archivo de ejemplo para GitHub Actions) que ejecute pruebas y migrations lint.

Pruebas

Unit tests para servicios y repositorios.

Tests de integración para endpoints críticos (crear factura, webhook de Stripe, tenant creation).

Documentación

README.md con instrucciones para ejecutar localmente.

API documentation (OpenAPI/Swagger).

Guía de cómo crear un tenant, cómo ejecutar migraciones en un schema nuevo, y cómo configurar Stripe/DGII (señalar que DGII real requiere credenciales).

Requisitos técnicos y librerías recomendadas

Backend (Go): usa fiber para router (elige uno y sé consistente), pgx o sqlx para DB, jwt-go o similar para JWT, fiber para middlewares si aplica. Evitar GORM. 

Frontend (React): CRA, Vite o Next.js (elige Vite para simplicidad), con React Router; panel administrativo y panel tenant. Consumir API REST y mostrar ejemplos de creación de factura, listado de clientes, y panel de configuración fiscal.

Jobs/background: usa Redis + worker (ej. asynq en Go) o una cola simple para firma de facturas, envío de correo y reintentos.

Logging y tracing: structured logging (zap/logrus) y ejemplo de metrics (Prometheus) básico.

Estructura mínima del repo (generar ejemplo)