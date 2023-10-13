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

-- PROCEDURE: public.create_client(bigint, smallint, text, bigint, text, character varying)
/*
Error codes:
    NoError = 0
*/
CREATE OR REPLACE PROCEDURE public.create_client(
    IN _created_by public.clients.created_by%TYPE,
    IN _status public.clients.status%TYPE,
    IN _status_comment public.clients.status_comment%TYPE,
    IN _app_id public.clients.app_id%TYPE,
    IN _user_agent public.clients.first_user_agent%TYPE,
    IN _ip public.clients.last_activity_ip%TYPE,
    OUT _id public.clients.id%TYPE,
    OUT err_code bigint,
    OUT err_msg text) AS $$
DECLARE
    _time timestamp(6) without time zone;
BEGIN
    _id := 0;
    err_code := 0; -- NoError
    err_msg := '';

    _time := (clock_timestamp() AT TIME ZONE 'UTC');
    INSERT INTO public.clients(created_at, created_by, updated_at, updated_by, status, status_updated_at, status_updated_by, status_comment,
            app_id, first_user_agent, last_user_agent, last_activity_time, last_activity_ip, _version_stamp, _timestamp)
        VALUES (_time, _created_by, _time, _created_by, _status, _time, _created_by, _status_comment, _app_id, _user_agent, _user_agent, _time, _ip, 1, _time)
        RETURNING id INTO _id;
END;
$$ LANGUAGE plpgsql;
