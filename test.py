from SortedSet import SortedSet

out_file = None


def print_all(com: str, ss: SortedSet):
    print("all", com, *ss, file=out_file)


def common_print(com: str, arg):
    if isinstance(arg, bool) and arg:
        print(com, "true", file=out_file)
        return
    if isinstance(arg, bool) and not arg:
        print(com, "false", file=out_file)
        return
    if arg is None:
        print(com, "None", file=out_file)
        return
    print(com, arg, file=out_file)


def use(ss: SortedSet, command: str, n: int = 0):
    if command == "add":
        ss.add(n)
        print_all(command, ss)
    if command == "discard":
        ss.discard(n)
        print_all(command, ss)
    if command == "get":
        try:
            common_print(command, ss[n])
        except IndexError:
            common_print(command, None)
    if command == "pop":
        try:
            common_print(command, ss.pop(n))
        except IndexError:
            common_print(command, None)
        print_all(command, ss)
    if command == "index":
        common_print(command, ss.index(n))
    if command == "index_right":
        common_print(command, ss.index_right(n))
    if command == "lt":
        common_print(command, ss.lt(n))
    if command == "le":
        common_print(command, ss.le(n))
    if command == "gt":
        common_print(command, ss.gt(n))
    if command == "ge":
        common_print(command, ss.ge(n))
    if command == "contains":
        common_print(command, n in ss)
    if command == "len":
        common_print(command, len(ss))


def run():
    with open("testdata/input.txt") as f:
        lines = f.readlines()
    ss = SortedSet()
    for line in lines + lines:
        command, *args = line.split()
        use(ss, command, *map(int, args))


if __name__ == "__main__":
    with open("testdata/output_py.txt", "w") as f:
        out_file = f
        run()
