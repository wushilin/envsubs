# envsubs
Substitude environment variable in your properties. Perfect for running in container and configure your settings.


# How does it work

envsubs takes 3 arguments.

First argument: Input file
Second argument: Output file
Third argument: Environment variable prefix

The program reads input file, for each line, looks tokens to replace with environment variable, and write the replaced

output into the output file

For example:

# InputFile

```
jdbc.url=jdbc:mysql://${MYSQL_HOST:localhost}
http.port=${myhttp.PORT}
https.port=${HTTPS_PORT}
```

# Environment variable
```bash
MY_TEST_ENV_MYSQL_HOST=192.168.44.2
MY_TEST_ENV_MYHTTP_PORT=80
MY_TEST_ENV_HTTPS_PORT=443
```

# Running

envsubs <inputfile> <outputfile> ENV_PREFIX

```bash
envsubs inputfile outputfile MY_TEST_ENV_
```

# Output file
```
jdbc.url=jdbc:mysql://192.168.44.2
http.port=80
https.port=443
```


# Syntax
Specifying the variable using:

`${VARIABLE:DEFAULT}`
format.

VARIABLE is looked up by environment variable name `ENV_PREFIX` + `VARIABLE`

If it is not found, then `ENV_PREFIX` + HEX_ + `UTF8HEX(VARIABLE)` is checked, and the value is hex decoded as UTF8, then will be used.

If still not found, then the DEFAULT value will be used.


