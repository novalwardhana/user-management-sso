CREATE TABLE public.role_has_permissions (
    id serial NOT NULL,
    role_id int4 NULL,
    permission_id int4 NULL,
    created_at timestamp with time zone NULL,
    updated_at timestamp with time zone NULL,
    CONSTRAINT role_has_permissions_pkey PRIMARY KEY (id)
);

CREATE INDEX role_has_permissions_role_id_index ON public.role_has_permissions USING btree (role_id);
CREATE INDEX role_has_permissions_permission_id_index ON public.role_has_permissions USING btree (permission_id);

ALTER TABLE public.role_has_permissions ADD CONSTRAINT role_has_permissions_role_id_foreign FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE;
ALTER TABLE public.role_has_permissions ADD CONSTRAINT role_has_permissions_permission_id_foreign FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE;