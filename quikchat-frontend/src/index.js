import { render } from './utils/dom.js';
import { isAuthenticated, logout } from './store/authStore.js';
import LoginPage from './pages/LoginPage.js';
import ChatPage from './pages/ChatPage.js';

const routes = {
    '/': ChatPage,
    '/login': LoginPage,
};

function router() {
    const path = window.location.hash.slice(1) || '/';
    
    if (!isAuthenticated() && path !== '/login') {
        window.location.hash = '/login';
        return;
    }

    if (isAuthenticated() && path === '/login') {
        window.location.hash = '/';
        return;
    }

    const page = routes[path] || routes['/']; // Default to chat page if authenticated
    render(page());
}

// Handle initial page load
document.addEventListener('DOMContentLoaded', router);

// Handle hash changes
window.addEventListener('hashchange', router);

// Handle logout globally
document.addEventListener('logout', () => {
    logout();
    window.location.hash = '/login';
});
```
