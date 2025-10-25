# 💰 Insertar Cuentas Contables

## 📋 Problema

La tabla `chart_of_accounts` existe pero está vacía. No se insertaron las cuentas al crear la tabla.

---

## ✅ Solución Rápida

### Opción 1: Insertar en TODOS los Tenants (Recomendado)

```bash
./scripts/insert-accounts-all-tenants.sh
```

**Esto va a:**
1. Buscar todos los tenants
2. Verificar si tienen la tabla `chart_of_accounts`
3. Si la tabla está vacía, insertar las 60+ cuentas
4. Si ya tiene cuentas, omitir ese tenant
5. Mostrar resumen

**Salida esperada:**
```
💰 Insertando cuentas en tenants sin datos...
==================================================
📋 Schemas de tenants encontrados:
tenant_demo

🔄 Procesando: tenant_demo
   ✅ Insertadas 60 cuentas exitosamente

==================================================
📊 Resumen:
   Total procesados: 1
   Exitosos: 1
   Omitidos: 0
   Fallidos: 0

✅ Cuentas insertadas correctamente en 1 tenant(s)
```

### Opción 2: Insertar en UN Tenant Específico

```bash
# Identificar el schema del tenant
psql -h localhost -U postgres -d vendix -c "
  SELECT schema_name FROM public.tenants;
"

# Insertar cuentas en ese tenant
psql -h localhost -U postgres -d vendix \
  -v schema=tenant_demo \
  -f scripts/insert-chart-of-accounts.sql
```

---

## 🔍 Verificar Antes

Primero verifica si las cuentas están o no:

```bash
./scripts/check-accounting-tables.sh
```

**Salida si falta insertar:**
```
✅ tenant_demo - Tabla existe con 0 cuentas
```

**Salida si ya están:**
```
✅ tenant_demo - Tabla existe con 60 cuentas
```

---

## 📝 ¿Por qué pasó esto?

Posibles causas:

1. **La migración se ejecutó parcialmente**
   - La tabla se creó pero el INSERT falló
   - Posible error en la transacción

2. **Se usó un script viejo**
   - El script de migración no incluía el INSERT

3. **Error de permisos**
   - El usuario de BD no tenía permisos para INSERT

---

## 🎯 Pasos Detallados

### 1. Verificar Estado Actual

```bash
# Ver cuántas cuentas hay en cada tenant
psql -h localhost -U postgres -d vendix << 'EOF'
SELECT 
    t.schema_name,
    COALESCE((
        SELECT count(*)::text 
        FROM information_schema.tables 
        WHERE table_schema = t.schema_name 
        AND table_name = 'chart_of_accounts'
    ), 'tabla no existe') as tabla_existe,
    COALESCE((
        SELECT count(*)::text 
        FROM (
            SELECT 1 
            FROM information_schema.tables 
            WHERE table_schema = t.schema_name 
            AND table_name = 'chart_of_accounts'
        ) x,
        LATERAL (
            SELECT count(*) as cnt 
            FROM chart_of_accounts
        ) y
        WHERE y.cnt IS NOT NULL
    ), '0') as num_cuentas
FROM public.tenants t
WHERE t.deleted_at IS NULL;
EOF
```

### 2. Insertar Cuentas

```bash
# Para todos los tenants
./scripts/insert-accounts-all-tenants.sh

# O para uno específico
psql -h localhost -U postgres -d vendix \
  -v schema=tenant_demo \
  -f scripts/insert-chart-of-accounts.sql
```

### 3. Verificar que Funcionó

```bash
# Opción A: Con script
./scripts/check-accounting-tables.sh

# Opción B: Manual
psql -h localhost -U postgres -d vendix -c "
  SELECT count(*) as total_cuentas 
  FROM tenant_demo.chart_of_accounts;
"
```

**Debe mostrar:** 60 cuentas

### 4. Verificar por Tipo

```bash
psql -h localhost -U postgres -d vendix -c "
  SELECT 
    account_type,
    count(*) as cantidad
  FROM tenant_demo.chart_of_accounts
  GROUP BY account_type
  ORDER BY account_type;
"
```

**Debe mostrar algo como:**
```
 account_type | cantidad 
--------------+----------
 ASSET        |       18
 EQUITY       |        4
 EXPENSE      |       26
 LIABILITY    |       10
 REVENUE      |        5
```

---

## 🧪 Probar en el Frontend

Después de insertar:

1. Ve a http://localhost:5173/accounting
2. Selecciona la pestaña "Plan de Cuentas"
3. Deberías ver una lista con 60+ cuentas organizadas por tipo

**Cuentas que deberías ver:**
```
1000 - ACTIVOS
1111 - Caja General
1112 - Banco - Cuenta Corriente
2121 - ITBIS por Pagar
3100 - Capital Social
4110 - Ventas
5100 - Costo de Ventas
6110 - Sueldos y Salarios
... y 52 más
```

---

## 🐛 Solución de Problemas

### Error: "relation does not exist"

**Causa:** La tabla no existe

**Solución:** Primero crear la tabla
```bash
./scripts/migrate-all-tenants-sql.sh
```

Luego insertar cuentas:
```bash
./scripts/insert-accounts-all-tenants.sh
```

### Error: "duplicate key value"

**Causa:** Las cuentas ya existen

**Solución:** No es necesario hacer nada, las cuentas ya están

```bash
# Verificar
psql -h localhost -U postgres -d vendix -c "
  SELECT count(*) FROM tenant_demo.chart_of_accounts;
"
```

### Error: "permission denied"

**Causa:** El usuario no tiene permisos

**Solución:** Usar usuario con permisos (postgres)
```bash
# Con usuario postgres
psql -h localhost -U postgres -d vendix \
  -v schema=tenant_demo \
  -f scripts/insert-chart-of-accounts.sql
```

### Las cuentas se insertan pero no aparecen

**Causa:** Puede ser un problema de search_path o schema incorrecto

**Solución:** Verificar el schema correcto
```bash
# Ver el schema_name del tenant
psql -h localhost -U postgres -d vendix -c "
  SELECT id, name, schema_name FROM public.tenants;
"

# Usar el schema_name exacto
psql -h localhost -U postgres -d vendix -c "
  SELECT count(*) FROM [schema_name_aqui].chart_of_accounts;
"
```

---

## 🔄 Reinsertar (Empezar de Nuevo)

Si quieres borrar las cuentas existentes y volver a insertar:

```bash
# ⚠️ ADVERTENCIA: Esto borrará todas las cuentas y asientos

# 1. Borrar cuentas existentes
psql -h localhost -U postgres -d vendix -c "
  TRUNCATE TABLE tenant_demo.chart_of_accounts CASCADE;
"

# 2. Insertar nuevamente
psql -h localhost -U postgres -d vendix \
  -v schema=tenant_demo \
  -f scripts/insert-chart-of-accounts.sql
```

---

## 📊 Resumen del Plan de Cuentas

El script inserta **60 cuentas** organizadas así:

| Tipo | Cantidad | Rango de Códigos |
|------|----------|------------------|
| **ACTIVOS** | 18 | 1000-1999 |
| - Caja y Bancos | 4 | 1110-1113 |
| - Cuentas por Cobrar | 3 | 1120-1122 |
| - Inventarios | 2 | 1130-1131 |
| - Activo Fijo | 5 | 1211-1215 |
| **PASIVOS** | 10 | 2000-2999 |
| - Cuentas por Pagar | 2 | 2110-2112 |
| - Impuestos | 2 | 2121-2122 |
| - Otros Pasivos | 2 | 2131, 2210 |
| **CAPITAL** | 4 | 3000-3999 |
| **INGRESOS** | 5 | 4000-4999 |
| **GASTOS** | 26 | 5000-6999 |
| - Costos | 2 | 5000-5100 |
| - Gastos Admin | 7 | 6110-6150 |
| - Gastos Ventas | 3 | 6210-6220 |
| - Gastos Financieros | 3 | 6310-6320 |

---

## ✅ Checklist

- [ ] 1. Verificar que la tabla existe: `./scripts/check-accounting-tables.sh`
- [ ] 2. Si la tabla tiene 0 cuentas, ejecutar: `./scripts/insert-accounts-all-tenants.sh`
- [ ] 3. Verificar que se insertaron: `./scripts/check-accounting-tables.sh` (debe mostrar 60+)
- [ ] 4. Refrescar el frontend y acceder a `/accounting`
- [ ] 5. Verificar que se muestran las cuentas en la pestaña "Plan de Cuentas"

---

## 🎉 Resultado Final

Después de ejecutar el script, deberías tener:

✅ **60 cuentas** en la base de datos  
✅ **Plan de cuentas completo** visible en el frontend  
✅ **Listo para crear asientos** contables  
✅ **Listo para generar reportes** financieros  

---

## 💡 Prevención

Para evitar este problema en el futuro:

1. **Al crear nuevos tenants**, usa el endpoint `/init`:
   ```bash
   curl -X POST http://localhost:8080/api/v1/public/tenants/{ID}/init
   ```

2. **Después de hacer pull**, verifica migraciones:
   ```bash
   ./scripts/check-accounting-tables.sh
   ```

3. **Usa siempre los scripts oficiales** para migraciones

---

¡Listo! Ahora deberías tener tu plan de cuentas completo y funcional. 🚀

