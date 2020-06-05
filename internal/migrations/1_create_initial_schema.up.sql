 
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS scraps
(
    id SERIAL PRIMARY KEY,
    uuid uuid DEFAULT uuid_generate_v4 (),
    refCode VARCHAR NOT NULL,
    fileUri VARCHAR NOT NULL,
    CONSTRAINT scraps_refCode_unique UNIQUE(refCode),
    CONSTRAINT scraps_uuid_unique UNIQUE(uuid),
    CONSTRAINT scraps_fileUri_unique UNIQUE(fileUri)
);