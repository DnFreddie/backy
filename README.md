# Baccky: Never Remember About Backups Again

## Project Overview
- **Create Temporary Configuration**: Easily create, modify, and revert configurations with one command, ideal for container/cloud usage.
- **Version Monitoring**: Track different versions of your backup files and receive email notifications about updates.
- **Dotfile Control Center**: Centralize your dotfiles, automate their deployment, and install packages based on configurations.
- **Selective Dotfile Copying**: Copy only the dotfiles for programs installed on your system, or install everything from scratch.
- **Centralized Backup Log**: Maintain a device-based backup log.
- **Tripwire Feature**: Verify file integrity and perform regular checks.
- **Future Plans**: Implement user profiles.

### Current Functionality

#### Backup 

It adds the paths "fiels and folders" to the  database for the daemon that will backup them daily

The backups are located in the *$HOME/.user_log/backups/*

Don't worry about adding duplicates 

```bash 
backup <filename direcroyName>
```

`-a` Creates the archive based on the specified files and outputs them into the current directory
```bash
backy backup -a <files and directories>
```

`-b` Backs up the given targets in a $HOME/.user_log/backups/ with the schema. Soon there will be a revert option for this.
```bash
backy backup -b <files and directories>
```
##### deamon 

It's starts the deamon that checks the **.backy.yaml** in the homedir for the **corne_time** variable
the default is set to **cron_time: "@daily"** but u can change it as u please  in the config  
 
```bash 

backy backup deamon 

```

*It's highly recomeded to add this to the startup commands so it will do regular backups*


#### Dotfiles
Recognize the repository or path to your dotfiles and copy all the executables. 

```bash
backy dot -p <path to the dotfiles or the GitHub URL>
```
This command copies configurations to the specified location, tracking them for easy reversion. Existing configurations are moved and backed up to *$HOME/.user_log/back_conf/*. To revert to a previous configuration by date, use:

```bash
backy dot revert
```

Add `-d` to remove the chosen backed-up configuration.

#### Tripwire
##### Add
Scan the specified repository, compute checksums, and store data about directories and files in a database. Paths are stored in *$HOME/.user_log/scan_paths.json*.

```bash
backy trip add -p <path>
```

##### Scan
Scan repositories specified in *$HOME/.user_log/scan_paths.json* and generate a CSV report detailing changes and new files.

```bash
backy trip scan
```

Use the `-c` flag to specify the CSV name. By default, it is named **trip_scan.csv**.

#### Config file 
For now the app looks for the `.backy.yaml` in two palaces 
- Home directory
- .config

```bash 

#Default configuration

#Email is for pushing notyfications(still developed)
email_creds:
  email: your_email@example.com
  passw: your_password
# This is the varaible for the chron daemon
cron_time: "@daily"
## This is the defualt value for the dotfiles 
config_path: ".config"

```

