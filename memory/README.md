# Memory profiling
This example demonstrates the use case of profiling a Go application with the memory profile. This application will
allocate memory and trigger the garbage collection before exiting. A memory profile will automatically be generate 
during the runtime of program.

## Run the example
1. Compile the binary.
    ```bash
   go build memory.go
   ```
2. Run the binary.
    ```bash
   ./memory
   ```
3. Let's have a look at the memory profile.
    ```bash
   go tool pprof mem.pprof 
   ```
   
   The output should look similar to the following:
   
   ```
   Showing nodes accounting for 25.39kB, 99.39% of 25.55kB total
   Dropped 15 nodes (cum <= 0.13kB)
   Showing top 10 nodes out of 45
         flat  flat%   sum%        cum   cum%
      12.75kB 49.91% 49.91%    12.75kB 49.91%  runtime.malg
          7kB 27.40% 77.31%    12.25kB 47.95%  runtime.allocm
       4.16kB 16.27% 93.58%     4.17kB 16.33%  time.LoadLocationFromTZData
       0.38kB  1.47% 95.05%     0.38kB  1.47%  os/signal.signal_enable
       0.28kB  1.10% 96.15%     0.28kB  1.10%  runtime.acquireSudog
       0.25kB  0.98% 97.13%     0.25kB  0.98%  runtime.allgadd
       0.20kB   0.8% 97.92%     0.58kB  2.26%  os/signal.Notify
       0.19kB  0.73% 98.65%     0.28kB  1.10%  runtime.gcBgMarkWorker
       0.11kB  0.43% 99.08%     0.69kB  2.69%  github.com/pkg/profile.Start.func10
       0.08kB  0.31% 99.39%     4.27kB 16.70%  log.(*Logger).Output
   ```
   We can see that the memory profile reports a total allocation size of just only 25kB. If we are looking at the 
   implementation of the application we just profiled we can see that the allocation size should be definitely more
   than what we see right now. This is because the go profiler shows the `inuse_space` as default sample index. This
   sample index type just shows the live in-use memory at the time when the sampling in the application happened.
   To see the total number of allocation we can switch to the `allocation_space` sample index. Luckily we do not need
   to change the application to generate a new profile, because this sample index is already existing in the profile.
   To see a list of all available sample indices run `go tool pprof` and the help will be printed.
4. Let's have a look at the `allocation_space` profile
    ```bash
   go tool pprof -alloc_space mem.pprof 
   ```
   
   The output should look similar to the following:
   
   ```bash
   Showing nodes accounting for 14940.44kB, 99.79% of 14972.20kB total
   Dropped 75 nodes (cum <= 74.86kB)
         flat  flat%   sum%        cum   cum%
   14255.98kB 95.22% 95.22% 14940.44kB 99.79%  main.allocate
     684.45kB  4.57% 99.79%   684.45kB  4.57%  main.makeByteSlice
            0     0% 99.79% 14950.87kB 99.86%  main.main
            0     0% 99.79% 14950.87kB 99.86%  runtime.main
   ```
   
   Now we can see that the application is allocating approximately 14. MB of memory during the runtime. Also we can see
   that `main.allocate` allocates the most memory.
   
5. To see in detail where `main.allocate` allocates the most memory we can have the annotated source of this method
    shown. To do so run the following sub-command inside `pprof`.
    ```
   list main.allocate
   ```
   
   This output of this sub-command should look similar to the following:
   ```
   Total: 14.62MB
   ROUTINE ======================== main.allocate in /Users/bjoern/repos/golang_profiling/memory/memory.go
      13.92MB    14.59MB (flat, cum) 99.79% of Total
            .          .     21:
            .          .     22:// allocate allocates count byte slices and returns the first slice allocated.
            .          .     23:func allocate() []byte {
            .          .     24:	var x [][]byte
            .          .     25:	for i := 0; i < count; i++ {
      13.92MB    14.59MB     26:		x = append(x, makeByteSlice())
            .          .     27:	}
            .          .     28:	return x[0]
            .          .     29:}
            .          .     30:
            .          .     31:// makeByteSlice returns a byte slice of a random length in the range [0, 16384).
   ```
   
   Thanks to this annotated output we can identify the exact line of code where the allocation is happening. We can also
   see that line 26 is responsible for allocation 13.92 MB of a total of 14.59 MB. Note the 14.59 MB are based on total
   size of allocations shown in the profile. The top output earlier shows that 75 nodes where dropped.
   The information gathered by the `list` sub-command can be extremely helpful when it comes to remove 
   unnecessary allocations.
   
 