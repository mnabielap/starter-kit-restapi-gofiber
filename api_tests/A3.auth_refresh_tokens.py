import sys
import os
sys.path.append(os.path.abspath(os.path.dirname(__file__)))
import utils

# Endpoint
URL = f"{utils.BASE_URL}/auth/refresh-tokens"

# Load Refresh Token
refresh_token = utils.load_config("refresh_token")

if not refresh_token:
    print(">> [ERROR] No refresh_token found in secrets.json. Please Login first.")
else:
    # Payload
    payload = {
        "refreshToken": refresh_token
    }

    # Send Request
    response = utils.send_and_print(
        url=URL,
        method="POST",
        body=payload,
        output_file=f"{os.path.splitext(os.path.basename(__file__))[0]}.json",
    )

    # Update secrets.json
    if response.status_code == 200:
        data = response.json()
        utils.save_config("access_token", data["tokens"]["access"])
        utils.save_config("refresh_token", data["tokens"]["refresh"])
        print(">> [INFO] Tokens refreshed and saved.")