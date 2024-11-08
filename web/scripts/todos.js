// Открытие и закрытие модального окна
function openModal() {
    document.getElementById("modal").style.display = "flex";
}

function closeModal() {
    document.getElementById("modal").style.display = "none";
}

// Асинхронное добавление задачи
function addTodo() {
    const title = document.getElementById("title").value;
    const description = document.getElementById("description").value;

    fetch('/todos', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({title, description})
    })
        .then(response => {
            if (response.ok) {
                return response.json(); // Ожидаем JSON с новой задачей
            } else {
                alert('Failed to add task');
            }
        })
        .then(newTodo => {
            if (newTodo) {
                // Обновляем список задач, добавляя новую задачу
                addTodoToList(newTodo);
                closeModal(); // Закрываем модальное окно
                document.getElementById("todoForm").reset(); // Сбрасываем форму
            }
        })
        .catch(error => console.error('Error:', error));
}


// Функция для добавления задачи в список на странице
function addTodoToList(todo) {
    const todoList = document.querySelector('.todo-list');
    const todoItem = document.createElement('div');
    todoItem.classList.add('todo-item');
    todoItem.setAttribute('data-id', todo.ID);

    todoItem.innerHTML = `
            <div class="todo-header">
                <input type="checkbox" onclick="toggleTodoStatus(${todo.ID})" ${todo.isDone ? 'checked' : ''} class="todo-status">
                <span class="todo-title">${todo.Title}</span>
                <div class="todo-actions">
                    <button onclick="openTodoMenu(event, ${todo.ID})" class="menu-btn">⋮</button>
                    <div class="todo-menu" id="menu-${todo.ID}">
                        <button onclick="editTodo(${todo.ID})">Edit</button>
                        <button onclick="deleteTodo(${todo.ID})">Delete</button>
                    </div>
                </div>
            </div>
            <p class="todo-description">${todo.Description}</p>
        `;

    todoList.appendChild(todoItem);
}

// Переключение статуса задачи (выполнено/не выполнено)
function toggleTodoStatus(id) {
    const checkbox = document.querySelector(`.todo-item[data-id='${id}'] .todo-status`);
    const isChecked = checkbox.checked;

    fetch(`/todos/${id}`, {
        method: 'PATCH',
        body: JSON.stringify({
            isDone: isChecked
        })
    })
        .then(response => {
            if (response.ok) {
                // Перезагружаем страницу, чтобы обновить статус задачи
                location.reload();
            } else {
                alert('Failed to update task status');
            }
        })
        .catch(error => console.error('Error:', error));
}

// Открытие и закрытие меню действий для задачи
function openTodoMenu(event, id) {
    event.stopPropagation();
    const menu = document.getElementById(`menu-${id}`);
    const isMenuOpen = menu.style.display === 'block';

    // Скрываем все открытые меню
    document.querySelectorAll('.todo-menu').forEach(menu => {
        menu.style.display = 'none';
    });

    // Отображаем выбранное меню, если оно было закрыто
    if (!isMenuOpen) {
        menu.style.display = 'block';
    }
}

// Закрытие всех открытых меню при клике вне их области
document.addEventListener('click', () => {
    document.querySelectorAll('.todo-menu').forEach(menu => {
        menu.style.display = 'none';
    });
});

// Функция для редактирования задачи
function editTodo(id) {
    const newTitle = prompt("Enter new title:");
    const newDescription = prompt("Enter new description:");

    if (newTitle && newDescription) {
        fetch(`/todos/${id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                title: newTitle,
                description: newDescription
            })
        })
            .then(response => {
                if (response.ok) {
                    location.reload();
                } else {
                    alert("Failed to update todo");
                }
            })
            .catch(error => console.error("Error:", error));
    }
}

// Функция для удаления задачи
function deleteTodo(id) {
    fetch(`/todos/${id}`, {
        method: 'DELETE'
    })
        .then(response => {
            if (response.ok) {
                location.reload();
            } else {
                alert('Failed to delete task');
            }
        })
        .catch(error => console.error('Error:', error));
}
