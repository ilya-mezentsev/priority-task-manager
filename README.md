# priority-task-manager

## Тестовый стенд для демонстрации работы ограничения нагрузки на основе ролей клиентов
В качестве системы ограничения доступа используется проект [mini-roles-manager](https://github.com/ilya-mezentsev/mini-roles-manager), а для обоснования принятия решений - теория систем массового обслуживания (СМО)

## Локальный запуск

1. Собираем и запускаем контейнеры
```bash
$ make containers-build
$ make containers-run
```

2. Запускаем сервисы
"Управленец" задач (см. [Описание](./source/manager/README.md))
```bash
$ make manager-run
```

Исполнитель (см. [Описание](./source/executor/README.md))
```bash
$ make executor-run
```

Экспорт метрик (см. [Описание](./source/metrics/README.md))
```bash
$ make metrics-run
```

3. Вспомогательные скрипты
Обертка над исполнителем задач
```bash
$ make task-executor-wrapper
```

Инструмент для нагрузочного тестирования
```bash
$ make locust
```

4. Далее проходим по адресу, сформированному locust-ом (скорее всего `http://0.0.0.0:8089`) и нажимаем на кнопку с текстом `Start swarming`

5. Наконец, по адресу `http://localhost:9090` можно наблюдать экспортируемые метрики
