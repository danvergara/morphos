Morphos Server
===============

![tests](https://github.com/danvergara/dblab/actions/workflows/test.yaml/badge.svg)

__Self-Hosted file converter server.__

## Table of contents

- [Overview](#overview)
- [Dependencies](#Dependencies)
- [Installation](#Installation)
- [Features](#features)
- [License](#license)

## Overview

Today we are forced to rely on third party services to convert files to other formats. This is a serious threat to our privacy, if we use such services to convert files with highly sensitive personal data. It can be used against us, sooner or later.
Morphos server aims to solve the problem mentioned above, by providing a self-hosted server to convert files privately. The project provides an user-friendly web UI.

## Dependencies

* [Go 1.21](https://go.dev/doc/devel/release#go1.21.0)
* [air for local development](https://github.com/cosmtrek/air)
* [Docker](https://docs.docker.com/engine/install/)
* [Make](https://www.gnu.org/software/make/)

## Installation

The project is written in Go 1.21.


1. You can run the project on bare metal (this uses air for live-reloading):

```
$ brew install cmake make pkg-config x265 libde265 libjpeg libtool aom
$ brew install libheif
$ make run
```

2. On a container (make sure docker is installed)

```
$ make docker-build
$ make docker-run
```

## Features

- Serves a nice web UI
- Simple installation (distributed as a Docker image)

## License
The MIT License (MIT). See [LICENSE](LICENSE) file for more details.
