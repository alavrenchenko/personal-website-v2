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

-- Database: email_notifier_notifications

CREATE DATABASE email_notifier_notifications
    WITH
    OWNER = postgres
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.UTF-8'
    LC_CTYPE = 'en_US.UTF-8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1
    IS_TEMPLATE = False;


-- Table: public.notifications
/*
Notification statuses:
    Unspecified = 0
	New         = 1
	Sending     = 2
	Sent        = 3
	SendFailed  = 4
	Deleting    = 5
	Deleted     = 6
*/
CREATE TABLE IF NOT EXISTS public.notifications
(
    id uuid NOT NULL,
    group_id bigint NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    created_by bigint NOT NULL,
    updated_at timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    updated_by bigint NOT NULL,
    status smallint NOT NULL,
    status_updated_at timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    status_updated_by bigint NOT NULL,
    status_comment text COLLATE pg_catalog."default",
    recipients text[] COLLATE pg_catalog."default" NOT NULL,
    subject text COLLATE pg_catalog."default" NOT NULL,
    body text COLLATE pg_catalog."default" NOT NULL,
    is_sent boolean NOT NULL DEFAULT FALSE,
    sent_at timestamp(6) without time zone,
    sent_by bigint,
    _version_stamp bigint NOT NULL,
    _timestamp timestamp(6) without time zone NOT NULL DEFAULT (clock_timestamp() AT TIME ZONE 'UTC'::text),
    CONSTRAINT notifications_pkey PRIMARY KEY (id),
    CONSTRAINT notifications_status_check CHECK (status >= 1 AND status <= 6)
)
TABLESPACE pg_default;

CREATE INDEX IF NOT EXISTS notifications_group_id_idx ON public.notifications (group_id);
CREATE INDEX IF NOT EXISTS notifications_created_at_idx ON public.notifications (created_at);
CREATE INDEX IF NOT EXISTS notifications_updated_at_idx ON public.notifications (updated_at);
CREATE INDEX IF NOT EXISTS notifications_status_idx ON public.notifications (status);
CREATE INDEX IF NOT EXISTS notifications_status_updated_at_idx ON public.notifications (status_updated_at);
CREATE INDEX IF NOT EXISTS notifications_is_sent_idx ON public.notifications (is_sent);
CREATE INDEX IF NOT EXISTS notifications_sent_at_idx ON public.notifications (sent_at);
