properties([
  buildDiscarder(
    logRotator(
      artifactDaysToKeepStr: '',
      artifactNumToKeepStr: '',
      daysToKeepStr: '',
      numToKeepStr: '10'
    )
  ),
  disableConcurrentBuilds(),
  disableResume(),
  pipelineTriggers([
    cron("H H * * *")
  ])
])
node("go") {
  stage("Checkout") {
    checkout scm
  }
  stage("Init") {
    sh 'make clean init test'
  }
  stage("darwin_amd64") {
    sh 'make -f Makefile.gen darwin_amd64'
  }
  stage("darwin_arm64") {
    sh 'make -f Makefile.gen darwin_arm64'
  }
  stage("dragonfly_amd64") {
    sh 'make -f Makefile.gen dragonfly_amd64'
  }
  stage("freebsd_386") {
    sh 'make -f Makefile.gen freebsd_386'
  }
  stage("freebsd_amd64") {
    sh 'make -f Makefile.gen freebsd_amd64'
  }
  stage("freebsd_arm6") {
    sh 'make -f Makefile.gen freebsd_arm6'
  }
  stage("freebsd_arm7") {
    sh 'make -f Makefile.gen freebsd_arm7'
  }
  stage("freebsd_arm64") {
    sh 'make -f Makefile.gen freebsd_arm64'
  }
  stage("freebsd_riscv64") {
    sh 'make -f Makefile.gen freebsd_riscv64'
  }
  stage("illumos_amd64") {
    sh 'make -f Makefile.gen illumos_amd64'
  }
  stage("linux_386") {
    sh 'make -f Makefile.gen linux_386'
  }
  stage("linux_amd64") {
    sh 'make -f Makefile.gen linux_amd64'
  }
  stage("linux_arm6") {
    sh 'make -f Makefile.gen linux_arm6'
  }
  stage("linux_arm7") {
    sh 'make -f Makefile.gen linux_arm7'
  }
  stage("linux_arm64") {
    sh 'make -f Makefile.gen linux_arm64'
  }
  stage("linux_mips") {
    sh 'make -f Makefile.gen linux_mips'
  }
  stage("linux_mips64") {
    sh 'make -f Makefile.gen linux_mips64'
  }
  stage("linux_mips64le") {
    sh 'make -f Makefile.gen linux_mips64le'
  }
  stage("linux_mipsle") {
    sh 'make -f Makefile.gen linux_mipsle'
  }
  stage("linux_ppc64") {
    sh 'make -f Makefile.gen linux_ppc64'
  }
  stage("linux_ppc64le") {
    sh 'make -f Makefile.gen linux_ppc64le'
  }
  stage("linux_riscv64") {
    sh 'make -f Makefile.gen linux_riscv64'
  }
  stage("linux_s390x") {
    sh 'make -f Makefile.gen linux_s390x'
  }
  stage("netbsd_386") {
    sh 'make -f Makefile.gen netbsd_386'
  }
  stage("netbsd_amd64") {
    sh 'make -f Makefile.gen netbsd_amd64'
  }
  stage("netbsd_arm6") {
    sh 'make -f Makefile.gen netbsd_arm6'
  }
  stage("netbsd_arm7") {
    sh 'make -f Makefile.gen netbsd_arm7'
  }
  stage("netbsd_arm64") {
    sh 'make -f Makefile.gen netbsd_arm64'
  }
  stage("openbsd_386") {
    sh 'make -f Makefile.gen openbsd_386'
  }
  stage("openbsd_amd64") {
    sh 'make -f Makefile.gen openbsd_amd64'
  }
  stage("openbsd_arm6") {
    sh 'make -f Makefile.gen openbsd_arm6'
  }
  stage("openbsd_arm7") {
    sh 'make -f Makefile.gen openbsd_arm7'
  }
  stage("openbsd_arm64") {
    sh 'make -f Makefile.gen openbsd_arm64'
  }
  stage("solaris_amd64") {
    sh 'make -f Makefile.gen solaris_amd64'
  }
  stage("windows_386") {
    sh 'make -f Makefile.gen windows_386'
  }
  stage("windows_amd64") {
    sh 'make -f Makefile.gen windows_amd64'
  }
  stage("windows_arm6") {
    sh 'make -f Makefile.gen windows_arm6'
  }
  stage("windows_arm7") {
    sh 'make -f Makefile.gen windows_arm7'
  }
  stage("windows_arm64") {
    sh 'make -f Makefile.gen windows_arm64'
  }
  stage("archiveArtifacts") {
    archiveArtifacts artifacts: 'dist/*'
  }
}