import json
from dataclasses import dataclass
from typing import Any, Dict


@dataclass
class MoveAction:
    from_id: str
    to_id: str
    quantity: int

    def serialize(self) -> Dict[str, Any]:
        return {
            "fromId": self.from_id,
            "toId": self.to_id,
            "quantity": self.quantity,
        }


@dataclass
class BuildAction:
    terrain_id: str
    terrain_type: int

    def serialize(self) -> Dict[str, Any]:
        return {
            "terrainId": self.terrain_id,
            "terrainType": self.terrain_type,
        }


@dataclass
class Action:
    action_type: int
    move: MoveAction | None
    build: BuildAction | None

    def serialize(self) -> Dict[str, Any]:
        sr: Dict[str, Any] = {
            "actionType": self.action_type,
        }

        if self.move is None:
            sr["move"] = None
        else:
            sr["move"] = self.move.serialize()

        if self.build is None:
            sr["build"] = None
        else:
            sr["build"] = self.build.serialize()

        return sr


def create_move_action(terrain_from_id: str, terrain_to_id: str, quantity_of_soldier: int) -> Action:
    """
    Create a move action to send to the server
    :param terrain_from_id: Source terrain id
    :param terrain_to_id: Target adjacent terrain id
    :param quantity_of_soldier: Quantity of soldier to send
    :return:
    """
    move = MoveAction(terrain_from_id, terrain_to_id, quantity_of_soldier)
    return Action(0, move, None)


def create_build_barricade_action(terrain_id: str) -> Action:
    """
    Create a build barricade action to send to the server
    :param terrain_id: Terrain id to build the barricade
    :return:
    """
    build = BuildAction(terrain_id, 0)
    return Action(1, None, build)


def create_build_factory_action(terrain_id: str) -> Action:
    """
    Create a build factory action to send to the server
    :param terrain_id: Terrain id to build the factory
    :return:
    """
    build = BuildAction(terrain_id, 1)
    return Action(1, None, build)


def create_demolish_action(terrain_id: str) -> Action:
    """
    Create a demolish action to send to the server
    :param terrain_id: Terrain id to destroy the building
    :return:
    """
    build = BuildAction(terrain_id, 2)
    return Action(1, None, build)
