# Mutex profiling
This example demonstrates the use case of profiling a Go application with the Mutex profile.
The application in this example uses for loops and mutexes to create artificial mutex contention. 
A Mutex profile is automatically generated in the root directory of the application.

## Run the example
1. Compile the binary.
    ```bash
   go build mutex.go
   ```
2. Run the binary.
    ```bash
   time ./mutex
   ```
3. By using the `time` command we can see how the program is running.
    ```
   2020/02/02 18:10:13 profile: mutex profiling enabled, mutex.pprof
   Start mutex contention
   Wait
   Sleep 5s
   2020/02/02 18:10:26 profile: mutex profiling disabled, mutex.pprof
   ./mutex  34.31s user 5.71s system 311% cpu 12.854 total
   ```
   
   We can see that the total runtime of this program is 12.854 seconds, 5 seconds of the total runtime
   are caused due to a sleep to demonstrate the capabilities of the mutex profile.
4. Let's have a look at the profile.
    ```
    go tool pprof mutex.pprof
   ``` 
   
   1. Show the top nodes with the `top` sub-command.
       ```
       (pprof) top
       Showing nodes accounting for 7.33s, 100% of 7.33s total
             flat  flat%   sum%        cum   cum%
            7.33s   100%   100%      7.33s   100%  sync.(*Mutex).Unlock
                0     0%   100%      7.33s   100%  main.main.func1
       ```
        The `flat` column shows the total time spent in the respective method excluding other method calls. 
        The `cum` columns shows the total time spent in the respective method including other method calls.
   2. We can clearly see that `sync.(*Mutex).Unlock` is causing the program to slow down. This function is 
   responsible for all the delay which is caused in the application.
   
        Furthermore we see that `main.main.func1` is not causing any slow down because the `flat` column is zero.
        However we can see in the `cum` column that a function call inside this method is causing the slow down.
    3. By using the `list` sub-command we can take a look at the annoted source of `main.main.func1` and see what
    exactly is causing the slow down.
        ```
       (pprof) list main.main
       Total: 7.33s
       ROUTINE ======================== main.main.func1 in /Users/bjoern/repos/golang_profiling/mutex/mutex.go
                0      7.33s (flat, cum)   100% of Total
                .          .     26:			mu.Lock()
                .          .     27:			defer mu.Unlock()
                .          .     28:			defer wg.Done()
                .          .     29:
                .          .     30:			items[i] = struct{}{}
                .      7.33s     31:		}(i)
                .          .     32:	}
                .          .     33:
                .          .     34:	fmt.Println("Wait")
                .          .     35:	wg.Wait()
                .          .     36:	fmt.Println("Sleep 5s")
       (pprof)
       ```
       We can clearly see that the go routine, which is called in line 31, is responsible for this slow down.
       However it's only indirect responsible since we don't have any entry in the `flat` column.
       
       To get even more insights we can also have the annotated source of `sync.(*Mutex).Unlock` shown.
       ```
       (pprof) list Unlock
       Total: 7.33s
       ROUTINE ======================== sync.(*Mutex).Unlock in /usr/local/Cellar/go/1.13.7/libexec/src/sync/mutex.go
            7.33s      7.33s (flat, cum)   100% of Total
                .          .    185:	// Fast path: drop lock bit.
                .          .    186:	new := atomic.AddInt32(&m.state, -mutexLocked)
                .          .    187:	if new != 0 {
                .          .    188:		// Outlined slow path to allow inlining the fast path.
                .          .    189:		// To hide unlockSlow during tracing we skip one extra frame when tracing GoUnblock.
            7.33s      7.33s    190:		m.unlockSlow(new)
                .          .    191:	}
                .          .    192:}
                .          .    193:
                .          .    194:func (m *Mutex) unlockSlow(new int32) {
                .          .    195:	if (new+mutexLocked)&mutexLocked == 0 {
       ```
       Now we see that unlocking the mutex is causing the slow down.
       
# Acknowledgements
This example is taken from [JBD](https://twitter.com/rakyll)'s post about mutex profiling in go.
https://rakyll.org/mutexprofile/

# License
See original license: https://rakyll.org/mutexprofile/

## Changes made
The program was adjusted to run as a standalone application. Additional logs were introduced to highlight the
what the application is doing. Furthermore an additional sleep was added to show that mutex profiling is only 
looking at mutex contention.