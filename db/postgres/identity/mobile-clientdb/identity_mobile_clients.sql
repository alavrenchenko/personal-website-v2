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

-- Database: identity_mobile_clients

CREATE DATABASE identity_mobile_clients
    WITH
    OWNER = postgres
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.UTF-8'
    LC_CTYPE = 'en_US.UTF-8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1
    IS_TEMPLATE = False;


-- Table: public.clients
/*
Client types:
    Unspecified = 0
    Web         = 1
    Mobile      = 2

Client statuses:
    Unspecified          = 0
    New                  = 1
    PendingApproval      = 2
    Active               = 3
    LockedOut            = 4
    TemporarilyLockedOut = 5
    Disabled             = 6
    Deleted              = 7

id:
increment: 1<<8 = 256
start: (1<<8)+2 = 258 // 00000001 00000010(Mobile), 2(00000010): Mobile Client
min_value: (1<<8)+2 = 258
max_value: (1<<63)-1 = 9223372036854775807 // ((256^8)/2)−1
max_count: (1<<55)-1 = 36028797018963967   // ((256^7)/2)−1, 9223372036854775807>>8

id examples:
258       // 00000001 00000010
+256: 514 // 00000010 00000010
+256: 770 // 00000011 00000010
and so on
*/
CREATE TABLE IF NOT EXISTS public.clients
(
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 256 START 258 MINVALUE 258 MAXVALUE 9223372036854775807 CACHE 1 ),
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
    last_activity_time timestamp(6) without time zone NOT NULL,
    last_activity_ip character varying(64) COLLATE pg_catalog."default" NOT NULL,
    _version_stamp bigint NOT NULL,
    _timestamp timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    CONSTRAINT clients_pkey PRIMARY KEY (id)
)
TABLESPACE pg_default;
