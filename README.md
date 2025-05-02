# Запуск проекта

Для корректной работы программы нужен файл `config.yaml` в корне проекта. Образец представлен в репозитории

### Запуск

```
docker-compose up --build todo-list-server
```

# Изменение базы данных psql

В данном проекте используется библиотека **Goose** для управления миграциями базы данных. Все миграции
хранятся в папке `migrations`

# Базовый URL

```
http://your-domain:3000/api/
```

## Структура таблицы `tasks`

```sql
CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    status TEXT CHECK (status IN ('new', 'in_progress', 'done')) DEFAULT 'new',
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);
```

# Эндпоинты

## 1. Создание задачи

POST /tasks
Создает новую задачу.

```json
{
  "title": "Название задачи",
  "description": "Описание задачи",       // необязательно
  "status": "new"                         // необязательно, по умолчанию "new"
}
```

## 2. Получение списка всех задач

GET /tasks
Возвращает список всех задач.

## 3. Обновление задачи

PUT /tasks/:id
Обновляет задачу с указанным id.

### Параметры URL

id - Идентификатор задачи

```json
{
  "title": "Новое название",           // необязательно
  "description": "Новое описание",     // необязательно
  "status": "in_progress"               // необязательно, должно быть одним из ['new', 'in_progress', 'done']
}
```

## 4. Удаление задачи

DELETE /tasks/:id
Удаляет задачу с указанным id.

### Параметры URL

id - Идентификатор задачи


