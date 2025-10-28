import { Navbar } from '../components/Navbar.js';
import { listFriends, listIncomingRequests, respondToRequest } from '../api/friends.js';

export function FriendsPage() {
    const page = document.createElement('div');
    page.className = 'h-screen flex flex-col';
    page.append(Navbar());

    const content = document.createElement('div');
    content.className = 'flex-grow p-8';
    content.innerHTML = `
        <h1 class="text-3xl font-bold mb-4">Friends & Requests</h1>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
            <div>
                <h2 class="text-2xl font-semibold mb-2">Friend Requests</h2>
                <div id="requests-list" class="space-y-2">Loading...</div>
            </div>
            <div>
                <h2 class="text-2xl font-semibold mb-2">Friends List</h2>
                <div id="friends-list" class="space-y-2">Loading...</div>
            </div>
        </div>
    `;

    page.append(content);

    const requestsList = page.querySelector('#requests-list');
    listIncomingRequests()
        .then(requests => {
            if (requests.length === 0) {
                requestsList.textContent = 'No new friend requests.';
                return;
            }
            requestsList.innerHTML = '';
            requests.forEach(req => {
                const reqEl = document.createElement('div');
                reqEl.className = 'p-2 border rounded flex justify-between items-center';
                reqEl.innerHTML = `<span>Request from user ID: ${req.sender_id}</span>`;
                
                const acceptBtn = document.createElement('button');
                acceptBtn.textContent = 'Accept';
                acceptBtn.className = 'bg-green-500 text-white px-2 py-1 rounded text-sm';
                acceptBtn.onclick = () => respondToRequest(req.id, 'accept').then(() => reqEl.remove());

                const declineBtn = document.createElement('button');
                declineBtn.textContent = 'Decline';
                declineBtn.className = 'bg-red-500 text-white px-2 py-1 rounded text-sm ml-2';
                declineBtn.onclick = () => respondToRequest(req.id, 'decline').then(() => reqEl.remove());

                const btnContainer = document.createElement('div');
                btnContainer.append(acceptBtn, declineBtn);
                reqEl.append(btnContainer);
                requestsList.append(reqEl);
            });
        })
        .catch(err => requestsList.textContent = `Error: ${err.message}`);

    const friendsList = page.querySelector('#friends-list');
    listFriends()
        .then(friends => {
            if (friends.length === 0) {
                friendsList.textContent = 'You have no friends yet.';
                return;
            }
            friendsList.innerHTML = '';
            friends.forEach(friend => {
                const friendEl = document.createElement('div');
                friendEl.className = 'p-2 border rounded';
                friendEl.textContent = `${friend.display_name} (@${friend.username})`;
                friendsList.append(friendEl);
            });
        })
        .catch(err => friendsList.textContent = `Error: ${err.message}`);

    return page;
}

