import { getAccessToken } from '../store/authStore.js';
import { handleResponse } from './auth.js';

const API_BASE_URL = '/api/v1';

export async function getProfile(username) {
    const response = await fetch(`${API_BASE_URL}/users/${username}`, {
        headers: {
            'Authorization': `Bearer ${getAccessToken()}`,
        },
    });
    return handleResponse(response);
}

export async function updateProfile(profileData) {
    const response = await fetch(`${API_BASE_URL}/users/me/profile`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${getAccessToken()}`,
        },
        body: JSON.stringify(profileData),
    });
    return handleResponse(response);
}

export async function blockUser(username) {
    const response = await fetch(`${API_BAPI_BASE_URLASE_URL}/users/${username}/block`, {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${getAccessToken()}`,
        },
    });
    return handleResponse(response);
}

export async function unblockUser(username) {
    const response = await fetch(`${API_BASE_URL}/users/${username}/unblock`, {
        method: 'DELETE',
        headers: {
            'Authorization': `Bearer ${getAccessToken()}`,
        },
    });
    return handleResponse(response);
}

