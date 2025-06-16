DO $$

    BEGIN

    create extension if not exists "uuid-ossp";

    create table if not exists attendance_periods (
        id uuid default uuid_generate_v4() not null primary key,
        start_date timestamp with time zone not null,
        end_date timestamp with time zone not null,
        created_by uuid null,
        updated_by uuid null,
        created_at timestamp with time zone not null default current_timestamp,
        updated_at timestamp with time zone not null default current_timestamp,
        deleted_at timestamp with time zone null
                                 );

    create index if not exists idx_attendance_periods_deleted_at on attendance_periods (deleted_at);

END $$;