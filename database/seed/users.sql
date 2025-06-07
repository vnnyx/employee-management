DO $$
BEGIN
  FOR i IN 1..100 LOOP
    INSERT INTO users (
      username,
      password,
      is_admin,
      salary,
      created_at,
      updated_at
    ) VALUES (
      format('employee_%s', i),
      crypt('password123', gen_salt('bf')),
      FALSE,
      trunc(random() * 5000000 + 3000000)::NUMERIC(12, 2),
      now(),
      now()
    );
  END LOOP;
END;
$$;

INSERT INTO users (
  username,
  password,
  is_admin,
  salary,
  created_at,
  updated_at
) VALUES (
  'admin',
  crypt('adminpass123', gen_salt('bf')),
  TRUE,
  10000000.00,
  now(),
  now()
);
