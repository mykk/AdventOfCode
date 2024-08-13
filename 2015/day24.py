import heapq
from math import prod
import file_reader


class Combination:
    def __init__(self, container: list[int], remaining_items: list[int]) -> None:
        self.container = container
        self.remaining_items = remaining_items

    def __lt__(self, other) -> bool:
        return sum(self.container) > sum(other.container)


def lhsCombinationBetter(lhs: Combination, rhs: Combination, target_load: int) -> bool:
    if lhs is None:
        return False
    assert sum(lhs.container) == target_load
    lhs_len = len(lhs.container)
    rhs_len = len(rhs.container)
    if sum(rhs.container) != target_load:
        if lhs_len <= rhs_len:
            return True
        return prod(lhs.container) <= prod(rhs.container)

    if lhs_len != rhs_len:
        return lhs_len < rhs_len
    return prod(lhs.container) <= prod(rhs.container)


def balanceLoad(
    items: list[int], containers: int, optimizeForFront: bool
) -> Combination:
    if containers == 1:
        return Combination(container=items.copy(), remaining_items=[])
    target_load: int = sum(items) // containers
    best_combo: Combination = None
    combinations = [
        Combination(container=[], remaining_items=sorted(items, key=lambda x: -x))
    ]

    heapq.heapify(combinations)
    while combinations:
        combination = heapq.heappop(combinations)
        if lhsCombinationBetter(best_combo, combination, target_load):
            continue
        new_remainder = combination.remaining_items.copy()

        for item in combination.remaining_items:
            new_remainder.remove(item)
            new_load = item + sum(combination.container)
            if new_load > target_load:
                continue
            elif new_load + sum(new_remainder) < target_load:
                break

            new_combo = Combination(
                container=combination.container + [item],
                remaining_items=new_remainder.copy(),
            )
            if lhsCombinationBetter(best_combo, new_combo, target_load):
                continue

            if new_load == target_load:
                combination.remaining_items.remove(item)
                if (
                    balanceLoad(combination.remaining_items, containers - 1, False)
                    is not None
                ):
                    if not optimizeForFront:
                        return new_combo
                    best_combo = new_combo
                break
            heapq.heappush(combinations, new_combo)
    return best_combo


def main():
    programInput = file_reader.getInput().splitlines()

    print(
        "part 1: ", prod(balanceLoad([int(x) for x in programInput], 3, True).container)
    )
    print(
        "part 2: ", prod(balanceLoad([int(x) for x in programInput], 4, True).container)
    )


example: list[int] = [x for x in range(1, 6)] + [x for x in range(7, 12)]

assert 99 == prod(balanceLoad(example, 3, True).container)
assert 44 == prod(balanceLoad(example, 4, True).container)

if __name__ == "__main__":
    main()
