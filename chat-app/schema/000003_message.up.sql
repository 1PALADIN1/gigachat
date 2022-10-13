CREATE TABLE messages (
    id serial not null unique,
    message varchar(255) not null,
    send_date_time timestamp,
    chat_id int references chats(id) on delete cascade not null,
    user_id int references users(id) on delete cascade not null
);