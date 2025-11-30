import sys
import os
sys.path.append(os.path.abspath(os.path.dirname(__file__)))
import utils

# Endpoint
URL = f"{utils.BASE_URL}/auth/login"

# Payload
payload = {
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

# Update secrets.json
if response.status_code == 200:
    data = response.json()
    
    utils.save_config("access_token", data["tokens"]["access"]["token"])
    utils.save_config("refresh_token", data["tokens"]["refresh"]["token"])
    utils.save_config("user_id", data["user"]["id"])
    
    print(">> [INFO] Tokens updated in secrets.json")