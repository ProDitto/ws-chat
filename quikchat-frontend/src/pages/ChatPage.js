import { Navbar } from '../components/Navbar.js';

function handleLogout() {
    const logoutEvent = new CustomEvent('logout');
    document.dispatchEvent(logoutEvent);
}

export function ChatPage() {
    const page = document.createElement('div');
    page.className = 'h-screen flex flex-col';
    page.append(Navbar());

    const content = document.createElement('div');
    content.className = 'flex-grow p-8';
    content.innerHTML = `
        <h1 class="text-3xl font-bold mb-4">Welcome to QuikChat!</h1>
        <p class="mb-6">This is the main chat interface. (Placeholder)</p>
    `;

    const logoutButton = document.createElement('button');
    logoutButton.textContent = 'Logout';
    logoutButton.className = 'bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 rounded';
    logoutButton.addEventListener('click', handleLogout);

    content.append(logoutButton);
    page.append(content);

    return page;
}

