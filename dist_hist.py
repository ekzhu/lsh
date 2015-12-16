'''
Plot the histogram of top-k distances
'''

import json, sys, collections
import numpy as np
import matplotlib
matplotlib.use("Agg")
import matplotlib.pyplot as plt

with open(sys.argv[1]) as f:
    data = json.load(f)

dists = collections.deque([])
for query_result in data:
    for neighbour in query_result["neighbours"]:
        dists.append(neighbour["distance"])
dists = np.array(list(dists))

plt.hist(dists, 50)
plt.savefig("distance_histogram.png")
