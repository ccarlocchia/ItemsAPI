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

## 🎛️ Visual playground

The static playground in `docs/` can be published with GitHub Pages. It supports:

- **Sample data** — try the catalog without running the Go API.
- **Live API** — enter the base URL of a deployed API and send item IDs to `POST /v1/api/items`. The URL is saved in your browser.

To publish it, open the repository's **Settings → Pages**, choose **Deploy from a branch**, select the `main` branch and the `/docs` folder, then save. GitHub Pages will provide the playground URL when deployment finishes.

The live API must be hosted separately; GitHub Pages only serves static files. The API enables cross-origin requests for the playground. When testing locally, start the API with `go run main.go` and use `http://localhost:8081` as the API base URL. A browser may block an HTTP API when the playground is served over HTTPS, so use an HTTPS API URL for the published Pages site.
