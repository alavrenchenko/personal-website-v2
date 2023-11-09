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

-- ../db/postgres/common/user-sessiondb/identity_user_sessions.sql

-- Database: identity_user_web_sessions

CREATE DATABASE identity_user_web_sessions
    WITH
    OWNER = postgres
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.UTF-8'
    LC_CTYPE = 'en_US.UTF-8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1
    IS_TEMPLATE = False;


-- Table: public.user_sessions
/*
User session types:
    Unspecified = 0
    Web         = 1
    Mobile      = 2

User session statuses:
    Unspecified = 0
    New         = 1
    Active      = 2
    Ended       = 3

id:
increment: 1<<8 = 256
start: (1<<8)+1 = 257 // 00000001 00000001(Web), 1(00000001): User's Web Session
min_value: (1<<8)+1 = 257
max_value: (1<<63)-1 = 9223372036854775807 // ((256^8)/2)-1
max_count: (1<<55)-1 = 36028797018963967   // ((256^7)/2)-1, 9223372036854775807>>8
exact_max_value: (36028797018963967*256)+1 = 9223372036854775553 // ((9223372036854775807>>8)<<8)+1, 1: User's Web Session

id examples:
257       // 00000001 00000001
+256: 513 // 00000010 00000001
+256: 769 // 00000011 00000001
and so on
*/
CREATE TABLE IF NOT EXISTS public.user_sessions
(
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 256 START 257 MINVALUE 257 MAXVALUE 9223372036854775807 CACHE 1 ),
    user_id bigint NOT NULL,
    client_id bigint NOT NULL,
    user_agent_id bigint NOT NULL,
    type smallint NOT NULL GENERATED ALWAYS AS (1) STORED,
    created_at timestamp(6) without time zone NOT NULL,
    created_by bigint NOT NULL,
    updated_at timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    updated_by bigint NOT NULL,
    status smallint NOT NULL,
    status_updated_at timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    status_updated_by bigint NOT NULL,
    status_comment text COLLATE pg_catalog."default",
    app_id bigint,
    start_time timestamp(6) without time zone,
    end_time timestamp(6) without time zone,
    first_ip character varying(64) COLLATE pg_catalog."default" NOT NULL,
    last_activity_time timestamp(6) without time zone,
    last_activity_ip character varying(64) COLLATE pg_catalog."default",
    _version_stamp bigint NOT NULL,
    _timestamp timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    CONSTRAINT user_sessions_pkey PRIMARY KEY (id),
    CONSTRAINT user_sessions_status_check CHECK (status >= 1 AND status <= 3)
)
TABLESPACE pg_default;

CREATE UNIQUE INDEX IF NOT EXISTS user_sessions_user_id_client_id_idx
    ON public.user_sessions (user_id, client_id)
    WHERE status <> 3;

CREATE UNIQUE INDEX IF NOT EXISTS user_sessions_user_agent_id_uidx
    ON public.user_sessions (user_agent_id)
    WHERE status <> 3;

CREATE INDEX IF NOT EXISTS user_sessions_user_id_idx ON public.user_sessions (user_id);
CREATE INDEX IF NOT EXISTS user_sessions_client_id_idx ON public.user_sessions (client_id);
CREATE INDEX IF NOT EXISTS user_sessions_user_agent_id_idx ON public.user_sessions (user_agent_id);
CREATE INDEX IF NOT EXISTS user_sessions_created_at_idx ON public.user_sessions (created_at);
CREATE INDEX IF NOT EXISTS user_sessions_updated_at_idx ON public.user_sessions (updated_at);
CREATE INDEX IF NOT EXISTS user_sessions_status_idx ON public.user_sessions (status);
CREATE INDEX IF NOT EXISTS user_sessions_status_updated_at_idx ON public.user_sessions (status_updated_at);
CREATE INDEX IF NOT EXISTS user_sessions_app_id_idx ON public.user_sessions (app_id);
CREATE INDEX IF NOT EXISTS user_sessions_first_ip_idx ON public.user_sessions (first_ip);
CREATE INDEX IF NOT EXISTS user_sessions_last_activity_time_idx ON public.user_sessions (last_activity_time);
CREATE INDEX IF NOT EXISTS user_sessions_last_activity_ip_idx ON public.user_sessions (last_activity_ip);
