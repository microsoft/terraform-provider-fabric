# Copyright Microsoft Corporation 2026
# SPDX-License-Identifier: MPL-2.0

import fabric.functions as fn

udf = fn.UserDataFunctions()


@udf.function()
def hello_fabric(name: str) -> str:
    return f"Hello {name}!"
