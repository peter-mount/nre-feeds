# Darwin
go library &amp; suite of microservices for handling the NRE DarwinD3 feeds

The main purpose of this project is to consume the feeds provided by National Rail Enquiries in real time and expose that information as a REST service which can be consumed by a client, usually a website.

https://departureboards.mobi/ is an example of one of these clients.

## Microservices

The project is split currently into 4 individual microservices, each in it's own container.

* darwinref handles the reference feed and provides access to data that doesn't change often, including stations, locations, train operators etc.
* darwintt handles the daily timetable feed
* darwind3 handles the real time feed which includes currently running trains and forecasts of their arrival times & delays.
* ldb provides a departure board service allowing the current status to be shown per station.

## Build Status

The builds are available for both AMD64 and ARM64v8 processors.

| Microservice | Architecture | Image | Build Status |
| :----------: | :----------: | ----- | ------------ |
| darwinref | amd64 |  | [![Build Status](http://jenkins.area51.onl/buildStatus/icon?job=UKRail/DarwinRef/microservice=darwinref,slave=AMD64)](http://jenkins.area51.onl/job/UKRail/DarwinRef/microservice=darwinref,slave=AMD64)
| darwinref | arm64v8 |  | [![Build Status](http://jenkins.area51.onl/buildStatus/icon?job=UKRail/DarwinRef/microservice=darwinref,slave=ARM64v8)](http://jenkins.area51.onl/job/UKRail/DarwinRef/microservice=darwinref,slave=ARM64v8)
| darwintt | amd64 |  | [![Build Status](http://jenkins.area51.onl/buildStatus/icon?job=UKRail/DarwinRef/microservice=darwintt,slave=AMD64)](http://jenkins.area51.onl/job/UKRail/DarwinRef/microservice=darwintt,slave=AMD64)
| darwintt | arm64v8 |  | [![Build Status](http://jenkins.area51.onl/buildStatus/icon?job=UKRail/DarwinRef/microservice=darwintt,slave=ARM64v8)](http://jenkins.area51.onl/job/UKRail/DarwinRef/microservice=darwintt,slave=ARM64v8)
| darwind3 | amd64 |  | [![Build Status](http://jenkins.area51.onl/buildStatus/icon?job=UKRail/DarwinRef/microservice=darwind3,slave=AMD64)](http://jenkins.area51.onl/job/UKRail/DarwinRef/microservice=darwind3,slave=AMD64)
| darwind3 | arm64v8 |  | [![Build Status](http://jenkins.area51.onl/buildStatus/icon?job=UKRail/DarwinRef/microservice=darwind3,slave=ARM64v8)](http://jenkins.area51.onl/job/UKRail/DarwinRef/microservice=darwind3,slave=ARM64v8)
| ldb | amd64 |  | [![Build Status](http://jenkins.area51.onl/buildStatus/icon?job=UKRail/DarwinRef/microservice=ldb,slave=AMD64)](http://jenkins.area51.onl/job/UKRail/DarwinRef/microservice=ldb,slave=AMD64)
| ldb | arm64v8 |  | [![Build Status](http://jenkins.area51.onl/buildStatus/icon?job=UKRail/DarwinRef/microservice=ldb,slave=ARM64v8)](http://jenkins.area51.onl/job/UKRail/DarwinRef/microservice=ldb,slave=ARM64v8)

## Volumes

The Docker image requires two mounts (the -v above):

### Database

We need to mount a volume to store the databases. This defaults to /database but that can be changed within the configuration.

If you don't mount /database then the database will be wiped when the container is stopped & destroyed. In the above above:

    -v /tmp/darwin/:/database

we are defining it as a local directory mount as /tmp/darwin/ but you could use a container. See the Docker documentation on other alternatives.

### Configuration

This is a simple yaml file which is expected to be at /config.yaml

In the above example we mounted the example config:

    -v $(pwd)/config-example.yaml:ro

Here $(pwd)/config-example.yaml presumes you are running this in the repository root. The :ro mounts it as read-only.

Again, there are alternatives to this method, as long as the file is accessible in the container.

## config.yaml

This is the main configuration file. config-example.yaml is fully documented on what each option is.
