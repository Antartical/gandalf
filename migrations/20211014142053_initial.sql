-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE apps_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1;
CREATE TABLE "public"."apps" (
    "id" bigint DEFAULT nextval('apps_id_seq') NOT NULL,
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    "uuid" uuid DEFAULT uuid_generate_v4(),
    "client_id" uuid DEFAULT uuid_generate_v4(),
    "client_secret" text NOT NULL,
    "name" text NOT NULL,
    "icon_uri" text,
    "redirect_uris" text[],
    "user_id" bigint,
    CONSTRAINT "apps_client_id_key" UNIQUE ("client_id"),
    CONSTRAINT "apps_pkey" PRIMARY KEY ("id"),
    CONSTRAINT "apps_uuid_key" UNIQUE ("uuid")
) WITH (oids = false);
CREATE INDEX "app_client_id" ON "public"."apps" USING btree ("client_id");
CREATE INDEX "app_uuid" ON "public"."apps" USING btree ("uuid");
CREATE INDEX "idx_apps_deleted_at" ON "public"."apps" USING btree ("deleted_at");


CREATE TABLE "public"."user_has_signin_on_app" (
    "user_id" bigint NOT NULL,
    "app_id" bigint NOT NULL,
    CONSTRAINT "user_has_signin_on_app_pkey" PRIMARY KEY ("user_id", "app_id")
) WITH (oids = false);


CREATE SEQUENCE users_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1;
CREATE TABLE "public"."users" (
    "id" bigint DEFAULT nextval('users_id_seq') NOT NULL,
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    "last_login" timestamptz,
    "uuid" uuid DEFAULT uuid_generate_v4(),
    "email" text NOT NULL,
    "password" text NOT NULL,
    "name" text NOT NULL,
    "surname" text NOT NULL,
    "birthday" timestamptz NOT NULL,
    "verified" boolean DEFAULT false,
    "phone" text,
    CONSTRAINT "users_email_key" UNIQUE ("email"),
    CONSTRAINT "users_pkey" PRIMARY KEY ("id"),
    CONSTRAINT "users_uuid_key" UNIQUE ("uuid")
) WITH (oids = false);

CREATE INDEX "idx_users_deleted_at" ON "public"."users" USING btree ("deleted_at");
CREATE INDEX "usr_email" ON "public"."users" USING btree ("email");
CREATE INDEX "usr_uuid" ON "public"."users" USING btree ("uuid");


ALTER TABLE ONLY "public"."apps" ADD CONSTRAINT "fk_users_apps" FOREIGN KEY (user_id) REFERENCES users(id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."user_has_signin_on_app" ADD CONSTRAINT "fk_user_has_signin_on_app_app" FOREIGN KEY (app_id) REFERENCES apps(id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."user_has_signin_on_app" ADD CONSTRAINT "fk_user_has_signin_on_app_user" FOREIGN KEY (user_id) REFERENCES users(id) NOT DEFERRABLE;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "apps";
DROP SEQUENCE IF EXISTS apps_id_seq;
DROP TABLE IF EXISTS "user_has_signin_on_app";
DROP TABLE IF EXISTS "users";
DROP SEQUENCE IF EXISTS users_id_seq;
-- +goose StatementEnd
