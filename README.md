# SSH Key Manager

| Command | Support |
|---------|:-------:|
| current |    ✔    |
| cat     |    ✔    |
| use     |    ✔    |
| add     |    ✔    |
| create  |    ✔    |
| delete  |    ✔    |
| list    |    ✔    |
| dav     |    ✘    |
| help    |    ✔    |

# Usage

Input:

```shell
skm h
```

Output:

```shell
NAME:                                                           
   skm - SSH Key Manager                                        
                                                                
USAGE:                                                          
   skm [global options] command [command options] [arguments...]
                                                                
VERSION:                                                        
   v1.0.0                                                       

COMMANDS:
   current, cur  Show current key name
   cat           Output public key file contents to the console
   use, u        Use a key
   add, a        Add a key
   create, c     Create a key
   delete, del   Delete a key
   list, ls      List all the keys
   dav           WebDav Setting
   help, h       Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```