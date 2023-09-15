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

-- PROCEDURE: public.create_and_start_logging_session(bigint, bigint, text)
/*
Logging session statuses:
    Started = 2

Error codes:
    NoError = 0
*/
CREATE OR REPLACE PROCEDURE public.create_and_start_logging_session(
    IN _app_id public.logging_sessions.app_id%TYPE,
    IN _created_by public.logging_sessions.created_by%TYPE,
    IN _status_comment public.logging_sessions.status_comment%TYPE,
    OUT _id public.logging_sessions.id%TYPE,
    OUT err_code bigint,
    OUT err_msg text) AS $$
DECLARE
    _time timestamp(6) without time zone;
BEGIN
    _id := 0;
    err_code := 0; -- NoError
    err_msg := '';

    _time := (clock_timestamp() AT TIME ZONE 'UTC');
    -- logging session status: Started(2)
    INSERT INTO public.logging_sessions(app_id, created_at, created_by, updated_at, updated_by, status, status_updated_at, status_updated_by, status_comment, start_time, _version_stamp, _timestamp)
        VALUES (_app_id, _time, _created_by, _time, _created_by, 2, _time, _created_by, _status_comment, _time, 1, _time)
        RETURNING id INTO _id;
END;
$$ LANGUAGE plpgsql;
