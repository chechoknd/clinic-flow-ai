# Manual de usuario de ClinicFlow AI

Este manual explica como usar los modulos actuales de ClinicFlow AI desde el navegador.

ClinicFlow AI es una herramienta comercial para clinicas privadas. Ayuda al equipo humano a registrar leads, hacer seguimiento manual y generar respuestas comerciales asistidas por AI para WhatsApp.

No es un sistema medico, no guarda historias clinicas, no diagnostica, no prescribe y no envia mensajes automaticos por WhatsApp.

## Estado actual y evolucion planeada

Status: Current manual covers implemented screens. Inbox AI manual conversation analysis is implemented as a first Smart Lead Inbox slice; broader inbox history and integrations remain planned.

La aplicacion actual permite gestionar leads, seguimientos, servicios, perfil de clinica, dashboard y respuestas asistidas por AI.

La evolucion se llama `Smart Lead Inbox` o `Inbox AI`. La primera version permite pegar una conversacion completa, analizarla con AI, detectar interes comercial, detectar objeciones, sugerir una respuesta, sugerir una proxima accion y crear un lead despues de revision humana. Actualizar leads existentes, historial persistente e integraciones externas siguen planeados.

Regla obligatoria:

```txt
AI sugiere -> humano revisa -> humano responde
```

La funcionalidad planeada no enviara mensajes automaticamente. El usuario seguira revisando, copiando y enviando manualmente la respuesta por el canal original.

## 1. Ingreso al sistema

Abre la aplicacion en:

```txt
http://localhost:4200
```

En la pantalla de ingreso escribe tu email y password.

Usuario demo local:

```txt
Email: admin@sonrisaviva.demo
Password: clinicflow123
```

Haz clic en `Entrar`.

Si los datos son correctos, veras el menu principal con:

- Atender hoy.
- Leads.
- Inbox AI.
- Seguimientos.
- Asistente AI.
- Servicios.
- Clinica.

Para cerrar sesion, usa el boton `Salir`.

## 2. Atender hoy

`Atender hoy` es la pantalla inicial de trabajo diario. Combina acciones prioritarias con metricas comerciales.

La cola de accion prioriza:

- Seguimientos vencidos.
- Seguimientos programados para hoy.
- Leads nuevos que requieren primer contacto.

Desde cada accion puedes ir a gestionar el seguimiento, responder con AI, abrir el pipeline o analizar una conversacion en Inbox AI.

La misma pantalla tambien muestra:

- Total de leads.
- Seguimientos pendientes para hoy.
- Seguimientos vencidos.
- Tasa de conversion.
- Pipeline por estado.
- Servicios mas consultados.


## Inbox AI

Inbox AI permite analizar una conversacion comercial pegada manualmente.

Uso recomendado:

1. Abre `Inbox AI`.
2. Selecciona el canal de origen.
3. Opcionalmente selecciona un servicio de referencia.
4. Pega la conversacion comercial.
5. Haz clic en `Analizar con AI`.
6. Revisa intencion, servicio detectado, objeciones, resumen, respuesta sugerida y siguiente accion.
7. Copia la respuesta sugerida solo si esta correcta.
8. Revisa los campos detectados antes de crear el lead.
9. Haz clic en la accion principal: `Crear lead revisado` o `Actualizar lead revisado`, segun si seleccionaste un lead existente.

No pegues historias clinicas, diagnosticos, recetas ni informacion medica sensible. La respuesta sugerida no se envia automaticamente.

## 3. Leads

El modulo de Leads es el CRM comercial simple. Sirve para registrar personas interesadas y llevar su seguimiento manual.

### Estados comerciales

Los leads pueden tener estos estados:

- Nuevo.
- Contactado.
- Interesado.
- Agendado.
- No Respondio.
- Perdido.
- Convertido.

La idea es mover cada lead por el pipeline segun avance la conversacion comercial.

### Crear un lead

En la seccion `Nuevo lead`, completa:

- `Nombre completo`: nombre de la persona interesada.
- `WhatsApp`: numero en formato internacional, por ejemplo `+573001234567`.
- `Servicio`: tratamiento o servicio de interes. Puede quedar sin definir.
- `Estado`: normalmente inicia como `Nuevo`.
- `Origen`: canal por donde llego el lead.
- `Proxima accion`: fecha y hora para volver a contactar.
- `Nota inicial`: resumen comercial de la consulta.

Luego haz clic en `Crear lead`.

Buenas practicas:

- Usa notas cortas y utiles.
- No escribas diagnosticos.
- No guardes informacion clinica sensible.
- Registra solo informacion comercial necesaria para seguimiento.

### Filtrar leads por estado

En la parte superior de la lista puedes seleccionar un estado:

- Nuevo.
- Contactado.
- Interesado.
- Agendado.
- No Respondio.
- Perdido.
- Convertido.

La lista mostrara solo los leads que pertenecen al estado seleccionado. Cada filtro muestra su conteo actual para ayudar a entender donde esta concentrado el trabajo.

### Editar seguimiento de un lead

En la lista, haz clic en `Editar` sobre un lead. Si el lead tiene una proxima accion vencida o para hoy, la fila aparece resaltada. Tambien puedes usar `Responder` para abrir el Asistente AI con contexto del lead.

Luego puedes cambiar:

- `Estado`.
- `Proxima accion`.
- `Nota de seguimiento`.

Haz clic en `Guardar seguimiento`.

Esto sirve para registrar lo que paso despues de hablar con el paciente potencial.

Ejemplos de notas comerciales validas:

- "Pregunta por disponibilidad en la tarde."
- "Quiere cotizacion aproximada de blanqueamiento."
- "Pidio que la contacten manana."
- "Se envio informacion comercial por WhatsApp."

Evita notas como:

- Diagnosticos.
- Prescripciones.
- Resultados clinicos.
- Historias medicas.

### Responder con AI desde un lead

Despues de seleccionar un lead, puedes usar el boton `Responder`.

Esto abre el `Asistente AI` con contexto del lead y del servicio, cuando existan.

La respuesta generada debe revisarla una persona antes de enviarla por WhatsApp.

## 4. Seguimientos

El modulo de Seguimientos muestra leads que tienen una proxima accion pendiente, agrupados por urgencia: `Vencidos`, `Hoy` y `Proximos`.

Sirve para que el equipo comercial no olvide recuperar conversaciones.

Cada tarjeta muestra:

- Nombre del lead.
- WhatsApp.
- Estado actual.
- Servicio relacionado.
- Fecha de proxima accion.

### Gestionar un seguimiento

Haz clic en `Gestionar` sobre un seguimiento.

En el panel derecho puedes:

- Cambiar el estado al completar.
- Elegir una nueva proxima accion.
- Escribir una nota.

### Completar seguimiento

Usa `Completar` cuando ya hiciste la accion pendiente.

Ejemplos:

- Llamaste al lead.
- Enviaste respuesta manual por WhatsApp.
- Confirmaste si sigue interesado.
- Cerraste la oportunidad como perdida o convertida.

### Reprogramar seguimiento

Usa `Reprogramar` cuando necesitas contactar de nuevo en otra fecha.

Debes seleccionar `Nueva proxima accion`.

Ejemplo:

- El lead pidio que le escriban el viernes.
- No respondio y quieres intentar de nuevo manana.
- Quieres confirmar decision despues de enviar informacion.

### Generar mensaje AI

Con un seguimiento seleccionado, puedes hacer clic en `Generar mensaje AI`.

El sistema genera un texto sugerido para retomar la conversacion.

Luego puedes usar `Copiar` y pegarlo manualmente en WhatsApp.

Importante:

- La aplicacion no envia mensajes automaticamente.
- La respuesta debe ser revisada antes de enviarse.
- No uses AI para diagnosticar o dar instrucciones medicas.

## 5. Catalogo de servicios

El modulo de Servicios define la oferta comercial de la clinica.

Esta informacion tambien alimenta el contexto del Asistente AI.

Cada servicio puede tener:

- Nombre.
- Descripcion comercial.
- Duracion estimada.
- Precio desde.
- Beneficios.
- Objeciones comunes.
- Estado activo o inactivo.

### Crear un servicio

Haz clic en `Nuevo servicio`.

Completa:

- `Nombre`: por ejemplo, "Blanqueamiento dental".
- `Descripcion comercial`: explicacion clara para pacientes.
- `Duracion min`: duracion aproximada.
- `Precio desde`: valor base orientativo.
- `Beneficios`: ventajas comerciales del servicio.
- `Objeciones comunes`: dudas frecuentes como precio, dolor, tiempo o seguridad.

Haz clic en `Crear servicio`.

### Editar un servicio

En una tarjeta de servicio, haz clic en `Editar`.

Actualiza los campos necesarios y haz clic en `Guardar cambios`.

Tambien puedes marcar o desmarcar `Servicio activo`.

### Eliminar un servicio

Selecciona el servicio con `Editar` y luego haz clic en `Eliminar`.

Usa esta accion con cuidado. Si solo quieres que no aparezca como disponible, suele ser mejor marcarlo como inactivo.

### Buenas practicas para servicios

Escribe informacion comercial, no clinica.

Correcto:

- "Ideal para mejorar la apariencia de la sonrisa."
- "Incluye valoracion inicial comercial y explicacion del proceso."
- "Precio desde segun evaluacion de la clinica."

Evita:

- Diagnosticos.
- Promesas medicas absolutas.
- Indicaciones clinicas personalizadas.
- Garantias de resultado.

## 6. Asistente AI

El Asistente AI genera respuestas comerciales listas para copiar a WhatsApp.

No es un chatbot autonomo. No envia mensajes. Solo sugiere texto.

### Tipos de ayuda

Puedes elegir:

- `Respuesta a pregunta`: para contestar dudas generales del paciente potencial.
- `Manejo de objecion`: para responder cuando alguien dice que es caro, tiene miedo, no tiene tiempo o quiere pensarlo.

### Generar una respuesta

1. Selecciona el tipo de ayuda.
2. Escribe el mensaje del paciente.
3. Haz clic en `Generar respuesta`.
4. Revisa la respuesta sugerida.
5. Haz clic en `Copiar`.
6. Pega manualmente el texto en WhatsApp si lo apruebas.

Ejemplos de mensajes que puedes pegar:

```txt
Me interesa el blanqueamiento, pero me parece costoso.
```

```txt
Quiero saber cuanto cuesta una limpieza dental.
```

```txt
Tengo miedo de que el tratamiento duela.
```

### Contexto activo

Cuando entras al Asistente AI desde un lead o seguimiento, puede aparecer el aviso `Contexto activo para esta respuesta`.

Eso significa que la respuesta puede usar contexto del lead o servicio seleccionado.

### Reglas de seguridad del Asistente AI

Antes de enviar una respuesta:

- Revisa que sea comercial y clara.
- Verifica que no diagnostique.
- Verifica que no prometa resultados medicos.
- Verifica que no indique medicamentos o tratamientos personalizados.
- Ajusta el tono si hace falta.

Si una persona pregunta algo medico, la respuesta debe orientar a una valoracion con la clinica, no resolver el caso clinico por chat.

## 7. Perfil de clinica

El modulo `Clinica` permite editar informacion comercial usada por el sistema y por las respuestas asistidas.

Puedes ver:

- Nombre.
- Tipo de clinica.
- Ciudad.
- WhatsApp.
- Telefono.
- Tono de comunicacion.
- Direccion.

### Editar datos comerciales

En el panel `Editar datos comerciales`, puedes actualizar:

- Nombre.
- Ciudad.
- WhatsApp.
- Telefono.
- Direccion.
- Tono de comunicacion.

Haz clic en `Guardar perfil`.

### Tonos de comunicacion

Los tonos disponibles son:

- Amable.
- Profesional.
- Cercano.
- Juvenil.
- Elegante.

El tono ayuda a que las respuestas asistidas se ajusten a la voz de la clinica.

Ejemplo:

- `Profesional`: mas formal y directo.
- `Cercano`: mas calido y conversacional.
- `Juvenil`: mas fresco, sin perder claridad.

## 8. Flujo recomendado de trabajo diario

1. Entrar al Dashboard para revisar pendientes.
2. Ir a Seguimientos y gestionar acciones del dia.
3. Crear nuevos leads desde el modulo Leads.
4. Actualizar estado y notas despues de cada conversacion.
5. Usar Asistente AI para redactar mejores respuestas.
6. Copiar respuestas aprobadas y enviarlas manualmente por WhatsApp.
7. Mantener actualizado el Catalogo de servicios.
8. Ajustar el Perfil de clinica cuando cambien datos comerciales.

## 9. Que no se debe hacer en ClinicFlow AI

No uses la plataforma para:

- Guardar historias clinicas.
- Registrar diagnosticos.
- Registrar prescripciones.
- Registrar imagenes medicas.
- Tomar decisiones clinicas.
- Enviar mensajes automaticos por WhatsApp.
- Reemplazar la valoracion de un profesional.

La plataforma esta hecha para atencion comercial asistida, no para atencion medica.

## 10. Problemas frecuentes

### No puedo entrar

Verifica:

- Email correcto.
- Password correcto.
- Backend corriendo.
- Base de datos activa.

En local, el usuario demo es:

```txt
admin@sonrisaviva.demo
clinicflow123
```

### El Dashboard no carga datos

Puede faltar:

- Backend activo.
- Sesion iniciada.
- Datos demo cargados.
- Leads registrados.

### No aparecen servicios al crear lead

Revisa el modulo `Servicios`.

Si no hay servicios, crea uno nuevo.

Si existen pero no aparecen, verifica que el backend este funcionando.

### El Asistente AI no responde

Posibles causas:

- Backend apagado.
- API key de proveedor AI ausente o invalida.
- Proveedor AI sin conexion.
- Mensaje bloqueado por reglas de seguridad.

### Un seguimiento no aparece

Los seguimientos dependen de `Proxima accion`.

Revisa que el lead tenga una fecha y hora de proxima accion.

## 11. Recomendaciones para el equipo comercial

- Registra cada lead apenas llegue.
- Mantén notas breves y comerciales.
- Agenda siempre una proxima accion si el lead sigue vivo.
- Cambia el estado despues de cada contacto.
- Usa AI para mejorar redaccion, no para reemplazar criterio humano.
- Revisa todo antes de enviarlo por WhatsApp.
- No guardes informacion clinica sensible.
