package tfgen

const tsUtilitiesFile = `
import * as pulumi from "@pulumi/pulumi";

export function getEnv(...vars: string[]): string | undefined {
    for (const v of vars) {
        const value = process.env[v];
        if (value) {
            return value;
        }
    }
    return undefined;
}

export function getEnvBoolean(...vars: string[]): boolean | undefined {
    const s = getEnv(...vars);
    if (s !== undefined) {
        // NOTE: these values are taken from https://golang.org/src/strconv/atob.go?s=351:391#L1, which is what
        // Terraform uses internally when parsing boolean values.
        if (["1", "t", "T", "true", "TRUE", "True"].find(v => v === s) !== undefined) {
            return true;
        }
        if (["0", "f", "F", "false", "FALSE", "False"].find(v => v === s) !== undefined) {
            return false;
        }
    }
    return undefined;
}

export function getEnvNumber(...vars: string[]): number | undefined {
    const s = getEnv(...vars);
    if (s !== undefined) {
        const f = parseFloat(s);
        if (!isNaN(f)) {
            return f;
        }
    }
    return undefined;
}

export function requireWithDefault<T>(req: () => T, def: T | undefined): T {
    try {
        return req();
    } catch (err) {
        if (def === undefined) {
            throw err;
        }
    }
    return def;
}

export function unwrap(val: pulumi.Input<any>): pulumi.Input<any> {
    // Bottom out at primitives.
    if (val === undefined || typeof val !== 'object') {
        return val;
    }

    // Recurse on outputs, promises, arrays, and objects.
    if (pulumi.Output.isInstance(val)) {
        return val.apply(unwrap);
    }
    if (val instanceof Promise || val instanceof Array) {
        return pulumi.output(val).apply(unwrap);
    }

    const array = Object.keys(val).map(k =>
        pulumi.output(unwrap(val[k])).apply(v => ({ key: k, value: v })));

    return pulumi.all(array).apply(keysAndValues => {
        const result: any = {};
        for (const kvp of keysAndValues) {
            result[kvp.key] = kvp.value;
        }
        return result;
    });
}
`
