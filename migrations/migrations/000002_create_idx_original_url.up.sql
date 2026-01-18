-- Создание индекса для original_url
create unique index idx_original_url on public.shorter (
    original_url
);
