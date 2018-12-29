[![Build Status](https://travis-ci.org/BakedSoftware/go-parameters.svg?branch=master)](https://travis-ci.org/BakedSoftware/go-parameters)


```
// usage:
//   1) parse json to parameters:
// parameters.MakeParsedReq(fn http.HandlerFunc)
//   2) get the parameters:
// params := parameters.GetParams(req)
// val := params.GetXXX("key")
```