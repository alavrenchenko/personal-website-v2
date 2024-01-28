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

-- FUNCTION: public.notification_exists(uuid)

CREATE OR REPLACE FUNCTION public.notification_exists(
    _id public.notifications.id%TYPE
) RETURNS boolean AS $$
BEGIN
    RETURN EXISTS (SELECT 1 FROM public.notifications WHERE id = _id LIMIT 1);
END;
$$ LANGUAGE plpgsql;
