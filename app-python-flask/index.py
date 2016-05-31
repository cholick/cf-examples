import os
from flask import Flask

app = Flask(__name__)

port = os.getenv('PORT', '3000')
host = '0.0.0.0'

@app.route('/')
def index():
    return "Success: Python + Flask"

if __name__ == "__main__":
    print "Listening on port [%s]\n" % port

    app.run(host=host, port=port)
