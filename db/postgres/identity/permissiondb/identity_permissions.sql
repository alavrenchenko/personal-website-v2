-- Copyright 2023 Alexey Lavrenchenko. All rights reserved.
--
-- Licensed under the Apache License, Version 2.0 (the "License");
-- you may not use this file except in compliance with the License.
-- You may obtain a copy of the License at
--
-- 	http://www.apache.org/licenses/LICENSE-2.0
--
-- Unless required by applicable law or agreed to in writing, software
-- distributed under the License is distributed on an "AS IS" BASIS,
-- WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
-- See the License for the specific language governing permissions and
-- limitations under the License.

-- Database: identity_permissions

CREATE DATABASE identity_permissions
    WITH
    OWNER = postgres
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.UTF-8'
    LC_CTYPE = 'en_US.UTF-8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1
    IS_TEMPLATE = False;


-- Table: public.permission_groups
/*
Permission group statuses:
    Unspecified = 0
    New         = 1
    Active      = 2
    Inactive    = 3
    Deleting    = 4
    Deleted     = 5
*/
CREATE TABLE IF NOT EXISTS public.permission_groups
(
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    name character varying(256) COLLATE pg_catalog."default" NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    created_by bigint NOT NULL,
    updated_at timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    updated_by bigint NOT NULL,
    status smallint NOT NULL,
    status_updated_at timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    status_updated_by bigint NOT NULL,
    status_comment text COLLATE pg_catalog."default",
    app_id bigint,
    app_group_id bigint,
    description text COLLATE pg_catalog."default" NOT NULL,
    _version_stamp bigint NOT NULL,
    _timestamp timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    CONSTRAINT permission_groups_pkey PRIMARY KEY (id),
    CONSTRAINT permission_groups_status_check CHECK (status >= 1 AND status <= 5)
)
TABLESPACE pg_default;

CREATE UNIQUE INDEX IF NOT EXISTS permission_groups_name_idx
    ON public.permission_groups (name)
    WHERE status <> 5;

CREATE UNIQUE INDEX IF NOT EXISTS permission_groups_name_lc_idx
    ON public.permission_groups (lower(name))
    WHERE status <> 5;

CREATE INDEX IF NOT EXISTS permission_groups_created_at_idx ON public.permission_groups (created_at);
CREATE INDEX IF NOT EXISTS permission_groups_updated_at_idx ON public.permission_groups (updated_at);
CREATE INDEX IF NOT EXISTS permission_groups_status_idx ON public.permission_groups (status);
CREATE INDEX IF NOT EXISTS permission_groups_status_updated_at_idx ON public.permission_groups (status_updated_at);
CREATE INDEX IF NOT EXISTS permission_groups_app_id_idx ON public.permission_groups (app_id);
CREATE INDEX IF NOT EXISTS permission_groups_app_group_id_idx ON public.permission_groups (app_group_id);

-- Table: public.permissions
/*
Permission statuses:
    Unspecified = 0
    New         = 1
    Active      = 2
    Inactive    = 3
    Deleting    = 4
    Deleted     = 5
*/
CREATE TABLE IF NOT EXISTS public.permissions
(
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    name character varying(256) COLLATE pg_catalog."default" NOT NULL,
    group_id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    created_by bigint NOT NULL,
    updated_at timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    updated_by bigint NOT NULL,
    status smallint NOT NULL,
    status_updated_at timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    status_updated_by bigint NOT NULL,
    status_comment text COLLATE pg_catalog."default",
    app_id bigint,
    app_group_id bigint,
    description text COLLATE pg_catalog."default" NOT NULL,
    _version_stamp bigint NOT NULL,
    _timestamp timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    CONSTRAINT permissions_pkey PRIMARY KEY (id),
    CONSTRAINT permissions_group_id_fkey FOREIGN KEY (group_id)
        REFERENCES public.permission_groups (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE RESTRICT,
    CONSTRAINT permissions_status_check CHECK (status >= 1 AND status <= 5)
)
TABLESPACE pg_default;

CREATE UNIQUE INDEX IF NOT EXISTS permissions_name_idx
    ON public.permissions (name)
    WHERE status <> 5;

CREATE UNIQUE INDEX IF NOT EXISTS permissions_name_lc_idx
    ON public.permissions (lower(name))
    WHERE status <> 5;

CREATE INDEX IF NOT EXISTS permissions_created_at_idx ON public.permissions (created_at);
CREATE INDEX IF NOT EXISTS permissions_updated_at_idx ON public.permissions (updated_at);
CREATE INDEX IF NOT EXISTS permissions_status_idx ON public.permissions (status);
CREATE INDEX IF NOT EXISTS permissions_status_updated_at_idx ON public.permissions (status_updated_at);
CREATE INDEX IF NOT EXISTS permissions_app_id_idx ON public.permissions (app_id);
CREATE INDEX IF NOT EXISTS permissions_app_group_id_idx ON public.permissions (app_group_id);

CREATE TABLE IF NOT EXISTS public.role_permissions
(
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    role_id bigint NOT NULL,
    permission_id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    created_by bigint NOT NULL,
    is_deleted boolean NOT NULL DEFAULT FALSE,
    deleted_at timestamp(6) without time zone,
    deleted_by bigint,
    _version_stamp bigint NOT NULL,
    _timestamp timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    CONSTRAINT role_permissions_pkey PRIMARY KEY (id),
    CONSTRAINT role_permissions_permission_id_fkey FOREIGN KEY (permission_id)
        REFERENCES public.permissions (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE RESTRICT
)
TABLESPACE pg_default;

CREATE UNIQUE INDEX IF NOT EXISTS role_permissions_role_id_permission_id_idx
    ON public.role_permissions (role_id, permission_id)
    WHERE is_deleted IS FALSE;

CREATE INDEX IF NOT EXISTS role_permissions_created_at_idx ON public.role_permissions (created_at);
CREATE INDEX IF NOT EXISTS role_permissions_is_deleted_idx ON public.role_permissions (is_deleted);
CREATE INDEX IF NOT EXISTS role_permissions_deleted_at_idx ON public.role_permissions (deleted_at);

CREATE TABLE IF NOT EXISTS public.deleted_role_permissions
(
    id bigint NOT NULL,
    role_id bigint NOT NULL,
    permission_id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    created_by bigint NOT NULL,
    is_deleted boolean NOT NULL,
    deleted_at timestamp(6) without time zone,
    deleted_by bigint,
    _version_stamp bigint NOT NULL,
    _timestamp timestamp(6) without time zone NOT NULL,
    CONSTRAINT deleted_role_permissions_pkey PRIMARY KEY (id),
    CONSTRAINT deleted_role_permissions_permission_id_fkey FOREIGN KEY (permission_id)
        REFERENCES public.permissions (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE RESTRICT,
    CONSTRAINT deleted_role_permissions_is_deleted_check CHECK (is_deleted IS TRUE)
)
TABLESPACE pg_default;

CREATE INDEX IF NOT EXISTS deleted_role_permissions_created_at_idx ON public.deleted_role_permissions (created_at);
CREATE INDEX IF NOT EXISTS deleted_role_permissions_deleted_at_idx ON public.deleted_role_permissions (deleted_at);
CREATE INDEX IF NOT EXISTS deleted_role_permissions_role_id_permission_id_idx ON public.deleted_role_permissions (role_id, permission_id);
