- [Git Setup](#git-setup)
- [GoHack](#gohack)
  - [About](#about)
  - [Supported Commands](#supported-commands)
    - [get](#get)
      - [Example 1](#example-1)
      - [Example 2](#example-2)
    - [undo](#undo)
      - [Example 1](#example-1-1)
      - [Example 2](#example-2-1)
    - [finish](#finish)
      - [Example 1](#example-1-2)
      - [Example 2](#example-2-2)
- [Running CBA Test Cases](#running-cba-test-cases)
  - [Method 1: Using existing environment](#method-1-using-existing-environment)
  - [Method 2: Using your VS Code devcontainer](#method-2-using-your-vs-code-devcontainer)
    - [Via the Creation Center](#via-the-creation-center)
    - [Via the Command line](#via-the-command-line)

## Git Setup
To configure Git within the container, put your ssh key (id_rsa) file in the **.devcontainer** directory (or any sub directory), and it will automatically be added when the container starts.

## GoHack

### About

There is a gohack.sh script included in root directory of the project. This script provides additional functionality for the gohack project https://github.com/rogpeppe/gohack. The script's purpose is to make working on module dependencies a more seamless process by adding `replace` directives within your **go.mod** file, and adding these modules to your VS Code workspace. 

### Supported Commands

#### get

```
./gohack.sh get [module_pattern]
```
Specifies a module(s) to hack. Any module in your go.mod containing the `module_pattern` the following actions will occur:
- the module's git repository will be cloned locally. The repository is accessible at **/home/ossm/gohack** within the container. This directory is also mapped to your host filesystem in **.devcontainer/gohack**
- the local repository will create a new branch called 'changeme' at the commit specified by your go.mod file
- your go.mod file will be modified to add the `replace` directive to point at the local repository
- the local repository will be added to your VS Code workspace. If the workspace is not already open, a new window will open that you should use instead. (Not applicable if not being run from within a VS Code devcontainer)

##### Example 1
Hacking a new module for the first time
```
[ossm@ossm_ca_dev ca]$ ./gohack.sh get api
======================================================
Executing ./gohack.sh get api
======================================================
gohack> Current Status:

-----------------------
gohack> Hacking module: github.ibm.com/Spectrum-Protect/sp-ossm-api
-----------------------
sp-ossm-api> Hacked module directory /home/ossm/gohack/github.ibm.com/Spectrum-Protect/sp-ossm-api
sp-ossm-api> Executing: 'gohack get -vcs -f github.ibm.com/Spectrum-Protect/sp-ossm-api'
creating github.ibm.com/Spectrum-Protect/sp-ossm-api@v0.0.0-20211013184913-eb7bda92dd41
github.ibm.com/Spectrum-Protect/sp-ossm-api => /home/ossm/gohack/github.ibm.com/Spectrum-Protect/sp-ossm-api
sp-ossm-api> Creating a temporary branch
Switched to a new branch 'changeme'
sp-ossm-api> Done

vscode> Updating workspace /home/ossm/devcontainer.code-workspace
vscode> Done.
```
\* Note: if this is your first time running the gohack script, you will see additional output installing dependencies.

##### Example 2
Hacking a module that has been hacked previously (and a local directory already exists)

Files will be moved to a temporary location to ensure no data is lost. Nothing will ever be deleted automatically. 

```
-----------------------
gohack> Hacking module: github.ibm.com/Spectrum-Protect/sp-ossm-api
-----------------------
sp-ossm-api> Hacked module directory /home/ossm/gohack/github.ibm.com/Spectrum-Protect/sp-ossm-api
sp-ossm-api> Git repository exists
sp-ossm-api> Module currently on branch 'changeme'
sp-ossm-api> Branch looks up to date
sp-ossm-api> Existing git repository looks clean but moving to /home/ossm/gohack/.tmp/sp-ossm-api-changeme-1640127140 just in case. Ensure nothing was lost before deleting it
sp-ossm-api> Executing: 'gohack get -vcs -f github.ibm.com/Spectrum-Protect/sp-ossm-api'
creating github.ibm.com/Spectrum-Protect/sp-ossm-api@v0.0.0-20211013184913-eb7bda92dd41
github.ibm.com/Spectrum-Protect/sp-ossm-api => /home/ossm/gohack/github.ibm.com/Spectrum-Protect/sp-ossm-api
sp-ossm-api> Creating a temporary branch
Switched to a new branch 'changeme'
sp-ossm-api> Done
```

#### undo

```
./gohack.sh undo [module_pattern]
```
Specifies a module(s) to undo the hack and revert changes to your go.mod file*. Any module in your go.mod containing the `module_pattern` the following actions will occur:
- the replace directive in your go.mod file will be removed
- the directory will be removed from your VS Code workspace

##### Example 1
Hacked module has no uncommitted changes
```
[ossm@ossm_ca_dev ca]$ ./gohack.sh undo api
======================================================
Executing ./gohack.sh undo api
======================================================
gohack> Current Status:
github.ibm.com/Spectrum-Protect/sp-ossm-api => /home/ossm/gohack/github.ibm.com/Spectrum-Protect/sp-ossm-api

-----------------------
gohack> Undoing module: github.ibm.com/Spectrum-Protect/sp-ossm-api
-----------------------
sp-ossm-api> Hacked module directory /home/ossm/gohack/github.ibm.com/Spectrum-Protect/sp-ossm-api
sp-ossm-api> Git repository exists
sp-ossm-api> Module currently on branch 'changeme'
sp-ossm-api> Branch looks up to date
sp-ossm-api> Executing: 'gohack undo github.ibm.com/Spectrum-Protect/sp-ossm-api'
dropped github.ibm.com/Spectrum-Protect/sp-ossm-api
sp-ossm-api> Done

vscode> Updating workspace /home/ossm/devcontainer.code-workspace
vscode> Done.
```


##### Example 2
Hacked module has uncommitted changes. Hack is still undone, but a warning is displayed to indicate you should finish committing/pushing your changes.
```
[ossm@ossm_ca_dev ca]$ ./gohack.sh undo api
======================================================
Executing ./gohack.sh undo api
======================================================
gohack> Current Status:
github.ibm.com/Spectrum-Protect/sp-ossm-api => /home/ossm/gohack/github.ibm.com/Spectrum-Protect/sp-ossm-api

-----------------------
gohack> Undoing module: github.ibm.com/Spectrum-Protect/sp-ossm-api
-----------------------
sp-ossm-api> Hacked module directory /home/ossm/gohack/github.ibm.com/Spectrum-Protect/sp-ossm-api
sp-ossm-api> Git repository exists
sp-ossm-api> git status:\n M pctx/protect_context.go
sp-ossm-api> Module has uncommitted/unstaged changes, aborting
sp-ossm-api> WARNING: git repository is not clean. Will undo the hack, but ensure you have push all necessary changes
sp-ossm-api> Executing: 'gohack undo github.ibm.com/Spectrum-Protect/sp-ossm-api'
dropped github.ibm.com/Spectrum-Protect/sp-ossm-api
sp-ossm-api> Done

vscode> Updating workspace /home/ossm/devcontainer.code-workspace
vscode> Done.
[ossm@ossm_ca_dev ca]$ 
```


#### finish
```
./gohack.sh finish [module_pattern]
```
Specifies a module(s) to finish after you have completed your modifications. For any module in your go.mod containing the `module_pattern` the following actions will occur:
- the go.mod file will be updated to reference the proper commit for modified module dependency
- the `replace` directive will be removed by running `gohack undo`
  
\* Note: If the repository has any unstaged changes, or no branch was checked an error will be displayed.

##### Example 1
Finishing a module that has all changed pushed to the remote repository

```
[ossm@ossm_ca_dev ca]$ ./gohack.sh finish api
======================================================
Executing ./gohack.sh finish api
======================================================
gohack> Current Status:
github.ibm.com/Spectrum-Protect/sp-ossm-api => /home/ossm/gohack/github.ibm.com/Spectrum-Protect/sp-ossm-api

-------------------------
gohack> Finishing module: github.ibm.com/Spectrum-Protect/sp-ossm-api
-------------------------
sp-ossm-api> Hacked module directory /home/ossm/gohack/github.ibm.com/Spectrum-Protect/sp-ossm-api
sp-ossm-api> Git repository exists
sp-ossm-api> Module currently on branch 'ramahin'
sp-ossm-api> Branch looks up to date, updating go.mod
sp-ossm-api> Executing: 'go get github.ibm.com/Spectrum-Protect/sp-ossm-api@ramahin'
sp-ossm-api> Executing: 'gohack undo github.ibm.com/Spectrum-Protect/sp-ossm-api'
dropped github.ibm.com/Spectrum-Protect/sp-ossm-api
sp-ossm-api> Done

vscode> Updating workspace /home/ossm/devcontainer.code-workspace
vscode> Done.
```

##### Example 2

Finishing a module that has not had all changes pushed to the remote repository. The pseudo version for the go.mod will be generated manually, and may not be accurate

```
[ossm@ossm_ca_dev ca]$ ./gohack.sh finish api
======================================================
Executing ./gohack.sh finish api
======================================================
gohack> Current Status:
github.ibm.com/Spectrum-Protect/sp-ossm-api => /home/ossm/gohack/github.ibm.com/Spectrum-Protect/sp-ossm-api

-------------------------
gohack> Finishing module: github.ibm.com/Spectrum-Protect/sp-ossm-api
-------------------------
sp-ossm-api> Hacked module directory /home/ossm/gohack/github.ibm.com/Spectrum-Protect/sp-ossm-api
sp-ossm-api> Git repository exists
sp-ossm-api> Module currently on branch 'changeme'
sp-ossm-api> WARNING: Remote does not contain branch 'changeme'. Ensure you push your changes pushing any changes in the main repo
sp-ossm-api> WARNING: Manually editing the go.mod file to reference pseudo-version: 
sp-ossm-api> Setting version for github.ibm.com/Spectrum-Protect/sp-ossm-api to '20211221235043-15b7e5e1f89e'
sp-ossm-api> Executing: 'gohack undo github.ibm.com/Spectrum-Protect/sp-ossm-api'
dropped github.ibm.com/Spectrum-Protect/sp-ossm-api
sp-ossm-api> Done

vscode> Updating workspace /home/ossm/devcontainer.code-workspace
cannot determine main module: go: github.ibm.com/Spectrum-Protect/sp-ossm-api@v0.0.0-20211221235043-15b7e5e1f89e: invalid version: unknown revision 15b7e5e1f89e
vscode> Done.
```

## Running CBA Test Cases

### Method 1: Using existing environment

This method is probably the easiest and does not require you to use VS Code and the OSSM container used for the test will be separate from the VS Code devcontainer. You will need to ensure compiled executables exist.

1. Ensure you have setup the server-dev-with-cba project  
   (https://github.ibm.com/Spectrum-Protect/server-dev-with-cba/wiki/Project-Setup)
2. Ensure your OSSM environment is setup properly  
   (https://github.ibm.com/Spectrum-Protect/server-dev-with-cba/wiki/OSSM#setup)
3. Clone the repo, or download a test case from https://github.ibm.com/Spectrum-Protect/server-dev-cba-tests/tree/current/tests/ossm
4. Run your test case via the [Command Line](https://github.ibm.com/Spectrum-Protect/server-dev-with-cba/wiki/CBA-Execution) or through the [Creation Center](https://github.ibm.com/Spectrum-Protect/server-dev-with-cba/wiki/Creation-Center)

### Method 2: Using your VS Code devcontainer

The steps to accomplish this vary depending on if you wish to run via the Creation Center or the command line. This section assumes you are familiar with both options. 

#### Via the Creation Center

1. Ensure your container is part of the **cba_network** docker network by specifying the **OSSM_NETWORK** in the **.devcontainer/.env** file. (otherwise other containers will not be able to communicate with this container)
2. Open the Creation Center in your browser
3. Navigate to the Environment page, copy/paste the OSSM product
4. Change the product name to `SpOssmAgent.ossm1.devcontainer`
5. Change the host to `ossm_ca_dev`
6. Navigate to the Test Data page, duplicate the Development Test Data, rename the duplicate to Development-Devcontainer
7. On the variables tab, update `SpOssmAgent1` to equal `SpOssmAgent.ossm1.devcontainer`
8. Navigate to the Execution page, and open the Run dialog
9. Select your Test Case, Environment, and your new Test Data. Verify that `SpOssmAgent.ossm1.devcontainer` is shown as the value for `SpOssmAgent1`
10. Test Case should now run within your devcontainer directly

#### Via the Command line

If you have performed the Creation Center steps above, with the Creation Center running you can run the `sync.sh|bat` script located in the **cba-execution** directory. This will export everything from the Creation Center into a usable format to run the test case. 

If you have not, follow these steps (or perform the above):
1. Open the **cba-execution/environments/dev.env** file, and duplicate the **SpOssmAgent.ossm1** json object. Ensure json syntax is correct.
2. Rename the key and the `productName` to `SpOssmAgent.ossm1.devcontainer`
3. Duplicate the **cba-execution/tests/dev.testdata** and name it **devcontainer.testdata**
4. Update `SpOssmAgent1` to equal `SpOssmAgent.ossm1.devcontainer`
5. Run your test case and add the parameter `-d tests/devcontainer.testdata`
