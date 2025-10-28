import { getAccessToken } from '../store/authStore.js';

let socket = null;
let messageListeners = [];

function connect() {
    const token = getAccessToken();
    if (!token || (socket && socket.readyState === WebSocket.OPEN)) {
        return;
    }

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const host = window.location.host;
    
    // Pass token as a query parameter for authentication
    socket = new WebSocket(`${protocol}//${host}/api/v1/ws?auth=${token}`);

    socket.onopen = () => {
        console.log('WebSocket connection established');
    };

    socket.onmessage = (event) => {
        const message = JSON.parse(event.data);
        messageListeners.forEach(listener => listener(message));
    };

    socket.onclose = () => {
        console.log('WebSocket connection closed. Attempting to reconnect...');
        setTimeout(connect, 5000); // Reconnect after 5 seconds
    };

    socket.onerror = (error) => {
        console.error('WebSocket error:', error);
        socket.close();
    };
}

function disconnect() {
    if (socket) {
        socket.onclose = null; // Prevent reconnection logic from firing
        socket.close();
        socket = null;
        console.log('WebSocket connection disconnected');
    }
}

function addMessageListener(callback) {
    messageListeners.push(callback);
}

function removeMessageListener(callback) {
    messageListeners = messageListeners.filter(listener => listener !== callback);
}

export const WebSocketService = {
    connect,
    disconnect,
    addMessageListener,
    removeMessageListener,
};

