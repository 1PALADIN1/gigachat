# О приложении
Бэкенд для мини мессенджера. Особенности:
+ можно общаться с другими пользователями в реальном времени (с помощью WebSocket);
+ можно создавать чаты с новыми пользователями (от 2-ух и более);
+ можно искать пользователей по никнейму;
+ регистрация новых пользователей;
+ автризация происходит с помощью JWT-токенов;
+ предоставляется готовый docker compose файл для разворачивания приложения в контейнере;
+ вся история переписки сохраняется в БД.

Также есть очень простой [клиент](https://github.com/1PALADIN1/gigachat_client_cc), который работает с данным приложением. Написан на [Cocos Creator](https://www.cocos.com/en/creator) с использованием [TypeScript](https://www.typescriptlang.org/).

В репозитории находится файл **postman_examples.json** с примерами запросов для Postman.

# Установка
## Настройка переменных окружения
Для начала нужно задать переменные в файле docker-compose.yaml (на данный момент стоят заглушки).

Для services.environment задать:
+ DB_PASSWORD - пароль от БД;
+ PASSWORD_HASH_SALT - соль для хэша пароля;
+ SINGING_KEY - ключ для подписи JWT-токена.

Для db.environment нужно задать пароль POSTGRES_PASSWORD.

## Запуск в Docker
Далее необходимо выполнить команду в терминале:

    docker-compose up --build gigachat-app


## Выполнение миграций
Для создания структуры БД необходимо выполнить миграции:

    migrate -path ./schema -database 'postgres://postgres:PASSWORD@localhost:5432/postgres?sslmode=disable' up

Вместо **PASSWORD** нужно указать пароль от БД. При необходимости заменить **localhost:5432**, если БД находится по другому адресу.

После успешной миграции **приложение готово к использованию!**

# Используемые инструменты
+ [Gorilla Mux](https://github.com/gorilla/mux) для роутинга REST API запросов.
+ [Gorilla WebSocket](https://github.com/gorilla/websocket) для работы с протоколом WebSocket.
+ [Jwt Go](https://github.com/golang-jwt/jwt) для генерации JWT-токенов.
+ [Sqlx](https://github.com/jmoiron/sqlx) для работы с БД.
+ Для юнит тестов использованы:
  + [Testify](https://github.com/stretchr/testify) для использования assert в тестах.
  + [Gomock](https://github.com/golang/mock) для генерации заглушек.
  + [Sqlx Mock](https://github.com/zhashkevych/go-sqlxmock) для симуляции поведения БД.
+ [Migrate](https://github.com/golang-migrate/migrate) для миграций БД.
+ [PostreSQL](https://www.postgresql.org/) в качестве СУБД.
