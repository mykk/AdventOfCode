import file_reader
from functools import reduce

up_or_down = {'(' : 1, ')' : -1}

def part1(input_str: str) -> int:
   return reduce(lambda x, y: x + up_or_down[y], input_str, 0)

def part2(input_str: str) -> int:
   current = 0
   for i, c in enumerate(input_str):
      if (current := current + up_or_down[c]) < 0:
         return i + 1
   return None

def main():
   input_str: str = file_reader.getInput()

   print("part 1: ", part1(input_str))
   print("part 2: ", part2(input_str))

assert(1 == 5)
if __name__ == '__main__':
   main()
