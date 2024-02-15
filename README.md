# CB-DevTool - **EXPERIMENTAL**
A CLI tool which allows you to use a DEV annotation where the following line will not be included in the build.
This is marked as experimental as this has only been tested against default Angular configurations. 
Custom configurations or other frameworks such as React may require the project to be altered.
```powershell
Usage: cb <dir>
```
## Project Plans
- Clean up code.
## Tag usage
### Single Line Example
```ts
api = "https://myliveapi.com/api"

constructor() {
    // @DEV
    api = "http://127.0.0.1:8080/api"
}
```
### Code Block Example
```ts
// @START-DEV - You can also add comments, by the way.
function myFunc() {
    let x = 5
    let y = 4

    console.log(x + y)
}
// @END-DEV
```
### Explanation
In the snippet above, I have marked my development API URL with the @DEV tag. This indicates to the tool that 
the line below should not be included in the build. A backup of this file will be created in the build and the 
original will be modified to remove all marked lines. After the build, the backup will be restored.
## How to use
The latest build is included in the **dist** folder for those who wish to run the application as-is.
To make usage of this tool as seamless as possible, add the application to your Path environment variables.
1. Open a command prompt window.
2. Navigate to your Angular project.
3. Mark above any lines that shouldn't be included in the build with the @DEV comment.
4. Execute **cb .** and wait for the project to build.
*Note: If you require a custom build command, edit your cb.json. Multiple commands should be separated by &&.*

## Contributions
Simply fork the project and create a pull request. Include a thorough description of your changes and it will be reviewed.