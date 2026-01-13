-- Создание таблицы shorter
create table if not exists public.shorter (
    url_id varchar(6) primary key unique ,
    original_url text not null
);
