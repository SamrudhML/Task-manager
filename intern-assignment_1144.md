# Intern Assignment — Task Management App on Kubernetes

**Team:** Stargate  
**Stack:** Go · Docker · Kubernetes · MongoDB  
**Duration:** ~1–2 weeks  
**Goal:** Get hands-on with the core tools we use every day by building and deploying a small but complete application end to end.

---

## Overview

You will build a **Task Management** application with three components:

| Component | Tech | What you build |
|-----------|------|----------------|
| **API** | Go | CRUD REST API for tasks |
| **Database** | MongoDB | Runs as a K8s Deployment |
| **UI** | HTML/JS or any simple frontend | Basic task board UI |

All three components must run as services inside a local Kubernetes cluster.

---

## Part 1 — Go REST API

Build a REST API in Go with the following endpoints:

```
POST   /tasks          – Create a new task
GET    /tasks          – List all tasks
GET    /tasks/{id}     – Get a single task by ID
PUT    /tasks/{id}     – Update a task (title, description, or status)
DELETE /tasks/{id}     – Delete a task
```

### Task Schema

```json
{
  "id":          "string (auto-generated)",
  "title":       "string (required)",
  "description": "string (optional)",
  "status":      "todo | in-progress | done",
  "created_at":  "timestamp",
  "updated_at":  "timestamp"
}
```

### Requirements

- Use the standard `net/http` package or a lightweight router like `chi` or `gorilla/mux`
- Connect to MongoDB using the official Go driver (`go.mongodb.org/mongo-driver`)
- Read MongoDB connection details (host, port, db name) from **environment variables** — no hardcoded values
- Return proper HTTP status codes (`200`, `201`, `400`, `404`, `500`)
- Respond with JSON for all endpoints

---

## Part 2 — Dockerise the API

- Write a `Dockerfile` for the Go API
- Use a **multi-stage build**: build stage with `golang` image, final stage with `alpine` or `scratch`
- The resulting image should be as small as possible
- Test locally with `docker build` and `docker run` before moving to K8s

---

## Part 3 — Kubernetes Deployment

Deploy all three components to a local K8s cluster (use [minikube](https://minikube.sigs.k8s.io/) or [kind](https://kind.sigs.k8s.io/)).

Write YAML manifests for each component:

### MongoDB

- `Deployment` — single replica MongoDB pod
- `Service` — ClusterIP to expose MongoDB internally
- `PersistentVolumeClaim` — so data survives pod restarts

### Go API

- `Deployment` — your API image, at least 1 replica
- `Service` — ClusterIP (or NodePort so you can hit it from outside the cluster)
- `ConfigMap` or `Secret` — to pass MongoDB connection details as env vars to the API pod

### UI

- `Deployment` — serve your UI (nginx with static files works fine)
- `Service` — NodePort or LoadBalancer so you can access it in a browser

### Suggested File Layout

```
k8s/
  mongo-deployment.yaml
  mongo-service.yaml
  mongo-pvc.yaml
  api-deployment.yaml
  api-service.yaml
  api-configmap.yaml
  ui-deployment.yaml
  ui-service.yaml
```

---

## Part 4 — Simple UI

Build a minimal UI that can:

- Display the list of tasks
- Create a new task (title + description)
- Mark a task as `in-progress` or `done`
- Delete a task

There is no strict requirement on the framework. Plain HTML + JS with `fetch()` calls is perfectly fine. The UI talks directly to the Go API.

---

## Deliverables

When you're done, your repo should contain:

```
/
├── main.go              (or cmd/, internal/ — your call)
├── go.mod
├── Dockerfile
├── k8s/
│   └── *.yaml
├── ui/
│   └── index.html (+ any JS/CSS)
└── README.md
```

### README must include

- How to run locally with Docker
- How to deploy to K8s (step-by-step commands)
- How to access the UI and API once deployed
- Any design decisions or trade-offs you made

---

## Evaluation Criteria

| Area | What we look at |
|------|----------------|
| **Go code** | Clean structure, error handling, proper use of interfaces |
| **Docker** | Multi-stage build, small image, no secrets baked in |
| **Kubernetes** | Correct resource types, env vars via ConfigMap/Secret, app actually runs |
| **API design** | Correct HTTP verbs, status codes, JSON responses |
| **README** | Clear enough that someone else can run it from scratch |

---

## Tips

- Start with the Go API working locally first, then Dockerise, then move to K8s — don't try to do it all at once
- Use `kubectl logs`, `kubectl describe pod`, and `kubectl exec` to debug issues in the cluster
- The MongoDB connection string will look like: `mongodb://mongo-service:27017` — the hostname is the K8s Service name
- If you get stuck, check the official docs: [Go](https://go.dev/doc/), [Docker](https://docs.docker.com/), [Kubernetes](https://kubernetes.io/docs/home/), [MongoDB Go Driver](https://www.mongodb.com/docs/drivers/go/current/)

---

*Questions? Ping us on Slack or drop by during standup. Good luck!*
