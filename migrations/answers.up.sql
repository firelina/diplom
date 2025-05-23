CREATE TABLE if not exists diplom.answers (
                         id UUID PRIMARY KEY,
                         user_id UUID REFERENCES diplom.users(id),
                         audio_answer_id UUID REFERENCES diplom.audio_answers(id),
                         text TEXT,
                         is_correct BOOLEAN
);