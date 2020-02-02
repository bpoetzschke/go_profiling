# CPU profiling
This example demonstrates the use case of profiling a golang application with the Block profile.
The application in this example uses for loops and channels to create some artificial blocking. A Block profile
is automatically generated in the root directory of the application.

## Run the example
1. Compile the binary.
    ```bash
   go build block.go
   ```
2. Run the binary.
    ```bash
   time ./block
   ```
3. By using the `time` command we can see how the program is running.
    ```
   2020/02/02 17:14:32 profile: block profiling enabled, block.pprof
   1048576
   2020/02/02 17:14:35 profile: block profiling disabled, block.pprof
   ./block  0.00s user 0.00s system 0% cpu 2.240 total
    ```
   
   We see that the program is running 2.24 seconds in total.
4. Let's have a look at the blocking profile why that is.
    1. Open the profiling tool.
        ```
        go tool pprof block.pprof
        ``` 
   
   2. Show the top nodes with the `top` sub-command.
        The output should look similar to the following:
        
        ```
        (pprof) top
        Showing nodes accounting for 21.62s, 100% of 21.62s total
        Dropped 8 nodes (cum <= 0.11s)
            flat  flat%   sum%        cum   cum%
          21.62s   100%   100%     21.62s   100%  runtime.chanrecv1
               0     0%   100%     19.57s 90.49%  main.generate
               0     0%   100%      2.06s  9.51%  main.main
               0     0%   100%      2.06s  9.51%  runtime.main
        ```
      
    3. We can see that `runtime.chanrecv1` is causing the application to be blocked for 21.61 seconds.
     The value we see is higher than the actual runtime which `time` reported because the application is using 
     go routines to compute the result.
    4. By using the `list` sub-command we can also look into the annotated source of our application to see which line
    in `main.generate` is causing the block.
        ```
       (pprof) list main.generate
       Total: 21.62s
       ROUTINE ======================== main.generate in /Users/bjoern/repos/golang_profiling/block/block.go
                0     19.57s (flat, cum) 90.49% of Total
                .          .     11:
                .          .     12:	"github.com/pkg/profile"
                .          .     13:)
                .          .     14:
                .          .     15:func generate(in <-chan int, out chan<- int) {
                .     19.57s     16:	i := <-in
                .          .     17:	time.Sleep(100 * time.Millisecond)
                .          .     18:	out <- i + i
                .          .     19:}
                .          .     20:
                .          .     21:func main() {
       ```
        We can see that line 16 in `block.go` is causing the block.
5. The block profile only analyses the code which is causing a block. Other operations which are causing the program
to run longer are not taken into account.
    1. Run the program now with the environment variable `SLEEP=true` to generate a new profile.
        ```
       time SLEEP=true ./block
       ``` 
       
       The output should look similar to the following:
       
       ```
       2020/02/02 17:30:08 profile: block profiling enabled, block.pprof
       1048576
       Sleep
       2020/02/02 17:30:20 profile: block profiling disabled, block.pprof
       SLEEP=true ./block  0.00s user 0.01s system 0% cpu 12.166 total
       ```
       
       We see now that the program runs in total 12.166 seconds. This is expected because with the environment variable
       set to true we put the program into sleep for 10 seconds.
       
   2. Let's have a look at the profile.
       ```
        go tool pprof block.pprof
       ```
   3. Show the top nodes with the `top` sub-command.
       ```
       (pprof) top
       Showing nodes accounting for 21.48s, 100% of 21.48s total
       Dropped 8 nodes (cum <= 0.11s)
             flat  flat%   sum%        cum   cum%
           21.48s   100%   100%     21.48s   100%  runtime.chanrecv1
                0     0%   100%     19.43s 90.45%  main.generate
                0     0%   100%      2.05s  9.55%  main.main
                0     0%   100%      2.05s  9.55%  runtime.main
       ```
      
      We still see that `runtime.chanrecv1` is causing the program to be blocked for a total of 21.48 seconds. The sleep
      we introduced is not taken into consideration when using the block profile.
   
# Acknowledgements
This example is taken from [Dave Cheney](https://twitter.com/davecheney)'s 
[`High Performance Go Workshop`](https://dave.cheney.net/high-performance-go-workshop/gophercon-2019.html) 
from GopherCon 2019.

# License
See original license: https://github.com/davecheney/high-performance-go-workshop#license-and-materials

## Changes made
The program was adjusted to include a environment variable to define whether the program should 
add an additional sleep step.