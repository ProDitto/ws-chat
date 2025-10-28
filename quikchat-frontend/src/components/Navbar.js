export function Navbar() {
    const nav = document.createElement('nav');
    nav.className = 'bg-gray-800 text-white p-4 flex justify-between items-center';

    const title = document.createElement('h1');
    title.className = 'text-xl font-bold';
    title.textContent = 'QuikChat';

    const links = document.createElement('div');
    links.innerHTML = `
        <a href="#/" class="px-3 py-2 rounded-md text-sm font-medium hover:bg-gray-700">Chat</a>
        <a href="#/friends" class="px-3 py-2 rounded-md text-sm font-medium hover:bg-gray-700">Friends</a>
        <a href="#/profile" class="px-3 py-2 rounded-md text-sm font-medium hover:bg-gray-700">My Profile</a>
    `;

    nav.append(title, links);
    return nav;
}

