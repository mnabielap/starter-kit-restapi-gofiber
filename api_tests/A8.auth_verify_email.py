import sys
import os
sys.path.append(os.path.abspath(os.path.dirname(__file__)))
import utils

# Ensure 'verify_token' is in secrets.json
token = utils.load_config("verify_token")

if not token:
    print(">> [ERROR] 'verify_token' not found in secrets.json.")
    print(">> Please check server logs/email, get the token, and add it to secrets.json with key 'verify_token'")
else:
    # Endpoint
    URL = f"{utils.BASE_URL}/auth/verify-email?token={token}"

    utils.send_and_print(
        url=URL,
        method="POST",
        output_file=f"{os.path.splitext(os.path.basename(__file__))[0]}.json",
    )