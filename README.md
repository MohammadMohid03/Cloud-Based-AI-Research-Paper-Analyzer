<p align="center">
  <h1 align="center">🔬 AI-Powered Research Paper Analyzer</h1>
  <p align="center">
    <strong>Upload · Analyze · Understand — Harness AI to unlock insights from research papers.</strong>
  </p>
  <p align="center">
    <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go" />
    <img src="https://img.shields.io/badge/React-18-61DAFB?style=for-the-badge&logo=react&logoColor=black" alt="React" />
    <img src="https://img.shields.io/badge/TypeScript-5-3178C6?style=for-the-badge&logo=typescript&logoColor=white" alt="TypeScript" />
    <img src="https://img.shields.io/badge/PostgreSQL-15-4169E1?style=for-the-badge&logo=postgresql&logoColor=white" alt="PostgreSQL" />
    <img src="https://img.shields.io/badge/AWS-Bedrock-FF9900?style=for-the-badge&logo=amazonaws&logoColor=white" alt="AWS" />
    <img src="https://img.shields.io/badge/Docker-Compose-2496ED?style=for-the-badge&logo=docker&logoColor=white" alt="Docker" />
    <img src="https://img.shields.io/badge/TailwindCSS-3-06B6D4?style=for-the-badge&logo=tailwindcss&logoColor=white" alt="Tailwind" />
  </p>
</p>

---

## 📖 Table of Contents

- [Overview](#-overview)
- [Architecture](#-architecture)
- [Tech Stack](#-tech-stack)
- [Features](#-features)
- [Prerequisites](#-prerequisites)
- [Quick Start — Docker Compose](#-quick-start--docker-compose)
- [Manual Setup](#-manual-setup)
- [API Documentation](#-api-documentation)
- [Database Schema](#-database-schema)
- [Environment Variables](#-environment-variables)
- [Project Structure](#-project-structure)
- [Screenshots](#-screenshots)
- [Contributing](#-contributing)
- [License](#-license)

---

## 🌐 Overview

The **AI-Powered Research Paper Analyzer** is a cloud-ready, full-stack web application that allows researchers, students, and academics to:

1. **Upload** research papers in PDF format.
2. **Analyze** them automatically using AI (Amazon Bedrock / Claude) to extract summaries, key findings, methodology, limitations, and future work.
3. **Chat** interactively with the paper — ask questions about specific sections, request deeper explanations, or compare findings.
4. **Organize** papers with bookmarks, search, and an intuitive dashboard.

Built as a semester project, the system demonstrates modern full-stack development, cloud-native architecture, RESTful API design, JWT authentication, and integration with AWS AI services.

---

## 🏗 Architecture

```
┌────────────────────┐         ┌────────────────────┐         ┌────────────────────┐
│                    │  HTTP   │                    │  SQL    │                    │
│   React Frontend   │◄───────►│   Go / Gin API     │◄───────►│   PostgreSQL 15    │
│   (TypeScript)     │  :3000  │   (REST + JWT)     │  :5432  │   (Data Store)     │
│   Tailwind CSS     │         │                    │         │                    │
└────────────────────┘         └────────┬───────────┘         └────────────────────┘
                                        │
                                        │ AWS SDK
                                        ▼
                               ┌────────────────────┐
                               │  AWS Cloud Services │
                               │  ┌──────────────┐  │
                               │  │ S3 (Storage)  │  │
                               │  └──────────────┘  │
                               │  ┌──────────────┐  │
                               │  │ Bedrock (AI)  │  │
                               │  └──────────────┘  │
                               └────────────────────┘
```

**Data Flow:**

1. User uploads a PDF via the React frontend.
2. The Go backend validates the file, stores it (local disk or S3), and creates a database record.
3. The backend sends the paper content to Amazon Bedrock for AI-powered analysis.
4. Analysis results (summary, key findings, etc.) are stored in PostgreSQL.
5. Users can view results on the dashboard or open an interactive chat session.

> For a detailed architecture deep-dive, see [ARCHITECTURE.md](ARCHITECTURE.md).  
> For AWS deployment instructions, see [AWS_DEPLOYMENT.md](AWS_DEPLOYMENT.md).

---

## 🛠 Tech Stack

| Layer        | Technology                           | Purpose                           |
| ------------ | ------------------------------------ | --------------------------------- |
| **Frontend** | React 18 + TypeScript                | Single-page application           |
| **Styling**  | Tailwind CSS 3                       | Utility-first CSS framework       |
| **Backend**  | Go 1.21+ with Gin framework         | High-performance REST API         |
| **ORM**      | GORM                                | Go ORM for PostgreSQL             |
| **Database** | PostgreSQL 15                        | Relational data storage           |
| **Auth**     | JWT (JSON Web Tokens)                | Stateless authentication          |
| **AI**       | Amazon Bedrock (Claude 3 Sonnet)     | Paper analysis & chat             |
| **Storage**  | Local disk / Amazon S3               | PDF file storage                  |
| **DevOps**   | Docker & Docker Compose              | Containerization & orchestration  |
| **Cloud**    | AWS (EC2, S3, RDS, CloudFront, etc.) | Production deployment             |

---

## ✨ Features

### Core Features
- 📄 **PDF Upload & Management** — Drag-and-drop upload with file validation and progress tracking.
- 🤖 **AI-Powered Analysis** — Automatic extraction of summary, key findings, methodology, limitations, future work, and keywords.
- 💬 **Interactive Chat** — Ask follow-up questions about any paper in a chat interface powered by AI.
- 🔍 **Search & Filter** — Full-text search across paper titles and metadata.
- 🔖 **Bookmarks** — Save and organize favorite papers for quick access.

### Technical Features
- 🔐 **JWT Authentication** — Secure sign-up, login, and token-based API access.
- 📊 **User Dashboard** — Overview of uploaded papers, analysis status, and recent activity.
- 🌙 **Dark Mode** — Toggle between light and dark themes.
- 📱 **Responsive Design** — Works on desktop, tablet, and mobile.
- 🐳 **Docker-Ready** — One-command setup with Docker Compose.
- ☁️ **Cloud-Native** — Designed for AWS deployment with S3, RDS, and Bedrock.
- 🧪 **Mock AI Provider** — Development without AWS credentials using built-in mock responses.

---

## 📋 Prerequisites

| Tool              | Version  | Required For         |
| ----------------- | -------- | -------------------- |
| **Docker**        | 20.10+   | Docker Compose setup |
| **Docker Compose**| 2.0+     | Docker Compose setup |
| **Go**            | 1.21+    | Manual backend setup |
| **Node.js**       | 20+      | Manual frontend setup|
| **npm**           | 9+       | Manual frontend setup|
| **PostgreSQL**    | 15+      | Manual database setup|

> **Tip:** For the quickest setup, you only need **Docker** and **Docker Compose**.

---

## 🚀 Quick Start — Docker Compose

Get the entire application running in under 2 minutes:

```bash
# 1. Clone the repository
git clone https://github.com/your-username/Cloud-based-AI-Powered-Research-Paper-Analyzer.git
cd Cloud-based-AI-Powered-Research-Paper-Analyzer

# 2. Create environment file
cp .env.example .env

# 3. Start all services
docker-compose up --build
```

Once all containers are healthy:

| Service    | URL                          |
| ---------- | ---------------------------- |
| Frontend   | http://localhost:3000        |
| Backend API| http://localhost:8080        |
| PostgreSQL | localhost:5432               |

**Default demo credentials:**

| Email              | Password   | Role  |
| ------------------ | ---------- | ----- |
| admin@example.com  | admin123   | Admin |
| demo@example.com   | demo1234   | User  |

To stop all services:
```bash
docker-compose down        # Stop containers
docker-compose down -v     # Stop and remove volumes (resets database)
```

---

## 🔧 Manual Setup

<details>
<summary><strong>Backend (Go + Gin)</strong></summary>

```bash
# 1. Navigate to the backend directory
cd backend

# 2. Install Go dependencies
go mod download

# 3. Set up environment variables
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=research_paper_analyzer
export JWT_SECRET=your-secret-key
export AI_PROVIDER=mock
export STORAGE_PROVIDER=local
export SERVER_PORT=8080
export ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173

# 4. Ensure PostgreSQL is running and the database exists
createdb research_paper_analyzer

# 5. (Optional) Seed the database
psql -U postgres -d research_paper_analyzer -f ../database/schema.sql

# 6. Run the backend server
go run cmd/server/main.go

# 7. (Optional) Hot-reload with Air
# Install: go install github.com/air-verse/air@latest
air
```

The API will be available at `http://localhost:8080`.

</details>

<details>
<summary><strong>Frontend (React + TypeScript)</strong></summary>

```bash
# 1. Navigate to the frontend directory
cd frontend

# 2. Install dependencies
npm install

# 3. Start the development server
npm run dev
```

The app will be available at `http://localhost:5173` (Vite default).

For a production build:
```bash
npm run build    # Output in dist/
npm run preview  # Preview production build locally
```

</details>

<details>
<summary><strong>Database (PostgreSQL)</strong></summary>

```bash
# Option A: Use Docker for PostgreSQL only
docker run -d \
  --name rpa-postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=research_paper_analyzer \
  -p 5432:5432 \
  postgres:15-alpine

# Option B: Use an existing PostgreSQL installation
createdb -U postgres research_paper_analyzer
psql -U postgres -d research_paper_analyzer -f database/schema.sql
```

</details>

---

## 📡 API Documentation

### Authentication

| Method | Endpoint             | Description               | Auth |
| ------ | -------------------- | ------------------------- | ---- |
| POST   | `/api/auth/register` | Register a new user       | ❌   |
| POST   | `/api/auth/login`    | Login & receive JWT token | ❌   |
| GET    | `/api/auth/me`       | Get current user profile  | ✅   |

### Papers

| Method | Endpoint              | Description                  | Auth |
| ------ | --------------------- | ---------------------------- | ---- |
| GET    | `/api/papers`         | List all papers (paginated)  | ✅   |
| POST   | `/api/papers`         | Upload a new paper (PDF)     | ✅   |
| GET    | `/api/papers/:id`     | Get paper details            | ✅   |
| DELETE | `/api/papers/:id`     | Delete a paper               | ✅   |
| GET    | `/api/papers/search`  | Search papers by query       | ✅   |

### Analysis

| Method | Endpoint                       | Description                | Auth |
| ------ | ------------------------------ | -------------------------- | ---- |
| POST   | `/api/papers/:id/analyze`      | Trigger AI analysis        | ✅   |
| GET    | `/api/papers/:id/analysis`     | Get analysis results       | ✅   |

### Chat

| Method | Endpoint                         | Description                | Auth |
| ------ | -------------------------------- | -------------------------- | ---- |
| POST   | `/api/papers/:id/chat`           | Create a new chat session  | ✅   |
| GET    | `/api/papers/:id/chat`           | List chat sessions         | ✅   |
| POST   | `/api/chat/:sessionId/messages`  | Send a message in chat     | ✅   |
| GET    | `/api/chat/:sessionId/messages`  | Get chat message history   | ✅   |

### Bookmarks

| Method | Endpoint                       | Description               | Auth |
| ------ | ------------------------------ | ------------------------- | ---- |
| GET    | `/api/bookmarks`               | List bookmarked papers    | ✅   |
| POST   | `/api/bookmarks`               | Bookmark a paper          | ✅   |
| DELETE | `/api/bookmarks/:id`           | Remove a bookmark         | ✅   |

### User

| Method | Endpoint                | Description               | Auth |
| ------ | ----------------------- | ------------------------- | ---- |
| PUT    | `/api/users/profile`    | Update user profile       | ✅   |
| PUT    | `/api/users/password`   | Change password           | ✅   |

### Health & Meta

| Method | Endpoint          | Description        | Auth |
| ------ | ----------------- | ------------------ | ---- |
| GET    | `/api/health`     | Health check       | ❌   |

**Authentication:** Include the JWT token in the `Authorization` header:
```
Authorization: Bearer <your-jwt-token>
```

---

## 🗄 Database Schema

The application uses 6 main tables:

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│    users     │     │    papers    │     │   analyses   │
├──────────────┤     ├──────────────┤     ├──────────────┤
│ id (PK)      │◄────│ user_id (FK) │     │ paper_id(FK) │
│ email        │     │ title        │────►│ summary      │
│ password_hash│     │ authors      │     │ key_findings │
│ first_name   │     │ abstract     │     │ methodology  │
│ last_name    │     │ file_path    │     │ limitations  │
│ role         │     │ status       │     │ keywords     │
│ created_at   │     │ created_at   │     │ confidence   │
└──────────────┘     └──────┬───────┘     └──────────────┘
        │                   │
        │                   │
        ▼                   ▼
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│  bookmarks   │     │chat_sessions │     │chat_messages │
├──────────────┤     ├──────────────┤     ├──────────────┤
│ user_id (FK) │     │ paper_id(FK) │     │session_id(FK)│
│ paper_id(FK) │     │ user_id (FK) │     │ role         │
│ note         │     │ title        │     │ content      │
│ created_at   │     │ is_active    │     │ tokens_used  │
└──────────────┘     └──────────────┘     └──────────────┘
```

> See [`database/schema.sql`](database/schema.sql) for the full SQL schema with indexes, constraints, and seed data.

---

## ⚙️ Environment Variables

| Variable               | Default                                        | Description                                  |
| ---------------------- | ---------------------------------------------- | -------------------------------------------- |
| `DB_HOST`              | `localhost`                                    | PostgreSQL host                              |
| `DB_PORT`              | `5432`                                         | PostgreSQL port                              |
| `DB_USER`              | `postgres`                                     | PostgreSQL username                          |
| `DB_PASSWORD`          | `postgres`                                     | PostgreSQL password                          |
| `DB_NAME`              | `research_paper_analyzer`                      | Database name                                |
| `SERVER_PORT`          | `8080`                                         | Backend server port                          |
| `ALLOWED_ORIGINS`      | `http://localhost:3000,http://localhost:5173`   | CORS allowed origins                         |
| `JWT_SECRET`           | `your-secret-key-change-in-production`         | Secret for signing JWT tokens                |
| `AI_PROVIDER`          | `mock`                                         | AI provider: `mock` or `bedrock`             |
| `STORAGE_PROVIDER`     | `local`                                        | Storage: `local` or `s3`                     |
| `AWS_REGION`           | `us-east-1`                                    | AWS region                                   |
| `AWS_ACCESS_KEY_ID`    | —                                              | AWS access key (for S3/Bedrock)              |
| `AWS_SECRET_ACCESS_KEY`| —                                              | AWS secret key (for S3/Bedrock)              |
| `S3_BUCKET`            | —                                              | S3 bucket name for file storage              |
| `BEDROCK_MODEL_ID`     | `anthropic.claude-3-sonnet-20240229-v1:0`      | Amazon Bedrock model identifier              |

---

## 📁 Project Structure

```
Cloud-based-AI-Powered-Research-Paper-Analyzer/
│
├── backend/                    # Go backend (REST API)
│   ├── cmd/
│   │   └── server/
│   │       └── main.go         # Application entry point
│   ├── internal/
│   │   ├── config/             # Configuration loading
│   │   ├── handlers/           # HTTP request handlers
│   │   ├── middleware/         # Auth, CORS, logging middleware
│   │   ├── models/            # GORM database models
│   │   ├── repository/        # Database access layer
│   │   ├── services/          # Business logic
│   │   │   ├── ai/            # AI provider (Bedrock + mock)
│   │   │   └── storage/       # Storage provider (S3 + local)
│   │   └── router/            # Gin route definitions
│   ├── uploads/                # Local file upload directory
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
│
├── frontend/                   # React frontend (SPA)
│   ├── public/
│   ├── src/
│   │   ├── api/               # Axios API client
│   │   ├── components/        # Reusable UI components
│   │   ├── contexts/          # React contexts (Auth, Theme)
│   │   ├── hooks/             # Custom React hooks
│   │   ├── pages/             # Page-level components
│   │   ├── types/             # TypeScript type definitions
│   │   ├── utils/             # Utility functions
│   │   ├── App.tsx
│   │   └── main.tsx
│   ├── Dockerfile
│   ├── package.json
│   ├── tailwind.config.js
│   ├── tsconfig.json
│   └── vite.config.ts
│
├── database/
│   └── schema.sql              # Raw SQL schema + seed data
│
├── docker-compose.yml          # Multi-service orchestration
├── .env.example                # Environment variable template
├── .gitignore                  # Git ignore rules
├── Makefile                    # Build automation
├── README.md                   # ← You are here
├── ARCHITECTURE.md             # System architecture docs
├── AWS_DEPLOYMENT.md           # AWS deployment guide
└── LICENSE                     # MIT License
```

---

## 📸 Screenshots

> Screenshots will be added as the UI is finalized.

| Dashboard | Paper Analysis | Chat Interface |
| --------- | -------------- | -------------- |
| *Coming soon* | *Coming soon* | *Coming soon* |

---

## 🤝 Contributing

Contributions are welcome! Follow these steps:

1. **Fork** the repository.
2. **Create** a feature branch:
   ```bash
   git checkout -b feature/amazing-feature
   ```
3. **Commit** your changes with clear messages:
   ```bash
   git commit -m "feat: add amazing feature"
   ```
4. **Push** to your branch:
   ```bash
   git push origin feature/amazing-feature
   ```
5. **Open a Pull Request** with a detailed description.

### Commit Convention

We follow [Conventional Commits](https://www.conventionalcommits.org/):

| Prefix     | Usage                        |
| ---------- | ---------------------------- |
| `feat:`    | New feature                  |
| `fix:`     | Bug fix                      |
| `docs:`    | Documentation only           |
| `style:`   | Formatting, no logic change  |
| `refactor:`| Code restructuring           |
| `test:`    | Adding or updating tests     |
| `chore:`   | Build process, dependencies  |

### Code Style

- **Go:** Follow [Effective Go](https://go.dev/doc/effective_go) guidelines. Run `golangci-lint` before committing.
- **TypeScript/React:** Follow the ESLint + Prettier configuration in the frontend project.

---

## 📄 License

This project is licensed under the **MIT License**.

```
MIT License

Copyright (c) 2025 AI-Powered Research Paper Analyzer

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

---

<p align="center">
  Made with ❤️ for Cloud Computing — Semester Project 2025
</p>
