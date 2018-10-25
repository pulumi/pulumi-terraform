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
import * as crypto from "crypto";
import * as fs from "fs";
import * as os from "os";
import * as Path from "path";
import { sprintf } from "sprintf-js";
import { v4 } from "uuid";
import { gzipSync } from "zlib";
import * as zlib from "zlib";

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
export function abs(input: number): number {
    return Math.abs(input);
}

/**
 * basename returns the last element of a path.
 *
 * @param path the full path
 */
export function basename(path: string): string {
    return Path.win32.basename(path);
}

/**
 * base64decode decodes base64-encoded data and returns the utf-8-encoded string. An
 * error is thrown if the argument is not valid base64 data.
 *
 * @param encoded valid base64-encoded string
 */
export function base64decode(encoded: string): string {
    const decoded = Buffer.from(encoded, "base64").toString("utf-8");
    if (decoded.length !== ((encoded.length / 4) * 3)) {
        throw new Error("failed to decode base64 data");
    }

    return decoded;
}

/**
 * base64encode returns a base64-encoded version of the given string.
 *
 * @param unencoded the string to encode
 */
export function base64encode(unencoded: string): string {
    return Buffer.from(unencoded).toString("base64");
}

/**
 * base64gzip compresses the given string with gzip and then returns the resulting data as
 * a base64-encoded string.
 *
 * @param unencoded the string to compress and encode
 */
export function base64gzip(unencoded: string): string {
    const compressed = gzipSync(unencoded, {
        level: zlib.constants.Z_DEFAULT_COMPRESSION,
    });
    return compressed.toString("base64");
}

export function base64sha256(input: string): string {
    const hash = crypto.createHash("sha256");
    hash.update(input);
    return hash.digest().toString("base64");
}

export function base64sha512(input: string): string {
    const hash = crypto.createHash("sha512");
    hash.update(input);
    return hash.digest().toString("base64");
}

/**
 * bcrypt returns the Blowfish encrypted hash of the string at the given cost. A default
 * cost of 10 will be used if not provided.
 *
 * @param password password to encrypt
 * @param cost defaults to 10 if unset
 */
export function bcrypt(password: string, cost?: number): string {
    return bcryptjs.hashSync(password, cost || 10);
}

/**
 * ceil returns the smallest integer value greater than or equal to the argument.
 *
 * @param input number of which to find the ceiling
 */
export function ceil(input: number): number {
    return Math.ceil(input);
}

export function chomp(input: string): string {
    return input.replace(/[\r\n]+$/, "");
}

/**
 * chunklist returns a list items chunked by size. For example:
 *
 * - chunklist(["id1", "id2", "id3"], 1) returns [["id1"], ["id2"], ["id3"]]
 * - chunklist(["id1", "id2", "id3"], 2) returns [["id1", "id2"], ["id3"]]
 *
 * @param input the list to chunk
 * @param size the size of each chunk
 */
export function chunklist<T>(input: T[], size: number): T[][] {
    const temp = input.slice(0);
    const results = [];

    while (temp.length) {
        results.push(temp.splice(0, size));
    }
    return results;
}

/**
 * coalesce returns the first non-empty string from the given arguments, or an empty string if all
 * of the arguments are undefined.
 *
 * @param first a potentially non-empty stringj
 * @param others other potentially non-empty values
 */
export function coalesce(first: string | undefined, ...others: (string | undefined)[]): string {
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

/**
 * coalescelist returns the first non-empty list from the given arguments, or an empty list if
 * all of the arguments are undefined or empty lists.
 *
 * @param first a potentially non-empty array
 * @param others other potentially non-empty arrays
 */
export function coalescelist<T>(first: T[] | undefined, ...others: (T[] | undefined)[]): T[] {
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

/**
 * compact removes empty string elements from a list. This can be useful in some cases,
 * for example when passing joined lists as module variables or when parsing module outputs.
 *
 * @param input a string array from which to remove empty values
 */
export function compact(input: string[]): string[] {
    return input.filter(x => x !== "");
}

/**
 * concat combines two or more lists into a single list.
 *
 * @param first the first list
 * @param others other lists
 */
export function concat<T>(first: T[], ...others: T[][]): T[] {
    let result = first;
    for (const entry of others) {
        result = result.concat(entry);
    }
    return result;
}

/**
 * contains returns true if an array contains the given element and returns false otherwise
 * @param haystack the array in which to search
 * @param needle the element for which to search
 */
export function contains<T>(haystack: T[], needle: T): boolean {
    return haystack.indexOf(needle) > -1;
}

export function dirname(path: string): string {
    return Path.win32.dirname(path);
}

export function distinct(input: string[]): string[] {
    return input.filter((elem, idx) => input.indexOf(elem) === idx);
}

export function element(inputList: string[], elementIndex: number) {
    if (inputList.length === 0) {
        throw new Error("list must not be empty");
    }
    if (elementIndex < 0) {
        throw new Error("elementIndex must be non-negative");
    }

    return inputList[elementIndex % inputList.length];
}

export function file(path: string): string {
    return fs.readFileSync(path, "utf-8");
}

export function floor(input: number): number {
    return Math.floor(input);
}

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
export function format(formatString: string, ...args: any[]): string {
    return sprintf(formatString, ...args);
}

export function formatlist(formatString: string, ...args: any[][]): string[] {
    if (args.length < 1) {
        throw new Error("At least one list must be passed to formatlist");
    }

    const lengths = args.map(x => x.length);
    if (!lengths.every(x => x === lengths[0])) {
        throw new Error("All lists passed to formatlist must be the same length");
    }

    const result: string[] = [];

    for (let i = 0; i < lengths[0]; i++) {
        const iterationArgs = args.map(x => x[i]);
        result.push(sprintf(formatString, ...iterationArgs));
    }

    return result;
}

export function indent(indentSize: number, stringToIndent: string) {
    return stringToIndent.replace(/(\n|\r\n)/mg, le => le + " ".repeat(indentSize));
}

export function index<T>(haystack: T[], needle: T): number {
    const idx = haystack.indexOf(needle);
    if (idx > -1) {
        return idx;
    }

    // TODO(jen20): improve the error message to report which needle
    throw new Error("Could not find needle in haystack");
}

export function join(separator: string, ...strings: string[][]) {
    const parts: string[] = [].concat.apply([], strings);

    return parts.join(separator);
}

export function jsonencode(input: any): string {
    return JSON.stringify(input);
}

export function keys(input: { [k: string]: any }): string[] {
    return Object.keys(input);
}

export function length(input: string | any[] | { [k: string]: any }): number {
    if (typeof input === "string") {
        return input.length;
    } else if (input instanceof Array) {
        return input.length;
    } else {
        return Object.keys(input).length;
    }
}

export function list<T>(...elements: T[]): T[] {
    return elements;
}

export function log(x: number, base: number): number {
    return Math.log(x) / Math.log(base);
}

export function lookup<T>(inputMap: { [k: string]: T }, key: string, defaultValue?: T): T {
    const val = inputMap[key];
    if (val === undefined) {
        if (defaultValue === undefined) {
            throw new Error("Cannot find key '" + key + "' in map");
        } else {
            return defaultValue;
        }
    }

    return val;
}

export function lower(input: string): string {
    return input.toLowerCase();
}

export function map(...keyValuePairs: any[]): { [k: string]: any } {
    if (keyValuePairs.length % 2 !== 0) {
        throw new Error("Number of parameters to map must be even");
    }

    const result: { [k: string]: any } = {};

    for (let i = 0; i < keyValuePairs.length; i += 2) {
        const key = keyValuePairs[i];
        result[key] = keyValuePairs[i + 1];
    }

    return result;
}

export function max(...numbers: number[]): number {
    return Math.max(...numbers);
}

export function merge(...maps: { [k: string]: any }[]): { [k: string]: any } {
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

export function min(...numbers: number[]): number {
    return Math.min(...numbers);
}

export function md5(input: string): string {
    const hash = crypto.createHash("md5");
    hash.update(input);
    return hash.digest().toString("hex");
}

export function pathexpand(path: string): string {
    const home = os.homedir();
    return path.replace(/^~(?=$|\/|\\)/, home);
}

export function pow(base: number, exponent: number): number {
    return Math.pow(base, exponent);
}

export function replace(input: string, toMatch: string | RegExp, replacement: string): string {
    return input.replace(toMatch, replacement);
}

export function sha1(input: string): string {
    const hash = crypto.createHash("sha1");
    hash.update(input);
    return hash.digest().toString("hex");
}

export function sha256(input: string): string {
    const hash = crypto.createHash("sha256");
    hash.update(input);
    return hash.digest().toString("hex");
}

export function sha512(input: string): string {
    const hash = crypto.createHash("sha512");
    hash.update(input);
    return hash.digest().toString("hex");
}

export function signum(input: number): number {
    if (input === 0) {
        return 0;
    }
    if (input < 0) {
        return -1;
    } else {
        return 1;
    }
}

export function slice(input: string[], start: number, end: number): string[] {
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

export function sort(input: string[]): string[] {
    return input.sort();
}

export function split(separator: string, input: string): string[] {
    return input.split(separator).filter(x => x !== "");
}

export function substr(input: string, offset: number, len: number): string {
    if (offset < 0) {
        offset += input.length;
    }

    if (len === -1) {
        len = input.length;
    } else if (len >= 0) {
        len += offset;
    } else {
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

export function transpose(inputMap: { [k: string]: string[] }): { [k: string]: string[] } {
    const tempMap: { [k: string]: Set<string> } = {};

    for (const entry of Object.keys(inputMap)) {
        for (const key of inputMap[entry]) {
            if (tempMap[key] === undefined) {
                tempMap[key] = new Set<string>();
            }

            tempMap[key].add(entry);
        }
    }

    const result: { [k: string]: string[] } = {};
    for (const entry of Object.keys(tempMap)) {
        result[entry] = Array.from(tempMap[entry]);
    }
    return result;
}

export function timestamp(): string {
    return (new Date()).toISOString();
}

export function title(input: string): string {
    return titleCase(input);
}

export function trimspace(input: string): string {
    return input.trim();
}

export function upper(input: string): string {
    return input.toUpperCase();
}

export function uuid(): string {
    return v4().toString();
}

export function urlencode(input: string): string {
    return encodeURIComponent(input);
}

export function values<T>(inputMap: { [k: string]: T }): T[] {
    const keyList = keys(inputMap);

    const result: T[] = [];
    for (const key of keyList) {
        result.push(inputMap[key]);
    }

    return result;
}

export function zipmap<T>(keyList: string[], valueList: T[]): { [k: string]: T } {
    if (keyList.length !== valueList.length) {
        throw new Error("key and value lists must be the same length");
    }

    const result: { [k: string]: T } = {};
    for (let i = 0; i < keyList.length; i++) {
        result[keyList[i]] = valueList[i];
    }

    return result;
}
