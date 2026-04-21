# Real Madrid News

**Práctica de Evaluación Continua 1 (PEC 1) - Redes y Sistemas Web**

## Descripción del Proyecto
Real Madrid News es un portal web de noticias dedicado a cubrir toda la actualidad del Real Madrid. Este proyecto engloba una interfaz de usuario atractiva desarrollada con HTML5 y CSS/SCSS, además de un servidor backend en Go (Golang) que sirve los archivos estáticos y maneja un formulario de contacto.

## Autores
*   **Grupo 11:** Jaime e Izan
*   **Asignatura:** Redes y Sistemas Web
*   **Año:** 2º, Semestre 2

## Tecnologías Utilizadas
*   **Frontend:** HTML5, CSS3 (código base procesado mediante preprocesador **SCSS** para asegurar una mejor estructura, variables y mantenibilidad).
*   **Backend:** **Go (Golang)** utilizando la biblioteca estándar `net/http`.
*   **Arquitectura del Backend:** Diseño dividido en capas (Presentación/Handlers, Lógica de Negocio/Services, y Datos/Repository).
*   **Almacenamiento:** Los datos del formulario de contacto se persisten de manera local en formato `JSONLine`.

## Estructura del Repositorio

```text
/
├── css/             # Archivos CSS finales compilados
├── img/             # Imágenes utilizadas en el portal (noticias, plantilla, etc.)
├── pages/           # Páginas HTML secundarias (noticias.html, plantilla.html, fichajes.html)
├── scss/            # Código fuente de estilos estructurado (BEM u otras metodologías)
├── index.html       # Página principal del portal
└── proyecto-go/     # Código y ejecutable del servidor backend en Go
    ├── cmd/         # Punto de entrada de la aplicación Go (main.go)
    ├── data/        # Almacenamiento de datos persistentes locales (contacts.jsonline)
    ├── internal/    # Lógica de la aplicación Go (database, handlers, services)
    ├── web/         # Plantillas HTML específicas del formulario y servidor
    ├── go.mod       # Fichero de dependencias de Go
    └── servidor.exe # Ejecutable compilado del servidor (Windows)
```

## Características Principales de la Web
1.  **Página de Inicio:** Resumen visual rápido de las últimas noticias, visualización de una parte de la plantilla y los rumores de fichajes más recientes con un diseño 'hero'.
2.  **Noticias:** Sección ampliada de artículos con la actualidad del club.
3.  **Plantilla:** Tarjetas con el listado de jugadores, imagen y posición en el campo.
4.  **Fichajes:** Panel para el seguimiento de rumores y traspasos con indicadores de estado ("Hecho", etc.).
5.  **Contacto Interactivo (Backend):** Un formulario de contacto real impulsado por el servidor de Go (`/contacto`) para que los usuarios puedan enviar rumores. El servidor atiende la petición POST y la almacena.

## Instrucciones de Instalación y Ejecución

Para poder visualizar correctamente el portal y que funcionen las rutas (como el formulario de contacto), es indispensable correr el servidor en Go incluido.

### Opción 1: Ejecutar el binario directo (Recomendado para Windows)
1. Abre tu Explorador de archivos o una terminal.
2. Navega hasta el interior de la carpeta `proyecto-go`.
3. Ejecuta el archivo **`servidor.exe`** (puedes hacer doble clic sobre él o ejecutar `.\servidor.exe` desde la consola).
4. El servidor arrancará en segundo plano.
5. Abre tu navegador web favorito y dirígete a: [http://localhost:8080](http://localhost:8080)

### Opción 2: Ejecutar desde el código fuente (Requiere Go instalado)
1. Abre una terminal.
2. Ubícate en la carpeta del backend:
   ```bash
   cd proyecto-go
   ```
3. Arranca el servidor ejecutando el código fuente:
   ```bash
   go run cmd/server/main.go
   ```
4. Abre tu navegador y dirígete a: [http://localhost:8080](http://localhost:8080)

> ⚠️ **Nota Importante:** Es fundamental ejecutar el servidor posicionándose **dentro** de la carpeta `proyecto-go`. El servidor está programado internamente (`http.Dir("../")`) para servir los archivos HTML subiendo un nivel en el árbol de directorios.

---
© 2026 Real Madrid News - PEC 1 - Redes y Sistemas Web
