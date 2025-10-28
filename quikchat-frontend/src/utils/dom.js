const appElement = document.getElementById('app');

export function render(component) {
    appElement.innerHTML = '';
    appElement.appendChild(component);
}
```
