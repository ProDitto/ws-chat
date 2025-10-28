import { getAccessToken } from '../store/authStore.js';
import { handleResponse } from './auth.js';

const API_BASE_URL = '/api/v1';

export async function fetchMessages(conversationId, cursor = null) {
    const url = new URL(`${API_BASE_URL}/conversations/${conversationId}/messages`);
    if (cursor) {
        url.searchParams.append('cursor', cursor);
    }
    url.searchParams.append('limit', 50);

    const response = await fetch(url, {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${getAccessToken()}`,
            'Content-Type': 'application/json',
        },
    });
    return handleResponse(response);
}

export async function sendMessage(conversationId, message) {
    const response = await fetch(`${API_BASE_URL}/conversations/${conversationId}/messages`, {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${getAccessToken()}`,
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(message),
    });
    return handleResponse(response);
}

export async function getPresignedUrl(filename) {
    const response = await fetch(`${API_BASE_URL}/media/presigned-url?filename=${encodeURIComponent(filename)}`, {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${getAccessToken()}`,
        },
    });
    return handleResponse(response);
}

export async function uploadFile(uploadUrl, file) {
    const response = await fetch(uploadUrl, {
        method: 'PUT',
        headers: {
            'Content-Type': file.type,
        },
        body: file,
    });
    if (!response.ok) {
        throw new Error('File upload failed');
    }
    // S3 PUT requests with a presigned URL return a 200 OK with an empty body on success.
    return true;
}

