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
tempImage = 'temp/' + imagePrefix + ':' + version

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

def dockerFile = { architecture, service -> "Dockerfile." + service + '.' + architecture }

// Build a service for a specific architecture
def buildArch = {
  architecture, service ->
    // Modify Dockerfile so the final image has the correct entrypoint
    sh 'sed "s/@@entrypoint@@/' + service + '/g" Dockerfile >' + dockerFile( architecture, service )

    sh 'docker build' +
      ' -t ' + dockerImage( service, architecture ) +
      ' -f ' + dockerFile( architecture, service ) +
      ' --build-arg skipTest=true' +
      ' --build-arg service=' + service +
      ' --build-arg arch=' + architecture +
      ' --build-arg goos=linux' +
      ' --build-arg goarch=' + goarch( architecture ) +
      ' --build-arg goarm=' + goarm( architecture ) +
      ' .'

    if( repository != '' ) {
      // Push all built images relevant docker repository
      sh 'docker push ' + dockerImage( service, architecture )
    } // repository != ''
}

def multiImage = { service -> dockerImage( service, '' ) }
// Deploy multi-arch image for a service
def multiArchService = {
  service -> {

    // Create/amend the manifest with our architectures
    manifests = architectures.collect { architecture -> dockerImage( service, architecture ) }
    sh 'docker manifest create -a ' + multiImage( service ) + ' ' + manifests.join(' ')

    // For each architecture annotate them to be correct
    architectures.each {
      architecture -> sh 'docker manifest annotate' +
        ' --os linux' +
        ' --arch ' + goarch( architecture ) +
        ' ' + multiImage( service ) +
        ' ' + dockerImage( service, architecture )
    }

    // Publish the manifest
    sh 'docker manifest push -p ' + multiImage( service )
  }
}

// Now build everything on one node
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

  services.each {
    service -> stage( 'Build ' + service ) {
      parallel (
        'amd64': {
          buildArch( "amd64", service )
        },
        'arm64v8': {
          buildArch( "arm64v8", service )
        }
      )
    }
  }

  // Stages valid only if we have a repository set
  if( repository != '' ) {
    stage( "Multiarch Image" ) {
      parallel(
        'darwinref': {
          multiArchService( 'darwinref' )
        },
        'darwintt': {
          multiArchService( 'darwintt' )
        },
        'darwind3': {
          multiArchService( 'darwind3' )
        },
        'ldb': {
          multiArchService( 'ldb' )
        }
      )
    }
  }

}
