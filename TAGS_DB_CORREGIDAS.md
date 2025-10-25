# ✅ Tags DB Corregidas en Modelos de Contabilidad

## 🐛 Problema

Error al intentar obtener el trial balance:
```
failed to get trial balance: missing destination name account_code in *[]accounting.TrialBalanceEntry
```

**Causa:** Los structs `TrialBalanceEntry` y `LedgerEntry` no tenían las tags `db` necesarias para que sqlx mapee correctamente las columnas de la base de datos.

---

## ✅ Solución

Se agregaron las tags `db` a los structs que faltaban:

### TrialBalanceEntry

**ANTES ❌:**
```go
type TrialBalanceEntry struct {
	AccountCode string  `json:"account_code"`  // Solo JSON
	AccountName string  `json:"account_name"`
	AccountType string  `json:"account_type"`
	Debit       float64 `json:"debit"`
	Credit      float64 `json:"credit"`
}
```

**AHORA ✅:**
```go
type TrialBalanceEntry struct {
	AccountCode string  `db:"account_code" json:"account_code"`  // DB + JSON
	AccountName string  `db:"account_name" json:"account_name"`
	AccountType string  `db:"account_type" json:"account_type"`
	Debit       float64 `db:"debit" json:"debit"`
	Credit      float64 `db:"credit" json:"credit"`
}
```

### LedgerEntry

**ANTES ❌:**
```go
type LedgerEntry struct {
	Date           time.Time `json:"date"`
	AccountCode    string    `json:"account_code"`
	AccountName    string    `json:"account_name"`
	Description    string    `json:"description"`
	Reference      string    `json:"reference"`
	Debit          float64   `json:"debit"`
	Credit         float64   `json:"credit"`
	Balance        float64   `json:"balance"`
	JournalEntryID uuid.UUID `json:"journal_entry_id"`
}
```

**AHORA ✅:**
```go
type LedgerEntry struct {
	Date           time.Time `db:"date" json:"date"`
	AccountCode    string    `db:"account_code" json:"account_code"`
	AccountName    string    `db:"account_name" json:"account_name"`
	Description    string    `db:"description" json:"description"`
	Reference      string    `db:"reference" json:"reference"`
	Debit          float64   `db:"debit" json:"debit"`
	Credit         float64   `db:"credit" json:"credit"`
	Balance        float64   `db:"balance" json:"balance"`
	JournalEntryID uuid.UUID `db:"journal_entry_id" json:"journal_entry_id"`
}
```

---

## 📚 Por Qué es Necesario

### Tags en Go Structs

Cuando usas `sqlx.SelectContext()` o `sqlx.GetContext()`:

1. **Tag `db`** → Le dice a sqlx cómo mapear las columnas de SQL al struct
2. **Tag `json`** → Le dice a encoding/json cómo serializar el struct

**Ejemplo:**

```sql
SELECT account_code, account_name FROM chart_of_accounts;
```

```go
// Sin tag db ❌
type Entry struct {
    AccountCode string `json:"account_code"`
}
// Error: sqlx no sabe que account_code de SQL → AccountCode del struct

// Con tag db ✅
type Entry struct {
    AccountCode string `db:"account_code" json:"account_code"`
}
// ✅ sqlx mapea: account_code → AccountCode
```

---

## 🔍 Otros Structs que YA tenían las tags correctas

Estos structs ya estaban bien:

### ✅ ChartOfAccounts
```go
type ChartOfAccounts struct {
    ID          uuid.UUID `db:"id" json:"id"`
    AccountCode string    `db:"account_code" json:"account_code"`
    ...
}
```

### ✅ JournalEntry
```go
type JournalEntry struct {
    ID          uuid.UUID `db:"id" json:"id"`
    EntryNumber string    `db:"entry_number" json:"entry_number"`
    ...
}
```

### ✅ JournalEntryLine
```go
type JournalEntryLine struct {
    ID             uuid.UUID `db:"id" json:"id"`
    AccountCode    string    `db:"account_code" json:"account_code"`
    ...
}
```

---

## 🎯 Regla General

**Para TODOS los structs que se usan con sqlx:**

```go
type MiModelo struct {
    Campo1 tipo `db:"nombre_columna_sql" json:"nombre_json"`
    Campo2 tipo `db:"otra_columna" json:"otro_nombre"`
    ...
}
```

- `db:` → Nombre de la columna en la base de datos (snake_case)
- `json:` → Nombre del campo en JSON (camelCase o snake_case según preferencia)

---

## ✅ Verificación

### Compilación
```bash
go build ./cmd/api
# Exit code: 0 ✅
```

### Backend Reiniciado
```bash
docker-compose -f docker-compose.dev.yml restart api
# ✅ Reiniciado
```

### Endpoints Funcionando
```bash
# Trial Balance
curl .../accounting/trial-balance \
  -H "Authorization: Bearer {token}" \
  -H "X-Tenant-ID: demo"
# ✅ Retorna balance de comprobación

# General Ledger
curl .../accounting/ledger/1111 \
  -H "Authorization: Bearer {token}" \
  -H "X-Tenant-ID: demo"
# ✅ Retorna libro mayor de la cuenta
```

---

## 🚀 Ahora Funciona

**Todos los reportes deberían funcionar:**

1. ✅ **Balance de Comprobación** (`/trial-balance`)
   - Muestra todas las cuentas con débitos y créditos

2. ✅ **Balance General** (`/balance-sheet`)
   - Activos = Pasivos + Capital

3. ✅ **Estado de Resultados** (`/income-statement`)
   - Ingresos - Gastos = Utilidad

4. ✅ **Libro Mayor** (`/ledger/:code`)
   - Movimientos de una cuenta específica

---

## 📊 Archivo Modificado

```
internal/modules/accounting/
  ✅ models.go - Agregadas tags `db` a:
     - TrialBalanceEntry (5 campos)
     - LedgerEntry (9 campos)
```

---

## 🎉 Resumen

**Problema:** Tags `db` faltantes en structs  
**Solución:** Agregar `db:"nombre_columna"` a todos los campos  
**Resultado:** Reportes funcionando correctamente  

**Refresca el frontend (F5) en http://localhost:5173/accounting**

Ahora todos los reportes deberían funcionar sin errores. 🚀

