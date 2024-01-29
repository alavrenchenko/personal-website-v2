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

-- FUNCTION: public.notification_exists(uuid)

CREATE OR REPLACE FUNCTION public.notification_exists(
    _id public.notifications.id%TYPE
) RETURNS boolean AS $$
BEGIN
    RETURN EXISTS (SELECT 1 FROM public.notifications WHERE id = _id LIMIT 1);
END;
$$ LANGUAGE plpgsql;

-- PROCEDURE: public.add_notification(uuid, bigint, timestamp without time zone, bigint, smallint, text, text[], text, text, timestamp without time zone)
/*
Notification statuses:
	New        = 1
	Sending    = 2
	Sent       = 3
	SendFailed = 4

Error codes:
    NoError                   = 0
    InvalidOperation          = 3
    NotificationAlreadyExists = 11001
*/
-- Minimum transaction isolation level: Read committed.
CREATE OR REPLACE PROCEDURE public.add_notification(
    IN _id public.notifications.id%TYPE,
    IN _group_id public.notifications.group_id%TYPE,
    IN _created_at public.notifications.created_at%TYPE,
    IN _created_by public.notifications.created_by%TYPE,
    IN _status public.notifications.status%TYPE,
    IN _status_comment public.notifications.status_comment%TYPE,
    IN _recipients public.notifications.recipients%TYPE,
    IN _subject public.notifications.subject%TYPE,
    IN _body public.notifications.body%TYPE,
    IN _sent_at public.notifications.sent_at%TYPE,
    OUT err_code bigint,
    OUT err_msg text) AS $$
DECLARE
    _time timestamp(6) without time zone;
    _is_sent boolean := FALSE;
    _sent_by bigint;
BEGIN
    err_code := 0; -- NoError
    err_msg := '';

    IF _sent_at IS NOT NULL THEN
        -- notification statuses: New(1), Sending(2), SendFailed(4)
        IF _status IN (1, 2, 4) THEN
            err_code := 3; -- InvalidOperation
            err_msg := format('invalid notification status (%s), sent_at isn''t null', _status);
            RETURN;
        END IF;

        _is_sent := TRUE;
        _sent_by := _created_by;
    ELSIF _status = 3 THEN -- notification status: Sent(3)
            err_code := 3; -- InvalidOperation
            err_msg := format('sent_at is null, notification status is %s', _status);
            RETURN;
    END IF;

    IF public.notification_exists(_id) THEN
        err_code := 11001; -- NotificationAlreadyExists
        err_msg := 'notification with the same id already exists';
        RETURN;
    END IF;

    _time := (clock_timestamp() AT TIME ZONE 'UTC');
    INSERT INTO public.notifications(id, group_id, created_at, created_by, updated_at, updated_by, status, status_updated_at,
            status_updated_by, status_comment, recipients, subject, body, is_sent, sent_at, sent_by, _version_stamp, _timestamp)
        VALUES (_id, _group_id, _created_at, _created_by, _time, _created_by, _status, _time, _created_by, _status_comment,
            _recipients, _subject, _body, _is_sent, _sent_at, _sent_by, 1, _time);

    EXCEPTION
        WHEN unique_violation THEN
            IF public.notification_exists(_id) THEN
                err_code := 11001; -- NotificationAlreadyExists
                err_msg := 'notification with the same id already exists';
                RETURN;
            END IF;
            RAISE;
END;
$$ LANGUAGE plpgsql;
