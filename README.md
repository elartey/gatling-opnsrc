# gatling-opnsrc
Stress Test tool written in Golang

Usage
========================================

./gatling-linux_64 --type='POST' --url 'http://acme.com' --rps 2 --objectType 'xml' --object '<xmlObject>Some data in here</xmlObject>' --numR 10 --headers 'Auth:SomeValue,X-Header:Foobar'


./gatling-linux_64 --type='GET' --url 'http://acme.com' --rps 2 --numR 10 -headers 'Auth:SomeValue,X-Header:Foobar'


Options
=======================================

- url:        "Url to stress test e.g. 'http://acme.com'."
- rps:        "Number of requests to make simultaneously."
- object:     "Custom object to post e.g. {'foo':'bar'}."
- objectType: "Type of object to post. e.g. 'xml' or 'json'."
- numR:       "Total number of requests to make."
- type:       "HTTP request type you'd like to make. Either 'GET' or 'POST'.")
- headers:    "Set HTTP headers. Format should be for example 'Auth:SomeToken,X-Header:Sugar'.
                Headers should be separated by commas.")


*** Headers are optional. Content-Type headers are set automatically based on object type specified.