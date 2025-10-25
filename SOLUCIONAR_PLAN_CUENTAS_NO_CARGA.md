# 🔧 Solucionar: Plan de Cuentas No Carga en Frontend

## ❌ Problema

La página de Contabilidad carga, pero la lista de cuentas aparece vacía o con error.

---

## 🔍 Diagnóstico Rápido

Ejecuta este script para diagnóstico automático:

```bash
./scripts/test-accounting-api.sh
```

Este script verificará:
1. ✅ Si el backend está corriendo
2. ✅ Si el tenant existe
3. ✅ Si la tabla existe
4. ✅ Si la tabla tiene datos
5. ✅ Si el endpoint responde

---

## 🎯 Causas Comunes y Soluciones

### 1. Backend NO está corriendo ❌

**Síntoma:** Error de red en el frontend, "Network Error" o "Failed to fetch"

**Verificar:**
```bash
curl http://localhost:8080/health
```

**Solución:**
```bash
# Opción A: Con Docker
make dev-up

# Opción B: Local
make run-api

# Opción C: Con air (hot reload)
cd /home/scorpion/desa/Projects/vendix
air
```

### 2. Tabla NO existe ❌

**Síntoma:** Error "relation does not exist" o tabla vacía en el diagnóstico

**Verificar:**
```bash
make check-accounts
```

**Solución:**
```bash
# Crear la tabla
make migrate-tenants
```

### 3. Tabla VACÍA (sin cuentas) ❌

**Síntoma:** Tabla existe con 0 cuentas

**Verificar:**
```bash
make check-accounts
# Muestra: ✅ tenant_demo - Tabla existe con 0 cuentas
```

**Solución:**
```bash
# Insertar las 60 cuentas
make insert-accounts
```

### 4. NO estás autenticado ❌

**Síntoma:** Error 401 Unauthorized en la consola del navegador (F12)

**Verificar:** Abrir consola del navegador (F12) → Tab Network → Ver la petición a `/accounting/accounts`

**Solución:**
```bash
# 1. Cerrar sesión y volver a iniciar sesión
# 2. Verificar que localStorage tiene:
#    - accessToken
#    - tenantId
```

En la consola del navegador:
```javascript
// Verificar
console.log(localStorage.getItem('accessToken'))
console.log(localStorage.getItem('tenantId'))

// Si son null, cierra sesión y vuelve a iniciar
```

### 5. Tenant ID incorrecto ❌

**Síntoma:** Error 401 o respuesta vacía

**Verificar:**
```bash
# Ver qué tenants existen
psql -h localhost -U postgres -d vendix -c "
  SELECT slug, schema_name FROM public.tenants;
"
```

**Solución:** Asegúrate de usar el slug correcto al iniciar sesión

### 6. Problema de CORS ❌

**Síntoma:** Error de CORS en la consola

**Verificar:** Ver en consola si hay mensajes de CORS

**Solución:**
```bash
# Verificar que el frontend esté en puerto 5173
# y el backend en puerto 8080

# El backend ya tiene CORS configurado para "*"
```

---

## 🧪 Pruebas Paso a Paso

### Paso 1: Verificar Backend

```bash
# Debe retornar status: ok
curl http://localhost:8080/health
```

### Paso 2: Verificar Base de Datos

```bash
# Ver cuántas cuentas hay
make check-accounts
```

**Resultado esperado:**
```
✅ tenant_demo - Tabla existe con 60 cuentas
```

**Si muestra 0 cuentas:**
```bash
make insert-accounts
```

### Paso 3: Verificar API con Token

```bash
# 1. Obtener token
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/tenant/auth/login \
  -H 'Content-Type: application/json' \
  -H 'X-Tenant-ID: demo' \
  -d '{"email":"admin@demo.com","password":"demo123"}' \
  | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)

echo "Token: $TOKEN"

# 2. Probar endpoint de cuentas
curl http://localhost:8080/api/v1/tenant/accounting/accounts \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Tenant-ID: demo"
```

**Debe retornar:** JSON con array de cuentas

### Paso 4: Verificar Frontend

1. Abre http://localhost:5173/accounting
2. Presiona F12 para abrir DevTools
3. Ve a la pestaña "Console"
4. Ve a la pestaña "Network"
5. Refresca la página (F5)
6. Busca la petición a `accounting/accounts`
7. Verifica:
   - Status code (debe ser 200)
   - Response (debe tener un array con cuentas)
   - Headers (debe tener Authorization y X-Tenant-ID)

---

## 🔧 Solución Completa

Si nada de lo anterior funciona, ejecuta esto:

```bash
# 1. Verificar estado
make check-accounts

# 2. Si falta tabla: migrar
make migrate-tenants

# 3. Si tabla está vacía: insertar
make insert-accounts

# 4. Verificar que ahora tiene datos
make check-accounts

# 5. Reiniciar backend
# Ctrl+C en la terminal del backend, luego:
make run-api

# 6. Refrescar frontend
# F5 en http://localhost:5173/accounting

# 7. Si aún falla: Cerrar sesión y volver a iniciar
# Ir a /login y entrar de nuevo
```

---

## 📋 Checklist de Diagnóstico

Marca cada ítem:

- [ ] Backend está corriendo (curl http://localhost:8080/health)
- [ ] Tabla existe (make check-accounts)
- [ ] Tabla tiene datos (make check-accounts muestra 60+ cuentas)
- [ ] Puedes autenticarte (login funciona)
- [ ] localStorage tiene accessToken (F12 → Console → `localStorage.getItem('accessToken')`)
- [ ] localStorage tiene tenantId (F12 → Console → `localStorage.getItem('tenantId')`)
- [ ] Network tab muestra petición a /accounting/accounts
- [ ] Petición tiene status 200
- [ ] Response tiene array de cuentas

---

## 🐛 Errores Específicos

### Error: "La tabla de contabilidad no existe"

```bash
make migrate-tenants
```

### Error: "Failed to fetch" o "Network Error"

```bash
# Backend no está corriendo
make run-api
```

### Error: "401 Unauthorized"

```javascript
// En consola del navegador (F12)
// Verificar token
console.log(localStorage.getItem('accessToken'))

// Si es null, cerrar sesión y volver a entrar
```

### Lista vacía sin error

```bash
# Tabla existe pero sin datos
make insert-accounts
```

### Error: "X-Tenant-ID header required"

```javascript
// En consola del navegador (F12)
console.log(localStorage.getItem('tenantId'))

// Si es null, volver a hacer login
```

---

## 💡 Verificación con curl

Si quieres verificar manualmente:

```bash
# 1. Login
curl -X POST http://localhost:8080/api/v1/tenant/auth/login \
  -H 'Content-Type: application/json' \
  -H 'X-Tenant-ID: demo' \
  -d '{"email":"admin@demo.com","password":"demo123"}'

# Copiar el access_token de la respuesta

# 2. Obtener cuentas
curl http://localhost:8080/api/v1/tenant/accounting/accounts \
  -H "Authorization: Bearer {PEGAR_TOKEN_AQUI}" \
  -H "X-Tenant-ID: demo" \
  | jq '.[0:3]'  # Muestra las primeras 3 cuentas
```

---

## 🎯 Solución Más Común

En el 90% de los casos, el problema es:

**La tabla está vacía**

```bash
make insert-accounts
```

Luego refresca el frontend (F5) y debería funcionar.

---

## 📞 Debug Avanzado

Si nada funciona, revisa los logs:

### Backend logs:
```bash
# Ver logs del contenedor (si usas Docker)
docker-compose logs -f api

# O logs directos si corre local
# Verás los logs en la terminal donde ejecutaste make run-api
```

### Frontend logs:
```bash
# Abrir DevTools (F12)
# Pestaña Console
# Buscar errores en rojo
```

### Base de datos:
```bash
# Conectar y verificar manualmente
psql -h localhost -U postgres -d vendix

# Dentro de psql:
\c vendix
SELECT schema_name FROM public.tenants;
SELECT count(*) FROM tenant_demo.chart_of_accounts;
\q
```

---

## ✅ Verificación Final

Cuando todo funcione, deberías ver:

**En el Frontend:**
- ✅ Tabla con 60+ filas
- ✅ Columnas: Código, Nombre, Tipo, Balance, Estado
- ✅ Filtro de búsqueda funcionando
- ✅ Cuentas por tipo con colores (azul=activo, rojo=pasivo, etc.)

**Primeras cuentas visibles:**
```
1000 - ACTIVOS
1100 - Activo Circulante
1110 - Caja y Bancos
1111 - Caja General
1112 - Banco - Cuenta Corriente
...
```

---

## 🚀 Comando Único de Solución

Si quieres solucionar todo de una vez:

```bash
# Esto ejecuta todo el proceso de corrección
make check-accounts && \
make migrate-tenants && \
make insert-accounts && \
make check-accounts && \
echo "✅ Todo listo! Refresca el frontend (F5)"
```

---

¿Aún no funciona? Ejecuta el diagnóstico y comparte la salida:

```bash
./scripts/test-accounting-api.sh > diagnostico.txt
cat diagnostico.txt
```

