import sys


def readFile(path: str) -> str:
    with open(path) as file:
        return file.read()


def getInput() -> str:
    try:
        fileName = sys.argv[1]
    except Exception:
        fileName = input("filename (hint you can provide filename as first arg): ")
    return readFile(fileName)
