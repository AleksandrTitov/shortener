-- Создание столбца deleted_flag
alter table public.shorter add column deleted_flag bool default false;

-- Обновляем существующие записи
update public.shorter set deleted_flag = false where deleted_flag is null;