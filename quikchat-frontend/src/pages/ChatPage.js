import { fetchMessages, sendMessage, getPresignedUrl, uploadFile } from '../api/messages.js';
import { MessageDisplay } from '../components/MessageDisplay.js';
import { chatStore } from '../store/chatStore.js';
import { WebSocketService } from '../services/WebSocketService.js';

// For demonstration purposes
const CURRENT_USER_ID = parseInt(localStorage.getItem('user_id'), 10);
const CONVERSATION_ID = 1; // Hardcoded for now

export function ChatPage() {
    const container = document.createElement('div');
    container.className = 'h-full w-full flex flex-col md:flex-row overflow-hidden';

    // --- Sidebar (Conversations List) ---
    const sidebar = document.createElement('div');
    sidebar.className = 'w-full md:w-1/4 bg-gray-800 text-white p-4 flex-col hidden md:flex'; // Hidden on mobile
    sidebar.innerHTML = `
        <h2 class="text-lg font-bold mb-4">Conversations</h2>
        <ul>
            <li class="p-2 rounded bg-gray-700 cursor-pointer">General Chat</li>
            <!-- More conversations would be listed here -->
        </ul>
    `;

    // --- Main Chat Area ---
    const chatArea = document.createElement('div');
    chatArea.className = 'w-full md:flex-1 flex flex-col bg-gray-700';

    // --- Chat Header ---
    const chatHeader = document.createElement('div');
    chatHeader.className = 'bg-gray-800 text-white p-4 font-bold';
    chatHeader.textContent = 'General Chat';

    // --- Messages Container ---
    const messagesContainer = document.createElement('div');
    messagesContainer.id = 'messages-container';
    messagesContainer.className = 'flex-1 p-4 overflow-y-auto flex flex-col-reverse'; // Reverse for new messages at bottom
    const messagesList = document.createElement('div');
    messagesContainer.appendChild(messagesList);

    // --- Message Input Form ---
    const form = document.createElement('form');
    form.className = 'p-4 bg-gray-800 flex items-center gap-2';
    const input = document.createElement('input');
    input.type = 'text';
    input.placeholder = 'Type a message...';
    input.className = 'flex-1 p-2 rounded bg-gray-600 text-white border border-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500';
    const sendButton = document.createElement('button');
    sendButton.type = 'submit';
    sendButton.textContent = 'Send';
    sendButton.className = 'bg-blue-500 hover:bg-blue-600 text-white p-2 rounded';

    const fileInput = document.createElement('input');
    fileInput.type = 'file';
    fileInput.id = 'file-upload';
    fileInput.accept = 'image/jpeg,image/png,image/gif,video/mp4';
    fileInput.className = 'hidden';

    const uploadButton = document.createElement('button');
    uploadButton.type = 'button';
    uploadButton.textContent = '📎';
    uploadButton.className = 'bg-gray-600 hover:bg-gray-500 text-white p-2 rounded';
    uploadButton.onclick = () => fileInput.click();

    form.append(uploadButton, fileInput, input, sendButton);
    chatArea.append(chatHeader, messagesContainer, form);
    container.append(sidebar, chatArea);

    let isLoading = false;
    let hasMoreMessages = true;
    let lastMessageId = null;

    const renderMessages = (messages, prepend = false) => {
        if (messages.length === 0 && prepend) {
            hasMoreMessages = false;
            return;
        }
        messages.forEach(msg => {
            const messageEl = MessageDisplay({
                message: msg,
                isOwnMessage: msg.sender_id === CURRENT_USER_ID
            });
            if (prepend) {
                messagesList.prepend(messageEl);
            } else {
                messagesList.appendChild(messageEl);
            }
        });
    };

    const loadMoreMessages = async () => {
        if (isLoading || !hasMoreMessages) return;
        isLoading = true;
        try {
            const cursor = chatStore.getState().messages[CONVERSATION_ID]?.[0]?.id;
            if (!cursor) {
                hasMoreMessages = false;
                return;
            }
            const olderMessages = await fetchMessages(CONVERSATION_ID, cursor);
            if (olderMessages.length > 0) {
                chatStore.prependMessages(CONVERSATION_ID, olderMessages);
            } else {
                hasMoreMessages = false;
            }
        } catch (error) {
            console.error('Failed to load older messages:', error);
        } finally {
            isLoading = false;
        }
    };

    messagesContainer.addEventListener('scroll', () => {
        // In reversed flex container, scrollTop is negative or 0.
        // We check if we are near the "top" which is visually the bottom of the scrollable area.
        if (messagesContainer.scrollHeight + messagesContainer.scrollTop - messagesContainer.clientHeight < 1) {
            loadMoreMessages();
        }
    });

    form.addEventListener('submit', async (e) => {
        e.preventDefault();
        const content = input.value.trim();
        if (content) {
            try {
                await sendMessage(CONVERSATION_ID, { type: 'text', content });
                input.value = '';
            } catch (error) {
                console.error('Failed to send message:', error);
                // TODO: Show error tooltip
            }
        }
    });

    fileInput.addEventListener('change', async (e) => {
        const file = e.target.files[0];
        if (!file) return;

        // Optional: Add file size check here (<= 25MB)

        try {
            const { upload_url, object_key } = await getPresignedUrl(file.name);
            await uploadFile(upload_url, file);
            await sendMessage(CONVERSATION_ID, { type: 'media', content: object_key });
        } catch (error) {
            console.error('Failed to upload file:', error);
            // TODO: Show error tooltip
        }
    });

    // Subscribe to chat store updates
    const unsubscribe = chatStore.subscribe(() => {
        const state = chatStore.getState();
        const currentMessages = state.messages[CONVERSATION_ID] || [];
        messagesList.innerHTML = ''; // Clear and re-render
        renderMessages(currentMessages);
    });

    // Initial load
    chatStore.setActiveConversation(CONVERSATION_ID);
    fetchMessages(CONVERSATION_ID)
        .then(messages => {
            chatStore.addMessage(CONVERSATION_ID, messages); // This will trigger the subscription
        })
        .catch(error => console.error('Failed to fetch initial messages:', error));

    // Connect WebSocket
    WebSocketService.connect()

    // Cleanup on component removal
    const observer = new MutationObserver((mutations, obs) => {
        if (!document.body.contains(container)) {
            unsubscribe();
            obs.disconnect();
        }
    });
    observer.observe(document.body, { childList: true, subtree: true });

    return container;
}
