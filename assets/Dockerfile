#
# Dockerfile
#
# Copyright (c) 2016-2017 Junpei Kawamoto
#
# This file is part of Roadie queue manager.
#
# Roadie queue manager is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# Roadie queue manager is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with Roadie queue manager. If not, see <http://www.gnu.org/licenses/>.
#

# This template is a dockerfile which installs apt packages.
#
FROM {{.Image}}
MAINTAINER Junpei Kawamoto <kawamoto.junpei@gmail.com>

ENV TERM vt100
ENV DEBIAN_FRONTEND noninteractive

# Install necessary packages
RUN apt-get update && \
    apt-get install -y apt-utils curl git gcc python-dev libffi-dev python-pip
RUN pip install -U pip gsutil
RUN echo "[GoogleCompute]\nservice_account = default\n[GSUtil]\nparallel_composite_upload_threshold = 150M" >> /etc/boto.cfg

# Install apt packages specified in the script file to be run.
{{range .APT}}
RUN apt-get install -y {{.}}
{{end}}

WORKDIR /data
ADD .roadie/entrypoint.sh /root/entrypoint.sh
ENTRYPOINT ["bash", "/root/entrypoint.sh"]
CMD [""]
