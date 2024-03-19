
create table actors
(
    id      serial
        primary key,
    name    varchar(50),
    gender  varchar(15),
    birhday varchar(15)
);

create table films
(
    id          serial
        primary key,
    name        varchar(150),
    description varchar(1000),
    date        varchar(8),
    rating      integer
);

create table actors_films
(
    actor_id integer not null
        references actors
            on update cascade,
    film_id  integer not null
        references films
            on update cascade on delete cascade,
    primary key (actor_id, film_id)
);

create table users
(
    id      serial
        primary key,
    role    varchar(15),
    api_key varchar(150)
);


INSERT INTO public.users (id, role, api_key) VALUES (1, 'admin', 'root');

