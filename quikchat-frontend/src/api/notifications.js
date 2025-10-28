import { getAccessToken } from '../store/authStore.js';

const API_BASE_URL = '/api/v1';

async function handleResponse(response) {
    if (!response.ok) {
        const error = await response.json().catch(() => ({ message: 'An unknown error occurred' }));
        throw new Error(error.message || 'API request failed');
    }
    if (response.status === 204) {
        return null;
    }
    return response.json();
}

export async function getNotifications(limit = 20, offset = 0) {
    const response = await fetch(`${API_BASE_URL}/notifications?limit=${limit}&offset=${offset}`, {
        headers: {
            'Authorization': `Bearer ${getAccessToken()}`,
        },
    });
    return handleResponse(response);
}

export async function markAsRead(notificationId) {
    const response = await fetch(`${API_BASE_URL}/notifications/${notificationId}/read`, {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${getAccessToken()}`,
        },
    });
    return handleResponse(response);
}

export async function markAllAsRead() {
    const response = await fetch(`${API_BASH_URL}/notifications/read-all`, {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${getAccessToken()}`,
        },
    });
    return handleResponse(response);
}

export async function getUnreadCount() {
    const response = await fetch(`${API_BASE_URL}/notifications/unread-count`, {
        headers: {
            'Authorization': `Bearer ${getAccessToken()}`,
        },
    });
    return handleResponse(response);
}

