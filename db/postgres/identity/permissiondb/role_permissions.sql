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
    _id bigint;
BEGIN
    IF array_length(_permission_ids, 1) IS NULL THEN
        RETURN FALSE;
    END IF;

    FOREACH _id IN ARRAY _permission_ids LOOP
        IF NOT EXISTS (SELECT 1 FROM public.role_permissions WHERE role_id = _role_id AND permission_id = _id AND is_deleted IS FALSE LIMIT 1) THEN
            RETURN FALSE;
        END IF;
    END LOOP;
    RETURN TRUE;
END;
$$ LANGUAGE plpgsql;
