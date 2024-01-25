Morphos Server
===============

![tests](https://github.com/danvergara/dblab/actions/workflows/test.yaml/badge.svg)

__Self-Hosted file converter server.__

## Table of contents

- [Overview](#overview)
- [Installation](#installation)
    - [Docker](#docker)
- [Features](#features)
- [Usage](#usage)
- [Supported Files](#supported-files-and-convert-matrix)
    - [Images](#images-x-images)
- [License](#license)

## Overview

Today we are forced to rely on third party services to convert files to other formats. This is a serious threat to our privacy, if we use such services to convert files with highly sensitive personal data. It can be used against us, sooner or later.
Morphos server aims to solve the problem mentioned above, by providing a self-hosted server to convert files privately. The project provides an user-friendly web UI.
For now, Morphos only supports images. Documents will be added soon.

## Installation

### Docker

```
docker run --rm -p 8080:8080 -v /tmp:/tmp ghcr.io/danvergara/morphos-server:latest
```

## Features

- Serves a nice web UI
- Simple installation (distributed as a Docker image)

## Usage

Run the server as mentioned above and open up your favorite browser. You'll see something like this:

<img src="screenshots/morphos.png"/>

Hit the file input section on the form to upload an image.

<img src="screenshots/upload_file_morphos.png"/>

You'll see the filed uploaded in the form.

<img src="screenshots/file_uploaded_morphos.png"/>

Then, you can select from a variety of other formats you can convert the current image to.

<img src="screenshots/select_options_morphos.png"/>

After hitting `Upload` button you will see a view like the one below, asking you to download the converted file.

<img src="screenshots/file_converted_morphos.png"/>

A modal will pop up with a preview of the converted image.

<img src="screenshots/modal_morphos.png"/>

## License
The MIT License (MIT). See [LICENSE](LICENSE) file for more details.
