module github.com/deniskaponchik/GoSoft/Unifi/GetAllcontrollerInfo

go 1.20

require github.com/unpoller/unifi v0.3.3

replace github.com/unpoller/unifi v0.3.3 => ../vendor/github.com/unpoller/unifi@v0.3.3

require golang.org/x/net v0.11.0 // indirect

replace golang.org/x/net v0.11.0 => ../vendor/golang.org/x/net@v0.11.0

replace golang.org/x/crypto v0.10.0 =>  ../vendor/golang.org/x/crypto@v0.10.0

replace golang.org/x/sys v0.9.0 =>  ../vendor/golang.org/x/sys@v0.9.0

replace golang.org/x/term v0.9.0 =>  ../vendor/golang.org/x/term@v0.9.0

replace golang.org/x/text v0.10.0 =>  ../vendor/golang.org/x/text@v0.10.0

replace github.com/stretchr/testify v1.8.4 => ../vendor/github.com/stretchr/testify@v1.8.4

replace github.com/davecgh/go-spew v1.1.1 => ../vendor/github.com/davecgh/go-spew@v1.1.1

replace github.com/pmezard/go-difflib v1.0.0 => ../vendor/github.com/pmezard/go-difflib@v1.0.0

replace gopkg.in/yaml.v3 v3.0.1 => ../vendor/gopkg.in/yaml.v3@v3.0.1
