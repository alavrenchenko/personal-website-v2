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

-- FUNCTION: public.user_role_assignment_exists(bigint, bigint)
/*
User role assignment statuses:
    Deleted = 5
*/
CREATE OR REPLACE FUNCTION public.user_role_assignment_exists(
    _user_id public.user_role_assignments.user_id%TYPE,
    _role_id public.user_role_assignments.role_id%TYPE
) RETURNS boolean AS $$
BEGIN
   -- user's role assignment status: Deleted(5)
    RETURN EXISTS (SELECT 1 FROM public.user_role_assignments WHERE user_id = _user_id AND role_id = _role_id AND status <> 5 LIMIT 1);
END;
$$ LANGUAGE plpgsql;

-- FUNCTION: public.is_role_assigned(bigint, bigint)
/*
User role assignment statuses:
    Deleted = 5
*/
CREATE OR REPLACE FUNCTION public.is_role_assigned(
    _user_id public.user_role_assignments.user_id%TYPE,
    _role_id public.user_role_assignments.role_id%TYPE
) RETURNS boolean AS $$
BEGIN
   -- user's role assignment status: Deleted(5)
    RETURN EXISTS (SELECT 1 FROM public.user_role_assignments WHERE user_id = _user_id AND role_id = _role_id AND status <> 5 LIMIT 1);
END;
$$ LANGUAGE plpgsql;

-- PROCEDURE: public.create_user_role_assignment(bigint, bigint, bigint, bigint, text)
/*
User role assignment statuses:
    Active = 2

Error codes:
    NoError                     = 0
    RoleAssignmentAlreadyExists = 13401
*/
-- Minimum transaction isolation level: Read committed.
CREATE OR REPLACE PROCEDURE public.create_user_role_assignment(
    IN _role_assignment_id public.user_role_assignments.role_assignment_id%TYPE,
    IN _user_id public.user_role_assignments.user_id%TYPE,
    IN _role_id public.user_role_assignments.role_id%TYPE,
    IN _created_by public.user_role_assignments.created_by%TYPE,
    IN _status_comment public.user_role_assignments.status_comment%TYPE,
    OUT _id public.user_role_assignments.id%TYPE,
    OUT err_code bigint,
    OUT err_msg text) AS $$
DECLARE
    _time timestamp(6) without time zone;
BEGIN
    _id := 0;
    err_code := 0; -- NoError
    err_msg := '';

    IF user_role_assignment_exists(_user_id, _role_id) THEN
        err_code := 13401; -- RoleAssignmentAlreadyExists
        err_msg := 'role assignment with the same params already exists';
        RETURN;
    END IF;

    _time := (clock_timestamp() AT TIME ZONE 'UTC');
    -- user's role assignment status: Active(2)
    INSERT INTO public.user_role_assignments(role_assignment_id, user_id, role_id, created_at, created_by, updated_at, updated_by, status, status_updated_at,
            status_updated_by, status_comment, _version_stamp, _timestamp)
        VALUES (_role_assignment_id, _user_id, _role_id, _time, _created_by, _time, _created_by, 2, _time, _created_by, _status_comment, 1, _time)
        RETURNING id INTO _id;

    EXCEPTION
        WHEN unique_violation THEN
            IF user_role_assignment_exists(_user_id, _role_id) THEN
                err_code := 13401; -- RoleAssignmentAlreadyExists
                err_msg := 'role assignment with the same params already exists';
                RETURN;
            END IF;
            RAISE;
END;
$$ LANGUAGE plpgsql;

-- PROCEDURE: public.delete_user_role_assignment(bigint, bigint, text)
/*
User role assignment statuses:
    Deleted  = 5

Error codes:
    NoError                = 0
    InvalidOperation       = 3
    RoleAssignmentNotFound = 13400
*/
-- Minimum transaction isolation level: Read committed.
CREATE OR REPLACE PROCEDURE public.delete_user_role_assignment(
    IN _id public.user_role_assignments.id%TYPE,
    IN _deleted_by public.user_role_assignments.updated_by%TYPE,
    IN _status_comment public.user_role_assignments.status_comment%TYPE,
    OUT err_code bigint,
    OUT err_msg text) AS $$
DECLARE
    _time timestamp(6) without time zone;
    _status public.user_role_assignments.status%TYPE;
BEGIN
    err_code := 0; -- NoError
    err_msg := '';

    SELECT status INTO _status FROM public.user_role_assignments WHERE id = _id LIMIT 1 FOR UPDATE;
    IF NOT FOUND THEN
        err_code := 13400; -- RoleAssignmentNotFound
        err_msg := 'user''s role assignment not found';
        RETURN;
    END IF;

    -- user's role assignment status: Deleted(5)
    IF _status = 5 THEN
        err_code := 3; -- InvalidOperation
        err_msg := 'user''s role assignment has already been deleted';
        RETURN;
    END IF;

    _time := (clock_timestamp() AT TIME ZONE 'UTC');
    -- user's role assignment status: Deleted(5)
    UPDATE public.user_role_assignments
        SET updated_at = _time, updated_by = _deleted_by, status = 5, status_updated_at = _time, status_updated_by = _deleted_by,
            status_comment = _status_comment, _version_stamp = _version_stamp + 1, _timestamp = _time
        WHERE id = _id;
END;
$$ LANGUAGE plpgsql;
