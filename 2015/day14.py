import file_reader

def parseInput(input_str: list) -> dict[str, tuple[int, int, int]]:
    parsed = {}
    for line in input_str:
        name, _, _, speed, _, _, duration, *_, rest, _ = line.split()
        parsed[name] = (int(speed), int(duration), int(rest))
    return parsed

def findDistance(speed: int, duration: int, rest: int, time: int) -> int:
    loops = int(time / (duration + rest))
    remainder = time - loops * (duration + rest)
    if remainder > duration: remainder = duration
    return loops * speed * duration + remainder * speed

def findFurthestTravelDistance(reindeers: dict[str, tuple[int, int, int]], travelTime: int) -> int:
    return max([findDistance(*reindeer, travelTime) for reindeer in reindeers.values()])

def findLeadingDeer(reindeers: dict[str, tuple[int, int, int]], travelTime: int) -> list[str]:
    times = {}
    for reindeer in reindeers:
        times.setdefault(findDistance(*reindeers[reindeer], travelTime), []).append(reindeer)
    return list(times[max(times)])

def findBestPoints(reindeers: dict[str, tuple[int, int, int]], travelTime: int) -> int:
    leaderboard = {reindeer : 0 for reindeer in reindeers}
    for current in range(1, travelTime + 1):
        leaders = findLeadingDeer(reindeers, current)
        for leader in leaders:
           leaderboard[leader] += 1
    return max(leaderboard.values())


def main():
   input_str = file_reader.getInput().splitlines()

   print("part 1: ", findFurthestTravelDistance(parseInput(input_str), 2503))
   print("part 2: ", findBestPoints(parseInput(input_str), 2503))

example = """Comet can fly 14 km/s for 10 seconds, but then must rest for 127 seconds.
Dancer can fly 16 km/s for 11 seconds, but then must rest for 162 seconds.""".split('\n')

parsedExample = parseInput(example)
assert(2 == len(parsedExample))
assert((14, 10, 127) == parsedExample["Comet"])
assert((16, 11, 162) == parsedExample["Dancer"])

assert(1120 == findDistance(*parsedExample["Comet"], 1000))
assert(1056 == findDistance(*parsedExample["Dancer"], 1000))
assert(1120 == findFurthestTravelDistance(parsedExample, 1000))
assert(689 == findBestPoints(parsedExample, 1000))

if __name__ == '__main__':
   main()
