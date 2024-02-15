# CB-DevTool - **EXPERIMENTAL**
A CLI tool which allows you to use a DEV annotation where the following line will not be included in the build.
This is marked as experimental as this has only been tested against default Angular configurations. 
Custom configurations or other frameworks such as React may require the project to be altered.
```powershell
Usage: cb <dir>
```
## How to use
The latest build is included in the **dist** folder for those who wish to run the application as-is.
To make usage of this tool as seamless as possible, add the application to your Path environment variables.
1. Open a command prompt window.
2. Navigate to your Angular project.
3. Execute **cb .** and wait for the project to build.
4. *Note: If you require a custom build command, edit your cb.json. Multiple commands should be separated by &&.*

## Contributions
Simply fork the project and create a pull request. Include a thorough description of your changes and it will be reviewed.