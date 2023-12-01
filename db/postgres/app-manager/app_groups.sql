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

-- FUNCTION: public.app_group_exists(character varying)
/*
App group statuses:
    Deleted = 5
*/
CREATE OR REPLACE FUNCTION public.app_group_exists(
    _name public.app_groups.name%TYPE
) RETURNS boolean AS $$
BEGIN
   -- app group status: Deleted(5)
    RETURN EXISTS (SELECT 1 FROM public.app_groups WHERE lower(name) = lower(_name) AND status <> 5 LIMIT 1);
END;
$$ LANGUAGE plpgsql;

-- PROCEDURE: public.create_app_group(character varying, smallint, character varying, bigint, text, character varying, text)
/*
App group statuses:
    Active = 2

Error codes:
    NoError               = 0
    AppGroupAlreadyExists = 11201
*/
-- Minimum transaction isolation level: Read committed.
CREATE OR REPLACE PROCEDURE public.create_app_group(
    IN _name public.app_groups.name%TYPE,
    IN _type public.app_groups.type%TYPE,
    IN _title public.app_groups.title%TYPE,
    IN _created_by public.app_groups.created_by%TYPE,
    IN _status_comment public.app_groups.status_comment%TYPE,
    IN _version public.app_groups.version%TYPE,
    IN _description public.app_groups.description%TYPE,
    OUT _id public.app_groups.id%TYPE,
    OUT err_code bigint,
    OUT err_msg text) AS $$
DECLARE
    _time timestamp(6) without time zone;
BEGIN
    _id := 0;
    err_code := 0; -- NoError
    err_msg := '';

    IF public.app_group_exists(_name) THEN
        err_code := 11201; -- AppGroupAlreadyExists
        err_msg := 'app group with the same name already exists';
        RETURN;
    END IF;

    _time := (clock_timestamp() AT TIME ZONE 'UTC');
    -- app group status: Active(2)
    INSERT INTO public.app_groups(name, type, title, created_at, created_by, updated_at, updated_by, status, status_updated_at, status_updated_by,
            status_comment, version, description, _version_stamp, _timestamp)
        VALUES (_name, _type, _title, _time, _created_by, _time, _created_by, 2, _time, _created_by, _status_comment, _version, _description, 1, _time)
        RETURNING id INTO _id;

    EXCEPTION
        WHEN unique_violation THEN
            IF public.app_group_exists(_name) THEN
                err_code := 11201; -- AppGroupAlreadyExists
                err_msg := 'app group with the same name already exists';
                RETURN;
            END IF;
            RAISE;
END;
$$ LANGUAGE plpgsql;

-- PROCEDURE: public.delete_app_group(bigint, bigint, text)
/*
App group statuses:
    Deleted = 5

App statuses:
    Deleted = 5

Error codes:
    NoError          = 0
    InvalidOperation = 3
    AppGroupNotFound = 11200
*/
-- Minimum transaction isolation level: Read committed.
CREATE OR REPLACE PROCEDURE public.delete_app_group(
    IN _id public.app_groups.id%TYPE,
    IN _deleted_by public.app_groups.updated_by%TYPE,
    IN _status_comment public.app_groups.status_comment%TYPE,
    OUT err_code bigint,
    OUT err_msg text) AS $$
DECLARE
    _time timestamp(6) without time zone;
    _status public.app_groups.status%TYPE;
BEGIN
    err_code := 0; -- NoError
    err_msg := '';

    SELECT status INTO _status FROM public.app_groups WHERE id = _id LIMIT 1 FOR UPDATE;
    IF NOT FOUND THEN
        err_code := 11200; -- AppGroupNotFound
        err_msg := 'app group not found';
        RETURN;
    END IF;

    -- app group status: Deleted(5)
    IF _status = 5 THEN
        err_code := 3; -- InvalidOperation
        err_msg := 'app group has already been deleted';
        RETURN;
    END IF;

    -- app status: Deleted(5)
    IF EXISTS (SELECT 1 FROM public.apps WHERE group_id = _id AND status <> 5 LIMIT 1) THEN
        err_code := 3; -- InvalidOperation
        err_msg := 'app group contains apps';
        RETURN;
    END IF;

    _time := (clock_timestamp() AT TIME ZONE 'UTC');
    -- app group status: Deleted(5)
    UPDATE public.app_groups
        SET updated_at = _time, updated_by = _deleted_by, status = 5, status_updated_at = _time, status_updated_by = _deleted_by,
            status_comment = _status_comment, _version_stamp = _version_stamp + 1, _timestamp = _time
        WHERE id = _id;
END;
$$ LANGUAGE plpgsql;
