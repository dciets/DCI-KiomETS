from typing import List

from actions import Action, create_move_action
from server import Game
import env


class Agent:
    def __init__(self):
        pass

    def update(self, game: Game) -> List[Action]:
        f"""
        Update function
        
        :param game: Game Data
        :return: {List[Action]}
        """
        player_name = env.PLAYER_NAME

        terrains = game.terrains()
        players = game.players()
        pipes = game.pipes()

        # Get the player
        player = [index for index in range(len(players)) if players[index].name() == player_name]

        if len(player) == 1:
            player_index = player[0]

            # Get the terrains of the player
            player_terrains = [(index, terrains[index]) for index in range(len(terrains)) if
                               terrains[index].owner_index() == player_index]
            player_terrains_index = [terrain[0] for terrain in player_terrains]

            orders = []
            for pipe in pipes:
                # If the first terrain of a connection belongs to the player and have at least one soldier
                if (pipe.first() in player_terrains_index and
                        pipe.second() not in player_terrains_index and
                        terrains[pipe.first()].number_of_soldier() > 0):
                    # Send one soldier to the end of the connection
                    orders.append(create_move_action(terrains[pipe.first()].id(), terrains[pipe.second()].id(), 1))

                # If the second terrain of a connection belongs to the player and have at least one soldier
                if (pipe.second() in player_terrains_index and
                        pipe.first() not in player_terrains_index and
                        terrains[pipe.second()].number_of_soldier() > 0):
                    # Send one soldier to the end of the connection
                    orders.append(create_move_action(terrains[pipe.second()].id(), terrains[pipe.first()].id(), 1))

            # Returns the order for each terrain of the player
            return orders

        # Returns no orders since the player was not found
        return []
