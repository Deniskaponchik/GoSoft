module github.com/deniskaponchik/GoSoft/Unifi/test_SOAP

go 1.20


require github.com/globusdigital/soap v1.4.0

replace github.com/globusdigital/soap v1.4.0 => ../vendor/github.com/globusdigital/soap@v1.4.0

replace github.com/stretchr/testify v1.8.4 => ../vendor/github.com/stretchr/testify@v1.8.4

replace github.com/pmezard/go-difflib v1.0.0 => ../vendor/github.com/pmezard/go-difflib@v1.0.0

replace github.com/davecgh/go-spew v1.1.1 => ../vendor/github.com/davecgh/go-spew@v1.1.1

replace gopkg.in/yaml.v3 v3.0.1 => ../vendor/gopkg.in/yaml.v3@v3.0.1

//require github.com/hooklift/gowsdl v0.5.0

//replace github.com/hooklift/gowsdl/cmd/gowsdl v0.5.0 =>  ../vendor/github.com/hooklift/gowsdl@v0.5.0

//require github.com/go-sql-driver/mysql v1.7.1

//replace github.com/go-sql-driver/mysql v1.7.1 => ../vendor/github.com/go-sql-driver/mysql@v1.7.1


