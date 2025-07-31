create table if not exists notifications (
    id serial primary key,
    user_id integer not null,
    from_id integer not null,
    
    type varchar(50) not null,
    message text not null,
    created_at timestamp with time zone default current_timestamp,
    read boolean default false
);