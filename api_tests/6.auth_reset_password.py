import sys
import os
sys.path.append(os.path.abspath(os.path.dirname(__file__)))
import utils

# Config
# Ensure you have saved 'reset_token' in secrets.json manually before running this
token = utils.load_config("reset_token") 

if not token:
    print(">> [ERROR] 'reset_token' not found in secrets.json.")
    print(">> Please check server logs/email, get the token, and add it to secrets.json with key 'reset_token'")
else:
    # Endpoint
    URL = f"{utils.BASE_URL}/auth/reset-password?token={token}"

    payload = {
        "password": "newPassword123"
    }

    utils.send_and_print(
        url=URL,
        method="POST",
        body=payload,
        output_file=f"{os.path.splitext(os.path.basename(__file__))[0]}.json",
    )