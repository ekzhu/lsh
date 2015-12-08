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
    return 0;
}
