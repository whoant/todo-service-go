create table todo_items(
    id int primary key auto_increment,
    title varchar(255) not null,
    image json,
    description text,
    status enum('Doing', 'Done', 'Deleted') default 'Doing',
    created_at datetime default now(),
    updated_at datetime default now() on update current_timestamp
);

insert into todo_items (title, image, description, status)
values ('Code', null, 'Mo ta ne', 'Ds');

EXPLAIN
select *
from todo_items
where title = 'sss';