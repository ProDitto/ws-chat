import { isAuthenticated, logout } from './store/authStore.js';
import { render } from './utils/dom.js';
import { LoginPage } from './pages/LoginPage.js';
import { ChatPage } from './pages/ChatPage.js';
import { ProfilePage } from './pages/ProfilePage.js';
import { FriendsPage } from './pages/FriendsPage.js';

const routes = {
    '/': ChatPage,
    '/login': LoginPage,
    '/profile': ProfilePage,
    '/friends': FriendsPage,
};

function router() {
    const path = window.location.hash.slice(1) || '/';
    const page = routes[path] || ChatPage; // Default to ChatPage if route not found

    if (path !== '/login' && !isAuthenticated()) {
        window.location.hash = '/login';
        return;
    }

    if (path === '/login' && isAuthenticated()) {
        window.location.hash = '/';
        return;
    }

    render(page);
}

document.addEventListener('DOMContentLoaded', router);
window.addEventListener('hashchange', router);

document.addEventListener('logout', () => {
    logout();
    window.location.hash = '/login';
});
