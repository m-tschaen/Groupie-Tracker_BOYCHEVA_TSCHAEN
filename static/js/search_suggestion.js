document.addEventListener('DOMContentLoaded', function () {
    const searchInput = document.querySelector('.search-input');
    if (!searchInput) return;

    const searchForm = searchInput.closest('.search-form');
    if (!searchForm) return;

    // On accroche la box au parent (pas dans la form) pour éviter overflow:hidden
    const host = searchForm.parentElement || document.body;
    if (getComputedStyle(host).position === 'static') host.style.position = 'relative';

    let suggestionsBox = host.querySelector(':scope > .suggestions-box');
    if (!suggestionsBox) {
        suggestionsBox = document.createElement('div');
        suggestionsBox.className = 'suggestions-box';
        host.appendChild(suggestionsBox);
    }

    // Positionner la box juste sous la barre
    function positionBox() {
        const formRect = searchForm.getBoundingClientRect();
        const hostRect = host.getBoundingClientRect();

        suggestionsBox.style.position = 'absolute';
        suggestionsBox.style.left = (formRect.left - hostRect.left) + 'px';
        suggestionsBox.style.top = (formRect.top - hostRect.top + formRect.height + 8) + 'px';
        suggestionsBox.style.width = formRect.width + 'px';
        suggestionsBox.style.zIndex = '9999';
    }

    window.addEventListener('resize', positionBox);
    positionBox();

    let allArtists = [];
    let currentFocus = -1;

    const artistCards = document.querySelectorAll('.artist-card');
    artistCards.forEach(card => {
        const name = card.querySelector('h3')?.textContent?.trim() || '';
        const ps = card.querySelectorAll('p');
        const members = ps[0]?.textContent?.trim() || '';
        const year = ps[1]?.textContent?.trim() || '';
        const id = card.getAttribute('href')?.split('/').pop() || '';

        if (name) allArtists.push({ name, members, year, id });
    });

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
            currentFocus = -1;
            return;
        }

        suggestionsBox.innerHTML = matches.map((artist, index) => `
            <div class="suggestion-item" data-index="${index}" data-name="${artist.name}">
                <div class="suggestion-name">${highlightMatch(artist.name, query)}</div>
                <div class="suggestion-info">${artist.members} • ${artist.year}</div>
            </div>
        `).join('');

        suggestionsBox.style.display = 'block';
        currentFocus = -1;

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
    }

    searchInput.addEventListener('input', e => showSuggestions(e.target.value));

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
            currentFocus = -1;
        }
    });

    document.addEventListener('click', function (e) {
        if (!searchForm.contains(e.target) && !suggestionsBox.contains(e.target)) {
            suggestionsBox.style.display = 'none';
            currentFocus = -1;
        }
    });

    searchInput.addEventListener('focus', function () {
        if (this.value.trim().length >= 1) showSuggestions(this.value);
    });
});
