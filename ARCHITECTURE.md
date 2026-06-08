# 🏛 System Architecture — AI-Powered Research Paper Analyzer

> This document describes the system architecture, component interactions, data flows, technology decisions, and scalability considerations.

---

## Table of Contents

- [High-Level Architecture](#high-level-architecture)
- [Component Diagram](#component-diagram)
- [Data Flow — Paper Upload & Analysis](#data-flow--paper-upload--analysis)
- [Authentication Flow](#authentication-flow)
- [Chat Flow](#chat-flow)
- [Technology Choices & Rationale](#technology-choices--rationale)
- [Scalability Considerations](#scalability-considerations)
- [Security Measures](#security-measures)

---

## High-Level Architecture

The system follows a classic **three-tier architecture** with a clear separation of concerns:

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         PRESENTATION TIER                              │
│                                                                         │
│   React 18 + TypeScript + Tailwind CSS                                 │
│   ┌──────────┐ ┌───────────┐ ┌──────────┐ ┌────────────┐              │
│   │Dashboard │ │Paper View │ │Chat UI   │ │Auth Pages  │              │
│   └──────────┘ └───────────┘ └──────────┘ └────────────┘              │
│                        │  Axios HTTP Client  │                         │
└────────────────────────┼─────────────────────┼─────────────────────────┘
                         │   REST API (JSON)   │
                         ▼                     ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                          APPLICATION TIER                               │
│                                                                         │
│   Go 1.21+ / Gin Framework                                             │
│   ┌──────────────┐ ┌──────────────┐ ┌──────────────┐                  │
│   │  Middleware   │ │   Handlers   │ │   Services   │                  │
│   │ ┌──────────┐ │ │ ┌──────────┐ │ │ ┌──────────┐ │                  │
│   │ │  Auth    │ │ │ │  Auth    │ │ │ │  Paper   │ │                  │
│   │ │  CORS    │ │ │ │  Paper   │ │ │ │  AI      │ │                  │
│   │ │  Logger  │ │ │ │  Chat    │ │ │ │  Chat    │ │                  │
│   │ │  Rate    │ │ │ │  Bookmark│ │ │ │  Storage │ │                  │
│   │ └──────────┘ │ │ └──────────┘ │ │ └──────────┘ │                  │
│   └──────────────┘ └──────────────┘ └──────┬───────┘                  │
│                                             │                          │
│   ┌──────────────────────────┐   ┌──────────┴───────────┐             │
│   │     Repository Layer     │   │  External Providers   │             │
│   │     (GORM + PostgreSQL)  │   │  ┌────────┐┌───────┐ │             │
│   │                          │   │  │Bedrock ││  S3   │ │             │
│   └──────────────────────────┘   │  │(AI)    ││(Files)│ │             │
│                                   │  └────────┘└───────┘ │             │
│                                   └──────────────────────┘             │
└─────────────────────────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                            DATA TIER                                    │
│                                                                         │
│   PostgreSQL 15                                                        │
│   ┌─────────┐ ┌─────────┐ ┌──────────┐ ┌──────────┐ ┌──────────────┐ │
│   │ users   │ │ papers  │ │analyses  │ │chat_sess.│ │chat_messages │ │
│   └─────────┘ └─────────┘ └──────────┘ └──────────┘ └──────────────┘ │
│   ┌─────────┐ ┌──────────────┐                                        │
│   │bookmarks│ │activity_logs │                                        │
│   └─────────┘ └──────────────┘                                        │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Component Diagram

### Frontend Components

```
src/
├── api/
│   └── client.ts              ← Axios instance with JWT interceptor
├── contexts/
│   ├── AuthContext.tsx         ← Global auth state + login/logout
│   └── ThemeContext.tsx        ← Dark/light theme toggle
├── pages/
│   ├── LoginPage.tsx           ← User authentication
│   ├── RegisterPage.tsx        ← User registration
│   ├── DashboardPage.tsx       ← Paper list + stats overview
│   ├── UploadPage.tsx          ← Drag-and-drop PDF upload
│   ├── PaperDetailPage.tsx     ← Analysis results view
│   ├── ChatPage.tsx            ← Interactive AI chat
│   └── BookmarksPage.tsx       ← Saved papers
└── components/
    ├── Navbar.tsx              ← Navigation bar
    ├── PaperCard.tsx           ← Paper preview card
    ├── AnalysisPanel.tsx       ← Analysis results display
    ├── ChatMessage.tsx         ← Individual chat bubble
    └── ProtectedRoute.tsx      ← Route guard for auth
```

### Backend Components

```
internal/
├── config/
│   └── config.go              ← Environment config loader
├── middleware/
│   ├── auth.go                ← JWT validation middleware
│   ├── cors.go                ← CORS configuration
│   └── logger.go              ← Request logging
├── models/
│   ├── user.go                ← User model + validation
│   ├── paper.go               ← Paper model
│   ├── analysis.go            ← Analysis model
│   ├── chat.go                ← ChatSession + ChatMessage models
│   └── bookmark.go            ← Bookmark model
├── repository/
│   ├── user_repo.go           ← User CRUD operations
│   ├── paper_repo.go          ← Paper CRUD + search
│   ├── analysis_repo.go       ← Analysis CRUD
│   ├── chat_repo.go           ← Chat session + message CRUD
│   └── bookmark_repo.go       ← Bookmark CRUD
├── services/
│   ├── auth_service.go        ← Registration, login, JWT generation
│   ├── paper_service.go       ← Paper management logic
│   ├── analysis_service.go    ← Orchestrates AI analysis
│   ├── chat_service.go        ← Chat logic + AI interaction
│   ├── ai/
│   │   ├── provider.go        ← AI provider interface
│   │   ├── bedrock.go         ← Amazon Bedrock implementation
│   │   └── mock.go            ← Mock AI for development
│   └── storage/
│       ├── provider.go        ← Storage provider interface
│       ├── s3.go              ← S3 implementation
│       └── local.go           ← Local filesystem implementation
├── handlers/
│   ├── auth_handler.go        ← Auth endpoint handlers
│   ├── paper_handler.go       ← Paper endpoint handlers
│   ├── chat_handler.go        ← Chat endpoint handlers
│   └── bookmark_handler.go    ← Bookmark endpoint handlers
└── router/
    └── router.go              ← Route registration + middleware setup
```

---

## Data Flow — Paper Upload & Analysis

```
    User                    Frontend               Backend                  AWS / Storage        Database
     │                         │                      │                         │                   │
     │  1. Select PDF          │                      │                         │                   │
     │─────────────────────►   │                      │                         │                   │
     │                         │  2. POST /api/papers │                         │                   │
     │                         │  (multipart/form)    │                         │                   │
     │                         │─────────────────────►│                         │                   │
     │                         │                      │  3. Validate file       │                   │
     │                         │                      │     (type, size)        │                   │
     │                         │                      │                         │                   │
     │                         │                      │  4. Store file ─────────►                   │
     │                         │                      │     (S3 or local)       │                   │
     │                         │                      │                         │                   │
     │                         │                      │  5. Create paper record ────────────────────►
     │                         │                      │     (status: uploaded)  │                   │
     │                         │                      │                         │                   │
     │                         │  6. Return paper ID  │                         │                   │
     │                         │◄─────────────────────│                         │                   │
     │                         │                      │                         │                   │
     │  7. Click "Analyze"     │                      │                         │                   │
     │─────────────────────►   │                      │                         │                   │
     │                         │ 8. POST /papers/:id/ │                         │                   │
     │                         │    analyze           │                         │                   │
     │                         │─────────────────────►│                         │                   │
     │                         │                      │  9. Update status ──────────────────────────►
     │                         │                      │     (processing)        │                   │
     │                         │                      │                         │                   │
     │                         │                      │ 10. Extract text        │                   │
     │                         │                      │     from PDF            │                   │
     │                         │                      │                         │                   │
     │                         │                      │ 11. Send to Bedrock ───►│                   │
     │                         │                      │     (AI analysis)       │                   │
     │                         │                      │                         │                   │
     │                         │                      │ 12. Receive results ◄───│                   │
     │                         │                      │                         │                   │
     │                         │                      │ 13. Save analysis ──────────────────────────►
     │                         │                      │     (status: analyzed)  │                   │
     │                         │                      │                         │                   │
     │                         │ 14. Return analysis  │                         │                   │
     │                         │◄─────────────────────│                         │                   │
     │  15. Display results    │                      │                         │                   │
     │◄───────────────────────│                      │                         │                   │
```

---

## Authentication Flow

```
    User                    Frontend                Backend                 Database
     │                         │                       │                       │
     │  1. Enter credentials   │                       │                       │
     │─────────────────────►   │                       │                       │
     │                         │  2. POST /auth/login  │                       │
     │                         │  {email, password}    │                       │
     │                         │──────────────────────►│                       │
     │                         │                       │  3. Find user ────────►
     │                         │                       │                       │
     │                         │                       │  4. Verify bcrypt hash│
     │                         │                       │     password_hash     │
     │                         │                       │                       │
     │                         │                       │  5. Generate JWT      │
     │                         │                       │     (24h expiry)      │
     │                         │                       │                       │
     │                         │  6. {token, user}     │                       │
     │                         │◄──────────────────────│                       │
     │                         │                       │                       │
     │                         │  7. Store token in    │                       │
     │                         │     localStorage      │                       │
     │                         │                       │                       │
     │  8. Redirect to         │                       │                       │
     │     Dashboard           │                       │                       │
     │◄───────────────────────│                       │                       │
     │                         │                       │                       │
     │  ─ ─ ─ Subsequent API Calls ─ ─ ─              │                       │
     │                         │                       │                       │
     │                         │  9. GET /api/papers   │                       │
     │                         │  Authorization:       │                       │
     │                         │  Bearer <token>       │                       │
     │                         │──────────────────────►│                       │
     │                         │                       │ 10. Validate JWT      │
     │                         │                       │     Extract user_id   │
     │                         │                       │                       │
     │                         │                       │ 11. Fetch data ───────►
     │                         │                       │                       │
     │                         │  12. Return data      │                       │
     │                         │◄──────────────────────│                       │
```

**JWT Token Structure:**

```json
{
  "header": {
    "alg": "HS256",
    "typ": "JWT"
  },
  "payload": {
    "user_id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
    "email": "user@example.com",
    "role": "user",
    "exp": 1700000000,
    "iat": 1699913600
  }
}
```

---

## Chat Flow

```
    User                   Frontend                Backend              AI Provider        Database
     │                        │                       │                      │                 │
     │ 1. Open chat for       │                       │                      │                 │
     │    paper               │                       │                      │                 │
     │───────────────────►    │                       │                      │                 │
     │                        │ 2. POST /papers/:id/  │                      │                 │
     │                        │    chat               │                      │                 │
     │                        │──────────────────────►│                      │                 │
     │                        │                       │ 3. Create session ───────────────────►│
     │                        │                       │                      │                 │
     │                        │ 4. {session_id}       │                      │                 │
     │                        │◄──────────────────────│                      │                 │
     │                        │                       │                      │                 │
     │ 5. Type question       │                       │                      │                 │
     │───────────────────►    │                       │                      │                 │
     │                        │ 6. POST /chat/:sid/   │                      │                 │
     │                        │    messages           │                      │                 │
     │                        │ {content: "..."}      │                      │                 │
     │                        │──────────────────────►│                      │                 │
     │                        │                       │ 7. Save user msg ───────────────────►│
     │                        │                       │                      │                 │
     │                        │                       │ 8. Build context:    │                 │
     │                        │                       │    - Paper content   │                 │
     │                        │                       │    - Analysis        │                 │
     │                        │                       │    - Chat history    │                 │
     │                        │                       │                      │                 │
     │                        │                       │ 9. Send prompt ─────►│                 │
     │                        │                       │                      │                 │
     │                        │                       │ 10. AI response ◄────│                 │
     │                        │                       │                      │                 │
     │                        │                       │ 11. Save assistant ──────────────────►│
     │                        │                       │     message          │                 │
     │                        │                       │                      │                 │
     │                        │ 12. {role: assistant, │                      │                 │
     │                        │      content: "..."}  │                      │                 │
     │                        │◄──────────────────────│                      │                 │
     │ 13. Display response   │                       │                      │                 │
     │◄──────────────────────│                       │                      │                 │
```

---

## Technology Choices & Rationale

### Why Go for the Backend?

| Factor         | Rationale                                                                                  |
| -------------- | ------------------------------------------------------------------------------------------ |
| **Performance**| Go compiles to native binaries — excellent latency and low memory footprint.                |
| **Concurrency**| Goroutines make handling concurrent PDF uploads and AI requests efficient.                  |
| **Simplicity** | Minimal syntax, fast compile times, great standard library.                                |
| **Gin**        | Most popular Go web framework — fast routing, middleware support, JSON binding.             |
| **GORM**       | Feature-rich ORM with auto-migration, hooks, and PostgreSQL-native support.                |

### Why React + TypeScript?

| Factor           | Rationale                                                                                |
| ---------------- | ---------------------------------------------------------------------------------------- |
| **Type Safety**  | TypeScript catches bugs at compile time — fewer runtime errors.                          |
| **Ecosystem**    | Largest component ecosystem — excellent libraries for file uploads, chat UIs, routing.   |
| **Vite**         | Sub-second hot module replacement during development.                                    |
| **Tailwind CSS** | Rapid UI development with utility classes — no context-switching to CSS files.           |

### Why PostgreSQL?

| Factor           | Rationale                                                                                |
| ---------------- | ---------------------------------------------------------------------------------------- |
| **Reliability**  | Battle-tested ACID-compliant relational database.                                        |
| **JSON Support** | Native JSONB columns for flexible schema fields (activity log metadata).                 |
| **Full-Text**    | `pg_trgm` extension enables fuzzy search on paper titles without external search engine. |
| **AWS RDS**      | Managed PostgreSQL on AWS with automated backups and failover.                           |

### Why Amazon Bedrock?

| Factor           | Rationale                                                                                |
| ---------------- | ---------------------------------------------------------------------------------------- |
| **Managed**      | No model hosting or GPU management — fully serverless.                                   |
| **Model Choice** | Access to Claude 3, Llama, Titan, and other foundation models.                           |
| **Security**     | Data stays within AWS — no third-party API calls.                                        |
| **Pay-per-Use**  | No upfront costs — pay only for tokens consumed.                                         |

### Why Docker Compose?

| Factor            | Rationale                                                                               |
| ----------------- | --------------------------------------------------------------------------------------- |
| **Reproducibility**| Identical environments across dev machines — "it works on my machine" solved.           |
| **One-Command**    | `docker-compose up` starts everything: database, backend, frontend.                    |
| **Isolation**      | Each service runs in its own container with defined resource limits.                    |

---

## Scalability Considerations

### Current Architecture (Semester Project)

- **Single instance** of each service behind Docker Compose.
- **Local file storage** on the backend container's filesystem.
- **Synchronous** AI analysis (user waits for completion).

### Production-Ready Scaling Path

```
                    ┌──────────────┐
                    │ CloudFront   │  ← CDN for static frontend
                    │ (CDN)        │
                    └──────┬───────┘
                           │
                    ┌──────┴───────┐
                    │    ALB       │  ← Application Load Balancer
                    │              │
                    └──────┬───────┘
                           │
              ┌────────────┼────────────┐
              ▼            ▼            ▼
        ┌──────────┐ ┌──────────┐ ┌──────────┐
        │ Backend  │ │ Backend  │ │ Backend  │  ← Auto Scaling Group
        │ (EC2)    │ │ (EC2)    │ │ (EC2)    │
        └────┬─────┘ └────┬─────┘ └────┬─────┘
             │             │             │
             └─────────────┼─────────────┘
                           │
                    ┌──────┴───────┐
                    │  RDS         │  ← Multi-AZ PostgreSQL
                    │  (Primary)   │
                    └──────┬───────┘
                           │
                    ┌──────┴───────┐
                    │  RDS         │  ← Read Replica
                    │  (Replica)   │
                    └──────────────┘
```

| Concern               | Solution                                                          |
| ---------------------- | ----------------------------------------------------------------- |
| **Horizontal Scaling** | Run multiple backend instances behind an ALB.                     |
| **File Storage**       | Switch from local disk to S3 (already implemented).               |
| **Async Processing**   | Use SQS to queue analysis jobs; process with Lambda or workers.   |
| **Database Scaling**   | RDS Multi-AZ for HA; read replicas for query-heavy workloads.     |
| **Caching**            | Add ElastiCache (Redis) for session tokens and hot data.          |
| **CDN**                | Serve frontend from S3 + CloudFront for global low latency.       |
| **Monitoring**         | CloudWatch metrics, alarms, and dashboards.                       |

---

## Security Measures

### Authentication & Authorization

| Measure                  | Implementation                                              |
| ------------------------ | ----------------------------------------------------------- |
| Password Hashing         | bcrypt with cost factor 10                                  |
| Token-Based Auth         | JWT with HS256, 24-hour expiry                              |
| Role-Based Access        | `user` and `admin` roles enforced at middleware level        |
| Route Protection         | Auth middleware applied to all `/api/*` routes except auth  |

### Data Security

| Measure                  | Implementation                                              |
| ------------------------ | ----------------------------------------------------------- |
| SQL Injection Prevention | Parameterized queries via GORM                              |
| XSS Prevention           | React's default output escaping + CSP headers               |
| CORS                     | Strict origin whitelist via `ALLOWED_ORIGINS`               |
| File Validation          | MIME type + extension checks on upload                      |
| File Size Limits         | Configurable max upload size (default: 20MB)                |
| Input Validation         | Request binding + validation tags on all models             |

### Infrastructure Security

| Measure                  | Implementation                                              |
| ------------------------ | ----------------------------------------------------------- |
| Environment Secrets      | All secrets in env vars, never committed to git             |
| Docker Isolation         | Each service in its own container with minimal base images  |
| Network Segmentation     | Docker bridge network — services communicate internally     |
| HTTPS                    | TLS termination at CloudFront / ALB in production           |
| IAM Least Privilege      | Specific IAM roles for Bedrock + S3 access only             |

---

<p align="center">
  <em>For deployment instructions, see <a href="AWS_DEPLOYMENT.md">AWS_DEPLOYMENT.md</a></em>
</p>
