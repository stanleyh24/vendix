# 🔍 Diagnóstico Frontend - Plan de Cuentas No Carga

## ✅ Backend Confirmado OK

El diagnóstico muestra que:
- ✅ Backend corriendo en http://localhost:8080
- ✅ Tabla `chart_of_accounts` existe
- ✅ Tiene 57 cuentas en la base de datos

**El problema está en el FRONTEND** 🎨

---

## 🔍 Paso 1: Abrir Consola del Navegador

1. Ve a http://localhost:5173/accounting
2. Presiona **F12** para abrir DevTools
3. Ve a la pestaña **Console**
4. Refresca la página (F5)
5. Lee los mensajes

---

## 📊 Mensajes que Deberías Ver

### Si TODO está bien:

```
Fetching accounts from /accounting/accounts...
Accounts response: [{...}, {...}, ...]
Loaded 57 accounts
```

### Si hay ERROR de autenticación:

```
Error fetching accounts: Error: Request failed with status code 401
```

**Solución:**
1. Cierra sesión
2. Vuelve a iniciar sesión en http://localhost:5173/login
3. Intenta de nuevo

### Si hay ERROR de red:

```
Error fetching accounts: Error: Network Error
```

**Solución:**
```bash
# El backend no está respondiendo
# Verificar que esté corriendo:
curl http://localhost:8080/health

# Si no responde, iniciar:
make run-api
```

### Si la respuesta no es un array:

```
Response data is not an array: {...}
```

**Solución:** Problema en el backend, verificar logs del servidor

---

## 🔍 Paso 2: Verificar Network Tab

En DevTools:

1. Ve a la pestaña **Network**
2. Refresca la página (F5)
3. Busca la petición `accounts`
4. Haz clic en ella
5. Revisa:

### Headers:
```
Request URL: http://localhost:8080/api/v1/tenant/accounting/accounts
Request Method: GET
Status Code: ???
```

**Debe tener:**
- `Authorization: Bearer xxxxxx...`
- `X-Tenant-ID: demo` (o tu tenant)

### Response:
Debe ser un JSON con array:
```json
[
  {
    "id": "...",
    "account_code": "1000",
    "account_name": "ACTIVOS",
    "account_type": "ASSET",
    "balance": 0,
    ...
  },
  ...
]
```

---

## 🎯 Causas Comunes y Soluciones

### Causa 1: No estás autenticado

**Síntomas:**
- Status Code: 401
- Mensaje: "Unauthorized"

**Verificar:**
```javascript
// En la consola del navegador (F12):
localStorage.getItem('accessToken')
localStorage.getItem('tenantId')
```

**Solución:**
```
1. Si son null → Cierra sesión y vuelve a iniciar
2. Si tienen valores → El token expiró, vuelve a iniciar sesión
```

### Causa 2: Tenant ID incorrecto

**Síntomas:**
- Status Code: 401 o 500
- No se envía el header X-Tenant-ID

**Verificar:**
```javascript
// En consola:
console.log(localStorage.getItem('tenantId'))
```

**Solución:**
```bash
# Ver qué tenants existen:
psql -h localhost -U postgres -d vendix -c "
  SELECT slug FROM public.tenants;
"

# Usa el slug correcto al hacer login
```

### Causa 3: CORS

**Síntomas:**
- Error de CORS en consola
- "Access-Control-Allow-Origin"

**Solución:**
Ya está configurado en el backend, pero verifica que el frontend esté en puerto 5173:
```bash
# En otra terminal
cd frontend
npm run dev
```

### Causa 4: API URL incorrecta

**Verificar:**
```javascript
// En consola del navegador
console.log(import.meta.env.VITE_API_URL)
```

**Solución:**
```bash
# Crear .env en frontend/
echo "VITE_API_URL=http://localhost:8080" > frontend/.env
```

---

## 🧪 Prueba Manual

### En la Consola del Navegador (F12 → Console):

```javascript
// 1. Verificar configuración
console.log('Token:', localStorage.getItem('accessToken'))
console.log('Tenant:', localStorage.getItem('tenantId'))

// 2. Probar manualmente la petición
import api from './lib/api'

api.get('/accounting/accounts')
  .then(response => {
    console.log('Success!', response.data)
    console.log('Number of accounts:', response.data.length)
  })
  .catch(error => {
    console.error('Error:', error)
    console.error('Response:', error.response)
  })
```

---

## 🔧 Soluciones Rápidas

### Solución 1: Reiniciar Todo

```bash
# 1. Detener frontend (Ctrl+C)
# 2. Limpiar caché
cd frontend
rm -rf node_modules/.vite
npm run dev

# 3. En el navegador:
# - Ctrl+Shift+R (hard refresh)
# - O: DevTools → Application → Clear Storage → Clear all
```

### Solución 2: Verificar localStorage

```javascript
// En consola del navegador (F12):

// Ver todo el estado
console.log('Auth State:', {
  accessToken: localStorage.getItem('accessToken'),
  tenantId: localStorage.getItem('tenantId')
})

// Si falta algo, hacer login de nuevo:
// 1. Ve a /login
// 2. Inicia sesión
// 3. Vuelve a /accounting
```

### Solución 3: Probar la API directamente

```bash
# Obtener token
curl -X POST http://localhost:8080/api/v1/tenant/auth/login \
  -H 'Content-Type: application/json' \
  -H 'X-Tenant-ID: demo' \
  -d '{"email":"admin@demo.com","password":"demo123"}' \
  | jq -r '.access_token'

# Copiar el token y probar:
curl http://localhost:8080/api/v1/tenant/accounting/accounts \
  -H "Authorization: Bearer TU_TOKEN_AQUI" \
  -H "X-Tenant-ID: demo" \
  | jq '.[0:3]'  # Muestra las primeras 3 cuentas
```

Si esto funciona, el problema es en el frontend (token o headers).

---

## 📋 Checklist de Diagnóstico

Ejecuta esto en orden:

- [ ] **1. Backend corriendo**
  ```bash
  curl http://localhost:8080/health
  # Debe retornar: {"status":"ok",...}
  ```

- [ ] **2. Frontend corriendo**
  ```bash
  curl http://localhost:5173
  # Debe retornar HTML
  ```

- [ ] **3. Cuentas en BD**
  ```bash
  make check-accounts
  # Debe mostrar: 57+ cuentas
  ```

- [ ] **4. Login funciona**
  - Ir a /login
  - Entrar con credenciales
  - Debe redirigir a /

- [ ] **5. Abrir /accounting**
  - DevTools (F12)
  - Console tab
  - Buscar mensajes de error

- [ ] **6. Network tab**
  - Ver petición a `/accounting/accounts`
  - Verificar Status Code
  - Verificar Headers
  - Verificar Response

---

## 🎯 Pasos para Solucionar AHORA

### Paso 1: Abrir DevTools

1. Ve a http://localhost:5173/accounting
2. Presiona **F12**
3. Tab **Console**
4. Refresca (F5)

### Paso 2: Leer el Error

Busca líneas en rojo que digan:
- "Error fetching accounts:"
- "401" o "403" → Problema de autenticación
- "Network Error" → Backend no responde
- "404" → Ruta incorrecta
- "500" → Error en el backend

### Paso 3: Aplicar la Solución

**Si dice "401 Unauthorized":**
```javascript
// En consola:
localStorage.clear()
// Luego ve a /login y vuelve a iniciar sesión
```

**Si dice "Network Error":**
```bash
# Verificar backend
curl http://localhost:8080/health

# Si no responde, reiniciar:
make run-api
```

**Si dice "tabla no existe":**
```bash
make migrate-tenants
make insert-accounts
```

**Si no muestra ningún error pero está vacía:**
```javascript
// En consola, ejecutar:
api.get('/accounting/accounts').then(r => console.log(r.data))

// Ver qué retorna
```

---

## 💡 Script de Diagnóstico Completo

He mejorado el código con mejor logging. Ahora verás en la consola:

```
Fetching accounts from /accounting/accounts...
Accounts response: [...]
Loaded 57 accounts
```

**Ejecuta esto:**

1. Abre http://localhost:5173/accounting
2. Presiona F12
3. Tab Console
4. Refresca (F5)
5. **Copia TODO lo que aparezca en rojo** y compártelo

---

## 🆘 Solución de Emergencia

Si nada funciona:

```bash
# 1. Detener todo
# Ctrl+C en backend
# Ctrl+C en frontend

# 2. Verificar que todo está en la BD
make check-accounts

# 3. Si falta, insertar:
make insert-accounts

# 4. Reiniciar backend
make run-api

# 5. En otra terminal, reiniciar frontend
cd frontend
npm run dev

# 6. En el navegador:
# - Ctrl+Shift+Delete → Borrar caché
# - O F12 → Application → Clear Storage → Clear all
# - Ir a /login
# - Iniciar sesión
# - Ir a /accounting
# - F12 → Console → Ver mensajes
```

---

## 📞 Información a Compartir

Si el problema persiste, comparte:

1. **Mensaje en consola del navegador** (F12 → Console)
2. **Respuesta del Network** (F12 → Network → accounts → Response)
3. **Headers enviados** (F12 → Network → accounts → Headers)
4. **Resultado de:**
   ```bash
   make check-accounts
   ./scripts/test-accounting-api.sh
   ```

---

## ✅ Resultado Esperado

Cuando funcione correctamente, deberías ver:

**En Console (F12):**
```
Fetching accounts from /accounting/accounts...
Accounts response: Array(57)
Loaded 57 accounts
```

**En la página:**
- Tabla con 57 filas
- Columnas: Código, Nombre, Tipo, Balance, Estado
- Cuentas como: 1000 ACTIVOS, 1111 Caja General, etc.

---

## 🚀 Comando de Verificación Única

```bash
# Ejecutar todo el diagnóstico
./scripts/test-accounting-api.sh
```

Esto te dirá exactamente dónde está el problema.

**Con la información del diagnóstico:**
- Backend: ✅ Corriendo
- Tabla: ✅ Existe
- Datos: ✅ 57 cuentas

**El problema es en el frontend. Ahora abre F12 y mira la consola del navegador para ver el error específico.** 🔍

