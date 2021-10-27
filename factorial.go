package main

func Factorial(n int) uint64 {

	var factVal uint64 = 1

	if n < 0 {

		panic("Factorial of 0 does not exist!")

	} else {

		for i := 1; i <= n; i++ {

			factVal *= uint64(i) // mismatched types int64 and int

		}

	}
	return factVal /* return from function*/
}
