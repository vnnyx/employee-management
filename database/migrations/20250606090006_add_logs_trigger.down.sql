DROP TRIGGER IF EXISTS trg_audit_attendance_periods ON attendance_periods;
DROP TRIGGER IF EXISTS trg_audit_payroll_summaries ON payroll_summaries;
DROP TRIGGER IF EXISTS trg_audit_payslips ON payslips;
DROP TRIGGER IF EXISTS trg_audit_payrolls ON payrolls;
DROP TRIGGER IF EXISTS trg_audit_reimbursements ON reimbursements;
DROP TRIGGER IF EXISTS trg_audit_overtimes ON overtimes;
DROP TRIGGER IF EXISTS trg_audit_attendances ON attendances;
DROP TRIGGER IF EXISTS trg_audit_users ON users;

DROP FUNCTION IF EXISTS fn_log_audit_changes;