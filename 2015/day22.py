from abc import ABC, abstractmethod
from collections import deque
from typing import Callable, List, Type
import math
import file_reader


class Sprite:
    pass


class Effect:
    pass


class Effect:
    def __init__(
        self, duration: int, effect: Callable[[Sprite, Sprite], None], effect_type: Type
    ) -> None:
        self._duration = duration
        self._effect = effect
        self._effect_type = effect_type
        self._on_finished = lambda x, y: None

    def SetOnSpellFinishCallback(
        self, onSpellFinished: Callable[[Sprite, Sprite], None]
    ) -> None:
        self._on_finished = onSpellFinished

    def ApplyEffect(self, hero: Sprite, boss: Sprite) -> None:
        if self._duration != 0:
            self._effect(hero, boss)
            self._duration -= 1
            if self._duration == 0:
                self._on_finished(hero, boss)

    def Copy(self) -> Effect:
        effect = Effect(
            duration=self.Duration, effect=self.Effect, effect_type=self.EffectType
        )
        effect.SetOnSpellFinishCallback(self._on_finished)
        return effect

    @property
    def EffectType(self) -> Type:
        return self._effect_type

    @property
    def Effect(self) -> Callable[[Sprite, Sprite], None]:
        return self._effect

    @property
    def Duration(self) -> int:
        return self._duration


class Sprite:
    def __init__(self, health: int, mana: int, armor: int, dmg: int) -> None:
        self._health = health
        self._mana = mana
        self._armor = armor
        self._dmg = dmg

    def Copy(self) -> Sprite:
        return Sprite(
            health=self.Health, mana=self.Mana, armor=self.Armor, dmg=self.Damage
        )

    def TakeMagicDmg(self, dmg: int) -> None:
        self._health -= dmg

    def TakeDmg(self, dmg: int) -> None:
        self._health -= 1 if self._armor >= dmg else (dmg - self._armor)

    def IncreaseMana(self, amount: int) -> None:
        self.DrainMana(-amount)

    def DrainMana(self, amount: int) -> None:
        self._mana -= amount

    def IncreaseArmor(self, amount: int) -> None:
        self._armor += amount

    @property
    def Health(self) -> int:
        return self._health

    @property
    def Damage(self) -> int:
        return self._dmg

    @property
    def Mana(self) -> int:
        return self._mana

    @property
    def Armor(self) -> int:
        return self._armor


class NotEnoughManaError(Exception):
    pass


class EffectAlreadyInPlace(Exception):
    pass


class Dead(Exception):
    pass


class Spell(ABC):
    @property
    @abstractmethod
    def ManaCost(self) -> int:
        raise NotImplementedError()

    @abstractmethod
    def _Cast(self, hero: Sprite, boss: Sprite) -> Effect:
        raise NotImplementedError()

    def CanCastSpell(self, effects: List[Effect]) -> bool:
        return not any(x.Duration != 0 and x.EffectType == type(self) for x in effects)

    def Cast(self, hero: Sprite, boss: Sprite, effects: List[Effect]) -> None:
        if hero.Mana < self.ManaCost:
            raise NotEnoughManaError("Need more mana")
        if not self.CanCastSpell(effects):
            raise EffectAlreadyInPlace("Adding same effect twice not allowed")
        hero.DrainMana(self.ManaCost)
        if effect := self._Cast(hero=hero, boss=boss):
            effects.append(effect)


class MagicMissile(Spell):
    @property
    def ManaCost(self):
        return 53

    def _Cast(self, hero: Sprite, boss: Sprite) -> Effect:
        boss.TakeMagicDmg(4)
        return None


class Drain(Spell):
    @property
    def ManaCost(self) -> int:
        return 73

    def _Cast(self, hero: Sprite, boss: Sprite) -> Effect:
        boss.TakeMagicDmg(2)
        hero.TakeMagicDmg(-2)
        return None


class Shield(Spell):
    @property
    def ManaCost(self) -> int:
        return 113

    def _Cast(self, hero: Sprite, boss: Sprite) -> Effect:
        hero.IncreaseArmor(7)
        effect = Effect(
            duration=6, effect=lambda hero, boss: None, effect_type=type(self)
        )
        effect.SetOnSpellFinishCallback(lambda hero, boss: hero.IncreaseArmor(-7))
        return effect


class Poison(Spell):
    @property
    def ManaCost(self) -> int:
        return 173

    def _Cast(self, hero: Sprite, boss: Sprite) -> Effect:
        return Effect(
            effect=lambda hero, boss: boss.TakeMagicDmg(3),
            duration=6,
            effect_type=type(self),
        )


class Recharge(Spell):
    @property
    def ManaCost(self) -> int:
        return 229

    def _Cast(self, hero: Sprite, boss: Sprite) -> Effect:
        return Effect(
            effect=lambda hero, boss: hero.IncreaseMana(101),
            duration=5,
            effect_type=type(self),
        )


SPELLS: List[Effect] = [MagicMissile(), Drain(), Shield(), Poison(), Recharge()]


def applyEffects(hero: Sprite, boss: Sprite, effects: List[Effect]) -> None:
    for effect in effects:
        effect.ApplyEffect(hero, boss)
    if hero.Health <= 0:
        raise Dead


def fightBossWithSpell(
    hero: Sprite, boss: Sprite, effects: List[Effect], spell: Spell
) -> tuple[Sprite, Sprite, List[Effect]]:
    applyEffects(hero=hero, boss=boss, effects=effects)
    if boss.Health <= 0:
        return (hero, boss, effects)

    spell.Cast(hero=hero, boss=boss, effects=effects)

    applyEffects(hero=hero, boss=boss, effects=effects)
    if boss.Health <= 0:
        return (hero, boss, effects)

    hero.TakeDmg(boss.Damage)
    if hero.Health <= 0:
        raise Dead

    return (hero, boss, effects)


def fightBoss(hero: Sprite, boss: Sprite, effects: List[Effect] = None) -> int:
    if effects is None:
        effects: List[Effect] = []
    min_cost = math.inf
    fight_plans = deque([(0, hero, boss, effects)])
    while fight_plans:
        cost, hero, boss, effects = fight_plans.popleft()

        for spell in SPELLS:
            if min_cost <= cost + spell.ManaCost:
                continue
            try:
                new_hero, new_boss, new_effects = fightBossWithSpell(
                    hero.Copy(), boss.Copy(), [x.Copy() for x in effects], spell
                )
                if new_boss.Health <= 0:
                    min_cost = cost + spell.ManaCost
                else:
                    fight_plans.append(
                        (cost + spell.ManaCost, new_hero, new_boss, new_effects)
                    )
            except (NotEnoughManaError, EffectAlreadyInPlace, Dead):
                continue

    return min_cost


def main() -> None:
    programInput = file_reader.getInput().splitlines()

    bosshealth = [
        int(x.removeprefix("Hit Points: "))
        for x in programInput
        if x.startswith("Hit Points")
    ][0]
    boosDamage = [
        int(x.removeprefix("Damage: ")) for x in programInput if x.startswith("Damage")
    ][0]

    print(
        "part 1: ",
        fightBoss(
            Sprite(health=50, mana=500, armor=0, dmg=0),
            Sprite(health=bosshealth, mana=0, armor=0, dmg=boosDamage),
        ),
    )

    hard_mode = [
        Effect(duration=-1, effect=lambda hero, boss: hero.TakeDmg(1), effect_type=None)
    ]
    print(
        "part 2: ",
        fightBoss(
            Sprite(health=50, mana=500, armor=0, dmg=0),
            Sprite(health=bosshealth, mana=0, armor=0, dmg=boosDamage),
            hard_mode,
        ),
    )


assert 226 == fightBoss(
    Sprite(health=10, mana=250, armor=0, dmg=0),
    Sprite(health=13, mana=0, armor=0, dmg=8),
)
assert 641 == fightBoss(
    Sprite(health=10, mana=250, armor=0, dmg=0),
    Sprite(health=14, mana=0, armor=0, dmg=8),
)

if __name__ == "__main__":
    main()
