CREATE TABLE IF NOT EXISTS public.trade_records (
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