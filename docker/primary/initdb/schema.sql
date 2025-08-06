-- ==============================================================
-- 1. Create replication role
-- ==============================================================

DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname='replicationuser') THEN
    CREATE ROLE replicationuser
      WITH REPLICATION LOGIN
      PASSWORD 'admin';
  ELSE
    ALTER ROLE replicationuser
      WITH ENCRYPTED PASSWORD 'admin';
  END IF;
END
$$;


