create table todo_items
(
    id          int auto_increment
        primary key,
    title       varchar(255)                                                not null,
    image       json                                                        null,
    description text                                                        null,
    status      enum ('Doing', 'Done', 'Deleted') default 'Doing'           null,
    created_at  datetime                          default CURRENT_TIMESTAMP null,
    updated_at  datetime                          default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP,
    user_id     int                                                         null,
    liked_count int                               default 0                 null,
    constraint todo_items_users_id_fk
        foreign key (user_id) references users (id)
);

create table users
(
    id         int auto_increment
        primary key,
    email      varchar(100)                       not null,
    password   varchar(100)                       not null,
    salt       varchar(50)                        null,
    last_name  varchar(100)                       not null,
    first_name varchar(100)                       null,
    phone      varchar(10)                        null,
    status     int      default 1                 null,
    role       varchar(10)                        null,
    created_at datetime default CURRENT_TIMESTAMP null,
    updated_at datetime default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP,
    constraint users_email_uindex
        unique (email)
);

create table user_like_items
(
    user_id    int                                null,
    item_id    int                                null,
    created_at datetime default CURRENT_TIMESTAMP null,
    constraint user_like_items_user_id_item_id_uindex
        unique (user_id, item_id)
);

create index user_like_items_item_id_index
    on user_like_items (item_id);

create index user_like_items_user_id_index
    on user_like_items (user_id);

