import json, sys, collections
import numpy as np
import scipy
from scipy.stats import gamma, norm
from scipy.integrate import quad
from scipy.constants import pi
import matplotlib.pyplot as plt

def _collision_probability(w, r):
    a = 1.0 - 2.0 * norm.cdf(- w / r)
    b = 2.0 / (np.sqrt(2.0 * pi) * w / r)
    c = 1.0 - np.exp(- (w * w) / (2.0 * r * r))
    return a - b * c

def _hash_probability(m, l, w, r):
    p = 1.0 - (1.0 - _collision_probability(w, r)**float(m))**float(l)
    if p < 0.0:
        print(m, l, w, r)
        raise ValueError()
    return p

def _recall(m, l, w, gamma_params):
    k = len(gamma_params)
    s = 0.0
    for i in range(k):
        shape, loc, scale = gamma_params[i]
        join_prob_func = lambda x : _hash_probability(m, l, w, np.sqrt(x)) * gamma.pdf(x, shape, loc, scale)
        prob, _ = quad(join_prob_func, 0.0, float('inf')) 
        s += prob
    return s / float(k)

def _selectivity(m, l, w, gamma_param):
    shape, loc, scale = gamma_param
    join_prob_func = lambda x : _hash_probability(m, l, w, np.sqrt(x)) * gamma.pdf(x, shape, loc, scale)
    prob, _ = quad(join_prob_func, 0.0, float('inf')) 
    return prob

def optimization(max_m, l, max_w, gamma_x, gamma_xk, min_recall):
    best_m = 0
    best_w = 0.0
    best_selectivity = float('inf')
    for m in range(1, max_m):
        # Search for the m and w that gives the smallest recall just above the min_recall
        # Use binary search
        right_bound = max_w
        left_bound = 0.0
        w = (right_bound - left_bound) / 2.0
        delta = float('inf')
        while delta > 1.0:
            recall = _recall(m, l, w, gamma_xk)
            if recall < min_recall:
                left_bound = w
            else:
                right_bound = w
            new_w = (right_bound - left_bound) / 2.0
            if new_w < 0.0:
                print(left_bound, right_bound, w, new_w, m, l, recall)
                raise ValueError()
            delta = np.abs(new_w - w)
            w = new_w
        selectivity = _selectivity(m, l, w, gamma_x) 
        print("Best for l = %d m = %d is w =  %f, selectivity = %f" % (l, m, w, selectivity))
        if selectivity < best_selectivity:
            best_selectivity = selectivity
            best_m = m
            best_w = w
    print("Best overall for l = %d is m = %d, w = %d" % (l, best_m, best_w))
    return best_m, best_w

def plot_gamma(param, d):
    shape, loc, scale = param
    size = len(d)
    x = np.linspace(0, int(np.max(d)), num=50)
    pdf_fitted = gamma.pdf(x, *param[:-2], loc=param[-2], scale=param[-1])*size
    plt.plot(x, pdf_fitted)
    plt.hist(d, 20, histtype="stepfilled", alpha=0.7)
    plt.xlim(xmin=0)
    plt.show()
    raw_input("Press ENTER")

all_pairs_sample = "./_image_all_pair_distance_sample"
dist_sample = "./_image_query_distance_sample"
max_k = 500
max_w = 0.0 
max_m = 32
min_recall = 0.5
c1 = 1000.0
c2 = 1000000.0

# Obtain distance distributions of any pair of points
with open(all_pairs_sample) as f:
    data = json.load(f)
dists = collections.deque([])
for query_result in data:
    for neighbour in query_result["neighbours"]:
        dists.append(neighbour["distance"])
dists = np.array(list(dists))
max_w = np.max(dists)
dists_squared = np.square(dists)
gamma_x = gamma.fit(dists_squared)
print("All pair distance distribution: ", gamma_x)
#plot_gamma(gamma_x, dists_squared)

# Obtain distance distributions of k nearest neighbours
with open(dist_sample) as f:
    data = json.load(f)
topks = []
for query_result in data:
    dists = collections.deque([])
    for neighbour in query_result["neighbours"]:
        dists.append(neighbour["distance"])
    dists = np.sort(dists)
    topks.append(dists)
topks = np.array(topks)
gamma_xk = []
for k in range(max_k):
    d = np.square(topks[:,k])
    params = gamma.fit(d)
    gamma_xk.append(params)
    print("k = %d distance distribution: " % k, params)
    #plot_gamma(params, d)

m = 4
l = 4
x = np.linspace(1.0, max_w, 100)
recalls = [_recall(m, l, w, gamma_xk) for w in x]
plt.plot(x, recalls)
plt.show()

ls = [1, 4, 8, 16, 32]
out = []
for l in ls:
    m, w = optimization(max_m, l, max_w, gamma_x, gamma_xk, min_recall)
    out.append({"L" : l, "M" : m, "W" : w})
with open(opt_param.json, 'w') as f:
    json.dump(f, out)
