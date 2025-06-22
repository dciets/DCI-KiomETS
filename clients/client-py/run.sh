#! /bin/bash
usage() {
  echo "Usage: $0 <bot_id> <bot_name>"
  echo "Example: $0 1234567890abcdef test"
  exit 1
}
if [ -z "$1" ]; then
    usage
fi
if [ -z "$2" ]; then
    usage
fi

# ./.venv/bin/python3 main.py $1 $2 1>/dev/null &
./.venv/bin/python3 main.py $1 $2 &