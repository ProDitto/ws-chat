import { Navbar } from '../components/Navbar.js';
import { MessageDisplay } from '../components/MessageDisplay.js';
import { chatStore } from '../store/chatStore.js';
import { WebSocketService } from '../services/WebSocketService.js';
import { fetchMessages, sendMessage, getPresignedUrl, uploadFile } from '../api/messages.js';

// Hardcoded for now. In a real app, this would come from a user store after login.
// For now, we assume user ID 1 is logged in and talking in conversation 1.
// You would need to run a SQL command to create these, e.g.:
// INSERT INTO conversations (id, type) VALUES (1, 'private');
// INSERT INTO conversation_participants (conversation_id, user_id) VALUES (1, 1);
// INSERT INTO conversation_participants (conversation_id, user_id) VALUES (1, 2);
const CURRENT_USER_ID = 1;
const CONVERSATION_ID = 1;

export function ChatPage() {
    const page = document.createElement('div');
    page.className = 'h-screen flex flex-col';
    page.appendChild(Navbar());

    const chatContainer = document.createElement('div');
    chatContainer.className = 'flex-1 flex flex-col p-4 overflow-hidden';

    const messagesArea = document.createElement('div');
    messagesArea.id = 'messages-area';
    messagesArea.className = 'flex-1 overflow-y-auto flex flex-col-reverse p-4 bg-gray-50';
    
    const messageList = document.createElement('div');
    messagesArea.appendChild(messageList);

    const form = document.createElement('form');
    form.className = 'mt-4 flex items-center gap-2';

    const fileInput = document.createElement('input');
    fileInput.type = 'file';
    fileInput.className = 'hidden';
    fileInput.id = 'file-input';
    fileInput.accept = 'image/jpeg,image/png,image/gif,video/mp4';

    const attachButton = document.createElement('button');
    attachButton.type = 'button';
    attachButton.innerHTML = `<svg class="w-6 h-6 text-gray-500" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 20 18"><path fill="currentColor" d="M13 5.5a.5.5 0 1 1-1 0 .5.5 0 0 1 1 0ZM7.565 7.423 4.5 14h11.518l-2.516-3.71L11 13 7.565 7.423Z"/><path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 1H2a1 1 0 0 0-1 1v14a1 1 0 0 0 1 1h16a1 1 0 0 0 1-1V2a1 1 0 0 0-1-1Z"/><path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 5.5a.5.5 0 1 1-1 0 .5.5 0 0 1 1 0Z"/></svg>`;
    attachButton.className = 'p-2 rounded-full bg-gray-200 hover:bg-gray-300';
    attachButton.onclick = () => fileInput.click();

    const textInput = document.createElement('input');
    textInput.type = 'text';
    textInput.placeholder = 'Type a message...';
    textInput.className = 'flex-1 p-2 border rounded-lg';
    textInput.autocomplete = 'off';

    const sendButton = document.createElement('button');
    sendButton.type = 'submit';
    sendButton.textContent = 'Send';
    sendButton.className = 'px-4 py-2 bg-blue-500 text-white rounded-lg';

    form.append(attachButton, fileInput, textInput, sendButton);
    chatContainer.append(messagesArea, form);
    page.appendChild(chatContainer);

    let isLoading = false;
    let hasMoreMessages = true;
    let oldestMessageId = null;

    const renderMessages = () => {
        const state = chatStore.getState();
        const messages = state.messages[CONVERSATION_ID] || [];
        messageList.innerHTML = '';
        messages.forEach(msg => {
            const isOwn = msg.sender_id === CURRENT_USER_ID;
            messageList.prepend(MessageDisplay({ message: msg, isOwnMessage: isOwn }));
        });
        if (messages.length > 0) {
            oldestMessageId = messages[0].id;
        }
    };

    const loadMoreMessages = async () => {
        if (isLoading || !hasMoreMessages) return;
        isLoading = true;
        try {
            const olderMessages = await fetchMessages(CONVERSATION_ID, oldestMessageId);
            if (olderMessages && olderMessages.length > 0) {
                chatStore.prependMessages(CONVERSATION_ID, olderMessages);
            } else {
                hasMoreMessages = false;
            }
        } catch (error) {
            console.error('Failed to load more messages:', error);
        } finally {
            isLoading = false;
        }
    };

    messagesArea.addEventListener('scroll', () => {
        if (messagesArea.scrollTop === 0) {
            loadMoreMessages();
        }
    });

    const handleNewMessage = (message) => {
        if (message.conversation_id === CONVERSATION_ID) {
            chatStore.addMessage(CONVERSATION_ID, message);
            if (messagesArea.scrollHeight - messagesArea.scrollTop - messagesArea.clientHeight < 200) {
                setTimeout(() => messagesArea.scrollTop = messagesArea.scrollHeight, 0);
            }
        }
    };

    form.addEventListener('submit', async (e) => {
        e.preventDefault();
        const content = textInput.value.trim();
        if (content) {
            try {
                await sendMessage(CONVERSATION_ID, { type: 'text', content });
                textInput.value = '';
            } catch (error) {
                console.error('Failed to send message:', error);
            }
        }
    });

    fileInput.addEventListener('change', async (e) => {
        const file = e.target.files[0];
        if (!file) return;
        try {
            const { upload_url, key } = await getPresignedUrl(file.name);
            await uploadFile(upload_url, file);
            await sendMessage(CONVERSATION_ID, { type: 'media', content: key });
        } catch (error) {
            console.error('Failed to upload file and send message:', error);
        }
    });

    const unsubscribe = chatStore.subscribe(renderMessages);
    WebSocketService.addMessageListener(handleNewMessage);
    
    chatStore.setActiveConversation(CONVERSATION_ID);
    loadMoreMessages().then(() => {
        setTimeout(() => messagesArea.scrollTop = messagesArea.scrollHeight, 0);
    });

    page.addEventListener('DOMNodeRemoved', () => {
        unsubscribe();
        WebSocketService.removeMessageListener(handleNewMessage);
    });

    return page;
}
