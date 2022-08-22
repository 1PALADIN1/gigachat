CREATE TABLE chats (
    id serial not null unique,
    title varchar(30) not null,
    description varchar(255)
);

CREATE TABLE users_chats (
    id serial not null unique,
    user_id int references users(id) on delete cascade not null,
    chat_id int references chats(id) on delete cascade not null
);