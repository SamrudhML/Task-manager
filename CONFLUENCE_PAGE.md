h1. Task Manager Project - Technical Overview

h2. 1) Project Summary

This project is a full-stack Task Manager application built as an intern assignment and deployed on Kubernetes.

It contains:
- A Go REST API for authentication and task CRUD operations.
- A MongoDB database for persistent storage.
- A static HTML/JavaScript UI served by Nginx.

The API uses JWT authentication. Each user only sees and modifies their own tasks.

h2. 2) Tech Stack

|| Layer || Technology ||
| Backend API | Go, chi router, MongoDB Go driver |
| Authentication | JWT (HS256), bcrypt password hashing |
| Database | MongoDB |
| Frontend | Vanilla HTML/CSS/JavaScript |
| Containerization | Docker (multi-stage build) |
| Orchestration | Kubernetes (Deployments, Services, ConfigMaps, PVC) |

h2. 3) Repository Structure

{code}
/
|- main.go
|- go.mod
|- Dockerfile
|- Dockerfile.nobase
|- task-api-linux
|- internals/
|  |- auth/
|  |  |- middleware.go
|  |- handlers/
|  |  |- auth_handler.go
|  |  |- task_handler.go
|  |- models/
|  |  |- task.go
|  |  |- user.go
|  |- repositories/
|  |  |- task_repository.go
|  |  |- user_repository.go
|- k8s/
|  |- api-configmap.yaml
|  |- api-deployment.yaml
|  |- api-service.yaml
|  |- mongo-deployment.yaml
|  |- mongo-pvc.yaml
|  |- mongo-service.yaml
|  |- ui-configmap.yaml
|  |- ui-deployment.yaml
|  |- ui-service.yaml
|- ui/
|  |- index.html
{code}

h2. 4) Architecture

h3. Runtime Architecture

{code}
[ Browser UI ] --HTTP--> [ task-ui-service :80 / NodePort 30001 ] --> [ Nginx serving index.html ]
                                         |
                                         | API calls (fetch)
                                         v
                          [ task-api-service :8080 / NodePort 30000 ] --> [ Go API ]
                                                                           |
                                                                           v
                                                    [ mongo-service :27017 (ClusterIP) ] --> [ MongoDB + PVC ]
{code}

h3. Request Flow
- User registers or logs in through /auth endpoints.
- API validates credentials and returns a JWT token.
- UI stores token in localStorage and sends Authorization: Bearer <token> for task operations.
- API middleware validates JWT and injects user context.
- Task repository filters all operations by user_id.

h2. 5) Backend API Details

h3. Health Endpoint
- GET /health
- Returns: {"status":"ok"}

h3. Auth Endpoints
- POST /auth/register
- POST /auth/login

Request body:
{code:json}
{
  "username": "string",
  "password": "string"
}
{code}

Response body:
{code:json}
{
  "token": "jwt-token",
  "username": "string"
}
{code}

Validation rules:
- Register: username length >= 3, password length >= 6.
- Passwords are hashed using bcrypt before storage.

h3. Task Endpoints (Protected)
All task endpoints require:
- Header: Authorization: Bearer <jwt>

Endpoints:
- POST /tasks
- GET /tasks
- GET /tasks/{id}
- PUT /tasks/{id}
- DELETE /tasks/{id}

Task model:
{code:json}
{
  "id": "string",
  "title": "string",
  "description": "string",
  "status": "todo | in-progress | done",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
{code}

Notes:
- status defaults to "todo" on create.
- user_id is stored in DB but hidden from API JSON responses.
- Deletion returns HTTP 204 No Content.

h2. 6) Data Layer

MongoDB collections:
- users
- tasks

Task ownership enforcement:
- Every task operation includes user_id in its MongoDB filter.
- This prevents cross-user access even if task IDs are known.

h2. 7) Configuration and Environment Variables

Configured in main.go and Kubernetes ConfigMap:
- MONGO_HOST (default: localhost)
- MONGO_PORT (default: 27017)
- MONGO_DB (default: taskdb)
- JWT_SECRET (default: replace-me-in-production)
- HTTP_ADDR (default: :8080)

Important:
- JWT_SECRET should be replaced with a secure value in non-local environments.
- For production, prefer a Kubernetes Secret over ConfigMap for JWT secret.

h2. 8) Docker

h3. Primary Image Build
The main Dockerfile uses a multi-stage build:
- Build stage: golang:1.26.4-alpine
- Runtime stage: alpine:latest

Build command:
{code:bash}
docker build -t task-api:local .
{code}

Run command:
{code:bash}
docker run --rm -p 8080:8080 \
  -e MONGO_HOST=host.docker.internal \
  -e MONGO_PORT=27017 \
  -e MONGO_DB=taskdb \
  -e JWT_SECRET=dev-secret \
  task-api:local
{code}

h3. Alternative Image
- Dockerfile.nobase runs a prebuilt binary (task-api-linux) from scratch.
- Useful for very small runtime images when a static binary already exists.

h2. 9) Kubernetes Deployment

h3. Deployed Components
- MongoDB Deployment + ClusterIP Service + PersistentVolumeClaim (10Gi)
- API Deployment + NodePort Service (30000) + ConfigMap
- UI Deployment + NodePort Service (30001) + ConfigMap-hosted index.html

h3. Apply manifests
{code:bash}
kubectl apply -f k8s/mongo-pvc.yaml
kubectl apply -f k8s/mongo-deployment.yaml
kubectl apply -f k8s/mongo-service.yaml
kubectl apply -f k8s/api-configmap.yaml
kubectl apply -f k8s/api-deployment.yaml
kubectl apply -f k8s/api-service.yaml
kubectl apply -f k8s/ui-configmap.yaml
kubectl apply -f k8s/ui-deployment.yaml
kubectl apply -f k8s/ui-service.yaml
{code}

h3. Verify status
{code:bash}
kubectl get pods
kubectl get svc
kubectl get configmap
kubectl get pvc
{code}

Expected service ports:
- task-api-service -> NodePort 30000
- task-ui-service -> NodePort 30001
- mongo-service -> ClusterIP 27017

h3. Access application
- UI: http://<node-ip>:30001
- API health: http://<node-ip>:30000/health

For minikube:
{code:bash}
minikube service task-ui-service --url
minikube service task-api-service --url
{code}

h2. 10) UI Behavior

The UI supports:
- Register and login.
- Create tasks.
- Filter tasks by status (all, todo, in-progress, done).
- Update task status.
- Delete tasks.
- Logout.

Implementation notes:
- API base URL is derived from browser host with port 8080 when served directly.
- In Kubernetes, UI is also available through the ConfigMap-mounted version in task-ui deployment.

h2. 11) Security and Operational Notes

Current good practices:
- Password hashing with bcrypt.
- JWT-protected task routes.
- Per-user data isolation in repository queries.

Recommended improvements:
- Move JWT_SECRET from ConfigMap to Secret.
- Add input validation for task status values on update.
- Restrict CORS origin in production.
- Pin MongoDB image version (avoid mongo:latest).

h2. 12) Known Trade-offs

- API currently uses in-code CORS middleware with permissive '*'.
- No refresh token flow; login token validity is fixed at 24 hours.
- UI is a single static page without framework/state library.
- No automated tests are included yet.

h2. 13) Quick Demo Script

1. Open UI endpoint.
2. Register a new user.
3. Create 2-3 tasks.
4. Move one task to in-progress and one to done.
5. Refresh and verify task persistence.
6. Logout and log in again to confirm session-based access.

h2. 14) File Pointers (for reviewers)

- Application bootstrap and routing: main.go
- JWT generation and middleware: internals/auth/middleware.go
- Auth handlers: internals/handlers/auth_handler.go
- Task handlers: internals/handlers/task_handler.go
- Task persistence logic: internals/repositories/task_repository.go
- User persistence logic: internals/repositories/user_repository.go
- Kubernetes manifests: k8s/
- UI source: ui/index.html
