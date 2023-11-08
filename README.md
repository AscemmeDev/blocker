# Test work 
###  ***A client for sending data in remote server with getting limits*** 

- [link to task conditions](https://gist.github.com/m1ome/2a083b0ee3a44f70c079fabba9e5247a
  )

- ***To run the code use***:
``make run``

- ***to run the build code:***
``make build``


__**This client starts the mock data 
generator and starts the process in
usecases every five seconds.
The process in the usecases has an 
easy sending logic. The data will be sent 
asynchronously, since the task did not include 
saving a queue of batches. This logic also 
starts the scheduler to resend the packet. 
the main task is not to send data as much as 
possible without loss**__