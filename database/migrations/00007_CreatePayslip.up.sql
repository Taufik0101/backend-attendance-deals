DO $$

BEGIN

    create extension if not exists "uuid-ossp";

    create table if not exists payslips (
        id uuid default uuid_generate_v4() not null primary key,
        user_id uuid constraint fk_payslip_users references users,
        payroll_id uuid constraint fk_payslip_payroll references payrolls,
        attendance_days integer not null,
        overtime_hours integer not null,
        reimbursement decimal not null,
        base_salary decimal not null,
        total_pay decimal not null,
        created_by uuid null,
        updated_by uuid null,
        created_at timestamp with time zone not null default current_timestamp,
        updated_at timestamp with time zone not null default current_timestamp,
        deleted_at timestamp with time zone null
                                 );

create index if not exists idx_payslips_deleted_at on payslips (deleted_at);

END $$;