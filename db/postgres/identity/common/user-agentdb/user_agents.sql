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

-- FUNCTION: public.user_agent_exists(bigint, bigint)
/*
User agent statuses:
    Deleted = 8
*/
CREATE OR REPLACE FUNCTION public.user_agent_exists(
    _user_id public.user_agents.user_id%TYPE,
    _client_id public.user_agents.client_id%TYPE
) RETURNS boolean AS $$
BEGIN
   -- user agent status: Deleted(8)
    RETURN EXISTS (SELECT 1 FROM public.user_agents WHERE user_id = _user_id AND client_id = _client_id AND status <> 8 LIMIT 1);
END;
$$ LANGUAGE plpgsql;

-- PROCEDURE: public.create_user_agent(bigint, bigint, bigint, smallint, text, bigint, text)
/*
Error codes:
    NoError                = 0
    UserAgentAlreadyExists = 12202
*/
-- Minimum transaction isolation level: Read committed.
CREATE OR REPLACE PROCEDURE public.create_user_agent(
    IN _user_id public.user_agents.user_id%TYPE,
    IN _client_id public.user_agents.client_id%TYPE,
    IN _created_by public.user_agents.created_by%TYPE,
    IN _status public.user_agents.status%TYPE,
    IN _status_comment public.user_agents.status_comment%TYPE,
    IN _app_id public.user_agents.app_id%TYPE,
    IN _user_agent public.user_agents.first_user_agent%TYPE,
    OUT _id public.user_agents.id%TYPE,
    OUT err_code bigint,
    OUT err_msg text) AS $$
DECLARE
    _time timestamp(6) without time zone;
BEGIN
    _id := 0;
    err_code := 0; -- NoError
    err_msg := '';

    IF public.user_agent_exists(_user_id, _client_id) THEN
        err_code := 12202; -- UserAgentAlreadyExists
        err_msg := 'user agent with the same params already exists';
        RETURN;
    END IF;

    _time := (clock_timestamp() AT TIME ZONE 'UTC');
    INSERT INTO public.user_agents(user_id, client_id, created_at, created_by, updated_at, updated_by, status, status_updated_at,
            status_updated_by, status_comment, app_id, first_user_agent, last_user_agent, _version_stamp, _timestamp)
        VALUES (_user_id, _client_id, _time, _created_by, _time, _created_by, _status, _time, _created_by, _status_comment,
            _app_id, _user_agent, _user_agent, 1, _time)
        RETURNING id INTO _id;

    EXCEPTION
        WHEN unique_violation THEN
            IF public.user_agent_exists(_user_id, _client_id) THEN
                err_code := 12202; -- UserAgentAlreadyExists
                err_msg := 'user agent with the same params already exists';
                RETURN;
            END IF;
            RAISE;
END;
$$ LANGUAGE plpgsql;

-- PROCEDURE: public.start_deleting_user_agent(bigint, bigint, text)
/*
User agent statuses:
    Deleting = 7
    Deleted  = 8

Error codes:
    NoError           = 0
    InvalidOperation  = 3
    UserAgentNotFound = 12200
*/
-- Minimum transaction isolation level: Read committed.
CREATE OR REPLACE PROCEDURE public.start_deleting_user_agent(
    IN _id public.user_agents.id%TYPE,
    IN _deleted_by public.user_agents.updated_by%TYPE,
    IN _status_comment public.user_agents.status_comment%TYPE,
    OUT err_code bigint,
    OUT err_msg text) AS $$
DECLARE
    _time timestamp(6) without time zone;
    _status public.user_agents.status%TYPE;
BEGIN
    err_code := 0; -- NoError
    err_msg := '';

    SELECT status INTO _status FROM public.user_agents WHERE id = _id LIMIT 1 FOR UPDATE;
    IF NOT FOUND THEN
        err_code := 12200; -- UserAgentNotFound
        err_msg := 'user agent not found';
        RETURN;
    END IF;

    -- user agent statuses: Deleting(7), Deleted(8)
    IF _status = 7 OR _status = 8 THEN
        err_code := 3; -- InvalidOperation
        err_msg := format('invalid user agent status (%s)', _status);
        RETURN;
    END IF;

    _time := (clock_timestamp() AT TIME ZONE 'UTC');
    -- user agent status: Deleting(7)
    UPDATE public.user_agents
        SET updated_at = _time, updated_by = _deleted_by, status = 7, status_updated_at = _time, status_updated_by = _deleted_by,
            status_comment = _status_comment, _version_stamp = _version_stamp + 1, _timestamp = _time
        WHERE id = _id;
END;
$$ LANGUAGE plpgsql;

-- PROCEDURE: public.delete_user_agent(bigint, bigint, text)
/*
User agent statuses:
    Deleting = 7
    Deleted  = 8

User agent session statuses:
    Deleted = 9

Error codes:
    NoError            = 0
    InvalidOperation   = 3
    UserAgentNotFound = 12200
*/
-- Minimum transaction isolation level: Read committed.
CREATE OR REPLACE PROCEDURE public.delete_user_agent(
    IN _id public.user_agents.id%TYPE,
    IN _deleted_by public.user_agents.updated_by%TYPE,
    IN _status_comment public.user_agents.status_comment%TYPE,
    OUT err_code bigint,
    OUT err_msg text) AS $$
DECLARE
    _time timestamp(6) without time zone;
    _status public.user_agents.status%TYPE;
BEGIN
    err_code := 0; -- NoError
    err_msg := '';

    SELECT status INTO _status FROM public.user_agents WHERE id = _id LIMIT 1 FOR UPDATE;
    IF NOT FOUND THEN
        err_code := 12200; -- UserAgentNotFound
        err_msg := 'user agent not found';
        RETURN;
    END IF;

    -- user agent status: Deleting(7)
    IF _status <> 7 THEN
        err_code := 3; -- InvalidOperation
        -- user agent status: Deleted(8)
        IF _status = 8 THEN
            err_msg := 'user agent has already been deleted';
        ELSE
            err_msg := format('invalid user agent status (%s)', _status);
        END IF;
        RETURN;
    END IF;

    -- user agent session status: Deleted(9)
    IF EXISTS (SELECT 1 FROM public.user_agent_sessions WHERE user_agent_id = _id AND status <> 9 LIMIT 1) THEN
        err_code := 3; -- InvalidOperation
        err_msg := 'user agent session exists';
        RETURN;
    END IF;

    _time := (clock_timestamp() AT TIME ZONE 'UTC');
    -- user agent status: Deleted(8)
    UPDATE public.user_agents
        SET updated_at = _time, updated_by = _deleted_by, status = 8, status_updated_at = _time, status_updated_by = _deleted_by,
            status_comment = _status_comment, _version_stamp = _version_stamp + 1, _timestamp = _time
        WHERE id = _id;
END;
$$ LANGUAGE plpgsql;
