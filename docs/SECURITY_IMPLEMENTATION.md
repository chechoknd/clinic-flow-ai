# Implementación de Seguridad - ClinicFlow AI

Este documento detalla las medidas de seguridad implementadas en el backend de ClinicFlow AI para garantizar la integridad de los datos, la protección contra ataques comunes y el aislamiento de inquilinos (clínicas).

## 1. Seguridad de Infraestructura y Servidor

Se han implementado protecciones a nivel de red y servidor para prevenir ataques de denegación de servicio (DoS) y abuso de recursos:

- **Timeouts del Servidor:** Configuración estricta de tiempos de espera para evitar conexiones lentas o inactivas que agoten los recursos:
    - `ReadHeaderTimeout`: 5s
    - `ReadTimeout`: 10s
    - `WriteTimeout`: 15s
    - `IdleTimeout`: 120s
- **Límite de Payloads:** Middleware global que restringe el tamaño del cuerpo de las solicitudes JSON a **1MB**, previniendo ataques por agotamiento de memoria.
- **Rate Limiting:** Implementación de un limitador de tasa basado en IP que permite un máximo de **60 peticiones por minuto**, protegiendo contra ataques de fuerza bruta y spam.

## 2. Cabeceras de Seguridad HTTP

Todas las respuestas del API incluyen cabeceras estándar de seguridad para proteger a los clientes (navegadores) y prevenir ataques de inyección:

- `X-Content-Type-Options: nosniff`: Evita que el navegador interprete archivos como un tipo MIME diferente al declarado.
- `X-Frame-Options: DENY`: Protege contra ataques de Clickjacking.
- `X-XSS-Protection: 1; mode=block`: Activa protecciones contra Cross-Site Scripting en navegadores compatibles.
- `Strict-Transport-Security (HSTS)`: Fuerza el uso de conexiones seguras (HTTPS).
- `Content-Security-Policy (CSP)`: Restringe el origen de los recursos permitidos.
- `Referrer-Policy`: Controla la información de referencia enviada en las peticiones.

## 3. Autenticación y Autorización

- **JWT (JSON Web Tokens):**
    - Firma segura mediante algoritmos HMAC (HS256).
    - **Validación Crítica:** El servidor no arranca si el secreto de firma (`JWT_SECRET`) tiene menos de 32 caracteres.
    - Claims incluyen `user_id`, `clinic_id` y `role` para un control granular.
- **RBAC (Role Based Access Control):** Middleware específico para validar roles de usuario (`superadmin`, `clinic_admin`, `assistant`) antes de acceder a recursos sensibles.
- **Multi-tenancy (Aislamiento):** Cada solicitud autenticada inyecta el `clinic_id` en el contexto. Todos los repositorios filtran obligatoriamente por este ID, garantizando que una clínica nunca pueda acceder a los datos de otra.

## 4. Validación de Negocio y Datos

- **Hashing de Contraseñas:** Uso de `bcrypt` con costo adaptativo para el almacenamiento seguro de credenciales.
- **Validación Estricta de Entradas:**
    - Teléfonos y WhatsApp validados bajo el estándar internacional **E.164**.
    - Emails validados mediante analizadores sintácticos robustos.
    - Listas blancas (Allow-lists) para estados de leads y tonos de comunicación.
- **Prevención de Inyección SQL:** Uso exclusivo de consultas preparadas (`prepared statements`) y parámetros posicionales en todas las interacciones con PostgreSQL.

## 5. Configuración y CORS

- **CORS (Cross-Origin Resource Sharing):** Configuración dinámica que permite restringir el acceso del API a dominios específicos (ej: el dominio del frontend Angular) mediante la variable `ALLOWED_ORIGINS`.
- **Sanitización de Errores:** Los errores internos del sistema (SQL, lógica interna) no se exponen al cliente. El API responde con códigos de error genéricos y controlados (`INTERNAL_ERROR`).
