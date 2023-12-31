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

-- Database: app_manager

CREATE DATABASE app_manager
    WITH
    OWNER = postgres
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.UTF-8'
    LC_CTYPE = 'en_US.UTF-8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1
    IS_TEMPLATE = False;


-- Table: public.app_groups
/*
App group types:
    Unspecified = 0
    Service     = 1

App group statuses:
    Unspecified = 0
    New         = 1
    Active      = 2
    Inactive    = 3
    Deleting    = 4
    Deleted     = 5
*/
CREATE TABLE IF NOT EXISTS public.app_groups
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
    version character varying(64) COLLATE pg_catalog."default" NOT NULL,
    description text COLLATE pg_catalog."default" NOT NULL,
    _version_stamp bigint NOT NULL,
    _timestamp timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    CONSTRAINT app_groups_pkey PRIMARY KEY (id),
    CONSTRAINT app_groups_type_check CHECK (type = 1),
    CONSTRAINT app_groups_status_check CHECK (status >= 1 AND status <= 5)
)
TABLESPACE pg_default;

CREATE UNIQUE INDEX IF NOT EXISTS app_groups_name_lc_idx
    ON public.app_groups (lower(name))
    WHERE status <> 5;

CREATE INDEX IF NOT EXISTS app_groups_type_idx ON public.app_groups (type);
CREATE INDEX IF NOT EXISTS app_groups_created_at_idx ON public.app_groups (created_at);
CREATE INDEX IF NOT EXISTS app_groups_updated_at_idx ON public.app_groups (updated_at);
CREATE INDEX IF NOT EXISTS app_groups_status_idx ON public.app_groups (status);
CREATE INDEX IF NOT EXISTS app_groups_status_updated_at_idx ON public.app_groups (status_updated_at);
CREATE INDEX IF NOT EXISTS app_groups_version_idx ON public.app_groups (version);

-- Table: public.apps
/*
App types:
    Unspecified = 0
    Service     = 1

App categories:
    Unspecified = 0
    Service     = 1

App statuses:
    Unspecified = 0
    New         = 1
    Active      = 2
    Inactive    = 3
    Deleting    = 4
    Deleted     = 5
*/
CREATE TABLE IF NOT EXISTS public.apps
(
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    name character varying(256) COLLATE pg_catalog."default" NOT NULL,
    group_id bigint NOT NULL,
    type smallint NOT NULL,
    title character varying(256) COLLATE pg_catalog."default" NOT NULL,
    category smallint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    created_by bigint NOT NULL,
    updated_at timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    updated_by bigint NOT NULL,
    status smallint NOT NULL,
    status_updated_at timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    status_updated_by bigint NOT NULL,
    status_comment text COLLATE pg_catalog."default",
    version character varying(64) COLLATE pg_catalog."default" NOT NULL,
    description text COLLATE pg_catalog."default" NOT NULL,
    _version_stamp bigint NOT NULL,
    _timestamp timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    CONSTRAINT apps_pkey PRIMARY KEY (id),
    CONSTRAINT apps_group_id_fkey FOREIGN KEY (group_id)
        REFERENCES public.app_groups (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE RESTRICT,
    CONSTRAINT apps_type_check CHECK (type = 1),
    CONSTRAINT apps_category_check CHECK (category = 1),
    CONSTRAINT apps_status_check CHECK (status >= 1 AND status <= 5)
)
TABLESPACE pg_default;

CREATE UNIQUE INDEX IF NOT EXISTS apps_name_lc_idx
    ON public.apps (lower(name))
    WHERE status <> 5;

CREATE INDEX IF NOT EXISTS apps_group_id_idx ON public.apps (group_id);
CREATE INDEX IF NOT EXISTS apps_type_idx ON public.apps (type);
CREATE INDEX IF NOT EXISTS apps_category_idx ON public.apps (category);
CREATE INDEX IF NOT EXISTS apps_created_at_idx ON public.apps (created_at);
CREATE INDEX IF NOT EXISTS apps_updated_at_idx ON public.apps (updated_at);
CREATE INDEX IF NOT EXISTS apps_status_idx ON public.apps (status);
CREATE INDEX IF NOT EXISTS apps_status_updated_at_idx ON public.apps (status_updated_at);
CREATE INDEX IF NOT EXISTS apps_version_idx ON public.apps (version);

-- Table: public.app_sessions
/*
App session statuses:
    Unspecified = 0
    New         = 1
    Active      = 2
    Ended       = 3
    Deleting    = 4
    Deleted     = 5
*/
CREATE TABLE IF NOT EXISTS public.app_sessions
(
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    app_id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    created_by bigint NOT NULL,
    updated_at timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    updated_by bigint NOT NULL,
    status smallint NOT NULL,
    status_updated_at timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    status_updated_by bigint NOT NULL,
    status_comment text COLLATE pg_catalog."default",
    start_time timestamp(6) without time zone,
    end_time timestamp(6) without time zone,
    _version_stamp bigint NOT NULL,
    _timestamp timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    CONSTRAINT app_sessions_pkey PRIMARY KEY (id),
    CONSTRAINT app_sessions_app_id_fkey FOREIGN KEY (app_id)
        REFERENCES public.apps (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE RESTRICT,
    CONSTRAINT app_sessions_status_check CHECK (status >= 1 AND status <= 5)
)
TABLESPACE pg_default;

CREATE INDEX IF NOT EXISTS app_sessions_app_id_idx ON public.app_sessions (app_id);
CREATE INDEX IF NOT EXISTS app_sessions_created_at_idx ON public.app_sessions (created_at);
CREATE INDEX IF NOT EXISTS app_sessions_updated_at_idx ON public.app_sessions (updated_at);
CREATE INDEX IF NOT EXISTS app_sessions_status_idx ON public.app_sessions (status);
CREATE INDEX IF NOT EXISTS app_sessions_status_updated_at_idx ON public.app_sessions (status_updated_at);
