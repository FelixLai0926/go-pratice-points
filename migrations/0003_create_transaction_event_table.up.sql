CREATE TABLE IF NOT EXISTS public.transaction_event (
    id SERIAL PRIMARY KEY,
    transaction_id UUID NOT NULL,
    event_type VARCHAR(50) NOT NULL,
    payload JSONB,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_transactions FOREIGN KEY (transaction_id) REFERENCES public.trade_records(transaction_id)
);
