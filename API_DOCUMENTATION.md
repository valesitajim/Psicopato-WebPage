# Documentación de la API de Usuarios

| Método | Ruta de Acceso | Descripción | Códigos de Estado (Respuestas) |
|---|---|---|---|
| GET | `/api/v1/usuarios` | Devuelve una lista con todos los usuarios registrados en el sistema. | **200 OK**: Éxito. Devuelve el array de usuarios en formato JSON. |
| GET | `/api/v1/usuarios/{id}` | Obtiene los detalles de un usuario específico mediante su identificador numérico. | **200 OK**: Éxito. Devuelve el objeto usuario.<br>**404 Not Found**: El usuario con ese ID no existe. |
| POST | `/api/v1/usuarios` | Crea un nuevo usuario en la base de datos (archivo JSONL). Se espera un JSON en el cuerpo de la petición con `nombre` y `email`. | **201 Created**: Creado con éxito. Devuelve el nuevo usuario y la cabecera `Location`.<br>**400 Bad Request**: El JSON está mal formado o tiene campos extra.<br>**422 Unprocessable Entity**: Faltan campos obligatorios. |
| PUT | `/api/v1/usuarios/{id}` | Actualiza completamente un usuario existente. | **200 OK**: Éxito. Devuelve el usuario actualizado.<br>**400 Bad Request**: El JSON enviado no es válido.<br>**404 Not Found**: El usuario a actualizar no existe. |
| DELETE | `/api/v1/usuarios/{id}` | Elimina a un usuario del sistema mediante su identificador. | **204 No Content**: Éxito al borrar. No devuelve cuerpo.<br>**404 Not Found**: El usuario no existe en la base de datos. |