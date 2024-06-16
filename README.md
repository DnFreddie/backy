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

#### Dotfiles
Recognize the repository or path to your dotfiles and copy all the executables. Ensure to change the default location from *Desktop* to *.config* in the code for proper functionality.

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
