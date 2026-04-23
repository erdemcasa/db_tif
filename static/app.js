// Fetch and display TIFs on the index page
async function fetchTifs() {
    const tbody = document.getElementById('tifsBody');
    if (!tbody) return;

    try {
        const response = await fetch('/api/tifs');
        if (!response.ok) throw new Error('Erreur réseau');
        
        const data = await response.json();
        
        tbody.innerHTML = ''; // clear loading

        if (data.length === 0) {
            tbody.innerHTML = '<tr><td colspan="4" style="text-align: center; color: var(--text-muted);">Aucune donnée pour le moment</td></tr>';
            return;
        }

        data.forEach(tif => {
            const tr = document.createElement('tr');
            
            // Format date
            const dateObj = new Date(tif.created_at);
            const dateStr = dateObj.toLocaleString('fr-FR', { 
                year: 'numeric', 
                month: 'long', 
                day: 'numeric',
                hour: '2-digit',
                minute: '2-digit'
            });

            tr.innerHTML = `
                <td>#${tif.id}</td>
                <td style="font-weight: 500;">${escapeHTML(tif.tif_label)}</td>
                <td>${escapeHTML(tif.author)}</td>
                <td style="color: var(--text-muted);">${dateStr}</td>
            `;
            tbody.appendChild(tr);
        });

    } catch (error) {
        tbody.innerHTML = '<tr><td colspan="4" style="text-align: center; color: var(--error-color);">Erreur lors du chargement des données</td></tr>';
        console.error('Erreur:', error);
    }
}

// Handle form submission on the add page
async function handleFormSubmit(e) {
    e.preventDefault();

    const tifLabelInput = document.getElementById('tif_label');
    const authorInput = document.getElementById('author');
    const errorAlert = document.getElementById('error-message');
    const successAlert = document.getElementById('success-message');
    const submitBtn = document.getElementById('submitBtn');

    const tifLabel = tifLabelInput.value.trim();
    const author = authorInput.value.trim();

    // Hide previous messages
    errorAlert.classList.add('hidden');
    successAlert.classList.add('hidden');

    // Validation
    if (!tifLabel.endsWith('tif')) {
        showError('Le TIF Label doit obligatoirement se terminer par "tif".');
        return;
    }
    if (!author) {
        showError('L\'auteur est obligatoire.');
        return;
    }

    // Disable button during request
    submitBtn.disabled = true;
    submitBtn.textContent = 'Enregistrement...';

    try {
        const response = await fetch('/api/tifs', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                tif_label: tifLabel,
                author: author
            })
        });

        if (response.ok) {
            successAlert.classList.remove('hidden');
            tifLabelInput.value = '';
            authorInput.value = '';
            setTimeout(() => {
                window.location.href = 'index.html';
            }, 1500);
        } else {
            const errText = await response.text();
            showError(`Erreur: ${errText}`);
        }
    } catch (error) {
        showError('Erreur de connexion au serveur.');
        console.error('Erreur:', error);
    } finally {
        submitBtn.disabled = false;
        submitBtn.textContent = 'Enregistrer';
    }

    function showError(msg) {
        errorAlert.textContent = msg;
        errorAlert.classList.remove('hidden');
    }
}

// Utils
function escapeHTML(str) {
    const div = document.createElement('div');
    div.textContent = str;
    return div.innerHTML;
}
