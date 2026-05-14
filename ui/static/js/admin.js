document.addEventListener('DOMContentLoaded', () => {
    cargarUsuarios();
});

async function cargarUsuarios() {
    const tbody = document.getElementById('tabla-usuarios');
    const template = document.getElementById('molde-usuario');
    tbody.innerHTML = '<tr><td colspan="5" style="text-align: center; padding: 20px;">Cargando...</td></tr>';

    try {
        const respuesta = await fetch('/api/v1/usuarios');
        const usuarios = await respuesta.json();
        tbody.innerHTML = ''; 

        usuarios.forEach(usuario => {
            const clon = template.content.cloneNode(true);
            clon.querySelector('.user-id').textContent = usuario.id;
            clon.querySelector('.user-nombre').textContent = usuario.name; 
            clon.querySelector('.user-email').textContent = usuario.email;
            clon.querySelector('.user-rol').textContent = "user";

            // BOTÓN BORRAR
            clon.querySelector('.btn-borrar').onclick = async () => {
                if (confirm(`¿Borrar a ${usuario.name}?`)) {
                    await fetch(`/api/v1/usuarios/${usuario.id}`, { method: 'DELETE' });
                    cargarUsuarios();
                }
            };

            // BOTÓN EDITAR (Aquí es donde ocurre la magia)
            const btnEditar = clon.querySelector('.btn-editar');
            btnEditar.onclick = () => {
                console.log("Editando usuario:", usuario.id); // Mira esto en la consola (F12)
                
                // Cambiamos el título y rellenamos campos
                document.getElementById('modal-titulo').textContent = "Editar Usuario";
                document.getElementById('edit-id').value = usuario.id;
                document.getElementById('add-nombre').value = usuario.name;
                document.getElementById('add-email').value = usuario.email;
                document.getElementById('add-password').required = false; // No obligatoria al editar
                
                // Mostramos el modal
                document.getElementById('modal-usuario').style.display = 'block';
            };

            tbody.appendChild(clon);
        });
    } catch (e) { console.error("Error al cargar:", e); }
}

// Lógica del Modal (Cerrar y Abrir para nuevo)
const modal = document.getElementById('modal-usuario');

document.getElementById('btn-nuevo-usuario').onclick = () => {
    document.getElementById('form-nuevo-usuario').reset();
    document.getElementById('edit-id').value = ""; // Vacío = Nuevo
    document.getElementById('modal-titulo').textContent = "Añadir Nuevo Usuario";
    document.getElementById('add-password').required = true;
    modal.style.display = 'block';
};

document.getElementById('btn-cancelar').onclick = () => modal.style.display = 'none';

document.getElementById('form-nuevo-usuario').onsubmit = async (e) => {
    e.preventDefault();
    const id = document.getElementById('edit-id').value;
    const datos = {
        name: document.getElementById('add-nombre').value,
        email: document.getElementById('add-email').value,
        password: document.getElementById('add-password').value
    };

    const metodo = id ? 'PUT' : 'POST';
    const url = id ? `/api/v1/usuarios/${id}` : '/api/v1/usuarios';

    const res = await fetch(url, {
        method: metodo,
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(datos)
    });

    if (res.ok) {
        alert(id ? "¡Actualizado! 🦆" : "¡Creado! 🦆");
        modal.style.display = 'none';
        cargarUsuarios();
    }
};