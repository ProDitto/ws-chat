import { login, register } from '../api/auth.js';
import { setTokens } from '../store/authStore.js';

export default function LoginPage() {
    const container = document.createElement('div');
    container.className = 'flex min-h-full flex-col justify-center px-6 py-12 lg:px-8 bg-gray-900 text-white';
    container.innerHTML = `
        <div class="sm:mx-auto sm:w-full sm:max-w-sm">
            <h2 id="form-title" class="mt-10 text-center text-2xl font-bold leading-9 tracking-tight">
                Sign in to your account
            </h2>
        </div>

        <div class="mt-10 sm:mx-auto sm:w-full sm:max-w-sm">
            <!-- Login Form -->
            <form id="login-form" class="space-y-6">
                <div>
                    <label for="login-username" class="block text-sm font-medium leading-6">Username</label>
                    <div class="mt-2">
                        <input id="login-username" name="username" type="text" required
                            class="block w-full rounded-md border-0 bg-white/5 py-1.5 shadow-sm ring-1 ring-inset ring-white/10 focus:ring-2 focus:ring-inset focus:ring-indigo-500 sm:text-sm sm:leading-6">
                    </div>
                </div>
                <div>
                    <label for="login-password" class="block text-sm font-medium leading-6">Password</label>
                    <div class="mt-2">
                        <input id="login-password" name="password" type="password" required
                            class="block w-full rounded-md border-0 bg-white/5 py-1.5 shadow-sm ring-1 ring-inset ring-white/10 focus:ring-2 focus:ring-inset focus:ring-indigo-500 sm:text-sm sm:leading-6">
                    </div>
                </div>
                <div>
                    <button type="submit"
                        class="flex w-full justify-center rounded-md bg-indigo-500 px-3 py-1.5 text-sm font-semibold leading-6 shadow-sm hover:bg-indigo-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-500">
                        Sign in
                    </button>
                </div>
                <p id="login-error" class="text-red-400 text-sm mt-2 text-center"></p>
            </form>

            <!-- Register Form (hidden by default) -->
            <form id="register-form" class="space-y-6 hidden">
                <div>
                    <label for="register-username" class="block text-sm font-medium leading-6">Username</label>
                    <div class="mt-2">
                        <input id="register-username" name="username" type="text" required
                            class="block w-full rounded-md border-0 bg-white/5 py-1.5 shadow-sm ring-1 ring-inset ring-white/10 focus:ring-2 focus:ring-inset focus:ring-indigo-500 sm:text-sm sm:leading-6">
                    </div>
                </div>
                <div>
                    <label for="register-password" class="block text-sm font-medium leading-6">Password</label>
                    <div class="mt-2">
                        <input id="register-password" name="password" type="password" required
                            class="block w-full rounded-md border-0 bg-white/5 py-1.5 shadow-sm ring-1 ring-inset ring-white/10 focus:ring-2 focus:ring-inset focus:ring-indigo-500 sm:text-sm sm:leading-6">
                    </div>
                     <p class="mt-2 text-xs text-gray-400">Password must be at least 8 characters long, with uppercase, lowercase, number, and special character.</p>
                </div>
                <div>
                    <button type="submit"
                        class="flex w-full justify-center rounded-md bg-indigo-500 px-3 py-1.5 text-sm font-semibold leading-6 shadow-sm hover:bg-indigo-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-500">
                        Create account
                    </button>
                </div>
                <p id="register-error" class="text-red-400 text-sm mt-2 text-center"></p>
            </form>

            <p class="mt-10 text-center text-sm text-gray-400">
                <a id="toggle-form" href="#" class="font-semibold leading-6 text-indigo-400 hover:text-indigo-300">
                    Not a member? Create an account
                </a>
            </p>
        </div>
    `;

    const loginForm = container.querySelector('#login-form');
    const registerForm = container.querySelector('#register-form');
    const toggleLink = container.querySelector('#toggle-form');
    const formTitle = container.querySelector('#form-title');
    const loginError = container.querySelector('#login-error');
    const registerError = container.querySelector('#register-error');

    let isLogin = true;

    toggleLink.addEventListener('click', (e) => {
        e.preventDefault();
        isLogin = !isLogin;
        loginForm.classList.toggle('hidden');
        registerForm.classList.toggle('hidden');
        formTitle.textContent = isLogin ? 'Sign in to your account' : 'Create a new account';
        toggleLink.innerHTML = isLogin ? 'Not a member? Create an account' : 'Already have an account? Sign in';
        loginError.textContent = '';
        registerError.textContent = '';
    });

    loginForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        loginError.textContent = '';
        const formData = new FormData(loginForm);
        const username = formData.get('username');
        const password = formData.get('password');

        try {
            const data = await login(username, password);
            setTokens(data.access_token, data.refresh_token);
            window.location.hash = '/';
        } catch (error) {
            loginError.textContent = error.message;
        }
    });

    registerForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        registerError.textContent = '';
        const formData = new FormData(registerForm);
        const username = formData.get('username');
        const password = formData.get('password');

        try {
            await register(username, password);
            // Automatically log in after successful registration
            const data = await login(username, password);
            setTokens(data.access_token, data.refresh_token);
            window.location.hash = '/';
        } catch (error) {
            registerError.textContent = error.message;
        }
    });

    return container;
}