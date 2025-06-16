DO $$

BEGIN

    create extension if not exists "uuid-ossp";

    create table if not exists payrolls (
        id uuid default uuid_generate_v4() not null primary key,
        period_id uuid constraint fk_payrol_attendance_period references attendance_periods,
        created_by uuid null,
        updated_by uuid null,
        created_at timestamp with time zone not null default current_timestamp,
        updated_at timestamp with time zone not null default current_timestamp,
        deleted_at timestamp with time zone null
                                 );

create index if not exists idx_payrolls_deleted_at on payrolls (deleted_at);

END $$;