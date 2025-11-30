import sys
import os
sys.path.append(os.path.abspath(os.path.dirname(__file__)))
import utils

# Config
user_id = utils.load_config("user_id")

if not user_id:
    print(">> [ERROR] 'user_id' not found in secrets.json.")
else:
    # Endpoint
    URL = f"{utils.BASE_URL}/users/{user_id}"

    # Auth Header
    token = utils.load_config("access_token")
    headers = {
        "Authorization": f"Bearer {token}"
    }

    utils.send_and_print(
        url=URL,
        method="DELETE",
        headers=headers,
        output_file=f"{os.path.splitext(os.path.basename(__file__))[0]}.json",
    )