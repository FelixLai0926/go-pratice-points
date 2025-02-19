CREATE TABLE IF NOT EXISTS public.account (
    user_id BIGINT NOT NULL PRIMARY KEY,
    available_balance NUMERIC(18,2) NOT NULL DEFAULT 0,
    reserved_balance NUMERIC(18,2) NOT NULL DEFAULT 0,
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS public.transaction (
    transaction_id UUID NOT NULL DEFAULT gen_random_uuid(),
    nonce BIGINT NOT NULL,
    from_account_id BIGINT NOT NULL,
    to_account_id BIGINT NOT NULL,
    amount NUMERIC(18,2) NOT NULL,
    status INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_transaction_id UNIQUE (transaction_id),
    CONSTRAINT fk_from_account FOREIGN KEY (from_account_id) REFERENCES public.account(user_id),
    CONSTRAINT fk_to_account FOREIGN KEY (to_account_id) REFERENCES public.account(user_id),
    PRIMARY KEY (from_account_id, nonce)
);

CREATE TABLE IF NOT EXISTS public.transaction_event (
    id SERIAL PRIMARY KEY,
    transaction_id UUID NOT NULL,
    event_type VARCHAR(50) NOT NULL,
    payload JSONB,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_transaction FOREIGN KEY (transaction_id) REFERENCES public.transaction(transaction_id)
);

INSERT INTO public.account(
user_id, available_balance, reserved_balance)
VALUES (1, 1000, 0);
