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

-- FUNCTION: public.role_exists(character varying)
/*
Role statuses:
    Deleted = 5
*/
CREATE OR REPLACE FUNCTION public.role_exists(
    _name public.roles.name%TYPE
) RETURNS boolean AS $$
BEGIN
   -- role status: Deleted(5)
    RETURN EXISTS (SELECT 1 FROM public.roles WHERE lower(name) = lower(_name) AND status <> 5 LIMIT 1);
END;
$$ LANGUAGE plpgsql;

-- PROCEDURE: public.create_role(character varying, smallint, character varying, bigint, text, bigint, bigint, text)
/*
Role statuses:
    Active = 2

Error codes:
    NoError           = 0
    RoleAlreadyExists = 11601
*/
-- Minimum transaction isolation level: Read committed.
CREATE OR REPLACE PROCEDURE public.create_role(
    IN _name public.roles.name%TYPE,
    IN _type public.roles.type%TYPE,
    IN _title public.roles.title%TYPE,
    IN _created_by public.roles.created_by%TYPE,
    IN _status_comment public.roles.status_comment%TYPE,
    IN _app_id public.roles.app_id%TYPE,
    IN _app_group_id public.roles.app_group_id%TYPE,
    IN _description public.roles.description%TYPE,
    OUT _id public.roles.id%TYPE,
    OUT err_code bigint,
    OUT err_msg text) AS $$
DECLARE
    _time timestamp(6) without time zone;
BEGIN
    _id := 0;
    err_code := 0; -- NoError
    err_msg := '';

    IF role_exists(_name) THEN
        err_code := 11601; -- RoleAlreadyExists
        err_msg := 'role with the same name already exists';
        RETURN;
    END IF;

    _time := (clock_timestamp() AT TIME ZONE 'UTC');
    -- role status: Active(2)
    INSERT INTO public.roles(name, type, title, created_at, created_by, updated_at, updated_by, status, status_updated_at, status_updated_by,
            status_comment, app_id, app_group_id, description, _version_stamp, _timestamp)
        VALUES (_name, _type, _title, _time, _created_by, _time, _created_by, 2, _time, _created_by, _status_comment, _app_id, _app_group_id, _description, 1, _time)
        RETURNING id INTO _id;

    INSERT INTO public.role_info(role_id, created_at, created_by, updated_at, updated_by, _version_stamp, _timestamp)
        VALUES (_id, _time, _created_by, _time, _created_by, 1, _time);
    
    EXCEPTION
        WHEN unique_violation THEN
            IF _id = 0 AND role_exists(_name) THEN
                err_code := 11601; -- RoleAlreadyExists
                err_msg := 'role with the same name already exists';
                RETURN;
            END IF;
            RAISE;
END;
$$ LANGUAGE plpgsql;
