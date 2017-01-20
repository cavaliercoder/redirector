# redirector

Redirector is a fast and lightweight HTTP server that serves a single purpose;
to redirect web clients from one URL to another. This is useful in the following
situations:

* URL shortening
* migrating from one URL layout to another
* vanity URLs

Redirector uses [BoltDB](https://github.com/boltdb/bolt), an embedded, high
performance key/value store to provide sub-millisecond responses, even when
managing millions of URL mappings.

This project is currently under active development and may require additional
contributions before it is production ready.


## Usage

```
# start web server
$ ./redirector serve

# add a URL mapping
$ ./redirector add \
	--key /abc123 \
	--dest http://my-site.com/some/path \
	--permanent

# test
$ curl -i http://localhost:8080/abc123
2017/01/04 12:06:32.531490 GET /abc123 191.262µs
HTTP/1.1 301 Moved Permanently
Content-Type: text/plain; charset=utf-8
Location: http://my-site.com/some/path
X-Content-Type-Options: nosniff
Date: Wed, 04 Jan 2017 04:06:32 GMT
Content-Length: 181

<html>
<head><title>301 Moved Permanently</title></head>
<body bgcolor="white">
<center><h1>301 Moved Permanently</h1></center>
<hr><center>redirector/1.1.1</center>
</body>
</html>

```


## License
Copyright (c) 2016 Ryan Armstrong

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
