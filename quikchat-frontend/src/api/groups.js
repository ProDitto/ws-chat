import { getAccessToken } from '../store/authStore.js';

const API_BASE_URL = '/api/v1';

async function handleResponse(response) {
    if (!response.ok) {
        const error = await response.json();
        throw new Error(error.message || 'Something went wrong');
    }
    return response.json();
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
    const response = await fetch(`${API_...
```
