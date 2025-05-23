CREATE TABLE if not exists diplom.scenarios (
                           id UUID PRIMARY KEY,
                           title TEXT NOT NULL,
                           status TEXT NOT NULL,
                           start_date TIMESTAMP,
                           end_date TIMESTAMP,
                           user_id UUID REFERENCES diplom.users(id)
);