import sys
import os
sys.path.append(os.path.abspath(os.path.dirname(__file__)))
import utils

# Endpoint
URL = f"{utils.BASE_URL}/auth/send-verification-email"

# Auth Header
token = utils.load_config("access_token")
headers = {
    "Authorization": f"Bearer {token}"
}

if not token:
    print(">> [ERROR] No access_token found. Please Login.")
else:
    utils.send_and_print(
        url=URL,
        method="POST",
        headers=headers,
        output_file=f"{os.path.splitext(os.path.basename(__file__))[0]}.json",
    )