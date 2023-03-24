-- SQLite
SELECT * FROM sqlite_schema
WHERE type='table'
ORDER BY name;
PRAGMA table_info(user_info);
.schema;
SELECT * FROM sqlite_master
WHERE type='table'
ORDER BY name;
