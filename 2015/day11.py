import file_reader

def generateNewPassword(currentPassword: str) -> str:
    currentPassword = [x for x in currentPassword]
    index = -1
    while (currentPassword[index] == 'z'):
        currentPassword[index] = 'a'
        index -= 1
    currentPassword[index] = chr(ord(currentPassword[index]) + 1)
    if currentPassword[index] in ['i', 'o', 'l']:
        currentPassword[index] = chr(ord(currentPassword[index]) + 1)
    return "".join(currentPassword)

def findPair(password: str, startPos: int, pairs: list[str]) -> int:
    for i in range(startPos, len(password) - 1):
        if password[i] == password[i+1] and password[i] not in pairs:
            if (i == 0 or password[i-1] != password[i]) and (i+1 == len(password) - 1 or password[i+1] !=  password[i+2]):
                pairs.append(password[i])
                return i
    return -1

def passwordIsValid(password: str) -> bool:
    threeStraight = False
    for i in range(len(password) - 2):
        nextOrd = ord(password[i]) + 1
        if password[i + 1] in ['j', 'p', 'm']:
            nextOrd += 1
        nextnextOrd = nextOrd + 1
        if password[i + 2] in ['j', 'p', 'm']:
            nextnextOrd += 1
        if ord(password[i + 1]) == nextOrd and ord(password[i + 2]) == nextnextOrd:
            threeStraight = True
            break

    if not threeStraight: return False
    for c in ['i', 'o', 'l']:
        if c in password: return False
    pairs = []
    findPair(password, findPair(password, 0, pairs) + 3, pairs)
    return len(pairs) > 1

def getNewPassword(password: str) -> str:
    password = generateNewPassword(password)
    while not passwordIsValid(password):
        password = generateNewPassword(password)
    return password

def main():
   input_str = file_reader.getInput()

   print("part 1: ", getNewPassword(input_str))
   print("part 2: ", getNewPassword(getNewPassword(input_str)))

assert('ab' == generateNewPassword('aa'))
assert('ba' == generateNewPassword('az'))
assert('caaa' == generateNewPassword('bzzz'))

assert(not passwordIsValid('hijklmmn'))
assert(not passwordIsValid('abbceffg'))
assert(not passwordIsValid('abbcegjk'))
assert(not passwordIsValid('ghijklmn'))
assert(passwordIsValid('ghjaabcc'))
assert(not passwordIsValid('ghjaabba'))

assert('ghjaabcc' == getNewPassword('ghijklmn'))

if __name__ == '__main__':
   main()
