CREATE TABLE public.users (
	id serial NOT NULL,
	name varchar(80) NOT NULL,
    username varchar(80) NOT NULL,
	email varchar(255) NOT NULL,
	"password" varchar(255) NOT NULL,
    is_active bool NOT NULL DEFAULT true,
	created_at timestamp with time zone NULL,
	updated_at timestamp with time zone NULL,
	deleted_at timestamp with time zone NULL,
    CONSTRAINT users_pkey PRIMARY KEY (id),
    CONSTRAINT users_username_unique UNIQUE (username),
	CONSTRAINT users_email_unique UNIQUE (email)
);

CREATE TYPE public.role_group AS ENUM (
    'superadmin',
    'admin',
    'operator',
    'user',
    'guest'
);

CREATE TABLE public.roles (
	id serial NOT NULL,
    code varchar(255) NOT NULL,
	"name" varchar(255) NOT NULL,
    "group" public.role_group NOT NULL,
    description text NULL,
	created_at timestamp with time zone NULL,
	updated_at timestamp with time zone NULL,
	deleted_at timestamp with time zone NULL,
	CONSTRAINT roles_pkey PRIMARY KEY (id),
    CONSTRAINT roles_code_unique UNIQUE (code),
    CONSTRAINT roles_name_unique UNIQUE (name)
);

CREATE TABLE public.permissions (
	id serial NOT NULL,
    code varchar(255) NOT NULL,
	"name" varchar(255) NOT NULL,
    description text NULL,
	created_at timestamp with time zone NULL,
	updated_at timestamp with time zone NULL,
	deleted_at timestamp with time zone NULL,
	CONSTRAINT permissions_pkey PRIMARY KEY (id),
    CONSTRAINT permissions_code_unique UNIQUE (code),
    CONSTRAINT permissions_name_unique UNIQUE (name)
);