# Benchmark

The benchmarks have been run on the sib-dev machine (`Intel(R) Xeon(R) CPU @ 2.00GHz`) using the command:

    go test -bench=. -benchmem -cpu 1,2,4

The results are:

    BenchmarkBufPool            	 2000000	       784 ns/op	    4864 B/op	       1 allocs/op
    BenchmarkBufPool-2          	 1000000	      1518 ns/op	    4864 B/op	       1 allocs/op
    BenchmarkBufPool-4          	 1000000	      1742 ns/op	    4865 B/op	       1 allocs/op
    BenchmarkBufWithoutPool     	  500000	      2354 ns/op	   16816 B/op	       6 allocs/op
    BenchmarkBufWithoutPool-2   	  500000	      3140 ns/op	   16816 B/op	       6 allocs/op
    BenchmarkBufWithoutPool-4   	  500000	      3775 ns/op	   16816 B/op	       6 allocs/op

Benchmarked on August 8th 2018.
