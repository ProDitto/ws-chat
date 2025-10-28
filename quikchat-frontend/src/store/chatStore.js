// A simple, in-memory store for chat state.
const state = {
    conversations: [],
    messages: {}, // { conversationId: [message1, message2] }
    activeConversationId: null,
};

const listeners = [];

function subscribe(listener) {
    listeners.push(listener);
    return function unsubscribe() {
        const index = listeners.indexOf(listener);
        if (index > -1) {
            listeners.splice(index, 1);
        }
    };
}

function notify() {
    listeners.forEach(listener => listener());
}

function addMessage(conversationId, message) {
    if (!state.messages[conversationId]) {
        state.messages[conversationId] = [];
    }
    // Avoid adding duplicate messages
    if (!state.messages[conversationId].some(m => m.id === message.id)) {
        state.messages[conversationId].push(message);
        notify();
    }
}

function prependMessages(conversationId, messages) {
    if (!state.messages[conversationId]) {
        state.messages[conversationId] = [];
    }
    state.messages[conversationId] = [...messages, ...state.messages[conversationId]];
    notify();
}

function setActiveConversation(conversationId) {
    state.activeConversationId = conversationId;
    if (!state.messages[conversationId]) {
        state.messages[conversationId] = [];
    }
    notify();
}

function getState() {
    return state;
}

export const chatStore = {
    subscribe,
    addMessage,
    prependMessages,
    setActiveConversation,
    getState,
};

