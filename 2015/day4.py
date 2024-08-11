import hashlib
import file_reader

def find_with_start(input_str: str, start: str) -> int:
    counter = -1
    while True:
        if str(hashlib.md5((input_str + str(counter := counter + 1)).encode('utf-8')).hexdigest()).startswith(start):
            return counter

def part1(input_str: str) -> int:
    return find_with_start(input_str, '00000')

def part2(input_str: str) -> int:
    return find_with_start(input_str, '000000')

def main():
   input_str: str = file_reader.getInput()

   print("part 1: ", part1(input_str))
   print("part 2: ", part2(input_str))

assert(609043 == part1("abcdef"))

if __name__ == '__main__':
   main()
