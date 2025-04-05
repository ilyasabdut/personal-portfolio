# Minimal Portfolio Website

A personal portfolio site built using **Go**, **HTML**, **HTMX**, **SQLite**, and **vanilla CSS** — inspired by [motherfuckingwebsite.com](http://motherfuckingwebsite.com/). No frameworks. Just clean, minimal, fast.

---

## 🌱 Project Development Stages

1. **Project Initialization**  
   - Set up project structure using Go, HTML, and CSS  
   - Built and ran a basic web server  
   ✅ Completed

2. **Data Layer Setup**  
   - Integrated SQLite as the database  
   - Defined schema for project entries  
   - Implemented data model and project retrieval logic in Go  
   ✅ Completed

3. **Dynamic Rendering**  
   - Replaced hardcoded HTML project list with dynamic content from the database  
   - Used Go templates to render the project list  
   ✅ Completed

4. **HTMX Integration**  
   - Added a search input field for real-time filtering  
   - Used HTMX to dynamically fetch filtered data from the backend  
   ✅ Completed

5. **Dockerization**  
   - Added Dockerfile and docker-compose.yaml  
   - Configured volume for SQLite persistence  
   ✅ Completed

6. **Code Cleanup & Refactor**  
   - Organized Go code into `handlers`, `utils`, and `database` packages  
   - Separated layout and partial templates for maintainability  
   ✅ Completed

7. **UI & Design Polishing**  
   - Added dark mode toggle  
   - Included SVG icons for social/contact links  
   - Improved semantic HTML and layout  
   ✅ Completed

8. **README and Documentation**  
   - Creating documentation with setup instructions  
   - Describing project features and usage  
   ✅ Completed

9. **Further Enhancements (Planned)**  
   - Admin UI to add/edit projects  
   - Tag or stack-based filtering  
   - Export/import database functionality  
   - Deployment to VPS or cloud hosting  
   ⏳ Upcoming

---

## ⚙️ Getting Started

### Prerequisites

- Go 1.21+
- SQLite3
- Docker (optional)

### Run locally

```bash
go run main.go
```

### Run with Docker

```bash
docker-compose up --build
```

---

## 🗃️ Project Structure

```
personal-portfolio/
│
├── database/
│   └── sqlite.go, models.go
│
├── handlers/
│   └── home.go, search.go
│
├── templates/
│   ├── layout.html
│   ├── index.html
│   └── partials/
│       └── project_list.html
│
├── static/
│   └── style.css
│
├── utils/
│   └── render.go
│
├── data/
│   └── data.sqlite
│
├── Dockerfile.backend
├── docker-compose.yml
└── main.go
```

---

## ✨ Features

- ⚡ Minimalistic layout and CSS
- 🌗 Dark mode toggle with localStorage persistence
- 📁 SQLite-backed project list
- 🔍 Real-time search using HTMX
- 🔧 Clean and modular Go backend
- 🐳 Dockerized for easy deployment
- 🖼️ SVG icons for contact section

---

## 📂 Notes

- `data.sqlite` is used for persistence.  
- You **can commit** the database if you want to share sample data — otherwise add it to `.gitignore`.  
- UI components are built manually using semantic HTML and lightweight CSS.  
- No frontend frameworks used (no Tailwind, Bootstrap, etc.).

---

## 📜 License

MIT — do whatever you want.
