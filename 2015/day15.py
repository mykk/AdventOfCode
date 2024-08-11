import file_reader

CALORIES = 'calories'
from functools import reduce
def parseInput(programInput: str) -> dict[str, dict[str, int]]:
    ingredients = {}
    for line in programInput:
        ingredient, flavors = [x.strip() for x in line.split(':')]
        flavors = [x.strip() for x in flavors.split(',')]
        ingredients[ingredient] = { x.split()[0] : int(x.split()[1]) for x in flavors}
    return ingredients

def countCalories(ingredients: dict[str, dict[str, int]], current: dict[str, int]) -> int:
    return sum(ingredients[ingredient][CALORIES] * amount for ingredient, amount in current.items())

def calculateCurrentFlavor(ingredients: dict[str, dict[str, int]], current: dict[str, int]) -> int:
    flavors = {}
    for ingredient, amount in current.items():
        for flavor, score in ingredients[ingredient].items():
            flavors.setdefault(flavor, 0)
            flavors[flavor] += score * amount
    return reduce(lambda x, y: x * y, (flavors[flavor] if flavors[flavor] > 0 else 0 for flavor in flavors if flavor != CALORIES), 1)

def findBestCombo(ingredients: dict[str, dict[str, int]], current_recipe: dict[str, int], badRecipe: callable) -> int:
    current_recipe = current_recipe.copy()
    current_ingredient = next((x for x in ingredients if x not in current_recipe))

    current_recipe[current_ingredient] = 0
    if len(current_recipe) == len(ingredients):
        current_recipe[current_ingredient] = 100 - sum(current_recipe.values())
        if badRecipe(ingredients, current_recipe):
            return 0
        return calculateCurrentFlavor(ingredients, current_recipe)

    best = 0
    for i in range(1, 101 - sum(current_recipe.values())):
        current_recipe[current_ingredient] = i
        flavor = findBestCombo(ingredients, current_recipe, badRecipe)
        if flavor > best or best == 0:
            best = flavor
    return best

def main():
   ingredients = parseInput(file_reader.getInput().splitlines())
   print("part 1: ", findBestCombo(ingredients, {}, lambda x, y: False))
   print("part 2: ", findBestCombo(ingredients, {}, lambda x, y: countCalories(x,y) != 500))

example = """Butterscotch: capacity -1, durability -2, flavor 6, texture 3, calories 8
Cinnamon: capacity 2, durability 3, flavor -2, texture -1, calories 3""".split('\n')
exampleParsed = parseInput(example)

assert(62842880 == findBestCombo(exampleParsed, {}, lambda x, y: False))
assert(57600000 == findBestCombo(exampleParsed, {}, lambda x, y: countCalories(x,y) != 500))

if __name__ == '__main__':
   main()
