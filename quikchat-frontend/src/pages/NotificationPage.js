import { getNotifications, markAsRead } from '../api/notifications.js';
import { NotificationItem } from '../components/NotificationItem.js';
import { notificationStore } from '../store/notificationStore.js';

export function NotificationPage() {
    const page = document.createElement('div');
    page.className = 'h-full flex flex-col bg-gray-800 text-white';

    page.innerHTML = `
        <header class="bg-gray-900 p-4 border-b border-gray-700">
            <h1 class="text-xl font-bold">Notifications</h1>
        </header>
        <main id="notification-list" class="flex-1 overflow-y-auto p-4">
            <p>Loading notifications...</p>
        </main>
    `;

    const notificationList = page.querySelector('#notification-list');

    async function loadNotifications() {
        try {
            const notifications = await getNotifications();
            notificationStore.setNotifications(notifications);
            renderNotifications(notifications);
        } catch (error) {
            notificationList.innerHTML = `<p class="text-red-400">Error loading notifications: ${error.message}</p>`;
        }
    }

    function renderNotifications(notifications) {
        notificationList.innerHTML = '';
        if (notifications.length === 0) {
            notificationList.innerHTML = '<p class="text-gray-400">No notifications yet.</p>';
            return;
        }
        notifications.forEach(notification => {
            const item = NotificationItem({ notification });
            item.addEventListener('click', async () => {
                if (!notification.read) {
                    try {
                        await markAsRead(notification.id);
                        notificationStore.markNotificationAsReadInState(notification.id);
                        item.classList.add('opacity-60');
                        item.querySelector('.bg-blue-500')?.remove();
                    } catch (error) {
                        console.error('Failed to mark notification as read:', error);
                    }
                }
            });
            notificationList.appendChild(item);
        });
    }

    const unsubscribe = notificationStore.subscribe(state => {
        renderNotifications(state.notifications);
    });

    // Cleanup subscription when page is removed
    page.addEventListener('DOMNodeRemoved', () => unsubscribe());

    loadNotifications();

    return page;
}

