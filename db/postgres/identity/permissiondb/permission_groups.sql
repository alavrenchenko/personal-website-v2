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

-- FUNCTION: public.permission_group_exists(character varying)
/*
Permission group statuses:
    Deleted = 5
*/
CREATE OR REPLACE FUNCTION public.permission_group_exists(
    _name public.permission_groups.name%TYPE
) RETURNS boolean AS $$
BEGIN
   -- permission status: Deleted(5)
    RETURN EXISTS (SELECT 1 FROM public.permission_groups WHERE lower(name) = lower(_name) AND status <> 5 LIMIT 1);
END;
$$ LANGUAGE plpgsql;

-- PROCEDURE: public.create_permission_group(character varying, bigint, text, bigint, bigint, text)
/*
Permission group statuses:
    Active = 2

Error codes:
    NoError                      = 0
    PermissionGroupAlreadyExists = 12001
*/
CREATE OR REPLACE PROCEDURE public.create_permission_group(
    IN _name public.permission_groups.name%TYPE,
    IN _created_by public.permission_groups.created_by%TYPE,
    IN _status_comment public.permission_groups.status_comment%TYPE,
    IN _app_id public.permission_groups.app_id%TYPE,
    IN _app_group_id public.permission_groups.app_group_id%TYPE,
    IN _description public.permission_groups.description%TYPE,
    OUT _id public.permission_groups.id%TYPE,
    OUT err_code bigint,
    OUT err_msg text) AS $$
DECLARE
    _time timestamp(6) without time zone;
BEGIN
    _id := 0;
    err_code := 0; -- NoError
    err_msg := '';

    IF public.permission_group_exists(_name) THEN
        err_code := 12001; -- PermissionGroupAlreadyExists
        err_msg := 'permission group with the same name already exists';
        RETURN;
    END IF;

    _time := (clock_timestamp() AT TIME ZONE 'UTC');
    -- permission group status: Active(2)
    INSERT INTO public.permission_groups(name, created_at, created_by, updated_at, updated_by, status, status_updated_at, status_updated_by,
            status_comment, app_id, app_group_id, description, _version_stamp, _timestamp)
        VALUES (_name, _time, _created_by, _time, _created_by, 2, _time, _created_by, _status_comment, _app_id, _app_group_id, _description, 1, _time)
        RETURNING id INTO _id;

    EXCEPTION
        WHEN unique_violation THEN
            IF public.permission_group_exists(_name) THEN
                err_code := 12001; -- PermissionGroupAlreadyExists
                err_msg := 'permission group with the same name already exists';
                RETURN;
            END IF;
            RAISE;
END;
$$ LANGUAGE plpgsql;
