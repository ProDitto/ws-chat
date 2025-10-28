/**
 * @param {{ notification: object }} props
 */
export function NotificationItem({ notification }) {
    const item = document.createElement('div');
    item.className = `p-3 border-b border-gray-700 flex items-start gap-3 ${notification.read ? 'opacity-60' : ''}`;

    const actorImageUrl = notification.actor?.profile_image_url || 'https://via.placeholder.com/40';

    item.innerHTML = `
        <img src="${actorImageUrl}" alt="${notification.actor?.username || 'System'}" class="w-10 h-10 rounded-full bg-gray-600">
        <div class="flex-1">
            <p class="text-gray-200">${notification.message}</p>
            <p class="text-xs text-gray-400 mt-1">${new Date(notification.created_at).toLocaleString()}</p>
        </div>
        ${!notification.read ? '<div class="w-2 h-2 bg-blue-500 rounded-full self-center"></div>' : ''}
    `;

    return item;
}

