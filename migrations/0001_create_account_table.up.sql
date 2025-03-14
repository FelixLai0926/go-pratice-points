CREATE TABLE IF NOT EXISTS public.account (
    user_id BIGINT NOT NULL PRIMARY KEY,
    available_balance NUMERIC(18,2) NOT NULL DEFAULT 0,
    reserved_balance NUMERIC(18,2) NOT NULL DEFAULT 0,
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
    CHECK (available_balance >= 0),
    CHECK (reserved_balance >= 0)
);