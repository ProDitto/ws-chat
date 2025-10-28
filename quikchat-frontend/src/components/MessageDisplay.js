/**
 * @param {object} props
 * @param {import('../types/types.js').Message} props.message
 * @param {boolean} props.isOwnMessage
 */
export function MessageDisplay({ message, isOwnMessage }) {
    const container = document.createElement('div');
    container.className = `flex items-start gap-2.5 my-2 ${isOwnMessage ? 'justify-end' : ''}`;

    const contentContainer = document.createElement('div');
    contentContainer.className = `flex flex-col w-full max-w-[320px] leading-1.5 p-4 border-gray-200 rounded-e-xl rounded-es-xl ${isOwnMessage ? 'bg-blue-600 text-white rounded-s-xl rounded-ee-none' : 'bg-gray-100 dark:bg-gray-700'}`;

    const senderName = document.createElement('span');
    senderName.className = 'text-sm font-semibold text-gray-900 dark:text-white';
    senderName.textContent = message.sender?.display_name || 'Unknown User';

    const messageContent = document.createElement('p');
    messageContent.className = 'text-sm font-normal py-2.5';
    messageContent.textContent = message.content;

    const timestamp = document.createElement('span');
    timestamp.className = `text-xs font-normal ${isOwnMessage ? 'text-blue-200' : 'text-gray-500 dark:text-gray-400'}`;
    timestamp.textContent = new Date(message.created_at).toLocaleTimeString();

    if (!isOwnMessage) {
        contentContainer.appendChild(senderName);
    }
    contentContainer.appendChild(messageContent);
    contentContainer.appendChild(timestamp);

    container.appendChild(contentContainer);

    return container;
}

