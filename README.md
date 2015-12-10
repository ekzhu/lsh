# CSC 2515 Project

## Build Instructions

### Flann

To build flann:

1.  Create subdirectory (e.g., `build`) in flann main directory
2.  `cd build && cmake ..`
3.  `make`
4.  `sudo make install` to install shared libraries (libflann\_cpp)

### OpenCV

To build OpenCV:

1.  Create subdirectory (e.g., `build`) in flann main directory
2.  `cd build && cmake ..`
3.  `make`
4.  `sudo make install` to install shared libraries (libopencv)

## To do

1.  Use LSH on small dataset using flann
2.  Use brute-force kNN on small dataset using OpenCV
3.  Write python script for calculating precision/recall on test results

