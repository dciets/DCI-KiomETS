#! /bin/bash
usage() {
  echo "Usage: $0 <bot_id> <bot_name> <mode> <session_id>"
  echo "Example: $0 1234567890abcdef test js 12345678"
  exit 1
}
if [ -z "$1" ]; then
    usage
fi
if [ -z "$2" ]; then
    usage
fi

if [ -z "$3" ]; then
    usage
fi

if [ -z "$4" ]; then
    usage
fi

if [ "$3" != "js" ] && [ "$3" != "py" ]; then
    echo "Mode must be either 'js' or 'py'."
    exit 1
fi

if [ "$3" == "js" ]; then
    export NODE_OPTIONS=--max-old-space-size=8192
    node ../client-js/src/index.js $1 $2"_"$4 1>/dev/null &
else
    export PYTHONUNBUFFERED=1
    ../client-py/.venv/bin/python3 ../client-py/main.py $1 $2"_"$4 1>/dev/null &
fi

# ./.venv/bin/python3 main.py $1 $2 &