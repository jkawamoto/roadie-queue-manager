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

# Install apt packages specified in the script file to be run.
{{with .APT}}
RUN apt-get update
{{range .}}
RUN apt-get install -y {{.}}
{{end}}
{{end}}

WORKDIR /data
ADD .roadie/entrypoint.sh /root/entrypoint.sh
ENTRYPOINT ["bash", "/root/entrypoint.sh"]
CMD [""]
