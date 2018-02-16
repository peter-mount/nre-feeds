/*

The problem: We have multiple builds with the same Dockerfile which we
select via --build-arg service= and ARG service inside the Dockerfile.

This works fine up until we use ENTRYPOINT or CMD which cannot use ${service}
as it's only valid within the Dockerfile.

So we could either:

1. name each service the same within each container, which is fine but later
   we would like to create non-docker builds - i.e. use docker to build but
   provide the binaries separately, just a case of extracting from the image.

2. Use a symlink to the binary. This works but the process table has the name
   of the symlink & the microservice sees that name as it's name not it's name.
   e.g. os.Args[0] == "/docker-entrypoint" when we would expect it to be "ldb".

Hence this wrapper which simply replaces itself with the microservice, ensuring
that it's name appears in os.Args[0] and in the process table.

The final output is about 9k so it doesn't add much to the final image size.

In the Dockerfile we have 2 stages. The first one is at the beginning as it
creates a small alpine based gcc environment

FROM alpine as gcc
RUN apk add --no-cache gcc musl-dev

Then when the test stage completes we run the second stage which gets this file,
replaces any instance of @@service@@ with the service name and compiles it.

FROM gcc as wrapper
ARG service
WORKDIR /work
ADD bin/docker .
RUN sed -i "s/@@service@@/${service}/g" main.c &&\
    gcc -o main -static main.c &&\
    strip main

Finally after the microservice has compiled we copy the wrapper into the
distribution:

COPY --from=wrapper /work/main /dest/docker-entrypoint

Note: We do that here & not in the final image as this ensures we have a single
layer & not a tiny one just for the wrapper ;-)

 */
#include <unistd.h>
#include <stdlib.h>
#include <string.h>

int main( const int argc, const char *argv[] ) {

  // Create new args string which contains any arguments but argv[0] is the new binary name
  char **args = malloc( sizeof( char* ) * argc );

  for( int i=1; i<argc; i++ ) {
    args[i] = (char *) argv[i];
  }

  args[0] = strdup( "@@service@@" );

  execvp( "/@@service@@", args );
}
