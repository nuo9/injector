# injector
func Inject(obj interface{}, paths []string, value interface{}) error {}

# example (see more in injector_test.go)
origin:
d := {
"Id": 10,
"Name": "2333",
"Bird": {
 "Flying": true
},
"Dog": {
 "Husky": {
  "IQ": 3
 }
}
}

execute: Inject(&d, []string{"Dog", "Husky", "IQ"}, 11)

result:
d := {
"Id": 10,
"Name": "2333",
"Bird": {
 "Flying": true
},
"Dog": {
 "Husky": {
  "IQ": 11
 }
}
}