# 🏋️ Guía de Integración Frontend: Módulo de Entrenamiento (Gestrym)

Este documento detalla la estructura y el flujo de trabajo para implementar la parte de entrenamiento deportivo en el frontend de Gestrym.

## 1. 🏗️ Modelos de Datos Clave

### Ejercicio (`Exercise`)
- `id`: Identificador único.
- `name`, `target`, `bodyPart`, `equipment`: Metadatos.
- `gifUrl`: Animación demostrativa.

### Entrenamiento (`Workout`)
- `id`: ID del workout.
- `exercises`: Lista de `WorkoutExercise` con sus respectivos `Sets`.

### Plan de Entrenamiento (`TrainingPlan`)
- `durationDays`: Duración (7, 30 días, etc.).
- `days`: Lista de `TrainingDay`.
- `assignedTo`: ID del usuario asignado.

---

## 2. 🚦 Flujos de Trabajo Principales

### A. Exploración y Búsqueda
- **Endpoint**: `GET /public/exercises`
- **Uso**: Mostrar catálogo con filtros por parte del cuerpo o equipo.

### B. Gestión de Planes (Entrenador)
- **Acción**: Crear plantilla -> Asignar a usuario.
- **Endpoint**: `POST /private/training-plans/:id/assign`.

### C. Seguimiento de Progreso (Usuario)
- **Acción**: Marcar día como completado.
- **Endpoint**: `PATCH /private/training-plans/:id/days/:dayId/complete`.

---

## 3. ⚠️ Manejo de Errores

La API utiliza códigos de estado HTTP estándar. Es importante manejarlos para mejorar la UX:

| Código | Significado | Acción Sugerida en Frontend |
|---|---|---|
| `400 Bad Request` | Datos inválidos (ej: DayNumber fuera de rango). | Mostrar mensaje de error de validación al usuario. |
| `401 Unauthorized`| Token JWT expirado o ausente. | Redirigir al login y limpiar el estado de Auth. |
| `403 Forbidden` | El usuario intenta ver un plan que no le pertenece. | Mostrar pantalla de "Acceso Denegado" o error 403. |
| `404 Not Found` | El Plan o Ejercicio solicitado no existe. | Mostrar "Recurso no encontrado" y botón de volver. |
| `500 Internal Error`| Error inesperado en el servidor. | Mostrar mensaje genérico: "Algo salió mal, intenta más tarde". |

---

## 4. ⏳ Estados de Carga (Loading States)

Para una experiencia premium, implementa los siguientes estados:

- **Skeletons**: Úsalos al cargar la lista de ejercicios (`/exercises`) o el detalle del plan. Evitan el salto de contenido.
- **Spinners con Mensaje**: Para la acción de **Adaptar Plan** (`/adapt`), utiliza un spinner con un texto tipo: *"Nuestra IA está optimizando tu rutina basada en tu progreso..."*. Esta operación puede tardar un poco más al procesar la lógica de clonación.
- **Optimistic Updates**: Al marcar un día como completado, actualiza el check en la UI inmediatamente antes de que la API responda. Si la API falla, revierte el cambio y muestra una notificación (Toast).

---

## 5. 🔌 API Reference (Endpoints)

| Método | Endpoint | Acceso |
|---|---|---|
| `GET` | `/public/exercises` | Público |
| `GET` | `/private/training-plans/user/:userId` | Privado (JWT) |
| `POST` | `/private/training-plans/adapt` | Privado (JWT) |
| `PATCH`| `/private/training-plans/:id/days/:dayId/complete` | Privado (JWT) |
