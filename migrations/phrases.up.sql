CREATE TABLE if not exists diplom.phrases (
                         id UUID PRIMARY KEY,
                         text TEXT NOT NULL,
                         type_id UUID REFERENCES diplom.phrase_types(id)
);