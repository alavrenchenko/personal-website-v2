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

-- FUNCTION: public.is_role_assigned(bigint, bigint, smallint)
/*
Role assignment statuses:
    Deleted = 5
*/
CREATE OR REPLACE FUNCTION public.is_role_assigned(
    _role_id public.role_assignments.role_id%TYPE,
    _assigned_to public.role_assignments.assigned_to%TYPE,
    _assignee_type public.role_assignments.assignee_type%TYPE
) RETURNS boolean AS $$
BEGIN
   -- role assignment status: Deleted(5)
    RETURN EXISTS (SELECT 1 FROM public.role_assignments WHERE role_id = _role_id AND assigned_to = _assigned_to AND assignee_type = _assignee_type AND status <> 5 LIMIT 1);
END;
$$ LANGUAGE plpgsql;

-- PROCEDURE: public.create_role_assignment(bigint, bigint, smallint, bigint, text, text)
/*
Role assignment statuses:
    Active = 2

Error codes:
    NoError             = 0
    RoleAlreadyAssigned = 13401
*/
CREATE OR REPLACE PROCEDURE public.create_role_assignment(
    IN _role_id public.role_assignments.role_id%TYPE,
    IN _assigned_to public.role_assignments.assigned_to%TYPE,
    IN _assignee_type public.role_assignments.assignee_type%TYPE,
    IN _created_by public.role_assignments.created_by%TYPE,
    IN _status_comment public.role_assignments.status_comment%TYPE,
    IN _description public.role_assignments.description%TYPE,
    OUT _id public.role_assignments.id%TYPE,
    OUT err_code bigint,
    OUT err_msg text) AS $$
DECLARE
    _time timestamp(6) without time zone;
BEGIN
    _id := 0;
    err_code := 0; -- NoError
    err_msg := '';

    IF is_role_assigned(_role_id, _assigned_to, _assignee_type) THEN
        err_code := 13401; -- RoleAlreadyAssigned
        err_msg := 'role already assigned';
        RETURN;
    END IF;

    _time := (clock_timestamp() AT TIME ZONE 'UTC');
    -- role assignment status: Active(2)
    INSERT INTO public.role_assignments(role_id, assigned_to, assignee_type, created_at, created_by, updated_at, updated_by, status, status_updated_at,
            status_updated_by, status_comment, description, _version_stamp, _timestamp)
        VALUES (_role_id, _assigned_to, _assignee_type, _time, _created_by, _time, _created_by, 2, _time, _created_by, _status_comment, _description, 1, _time)
        RETURNING id INTO _id;

    EXCEPTION
        WHEN unique_violation THEN
            IF is_role_assigned(_role_id, _assigned_to, _assignee_type) THEN
                err_code := 13401; -- RoleAlreadyAssigned
                err_msg := 'role already assigned';
                RETURN;
            END IF;
            RAISE;
END;
$$ LANGUAGE plpgsql;

-- PROCEDURE: public.start_deleting_role_assignment(bigint, bigint, text)
/*
Role assignment statuses:
    Deleting = 4
    Deleted  = 5

Error codes:
    NoError                = 0
    InvalidOperation       = 3
    RoleAssignmentNotFound = 13400
*/
CREATE OR REPLACE PROCEDURE public.start_deleting_role_assignment(
    IN _id public.role_assignments.id%TYPE,
    IN _deleted_by public.role_assignments.updated_by%TYPE,
    IN _status_comment public.role_assignments.status_comment%TYPE,
    OUT err_code bigint,
    OUT err_msg text) AS $$
DECLARE
    _time timestamp(6) without time zone;
    _status public.role_assignments.status%TYPE;
BEGIN
    err_code := 0; -- NoError
    err_msg := '';

    SELECT status INTO _status FROM public.role_assignments WHERE id = _id LIMIT 1 FOR UPDATE;
    IF NOT FOUND THEN
        err_code := 13400; -- RoleAssignmentNotFound
        err_msg := 'role assignment not found';
        RETURN;
    END IF;

    -- role assignment statuses: Deleting(4), Deleted(5)
    IF _status = 4 OR _status = 5 THEN
        err_code := 3; -- InvalidOperation
        err_msg := format('invalid role assignment status (%s)', _status);
        RETURN;
    END IF;

    _time := (clock_timestamp() AT TIME ZONE 'UTC');
    -- role assignment status: Deleting(4)
    UPDATE public.role_assignments
        SET updated_at = _time, updated_by = _deleted_by, status = 4, status_updated_at = _time, status_updated_by = _deleted_by,
            status_comment = _status_comment, _version_stamp = _version_stamp + 1, _timestamp = _time
        WHERE id = _id;
END;
$$ LANGUAGE plpgsql;

-- PROCEDURE: public.delete_role_assignment(bigint, bigint, text)
/*
Role assignment statuses:
    Deleted = 5

Error codes:
    NoError                = 0
    InternalError          = 2
    InvalidOperation       = 3
    RoleAssignmentNotFound = 13400
*/
CREATE OR REPLACE PROCEDURE public.delete_role_assignment(
    IN _id public.role_assignments.id%TYPE,
    IN _deleted_by public.role_assignments.updated_by%TYPE,
    IN _status_comment public.role_assignments.status_comment%TYPE,
    OUT err_code bigint,
    OUT err_msg text) AS $$
DECLARE
    _time timestamp(6) without time zone;
    _status public.role_assignments.status%TYPE;
BEGIN
    err_code := 0; -- NoError
    err_msg := '';

    SELECT status INTO _status FROM public.role_assignments WHERE id = _id LIMIT 1 FOR UPDATE;
    IF NOT FOUND THEN
        err_code := 13400; -- RoleAssignmentNotFound
        err_msg := 'role assignment not found';
        RETURN;
    END IF;

    -- role assignment status: Deleting(4)
    IF _status <> 4 THEN
        err_code := 3; -- InvalidOperation
        -- role assignment status: Deleted(5)
        IF _status = 5 THEN
            err_msg := 'role assignment has already been deleted';
        ELSE
            err_msg := format('invalid role assignment status (%s)', _status);
        END IF;
        RETURN;
    END IF;

    _time := (clock_timestamp() AT TIME ZONE 'UTC');
    -- role assignment status: Deleted(5)
    UPDATE public.role_assignments
        SET updated_at = _time, updated_by = _deleted_by, status = 5, status_updated_at = _time, status_updated_by = _deleted_by,
            status_comment = _status_comment, _version_stamp = _version_stamp + 1, _timestamp = _time
        WHERE id = _id;

    INSERT INTO public.deleted_role_assignments SELECT * FROM public.role_assignments WHERE id = _id LIMIT 1;
    -- an insertion check
    IF NOT FOUND THEN
        err_code := 2; -- InternalError
        err_msg := 'role assignment wasn''t inserted into deleted_role_assignments';
    END IF;
END;
$$ LANGUAGE plpgsql;
