# URLs Checker

Yet another URLs checker. It interprets all incoming parameters as a URL and sends
GET requests on each of them, reads whole response body, and prints summary in the STDOUT in a view:
```
URL1 bodySize1
URL2 bodySize2
.....
```
Output is sorted by the size of the response body.

If some of the requests failed an error message will be printed in STDERR

### How to run

```
> go run main.go http://yandex.ru http://google.com http://avito.ru

http://yandex.ru 6738
http://google.com 16856
http://avito.ru 902250

```

### Build with Bazel

**Build binary**

```
> bazel build //:url-checker
```

**Build pkg/checker library**

```
> bazel build //pkg/checker:pkg_checker_library
```