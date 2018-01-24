# Darwin
go library &amp; Application for handling the NRE DarwinD3 feeds

The main purpose of this project is to consume the feeds provided by National Rail Enquiries in real time and expose that information as a REST service which can be consumed by a client, usually a website.

https://departureboards.mobi/ is an example of one of these clients.

## Running

Not yet available but a pre-built image will be available on Docker Hub.
For now see Building below.

## Building

As this is intended to be run within a docker container, the build is entirely within the Dockerfile. Just clone this repository and from the base directory run:

    docker build -t mytag .

To run that built image you can easily run it with:

    docker run -it --rm \
      -v /tmp/darwin/:/database \
      -v $(pwd)/config-example.yaml:ro \
      -p 8081:8081 \
      mytag

See below for more details on this.

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

This is the main configuration file. config-example.yaml is fully documented on what each option is, but it's broken into several sections:

* server configures the embedded web server
* database configures the database
* ftp holds your credentials to the NRE FTP server for retrieving the daily reference files.
* statistics Allows for monitoring of the components of the application.

### database

This can be kept as it is. By default if nothing is set for path then /database will be used. If this section can be absent & all defaults will be used.

### server

This configures the webserver. If absent then port 8080 will be used.

### statistics

This controls the statistics used to monitor the application. If not set then
the monitoring is hidden.

* log set to true then once a minute statistics are logged to the console. This is useful when used with docker as you can export the logs to a remote service like AWS CloudWatch and extract those values.

* schedule defines how often the statistics are captured & reported. By default this is once per minute.
