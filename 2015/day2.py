import file_reader
from functools import reduce

def part1(dimensions: list[tuple[int, int, int]]) -> int:
   calc = lambda w, h, l: 2*l*w + 2*w*h + 2*h*l + min([w*h, h*l, l*w])
   return reduce(lambda x, dims: x + calc(*dims), dimensions, 0)

def part2(dimensions: list[tuple[int, int, int]]) -> int:
   calc = lambda w, h, l: sum(sorted([w, h, l])[0:2]) * 2 + w*h*l
   return reduce(lambda x, dims: x + calc(*dims), dimensions, 0)

def main():
   input_str: str = file_reader.getInput()
   dimensions = [(int(w), int(h), int(l)) for w, h, l in [x.split('x') for x in input_str.splitlines()]]

   print("part 1: ", part1(dimensions))
   print("part 2: ", part2(dimensions))

if __name__ == '__main__':
   main()