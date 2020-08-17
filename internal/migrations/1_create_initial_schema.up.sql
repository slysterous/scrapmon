
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE status AS ENUM ('pending', 'ongoing', 'success','failure');
CREATE TABLE IF NOT EXISTS screenshots
(
    id SERIAL PRIMARY KEY,
    uuid uuid DEFAULT uuid_generate_v4 (),
    refCode VARCHAR NOT NULL,
    codeCreatedAt DATE NOT NULL,
    fileUri VARCHAR,
    downloadStatus status DEFAULT 'pending',
    CONSTRAINT scraps_refCode_unique UNIQUE(refCode),
    CONSTRAINT scraps_uuid_unique UNIQUE(uuid),
    CONSTRAINT scraps_fileUri_unique UNIQUE(fileUri)
);