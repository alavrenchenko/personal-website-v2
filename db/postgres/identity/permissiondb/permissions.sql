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

-- FUNCTION: public.permission_exists(character varying)
/*
Permission statuses:
    Deleted = 5
*/
CREATE OR REPLACE FUNCTION public.permission_exists(
    _name public.permissions.name%TYPE
) RETURNS boolean AS $$
BEGIN
   -- permission status: Deleted(5)
    RETURN EXISTS (SELECT 1 FROM public.permissions WHERE lower(name) = lower(_name) AND status <> 5 LIMIT 1);
END;
$$ LANGUAGE plpgsql;

-- PROCEDURE: public.create_permission(character varying, bigint, bigint, text, bigint, bigint, text)
/*
Permission statuses:
    Active = 2

Permission group statuses:
    Active = 2

Error codes:
    NoError                 = 0
    InvalidOperation        = 3
    PermissionAlreadyExists = 11801
    PermissionGroupNotFound = 12000
*/
-- Minimum transaction isolation level: Read committed.
CREATE OR REPLACE PROCEDURE public.create_permission(
    IN _name public.permissions.name%TYPE,
    IN _group_id public.permissions.group_id%TYPE,
    IN _created_by public.permissions.created_by%TYPE,
    IN _status_comment public.permissions.status_comment%TYPE,
    IN _app_id public.permissions.app_id%TYPE,
    IN _app_group_id public.permissions.app_group_id%TYPE,
    IN _description public.permissions.description%TYPE,
    OUT _id public.permissions.id%TYPE,
    OUT err_code bigint,
    OUT err_msg text) AS $$
DECLARE
    _time timestamp(6) without time zone;
    _status public.permission_groups.status%TYPE;
BEGIN
    _id := 0;
    err_code := 0; -- NoError
    err_msg := '';

    IF public.permission_exists(_name) THEN
        err_code := 11801; -- PermissionAlreadyExists
        err_msg := 'permission with the same name already exists';
        RETURN;
    END IF;

    SELECT status INTO _status FROM public.permission_groups WHERE id = _group_id LIMIT 1 FOR UPDATE;
    IF NOT FOUND THEN
        err_code := 12000; -- PermissionGroupNotFound
        err_msg := 'permission group not found';
        RETURN;
    END IF;

    -- permission group status: Active(2)
    IF _status <> 2 THEN
        err_code := 3; -- InvalidOperation
        err_msg := format('invalid permission group status (%s)', _status);
        RETURN;
    END IF;

    _time := (clock_timestamp() AT TIME ZONE 'UTC');
    -- permission status: Active(2)
    INSERT INTO public.permissions(name, group_id, created_at, created_by, updated_at, updated_by, status, status_updated_at, status_updated_by,
            status_comment, app_id, app_group_id, description, _version_stamp, _timestamp)
        VALUES (_name, _group_id, _time, _created_by, _time, _created_by, 2, _time, _created_by, _status_comment, _app_id, _app_group_id, _description, 1, _time)
        RETURNING id INTO _id;

    EXCEPTION
        WHEN unique_violation THEN
            IF public.permission_exists(_name) THEN
                err_code := 11801; -- PermissionAlreadyExists
                err_msg := 'permission with the same name already exists';
                RETURN;
            END IF;
            RAISE;
END;
$$ LANGUAGE plpgsql;

-- PROCEDURE: public.start_deleting_permission(bigint, bigint, text)
/*
Permission statuses:
    Deleting = 4
    Deleted  = 5

Error codes:
    NoError            = 0
    InvalidOperation   = 3
    PermissionNotFound = 11800
*/
-- Minimum transaction isolation level: Read committed.
CREATE OR REPLACE PROCEDURE public.start_deleting_permission(
    IN _id public.permissions.id%TYPE,
    IN _deleted_by public.permissions.updated_by%TYPE,
    IN _status_comment public.permissions.status_comment%TYPE,
    OUT err_code bigint,
    OUT err_msg text) AS $$
DECLARE
    _time timestamp(6) without time zone;
    _status public.permissions.status%TYPE;
BEGIN
    err_code := 0; -- NoError
    err_msg := '';

    SELECT status INTO _status FROM public.permissions WHERE id = _id LIMIT 1 FOR UPDATE;
    IF NOT FOUND THEN
        err_code := 11800; -- PermissionNotFound
        err_msg := 'permission not found';
        RETURN;
    END IF;

    -- permission statuses: Deleting(4), Deleted(5)
    IF _status = 4 OR _status = 5 THEN
        err_code := 3; -- InvalidOperation
        err_msg := format('invalid permission status (%s)', _status);
        RETURN;
    END IF;

    _time := (clock_timestamp() AT TIME ZONE 'UTC');
    -- permission status: Deleting(4)
    UPDATE public.permissions
        SET updated_at = _time, updated_by = _deleted_by, status = 4, status_updated_at = _time, status_updated_by = _deleted_by,
            status_comment = _status_comment, _version_stamp = _version_stamp + 1, _timestamp = _time
        WHERE id = _id;
END;
$$ LANGUAGE plpgsql;

-- PROCEDURE: public.delete_permission(bigint, bigint, text)
/*
Permission statuses:
    Deleting = 4
    Deleted  = 5

Error codes:
    NoError            = 0
    InvalidOperation   = 3
    PermissionNotFound = 11800
*/
-- Minimum transaction isolation level: Serializable.
CREATE OR REPLACE PROCEDURE public.delete_permission(
    IN _id public.permissions.id%TYPE,
    IN _deleted_by public.permissions.updated_by%TYPE,
    IN _status_comment public.permissions.status_comment%TYPE,
    OUT err_code bigint,
    OUT err_msg text) AS $$
DECLARE
    _time timestamp(6) without time zone;
    _status public.permissions.status%TYPE;
BEGIN
    err_code := 0; -- NoError
    err_msg := '';

    SELECT status INTO _status FROM public.permissions WHERE id = _id LIMIT 1 FOR UPDATE;
    IF NOT FOUND THEN
        err_code := 11800; -- PermissionNotFound
        err_msg := 'permission not found';
        RETURN;
    END IF;

    -- permission status: Deleting(4)
    IF _status <> 4 THEN
        err_code := 3; -- InvalidOperation
        -- permission status: Deleted(5)
        IF _status = 5 THEN
            err_msg := 'permission has already been deleted';
        ELSE
            err_msg := format('invalid permission status (%s)', _status);
        END IF;
        RETURN;
    END IF;

    IF EXISTS (SELECT 1 FROM public.role_permissions WHERE permission_id = _id AND is_deleted IS FALSE LIMIT 1) THEN
        err_code := 3; -- InvalidOperation
        err_msg := 'permission is granted';
        RETURN;
    END IF;

    _time := (clock_timestamp() AT TIME ZONE 'UTC');
    -- permission status: Deleted(5)
    UPDATE public.permissions
        SET updated_at = _time, updated_by = _deleted_by, status = 5, status_updated_at = _time, status_updated_by = _deleted_by,
            status_comment = _status_comment, _version_stamp = _version_stamp + 1, _timestamp = _time
        WHERE id = _id;
END;
$$ LANGUAGE plpgsql;
