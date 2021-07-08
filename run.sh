#! /bin/bash

# Copyright 2017 George S. Kong. All rights reserved.
# Use of this source code is governed by a license that can be found in the LICENSE.txt file.

# run a development server locally
#
# look for environment variable DM_DEV_SECRETS_TOML, specifying a customized dev secrets file.
# if the environment variable doesn't exist, use the example secrets file.

secretFile=${DM_DEV_SECRETS_TOML:-config/example-secrets.toml}

./dm config/base.toml config/dev.toml $secretFile
