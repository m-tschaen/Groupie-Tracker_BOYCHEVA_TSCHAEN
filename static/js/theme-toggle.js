document.addEventListener('DOMContentLoaded', function() {
    const toggleBtn = document.createElement('button');
    toggleBtn.className = 'theme-toggle';
    toggleBtn.setAttribute('aria-label', 'Changer de thème');
    toggleBtn.innerHTML = `
        <svg class="moon-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor">
            <path d="M21.64,13a1,1,0,0,0-1.05-.14,8.05,8.05,0,0,1-3.37.73A8.15,8.15,0,0,1,9.08,5.49a8.59,8.59,0,0,1,.25-2A1,1,0,0,0,8,2.36,10.14,10.14,0,1,0,22,14.05,1,1,0,0,0,21.64,13Zm-9.5,6.69A8.14,8.14,0,0,1,7.08,5.22v.27A10.15,10.15,0,0,0,17.22,15.63a9.79,9.79,0,0,0,2.1-.22A8.11,8.11,0,0,1,12.14,19.73Z"/>
        </svg>
        <svg class="sun-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor">
            <path d="M12,18a6,6,0,1,1,6-6A6,6,0,0,1,12,18Zm0-10a4,4,0,1,0,4,4A4,4,0,0,0,12,8Z"/>
            <path d="M12,4a1,1,0,0,0,1-1V1a1,1,0,0,0-2,0V3A1,1,0,0,0,12,4Z"/>
            <path d="M21,11H19a1,1,0,0,0,0,2h2a1,1,0,0,0,0-2Z"/>
            <path d="M5,11H3a1,1,0,0,0,0,2H5a1,1,0,0,0,0-2Z"/>
            <path d="M12,20a1,1,0,0,0-1,1v2a1,1,0,0,0,2,0V21A1,1,0,0,0,12,20Z"/>
            <path d="M16.89,6.34a1,1,0,0,0,.7-.29l1.42-1.42a1,1,0,1,0-1.41-1.41L16.18,4.64a1,1,0,0,0,.71,1.7Z"/>
            <path d="M5.64,17.05a1,1,0,0,0-.71.29L3.51,18.75a1,1,0,0,0,1.41,1.41l1.42-1.41a1,1,0,0,0-.7-1.7Z"/>
            <path d="M4.93,6.34A1,1,0,0,0,5.64,5.05L4.22,3.63A1,1,0,0,0,2.81,5.05L4.22,6.46A1,1,0,0,0,4.93,6.34Z"/>
            <path d="M18.36,17.05a1,1,0,0,0-.7,1.7l1.41,1.41a1,1,0,0,0,1.41-1.41l-1.41-1.41A1,1,0,0,0,18.36,17.05Z"/>
        </svg>
    `;
    
    document.body.appendChild(toggleBtn);
    
    const savedTheme = localStorage.getItem('theme');
    if (savedTheme === 'dark') {
        document.body.classList.add('dark-theme');
    }
    
    toggleBtn.addEventListener('click', function() {
        document.body.classList.toggle('dark-theme');
        
        if (document.body.classList.contains('dark-theme')) {
            localStorage.setItem('theme', 'dark');
        } else {
            localStorage.setItem('theme', 'light');
        }
    });
});