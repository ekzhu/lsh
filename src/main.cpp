#include <fstream>
#include <iostream>
#include <memory>
// #include <iomanip>
#include <vector>

using namespace std;

using flann::Matrix;

namespace {
// Constants.
const int kWidth = 32;
const int kHeight = 32;
const int kChannels = 3;
const int kSize = kWidth * kHeight * kChannels;
}  // namespace


// Reads input dataset into a matrix.
Matrix<int> ReadData(const string& filename) {
    std::ifstream is(filename, std::ifstream::binary);


    // Read entire file into local vector.
    std::vector<char> buffer(
        std::istreambuf_iterator<char>(is), 
        std::istreambuf_iterator<char>());

    cout << "Read: " << buffer.size() << " values." << endl;

    // Convert into matrix.
    int* data[] = new int[buffer.size()];
    for (int i = 0; i < buffer.size(); i++) {
        (*data)[i] = buffer[i] & 0x0000FF;
    }

    
    
/*
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
*/


    return Matrix<int>(data, buffer.size() / kSize, kSize);
}


int main(int argc, char *argv[]) {
    Matrix<int> dataset = ReadData(argv[1]);

    // int nn = 3;

    // Matrix<float> dataset;
    // Matrix<float> query;
    // load_from_file(dataset, "dataset.hdf5","dataset");
    // load_from_file(query, "dataset.hdf5","query");

    Matrix<int> indices(new int[dataset.rows * dataset.rows], dataset.rows, dataset.rows);
    Matrix<float> dists(new float[dataset.rows * dataset.rows], dataset.rows, dataset.rows);

    // construct an randomized kd-tree index using 4 kd-trees
    // Index<L2<float> > index(dataset, flann::KDTreeIndexParams(4));
    Index<L2<float>> index(dataset, flann::LinearIndexParams());
    index.buildIndex();                                                                                               

    // do a knn search, using 128 checks
    index.knnSearch(dataset, indices, dists, nn, flann::SearchParams(CHECKS_UNLIMITED));

    // flann::save_to_file(indices,"result.hdf5","result");

    delete[] dataset.ptr();
    // delete[] query.ptr();
    delete[] indices.ptr();
    delete[] dists.ptr();


    return 0;
}

