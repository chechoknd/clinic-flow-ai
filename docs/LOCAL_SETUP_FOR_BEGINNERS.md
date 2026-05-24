# Levantar ClinicFlow AI en local

Esta guia es el camino simple para abrir la aplicacion en tu navegador.

Vas a usar 3 cosas:

- PostgreSQL con Docker.
- Backend Go.
- Frontend Angular.

Al final abriras:

```txt
http://localhost:4200
```

## Paso 1. Estar en la raiz del proyecto

Abre una terminal y entra a la carpeta del proyecto:

```bash
cd /home/toro/Documentos/workspace/monolith/clinic-flow-ai
```

Verifica que estas en la raiz:

```bash
ls
```

Debes ver archivos y carpetas como:

```txt
AGENTS.md
README.md
apps
database
docker-compose.yml
docs
```

## Paso 2. Levantar la base de datos

Ejecuta:

```bash
docker compose up -d postgres
```

Verifica que PostgreSQL quedo activo:

```bash
docker compose ps
```

Debe aparecer el servicio `postgres` en estado activo.

## Paso 3. Crear las tablas

Ejecuta:

```bash
cd apps/backend-go
DATABASE_URL="postgres://clinicflow:clinicflow@localhost:5432/clinicflow_db?sslmode=disable" go run ./cmd/migrate
cd ../..
```

Si no aparece un error, puedes seguir.

## Paso 4. Cargar datos demo

Ejecuta estos dos comandos desde la raiz del proyecto:

```bash
docker compose cp database/seeds/20260521000100_demo_core_data.sql postgres:/tmp/demo_core_data.sql
```

```bash
docker compose exec -T postgres psql -U clinicflow -d clinicflow_db -f /tmp/demo_core_data.sql
```

Esto crea el usuario demo y datos iniciales.

## Paso 5. Levantar el backend

Abre una terminal nueva.

Entra al backend:

```bash
cd /home/toro/Documentos/workspace/monolith/clinic-flow-ai/apps/backend-go
```

Si quieres usar tu archivo `.env` real, ejecuta:

```bash
set -a
source .env
set +a
```

Luego levanta el backend en el puerto que usa el frontend:

```bash
HTTP_ADDR=":18080" go run ./cmd/api
```

Deja esta terminal abierta.

Para verificar el backend, abre otra terminal y ejecuta:

```bash
curl http://127.0.0.1:18080/healthz
```

Debe responder:

```json
{"status":"ok"}
```

Luego ejecuta:

```bash
curl http://127.0.0.1:18080/readyz
```

Debe responder:

```json
{"status":"ready"}
```

## Paso 6. Levantar el frontend

Abre otra terminal nueva.

Entra al frontend:

```bash
cd /home/toro/Documentos/workspace/monolith/clinic-flow-ai/apps/frontend-angular
```

Instala dependencias si es la primera vez:

```bash
npm install
```

Levanta Angular:

```bash
npm start
```

Deja esta terminal abierta.

Cuando veas algo como `http://localhost:4200`, abre esa URL en el navegador.

## Paso 7. Entrar a la app

Abre:

```txt
http://localhost:4200
```

Usuario demo:

```txt
Email: admin@sonrisaviva.demo
Password: clinicflow123
```

## Paso 8. Probar la app

Prueba en este orden:

1. Dashboard.
2. Servicios.
3. Leads.
4. Seguimientos.
5. Asistente AI.
6. Clinica.

## Paso 9. Detener todo

Para detener backend o frontend, ve a cada terminal y presiona:

```txt
Ctrl + C
```

Para detener PostgreSQL:

```bash
docker compose down
```

## Comandos rapidos

Terminal 1, base de datos:

```bash
cd /home/toro/Documentos/workspace/monolith/clinic-flow-ai
docker compose up -d postgres
```

Terminal 2, backend:

```bash
cd /home/toro/Documentos/workspace/monolith/clinic-flow-ai/apps/backend-go
set -a
source .env
set +a
HTTP_ADDR=":18080" go run ./cmd/api
```

Terminal 3, frontend:

```bash
cd /home/toro/Documentos/workspace/monolith/clinic-flow-ai/apps/frontend-angular
npm start
```

Navegador:

```txt
http://localhost:4200
```

## Si algo falla

### No puedo iniciar sesion

Revisa que el backend este prendido:

```bash
curl http://127.0.0.1:18080/readyz
```

Si no responde `ready`, revisa PostgreSQL:

```bash
docker compose ps
```

### El frontend no abre

Revisa que `npm start` siga corriendo en la terminal del frontend.

### El backend no abre

Revisa que `go run ./cmd/api` siga corriendo en la terminal del backend.

### El Asistente AI falla

Si usas DeepSeek real, revisa que tu archivo esté en:

```txt
apps/backend-go/.env
```

Y que tenga:

```txt
AI_PROVIDER=deepseek
DEEPSEEK_API_KEY=tu-key
```

No compartas tu API key.
