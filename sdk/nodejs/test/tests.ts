// Copyright 2016-2018, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import * as bcryptjs from "bcryptjs";
import { assert } from "chai";
import * as os from "os";
import * as path from "path";
import { gunzipSync } from "zlib";
import * as sut from "../src/index";

const validator = require("validator");

describe("abs", () => {
    it("works", () => {
        assert.equal(sut.abs(1), 1);
        assert.equal(sut.abs(-1), 1);
        assert.equal(sut.abs(-3.14), 3.14);
        assert.equal(sut.abs(42.001), 42.001);
    });
});

describe("basename", () => {
    it("works for unix paths", () => {
        assert.equal(sut.basename("/foo/bar/baz"), "baz");
    });

    it("works for Windows paths", () => {
        assert.equal(sut.basename("C:\\Windows\\system32\\kernel32.dll"), "kernel32.dll");
    });
});

describe("base64decode", () => {
    it("works with valid base64-encoded input", () => {
        assert.equal(sut.base64decode("YWJjMTIzIT8kKiYoKSctPUB+"), "abc123!?$*&()'-=@~");
    });

    it("throws with invalid base64-encoded input", () => {
        assert.throws(() => sut.base64decode("this-is-an-invalid-base64-data"));
    });
});

describe("base64encode", () => {
    it("works", () => {
        assert.equal(sut.base64encode("abc123!?$*&()'-=@~"), "YWJjMTIzIT8kKiYoKSctPUB+");
    });
});

describe("base64gzip", () => {
    it("works", () => {
        const zipped = sut.base64gzip("test");

        assert.equal(gunzipSync(Buffer.from(zipped, "base64")).toString("utf-8"), "test");
    });
});

describe("base64sha256", () => {
    it("works", () => {
        assert.equal(sut.base64sha256("test"), "n4bQgYhMfWWaL+qgxVrQFaO/TxsrC4Is0V1sFbDwCgg=");
    });

    it("is different to base64encode(sha256(input))", () => {
        assert.equal(sut.base64encode(sut.sha256("test")),
            "OWY4NmQwODE4ODRjN2Q2NTlhMmZlYWEwYzU1YWQwMTVhM2JmNGYxYjJiMGI4MjJjZDE1ZDZjMTViMGYwMGEwOA==");
    });
});

describe("base64sha512", () => {
    it("works", () => {
        assert.equal(sut.base64sha512("test"),
            "7iaw3Ur350mqGo7jwQrpkj9hiYB3Lkc/iBml1JQODbJ6wYX4oOHV+E+IvIh/1nsUNzLDBMxfqa2Ob1f1ACio/w==");
    });

    it("is different to base64encode(sha512(input))", () => {
        assert.equal(sut.base64encode(sut.sha512("test")),
            "ZWUyNmIwZGQ0YWY3ZTc0OWFhMWE4ZWUzYzEwYWU5OTIzZjYxODk4MDc3MmU0NzNmODgxOWE1ZDQ5NDBlMGRiMjdhYzE4NWY4YTBlM" +
            "WQ1Zjg0Zjg4YmM4ODdmZDY3YjE0MzczMmMzMDRjYzVmYTlhZDhlNmY1N2Y1MDAyOGE4ZmY=");
    });
});

describe("bcrypt", () => {
    it("works with default cost", () => {
        const hash = sut.bcrypt("test");
        assert.isTrue(bcryptjs.compareSync("test", hash));
    }).timeout(8000);

    it("works with overridden cost", () => {
        const hash = sut.bcrypt("test", 5);
        assert.isTrue(bcryptjs.compareSync("test", hash));
    }).timeout(5000);
});

describe("ceil", () => {
    it("works for positive numbers", () => {
        assert.equal(sut.ceil(1.2), 2);
    });

    it("works for negative numbers", () => {
        assert.equal(sut.ceil(-1.2), -1);
    });
});

describe("chomp", () => {
    it("does not affect strings with no trailing newlines", () => {
        assert.equal(sut.chomp("hello\ncruel\nworld"), "hello\ncruel\nworld");
    });

    it("does not affect leading newlines", () => {
        assert.equal(sut.chomp("\n\nhello\ncruel\nworld"), "\n\nhello\ncruel\nworld");
    });

    it("works for strings with one UNIX trailing newline", () => {
        assert.equal(sut.chomp("hello\ncruel\nworld\n"), "hello\ncruel\nworld");
    });

    it("works for strings with multiple UNIX trailing newlines", () => {
        assert.equal(sut.chomp("hello\ncruel\nworld\n\n\n"), "hello\ncruel\nworld");
    });

    it("works for strings with one Windows trailing newline", () => {
        assert.equal(sut.chomp("hello\r\ncruel\r\nworld\r\n"), "hello\r\ncruel\r\nworld");
    });

    it("works for strings with multiple Windows trailing newline", () => {
        assert.equal(sut.chomp("hello\r\ncruel\r\nworld\r\n\r\n"), "hello\r\ncruel\r\nworld");
    });
});

describe("chunklist", () => {
    it("chunks single items", () => {
        assert.deepEqual(sut.chunklist(["id1", "id2", "id3"], 1), [["id1"], ["id2"], ["id3"]]);
    });

    it("gives chunks of equal size with remainder in last chunk", () => {
        assert.deepEqual(sut.chunklist(["id1", "id2", "id3", "id4", "id5"], 2),
            [["id1", "id2"], ["id3", "id4"], ["id5"]]);
    });

    it("works for other types", () => {
        assert.deepEqual(sut.chunklist([1, 2, 3, 4, 5], 2), [[1, 2], [3, 4], [5]]);
    });
});

describe("coalesce", () => {
    it("passes Terraform tests", () => {
        assert.equal(sut.coalesce("first", "second", "third"), "first");
        assert.equal(sut.coalesce("", "second", "third"), "second");
        assert.equal(sut.coalesce("", "", ""), "");
    });
});

describe("coalescelist", () => {
    it("passes Terraform tests", () => {
        assert.deepEqual(sut.coalescelist(["first"], ["second"], ["third"]), ["first"]);
        assert.deepEqual(sut.coalescelist([], ["second"], ["third"]), ["second"]);
        assert.deepEqual(sut.coalescelist([], [], []), []);
    });
});

describe("compact", () => {
    it("passes Terraform tests", () => {
        assert.deepEqual(sut.compact(["", "", "first", "second"]), ["first", "second"]);
        assert.deepEqual(sut.compact(["", "", "first", "second", ""]), ["first", "second"]);
        assert.deepEqual(sut.compact([""]), []);
    });
});

describe("concat", () => {
    it("passes Terraform tests", () => {
        assert.deepEqual(sut.concat(["", "first", ""]),
            ["", "first", ""]);
        assert.deepEqual(sut.concat(["", "first", ""], ["second", "third"]),
            ["", "first", "", "second", "third"]);
        assert.deepEqual(sut.concat(["first"], ["second"], ["third"]),
            ["first", "second", "third"]);
        assert.deepEqual(sut.concat([{key: "value"}], [{key2: "value2"}]),
            [{key: "value"}, {key2: "value2"}]);
    });
});

describe("contains", () => {
    it("passes Terraform tests", () => {
        assert.equal(sut.contains([1, 2, 3, 4, 5], 6), false);
        assert.equal(sut.contains([1, 2, 3, 4, 5], 1), true);
        assert.equal(sut.contains(["first", "", "second", "third"], ""), true);
        assert.equal(sut.contains(["first", "", "second", "third"], "second"), true);
        assert.equal(sut.contains(["first", "", "second", "third"], "fourth"), false);
    });
});

describe("dirname", () => {
    it("works for unix paths", () => {
        assert.equal(sut.dirname("/foo/bar/baz"), "/foo/bar");
    });

    it("works for Windows paths", () => {
        assert.equal(sut.dirname("C:\\Windows\\system32\\kernel32.dll"), "C:\\Windows\\system32");
    });
});

describe("distinct", () => {
    it("passes Terraform tests", () => {
        assert.deepEqual(sut.distinct(["a", "a", "b"]), ["a", "b"]);
        assert.deepEqual(sut.distinct(["a", "b", "c"]), ["a", "b", "c"]);
    });
});

describe("element", () => {
    it("passes Terraform tests", () => {
        assert.equal(sut.element(["a", "b"], 0), "a");
        assert.equal(sut.element(["a", "b"], 1), "b");
        assert.equal(sut.element(["a", "b"], 2), "a");

        assert.throws(() => sut.element([], 1));
        assert.throws(() => sut.element(["a", "b"], -1));
    });
});

describe("file", () => {
    it("reads a file", () => {
        const testFilePath = path.join(__dirname, "testdata", "test-file.txt");
        assert.equal(sut.file(testFilePath), "Hello World\n");
    });
});

describe("floor", () => {
    it("passes Terraform tests", () => {
        assert.equal(sut.floor(-1.3), -2);
        assert.equal(sut.floor(1.7), 1);
    });
});

describe("format", () => {
    it("passes Terraform tests", () => {
        assert.equal(sut.format("hello"), "hello");
        assert.equal(sut.format("hello %s", "world"), "hello world");
        assert.equal(sut.format("hello %d", 42), "hello 42");
        assert.equal(sut.format("hello %05d", 42), "hello 00042");
        assert.equal(sut.format("hello %05d", 12345), "hello 12345");
    });
});

describe("formatlist", () => {
    it("passes Terraform tests", () => {
        assert.throws(() => sut.formatlist("hello"));
        assert.throws(() => sut.formatlist("hello %s", ["world"], ["world2", "world3"]));

        assert.deepEqual(sut.formatlist("<%s>", ["A", "B"]), ["<A>", "<B>"]);
        assert.deepEqual(sut.formatlist("%s=%d", ["A", "B", "C"], [1, 2, 3]), ["A=1", "B=2", "C=3"]);

        assert.deepEqual(sut.formatlist("%s", [], []), []);
    });
});

describe("indent", () => {
    it("passes Terraform tests", () => {
        assert.equal(sut.indent(4, "Fleas:\nAdam\nHad'em\n\nE.E. Cummings"),
            "Fleas:\n    Adam\n    Had'em\n    \n    E.E. Cummings");
        assert.equal(sut.indent(4, "oneliner"), "oneliner");
        assert.equal(sut.indent(8, "#!/usr/bin/env bash\ndate\npwd\n"),
            "#!/usr/bin/env bash\n        date\n        pwd\n        ");
    });

    it("works with Windows strings", () => {
        assert.equal(sut.indent(8, "#!/usr/bin/env bash\r\ndate\r\npwd\r\n"),
            "#!/usr/bin/env bash\r\n        date\r\n        pwd\r\n        ");
    });
});

describe("index", () => {
    it("passes Terraform tests", () => {
        assert.throws(() => sut.index(["notfoo", "stillnotfoo", "bar"], "foo"));
        assert.equal(sut.index(["foo"], "foo"), 0);
        assert.equal(sut.index(["foo", "spam", "bar", "eggs"], "bar"), 2);
    });
});

describe("join", () => {
    it("passes Terraform tests", () => {
        assert.equal(sut.join(", ", ["hello", "world"]), "hello, world");
        assert.equal(sut.join(".", ["10", "11", "12", "13"]), "10.11.12.13");
    });

    it("works with multiple lists", () => {
        assert.equal(sut.join(", ", ["hello", "world"], ["goodbye", "world"]),
            "hello, world, goodbye, world");
    });
});

describe("jsonencode", () => {
    it("works for a string", () => {
        assert.equal(sut.jsonencode("test"), "\"test\"");
    });

    it("works for an array", () => {
        assert.equal(sut.jsonencode(["hello", "world"]), "[\"hello\",\"world\"]");
    });

    it("works with objects", () => {
        assert.equal(sut.jsonencode({key: "value", key2: ["array1", "array2"]}),
            "{\"key\":\"value\",\"key2\":[\"array1\",\"array2\"]}");
    });
});

describe("keys", () => {
    it("passes Terraform tests", () => {
        assert.deepEqual(sut.keys({bar: "baz", qux: "quack"}), ["bar", "qux"]);
    });
});

describe("length", () => {
    it("returns number of characters in a string", () => {
        assert.equal(sut.length("hello"), 5);
    });
    it("returns the number of elements in an array", () => {
        assert.equal(sut.length(["a", "b", "c"]), 3);
    });
    it("returns the number of keys in an object", () => {
        assert.equal(sut.length({foo: "bar", baz: "qux"}), 2);
    });
});

describe("list", () => {
    it("passes Terraform tests", () => {
        assert.deepEqual(sut.list(), []);
        assert.deepEqual(sut.list("hello"), ["hello"]);
        assert.deepEqual(sut.list("hello", ...["hello", "world"]), ["hello", "hello", "world"]);
        assert.deepEqual(sut.list({key: "value1"}, {key: "value2"}), [{key: "value1"}, {key: "value2"}]);
    });
});

describe("log", () => {
    it("passes Terraform tests", () => {
        assert.equal(sut.log(1, 10), 0);
        assert.equal(sut.log(10, 10), 1);
        assert.equal(sut.log(0, 10), -Infinity);
        assert.equal(sut.log(10, 0), 0);
    });
});

describe("lookup", () => {
    it("passes Terraform tests", () => {
        const exampleMap = {
            bar: "baz",
        };

        assert.equal(sut.lookup(exampleMap, "bar"), "baz");
        assert.throws(() => sut.lookup(exampleMap, "baz"));
        assert.equal(sut.lookup(exampleMap, "baz", "foo"), "foo");
    });
});

describe("lower", () => {
    it("passes Terraform tests", () => {
        assert.equal(sut.lower("HELLO"), "hello");
        assert.equal(sut.lower("hello"), "hello");
        assert.equal(sut.lower(""), "");
    });
});

describe("map", () => {
    it("passes Terraform tests", () => {
        assert.deepEqual(sut.map("key1", "value1", "key2", "value2"), {key1: "value1", key2: "value2"});
        assert.deepEqual(sut.map("key1", 1, "key2", 2), {key1: 1, key2: 2});
    });
});

describe("max", () => {
    it("passes Terraform tests", () => {
        assert.equal(sut.max(-1, 0, 1), 1);
        assert.equal(sut.max(-1, -2), -1);
        assert.equal(sut.max(-1, -2, 42), 42);
    });
});

describe("merge", () => {
    it("passes Terraform tests", () => {
        assert.deepEqual(sut.merge({a: "b"}, {c: "d"}), {a: "b", c: "d"});
        assert.deepEqual(sut.merge({a: "b"}, {a: "d"}), {a: "d"});
        assert.deepEqual(sut.merge({a: 1}, {a: "hello"}), {a: "hello"});
    });
});

describe("min", () => {
    it("passes Terraform tests", () => {
        assert.equal(sut.min(-1, 0, 1), -1);
        assert.equal(sut.min(-1, -2), -2);
        assert.equal(sut.min(-1, -2, 42), -2);
    });
});

describe("md5", () => {
    it("passes Terraform tests", () => {
        assert.equal(sut.md5("tada"), "ce47d07243bb6eaf5e1322c81baf9bbf");
        assert.equal(sut.md5(" tada "), "aadf191a583e53062de2d02c008141c4");
        assert.equal(sut.md5(""), "d41d8cd98f00b204e9800998ecf8427e");
    });
});

describe("pathexpand", () => {
    it("passes Terraform tests", () => {
        const homedir = os.homedir();

        assert.equal(sut.pathexpand("~/hello.txt"), `${homedir}/hello.txt`);
        assert.equal(sut.pathexpand("~/a/test/file"), `${homedir}/a/test/file`);
        assert.equal(sut.pathexpand("/a/test/file"), `/a/test/file`);
        assert.equal(sut.pathexpand("/"), `/`);
    });
});

describe("pow", () => {
    it("passes Terraform tests", () => {
        assert.equal(sut.pow(3, 2), 9);
        assert.equal(sut.pow(4, 0), 1);
    });
});

describe("replace", () => {
    it("passes Terraform tests", () => {
        assert.equal(sut.replace("hello", "hel", "bel"), "bello");
        assert.equal(sut.replace("hello", "nope", "bel"), "hello");
        assert.equal(sut.replace("hello", /l/g, "L"), "heLLo");
        assert.equal(sut.replace("hello", /(l)/, "$1"), "hello");
        assert.equal(sut.replace("helo", /(l)/, "$1$1"), "hello");
    });
});

describe("sha1", () => {
    it("passes Terraform tests", () => {
        assert.equal(sut.sha1("test"), "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3");
    });
});

describe("sha256", () => {
    it("passes Terraform tests", () => {
        assert.equal(sut.sha256("test"), "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08");
    });
});

describe("sha512", () => {
    it("passes Terraform tests", () => {
        assert.equal(sut.sha512("test"),
            "ee26b0dd4af7e749aa1a8ee3c10ae9923f618980772e473f8819a5d4940e0db27ac185f8a0e1d5f84f88bc887fd67b1437" +
            "32c304cc5fa9ad8e6f57f50028a8ff");
    });
});

describe("signum", () => {
    it("passes Terraform tests", () => {
        assert.equal(sut.signum(-10), -1);
        assert.equal(sut.signum(0), 0);
        assert.equal(sut.signum(10), 1);
    });
});

describe("slice", () => {
    it("passes Terraform tests", () => {
        assert.deepEqual(sut.slice(["a", "b", "c"], 1, 1), []);
        assert.deepEqual(sut.slice(["a", "b", "c"], 1, 2), ["b"]);
        assert.deepEqual(sut.slice(["a", "b", "c"], 0, 2), ["a", "b"]);

        assert.throws(() => sut.slice(["a"], -1, 0));
        assert.throws(() => sut.slice(["a", "b", "c"], 2, 1));
        assert.throws(() => sut.slice(["a", "b", "c"], 1, 4));
    });
});

describe("sort", () => {
    it("passes Terraform tests", () => {
        assert.deepEqual(sut.sort(["c", "a", "b"]), ["a", "b", "c"]);
    });
});

describe("split", () => {
    it("passes Terraform tests", () => {
        assert.deepEqual(sut.split(",", "a,,b"), ["a", "b"]);
        assert.deepEqual(sut.split(",", "a,b,"), ["a", "b"]);
        assert.deepEqual(sut.split(",", ""), []);
    });
});

describe("substr", () => {
    it("passes Terraform tests", () => {
        assert.equal(sut.substr("foobar", 0, 0), "");
        assert.equal(sut.substr("foobar", 0, -1), "foobar");
        assert.equal(sut.substr("foobar", 0, 3), "foo");
        assert.equal(sut.substr("foobar", 3, 3), "bar");
        assert.equal(sut.substr("foobar", -3, 3), "bar");
        assert.equal(sut.substr("", 0, 0), "");
        assert.throws(() => sut.substr("foo", -4, -1));
        assert.throws(() => sut.substr("", 1, 0));
        assert.throws(() => sut.substr("", 0, 1));
        assert.throws(() => sut.substr("", 0, -2));
    });
});

describe("timestamp", () => {
    it("passes Terraform tests", () => {
        assert.isTrue(validator.isRFC3339(sut.timestamp()));
    });
});

describe("title", () => {
    it("passes Terraform tests", () => {
        assert.equal(sut.title(""), "");
        assert.equal(sut.title("hello"), "Hello");
        assert.equal(sut.title("hello world"), "Hello World");
    });
});

describe("transpose", () => {
    it("passes Terraform tests", () => {
        const inputMap = {
            key1: ["a", "b"],
            key2: ["a", "b", "c"],
            key3: ["c"],
            key4: [],
        };
        const expectedMap = {
            "a": ["key1", "key2"],
            "b": ["key1", "key2"],
            "c": ["key2", "key3"],
        };

        assert.deepEqual(sut.transpose(inputMap), expectedMap);
    });
});

describe("trimspace", () => {
    it("passes Terraform tests", () => {
        assert.equal(sut.trimspace(" test "), "test");
    });
});

describe("upper", () => {
    it("passes Terraform tests", () => {
        assert.equal(sut.upper("HELLO"), "HELLO");
        assert.equal(sut.upper("hello"), "HELLO");
    });
});

describe("urlencode", () => {
    it("passes Terraform tests", () => {
        assert.equal(sut.urlencode("abc123-_"), "abc123-_");
        assert.equal(sut.urlencode("foo:bar@localhost?foo=bar&bar=baz"),
            "foo%3Abar%40localhost%3Ffoo%3Dbar%26bar%3Dbaz");
        assert.equal(sut.urlencode("mailto:email?subject=this+is+my+subject"),
            "mailto%3Aemail%3Fsubject%3Dthis%2Bis%2Bmy%2Bsubject");
        assert.equal(sut.urlencode("foo/bar"), "foo%2Fbar");
    });
});

describe("uuid", () => {
    it("passes Terraform tests", () => {
        const validUUID = /^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i;

        const generated1 = sut.uuid();
        const generated2 = sut.uuid();

        assert.match(generated1, validUUID);
        assert.match(generated2, validUUID);
        assert.notEqual(generated1, generated2);
    });
});

describe("values", () => {
    it("passes Terraform tests", () => {
        const map = {
            key: "value",
            key2: "value2",
            key3: "value3",
        };

        assert.deepEqual(sut.values(map), ["value", "value2", "value3"]);
    });
});

describe("zipmap", () => {
    it("passes Terraform tests", () => {
        const keyList = ["Hello", "World"];
        const valueList = ["foo", "bar"];

        assert.deepEqual(sut.zipmap(keyList, valueList), {"Hello": "foo", "World": "bar"});

        assert.throws(() => sut.zipmap(keyList, ["foo"]));
    });
});
