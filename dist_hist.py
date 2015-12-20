'''
Plot the histogram of top-k distances
'''

import json, sys, collections
import numpy as np
import scipy
from scipy.stats import gamma
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
    param = gamma.fit(topks[:,k])
    x = scipy.arange(len(topks[:,k])) 
    pdf_fitted = gamma.pdf(x, *param[:-2], loc=param[-2], scale=param[-1])*len(topks[:,k])
    #axes[k].hist(topks[:,k], 20)
    axes[k].plot(pdf_fitted)
    axes[k].set_xlim(xmin=0)
plt.savefig("kth_dist.png")
