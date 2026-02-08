BEGIN;

CREATE TABLE IF NOT EXISTS gender
(
    id serial NOT NULL,
    name character varying(255) NOT NULL,
    CONSTRAINT gender_pkey PRIMARY KEY (id),
    CONSTRAINT gender_name_key UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS logos
(
    id serial NOT NULL,
    mime_type text NOT NULL,
    logo bytea NOT NULL,
    CONSTRAINT logos_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS match
(
    id serial NOT NULL,
    team1_id integer,
    team2_id integer,
    court_number integer,
    CONSTRAINT match_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS person
(
    id serial NOT NULL,
    name character varying(255) NOT NULL,
    CONSTRAINT person_pkey PRIMARY KEY (id),
    CONSTRAINT person_name_key UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS round_tournament
(
    tournament_id integer NOT NULL,
    match_id integer NOT NULL,
    round_number integer NOT NULL,
    CONSTRAINT round_tournament_pkey PRIMARY KEY (tournament_id, match_id)
);

CREATE TABLE IF NOT EXISTS sports_center
(
    id serial NOT NULL,
    logo_id integer,
    name character varying(255) COLLATE pg_catalog."default",
    CONSTRAINT sports_center_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS team
(
    id serial NOT NULL,
    person1_id integer,
    person2_id integer,
    gender_id integer DEFAULT 1,
    CONSTRAINT team_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS tournament
(
    id serial NOT NULL,
    tournament_date timestamp without time zone NOT NULL,
    tournament_type_id integer,
    user_id integer NOT NULL,
    CONSTRAINT tournament_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS tournament_type
(
    id serial NOT NULL,
    name character varying(255) NOT NULL,
    CONSTRAINT tournament_type_pkey PRIMARY KEY (id),
    CONSTRAINT tournament_type_name_key UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS users
(
    id bigint NOT NULL,
    name character varying(255) NOT NULL,
    sports_center_id integer NOT NULL,
    CONSTRAINT users_pkey PRIMARY KEY (id)
);

ALTER TABLE IF EXISTS match
    ADD CONSTRAINT match_team1_id_fkey FOREIGN KEY (team1_id)
    REFERENCES team (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;


ALTER TABLE IF EXISTS match
    ADD CONSTRAINT match_team2_id_fkey FOREIGN KEY (team2_id)
    REFERENCES team (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;


ALTER TABLE IF EXISTS round_tournament
    ADD CONSTRAINT round_tournament_match_id_fkey FOREIGN KEY (match_id)
    REFERENCES match (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;


ALTER TABLE IF EXISTS round_tournament
    ADD CONSTRAINT round_tournament_tournament_id_fkey FOREIGN KEY (tournament_id)
    REFERENCES tournament (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;


ALTER TABLE IF EXISTS sports_center
    ADD CONSTRAINT sports_center_logo_id_fkey FOREIGN KEY (logo_id)
    REFERENCES logos (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;


ALTER TABLE IF EXISTS team
    ADD CONSTRAINT team_gender_id_fkey FOREIGN KEY (gender_id)
    REFERENCES gender (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;


ALTER TABLE IF EXISTS team
    ADD CONSTRAINT team_person1_id_fkey FOREIGN KEY (person1_id)
    REFERENCES person (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;


ALTER TABLE IF EXISTS team
    ADD CONSTRAINT team_person2_id_fkey FOREIGN KEY (person2_id)
    REFERENCES person (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;


ALTER TABLE IF EXISTS tournament
    ADD CONSTRAINT tournament_tournament_type_id_fkey FOREIGN KEY (tournament_type_id)
    REFERENCES tournament_type (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;


ALTER TABLE IF EXISTS tournament
    ADD CONSTRAINT tournament_user_id_fkey FOREIGN KEY (user_id)
    REFERENCES users (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;


ALTER TABLE IF EXISTS users
    ADD CONSTRAINT users_sports_center_id_fkey FOREIGN KEY (sports_center_id)
    REFERENCES sports_center (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;

END;