# db

## Dump

```shell
pg_dump <db_name> --inserts -U <username> --create --no-comments > dump.sql
```

```shell
mysqldump <db_name> -h 127.0.0.1 -P 3001 -u <username> -p --insert-ignore=TRUE > dump_prod.sql
```

## Remove DB

### New

```sql
DROP DATABASE "name" WITH (FORCE);
```

### old

```sql
ALTER DATABASE <db_name> CONNECTION LIMIT 0;
SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = '<db_name>';
DROP database <db_name>;
```

## Check extensions

```sql
SELECT *
FROM pg_extension;
SELECT * FROM pg_available_extensions;
```

## Add extension

```sql
create EXTENSION <extension>;
```

## Current number of connections

```sql
SELECT sum(numbackends) FROM pg_stat_database;
```

## Full dump flow

### Create dump

```shell
pg_dump -Z0 --inserts --host <host> -U <username <db> --no-comments --no-acl --no-privileges --no-owner -F d -b -v -f <dump_file>
tar -cf - <dump_file> | pigz -p 10 > <dump_file>.tar.gz
scp -i ${HOME}/.ssh/prod_main_ssh_key_id_rsa <dump_file>.tar.gz <username>@<host>:/bks/
mkdir -p <dump_file> && pigz -p 10 -dc <dump_file>.tar.gz | tar -C <dump_file> --strip-components 1 -xf -
pg_restore -h localhost -U <username> -d <db> -j 10 -v -F d <dump_file>
```

### Restore

```shell
pg_restore -h <host> -U <username> -d <db> -j 10 -v -F d <dump_file>
```
