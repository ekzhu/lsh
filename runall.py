import argparse, sys, os, time, subprocess, logging, json
logging.basicConfig(format='%(asctime)s %(message)s', level=logging.INFO)

go_progs = {"gist" :  {"simple" : "simple_gist.go",
                       "forest" : "forest_gist.go",
                       "knn"    : "knn_gist.go"},
            "image" : {"simple" : "simple_image.go",
                       "forest" : "forest_image.go",
                       "knn"    : "knn_image.go"}}

def result_name(outdir, data_type, index_type, *args):
    if len(args) % 2 != 0:
        raise ValueError("Incorrect number of arguments")
    f = "%s_%s" % (index_type, data_type)
    if len(args) > 0:
        for i in range(len(args)/2):
            f += "_%s_%s" % (args[2*i], str(args[2*i+1]))
    return os.path.join(outdir, f)

def analysis_name(outdir, data_type, index_type, *args):
    return result_name(outdir, data_type, index_type, *args) + "_analysis"

def run_simple_var_l(outs, data_type, m, ls, w, k):
    prog = go_progs[data_type]["simple"]
    for out, l in zip(outs, ls):
        logging.info("Running simple lsh l = " + str(l))
        p = subprocess.Popen(["go", "run", prog, 
            "-m", str(m), "-l", str(l), "-w", str(w), "-k", str(k),
            "-o", out])
        p.wait()

def run_knn(out, data_type, k):
    prog = go_progs[data_type]["knn"]
    logging.info("Running knn")
    p = subprocess.Popen(["go", "run", prog, "-k", str(k),
        "-o", out])
    p.wait()

def run_analysis(outs, results, ground_truth):
    prog = "analysis.go"
    for result, out in zip(results, outs):
        logging.info("Running analysis for " + result)
        p = subprocess.Popen(["go", "run", prog, 
            "-r", result, "-g", ground_truth,
            "-o", out])
        p.wait()

if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--data-type", type=str, default="gist")
    parser.add_argument("--default-w", type=float, default=1.0)
    parser.add_argument("--default-l", type=int, default=8)
    parser.add_argument("--default-m", type=int, default=5)
    parser.add_argument("--default-k", type=int, default=20)
    parser.add_argument("outdir", type=str)
    args = parser.parse_args(sys.argv[1:])
    outdir = args.outdir
    data_type = args.data_type
    m, l, w, k = args.default_m, args.default_l, args.default_w, args.default_k
    if not os.path.exists(outdir):
        os.mkdir(outdir)

    # Run experiments
    ls = [1,] + range(2, 129, 4)
    simple_results = [result_name(outdir, data_type, "simple", "l", l)
                   for l in ls]
    run_simple_var_l(simple_results, data_type, m, ls, w, k)
    knn_result = result_name(outdir, data_type, "knn")
    run_knn(knn_result, data_type, k)

    # Run analysis
    simple_analysis = [analysis_name(outdir, data_type, "simple", "l", l)
                       for l in ls]
    run_analysis(simple_analysis, simple_results, knn_result)

    # Save meta file for plotting
    metafile = os.path.join(outdir, ".meta")
    meta = {"simple_analysis" : simple_analysis,
            "ls" : ls,
            "m" : m,
            "w" : w,
            "k" : k}
    with open(metafile, 'w') as f:
        json.dump(meta, f)

