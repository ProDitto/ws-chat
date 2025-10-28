const appElement = document.getElementById('app');

export function render(component) {
    if (!appElement) {
        console.error("Root element #app not found");
        return;
    }
    appElement.innerHTML = '';
    appElement.append(component());
}

