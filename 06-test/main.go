package main

func sum(vals []int) (result int) {
	for _, val := range vals {
		result += val
	}
	return
}
