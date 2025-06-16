DO $$

    BEGIN

    create extension if not exists "uuid-ossp";

    create table if not exists audit_logs (
        id uuid default uuid_generate_v4() not null primary key,
        user_id uuid constraint fk_audit_log_users references users,
        period_id uuid constraint fk_attendance_period references attendance_periods,
        action varchar(255) not null,
        entity varchar(255) not null,
        entity_id uuid not null,
        ip_address varchar(255) not null,
        request_id varchar(255) not null,
        created_by uuid null,
        updated_by uuid null,
        created_at timestamp with time zone not null default current_timestamp,
        updated_at timestamp with time zone not null default current_timestamp,
        deleted_at timestamp with time zone null
                                 );

create index if not exists idx_audit_logs_deleted_at on audit_logs (deleted_at);

END $$;