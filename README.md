# BAREBONES BLOG AGRREGATOR
#### Hello and welcome!
This serves as a practice project for a course in Boot.dev. 
I personally don't recommend using this in real production but, you do you.
>[!IMPORTANT]
> Made this using WSL for Windows specifically Ubuntu. Everything in this repo will only work on Linux.

## Prerequisites
Here's everything you will need to install be able to build this:
- Go v1.27 or later
- Postgres v15 or later

### Install Go
Here's an easy installation using [webi](https://webinstall.dev/golang/)
```
curl.exe https://webi.ms/golang | powershell
```
Or you could download it from Go's [official download](https://go.dev/doc/install) page

### Install Postgres
This came from [Microsoft](https://learn.microsoft.com/en-us/windows/wsl/tutorials/wsl-database#install-postgresql) Docs which I've gone ahead and put here for you:
```
sudo apt update
sudo apt install postgresql postgresql-contrib
```
and then run this command to make sure it installed:
```
psql --version
```

# Postgres configuration
Alright, now that you've installed both Go and Postgres, lets go ahead and configure with some settings

First, start the server
```
sudo service postgresql start
```
Connect to it using the `psql` client.
While there are other ones, this the the simplest.
```
sudo -u postgres psql
```
To be able to interact with Postgres, you'll have to use its shell. Which you should be in at this point.

Create a new database and give it a name
```
CREATE DATABASE bloggregator;
```
Connect to the one you just created
```
\c bloggregator
```
And see something like this:
```
bloggregator=#
```
Finally, set your password.
```
ALTER USER postgres PASSWORD 'angregator';
```
Use this command to quit from the shell
```
\q
```

# Build and Install
Why do we need to build and then install? Because `build` will compile the program into binary and `install` will make it possible to run the program everywhere
```
go build
go install
```
You should be able to run the program now by only using 
```
blogaggre
```