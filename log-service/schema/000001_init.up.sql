CREATE TABLE logs (
    id INT NOT NULL,
    level varchar(50),
    source varchar(50),
    message varchar(255),
    PRIMARY KEY (id)
);