> **_NOTE_**: this is a work in progress. main is most likely not updated. check the `not-yet-working` or `dev` branches to see what I'm currently up to. 

# distributed-task-scheduler
A distributed task scheduler application written in Go using its concurrency features. When complete, it will have a load manager, task queue, databases and more.

### TODOs

#### V1

- [ ] Add logs to every major function
- [ ] Add logs for every database transaction and generated SQL code
- [ ] Build Load Balancer (LB):
    - [ ] Learn a Round Robin implementation 
    - [ ] Figure out how to distribute the HTTP requests
    - [ ] Figure out how to simulate a distributed environment with the few computers I have available
    - [ ] Write tests for the LB if possible
- [ ] Build the Task Queue (TQ):
    - [ ] Learn more about it, and how a possible implementation could look like
