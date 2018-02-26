// Repository name use, must end with / or be '' for none
repository= 'area51/'

// image prefix
imagePrefix = 'nre-feeds'

// The image version, master branch is latest in docker
version=BRANCH_NAME
if( version == 'master' ) {
  version = 'latest'
}

// The architectures to build, in format recognised by docker
architectures = [ 'amd64', 'arm64v8' ]

// The services to build
services = [ 'darwinref', 'darwintt', 'darwind3', 'ldb' ]

// Temp docker image name
tempImage = 'temp/' + imagePrefix + version

// The docker image name
// architecture can be '' for multiarch images
def dockerImage = {
  service, architecture -> repository + imagePrefix +
    ':' + service +
    ( architecture=='' ? '' : ('-' + architecture) ) +
    '-' + version
}

// The go arch
def goarch = {
  architecture -> switch( architecture ) {
    case 'amd64':
      return 'amd64'
    case 'arm32v6':
    case 'arm32v7':
      return 'arm'
    case 'arm64v8':
      return 'arm64'
    default:
      return architecture
  }
}

// goarm is for arm32 only
def goarm = {
  architecture -> switch( architecture ) {
    case 'arm32v6':
      return '6'
    case 'arm32v7':
      return '7'
    default:
      return ''
  }
}

// Build properties
properties([
  buildDiscarder(logRotator(artifactDaysToKeepStr: '', artifactNumToKeepStr: '', daysToKeepStr: '', numToKeepStr: '10')),
  disableConcurrentBuilds(),
  disableResume()
])

def buildArch = {
  architecture ->
    services.each {
      service -> stage( service + ' ' + architecture ) {
        // Modify Dockerfile so the final image has the correct entrypoint
        dockerFile = "Dockerfile." + service
        sh 'sed "s/@@entrypoint@@/' + service + '/g" Dockerfile >' + dockerFile

        sh 'docker build' +
          ' -t ' + dockerImage( service, architecture ) +
          ' -f ' + dockerFile +
          ' --build-arg skipTest=true' +
          ' --build-arg service=' + service +
          ' --build-arg arch=' + architecture +
          ' --build-arg goos=linux' +
          ' --build-arg goarch=' + goarch( architecture ) +
          ' --build-arg goarm=' + goarm( architecture ) +
          ' .'
      }
    }

    if( repository != '' ) {
      // Push all built images relevant docker repository
      stage( 'Publish ' + architecture + ' images' ) {
        services.each {
          service -> sh 'docker push ' + dockerImage( service, architecture )
        }
      }
    } // repository != ''
}

node('AMD64') {
  stage("Checkout") {
    checkout scm
  }

  // Prepare the go base image with the source and libraries
  stage("Prepare Build") {
    // Ensure we have current versions of each base image
    sh 'docker pull golang:alpine'

    // Run up to the source target
    sh 'docker build -t ' + tempImage + ' --target source .'
  }

  // Run unit tests
  stage("Run Tests") {
    sh 'docker build -t ' + tempImage + ' --target test .'
  }
}

parallel (
  'amd64': {
    node('AMD64') {
      build( "amd64" )
    }
  },
  'arm64v8': {
    node('AMD64') {
      build( "arm64v8" )
    }
  }
)

node('AMD64') {
  // Stages valid only if we have a repository set
  if( repository != '' ) {
    // Experimental: Create multi-arch images
    services.each {
      service -> stage( 'Publish ' + service + ' MultiArch image' ) {
        // The manifest to publish
        multiImage = dockerImage( service, '' )

        // Create/amend the manifest with our architectures
        manifests = architectures.collect { architecture -> dockerImage( service, architecture ) }
        sh 'docker manifest create -a ' + multiImage + ' ' + manifests.join(' ')

        // For each architecture annotate them to be correct
        architectures.each {
          architecture -> sh 'docker manifest annotate' +
            ' --os linux' +
            ' --arch ' + goarch( architecture ) +
            ' ' + multiImage +
            ' ' + dockerImage( service, architecture )
        }

        // Publish the manifest
        sh 'docker manifest push -p ' + multiImage
      }
    }
  }

}
