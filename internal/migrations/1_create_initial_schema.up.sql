
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE status AS ENUM ('pending', 'ongoing','notfound','success','failure');
CREATE TABLE IF NOT EXISTS scraps
(
    id SERIAL PRIMARY KEY,
    uuid uuid DEFAULT uuid_generate_v4 (),
    refCode VARCHAR NOT NULL,
    codeCreatedAt TIMESTAMP NOT NULL,
    fileUri VARCHAR,
    downloadStatus status DEFAULT 'pending',
    CONSTRAINT screenshots_refCode_unique UNIQUE(refCode),
    CONSTRAINT screenshots_uuid_unique UNIQUE(uuid)
);