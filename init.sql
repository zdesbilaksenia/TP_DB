CREATE
    EXTENSION IF NOT EXISTS citext;

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
    created  timestamp with time zone default now()
);

CREATE
    OR REPLACE FUNCTION add_post() RETURNS TRIGGER AS
$$
BEGIN
    UPDATE forum
    SET posts = posts + 1
    WHERE Slug = NEW.forum;
    NEW
        .
        path
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
    slug    citext unique,
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

create function update_vote() returns trigger as
$$
BEGIN
    IF
        OLD.voice <> NEW.voice THEN
        UPDATE thread
        SET votes=(votes + NEW.Voice * 2)
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