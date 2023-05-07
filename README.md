> **_NOTE_**: this is a work in progress. main is most likely not updated. check the `not-yet-working` or `dev` branches to see what I'm currently up to. 

# distributed-task-scheduler
A distributed task scheduler application written in Go using its concurrency features. When complete, it will have a load manager, task queue, databases and more.

### TODOs

#### V1

Misc. 
- [ ] Add logs to every major function
- [ ] Add logs for every database transaction and generated SQL code
- [ ] Separate into different directories if needed

(1) Meta
- [ ] Better oganize the directories
- [ ] Add the different makefile configurations for the components (Load Balancer, Remote Servers, etc.)

(2) Load Balancer
- [ ] Build Load Balancer (LB):
    - [ ] Fix Load Balancer module code
    - [ ] Fix loadBalancerMain.go
    - [ ] Learn a Round Robin implementation 
    - [ ] Figure out how to distribute the HTTP requests
    - [ ] Figure out how to simulate a distributed environment with the few computers I have available
    - [ ] Write tests for the LB if possible

(3) Task Queue
- [ ] Build the Task Queue (TQ):
    - [ ] Learn more about it, and how a possible implementation could look like
