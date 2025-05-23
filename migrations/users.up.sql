CREATE TABLE if not exists diplom.users   (
                       id UUID PRIMARY KEY,
                       name TEXT NOT NULL,
                       login TEXT NOT NULL,
                       password TEXT NOT NULL,
                       role BOOLEAN NOT NULL
);