import math
import file_reader

def parseItems(items: str) -> dict[str, dict[str, int]]:
    return {name: {"cost" : int(cost), "dmg" : int(dmg), "armor" : int(armor)} for name, cost, dmg, armor in (item.split() for item in items)}

def fightBoss(hero: dict[str, int], boss: dict[str, dict]) -> bool:
    if boss["health"] <= 0: return True
    if hero["health"] <= 0: return False
    dmgToBoss = hero["dmg"] - boss["armor"] if hero["dmg"] > boss["armor"] else 1
    dmgToHero = boss["dmg"] - hero["armor"] if boss["dmg"] > hero["armor"] else 1
    return math.ceil(hero["health"] / dmgToHero) >= math.ceil(boss["health"] / dmgToBoss)

assert(fightBoss({"health" : 8, "dmg" : 5, "armor" : 5}, {"health" : 12, "dmg" : 7, "armor" : 2}))

def appendItem(hero: dict[str, int], item: dict[str, int]) -> int:
    if item is not None:
        hero["dmg"] += item["dmg"]
        hero["armor"] += item["armor"]
        return item["cost"]
    return 0

def buyItemAndFight(
        shop: dict[str, dict[str, tuple[int, int, int]]], 
        hero: dict[str, int], boss: dict[str, int], 
        item: dict[str, int],
        forwardParams: tuple[callable, callable, int]
        ) -> int:
    tempHero = hero.copy()
    cost = appendItem(tempHero, item)
    return cost + shopAndFight(shop, tempHero, boss, *forwardParams)

def shopAndFight(
        shop: dict[str, dict[str, tuple[int, int, int]]], 
        hero: dict[str, int], boss: dict[str, int],
        comparePriceCallback: callable,
        fightBossCallback: callable,
        initValue: int) -> int:
    forwardParams = (comparePriceCallback, fightBossCallback, initValue)
    if "weapons" in shop:
        bestCost = initValue
        for weapon in shop["weapons"]:
            cost = buyItemAndFight({"armor" : shop["armor"], "rings" : shop["rings"]}, hero, boss, shop["weapons"][weapon], forwardParams)
            if comparePriceCallback(cost, bestCost): bestCost = cost
        return bestCost
    elif "armor" in shop:
        bestCost = initValue
        for armor in list(shop["armor"]) + [None]:
            cost = buyItemAndFight({"rings" : shop["rings"]}, hero, boss, shop["armor"].get(armor, None), forwardParams)
            if comparePriceCallback(cost, bestCost): bestCost = cost
        return bestCost
    bestCost = initValue
    for ring1 in list(shop["rings"]) + [None]:
        for ring2 in list(shop["rings"]) + [None]:
            if ring1 == ring2 and ring1 is not None: continue
            tempHero = hero.copy()
            cost = appendItem(tempHero, shop["rings"].get(ring1, None)) + appendItem(tempHero, shop["rings"].get(ring2, None))
            if comparePriceCallback(cost, bestCost) and fightBossCallback(tempHero, boss):
                 bestCost = cost
    return bestCost

def main(weapons, armor, rings):
    weapons = parseItems(weapons)
    armor = parseItems(armor)
    rings = parseItems(rings)

    programInput = file_reader.getInput().splitlines()

    bosshealth = [int(x.removeprefix("Hit Points: ")) for x in programInput if x.startswith("Hit Points")][0]
    boosDamage = [int(x.removeprefix("Damage: ")) for x in programInput if x.startswith("Damage")][0]
    bossArmor = [int(x.removeprefix("Armor: ")) for x in programInput if x.startswith("Armor")][0]

    cheapestWin = shopAndFight({"weapons" : weapons, "armor" : armor, "rings" : rings},
                 {"dmg" : 0, "armor" : 0, "health" : 100},
                 {"dmg" : boosDamage, "armor" : bossArmor, "health" : bosshealth},
                 lambda cost, bestCost: cost < bestCost,
                 lambda hero, boss: fightBoss(hero, boss),
                 math.inf)
    
    expensiveLoss = shopAndFight({"weapons" : weapons, "armor" : armor, "rings" : rings},
                                 {"dmg" : 0, "armor" : 0, "health" : 100},
                                 {"dmg" : boosDamage, "armor" : bossArmor, "health" : bosshealth},
                                 lambda cost, bestCost: cost > bestCost,
                                 lambda hero, boss: not fightBoss(hero, boss),
                                 -math.inf)
    print("part 1: ", cheapestWin)
    print("part 1: ", expensiveLoss)

if __name__ == '__main__':
    weapons = """Dagger        8     4       0
    Shortsword   10     5       0
    Warhammer    25     6       0
    Longsword    40     7       0
    Greataxe     74     8       0""".split('\n')

    armor = """Leather      13     0       1
    Chainmail    31     0       2
    Splintmail   53     0       3
    Bandedmail   75     0       4
    Platemail   102     0       5""".split('\n')

    rings = """Damage+1    25     1       0
    Damage+2    50     2       0
    Damage+3   100     3       0
    Defense+1   20     0       1
    Defense+2   40     0       2
    Defense+3   80     0       3""".split('\n')

    main(weapons, armor, rings)

