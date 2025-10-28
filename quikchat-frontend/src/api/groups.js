import { getAccessToken } from '../store/authStore.js';

const API_BASE_URL = '/api/v1';

async function handleResponse(response) {
    if (!response.ok) {
        const error = await response.json().catch(() => ({}));
        throw new Error(error.message || 'Something went wrong');
    }
    // 204 = No Content → return nothing
    return response.status === 204 ? null : response.json();
}

export async function createGroup(name, tag, description) {
    const response = await fetch(`${API_BASE_URL}/groups`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${getAccessToken()}`
        },
        body: JSON.stringify({ name, tag, description })
    });
    return handleResponse(response);
}

export async function searchGroups(tag) {
    const response = await fetch(`${API_BASE_URL}/groups/search?tag=${encodeURIComponent(tag)}`, {
        headers: { 'Authorization': `Bearer ${getAccessToken()}` }
    });
    return handleResponse(response);
}

export async function getGroupDetails(groupID) {
    const response = await fetch(`${API_BASE_URL}/groups/${groupID}`, {
        headers: { 'Authorization': `Bearer ${getAccessToken()}` }
    });
    return handleResponse(response);
}

export async function joinGroup(groupID) {
    const response = await fetch(`${API_BASE_URL}/groups/${groupID}/join`, {
        method: 'POST',
        headers: { 'Authorization': `Bearer ${getAccessToken()}` }
    });
    return handleResponse(response);
}

export async function leaveGroup(groupID) {
    const response = await fetch(`${API_BASE_URL}/groups/${groupID}/leave`, {
        method: 'POST',
        headers: { 'Authorization': `Bearer ${getAccessToken()}` }
    });
    return handleResponse(response);
}
