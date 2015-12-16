import json, sys, argparse, os
import numpy as np
import matplotlib
matplotlib.use("Agg")
import matplotlib.pyplot as plt

def get_analysis(name):
    error_ratios = []
    recalls = []
    times = []
    for analysis_file in meta[name]:
        with open(analysis_file) as f:
            analysis = json.load(f)
        error_ratios.append(np.mean(analysis["errorratios"]))
        recalls.append(np.mean(analysis["recalls"]))
        times.append(np.percentile(analysis["times"], 90))
    return error_ratios, recalls, times

if __name__ == "__main__":
    analysis_dir = sys.argv[1]
    metafile = os.path.join(analysis_dir, ".meta")
    with open(metafile) as f:
        meta = json.load(f)

    ls = meta["ls"]
    possible_labels = ["simple_analysis", "forest_analysis", "multiprobe_analysis"]
    labels = []
    error_ratios = []
    recalls = []
    times = []
    for label in possible_labels:
        if label in meta:
            e, r, t = get_analysis(label)
            error_ratios.append(e)
            recalls.append(r)
            times.append(t)
            labels.append(label)

    # Plot recall
    fig, axes = plt.subplots(1, 1)
    for recall, label in zip(recalls, labels):
        axes.plot(ls, recall, label=label)
    axes.set_xlabel("Number of hash tables")
    axes.set_ylabel("Recall")
    fig.savefig("recall.png")
    plt.close()

    # Plot error_ratio
    fig, axes = plt.subplots(1, 1)
    for error_ratio, label in zip(error_ratios, labels):
        axes.plot(ls, error_ratio, label=label)
    axes.set_xlabel("Number of hash tables")
    axes.set_ylabel("Error ratio")
    fig.savefig("error_ratio.png")
    plt.close()

    # Plot time
    fig, axes = plt.subplots(1, 1)
    for time, label in zip(times, labels):
        axes.plot(ls, time, label=label)
    axes.set_xlabel("Number of hash tables")
    axes.set_ylabel("90 percentil query time (ms)")
    fig.savefig("time.png")
    plt.close()

