# mysql-benchmark

## sequence table schema
Database : test_seq

```
CREATE DATABASE test_seq;
use test_seq;
CREATE TABLE IF NOT EXISTS sequence (
  id bigint(20) unsigned DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE utf8mb4_bin ROW_FORMAT=DYNAMIC;
INSERT IGNORE INTO sequence VALUES (0);
```

### GRANT
```
GRANT ALL PRIVILEGES ON test_seq.* TO root@'127.0.0.1';
```