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
User groups:
    Unspecified    = 0
    AnonymousUsers = 1
    SystemUsers    = 2
    Administrators = 3
    StandardUsers  = 4

User statuses:
    Unspecified          = 0
    New                  = 1
    PendingApproval      = 2
    Active               = 3
    LockedOut            = 4
    TemporarilyLockedOut = 5
    Disabled             = 6
    Deleted              = 7
*/
CREATE TABLE IF NOT EXISTS public.users
(
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    name character varying(256) COLLATE pg_catalog."default" NOT NULL,
    "group" smallint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    created_by bigint NOT NULL,
    updated_at timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    updated_by bigint NOT NULL,
    status smallint NOT NULL,
    status_updated_at timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    status_updated_by bigint NOT NULL,
    status_comment text COLLATE pg_catalog."default",
    email character varying(512) COLLATE pg_catalog."default",
    last_sign_in_time timestamp(6) without time zone NOT NULL,
    last_sign_in_ip character varying(64) COLLATE pg_catalog."default" NOT NULL,
    last_sign_out_time timestamp(6) without time zone NOT NULL,
    last_activity_time timestamp(6) without time zone NOT NULL,
    last_activity_ip character varying(64) COLLATE pg_catalog."default" NOT NULL,
    _version_stamp bigint NOT NULL,
    _timestamp timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    CONSTRAINT users_pkey PRIMARY KEY (id),
    CONSTRAINT users_name_key UNIQUE (name)
)
TABLESPACE pg_default;

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
    first_name character varying(512) COLLATE pg_catalog."default" NOT NULL,
    last_name character varying(512) COLLATE pg_catalog."default" NOT NULL,
    display_name character varying(512) COLLATE pg_catalog."default" NOT NULL,
    birth_date timestamp(6) without time zone,
    gender smallint NOT NULL,
    is_deleted boolean NOT NULL DEFAULT false,
    deleted_at timestamp(6) without time zone,
    deleted_by bigint,
    _version_stamp bigint NOT NULL,
    _timestamp timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    CONSTRAINT personal_info_pkey PRIMARY KEY (id),
    CONSTRAINT personal_info_user_id_fkey FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE RESTRICT
)
TABLESPACE pg_default;