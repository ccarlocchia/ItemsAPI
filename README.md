# 🧩 Items API

This project implements a simple **Items API** in Go that serves product information from a local JSON file (`mocked_items.json`).  
It follows a clean, layered architecture with three main components: **Handler**, **Service**, and **Repository**.

---

## 📁 Project Structure

```
project/
├── mocked_items.json # Mocked data source (sample items)
├── internal/
│ ├── handlers/
│ │ └── handlers.go
│ ├── services/
│ │ └── items_service.go
│ └── repositories/
│ └── items_repository.go
└── go.mod / go.sum
```

---

## 🧠 Architecture Overview

```
┌────────────┐
│ Handler │ → Handles HTTP requests (Gin)
└────┬───────┘
│
┌────▼───────┐
│ Service │ → Business logic, delegates to repository
└────┬───────┘
│
┌────▼───────┐
│ Repository │ → Loads and retrieves items from JSON
└────────────┘

```

---

## 🚀 How It Works

### 1. Repository Layer
**File:** `internal/repositories/items_repository.go`
- Loads all items from `mocked_items.json` when initialized.
- Provides the method:
  ```
  GetItemsByIDs(ctx context.Context, ids []string) []model.ItemDTO
  which returns a list of items matching the provided IDs.

### 2. Service Layer

File: `internal/services/items_service.go`

Acts as a middle layer between the handler and repository.

Contains minimal business logic and delegates data fetching to the repository.

### 3. Handler Layer

File: `internal/handlers/handlers.go`

Defines the HTTP endpoint:

`POST /v1/api/items`


Expects a JSON body like:

```
{
    "data": [
        { "id": "123abc" },
        { "id": "456def" }
    ]
}
```


Returns:

- `200 OK → Matching items found`
- `400 Bad Request → Invalid or empty request body`
- `404 Not Found → No items found for given IDs`

---

## ⚙️ Setup & Run
Prerequisites

Go 1.21 or newer

`mocked_items.json` present in the project root

Install dependencies

```
go mod tidy
```

### Run the API
```
go run main.go
```

### Example request

```
curl -X POST http://localhost:8080/v1/api/items \
-H "Content-Type: application/json" \
-d '{"data": [{"id": "123abc"}]}'
```
