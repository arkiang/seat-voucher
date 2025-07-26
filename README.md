
# Bookcabin Voucher App

A fullstack app for generating crew seat vouchers based on flight and aircraft information. A take home project by Book Cabin.
Built with:

- **Frontend**: Vite + React + Material UI (MUI)
- **Backend**: Go (Gin + Gorm), SQLite
- **Dev Stack**: Docker Compose (multi-service)

---

## Run with Docker (Recommended)

### 1. Prerequisites

Docker is required. Please install Docker Desktop before proceeding: https://www.docker.com/products/docker-desktop

Build and Run

```bash
docker compose up --build
```

This will:

- Build and run the Go backend on http://localhost:8081  
- Build and serve the React frontend on http://localhost:3000  

---

### 2. Access the App

Open your browser:

http://localhost:3000

Try generating vouchers by entering:
- Name
- ID
- Flight Number (e.g. `JT692`)
- Date (`DD-MM-YYYY`)
- Aircraft Type (e.g. `Airbus 320`)

---

## Author

Apriyanto Arkiang — Backend engineer with 8+ years of experience.

---
