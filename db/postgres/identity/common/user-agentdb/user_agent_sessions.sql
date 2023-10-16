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

-- PROCEDURE: public.create_and_start_user_agent_session(bigint, bigint, bigint, bigint, bigint, text, character varying)
/*
User agent statuses:
    Active = 2

User agent session statuses:
    Active = 2

Error codes:
    NoError           = 0
    InvalidOperation  = 3
    UserAgentNotFound = 11400
*/
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

    -- user agent status: Active(2)
    IF _status <> 2 THEN
        err_code := 3; -- InvalidOperation
        err_msg := format('invalid user agent status (%s)', _status);
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
END;
$$ LANGUAGE plpgsql;

-- PROCEDURE: public.terminate_user_agent_session(bigint, bigint, text)
/*
User agent session statuses:
    Active               = 2
    Ended                = 4
    LockedOut            = 5
    TemporarilyLockedOut = 6

Error codes:
    NoError                  = 0
    InvalidOperation         = 3
    UserAgentSessionNotFound = 12400
*/
CREATE OR REPLACE PROCEDURE public.terminate_user_agent_session(
    IN _id public.user_agent_sessions.id%TYPE,
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
        err_code := 12400; -- UserAgentSessionNotFound
        err_msg := 'user agent session not found';
        RETURN;
    END IF;

    -- user agent session statuses: Active(2), LockedOut(5), TemporarilyLockedOut(6)
    IF _status <> 2 AND _status <> 5 AND _status <> 6 THEN
        err_code := 3; -- InvalidOperation
        err_msg := format('invalid user agent session status (%s)', _status);
        RETURN;
    END IF;

    _time := (clock_timestamp() AT TIME ZONE 'UTC');
    -- user agent session status: Ended(4)
    UPDATE public.user_agent_sessions
        SET updated_at = _time, updated_by = _updated_by, status = 4, status_updated_at = _time, status_updated_by = _updated_by,
            status_comment = _status_comment, _version_stamp = _version_stamp + 1, _timestamp = _time
        WHERE id = _id;
END;
$$ LANGUAGE plpgsql;
