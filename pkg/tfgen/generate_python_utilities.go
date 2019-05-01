package tfgen

const pyUtilitiesFile = `
import os
import pkg_resources

def get_env(*args):
    for v in args:
        value = os.getenv(v)
        if value is not None:
            return value
    return None

def get_env_bool(*args):
    str = get_env(*args)
    if str is not None:
        # NOTE: these values are taken from https://golang.org/src/strconv/atob.go?s=351:391#L1, which is what
        # Terraform uses internally when parsing boolean values.
        if str in ["1", "t", "T", "true", "TRUE", "True"]:
            return True
        if str in ["0", "f", "F", "false", "FALSE", "False"]:
            return False
    return None

def get_env_int(*args):
    str = get_env(*args)
    if str is not None:
        try:
            return int(str)
        except:
            return None
    return None

def get_env_float(*args):
    str = get_env(*args)
    if str is not None:
        try:
            return float(str)
        except:
            return None
    return None

def require_with_default(req, default):
    try:
        return req()
    except:
        if default is not None:
            return default
        raise

def get_version():
    # __name__ is set to the fully-qualified name of the current module, In our case, it will be
    # <some module>.utilities. <some module> is the module we want to query the version for.
    root_package, *rest = __name__.split('.')

    # pkg_resources uses setuptools to inspect the set of installed packages. We use it here to ask
    # for the currently installed version of the root package (i.e. us) and get its version.
    return pkg_resources.require(root_package)[0].version
`
