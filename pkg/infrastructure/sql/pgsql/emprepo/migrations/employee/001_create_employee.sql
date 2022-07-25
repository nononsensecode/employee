-- This is a sample migration.

create table employees(
  id serial primary key,
  name varchar not null,
  age smallint not null
);

---- create above / drop below ----

drop table employees;

