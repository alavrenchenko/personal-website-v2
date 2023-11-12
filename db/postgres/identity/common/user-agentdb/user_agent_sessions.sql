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

-- FUNCTION: public.user_agent_session_exists(bigint, bigint, bigint)
/*
User agent session statuses:
    Deleted = 9
*/
CREATE OR REPLACE FUNCTION public.user_agent_session_exists(
    _user_id public.user_agent_sessions.user_id%TYPE,
    _client_id public.user_agent_sessions.client_id%TYPE,
    _user_agent_id public.user_agent_sessions.user_agent_id%TYPE
) RETURNS boolean AS $$
BEGIN
    -- user agent session status: Deleted(9)
    RETURN EXISTS (SELECT 1 FROM public.user_agent_sessions WHERE (user_id = _user_id AND client_id = _client_id OR user_agent_id = _user_agent_id) AND status <> 9 LIMIT 1);
END;
$$ LANGUAGE plpgsql;

-- PROCEDURE: public.create_and_start_user_agent_session(bigint, bigint, bigint, bigint, bigint, text, character varying)
/*
User agent statuses:
    Active = 3

User agent session statuses:
    Active = 2

Error codes:
    NoError                       = 0
    InvalidOperation              = 3
    UserAgentNotFound             = 11400
    UserAgentSessionAlreadyExists = 12602
*/
-- Minimum transaction isolation level: Read committed.
CREATE OR REPLACE PROCEDURE public.create_and_start_user_agent_session(
    IN _user_id public.user_agent_sessions.user_id%TYPE,
    IN _client_id public.user_agent_sessions.client_id%TYPE,
    IN _user_agent_id public.user_agent_sessions.user_agent_id%TYPE,
    IN _user_session_id public.user_agent_sessions.user_session_id%TYPE,
    IN _created_by public.user_agent_sessions.created_by%TYPE,
    IN _status_comment public.user_agent_sessions.status_comment%TYPE,
    IN _ip public.user_agent_sessions.first_sign_in_ip%TYPE,
    OUT _id public.user_agent_sessions.id%TYPE,
    OUT err_code bigint,
    OUT err_msg text) AS $$
DECLARE
    _time timestamp(6) without time zone;
    _status public.user_agents.status%TYPE;
BEGIN
    _id := 0;
    err_code := 0; -- NoError
    err_msg := '';

    SELECT status INTO _status FROM public.user_agents WHERE id = _user_agent_id LIMIT 1 FOR UPDATE;
    IF NOT FOUND THEN
        err_code := 11400; -- UserAgentNotFound
        err_msg := 'user agent not found';
        RETURN;
    END IF;

    -- user agent status: Active(3)
    IF _status <> 3 THEN
        err_code := 3; -- InvalidOperation
        err_msg := format('invalid user agent status (%s)', _status);
        RETURN;
    END IF;

    IF public.user_agent_session_exists(_user_id, _client_id, _user_agent_id) THEN
        err_code := 12602; -- UserAgentSessionAlreadyExists
        err_msg := 'user agent session with the same params already exists';
        RETURN;
    END IF;

    _time := (clock_timestamp() AT TIME ZONE 'UTC');
    -- user agent session status: Active(2)
    INSERT INTO public.user_agent_sessions(user_id, client_id, user_agent_id, user_session_id, created_at, created_by, updated_at, updated_by, status,
            status_updated_at, status_updated_by, status_comment, first_sign_in_time, first_sign_in_ip, last_sign_in_time, last_sign_in_ip,
            last_activity_time, last_activity_ip, _version_stamp, _timestamp)
        VALUES (_user_id, _client_id, _user_agent_id, user_session_id, _time, _created_by, _time, _created_by, 2, _time, _created_by, _status_comment,
            _time, _ip, _time, _ip, _time, _ip, 1, _time)
        RETURNING id INTO _id;

    EXCEPTION
        WHEN unique_violation THEN
            IF public.user_agent_session_exists(_user_id, _client_id, _user_agent_id) THEN
                err_code := 12602; -- UserAgentSessionAlreadyExists
                err_msg := 'user agent session with the same params already exists';
                RETURN;
            END IF;
            RAISE;
END;
$$ LANGUAGE plpgsql;

-- PROCEDURE: public.start_user_agent_session(bigint, bigint, text, character varying, bigint)
/*
User agent statuses:
    Active = 3

User agent session statuses:
    New       = 1
    Active    = 2
    SignedOut = 3
    Ended     = 4

Error codes:
    NoError                  = 0
    InvalidOperation         = 3
    UserAgentNotFound        = 11400
    UserAgentSessionNotFound = 12600
*/
-- Minimum transaction isolation level: Read committed.
CREATE OR REPLACE PROCEDURE public.start_user_agent_session(
    IN _id public.user_agent_sessions.id%TYPE,
    IN _user_session_id public.user_agent_sessions.user_session_id%TYPE,
    IN _status_comment public.user_agent_sessions.status_comment%TYPE,
    IN _ip public.user_agent_sessions.last_sign_in_ip%TYPE,
    IN _operation_user_id public.user_agent_sessions.updated_by%TYPE,
    OUT err_code bigint,
    OUT err_msg text) AS $$
DECLARE
    _time timestamp(6) without time zone;
    _ua_id public.user_agents.id%TYPE;
    _uas_status public.user_agent_sessions.status%TYPE;
    _ua_status public.user_agents.status%TYPE;
    _ua_first_sign_in_time public.user_agents.status%TYPE;
BEGIN
    err_code := 0; -- NoError
    err_msg := '';

    SELECT user_agent_id, status INTO _ua_id, _uas_status FROM public.user_agent_sessions WHERE id = _id LIMIT 1 FOR UPDATE;
    IF NOT FOUND THEN
        err_code := 12600; -- UserAgentSessionNotFound
        err_msg := 'user agent session not found';
        RETURN;
    END IF;

    -- user agent session status: New(1), SignedOut(3), Ended(4)
    IF _uas_status <> 1 AND _uas_status <> 3 AND _uas_status <> 4 THEN
        err_code := 3; -- InvalidOperation
        err_msg := format('invalid user agent session status (%s)', _uas_status);
        RETURN;
    END IF;

    SELECT status, first_sign_in_time INTO _ua_status, _ua_first_sign_in_time FROM public.user_agents WHERE id = _ua_id LIMIT 1 FOR UPDATE;
    IF NOT FOUND THEN
        err_code := 11400; -- UserAgentNotFound
        err_msg := 'user agent not found';
        RETURN;
    END IF;

    -- user agent status: Active(3)
    IF _status <> 3 THEN
        err_code := 3; -- InvalidOperation
        err_msg := format('invalid user agent status (%s)', _status);
        RETURN;
    END IF;

    _time := (clock_timestamp() AT TIME ZONE 'UTC');
    -- user agent session status: New(1)
    IF _uas_status = 1 THEN
        -- user agent session status: Active(2)
        UPDATE public.user_agent_sessions
            SET user_session_id = _user_session_id, updated_at = _time, updated_by = _operation_user_id, status = 2, status_updated_at = _time,
                status_updated_by = _operation_user_id, status_comment = _status_comment, first_sign_in_time = _time, first_sign_in_ip = _ip,
                last_sign_in_time = _time, last_sign_in_ip = _ip, last_activity_time = _time, last_activity_ip = _ip,
                _version_stamp = _version_stamp + 1, _timestamp = _time
            WHERE id = _id;
    ELSE
        -- user agent session status: Active(2)
        UPDATE public.user_agent_sessions
            SET user_session_id = _user_session_id, updated_at = _time, updated_by = _operation_user_id, status = 2, status_updated_at = _time,
                status_updated_by = _operation_user_id, status_comment = _status_comment, last_sign_in_time = _time, last_sign_in_ip = _ip,
                last_activity_time = _time, last_activity_ip = _ip, _version_stamp = _version_stamp + 1, _timestamp = _time
            WHERE id = _id;
    END IF;

    IF _ua_first_sign_in_time IS NOT NULL THEN
        UPDATE public.user_agents
            SET updated_at = _time, updated_by = _operation_user_id, last_sign_in_time = _time, last_sign_in_ip = _ip, last_activity_time = _time,
                last_activity_ip = _ip, _version_stamp = _version_stamp + 1, _timestamp = _time
            WHERE id = _ua_id;
    ELSE
        UPDATE public.user_agents
            SET updated_at = _time, updated_by = _operation_user_id, first_sign_in_time = _time, first_sign_in_ip = _ip, last_sign_in_time = _time,
                last_sign_in_ip = _ip, last_activity_time = _time, last_activity_ip = _ip, _version_stamp = _version_stamp + 1, _timestamp = _time
            WHERE id = _ua_id;
    END IF;
END;
$$ LANGUAGE plpgsql;

-- PROCEDURE: public.terminate_user_agent_session(bigint, boolean, bigint, text)
/*
User agent session statuses:
    Active    = 2
    SignedOut = 3
    Ended     = 4

Error codes:
    NoError                  = 0
    InvalidOperation         = 3
    UserAgentSessionNotFound = 12600
*/
CREATE OR REPLACE PROCEDURE public.terminate_user_agent_session(
    IN _id public.user_agent_sessions.id%TYPE,
    IN _signed_out boolean,
    IN _updated_by public.user_agent_sessions.updated_by%TYPE,
    IN _status_comment public.user_agent_sessions.status_comment%TYPE,
    OUT err_code bigint,
    OUT err_msg text) AS $$
DECLARE
    _time timestamp(6) without time zone;
    _status public.user_agent_sessions.status%TYPE;
BEGIN
    err_code := 0; -- NoError
    err_msg := '';

    SELECT status INTO _status FROM public.user_agent_sessions WHERE id = _id LIMIT 1 FOR UPDATE;
    IF NOT FOUND THEN
        err_code := 12600; -- UserAgentSessionNotFound
        err_msg := 'user agent session not found';
        RETURN;
    END IF;

    -- user agent session status: Active(2)
    IF _status <> 2 THEN
        err_code := 3; -- InvalidOperation
        err_msg := format('invalid user agent session status (%s)', _status);
        RETURN;
    END IF;

    IF _signed_out THEN
        -- user agent session status: SignedOut(3)
        _status := 3;
    ELSE
        -- user agent session status: Ended(4)
        _status := 4;
    END IF;

    _time := (clock_timestamp() AT TIME ZONE 'UTC');
    UPDATE public.user_agent_sessions
        SET updated_at = _time, updated_by = _updated_by, status = _status, status_updated_at = _time, status_updated_by = _updated_by,
            status_comment = _status_comment, _version_stamp = _version_stamp + 1, _timestamp = _time
        WHERE id = _id;
END;
$$ LANGUAGE plpgsql;

-- PROCEDURE: public.start_deleting_user_agent_session(bigint, bigint, text)
/*
User agent session statuses:
    Deleting = 8
    Deleted  = 9

Error codes:
    NoError                  = 0
    InvalidOperation         = 3
    UserAgentSessionNotFound = 12600
*/
-- Minimum transaction isolation level: Read committed.
CREATE OR REPLACE PROCEDURE public.start_deleting_user_agent_session(
    IN _id public.user_agent_sessions.id%TYPE,
    IN _deleted_by public.user_agent_sessions.updated_by%TYPE,
    IN _status_comment public.user_agent_sessions.status_comment%TYPE,
    OUT err_code bigint,
    OUT err_msg text) AS $$
DECLARE
    _time timestamp(6) without time zone;
    _status public.user_agent_sessions.status%TYPE;
BEGIN
    err_code := 0; -- NoError
    err_msg := '';

    SELECT status INTO _status FROM public.user_agent_sessions WHERE id = _id LIMIT 1 FOR UPDATE;
    IF NOT FOUND THEN
        err_code := 12600; -- UserAgentSessionNotFound
        err_msg := 'user agent session not found';
        RETURN;
    END IF;

    -- user agent session statuses: Deleting(8), Deleted(9)
    IF _status = 8 OR _status = 9 THEN
        err_code := 3; -- InvalidOperation
        err_msg := format('invalid user agent session status (%s)', _status);
        RETURN;
    END IF;

    _time := (clock_timestamp() AT TIME ZONE 'UTC');
    -- user agent session status: Deleting(8)
    UPDATE public.user_agent_sessions
        SET updated_at = _time, updated_by = _deleted_by, status = 8, status_updated_at = _time, status_updated_by = _deleted_by,
            status_comment = _status_comment, _version_stamp = _version_stamp + 1, _timestamp = _time
        WHERE id = _id;
END;
$$ LANGUAGE plpgsql;
