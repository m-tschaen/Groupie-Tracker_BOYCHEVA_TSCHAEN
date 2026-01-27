function toggleContinent(continentId) {
            const content = document.getElementById(continentId);
            const icon = event.currentTarget.querySelector('.toggle-icon');
            
            if (content.style.display === 'none' || content.style.display === '') {
                content.style.display = 'block';
                icon.textContent = '▲';
            } else {
                content.style.display = 'none';
                icon.textContent = '▼';
            }
        }
        
        function toggleDetails(locationId) {
            const details = document.getElementById(locationId);
            
            if (details.style.display === 'none' || details.style.display === '') {
                details.style.display = 'block';
            } else {
                details.style.display = 'none';
            }
        }
        
        // Fonction pour scroller vers une ville quand on clique sur un marqueur
        function scrollToLocation(locationId) {
            event.stopPropagation();
            
            // Trouver l'élément dans la liste
            const locationElement = document.getElementById(locationId);
            if (locationElement) {
                // Ouvrir le continent parent
                const continentContent = locationElement.closest('.continent-content');
                if (continentContent && continentContent.style.display === 'none') {
                    const continentTitle = continentContent.previousElementSibling;
                    if (continentTitle) {
                        continentTitle.click();
                    }
                }
                
                // Ouvrir les détails de la ville
                if (locationElement.style.display === 'none' || locationElement.style.display === '') {
                    locationElement.parentElement.click();
                }
                
                // Scroll vers l'élément
                locationElement.parentElement.scrollIntoView({ 
                    behavior: 'smooth', 
                    block: 'center' 
                });
                
                // Effet de highlight
                locationElement.parentElement.style.background = 'rgba(143, 188, 143, 0.3)';
                setTimeout(() => {
                    locationElement.parentElement.style.background = '';
                }, 2000);
            }
        }
        
        // Système de zoom sur la carte
        let zoomLevel = 1;
        const mapContainer = document.getElementById('map-container');
        const mapImage = document.getElementById('map-image');
        const mapSvg = document.getElementById('map-svg');
        
        function zoomIn() {
            zoomLevel = Math.min(zoomLevel + 0.5, 4);
            updateZoom();
        }
        
        function zoomOut() {
            zoomLevel = Math.max(zoomLevel - 0.5, 1);
            updateZoom();
        }
        
        function resetZoom() {
            zoomLevel = 1;
            updateZoom();
            mapContainer.scrollTop = 0;
            mapContainer.scrollLeft = 0;
        }
        
        function updateZoom() {
            mapImage.style.transform = `scale(${zoomLevel})`;
            mapSvg.style.transform = `scale(${zoomLevel})`;
            
            // Ajuster la taille du conteneur pour permettre le scroll
            const newWidth = 100 * zoomLevel;
            const newHeight = (600 / 1200 * 100) * zoomLevel; // Maintenir le ratio de la carte
            
            mapImage.style.width = newWidth + '%';
            mapSvg.style.width = newWidth + '%';
        }
        
        // Permettre le drag pour déplacer la carte quand zoomée
        let isDragging = false;
        let startX, startY, scrollLeft, scrollTop;
        
        mapContainer.addEventListener('mousedown', (e) => {
            if (zoomLevel > 1) {
                isDragging = true;
                mapContainer.style.cursor = 'grabbing';
                startX = e.pageX - mapContainer.offsetLeft;
                startY = e.pageY - mapContainer.offsetTop;
                scrollLeft = mapContainer.scrollLeft;
                scrollTop = mapContainer.scrollTop;
            }
        });
        
        mapContainer.addEventListener('mouseleave', () => {
            isDragging = false;
            if (zoomLevel > 1) {
                mapContainer.style.cursor = 'grab';
            }
        });
        
        mapContainer.addEventListener('mouseup', () => {
            isDragging = false;
            if (zoomLevel > 1) {
                mapContainer.style.cursor = 'grab';
            }
        });
        
        mapContainer.addEventListener('mousemove', (e) => {
            if (!isDragging) return;
            e.preventDefault();
            const x = e.pageX - mapContainer.offsetLeft;
            const y = e.pageY - mapContainer.offsetTop;
            const walkX = (x - startX) * 2;
            const walkY = (y - startY) * 2;
            mapContainer.scrollLeft = scrollLeft - walkX;
            mapContainer.scrollTop = scrollTop - walkY;
        });
        
        // Mettre à jour le curseur selon le zoom
        function updateCursor() {
            if (zoomLevel > 1) {
                mapContainer.style.cursor = 'grab';
            } else {
                mapContainer.style.cursor = 'default';
            }
        }
        
        // Appeler updateCursor après chaque zoom
        const originalUpdateZoom = updateZoom;
        updateZoom = function() {
            originalUpdateZoom();
            updateCursor();
        };