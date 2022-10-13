CREATE TABLE logs (
    id int NOT NULL auto_increment primary key,
    level varchar(50),
    source varchar(50),
    message varchar(255),
    created_at datetime default CURRENT_TIMESTAMP
);