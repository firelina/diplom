CREATE TABLE if not exists diplom.audio_phrases (
                               id UUID PRIMARY KEY,
                               path_to_audio TEXT NOT NULL,
                               phrase_id UUID REFERENCES diplom.phrases(id),
                               accent TEXT,
                               noise INT
);