
# Execute a command and abort if it fails
function execute() {
  CMD=$@
  echo $CMD
  $CMD || exit $?
}

# Resolve the goarch for the architecture
function goarch() {
  ARCH=$1

  case $ARCH in
    amd64)
    echo amd64
    ;;
    arm32v6)
    echo arm
    ;;
    arm32v7)
    echo arm
    ;;
    arm64v8)
    echo arm64
    ;;
    *)
    echo "Unsupported architecture $ARCH"
    exit 1
    ;;
  esac
}

# Resolve the goarm value for the architecture.
function goarm() {
  ARCH=$1

  # Resolve the architecture
  case $ARCH in
    arm32v6)
    echo 6
    ;;
    arm32v7)
    echo 7
    ;;
    *)
    echo
    ;;
  esac
}
