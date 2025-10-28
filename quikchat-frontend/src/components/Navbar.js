import { notificationStore } from '../store/notificationStore.js';
import { getUnreadCount } from '../api/notifications.js';

export function Navbar() {
    const nav = document.createElement('nav');
    nav.className = 'bg-gray-900 text-white p-4 flex justify-between items-center';

    nav.innerHTML = `
        <a href="#" class="text-2xl font-bold">QuikChat</a>
        <div class="flex items-center space-x-4">
            <a href="#/chat" class="hover:text-gray-300">Chat</a>
            <a href="#/friends" class="hover:text-gray-300">Friends</a>
            <a href="#/notifications" id="notifications-link" class="relative hover:text-gray-300">
                <span>Notifications</span>
                <span id="notification-badge" class="absolute -top-1 -right-2 bg-red-600 text-white text-xs rounded-full h-4 w-4 flex items-center justify-center hidden"></span>
            </a>
            <a href="#/profile" class="hover:text-gray-300">Profile</a>
            <button id="logout-btn" class="bg-red-600 px-3 py-1 rounded hover:bg-red-700">Logout</button>
        </div>
    `;

    const logoutBtn = nav.querySelector('#logout-btn');
    logoutBtn.addEventListener('click', () => {
        // Dispatch a custom event that the router can listen for
        window.dispatchEvent(new CustomEvent('logout'));
    });

    const badge = nav.querySelector('#notification-badge');

    function updateBadge(count) {
        if (count > 0) {
            badge.textContent = count > 9 ? '9+' : count;
            badge.classList.remove('hidden');
        } else {
            badge.classList.add('hidden');
        }
    }

    const unsubscribe = notificationStore.subscribe(state => {
        updateBadge(state.unreadCount);
    });

    // Cleanup subscription when navbar is removed (though it's unlikely)
    nav.addEventListener('DOMNodeRemoved', () => unsubscribe());

    // Initial fetch
    getUnreadCount().then(data => {
        notificationStore.setUnreadCount(data.count);
    }).catch(err => console.error("Failed to fetch unread count", err));

    return nav;
}
