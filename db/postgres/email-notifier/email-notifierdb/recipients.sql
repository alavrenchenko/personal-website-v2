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

-- FUNCTION: public.recipient_exists(bigint, character varying)

CREATE OR REPLACE FUNCTION public.recipient_exists(
    _notif_group_id public.recipients.notif_group_id%TYPE,
    _email public.recipients.email%TYPE
) RETURNS boolean AS $$
BEGIN
    RETURN EXISTS (SELECT 1 FROM public.recipients WHERE notif_group_id = _notif_group_id AND lower(email) = lower(_email) AND is_deleted IS FALSE LIMIT 1);
END;
$$ LANGUAGE plpgsql;

-- PROCEDURE: public.create_recipient(bigint, smallint, bigint, character varying, character varying, character varying)
/*
Notification group statuses:
    Active = 2

Error codes:
    NoError                   = 0
    InvalidOperation          = 3
    NotificationGroupNotFound = 11200
    RecipientAlreadyExists    = 11401
*/
-- Minimum transaction isolation level: Read committed.
CREATE OR REPLACE PROCEDURE public.create_recipient(
    IN _notif_group_id public.recipients.notif_group_id%TYPE,
    IN _type public.recipients.type%TYPE,
    IN _created_by public.recipients.created_by%TYPE,
    IN _name public.recipients.name%TYPE,
    IN _email public.recipients.email%TYPE,
    IN _addr public.recipients.addr%TYPE,
    OUT _id public.recipients.id%TYPE,
    OUT err_code bigint,
    OUT err_msg text) AS $$
DECLARE
    _time timestamp(6) without time zone;
    _status public.notification_groups.status%TYPE;
BEGIN
    _id := 0;
    err_code := 0; -- NoError
    err_msg := '';

    SELECT status INTO _status FROM public.notification_groups WHERE id = _notif_group_id LIMIT 1 FOR UPDATE;
    IF NOT FOUND THEN
        err_code := 11200; -- NotificationGroupNotFound
        err_msg := 'notification group not found';
        RETURN;
    END IF;

    -- notification group status: Active(2)
    IF _status <> 2 THEN
        err_code := 3; -- InvalidOperation
        err_msg := format('invalid notification group status (%s)', _status);
        RETURN;
    END IF;

    IF public.recipient_exists(_notif_group_id, _email) THEN
        err_code := 11401; -- RecipientAlreadyExists
        err_msg := 'notification recipient with the same params already exists';
        RETURN;
    END IF;

    _time := (clock_timestamp() AT TIME ZONE 'UTC');
    INSERT INTO public.recipients(notif_group_id, type, created_at, created_by, name, email, addr, _version_stamp, _timestamp)
        VALUES (_notif_group_id, _type, _time, _created_by, _name, _email, _addr, 1, _time)
        RETURNING id INTO _id;

    EXCEPTION
        WHEN unique_violation THEN
            IF public.recipient_exists(_notif_group_id, _email) THEN
                err_code := 11401; -- RecipientAlreadyExists
                err_msg := 'notification recipient with the same params already exists';
                RETURN;
            END IF;
            RAISE;
END;
$$ LANGUAGE plpgsql;
