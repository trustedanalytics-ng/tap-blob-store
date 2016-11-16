# Copyright (c) 2016 Intel Corporation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
#!/bin/bash
mkdir -p $MINIO_STORAGE_PATH
exec /opt/app/minio server $MINIO_STORAGE_PATH | sed -u "s/$MINIO_ACCESS_KEY/\$MINIO_ACCESS_KEY/g; s/$MINIO_SECRET_KEY/\$MINIO_SECRET_KEY/g"
