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

-- FUNCTION: public.is_permission_granted(bigint, bigint)

CREATE OR REPLACE FUNCTION public.is_permission_granted(
    _role_id public.role_permissions.role_id%TYPE,
    _permission_id public.role_permissions.permission_id%TYPE
) RETURNS boolean AS $$
BEGIN
    RETURN EXISTS (SELECT 1 FROM public.role_permissions WHERE role_id = _role_id AND permission_id = _permission_id AND is_deleted IS FALSE LIMIT 1);
END;
$$ LANGUAGE plpgsql;

-- FUNCTION: public.are_permissions_granted(bigint, bigint[])

CREATE OR REPLACE FUNCTION public.are_permissions_granted(
    _role_id public.role_permissions.role_id%TYPE,
    _permission_ids bigint[]
) RETURNS boolean AS $$
DECLARE
    _permission_id bigint;
BEGIN
    IF array_length(_permission_ids, 1) IS NULL THEN
        RETURN FALSE;
    END IF;

    FOREACH _permission_id IN ARRAY _permission_ids LOOP
        IF NOT EXISTS (SELECT 1 FROM public.role_permissions WHERE role_id = _role_id AND permission_id = _permission_id AND is_deleted IS FALSE LIMIT 1) THEN
            RETURN FALSE;
        END IF;
    END LOOP;
    RETURN TRUE;
END;
$$ LANGUAGE plpgsql;

-- PROCEDURE: public.grant_permissions(bigint, bigint[], bigint)
/*
Permission statuses:
    Active = 2

Error codes:
    NoError                  = 0
    InvalidOperation         = 3
    InvalidData              = 1000
    PermissionNotFound       = 11800
    PermissionAlreadyGranted = 11902
*/
-- Minimum transaction isolation level: Serializable.
CREATE OR REPLACE PROCEDURE public.grant_permissions(
    IN _role_id public.role_permissions.role_id%TYPE,
    IN _permission_ids bigint[],
    IN _operation_user_id public.role_permissions.created_by%TYPE,
    OUT err_code bigint,
    OUT err_msg text) AS $$
DECLARE
    _time timestamp(6) without time zone;
    _permission_id bigint;
    _status public.permissions.status%TYPE;
BEGIN
    err_code := 0; -- NoError
    err_msg := '';

    IF array_length(_permission_ids, 1) IS NULL THEN
        err_code := 1000; -- InvalidData
        err_msg := 'number of permission ids is 0';
        RETURN;
    END IF;

    _time := (clock_timestamp() AT TIME ZONE 'UTC');

    FOREACH _permission_id IN ARRAY _permission_ids LOOP
        SELECT status INTO _status FROM public.permissions WHERE id = _permission_id LIMIT 1;
        IF NOT FOUND THEN
            err_code := 11800; -- PermissionNotFound
            err_msg := format('permission (%s) not found', _permission_id);
            RETURN;
        END IF;

        -- permission status: Active(2)
        IF _status <> 2 THEN
            err_code := 3; -- InvalidOperation
            err_msg := format('invalid permission (%s) status (%s)', _permission_id, _status);
            RETURN;
        END IF;

        IF public.is_permission_granted(_role_id, _permission_id) THEN
            err_code := 11902; -- PermissionAlreadyGranted
            err_msg := format('permission (%s) already granted to the role', _permission_id);
            RETURN;
        END IF;

        INSERT INTO public.role_permissions(role_id, permission_id, created_at, created_by, _version_stamp, _timestamp)
            VALUES (_role_id, _permission_id, _time, _operation_user_id, 1, _time);
    END LOOP;

    EXCEPTION
        WHEN unique_violation THEN
            IF _permission_id IS NOT NULL AND public.is_permission_granted(_role_id, _permission_id) THEN
                err_code := 11902; -- PermissionAlreadyGranted
                err_msg := format('permission (%s) already granted to the role', _permission_id);
                RETURN;
            END IF;
            RAISE;
END;
$$ LANGUAGE plpgsql;

-- PROCEDURE: public.revoke_permissions(bigint, bigint[], bigint)
/*
Error codes:
    NoError              = 0
    InternalError        = 2
    InvalidData          = 1000
    PermissionNotGranted = 11903
*/
-- Minimum transaction isolation level: Read committed.
CREATE OR REPLACE PROCEDURE public.revoke_permissions(
    IN _role_id public.role_permissions.role_id%TYPE,
    IN _permission_ids bigint[],
    IN _operation_user_id public.role_permissions.deleted_by%TYPE,
    OUT err_code bigint,
    OUT err_msg text) AS $$
DECLARE
    _time timestamp(6) without time zone;
    _permission_id bigint;
    _id public.role_permissions.id%TYPE;
BEGIN
    err_code := 0; -- NoError
    err_msg := '';

    IF array_length(_permission_ids, 1) IS NULL THEN
        err_code := 1000; -- InvalidData
        err_msg := 'number of permission ids is 0';
        RETURN;
    END IF;

    _time := (clock_timestamp() AT TIME ZONE 'UTC');

    -- _permission_ids need to be sorted to avoid deadlocks
    FOR _permission_id IN SELECT DISTINCT unnest(_permission_ids) ORDER BY 1 LOOP
        SELECT id INTO _id FROM public.role_permissions WHERE role_id = _role_id AND permission_id = _permission_id AND is_deleted IS FALSE LIMIT 1 FOR UPDATE;
        IF NOT FOUND THEN
            err_code := 11903; -- PermissionNotGranted
            err_msg := format('permission (%s) not granted to the role', _permission_id);
            RETURN;
        END IF;

        UPDATE public.role_permissions
            SET is_deleted = TRUE, deleted_at = _time, deleted_by = _operation_user_id, _version_stamp = _version_stamp + 1, _timestamp = _time
            WHERE id = _id;

        INSERT INTO public.deleted_role_permissions SELECT * FROM public.role_permissions WHERE id = _id LIMIT 1;
        -- an insertion check
        IF NOT FOUND THEN
            err_code := 2; -- InternalError
            err_msg := 'role permission info wasn''t inserted into deleted_role_permissions';
            RETURN;
        END IF;
    END LOOP;
END;
$$ LANGUAGE plpgsql;

-- PROCEDURE: public.revoke_all_permissions(bigint, bigint)
/*
Error codes:
    NoError       = 0
    InternalError = 2
*/
-- Minimum transaction isolation level: Read committed.
CREATE OR REPLACE PROCEDURE public.revoke_all_permissions(
    IN _role_id public.role_permissions.role_id%TYPE,
    IN _operation_user_id public.role_permissions.deleted_by%TYPE,
    OUT err_code bigint,
    OUT err_msg text) AS $$
DECLARE
    _time timestamp(6) without time zone;
    -- rows need to be sorted to avoid deadlocks
    _rpc CURSOR FOR
        SELECT id FROM public.role_permissions
            WHERE role_id = _role_id AND is_deleted IS FALSE
            ORDER BY permission_id
            FOR UPDATE;
BEGIN
    err_code := 0; -- NoError
    err_msg := '';
    _time := (clock_timestamp() AT TIME ZONE 'UTC');

    FOR _item IN _rpc LOOP
        UPDATE public.role_permissions
            SET is_deleted = TRUE, deleted_at = _time, deleted_by = _operation_user_id, _version_stamp = _version_stamp + 1, _timestamp = _time
            WHERE CURRENT OF _rpc;

        INSERT INTO public.deleted_role_permissions SELECT * FROM public.role_permissions WHERE id = _item.id LIMIT 1;
        -- an insertion check
        IF NOT FOUND THEN
            err_code := 2; -- InternalError
            err_msg := 'role permission info wasn''t inserted into deleted_role_permissions';
            RETURN;
        END IF;
    END LOOP;
END;
$$ LANGUAGE plpgsql;
