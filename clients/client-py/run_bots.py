from argparse import ArgumentParser
import env
import subprocess
import os
import requests
from uuid import uuid4

args = ArgumentParser(description="Run bots in parallel")
args.add_argument("bots_amount", type=int, help="Number of bots to run")
parsed_args = args.parse_args()

if parsed_args.bots_amount <= 0:
    raise ValueError("The number of bots must be a positive integer.")
game_settings = {
    "mapSize": 10,
    "soldierSpeed": 1,
    "soldierCreationSpeed": 1,
    "terrainChangeSpeed": 1,
    "gameLength": 120,
}
res = requests.put(f"http://{env.HOST}:8080/api/game", json=game_settings)
res.raise_for_status()  # Ensure the request was successful
res = requests.post(f"http://{env.HOST}:8080/api/start")
if res.status_code in [200, 409]:
    print("Game started successfully or already running.")

for i in range(parsed_args.bots_amount):
    bot_name = f"bot_{i+1}_{uuid4().hex}"
    res = requests.post(f"http://{env.HOST}:8080/api/agent", json={"name": bot_name})
    res.raise_for_status()  # Ensure the request was successful
    bot_id = res.json().get("UID")
    command = ["/bin/bash", "run.sh", bot_id, bot_name]

    # Run the bot in a new process
    subprocess.Popen(command, cwd=os.path.dirname(os.path.abspath(__file__)))

while True:
    try:
        s = input("Press any key to print status or Ctrl+C to stop all bots...\n")
        if s:
            subprocess.Popen(
                ["/bin/bash", "status.sh"],
                cwd=os.path.dirname(os.path.abspath(__file__)),
            )
    except KeyboardInterrupt:
        print("Stopping all bots...")
        subprocess.Popen(
            ["/bin/bash", "stop.sh"], cwd=os.path.dirname(os.path.abspath(__file__))
        )
        break  # Exit the loop on keyboard interrupt
