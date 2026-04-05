-- Users Table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'Viewer', -- Viewer, Analyst, Admin
    status TEXT NOT NULL DEFAULT 'Active', -- Active, Inactive
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE -- For soft delete
);

-- Financial Records Table
CREATE TABLE IF NOT EXISTS financial_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    amount NUMERIC(15, 2) NOT NULL,
    type TEXT NOT NULL, -- Income, Expense
    category TEXT NOT NULL,
    date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE -- For soft delete
);

-- Role Permissions Table (Dynamic RBAC)
CREATE TABLE IF NOT EXISTS role_permissions (
    role TEXT NOT NULL,
    permission TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (role, permission)
);

-- Index for soft delete queries
CREATE UNIQUE INDEX idx_users_username_active ON users(username) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_deleted_at ON users(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_financial_records_deleted_at ON financial_records(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_financial_records_user_id ON financial_records(user_id);
CREATE INDEX idx_role_permissions_role ON role_permissions(role);

-- Initial Permissions
INSERT INTO role_permissions (role, permission) VALUES
('Admin', 'user::create'),
('Admin', 'user::update'),
('Admin', 'user::delete'),
('Admin', 'user::list'),
('Admin', 'record::create'),
('Admin', 'record::update'),
('Admin', 'record::delete'),
('Admin', 'record::list'),
('Admin', 'summary::view'),
('Analyst', 'record::list'),
('Analyst', 'summary::view'),
('Viewer', 'record::list'),
('Viewer', 'summary::view')
ON CONFLICT (role, permission) DO NOTHING;
