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

INSERT INTO public.app_groups(name, type, created_at, created_by, updated_at, updated_by, status, status_updated_at, status_updated_by, status_comment, version, description, _version_stamp, _timestamp)
    VALUES ('App Group 1', 1, (now() AT TIME ZONE 'UTC'), 1, (now() AT TIME ZONE 'UTC'), 1, 1, (now() AT TIME ZONE 'UTC'), 1, NULL, '1.0.0', 'App Group 1', 1, (now() AT TIME ZONE 'UTC')),
    ('App Group 2', 1, (now() AT TIME ZONE 'UTC'), 1, (now() AT TIME ZONE 'UTC'), 1, 2, (now() AT TIME ZONE 'UTC'), 1, NULL, '1.0.0', 'App Group 2', 1, (now() AT TIME ZONE 'UTC')),
    ('App Group 3', 1, (now() AT TIME ZONE 'UTC'), 1, (now() AT TIME ZONE 'UTC'), 1, 3, (now() AT TIME ZONE 'UTC'), 1, 'test', '1.0.0', 'App Group 3', 1, (now() AT TIME ZONE 'UTC'));

INSERT INTO public.apps(group_id, name, type, category, created_at, created_by, updated_at, updated_by, status, status_updated_at, status_updated_by, status_comment, version, description, _version_stamp, _timestamp)
	VALUES (2, 'App 1', 1, 1, (now() AT TIME ZONE 'UTC'), 1, (now() AT TIME ZONE 'UTC'), 1, 1, (now() AT TIME ZONE 'UTC'), 1, NULL, '1.0.0', 'App 1', 1, (now() AT TIME ZONE 'UTC')),
    (2, 'App 2', 1, 1, (now() AT TIME ZONE 'UTC'), 1, (now() AT TIME ZONE 'UTC'), 1, 2, (now() AT TIME ZONE 'UTC'), 1, NULL, '1.0.0', 'App 2', 1, (now() AT TIME ZONE 'UTC')),
    (3, 'App 3', 1, 1, (now() AT TIME ZONE 'UTC'), 1, (now() AT TIME ZONE 'UTC'), 1, 3, (now() AT TIME ZONE 'UTC'), 1, NULL, '1.0.0', 'App 3', 1, (now() AT TIME ZONE 'UTC'));

