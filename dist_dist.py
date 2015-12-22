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

def load_all_pair_sample(datafile):
    with open(datafile) as f:
        data = json.load(f)
    dists = collections.deque([])
    for query_result in data:
        for neighbour in query_result["neighbours"]:
            dists.append(neighbour["distance"])
    dists = np.array(list(dists))
    dists_squared = np.square(dists)
    gamma_x = gamma.fit(dists_squared)
    print("Distance-squared distribution: ", gamma_x)
    return gamma_x, dists_squared

def load_topk_sample(datafile):
    with open(datafile) as f:
        data = json.load(f)
    topks = collections.deque([])
    for query_result in data:
        dists = collections.deque([])
        for neighbour in query_result["neighbours"]:
            dists.append(neighbour["distance"])
        dists_squared = np.square(np.sort(list(dists)))
        topks.append(dists_squared)
    topks = np.array(list(topks))
    gamma_xk = []
    for i in ks:
        params = gamma.fit(topks[:,i])
        gamma_xk.append(params)
        print("k = %d distance-squared distribution: " % i, params)
    return gamma_xk, topks 

all_pairs_sample = "./_image_all_pair_distance_sample"
topk_sample = "./_image_query_distance_sample"
ks = [10, 50, 200]
max_w = 15000.0 
max_m = 12
required_recall = 0.5

gamma_x, dists_squared = load_all_pair_sample(all_pairs_sample)
gamma_xk, topk_dists_squared = load_topk_sample(topk_sample)
max_x = np.max(dists_squared)


fig, axes = plt.subplots(1, 2, figsize=(10, 5), sharex=True)

# Plot all pair distance distribution
x = np.linspace(0.0, max_x/2.0, num=100)
pdf = gamma.pdf(x, gamma_x[0], gamma_x[1], gamma_x[2]) 
axes[0].plot(x, pdf)
axes[0].grid()
axes[0].set_ylabel("Probability")
axes[0].set_xlabel("Sqaured L2 distance") 

# Plot kth nearest neighbour distance distribution
for i, k in enumerate(ks):
    shape, loc, scale = gamma_xk[i]
    pdf = gamma.pdf(x, shape, loc, scale) 
    axes[1].plot(x, pdf, label="%d-NN" % k)
axes[1].legend()
axes[1].set_ylabel("Probability")
axes[1].set_xlabel("Sqaured L2 distance") 
axes[1].grid()

plt.savefig("dist_dist.png")
plt.close()
