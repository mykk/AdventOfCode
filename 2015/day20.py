from sympy import divisors
from itertools import count
import file_reader

def getPresentCount(house: int, multiplier: int, limit: int) -> int:
    visitors = divisors(house)
    return sum(i * multiplier for i in visitors if limit == 0 or house / i <= limit)

def findLowestHouseNumber(presents: int, multiplier: int, limit: int) -> int:
    return next(x for x in count(start=1) if getPresentCount(house=x, multiplier=multiplier, limit=limit) >= presents)

def main():
    programInput = file_reader.getInput().splitlines()
    presentCount = int(programInput[0])
    print("part 1: ", findLowestHouseNumber(presentCount, 10, 0))
    print("part 1: ", findLowestHouseNumber(presentCount, 11, 50))

assert(130 == getPresentCount(9, 10, 0))
assert(150 == getPresentCount(8, 10, 0))
assert(70 == getPresentCount(4, 10, 0))
assert(310 == getPresentCount(16, 10, 0))

assert(4 == findLowestHouseNumber(70, 10, 0))
assert(8 == findLowestHouseNumber(150, 10, 0))
assert(8 == findLowestHouseNumber(130, 10, 0))
assert(12 == findLowestHouseNumber(240, 10, 0))

if __name__ == '__main__':
   main()