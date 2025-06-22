#! /bin/bash
for pid in $(ps aux | grep "./.venv/bin/python3 main" | grep -v "grep" | cut -d " " -f 6)
	do kill -9 $pid
done
for pid in $(ps aux | grep "./.venv/bin/python3 main" | grep -v "grep" | cut -d " " -f 5)
	do kill -9 $pid
done
