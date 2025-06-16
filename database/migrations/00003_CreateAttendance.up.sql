DO $$

    BEGIN

    create extension if not exists "uuid-ossp";

    create table if not exists attendances (
        id uuid default uuid_generate_v4() not null primary key,
        user_id uuid constraint fk_attendance_user references users,
        period_id uuid constraint fk_attendance_period references attendance_periods,
        date_in timestamp with time zone not null,
        date_out timestamp with time zone null,
        ip_address varchar(255) not null,
        created_by uuid null,
        updated_by uuid null,
        created_at timestamp with time zone not null default current_timestamp,
        updated_at timestamp with time zone not null default current_timestamp,
        deleted_at timestamp with time zone null
                                 );

    create index if not exists idx_attendances_deleted_at on attendances (deleted_at);

END $$;