# 🏋️‍♂️ Gestrym - Training & Nutrition Microservice

[![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Framework](https://img.shields.io/badge/Framework-Gin%20Gonic-008080?style=flat&logo=go)](https://gin-gonic.com/)
[![ORM](https://img.shields.io/badge/ORM-GORM-blue?style=flat)](https://gorm.io/)
[![Database](https://img.shields.io/badge/Database-PostgreSQL-336791?style=flat&logo=postgresql)](https://www.postgresql.org/)
[![Storage](https://img.shields.io/badge/Storage-MinIO%20%2F%20S3-C72C48?style=flat&logo=minio)](https://min.io/)
[![API Docs](https://img.shields.io/badge/Swagger-OpenAPI%203.0-85EA2D?style=flat&logo=swagger)](https://swagger.io/)

**Gestrym Training & Nutrition Service** es un microservicio robusto desarrollado en **Golang** bajo **Arquitectura Hexagonal (Puertos y Adaptadores)**. Constituye el núcleo motor de la plataforma de fitness **Gestrym**, encargado de gestionar rutinas de entrenamiento, planes semanales/mensuales, catálogo de ejercicios, seguimiento nutricional (cálculo de macros y TDEE), importación masiva mediante pipelines ETL concurrentes e integración inteligente con servicios de almacenamiento multimedia y generación por IA.

---

## 📌 Tabla de Contenidos

1. [✨ Características Principales](#-características-principales)
2. [🏗️ Arquitectura del Proyecto](#️-arquitectura-del-proyecto)
3. [📁 Estructura del Repositorio](#-estructura-del-repositorio)
4. [🛠️ Tecnologías y Dependencias](#️-tecnologías-y-dependencias)
5. [⚙️ Variables de Entorno](#️-variables-de-entorno)
6. [🚀 Instalación y Ejecución](#-instalación-y-ejecución)
7. [🥗 Pipeline ETL Concurrentes (Alimentos)](#-pipeline-etl-concurrentes-alimentos)
8. [🌐 API Endpoints y Documentación Swagger](#-api-endpoints-y-documentación-swagger)
9. [🔐 Seguridad y Autenticación (RBAC)](#-seguridad-y-autenticación-rbac)
10. [🤖 Integración con IA y Adaptación Inteligente](#-integración-con-ia-y-adaptación-inteligente)

---

## ✨ Características Principales

### 🏋️ Module de Entrenamiento (Training)
- **Planes de Entrenamiento Multidía (7 a 30 días):** Creación, personalización, clonación profunda y asignación de planes por parte de entrenadores o por los propios usuarios.
- **Seguimiento de Progreso:** Marcado de días completados (`IsCompleted`) y cálculo dinámico de tasa de avance.
- **Adaptación Progresiva por IA:** Algoritmo que analiza el cumplimiento (>80%) y clona automáticamente los planes incrementando la sobrecarga progresiva.
- **Gestión de Rutinas (Workouts) y Series (Sets):** Estructura jerárquica con precarga eficiente (`Preload`) para evitar problemas de consultas N+1.
- **Catálogo de Ejercicios:** Filtros por grupo muscular (`bodyPart`), objetivo (`target`), e importación desde ExerciseDB / RapidAPI.

### 🥗 Módulo Nutricional (Nutrition)
- **Catálogo Masivo de Alimentos:** Búsqueda paginada y detallada de información nutricional (Calorías, Proteínas, Carbohidratos, Grasas).
- **Generación de Planes Nutricionales:** Cálculo automatizado de TDEE y distribución de macronutrientes basado en la fórmula de Mifflin-St Jeor según objetivos (Pérdida de peso, Ganancia muscular, Mantenimiento).
- **Integración Multimedia:** Vinculación automática con imágenes de alta calidad obtenidas vía Pexels API y almacenadas en MinIO.

### ⚡ Pipeline ETL (Extract, Transform, Load)
- Pipeline batch en Go (`/internal/etl`) para importar masivamente el catálogo de alimentos de la **USDA FoodData Central**, enriquecerlo con imágenes y realizar *upserts* en PostgreSQL con concurrencia mediante Worker Pools y reintentos con backoff exponencial.

### 📦 Integración con File Service
- Comunicación interna mediante `FileStorageAdapter` con el `file-service` de Gestrym para la gestión unificada de archivos mediante `CollectionID`.

---

## 🏗️ Arquitectura del Proyecto

El proyecto implementa una variante optimizada de **Arquitectura Hexagonal (Ports & Adapters)** orientada a microservicios:

```
                          +-----------------------------------+
                          |        HTTP Interfaces            |
                          |  (Gin Handlers & Middlewares)     |
                          +-----------------+-----------------+
                                            |
                                            v
                          +-----------------+-----------------+
                          |        Application Layer          |
                          |      (Use Cases & DTO Mappers)    |
                          +-----------------+-----------------+
                                            |
                                            v
                          +-----------------+-----------------+
                          |          Domain Layer             |
                          |  (Interfaces / Core Contracts)    |
                          +-----------------+-----------------+
                                            |
                                            v
                          +-----------------+-----------------+
                          |       Infrastructure Layer        |
                          | (GORM Repositories, HTTP Adapters)|
                          +-----------------------------------+
```

### ⚠️ Regla Personalizada Destacada: `common/models`
A diferencia del patrón hexagonal estricto donde las entidades de dominio y los modelos de persistencia se separan completamente agregando mapeadores redundantes:
- **Todos los modelos de base de datos de la aplicación se centralizan en `src/common/models`**.
- Estos modelos de datos (estructurados para **GORM** y serialización **JSON**) son compartidos transversalmente por todas las capas, reduciendo el código repetitivo (*boilerplate*) y garantizando una única fuente de verdad.

---

## 📁 Estructura del Repositorio

```
gestrym-training/
├── cmd/                          # Puntos de entrada para ejecutables y scripts CLI
│   ├── etl-foods/                # Executable para el pipeline ETL de alimentos (USDA)
│   └── import_exercises/         # Executable para la importación de ejercicios (ExerciseDB)
├── deployment/                   # Configuraciones de despliegue y archivos YAML locales
│   └── env_local.yaml            # Configuración de variables para entorno local
├── docs/                         # Documentación de la API autogenerada por Swagger
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── internal/                     # Código interno no exportable
│   └── etl/                      # Pipeline ETL (Extractor, Transformer, Loader, Pipeline)
├── src/                          # Código fuente principal de la aplicación
│   ├── app.go                    # Inicializador de la aplicación, DI y servidor Gin
│   ├── common/                   # Componentes compartidos
│   │   ├── config/               # Carga de configuraciones y YAML
│   │   ├── middleware/           # Middlewares de seguridad (JWT, RBAC)
│   │   ├── models/               # Modelos de datos GORM compartidos (Única fuente de verdad)
│   │   ├── routes/               # Registro central de rutas HTTP
│   │   ├── shared/               # Adaptadores HTTP externos (FileStorageAdapter)
│   │   └── utils/                # Utilidades de respuesta JSON y errores
│   ├── nutrition/                # Módulo Nutricional (Hexagonal)
│   │   ├── application/          # Use Cases (GenerateNutritionPlan, ImportFoods, SearchFoods)
│   │   ├── domain/               # Interfaces de repositorios nutricionales
│   │   ├── infrastructure/       # Repositorios GORM e integraciones
│   │   └── interfaces/           # Handlers HTTP para nutrición
│   └── training/                 # Módulo de Entrenamiento (Hexagonal)
│       ├── application/          # Use Cases (CreatePlan, AssignPlan, AdaptPlan, etc.)
│       ├── domain/               # Interfaces de repositorios de entrenamiento
│       ├── infrastructure/       # Repositorios GORM e integraciones
│       └── interfaces/           # Handlers HTTP para entrenamiento y planes
├── AI_MEMORY.md                  # Reglas de contexto y guía de arquitectura para IA
├── FRONTEND_TRAINING_GUIDE.md    # Guía técnica para consumo desde el frontend (React)
├── go.mod                        # Módulo Go y declaración de dependencias
├── go.sum                        # Checksums de dependencias Go
└── main.go                       # Punto de entrada principal del servidor HTTP
```

---

## 🛠️ Tecnologías y Dependencias

- **Lenguaje:** [Go 1.22+](https://golang.org/)
- **Framework Web:** [Gin Gonic (`github.com/gin-gonic/gin`)](https://github.com/gin-gonic/gin)
- **ORM:** [GORM (`gorm.io/gorm`)](https://gorm.io/) con driver PostgreSQL (`gorm.io/driver/postgres`)
- **Autenticación:** JWT (`github.com/golang-jwt/jwt/v5`)
- **Documentación API:** Swagger / Swaggo (`github.com/swaggo/swag`, `github.com/swaggo/gin-swagger`)
- **Configuración:** Viper (`github.com/spf13/viper`)
- **APIs Externas Integradas:**
  - **USDA FoodData Central API:** Fuente del catálogo nutricional.
  - **Pexels API:** Imágenes de alta calidad para alimentos.
  - **ExerciseDB (RapidAPI):** Importación del catálogo de ejercicios.

---

## ⚙️ Variables de Entorno

A continuación se detallan las variables requeridas en el entorno o en `deployment/env_local.yaml`:

| Variable | Descripción | Ejemplo / Valor por Defecto |
|---|---|---|
| `PORT` | Puerto de escucha del servicio | `8082` |
| `DB_HOST` | Host de PostgreSQL | `localhost` |
| `DB_PORT` | Puerto de PostgreSQL | `5432` |
| `DB_USER` | Usuario de la BD | `postgres` |
| `DB_PASSWORD` | Contraseña de la BD | `secret` |
| `DB_NAME` | Nombre de la base de datos | `gestrym_training` |
| `JWT_SECRET` | Clave secreta para firmar/verificar JWT | `your-secret-key` |
| `STORAGE_SERVICE_URL` | URL base del `file-service` interno | `http://localhost:8081` |
| `STORAGE_SERVICE_API_KEY` | X-API-Key para autenticación interna con storage | `internal-key-xxxx` |
| `USDA_API_KEY` | Clave de API de USDA FoodData Central | `DEMO_KEY` |
| `PEXELS_API_KEY` | Clave de API de Pexels para imágenes | `your-pexels-key` |
| `RAPID_API_KEY` | Clave de RapidAPI para ExerciseDB | `your-rapidapi-key` |

---

## 🚀 Instalación y Ejecución

### 1. Clonar el repositorio
```bash
git clone https://github.com/tu-usuario/gestrym.trining.back.git
cd gestrym.trining.back
```

### 2. Instalar dependencias
```bash
go mod download
```

### 3. Ejecutar en modo desarrollo local
El parámetro `--local=true` carga la configuración desde `deployment/env_local.yaml`:
```bash
go run main.go --local=true
```

### 4. Construir ejecutable de producción
```bash
go build -o gestrym-training main.go
./gestrym-training
```

---

## 🥗 Pipeline ETL Concurrentes (Alimentos)

El proyecto incluye un pipeline ETL robusto ubicado en `/cmd/etl-foods/main.go` para poblar la base de datos con información nutricional oficial de la USDA y vincular imágenes de Pexels automáticamente.

### Ejecución del ETL:
```bash
go run cmd/etl-foods/main.go
```

**Características del Pipeline:**
- **Extractor:** Paginación automática sobre la API de USDA.
- **Transformer:** Normalización de nombres de alimentos, extracción de macros por 100g y descarga de imágenes vía Pexels API.
- **Loader:** Transmisión de imágenes al `file-service` (MinIO) y *Upsert* directo en PostgreSQL evitando duplicados por nombre.
- **Concurrencia:** Utiliza Worker Pools con `goroutines` y `channels` para paralelizar el fetch de imágenes y la transformación.

---

## 🌐 API Endpoints y Documentación Swagger

### 📖 Swagger UI Interactivo
Con el servidor en ejecución local, ingresa en tu navegador a:
👉 **`http://localhost:8082/gestrym-training/swagger/index.html`**

Para regenerar la documentación de Swagger tras cambios en los handlers:
```bash
swag init
```

### 📌 Resumen de Rutas Principales

#### 🟢 Públicas (`/gestrym-training/public/`)
| Método | Endpoint | Descripción |
|---|---|---|
| `GET` | `/public/exercises` | Listar catálogo de ejercicios (Filtros: `bodyPart`, `target`) |
| `GET` | `/public/exercises/:id` | Obtener detalle de un ejercicio por ID |
| `POST` | `/public/exercises/import` | Disparar importación manual de ejercicios |
| `GET` | `/public/workouts/:id/full` | Obtener estructura completa de rutina (Workouts + Sets) |
| `GET` | `/public/foods` | Buscar alimentos en catálogo (Paginación: `search`, `page`, `limit`) |
| `GET` | `/public/foods/:id` | Detalle nutricional de un alimento |
| `POST` | `/public/foods/import` | Disparar importación manual desde USDA y Pexels |

#### 🔒 Privadas (`/gestrym-training/private/`) - *Requieren Bearer JWT Header*
| Método | Endpoint | Roles Permitidos | Descripción |
|---|---|---|---|
| `POST` | `/private/training-plans` | ENTRENADOR, USUARIO | Crear nuevo plan de entrenamiento |
| `GET` | `/private/training-plans/:id` | ENTRENADOR, USUARIO | Obtener plan por ID (RBAC: solo propios o asignados) |
| `GET` | `/private/training-plans/user/:userId` | ENTRENADOR, USUARIO | Obtener todos los planes de un usuario |
| `POST` | `/private/training-plans/:id/assign` | **SOLO ENTRENADOR** | Asignar plan a cliente (autoclonado si ya existe) |
| `POST` | `/private/training-plans/:id/days` | ENTRENADOR, USUARIO | Agregar día de entrenamiento a un plan |
| `POST` | `/private/training-plans/:id/clone` | ENTRENADOR, USUARIO | Clonar una plantilla a un usuario |
| `POST` | `/private/training-plans/adapt` | USUARIO | Adaptar plan actual por IA según cumplimiento |
| `PATCH` | `/private/training-plans/:id/days/:dayId/complete` | ENTRENADOR, USUARIO | Marcar día como completado o pendiente |
| `POST` | `/private/nutrition-plans/generate` | USUARIO | Generar plan nutricional personalizado (TDEE/Macros) |

#### 🤖 Internas (`/gestrym-training/internal/`) - *Requieren X-API-Key Header*
| Método | Endpoint | Descripción |
|---|---|---|
| `POST` | `/internal/training-plans/ai` | Recibe y almacena planes generados por el servicio externo de IA |

---

## 🔐 Seguridad y Autenticación (RBAC)

1. **JSON Web Tokens (JWT):** Todas las rutas bajo `/private` se validan contra el middleware `SetupJWTMiddleware()`, extrayendo la identidad del usuario y sus roles.
2. **Role-Based Access Control (RBAC):**
   - **Cliente (Role ID 4):** Solo puede consultar sus propios planes asignados o su plan nutricional personal.
   - **Entrenador / Admin (Role ID 2 y 1):** Tienen privilegios para crear plantillas, asignar planes a cualquier cliente y gestionar la estructura global.
3. **Internal Security:** Rutas bajo `/internal` protegidas con la cabecera `X-API-Key` para comunicación microservicio a microservicio.

---

## 🤖 Integración con IA y Adaptación Inteligente

Gestrym Training cuenta con soporte nativo para **Sobrecarga Progresiva Adaptativa por IA**:
- Endpoint `POST /private/training-plans/adapt`: Analiza la tasa de finalización de un plan activo. Si la persona completó más del **80%** de los días de entrenamiento, el sistema automáticamente:
  1. Realiza un clon profundo del plan y sus días asociados.
  2. Ajusta el título a *"Adaptación de [Nombre]"*.
  3. Prepara el plan para el siguiente nivel de sobrecarga progresiva.
- Endpoint `POST /internal/training-plans/ai`: Permite que un agente o servicio externo de IA (AI-Brain) inyecte directamente planes de entrenamiento estructurados listos para ser consumidos por el cliente.

---

## 📄 Licencia

Este proyecto es propiedad privada de **Gestrym** - Todos los derechos reservados.
