CREATE TABLE chats (
    id serial not null unique,
    title varchar(100) not null,
    description varchar(255)
);

CREATE TABLE users_chats (
    id serial not null unique,
    user_id int references users(id) on delete cascade not null,
    chat_id int references chats(id) on delete cascade not null
);

CREATE UNIQUE INDEX users_chats_unique_index
ON users_chats (user_id, chat_id);