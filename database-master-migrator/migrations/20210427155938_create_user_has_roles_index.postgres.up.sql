CREATE INDEX user_has_roles_user_id_index ON public.user_has_roles USING btree (user_id);
CREATE INDEX user_has_roles_role_id_index ON public.user_has_roles USING btree (role_id);

ALTER TABLE public.user_has_roles ADD CONSTRAINT user_has_roles_user_id_foreign FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE public.user_has_roles ADD CONSTRAINT user_has_roles_role_id_foreign FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE;