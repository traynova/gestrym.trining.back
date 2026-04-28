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
- `isAiGenerated`: Boolean que indica si el plan fue creado por inteligencia artificial.

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

| Método | Endpoint | Acceso | Descripción |
|---|---|---|---|
| `GET` | `/public/exercises` | Público | Lista de ejercicios. |
| `GET` | `/private/training-plans/user/:userId` | Privado (JWT) | Obtiene los planes de un usuario. |
| `POST` | `/private/training-plans/:id/assign` | Privado (TRAINER) | Asigna un plan a un usuario y registra el historial. |
| `POST` | `/private/training-plans/:id/clone` | Privado (JWT) | Clona un template y registra el historial de asignación. |
| `POST` | `/private/training-plans/adapt` | Privado (JWT) | Adapta el plan activo usando el motor de IA. |
| `PATCH`| `/private/training-plans/:id/days/:dayId/complete` | Privado (JWT) | Marca un día de entrenamiento como completado. |
| `POST` | `/internal/training-plans/ai` | **Interno (API Key)** | (Uso exclusivo para `ai-service`) Registra un nuevo plan generado por IA. |

---

## 6. 🧠 Flujo de Inteligencia Artificial (Adaptación)

Gestrym incluye un motor de evaluación y adaptación de rutinas basado en IA. El frontend debe manejar este flujo de la siguiente manera:

### Endpoint: `POST /private/training-plans/adapt`
Este endpoint evalúa automáticamente el **último plan asignado al usuario** y toma decisiones basándose en el porcentaje de progreso (`completionRate`):

*   🟢 **Progreso >= 80% (Excelente)**: La IA genera un **nuevo plan clonado** con mayor intensidad (añadiendo "(Adapted - High Intensity)" al nombre) e incrementa la dificultad de los días enfocándose en la sobrecarga progresiva.
*   🟡 **Progreso >= 50% (Aceptable)**: La IA determina que el usuario va por buen camino y recomienda mantener el plan actual sin cambios bruscos. Retorna el plan original.
*   🔴 **Progreso < 50% (Bajo)**: La IA identifica problemas de adherencia y sugiere tomar un descanso o elegir un plan más ligero. Retorna el plan original.

**Manejo en la UI:**
La respuesta siempre incluirá un campo `recommendation` con el mensaje generado por la IA y el objeto `data` (con el plan adaptado o el mismo plan). 
1. Muestra un "Spinner Mágico" o estado de carga tipo "Nuestra IA está evaluando tu rendimiento..." al llamar este endpoint.
2. Muestra la frase contenida en `recommendation` en un Toast, Modal o tarjeta resaltada para darle feedback inmediato al usuario sobre lo que la IA pensó.

---

## 7. 🔗 Historial de Asignaciones y Clonación

Para el rol de Entrenador (Trainer):
- **Clonación de Templates (`/clone`)**: Toma un plan maestro y le genera una copia limpia (con todos los días en `IsCompleted: false`) para un usuario específico.
- **Asignación y Trazabilidad (`/assign`)**: Ambos endpoints ahora registran automáticamente el historial de asignaciones en la base de datos, almacenando quién asignó el plan y la fecha de inicio. La UI de entrenadores podría aprovechar este estado para saber qué rutinas enviaron y cuándo.
