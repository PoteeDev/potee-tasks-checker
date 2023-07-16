import sys
import time
import random


address = sys.argv[1]
input_file = sys.argv[2]

with open(input_file) as f:
    data = f.read()
wait = random.randint(1,10)/10
time.sleep(wait)
print(f"{wait} send '{data}' to '{address}'")
