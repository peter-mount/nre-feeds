// Repository name use, must end with / or be '' for none
repository= 'area51/'
// Disable deployment until refactor is complete
//repository=''

// image prefix
imagePrefix = 'nre-feeds'

// The git repo / package prefix
gitRepoPrefix = 'github.com/peter-mount/nre-feeds/'

// The image version, master branch is latest in docker
version=BRANCH_NAME
if( version == 'master' ) {
  version = 'latest'
}

// The architectures to build, in format recognised by docker
architectures = [ 'amd64', 'arm64v8' ]

// The services to build
services = [ 'darwinref', 'darwintt', 'darwind3', 'ldb', 'darwinkb' ]

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

// The multi arch image name
def multiImage = { service -> repository + imagePrefix + ':' + service + '-' + version }

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
  disableResume(),
  pipelineTriggers([
    cron('H H * * *')
  ])
])

// Run tests against a suite or library
def runTest = {
  test -> sh 'docker run -i --rm ' + tempImage + ' go test -v ' + gitRepoPrefix + test
}

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

def manifests = {
  service -> manifests = architectures.collect { architecture -> dockerImage( service, architecture ) }
  manifests.join(' ')
}

// Deploy multi-arch image for a service
def multiArchService = {
  service ->
    // Create/amend the manifest with our architectures
    sh 'docker manifest create -a ' + multiImage( service ) + ' ' + manifests( service )

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


// Now build everything on one node
node('AMD64') {
  stage( "Checkout" ) {
    checkout scm

    // Prepare the go base image with the source and libraries
    sh 'docker pull golang:alpine'

    // Run up to the source target so libraries are checked out
    sh 'docker build -t ' + tempImage + ' --target source .'
  }

  // Run unit tests
  stage("Run Tests") {
    parallel (
      'darwind3': { runTest( 'darwind3' ) },
      'darwinref': { runTest( 'darwinref' ) },
      'ldb': { runTest( 'ldb' ) },
      'util': { runTest( 'util' ) },
    )
  }

  // Run issue tests separately as these will grow over time
  stage( "Test Issues" ) {
    runTest( 'issues' )
  }

  services.each {
    service -> stage( service ) {
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
        },
        'darwinkb': {
          multiArchService( 'darwinkb' )
        }
      )
    }
  }

}
