#!/usr/bin/python3

import json
import sys

server = "server.mime"

if __name__ == "__main__":
    with open("src.json", "r") as f:
        data = json.loads(f.read())
        flipped = {}
        for m, l in data.items():
            for i in l:
                flipped[i.replace("*", ".*", 1) + "$"] = m
        with open("out.sql", "w") as r:
            r.write(f"TRUNCATE {server};\n")
            r.writelines([f"INSERT INTO {server} (extension, mimetype, \"index\") VALUES ('{ext}', '{typ}', {index});\n" for index, (ext, typ) in enumerate(flipped.items())])
            print("generated SQL")


# json joinked from https://github.com/broofa/mime