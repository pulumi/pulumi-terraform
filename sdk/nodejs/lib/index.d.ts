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
export declare function abs(input: number): number;
/**
 * basename returns the last element of a path.
 *
 * @param path the full path
 */
export declare function basename(path: string): string;
/**
 * base64decode decodes base64-encoded data and returns the utf-8-encoded string. An
 * error is thrown if the argument is not valid base64 data.
 *
 * @param encoded valid base64-encoded string
 */
export declare function base64decode(encoded: string): string;
/**
 * base64encode returns a base64-encoded version of the given string.
 *
 * @param unencoded the string to encode
 */
export declare function base64encode(unencoded: string): string;
/**
 * base64gzip compresses the given string with gzip and then returns the resulting data as
 * a base64-encoded string.
 *
 * @param unencoded the string to compress and encode
 */
export declare function base64gzip(unencoded: string): string;
export declare function base64sha256(input: string): string;
export declare function base64sha512(input: string): string;
/**
 * bcrypt returns the Blowfish encrypted hash of the string at the given cost. A default
 * cost of 10 will be used if not provided.
 *
 * @param password password to encrypt
 * @param cost defaults to 10 if unset
 */
export declare function bcrypt(password: string, cost?: number): string;
/**
 * ceil returns the smallest integer value greater than or equal to the argument.
 *
 * @param input number of which to find the ceiling
 */
export declare function ceil(input: number): number;
export declare function chomp(input: string): string;
/**
 * chunklist returns a list items chunked by size. For example:
 *
 * - chunklist(["id1", "id2", "id3"], 1) returns [["id1"], ["id2"], ["id3"]]
 * - chunklist(["id1", "id2", "id3"], 2) returns [["id1", "id2"], ["id3"]]
 *
 * @param input the list to chunk
 * @param size the size of each chunk
 */
export declare function chunklist<T>(input: T[], size: number): T[][];
/**
 * coalesce returns the first non-empty string from the given arguments, or an empty string if all
 * of the arguments are undefined.
 *
 * @param first a potentially non-empty stringj
 * @param others other potentially non-empty values
 */
export declare function coalesce(first: string | undefined, ...others: (string | undefined)[]): string;
/**
 * coalescelist returns the first non-empty list from the given arguments, or an empty list if
 * all of the arguments are undefined or empty lists.
 *
 * @param first a potentially non-empty array
 * @param others other potentially non-empty arrays
 */
export declare function coalescelist<T>(first: T[] | undefined, ...others: (T[] | undefined)[]): T[];
/**
 * compact removes empty string elements from a list. This can be useful in some cases,
 * for example when passing joined lists as module variables or when parsing module outputs.
 *
 * @param input a string array from which to remove empty values
 */
export declare function compact(input: string[]): string[];
/**
 * concat combines two or more lists into a single list.
 *
 * @param first the first list
 * @param others other lists
 */
export declare function concat<T>(first: T[], ...others: T[][]): T[];
/**
 * contains returns true if an array contains the given element and returns false otherwise
 * @param haystack the array in which to search
 * @param needle the element for which to search
 */
export declare function contains<T>(haystack: T[], needle: T): boolean;
export declare function dirname(path: string): string;
export declare function distinct(input: string[]): string[];
export declare function element(inputList: string[], elementIndex: number): string;
export declare function file(path: string): string;
export declare function floor(input: number): number;
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
export declare function format(formatString: string, ...args: any[]): string;
export declare function formatlist(formatString: string, ...args: any[][]): string[];
export declare function indent(indentSize: number, stringToIndent: string): string;
export declare function index<T>(haystack: T[], needle: T): number;
export declare function join(separator: string, ...strings: string[][]): string;
export declare function jsonencode(input: any): string;
export declare function keys(input: {
    [k: string]: any;
}): string[];
export declare function length(input: string | any[] | {
    [k: string]: any;
}): number;
export declare function list<T>(...elements: T[]): T[];
export declare function log(x: number, base: number): number;
export declare function lookup<T>(inputMap: {
    [k: string]: T;
}, key: string, defaultValue?: T): T;
export declare function lower(input: string): string;
export declare function map(...keyValuePairs: any[]): {
    [k: string]: any;
};
export declare function max(...numbers: number[]): number;
export declare function merge(...maps: {
    [k: string]: any;
}[]): {
    [k: string]: any;
};
export declare function min(...numbers: number[]): number;
export declare function md5(input: string): string;
export declare function pathexpand(path: string): string;
export declare function pow(base: number, exponent: number): number;
export declare function replace(input: string, toMatch: string | RegExp, replacement: string): string;
export declare function sha1(input: string): string;
export declare function sha256(input: string): string;
export declare function sha512(input: string): string;
export declare function signum(input: number): number;
export declare function slice(input: string[], start: number, end: number): string[];
export declare function sort(input: string[]): string[];
export declare function split(separator: string, input: string): string[];
export declare function substr(input: string, offset: number, len: number): string;
export declare function transpose(inputMap: {
    [k: string]: string[];
}): {
    [k: string]: string[];
};
export declare function timestamp(): string;
export declare function title(input: string): string;
export declare function trimspace(input: string): string;
export declare function upper(input: string): string;
export declare function uuid(): string;
export declare function urlencode(input: string): string;
export declare function values<T>(inputMap: {
    [k: string]: T;
}): T[];
export declare function zipmap<T>(keyList: string[], valueList: T[]): {
    [k: string]: T;
};
