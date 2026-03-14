create or replace function set_updated_at() returns trigger as $$
begin
    new.updated_at = now();
    return new;
end;
$$ language plpgsql;

create table if not exists events (
    id         uuid primary key,
    name       text not null,
    metadata   jsonb not null,
    content    jsonb not null,
    status     text not null default 'draft',
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

drop trigger if exists trg_events_updated_at on events;
create trigger trg_events_updated_at
before update on events
for each row
execute function set_updated_at();

create table if not exists users (
    id uuid primary key,
    email text not null unique,
    password_hash text not null,
    password_salt text not null,
    role text not null
);