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
