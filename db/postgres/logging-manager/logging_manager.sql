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

-- Database: logging_manager

CREATE DATABASE logging_manager
    WITH
    OWNER = postgres
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.UTF-8'
    LC_CTYPE = 'en_US.UTF-8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1
    IS_TEMPLATE = False;


-- Table: public.logging_sessions
/*
Logging session statuses:
    Unspecified = 0
    New         = 1
    Started     = 2
    Deleting    = 4
    Deleted     = 5
*/
CREATE TABLE IF NOT EXISTS public.logging_sessions
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
    _version_stamp bigint NOT NULL,
    _timestamp timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    CONSTRAINT logging_sessions_pkey PRIMARY KEY (id),
    CONSTRAINT logging_sessions_status_check CHECK (status IN (1, 2, 4, 5))
)
TABLESPACE pg_default;

CREATE INDEX IF NOT EXISTS logging_sessions_app_id_idx ON public.logging_sessions (app_id);
CREATE INDEX IF NOT EXISTS logging_sessions_created_at_idx ON public.logging_sessions (created_at);
CREATE INDEX IF NOT EXISTS logging_sessions_updated_at_idx ON public.logging_sessions (updated_at);
CREATE INDEX IF NOT EXISTS logging_sessions_status_idx ON public.logging_sessions (status);
CREATE INDEX IF NOT EXISTS logging_sessions_status_updated_at_idx ON public.logging_sessions (status_updated_at);
