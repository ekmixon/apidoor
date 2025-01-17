BEGIN;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS public.apiinfo
(
    id serial primary key ,
    name text NOT NULL,
    source text NOT NULL,
    description text NOT NULL,
    thumbnail text NOT NULL,
    swagger_url text NOT NULL
)
WITH (
    OIDS = FALSE
);

COMMENT ON TABLE public.apiinfo
    IS 'Store information of products(API).';

CREATE TABLE IF NOT EXISTS public.apiuser
(
    id serial primary key ,
    account_id VARCHAR(32) not null unique,
    email_address TEXT not null,
    login_password_hash TEXT not null,  /* pgcryptoのcrypt関数を使用 */
    name TEXT,
    permission_flag VARCHAR(2) not null default '00',
    created_at TIMESTAMP,
    updated_at TIMESTAMP
)
WITH (
    OIDS = FALSE
);

COMMENT ON TABLE public.apiuser
    IS 'Store management-api users.';

CREATE TABLE IF NOT EXISTS public.log_list
(
    id serial primary key ,
    run_date timestamp with time zone NOT NULL,
    api_key text NOT NULL,
    api_path text NOT NULL,
    custom_log jsonb
);

COMMENT ON TABLE public.log_list
    IS 'Table to save log of gateway.';
END;


