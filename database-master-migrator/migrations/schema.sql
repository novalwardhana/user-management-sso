--
-- PostgreSQL database dump
--

-- Dumped from database version 11.10 (Ubuntu 11.10-1.pgdg20.04+1)
-- Dumped by pg_dump version 11.10 (Ubuntu 11.10-1.pgdg20.04+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


--
-- Name: role_group; Type: TYPE; Schema: public; Owner: noval
--

CREATE TYPE public.role_group AS ENUM (
    'superadmin',
    'admin',
    'operator',
    'user',
    'guest'
);


ALTER TYPE public.role_group OWNER TO noval;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: permissions; Type: TABLE; Schema: public; Owner: noval
--

CREATE TABLE public.permissions (
    id integer NOT NULL,
    code character varying(255) NOT NULL,
    name character varying(255) NOT NULL,
    description text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.permissions OWNER TO noval;

--
-- Name: permissions_id_seq; Type: SEQUENCE; Schema: public; Owner: noval
--

CREATE SEQUENCE public.permissions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.permissions_id_seq OWNER TO noval;

--
-- Name: permissions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: noval
--

ALTER SEQUENCE public.permissions_id_seq OWNED BY public.permissions.id;


--
-- Name: role_has_permissions; Type: TABLE; Schema: public; Owner: noval
--

CREATE TABLE public.role_has_permissions (
    id integer NOT NULL,
    role_id integer,
    permission_id integer,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.role_has_permissions OWNER TO noval;

--
-- Name: role_has_permissions_id_seq; Type: SEQUENCE; Schema: public; Owner: noval
--

CREATE SEQUENCE public.role_has_permissions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.role_has_permissions_id_seq OWNER TO noval;

--
-- Name: role_has_permissions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: noval
--

ALTER SEQUENCE public.role_has_permissions_id_seq OWNED BY public.role_has_permissions.id;


--
-- Name: roles; Type: TABLE; Schema: public; Owner: noval
--

CREATE TABLE public.roles (
    id integer NOT NULL,
    code character varying(255) NOT NULL,
    name character varying(255) NOT NULL,
    "group" public.role_group NOT NULL,
    description text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.roles OWNER TO noval;

--
-- Name: roles_id_seq; Type: SEQUENCE; Schema: public; Owner: noval
--

CREATE SEQUENCE public.roles_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.roles_id_seq OWNER TO noval;

--
-- Name: roles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: noval
--

ALTER SEQUENCE public.roles_id_seq OWNED BY public.roles.id;


--
-- Name: schema_migration; Type: TABLE; Schema: public; Owner: noval
--

CREATE TABLE public.schema_migration (
    version character varying(14) NOT NULL
);


ALTER TABLE public.schema_migration OWNER TO noval;

--
-- Name: sso_services; Type: TABLE; Schema: public; Owner: noval
--

CREATE TABLE public.sso_services (
    id integer NOT NULL,
    email character varying(80) NOT NULL,
    unique_code character varying(80) NOT NULL,
    domain character varying(80) NOT NULL,
    secret character varying(80) NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.sso_services OWNER TO noval;

--
-- Name: sso_services_id_seq; Type: SEQUENCE; Schema: public; Owner: noval
--

CREATE SEQUENCE public.sso_services_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.sso_services_id_seq OWNER TO noval;

--
-- Name: sso_services_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: noval
--

ALTER SEQUENCE public.sso_services_id_seq OWNED BY public.sso_services.id;


--
-- Name: user_has_roles; Type: TABLE; Schema: public; Owner: noval
--

CREATE TABLE public.user_has_roles (
    id integer NOT NULL,
    user_id integer,
    role_id integer,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.user_has_roles OWNER TO noval;

--
-- Name: user_has_roles_id_seq; Type: SEQUENCE; Schema: public; Owner: noval
--

CREATE SEQUENCE public.user_has_roles_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.user_has_roles_id_seq OWNER TO noval;

--
-- Name: user_has_roles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: noval
--

ALTER SEQUENCE public.user_has_roles_id_seq OWNED BY public.user_has_roles.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: noval
--

CREATE TABLE public.users (
    id integer NOT NULL,
    name character varying(80) NOT NULL,
    username character varying(80) NOT NULL,
    email character varying(255) NOT NULL,
    password character varying(255) NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    user_uuid uuid
);


ALTER TABLE public.users OWNER TO noval;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: noval
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO noval;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: noval
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: permissions id; Type: DEFAULT; Schema: public; Owner: noval
--

ALTER TABLE ONLY public.permissions ALTER COLUMN id SET DEFAULT nextval('public.permissions_id_seq'::regclass);


--
-- Name: role_has_permissions id; Type: DEFAULT; Schema: public; Owner: noval
--

ALTER TABLE ONLY public.role_has_permissions ALTER COLUMN id SET DEFAULT nextval('public.role_has_permissions_id_seq'::regclass);


--
-- Name: roles id; Type: DEFAULT; Schema: public; Owner: noval
--

ALTER TABLE ONLY public.roles ALTER COLUMN id SET DEFAULT nextval('public.roles_id_seq'::regclass);


--
-- Name: sso_services id; Type: DEFAULT; Schema: public; Owner: noval
--

ALTER TABLE ONLY public.sso_services ALTER COLUMN id SET DEFAULT nextval('public.sso_services_id_seq'::regclass);


--
-- Name: user_has_roles id; Type: DEFAULT; Schema: public; Owner: noval
--

ALTER TABLE ONLY public.user_has_roles ALTER COLUMN id SET DEFAULT nextval('public.user_has_roles_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: noval
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: permissions permissions_code_unique; Type: CONSTRAINT; Schema: public; Owner: noval
--

ALTER TABLE ONLY public.permissions
    ADD CONSTRAINT permissions_code_unique UNIQUE (code);


--
-- Name: permissions permissions_name_unique; Type: CONSTRAINT; Schema: public; Owner: noval
--

ALTER TABLE ONLY public.permissions
    ADD CONSTRAINT permissions_name_unique UNIQUE (name);


--
-- Name: permissions permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: noval
--

ALTER TABLE ONLY public.permissions
    ADD CONSTRAINT permissions_pkey PRIMARY KEY (id);


--
-- Name: role_has_permissions role_has_permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: noval
--

ALTER TABLE ONLY public.role_has_permissions
    ADD CONSTRAINT role_has_permissions_pkey PRIMARY KEY (id);


--
-- Name: roles roles_code_unique; Type: CONSTRAINT; Schema: public; Owner: noval
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_code_unique UNIQUE (code);


--
-- Name: roles roles_name_unique; Type: CONSTRAINT; Schema: public; Owner: noval
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_name_unique UNIQUE (name);


--
-- Name: roles roles_pkey; Type: CONSTRAINT; Schema: public; Owner: noval
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_pkey PRIMARY KEY (id);


--
-- Name: sso_services sso_services_domain_unique; Type: CONSTRAINT; Schema: public; Owner: noval
--

ALTER TABLE ONLY public.sso_services
    ADD CONSTRAINT sso_services_domain_unique UNIQUE (domain);


--
-- Name: sso_services sso_services_pkey; Type: CONSTRAINT; Schema: public; Owner: noval
--

ALTER TABLE ONLY public.sso_services
    ADD CONSTRAINT sso_services_pkey PRIMARY KEY (id);


--
-- Name: user_has_roles user_has_roles_pkey; Type: CONSTRAINT; Schema: public; Owner: noval
--

ALTER TABLE ONLY public.user_has_roles
    ADD CONSTRAINT user_has_roles_pkey PRIMARY KEY (id);


--
-- Name: users users_email_unique; Type: CONSTRAINT; Schema: public; Owner: noval
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_unique UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: noval
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users users_username_unique; Type: CONSTRAINT; Schema: public; Owner: noval
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_unique UNIQUE (username);


--
-- Name: role_has_permissions_permission_id_index; Type: INDEX; Schema: public; Owner: noval
--

CREATE INDEX role_has_permissions_permission_id_index ON public.role_has_permissions USING btree (permission_id);


--
-- Name: role_has_permissions_role_id_index; Type: INDEX; Schema: public; Owner: noval
--

CREATE INDEX role_has_permissions_role_id_index ON public.role_has_permissions USING btree (role_id);


--
-- Name: schema_migration_version_idx; Type: INDEX; Schema: public; Owner: noval
--

CREATE UNIQUE INDEX schema_migration_version_idx ON public.schema_migration USING btree (version);


--
-- Name: user_has_roles_role_id_index; Type: INDEX; Schema: public; Owner: noval
--

CREATE INDEX user_has_roles_role_id_index ON public.user_has_roles USING btree (role_id);


--
-- Name: user_has_roles_user_id_index; Type: INDEX; Schema: public; Owner: noval
--

CREATE INDEX user_has_roles_user_id_index ON public.user_has_roles USING btree (user_id);


--
-- Name: role_has_permissions role_has_permissions_permission_id_foreign; Type: FK CONSTRAINT; Schema: public; Owner: noval
--

ALTER TABLE ONLY public.role_has_permissions
    ADD CONSTRAINT role_has_permissions_permission_id_foreign FOREIGN KEY (permission_id) REFERENCES public.permissions(id) ON DELETE CASCADE;


--
-- Name: role_has_permissions role_has_permissions_role_id_foreign; Type: FK CONSTRAINT; Schema: public; Owner: noval
--

ALTER TABLE ONLY public.role_has_permissions
    ADD CONSTRAINT role_has_permissions_role_id_foreign FOREIGN KEY (role_id) REFERENCES public.roles(id) ON DELETE CASCADE;


--
-- Name: user_has_roles user_has_roles_role_id_foreign; Type: FK CONSTRAINT; Schema: public; Owner: noval
--

ALTER TABLE ONLY public.user_has_roles
    ADD CONSTRAINT user_has_roles_role_id_foreign FOREIGN KEY (role_id) REFERENCES public.roles(id) ON DELETE CASCADE;


--
-- Name: user_has_roles user_has_roles_user_id_foreign; Type: FK CONSTRAINT; Schema: public; Owner: noval
--

ALTER TABLE ONLY public.user_has_roles
    ADD CONSTRAINT user_has_roles_user_id_foreign FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

