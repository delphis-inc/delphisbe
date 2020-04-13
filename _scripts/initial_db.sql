CREATE TABLE discussions (
    id varchar(32) PRIMARY KEY,
    created_at timestamp with time zone CONSTRAINT not null,
    updated_at timestamp with time zone CONSTRAINT not null,
    deleted_at timestamp with time zone,
    title varchar(256) CONSTRAINT not null,
    anonymity_type varchar(32) not null,
    moderator_id varchar(32)
    
)