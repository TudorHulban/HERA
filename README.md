# HERA - Layer for RDBMS
(SQLite) Schema

```sql
drop table users

create table users (
id integer primary key,
first_name text NOT NULL,
last_name text NOT NULL,
role integer NOT NULL)
```

```sql
insert into users(first_name, last_name, role) values("john", "doe",1)
```
