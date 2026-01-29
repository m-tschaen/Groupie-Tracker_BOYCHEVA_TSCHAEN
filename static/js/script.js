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
        
        function scrollToLocation(locationId) {
            event.stopPropagation();
            
            const locationElement = document.getElementById(locationId);
            if (locationElement) {
                const continentContent = locationElement.closest('.continent-content');
                if (continentContent && continentContent.style.display === 'none') {
                    const continentTitle = continentContent.previousElementSibling;
                    if (continentTitle) {
                        continentTitle.click();
                    }
                }
                
                if (locationElement.style.display === 'none' || locationElement.style.display === '') {
                    locationElement.parentElement.click();
                }
                
                locationElement.parentElement.scrollIntoView({ 
                    behavior: 'smooth', 
                    block: 'center' 
                });
                
                locationElement.parentElement.style.background = 'rgba(143, 188, 143, 0.3)';
                setTimeout(() => {
                    locationElement.parentElement.style.background = '';
                }, 2000);
            }
        }
        
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
            
            const newWidth = 100 * zoomLevel;
            const newHeight = (600 / 1200 * 100) * zoomLevel; 
            
            mapImage.style.width = newWidth + '%';
            mapSvg.style.width = newWidth + '%';
        }
        
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
        
        function updateCursor() {
            if (zoomLevel > 1) {
                mapContainer.style.cursor = 'grab';
            } else {
                mapContainer.style.cursor = 'default';
            }
        }
    
        const originalUpdateZoom = updateZoom;
        updateZoom = function() {
            originalUpdateZoom();
            updateCursor();
        };