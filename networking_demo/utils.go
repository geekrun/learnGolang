package networking_demo

import "math/rand"

// shuffle shuffles the elements of an array in place
func shuffle(array []uint16) []uint16 {
	for i := range array { //run the loop till the range of array
		j := rand.Intn(i + 1)                   //choose any random number
		array[i], array[j] = array[j], array[i] //swap the random element with current element
	}
	return array
}
