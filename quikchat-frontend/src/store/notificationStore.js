let state = {
    notifications: [],
    unreadCount: 0,
};

const listeners = new Set();

function subscribe(listener) {
    listeners.add(listener);
    return () => listeners.delete(listener);
}

function notify() {
    for (const listener of listeners) {
        listener(state);
    }
}

function setNotifications(notifications) {
    state.notifications = notifications;
    notify();
}

function addNotification(notification) {
    state.notifications = [notification, ...state.notifications];
    state.unreadCount++;
    notify();
}

function setUnreadCount(count) {
    state.unreadCount = count;
    notify();
}

function decrementUnreadCount() {
    if (state.unreadCount > 0) {
        state.unreadCount--;
    }
    notify();
}

function markNotificationAsReadInState(notificationId) {
    const notification = state.notifications.find(n => n.id === notificationId);
    if (notification && !notification.read) {
        notification.read = true;
        decrementUnreadCount();
    }
    notify();
}

function getState() {
    return state;
}

export const notificationStore = {
    subscribe,
    setNotifications,
    addNotification,
    setUnreadCount,
    decrementUnreadCount,
    markNotificationAsReadInState,
    getState,
};

