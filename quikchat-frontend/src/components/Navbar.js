import { getUnreadCount } from '../api/notifications.js';
import { notificationStore } from '../store/notificationStore.js';

export function Navbar() {
    const nav = document.createElement('nav');
    nav.className = 'bg-gray-800 text-white p-4 flex justify-between items-center flex-wrap gap-2';

    const title = document.createElement('h1');
    title.className = 'text-xl font-bold';
    title.textContent = 'QuikChat';

    const links = document.createElement('div');
    links.className = 'flex items-center gap-4';

    const chatLink = document.createElement('a');
    chatLink.href = '#/';
    chatLink.textContent = 'Chat';
    chatLink.className = 'hover:text-gray-300';

    const friendsLink = document.createElement('a');
    friendsLink.href = '#/friends';
    friendsLink.textContent = 'Friends';
    friendsLink.className = 'hover:text-gray-300';

    const notificationsLink = document.createElement('a');
    notificationsLink.href = '#/notifications';
    notificationsLink.textContent = 'Notifications';
    notificationsLink.className = 'hover:text-gray-300 relative';

    const notificationBadge = document.createElement('span');
    notificationBadge.id = 'notification-badge';
    notificationBadge.className = 'absolute -top-2 -right-2 bg-red-500 text-white text-xs rounded-full h-4 w-4 flex items-center justify-center hidden';
    notificationsLink.appendChild(notificationBadge);

    const profileLink = document.createElement('a');
    profileLink.href = '#/profile';
    profileLink.textContent = 'Profile';
    profileLink.className = 'hover:text-gray-300';

    const logoutButton = document.createElement('button');
    logoutButton.textContent = 'Logout';
    logoutButton.className = 'bg-red-500 hover:bg-red-600 px-3 py-1 rounded';
    logoutButton.onclick = () => {
        document.dispatchEvent(new CustomEvent('logout'));
    };

    links.append(chatLink, friendsLink, notificationsLink, profileLink, logoutButton);
    nav.append(title, links);

    function updateBadge(count) {
        if (count > 0) {
            notificationBadge.textContent = count > 9 ? '9+' : count;
            notificationBadge.classList.remove('hidden');
        } else {
            notificationBadge.classList.add('hidden');
        }
    }

    // Initial fetch and subscribe
    getUnreadCount().then(data => notificationStore.setUnreadCount(data.unread_count));
    const unsubscribe = notificationStore.subscribe(state => updateBadge(state.unreadCount));

    // Cleanup subscription when navbar is removed (though it's unlikely)
    const observer = new MutationObserver((mutations, obs) => {
        if (!document.body.contains(nav)) {
            unsubscribe();
            obs.disconnect();
        }
    });
    observer.observe(document.body, { childList: true, subtree: true });


    return nav;
}

