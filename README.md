# 🌐 Sentrinet

**Sentrinet** is a full-stack network intelligence platform that brings together real-time scanning, monitoring, and analytics for modern infrastructure.  
It uses a modular **Golang backend**, a **Next.js frontend**, and **Prometheus** metrics collection—deployed with **Render** and **Vercel**, with continuous delivery through **GitHub Actions**.

---

## 🚀 Overview

Sentrinet helps detect network activity, schedule automated scans, and visualize metrics in real time.  
It is designed for security engineers, DevOps teams, and researchers who need fast and reliable insights into distributed systems.

---

## 🧰 Tech Stack

| Layer | Technology |
|-------|-------------|
| **Frontend** | [Next.js 14](https://nextjs.org) • TypeScript • TailwindCSS |
| **Backend** | [Golang](https://golang.org) • Clean architecture in `/internal` modules |
| **Database** | SQLite (default) |
| **Metrics** | Prometheus |
| **Deployments** | Render (backend) • Vercel (frontend) |
| **CI/CD** | GitHub Actions (auto-deploy on push to `main`) |

---

## ⚙️ Local Setup

### 1. Clone the repository
```bash
git clone https://github.com/<your-username>/sentrinet.git
cd sentrinet
```

### 2. Run the Backend
```bash
cd backend
go mod tidy
go run ./server/main.go
```
Default port: `http://localhost:8080`

### 3. Run the Frontend
```bash
cd frontend/sentriface
npm install
npm run dev
```
Frontend: `http://localhost:3000`

---

## 🔄 Continuous Integration & Deployment
The repository uses GitHub Actions to automatically deploy both the backend and frontend whenever new commits are pushed to the `main` branch.

**Required Secrets**

Add these in **GitHub → Settings → Secrets and variables → Actions**:

| Secret | Description |
|--------|-------------|
| RENDER_API_KEY | Add these in GitHub → Settings → Secrets and variables → Actions: |
| VERCEL_TOKEN |Your Render API key (from Render dashboard → API Keys) |

## 🔐 Authentication Flow

1. User registers through `/register`
2. Backend returns a JWT token
3. Token stored in browser by `AuthContext`
4. `authFetch` automatically attaches token for protected requests
5. 401 responses trigger auto logout

## 🧱 Key Features
* 🔐 Secure authentication and session handling
* 🧩 Modular Go backend for easy scaling
* 📡 Real-time network scan scheduling
* 📊 Prometheus metrics endpoint (`/metrics`)
* 💅 Clean Next.js + TailwindCSS frontend
* ⚙️ Full CI/CD via GitHub Actions
* 🚀 One-click deployments on Render and Vercel

## 📊 Monitoring
* Prometheus configured through `backend/prometheus.yml`
* Exposes live metrics at `/metrics`
* Compatible with Grafana dashboards