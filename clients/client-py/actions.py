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
    move: MoveAction or None
    build: BuildAction or None

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
    move = MoveAction(terrain_from_id, terrain_to_id, quantity_of_soldier)
    return Action(0, move, None)


def create_build_barricade_action(terrain_id: str) -> Action:
    build = BuildAction(terrain_id, 0)
    return Action(1, None, build)


def create_build_factory_action(terrain_id: str) -> Action:
    build = BuildAction(terrain_id, 1)
    return Action(1, None, build)


def create_demolish_action(terrain_id: str) -> Action:
    build = BuildAction(terrain_id, 2)
    return Action(1, None, build)
