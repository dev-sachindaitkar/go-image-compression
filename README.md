# Project Purpose

The purpose of this project is to provide a high-performance, concurrent API Gateway designed specifically to handle high-throughput image uploads, streaming data processing, and on-the-fly image compression. By leveraging Go's native concurrency primitives, the service efficiently offloads intensive I/O and CPU-bound optimization tasks, ensuring minimal memory overhead before securely forwarding compressed assets and metadata downstream to our orchestration and AI analysis layers.

## Repository Structure

The project is organized into a mono-repo structure containing both the high-performance Go backend service and the frontend orchestration layer.

```text
├── backend/                # Go (Golang) API Gateway & Processing Service
│   ├── cmd/                # Entry points for the application
│   ├── internal/           # Private application and business logic
│   └── pkg/                # Explicit library code usable by other services
├── frontend/               # Next.js Orchestration & UI Layer
│   ├── src/                # Frontend application source code
│   └── public/             # Static assets
├── README.md               # Project documentation
└── .gitignore              # Git ignore rules
```

## Getting Started

Follow these steps to set up and run both the backend and frontend services locally.

### Prerequisites

Ensure you have the following tools installed on your machine:
* **Go** (v1.22 or higher)
* **Node.js** (v18.0 or higher) & **npm** (or yarn/pnpm)
* **Git**

### Installation & Local Setup

1. **Clone the repository:**
   ```bash
   git clone https://github.com/dev-sachindaitkar/go-image-compression.git
   ```
3. **Set up the Backend (Go):**
    ```bash 
    cd backend
    go mod download
    go run cmd/gateway/main.go
    ```
    The backend gateway will be accessible at http://localhost:8080.
4. **Set up the Frontend (Next.js):**
   ```bash
   cd frontend
   npm install
   npm run dev
   ```
   The frontend application will be accessible at http://localhost:3000.
## Configuration & Environment Variables

Both the backend and frontend services require configuration via `.env` files to communicate correctly. Create a `.env` file in the root of their respective directories.

### Backend Configuration (`backend/.env`)

```env
PORT=8080
ENV=development

# Concurrency & Performance Settings
MAX_WORKERS=10
MAX_UPLOAD_SIZE_MB=20
RATE_LIMIT_FPS=50

# Downstream Redirection
NEXTJS_ORCHESTRATION_URL=http://localhost:3000
```

### Frontend Configuration (`frontend/.env`)
```env
NEXT_PUBLIC_API_GATEWAY_URL=http://localhost:8080
NODE_ENV=development
```

### Frontend 
<img width="1486" height="847" alt="image" src="https://github.com/user-attachments/assets/2ba6372c-9bab-44a1-afdb-5f364f7c3085" />
