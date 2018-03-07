import os
import time
import sys
from flask import Flask

app = Flask(__name__)

port = os.getenv('PORT', '3000')
host = '0.0.0.0'

count = 0

@app.route('/')
def index():
    global count
    count = count +1
    print '{"success": false, "hello": "world", "foo": %d}' % count
    sys.stdout.flush()
    return "Success: Python + Flask"

@app.route('/crash')
def crash():
    sys.exit(1)

@app.route('/5')
def five():
    time.sleep(5)
    return "Success: Python + Flask"

@app.route('/9')
def nine():
    time.sleep(9)
    return "Success: Python + Flask"

@app.route('/10')
def ten():
    time.sleep(10)
    return "Success: Python + Flask"

@app.route('/15')
def fifteen():
    time.sleep(15)
    return "Success: Python + Flask"

@app.route('/login/ensure_availability')
def nope():
    time.sleep(9999)
    return "Success: Python + Flask"


if __name__ == "__main__":
    print os.environ.get("VCAP_SERVICES")
    print "Listening on port [%s]\n" % port

    app.run(host=host, port=int(port))
