# PCDeploy - A work in progress
## Project Explanation:
PC Deploy will be a control center for PC Deployment. It will have the following responsibilities:
* Listen for client connection and accept them
* Listen for updates from each client and be able to save them as a log
* Be able to choose which computer to "look" at.
    - Use channels to choose which client is in the "foreground". All the rest are in the background
* Be able to send commands to a client.

### 3 main components:

#### Control Center
- This is where the shell lives. It parses commands and handles them. From the control center we can do the following. This is what the user views.
* List of commands:
    - create <option>: create client OR create package
    - show <option>: show connections, clients, packages...
    - select <option>: select a connection, client...
    - kill <connection>: Kills the mentioned connections.
    - <ctrl-z>: Put current connection in background.
    - help <option>
    - quit/exit
* View status of files. If some files have not been updated, they will be displayed in red.

The client will have the following responsibilities:
* Try to connect to control center
* If connection is through, it will start the powershell scripts to setup the PC
* Send updates to the control center on the scripts

#### Server
- The server will listen for upcoming connections from a client. Once it does, it will task the client wiht starting the powershell scripts. It will receive logs which will be saved to a file named `<PC-NAME>_logs.txt`.
- The server will also be able to execute code on the remote computer.
- The server will handle all connections, foreground and background. Users can put connections in background.

#### Client
- The client will be run on the remote computer thanks to the unattend.xml on a USB connected to the remote computer.
- It will try to connect to the server for 10 minutes max. Once it has connected, it will wait for the okay to start running the scripts.
- Once scripts have started, it will send updates to the server
- There are possibly more responsibility for the client down the line, like create a task via the task scheduler etc...

##### Commands Explanation and Workflow
1. Create:
    * Users can choose to create a client or package. Let's walk through both:
        - package: A couple of questions will need to be asked and answered. Most are optional, as default values are filled.  
            * **Creates the unattend.xml file based on answers and writes to USB**.
            * Writes the scripts and client.exe to the USB as well.
        - client: Only the questions change. All questions have to be answered.

## Up next:
* Write channels and switch between connections in foreground/background.
* Close connections gracefully
* More security on accepting clients. Only clients that send a specific request are allowed to connect

## Sending command to put connection in foreground/background:
* Whatever the current connection is, sending a command `background` or `CTRL+z` should put it in "background"
* Any connection in background should log updates rather than print them to console.
* To pick a connection to put in foreground, I can select the index from the Clients array
* Each connection will have a channel that listens for the status change command.
* When a new connection is detected, it automatically becomes the "foreground" connection
* How does the server handle this though?

# VERY IMPORTANT - CURRENT FOCUS:
1.  **CURRENT MAIN FOCUS** (Will make program usable, possibly v1.0): 
    1. Be able to create clients,
    2. Create packages (unattend.xml) to be written to USB as well as the scripts. 
    3. Have the client connect to the server
    4. Have the client start scripts and update server
    5. Server be able to write logs to text file.
    1. Create packages and write to usb
    2. Config maintenance (update agent installers..., scripts, edit default config)
2. **FUTURE FOCUS**:
    1. Handle connections better: background, kill, foreground...
    2. Looks: colors and showing as tables and stuff
    3. Show status of files.
    4. Creating default configs for clients as well as templates for employee per client.
3. **MORE TO COME**

