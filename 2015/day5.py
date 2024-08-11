import file_reader

VOWELS = ["a", "e", "i", "o", "u"]
BAD_SUBSTRINGS = ["ab", "cd", "pq", "xy"]

def niceString(line:str) -> bool:
   if sum([1 for c in line if c in VOWELS]) < 3:
      return False
   double = False
   for i in range(len(line) - 1):
      if line[i] == line[i+1]:
         double = True
         break
   if not double: return False
   for bad in BAD_SUBSTRINGS:
      if line.find(bad) != -1: return False
   return True

def niceString2(line:str) -> bool:
   goodPair = False
   for i in range(len(line) - 1):
      pair = line[i] + line[i+1]
      if line.find(pair, i + 2) != -1:
         goodPair = True
   if not goodPair: return False
   for i in range(len(line) - 2):
      if line[i] == line[i+2]:
         return True
   return False

def main():
   input_str: str = file_reader.getInput()

   print("part 1: ", niceString(input_str))
   print("part 2: ", niceString2(input_str))

assert(niceString2("qjhvhtzxzqqjkmpb"))
assert(niceString2("xxyxx"))
assert(not niceString2("uurcxstgmygtbstg"))
assert(not niceString2("ieodomkazucvgmuy"))

if __name__ == '__main__':
   main()
