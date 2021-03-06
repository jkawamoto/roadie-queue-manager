# Roadie Queue Manager
[![GPLv3](https://img.shields.io/badge/license-GPLv3-blue.svg)](https://www.gnu.org/copyleft/gpl.html)
[![Build Status](https://travis-ci.org/jkawamoto/roadie-queue-manager.svg?branch=master)](https://travis-ci.org/jkawamoto/roadie-queue-manager)
[![Code Climate](https://codeclimate.com/github/jkawamoto/roadie-queue-manager/badges/gpa.svg)](https://codeclimate.com/github/jkawamoto/roadie-queue-manager)
[![Release](https://img.shields.io/badge/release-0.2.3-brightgreen.svg)](https://github.com/jkawamoto/roadie/releases/tag/v0.2.3)

A helper tool of Roadie to execute your scripts with a queue.

## Description
Roadie Queue Manager is a helper tool of [Roadie](https://jkawamoto.github.io/roadie/).
This tool checks a queue implemented on [Google Cloud Datastore](https://cloud.google.com/datastore/)
and executes scripts enqueued in the queue.

## Usage
```shell
$ roadie-queue-manager <project ID> <queue name>
```

## License
This software is released under The GNU General Public License Version 3,
see [COPYING](COPYING) and [LICENSES](LICENSES.md) for more detail.
