DROP TABLE IF EXISTS public.users;
CREATE TABLE public.users
(
    created_at  timestamp with time zone DEFAULT timezone('UTC'::text, now()) NOT NULL,
    modified_at timestamp with time zone DEFAULT timezone('UTC'::text, now()) NOT NULL,
    id          serial                                                        NOT NULL,
    name        text                                                          NOT NULL
);
