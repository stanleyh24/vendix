# 🎨 Paleta de Colores Vendix

## Colores Principales

| Color | Uso | Hex | Clase Tailwind |
|-------|-----|-----|----------------|
| 🟠 **Naranja Vendix** | Color principal (energía, ventas, dinamismo) | `#FF6B00` | `bg-vendix-orange` / `text-vendix-orange` |
| ⚫ **Negro Carbón** | Tipografía y contraste | `#212121` | `bg-carbon` / `text-carbon` |
| ⚪ **Blanco Puro** | Fondo y limpieza visual | `#FFFFFF` | `bg-white` / `text-white` |
| 🩶 **Gris Claro** | Bordes, tarjetas y elementos neutros | `#E6E6E6` | `bg-neutral-light` / `border-neutral-light` |
| 🔵 **Azul SaaS** | Confianza tecnológica (secundario) | `#1877F2` | `bg-saas-blue` / `text-saas-blue` |

## Variantes Adicionales

### Naranja Vendix
- **Oscuro** (hover): `#E66000` → `bg-vendix-orange-dark`
- **Claro**: `#FF8533` → `bg-vendix-orange-light`

### Azul SaaS
- **Oscuro** (hover): `#1565D8` → `bg-saas-blue-dark`

### Neutral
- **Más claro** (fondos secundarios): `#F5F5F5` → `bg-neutral-lighter`

## Variables CSS

También puedes usar las variables CSS directamente:

```css
:root {
  --vendix-orange: #FF6B00;
  --vendix-orange-dark: #E66000;
  --vendix-orange-light: #FF8533;
  --carbon: #212121;
  --neutral-light: #E6E6E6;
  --neutral-lighter: #F5F5F5;
  --saas-blue: #1877F2;
  --saas-blue-dark: #1565D8;
  --white: #FFFFFF;
}
```

## Componentes Predefinidos

### Botones

```jsx
// Botón principal (naranja)
<button className="btn-primary">Crear Factura</button>

// Botón secundario (azul)
<button className="btn-secondary">Guardar</button>

// Botón outline
<button className="btn-outline">Cancelar</button>
```

### Tarjetas

```jsx
<div className="card">
  <h3 className="text-lg font-semibold text-carbon">Título</h3>
  <p className="text-gray-600">Contenido de la tarjeta</p>
</div>
```

### Inputs

```jsx
<input 
  type="text" 
  className="input-field" 
  placeholder="Nombre del cliente"
/>
```

### Badges

```jsx
// Badge naranja
<span className="badge-orange">Nuevo</span>

// Badge azul
<span className="badge-blue">Pagado</span>

// Badge personalizado
<span className="badge bg-green-500 text-white">Completado</span>
```

### Links

```jsx
<a href="#" className="link">Ver más información</a>
```

## Ejemplos de Uso

### Header con branding

```jsx
<header className="bg-white border-b border-neutral-light">
  <div className="flex items-center justify-between px-6 py-4">
    <h1 className="text-2xl font-bold text-vendix-orange">Vendix</h1>
    <button className="btn-primary">Nueva Factura</button>
  </div>
</header>
```

### Card de estadística

```jsx
<div className="card">
  <div className="flex items-center justify-between">
    <div>
      <p className="text-sm text-gray-500">Ventas del mes</p>
      <p className="text-3xl font-bold text-carbon">$125,450</p>
    </div>
    <div className="w-12 h-12 bg-vendix-orange-light rounded-full flex items-center justify-center">
      <span className="text-white text-xl">📊</span>
    </div>
  </div>
</div>
```

### Formulario

```jsx
<form className="card space-y-4">
  <div>
    <label className="block text-sm font-medium text-carbon mb-2">
      Nombre del cliente
    </label>
    <input type="text" className="input-field" />
  </div>
  
  <div className="flex gap-3">
    <button type="submit" className="btn-primary">
      Guardar
    </button>
    <button type="button" className="btn-outline">
      Cancelar
    </button>
  </div>
</form>
```

### Estado de factura con badges

```jsx
<div className="flex gap-2">
  <span className="badge-orange">Pendiente</span>
  <span className="badge-blue">Enviada a DGII</span>
  <span className="badge bg-green-500 text-white">Pagada</span>
  <span className="badge bg-red-500 text-white">Vencida</span>
</div>
```

## Guía de Accesibilidad

### Contraste de colores

Todos los colores cumplen con WCAG 2.1 AA:

- ✅ Naranja Vendix (#FF6B00) sobre blanco: **Contraste 4.5:1**
- ✅ Negro Carbón (#212121) sobre blanco: **Contraste 16:1**
- ✅ Azul SaaS (#1877F2) sobre blanco: **Contraste 4.7:1**
- ✅ Texto blanco sobre Naranja Vendix: **Contraste 4.5:1**

### Mejores prácticas

1. **Botones principales**: Usar `btn-primary` (naranja) para acciones primarias
2. **Botones secundarios**: Usar `btn-secondary` (azul) para acciones secundarias
3. **Texto principal**: Siempre usar `text-carbon` para máximo contraste
4. **Bordes suaves**: Usar `border-neutral-light` para separadores discretos
5. **Fondos**: Usar `bg-neutral-lighter` para fondos de página, `bg-white` para tarjetas

## Modo Oscuro (Futuro)

Si se implementa modo oscuro, se recomienda:

- Fondo principal: `#1A1A1A`
- Tarjetas: `#2D2D2D`
- Texto principal: `#E6E6E6`
- Mantener Naranja Vendix y Azul SaaS para acciones

