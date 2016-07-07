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

FROM tapimages.us.enableiot.com:8080/tap-base-binary:binary-jessie
MAINTAINER Pawel Pieniazek <pawel.pieniazek@intel.com>


RUN mkdir -p /opt/app
ADD build/tapng-blob-store /opt/app

RUN chmod +x /opt/app/tapng-blob-store

WORKDIR /opt/app/

ENV BLOB_STORE_PORT "8084"
EXPOSE 8084

ENTRYPOINT ["/opt/app/tapng-blob-store"]
CMD [""]
