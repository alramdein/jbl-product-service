CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name  VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(100) NOT NULL, -- since one email can be either generator or contirbutor
    password TEXT NOT NULL,
    role_id UUID NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    FOREIGN KEY (role_id) REFERENCES roles(id)
);

CREATE TABLE referral_links (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    generator_id UUID NOT NULL,
    code VARCHAR(100) NOT NULL UNIQUE,
    expired_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    FOREIGN KEY (generator_id) REFERENCES users(id)
);

CREATE TABLE contributions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    referral_link_id UUID NOT NULL,
    contributor_id UUID NOT NULL,
    accessed_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    FOREIGN KEY (referral_link_id) REFERENCES referral_links(id),
    FOREIGN KEY (contributor_id) REFERENCES users(id)
);

INSERT INTO roles
(id, "name", created_at, updated_at, deleted_at) VALUES
('61638ab8-f516-42d4-b5a3-d95e6de69684'::uuid, 'generator', timezone('utc', now()), timezone('utc', now()), NULL),
('bc5fa404-a035-4664-a4ab-10cdc347cf64'::uuid, 'contributor', timezone('utc', now()), timezone('utc', now()), NULL);
