-- Создание столбца user_id
alter table public.shorter add column user_id uuid;

-- Создание индекса для user_id
create unique index idx_user_id on public.shorter (
    user_id
);
