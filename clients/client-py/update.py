from typing import List

from actions import Action, create_move_action
from server import Game


class Bot:
    def __init__(self):
        pass

    def update(self, game: Game) -> List[Action]:
        return [
            create_move_action('1', '2', 3)
        ]

