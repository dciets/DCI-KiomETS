#! /bin/bash
usage() {
  echo "Usage: $0 <session_id>"
  echo "This script stops all bots running in the session."
  exit 1
}
if [ -z "$1" ]; then
    usage
fi

for pid in $(ps aux | grep "$1" | grep -v "grep" | grep -v "stop.sh $1" | cut -d " " -f 6)
	do kill -9 $pid
done
for pid in $(ps aux | grep "$1" | grep -v "grep" | grep -v "stop.sh $1" | cut -d " " -f 5)
	do kill -9 $pid
done
