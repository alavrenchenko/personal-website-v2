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

-- FUNCTION: public.username_exists(character varying)
/*
User statuses:
    Deleted = 8
*/
CREATE OR REPLACE FUNCTION public.username_exists(
    _name public.users.name%TYPE
) RETURNS boolean AS $$
BEGIN
   -- user status: Deleted(8)
    RETURN EXISTS (SELECT 1 FROM public.users WHERE lower(name) = lower(_name) AND status <> 8 LIMIT 1);
END;
$$ LANGUAGE plpgsql;

-- FUNCTION: public.user_email_exists(character varying)
/*
User statuses:
    Deleted = 8
*/
CREATE OR REPLACE FUNCTION public.user_email_exists(
    _email public.users.email%TYPE
) RETURNS boolean AS $$
BEGIN
   -- user status: Deleted(8)
    RETURN EXISTS (SELECT 1 FROM public.users WHERE lower(email) = lower(_email) AND status <> 8 LIMIT 1);
END;
$$ LANGUAGE plpgsql;

-- PROCEDURE: public.create_user(smallint, bigint, bigint, smallint, text, character varying, character varying, character varying, character varying, timestamp without time zone, smallint)
/*
User statuses:
    Active = 3

Error codes:
    NoError                = 0
    UserEmailAlreadyExists = 11002
*/
-- Minimum transaction isolation level: Read committed.
CREATE OR REPLACE PROCEDURE public.create_user(
    IN _type public.users.type%TYPE,
    IN _group public.users.group%TYPE,
    IN _created_by public.users.created_by%TYPE,
    IN _status public.users.status%TYPE,
    IN _status_comment public.users.status_comment%TYPE,
    IN _email public.users.email%TYPE,
    IN _first_name public.personal_info.first_name%TYPE,
    IN _last_name public.personal_info.last_name%TYPE,
    IN _display_name public.personal_info.display_name%TYPE,
    IN _birth_date public.personal_info.birth_date%TYPE,
    IN _gender public.personal_info.gender%TYPE,
    OUT _id public.users.id%TYPE,
    OUT err_code bigint,
    OUT err_msg text) AS $$
DECLARE
    _time timestamp(6) without time zone;
BEGIN
    _id := 0;
    err_code := 0; -- NoError
    err_msg := '';

    IF _email IS NOT NULL AND public.user_email_exists(_email) THEN
        err_code := 11002; -- UserEmailAlreadyExists
        err_msg := 'user with the same email already exists';
        RETURN;
    END IF;

    _time := (clock_timestamp() AT TIME ZONE 'UTC');
    -- user status: Active(3)
    INSERT INTO public.users(type, "group", created_at, created_by, updated_at, updated_by, status, status_updated_at, status_updated_by,
            status_comment, email, _version_stamp, _timestamp)
        VALUES (_type, _group, _time, _created_by, _time, _created_by, 3, _time, _created_by, _status_comment, _email, 1, _time)
        RETURNING id INTO _id;

    INSERT INTO public.personal_info(user_id, created_at, created_by, updated_at, updated_by, first_name, last_name, display_name,
            birth_date, gender, _version_stamp, _timestamp)
        VALUES (_id, _time, _created_by, _time, _created_by, _first_name, _last_name, _display_name, _birth_date, _gender, 1, _time);

    EXCEPTION
        WHEN unique_violation THEN
            IF _id = 0 AND _email IS NOT NULL AND public.user_email_exists(_email) THEN
                err_code := 11002; -- UserEmailAlreadyExists
                err_msg := 'user with the same email already exists';
                RETURN;
            END IF;
            RAISE;
END;
$$ LANGUAGE plpgsql;
