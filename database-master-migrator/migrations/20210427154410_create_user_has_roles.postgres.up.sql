CREATE TABLE public.user_has_roles (
    id serial NOT NULL,
    user_id int4 NULL,
    role_id int4 NULL,
    created_at timestamp with time zone NULL,
    updated_at timestamp with time zone NULL,
    CONSTRAINT user_has_roles_pkey PRIMARY KEY (id)
)