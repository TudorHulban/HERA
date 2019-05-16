# HERA - Layer for RDBMS
(SQLite) Schema - needs SQLite support (database file created at first execution, ex. db.Exec(DDL)):
```bash
sudo apt-get install sqlite3 libsqlite3-dev
```

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

```sql
drop table roles

create table roles (
id integer primary key,
code text not null,
description text not null,
enabled text not null)
```

```sql
insert into roles(code, description, enabled) values("ADMIN", "Full rights", "Y");
insert into roles(code, description, enabled) values("USER", "Some rights", "Y");
insert into roles(code, description, enabled) values("GUEST", "Few rights", "Y");
```

