# gatling-opnsrc
Stress Test tool written in Golang

Usage:

./gatling-linux_64 --type='POST' --url='http://acme.com' --rps=2 --objectType='xml' --object='<xmlObject>Some data in here</xmlObject>' --numR=10 --headers='Auth:SomeValue X-Header:Foobar'

./gatling-linux_64 --type='POST' --url='http://acme.com' --rps=2 --objectType='json' --object='{"data": "some-value-here"}' --numR=10 --headers='Auth:SomeValue X-Header:Foobar'

Headers are optional. Content-Type headers are set automatically based on object type specified.