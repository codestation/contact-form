-- +migrate Up
-- +migrate StatementBegin
create or replace function notify_event() returns trigger as $$
begin
    if (tg_op = 'INSERT') then
        perform pg_notify('contactform.newtask', row_to_json(NEW)::text);
    end if;

    return null;
end;
$$ language plpgsql;
-- +migrate StatementEnd

-- +migrate Down
drop function if exists notify_event();
