DO $$
    DECLARE
admin_id UUID;

BEGIN
DROP TYPE IF EXISTS role_types;
CREATE TYPE role_types AS enum ('admin', 'employee');

create extension if not exists "uuid-ossp";
        CREATE EXTENSION IF NOT EXISTS pgcrypto;

create table if not exists users (
    id uuid default uuid_generate_v4() not null primary key,
    username varchar(255) not null,
    password text null,
    role    role_types                    NOT NULL DEFAULT 'employee',
    salary integer null,
    created_by uuid null,
    updated_by uuid null,
    created_at timestamp with time zone not null default current_timestamp,
    updated_at timestamp with time zone not null default current_timestamp,
    deleted_at timestamp with time zone null
                             );

create index if not exists idx_users_deleted_at on users (deleted_at);

INSERT INTO users(id, username, password, role, created_at, updated_at)
VALUES (
           gen_random_uuid(),
           'admin',
           crypt('password123', gen_salt('bf')),
           'admin',
           NOW(),
           NOW()
       ) RETURNING id INTO admin_id;

        FOR i IN 1..100 LOOP
                INSERT INTO users (id, username, password, role, created_by, updated_by, created_at, updated_at, salary)
                VALUES (
                    gen_random_uuid(),
                    'employee' || LPAD(i::text, 3, '0'),
                    crypt('password123', gen_salt('bf')),
                    'employee',
                    admin_id,
                    admin_id,
                    NOW(),
                    NOW(),
                    ROUND((3000000 + random() * 2000000)::numeric, 2)
                );
    END LOOP;
END $$;