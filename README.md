# О приложении
Бэкенд для мини мессенджера. Представляет из себя набор из 3-ёх микросервисов:
+ chat-app - микросервис для чата;
+ auth-service - микросервис авторизации пользователей;
+ log-service - микросервис логирования.

Общение между сервисами происходит с помощью gRPC.

![Схема](./docs/scheme.png)

Реализованный функционал:
+ можно общаться с другими пользователями в реальном времени (с помощью WebSocket);
+ можно создавать чаты с новыми пользователями (от 2-ух и более);
+ можно искать пользователей по никнейму;
+ регистрация новых пользователей;
+ авторизация пользователей с помощью JWT-токенов;
+ вся история переписки сохраняется в БД;
+ логирование реализовано отдельным микросервисом, логи пишутся в БД.

Также есть очень простой [клиент](https://github.com/1PALADIN1/gigachat_client_cc), который работает с данным приложением. Написан на [Cocos Creator](https://www.cocos.com/en/creator) с использованием [TypeScript](https://www.typescriptlang.org/).

В репозитории находится файл **postman_examples.json** с примерами запросов для Postman.


# Установка
## Настройка переменных окружения
Если требуется использовать проект в продакшн среде, то для начала нужно задать переменные окружения в файле project/docker-compose.yml (на данный момент стоят заглушки).


## Запуск в Docker (Makefile)
Проект предоставляется с Makefile, в котором реализованы основые команды для работы с приложением. Makefile настроен для использования в **Mac** или **Linux**.

Чтобы запустить проект, необходимо выполнить следующие команды в терминале:

    cd /project
    make up-build

Команда соберёт каждый микросервис и развернёт их в Docker.

Чтобы остановить сервер, используйте команду:

    make down


## Выполнение миграций
Для создания структуры БД необходимо выполнить миграции (выполнять из папки /project):

    migrate -path ./../chat-app/schema -database 'postgres://postgres:123456psql@localhost:5432/chat_db?sslmode=disable' up

И для таблицы логов:

    migrate -path ./../log-service/schema -database 'mysql://root:123456mysql@tcp(localhost:3306)/logs_db' up

Если пароли и адреса к базам данным были изменены, то в указанных выше командах тоже необходимо их поменять.

После успешной миграции **приложение готово к использованию!**

# Используемые инструменты
+ Go v1.19
+ [Gorilla Mux](https://github.com/gorilla/mux) для роутинга REST API запросов.
+ [Gorilla WebSocket](https://github.com/gorilla/websocket) для работы с протоколом WebSocket.
+ [Jwt Go](https://github.com/golang-jwt/jwt) для генерации JWT-токенов.
+ [Sqlx](https://github.com/jmoiron/sqlx) для работы с БД.
+ Для юнит тестов использованы:
  + [Testify](https://github.com/stretchr/testify) для использования assert в тестах.
  + [Gomock](https://github.com/golang/mock) для генерации заглушек.
  + [Sqlx Mock](https://github.com/zhashkevych/go-sqlxmock) для симуляции поведения БД.
+ [Migrate](https://github.com/golang-migrate/migrate) для миграций БД.
+ [PostreSQL](https://www.postgresql.org/) и [MySQL](https://www.mysql.com/) в качестве БД.
