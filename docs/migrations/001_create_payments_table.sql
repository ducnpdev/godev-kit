-- Create payments table
CREATE TABLE IF NOT EXISTS payments (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'VND',
    payment_type VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    meter_number VARCHAR(50) NOT NULL,
    customer_code VARCHAR(50) NOT NULL,
    description TEXT,
    transaction_id VARCHAR(100) UNIQUE NOT NULL,
    payment_method VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create payment_history table
CREATE TABLE IF NOT EXISTS payment_history (
    id BIGSERIAL PRIMARY KEY,
    payment_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    payment_type VARCHAR(20) NOT NULL,
    meter_number VARCHAR(50) NOT NULL,
    customer_code VARCHAR(50) NOT NULL,
    description TEXT,
    transaction_id VARCHAR(100) NOT NULL,
    payment_method VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_payments_user_id ON payments(user_id);
CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status);
CREATE INDEX IF NOT EXISTS idx_payments_created_at ON payments(created_at);
CREATE INDEX IF NOT EXISTS idx_payments_transaction_id ON payments(transaction_id);

CREATE INDEX IF NOT EXISTS idx_payment_history_payment_id ON payment_history(payment_id);
CREATE INDEX IF NOT EXISTS idx_payment_history_user_id ON payment_history(user_id);
CREATE INDEX IF NOT EXISTS idx_payment_history_status ON payment_history(status);
CREATE INDEX IF NOT EXISTS idx_payment_history_created_at ON payment_history(created_at);

-- Add foreign key constraints
ALTER TABLE payments ADD CONSTRAINT fk_payments_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE payment_history ADD CONSTRAINT fk_payment_history_payment_id FOREIGN KEY (payment_id) REFERENCES payments(id) ON DELETE CASCADE;
ALTER TABLE payment_history ADD CONSTRAINT fk_payment_history_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE; 