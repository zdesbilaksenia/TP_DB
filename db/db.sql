CREATE
    EXTENSION IF NOT EXISTS citext;

drop table if exists thread, forum, users, votes, post, nickname_forum CASCADE;

CREATE
    UNLOGGED TABLE users
(
    nickname citext collate "C" not null
        constraint users_pk
            primary key,
    fullname text,
    about    text,
    email    citext unique
);

CREATE
    UNLOGGED TABLE forum
(
    title   text,
    "user"  citext
        constraint forum_users_nickname_fk
            references users,
    slug    citext not null
        constraint forum_pk
            primary key,
    posts   bigint default 0,
    threads bigint default 0
);

CREATE
    UNLOGGED TABLE post
(
    id       bigserial
        constraint post_pkey
            primary key,
    parent   bigint                   default 0,
    author   citext
        constraint post_users_nickname_fk
            references users,
    message  text,
    isedited boolean                  default false,
    forum    citext
        constraint post_forum_slug_fk
            references forum,
    thread   integer,
    created  timestamp with time zone default now(),
    path     INT[]                    DEFAULT ARRAY []::INTEGER[]
);

CREATE
    OR REPLACE FUNCTION add_post() RETURNS TRIGGER AS
$$
BEGIN
    UPDATE forum
    SET posts = posts + 1
    WHERE Slug = NEW.forum;
    NEW.path
        = (SELECT path FROM post WHERE id = NEW.parent LIMIT 1) || NEW.id;
    RETURN NEW;
END
$$
    LANGUAGE 'plpgsql';

CREATE TRIGGER add_post
    BEFORE INSERT
    ON post
    FOR EACH ROW
EXECUTE PROCEDURE add_post();

CREATE
    UNLOGGED TABLE thread
(
    id      serial
        constraint thread_pk
            primary key,
    title   text,
    author  citext
        constraint thread_users_nickname_fk
            references users,
    forum   citext
        constraint thread_forum_slug_fk
            references forum,
    message text,
    votes   integer                  default 0,
    slug    citext,
    created timestamp with time zone default now()
);

CREATE
    OR REPLACE FUNCTION add_thread() RETURNS TRIGGER AS
$$
BEGIN
    UPDATE forum
    SET threads = threads + 1
    WHERE NEW.Forum = slug;
    RETURN NULL;
END
$$
    LANGUAGE 'plpgsql';

CREATE TRIGGER create_thread_trigger
    AFTER INSERT
    ON thread
    FOR EACH ROW
EXECUTE PROCEDURE add_thread();

CREATE
    UNLOGGED TABLE votes
(
    id       bigserial
        constraint votes_pkey
            primary key,
    nickname citext
        constraint votes_nickname_fkey
            references users,
    voice    integer,
    thread   integer not null
        constraint votes_thread_fkey
            references thread,
    constraint votes_thread_nickname_key
        unique (thread, nickname)
);

create
    or replace function add_vote() returns trigger as
$$
BEGIN
    UPDATE thread
    SET votes=(votes + NEW.voice)
    WHERE id = NEW.thread;
    RETURN NEW;
END
$$
    LANGUAGE 'plpgsql';

CREATE TRIGGER add_vote
    AFTER INSERT
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE add_vote();

create or replace function update_vote() returns trigger as
$$
BEGIN
    IF OLD.voice <> NEW.voice THEN
        UPDATE thread
        SET votes = votes - OLD.voice + NEW.voice
        WHERE id = NEW.thread;
    END IF;
    RETURN NEW;
END
$$
    LANGUAGE 'plpgsql';

CREATE TRIGGER update_vote
    AFTER UPDATE
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE update_vote();

CREATE UNLOGGED TABLE nickname_forum
(
    nickname citext collate "C"
        constraint forum_users_nickname_fkey
            references users,
    fullname text,
    about    text,
    email    citext,
    forum    citext
        constraint forum_users_forum_slug
            references forum,
    constraint forum_users_forum_nickname_key
        unique (forum, nickname)
);

create or replace function add_post_user() returns trigger as
$$
DECLARE
    author_nickname CITEXT;
    author_fullname TEXT;
    author_about    TEXT;
    author_email    CITEXT;
BEGIN
    SELECT nickname, fullname, about, email
    FROM users
    WHERE nickname = NEW.author
    INTO author_nickname, author_fullname, author_about, author_email;

    INSERT INTO nickname_forum (nickname, fullname, about, email, forum)
    VALUES (author_nickname, author_fullname, author_about, author_email, NEW.forum)
    ON CONFLICT DO NOTHING;

    RETURN NEW;
END
$$
    LANGUAGE 'plpgsql';

CREATE TRIGGER add_post_user
    AFTER INSERT
    ON post
    FOR EACH ROW
EXECUTE PROCEDURE add_post_user();



create or replace function add_thread_user() returns trigger
as
$$
DECLARE
    author_nickname CITEXT;
    author_fullname TEXT;
    author_about    TEXT;
    author_email    CITEXT;
BEGIN
    SELECT nickname, fullname, about, email
    FROM users
    WHERE nickname = NEW.author
    INTO author_nickname, author_fullname, author_about, author_email;

    INSERT INTO nickname_forum (nickname, fullname, about, email, forum)
    VALUES (author_nickname, author_fullname, author_about, author_email, NEW.forum)
    ON CONFLICT DO NOTHING;

    RETURN NEW;
END
$$
    LANGUAGE 'plpgsql';

CREATE TRIGGER add_thread_user
    AFTER INSERT
    ON thread
    FOR EACH ROW
EXECUTE PROCEDURE add_thread_user();

CREATE INDEX IF NOT EXISTS for_user_nickname ON users USING hash (nickname);
CREATE INDEX IF NOT EXISTS for_user_email ON users USING hash (email);
CREATE INDEX IF NOT EXISTS for_forum_slug ON forum USING hash (slug);
CREATE INDEX IF NOT EXISTS for_thread_slug ON thread USING hash (slug);
CREATE INDEX IF NOT EXISTS for_thread_forum ON thread USING hash (forum);
CREATE INDEX IF NOT EXISTS for_post_thread ON post USING hash (thread);
CREATE INDEX IF NOT EXISTS for_thread_created ON thread (created);
CREATE INDEX IF NOT EXISTS for_thread_created_forum ON thread (forum, created);
CREATE INDEX IF NOT EXISTS for_post_path_single ON post ((path[1]));
CREATE INDEX IF NOT EXISTS for_post_id_path_single on post (id, (path[1]));
CREATE INDEX IF NOT EXISTS for_post_path ON post (path);
CREATE UNIQUE INDEX IF NOT EXISTS for_votes_nickname_thread_nickname on votes (thread, nickname);
CREATE INDEX for_nickname_forum ON nickname_forum USING hash (nickname);
CREATE INDEX for_nickname_forum_nickname ON nickname_forum (forum, nickname);

vacuum;
vacuum analyze;
