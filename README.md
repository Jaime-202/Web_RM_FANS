# 🏟️ Real Madrid News — Portal de Noticias

**PEC 1 — Redes y Sistemas Web**  
**Grupo 11 | Jaime & Izan**

Portal de noticias dinámico sobre el Real Madrid, desarrollado con Go en el backend y CSS/HTML en el frontend. Incluye un sistema completo de usuarios con roles, gestión de contenido y valoraciones por estrellas.

---

## 📋 Tabla de Contenidos

- [Descripción](#descripción)
- [Tecnologías](#tecnologías)
- [Estructura del Proyecto](#estructura-del-proyecto)
- [Funcionalidades](#funcionalidades)
- [Roles de Usuario](#roles-de-usuario)
- [Cómo Ejecutar](#cómo-ejecutar)
- [Credenciales de Prueba](#credenciales-de-prueba)
- [Rutas Disponibles](#rutas-disponibles)

---

## 📖 Descripción

Real Madrid News es una aplicación web dinámica que permite consultar noticias, ver la plantilla completa de jugadores y seguir el mercado de fichajes y rumores del Real Madrid.

El sistema cuenta con autenticación de usuarios, dos niveles de acceso (Administrador y Usuario registrado) y persistencia de datos mediante archivos JSON, sin necesidad de una base de datos externa.

---

## ⚙️ Tecnologías

| Capa | Tecnología |
|---|---|
| **Backend** | Go 1.21+ (`net/http`, `html/template`) |
| **Frontend** | HTML5, CSS3, SCSS |
| **Persistencia** | Archivos JSON (sin base de datos externa) |
| **Sesiones** | Gestión en memoria con cookies `HttpOnly` |
| **Dependencias** | `github.com/google/uuid` |

---

## 📁 Estructura del Proyecto

```
RSW-PEC1_Grupo11-Jaime-Izan/
│
├── css/                          # CSS compilado (estilos.css)
├── scss/                         # Código fuente SCSS
├── img/                          # Imágenes del portal
│   └── plantilla/                # Fotos de los jugadores
│
└── proyecto-go/                  # Aplicación Go (backend completo)
    │
    ├── cmd/server/
    │   └── main.go               # Punto de entrada: configura rutas y arranca el servidor
    │
    ├── data/                     # Base de datos en JSON (se genera automáticamente)
    │   ├── users.json            # Usuarios registrados
    │   ├── news.json             # Noticias publicadas
    │   ├── players.json          # Plantilla de jugadores
    │   ├── transfers.json        # Fichajes y rumores
    │   ├── ratings.json          # Valoraciones de usuarios
    │   └── contacts.jsonline     # Mensajes del buzón de rumores
    │
    ├── internal/
    │   ├── models/               # Estructuras de datos
    │   │   ├── user.go           # Usuario (ID, nombre, contraseña, rol)
    │   │   ├── news.go           # Noticia (título, contenido, imagen, valoración)
    │   │   ├── player.go         # Jugador (nombre, posición, imagen)
    │   │   ├── transfer.go       # Fichaje (jugador, equipos, estado, descripción)
    │   │   ├── rating.go         # Valoración (usuario, noticia, puntuación 1-5)
    │   │   └── contact.go        # Mensaje de contacto
    │   │
    │   ├── database/             # Repositorios de acceso a los JSON
    │   │   ├── user_repository.go
    │   │   ├── news_repository.go
    │   │   ├── player_repository.go
    │   │   ├── transfer_repository.go
    │   │   ├── rating_repository.go
    │   │   └── json_repository.go
    │   │
    │   ├── handlers/             # Controladores HTTP (uno por sección)
    │   │   ├── home_handler.go   # Página de inicio
    │   │   ├── auth_handler.go   # Login, registro y logout
    │   │   ├── news_handler.go   # CRUD de noticias + valoraciones
    │   │   ├── player_handler.go # CRUD de la plantilla
    │   │   ├── transfer_handler.go # CRUD de fichajes y rumores
    │   │   └── contact_handler.go  # Formulario de contacto
    │   │
    │   ├── middleware/
    │   │   └── auth.go           # Middleware de autenticación y roles
    │   │
    │   ├── services/
    │   │   └── contact_service.go # Lógica del formulario de contacto
    │   │
    │   └── session/
    │       └── store.go          # Gestión de sesiones en memoria
    │
    └── web/templates/            # Plantillas HTML dinámicas
        ├── layout.html           # Cabecera y pie de página compartidos (con gestión de sesión)
        ├── index.html            # Página de inicio (preview de noticias, jugadores, fichajes)
        ├── news_list.html        # Lista completa de noticias
        ├── news_detail.html      # Vista de una noticia + sistema de estrellas
        ├── news_form.html        # Formulario para crear/editar noticias (admin)
        ├── plantilla.html        # Plantilla completa agrupada por posición
        ├── player_form.html      # Formulario para añadir/editar jugadores (admin)
        ├── fichajes.html         # Mercado de fichajes y rumores
        ├── transfer_form.html    # Formulario para gestionar fichajes (admin)
        ├── contact.html          # Buzón de rumores para fans
        ├── login.html            # Página de inicio de sesión
        └── register.html         # Página de registro
```

---

## ✨ Funcionalidades

### 🌐 Para todos los visitantes
- Consultar las últimas noticias del club
- Ver la plantilla completa de jugadores (25 jugadores) agrupados por posición: Porteros, Defensas, Centrocampistas y Delanteros
- Seguir el mercado de fichajes y rumores con su estado actualizado
- Enviar rumores y exclusivas a través del buzón de contacto

### ⭐ Para usuarios registrados
- Todo lo anterior, más:
- **Valorar noticias** con un sistema interactivo de 1 a 5 estrellas
- La puntuación media de cada noticia se actualiza automáticamente en tiempo real
- Cada usuario solo puede votar una vez por noticia (puede cambiar su voto)

### 🛡️ Para el administrador
- Todo lo anterior, más:
- **Noticias**: Crear, editar y borrar artículos de noticias
- **Plantilla**: Añadir, editar y eliminar jugadores de la plantilla
- **Fichajes**: Añadir nuevos rumores, editar su información y cambiar el **estado** (`Rumor` → `Contrato verbal` → `Hecho`), o eliminarlos

---

## 👤 Roles de Usuario

| Rol | Capacidades |
|---|---|
| **Visitante** | Leer noticias, ver plantilla y fichajes, enviar mensajes de contacto |
| **Usuario registrado** | Todo lo anterior + valorar noticias con estrellas (1-5) |
| **Administrador** | Todo lo anterior + CRUD completo de noticias, jugadores y fichajes |

---

## 🚀 Cómo Ejecutar

### Requisitos previos
- [Go 1.21 o superior](https://go.dev/dl/) instalado en el sistema

### Pasos

**1. Clonar o descargar el proyecto**
```bash
git clone <url-del-repositorio>
cd RSW-PEC1_Grupo11-Jaime-Izan
```

**2. Instalar dependencias**
```bash
cd proyecto-go
go mod tidy
```

**3. Arrancar el servidor**
```bash
# Desde dentro de la carpeta proyecto-go
go run ./cmd/server/main.go
```

**4. Abrir en el navegador**
```
http://localhost:8080
```

> **Importante:** El servidor debe ejecutarse siempre **desde dentro de la carpeta `proyecto-go`**, ya que las rutas a los archivos estáticos (`/css/`, `/img/`) y a los datos (`data/`) son relativas a esa carpeta.

---

## 🔑 Credenciales de Prueba

### Administrador (preconfigurado)
| Campo | Valor |
|---|---|
| **Usuario** | `admin` |
| **Contraseña** | `admin` |

### Usuario normal
Crea tu propia cuenta yendo a [http://localhost:8080/register](http://localhost:8080/register) e introduciendo cualquier nombre de usuario y contraseña.

---

## 🗺️ Rutas Disponibles

### Públicas (sin sesión requerida)
| Ruta | Descripción |
|---|---|
| `GET /` | Página de inicio con últimas noticias, jugadores destacados y fichajes |
| `GET /noticias` | Lista completa de noticias |
| `GET /noticia?id=X` | Detalle de una noticia |
| `GET /plantilla` | Plantilla completa de jugadores por posición |
| `GET /fichajes` | Mercado de fichajes y rumores |
| `GET /contacto` | Formulario de buzón de rumores |
| `GET/POST /login` | Inicio de sesión |
| `GET/POST /register` | Registro de nuevo usuario |
| `GET /logout` | Cerrar sesión |

### Usuarios registrados
| Ruta | Descripción |
|---|---|
| `POST /news/rate` | Enviar valoración (1-5 estrellas) a una noticia |

### Solo Administradores
| Ruta | Descripción |
|---|---|
| `GET/POST /admin/news/create` | Crear nueva noticia |
| `GET/POST /admin/news/edit?id=X` | Editar noticia existente |
| `GET /admin/news/delete?id=X` | Eliminar noticia |
| `GET/POST /admin/player/create` | Añadir jugador a la plantilla |
| `GET/POST /admin/player/edit?id=X` | Editar datos de un jugador |
| `GET /admin/player/delete?id=X` | Eliminar jugador |
| `GET/POST /admin/transfer/create` | Añadir fichaje o rumor |
| `GET/POST /admin/transfer/edit?id=X` | Editar estado o datos de un fichaje |
| `GET /admin/transfer/delete?id=X` | Eliminar fichaje |

---

## 🔒 Seguridad

- Las contraseñas se almacenan con hash **SHA-256** (nunca en texto plano)
- Las sesiones usan cookies `HttpOnly` para evitar acceso desde JavaScript
- Los middlewares verifican el rol del usuario antes de permitir acceso a rutas protegidas
- Los tokens de sesión se generan con UUIDs únicos

---

*© 2026 Real Madrid News — PEC 1, Redes y Sistemas Web*
