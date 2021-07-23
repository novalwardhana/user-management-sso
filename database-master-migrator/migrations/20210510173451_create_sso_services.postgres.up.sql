CREATE TABLE public.sso_services (
    id serial NOT NULL,
    email varchar(80) NOT NULL,
    unique_code varchar(80) NOT NULL,
    domain varchar(80) NOT NULL,
    "secret" varchar(80) NOT NULL,
    created_at timestamp with time zone NULL,
    updated_at timestamp with time zone NULL,
    deleted_at timestamp with time zone NULL,
    CONSTRAINT sso_services_pkey PRIMARY KEY (id),
    CONSTRAINT sso_services_domain_unique UNIQUE (domain)
);