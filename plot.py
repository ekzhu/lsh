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

    # var L experiments
    metafile = os.path.join(args.varlout, ".meta")
    with open(metafile) as f:
        varlmeta = json.load(f)
    ls = varlmeta["Ls"]
    ms = varlmeta["Ms"]
    ws = varlmeta["Ws"]
    varl_analysis = get_analysis(varlmeta)
    
    # var T experiments
    metafile = os.path.join(args.vartout, ".meta")
    with open(metafile) as f:
        vartmeta = json.load(f)
    ts = vartmeta["Ts"]
    vart_analysis = get_analysis(vartmeta)
    
    
    # Plot recall
    fig, axes = plt.subplots(1, 2, figsize=(10, 5), sharey=True)
    #axes[0].set_ylim(0.5, 1.0)
    axes[0].set_xscale('log', basex=2)
    axes[0].grid()
    for label in varl_analysis:
        recall = varl_analysis[label]["recalls"]
        axes[0].plot(ls, recall, label=label, marker="+")
    axes[0].set_xlabel("Number of hash tables")
    axes[0].set_ylabel("Recall")
    axes[0].legend(loc="lower right")
    axes[0].set_title("M = %d, W = %d, T = %d" % (ms[0], ws[0], varlmeta["T"]))
    axes[1].grid()
    axes[1].set_xscale('log', basex=2)
    for label in vart_analysis:
        recall = vart_analysis[label]["recalls"]
        axes[1].plot(ts, recall, label=label, marker="+")
    axes[1].set_xlabel("Number of probes")
    axes[1].set_ylabel("Recall")
    axes[1].set_title("M = %d, L = %d, W = %d" % (vartmeta["M"], vartmeta["L"], vartmeta["W"]))
    fig.savefig("recall.png")
    plt.close()

    # Plot time
    fig, axes = plt.subplots(1, 2, figsize=(10, 5), sharey=True)
    axes[0].set_xscale('log', basex=2)
    axes[0].grid()
    for label in varl_analysis:
        times = varl_analysis[label]["times"]
        axes[0].plot(ls, times, label=label, marker="+")
    axes[0].set_xlabel("Number of hash tables")
    axes[0].set_ylabel("90 percentil query time (ms)")
    axes[0].legend(loc="upper left")
    axes[0].set_title("T = %d" % (varlmeta["T"]))
    axes[1].grid()
    axes[1].set_xscale('log', basex=2)
    for label in vart_analysis:
        times = vart_analysis[label]["times"]
        axes[1].plot(ts, times, label=label, marker="+")
    axes[1].set_xlabel("Number of probes")
    axes[1].set_ylabel("90 percentil query time (ms)")
    axes[1].set_title("M = %d, L = %d, W = %d" % (vartmeta["M"], vartmeta["L"], vartmeta["W"]))
    fig.savefig("time.png")
    plt.close()

