CREATE TABLE if not exists diplom.phrase_streams (
                                id UUID PRIMARY KEY,
                                audio_phrase_id UUID REFERENCES diplom.audio_phrases(id),
                                scenario_id UUID REFERENCES diplom.scenarios(id),
                                answer_id UUID REFERENCES diplom.answers(id),
                                phrase_id UUID REFERENCES diplom.phrases(id),
                                status TEXT NOT NULL
);