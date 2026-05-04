// MVC Structure for Fichajes CRUD

class TransferModel {
    #apiUrl = '/api/transfers';

    async getAll() {
        const response = await fetch(this.#apiUrl);
        if (!response.ok) throw new Error('Error al obtener fichajes');
        return await response.json();
    }

    async create(data) {
        const response = await fetch(this.#apiUrl, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });
        if (!response.ok) throw new Error('Error al crear fichaje');
        return await response.json();
    }

    async update(data) {
        const response = await fetch(this.#apiUrl, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });
        if (!response.ok) throw new Error('Error al actualizar fichaje');
        return await response.json();
    }

    async delete(id) {
        const response = await fetch(`${this.#apiUrl}?id=${id}`, {
            method: 'DELETE'
        });
        if (!response.ok) throw new Error('Error al eliminar fichaje');
        return await response.json();
    }
}

class TransferView {
    #container;
    #template;
    #modal;
    #form;
    #modalTitle;
    #isAdmin;

    constructor(isAdmin) {
        this.#isAdmin = isAdmin;
        this.#container = document.getElementById('transfers-container');
        this.#template = document.getElementById('transfer-template');
        
        if (this.#isAdmin) {
            this.#modal = document.getElementById('transfer-modal');
            this.#form = document.getElementById('transfer-form');
            this.#modalTitle = document.getElementById('modal-title');
        }
    }

    renderList(transfers) {
        this.#container.textContent = ''; // Limpiar contenedor de forma segura (evita innerHTML)
        
        if (!transfers || transfers.length === 0) {
            const p = document.createElement('p');
            p.textContent = 'No hay fichajes ni rumores activos.';
            this.#container.appendChild(p);
            return;
        }

        transfers.forEach(transfer => {
            const clone = this.#template.content.cloneNode(true);
            
            // Asignar ID al card padre para manipular luego si es necesario
            const card = clone.querySelector('.transfer-card');
            if(card) card.dataset.id = transfer.ID || transfer.id;

            const img = clone.querySelector('.transfer-card__img');
            if(img) {
                img.src = transfer.ImageURL || transfer.image_url || '';
                img.alt = transfer.PlayerName || transfer.player_name || '';
            }

            const name = clone.querySelector('.transfer-card__name');
            if(name) name.textContent = transfer.PlayerName || transfer.player_name || 'Desconocido';

            const route = clone.querySelector('.transfer-card__route');
            if(route) route.textContent = `${transfer.FromTeam || transfer.from_team || ''} → ${transfer.ToTeam || transfer.to_team || ''}`;

            const desc = clone.querySelector('.transfer-card__desc');
            if(desc) desc.textContent = transfer.Description || transfer.description || '';

            const status = clone.querySelector('.transfer-card__status');
            const statusVal = transfer.Status || transfer.status || '';
            if(status) {
                status.textContent = statusVal;
                if (statusVal === 'Hecho') {
                    status.classList.add('transfer-card__status_done');
                } else {
                    status.classList.add('transfer-card__status_rumor');
                }
            }

            if (this.#isAdmin) {
                const btnEdit = clone.querySelector('.btn-edit');
                if(btnEdit) {
                    btnEdit.addEventListener('click', () => {
                        const evt = new CustomEvent('editTransfer', { detail: transfer });
                        document.dispatchEvent(evt);
                    });
                }

                const btnDelete = clone.querySelector('.btn-delete');
                if(btnDelete) {
                    btnDelete.addEventListener('click', () => {
                        if (confirm(`¿Seguro que quieres borrar el fichaje de ${transfer.PlayerName || transfer.player_name}?`)) {
                            const evt = new CustomEvent('deleteTransfer', { detail: transfer.ID || transfer.id });
                            document.dispatchEvent(evt);
                        }
                    });
                }
            }

            this.#container.appendChild(clone);
        });
    }

    showError(msg) {
        this.#container.textContent = '';
        const p = document.createElement('p');
        p.textContent = msg;
        p.classList.add('error-message'); // Asume que se estilará, en vez de usar style
        this.#container.appendChild(p);
    }

    showLoading() {
        this.#container.textContent = '';
        const p = document.createElement('p');
        p.textContent = 'Cargando fichajes...';
        this.#container.appendChild(p);
    }

    // Modal control
    openModal(transfer = null) {
        if (!this.#isAdmin) return;
        
        this.#modal.classList.remove('hidden');
        if (transfer) {
            this.#modalTitle.textContent = 'Editar Fichaje';
            this.#form.elements['id'].value = transfer.ID || transfer.id;
            this.#form.elements['player_name'].value = transfer.PlayerName || transfer.player_name;
            this.#form.elements['from_team'].value = transfer.FromTeam || transfer.from_team;
            this.#form.elements['to_team'].value = transfer.ToTeam || transfer.to_team;
            this.#form.elements['status'].value = transfer.Status || transfer.status;
            this.#form.elements['description'].value = transfer.Description || transfer.description;
            this.#form.elements['image_url'].value = transfer.ImageURL || transfer.image_url;
        } else {
            this.#modalTitle.textContent = 'Añadir Fichaje';
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

class TransferController {
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
        await this.#loadTransfers();
    }

    async #loadTransfers() {
        try {
            this.#view.showLoading();
            const transfers = await this.#model.getAll();
            this.#view.renderList(transfers);
        } catch (error) {
            console.error(error);
            this.#view.showError('Ocurrió un error al cargar los fichajes.');
        }
    }

    #setupAdminEvents() {
        const btnCreate = document.getElementById('btn-create-transfer');
        if (btnCreate) {
            btnCreate.addEventListener('click', () => this.#view.openModal());
        }

        const btnCancel = document.getElementById('btn-cancel-modal');
        if (btnCancel) {
            btnCancel.addEventListener('click', () => this.#view.closeModal());
        }

        const form = document.getElementById('transfer-form');
        if (form) {
            form.addEventListener('submit', async (event) => {
                event.preventDefault();
                const formData = new FormData(event.target);
                const data = Object.fromEntries(formData.entries());

                try {
                    if (data.id === '') {
                        // Create
                        await this.#model.create(data);
                    } else {
                        // Update
                        await this.#model.update(data);
                    }
                    this.#view.closeModal();
                    await this.#loadTransfers();
                } catch (error) {
                    console.error(error);
                    alert('Error al guardar el fichaje');
                }
            });
        }

        // Custom events triggered by View
        document.addEventListener('editTransfer', (e) => {
            this.#view.openModal(e.detail);
        });

        document.addEventListener('deleteTransfer', async (e) => {
            try {
                await this.#model.delete(e.detail);
                await this.#loadTransfers();
            } catch (error) {
                console.error(error);
                alert('Error al eliminar el fichaje');
            }
        });
    }
}

// Inicialización de la app
document.addEventListener('DOMContentLoaded', () => {
    const isAdmin = window.appConfig?.isAdmin || false;
    const model = new TransferModel();
    const view = new TransferView(isAdmin);
    new TransferController(model, view, isAdmin);
});
