from SortedSet import SortedSet
from dataclasses import dataclass
from typing import Any, Iterable, Optional
import json


@dataclass
class Output:
    method: str
    contents: list
    arg: Optional[Any] = None
    result: Optional[Any] = None


@dataclass
class Operation:
    method: str
    arg: Optional[Any] = None

    def run(self, ss: SortedSet | None = None) -> Output:
        if self.method == "init":
            if self.arg is None:
                ss = SortedSet()
            else:
                if isinstance(self.arg, Iterable):
                    ss = SortedSet(self.arg)
                else:
                    raise ValueError("arg is not iterable")
            return Output(self.method, [s for s in ss], self.arg, ss)

        if ss is None:
            raise ValueError

        try:
            result = getattr(ss, self.method)(self.arg)
        except IndexError:
            return Output(self.method, [s for s in ss], self.arg, "index error")
        return Output(self.method, [s for s in ss], self.arg, result)


class OperationEncoder(json.JSONEncoder):
    def default(self, o):
        if isinstance(o, Operation):
            return {"method": o.method, "arg": o.arg}
        if isinstance(o, Output):
            return {
                "method": o.method,
                "arg": o.arg,
                "result": o.result,
                "contents": o.contents,
            }
        if isinstance(o, SortedSet):
            return o.a
        return super().default(o)


def main() -> None:
    with open("testdata/input.json") as f:
        row_data = json.load(f)
        test_cases: dict[str, list[Operation]] = dict()
        for name, row_operations in row_data.items():
            operations = [Operation(**op) for op in row_operations]
            test_cases[name] = operations

    all_outputs: dict[str, list[Output]] = dict()
    for name, operations in test_cases.items():
        ss: SortedSet | None = None
        init_count = 0
        results: list[Output] = []
        for op in operations:
            if op.method == "init":
                if init_count > 0:
                    raise ValueError("multiple init")
                result = op.run()
                ss = result.result
                init_count += 1
            else:
                result = op.run(ss)
            results.append(result)
        all_outputs[name] = results

    with open("testdata/output_py.json", "w") as f:
        json.dump(all_outputs, f, cls=OperationEncoder, indent=2)


if __name__ == "__main__":
    main()
