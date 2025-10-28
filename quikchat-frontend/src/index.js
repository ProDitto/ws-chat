import { LoginPage } from './pages/LoginPage.js';
import { ChatPage } from './pages/ChatPage.js';
import { ProfilePage } from './pages/ProfilePage.js';
import { FriendsPage } from './pages/FriendsPage.js';
import { NotificationPage } from './pages/NotificationPage.js';
import { render } from './utils/dom.js';
import { isAuthenticated, logout } from './store/authStore.js';
import { WebSocketService } from './services/WebSocketService.js';
import { Navbar } from './components/Navbar.js';

const routes = {
    '/': ChatPage,
    '/login': LoginPage,
    '/profile': ProfilePage,
    '/friends': FriendsPage,
    '/notifications': NotificationPage,
};

const protectedRoutes = ['/', '/profile', '/friends', '/notifications'];

function router() {
    const path = window.location.hash.slice(1) || '/';
    const appElement = document.getElementById('app');
    appElement.innerHTML = ''; // Clear the app container

    const isAuth = isAuthenticated();
    const isProtectedRoute = protectedRoutes.includes(path);

    if (isProtectedRoute && !isAuth) {
        window.location.hash = '/login';
        return;
    }

    if (path === '/login' && isAuth) {
        window.location.hash = '/';
        return;
    }

    if (isAuth) {
        appElement.appendChild(Navbar());
        WebSocketService.connect();
    } else {
        WebSocketService.disconnect();
    }

    const pageContainer = document.createElement('div');
    pageContainer.className = 'page-container'; // Add a class for styling if needed
    const page = routes[path] || routes['/']; // Fallback to a default page
    pageContainer.appendChild(page());
    appElement.appendChild(pageContainer);
}

window.addEventListener('hashchange', router);
window.addEventListener('DOMContentLoaded', router);

// Custom event listener for logout
window.addEventListener('logout', () => {
    logout();
    WebSocketService.disconnect();
    window.location.hash = '/login';
});
