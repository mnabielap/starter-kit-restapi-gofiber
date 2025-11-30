import sys
import os
sys.path.append(os.path.abspath(os.path.dirname(__file__)))
import utils

# Endpoint
URL = f"{utils.BASE_URL}/auth/register"

# Payload
payload = {
    "name": "New User",
    "email": "test@example.com",
    "password": "password123"
}

# Send Request
response = utils.send_and_print(
    url=URL,
    method="POST",
    body=payload,
    output_file=f"{os.path.splitext(os.path.basename(__file__))[0]}.json",
)

# Save Tokens and User ID to secrets.json for other scripts to use
if response.status_code == 201:
    data = response.json()
    
    # Save Access Token
    access_token = data.get("tokens", {}).get("access", {}).get("token")
    if access_token:
        utils.save_config("access_token", access_token)
        print(">> [INFO] Access Token saved to secrets.json")

    # Save Refresh Token
    refresh_token = data.get("tokens", {}).get("refresh", {}).get("token")
    if refresh_token:
        utils.save_config("refresh_token", refresh_token)
        print(">> [INFO] Refresh Token saved to secrets.json")

    # Save User ID
    user_id = data.get("user", {}).get("id")
    if user_id:
        utils.save_config("user_id", user_id)
        print(">> [INFO] User ID saved to secrets.json")