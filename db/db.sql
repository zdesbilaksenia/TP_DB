create extension if not exists citext;

drop table if exists thread, forum, users, votes, post, nickname_forum cascade;

create unlogged table users
(
    nickname citext collate "C" not null
        constraint users_pk primary key,
    fullname text,
    about    text,
    email    citext unique
);

create unlogged table forum
(
    title   text,
    "user"  citext
        constraint forum_users_nickname_fk references users,
    slug    citext not null
        constraint forum_pk primary key,
    posts   bigint default 0,
    threads bigint default 0
);

create unlogged table post
(
    id       bigserial
        constraint post_pkey primary key,
    parent   bigint                   default 0,
    author   citext
        constraint post_users_nickname_fk references users,
    message  text,
    isedited boolean                  default false,
    forum    citext
        constraint post_forum_slug_fk references forum,
    thread   integer,
    created  timestamp with time zone default now(),
    path     bigint[]                 default array []::bigint[]
);

create or replace function add_post() returns trigger as
$$
begin
    update forum
    set posts = posts + 1
    where slug = new.forum;
    new.path = (select path from post where id = new.parent limit 1) || new.id;
    return new;
end
$$ language 'plpgsql';

create trigger add_post
    before insert
    on post
    for each row
execute procedure add_post();

create unlogged table thread
(
    id      serial
        constraint thread_pk primary key,
    title   text,
    author  citext
        constraint thread_users_nickname_fk references users,
    forum   citext
        constraint thread_forum_slug_fk references forum,
    message text,
    votes   integer                  default 0,
    slug    citext,
    created timestamp with time zone default now()
);

create or replace function add_thread() returns trigger as
$$
begin
    update forum
    set threads = threads + 1
    where new.forum = slug;
    return null;
end
$$ language 'plpgsql';

create trigger add_thread
    after insert
    on thread
    for each row
execute procedure add_thread();

create unlogged table votes
(
    id       bigserial
        constraint votes_pkey
            primary key,
    nickname citext
        constraint votes_nickname_fkey
            references users,
    voice    integer,
    thread   integer not null
        constraint votes_thread_fkey references thread,
    constraint votes_thread_nickname_key
        unique (thread, nickname)
);

create or replace function add_vote() returns trigger as
$$
begin
    update thread
    set votes=(votes + new.voice)
    where id = new.thread;
    return new;
end
$$ language 'plpgsql';

create trigger add_vote
    after insert
    on votes
    for each row
execute procedure add_vote();

create or replace function update_vote() returns trigger as
$$
begin
    if old.voice <> new.voice then
        update thread
        set votes = votes - old.voice + new.voice
        where id = new.thread;
    end if;
    return new;
end
$$ language 'plpgsql';

create trigger update_vote
    after update
    on votes
    for each row
execute procedure update_vote();

create unlogged table nickname_forum
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
declare
    author_nickname citext;
    author_fullname text;
    author_about    text;
    author_email    citext;
begin
    select nickname, fullname, about, email
    from users
    where nickname = new.author
    into author_nickname, author_fullname, author_about, author_email;

    insert into nickname_forum (nickname, fullname, about, email, forum)
    values (author_nickname, author_fullname, author_about, author_email, new.forum)
    on conflict do nothing;

    return new;
end
$$ language 'plpgsql';

create trigger add_post_user
    after insert
    on post
    for each row
execute procedure add_post_user();

create or replace function add_thread_user() returns trigger as
$$
declare
    author_nickname citext;
    author_fullname text;
    author_about    text;
    author_email    citext;
begin
    select nickname, fullname, about, email
    from users
    where nickname = new.author
    into author_nickname, author_fullname, author_about, author_email;

    insert into nickname_forum (nickname, fullname, about, email, forum)
    values (author_nickname, author_fullname, author_about, author_email, new.forum)
    on conflict do nothing;

    return new;
end
$$ language 'plpgsql';

create trigger add_thread_user
    after insert
    on thread
    for each row
execute procedure add_thread_user();

create index if not exists for_user_nickname on users using hash (nickname);
create index if not exists for_user_email on users using hash (email);
create index if not exists for_forum_slug on forum using hash (slug);
create index if not exists for_thread_slug on thread using hash (slug);
create index if not exists for_thread_forum on thread using hash (forum);
create index if not exists for_post_thread on post using hash (thread);
create index if not exists for_thread_created on thread (created);
create index if not exists for_thread_created_forum on thread (forum, created);
create index if not exists for_post_path_single on post ((path[1]));
create index if not exists for_post_id_path_single on post (id, (path[1]));
create index if not exists for_post_path on post (path);
create unique index if not exists for_votes_nickname_thread_nickname on votes (thread, nickname);
create index if not exists for_nickname_forum on nickname_forum using hash (nickname);
create index if not exists for_nickname_forum_nickname on nickname_forum (forum, nickname);

vacuum;
vacuum analyze;
