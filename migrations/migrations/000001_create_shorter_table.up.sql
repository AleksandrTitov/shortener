-- Создание таблицы shorter
create table if not exists public.shorter (
    url_id varchar(6) primary key unique ,
    short_url text not null unique,
    original_url text not null,
    created_at timestamp default now()
);
