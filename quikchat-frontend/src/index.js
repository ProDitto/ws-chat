import { LoginPage } from './pages/LoginPage.js';
import { ChatPage } from './pages/ChatPage.js';
import { ProfilePage } from './pages/ProfilePage.js';
import { FriendsPage } from './pages/FriendsPage.js';
import { render } from './utils/dom.js';
import { isAuthenticated, logout } from './store/authStore.js';
import { WebSocketService } from './services/WebSocketService.js';

const routes = {
    '/': ChatPage,
    '/login': LoginPage,
    '/profile': ProfilePage,
    '/friends': FriendsPage,
};

const router = () => {
    const path = window.location.hash.slice(1) || '/';

    if (isAuthenticated()) {
        WebSocketService.connect(); // Connect WebSocket if authenticated
        if (path === '/login') {
            window.location.hash = '/';
            return;
        }
        const page = routes[path] || routes['/'];
        render(page());
    } else {
        WebSocketService.disconnect(); // Disconnect if not authenticated
        if (path !== '/login') {
            window.location.hash = '/login';
            return;
        }
        render(LoginPage());
    }
};

window.addEventListener('DOMContentLoaded', router);
window.addEventListener('hashchange', router);

// Custom event for logout
window.addEventListener('logout', () => {
    logout();
    WebSocketService.disconnect(); // Also disconnect on logout event
    window.location.hash = '/login';
});

