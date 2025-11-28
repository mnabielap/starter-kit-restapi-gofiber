import sys
import os
sys.path.append(os.path.abspath(os.path.dirname(__file__)))
import utils

# Endpoint
URL = f"{utils.BASE_URL}/users"

# Auth Header
token = utils.load_config("access_token")
headers = {
    "Authorization": f"Bearer {token}"
}

# Payload
payload = {
    "name": "Admin Created User",
    "email": "admin_created@example.com",
    "password": "password123",
    "role": "user"
}

if not token:
    print(">> [ERROR] No access_token found.")
else:
    utils.send_and_print(
        url=URL,
        method="POST",
        headers=headers,
        body=payload,
        output_file=f"{os.path.splitext(os.path.basename(__file__))[0]}.json",
    )