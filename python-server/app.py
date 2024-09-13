from flask import Flask, request
import json

app = Flask(__name__)

@app.route('/log_summary', methods=['POST'])
def log_summary():
    data = request.get_json()
    with open("vehicle_summary.txt", "a") as file:
        file.write(json.dumps(data) + "\n")
    return {"status": "success"}, 200

if __name__ == '__main__':
    app.run(host="0.0.0.0", port=5000)
