# Morphos Server

## Description

Self-Hosted file converter server.

## Goals of the project

The goals is the creation of a Self-Hosted file converter server to be sovereing on the file manaement. 
The most of the file converter websites claims tha they never keep our files longer than a short period of time, but we can't be so sure if they really keep that promise. Even if so, they can read and process the content of our files in that period of time. We are better off if our files never get the cloud (other people's computers), specially if our files hold very private or legal data. Morphos aims to solve that by providing a highly portable and configurable file converter server.

## Distribution

The project is meant to be distributed by Golang binary embeding all the required assets.
A Container image can be provided as well, but for the most of the cases the compiled binary will be enough, so the lack of container tecnologies installed on the host machine will not be an impediment to run the project.

## Architecture

The project is made up of a Golang server responsible of converting the uploaded files and the fronted composed by Golang templates and a little of HTMX reponsible of uploading files and the user interaction.

## MVP

Just a single view containing a file-select input and a menu showing some dropdowns buttons holding the supported file formats.
