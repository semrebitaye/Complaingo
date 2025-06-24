ALTER TABLE complaints
  DROP CONSTRAINT complaints_status_check,
  ADD CONSTRAINT complaints_status_check CHECK (
    status IN ('Accepted', 'Resolved', 'Rejected', 'Created')
  );
