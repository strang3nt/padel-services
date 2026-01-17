CREATE TABLE IF NOT EXISTS tournament_type (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS gender (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

INSERT INTO gender (name) VALUES ('Male'), ('Female') ON CONFLICT (name) DO NOTHING;
INSERT INTO tournament_type (name) VALUES ('Rodeo') ON CONFLICT (name) DO NOTHING;

CREATE TABLE IF NOT EXISTS person (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS tournament (
    id SERIAL PRIMARY KEY,
    tournament_date TIMESTAMP NOT NULL,
    tournament_type_id INT REFERENCES tournament_type(id)
);

CREATE TABLE IF NOT EXISTS team (
    id SERIAL PRIMARY KEY,
    person1_id INT REFERENCES person(id),
    person2_id INT REFERENCES person(id),
    gender_id INT REFERENCES gender(id) DEFAULT 1
);

CREATE TABLE IF NOT EXISTS match (
    id SERIAL PRIMARY KEY,
    team1_id INT REFERENCES team(id),
    team2_id INT REFERENCES team(id),
    court_number INT NOT NULL
);

CREATE TABLE IF NOT EXISTS round_tournament (
    tournament_id INT REFERENCES tournament(id),
    match_id INT REFERENCES match(id),
    round_number INT NOT NULL,
    PRIMARY KEY (tournament_id, match_id)
);