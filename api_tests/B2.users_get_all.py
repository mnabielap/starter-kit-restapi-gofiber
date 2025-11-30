import sys
import os
sys.path.append(os.path.abspath(os.path.dirname(__file__)))
import utils

# Endpoint (with pagination query params)
URL = f"{utils.BASE_URL}/users?limit=10&page=1"

# Auth Header
token = utils.load_config("access_token")
headers = {
    "Authorization": f"Bearer {token}"
}

if not token:
    print(">> [ERROR] No access_token found.")
else:
    utils.send_and_print(
        url=URL,
        method="GET",
        headers=headers,
        output_file=f"{os.path.splitext(os.path.basename(__file__))[0]}.json",
    )