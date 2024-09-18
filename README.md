Morphos Server
===============

![tests](https://github.com/danvergara/morphos/actions/workflows/test.yml/badge.svg)
[![Release](https://img.shields.io/github/release/danvergara/morphos.svg?label=Release)](https://github.com/danvergara/morphos/releases)

<p align="center">
  <img style="float: right;" src="screenshots/morphos.jpg" alt="morphos logo"/  width=200>
</p>

__Self-Hosted file converter server.__

## Table of contents

- [Overview](#overview)
- [Installation](#installation)
    - [Docker](#docker)
- [Features](#features)
- [Usage](#usage)
- [Supported Files](#supported-files-and-convert-matrix)
    - [Images To Images](#images-x-images)
    - [Images To Documents](#images-x-documents)
    - [Documents To Images](#documents-x-images)
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

### HTML form

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

### API

You can consume morphos through an API, so other systems can integrate with it.

##### Endpoints

`GET /api/v1/formats`

This returns a JSON that shows the supported formats at the moment.

e.g.

```
{"documents": ["docx", "xls"], "image": ["png", "jpeg"]}
```

`POST /api/v1/upload`

This is the endpoint that converts files to a desired format. It is basically a multipart form data in a POST request. The API simply writes the converted files to the response body.

e.g.

```
 curl -F 'targetFormat=epub' -F 'uploadFile=@/path/to/file/foo.pdf' localhost:8080/api/v1/upload --output foo.epub
```
The form fields are:

* targetFormat: the target format the file will be converted to
* uploadFile: The path to the file that is going to be converted

### Configuration

The configuration is only done by the environment varibles shown below.

* `MORPHOS_PORT` changes the port the server will listen to (default is `8080`)
* `MORPHOS_UPLOAD_PATH` defines the temporary path the files will be stored on disk (default is `/tmp`)

## Supported Files And Convert Matrix

### Images X Images

|       |  PNG  |  JPEG  |  GIF  |  WEBP  |  TIFF  |  BMP  |  AVIF  |
|-------|-------|--------|-------|--------|--------|-------|--------|  
|  PNG  |       |   ✅   |  ✅   |   ✅   |   ✅   |  ✅   |   ✅   |
|  JPEG |  ✅   |        |  ✅   |   ✅   |   ✅   |  ✅   |   ✅   |  
|  GIF  |  ✅   |   ✅   |       |   ✅   |   ✅   |  ✅   |   ✅   | 
|  WEBP |  ✅   |   ✅   |  ✅   |        |   ✅   |  ✅   |   ✅   |
|  TIFF |  ✅   |   ✅   |  ✅   |   ✅   |        |  ✅   |   ✅   |
|  BMP  |  ✅   |   ✅   |  ✅   |   ✅   |   ✅   |       |   ✅   |
|  AVIF |  ✅   |   ✅   |  ✅   |   ✅   |   ✅   |  ✅   |        |

### Images X Documents

|       |  PDF  |
|-------|-------|
|  PNG  |  ✅   |
|  JPEG |  ✅   |
|  GIF  |  ✅   |
|  WEBP |  ✅   |
|  TIFF |  ✅   |
|  BMP  |  ✅   |
|  AVIF |       |

## Documents X Images

|     | PNG | JPEG | GIF | WEBP | TIFF | BMP |  AVIF | 
| --- | --- | ---- | --- | ---- | ---- | --- | ----  |
| PDF | ✅  | ✅   | ✅  | ✅   | ✅   | ✅  |       |

## Documents X Documents

|      | DOCX | PDF | XLSX | CSV |
| ---- | ---- | --- | ---- | --- |
| PDF  | ✅   |     |      |     |
| DOCX |      | ✅  |      |     |
| CSV  |      |     |  ✅  |     |
| XLSX |      |     |      | ✅  |

## Ebooks X Ebooks

|      | MOBI | EPUB |
| ---- | ---- | --- |
| EPUB | ✅   |     |
| MOBI |      | ✅  |


## Documents X Ebooks

|      | EPUB | MOBI |
| ---- | ---- | ---  |
| PDF  | ✅   | ✅   |
| DOCX |      |      |
| CSV  |      |      |
| XLSX |      |      |

## Ebooks X Documents

|      | PDF  |
| ---- | ---- |
| EPUB |  ✅  |
| MOBI |  ✅  |

## License
The MIT License (MIT). See [LICENSE](LICENSE) file for more details.
