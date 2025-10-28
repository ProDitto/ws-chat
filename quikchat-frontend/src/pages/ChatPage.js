export default function ChatPage() {
    const handleLogout = () => {
        document.dispatchEvent(new CustomEvent('logout'));
    };

    const container = document.createElement('div');
    container.className = 'flex flex-col items-center justify-center h-full bg-gray-800 text-white';
    container.innerHTML = `
        <div class="text-center">
            <h1 class="text-4xl font-bold mb-4">Welcome to QuikChat</h1>
            <p class="text-lg text-gray-300 mb-8">You are logged in. Chat functionality coming soon!</p>
            <button id="logout-button" class="px-4 py-2 bg-red-600 hover:bg-red-700 rounded-md text-white font-semibold">
                Logout
            </button>
        </div>
    `;

    container.querySelector('#logout-button').addEventListener('click', handleLogout);

    return container;
}
```
