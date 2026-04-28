# 🧠 AI Context & Guidelines (AI_MEMORY)

## 1. 📌 Project Overview
**Gestrym** is a comprehensive fitness platform designed to manage training (workouts, exercises, sets, progress) and nutrition (planned). It tracks user progress, records statistics, and features authentication, notifications, and file storage. 

The backend system is built as a set of **Microservices in Golang**.

## 2. 🏗️ Architecture Explanation
The project follows **Hexagonal Architecture (Ports and Adapters)**, adapted slightly for our specific microservice needs to reduce boilerplate:
- **Domain (Core):** Pure business logic, independent of external frameworks.
- **Application (Use Cases):** Connects external inputs (via ports) with domain logic. Handlers and use cases live here.
- **Infrastructure:** Frameworks, routing (Gin), database connections (GORM/PostgreSQL), external API integrations.

⚠️ **CRITICAL CUSTOM RULE - `common/models`** ⚠️
Unlike strict hexagonal architecture where domain entities are completely separated from DB models, **all database models (structs used for persistence, migrations, and cross-layer data representation) are located in `common/models`**.
- This package is shared across all layers and (when applicable) across services to avoid boilerplate mapping where it isn't beneficial.
- These structs are structurally used for ORM (GORM) and migrations. 
- AI assistants should use them directly across the application instead of creating duplicate domain-only structs unless the complexity of the domain heavily dictates otherwise.

## 3. 🧩 Microservices Overview
- **auth-service** *(implemented)*: Handles user authentication, JWT token generation, role mapping (clients, trainers, admins), and registration.
- **notification-service** *(implemented)*: Sends emails (e.g., via Brevo) and system alerts. Handles templates like recovery emails.
- **file-service** *(implemented)*: Manages multimedia uploads (images, workout videos) securely.
- **training-service** *(in progress)*: The core of the platform. Handles fitness: Exercises (catalog), workouts (structure & execution), and progress tracking.
- **nutrition-service** *(integrated)*: Initial implementation integrated within the project. Manages Food catalogs and Nutritional tracking.

## 4. 💻 Coding Standards
- **Golang Best Practices:** Follow standard Go idioms (Effective Go). Use descriptive variable names, handle errors explicitly without swallowing them, and rely on early returns (guard clauses).
- **Clean Architecture:** Keep HTTP Handlers dumb. They should parse requests, call a Use Case, and format the response. Do not place business logic in Handlers.
- **Separation of Concerns:** Repositories handle database interactions only. Use Cases orchestrate business rules. HTTP/Controllers handle web mechanics.
- **Dependency Injection:** Pass dependencies (like repositories) into Use Cases, and Use Cases into HTTP Handlers via constructors.

## 5. 🎯 Domain Guidelines
- **Entity Definitions:** Define entities cleanly based on real-world fitness logic. An `Exercise` is a dictionary/catalog item. A `Workout` is a planned or executed session. A `WorkoutSet` belongs to a `Workout` and logs reps/weights.
- **Responsibilities:** 
  - *Infrastructure/Repositories* construct and execute GORM queries.
  - *Application/Use Cases* validate business rules natively (e.g., "Cannot log a set for a future date").
  - *Domain/Models* (in `common/models`) define the data schema strictly with GORM and JSON serialization tags.

## 6. 🤖 AI Instructions
When generating code for Gestrym:
1. **Always respect the existing architecture:** Handlers -> Use Cases -> Repositories.
2. **Reuse Models:** Always import and use structs from `common/models` for DB interactions and domain representations. **Do not create new entity files under the domain layer** if a shared DB model will suffice.
3. **Avoid Boilerplate:** Don't write extensive mappers between DB models and Domain models. Lean on `common/models`.
4. **Context Maintenance:** Keep in mind the context of the specific microservice being worked on. Do not import or mix logic from `auth-service` into `training-service` arbitrarily; assume communication via APIs or shared libs.
5. **No Hallucinations:** Use `gin-gonic/gin` for routing and `gorm.io/gorm` for ORM. Check `src/common` for standard utilities before inventing new ones.
6. **No External Runtime Dependencies for Catalogs:** Avoid relying on external APIs (like ExerciseDB) at runtime. Catalog data must be imported into our own database (via scripts or admin endpoints) and served internally. Always build an adapter to fetch, an application usecase to map/deduplicate, and a repository to store using `common/models` via GORM.

## 7. 🚀 Future Expansion Notes
- **Nutrition Module:** Will be added later to track diets, macros, and meal plans. Models should remain loosely coupled so nutrition can optionally tie into training (e.g., caloric goals vs. workout expenditure).
- **Advanced Analytics:** Data scaling will require efficient queries and possibly event sourcing for stats. Always anticipate potential N+1 query traps.
- **Multi-tenant Support:** Plan for a potential shift where gym centers/trainers autonomously manage their own clients. Always consider a `ClientID` or `TrainerID` foreign key presence in core entities.

## 8. 🌐 API Endpoints & Frontend Consumption
The backend exposes RESTful APIs using standard JSON payloads. All routes are documented via `swaggo/swag`.

**Available Endpoints (`training-service`):**
- **`GET /gestrym-training/public/exercises`**: Retrieves all exercises. Accepts `?bodyPart=` and `?target=` query filters.
- **`GET /gestrym-training/public/exercises/:id`**: Retrieves details of a specific exercise by its unique GORM ID.
- **`POST /gestrym-training/public/exercises/import`**: Triggers the manual import/sync process for exercises (ExerciseDB).
- **`GET /gestrym-training/public/workouts/:id/full`**: Retrieves a complete workout structure, including exercises and sets, optimized for frontend rendering (React). Contains nested `WorkoutExercise` and `WorkoutSet` data.
- **`GET /gestrym-training/public/foods`**: Searches the food catalog. Supports `?search=`, `?page=`, and `?limit=` (default 1/10). Returns results with categories and total count.
- **`GET /gestrym-training/public/foods/:id`**: Retrieves specific nutritional details for a food item.
- **`POST /gestrym-training/public/foods/import`**: Triggers the manual import process from **USDA FoodData Central** and fetches images from **Pexels**.

## 9. 📦 File Storage Integration
Training entities (like `Exercises` and `Foods`) are linked to multimedia files through a `CollectionID`. 
- **Storage Workflow**: When importing or creating entities with files, the `training-service` communicates with the `file-service` internally.
- **`FileStorageAdapter`**: Used to upload files (from URLs or Readers) to the storage service.
- **Collection-based group**: Multiple files (images, videos, gifs) per entity are grouped under the same `CollectionID`.
- **Environment Variables**:
  - `STORAGE_SERVICE_URL`: Endpoint of the file-service.
  - `STORAGE_SERVICE_API_KEY`: X-API-Key for internal authentication.
  - `USDA_API_KEY`: Key for USDA FoodData Central API.
  - `PEXELS_API_KEY`: Key for Pexels Image API.
  - `RAPID_API_KEY`: Key for ExerciseDB (RapidAPI).

## 10. 🥗 Nutrition & Workout Modeling
- **Food Catalog**: Foods are imported from USDA and stored locally. Mapped nutrients: Calories, Protein, Carbs, Fats.
- **Image Management**: 
    - Food images are fetched from **Pexels** during import using normalized food names.
    - Images are uploaded to MinIO via `file-service`.
    - `Food` model stores `ImageURL` (MinIO link) and `CollectionID`.
- **Optimization**:
    - **N+1 Avoidance**: Use GORM `Preload` for hierarchical data (Workouts -> Exercises -> Sets).
    - **Pagination**: Compulsory for food search and exercise listings.
- **Frontend DTOs**: Use specialized DTOs (e.g., `WorkoutFullResponse`) to assemble nested structures.

## 11. 🚀 Batch ETL Pipeline (Foods)
To populate the food catalog at scale, a standalone ETL (Extract, Transform, Load) pipeline is implemented in `/internal/etl`.
- **Extractor**: Fetches raw data from USDA FoodData Central. Supports pagination and query-based extraction.
- **Transformer**: Normalizes food names, extracts nutritional macros (Calories, Protein, Carbs, Fats), and fetches high-quality images via Pexels.
- **Loader**: Uploads images to MinIO via streaming and performs bulk/upsert operations in PostgreSQL via GORM.
- **Concurrency & Reliability**:
    - **Worker Pool**: Uses goroutines and channels (`jobs` and `results`) to parallelize the Transform and Image Fetching stages.
    - **Retries**: Implements a 3-attempt retry mechanism with exponential backoff for external API calls.
    - **Deduplication**: Ensures no duplicate food names are inserted into the database.
- **Execution**: Run via CLI: `go run cmd/etl-foods/main.go`.

---
*Last updated: 2026-04-22 (Implemented AI Plan Adaptation & Nutritional Plan Generation)*

**Swagger Documentation:**
- Swagger definitions live within the `docs/` folder.
- They are auto-generated from `// @Summary`, `// @Description`, etc., annotations above specific Gin Handler functions.
- Always run `swag init` when altering or building new endpoints to maintain a fresh contract for the frontend.
- Frontend developers can browse the interactive Swagger UI dynamically at **`GET /gestrym-training/swagger/index.html`** when the server is running locally.

---

## 12. 🗓️ Training Plans Module

*Last updated: 2026-04-19 (Implemented Training Plans: plans, days, assignment, RBAC)*

### Models (in `common/models`)

| Model | Description |
|---|---|
| `TrainingPlan` | Weekly/monthly/custom fitness plan. Has `AssignedTo` (nullable), `CreatedBy`, `IsTemplate`, `DurationDays`. |
| `TrainingDay` | One day within a plan. Links to `WorkoutID`. Has `DayNumber` (1..N) and `Notes`. |
| `TrainingPlanAssignment` | (Future) Explicit assignment record with `AssignedBy` (trainer), `UserID`, `StartDate`. |
| `NutritionPlan` | Stores daily caloric and macro goals based on user objective (loss, gain, maintenance). |

### Architecture Location

```
src/common/models/
  TrainingPlan.go
  TrainingDay.go
  TrainingPlanAssignment.go

src/training/domain/interfaces/
  TrainingPlanRepository.go
  TrainingDayRepository.go

src/training/infrastructure/repositories/
  TrainingPlanRepositoryImpl.go      ← GORM, Preload chains to avoid N+1
  TrainingDayRepositoryImpl.go

src/training/application/usecases/
  CreateTrainingPlanUseCase.go
  AssignTrainingPlanUseCase.go       ← auto-clones if plan already assigned
  GetTrainingPlanUseCase.go          ← RBAC built-in (user can only see own plans)
  GetUserTrainingPlansUseCase.go
  UpdateDayCompletionUseCase.go      ← updates `IsCompleted` flag
  AdaptTrainingPlanUseCase.go        ← **New**: Evaluates completion rate and auto-clones for intensity level-up.
  TrainingPlanMapper.go              ← shared DTO mappers

src/training/application/dtos/
  TrainingPlanDTO.go                 ← request + response structs

src/nutrition/application/usecases/
  GenerateNutritionPlanUseCase.go    ← **New**: Calculates TDEE/Macros using Mifflin-St Jeor formula.

src/training/interfaces/http/handlers/
  TrainingPlanHandler.go             ← Added `/adapt` endpoint

src/nutrition/interfaces/http/handlers/
  NutritionPlanHandler.go            ← **New**: Handles `/nutrition-plans/generate`
```

### Endpoints (`/gestrym-training/private/training-plans`)

| Method | Path | Role | Description |
|---|---|---|---|
| POST | `/` | TRAINER, USER | Create a training plan |
| GET | `/:id` | TRAINER, USER | Get plan by ID (users only see own) |
| GET | `/user/:userId` | TRAINER, USER | Get all plans for a user |
| POST | `/adapt` | USER | Adapt latest plan based on completion progress |
| POST | `/:id/assign` | **TRAINER ONLY** | Assign plan to user (auto-clones if needed) |
| POST | `/:id/days` | TRAINER, USER | Add a workout day to a plan |
| POST | `/:id/clone` | TRAINER, USER | Clone a template plan to a user |
| PATCH | `/:id/days/:dayId/complete` | TRAINER, USER | Mark a training day as completed/not |

### RBAC Rules
- All routes under `/private` require JWT (via `SetupJWTMiddleware()`).
- `/assign` additionally requires `RequireRoles(RoleCoach, RoleAdmin)`.
- For GET endpoints, users (RoleCliente = 4) can only access plans where `AssignedTo = their own userID`.
- For `/complete`, only the assigned user or a trainer/admin can update progress.

### Key Business Rules
- **DayNumber validation**: `AddTrainingDayUseCase` ensures `1 ≤ dayNumber ≤ plan.DurationDays`.
- **Auto-clone on re-assign**: If a trainer tries to assign a plan already assigned to another user, the system clones the plan + all its days before assigning.
- **Deep Cloning**: `CloneTrainingPlanUseCase` performs a deep copy of the plan and all its associated days.
- **Progress Tracking**: Each day has an `IsCompleted` flag. Clones start with all days as uncompleted.
- **AI Adaptation**: `/adapt` analyzes the latest plan. If completion > 80%, it clones the plan with an "Adapted" tag for progressive overload.
- **Template support**: `IsTemplate = true` makes a plan reusable. Templates have `AssignedTo = null`.

### Nutrition Plans (`/gestrym-training/private/nutrition-plans`)

| Method | Path | Role | Description |
|---|---|---|---|
| POST | `/generate` | USER | Generate macro/caloric goals based on weight, height, age, and objective. |

**Objective logic**:
- `weight_loss`: TDEE - 500 kcal
- `muscle_gain`: TDEE + 300 kcal
- `maintenance`: TDEE
- **Macros**: 2g Protein/kg, 0.8g Fat/kg, balance in Carbs.

### Future Considerations
- `TrainingPlanAssignment` model is ready to store `StartDate` for assignment history.
- `AssignTrainingPlanUseCase` accepts `startDate` but doesn't persist it yet (marked with `_ startDate`).
- Future AI plan generation should create a `TrainingPlan` with `IsTemplate = false` and set `CreatedBy = AI_AGENT_ID` or similar.
- Plan cloning (`CloneTrainingPlanUseCase`) can be built on top of the existing clone logic in `AssignTrainingPlanUseCase`.

## 13. 🤖 AI Agent Persona & Prompt Context
You are a senior Golang backend developer specialized in microservices and hexagonal architecture.

I already have a training-service with:

* exercises
* workouts
* foods (read-only catalog)

I need to EXTEND it to support **Training Plans (weekly/monthly)** and assignments.

---

### ⚠️ IMPORTANT RULES

* Do NOT modify existing project structure
* Models are in common/models → DO NOT move or duplicate
* Follow hexagonal architecture strictly
* Use GORM
* Use dependency injection

---

### 🎯 GOAL

Implement Training Plans:

* Multi-day plans (7–30 days)
* Assign plans to users
* Trainer can assign plans
* User can view plans

---

### 🧱 MODELS (use common/models)

TrainingPlan:

* ID
* Name
* Description
* DurationDays
* CreatedBy
* IsTemplate
* CreatedAt

TrainingPlanAssignment:

* ID
* TrainingPlanID
* UserID
* AssignedBy
* StartDate

TrainingDay:

* ID
* TrainingPlanID
* DayNumber
* WorkoutID
* Notes

---

### ⚙️ USE CASES

* CreateTrainingPlanUseCase
* AssignTrainingPlanUseCase
* GetTrainingPlanUseCase
* GetUserTrainingPlansUseCase
* AddTrainingDayUseCase

---

### 🗄️ REPOSITORIES

TrainingPlanRepository:

* Create
* FindByID
* FindByUserID

TrainingDayRepository:

* Create
* FindByPlanID

AssignmentRepository:

* Assign
* FindByUserID

---

### 🌐 ENDPOINTS

POST /gestrym-training/private/training-plans
GET  /gestrym-training/private/training-plans/:id
GET  /gestrym-training/private/training-plans/user/:userId
POST /gestrym-training/private/training-plans/:id/assign
POST /gestrym-training/private/training-plans/:id/days

---

### 🔐 AUTH

* TRAINER → can assign
* USER → can only see own plans

---

### 📦 RESPONSE

Return nested structure (frontend-ready):

{
id,
name,
durationDays,
days: [
{
dayNumber,
workout: {...}
}
]
}

---

### 🚀 BONUS

* Add endpoint to clone template plan
* Prepare for AI-generated plans

---

### OUTPUT

* Repositories
* Use cases
* Handlers
* DTOs
* Clean structure

