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

-- ../db/postgres/common/user-agentdb/identity_user_agents.sql

-- Database: identity_mobile_user_agents

CREATE DATABASE identity_mobile_user_agents
    WITH
    OWNER = postgres
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.UTF-8'
    LC_CTYPE = 'en_US.UTF-8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1
    IS_TEMPLATE = False;


-- Table: public.user_agents
/*
User agent types:
    Unspecified = 0
    Web         = 1
    Mobile      = 2

User agent statuses:
    Unspecified          = 0
    New                  = 1
    PendingApproval      = 2
    Active               = 3
    LockedOut            = 4
    TemporarilyLockedOut = 5
    Disabled             = 6
    Deleting             = 7
    Deleted              = 8

id:
increment: 1<<8 = 256
start: (1<<8)+2 = 258 // 00000001 00000010(Mobile), 2(00000010): Mobile User Agent
min_value: (1<<8)+2 = 258
max_value: (1<<63)-1 = 9223372036854775807 // ((256^8)/2)-1
max_count: (1<<55)-1 = 36028797018963967   // ((256^7)/2)-1, 9223372036854775807>>8
exact_max_value: (36028797018963967*256)+2 = 9223372036854775554 // ((9223372036854775807>>8)<<8)+2, 2: Mobile User Agent

id examples:
258       // 00000001 00000010
+256: 514 // 00000010 00000010
+256: 770 // 00000011 00000010
and so on
*/
CREATE TABLE IF NOT EXISTS public.user_agents
(
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 256 START 258 MINVALUE 258 MAXVALUE 9223372036854775807 CACHE 1 ),
    user_id bigint NOT NULL,
    client_id bigint NOT NULL,
    type smallint NOT NULL GENERATED ALWAYS AS (2) STORED,
    created_at timestamp(6) without time zone NOT NULL,
    created_by bigint NOT NULL,
    updated_at timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    updated_by bigint NOT NULL,
    status smallint NOT NULL,
    status_updated_at timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    status_updated_by bigint NOT NULL,
    status_comment text COLLATE pg_catalog."default",
    app_id bigint,
    first_user_agent text COLLATE pg_catalog."default",
    last_user_agent text COLLATE pg_catalog."default",
    first_sign_in_time timestamp(6) without time zone,
    first_sign_in_ip character varying(64) COLLATE pg_catalog."default",
    last_sign_in_time timestamp(6) without time zone,
    last_sign_in_ip character varying(64) COLLATE pg_catalog."default",
    last_sign_out_time timestamp(6) without time zone,
    last_activity_time timestamp(6) without time zone,
    last_activity_ip character varying(64) COLLATE pg_catalog."default",
    _version_stamp bigint NOT NULL,
    _timestamp timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    CONSTRAINT user_agents_pkey PRIMARY KEY (id),
    CONSTRAINT user_agents_status_check CHECK (status >= 1 AND status <= 8)
)
TABLESPACE pg_default;

CREATE INDEX IF NOT EXISTS user_agents_user_id_idx ON public.user_agents (user_id);
CREATE INDEX IF NOT EXISTS user_agents_client_id_idx ON public.user_agents (client_id);
CREATE INDEX IF NOT EXISTS user_agents_created_at_idx ON public.user_agents (created_at);
CREATE INDEX IF NOT EXISTS user_agents_updated_at_idx ON public.user_agents (updated_at);
CREATE INDEX IF NOT EXISTS user_agents_status_idx ON public.user_agents (status);
CREATE INDEX IF NOT EXISTS user_agents_status_updated_at_idx ON public.user_agents (status_updated_at);
CREATE INDEX IF NOT EXISTS user_agents_app_id_idx ON public.user_agents (app_id);
CREATE INDEX IF NOT EXISTS user_agents_first_sign_in_time_idx ON public.user_agents (first_sign_in_time);
CREATE INDEX IF NOT EXISTS user_agents_first_sign_in_ip_idx ON public.user_agents (first_sign_in_ip);
CREATE INDEX IF NOT EXISTS user_agents_last_sign_in_time_idx ON public.user_agents (last_sign_in_time);
CREATE INDEX IF NOT EXISTS user_agents_last_sign_in_ip_idx ON public.user_agents (last_sign_in_ip);
CREATE INDEX IF NOT EXISTS user_agents_last_activity_time_idx ON public.user_agents (last_activity_time);
CREATE INDEX IF NOT EXISTS user_agents_last_activity_ip_idx ON public.user_agents (last_activity_ip);

-- Table: public.user_agent_sessions
/*
User agent session types:
    Unspecified = 0
    Web         = 1
    Mobile      = 2

User agent session statuses:
    Unspecified          = 0
    New                  = 1
    Active               = 2
    SignedOut            = 3
    Ended                = 4
    LockedOut            = 5
    TemporarilyLockedOut = 6
    Disabled             = 7
    Deleting             = 8
    Deleted              = 9

id:
increment: 1<<8 = 256
start: (1<<8)+2 = 258 // 00000001 00000010(Mobile), 2(00000010): Mobile User Agent Session
min_value: (1<<8)+2 = 258
max_value: (1<<63)-1 = 9223372036854775807 // ((256^8)/2)-1
max_count: (1<<55)-1 = 36028797018963967   // ((256^7)/2)-1, 9223372036854775807>>8
exact_max_value: (36028797018963967*256)+2 = 9223372036854775554 // ((9223372036854775807>>8)<<8)+2, 2: Mobile User Agent Session

id examples:
258       // 00000001 00000010
+256: 514 // 00000010 00000010
+256: 770 // 00000011 00000010
and so on
*/
CREATE TABLE IF NOT EXISTS public.user_agent_sessions
(
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 256 START 258 MINVALUE 258 MAXVALUE 9223372036854775807 CACHE 1 ),
    user_id bigint NOT NULL,
    client_id bigint NOT NULL,
    user_agent_id bigint NOT NULL,
    type smallint NOT NULL GENERATED ALWAYS AS (2) STORED,
    user_session_id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    created_by bigint NOT NULL,
    updated_at timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    updated_by bigint NOT NULL,
    status smallint NOT NULL,
    status_updated_at timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    status_updated_by bigint NOT NULL,
    status_comment text COLLATE pg_catalog."default",
    first_sign_in_time timestamp(6) without time zone,
    first_sign_in_ip character varying(64) COLLATE pg_catalog."default",
    last_sign_in_time timestamp(6) without time zone,
    last_sign_in_ip character varying(64) COLLATE pg_catalog."default",
    last_sign_out_time timestamp(6) without time zone,
    last_activity_time timestamp(6) without time zone,
    last_activity_ip character varying(64) COLLATE pg_catalog."default",
    _version_stamp bigint NOT NULL,
    _timestamp timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    CONSTRAINT user_agent_sessions_pkey PRIMARY KEY (id),
    CONSTRAINT user_agent_sessions_user_agent_id_fkey FOREIGN KEY (user_agent_id)
        REFERENCES public.user_agents (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE RESTRICT,
    CONSTRAINT user_agent_sessions_status_check CHECK (status >= 1 AND status <= 9)
)
TABLESPACE pg_default;

CREATE UNIQUE INDEX IF NOT EXISTS user_agent_sessions_user_id_client_id_idx
    ON public.user_agent_sessions (user_id, client_id)
    WHERE status <> 9;

CREATE UNIQUE INDEX IF NOT EXISTS user_agent_sessions_user_agent_id_uidx
    ON public.user_agent_sessions (user_agent_id)
    WHERE status <> 9;

CREATE INDEX IF NOT EXISTS user_agent_sessions_user_id_idx ON public.user_agent_sessions (user_id);
CREATE INDEX IF NOT EXISTS user_agent_sessions_client_id_idx ON public.user_agent_sessions (client_id);
CREATE INDEX IF NOT EXISTS user_agent_sessions_user_agent_id_idx ON public.user_agent_sessions (user_agent_id);
CREATE INDEX IF NOT EXISTS user_agent_sessions_created_at_idx ON public.user_agent_sessions (created_at);
CREATE INDEX IF NOT EXISTS user_agent_sessions_updated_at_idx ON public.user_agent_sessions (updated_at);
CREATE INDEX IF NOT EXISTS user_agent_sessions_status_idx ON public.user_agent_sessions (status);
CREATE INDEX IF NOT EXISTS user_agent_sessions_status_updated_at_idx ON public.user_agent_sessions (status_updated_at);
CREATE INDEX IF NOT EXISTS user_agent_sessions_first_sign_in_time_idx ON public.user_agent_sessions (first_sign_in_time);
CREATE INDEX IF NOT EXISTS user_agent_sessions_first_sign_in_ip_idx ON public.user_agent_sessions (first_sign_in_ip);
CREATE INDEX IF NOT EXISTS user_agent_sessions_last_sign_in_time_idx ON public.user_agent_sessions (last_sign_in_time);
CREATE INDEX IF NOT EXISTS user_agent_sessions_last_sign_in_ip_idx ON public.user_agent_sessions (last_sign_in_ip);
CREATE INDEX IF NOT EXISTS user_agent_sessions_last_activity_time_idx ON public.user_agent_sessions (last_activity_time);
CREATE INDEX IF NOT EXISTS user_agent_sessions_last_activity_ip_idx ON public.user_agent_sessions (last_activity_ip);
