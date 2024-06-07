
# go-create-project





This tool automates the setup of a Go backend project using Gin and GORM with PostgreSQL database.

## Installation

To use this tool, make sure you have Go installed on your system. Then, follow these steps:

1. Clone this repository or copy the `setup.go` script to your local machine.

2. Run the following commands to install the CLI tool:

```bash
sudo chmod +x ./setup.sh
```
```bash
./setup.sh
```
## Usage
To create a new Go project with Gin and GORM, run the following command:

```bash
go-create-project new <project-name>
```

Replace ```<project-name>``` with your desired project name.

This command will set up a new project directory, initialize the Go module, install dependencies, create a .env file, and generate the necessary files for a Gin and GORM-based project using PostgreSQL.

Make sure to set the ```DATABASE_URL``` environment variable in the .env file with your PostgreSQL connection string before running the project.

## Environment Variables
The .env file in your project directory should contain the following environment variable:

```bash
DATABASE_URL=host=localhost user=postgres dbname=<project-name>_db sslmode=disable password=yourpassword
```
Make sure <project-name> is your actual project name after creating the project.
