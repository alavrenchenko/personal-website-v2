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

-- PROCEDURE: public.create_permission(bigint, character varying, bigint, text, text)
/*
Permission statuses:
    Active = 2

Permission group statuses:
    Active = 2

Error codes:
    NoError                 = 0
    InvalidOperation        = 3
    PermissionGroupNotFound = 12200
*/
CREATE OR REPLACE PROCEDURE public.create_permission(
    IN _group_id public.permissions.group_id%TYPE,
    IN _name public.permissions.name%TYPE,
    IN _created_by public.permissions.created_by%TYPE,
    IN _status_comment public.permissions.status_comment%TYPE,
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

    SELECT status INTO _status FROM public.permission_groups WHERE id = _group_id LIMIT 1 FOR UPDATE;
    IF NOT FOUND THEN
        err_code := 12200; -- PermissionGroupNotFound
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
    INSERT INTO public.permissions(group_id, name, created_at, created_by, updated_at, updated_by, status, status_updated_at, status_updated_by,
            status_comment, description, _version_stamp, _timestamp)
        VALUES (_group_id, _name, _time, _created_by, _time, _created_by, 2, _time, _created_by, _status_comment, _description, 1, _time)
        RETURNING id INTO _id;
END;
$$ LANGUAGE plpgsql;