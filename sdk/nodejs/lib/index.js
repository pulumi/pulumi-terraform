"use strict";
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
Object.defineProperty(exports, "__esModule", { value: true });
const bcryptjs = require("bcryptjs");
const crypto = require("crypto");
const fs = require("fs");
const os = require("os");
const Path = require("path");
const sprintf_js_1 = require("sprintf-js");
const uuid_1 = require("uuid");
const zlib_1 = require("zlib");
const zlib = require("zlib");
const titleCase = require("title-case");
/**
 * abs returns the absolute value of a given number. See also the signum function.
 *
 * Examples:
 *     - `abs(1)` returns `1`
 *     - `abs(-1)` returns `1`
 *     - `abs(-3.14)` returns `3.14`
 *
 * @param input number for which to return the absolute value
 */
function abs(input) {
    return Math.abs(input);
}
exports.abs = abs;
/**
 * basename returns the last element of a path.
 *
 * @param path the full path
 */
function basename(path) {
    return Path.win32.basename(path);
}
exports.basename = basename;
/**
 * base64decode decodes base64-encoded data and returns the utf-8-encoded string. An
 * error is thrown if the argument is not valid base64 data.
 *
 * @param encoded valid base64-encoded string
 */
function base64decode(encoded) {
    const decoded = Buffer.from(encoded, "base64").toString("utf-8");
    if (decoded.length !== ((encoded.length / 4) * 3)) {
        throw new Error("failed to decode base64 data");
    }
    return decoded;
}
exports.base64decode = base64decode;
/**
 * base64encode returns a base64-encoded version of the given string.
 *
 * @param unencoded the string to encode
 */
function base64encode(unencoded) {
    return Buffer.from(unencoded).toString("base64");
}
exports.base64encode = base64encode;
/**
 * base64gzip compresses the given string with gzip and then returns the resulting data as
 * a base64-encoded string.
 *
 * @param unencoded the string to compress and encode
 */
function base64gzip(unencoded) {
    const compressed = zlib_1.gzipSync(unencoded, {
        level: zlib.constants.Z_DEFAULT_COMPRESSION,
    });
    return compressed.toString("base64");
}
exports.base64gzip = base64gzip;
function base64sha256(input) {
    const hash = crypto.createHash("sha256");
    hash.update(input);
    return hash.digest().toString("base64");
}
exports.base64sha256 = base64sha256;
function base64sha512(input) {
    const hash = crypto.createHash("sha512");
    hash.update(input);
    return hash.digest().toString("base64");
}
exports.base64sha512 = base64sha512;
/**
 * bcrypt returns the Blowfish encrypted hash of the string at the given cost. A default
 * cost of 10 will be used if not provided.
 *
 * @param password password to encrypt
 * @param cost defaults to 10 if unset
 */
function bcrypt(password, cost) {
    return bcryptjs.hashSync(password, cost || 10);
}
exports.bcrypt = bcrypt;
/**
 * ceil returns the smallest integer value greater than or equal to the argument.
 *
 * @param input number of which to find the ceiling
 */
function ceil(input) {
    return Math.ceil(input);
}
exports.ceil = ceil;
function chomp(input) {
    return input.replace(/[\r\n]+$/, "");
}
exports.chomp = chomp;
/**
 * chunklist returns a list items chunked by size. For example:
 *
 * - chunklist(["id1", "id2", "id3"], 1) returns [["id1"], ["id2"], ["id3"]]
 * - chunklist(["id1", "id2", "id3"], 2) returns [["id1", "id2"], ["id3"]]
 *
 * @param input the list to chunk
 * @param size the size of each chunk
 */
function chunklist(input, size) {
    const temp = input.slice(0);
    const results = [];
    while (temp.length) {
        results.push(temp.splice(0, size));
    }
    return results;
}
exports.chunklist = chunklist;
/**
 * coalesce returns the first non-empty string from the given arguments, or an empty string if all
 * of the arguments are undefined.
 *
 * @param first a potentially non-empty stringj
 * @param others other potentially non-empty values
 */
function coalesce(first, ...others) {
    if (first !== undefined && first !== "") {
        return first;
    }
    for (const entry of others) {
        if (entry !== undefined && entry !== "") {
            return entry;
        }
    }
    return "";
}
exports.coalesce = coalesce;
/**
 * coalescelist returns the first non-empty list from the given arguments, or an empty list if
 * all of the arguments are undefined or empty lists.
 *
 * @param first a potentially non-empty array
 * @param others other potentially non-empty arrays
 */
function coalescelist(first, ...others) {
    if (first !== undefined && first.length !== 0) {
        return first;
    }
    for (const entry of others) {
        if (entry !== undefined && entry.length !== 0) {
            return entry;
        }
    }
    return [];
}
exports.coalescelist = coalescelist;
/**
 * compact removes empty string elements from a list. This can be useful in some cases,
 * for example when passing joined lists as module variables or when parsing module outputs.
 *
 * @param input a string array from which to remove empty values
 */
function compact(input) {
    return input.filter(x => x !== "");
}
exports.compact = compact;
/**
 * concat combines two or more lists into a single list.
 *
 * @param first the first list
 * @param others other lists
 */
function concat(first, ...others) {
    let result = first;
    for (const entry of others) {
        result = result.concat(entry);
    }
    return result;
}
exports.concat = concat;
/**
 * contains returns true if an array contains the given element and returns false otherwise
 * @param haystack the array in which to search
 * @param needle the element for which to search
 */
function contains(haystack, needle) {
    return haystack.indexOf(needle) > -1;
}
exports.contains = contains;
function dirname(path) {
    return Path.win32.dirname(path);
}
exports.dirname = dirname;
function distinct(input) {
    return input.filter((elem, idx) => input.indexOf(elem) === idx);
}
exports.distinct = distinct;
function element(inputList, elementIndex) {
    if (inputList.length === 0) {
        throw new Error("list must not be empty");
    }
    if (elementIndex < 0) {
        throw new Error("elementIndex must be non-negative");
    }
    return inputList[elementIndex % inputList.length];
}
exports.element = element;
function file(path) {
    return fs.readFileSync(path, "utf-8");
}
exports.file = file;
function floor(input) {
    return Math.floor(input);
}
exports.floor = floor;
/**
 * format returns a string formatted using the given base string and arguments.
 *
 * *NOTE*: This is not a 1-1 map of the Terraform `format` function, since it uses the
 * node.js `sprintf-js` library instead of the implementation of `fmt.Sprintf` in the Go
 * standard library. However, it is reasonably close in semantic and will work for a
 * large number of simple use cases.
 *
 * @param formatString a format string
 * @param args the arguments to be applied to the format string
 */
function format(formatString, ...args) {
    return sprintf_js_1.sprintf(formatString, ...args);
}
exports.format = format;
function formatlist(formatString, ...args) {
    if (args.length < 1) {
        throw new Error("At least one list must be passed to formatlist");
    }
    const lengths = args.map(x => x.length);
    if (!lengths.every(x => x === lengths[0])) {
        throw new Error("All lists passed to formatlist must be the same length");
    }
    const result = [];
    for (let i = 0; i < lengths[0]; i++) {
        const iterationArgs = args.map(x => x[i]);
        result.push(sprintf_js_1.sprintf(formatString, ...iterationArgs));
    }
    return result;
}
exports.formatlist = formatlist;
function indent(indentSize, stringToIndent) {
    return stringToIndent.replace(/(\n|\r\n)/mg, le => le + " ".repeat(indentSize));
}
exports.indent = indent;
function index(haystack, needle) {
    const idx = haystack.indexOf(needle);
    if (idx > -1) {
        return idx;
    }
    // TODO(jen20): improve the error message to report which needle
    throw new Error("Could not find needle in haystack");
}
exports.index = index;
function join(separator, ...strings) {
    const parts = [].concat.apply([], strings);
    return parts.join(separator);
}
exports.join = join;
function jsonencode(input) {
    return JSON.stringify(input);
}
exports.jsonencode = jsonencode;
function keys(input) {
    return Object.keys(input);
}
exports.keys = keys;
function length(input) {
    if (typeof input === "string") {
        return input.length;
    }
    else if (input instanceof Array) {
        return input.length;
    }
    else {
        return Object.keys(input).length;
    }
}
exports.length = length;
function list(...elements) {
    return elements;
}
exports.list = list;
function log(x, base) {
    return Math.log(x) / Math.log(base);
}
exports.log = log;
function lookup(inputMap, key, defaultValue) {
    const val = inputMap[key];
    if (val === undefined) {
        if (defaultValue === undefined) {
            throw new Error("Cannot find key '" + key + "' in map");
        }
        else {
            return defaultValue;
        }
    }
    return val;
}
exports.lookup = lookup;
function lower(input) {
    return input.toLowerCase();
}
exports.lower = lower;
function map(...keyValuePairs) {
    if (keyValuePairs.length % 2 !== 0) {
        throw new Error("Number of parameters to map must be even");
    }
    const result = {};
    for (let i = 0; i < keyValuePairs.length; i += 2) {
        const key = keyValuePairs[i];
        result[key] = keyValuePairs[i + 1];
    }
    return result;
}
exports.map = map;
function max(...numbers) {
    return Math.max(...numbers);
}
exports.max = max;
function merge(...maps) {
    const output = maps[0];
    for (const x of maps) {
        for (const y in x) {
            if (x.hasOwnProperty(y)) {
                output[y] = x[y];
            }
        }
    }
    return output;
}
exports.merge = merge;
function min(...numbers) {
    return Math.min(...numbers);
}
exports.min = min;
function md5(input) {
    const hash = crypto.createHash("md5");
    hash.update(input);
    return hash.digest().toString("hex");
}
exports.md5 = md5;
function pathexpand(path) {
    const home = os.homedir();
    return path.replace(/^~(?=$|\/|\\)/, home);
}
exports.pathexpand = pathexpand;
function pow(base, exponent) {
    return Math.pow(base, exponent);
}
exports.pow = pow;
function replace(input, toMatch, replacement) {
    return input.replace(toMatch, replacement);
}
exports.replace = replace;
function sha1(input) {
    const hash = crypto.createHash("sha1");
    hash.update(input);
    return hash.digest().toString("hex");
}
exports.sha1 = sha1;
function sha256(input) {
    const hash = crypto.createHash("sha256");
    hash.update(input);
    return hash.digest().toString("hex");
}
exports.sha256 = sha256;
function sha512(input) {
    const hash = crypto.createHash("sha512");
    hash.update(input);
    return hash.digest().toString("hex");
}
exports.sha512 = sha512;
function signum(input) {
    if (input === 0) {
        return 0;
    }
    if (input < 0) {
        return -1;
    }
    else {
        return 1;
    }
}
exports.signum = signum;
function slice(input, start, end) {
    if (start < 0) {
        throw new Error("start must be greater than or equal to 0");
    }
    if (end > input.length) {
        throw new Error("end must be less than or equal to the length of input");
    }
    if (start > end) {
        throw new Error("start must be less than or equal to end");
    }
    return input.slice(start, end);
}
exports.slice = slice;
function sort(input) {
    return input.sort();
}
exports.sort = sort;
function split(separator, input) {
    return input.split(separator).filter(x => x !== "");
}
exports.split = split;
function substr(input, offset, len) {
    if (offset < 0) {
        offset += input.length;
    }
    if (len === -1) {
        len = input.length;
    }
    else if (len >= 0) {
        len += offset;
    }
    else {
        throw new Error("len should be -1, 0 or positive");
    }
    if (offset > input.length || offset < 0) {
        throw new Error("offset may not be larger than the length of input");
    }
    if (len > input.length) {
        throw new Error("offset+length may noe be larger than the length of input");
    }
    return input.substr(offset, len);
}
exports.substr = substr;
function transpose(inputMap) {
    const tempMap = {};
    for (const entry of Object.keys(inputMap)) {
        for (const key of inputMap[entry]) {
            if (tempMap[key] === undefined) {
                tempMap[key] = new Set();
            }
            tempMap[key].add(entry);
        }
    }
    const result = {};
    for (const entry of Object.keys(tempMap)) {
        result[entry] = Array.from(tempMap[entry]);
    }
    return result;
}
exports.transpose = transpose;
function timestamp() {
    return (new Date()).toISOString();
}
exports.timestamp = timestamp;
function title(input) {
    return titleCase(input);
}
exports.title = title;
function trimspace(input) {
    return input.trim();
}
exports.trimspace = trimspace;
function upper(input) {
    return input.toUpperCase();
}
exports.upper = upper;
function uuid() {
    return uuid_1.v4().toString();
}
exports.uuid = uuid;
function urlencode(input) {
    return encodeURIComponent(input);
}
exports.urlencode = urlencode;
function values(inputMap) {
    const keyList = keys(inputMap);
    const result = [];
    for (const key of keyList) {
        result.push(inputMap[key]);
    }
    return result;
}
exports.values = values;
function zipmap(keyList, valueList) {
    if (keyList.length !== valueList.length) {
        throw new Error("key and value lists must be the same length");
    }
    const result = {};
    for (let i = 0; i < keyList.length; i++) {
        result[keyList[i]] = valueList[i];
    }
    return result;
}
exports.zipmap = zipmap;
