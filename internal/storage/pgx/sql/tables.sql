CREATE TABLE IF NOT EXISTS leads
(

    id     SERIAL PRIMARY KEY,
    name   TEXT        not null,
    email  TEXT unique not null,
    phone  TEXT unique not null,
    source TEXT        not null


)