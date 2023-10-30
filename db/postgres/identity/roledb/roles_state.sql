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

-- PROCEDURE: public.start_assigning_role(uuid, bigint, bigint)
/*
Role statuses:
    Active = 2

Error codes:
    NoError          = 0
    InvalidOperation = 3
    RoleNotFound     = 11600
*/
-- Minimum transaction isolation level: Read committed.
CREATE OR REPLACE PROCEDURE public.start_assigning_role(
    IN _operation_id public.new_role_assignments.operation_id%TYPE,
    IN _role_id public.new_role_assignments.role_id%TYPE,
    IN _operation_user_id public.new_role_assignments.created_by%TYPE,
    OUT err_code bigint,
    OUT err_msg text) AS $$
DECLARE
    _time timestamp(6) without time zone;
    _status public.roles.status%TYPE;
BEGIN
    err_code := 0; -- NoError
    err_msg := '';

    SELECT status INTO _status FROM public.roles WHERE id = _role_id LIMIT 1 FOR UPDATE;
    IF NOT FOUND THEN
        err_code := 11600; -- RoleNotFound
        err_msg := 'role not found';
        RETURN;
    END IF;

    -- role status: Active(2)
    IF _status <> 2 THEN
        err_code := 3; -- InvalidOperation
        err_msg := format('invalid role status (%s)', _status);
        RETURN;
    END IF;

    _time := (clock_timestamp() AT TIME ZONE 'UTC');
    INSERT INTO public.new_role_assignments(operation_id, role_id, created_at, created_by, _version_stamp, _timestamp)
        VALUES (_operation_id, _role_id, _time, _operation_user_id, 1, _time);
END;
$$ LANGUAGE plpgsql;

-- PROCEDURE: public.finish_assigning_role(uuid, boolean, bigint)
/*
Error codes:
    NoError                = 0
    InternalError          = 2
    InvalidOperation       = 3
    RoleInfoNotFound       = 11602
    RoleAssignmentNotFound = 13400
*/
-- Minimum transaction isolation level: Read committed.
CREATE OR REPLACE PROCEDURE public.finish_assigning_role(
    IN _operation_id public.new_role_assignments.operation_id%TYPE,
    IN _succeeded boolean,
    IN _operation_user_id public.new_role_assignments.created_by%TYPE,
    OUT err_code bigint,
    OUT err_msg text) AS $$
DECLARE
    _time timestamp(6) without time zone;
    _role_id public.new_role_assignments.role_id%TYPE;
    _is_deleted public.role_info.is_deleted%TYPE;
BEGIN
    err_code := 0; -- NoError
    err_msg := '';

    SELECT role_id INTO _role_id FROM public.new_role_assignments WHERE operation_id = _operation_id LIMIT 1 FOR UPDATE;
    IF NOT FOUND THEN
        err_code := 13400; -- RoleAssignmentNotFound
        err_msg := 'new role assignment not found';
        RETURN;
    END IF;

    IF _succeeded THEN
        SELECT is_deleted INTO _is_deleted FROM public.role_info WHERE role_id = _role_id LIMIT 1 FOR UPDATE;
        IF NOT FOUND THEN
            err_code := 11602; -- RoleInfoNotFound
            err_msg := 'role info not found';
            RETURN;
        END IF;

        IF _is_deleted THEN
            err_code := 3; -- InvalidOperation
            err_msg := 'role info is deleted';
            RETURN;
        END IF;

        _time := (clock_timestamp() AT TIME ZONE 'UTC');
        UPDATE public.role_info
            SET updated_at = _time, updated_by = _operation_user_id, active_assignment_count = active_assignment_count + 1, existing_assignment_count = existing_assignment_count + 1,
                _version_stamp = _version_stamp + 1, _timestamp = _time
            WHERE role_id = _role_id;
        -- an update check
        IF NOT FOUND THEN
            err_code := 2; -- InternalError
            err_msg := 'role info wasn''t updated';
            RETURN;
        END IF;
    END IF;

    DELETE FROM public.new_role_assignments WHERE operation_id = _operation_id;
    -- a deletion check
    IF NOT FOUND THEN
        err_code := 2; -- InternalError
        err_msg := 'new role assignment wasn''t deleted';
    END IF;
END;
$$ LANGUAGE plpgsql;


-- PROCEDURE: public.decr_role_assignments(bigint, bigint)
/*
Error codes:
    NoError          = 0
    InternalError    = 2
    InvalidOperation = 3
    RoleInfoNotFound = 11602
*/
-- Minimum transaction isolation level: Read committed.
CREATE OR REPLACE PROCEDURE public.decr_role_assignments(
    IN _role_id public.role_info.role_id%TYPE,
    IN _operation_user_id public.role_info.updated_by%TYPE,
    OUT err_code bigint,
    OUT err_msg text) AS $$
DECLARE
    _time timestamp(6) without time zone;
    _is_deleted public.role_info.is_deleted%TYPE;
BEGIN
    err_code := 0; -- NoError
    err_msg := '';
    
    SELECT is_deleted INTO _is_deleted FROM public.role_info WHERE role_id = _role_id LIMIT 1 FOR UPDATE;
    IF NOT FOUND THEN
        err_code := 11602; -- RoleInfoNotFound
        err_msg := 'role info not found';
        RETURN;
    END IF;

    IF _is_deleted THEN
        err_code := 3; -- InvalidOperation
        err_msg := 'role info is deleted';
        RETURN;
    END IF;

    _time := (clock_timestamp() AT TIME ZONE 'UTC');
    UPDATE public.role_info
        SET updated_at = _time, updated_by = _operation_user_id, active_assignment_count = active_assignment_count - 1, existing_assignment_count = existing_assignment_count - 1,
            _version_stamp = _version_stamp + 1, _timestamp = _time
        WHERE role_id = _role_id;
    -- an update check
    IF NOT FOUND THEN
        err_code := 2; -- InternalError
        err_msg := 'role info wasn''t updated';
    END IF;
END;
$$ LANGUAGE plpgsql;

-- PROCEDURE: public.incr_active_role_assignments(bigint, bigint)
/*
Error codes:
    NoError          = 0
    InternalError    = 2
    InvalidOperation = 3
    RoleInfoNotFound = 11602
*/
-- Minimum transaction isolation level: Read committed.
CREATE OR REPLACE PROCEDURE public.incr_active_role_assignments(
    IN _role_id public.role_info.role_id%TYPE,
    IN _operation_user_id public.role_info.updated_by%TYPE,
    OUT err_code bigint,
    OUT err_msg text) AS $$
DECLARE
    _time timestamp(6) without time zone;
    _is_deleted public.role_info.is_deleted%TYPE;
BEGIN
    err_code := 0; -- NoError
    err_msg := '';
    
    SELECT is_deleted INTO _is_deleted FROM public.role_info WHERE role_id = _role_id LIMIT 1 FOR UPDATE;
    IF NOT FOUND THEN
        err_code := 11602; -- RoleInfoNotFound
        err_msg := 'role info not found';
        RETURN;
    END IF;

    IF _is_deleted THEN
        err_code := 3; -- InvalidOperation
        err_msg := 'role info is deleted';
        RETURN;
    END IF;

    _time := (clock_timestamp() AT TIME ZONE 'UTC');
    UPDATE public.role_info
        SET updated_at = _time, updated_by = _operation_user_id, active_assignment_count = active_assignment_count + 1,
            _version_stamp = _version_stamp + 1, _timestamp = _time
        WHERE role_id = _role_id;
    -- an update check
    IF NOT FOUND THEN
        err_code := 2; -- InternalError
        err_msg := 'role info wasn''t updated';
    END IF;
END;
$$ LANGUAGE plpgsql;

-- PROCEDURE: public.decr_active_role_assignments(bigint, bigint)
/*
Error codes:
    NoError          = 0
    InternalError    = 2
    InvalidOperation = 3
    RoleInfoNotFound = 11602
*/
-- Minimum transaction isolation level: Read committed.
CREATE OR REPLACE PROCEDURE public.decr_active_role_assignments(
    IN _role_id public.role_info.role_id%TYPE,
    IN _operation_user_id public.role_info.updated_by%TYPE,
    OUT err_code bigint,
    OUT err_msg text) AS $$
DECLARE
    _time timestamp(6) without time zone;
    _is_deleted public.role_info.is_deleted%TYPE;
BEGIN
    err_code := 0; -- NoError
    err_msg := '';
    
    SELECT is_deleted INTO _is_deleted FROM public.role_info WHERE role_id = _role_id LIMIT 1 FOR UPDATE;
    IF NOT FOUND THEN
        err_code := 11602; -- RoleInfoNotFound
        err_msg := 'role info not found';
        RETURN;
    END IF;

    IF _is_deleted THEN
        err_code := 3; -- InvalidOperation
        err_msg := 'role info is deleted';
        RETURN;
    END IF;

    _time := (clock_timestamp() AT TIME ZONE 'UTC');
    UPDATE public.role_info
        SET updated_at = _time, updated_by = _operation_user_id, active_assignment_count = active_assignment_count - 1,
            _version_stamp = _version_stamp + 1, _timestamp = _time
        WHERE role_id = _role_id;
    -- an update check
    IF NOT FOUND THEN
        err_code := 2; -- InternalError
        err_msg := 'role info wasn''t updated';
    END IF;
END;
$$ LANGUAGE plpgsql;

-- PROCEDURE: public.decr_existing_role_assignments(bigint, bigint)
/*
Error codes:
    NoError          = 0
    InternalError    = 2
    InvalidOperation = 3
    RoleInfoNotFound = 11602
*/
-- Minimum transaction isolation level: Read committed.
CREATE OR REPLACE PROCEDURE public.decr_existing_role_assignments(
    IN _role_id public.role_info.role_id%TYPE,
    IN _operation_user_id public.role_info.updated_by%TYPE,
    OUT err_code bigint,
    OUT err_msg text) AS $$
DECLARE
    _time timestamp(6) without time zone;
    _is_deleted public.role_info.is_deleted%TYPE;
BEGIN
    err_code := 0; -- NoError
    err_msg := '';
    
    SELECT is_deleted INTO _is_deleted FROM public.role_info WHERE role_id = _role_id LIMIT 1 FOR UPDATE;
    IF NOT FOUND THEN
        err_code := 11602; -- RoleInfoNotFound
        err_msg := 'role info not found';
        RETURN;
    END IF;

    IF _is_deleted THEN
        err_code := 3; -- InvalidOperation
        err_msg := 'role info is deleted';
        RETURN;
    END IF;

    _time := (clock_timestamp() AT TIME ZONE 'UTC');
    UPDATE public.role_info
        SET updated_at = _time, updated_by = _operation_user_id, existing_assignment_count = existing_assignment_count - 1,
            _version_stamp = _version_stamp + 1, _timestamp = _time
        WHERE role_id = _role_id;
    -- an update check
    IF NOT FOUND THEN
        err_code := 2; -- InternalError
        err_msg := 'role info wasn''t updated';
    END IF;
END;
$$ LANGUAGE plpgsql;
