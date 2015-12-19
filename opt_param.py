import numpy as np
from scipy.stats import norm
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
    return 1 - (1 - _collision_probability(w, r)**float(m))**float(l)



m = 10
l = 25
w = 1000.0
min_r = 0.0
max_r = 5000.0
rs = np.linspace(min_r, max_r, num=100)
probs = [_collision_probability(w, r) for r in rs]
plt.plot(rs, probs)
plt.show()
probs = [_hash_probability(m, l, w, r) for r in rs] 
plt.plot(rs, probs)
plt.show()
