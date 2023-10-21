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

-- Database: identity_role_assignments

CREATE DATABASE identity_role_assignments
    WITH
    OWNER = postgres
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.UTF-8'
    LC_CTYPE = 'en_US.UTF-8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1
    IS_TEMPLATE = False;


-- Table: public.role_assignments
/*
Assignee types:
    Unspecified = 0
    User        = 1
    Group       = 2

Role assignment statuses:
    Unspecified = 0
    New         = 1
    Active      = 2
    Inactive    = 3
    Deleting    = 4
    Deleted     = 5
*/
CREATE TABLE IF NOT EXISTS public.role_assignments
(
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    role_id bigint NOT NULL,
    assigned_to bigint NOT NULL,
    assignee_type smallint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    created_by bigint NOT NULL,
    updated_at timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    updated_by bigint NOT NULL,
    status smallint NOT NULL,
    status_updated_at timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    status_updated_by bigint NOT NULL,
    status_comment text COLLATE pg_catalog."default",
    description text COLLATE pg_catalog."default",
    _version_stamp bigint NOT NULL,
    _timestamp timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    CONSTRAINT role_assignments_pkey PRIMARY KEY (id),
    CONSTRAINT role_assignments_assignee_type_check CHECK (assignee_type = 1 OR assignee_type = 2),
    CONSTRAINT role_assignments_status_check CHECK (status >= 1 AND status <= 5)
)
TABLESPACE pg_default;

CREATE UNIQUE INDEX IF NOT EXISTS role_assignments_role_id_assigned_to_assignee_type_idx
    ON public.role_assignments (role_id, assigned_to, assignee_type)
    WHERE status <> 5;

CREATE INDEX IF NOT EXISTS role_assignments_role_id_idx ON public.role_assignments (role_id);
CREATE INDEX IF NOT EXISTS role_assignments_assigned_to_idx ON public.role_assignments (assigned_to);
CREATE INDEX IF NOT EXISTS role_assignments_assignee_type_idx ON public.role_assignments (assignee_type);
CREATE INDEX IF NOT EXISTS role_assignments_created_at_idx ON public.role_assignments (created_at);
CREATE INDEX IF NOT EXISTS role_assignments_updated_at_idx ON public.role_assignments (updated_at);
CREATE INDEX IF NOT EXISTS role_assignments_status_idx ON public.role_assignments (status);
CREATE INDEX IF NOT EXISTS role_assignments_status_updated_at_idx ON public.role_assignments (status_updated_at);

-- Table: public.deleted_role_assignments
/*
Assignee types:
    Unspecified = 0
    User        = 1
    Group       = 2

Role assignment statuses:
    Unspecified = 0
    New         = 1
    Active      = 2
    Inactive    = 3
    Deleting    = 4
    Deleted     = 5
*/
CREATE TABLE IF NOT EXISTS public.deleted_role_assignments
(
    id bigint NOT NULL,
    role_id bigint NOT NULL,
    assigned_to bigint NOT NULL,
    assignee_type smallint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    created_by bigint NOT NULL,
    updated_at timestamp(6) without time zone NOT NULL,
    updated_by bigint NOT NULL,
    status smallint NOT NULL,
    status_updated_at timestamp(6) without time zone NOT NULL,
    status_updated_by bigint NOT NULL,
    status_comment text COLLATE pg_catalog."default",
    description text COLLATE pg_catalog."default",
    _version_stamp bigint NOT NULL,
    _timestamp timestamp(6) without time zone NOT NULL,
    CONSTRAINT deleted_role_assignments_pkey PRIMARY KEY (id),
    CONSTRAINT deleted_role_assignments_assignee_type_check CHECK (assignee_type = 1 OR assignee_type = 2),
    CONSTRAINT deleted_role_assignments_status_check CHECK (status = 5)
)
TABLESPACE pg_default;

CREATE INDEX IF NOT EXISTS deleted_role_assignments_role_id_idx ON public.deleted_role_assignments (role_id);
CREATE INDEX IF NOT EXISTS deleted_role_assignments_assigned_to_idx ON public.deleted_role_assignments (assigned_to);
CREATE INDEX IF NOT EXISTS deleted_role_assignments_assignee_type_idx ON public.deleted_role_assignments (assignee_type);
CREATE INDEX IF NOT EXISTS deleted_role_assignments_created_at_idx ON public.deleted_role_assignments (created_at);
CREATE INDEX IF NOT EXISTS deleted_role_assignments_updated_at_idx ON public.deleted_role_assignments (updated_at);
CREATE INDEX IF NOT EXISTS deleted_role_assignments_status_idx ON public.deleted_role_assignments (status);
CREATE INDEX IF NOT EXISTS deleted_role_assignments_status_updated_at_idx ON public.deleted_role_assignments (status_updated_at);

CREATE INDEX IF NOT EXISTS deleted_role_assignments_role_id_assigned_to_assignee_type_idx
    ON public.deleted_role_assignments (role_id, assigned_to, assignee_type);
