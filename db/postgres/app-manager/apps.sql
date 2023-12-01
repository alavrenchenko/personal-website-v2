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

-- FUNCTION: public.app_exists(character varying)
/*
App statuses:
    Deleted = 5
*/
CREATE OR REPLACE FUNCTION public.app_exists(
    _name public.apps.name%TYPE
) RETURNS boolean AS $$
BEGIN
   -- app status: Deleted(5)
    RETURN EXISTS (SELECT 1 FROM public.apps WHERE lower(name) = lower(_name) AND status <> 5 LIMIT 1);
END;
$$ LANGUAGE plpgsql;

-- PROCEDURE: public.create_app(character varying, bigint, smallint, character varying, smallint, bigint, text, character varying, text)
/*
App statuses:
    Active = 2

App group statuses:
    Active = 2

Error codes:
    NoError          = 0
    AppAlreadyExists = 11001
    AppGroupNotFound = 11200
*/
-- Minimum transaction isolation level: Read committed.
CREATE OR REPLACE PROCEDURE public.create_app(
    IN _name public.apps.name%TYPE,
    IN _group_id public.apps.group_id%TYPE,
    IN _type public.apps.type%TYPE,
    IN _title public.apps.title%TYPE,
    IN _category public.apps.category%TYPE,
    IN _created_by public.apps.created_by%TYPE,
    IN _status_comment public.apps.status_comment%TYPE,
    IN _version public.apps.version%TYPE,
    IN _description public.apps.description%TYPE,
    OUT _id public.apps.id%TYPE,
    OUT err_code bigint,
    OUT err_msg text) AS $$
DECLARE
    _time timestamp(6) without time zone;
    _status public.app_groups.status%TYPE;
BEGIN
    _id := 0;
    err_code := 0; -- NoError
    err_msg := '';

    IF public.app_exists(_name) THEN
        err_code := 11001; -- AppAlreadyExists
        err_msg := 'app with the same name already exists';
        RETURN;
    END IF;

    SELECT status INTO _status FROM public.app_groups WHERE id = _group_id LIMIT 1 FOR UPDATE;
    IF NOT FOUND THEN
        err_code := 11200; -- AppGroupNotFound
        err_msg := 'app group not found';
        RETURN;
    END IF;

    -- app group status: Active(2)
    IF _status <> 2 THEN
        err_code := 3; -- InvalidOperation
        err_msg := format('invalid app group status (%s)', _status);
        RETURN;
    END IF;

    _time := (clock_timestamp() AT TIME ZONE 'UTC');
    -- app status: Active(2)
    INSERT INTO public.apps(name, group_id, type, title, category, created_at, created_by, updated_at, updated_by, status, status_updated_at,
            status_updated_by, status_comment, version, description, _version_stamp, _timestamp)
        VALUES (_name, _group_id, _type, _title, _category, _time, _created_by, _time, _created_by, 2, _time, _created_by, _status_comment,
            _version, _description, 1, _time)
        RETURNING id INTO _id;

    EXCEPTION
        WHEN unique_violation THEN
            IF public.app_exists(_name) THEN
                err_code := 11001; -- AppAlreadyExists
                err_msg := 'app with the same name already exists';
                RETURN;
            END IF;
            RAISE;
END;
$$ LANGUAGE plpgsql;
