CREATE EXTENSION IF NOT EXISTS pgcrypto;
SET intervalstyle = 'iso_8601';

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    is_admin BOOLEAN DEFAULT FALSE,
    salary NUMERIC(12, 2) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    created_by UUID,
    updated_by UUID,
    ip_address TEXT
);

CREATE TABLE attendance_periods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    UNIQUE(start_date, end_date),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    created_by UUID,
    updated_by UUID,
    ip_address TEXT
);

CREATE TABLE attendances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    attendance_date DATE NOT NULL,
    UNIQUE(user_id, attendance_date),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    created_by UUID,
    updated_by UUID,
    ip_address TEXT
);

CREATE TABLE overtimes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    overtime_date DATE NOT NULL,
    overtime_hours INTERVAL NOT NULL,
    UNIQUE(user_id, overtime_date),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    created_by UUID,
    updated_by UUID,
    ip_address TEXT
);

CREATE TABLE reimbursements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    amount NUMERIC(12, 2) NOT NULL,
    description TEXT,
    reimbursement_date DATE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    created_by UUID,
    updated_by UUID,
    ip_address TEXT
);

CREATE TABLE payrolls (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    period_id UUID REFERENCES attendance_periods(id) UNIQUE,
    run_by UUID REFERENCES users(id),
    run_at TIMESTAMPTZ DEFAULT now(),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    created_by UUID,
    updated_by UUID,
    ip_address TEXT
);

CREATE TABLE payslips (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    payroll_id UUID REFERENCES payrolls(id),
    base_salary NUMERIC(12, 2),
    attendance_days INTEGER,
    overtime_hours INTERVAL,
    overtime_pay NUMERIC(12, 2),
    reimbursement_total NUMERIC(12, 2),
    total_take_home NUMERIC(12, 2),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    created_by UUID,
    updated_by UUID,
    ip_address TEXT
);

CREATE TABLE payroll_summaries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payroll_id UUID REFERENCES payrolls(id),
    total_take_home NUMERIC(14, 2),
    generated_by UUID REFERENCES users(id),
    generated_at TIMESTAMPTZ DEFAULT now(),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    created_by UUID,
    updated_by UUID,
    ip_address TEXT
);

CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    table_name TEXT NOT NULL,
    record_id UUID NOT NULL,
    action TEXT NOT NULL,
    changed_by UUID REFERENCES users(id),
    ip_address TEXT,
    request_id TEXT,
    old_data JSONB,
    new_data JSONB,
    created_at TIMESTAMPTZ DEFAULT now()
);

