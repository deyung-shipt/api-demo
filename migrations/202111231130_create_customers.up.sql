create table if not exists customers (
                                         id bigserial primary key,
                                         email varchar not null,
                                         first_name varchar not null,
                                         last_name varchar not null,
                                         created_at timestamp with time zone not null,
                                         updated_at timestamp with time zone not null
);
create unique index customers_email_idx on customers (email);