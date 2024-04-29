package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	lastVal := r.Int()
	vals := make([]int, 0, 100000000)
	for len(vals) < cap(vals) {
		lastVal += r.Intn(23)
		vals = append(vals, lastVal)
	}
	searchCount := vals[len(vals)-1] - vals[0]
	startTimeBin := time.Now()
	for i := vals[0]; i <= vals[len(vals)-1]; i++ {
		_ = macnealeSearch(vals, i)
	}
	durBin := time.Since(startTimeBin)

	startTime := time.Now()
	for i := vals[0]; i <= vals[len(vals)-1]; i++ {
		_ = prollyBinSearch(vals, i)
	}
	dur := time.Since(startTime)

	fmt.Printf("Searches performed: %d. On slice size: %d\n", searchCount, len(vals))
	fmt.Printf("ProllySearch: %v\n", dur)
	fmt.Printf("BinarySearch: %v\n", durBin)
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

		if high-low > 32 {
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

			neiborhoodSize := (high - low) >> 7 // 1/128th of the current slice size.
			if neiborhoodSize < 32 {
				neiborhoodSize = 32
				if neiborhoodSize > high-low {
					neiborhoodSize = high - low
				}
			}

			if slice[estIdx] > target {
				// We overshot, so search left
				high = estIdx - 1
				mayLow := high - neiborhoodSize
				if mayLow >= low && slice[mayLow] <= target {
					low = mayLow
				} // else, our estimate was way off. handle???
			} else {
				// We undershot, so search right
				low = estIdx + 1
				mayHigh := low + neiborhoodSize
				if mayHigh <= high && slice[mayHigh] >= target {
					high = mayHigh
				} // else, our estimate was way off. handle???
			}
		} else {
			// low and high are close, so just do a linear search
			low++
			high--
		}
	}
	return -1
}

// macnealeSearch vanilla style!
func macnealeSearch(slice []int, target int) int {
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
