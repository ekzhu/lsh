import json, sys, collections
import numpy as np
import scipy
from scipy.stats import gamma, norm
from scipy.integrate import quad
from scipy.constants import pi
import matplotlib.pyplot as plt

def _integration(a, b, f, p):
    area = 0.0
    x = a
    while x < b:
        area += f(x+0.5*p)*p
        x += p
    return area

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

def _recall(m, l, w, gamma_params, base, max_x):
    k = len(gamma_params)
    s = 0.0
    for i in range(k):
        shape, loc, scale = gamma_params[i]
        join_prob_func = lambda x : _hash_probability(m, l, w, np.sqrt(x)) * gamma.pdf(x, shape, loc, scale) * base
        prob, _ = quad(join_prob_func, 0.0, max_x) 
        s += prob
    return s / float(k)

def _selectivity(m, l, w, gamma_param, base, max_x):
    shape, loc, scale = gamma_param
    join_prob_func = lambda x : _hash_probability(m, l, w, np.sqrt(x)) * gamma.pdf(x, shape, loc, scale) * base
    prob, _ = quad(join_prob_func, 0.0, max_x) 
    return prob

def optimization(max_m, l, max_w, max_x, gamma_x, gamma_xk, required_recall, base_x, base_xk):
    best_m = 0
    best_w = 0.0
    best_selectivity = float('inf')
    for m in range(1, max_m):
        # Search for the m and w that gives the smallest recall just above the required_recall
        # Use binary search
        right_bound = max_w
        left_bound = 0.0
        w = (right_bound - left_bound) / 2.0
        delta = float('inf')
        while delta > 1.0:
            recall = _recall(m, l, w, gamma_xk, base_xk, max_x)
            print("recall", recall)
            if recall < required_recall:
                left_bound = w
            else:
                right_bound = w
            new_w = (right_bound - left_bound) / 2.0
            if new_w < 0.0:
                print(left_bound, right_bound, w, new_w, m, l, recall)
                raise ValueError()
            delta = np.abs(new_w - w)
            w = new_w
        if recall < required_recall - 0.01:
            print("Failed for l = %d m = %d is w =  %f, recall = %f" % (l, m, w, recall))
            continue
        selectivity = _selectivity(m, l, w, gamma_x, base_x, max_x) 
        print("Best for l = %d m = %d is w =  %f, recall = %f, selectivity = %f" % (l, m, w, recall, selectivity))
        if selectivity < best_selectivity:
            best_selectivity = selectivity
            best_m = m
            best_w = w
    print("Best overall for l = %d is m = %d, w = %d" % (l, best_m, best_w))
    return best_m, best_w

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
    for i in range(k):
        params = gamma.fit(topks[:,i])
        gamma_xk.append(params)
        print("k = %d distance-squared distribution: " % i, params)
    return gamma_xk, topks 

all_pairs_sample = "./_image_all_pair_distance_sample"
topk_sample = "./_image_query_distance_sample"
k = 50
max_w = 15000 
max_m = 8
required_recall = 0.9
dataset_size = 10000
output = "opt_param_k_%d_recall_%.2f.json" % (k, required_recall)
base_x = 1.0
base_xk = 1.0

gamma_x, dists_squared = load_all_pair_sample(all_pairs_sample)
gamma_xk, topk_dists_squared = load_topk_sample(topk_sample)
max_x = np.max(dists_squared)
ls = [1, 4, 8, 16, 32]
out = []
for l in ls:
    m, w = optimization(max_m, l, max_w, max_x, 
            gamma_x, gamma_xk, required_recall, base_x, base_xk)
    out.append({"L" : l, "M" : m, "W" : w})
with open("opt_param.json", 'w') as f:
    json.dump(f, out)
