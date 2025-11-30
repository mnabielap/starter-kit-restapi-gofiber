import sys
import os
sys.path.append(os.path.abspath(os.path.dirname(__file__)))
import utils

# Endpoint
URL = f"{utils.BASE_URL}/auth/logout"

# Load Refresh Token
refresh_token = utils.load_config("refresh_token")

if not refresh_token:
    print(">> [ERROR] No refresh_token found in secrets.json.")
else:
    payload = {
        "refreshToken": refresh_token
    }

    utils.send_and_print(
        url=URL,
        method="POST",
        body=payload,
        output_file=f"{os.path.splitext(os.path.basename(__file__))[0]}.json",
    )