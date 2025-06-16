DO $$

BEGIN

    create extension if not exists "uuid-ossp";

    create table if not exists overtimes (
        id uuid default uuid_generate_v4() not null primary key,
        user_id uuid constraint fk_overtime_user references users,
        period_id uuid constraint fk_attendance_period references attendance_periods,
        date timestamp with time zone not null,
        hours integer not null,
        ip_address varchar(255) not null,
        created_by uuid null,
        updated_by uuid null,
        created_at timestamp with time zone not null default current_timestamp,
        updated_at timestamp with time zone not null default current_timestamp,
        deleted_at timestamp with time zone null
                           );

create index if not exists idx_overtimes_deleted_at on overtimes (deleted_at);

END $$;