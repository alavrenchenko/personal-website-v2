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

-- Database: identity_users

CREATE DATABASE identity_users
    WITH
    OWNER = postgres
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.UTF-8'
    LC_CTYPE = 'en_US.UTF-8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1
    IS_TEMPLATE = False;


-- Table: public.users
/*
User types:
    Unspecified = 0
    User        = 1
    SystemUser  = 2

User statuses:
    Unspecified          = 0
    New                  = 1
    PendingApproval      = 2
    Active               = 3
    LockedOut            = 4
    TemporarilyLockedOut = 5
    Disabled             = 6
    Deleting             = 7
    Deleted              = 8
*/
CREATE TABLE IF NOT EXISTS public.users
(
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    name character varying(256) COLLATE pg_catalog."default",
    type smallint NOT NULL,
    "group" bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    created_by bigint NOT NULL,
    updated_at timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    updated_by bigint NOT NULL,
    status smallint NOT NULL,
    status_updated_at timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    status_updated_by bigint NOT NULL,
    status_comment text COLLATE pg_catalog."default",
    email character varying(512) COLLATE pg_catalog."default",
    first_sign_in_time timestamp(6) without time zone,
    first_sign_in_ip character varying(64) COLLATE pg_catalog."default",
    last_sign_in_time timestamp(6) without time zone,
    last_sign_in_ip character varying(64) COLLATE pg_catalog."default",
    last_sign_out_time timestamp(6) without time zone,
    last_activity_time timestamp(6) without time zone,
    last_activity_ip character varying(64) COLLATE pg_catalog."default",
    _version_stamp bigint NOT NULL,
    _timestamp timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    CONSTRAINT users_pkey PRIMARY KEY (id),
    CONSTRAINT users_type_check CHECK (type = 1 OR type = 2),
    CONSTRAINT users_status_check CHECK (status >= 1 AND status <= 8)
)
TABLESPACE pg_default;

CREATE UNIQUE INDEX IF NOT EXISTS users_name_idx
    ON public.users (lower(name))
    WHERE status <> 8;

CREATE UNIQUE INDEX IF NOT EXISTS users_email_idx
    ON public.users (lower(email))
    WHERE status <> 8;

CREATE INDEX IF NOT EXISTS users_type_idx ON public.users (type);
CREATE INDEX IF NOT EXISTS users_group_idx ON public.users ("group");
CREATE INDEX IF NOT EXISTS users_created_at_idx ON public.users (created_at);
CREATE INDEX IF NOT EXISTS users_updated_at_idx ON public.users (updated_at);
CREATE INDEX IF NOT EXISTS users_status_idx ON public.users (status);
CREATE INDEX IF NOT EXISTS users_status_updated_at_idx ON public.users (status_updated_at);
CREATE INDEX IF NOT EXISTS users_first_sign_in_time_idx ON public.users (first_sign_in_time);
CREATE INDEX IF NOT EXISTS users_first_sign_in_ip_idx ON public.users (first_sign_in_ip);
CREATE INDEX IF NOT EXISTS users_last_sign_in_time_idx ON public.users (last_sign_in_time);
CREATE INDEX IF NOT EXISTS users_last_sign_in_ip_idx ON public.users (last_sign_in_ip);
CREATE INDEX IF NOT EXISTS users_last_activity_time_idx ON public.users (last_activity_time);
CREATE INDEX IF NOT EXISTS users_last_activity_ip_idx ON public.users (last_activity_ip);

-- Table: public.personal_info
/*
User genders:
    Unspecified = 0
    Unknown     = 1
    Female      = 2
    Male        = 3
    Other       = 4
*/
CREATE TABLE IF NOT EXISTS public.personal_info
(
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    user_id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    created_by bigint NOT NULL,
    updated_at timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    updated_by bigint NOT NULL,
    is_deleted boolean NOT NULL DEFAULT FALSE,
    deleted_at timestamp(6) without time zone,
    deleted_by bigint,
    first_name character varying(512) COLLATE pg_catalog."default" NOT NULL,
    last_name character varying(512) COLLATE pg_catalog."default" NOT NULL,
    display_name character varying(512) COLLATE pg_catalog."default" NOT NULL,
    birth_date timestamp(6) without time zone,
    gender smallint NOT NULL,
    _version_stamp bigint NOT NULL,
    _timestamp timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    CONSTRAINT personal_info_pkey PRIMARY KEY (id),
    CONSTRAINT personal_info_user_id_key UNIQUE (user_id),
    CONSTRAINT personal_info_user_id_fkey FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE RESTRICT
)
TABLESPACE pg_default;

CREATE INDEX IF NOT EXISTS personal_info_created_at_idx ON public.personal_info (created_at);
CREATE INDEX IF NOT EXISTS personal_info_updated_at_idx ON public.personal_info (updated_at);
CREATE INDEX IF NOT EXISTS personal_info_is_deleted_idx ON public.personal_info (is_deleted);
CREATE INDEX IF NOT EXISTS personal_info_deleted_at_idx ON public.personal_info (deleted_at);
CREATE INDEX IF NOT EXISTS personal_info_first_name_idx ON public.personal_info (first_name);
CREATE INDEX IF NOT EXISTS personal_info_last_name_idx ON public.personal_info (last_name);
CREATE INDEX IF NOT EXISTS personal_info_birth_date_idx ON public.personal_info (birth_date);
CREATE INDEX IF NOT EXISTS personal_info_gender_idx ON public.personal_info (gender);

-- Table: public.user_role_assignments
/*
User role assignment statuses:
    Unspecified = 0
    New         = 1
    Active      = 2
    Inactive    = 3
    Deleting    = 4
    Deleted     = 5
*/
CREATE TABLE IF NOT EXISTS public.user_role_assignments
(
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    role_assignment_id bigint NOT NULL,
    user_id bigint NOT NULL,
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
    CONSTRAINT user_role_assignments_pkey PRIMARY KEY (id),
    CONSTRAINT user_role_assignments_role_assignment_id_key UNIQUE (role_assignment_id),
    CONSTRAINT user_role_assignments_user_id_fkey FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE RESTRICT,
    CONSTRAINT user_role_assignments_status_check CHECK (status >= 1 AND status <= 5)
)
TABLESPACE pg_default;

CREATE UNIQUE INDEX IF NOT EXISTS user_role_assignments_user_id_role_id_idx
    ON public.user_role_assignments (user_id, role_id)
    WHERE status <> 5;

CREATE INDEX IF NOT EXISTS user_role_assignments_user_id_idx ON public.user_role_assignments (user_id);
CREATE INDEX IF NOT EXISTS user_role_assignments_role_id_idx ON public.user_role_assignments (role_id);
CREATE INDEX IF NOT EXISTS user_role_assignments_created_at_idx ON public.user_role_assignments (created_at);
CREATE INDEX IF NOT EXISTS user_role_assignments_updated_at_idx ON public.user_role_assignments (updated_at);
CREATE INDEX IF NOT EXISTS user_role_assignments_status_idx ON public.user_role_assignments (status);
CREATE INDEX IF NOT EXISTS user_role_assignments_status_updated_at_idx ON public.user_role_assignments (status_updated_at);
