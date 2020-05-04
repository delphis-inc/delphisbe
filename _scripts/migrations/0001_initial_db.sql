CREATE TABLE IF NOT EXISTS discussions (
    id varchar(36) PRIMARY KEY,
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone not null,
    deleted_at timestamp with time zone,
    title varchar(256) not null,
    anonymity_type varchar(36) not null,
    moderator_id varchar(36)
);

CREATE TABLE IF NOT EXISTS moderators (
    id varchar(36) PRIMARY KEY,
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone not null,
    deleted_at timestamp with time zone,
    user_profile_id varchar(36) not null
);

CREATE TABLE IF NOT EXISTS participants (
    id varchar(36) PRIMARY KEY,
    participant_id int not null,
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone not null,
    deleted_at timestamp with time zone,
    discussion_id varchar(36) not null,
    viewer_id varchar(36) not null,
    user_id varchar(36)
);

CREATE TABLE IF NOT EXISTS post_contents (
    id varchar(36) PRIMARY KEY,
    content text
);

CREATE TABLE IF NOT EXISTS posts (
    id varchar(36) PRIMARY KEY,
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone not null,
    deleted_at timestamp with time zone,
    deleted_reason_code varchar(16),
    discussion_id varchar(36),
    participant_id varchar(36),
    post_content_id varchar(36) not null
);

CREATE TABLE IF NOT EXISTS social_infos (
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone not null,
    deleted_at timestamp with time zone,
    user_profile_id varchar(36),
    network varchar(16) not null,
    access_token varchar(256),
    access_token_secret varchar(256),
    user_id varchar(36),
    profile_image_url text,
    screen_name varchar(64),
    is_verified boolean,
    PRIMARY KEY (network, user_id)
);

CREATE TABLE IF NOT EXISTS user_profiles (
    id varchar(36) PRIMARY KEY,
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone not null,
    deleted_at timestamp with time zone,
    display_name varchar(128) not null,
    user_id varchar(36),
    twitter_handle varchar(36)
);

CREATE TABLE IF NOT EXISTS users (
    id varchar(36) PRIMARY KEY,
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone not null,
    deleted_at timestamp with time zone   
);

CREATE TABLE IF NOT EXISTS viewers (
    id varchar(36) PRIMARY KEY,
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone not null,
    deleted_at timestamp with time zone,
    last_viewed timestamp with time zone,
    last_viewed_post_id varchar(36),
    discussion_id varchar(36),
    user_id varchar(36)
);