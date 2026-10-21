UPDATE tasks
SET position = numbered.position
FROM (
  SELECT id, ROW_NUMBER() OVER (ORDER BY created_at, id) AS position
  FROM tasks
) AS numbered
WHERE tasks.id = numbered.id;
