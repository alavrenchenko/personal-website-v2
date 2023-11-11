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

-- FUNCTION: public.user_session_exists(bigint, bigint, bigint)
/*
User session statuses:
    Ended   = 3
    Deleted = 5
*/
CREATE OR REPLACE FUNCTION public.user_session_exists(
    _user_id public.user_sessions.user_id%TYPE,
    _client_id public.user_sessions.client_id%TYPE,
    _user_agent_id public.user_sessions.user_agent_id%TYPE
) RETURNS boolean AS $$
BEGIN
    -- user's session status: Ended(3), Deleted(5)
    RETURN EXISTS (SELECT 1 FROM public.user_sessions WHERE (user_id = _user_id AND client_id = _client_id OR user_agent_id = _user_agent_id) AND status <> 3 AND status <> 5 LIMIT 1);
END;
$$ LANGUAGE plpgsql;

-- PROCEDURE: public.create_and_start_user_session(bigint, bigint, bigint, bigint, text, character varying)
/*
User session statuses:
    Active = 2

Error codes:
    NoError                  = 0
    UserSessionAlreadyExists = 12402
*/
-- Minimum transaction isolation level: Read committed.
CREATE OR REPLACE PROCEDURE public.create_and_start_user_session(
    IN _user_id public.user_sessions.user_id%TYPE,
    IN _client_id public.user_sessions.client_id%TYPE,
    IN _user_agent_id public.user_sessions.user_agent_id%TYPE,
    IN _created_by public.user_sessions.created_by%TYPE,
    IN _status_comment public.user_sessions.status_comment%TYPE,
    IN _app_id public.user_sessions.app_id%TYPE,
    IN _first_ip public.user_sessions.first_ip%TYPE,
    OUT _id public.user_sessions.id%TYPE,
    OUT err_code bigint,
    OUT err_msg text) AS $$
DECLARE
    _time timestamp(6) without time zone;
BEGIN
    _id := 0;
    err_code := 0; -- NoError
    err_msg := '';

    IF public.user_session_exists(_user_id, _client_id, _user_agent_id) THEN
        err_code := 12402; -- UserSessionAlreadyExists
        err_msg := 'user''s session with the same params already exists';
        RETURN;
    END IF;

    _time := (clock_timestamp() AT TIME ZONE 'UTC');
    -- user's session status: Active(2)
    INSERT INTO public.user_sessions(user_id, client_id, user_agent_id, created_at, created_by, updated_at, updated_by, status, status_updated_at,
            status_updated_by, status_comment, app_id, start_time, first_ip, last_activity_time, last_activity_ip, _version_stamp, _timestamp)
        VALUES (_user_id, _client_id, _user_agent_id, _time, _created_by, _time, _created_by, 2, _time, _created_by, _status_comment,
            _app_id, _time, _first_ip, _time, _first_ip, 1, _time)
        RETURNING id INTO _id;

    EXCEPTION
        WHEN unique_violation THEN
            IF public.user_session_exists(_user_id, _client_id, _user_agent_id) THEN
                err_code := 12402; -- UserSessionAlreadyExists
                err_msg := 'user''s session with the same params already exists';
                RETURN;
            END IF;
            RAISE;
END;
$$ LANGUAGE plpgsql;

-- PROCEDURE: public.terminate_user_session(bigint, bigint, text)
/*
User session statuses:
    Active = 2
    Ended  = 3

Error codes:
    NoError            = 0
    InvalidOperation   = 3
    UserSessionNotFound = 12400
*/
-- Minimum transaction isolation level: Read committed.
CREATE OR REPLACE PROCEDURE public.terminate_user_session(
    IN _id public.user_sessions.id%TYPE,
    IN _updated_by public.user_sessions.updated_by%TYPE,
    IN _status_comment public.user_sessions.status_comment%TYPE,
    OUT err_code bigint,
    OUT err_msg text) AS $$
DECLARE
    _time timestamp(6) without time zone;
    _status public.user_sessions.status%TYPE;
BEGIN
    err_code := 0; -- NoError
    err_msg := '';

    SELECT status INTO _status FROM public.user_sessions WHERE id = _id LIMIT 1 FOR UPDATE;
    IF NOT FOUND THEN
        err_code := 12400; -- UserSessionNotFound
        err_msg := 'user''s session not found';
        RETURN;
    END IF;

    -- user's session status: Active(2)
    IF _status <> 2 THEN
        err_code := 3; -- InvalidOperation
        err_msg := format('invalid user session status (%s)', _status);
        RETURN;
    END IF;

    _time := (clock_timestamp() AT TIME ZONE 'UTC');
    -- user's session status: Ended(3)
    UPDATE public.user_sessions
        SET updated_at = _time, updated_by = _updated_by, status = 3, status_updated_at = _time, status_updated_by = _updated_by,
            status_comment = _status_comment, end_time = _time, _version_stamp = _version_stamp + 1, _timestamp = _time
        WHERE id = _id;
END;
$$ LANGUAGE plpgsql;
