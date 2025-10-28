import { getAccessToken } from '../store/authStore.js';
import { chatStore } from '../store/chatStore.js';
import { notificationStore } from '../store/notificationStore.js';

let socket = null;
let reconnectInterval = 5000;
let shouldReconnect = false;

function connect() {
    const token = getAccessToken();
    if (!token) {
        console.log("WebSocket: No access token found.");
        return;
    }

    shouldReconnect = true;
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const url = `${protocol}//${window.location.host}/api/v1/ws?auth=${token}`;

    socket = new WebSocket(url);

    socket.onopen = () => {
        console.log("WebSocket connection established.");
    };

    socket.onmessage = (event) => {
        try {
            const data = JSON.parse(event.data);
            if (data.type === 'new_message') {
                const message = data.payload;
                chatStore.addMessage(message.conversation_id, message);
            } else if (data.type === 'new_notification') {
                const notification = data.payload;
                notificationStore.addNotification(notification);
            }
        } catch (error) {
            console.error("Error parsing WebSocket message:", error);
        }
    };

    socket.onclose = () => {
        console.log("WebSocket connection closed.");
        if (shouldReconnect) {
            setTimeout(connect, reconnectInterval);
        }
    };

    socket.onerror = (error) => {
        console.error("WebSocket error:", error);
        socket.close();
    };
}

function disconnect() {
    shouldReconnect = false;
    if (socket) {
        socket.close();
        socket = null;
    }
    console.log("WebSocket disconnected by client.");
}

export const WebSocketService = {
    connect,
    disconnect,
};
