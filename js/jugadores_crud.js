// MVC Structure for Jugadores (Plantilla) CRUD

class PlayerModel {
    #apiUrl = '/api/players';

    async getAll() {
        const response = await fetch(this.#apiUrl);
        if (!response.ok) throw new Error('Error al obtener jugadores');
        return await response.json();
    }

    async create(data) {
        const response = await fetch(this.#apiUrl, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });
        if (!response.ok) throw new Error('Error al crear jugador');
        return await response.json();
    }

    async update(data) {
        const response = await fetch(this.#apiUrl, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });
        if (!response.ok) throw new Error('Error al actualizar jugador');
        return await response.json();
    }

    async delete(id) {
        const response = await fetch(`${this.#apiUrl}?id=${id}`, {
            method: 'DELETE'
        });
        if (!response.ok) throw new Error('Error al eliminar jugador');
        return await response.json();
    }
}

class PlayerView {
    #container;
    #playerTemplate;
    #groupTemplate;
    #modal;
    #form;
    #modalTitle;
    #isAdmin;
    
    // Posiciones en orden de visualización
    #positions = [
        { key: 'Portero', title: 'Porteros' },
        { key: 'Defensa', title: 'Defensas' },
        { key: 'Centrocampista', title: 'Centrocampistas' },
        { key: 'Delantero', title: 'Delanteros' }
    ];

    constructor(isAdmin) {
        this.#isAdmin = isAdmin;
        this.#container = document.getElementById('players-container');
        this.#playerTemplate = document.getElementById('player-template');
        this.#groupTemplate = document.getElementById('position-group-template');
        
        if (this.#isAdmin) {
            this.#modal = document.getElementById('player-modal');
            this.#form = document.getElementById('player-form');
            this.#modalTitle = document.getElementById('modal-title');
        }
    }

    renderList(players) {
        this.#container.textContent = ''; 
        
        if (!players || players.length === 0) {
            const p = document.createElement('p');
            p.textContent = 'No hay jugadores en la plantilla.';
            p.style.textAlign = 'center'; // Solo para placeholder genérico
            p.style.color = '#666';
            p.style.marginTop = '40px';
            this.#container.appendChild(p);
            return;
        }

        // Agrupar jugadores
        const groupedPlayers = {};
        this.#positions.forEach(pos => {
            groupedPlayers[pos.key] = [];
        });

        players.forEach(player => {
            const pos = player.Position || player.position;
            if (groupedPlayers[pos]) {
                groupedPlayers[pos].push(player);
            } else {
                // Posición desconocida, crear categoría si no existe (no debería pasar normalmente)
                if (!groupedPlayers[pos]) groupedPlayers[pos] = [];
                groupedPlayers[pos].push(player);
            }
        });

        // Renderizar por grupos
        this.#positions.forEach(posDef => {
            const groupPlayers = groupedPlayers[posDef.key];
            if (groupPlayers && groupPlayers.length > 0) {
                const groupClone = this.#groupTemplate.content.cloneNode(true);
                const groupDiv = groupClone.querySelector('.position-group');
                const groupTitle = groupClone.querySelector('.position-title');
                const groupGrid = groupClone.querySelector('.position-grid');

                groupDiv.dataset.position = posDef.key;
                groupTitle.textContent = posDef.title;
                groupDiv.classList.remove('hidden');

                groupPlayers.forEach(player => {
                    const playerClone = this.#playerTemplate.content.cloneNode(true);
                    const card = playerClone.querySelector('.player-mini');
                    if(card) card.dataset.id = player.ID || player.id;

                    const img = playerClone.querySelector('.player-mini__img');
                    if(img) {
                        img.src = player.ImageURL || player.image_url || '';
                        img.alt = player.Name || player.name || '';
                    }

                    const name = playerClone.querySelector('.player-mini__name');
                    if(name) name.textContent = player.Name || player.name || 'Desconocido';

                    if (this.#isAdmin) {
                        const btnEdit = playerClone.querySelector('.btn-edit');
                        if(btnEdit) {
                            btnEdit.addEventListener('click', () => {
                                const evt = new CustomEvent('editPlayer', { detail: player });
                                document.dispatchEvent(evt);
                            });
                        }

                        const btnDelete = playerClone.querySelector('.btn-delete');
                        if(btnDelete) {
                            btnDelete.addEventListener('click', () => {
                                if (confirm(`¿Seguro que quieres borrar a ${player.Name || player.name}?`)) {
                                    const evt = new CustomEvent('deletePlayer', { detail: player.ID || player.id });
                                    document.dispatchEvent(evt);
                                }
                            });
                        }
                    }

                    groupGrid.appendChild(playerClone);
                });

                this.#container.appendChild(groupClone);
            }
        });
    }

    showError(msg) {
        this.#container.textContent = '';
        const p = document.createElement('p');
        p.textContent = msg;
        p.classList.add('error-message');
        this.#container.appendChild(p);
    }

    showLoading() {
        this.#container.textContent = '';
        const p = document.createElement('p');
        p.textContent = 'Cargando jugadores...';
        this.#container.appendChild(p);
    }

    // Modal control
    openModal(player = null) {
        if (!this.#isAdmin) return;
        
        this.#modal.classList.remove('hidden');
        if (player) {
            this.#modalTitle.textContent = 'Editar Jugador';
            this.#form.elements['id'].value = player.ID || player.id;
            this.#form.elements['name'].value = player.Name || player.name;
            this.#form.elements['position'].value = player.Position || player.position;
            this.#form.elements['image_url'].value = player.ImageURL || player.image_url;
        } else {
            this.#modalTitle.textContent = 'Añadir Jugador';
            this.#form.reset();
            this.#form.elements['id'].value = '';
        }
    }

    closeModal() {
        if (!this.#isAdmin) return;
        this.#modal.classList.add('hidden');
        this.#form.reset();
    }
}

class PlayerController {
    #model;
    #view;
    #isAdmin;

    constructor(model, view, isAdmin) {
        this.#model = model;
        this.#view = view;
        this.#isAdmin = isAdmin;

        this.#init();
    }

    async #init() {
        if (this.#isAdmin) {
            this.#setupAdminEvents();
        }
        await this.#loadPlayers();
    }

    async #loadPlayers() {
        try {
            this.#view.showLoading();
            const players = await this.#model.getAll();
            this.#view.renderList(players);
        } catch (error) {
            console.error(error);
            this.#view.showError('Ocurrió un error al cargar la plantilla.');
        }
    }

    #setupAdminEvents() {
        const btnCreate = document.getElementById('btn-create-player');
        if (btnCreate) {
            btnCreate.addEventListener('click', () => this.#view.openModal());
        }

        const btnCancel = document.getElementById('btn-cancel-modal');
        if (btnCancel) {
            btnCancel.addEventListener('click', () => this.#view.closeModal());
        }

        const form = document.getElementById('player-form');
        if (form) {
            form.addEventListener('submit', async (event) => {
                event.preventDefault();
                const formData = new FormData(event.target);
                const data = Object.fromEntries(formData.entries());

                try {
                    if (data.id === '') {
                        await this.#model.create(data);
                    } else {
                        await this.#model.update(data);
                    }
                    this.#view.closeModal();
                    await this.#loadPlayers();
                } catch (error) {
                    console.error(error);
                    alert('Error al guardar el jugador');
                }
            });
        }

        document.addEventListener('editPlayer', (e) => {
            this.#view.openModal(e.detail);
        });

        document.addEventListener('deletePlayer', async (e) => {
            try {
                await this.#model.delete(e.detail);
                await this.#loadPlayers();
            } catch (error) {
                console.error(error);
                alert('Error al eliminar el jugador');
            }
        });
    }
}

// Inicialización
document.addEventListener('DOMContentLoaded', () => {
    const isAdmin = window.appConfig?.isAdmin || false;
    const model = new PlayerModel();
    const view = new PlayerView(isAdmin);
    new PlayerController(model, view, isAdmin);
});
