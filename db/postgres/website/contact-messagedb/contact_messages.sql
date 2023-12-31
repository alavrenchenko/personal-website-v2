-- Copyright 2024 Alexey Lavrenchenko. All rights reserved.
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

-- PROCEDURE: public.create_contact_message(bigint, text, character varying, character varying, text)
/*
Contact message statuses:
    New = 1

Error codes:
    NoError = 0
*/
-- Minimum transaction isolation level: Read committed.
CREATE OR REPLACE PROCEDURE public.create_contact_message(
    IN _created_by public.contact_messages.created_by%TYPE,
    IN _status_comment public.contact_messages.status_comment%TYPE,
    IN _name public.contact_messages.name%TYPE,
    IN _email public.contact_messages.email%TYPE,
    IN _message public.contact_messages.message%TYPE,
    OUT _id public.contact_messages.id%TYPE,
    OUT err_code bigint,
    OUT err_msg text) AS $$
DECLARE
    _time timestamp(6) without time zone;
BEGIN
    _id := 0;
    err_code := 0; -- NoError
    err_msg := '';

    _time := (clock_timestamp() AT TIME ZONE 'UTC');
    -- contact message status: New(1)
    INSERT INTO public.contact_messages(created_at, created_by, updated_at, updated_by, status, status_updated_at, status_updated_by,
            status_comment, name, email, message, _version_stamp, _timestamp)
        VALUES (_time, _created_by, _time, _created_by, 1, _time, _created_by, _status_comment, _name, _email, _message, 1, _time)
        RETURNING id INTO _id;
END;
$$ LANGUAGE plpgsql;
