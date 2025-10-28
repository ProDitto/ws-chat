import { Navbar } from '../components/Navbar.js';
import { getProfile } from '../api/users.js';

export function ProfilePage() {
    const page = document.createElement('div');
    page.className = 'h-screen flex flex-col';
    page.append(Navbar());

    const content = document.createElement('div');
    content.className = 'flex-grow p-8';
    content.innerHTML = `<h1 class="text-3xl font-bold mb-4">User Profile</h1>`;

    const profileContainer = document.createElement('div');
    profileContainer.id = 'profile-container';
    profileContainer.textContent = 'Loading profile...';
    content.append(profileContainer);

    // For simplicity, this page shows the logged-in user's profile.
    // A real implementation would use the hash for other users.
    getProfile('me').catch(() => getProfile('testuser')) // Fallback for dev
        .then(user => {
            profileContainer.innerHTML = `
                <div class="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4">
                    <p><strong>Username:</strong> ${user.username}</p>
                    <p><strong>Display Name:</strong> ${user.display_name}</p>
                    <p><strong>Joined:</strong> ${new Date(user.created_at).toLocaleDateString()}</p>
                    <div class="mt-4">
                        <button class="bg-blue-500 text-white px-4 py-2 rounded mr-2">Edit Profile</button>
                        <button class="bg-green-500 text-white px-4 py-2 rounded mr-2">Add Friend</button>
                        <button class="bg-red-500 text-white px-4 py-2 rounded">Block User</button>
                    </div>
                </div>
            `;
        })
        .catch(err => {
            profileContainer.textContent = `Error loading profile: ${err.message}`;
        });

    page.append(content);
    return page;
}

