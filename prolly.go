package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const binarySizeCutoff = 256

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: prollySearch <size>")
		return
	}

	size, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Invalid size argument:", err)
		return
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	lastVal := r.Int()
	vals := make([]int, 0, size)
	for len(vals) < cap(vals) {
		lastVal += r.Intn(23)
		vals = append(vals, lastVal)
	}

	searchCount := vals[len(vals)-1] - vals[0]
	startTimeBin := time.Now()
	for i := vals[0]; i <= vals[len(vals)-1]; i++ {
		_ = binarySearch(vals, i)
	}
	durBin := time.Since(startTimeBin)

	startTime := time.Now()
	for i := vals[0]; i <= vals[len(vals)-1]; i++ {
		_ = prollyBinSearch(vals, i)
	}
	dur := time.Since(startTime)

	fmt.Printf("%d,%d,%d,%d\n", len(vals), searchCount, dur.Nanoseconds(), durBin.Nanoseconds())
}

func prollyBinSearch(slice []int, target int) int {
	if len(slice) == 0 {
		return -1
	}

	low := 0
	high := len(slice) - 1
	for low <= high {
		if slice[low] > target {
			return -1
		} else if slice[low] == target {
			return low
		}
		if slice[high] < target {
			return -1
		} else if slice[high] == target {
			return high
		}

		if high-low > binarySizeCutoff {
			// Determine the estimated position of the target in the slice, as a float from 0 to 1.
			minVal := slice[low]
			maxVal := slice[high]
			shiftedTarget := target - minVal
			shiftedMax := maxVal - minVal

			est := float64(shiftedTarget) / float64(shiftedMax)
			estIdx := int(float64(high-low) * est)
			estIdx += low

			if estIdx >= len(slice) {
				estIdx = len(slice) - 1
			}

			if slice[estIdx] == target {
				return estIdx // bulls-eye!
			}

			// When we miss the target, we know that we are pretty close based on the assumption of distribution.
			// Therefore, unlike a binary search where we consider everything on the left or right, we instead do
			// a scan in the appropriate direction using a widening scope. When all is said and done, low and high
			// will be set to values which are pretty close to the guess.
			widenScope := 16
			if slice[estIdx] > target {
				// We overshot, so search left
				high = estIdx - 1
				newLow := high - widenScope
				for newLow > low && slice[newLow] > target {
					high = newLow // just verified that newLow is higher than target
					widenScope <<= 2
					newLow = high - widenScope
				}
				if newLow > low {
					low = newLow
				}
			} else {
				// We undershot, so search right
				low = estIdx + 1
				newHigh := low + widenScope
				for newHigh < high && slice[newHigh] < target {
					low = newHigh // just verified that newHigh is lower than target
					widenScope <<= 2
					newHigh = low + widenScope
				}
				if newHigh < high {
					high = newHigh
				}
			}
		} else {
			// Fall back to binary search
			for low <= high {
				mid := low + (high-low)/2
				if slice[mid] == target {
					return mid // Found
				} else if slice[mid] < target {
					low = mid + 1 // Search right half
				} else {
					high = mid - 1 // Search left half
				}
			}
			return -1 // Not found
		}
	}
	return -1
}

// binarySearch vanilla style!
func binarySearch(slice []int, target int) int {
	low := 0
	high := len(slice) - 1
	for low <= high {
		mid := low + (high-low)/2
		if slice[mid] == target {
			return mid // Found
		} else if slice[mid] < target {
			low = mid + 1 // Search right half
		} else {
			high = mid - 1 // Search left half
		}
	}
	return -1 // Not found
}
