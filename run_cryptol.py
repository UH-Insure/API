#!/usr/bin/env python3
import sys
import json
from cryptol import connect

def main():
    try:
        req = json.load(sys.stdin)
        code = req.get("code", "").strip()
        if not code:
            print(json.dumps({"stdout": "", "stderr": "No Cryptol code provided"}))
            return

        # Connect to the Cryptol remote API server
        conn = connect(url="http://localhost:8080", reset_server=True)

        # Evaluate code
        result = conn.eval(code)

        print(json.dumps({"stdout": str(result), "stderr": ""}))
    except Exception as e:
        print(json.dumps({"stdout": "", "stderr": str(e)}))

if __name__ == "__main__":
    main()