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

-- Database: identity

CREATE DATABASE identity
    WITH
    OWNER = postgres
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.UTF-8'
    LC_CTYPE = 'en_US.UTF-8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1
    IS_TEMPLATE = False;


-- Table: public.roles
/*
Role types:
    Unspecified = 0
    System      = 1
    Service     = 2

Role statuses:
    Unspecified = 0
    New         = 1
    Active      = 2
    Inactive    = 3
    Deleting    = 4
    Deleted     = 5
*/
CREATE TABLE IF NOT EXISTS public.roles
(
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    name character varying(256) COLLATE pg_catalog."default" NOT NULL,
    type smallint NOT NULL,
    title character varying(256) COLLATE pg_catalog."default" NOT NULL,
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
    CONSTRAINT roles_pkey PRIMARY KEY (id),
    CONSTRAINT roles_type_check CHECK (type = 1 OR type = 2),
    CONSTRAINT roles_status_check CHECK (status >= 1 AND status <= 5)
)
TABLESPACE pg_default;

CREATE UNIQUE INDEX IF NOT EXISTS roles_name_idx
    ON public.roles (lower(name))
    WHERE status <> 5;

CREATE INDEX IF NOT EXISTS roles_type_idx ON public.roles (type);
CREATE INDEX IF NOT EXISTS roles_created_at_idx ON public.roles (created_at);
CREATE INDEX IF NOT EXISTS roles_updated_at_idx ON public.roles (updated_at);
CREATE INDEX IF NOT EXISTS roles_status_idx ON public.roles (status);
CREATE INDEX IF NOT EXISTS roles_status_updated_at_idx ON public.roles (status_updated_at);
CREATE INDEX IF NOT EXISTS roles_app_id_idx ON public.roles (app_id);
CREATE INDEX IF NOT EXISTS roles_app_group_id_idx ON public.roles (app_group_id);

-- Table: public.role_info

CREATE TABLE IF NOT EXISTS public.role_info
(
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    role_id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    created_by bigint NOT NULL,
    updated_at timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    updated_by bigint NOT NULL,
    is_deleted boolean NOT NULL DEFAULT false,
    deleted_at timestamp(6) without time zone,
    deleted_by bigint,
    active_assignment_count bigint NOT NULL,
    existing_assignment_count bigint NOT NULL,
    _version_stamp bigint NOT NULL,
    _timestamp timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    CONSTRAINT role_info_pkey PRIMARY KEY (id),
    CONSTRAINT role_info_role_id_key UNIQUE (role_id),
    CONSTRAINT role_info_role_id_fkey FOREIGN KEY (role_id)
        REFERENCES public.roles (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE RESTRICT,
    CONSTRAINT role_info_active_assignment_count_check CHECK (active_assignment_count >= 0),
    CONSTRAINT role_info_existing_assignment_count_check CHECK (existing_assignment_count >= 0)
)
TABLESPACE pg_default;

CREATE INDEX IF NOT EXISTS role_info_created_at_idx ON public.role_info (created_at);
CREATE INDEX IF NOT EXISTS role_info_updated_at_idx ON public.role_info (updated_at);
CREATE INDEX IF NOT EXISTS role_info_is_deleted_idx ON public.role_info (is_deleted);
CREATE INDEX IF NOT EXISTS role_info_deleted_at_idx ON public.role_info (deleted_at);

-- Table: public.new_role_assignments

CREATE TABLE IF NOT EXISTS public.new_role_assignments
(
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    operation_id uuid NOT NULL,
    role_id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    created_by bigint NOT NULL,
    _version_stamp bigint NOT NULL,
    _timestamp timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    CONSTRAINT new_role_assignments_pkey PRIMARY KEY (id),
    CONSTRAINT new_role_assignments_operation_id_key UNIQUE (operation_id)
)
TABLESPACE pg_default;

-- Table: public.group_role_assignments
/*
Group role assignment statuses:
    Unspecified = 0
    New         = 1
    Active      = 2
    Inactive    = 3
    Deleting    = 4
    Deleted     = 5
*/
CREATE TABLE IF NOT EXISTS public.group_role_assignments
(
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    role_assignment_id bigint NOT NULL,
    "group" smallint NOT NULL,
    role_id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    created_by bigint NOT NULL,
    updated_at timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    updated_by bigint NOT NULL,
    status smallint NOT NULL,
    status_updated_at timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    status_updated_by bigint NOT NULL,
    status_comment text COLLATE pg_catalog."default",
    _version_stamp bigint NOT NULL,
    _timestamp timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    CONSTRAINT group_role_assignments_pkey PRIMARY KEY (id),
    CONSTRAINT group_role_assignments_role_assignment_id_key UNIQUE (role_assignment_id),
    CONSTRAINT group_role_assignments_status_check CHECK (status >= 1 AND status <= 5)
)
TABLESPACE pg_default;

CREATE UNIQUE INDEX IF NOT EXISTS group_role_assignments_group_role_id_idx
    ON public.group_role_assignments ("group", role_id)
    WHERE status <> 5;

CREATE INDEX IF NOT EXISTS group_role_assignments_group_idx ON public.group_role_assignments ("group");
CREATE INDEX IF NOT EXISTS group_role_assignments_role_id_idx ON public.group_role_assignments (role_id);
CREATE INDEX IF NOT EXISTS group_role_assignments_created_at_idx ON public.group_role_assignments (created_at);
CREATE INDEX IF NOT EXISTS group_role_assignments_updated_at_idx ON public.group_role_assignments (updated_at);
CREATE INDEX IF NOT EXISTS group_role_assignments_status_idx ON public.group_role_assignments (status);
CREATE INDEX IF NOT EXISTS group_role_assignments_status_updated_at_idx ON public.group_role_assignments (status_updated_at);

-- Table: public.permission_groups
/*
Permission group statuses:
    Unspecified = 0
    New         = 1
    Active      = 2
    Inactive    = 3
    Deleted     = 4
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
    description text COLLATE pg_catalog."default" NOT NULL,
    _version_stamp bigint NOT NULL,
    _timestamp timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    CONSTRAINT permission_groups_pkey PRIMARY KEY (id),
    CONSTRAINT permission_groups_name_key UNIQUE (name)
)
TABLESPACE pg_default;

-- Table: public.permissions
/*
Permission statuses:
    Unspecified = 0
    New         = 1
    Active      = 2
    Inactive    = 3
    Deleted     = 4
*/
CREATE TABLE IF NOT EXISTS public.permissions
(
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    group_id bigint NOT NULL,
    name character varying(256) COLLATE pg_catalog."default" NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    created_by bigint NOT NULL,
    updated_at timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    updated_by bigint NOT NULL,
    status smallint NOT NULL,
    status_updated_at timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    status_updated_by bigint NOT NULL,
    status_comment text COLLATE pg_catalog."default",
    description text COLLATE pg_catalog."default" NOT NULL,
    _version_stamp bigint NOT NULL,
    _timestamp timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    CONSTRAINT permissions_pkey PRIMARY KEY (id),
    CONSTRAINT permissions_name_key UNIQUE (name),
    CONSTRAINT permissions_group_id_fkey FOREIGN KEY (group_id)
        REFERENCES public.permission_groups (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE RESTRICT
)
TABLESPACE pg_default;
