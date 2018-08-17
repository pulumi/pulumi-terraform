package tfgen

const pyUtilitiesFile = `
import os

def get_env(*args):
    for v in args:
        value = os.getenv(v)
        if value is not None:
            return value
    return None

def get_env_bool(*args):
    str = get_env(*args)
    if str is not None:
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
`