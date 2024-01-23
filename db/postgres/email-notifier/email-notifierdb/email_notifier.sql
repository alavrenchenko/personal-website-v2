-- Copyright 2024 Alexey Lavrenchenko. All rights reserved.
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

-- Database: email_notifier

CREATE DATABASE email_notifier
    WITH
    OWNER = postgres
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.UTF-8'
    LC_CTYPE = 'en_US.UTF-8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1
    IS_TEMPLATE = False;


-- Table: public.notification_groups
/*
Notification group statuses:
    Unspecified = 0
    New         = 1
    Active      = 2
    Inactive    = 3
    Deleting    = 4
    Deleted     = 5
*/
CREATE TABLE IF NOT EXISTS public.notification_groups
(
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    name character varying(256) COLLATE pg_catalog."default" NOT NULL,
    title character varying(256) COLLATE pg_catalog."default" NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    created_by bigint NOT NULL,
    updated_at timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    updated_by bigint NOT NULL,
    status smallint NOT NULL,
    status_updated_at timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    status_updated_by bigint NOT NULL,
    status_comment text COLLATE pg_catalog."default",
    description text COLLATE pg_catalog."default" NOT NULL,
    sender_name character varying(128) COLLATE pg_catalog."default",
    sender_email character varying(512) COLLATE pg_catalog."default" NOT NULL,
    sender_addr character varying(640) COLLATE pg_catalog."default" NOT NULL,
    _version_stamp bigint NOT NULL,
    _timestamp timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    CONSTRAINT notification_groups_pkey PRIMARY KEY (id),
    CONSTRAINT notification_groups_status_check CHECK (status >= 1 AND status <= 5)
)
TABLESPACE pg_default;

CREATE UNIQUE INDEX IF NOT EXISTS notification_groups_name_lc_idx
    ON public.notification_groups (lower(name))
    WHERE status <> 5;

CREATE INDEX IF NOT EXISTS notification_groups_created_at_idx ON public.notification_groups (created_at);
CREATE INDEX IF NOT EXISTS notification_groups_updated_at_idx ON public.notification_groups (updated_at);
CREATE INDEX IF NOT EXISTS notification_groups_status_idx ON public.notification_groups (status);
CREATE INDEX IF NOT EXISTS notification_groups_status_updated_at_idx ON public.notification_groups (status_updated_at);
CREATE INDEX IF NOT EXISTS notification_groups_sender_email_lc_idx ON public.notification_groups (lower(sender_email));
CREATE INDEX IF NOT EXISTS notification_groups_sender_addr_lc_idx ON public.notification_groups (lower(sender_addr));

-- Table: public.recipients
/*
Recipient types:
    Unspecified = 0
    CC          = 2
    BCC         = 3
*/
CREATE TABLE IF NOT EXISTS public.recipients
(
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    notif_group_id bigint NOT NULL,
    type smallint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    created_by bigint NOT NULL,
    is_deleted boolean NOT NULL DEFAULT FALSE,
    deleted_at timestamp(6) without time zone,
    deleted_by bigint,
    name character varying(128) COLLATE pg_catalog."default",
    email character varying(512) COLLATE pg_catalog."default" NOT NULL,
    addr character varying(640) COLLATE pg_catalog."default" NOT NULL,
    _version_stamp bigint NOT NULL,
    _timestamp timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    CONSTRAINT recipients_pkey PRIMARY KEY (id),
    CONSTRAINT recipients_notif_group_id_fkey FOREIGN KEY (notif_group_id)
        REFERENCES public.notification_groups (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE RESTRICT,
    CONSTRAINT recipients_type_check CHECK (type = 2 OR type = 3)
)
TABLESPACE pg_default;

CREATE UNIQUE INDEX IF NOT EXISTS recipients_notif_group_id_email_lc_idx
    ON public.recipients (notif_group_id, lower(email))
    WHERE is_deleted IS FALSE;

CREATE INDEX IF NOT EXISTS recipients_notif_group_id_idx ON public.recipients (notif_group_id);
CREATE INDEX IF NOT EXISTS recipients_type_idx ON public.recipients (type);
CREATE INDEX IF NOT EXISTS recipients_created_at_idx ON public.recipients (created_at);
CREATE INDEX IF NOT EXISTS recipients_is_deleted_idx ON public.recipients (is_deleted);
CREATE INDEX IF NOT EXISTS recipients_deleted_at_idx ON public.recipients (deleted_at);
CREATE INDEX IF NOT EXISTS recipients_email_lc_idx ON public.recipients (lower(email));
CREATE INDEX IF NOT EXISTS recipients_addr_lc_idx ON public.recipients (lower(addr));
