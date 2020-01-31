# CPU profiling
This example demonstrates the use case of profiling a golang application with the CPU profile.
The application in this example takes a path to a text file and counts the words in it. A CPU profile
is automatically generated in the root directory of the application.

## Run the example
1. Compile the binary:
   
   ```
   go buid cpu.go
   ```  
   The binary will be located in the root directory with the name `cpu`.
2. Run the binary with the text file where we want to count the words.

    ```
    time ./cpu <path_to_text_file>
   ```
3. By using time we will know how much time the application will spent counting the words in the text file. You will
    see a output similar to the following.
    
    ```bash
   2020/01/31 11:18:37 profile: cpu profiling enabled, cpu.pprof
   "moby.txt": 181276 words
   2020/01/31 11:18:38 profile: cpu profiling disabled, cpu.pprof
   ./cpu moby.txt  0.49s user 0.69s system 87% cpu 1.352 total
   ```
4. Now let's compare the result with the native word count application, `wc`.
    
    ```bash
       22333  215830 1276201 moby.txt
    wc moby.txt  0.01s user 0.00s system 87% cpu 0.009 total
   ```
    We can see two differences between our implementation and the implementation of `wc`.
    * Different count of words, we can neglect this.
    * `wc` runs much faster
    
5. To figure out why our application is slow we will analyze the CPU profile which we already gathered.

    1. To analyze the profile run the following command.
        ```bash
        go tool pprof cpu.pprof
        ```
    2. Show top entries with `top`
        ```
       (pprof) top
       Showing nodes accounting for 700ms, 100% of 700ms total
             flat  flat%   sum%        cum   cum%
            690ms 98.57% 98.57%      690ms 98.57%  syscall.syscall
             10ms  1.43%   100%      700ms   100%  internal/poll.(*FD).Read
                0     0%   100%      700ms   100%  main.main
                0     0%   100%      700ms   100%  main.readbyte
                0     0%   100%      700ms   100%  main.v0
                0     0%   100%      700ms   100%  os.(*File).Read
                0     0%   100%      700ms   100%  os.(*File).read
                0     0%   100%      700ms   100%  runtime.main
                0     0%   100%      690ms 98.57%  syscall.Read
                0     0%   100%      690ms 98.57%  syscall.read
       ```
   3. Printing out the top entries reveals that `syscall.syscall` is using most of the time. This is caused by calling 
   `main.readbyte`. This is not surprising since `syscall.syscall` is a expensive and slow operation.
        
   4. To improve the application we can introduce buffered reader which will help us improving the performance of the 
   application.
6. The sample application already has the `bufio.Reader` implemented. Run the application again, but this time with 
    the `VERSION` env var set to `1`.
    ```bash
   time VERSION=1 ./cpu <path_to_text_file>
   ``` 
7. Now we can see that the performance of the program increased by a magnitude just by changing to a buffered reader.
    ```
    2020/01/31 13:46:11 profile: cpu profiling enabled, cpu.pprof
    "moby.txt": 181276 words
    2020/01/31 13:46:11 profile: cpu profiling disabled, cpu.pprof
    VERSION=1 ./cpu moby.txt  0.04s user 0.01s system 18% cpu 0.247 total
   ```
   The `top` sub-command in the go profiler also shows that we got rid of `syscall.syscall`. The program now only spents
   a small amount of time to allocate memory and to read the file.
   ```
   (pprof) top
   Showing nodes accounting for 30ms, 100% of 30ms total
         flat  flat%   sum%        cum   cum%
         20ms 66.67% 66.67%       20ms 66.67%  runtime.mallocgc
         10ms 33.33%   100%       10ms 33.33%  bufio.(*Reader).Read
            0     0%   100%       30ms   100%  main.main
            0     0%   100%       30ms   100%  main.readbyte
            0     0%   100%       30ms   100%  main.v1
            0     0%   100%       30ms   100%  runtime.main
            0     0%   100%       20ms 66.67%  runtime.newobject
   ```
    
# Acknowledgements
This example is taken from [Dave Cheney](https://twitter.com/davecheney)'s 
[`High Performance Go Workshop`](https://dave.cheney.net/high-performance-go-workshop/gophercon-2019.html) 
from GopherCon 2019.

# License
See original license: https://github.com/davecheney/high-performance-go-workshop#license-and-materials

## Changes made
Changes were made to support a environment var to select between buffered reader and regular reader.