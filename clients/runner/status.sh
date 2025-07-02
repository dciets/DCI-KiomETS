usage() {
  echo "Usage: $0 <session_id>"
  echo "This script checks the status of bots running in the session."
  exit 1
}
if [ -z "$1" ]; then
    usage
fi

ps aux | grep "$1" | grep -v "grep" | grep -v "status.sh $1"