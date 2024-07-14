from typing import List

from actions import Action, create_move_action
from server import Game


class Agent:
    def __init__(self):
        pass

    def update(self, game: Game) -> List[Action]:
        player_name = 'test'
        terrains = game.terrains()
        players = game.players()
        pipes = game.pipes()
        player = [index for index in range(len(players)) if players[index].name() == player_name]
        if len(player) == 1:
            player_index = player[0]
            player_terrains = [(index, terrains[index]) for index in range(len(terrains)) if terrains[index].owner_index() == player_index]
            player_terrains_index = [terrain[0] for terrain in player_terrains]
            orders = []
            for pipe in pipes:
                if pipe.first() in player_terrains_index and pipe.second() not in player_terrains_index and terrains[pipe.first()].number_of_soldier() > 0:
                    orders.append(create_move_action(terrains[pipe.first()].id(), terrains[pipe.second()].id(), 1))
                if pipe.second() in player_terrains_index and pipe.first() not in player_terrains_index and terrains[pipe.second()].number_of_soldier() > 0:
                    orders.append(create_move_action(terrains[pipe.second()].id(), terrains[pipe.first()].id(), 1))

            return orders
        return []

