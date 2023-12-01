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

-- FUNCTION: public.app_session_exists(bigint)
/*
App session statuses:
    Ended   = 3
    Deleted = 5
*/
CREATE OR REPLACE FUNCTION public.app_session_exists(
    _app_id public.app_sessions.app_id%TYPE
) RETURNS boolean AS $$
BEGIN
    -- app session statuses: Ended(3), Deleted(5)
    RETURN EXISTS (SELECT 1 FROM public.app_sessions WHERE app_id = _app_id AND status <> 3 AND status <> 5 LIMIT 1);
END;
$$ LANGUAGE plpgsql;

-- PROCEDURE: public.create_and_start_app_session(bigint, bigint, text)
/*
App statuses:
    Active = 2

App session statuses:
    Active = 2

Error codes:
    NoError          = 0
    InvalidOperation = 3
    AppNotFound      = 11000
*/
CREATE OR REPLACE PROCEDURE public.create_and_start_app_session(
    IN _app_id public.app_sessions.app_id%TYPE,
    IN _created_by public.app_sessions.created_by%TYPE,
    IN _status_comment public.app_sessions.status_comment%TYPE,
    OUT _id public.app_sessions.id%TYPE,
    OUT err_code bigint,
    OUT err_msg text) AS $$
DECLARE
    _time timestamp(6) without time zone;
    _status public.apps.status%TYPE;
BEGIN
    _id := 0;
    err_code := 0; -- NoError
    err_msg := '';

    SELECT status INTO _status FROM public.apps WHERE id = _app_id LIMIT 1 FOR UPDATE;
    IF NOT FOUND THEN
        err_code := 11000; -- AppNotFound
        err_msg := 'app not found';
        RETURN;
    END IF;

    -- app status: Active(2)
    IF _status <> 2 THEN
        err_code := 3; -- InvalidOperation
        err_msg := format('invalid app status (%s)', _status);
        RETURN;
    END IF;

    _time := (clock_timestamp() AT TIME ZONE 'UTC');
    -- app session status: Active(2)
    INSERT INTO public.app_sessions(app_id, created_at, created_by, updated_at, updated_by, status, status_updated_at, status_updated_by, status_comment, start_time, _version_stamp, _timestamp)
        VALUES (_app_id, _time, _created_by, _time, _created_by, 2, _time, _created_by, _status_comment, _time, 1, _time)
        RETURNING id INTO _id;
END;
$$ LANGUAGE plpgsql;

-- PROCEDURE: public.terminate_app_session(bigint, bigint, text)
/*
App session statuses:
    Active = 2
    Ended  = 3

Error codes:
    NoError            = 0
    InvalidOperation   = 3
    AppSessionNotFound = 11400
*/
CREATE OR REPLACE PROCEDURE public.terminate_app_session(
    IN _id public.app_sessions.id%TYPE,
    IN _updated_by public.app_sessions.updated_by%TYPE,
    IN _status_comment public.app_sessions.status_comment%TYPE,
    OUT err_code bigint,
    OUT err_msg text) AS $$
DECLARE
    _time timestamp(6) without time zone;
    _status public.app_sessions.status%TYPE;
BEGIN
    err_code := 0; -- NoError
    err_msg := '';

    SELECT status INTO _status FROM public.app_sessions WHERE id = _id LIMIT 1 FOR UPDATE;
    IF NOT FOUND THEN
        err_code := 11400; -- AppSessionNotFound
        err_msg := 'app session not found';
        RETURN;
    END IF;

    -- app session status: Active(2)
    IF _status <> 2 THEN
        err_code := 3; -- InvalidOperation
        err_msg := format('invalid app session status (%s)', _status);
        RETURN;
    END IF;

    _time := (clock_timestamp() AT TIME ZONE 'UTC');
    -- app session status: Ended(3)
    UPDATE public.app_sessions
        SET updated_at = _time, updated_by = _updated_by, status = 3, status_updated_at = _time, status_updated_by = _updated_by,
            status_comment = _status_comment, end_time = _time, _version_stamp = _version_stamp + 1, _timestamp = _time
        WHERE id = _id;
END;
$$ LANGUAGE plpgsql;
