CREATE OR REPLACE FUNCTION fn_log_audit_changes() RETURNS TRIGGER AS $$
DECLARE
  user_id UUID := null;
  req_id TEXT := null;
BEGIN
  BEGIN
    user_id := current_setting('app.current_user', true)::UUID;
  EXCEPTION WHEN OTHERS THEN
    user_id := null;
  END;

  BEGIN
    req_id := current_setting('app.request_id', true);
  EXCEPTION WHEN OTHERS THEN
    req_id := null;
  END;

  INSERT INTO audit_logs (
    table_name,
    record_id,
    action,
    changed_by,
    ip_address,
    request_id,
    old_data,
    new_data,
    created_at
  )
  VALUES (
    TG_TABLE_NAME,
    COALESCE(NEW.id, OLD.id),
    TG_OP,
    user_id,
    inet_client_addr(),
    req_id,
    CASE WHEN TG_OP IN ('UPDATE', 'DELETE') THEN to_jsonb(OLD) ELSE NULL END,
    CASE WHEN TG_OP IN ('INSERT', 'UPDATE') THEN to_jsonb(NEW) ELSE NULL END,
    now()
  );

  RETURN CASE
    WHEN TG_OP = 'DELETE' THEN OLD
    ELSE NEW
  END;
END;
$$ LANGUAGE plpgsql;

-- Users
CREATE TRIGGER trg_audit_users
AFTER INSERT OR UPDATE OR DELETE ON users
FOR EACH ROW EXECUTE FUNCTION fn_log_audit_changes();

-- Attendance
CREATE TRIGGER trg_audit_attendances
AFTER INSERT OR UPDATE OR DELETE ON attendances
FOR EACH ROW EXECUTE FUNCTION fn_log_audit_changes();

-- Overtimes
CREATE TRIGGER trg_audit_overtimes
AFTER INSERT OR UPDATE OR DELETE ON overtimes
FOR EACH ROW EXECUTE FUNCTION fn_log_audit_changes();

-- Reimbursements
CREATE TRIGGER trg_audit_reimbursements
AFTER INSERT OR UPDATE OR DELETE ON reimbursements
FOR EACH ROW EXECUTE FUNCTION fn_log_audit_changes();

-- Payrolls
CREATE TRIGGER trg_audit_payrolls
AFTER INSERT OR UPDATE OR DELETE ON payrolls
FOR EACH ROW EXECUTE FUNCTION fn_log_audit_changes();

-- Payslips
CREATE TRIGGER trg_audit_payslips
AFTER INSERT OR UPDATE OR DELETE ON payslips
FOR EACH ROW EXECUTE FUNCTION fn_log_audit_changes();

-- Payroll Summaries
CREATE TRIGGER trg_audit_payroll_summaries
AFTER INSERT OR UPDATE OR DELETE ON payroll_summaries
FOR EACH ROW EXECUTE FUNCTION fn_log_audit_changes();

-- Attendance Periods
CREATE TRIGGER trg_audit_attendance_periods
AFTER INSERT OR UPDATE OR DELETE ON attendance_periods
FOR EACH ROW EXECUTE FUNCTION fn_log_audit_changes();