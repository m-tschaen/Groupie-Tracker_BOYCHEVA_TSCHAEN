document.addEventListener('DOMContentLoaded', function () {
    const searchInput = document.querySelector('.search-input');
    const searchForm = document.querySelector('.search-form');
    
    if (!searchInput || !searchForm) {
        console.error("Erreur : L'input ou le formulaire de recherche est introuvable.");
        return;
    }

    // Définition de l'hôte pour la box (le parent du formulaire)
    const host = searchForm.parentElement || document.body;
    if (getComputedStyle(host).position === 'static') {
        host.style.position = 'relative';
    }

    // Création ou récupération de la box de suggestions
    let suggestionsBox = host.querySelector(':scope > .suggestions-box');
    if (!suggestionsBox) {
        suggestionsBox = document.createElement('div');
        suggestionsBox.className = 'suggestions-box';
        host.appendChild(suggestionsBox);
    }

    // Fonction de positionnement précise
    function positionBox() {
        const formRect = searchForm.getBoundingClientRect();
        const hostRect = host.getBoundingClientRect();

        suggestionsBox.style.position = 'absolute';
        // Aligné sur la gauche du formulaire
        suggestionsBox.style.left = (formRect.left - hostRect.left) + 'px';
        // Placé 8px sous le formulaire
        suggestionsBox.style.top = (formRect.top - hostRect.top + formRect.height + 8) + 'px';
        suggestionsBox.style.width = formRect.width + 'px';
        suggestionsBox.style.zIndex = '9999';
    }

    window.addEventListener('resize', positionBox);

    // Extraction des données depuis les cartes HTML
    let allArtists = [];
    const artistCards = document.querySelectorAll('.artist-card');
    
    artistCards.forEach(card => {
        const name = card.querySelector('h3')?.textContent?.trim() || '';
        const ps = card.querySelectorAll('p');
        // ps[0] contient "X membres", ps[1] contient l'année
        const members = ps[0]?.textContent?.trim() || '';
        const year = ps[1]?.textContent?.trim() || '';
        const href = card.getAttribute('href') || '';
        const id = href.split('/').pop() || '';

        if (name) {
            allArtists.push({ name, members, year, id });
        }
    });

    console.log(`${allArtists.length} artistes chargés pour les suggestions.`);

    let currentFocus = -1;

    function highlightMatch(text, query) {
        const safe = query.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
        const regex = new RegExp(`(${safe})`, 'gi');
        return text.replace(regex, '<strong>$1</strong>');
    }

    function showSuggestions(query) {
        positionBox();

        if (!query || query.trim().length < 1) {
            suggestionsBox.innerHTML = '';
            suggestionsBox.style.display = 'none';
            return;
        }

        const q = query.toLowerCase();
        const matches = allArtists
            .filter(a =>
                a.name.toLowerCase().includes(q) ||
                a.members.toLowerCase().includes(q) ||
                a.year.toLowerCase().includes(q)
            )
            .slice(0, 8);

        if (matches.length === 0) {
            suggestionsBox.innerHTML = '<div class="suggestion-item no-result">Aucun résultat</div>';
            suggestionsBox.style.display = 'block';
            return;
        }

        suggestionsBox.innerHTML = matches.map((artist, index) => `
            <div class="suggestion-item" data-index="${index}" data-name="${artist.name}">
                <div class="suggestion-name">${highlightMatch(artist.name, query)}</div>
                <div class="suggestion-info">${artist.members} • ${artist.year}</div>
            </div>
        `).join('');

        suggestionsBox.style.display = 'block';
        
        // Réattacher les événements de clic
        suggestionsBox.querySelectorAll('.suggestion-item:not(.no-result)').forEach(item => {
            item.addEventListener('click', function () {
                searchInput.value = this.getAttribute('data-name') || '';
                suggestionsBox.style.display = 'none';
                searchForm.submit();
            });
        });
    }

    function setActive(items) {
        items.forEach(i => i.classList.remove('active'));
        if (items.length === 0) return;
        if (currentFocus >= items.length) currentFocus = 0;
        if (currentFocus < 0) currentFocus = items.length - 1;
        items[currentFocus]?.classList.add('active');
        items[currentFocus]?.scrollIntoView({ block: 'nearest' });
    }

    // Événements clavier et souris
    searchInput.addEventListener('input', e => {
        currentFocus = -1;
        showSuggestions(e.target.value);
    });

    searchInput.addEventListener('keydown', function (e) {
        const items = suggestionsBox.querySelectorAll('.suggestion-item:not(.no-result)');

        if (e.key === 'ArrowDown') {
            e.preventDefault();
            currentFocus++;
            setActive(items);
        } else if (e.key === 'ArrowUp') {
            e.preventDefault();
            currentFocus--;
            setActive(items);
        } else if (e.key === 'Enter') {
            if (currentFocus > -1 && items[currentFocus]) {
                e.preventDefault();
                items[currentFocus].click();
            }
        } else if (e.key === 'Escape') {
            suggestionsBox.style.display = 'none';
        }
    });

    document.addEventListener('click', function (e) {
        if (!searchForm.contains(e.target) && !suggestionsBox.contains(e.target)) {
            suggestionsBox.style.display = 'none';
        }
    });

    searchInput.addEventListener('focus', function () {
        if (this.value.trim().length >= 1) showSuggestions(this.value);
    });
});