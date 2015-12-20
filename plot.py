import json, sys, argparse, os
import numpy as np
import matplotlib
matplotlib.use("Agg")
import matplotlib.pyplot as plt

def get_analysis(meta):
    a = {}
    for analysis_result in meta["analysis_results"]:
        label = analysis_result["algorithm"]
        result_files = analysis_result["result_files"]
        error_ratios = []
        recalls = []
        times = []
        for result_file in result_files:
            with open(result_file) as f:
                analysis = json.load(f)
            error_ratios.append(np.mean(analysis["errorratios"]))
            recalls.append(np.mean(analysis["recalls"]))
            times.append(np.percentile(analysis["times"], 90))
        a[label] = {"error_ratios" : error_ratios,
                    "recalls" : recalls,
                    "times" : times}
    return a

if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("varlout")
    parser.add_argument("vartout")
    args = parser.parse_args(sys.argv[1:])

    # Plot var L experiments
    metafile = os.path.join(args.varlout, ".meta")
    with open(metafile) as f:
        meta = json.load(f)
    ls = meta["Ls"]
    analysis = get_analysis(meta)
    
    # Plot recall
    fig, axes = plt.subplots(1, 1)
    for label in analysis:
        recall = analysis[label]["recalls"]
        axes.plot(ls, recall, label=label)
    axes.set_xlabel("Number of hash tables")
    axes.set_ylabel("Recall")
    axes.legend()
    fig.savefig("recall_var_l.png")
    plt.close()

    # Plot error_ratio
    fig, axes = plt.subplots(1, 1)
    for label in analysis:
        error_ratios = analysis[label]["error_ratios"]
        axes.plot(ls, error_ratios, label=label)
    axes.set_xlabel("Number of hash tables")
    axes.set_ylabel("Error ratio")
    axes.legend()
    fig.savefig("error_ratio_var_l.png")
    plt.close()

    # Plot time
    fig, axes = plt.subplots(1, 1)
    for label in analysis:
        times = analysis[label]["times"]
        axes.plot(ls, times, label=label)
    axes.set_xlabel("Number of hash tables")
    axes.set_ylabel("90 percentil query time (ms)")
    axes.legend()
    fig.savefig("time_var_l.png")
    plt.close()

    # Plot var T experiments
    metafile = os.path.join(args.vartout, ".meta")
    with open(metafile) as f:
        meta = json.load(f)
    ts = meta["Ts"]
    analysis = get_analysis(meta)
    
    # Plot recall
    fig, axes = plt.subplots(1, 1)
    for label in analysis:
        recall = analysis[label]["recalls"]
        axes.plot(ts, recall, label=label)
    axes.set_xlabel("Number of probes")
    axes.set_ylabel("Recall")
    axes.legend()
    fig.savefig("recall_var_t.png")
    plt.close()

    # Plot error_ratio
    fig, axes = plt.subplots(1, 1)
    for label in analysis:
        error_ratios = analysis[label]["error_ratios"]
        axes.plot(ts, error_ratios, label=label)
    axes.set_xlabel("Number of probes")
    axes.set_ylabel("Error ratio")
    axes.legend()
    fig.savefig("error_ratio_var_t.png")
    plt.close()

    # Plot time
    fig, axes = plt.subplots(1, 1)
    for label in analysis:
        times = analysis[label]["times"]
        axes.plot(ts, times, label=label)
    axes.set_xlabel("Number of probes")
    axes.set_ylabel("90 percentil query time (ms)")
    axes.legend()
    fig.savefig("time_var_t.png")
    plt.close()

