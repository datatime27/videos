package main

import "fmt"
import "math/rand"

func create_shuffled_numbers(num_prisoners int) []int{
    a := []int{}
    for i := 0; i < num_prisoners; i++{
        a = append(a,i)
    }
    rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
    return a
}
    
func single_prisoner_attempt(prisoner_number int, boxes []int) bool{
    box_to_open := prisoner_number
    for i := 0; i < len(boxes)/2; i++{
        if prisoner_number == boxes[box_to_open]{
            return true
        }
        box_to_open = boxes[box_to_open]
    }
    return false
}

func run_simulation(num_prisoners int) bool{
    boxes := create_shuffled_numbers(num_prisoners)
    for prisoner_number := 0; prisoner_number < num_prisoners; prisoner_number++{
        if single_prisoner_attempt(prisoner_number, boxes) == false {
            return false
        }
    }
    return true
}

const num_prisoners = 10
func main() {    
    var successes, num_total_simulations int64
    
    // Run 1,000,000 experiments
    for i:=0; i < 1000000; i++ {
        if run_simulation(num_prisoners) == true {
            successes++
        }
        num_total_simulations++
        fmt.Printf("%d Prisoners: Successful Simulations: %d out of %d %.2f:1 %.4f\n", 
            num_prisoners, 
            successes, 
            num_total_simulations, 
            float32(num_total_simulations)/float32(successes), 
            100.0*float32(successes)/float32(num_total_simulations))
    }
}
