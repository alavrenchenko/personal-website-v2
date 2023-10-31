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

    IF public.role_exists(_name) THEN
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
            IF _id = 0 AND public.role_exists(_name) THEN
                err_code := 11601; -- RoleAlreadyExists
                err_msg := 'role with the same name already exists';
                RETURN;
            END IF;
            RAISE;
END;
$$ LANGUAGE plpgsql;

-- PROCEDURE: public.delete_role(bigint, bigint, text)
/*
Role assignment statuses:
    Deleted = 5

Error codes:
    NoError          = 0
    InternalError    = 2
    InvalidOperation = 3
    RoleNotFound     = 11600
    RoleInfoNotFound = 11602
*/
-- Minimum transaction isolation level: Read committed.
CREATE OR REPLACE PROCEDURE public.delete_role(
    IN _id public.roles.id%TYPE,
    IN _deleted_by public.roles.updated_by%TYPE,
    IN _status_comment public.roles.status_comment%TYPE,
    OUT err_code bigint,
    OUT err_msg text) AS $$
DECLARE
    _time timestamp(6) without time zone;
    _status public.roles.status%TYPE;
    _is_role_info_deleted public.role_info.is_deleted%TYPE;
    _existing_assignment_count public.role_info.existing_assignment_count%TYPE;
BEGIN
    err_code := 0; -- NoError
    err_msg := '';

    SELECT status INTO _status FROM public.roles WHERE id = _id LIMIT 1 FOR UPDATE;
    IF NOT FOUND THEN
        err_code := 11600; -- RoleNotFound
        err_msg := 'role not found';
        RETURN;
    END IF;

    -- role status: Deleted(5)
    IF _status = 5 THEN
        err_code := 3; -- InvalidOperation
        err_msg := 'role has already been deleted';
        RETURN;
    END IF;

    IF EXISTS (SELECT 1 FROM public.new_role_assignments WHERE role_id = _id LIMIT 1) THEN
        err_code := 3; -- InvalidOperation
        err_msg := 'role is assigned';
        RETURN;
    END IF;

    SELECT is_deleted, existing_assignment_count INTO _is_role_info_deleted, _existing_assignment_count FROM public.role_info WHERE role_id = _role_id LIMIT 1 FOR UPDATE;
    IF NOT FOUND THEN
        -- internal error
        err_code := 11602; -- RoleInfoNotFound
        err_msg := 'role info not found';
        RETURN;
    END IF;

    IF _is_role_info_deleted THEN
        err_code := 2; -- InternalError
        err_msg := 'role info is deleted';
        RETURN;
    END IF;

    IF _existing_assignment_count > 0 THEN
        err_code := 3; -- InvalidOperation
        err_msg := 'role is assigned';
        RETURN;
    END IF;

    _time := (clock_timestamp() AT TIME ZONE 'UTC');
    -- role status: Deleted(5)
    UPDATE public.roles
        SET updated_at = _time, updated_by = _deleted_by, status = 5, status_updated_at = _time, status_updated_by = _deleted_by,
            status_comment = _status_comment, _version_stamp = _version_stamp + 1, _timestamp = _time
        WHERE id = _id;

    UPDATE public.role_info
        SET updated_at = _time, updated_by = _deleted_by, is_deleted = TRUE, deleted_at = _time, deleted_by = _deleted_by, _version_stamp = _version_stamp + 1, _timestamp = _time
    WHERE role_id = _role_id;
END;
$$ LANGUAGE plpgsql;
