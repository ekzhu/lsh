#include <iostream>
#include <fstream>
#include <memory>
#include <iomanip>

using std::cin;




int main(int argc, char *argv[]) {
    std::ifstream is(argv[1], std::ifstream::binary);

    const int kWidth = 32;
    const int kHeight = 32;
    const int kChannels = 3;
    const int kSize = kWidth * kHeight * kChannels;

    if (is) {
        char* buffer = new char[kWidth*kHeight*kChannels];
	is.read(buffer, kWidth * kHeight * kChannels);
	for (int i = 0; i < kSize; i++) {
	    std::cout << (int) (buffer[i] & 0x0000FF);
	    if (i % kWidth == 31) {
	        std::cout << std::endl;
	    } else {
	        std::cout << " ";
	    }
	}

	delete[] buffer;

    }


    int nn = 3;

    Matrix<float> dataset;
    Matrix<float> query;
    load_from_file(dataset, "dataset.hdf5","dataset");
    load_from_file(query, "dataset.hdf5","query");

    Matrix<int> indices(new int[query.rows*nn], query.rows, nn);
    Matrix<float> dists(new float[query.rows*nn], query.rows, nn);

    // construct an randomized kd-tree index using 4 kd-trees
    Index<L2<float> > index(dataset, flann::KDTreeIndexParams(4));
    index.buildIndex();                                                                                               

    // do a knn search, using 128 checks
    index.knnSearch(query, indices, dists, nn, flann::SearchParams(128));

    flann::save_to_file(indices,"result.hdf5","result");

    delete[] dataset.ptr();
    delete[] query.ptr();
    delete[] indices.ptr();
    delete[] dists.ptr();


    return 0;
}
