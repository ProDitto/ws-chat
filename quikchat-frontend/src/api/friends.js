import { getAccessToken } from '../store/authStore.js';
import { handleResponse } from './auth.js';

const API_BASE_URL = '/api/v1';

export async function listFriends() {
    const response = await fetch(`${API_BASE_URL}/friends`, {
        headers: { 'Authorization': `Bearer ${getAccessToken()}` },
    });
    return handleResponse(response);
}

export async function listIncomingRequests() {
    const response = await fetch(`${API_BASE_URL}/friends/requests`, {
        headers: { 'Authorization': `Bearer ${getAccessToken()}` },
    });
    return handleResponse(response);
}

export async function sendFriendRequest(username) {
    const response = await fetch(`${API_BASE_URL}/friends/requests`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${getAccessToken()}`,
        },
        body: JSON.stringify({ username }),
    });
    return handleResponse(response);
}

export async function respondToRequest(requestID, action) {
    const response = await fetch(`${API_BASE_URL}/friends/requests/${requestID}`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${getAccessToken()}`,
        },
        body: JSON.stringify({ action }),
    });
    return handleResponse(response);
}

export async function unfriend(username) {
    const response = await fetch(`${API_BASE_URL}/friends/${username}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${getAccessToken()}` },
    });
    return handleResponse(response);
}

