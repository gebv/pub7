# pub7

# Features

* to collect user data
* communicate according to the script

# Status

Not Ready for production

# Examples

## example 1

Diagram

![example1](https://s3.amazonaws.com/idheap/ss/pub7_example.png)

Script

``` toml
[[nodes]]
id = "start"
text = "Hi, i am robot."
next = "q_what_name"
before = '''
if #ctx:get("name") > 0 then
    ctx:redirect("q1")
end
'''
transit = true


[[nodes]]
id = "setname"
next = "q_what_name"
transit = true

[[nodes]]
id = "q_what_name"
text = "What's your name?"
next = "h_what_name"

[[nodes]]
id = "h_what_name"
handler = '''
if ctx:text() == "" then
    print("empty name")
    ctx:redirect("q_what_name")
end
'''
next = "q1"
param = "name"
transit = true


[[nodes]]
id = "q1"
text = "{{.name}} how much 19+15 is?"
next = "h_q1"

[[nodes]]
id = "h_q1"
handler = '''
print("h_q1", ctx:text())
if ctx:text() == "34" then
    ctx:redirect("finish")
else 
    ctx:send("No, try again")
    ctx:redirect("q1")
end
'''
transit = true

[[nodes]]
id = "finish"
text = "Right! Again? /start"
```

# Quick start

``` shell
source .env # tarantool config

broombot -file scripts/example.toml run
```