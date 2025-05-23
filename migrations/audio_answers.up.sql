CREATE TABLE if not exists diplom.audio_answers (
                               id UUID PRIMARY KEY,
                               path_to_audio TEXT NOT NULL,
                               record_time TIMESTAMP NOT NULL
);