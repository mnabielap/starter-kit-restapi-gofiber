import sys
import os
sys.path.append(os.path.abspath(os.path.dirname(__file__)))
import utils

# Endpoint
URL = f"{utils.BASE_URL}/auth/forgot-password"

payload = {
    "email": "test@example.com"
}

utils.send_and_print(
    url=URL,
    method="POST",
    body=payload,
    output_file=f"{os.path.splitext(os.path.basename(__file__))[0]}.json",
)