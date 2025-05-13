document.getElementById('roleForm').addEventListener('submit', async (event) => {
    event.preventDefault();

    const formData = new FormData(event.target);
    const rawResponse = await fetch('/fgw/roles/add', {
        method: 'POST',
        body: formData
    });

    const contentType = rawResponse.headers.get('content-type');

    if (rawResponse.status === 422 && contentType.includes('application/json')) {
        const result = await rawResponse.json();
        Object.entries(result.errors).forEach(([field, message]) => {
            document.getElementById(`${field}-error`).alert(message);
        });
    } else if (rawResponse.ok) {
        alert('Роль успешно добавлена!');
        window.location.href = '/fgw/roles';
    } else {
        console.error('Ошибка при обработке запроса.', rawResponse.statusText);
    }
});

document.getElementById('employeeForm').addEventListener('submit', async (event) => {
    event.preventDefault();

    const formData = new FormData(event.target);
    const rawResponse = await fetch('/fgw/employees/add', {
        method: 'POST',
        body: formData
    });

    const contentType = rawResponse.headers.get('content-type');

    if (rawResponse.status === 422 && contentType.includes('application/json')) {
        const result = await rawResponse.json();
        Object.entries(result.errors).forEach(([field, message]) => {
            document.getElementById(`${field}-error`).alert(message);
        });
    } else if (rawResponse.ok) {
        alert('Сотрудник успешно добавлен!');
        window.location.href = '/fgw/employees';
    } else {
        console.error('Ошибка при обработке запроса.', rawResponse.statusText);
    }
});