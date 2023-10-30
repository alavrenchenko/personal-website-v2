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

-- Database: identity_user_groups

CREATE DATABASE identity_user_groups
    WITH
    OWNER = postgres
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.UTF-8'
    LC_CTYPE = 'en_US.UTF-8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1
    IS_TEMPLATE = False;


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
    "group" bigint NOT NULL,
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
