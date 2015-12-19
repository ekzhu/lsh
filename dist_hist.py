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

max_k = 20
topks = []
for query_result in data:
    dists = collections.deque([])
    for neighbour in query_result["neighbours"]:
        dists.append(neighbour["distance"])
    dists = np.sort(dists)
    topks.append(dists)
topks = np.array(topks)

fig, axes = plt.subplots(max_k, 1, figsize=(3, max_k*3))
for k in range(max_k):
    axes[k].hist(topks[:,k], 20)
plt.savefig("kth_dist.png")
