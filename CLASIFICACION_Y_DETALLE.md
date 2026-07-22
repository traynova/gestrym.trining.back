# 📊 Clasificación Detallada y Auditoría del Repositorio (`gestrym.trining.back`)

Este documento ofrece una **clasificación exhaustiva, arquitectura y mapa de componentes** del microservicio de **Entrenamiento y Nutrición de Gestrym**. Su objetivo es brindar a desarrolladores e inteligencias artificiales una visión completa de lo que contiene el repositorio, cómo interactúan sus piezas y dónde reside cada responsabilidad.

---

## 📌 1. Clasificación General del Proyecto

- **Nombre del Microservicio:** `gestrym-training` (Repositorio: `gestrym.trining.back`)
- **Lenguaje / Versión:** Go (Golang) 1.22+
- **Propósito:** Gestionar el catálogo de ejercicios, rutinas de entrenamiento, planes multidía, seguimiento del progreso del usuario, catálogo de alimentos, cálculo de planes nutricionales, canalizaciones ETL de ingesta masiva y comunicación con el microservicio de archivos (`file-service`).
- **Estilo de Arquitectura:** **Arquitectura Hexagonal (Puertos y Adaptadores)** con un modelo de persistencia/dominio compartido denominado `common/models`.

---

## 🗂️ 2. Mapa y Clasificación de Componentes por Capas

Below is the structured classification of all directories and files within the codebase:

```
src/
├── app.go                          # Inicializador del servidor Gin, DI y montaje de rutas
├── common/                         # Núcleo compartido entre módulos
│   ├── config/                     # Carga de configuraciones desde YAML / ENV con Viper
│   ├── middleware/                 # Control de acceso (JWT, RBAC)
│   ├── models/                     # ÚNICA FUENTE DE VERDAD: Modelos DB y DTOs transversales
│   ├── routes/                     # Enrutador central Gin (/gestrym-training)
│   ├── shared/                     # Adaptadores cliente HTTP externos (FileStorageAdapter)
│   └── utils/                      # Formateador de respuestas HTTP JSON (response.go)
├── training/                       # Módulo de Entrenamientos y Planes (Hexagonal)
│   ├── application/                # Casos de Uso y Mapeadores DTO
│   ├── domain/                     # Contratos/Interfaces de repositorios
│   ├── infrastructure/             # Implementaciones de Repositorios GORM
│   └── interfaces/                 # Controladores HTTP (Gin Handlers)
└── nutrition/                      # Módulo de Nutrición y Alimentos (Hexagonal)
    ├── application/                # Casos de Uso (Búsqueda, Importación, Generación)
    ├── domain/                     # Contratos/Interfaces de repositorios
    ├── infrastructure/             # Implementaciones de Repositorios GORM
    └── interfaces/                 # Controladores HTTP (Gin Handlers)

internal/
└── etl/                            # Pipeline ETL de ingesta masiva (USDA + Pexels + MinIO)
    ├── extractor/                  # Cliente HTTP para USDA FoodData Central
    ├── transformer/                # Normalizador de macros y buscador Pexels
    ├── loader/                     # Carga masiva en DB e imágenes a MinIO
    └── pipeline/                   # Worker pool concurrente (Goroutines/Channels)

cmd/
├── etl-foods/                      # Entrypoint CLI para ejecutar el pipeline de alimentos
└── import_exercises/               # Entrypoint CLI para importación de ejercicios
```

---

## 🧩 3. Detalle de Archivos por Módulo

### 🏢 A. Modelo de Datos Centralizado (`src/common/models/`)
Todos los modelos de la base de datos están unificados en este directorio para actuar como entidad única entre las capas de dominio, persistencia y API:

1. **`Exercise.go`**: Entidad de catálogo de ejercicios (nombre, grupo muscular, equipo, GIFs/videos de demostración, `CollectionID`).
2. **`Workout.go`**: Agrupador de ejercicios para una sesión (ej. "Día de Pecho y Tríceps").
3. **`WorkoutExercise.go`**: Tabla intermedia que vincula un `Workout` con un `Exercise`, especificando orden y descanso.
4. **`WorkoutSet.go`**: Registra cada serie dentro de un ejercicio (repeticiones, peso en kg/lbs, RPE, tipo de serie).
5. **`TrainingPlan.go`**: Plan de entrenamiento estructurado (duración en días, `CreatedBy`, `AssignedTo`, flag `IsTemplate`, tags de IA).
6. **`TrainingDay.go`**: Vincula un día específico dentro de un `TrainingPlan` con un `Workout` (ej. Día 1, Día 2), incluyendo notas y estado `IsCompleted`.
7. **`TrainingPlanAssignment.go`**: Historial de asignaciones de planes hechas por entrenadores a clientes.
8. **`FoodCategory.go`**: Categorías de alimentos (Frutas, Carnes, Cereales, etc.).
9. **`Food.go`**: Alimentos del catálogo con sus valores nutricionales por 100g (Calorías, Proteínas, Carbohidratos, Grasas), `ImageURL` y `CollectionID`.
10. **`NutritionPlan.go`**: Plan de objetivos nutricionales del usuario (calorías objetivo, gramos de proteína, carbohidratos y grasas según su meta).

---

### 🏋️ B. Módulo de Entrenamientos (`src/training/`)

#### 🛠️ Casos de Uso (`src/training/application/usecases/`)
- **`CreateTrainingPlanUseCase.go`**: Crea la cabecera de un plan de entrenamiento (plantilla o personalizado).
- **`AddTrainingDayUseCase.go`**: Añade un día estructurado (`TrainingDay`) asociándolo a una rutina (`Workout`).
- **`AssignTrainingPlanUseCase.go`**: Asigna un plan a un cliente. Si el plan ya está asignado a otro usuario, activa un **autoclonado inteligente** para evitar colisiones de datos.
- **`CloneTrainingPlanUseCase.go`**: Realiza una copia profunda (*deep clone*) de un plan y de todos sus días y rutinas vinculadas.
- **`AdaptTrainingPlanUseCase.go`**: Analiza el avance del usuario en su plan actual. Si supera el **80% de días completados**, clona el plan e incrementa la dificultad para sobrecarga progresiva.
- **`UpdateDayCompletionUseCase.go`**: Marca o desmarca un día de entrenamiento como completado (`IsCompleted = true/false`).
- **`CreateTrainingPlanFromAIUseCase.go`**: Procesa e inserta en la base de datos un plan recibido desde un servicio externo de IA.
- **`GetTrainingPlanUseCase.go`**: Obtiene un plan por ID con validación estricta de RBAC (los clientes solo ven sus propios planes).
- **`GetUserTrainingPlansUseCase.go`**: Retorna el listado de planes asignados a un cliente específico.
- **`GetWorkoutFullUseCase.go`**: Ensambla una rutina completa con sus ejercicios y series optimizada para renderizado en el frontend.
- **`ImportExercisesUseCase.go`**: Triggers de importación desde APIs de terceros (ExerciseDB).
- **`TrainingPlanMapper.go`**: Transforma entidades GORM a estructuras DTO limpias para las respuestas JSON de la API.

#### 🗄️ Repositorios (`src/training/infrastructure/repositories/`)
- **`TrainingPlanRepositoryImpl.go`**: Consultas GORM con `Preload` para cargar planes junto con sus días y rutinas asociadas sin incurrir en problemas N+1.
- **`TrainingDayRepositoryImpl.go`**: Operaciones CRUD para los días individuales de un plan.
- **`ExerciseRepositoryImpl.go`**: Búsquedas, filtrado por músculo/target y paginación de ejercicios.
- **`WorkoutRepositoryImpl.go`**: Persistencia de rutinas y ejercicios de rutinas.

#### 🌐 HTTP Handlers (`src/training/interfaces/http/handlers/`)
- **`TrainingPlanHandler.go`**: Expone los endpoints de planes, asignación, clonación, marcado de días, adaptación por IA y comunicación con el AI Service.
- **`ExerciseHandler.go`**: Expone las rutas de consulta e importación del catálogo de ejercicios.
- **`WorkoutHandler.go`**: Expone endpoints para consultar la estructura completa de rutinas.

---

### 🥗 C. Módulo Nutricional (`src/nutrition/`)

#### 🛠️ Casos de Uso (`src/nutrition/application/usecases/`)
- **`GenerateNutritionPlanUseCase.go`**: Calcula las necesidades calóricas diarias (TDEE) mediante la ecuación de **Mifflin-St Jeor** tomando edad, peso, altura y objetivo (`weight_loss`, `muscle_gain`, `maintenance`), distribuyendo los macronutrientes (2g proteína/kg, 0.8g grasa/kg y el resto en carbohidratos).
- **`SearchFoodsUseCase.go`**: Búsqueda paginada en el catálogo de alimentos con filtros por texto.
- **`GetFoodByIDUseCase.go`**: Obtiene la ficha nutricional completa de un alimento por su ID.
- **`ImportFoodsWithImagesUseCase.go`**: Ejecuta la ingesta manual de alimentos vinculando imágenes mediante Pexels.

#### 🌐 HTTP Handlers (`src/nutrition/interfaces/http/handlers/`)
- **`NutritionPlanHandler.go`**: Endpoint `POST /private/nutrition-plans/generate`.
- **`FoodHandler.go`**: Endpoints `GET /public/foods`, `GET /public/foods/:id` y `POST /public/foods/import`.

---

### ⚡ D. Pipeline ETL y Scripts CLI (`internal/etl/` y `cmd/`)

Ubicado en `internal/etl`, es una arquitectura distribuida de procesamiento en lote para ingesta de datos a gran escala:

1. **`extractor/usda_extractor.go`**: Realiza peticiones paginadas a la API oficial de la **USDA FoodData Central**, extrayendo información nutricional cruda.
2. **`transformer/food_transformer.go`**:
   - Limpia y normaliza los nombres de los alimentos.
   - Extrae macronutrientes ajustados a porciones estándar de 100 gramos.
   - Consume **Pexels API** para buscar la mejor foto en alta resolución correspondiente al alimento.
3. **`loader/db_loader.go`**:
   - Transmite las imágenes descargadas directamente hacia el `file-service` (MinIO) en streaming.
   - Realiza inserciones o actualizaciones (*upsert*) en PostgreSQL vía GORM garantizando que no existan duplicados por nombre.
4. **`pipeline/pipeline.go`**:
   - Implementa el patrón **Worker Pool** utilizando Goroutines y Channels en Go.
   - Mantiene una cola concurrente de trabajos de transformación y carga masiva para maximizar el rendimiento de la red y la BD.
5. **`cmd/etl-foods/main.go`**: Punto de entrada CLI para ejecutar este proceso en segundo plano.

---

## 🔐 4. Esquema de Seguridad y Roles (RBAC)

El acceso a la API está regulado en `src/common/middleware/`:

- **JWT Middleware (`jwt_middleware.go`):** Valida la firma del token `Bearer` presente en la cabecera `Authorization`, inyectando el `UserID` y su `RoleID` en el contexto de Gin.
- **Roles Definidos:**
  - **`RoleAdmin` (ID: 1):** Acceso total al sistema.
  - **`RoleCoach` (ID: 2):** Puede crear plantillas globalmente y asignar planes de entrenamiento a cualquier usuario.
  - **`RoleCliente` (ID: 4):** Solo puede visualizar y actualizar sus propios planes y generar sus metas nutricionales.
- **Protección Interna:** Los endpoints bajo `/internal/` requieren la cabecera `X-API-Key` para autenticar comunicaciones de servicio a servicio (por ejemplo, llamadas desde `file-service` o el microservicio de IA).

---

## ⚙️ 5. Flujos de Trabajo Típicos

### 🔄 Flujo 1: Asignación y Adaptación de un Plan de Entrenamiento
```
[Cliente / Entrenador] 
        │
        ▼ (POST /private/training-plans/:id/assign)
[TrainingPlanHandler] ──> Valida JWT y Rol (Entrenador/Admin)
        │
        ▼
[AssignTrainingPlanUseCase] ──> ¿El plan pertenece a otro usuario?
        │                             │
        ├── SI                        └── NO
        ▼                                 ▼
[CloneTrainingPlanUseCase] ────────> Asigna directamente
(Copia profunda de Plan + Días)
        │
        ▼
[PostgreSQL Database] (Guarda asignación en common/models)
```

### 🤖 Flujo 2: Adaptación de Plan por IA (Sobrecarga Progresiva)
```
[Cliente de React] ──> POST /private/training-plans/adapt
        │
        ▼
[AdaptTrainingPlanUseCase] ──> Calcula % de días completados
        │
        ├── Si % < 80% ──> Devuelve mensaje de continuar plan actual
        │
        └── Si % >= 80% ──> Invoca CloneTrainingPlanUseCase()
                                  │
                                  ▼
                         Crea versión "Adaptada" 
                         con incremento de volumen/intensidad
```

---

## 📝 6. Conclusión y Buenas Prácticas para Futuros Desarrollos

1. **Garantizar la regla de `common/models`:** Nunca instanciar modelos de dominio duplicados en capas internas si el struct en `common/models` satisface la estructura de la base de datos.
2. **Consultas Eficientes:** Al agregar relaciones entre `TrainingPlan`, `TrainingDay` y `Workout`, utilizar siempre `Preload` en las consultas de GORM para evitar la degradación por consultas N+1.
3. **Documentación Swagger:** Cada nuevo handler creado en Gin debe contener sus anotaciones `@Summary`, `@Description`, `@Tags`, `@Param` y `@Success` para mantener la especificación de Swagger siempre actualizada.
