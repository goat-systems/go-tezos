# Pull Request Guide
## Description
Go Tezos is a project that I beleive could have some critical data driven usage. Because of that I would like to maintain
the project in an organized fashion. Go Tezos welcomes all pull requests, so each pull request will be thoroughly reviewed. 
In order to make the reviewing process easier, please do your best to conform to these pull request rules. 

## Commits
1. Each pull request should only have 1 commit. 
    * If your pull request has more than 1 commit, please squash them. 
2. Each commit message should have an obvious commit title, with a desciptive body. 

## Pull Requests
1. When making a pull request please add a thorough description of what was changed. 
2. Please provide a unit test for any new exported functions. 
3. Your pull request should only have one objective. 
    * If you have multiple objectives, please make multiple pull requests. 

## Code Etiquette
1. If you are adding new functions and structures, please do not export them if they are not explicitly used by the end user.
2. Remember to add proper error handling 
    * See this [link](https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully) on how to properly handle errors
    * Never exit as a result of an error, bubble the error up and let the end user handle
3. Please do not add logging or print statements, logging should be in the application and not the library. 