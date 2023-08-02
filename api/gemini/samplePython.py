import requests
import json
import base64
import hmac
import hashlib
import datetime, time

url = "https://api.gemini.com/v1/mytrades"
gemini_api_key = "mykey"
gemini_api_secret = "1234abcd".encode()

t = datetime.datetime.now()
payload_nonce =  str(int(time.mktime(t.timetuple())*1000))
payload =  {"request": "/v1/mytrades", "nonce": payload_nonce}
encoded_payload = json.dumps(payload).encode()
b64 = base64.b64encode(encoded_payload)
signature = hmac.new(gemini_api_secret, b64, hashlib.sha384).hexdigest()

request_headers = {
    'Content-Type': "text/plain",
    'Content-Length': "0",
    'X-GEMINI-APIKEY': gemini_api_key,
    'X-GEMINI-PAYLOAD': b64,
    'X-GEMINI-SIGNATURE': signature,
    'Cache-Control': "no-cache"
    }

response = requests.post(url, headers=request_headers)

my_trades = response.json()
print(my_trades)