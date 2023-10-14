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

-- Database: identity_{type}_user_agents

CREATE DATABASE identity_{type}_user_agents
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
    Deleted              = 7
*/
CREATE TABLE IF NOT EXISTS public.user_agents
(
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT {increment} START {start} MINVALUE {min_value} MAXVALUE 9223372036854775807 CACHE 1 ),
    user_id bigint NOT NULL,
    client_id bigint NOT NULL,
    type smallint NOT NULL GENERATED ALWAYS AS ({type}) STORED,
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
    CONSTRAINT user_agents_pkey PRIMARY KEY (id)
)
TABLESPACE pg_default;

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
    LockedOut            = 4
    TemporarilyLockedOut = 5
    Disabled             = 6
    Ended                = 7
*/
CREATE TABLE IF NOT EXISTS public.user_agent_sessions
(
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT {increment} START {start} MINVALUE {min_value} MAXVALUE 9223372036854775807 CACHE 1 ),
    user_id bigint NOT NULL,
    client_id bigint NOT NULL,
    user_agent_id bigint NOT NULL,
    type smallint NOT NULL GENERATED ALWAYS AS ({type}) STORED,
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
    first_sign_in_time timestamp(6) without time zone NOT NULL,
    first_sign_in_ip character varying(64) COLLATE pg_catalog."default" NOT NULL,
    last_sign_in_time timestamp(6) without time zone NOT NULL,
    last_sign_in_ip character varying(64) COLLATE pg_catalog."default" NOT NULL,
    last_sign_out_time timestamp(6) without time zone NOT NULL,
    last_activity_time timestamp(6) without time zone NOT NULL,
    last_activity_ip character varying(64) COLLATE pg_catalog."default" NOT NULL,
    _version_stamp bigint NOT NULL,
    _timestamp timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    CONSTRAINT user_agent_sessions_pkey PRIMARY KEY (id),
    CONSTRAINT user_agent_sessions_user_agent_id_fkey FOREIGN KEY (user_agent_id)
        REFERENCES public.user_agents (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE RESTRICT
)
TABLESPACE pg_default;